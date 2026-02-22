package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"pos-saas/internal/domain"
)

type testPaymobClient struct {
	authToken   string
	orderID     int64
	paymentKey  string
	verifyHMAC  bool
	authErr     error
	createErr   error
	paymentErr  error
	iframeBase  string
	orderCalls  int
	keyCalls    int
	verifyCalls int
}

func (c *testPaymobClient) Authenticate() (string, error) {
	if c.authErr != nil {
		return "", c.authErr
	}
	if strings.TrimSpace(c.authToken) == "" {
		return "mock_auth", nil
	}
	return c.authToken, nil
}

func (c *testPaymobClient) CreateOrder(authToken string, amountCents int, currency string, merchantOrderID string) (int64, error) {
	c.orderCalls++
	if c.createErr != nil {
		return 0, c.createErr
	}
	if c.orderID <= 0 {
		c.orderID = 1001
	}
	return c.orderID, nil
}

func (c *testPaymobClient) GetPaymentKey(authToken string, orderID int64, amountCents int, currency string, integrationID int, billingData map[string]string, expiration int) (string, error) {
	c.keyCalls++
	if c.paymentErr != nil {
		return "", c.paymentErr
	}
	if strings.TrimSpace(c.paymentKey) == "" {
		return "mock_payment_key", nil
	}
	return c.paymentKey, nil
}

func (c *testPaymobClient) RefundTransaction(authToken string, transactionID int64, amountCents int) error {
	return nil
}

func (c *testPaymobClient) GetIframeURL(paymentKey string) string {
	base := strings.TrimSpace(c.iframeBase)
	if base == "" {
		base = "/api/v1/mock/paymob/iframe?payment_token="
	}
	return base + paymentKey
}

func (c *testPaymobClient) VerifyHMAC(payload []byte, receivedHMAC string) bool {
	c.verifyCalls++
	return c.verifyHMAC
}

type testPaymentRepo struct {
	nextID          int64
	payments        map[int64]*domain.Payment
	orderToPayment  map[int64]int64
	lastLimit       int
	lastOffset      int
	updateStatusCnt int
	billingByTenant map[int64]*domain.TenantBillingInfo
}

func newTestPaymentRepo() *testPaymentRepo {
	return &testPaymentRepo{
		nextID:          1,
		payments:        map[int64]*domain.Payment{},
		orderToPayment:  map[int64]int64{},
		billingByTenant: map[int64]*domain.TenantBillingInfo{},
	}
}

func (r *testPaymentRepo) CreatePayment(ctx context.Context, payment *domain.Payment) (int64, error) {
	if payment == nil {
		return 0, errors.New("payment is nil")
	}

	id := r.nextID
	r.nextID++

	clone := *payment
	clone.ID = id
	r.payments[id] = &clone
	return id, nil
}

func (r *testPaymentRepo) GetPaymentByID(ctx context.Context, id int64) (*domain.Payment, error) {
	p, ok := r.payments[id]
	if !ok {
		return nil, fmt.Errorf("payment not found")
	}
	clone := *p
	return &clone, nil
}

func (r *testPaymentRepo) GetPaymentByMerchantOrderID(ctx context.Context, merchantOrderID string) (*domain.Payment, error) {
	for _, p := range r.payments {
		if p.MerchantOrderID == merchantOrderID {
			clone := *p
			return &clone, nil
		}
	}
	return nil, fmt.Errorf("payment not found")
}

func (r *testPaymentRepo) GetPaymentByPaymobOrderID(ctx context.Context, paymobOrderID int64) (*domain.Payment, error) {
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

func (r *testPaymentRepo) UpdatePaymentOrder(ctx context.Context, id int64, paymobOrderID int64, authToken string) error {
	p, ok := r.payments[id]
	if !ok {
		return fmt.Errorf("payment not found")
	}
	p.PaymobOrderID = &paymobOrderID
	r.orderToPayment[paymobOrderID] = id
	return nil
}

func (r *testPaymentRepo) UpdatePaymentStatus(ctx context.Context, id int64, status string, txnID *int64, cardLastFour *string, cardBrand *string, errorMsg *string, paidAt *time.Time) error {
	p, ok := r.payments[id]
	if !ok {
		return fmt.Errorf("payment not found")
	}
	r.updateStatusCnt++
	p.Status = status
	p.PaymobTransactionID = txnID
	p.CardLastFour = cardLastFour
	p.CardBrand = cardBrand
	p.ErrorMessage = errorMsg
	p.PaidAt = paidAt
	return nil
}

func (r *testPaymentRepo) GetPaymentsByTenantID(ctx context.Context, tenantID int64, limit, offset int) ([]domain.Payment, int, error) {
	r.lastLimit = limit
	r.lastOffset = offset

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

func (r *testPaymentRepo) GetPaymentsByTenantCursor(ctx context.Context, tenantID int64, cursorID int64, limit int) ([]domain.Payment, *int64, error) {
	if limit <= 0 {
		limit = 20
	}

	all := make([]domain.Payment, 0)
	for _, p := range r.payments {
		if p.TenantID != tenantID {
			continue
		}
		if cursorID > 0 && p.ID >= cursorID {
			continue
		}
		all = append(all, *p)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].ID > all[j].ID
	})

	var nextCursor *int64
	if len(all) > limit {
		lastID := all[limit-1].ID
		nextCursor = &lastID
		all = all[:limit]
	}

	return all, nextCursor, nil
}

func (r *testPaymentRepo) UpdateTenantSubscription(ctx context.Context, tenantID int64, status string, expiresAt time.Time, lastPaymentAt time.Time, lastPaymentAmount int, currency string, planID string) error {
	return nil
}

func (r *testPaymentRepo) GetTenantBillingInfo(ctx context.Context, tenantID int64) (*domain.TenantBillingInfo, error) {
	info, ok := r.billingByTenant[tenantID]
	if !ok {
		return &domain.TenantBillingInfo{}, nil
	}
	clone := *info
	return &clone, nil
}

type testSubscriptionRepo struct {
	infoByTenant map[int64]*domain.TenantSubscriptionInfo
	history      []*domain.CreateHistoryRecordParams
}

func newTestSubscriptionRepo() *testSubscriptionRepo {
	return &testSubscriptionRepo{
		infoByTenant: map[int64]*domain.TenantSubscriptionInfo{},
		history:      make([]*domain.CreateHistoryRecordParams, 0),
	}
}

func (r *testSubscriptionRepo) GetTenantSubscriptionInfo(ctx context.Context, tenantID int64) (*domain.TenantSubscriptionInfo, error) {
	info, ok := r.infoByTenant[tenantID]
	if !ok {
		return nil, fmt.Errorf("tenant not found")
	}
	clone := *info
	return &clone, nil
}

func (r *testSubscriptionRepo) GetTenantSubscriptionBySlug(ctx context.Context, slug string) (*domain.TenantSubscriptionInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *testSubscriptionRepo) UpdateSubscriptionStatus(ctx context.Context, tenantID int64, status string) error {
	return nil
}

func (r *testSubscriptionRepo) ExtendTrial(ctx context.Context, params *domain.ExtendTrialParams) error {
	return nil
}

func (r *testSubscriptionRepo) DisableTenant(ctx context.Context, params *domain.DisableTenantParams) error {
	return nil
}

func (r *testSubscriptionRepo) EnableTenant(ctx context.Context, params *domain.EnableTenantParams) error {
	return nil
}

func (r *testSubscriptionRepo) ScheduleDeletion(ctx context.Context, params *domain.DeleteTenantParams) error {
	return nil
}

func (r *testSubscriptionRepo) CancelDeletion(ctx context.Context, tenantID int64, performedBy int64) error {
	return nil
}

func (r *testSubscriptionRepo) GetSubscriptionHistory(ctx context.Context, tenantID int64, page int, pageSize int) ([]domain.SubscriptionHistory, int64, error) {
	return []domain.SubscriptionHistory{}, 0, nil
}

func (r *testSubscriptionRepo) CreateHistoryRecord(ctx context.Context, params *domain.CreateHistoryRecordParams) error {
	if params != nil {
		r.history = append(r.history, params)
	}
	return nil
}

func (r *testSubscriptionRepo) GetExpiredTrials(ctx context.Context) ([]domain.TenantSubscriptionInfo, error) {
	return []domain.TenantSubscriptionInfo{}, nil
}

func (r *testSubscriptionRepo) GetTenantsScheduledForDeletion(ctx context.Context) ([]domain.TenantSubscriptionInfo, error) {
	return []domain.TenantSubscriptionInfo{}, nil
}

func (r *testSubscriptionRepo) CheckTenantAccess(ctx context.Context, tenantID int64) (canAccess bool, blockReason string, err error) {
	return true, "", nil
}

func TestPaymentUseCaseInitiatePaymentCreatesPendingRecord(t *testing.T) {
	repo := newTestPaymentRepo()
	subRepo := newTestSubscriptionRepo()
	client := &testPaymobClient{
		authToken:  "auth_123",
		orderID:    909090,
		paymentKey: "pk_abc",
		verifyHMAC: true,
	}

	uc := NewPaymentUseCase(repo, subRepo, nil, client, nil)

	resp, err := uc.InitiatePayment(
		context.Background(),
		1,
		7,
		"monthly",
		map[string]string{
			"first_name":   "Test",
			"last_name":    "User",
			"email":        "test@example.com",
			"phone_number": "+201000000000",
		},
	)
	if err != nil {
		t.Fatalf("InitiatePayment returned error: %v", err)
	}
	if resp.PaymentID <= 0 {
		t.Fatalf("expected payment id > 0, got %d", resp.PaymentID)
	}
	if resp.PaymobOrderID != 909090 {
		t.Fatalf("expected paymob order id 909090, got %d", resp.PaymobOrderID)
	}
	if !strings.Contains(resp.IframeURL, "pk_abc") {
		t.Fatalf("expected iframe URL to contain payment key, got %s", resp.IframeURL)
	}

	created, err := repo.GetPaymentByID(context.Background(), resp.PaymentID)
	if err != nil {
		t.Fatalf("failed to load created payment: %v", err)
	}
	if created.Status != domain.PaymentStatusPending {
		t.Fatalf("expected pending status, got %s", created.Status)
	}
	if created.PaymobOrderID == nil || *created.PaymobOrderID != 909090 {
		t.Fatalf("expected paymob order id to be stored, got %+v", created.PaymobOrderID)
	}
}

func TestPaymentUseCaseInitiatePaymentRejectsInvalidPlan(t *testing.T) {
	repo := newTestPaymentRepo()
	subRepo := newTestSubscriptionRepo()
	client := &testPaymobClient{verifyHMAC: true}
	uc := NewPaymentUseCase(repo, subRepo, nil, client, nil)

	_, err := uc.InitiatePayment(context.Background(), 1, 1, "invalid-plan", map[string]string{})
	if err == nil {
		t.Fatal("expected error for invalid plan")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "invalid plan") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPaymentUseCaseHandleWebhookDuplicateSuccessIsIdempotent(t *testing.T) {
	repo := newTestPaymentRepo()
	subRepo := newTestSubscriptionRepo()
	client := &testPaymobClient{verifyHMAC: true}
	uc := NewPaymentUseCase(repo, subRepo, nil, client, nil)

	paymentID, err := repo.CreatePayment(context.Background(), &domain.Payment{
		TenantID:    1,
		AmountCents: 50000,
		Currency:    "EGP",
		Status:      domain.PaymentStatusSuccess,
		PlanID:      "monthly",
	})
	if err != nil {
		t.Fatalf("CreatePayment failed: %v", err)
	}
	if err := repo.UpdatePaymentOrder(context.Background(), paymentID, 12345, "auth"); err != nil {
		t.Fatalf("UpdatePaymentOrder failed: %v", err)
	}

	payload := webhookPayloadForTest(12345, 50000, true, false, false)
	raw, _ := json.Marshal(payload)
	if err := uc.HandleWebhook(context.Background(), raw, payload, "ok"); err != nil {
		t.Fatalf("expected idempotent success, got error: %v", err)
	}
}

func TestPaymentUseCaseHandleWebhookRejectsAmountMismatch(t *testing.T) {
	repo := newTestPaymentRepo()
	subRepo := newTestSubscriptionRepo()
	client := &testPaymobClient{verifyHMAC: true, orderID: 777}
	uc := NewPaymentUseCase(repo, subRepo, nil, client, nil)

	initResp, err := uc.InitiatePayment(context.Background(), 1, 1, "monthly", map[string]string{})
	if err != nil {
		t.Fatalf("InitiatePayment failed: %v", err)
	}

	payload := webhookPayloadForTest(initResp.PaymobOrderID, 100, true, false, false)
	raw, _ := json.Marshal(payload)

	err = uc.HandleWebhook(context.Background(), raw, payload, "ok")
	if err == nil {
		t.Fatal("expected amount mismatch error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "amount mismatch") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPaymentUseCaseHandleWebhookFailedPaymentMarksFailed(t *testing.T) {
	repo := newTestPaymentRepo()
	subRepo := newTestSubscriptionRepo()
	client := &testPaymobClient{verifyHMAC: true, orderID: 555}
	uc := NewPaymentUseCase(repo, subRepo, nil, client, nil)

	initResp, err := uc.InitiatePayment(context.Background(), 1, 1, "monthly", map[string]string{})
	if err != nil {
		t.Fatalf("InitiatePayment failed: %v", err)
	}

	payload := webhookPayloadForTest(initResp.PaymobOrderID, 50000, false, false, true)
	raw, _ := json.Marshal(payload)
	if err := uc.HandleWebhook(context.Background(), raw, payload, "ok"); err != nil {
		t.Fatalf("HandleWebhook failed payload returned error: %v", err)
	}

	p, err := repo.GetPaymentByID(context.Background(), initResp.PaymentID)
	if err != nil {
		t.Fatalf("GetPaymentByID failed: %v", err)
	}
	if p.Status != domain.PaymentStatusFailed {
		t.Fatalf("expected failed status, got %s", p.Status)
	}
	if repo.updateStatusCnt == 0 {
		t.Fatal("expected UpdatePaymentStatus to be called")
	}
}

func TestPaymentUseCaseGetPlansIncludesMonthlyPlan(t *testing.T) {
	uc := NewPaymentUseCase(newTestPaymentRepo(), newTestSubscriptionRepo(), nil, &testPaymobClient{verifyHMAC: true}, nil)

	plans := uc.GetPlans()
	if len(plans) == 0 {
		t.Fatal("expected at least one plan")
	}

	foundMonthly := false
	for _, plan := range plans {
		if plan.ID == "monthly" {
			foundMonthly = true
			break
		}
	}
	if !foundMonthly {
		t.Fatalf("expected monthly plan in list, got %+v", plans)
	}
}

func TestPaymentUseCaseGetPaymentHistoryAppliesPaginationBounds(t *testing.T) {
	repo := newTestPaymentRepo()
	subRepo := newTestSubscriptionRepo()
	client := &testPaymobClient{verifyHMAC: true}
	uc := NewPaymentUseCase(repo, subRepo, nil, client, nil)

	for i := 0; i < 3; i++ {
		_, err := repo.CreatePayment(context.Background(), &domain.Payment{
			TenantID:    1,
			AmountCents: 50000,
			Currency:    "EGP",
			Status:      domain.PaymentStatusPending,
			PlanID:      "monthly",
		})
		if err != nil {
			t.Fatalf("CreatePayment failed: %v", err)
		}
	}

	payments, total, err := uc.GetPaymentHistory(context.Background(), 1, 0, 500)
	if err != nil {
		t.Fatalf("GetPaymentHistory returned error: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total=3, got %d", total)
	}
	if len(payments) != 3 {
		t.Fatalf("expected 3 payments, got %d", len(payments))
	}
	if repo.lastLimit != 20 || repo.lastOffset != 0 {
		t.Fatalf("expected normalized pagination limit=20 offset=0, got limit=%d offset=%d", repo.lastLimit, repo.lastOffset)
	}
}

func TestPaymentUseCaseGetCurrentSubscriptionReturnsPlanAndBilling(t *testing.T) {
	repo := newTestPaymentRepo()
	subRepo := newTestSubscriptionRepo()
	client := &testPaymobClient{verifyHMAC: true}
	uc := NewPaymentUseCase(repo, subRepo, nil, client, nil)

	exp := time.Now().Add(7 * 24 * time.Hour)
	lastPayment := time.Now().Add(-24 * time.Hour)
	subRepo.infoByTenant[1] = &domain.TenantSubscriptionInfo{
		TenantID:           1,
		SubscriptionStatus: domain.StatusActive,
		DaysRemaining:      7,
		CanAccess:          true,
	}
	repo.billingByTenant[1] = &domain.TenantBillingInfo{
		SubscriptionPlan:      "monthly",
		SubscriptionExpiresAt: &exp,
		LastPaymentAt:         &lastPayment,
	}

	resp, err := uc.GetCurrentSubscription(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCurrentSubscription returned error: %v", err)
	}
	if resp.SubscriptionStatus != domain.StatusActive {
		t.Fatalf("expected active subscription status, got %s", resp.SubscriptionStatus)
	}
	if resp.Plan.ID != "monthly" {
		t.Fatalf("expected monthly plan id, got %s", resp.Plan.ID)
	}
	if resp.SubscriptionExpiresAt == nil {
		t.Fatal("expected subscription_expires_at to be set")
	}
}

func webhookPayloadForTest(orderID int64, amountCents int, success bool, pending bool, errorOccured bool) *domain.PaymobWebhookPayload {
	payload := &domain.PaymobWebhookPayload{
		Type: "TRANSACTION",
		HMAC: "mock_hmac",
	}
	payload.Obj.ID = time.Now().UnixNano()
	payload.Obj.AmountCents = amountCents
	payload.Obj.Success = success
	payload.Obj.Currency = "EGP"
	payload.Obj.IsAuth = success
	payload.Obj.IsCapture = success
	payload.Obj.IsRefunded = false
	payload.Obj.IsVoided = false
	payload.Obj.Is3DSecure = false
	payload.Obj.IsStandalonePayment = true
	payload.Obj.HasParentTransaction = false
	payload.Obj.ErrorOccured = errorOccured
	payload.Obj.Pending = pending
	payload.Obj.IntegrationID = 0
	payload.Obj.Owner = 1
	payload.Obj.CreatedAt = time.Now().Format(time.RFC3339)
	payload.Obj.Order.ID = orderID
	payload.Obj.Order.MerchantOrderID = fmt.Sprintf("tenant_1_monthly_%d", orderID)
	payload.Obj.SourceData.Pan = "5123456789011111"
	payload.Obj.SourceData.SubType = "mastercard"
	payload.Obj.SourceData.Type = "card"

	return payload
}
