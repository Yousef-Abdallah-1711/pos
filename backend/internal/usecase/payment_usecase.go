package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/commerce"
	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/database"
	"pos-saas/internal/repository"
	"pos-saas/internal/service"
)

// PaymentUseCase handles payment business logic.
type PaymentUseCase struct {
	paymentRepo            repository.PaymentRepository
	subscriptionRepo       repository.SubscriptionRepository
	planRepo               repository.PlanRepository
	paymobClient           service.PaymobClient
	plans                  map[string]domain.SubscriptionPlan
	integrationID          int
	paymentKeyExpirySecond int
	db                     *sql.DB
}

var errWebhookAlreadyProcessed = errors.New("webhook already processed")

// NewPaymentUseCase creates a payment use case.
func NewPaymentUseCase(
	paymentRepo repository.PaymentRepository,
	subscriptionRepo repository.SubscriptionRepository,
	planRepo repository.PlanRepository,
	paymobClient service.PaymobClient,
	db *sql.DB,
) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo:            paymentRepo,
		subscriptionRepo:       subscriptionRepo,
		planRepo:               planRepo,
		paymobClient:           paymobClient,
		plans:                  buildDefaultPlans(),
		integrationID:          getEnvInt("PAYMOB_INTEGRATION_ID", 0),
		paymentKeyExpirySecond: 3600,
		db:                     db,
	}
}

// GetPlans returns available plans.
func (uc *PaymentUseCase) GetPlans() []domain.SubscriptionPlan {
	return uc.listPlans(context.Background())
}

// InitiatePayment starts Paymob checkout for a tenant.
func (uc *PaymentUseCase) InitiatePayment(ctx context.Context, tenantID int64, userID int64, planID string, billingData map[string]string) (*domain.InitiatePaymentResponse, error) {
	log.Printf("[PaymentUseCase] InitiatePayment - tenant_id: %d, plan_id: %s\n", tenantID, planID)
	if err := uc.requirePaymentInfra(); err != nil {
		return nil, err
	}

	if tenantID <= 0 {
		return nil, fmt.Errorf("invalid tenant id")
	}

	if strings.TrimSpace(planID) == "" {
		planID = "monthly"
	}

	plan, ok := uc.getPlanByID(ctx, planID)
	if !ok {
		return nil, fmt.Errorf("invalid plan id: %s", planID)
	}

	merchantOrderID := fmt.Sprintf("tenant_%d_%s_%d", tenantID, planID, time.Now().UnixNano())

	var userIDPtr *int64
	if userID > 0 {
		userIDPtr = &userID
	}

	billingName := strings.TrimSpace(strings.Join([]string{billingData["first_name"], billingData["last_name"]}, " "))
	if billingName == "" {
		billingName = "NA"
	}
	billingEmail := strings.TrimSpace(billingData["email"])
	if billingEmail == "" {
		billingEmail = "unknown@example.com"
	}
	billingPhone := strings.TrimSpace(billingData["phone_number"])
	if billingPhone == "" {
		billingPhone = "+200000000000"
	}

	payment := &domain.Payment{
		TenantID:        tenantID,
		UserID:          userIDPtr,
		AmountCents:     plan.PriceCents,
		Currency:        plan.Currency,
		Status:          domain.PaymentStatusPending,
		BillingName:     &billingName,
		BillingEmail:    &billingEmail,
		BillingPhone:    &billingPhone,
		PlanID:          plan.ID,
		MerchantOrderID: merchantOrderID,
		Metadata: map[string]interface{}{
			"plan_id": plan.ID,
			"source":  "dashboard",
		},
	}

	paymentID, err := uc.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create pending payment: %w", err)
	}

	authToken, err := uc.paymobClient.Authenticate()
	if err != nil {
		_ = uc.markPaymentFailed(ctx, tenantID, paymentID, err.Error())
		return nil, fmt.Errorf("failed to authenticate with paymob: %w", err)
	}

	paymobOrderID, err := uc.paymobClient.CreateOrder(authToken, plan.PriceCents, plan.Currency, merchantOrderID)
	if err != nil {
		_ = uc.markPaymentFailed(ctx, tenantID, paymentID, err.Error())
		return nil, fmt.Errorf("failed to create paymob order: %w", err)
	}

	paymentKey, err := uc.paymobClient.GetPaymentKey(
		authToken,
		paymobOrderID,
		plan.PriceCents,
		plan.Currency,
		uc.integrationID,
		uc.normalizeBillingData(billingData),
		uc.paymentKeyExpirySecond,
	)
	if err != nil {
		_ = uc.markPaymentFailed(ctx, tenantID, paymentID, err.Error())
		return nil, fmt.Errorf("failed to get paymob payment key: %w", err)
	}

	if err := uc.paymentRepo.UpdatePaymentOrder(ctx, tenantID, paymentID, paymobOrderID, authToken); err != nil {
		_ = uc.markPaymentFailed(ctx, tenantID, paymentID, err.Error())
		return nil, fmt.Errorf("failed to update payment order details: %w", err)
	}

	uc.writeSubscriptionHistory(ctx, &domain.CreateHistoryRecordParams{
		TenantID:    tenantID,
		Action:      domain.ActionPaymentInitiated,
		PerformedBy: userIDPtr,
		Metadata: map[string]interface{}{
			"payment_id":        paymentID,
			"merchant_order_id": merchantOrderID,
			"paymob_order_id":   paymobOrderID,
			"amount_cents":      plan.PriceCents,
			"currency":          plan.Currency,
			"plan_id":           plan.ID,
		},
	})

	return &domain.InitiatePaymentResponse{
		PaymentID:       paymentID,
		IframeURL:       uc.paymobClient.GetIframeURL(paymentKey),
		PaymobOrderID:   paymobOrderID,
		MerchantOrderID: merchantOrderID,
	}, nil
}

// HandleWebhook processes Paymob webhook callbacks.
func (uc *PaymentUseCase) HandleWebhook(ctx context.Context, rawPayload []byte, payload *domain.PaymobWebhookPayload, receivedHMAC string) error {
	if payload == nil {
		return fmt.Errorf("payload is required")
	}

	if strings.TrimSpace(receivedHMAC) == "" {
		receivedHMAC = payload.HMAC
	}

	if !uc.paymobClient.VerifyHMAC(rawPayload, receivedHMAC) {
		return fmt.Errorf("invalid hmac signature")
	}

	if payload.Obj.Order.ID <= 0 {
		return fmt.Errorf("missing paymob order id")
	}

	if err := uc.requirePaymentInfra(); err != nil {
		return err
	}

	payment, err := uc.paymentRepo.GetPaymentByPaymobOrderID(ctx, payload.Obj.Order.ID, payload.Obj.Order.MerchantOrderID)
	if err != nil {
		return fmt.Errorf("payment not found for paymob order %d: %w", payload.Obj.Order.ID, err)
	}

	if payment.Status == domain.PaymentStatusSuccess && !payload.Obj.IsRefunded {
		log.Printf("[PaymentUseCase] HandleWebhook idempotent skip - payment_id: %d\n", payment.ID)
		return nil
	}
	if payment.Status == domain.PaymentStatusRefunded && payload.Obj.IsRefunded {
		log.Printf("[PaymentUseCase] HandleWebhook refund idempotent skip - payment_id: %d\n", payment.ID)
		return nil
	}

	// A refunded payment is terminal: a replayed or late "success" event must never re-activate
	// the subscription it was refunded for.
	if payment.Status == domain.PaymentStatusRefunded {
		log.Printf("[PaymentUseCase] HandleWebhook ignored - payment_id: %d is refunded\n", payment.ID)
		return nil
	}

	planID := strings.TrimSpace(payment.PlanID)
	if planID == "" {
		planID = "monthly"
	}

	plan, ok := uc.getPlanByID(ctx, planID)
	if !ok {
		return fmt.Errorf("unknown plan id for payment: %s", planID)
	}

	// A signed callback settles the immutable payment snapshot, not today's offer.
	// Catalog price changes must not invalidate previously initiated payments.
	if payment.AmountCents <= 0 || payload.Obj.AmountCents != payment.AmountCents {
		return fmt.Errorf("amount mismatch: expected %d got %d", payment.AmountCents, payload.Obj.AmountCents)
	}

	if !strings.EqualFold(payload.Obj.Currency, payment.Currency) {
		return fmt.Errorf("currency mismatch: expected %s got %s", payment.Currency, payload.Obj.Currency)
	}

	txnID := payload.Obj.ID
	var txnIDPtr *int64
	if txnID > 0 {
		txnIDPtr = &txnID
	}

	cardLast4 := extractLastFour(payload.Obj.SourceData.Pan)
	cardBrand := strings.TrimSpace(payload.Obj.SourceData.SubType)

	var cardLast4Ptr *string
	var cardBrandPtr *string
	if cardLast4 != "" {
		cardLast4Ptr = &cardLast4
	}
	if cardBrand != "" {
		cardBrandPtr = &cardBrand
	}
	if payload.Obj.IsRefunded {
		if payment.PaymobTransactionID != nil && *payment.PaymobTransactionID != payload.Obj.ID {
			return fmt.Errorf("refund transaction mismatch for payment %d", payment.ID)
		}
		return uc.applyRefundWebhook(ctx, payment, txnIDPtr, cardLast4Ptr, cardBrandPtr)
	}

	// Validate the signed callback's order amount and currency before persisting
	// this terminal state, so inconsistent provider data cannot void a payment.
	if payload.Obj.IsVoided {
		if uc.db == nil {
			log.Printf("[PaymentUseCase] HandleWebhook void acknowledged without database - payment_id: %d\n", payment.ID)
			return nil
		}
		return uc.applyVoidWebhook(ctx, payment, payload.Obj.Order.ID, payload.Obj.ID)
	}

	isSuccess := payload.Obj.Success && !payload.Obj.Pending && !payload.Obj.ErrorOccured
	if payload.Obj.Pending {
		log.Printf("[PaymentUseCase] HandleWebhook pending - payment_id: %d\n", payment.ID)
		return nil
	}

	if isSuccess {
		if uc.db == nil {
			return fmt.Errorf("payment database is not configured")
		}

		paidAt := time.Now()
		expiresAt := paidAt.AddDate(0, periodMonthsFromPlan(plan), 0)
		err := database.WithTenantTx(ctx, uc.db, payment.TenantID, 0, func(tx *sql.Tx) error {
			if err := uc.updatePaymentStatusTx(tx, ctx, payment.TenantID, payment.ID, domain.PaymentStatusSuccess, txnIDPtr, cardLast4Ptr, cardBrandPtr, nil, &paidAt); err != nil {
				if errors.Is(err, errWebhookAlreadyProcessed) {
					log.Printf("[PaymentUseCase] HandleWebhook idempotent tx skip - payment_id: %d\n", payment.ID)
					return errWebhookAlreadyProcessed
				}
				return fmt.Errorf("failed to mark payment success: %w", err)
			}
			return uc.updateTenantSubscriptionTx(tx, ctx, payment.TenantID, domain.StatusActive, expiresAt, paidAt, payment.AmountCents, payment.Currency, payment.PlanID)
		})
		if errors.Is(err, errWebhookAlreadyProcessed) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to commit payment success transaction: %w", err)
		}

		uc.writeSubscriptionHistory(ctx, &domain.CreateHistoryRecordParams{
			TenantID:    payment.TenantID,
			Action:      domain.ActionPaymentSuccess,
			PerformedBy: payment.UserID,
			NewStatus:   stringValuePtr(domain.StatusActive),
			Metadata: map[string]interface{}{
				"payment_id":            payment.ID,
				"paymob_order_id":       payload.Obj.Order.ID,
				"paymob_transaction_id": payload.Obj.ID,
				"amount_cents":          payment.AmountCents,
				"currency":              payment.Currency,
			},
		})
		uc.writeSubscriptionHistory(ctx, &domain.CreateHistoryRecordParams{
			TenantID:    payment.TenantID,
			Action:      domain.ActionSubscriptionActivated,
			PerformedBy: payment.UserID,
			NewStatus:   stringValuePtr(domain.StatusActive),
			Metadata: map[string]interface{}{
				"payment_id":      payment.ID,
				"expires_at":      expiresAt.Format(time.RFC3339),
				"merchant_order":  payment.MerchantOrderID,
				"subscription_id": payment.PlanID,
			},
		})

		return nil
	}

	errMsg := "payment failed"
	if payload.Obj.Pending {
		errMsg = "payment is pending"
	} else if payload.Obj.ErrorOccured {
		errMsg = "paymob reported transaction error"
	}

	if uc.db == nil {
		return uc.paymentRepo.UpdatePaymentStatus(ctx, payment.TenantID, payment.ID, domain.PaymentStatusFailed, txnIDPtr, cardLast4Ptr, cardBrandPtr, &errMsg, nil)
	}
	err = database.WithTenantTx(ctx, uc.db, payment.TenantID, 0, func(tx *sql.Tx) error {
		return uc.updatePaymentStatusTx(tx, ctx, payment.TenantID, payment.ID, domain.PaymentStatusFailed, txnIDPtr, cardLast4Ptr, cardBrandPtr, &errMsg, nil)
	})
	if errors.Is(err, errWebhookAlreadyProcessed) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to commit payment failure: %w", err)
	}

	uc.writeSubscriptionHistory(ctx, &domain.CreateHistoryRecordParams{
		TenantID:    payment.TenantID,
		Action:      domain.ActionPaymentFailed,
		PerformedBy: payment.UserID,
		Metadata: map[string]interface{}{
			"payment_id":            payment.ID,
			"paymob_order_id":       payload.Obj.Order.ID,
			"paymob_transaction_id": payload.Obj.ID,
			"amount_cents":          payment.AmountCents,
			"currency":              payment.Currency,
			"error":                 errMsg,
		},
	})

	return nil
}

func (uc *PaymentUseCase) applyVoidWebhook(ctx context.Context, payment *domain.Payment, orderID, transactionID int64) error {
	if uc.db == nil {
		return fmt.Errorf("payment database is not configured")
	}
	err := database.WithTenantTx(ctx, uc.db, payment.TenantID, 0, func(tx *sql.Tx) error {
		return uc.updatePaymentStatusTx(tx, ctx, payment.TenantID, payment.ID, domain.PaymentStatusVoided, &transactionID, nil, nil, nil, nil)
	})
	if errors.Is(err, errWebhookAlreadyProcessed) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to commit payment void: %w", err)
	}
	uc.writeSubscriptionHistory(ctx, &domain.CreateHistoryRecordParams{
		TenantID:    payment.TenantID,
		Action:      domain.ActionPaymentVoided,
		PerformedBy: payment.UserID,
		Metadata: map[string]interface{}{
			"payment_id":            payment.ID,
			"paymob_order_id":       orderID,
			"paymob_transaction_id": transactionID,
		},
	})
	return nil
}

// GetPaymentHistory returns paginated payment history.
func (uc *PaymentUseCase) GetPaymentHistory(ctx context.Context, tenantID int64, page int, pageSize int) ([]domain.Payment, int, error) {
	if err := uc.requirePaymentInfra(); err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	return uc.paymentRepo.GetPaymentsByTenantID(ctx, tenantID, pageSize, offset)
}

// GetPaymentHistoryByCursor returns cursor-based payment history ordered by newest first.
func (uc *PaymentUseCase) GetPaymentHistoryByCursor(ctx context.Context, tenantID int64, cursorID int64, limit int) ([]domain.Payment, *int64, error) {
	if err := uc.requirePaymentInfra(); err != nil {
		return nil, nil, err
	}

	if limit < 1 || limit > 100 {
		limit = 20
	}
	if cursorID < 0 {
		cursorID = 0
	}

	return uc.paymentRepo.GetPaymentsByTenantCursor(ctx, tenantID, cursorID, limit)
}

// GetCurrentSubscription returns payment-extended subscription details.
func (uc *PaymentUseCase) GetCurrentSubscription(ctx context.Context, tenantID int64) (*domain.CurrentSubscriptionResponse, error) {
	if err := uc.requirePaymentInfra(); err != nil {
		return nil, err
	}

	subscription, err := uc.subscriptionRepo.GetTenantSubscriptionInfo(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to load subscription: %w", err)
	}

	billingInfo, err := uc.paymentRepo.GetTenantBillingInfo(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to load billing info: %w", err)
	}

	planID := strings.TrimSpace(billingInfo.SubscriptionPlan)
	if planID == "" {
		planID = "monthly"
	}
	plans := uc.listPlans(ctx)
	planByID := make(map[string]domain.SubscriptionPlan, len(plans))
	for _, p := range plans {
		planByID[p.ID] = p
	}

	plan, ok := planByID[planID]
	if !ok && uc.planRepo != nil {
		// Historical subscriptions retain their legacy catalog price and identity.
		if historical, err := uc.planRepo.GetByCode(ctx, planID); err == nil && historical != nil {
			plan = *historical
			ok = true
		}
	}
	if !ok {
		plan = uc.defaultPlan()
	}

	daysRemaining := subscription.DaysRemaining
	if strings.EqualFold(subscription.SubscriptionStatus, domain.StatusActive) && billingInfo.SubscriptionExpiresAt != nil {
		daysRemaining = daysUntil(*billingInfo.SubscriptionExpiresAt)
	}

	currentPlan := toPlanSummary(plan)
	availablePlans := make([]domain.SubscriptionPlanSummary, 0, len(plans))
	for _, p := range plans {
		availablePlans = append(availablePlans, toPlanSummary(p))
	}

	return &domain.CurrentSubscriptionResponse{
		PreferredBillingPeriod:  subscription.PreferredBillingPeriod,
		CheckoutAvailable:       commerce.Current().CheckoutAvailable,
		SubscriptionStatus:      subscription.SubscriptionStatus,
		TrialEndDate:            subscription.TrialEndDate,
		SubscriptionExpiresAt:   billingInfo.SubscriptionExpiresAt,
		SubscriptionCancelledAt: billingInfo.SubscriptionCancelledAt,
		CancelAtPeriodEnd:       billingInfo.SubscriptionCancelledAt != nil,
		LastPaymentAt:           billingInfo.LastPaymentAt,
		DaysRemaining:           daysRemaining,
		CanAccess:               subscription.CanAccess,
		BlockReason:             subscription.BlockReason,
		Plan:                    currentPlan,
		CurrentPlan:             currentPlan,
		AvailablePlans:          availablePlans,
	}, nil
}

// CancelSubscription schedules cancellation at period end for an active tenant subscription.
func (uc *PaymentUseCase) CancelSubscription(
	ctx context.Context,
	tenantID int64,
	userID int64,
	reason *string,
	ipAddress string,
	userAgent string,
) (*domain.CancelSubscriptionResponse, error) {
	if err := uc.requirePaymentInfra(); err != nil {
		return nil, err
	}
	if uc.db == nil {
		return nil, fmt.Errorf("payment database is not configured")
	}
	if tenantID <= 0 {
		return nil, fmt.Errorf("invalid tenant id")
	}

	subscription, err := uc.subscriptionRepo.GetTenantSubscriptionInfo(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to load subscription: %w", err)
	}
	if !strings.EqualFold(subscription.SubscriptionStatus, domain.StatusActive) {
		return nil, fmt.Errorf("only active subscriptions can be cancelled")
	}

	billingInfo, err := uc.paymentRepo.GetTenantBillingInfo(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to load billing info: %w", err)
	}
	if billingInfo.SubscriptionExpiresAt == nil {
		return nil, fmt.Errorf("subscription expiry date is missing")
	}

	if billingInfo.SubscriptionCancelledAt != nil {
		return &domain.CancelSubscriptionResponse{
			SubscriptionExpiresAt: billingInfo.SubscriptionExpiresAt,
			CancelledAt:           *billingInfo.SubscriptionCancelledAt,
			Message:               "Subscription is already set to cancel at period end",
		}, nil
	}

	cancelledAt := time.Now().UTC()
	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start cancellation transaction: %w", err)
	}
	defer tx.Rollback()

	if err := uc.setSubscriptionCancellationTx(tx, ctx, tenantID, cancelledAt); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit cancellation transaction: %w", err)
	}

	performedBy := optionalInt64Ptr(userID)
	uc.writeSubscriptionHistory(ctx, &domain.CreateHistoryRecordParams{
		TenantID:    tenantID,
		Action:      domain.ActionSubscriptionCancelled,
		PerformedBy: performedBy,
		OldStatus:   stringValuePtr(domain.StatusActive),
		NewStatus:   stringValuePtr(domain.StatusActive),
		Reason:      normalizeOptionalString(reason),
		IPAddress:   normalizeOptionalString(optionalStringPtr(ipAddress)),
		UserAgent:   normalizeOptionalString(optionalStringPtr(userAgent)),
		Metadata: map[string]interface{}{
			"cancelled_at":            cancelledAt.Format(time.RFC3339),
			"subscription_expires_at": billingInfo.SubscriptionExpiresAt.Format(time.RFC3339),
		},
	})

	return &domain.CancelSubscriptionResponse{
		SubscriptionExpiresAt: billingInfo.SubscriptionExpiresAt,
		CancelledAt:           cancelledAt,
		Message:               "Subscription will remain active until the current period ends",
	}, nil
}

// RefundPayment performs a full refund and reverts tenant status to trial/expired.
func (uc *PaymentUseCase) RefundPayment(
	ctx context.Context,
	paymentID int64,
	performedBy int64,
	reason *string,
	ipAddress string,
	userAgent string,
) (*domain.RefundPaymentResponse, error) {
	if err := uc.requirePaymentInfra(); err != nil {
		return nil, err
	}
	if uc.db == nil {
		return nil, fmt.Errorf("payment database is not configured")
	}
	if paymentID <= 0 {
		return nil, fmt.Errorf("invalid payment id")
	}

	tenantID, err := uc.paymentRepo.GetPaymentTenantID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve payment tenant: %w", err)
	}
	payment, err := uc.paymentRepo.GetPaymentByID(ctx, tenantID, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to load payment: %w", err)
	}
	if payment.Status == domain.PaymentStatusRefunded {
		return nil, fmt.Errorf("payment is already refunded")
	}
	if payment.Status != domain.PaymentStatusSuccess {
		return nil, fmt.Errorf("only successful payments can be refunded")
	}
	if payment.PaymobTransactionID == nil || *payment.PaymobTransactionID <= 0 {
		return nil, fmt.Errorf("paymob transaction id is missing")
	}

	authToken, err := uc.paymobClient.Authenticate()
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with paymob: %w", err)
	}

	subscriptionInfo, err := uc.subscriptionRepo.GetTenantSubscriptionInfo(ctx, payment.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to load tenant subscription: %w", err)
	}
	nextStatus := resolveStatusAfterRefund(subscriptionInfo.TrialEndDate)

	refundedAt := time.Now().UTC()
	paidAt := payment.PaidAt
	if paidAt == nil {
		paidAt = &refundedAt
	}
	claimed, err := uc.claimRefund(ctx, tenantID, payment.ID)
	if err != nil {
		return nil, err
	}
	if !claimed {
		return nil, fmt.Errorf("payment refund is already pending reconciliation")
	}
	if err := uc.paymobClient.RefundTransaction(authToken, *payment.PaymobTransactionID, payment.AmountCents); err != nil {
		if persistErr := uc.recordRefundProviderError(ctx, tenantID, payment.ID, err.Error()); persistErr != nil {
			return nil, fmt.Errorf("failed to refund paymob transaction: %w; failed to persist reconciliation error: %v", err, persistErr)
		}
		return nil, fmt.Errorf("failed to refund paymob transaction; payment remains pending reconciliation: %w", err)
	}

	var finalized bool
	err = database.WithTenantTx(ctx, uc.db, tenantID, 0, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `
			UPDATE payments
			SET status = 'refunded', refund_state = 'refunded', refund_error = NULL, updated_at = NOW()
			WHERE id = $1 AND tenant_id = $2 AND status = 'success' AND refund_state = 'pending'
		`, payment.ID, tenantID)
		if err != nil {
			return fmt.Errorf("failed to finalize provider refund: %w", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to read refund finalization result: %w", err)
		}
		if rowsAffected == 0 {
			var status, refundState string
			if err := tx.QueryRowContext(ctx, `SELECT status, refund_state FROM payments WHERE id = $1 AND tenant_id = $2`, payment.ID, tenantID).Scan(&status, &refundState); err != nil {
				return fmt.Errorf("failed to inspect refund after provider success: %w", err)
			}
			if status == domain.PaymentStatusRefunded && refundState == "refunded" {
				return errWebhookAlreadyProcessed
			}
			return fmt.Errorf("provider refund succeeded but local payment remains pending reconciliation")
		}
		if err := uc.clearSubscriptionAfterRefundTx(tx, ctx, payment.TenantID, nextStatus); err != nil {
			return err
		}
		finalized = true
		return nil
	})
	if errors.Is(err, errWebhookAlreadyProcessed) {
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("provider refund succeeded but local reconciliation failed: %w", err)
	}
	if finalized {
		uc.writeSubscriptionHistory(ctx, &domain.CreateHistoryRecordParams{
			TenantID:    payment.TenantID,
			Action:      domain.ActionPaymentRefunded,
			PerformedBy: optionalInt64Ptr(performedBy),
			OldStatus:   stringValuePtr(domain.StatusActive),
			NewStatus:   stringValuePtr(nextStatus),
			Reason:      normalizeOptionalString(reason),
			IPAddress:   normalizeOptionalString(optionalStringPtr(ipAddress)),
			UserAgent:   normalizeOptionalString(optionalStringPtr(userAgent)),
			Metadata: map[string]interface{}{
				"payment_id": payment.ID, "paymob_transaction_id": *payment.PaymobTransactionID,
				"amount_cents": payment.AmountCents, "currency": payment.Currency,
				"refunded_at": refundedAt.Format(time.RFC3339),
			},
		})
	}

	return &domain.RefundPaymentResponse{
		PaymentID:          payment.ID,
		TenantID:           payment.TenantID,
		Status:             domain.PaymentStatusRefunded,
		SubscriptionStatus: nextStatus,
		RefundedAt:         refundedAt,
		AmountCents:        payment.AmountCents,
		Currency:           payment.Currency,
	}, nil
}

func (uc *PaymentUseCase) claimRefund(ctx context.Context, tenantID, paymentID int64) (bool, error) {
	var claimed bool
	err := database.WithTenantTx(ctx, uc.db, tenantID, 0, func(tx *sql.Tx) error {
		var status, refundState string
		err := tx.QueryRowContext(ctx, `SELECT status, refund_state FROM payments WHERE id = $1 AND tenant_id = $2 FOR UPDATE`, paymentID, tenantID).Scan(&status, &refundState)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("payment not found")
			}
			return fmt.Errorf("failed to lock payment for refund: %w", err)
		}
		if status != domain.PaymentStatusSuccess || refundState != "none" {
			return fmt.Errorf("payment is no longer refundable")
		}
		result, err := tx.ExecContext(ctx, `UPDATE payments SET refund_state = 'pending', refund_error = NULL, updated_at = NOW() WHERE id = $1 AND tenant_id = $2`, paymentID, tenantID)
		if err != nil {
			return fmt.Errorf("failed to persist pending refund: %w", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to read pending refund result: %w", err)
		}
		claimed = rowsAffected == 1
		return nil
	})
	return claimed, err
}

func (uc *PaymentUseCase) recordRefundProviderError(ctx context.Context, tenantID, paymentID int64, message string) error {
	return database.WithTenantTx(ctx, uc.db, tenantID, 0, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE payments SET refund_error = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3 AND refund_state = 'pending'`, message, paymentID, tenantID)
		return err
	})
}

func (uc *PaymentUseCase) applyRefundWebhook(ctx context.Context, payment *domain.Payment, txnID *int64, cardLastFour, cardBrand *string) error {
	info, err := uc.subscriptionRepo.GetTenantSubscriptionInfo(ctx, payment.TenantID)
	if err != nil {
		return fmt.Errorf("failed to load subscription for refund callback: %w", err)
	}
	refundedAt := time.Now().UTC()
	paidAt := payment.PaidAt
	if paidAt == nil {
		paidAt = &refundedAt
	}
	nextStatus := resolveStatusAfterRefund(info.TrialEndDate)
	var previousStatus, refundState string
	err = database.WithTenantTx(ctx, uc.db, payment.TenantID, 0, func(tx *sql.Tx) error {
		if err := tx.QueryRowContext(ctx, `SELECT status, refund_state FROM payments WHERE id = $1 AND tenant_id = $2 FOR UPDATE`, payment.ID, payment.TenantID).Scan(&previousStatus, &refundState); err != nil {
			return fmt.Errorf("failed to lock payment for refund callback: %w", err)
		}
		if previousStatus == domain.PaymentStatusRefunded {
			return errWebhookAlreadyProcessed
		}
		if previousStatus != domain.PaymentStatusSuccess && previousStatus != domain.PaymentStatusPending && previousStatus != domain.PaymentStatusFailed {
			return fmt.Errorf("payment status %q cannot be reconciled as refunded", previousStatus)
		}
		if err := uc.updatePaymentStatusTx(tx, ctx, payment.TenantID, payment.ID, domain.PaymentStatusRefunded, txnID, cardLastFour, cardBrand, nil, paidAt); err != nil {
			return err
		}
		if previousStatus != domain.PaymentStatusSuccess {
			return nil
		}
		return uc.clearSubscriptionAfterRefundTx(tx, ctx, payment.TenantID, nextStatus)
	})
	if errors.Is(err, errWebhookAlreadyProcessed) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to reconcile refunded payment callback: %w", err)
	}
	if previousStatus != domain.PaymentStatusSuccess {
		return nil
	}
	uc.writeSubscriptionHistory(ctx, &domain.CreateHistoryRecordParams{
		TenantID: payment.TenantID, Action: domain.ActionPaymentRefunded, PerformedBy: payment.UserID,
		OldStatus: stringValuePtr(domain.StatusActive), NewStatus: stringValuePtr(nextStatus),
		Metadata: map[string]interface{}{"payment_id": payment.ID, "paymob_transaction_id": txnID, "refunded_at": refundedAt.Format(time.RFC3339)},
	})
	return nil
}

func (uc *PaymentUseCase) listPlans(ctx context.Context) []domain.SubscriptionPlan {
	planByID := copyPlanMap(uc.plans)

	if uc.planRepo != nil {
		repoPlans, err := uc.planRepo.GetAllActive(ctx)
		if err != nil {
			log.Printf("[PaymentUseCase] WARN: failed to load plans from DB, using fallback plans: %v", err)
		} else if len(repoPlans) > 0 {
			planByID = make(map[string]domain.SubscriptionPlan, len(repoPlans))
			for _, plan := range repoPlans {
				normalizedPlan, ok := normalizePlan(plan)
				if !ok {
					continue
				}
				planByID[normalizedPlan.ID] = normalizedPlan
			}
			if len(planByID) == 0 {
				planByID = copyPlanMap(uc.plans)
			}
		}
	}

	plans := make([]domain.SubscriptionPlan, 0, len(planByID))
	for _, plan := range planByID {
		plans = append(plans, plan)
	}
	sort.Slice(plans, func(i, j int) bool {
		if plans[i].PriceCents == plans[j].PriceCents {
			return plans[i].ID < plans[j].ID
		}
		return plans[i].PriceCents < plans[j].PriceCents
	})

	return plans
}

func (uc *PaymentUseCase) getPlanByID(ctx context.Context, planID string) (domain.SubscriptionPlan, bool) {
	targetID := strings.ToLower(strings.TrimSpace(planID))
	if targetID == "" {
		targetID = "monthly"
	}

	for _, plan := range uc.listPlans(ctx) {
		if strings.EqualFold(plan.ID, targetID) {
			return plan, true
		}
	}
	// Compatibility aliases apply only when no explicit legacy plan was supplied.
	for _, plan := range uc.listPlans(ctx) {
		if (targetID == "monthly" && plan.ID == commerce.Current().MonthlyPlanID) || (targetID == "annual" && plan.ID == commerce.Current().AnnualPlanID) {
			return plan, true
		}
	}
	return domain.SubscriptionPlan{}, false
}

func (uc *PaymentUseCase) defaultPlan() domain.SubscriptionPlan {
	monthlyPlan, ok := uc.plans["monthly"]
	if ok {
		return monthlyPlan
	}
	for _, plan := range uc.plans {
		return plan
	}
	return domain.SubscriptionPlan{
		ID:         "monthly",
		Name:       "Pro Plan",
		NameAr:     "Pro Plan",
		Price:      commerce.Current().MonthlyPriceCents / 100,
		PriceCents: commerce.Current().MonthlyPriceCents,
		Currency:   "EGP",
		Period:     "month",
	}
}

func (uc *PaymentUseCase) markPaymentFailed(ctx context.Context, tenantID, paymentID int64, message string) error {
	msg := strings.TrimSpace(message)
	if msg == "" {
		msg = "payment initiation failed"
	}
	if uc.paymentRepo == nil {
		return fmt.Errorf("payment repository is not configured")
	}
	return uc.paymentRepo.UpdatePaymentStatus(ctx, tenantID, paymentID, domain.PaymentStatusFailed, nil, nil, nil, &msg, nil)
}

func (uc *PaymentUseCase) requirePaymentInfra() error {
	if uc.paymobClient == nil {
		return fmt.Errorf("payment client is not configured")
	}
	if uc.paymentRepo == nil {
		return fmt.Errorf("payment repository is not configured")
	}
	if uc.subscriptionRepo == nil {
		return fmt.Errorf("subscription repository is not configured")
	}
	return nil
}

func (uc *PaymentUseCase) updatePaymentStatusTx(
	tx *sql.Tx,
	ctx context.Context,
	tenantID int64,
	id int64,
	status string,
	txnID *int64,
	cardLastFour *string,
	cardBrand *string,
	errorMsg *string,
	paidAt *time.Time,
) error {
	query := `
		UPDATE payments
		SET status = $1,
			paymob_transaction_id = COALESCE($2, paymob_transaction_id),
			card_last_four = COALESCE($3, card_last_four),
			card_brand = COALESCE($4, card_brand),
			error_message = $5,
			paid_at = $6,
			updated_at = NOW()
		WHERE id = $7 AND tenant_id = $8
	`
	args := []interface{}{status, txnID, cardLastFour, cardBrand, errorMsg, paidAt, id, tenantID}
	switch status {
	case domain.PaymentStatusSuccess:
		query += ` AND status IN ('pending', 'failed')`
	case domain.PaymentStatusFailed:
		query += ` AND status = 'pending'`
	case domain.PaymentStatusVoided:
		query += ` AND status IN ('pending', 'failed')`
	case domain.PaymentStatusRefunded:
		query = `
			UPDATE payments
			SET status = $1, refund_state = 'refunded', refund_error = NULL,
				paymob_transaction_id = COALESCE($2, paymob_transaction_id),
				card_last_four = COALESCE($3, card_last_four),
				card_brand = COALESCE($4, card_brand), error_message = $5,
				paid_at = $6, updated_at = NOW()
			WHERE id = $7 AND tenant_id = $8 AND status IN ('pending', 'failed', 'success')
		`
	default:
		return fmt.Errorf("unsupported transactional payment status: %s", status)
	}
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read affected rows: %w", err)
	}
	if rowsAffected == 0 {
		var currentStatus string
		if err := tx.QueryRowContext(ctx, `SELECT status FROM payments WHERE id = $1 AND tenant_id = $2`, id, tenantID).Scan(&currentStatus); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("payment not found")
			}
			return fmt.Errorf("failed to read payment status after refund update: %w", err)
		}
		if currentStatus == domain.PaymentStatusRefunded || currentStatus == domain.PaymentStatusSuccess || currentStatus == domain.PaymentStatusVoided || currentStatus == domain.PaymentStatusFailed {
			return errWebhookAlreadyProcessed
		}
		return fmt.Errorf("payment status changed before webhook transition")
	}

	return nil
}

func (uc *PaymentUseCase) updateTenantSubscriptionTx(
	tx *sql.Tx,
	ctx context.Context,
	tenantID int64,
	status string,
	expiresAt time.Time,
	lastPaymentAt time.Time,
	lastPaymentAmount int,
	currency string,
	planID string,
) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE tenants
		SET subscription_status = $1,
			subscription_plan = $2,
			subscription_expires_at = $3,
			last_payment_at = $4,
			last_payment_amount = $5,
			last_payment_currency = $6,
			updated_at = NOW()
		WHERE id = $7
	`, status, planID, expiresAt, lastPaymentAt, lastPaymentAmount, currency, tenantID)
	if err != nil {
		return fmt.Errorf("failed to update tenant subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tenant not found")
	}

	return nil
}

func (uc *PaymentUseCase) setSubscriptionCancellationTx(
	tx *sql.Tx,
	ctx context.Context,
	tenantID int64,
	cancelledAt time.Time,
) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE tenants
		SET subscription_cancelled_at = $1,
			updated_at = NOW()
		WHERE id = $2
		  AND subscription_status = $3
	`, cancelledAt, tenantID, domain.StatusActive)
	if err != nil {
		return fmt.Errorf("failed to set subscription cancellation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tenant not found or subscription is not active")
	}
	return nil
}

func (uc *PaymentUseCase) clearSubscriptionAfterRefundTx(
	tx *sql.Tx,
	ctx context.Context,
	tenantID int64,
	nextStatus string,
) error {
	query := `
		UPDATE tenants
		SET subscription_status = $1,
			subscription_expires_at = NULL,
			subscription_cancelled_at = NULL,
			updated_at = NOW()
		WHERE id = $2
	`
	if nextStatus == domain.StatusTrial {
		query = `
			UPDATE tenants
			SET subscription_status = $1,
				subscription_plan = 'trial',
				subscription_expires_at = NULL,
				subscription_cancelled_at = NULL,
				updated_at = NOW()
			WHERE id = $2
		`
	}
	result, err := tx.ExecContext(ctx, query, nextStatus, tenantID)
	if err != nil {
		return fmt.Errorf("failed to update tenant status after refund: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tenant not found")
	}
	return nil
}

func (uc *PaymentUseCase) normalizeBillingData(input map[string]string) map[string]string {
	output := map[string]string{
		"email":           "unknown@example.com",
		"first_name":      "NA",
		"last_name":       "NA",
		"phone_number":    "+200000000000",
		"apartment":       "NA",
		"floor":           "NA",
		"building":        "NA",
		"street":          "NA",
		"city":            "Cairo",
		"state":           "Cairo",
		"country":         "EG",
		"shipping_method": "NA",
		"postal_code":     "NA",
	}

	for key, value := range input {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			output[key] = trimmed
		}
	}

	return output
}

func (uc *PaymentUseCase) writeSubscriptionHistory(ctx context.Context, params *domain.CreateHistoryRecordParams) {
	if uc.subscriptionRepo == nil || params == nil {
		return
	}
	if err := uc.subscriptionRepo.CreateHistoryRecord(ctx, params); err != nil {
		log.Printf("[PaymentUseCase] WARN: failed to write subscription history action %s: %v\n", params.Action, err)
	}
}

func buildDefaultPlans() map[string]domain.SubscriptionPlan {
	offer := commerce.Current()
	features := []string{"Cashier POS", "Kitchen display", "Cash ordering", "Theme configuration"}
	featuresAr := []string{"نقاط البيع", "شاشة المطبخ", "الطلبات النقدية", "إعداد القالب"}
	return map[string]domain.SubscriptionPlan{
		"monthly": {ID: offer.MonthlyPlanID, Name: "SERVIO Monthly", NameAr: "سيرفيو شهري", Price: offer.MonthlyPriceCents / 100, PriceCents: offer.MonthlyPriceCents, Currency: offer.Currency, Period: "month", Features: features, FeaturesAr: featuresAr},
		"annual":  {ID: offer.AnnualPlanID, Name: "SERVIO Annual", NameAr: "سيرفيو سنوي", Price: offer.AnnualPriceCents / 100, PriceCents: offer.AnnualPriceCents, Currency: offer.Currency, Period: "year", Features: features, FeaturesAr: featuresAr},
	}
}

func copyPlanMap(input map[string]domain.SubscriptionPlan) map[string]domain.SubscriptionPlan {
	output := make(map[string]domain.SubscriptionPlan, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

func normalizePlan(plan domain.SubscriptionPlan) (domain.SubscriptionPlan, bool) {
	plan.ID = strings.ToLower(strings.TrimSpace(plan.ID))
	if plan.ID == "" {
		return domain.SubscriptionPlan{}, false
	}

	if plan.PriceCents <= 0 && plan.Price > 0 {
		plan.PriceCents = plan.Price * 100
	}
	if plan.Price <= 0 && plan.PriceCents > 0 {
		plan.Price = plan.PriceCents / 100
	}

	plan.Currency = strings.ToUpper(strings.TrimSpace(plan.Currency))
	if plan.Currency == "" {
		plan.Currency = "EGP"
	}

	plan.Period = strings.ToLower(strings.TrimSpace(plan.Period))
	if plan.Period == "" {
		plan.Period = "month"
	}

	if plan.Features == nil {
		plan.Features = []string{}
	}
	if plan.FeaturesAr == nil {
		plan.FeaturesAr = []string{}
	}

	return plan, true
}

func toPlanSummary(plan domain.SubscriptionPlan) domain.SubscriptionPlanSummary {
	return domain.SubscriptionPlanSummary{
		ID:         plan.ID,
		Name:       plan.Name,
		Price:      plan.Price,
		PriceCents: plan.PriceCents,
		Currency:   plan.Currency,
		Period:     plan.Period,
	}
}

func periodMonthsFromPlan(plan domain.SubscriptionPlan) int {
	period := strings.ToLower(strings.TrimSpace(plan.Period))
	switch period {
	case "year", "annual":
		return 12
	case "month", "monthly":
		return 1
	}

	if strings.HasSuffix(period, "-month") {
		parts := strings.Split(period, "-")
		if len(parts) > 0 {
			months, err := strconv.Atoi(parts[0])
			if err == nil && months > 0 {
				return months
			}
		}
	}

	return 1
}

func resolveStatusAfterRefund(trialEndDate *time.Time) string {
	now := time.Now().UTC()
	if trialEndDate != nil && trialEndDate.After(now) {
		return domain.StatusTrial
	}
	return domain.StatusExpired
}

func optionalInt64Ptr(value int64) *int64 {
	if value <= 0 {
		return nil
	}
	return &value
}

func optionalStringPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func extractLastFour(pan string) string {
	cleaned := strings.TrimSpace(pan)
	if cleaned == "" {
		return ""
	}

	cleaned = strings.ReplaceAll(cleaned, "-", "")
	if len(cleaned) < 4 {
		return cleaned
	}
	return cleaned[len(cleaned)-4:]
}

func daysUntil(target time.Time) int {
	diff := target.Sub(time.Now()).Hours() / 24
	return int(math.Ceil(diff))
}

func getEnvString(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func stringValuePtr(value string) *string {
	return &value
}
