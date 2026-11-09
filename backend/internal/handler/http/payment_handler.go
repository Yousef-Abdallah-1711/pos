package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/service"
	"pos-saas/internal/usecase"
)

// PaymentHandler handles subscription payment endpoints.
type PaymentHandler struct {
	paymentUseCase *usecase.PaymentUseCase
	webhookMonitor *service.PaymentWebhookMonitor
	retryQueue     *service.WebhookRetryQueue
	errorTracker   service.ErrorTracker
}

// NewPaymentHandler creates a payment handler.
func NewPaymentHandler(paymentUseCase *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
		errorTracker:   service.NewNoopTracker(),
	}
}

// WithObservability wires optional monitoring/reliability dependencies.
func (h *PaymentHandler) WithObservability(
	webhookMonitor *service.PaymentWebhookMonitor,
	retryQueue *service.WebhookRetryQueue,
	errorTracker service.ErrorTracker,
) *PaymentHandler {
	h.webhookMonitor = webhookMonitor
	h.retryQueue = retryQueue
	if errorTracker != nil {
		h.errorTracker = errorTracker
	}
	return h
}

// GetPlans handles GET /api/v1/subscription/plans.
func (h *PaymentHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PaymentHandler] GetPlans called\n")
	plans := h.paymentUseCase.GetPlans()
	OKResponse(w, map[string]interface{}{
		"plans": plans,
	})
}

// InitiatePayment handles POST /api/v1/subscription/initiate.
// This route must be auth-protected but not blocked by tenant status middleware.
func (h *PaymentHandler) InitiatePayment(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PaymentHandler] InitiatePayment called\n")

	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		BadRequestResponse(w, "tenant context not found")
		return
	}

	userID := middleware.GetUserID(r)
	// Impersonation requests may not carry a real dashboard user id that exists
	// in this tenant DB; keep payment.user_id nullable in that case.
	if strings.TrimSpace(r.Header.Get("X-Impersonation-Token")) != "" {
		userID = 0
	}

	var req domain.InitiatePaymentRequest
	if !DecodeJSONBody(w, r, &req) {
		return
	}

	if strings.TrimSpace(req.PlanID) == "" {
		req.PlanID = "monthly"
	}

	claims := middleware.GetUserClaims(r)
	billingData := map[string]string{
		"first_name": "NA",
		"last_name":  "NA",
	}
	if claims != nil {
		billingData["email"] = claims.Email
	}

	response, err := h.paymentUseCase.InitiatePayment(r.Context(), tenantID, userID, req.PlanID, billingData)
	if err != nil {
		log.Printf("[PaymentHandler] InitiatePayment error: %v\n", err)

		errMsg := strings.ToLower(strings.TrimSpace(err.Error()))
		if strings.Contains(errMsg, "invalid plan") ||
			strings.Contains(errMsg, "invalid tenant") ||
			strings.Contains(errMsg, "not configured") {
			BadRequestResponse(w, err.Error())
		} else {
			InternalErrorResponse(w, "failed to initiate payment")
		}
		return
	}

	OKResponse(w, response)
}

// HandlePaymobWebhook handles POST /api/v1/webhook/paymob.
// Security policy:
//   - 401 for invalid HMAC signature
//   - 400 for invalid webhook payload fields
//   - 500 for unexpected internal failures
func (h *PaymentHandler) HandlePaymobWebhook(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PaymentHandler] HandlePaymobWebhook called\n")
	startedAt := time.Now()

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[PaymentHandler] webhook body read error: %v\n", err)
		h.trackWebhookEvent(r, 0, "invalid_body", false, time.Since(startedAt), err.Error(), nil)
		OKResponse(w, map[string]interface{}{"received": true})
		return
	}

	var payload domain.PaymobWebhookPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		log.Printf("[PaymentHandler] webhook payload decode error: %v\n", err)
		h.trackWebhookEvent(r, 0, "invalid_payload_json", false, time.Since(startedAt), err.Error(), nil)
		OKResponse(w, map[string]interface{}{"received": true})
		return
	}

	paymobOrderID := payload.Obj.Order.ID
	receivedHMAC := strings.TrimSpace(r.URL.Query().Get("hmac"))
	if receivedHMAC == "" {
		receivedHMAC = strings.TrimSpace(payload.HMAC)
	}

	if err := h.paymentUseCase.HandleWebhook(r.Context(), rawBody, &payload, receivedHMAC); err != nil {
		log.Printf("[PaymentHandler] webhook processing error: %v\n", err)
		errMsg := strings.ToLower(strings.TrimSpace(err.Error()))
		processing := time.Since(startedAt)
		isRetryAttempt := strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Webhook-Retry-Attempt")), "true")

		switch {
		case strings.Contains(errMsg, "invalid hmac signature"):
			log.Printf("[PaymentWebhook] hmac_invalid order_id=%d processing_time_ms=%d source=%s",
				paymobOrderID,
				processing.Milliseconds(),
				webhookSource(r),
			)
			h.errorTracker.CaptureMessage("payment webhook invalid hmac signature", "warning")
			h.trackWebhookEvent(r, paymobOrderID, "hmac_invalid", false, processing, err.Error(), nil)
			UnauthorizedResponse(w, "invalid webhook signature")
		case strings.Contains(errMsg, "amount mismatch"),
			strings.Contains(errMsg, "currency mismatch"),
			strings.Contains(errMsg, "missing paymob order id"),
			strings.Contains(errMsg, "payload is required"):
			log.Printf("[PaymentWebhook] validation_rejected order_id=%d processing_time_ms=%d source=%s error=%s",
				paymobOrderID,
				processing.Milliseconds(),
				webhookSource(r),
				strings.TrimSpace(err.Error()),
			)
			h.trackWebhookEvent(r, paymobOrderID, "validation_rejected", true, processing, err.Error(), nil)
			BadRequestResponse(w, "invalid webhook payload")
		default:
			var retryRecordID *int64
			endpoint := r.URL.Path
			if r.URL.RawQuery != "" {
				endpoint += "?" + r.URL.RawQuery
			}
			if !isRetryAttempt && h.retryQueue != nil {
				id, enqueueErr := h.retryQueue.Enqueue(r.Context(), endpoint, rawBody, err.Error())
				if enqueueErr != nil {
					log.Printf("[PaymentWebhook] retry enqueue failed order_id=%d error=%v", paymobOrderID, enqueueErr)
					h.errorTracker.CaptureError(enqueueErr, map[string]interface{}{
						"component":       "payment_webhook",
						"action":          "enqueue_retry",
						"paymob_order_id": paymobOrderID,
					})
				} else {
					retryRecordID = &id
				}
			}

			log.Printf("[PaymentWebhook] processing_failed order_id=%d processing_time_ms=%d source=%s retry_record_id=%v error=%s",
				paymobOrderID,
				processing.Milliseconds(),
				webhookSource(r),
				retryRecordID,
				strings.TrimSpace(err.Error()),
			)
			h.errorTracker.CaptureError(err, map[string]interface{}{
				"component":       "payment_webhook",
				"paymob_order_id": paymobOrderID,
				"source":          webhookSource(r),
				"retry_attempt":   isRetryAttempt,
			})
			h.trackWebhookEvent(r, paymobOrderID, "processing_failed", true, processing, err.Error(), retryRecordID)
			InternalErrorResponse(w, "failed to process webhook")
		}
		return
	}

	status := "success"
	if !payload.Obj.Success || payload.Obj.Pending || payload.Obj.ErrorOccured {
		status = "failed"
	}
	processing := time.Since(startedAt)
	log.Printf("[PaymentWebhook] processed order_id=%d status=%s hmac_valid=true processing_time_ms=%d source=%s",
		paymobOrderID,
		status,
		processing.Milliseconds(),
		webhookSource(r),
	)
	h.trackWebhookEvent(r, paymobOrderID, status, true, processing, "", nil)

	OKResponse(w, map[string]interface{}{"received": true})
}

// GetPaymentHistory handles GET /api/v1/subscription/payments.
func (h *PaymentHandler) GetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PaymentHandler] GetPaymentHistory called\n")

	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		BadRequestResponse(w, "tenant context not found")
		return
	}

	cursorRaw := strings.TrimSpace(r.URL.Query().Get("cursor"))
	if cursorRaw != "" {
		cursorID, err := strconv.ParseInt(cursorRaw, 10, 64)
		if err != nil || cursorID < 0 {
			BadRequestResponse(w, "invalid cursor")
			return
		}

		limit := GetQueryParamInt(r, "limit", 20)
		payments, nextCursor, err := h.paymentUseCase.GetPaymentHistoryByCursor(r.Context(), tenantID, cursorID, limit)
		if err != nil {
			log.Printf("[PaymentHandler] GetPaymentHistory (cursor) error: %v\n", err)
			InternalErrorResponse(w, "failed to fetch payment history")
			return
		}

		response := map[string]interface{}{
			"payments": payments,
			"cursor":   cursorID,
			"limit":    limit,
			"has_more": nextCursor != nil,
		}
		if nextCursor != nil {
			response["next_cursor"] = *nextCursor
		}
		OKResponse(w, response)
		return
	}

	page := GetQueryParamInt(r, "page", 1)
	pageSize := GetQueryParamInt(r, "page_size", 20)

	payments, total, err := h.paymentUseCase.GetPaymentHistory(r.Context(), tenantID, page, pageSize)
	if err != nil {
		log.Printf("[PaymentHandler] GetPaymentHistory error: %v\n", err)
		InternalErrorResponse(w, "failed to fetch payment history")
		return
	}

	OKResponse(w, map[string]interface{}{
		"payments":  payments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetCurrentSubscription handles GET /api/v1/subscription/current.
func (h *PaymentHandler) GetCurrentSubscription(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PaymentHandler] GetCurrentSubscription called\n")

	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		BadRequestResponse(w, "tenant context not found")
		return
	}

	response, err := h.paymentUseCase.GetCurrentSubscription(r.Context(), tenantID)
	if err != nil {
		log.Printf("[PaymentHandler] GetCurrentSubscription error: %v\n", err)
		InternalErrorResponse(w, "failed to fetch current subscription")
		return
	}

	OKResponse(w, response)
}

// CancelSubscription handles POST /api/v1/subscription/cancel.
func (h *PaymentHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PaymentHandler] CancelSubscription called\n")

	tenantID := middleware.GetTenantID(r)
	if tenantID == 0 {
		BadRequestResponse(w, "tenant context not found")
		return
	}

	userID := middleware.GetUserID(r)
	var req domain.CancelSubscriptionRequest
	if !DecodeJSONBody(w, r, &req) {
		return
	}

	response, err := h.paymentUseCase.CancelSubscription(
		r.Context(),
		tenantID,
		userID,
		normalizeOptionalString(optionalStringPtr(req.Reason)),
		getPaymentRequestIPAddress(r),
		r.UserAgent(),
	)
	if err != nil {
		log.Printf("[PaymentHandler] CancelSubscription error: %v\n", err)
		errMsg := strings.ToLower(strings.TrimSpace(err.Error()))
		switch {
		case strings.Contains(errMsg, "only active subscriptions"),
			strings.Contains(errMsg, "missing"),
			strings.Contains(errMsg, "invalid tenant"):
			ConflictResponse(w, err.Error())
		case strings.Contains(errMsg, "not configured"):
			BadRequestResponse(w, err.Error())
		default:
			InternalErrorResponse(w, "failed to cancel subscription")
		}
		return
	}

	OKResponse(w, response)
}

// RefundPayment handles POST /api/v1/super-admin/payments/{id}/refund.
func (h *PaymentHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PaymentHandler] RefundPayment called\n")

	superAdminID := middleware.GetSuperAdminID(r)
	if superAdminID == 0 {
		UnauthorizedResponse(w, "super admin not authenticated")
		return
	}

	paymentID, err := parsePathInt64(r, "id")
	if err != nil {
		BadRequestResponse(w, "invalid payment id")
		return
	}

	var req struct {
		Reason string `json:"reason,omitempty"`
	}
	if !DecodeJSONBody(w, r, &req) {
		return
	}

	response, err := h.paymentUseCase.RefundPayment(
		r.Context(),
		paymentID,
		superAdminID,
		normalizeOptionalString(optionalStringPtr(req.Reason)),
		getPaymentRequestIPAddress(r),
		r.UserAgent(),
	)
	if err != nil {
		log.Printf("[PaymentHandler] RefundPayment error: %v\n", err)
		errMsg := strings.ToLower(strings.TrimSpace(err.Error()))
		switch {
		case strings.Contains(errMsg, "only successful payments"),
			strings.Contains(errMsg, "already refunded"),
			strings.Contains(errMsg, "no longer refundable"):
			ConflictResponse(w, err.Error())
		case strings.Contains(errMsg, "invalid payment id"),
			strings.Contains(errMsg, "missing"),
			strings.Contains(errMsg, "not found"):
			BadRequestResponse(w, err.Error())
		default:
			InternalErrorResponse(w, "failed to refund payment")
		}
		return
	}

	OKResponse(w, response)
}

func (h *PaymentHandler) trackWebhookEvent(
	r *http.Request,
	paymobOrderID int64,
	status string,
	hmacValid bool,
	processing time.Duration,
	errorMsg string,
	retryQueueRecordID *int64,
) {
	if h.webhookMonitor == nil {
		return
	}

	h.webhookMonitor.Record(service.PaymentWebhookEvent{
		Timestamp:          time.Now().UTC(),
		PaymobOrderID:      paymobOrderID,
		Status:             status,
		HMACValid:          hmacValid,
		ProcessingTimeMS:   processing.Milliseconds(),
		Source:             webhookSource(r),
		Error:              strings.TrimSpace(errorMsg),
		RetryQueueRecordID: retryQueueRecordID,
	})
}

func webhookSource(r *http.Request) string {
	if r == nil {
		return "webhook"
	}
	if strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Webhook-Retry-Attempt")), "true") {
		return "retry_worker"
	}
	return "webhook"
}

func parsePathInt64(r *http.Request, key string) (int64, error) {
	raw := strings.TrimSpace(r.PathValue(key))
	if raw == "" {
		return 0, http.ErrMissingFile
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func getPaymentRequestIPAddress(r *http.Request) string {
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		return xff
	}
	if xri := strings.TrimSpace(r.Header.Get("X-Real-IP")); xri != "" {
		return xri
	}
	return strings.TrimSpace(r.RemoteAddr)
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
