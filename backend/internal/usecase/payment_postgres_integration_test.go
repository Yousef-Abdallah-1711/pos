package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestPaymentPostgresDuplicateSuccessCallbacksAreIdempotent(t *testing.T) {
	db, uc, tenantID, paymentID, orderID, subRepo := openPaymentPostgresFixture(t)
	payload := webhookPayloadForTest(orderID, 50000, true, false, false)
	raw, _ := json.Marshal(payload)

	const callbacks = 12
	var group sync.WaitGroup
	errs := make(chan error, callbacks)
	for range callbacks {
		group.Add(1)
		go func() {
			defer group.Done()
			errs <- uc.HandleWebhook(context.Background(), raw, payload, "valid-test-signature")
		}()
	}
	group.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Errorf("duplicate callback returned error: %v", err)
		}
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusSuccess, domain.StatusActive)
	if len(subRepo.history) != 2 {
		t.Fatalf("expected two success/activation history entries, got %d", len(subRepo.history))
	}
}

func TestPaymentPostgresOutOfOrderFailureThenSuccessStaysSuccessful(t *testing.T) {
	db, uc, tenantID, paymentID, orderID, _ := openPaymentPostgresFixture(t)
	failure := webhookPayloadForTest(orderID, 50000, false, false, true)
	if err := handlePostgresWebhook(t, uc, failure); err != nil {
		t.Fatalf("initial failure callback: %v", err)
	}
	if err := handlePostgresWebhook(t, uc, webhookPayloadForTest(orderID, 50000, true, false, false)); err != nil {
		t.Fatalf("later success callback: %v", err)
	}
	if err := handlePostgresWebhook(t, uc, failure); err != nil {
		t.Fatalf("delayed failure callback: %v", err)
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusSuccess, domain.StatusActive)
}

func TestPaymentPostgresPendingCallbackLeavesPaymentPendingUntilSuccess(t *testing.T) {
	db, uc, tenantID, paymentID, orderID, subRepo := openPaymentPostgresFixture(t)
	if err := handlePostgresWebhook(t, uc, webhookPayloadForTest(orderID, 50000, false, true, false)); err != nil {
		t.Fatalf("pending callback: %v", err)
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusPending, domain.StatusTrial)
	if len(subRepo.history) != 0 {
		t.Fatalf("pending callback wrote subscription history: %+v", subRepo.history)
	}

	if err := handlePostgresWebhook(t, uc, webhookPayloadForTest(orderID, 50000, true, false, false)); err != nil {
		t.Fatalf("success callback after pending event: %v", err)
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusSuccess, domain.StatusActive)
	if len(subRepo.history) != 2 {
		t.Fatalf("expected payment success and activation history, got %d entries", len(subRepo.history))
	}
}

func TestPaymentPostgresInvalidSignatureDoesNotMutateState(t *testing.T) {
	db, uc, tenantID, paymentID, orderID, _ := openPaymentPostgresFixture(t)
	uc.paymobClient = &postgresRacePaymobClient{validHMAC: false}
	err := handlePostgresWebhook(t, uc, webhookPayloadForTest(orderID, 50000, true, false, false))
	if err == nil {
		t.Fatal("invalid callback signature unexpectedly succeeded")
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusPending, domain.StatusTrial)
}

func TestPaymentPostgresSuccessRollsBackWhenSubscriptionWriteFails(t *testing.T) {
	db, uc, tenantID, paymentID, orderID, _ := openPaymentPostgresFixture(t)
	constraint := fmt.Sprintf("closure_reject_active_%d", tenantID)
	if _, err := db.Exec(`ALTER TABLE tenants ADD CONSTRAINT ` + constraint + ` CHECK (id <> ` + fmt.Sprint(tenantID) + ` OR subscription_status <> 'active')`); err != nil {
		t.Fatalf("install rollback constraint: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Exec(`ALTER TABLE tenants DROP CONSTRAINT IF EXISTS ` + constraint) })

	err := handlePostgresWebhook(t, uc, webhookPayloadForTest(orderID, 50000, true, false, false))
	if err == nil {
		t.Fatal("success callback should fail when subscription update is rejected")
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusPending, domain.StatusTrial)
}

func TestPaymentPostgresConcurrentRefundCallsProviderOnce(t *testing.T) {
	db, uc, tenantID, paymentID, _, subRepo := openPaymentPostgresFixture(t)
	transactionID := int64(880011)
	if _, err := db.Exec(`UPDATE payments SET status = 'success', paymob_transaction_id = $1 WHERE id = $2`, transactionID, paymentID); err != nil {
		t.Fatalf("prepare successful payment: %v", err)
	}
	subRepo.infoByTenant[tenantID] = &domain.TenantSubscriptionInfo{SubscriptionStatus: domain.StatusActive}
	provider := &postgresRacePaymobClient{authBarrier: make(chan struct{}), validHMAC: true}
	uc.paymobClient = provider

	var group sync.WaitGroup
	errs := make(chan error, 2)
	for range 2 {
		group.Add(1)
		go func() {
			defer group.Done()
			_, err := uc.RefundPayment(context.Background(), paymentID, 77, nil, "", "")
			errs <- err
		}()
	}
	group.Wait()
	close(errs)
	successes := 0
	for err := range errs {
		if err == nil {
			successes++
		} else if err.Error() != "payment is no longer refundable" {
			t.Errorf("unexpected refund result: %v", err)
		}
	}
	if successes != 1 || provider.refundCalls.Load() != 1 {
		t.Fatalf("want one refund commit and one provider call; commits=%d provider_calls=%d", successes, provider.refundCalls.Load())
	}
	if err := handlePostgresWebhook(t, uc, webhookPayloadForTest(880010, 50000, true, false, false)); err != nil {
		t.Fatalf("delayed success after refund: %v", err)
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusRefunded, domain.StatusExpired)
}

func openPaymentPostgresFixture(t *testing.T) (*sql.DB, *PaymentUseCase, int64, int64, int64, *testSubscriptionRepo) {
	t.Helper()
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}

	var tenantID int64
	err = db.QueryRowContext(ctx, `INSERT INTO tenants (name, slug, domain) VALUES ('closure-payment-tenant', 'closure-payment-tenant', 'closure-payment.example.test') RETURNING id`).Scan(&tenantID)
	if err != nil {
		db.Close()
		t.Fatalf("insert synthetic tenant: %v", err)
	}
	paymentRepo := repository.NewPaymentRepository(db)
	orderID := int64(880010)
	paymentID, err := paymentRepo.CreatePayment(ctx, &domain.Payment{
		TenantID: tenantID, AmountCents: 50000, Currency: "EGP", Status: domain.PaymentStatusPending,
		PlanID: "monthly", MerchantOrderID: "closure-payment-order",
	})
	if err == nil {
		err = paymentRepo.UpdatePaymentOrder(ctx, paymentID, orderID, "synthetic-auth-token")
	}
	if err != nil {
		_, _ = db.Exec(`DELETE FROM tenants WHERE id = $1`, tenantID)
		db.Close()
		t.Fatalf("create synthetic pending payment: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.Exec(`DELETE FROM tenants WHERE id = $1`, tenantID)
		_ = db.Close()
	})

	subRepo := newTestSubscriptionRepo()
	uc := NewPaymentUseCase(paymentRepo, subRepo, nil, &postgresRacePaymobClient{validHMAC: true}, db)
	return db, uc, tenantID, paymentID, orderID, subRepo
}

func handlePostgresWebhook(t *testing.T, uc *PaymentUseCase, payload *domain.PaymobWebhookPayload) error {
	t.Helper()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal callback: %v", err)
	}
	return uc.HandleWebhook(context.Background(), raw, payload, "valid-test-signature")
}

func assertPaymentAndTenantState(t *testing.T, db *sql.DB, paymentID, tenantID int64, paymentStatus, tenantStatus string) {
	t.Helper()
	var actualPaymentStatus, actualTenantStatus string
	err := db.QueryRow(`SELECT p.status, t.subscription_status FROM payments p JOIN tenants t ON t.id = p.tenant_id WHERE p.id = $1 AND t.id = $2`, paymentID, tenantID).Scan(&actualPaymentStatus, &actualTenantStatus)
	if err != nil {
		t.Fatalf("read payment/tenant state: %v", err)
	}
	if actualPaymentStatus != paymentStatus || actualTenantStatus != tenantStatus {
		t.Fatalf("unexpected payment/tenant state: payment=%q tenant=%q, want %q/%q", actualPaymentStatus, actualTenantStatus, paymentStatus, tenantStatus)
	}
}

type postgresRacePaymobClient struct {
	authCalls   atomic.Int32
	refundCalls atomic.Int32
	authBarrier chan struct{}
	barrierOnce sync.Once
	validHMAC   bool
}

func (c *postgresRacePaymobClient) Authenticate() (string, error) {
	if c.authBarrier != nil {
		if c.authCalls.Add(1) >= 2 {
			c.barrierOnce.Do(func() { close(c.authBarrier) })
		}
		<-c.authBarrier
	}
	return "synthetic-auth", nil
}

func (*postgresRacePaymobClient) CreateOrder(string, int, string, string) (int64, error) {
	return 1, nil
}
func (*postgresRacePaymobClient) GetPaymentKey(string, int64, int, string, int, map[string]string, int) (string, error) {
	return "synthetic-key", nil
}
func (c *postgresRacePaymobClient) RefundTransaction(string, int64, int) error {
	c.refundCalls.Add(1)
	return nil
}
func (*postgresRacePaymobClient) GetIframeURL(string) string       { return "about:blank" }
func (c *postgresRacePaymobClient) VerifyHMAC([]byte, string) bool { return c.validHMAC }
