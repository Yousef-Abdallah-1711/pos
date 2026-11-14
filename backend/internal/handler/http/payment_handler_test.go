//go:build legacy
// +build legacy

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/middleware"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/usecase"
)

type handlerTestPaymobClient struct {
	orderID    int64
	paymentKey string
	verifyHMAC bool
}

func (c *handlerTestPaymobClient) Authenticate() (string, error) {
	return "handler_auth_token", nil
}

func (c *handlerTestPaymobClient) CreateOrder(authToken string, amountCents int, currency string, merchantOrderID string) (int64, error) {
	if c.orderID <= 0 {
		c.orderID = 123456
	}
	return c.orderID, nil
}

func (c *handlerTestPaymobClient) GetPaymentKey(authToken string, orderID int64, amountCents int, currency string, integrationID int, billingData map[string]string, expiration int) (string, error) {
	if strings.TrimSpace(c.paymentKey) == "" {
		return "handler_payment_key", nil
	}
	return c.paymentKey, nil
}

func (c *handlerTestPaymobClient) RefundTransaction(authToken string, transactionID int64, amountCents int) error {
	return nil
}

func (c *handlerTestPaymobClient) GetIframeURL(paymentKey string) string {
	return "/api/v1/mock/paymob/iframe?payment_token=" + paymentKey
}

func (c *handlerTestPaymobClient) VerifyHMAC(payload []byte, receivedHMAC string) bool {
	return c.verifyHMAC
}

type handlerTestPaymentRepo struct {
	nextID         int64
	payments       map[int64]*domain.Payment
	orderToPayment map[int64]int64
}

func newHandlerTestPaymentRepo() *handlerTestPaymentRepo {
	return &handlerTestPaymentRepo{
		nextID:         1,
		payments:       map[int64]*domain.Payment{},
		orderToPayment: map[int64]int64{},
	}
}

func (r *handlerTestPaymentRepo) CreatePayment(ctx context.Context, payment *domain.Payment) (int64, error) {
	id := r.nextID
	r.nextID++
	clone := *payment
	clone.ID = id
	r.payments[id] = &clone
	return id, nil
}

func (r *handlerTestPaymentRepo) GetPaymentByID(ctx context.Context, tenantID, id int64) (*domain.Payment, error) {
	p, ok := r.payments[id]
	if !ok || p.TenantID != tenantID {
		return nil, fmt.Errorf("payment not found")
	}
	clone := *p
	return &clone, nil
}

func (r *handlerTestPaymentRepo) GetPaymentTenantID(ctx context.Context, id int64) (int64, error) {
	p, ok := r.payments[id]
	if !ok {
		return 0, fmt.Errorf("payment not found")
	}
	return p.TenantID, nil
}

func (r *handlerTestPaymentRepo) GetPaymentByMerchantOrderID(ctx context.Context, tenantID int64, merchantOrderID string) (*domain.Payment, error) {
	for _, p := range r.payments {
		if p.TenantID == tenantID && p.MerchantOrderID == merchantOrderID {
			clone := *p
			return &clone, nil
		}
	}
	return nil, fmt.Errorf("payment not found")
}

func (r *handlerTestPaymentRepo) GetPaymentByPaymobOrderID(ctx context.Context, paymobOrderID int64) (*domain.Payment, error) {
	id, ok := r.orderToPayment[paymobOrderID]
	if !ok {
		return nil, fmt.Errorf("payment not found")
	}
	p, ok := r.payments[id]
	if !ok {
		return nil, fmt.Errorf("payment not found")
	}
	clone := *p
	return &clone, nil
}

func (r *handlerTestPaymentRepo) UpdatePaymentOrder(ctx context.Context, tenantID, id int64, paymobOrderID int64, authToken string) error {
	p, ok := r.payments[id]
	if !ok || p.TenantID != tenantID {
		return fmt.Errorf("payment not found")
	}
	p.PaymobOrderID = &paymobOrderID
	r.orderToPayment[paymobOrderID] = id
	return nil
}

func (r *handlerTestPaymentRepo) UpdatePaymentStatus(ctx context.Context, tenantID, id int64, status string, txnID *int64, cardLastFour *string, cardBrand *string, errorMsg *string, paidAt *time.Time) error {
	p, ok := r.payments[id]
	if !ok || p.TenantID != tenantID {
		return fmt.Errorf("payment not found")
	}
	p.Status = status
	p.PaymobTransactionID = txnID
	p.CardLastFour = cardLastFour
	p.CardBrand = cardBrand
	p.ErrorMessage = errorMsg
	p.PaidAt = paidAt
	return nil
}

func (r *handlerTestPaymentRepo) GetPaymentsByTenantID(ctx context.Context, tenantID int64, limit, offset int) ([]domain.Payment, int, error) {
	all := make([]domain.Payment, 0)
	for _, p := range r.payments {
		if p.TenantID == tenantID {
			all = append(all, *p)
		}
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].ID < all[j].ID
	})
	total := len(all)
	if offset >= total {
		return []domain.Payment{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

func (r *handlerTestPaymentRepo) UpdateTenantSubscription(ctx context.Context, tenantID int64, status string, expiresAt time.Time, lastPaymentAt time.Time, lastPaymentAmount int, currency string, planID string) error {
	return nil
}

func (r *handlerTestPaymentRepo) GetTenantBillingInfo(ctx context.Context, tenantID int64) (*domain.TenantBillingInfo, error) {
	return &domain.TenantBillingInfo{
		SubscriptionPlan: "monthly",
	}, nil
}

type handlerTestSubscriptionRepo struct{}

func (r *handlerTestSubscriptionRepo) GetTenantSubscriptionInfo(ctx context.Context, tenantID int64) (*domain.TenantSubscriptionInfo, error) {
	return &domain.TenantSubscriptionInfo{
		TenantID:           tenantID,
		SubscriptionStatus: domain.StatusTrial,
		DaysRemaining:      7,
		CanAccess:          true,
	}, nil
}

func (r *handlerTestSubscriptionRepo) GetTenantSubscriptionBySlug(ctx context.Context, slug string) (*domain.TenantSubscriptionInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *handlerTestSubscriptionRepo) UpdateSubscriptionStatus(ctx context.Context, tenantID int64, status string) error {
	return nil
}

func (r *handlerTestSubscriptionRepo) ExtendTrial(ctx context.Context, params *domain.ExtendTrialParams) error {
	return nil
}

func (r *handlerTestSubscriptionRepo) DisableTenant(ctx context.Context, params *domain.DisableTenantParams) error {
	return nil
}

func (r *handlerTestSubscriptionRepo) EnableTenant(ctx context.Context, params *domain.EnableTenantParams) error {
	return nil
}

func (r *handlerTestSubscriptionRepo) ScheduleDeletion(ctx context.Context, params *domain.DeleteTenantParams) error {
	return nil
}

func (r *handlerTestSubscriptionRepo) CancelDeletion(ctx context.Context, tenantID int64, performedBy int64) error {
	return nil
}

func (r *handlerTestSubscriptionRepo) GetSubscriptionHistory(ctx context.Context, tenantID int64, page int, pageSize int) ([]domain.SubscriptionHistory, int64, error) {
	return []domain.SubscriptionHistory{}, 0, nil
}

func (r *handlerTestSubscriptionRepo) CreateHistoryRecord(ctx context.Context, params *domain.CreateHistoryRecordParams) error {
	return nil
}

func (r *handlerTestSubscriptionRepo) GetExpiredTrials(ctx context.Context) ([]domain.TenantSubscriptionInfo, error) {
	return []domain.TenantSubscriptionInfo{}, nil
}

func (r *handlerTestSubscriptionRepo) GetTenantsScheduledForDeletion(ctx context.Context) ([]domain.TenantSubscriptionInfo, error) {
	return []domain.TenantSubscriptionInfo{}, nil
}

func (r *handlerTestSubscriptionRepo) CheckTenantAccess(ctx context.Context, tenantID int64) (canAccess bool, blockReason string, err error) {
	return true, "", nil
}

func newPaymentHandlerForTest(t *testing.T) (*PaymentHandler, *usecase.PaymentUseCase, *jwt.TokenService, *handlerTestPaymobClient) {
	t.Helper()

	paymobClient := &handlerTestPaymobClient{
		orderID:    98765,
		paymentKey: "handler_key",
		verifyHMAC: true,
	}
	paymentRepo := newHandlerTestPaymentRepo()
	subscriptionRepo := &handlerTestSubscriptionRepo{}
	uc := usecase.NewPaymentUseCase(paymentRepo, subscriptionRepo, nil, paymobClient, nil)

	tokenService, err := jwt.NewTokenService("test-secret", "1h")
	if err != nil {
		t.Fatalf("failed to create token service: %v", err)
	}

	return NewPaymentHandler(uc), uc, tokenService, paymobClient
}

func protectedRoute(tokenService *jwt.TokenService, next http.HandlerFunc) http.Handler {
	return middleware.AuthMiddleware(tokenService)(middleware.TenantContextMiddleware(http.HandlerFunc(next)))
}

func authHeader(t *testing.T, tokenService *jwt.TokenService, tenantID int, restaurantID int) string {
	t.Helper()
	rest := restaurantID
	token, err := tokenService.GenerateToken(7, tenantID, &rest, "user@example.com", "admin")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	return "Bearer " + token
}

func TestPaymentHandlerGetPlansReturns200(t *testing.T) {
	handler, _, _, _ := newPaymentHandlerForTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscription/plans", nil)
	rec := httptest.NewRecorder()

	handler.GetPlans(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"plans"`) {
		t.Fatalf("expected response to contain plans, got: %s", rec.Body.String())
	}
}

func TestPaymentHandlerInitiatePaymentReturns401WithoutJWT(t *testing.T) {
	handler, _, tokenService, _ := newPaymentHandlerForTest(t)

	wrapped := protectedRoute(tokenService, handler.InitiatePayment)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/subscription/initiate", bytes.NewBufferString(`{"plan_id":"monthly"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without JWT, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestPaymentHandlerInitiatePaymentReturns200WithJWT(t *testing.T) {
	handler, _, tokenService, _ := newPaymentHandlerForTest(t)

	wrapped := protectedRoute(tokenService, handler.InitiatePayment)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/subscription/initiate", bytes.NewBufferString(`{"plan_id":"monthly"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader(t, tokenService, 1, 1))
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with valid JWT, got %d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"iframe_url"`) {
		t.Fatalf("expected iframe_url in response, got: %s", rec.Body.String())
	}
}

func TestPaymentHandlerHandlePaymobWebhookReturns200(t *testing.T) {
	handler, uc, _, _ := newPaymentHandlerForTest(t)

	initResp, err := uc.InitiatePayment(context.Background(), 1, 7, "monthly", map[string]string{"email": "user@example.com"})
	if err != nil {
		t.Fatalf("InitiatePayment setup failed: %v", err)
	}

	payload := map[string]interface{}{
		"type": "TRANSACTION",
		"hmac": "mock_hmac",
		"obj": map[string]interface{}{
			"id":                     time.Now().UnixNano(),
			"amount_cents":           75000,
			"success":                false,
			"currency":               "EGP",
			"is_auth":                false,
			"is_capture":             false,
			"is_refunded":            false,
			"is_voided":              false,
			"is_3d_secure":           false,
			"is_standalone_payment":  true,
			"has_parent_transaction": false,
			"error_occured":          true,
			"pending":                false,
			"integration_id":         0,
			"owner":                  1,
			"created_at":             time.Now().Format(time.RFC3339),
			"order": map[string]interface{}{
				"id":                initResp.PaymobOrderID,
				"merchant_order_id": initResp.MerchantOrderID,
			},
			"source_data": map[string]interface{}{
				"pan":      "5123456789011111",
				"sub_type": "mastercard",
				"type":     "card",
			},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhook/paymob?hmac=mock_hmac", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.HandlePaymobWebhook(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestPaymentHandlerGetPaymentHistoryReturns200WithPagination(t *testing.T) {
	handler, _, tokenService, _ := newPaymentHandlerForTest(t)

	initiateWrapped := protectedRoute(tokenService, handler.InitiatePayment)
	initReq := httptest.NewRequest(http.MethodPost, "/api/v1/subscription/initiate", bytes.NewBufferString(`{"plan_id":"monthly"}`))
	initReq.Header.Set("Content-Type", "application/json")
	initReq.Header.Set("Authorization", authHeader(t, tokenService, 1, 1))
	initRec := httptest.NewRecorder()
	initiateWrapped.ServeHTTP(initRec, initReq)
	if initRec.Code != http.StatusOK {
		t.Fatalf("setup initiate expected 200, got %d body=%s", initRec.Code, initRec.Body.String())
	}

	wrapped := protectedRoute(tokenService, handler.GetPaymentHistory)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscription/payments?page=1&page_size=20", nil)
	req.Header.Set("Authorization", authHeader(t, tokenService, 1, 1))
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for payment history, got %d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"payments"`) {
		t.Fatalf("expected payments key in response, got: %s", rec.Body.String())
	}
}
