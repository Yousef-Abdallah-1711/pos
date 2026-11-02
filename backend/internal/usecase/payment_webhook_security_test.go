package usecase

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"pos-saas/internal/domain"
	"pos-saas/internal/service"
)

// SEC-G / SEC-015: webhook authentication, replay and state-transition safety.
// No live payments: deterministic payloads against in-memory repositories.

func pendingPaymentWithOrder(t *testing.T, repo *testPaymentRepo, status string, paymobOrderID int64) int64 {
	t.Helper()
	uid := int64(1)
	id, err := repo.CreatePayment(context.Background(), &domain.Payment{
		TenantID:    1,
		UserID:      &uid,
		AmountCents: 50000,
		Currency:    "EGP",
		Status:      status,
		PlanID:      "monthly",
	})
	if err != nil {
		t.Fatalf("CreatePayment: %v", err)
	}
	if err := repo.UpdatePaymentOrder(context.Background(), id, paymobOrderID, "auth"); err != nil {
		t.Fatalf("UpdatePaymentOrder: %v", err)
	}
	return id
}

func TestWebhookSecurity_InvalidSignatureIsRejectedWithoutStateChange(t *testing.T) {
	repo := newTestPaymentRepo()
	uc := NewPaymentUseCase(repo, newTestSubscriptionRepo(), nil, &testPaymobClient{verifyHMAC: false}, nil)
	id := pendingPaymentWithOrder(t, repo, domain.PaymentStatusPending, 9001)

	payload := webhookPayloadForTest(9001, 50000, true, false, false)
	raw, _ := json.Marshal(payload)
	err := uc.HandleWebhook(context.Background(), raw, payload, "forged")
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "invalid hmac signature") {
		t.Fatalf("want invalid hmac signature, got %v", err)
	}
	p, _ := repo.GetPaymentByID(context.Background(), id)
	if p.Status != domain.PaymentStatusPending || repo.updateStatusCnt != 0 {
		t.Fatalf("state must not change: status=%s updates=%d", p.Status, repo.updateStatusCnt)
	}
}

func TestWebhookSecurity_MissingSignatureIsRejectedByLiveClient(t *testing.T) {
	repo := newTestPaymentRepo()
	live := service.NewRealPaymobClient("api-key", 1, "unit-test-hmac-secret", "")
	uc := NewPaymentUseCase(repo, newTestSubscriptionRepo(), nil, live, nil)
	pendingPaymentWithOrder(t, repo, domain.PaymentStatusPending, 9002)

	payload := webhookPayloadForTest(9002, 50000, true, false, false)
	payload.HMAC = ""
	raw, _ := json.Marshal(payload)
	for name, sig := range map[string]string{"empty": "", "whitespace": "   ", "wrong": "deadbeef"} {
		err := uc.HandleWebhook(context.Background(), raw, payload, sig)
		if err == nil || !strings.Contains(strings.ToLower(err.Error()), "invalid hmac signature") {
			t.Errorf("%s signature must be rejected, got %v", name, err)
		}
	}
	if repo.updateStatusCnt != 0 {
		t.Fatalf("no state change allowed, updates=%d", repo.updateStatusCnt)
	}
}

func TestWebhookSecurity_UnknownPaymobOrderIsRejected(t *testing.T) {
	repo := newTestPaymentRepo()
	uc := NewPaymentUseCase(repo, newTestSubscriptionRepo(), nil, &testPaymobClient{verifyHMAC: true}, nil)
	payload := webhookPayloadForTest(424242, 50000, true, false, false)
	raw, _ := json.Marshal(payload)
	if err := uc.HandleWebhook(context.Background(), raw, payload, "ok"); err == nil {
		t.Fatal("event for an unknown paymob order must be rejected")
	}
	if repo.updateStatusCnt != 0 {
		t.Fatal("no state change allowed")
	}
}

// A late "failed" event must never downgrade a payment that already succeeded.
func TestWebhookSecurity_FailureAfterSuccessDoesNotDowngrade(t *testing.T) {
	repo := newTestPaymentRepo()
	uc := NewPaymentUseCase(repo, newTestSubscriptionRepo(), nil, &testPaymobClient{verifyHMAC: true}, nil)
	id := pendingPaymentWithOrder(t, repo, domain.PaymentStatusSuccess, 9003)

	payload := webhookPayloadForTest(9003, 50000, false, false, true)
	raw, _ := json.Marshal(payload)
	if err := uc.HandleWebhook(context.Background(), raw, payload, "ok"); err != nil {
		t.Fatalf("late failure for a succeeded payment should be a no-op, got %v", err)
	}
	p, _ := repo.GetPaymentByID(context.Background(), id)
	if p.Status != domain.PaymentStatusSuccess || repo.updateStatusCnt != 0 {
		t.Fatalf("status must stay success: status=%s updates=%d", p.Status, repo.updateStatusCnt)
	}
}

// Replaying a success event for a refunded payment must not re-activate the subscription.
func TestWebhookSecurity_ReplayedSuccessForRefundedPaymentIsIgnored(t *testing.T) {
	repo := newTestPaymentRepo()
	subRepo := newTestSubscriptionRepo()
	uc := NewPaymentUseCase(repo, subRepo, nil, &testPaymobClient{verifyHMAC: true}, nil)
	id := pendingPaymentWithOrder(t, repo, domain.PaymentStatusRefunded, 9004)

	payload := webhookPayloadForTest(9004, 50000, true, false, false)
	raw, _ := json.Marshal(payload)
	for i := 0; i < 3; i++ {
		if err := uc.HandleWebhook(context.Background(), raw, payload, "ok"); err != nil {
			t.Fatalf("replay %d: want silent ignore, got %v", i, err)
		}
	}
	p, _ := repo.GetPaymentByID(context.Background(), id)
	if p.Status != domain.PaymentStatusRefunded || repo.updateStatusCnt != 0 {
		t.Fatalf("refunded payment must stay refunded: status=%s updates=%d", p.Status, repo.updateStatusCnt)
	}
}

// Refund / void callbacks are acknowledged without activating anything.
func TestWebhookSecurity_RefundAndVoidCallbacksNeverActivate(t *testing.T) {
	for name, mutate := range map[string]func(*domain.PaymobWebhookPayload){
		"refunded": func(p *domain.PaymobWebhookPayload) { p.Obj.IsRefunded = true },
		"voided":   func(p *domain.PaymobWebhookPayload) { p.Obj.IsVoided = true },
	} {
		repo := newTestPaymentRepo()
		uc := NewPaymentUseCase(repo, newTestSubscriptionRepo(), nil, &testPaymobClient{verifyHMAC: true}, nil)
		id := pendingPaymentWithOrder(t, repo, domain.PaymentStatusPending, 9005)

		payload := webhookPayloadForTest(9005, 50000, true, false, false)
		mutate(payload)
		raw, _ := json.Marshal(payload)
		if err := uc.HandleWebhook(context.Background(), raw, payload, "ok"); err != nil {
			t.Fatalf("%s: want acknowledged, got %v", name, err)
		}
		p, _ := repo.GetPaymentByID(context.Background(), id)
		if p.Status != domain.PaymentStatusPending || repo.updateStatusCnt != 0 {
			t.Fatalf("%s must not change state: status=%s updates=%d", name, p.Status, repo.updateStatusCnt)
		}
	}
}

// Duplicate failed events are harmless and idempotent.
func TestWebhookSecurity_DuplicateFailedEventsStayFailed(t *testing.T) {
	repo := newTestPaymentRepo()
	uc := NewPaymentUseCase(repo, newTestSubscriptionRepo(), nil, &testPaymobClient{verifyHMAC: true}, nil)
	id := pendingPaymentWithOrder(t, repo, domain.PaymentStatusPending, 9006)

	payload := webhookPayloadForTest(9006, 50000, false, false, true)
	raw, _ := json.Marshal(payload)
	for i := 0; i < 3; i++ {
		if err := uc.HandleWebhook(context.Background(), raw, payload, "ok"); err != nil {
			t.Fatalf("duplicate failed event %d: %v", i, err)
		}
	}
	p, _ := repo.GetPaymentByID(context.Background(), id)
	if p.Status != domain.PaymentStatusFailed {
		t.Fatalf("want failed, got %s", p.Status)
	}
}

// Conflicting amount/currency on a legitimate order is rejected (state-transition integrity).
func TestWebhookSecurity_CurrencyMismatchIsRejected(t *testing.T) {
	repo := newTestPaymentRepo()
	uc := NewPaymentUseCase(repo, newTestSubscriptionRepo(), nil, &testPaymobClient{verifyHMAC: true}, nil)
	pendingPaymentWithOrder(t, repo, domain.PaymentStatusPending, 9007)

	payload := webhookPayloadForTest(9007, 50000, true, false, false)
	payload.Obj.Currency = "USD"
	raw, _ := json.Marshal(payload)
	err := uc.HandleWebhook(context.Background(), raw, payload, "ok")
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "currency mismatch") {
		t.Fatalf("want currency mismatch, got %v", err)
	}
	if repo.updateStatusCnt != 0 {
		t.Fatal("no state change allowed")
	}
}
