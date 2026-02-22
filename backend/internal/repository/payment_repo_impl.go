package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"pos-saas/internal/domain"
)

// PaymentRepositoryImpl implements PaymentRepository.
type PaymentRepositoryImpl struct {
	db *sql.DB
}

// NewPaymentRepository creates a payment repository.
func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &PaymentRepositoryImpl{db: db}
}

// CreatePayment creates a new payment row.
func (r *PaymentRepositoryImpl) CreatePayment(ctx context.Context, payment *domain.Payment) (int64, error) {
	log.Printf("[PaymentRepo] CreatePayment - tenant_id: %d, status: %s\n", payment.TenantID, payment.Status)

	metadataJSON := "{}"
	if payment.Metadata != nil {
		if bytes, err := json.Marshal(payment.Metadata); err == nil {
			metadataJSON = string(bytes)
		}
	}

	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO payments (
			tenant_id, user_id, amount_cents, currency, status,
			paymob_order_id, paymob_auth_token,
			payment_method, card_last_four, card_brand,
			billing_name, billing_email, billing_phone,
			plan_id, plan_period_start, plan_period_end,
			error_message, error_code, merchant_order_id, metadata,
			created_at, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7,
			$8, $9, $10,
			$11, $12, $13,
			$14, $15, $16,
			$17, $18, $19, $20::jsonb,
			NOW(), NOW()
		)
		RETURNING id
	`,
		payment.TenantID,
		payment.UserID,
		payment.AmountCents,
		payment.Currency,
		payment.Status,
		payment.PaymobOrderID,
		nil,
		payment.PaymentMethod,
		payment.CardLastFour,
		payment.CardBrand,
		payment.BillingName,
		payment.BillingEmail,
		payment.BillingPhone,
		payment.PlanID,
		payment.PlanPeriodStart,
		payment.PlanPeriodEnd,
		payment.ErrorMessage,
		payment.ErrorCode,
		payment.MerchantOrderID,
		metadataJSON,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create payment: %w", err)
	}

	return id, nil
}

// GetPaymentByID returns payment by internal ID.
func (r *PaymentRepositoryImpl) GetPaymentByID(ctx context.Context, id int64) (*domain.Payment, error) {
	row := r.db.QueryRowContext(ctx, selectPaymentQuery+` WHERE id = $1`, id)
	payment, err := scanPayment(row)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// GetPaymentByMerchantOrderID returns payment by merchant order ID.
func (r *PaymentRepositoryImpl) GetPaymentByMerchantOrderID(ctx context.Context, merchantOrderID string) (*domain.Payment, error) {
	row := r.db.QueryRowContext(ctx, selectPaymentQuery+` WHERE merchant_order_id = $1`, merchantOrderID)
	payment, err := scanPayment(row)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// GetPaymentByPaymobOrderID returns payment by Paymob order ID.
func (r *PaymentRepositoryImpl) GetPaymentByPaymobOrderID(ctx context.Context, paymobOrderID int64) (*domain.Payment, error) {
	row := r.db.QueryRowContext(ctx, selectPaymentQuery+` WHERE paymob_order_id = $1 ORDER BY created_at DESC LIMIT 1`, paymobOrderID)
	payment, err := scanPayment(row)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// UpdatePaymentOrder stores Paymob order/auth refs after initiation.
func (r *PaymentRepositoryImpl) UpdatePaymentOrder(ctx context.Context, id int64, paymobOrderID int64, authToken string) error {
	log.Printf("[PaymentRepo] UpdatePaymentOrder - payment_id: %d, paymob_order_id: %d\n", id, paymobOrderID)

	result, err := r.db.ExecContext(ctx, `
		UPDATE payments
		SET paymob_order_id = $1,
			paymob_auth_token = $2,
			updated_at = NOW()
		WHERE id = $3
	`, paymobOrderID, authToken, id)
	if err != nil {
		return fmt.Errorf("failed to update payment order details: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// UpdatePaymentStatus updates payment status and transaction details.
func (r *PaymentRepositoryImpl) UpdatePaymentStatus(ctx context.Context, id int64, status string, txnID *int64, cardLastFour *string, cardBrand *string, errorMsg *string, paidAt *time.Time) error {
	log.Printf("[PaymentRepo] UpdatePaymentStatus - payment_id: %d, status: %s\n", id, status)

	result, err := r.db.ExecContext(ctx, `
		UPDATE payments
		SET status = $1,
			paymob_transaction_id = COALESCE($2, paymob_transaction_id),
			card_last_four = COALESCE($3, card_last_four),
			card_brand = COALESCE($4, card_brand),
			error_message = $5,
			paid_at = $6,
			updated_at = NOW()
		WHERE id = $7
	`, status, txnID, cardLastFour, cardBrand, errorMsg, paidAt, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// GetPaymentsByTenantID returns paginated payments and total count.
func (r *PaymentRepositoryImpl) GetPaymentsByTenantID(ctx context.Context, tenantID int64, limit, offset int) ([]domain.Payment, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM payments
		WHERE tenant_id = $1
	`, tenantID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT
		`+selectPaymentColumns+`
		FROM payments
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, tenantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query payments: %w", err)
	}
	defer rows.Close()

	payments := make([]domain.Payment, 0, limit)
	for rows.Next() {
		payment, err := scanPayment(rows)
		if err != nil {
			return nil, 0, err
		}
		payments = append(payments, *payment)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate payments: %w", err)
	}

	return payments, total, nil
}

// GetPaymentsByTenantCursor returns payments using cursor pagination ordered by newest first.
func (r *PaymentRepositoryImpl) GetPaymentsByTenantCursor(ctx context.Context, tenantID int64, cursorID int64, limit int) ([]domain.Payment, *int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if cursorID < 0 {
		cursorID = 0
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT
		`+selectPaymentColumns+`
		FROM payments
		WHERE tenant_id = $1
		  AND ($2 = 0 OR id < $2)
		ORDER BY id DESC
		LIMIT $3
	`, tenantID, cursorID, limit+1)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query payments by cursor: %w", err)
	}
	defer rows.Close()

	payments := make([]domain.Payment, 0, limit+1)
	for rows.Next() {
		payment, err := scanPayment(rows)
		if err != nil {
			return nil, nil, err
		}
		payments = append(payments, *payment)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to iterate payments by cursor: %w", err)
	}

	var nextCursor *int64
	if len(payments) > limit {
		lastReturned := payments[limit-1].ID
		nextCursor = &lastReturned
		payments = payments[:limit]
	}

	return payments, nextCursor, nil
}

// UpdateTenantSubscription updates paid subscription fields on tenant row.
func (r *PaymentRepositoryImpl) UpdateTenantSubscription(ctx context.Context, tenantID int64, status string, expiresAt time.Time, lastPaymentAt time.Time, lastPaymentAmount int, currency string, planID string) error {
	log.Printf("[PaymentRepo] UpdateTenantSubscription - tenant_id: %d, status: %s\n", tenantID, status)

	result, err := r.db.ExecContext(ctx, `
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

// GetTenantBillingInfo returns payment-related subscription fields from tenants.
func (r *PaymentRepositoryImpl) GetTenantBillingInfo(ctx context.Context, tenantID int64) (*domain.TenantBillingInfo, error) {
	info := &domain.TenantBillingInfo{
		SubscriptionPlan:    "monthly",
		LastPaymentCurrency: "EGP",
	}

	var subscriptionPlan sql.NullString
	var expiresAt sql.NullTime
	var cancelledAt sql.NullTime
	var lastPaymentAt sql.NullTime
	var lastAmount sql.NullInt64
	var lastCurrency sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT subscription_plan, subscription_expires_at, subscription_cancelled_at, last_payment_at, last_payment_amount, last_payment_currency
		FROM tenants
		WHERE id = $1
	`, tenantID).Scan(&subscriptionPlan, &expiresAt, &cancelledAt, &lastPaymentAt, &lastAmount, &lastCurrency)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tenant not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant billing info: %w", err)
	}

	if subscriptionPlan.Valid && subscriptionPlan.String != "" {
		info.SubscriptionPlan = subscriptionPlan.String
	}
	if expiresAt.Valid {
		info.SubscriptionExpiresAt = &expiresAt.Time
	}
	if cancelledAt.Valid {
		info.SubscriptionCancelledAt = &cancelledAt.Time
	}
	if lastPaymentAt.Valid {
		info.LastPaymentAt = &lastPaymentAt.Time
	}
	if lastAmount.Valid {
		value := int(lastAmount.Int64)
		info.LastPaymentAmount = &value
	}
	if lastCurrency.Valid && lastCurrency.String != "" {
		info.LastPaymentCurrency = lastCurrency.String
	}

	return info, nil
}

const selectPaymentColumns = `
			id, tenant_id, user_id, amount_cents, currency, status,
			paymob_order_id, paymob_transaction_id, payment_method, card_last_four, card_brand,
			billing_name, billing_email, billing_phone, plan_id, plan_period_start, plan_period_end,
			merchant_order_id, error_message, error_code, metadata::text, created_at, updated_at, paid_at
`

const selectPaymentQuery = `
		SELECT
` + selectPaymentColumns + `
		FROM payments
`

type paymentScanner interface {
	Scan(dest ...interface{}) error
}

func scanPayment(scanner paymentScanner) (*domain.Payment, error) {
	var payment domain.Payment
	var userID sql.NullInt64
	var paymobOrderID sql.NullInt64
	var paymobTransactionID sql.NullInt64
	var paymentMethod sql.NullString
	var cardLastFour sql.NullString
	var cardBrand sql.NullString
	var billingName sql.NullString
	var billingEmail sql.NullString
	var billingPhone sql.NullString
	var planPeriodStart sql.NullTime
	var planPeriodEnd sql.NullTime
	var errorMessage sql.NullString
	var errorCode sql.NullString
	var metadataRaw sql.NullString
	var paidAt sql.NullTime

	err := scanner.Scan(
		&payment.ID,
		&payment.TenantID,
		&userID,
		&payment.AmountCents,
		&payment.Currency,
		&payment.Status,
		&paymobOrderID,
		&paymobTransactionID,
		&paymentMethod,
		&cardLastFour,
		&cardBrand,
		&billingName,
		&billingEmail,
		&billingPhone,
		&payment.PlanID,
		&planPeriodStart,
		&planPeriodEnd,
		&payment.MerchantOrderID,
		&errorMessage,
		&errorCode,
		&metadataRaw,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&paidAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan payment row: %w", err)
	}

	if userID.Valid {
		value := userID.Int64
		payment.UserID = &value
	}
	if paymobOrderID.Valid {
		value := paymobOrderID.Int64
		payment.PaymobOrderID = &value
	}
	if paymobTransactionID.Valid {
		value := paymobTransactionID.Int64
		payment.PaymobTransactionID = &value
	}
	if paymentMethod.Valid {
		value := paymentMethod.String
		payment.PaymentMethod = &value
	}
	if cardLastFour.Valid {
		value := cardLastFour.String
		payment.CardLastFour = &value
	}
	if cardBrand.Valid {
		value := cardBrand.String
		payment.CardBrand = &value
	}
	if billingName.Valid {
		value := billingName.String
		payment.BillingName = &value
	}
	if billingEmail.Valid {
		value := billingEmail.String
		payment.BillingEmail = &value
	}
	if billingPhone.Valid {
		value := billingPhone.String
		payment.BillingPhone = &value
	}
	if planPeriodStart.Valid {
		value := planPeriodStart.Time
		payment.PlanPeriodStart = &value
	}
	if planPeriodEnd.Valid {
		value := planPeriodEnd.Time
		payment.PlanPeriodEnd = &value
	}
	if errorMessage.Valid {
		value := errorMessage.String
		payment.ErrorMessage = &value
	}
	if errorCode.Valid {
		value := errorCode.String
		payment.ErrorCode = &value
	}
	if paidAt.Valid {
		value := paidAt.Time
		payment.PaidAt = &value
	}
	if metadataRaw.Valid && metadataRaw.String != "" {
		_ = json.Unmarshal([]byte(metadataRaw.String), &payment.Metadata)
	}

	return &payment, nil
}
