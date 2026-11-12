package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

func TestPaymentPostgresVoidIsDurableAndRejectsDelayedSuccess(t *testing.T) {
	db, uc, tenantID, paymentID, orderID, subRepo := openPaymentPostgresFixture(t)
	invalidVoid := webhookPayloadForTest(orderID, 1, false, false, false)
	invalidVoid.Obj.IsVoided = true
	if err := handlePostgresWebhook(t, uc, invalidVoid); err == nil {
		t.Fatal("void callback with mismatched amount should be rejected")
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusPending, domain.StatusTrial)

	void := webhookPayloadForTest(orderID, 50000, false, false, false)
	void.Obj.IsVoided = true
	if err := handlePostgresWebhook(t, uc, void); err != nil {
		t.Fatalf("provider void callback: %v", err)
	}
	if err := handlePostgresWebhook(t, uc, void); err != nil {
		t.Fatalf("duplicate provider void callback: %v", err)
	}
	if err := handlePostgresWebhook(t, uc, webhookPayloadForTest(orderID, 50000, true, false, false)); err != nil {
		t.Fatalf("delayed success callback after void: %v", err)
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusVoided, domain.StatusTrial)
	if len(subRepo.history) != 1 || subRepo.history[0].Action != domain.ActionPaymentVoided {
		t.Fatalf("void replay or delayed success wrote unexpected history: %+v", subRepo.history)
	}
}

func TestPaymentPostgresVoidCallbackUsesRestrictedTenantTransaction(t *testing.T) {
	runtimeURL := os.Getenv("POS_APP_DATABASE_URL")
	if runtimeURL == "" {
		t.Skip("set POS_APP_DATABASE_URL to the restricted pos_app role")
	}
	db, uc, tenantID, paymentID, orderID, _ := openPaymentPostgresFixture(t)
	runtimeDB, err := sql.Open("postgres", runtimeURL)
	if err != nil {
		t.Fatalf("open restricted runtime database: %v", err)
	}
	t.Cleanup(func() { _ = runtimeDB.Close() })
	if err := runtimeDB.PingContext(context.Background()); err != nil {
		t.Fatalf("connect restricted runtime database: %v", err)
	}
	uc.db = runtimeDB
	uc.paymentRepo = repository.NewPaymentRepository(runtimeDB)

	void := webhookPayloadForTest(orderID, 50000, false, false, false)
	void.Obj.IsVoided = true
	if err := handlePostgresWebhook(t, uc, void); err != nil {
		t.Fatalf("restricted-role provider void callback: %v", err)
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusVoided, domain.StatusTrial)
}

func TestPaymentPostgresConcurrentVoidAndSuccessHaveSingleWinner(t *testing.T) {
	db, uc, tenantID, paymentID, orderID, subRepo := openPaymentPostgresFixture(t)
	void := webhookPayloadForTest(orderID, 50000, false, false, false)
	void.Obj.IsVoided = true
	success := webhookPayloadForTest(orderID, 50000, true, false, false)
	start := make(chan struct{})
	errs := make(chan error, 2)
	var group sync.WaitGroup
	for _, callback := range []*domain.PaymobWebhookPayload{void, success} {
		group.Add(1)
		go func(payload *domain.PaymobWebhookPayload) {
			defer group.Done()
			<-start
			errs <- handlePostgresWebhook(t, uc, payload)
		}(callback)
	}
	close(start)
	group.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Errorf("concurrent callback returned error: %v", err)
		}
	}

	var paymentStatus, tenantStatus string
	if err := db.QueryRow(`SELECT p.status, t.subscription_status FROM payments p JOIN tenants t ON t.id = p.tenant_id WHERE p.id = $1 AND p.tenant_id = $2`, paymentID, tenantID).Scan(&paymentStatus, &tenantStatus); err != nil {
		t.Fatalf("read concurrent callback result: %v", err)
	}
	switch {
	case paymentStatus == domain.PaymentStatusVoided && tenantStatus == domain.StatusTrial:
		if len(subRepo.history) != 1 || subRepo.history[0].Action != domain.ActionPaymentVoided {
			t.Fatalf("void winner recorded inconsistent history: %+v", subRepo.history)
		}
	case paymentStatus == domain.PaymentStatusSuccess && tenantStatus == domain.StatusActive:
		if len(subRepo.history) != 2 {
			t.Fatalf("success winner should record activation exactly once: %+v", subRepo.history)
		}
	default:
		t.Fatalf("void and success callbacks produced invalid mixed state payment=%q tenant=%q", paymentStatus, tenantStatus)
	}
}

func TestPaymentRepositoryUsesForcedTenantRLSWithRuntimeRole(t *testing.T) {
	runtimeURL := os.Getenv("POS_APP_DATABASE_URL")
	if runtimeURL == "" {
		t.Skip("set POS_APP_DATABASE_URL to the restricted pos_app role")
	}
	db, _, tenantID, paymentID, _, _ := openPaymentPostgresFixture(t)
	runtimeDB, err := sql.Open("postgres", runtimeURL)
	if err != nil {
		t.Fatalf("open restricted runtime database: %v", err)
	}
	t.Cleanup(func() { _ = runtimeDB.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := runtimeDB.PingContext(ctx); err != nil {
		t.Fatalf("connect restricted runtime database: %v", err)
	}

	var unscopedCount int
	if err := runtimeDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM payments`).Scan(&unscopedCount); err != nil {
		t.Fatalf("unscoped runtime count should fail closed: %v", err)
	}
	if unscopedCount != 0 {
		t.Fatalf("unscoped runtime query exposed %d payment rows", unscopedCount)
	}
	var aggregatePayments int64
	if err := runtimeDB.QueryRowContext(ctx, `SELECT total_payments FROM get_payment_metrics(NULL)`).Scan(&aggregatePayments); err != nil {
		t.Fatalf("restricted runtime aggregate payment metrics: %v", err)
	}
	if aggregatePayments < 1 {
		t.Fatalf("aggregate payment metrics omitted seeded test payment: %d", aggregatePayments)
	}
	deleteTx, err := runtimeDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin restricted runtime delete check: %v", err)
	}
	if _, err := deleteTx.ExecContext(ctx, `SELECT set_config('app.tenant_id', $1, true)`, fmt.Sprint(tenantID)); err != nil {
		_ = deleteTx.Rollback()
		t.Fatalf("set restricted runtime tenant context: %v", err)
	}
	if _, err := deleteTx.ExecContext(ctx, `DELETE FROM payments WHERE id = $1`, paymentID); err == nil {
		_ = deleteTx.Rollback()
		t.Fatal("restricted runtime role unexpectedly deleted a payment row")
	}
	_ = deleteTx.Rollback()

	runtimeRepo := repository.NewPaymentRepository(runtimeDB)
	payment, err := runtimeRepo.GetPaymentByID(ctx, tenantID, paymentID)
	if err != nil || payment.ID != paymentID {
		t.Fatalf("scoped runtime payment read: payment=%+v err=%v", payment, err)
	}
	if _, err := runtimeRepo.GetPaymentByID(ctx, tenantID+1000, paymentID); err == nil {
		t.Fatal("foreign tenant read unexpectedly returned a payment")
	}
	if err := runtimeRepo.UpdatePaymentStatus(ctx, tenantID+1000, paymentID, domain.PaymentStatusFailed, nil, nil, nil, nil, nil); err == nil {
		t.Fatal("foreign tenant update unexpectedly succeeded")
	}

	newPaymentID, err := runtimeRepo.CreatePayment(ctx, &domain.Payment{
		TenantID: tenantID, AmountCents: 50000, Currency: "EGP", Status: domain.PaymentStatusPending,
		PlanID: "monthly", MerchantOrderID: fmt.Sprintf("closure-payment-runtime-%d", time.Now().UnixNano()),
	})
	if err != nil {
		t.Fatalf("scoped runtime insert: %v", err)
	}
	if err := runtimeRepo.UpdatePaymentOrder(ctx, tenantID, newPaymentID, 880099, "runtime-auth"); err != nil {
		t.Fatalf("scoped runtime order update: %v", err)
	}
	var storedTenant int64
	if err := db.QueryRowContext(ctx, `SELECT tenant_id FROM payments WHERE id = $1`, newPaymentID).Scan(&storedTenant); err != nil {
		t.Fatalf("read runtime insert with admin role: %v", err)
	}
	if storedTenant != tenantID {
		t.Fatalf("runtime insert tenant=%d, want %d", storedTenant, tenantID)
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
	if _, err := db.Exec(`ALTER TABLE tenants DROP CONSTRAINT ` + constraint); err != nil {
		t.Fatalf("remove temporary failure condition: %v", err)
	}
	if err := handlePostgresWebhook(t, uc, webhookPayloadForTest(orderID, 50000, true, false, false)); err != nil {
		t.Fatalf("retry provider-success callback after local recovery: %v", err)
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusSuccess, domain.StatusActive)
}

func TestPaymentPostgresCallbackRecoversProviderOrderNotCommittedLocally(t *testing.T) {
	db, uc, tenantID, paymentID, _, _ := openPaymentPostgresFixture(t)
	if _, err := db.Exec(`UPDATE payments SET paymob_order_id = NULL, paymob_auth_token = NULL WHERE id = $1`, paymentID); err != nil {
		t.Fatalf("simulate provider order not persisted locally: %v", err)
	}
	payload := webhookPayloadForTest(880021, 50000, true, false, false)
	payload.Obj.Order.MerchantOrderID = "closure-payment-order"
	if err := handlePostgresWebhook(t, uc, payload); err != nil {
		t.Fatalf("provider-success callback should resolve by durable merchant order id: %v", err)
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusSuccess, domain.StatusActive)
}

func TestPaymentPostgresRefundCallbackReconcilesAndReplaysIdempotently(t *testing.T) {
	db, uc, tenantID, paymentID, orderID, subRepo := openPaymentPostgresFixture(t)
	transactionID := int64(880012)
	if _, err := db.Exec(`UPDATE payments SET status = 'success', paymob_transaction_id = $1 WHERE id = $2`, transactionID, paymentID); err != nil {
		t.Fatalf("prepare successful payment: %v", err)
	}
	subRepo.infoByTenant[tenantID] = &domain.TenantSubscriptionInfo{SubscriptionStatus: domain.StatusActive}
	payload := webhookPayloadForTest(orderID, 50000, true, false, false)
	payload.Obj.ID = transactionID
	payload.Obj.IsRefunded = true
	for attempt := 0; attempt < 2; attempt++ {
		if err := handlePostgresWebhook(t, uc, payload); err != nil {
			t.Fatalf("refund callback replay %d: %v", attempt+1, err)
		}
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusRefunded, domain.StatusExpired)
	if len(subRepo.history) != 1 {
		t.Fatalf("refund callback replay should write history once, got %d entries", len(subRepo.history))
	}
}

func TestPaymentPostgresRefundCallbackBeforeSuccessDoesNotActivateSubscription(t *testing.T) {
	db, uc, tenantID, paymentID, orderID, subRepo := openPaymentPostgresFixture(t)
	subRepo.infoByTenant[tenantID] = &domain.TenantSubscriptionInfo{SubscriptionStatus: domain.StatusTrial}
	payload := webhookPayloadForTest(orderID, 50000, false, false, false)
	payload.Obj.IsRefunded = true
	if err := handlePostgresWebhook(t, uc, payload); err != nil {
		t.Fatalf("refund callback for provider-refunded pending payment: %v", err)
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusRefunded, domain.StatusTrial)
}

func TestPaymentPostgresProviderRefundRecoversAfterLocalCommitFailure(t *testing.T) {
	db, uc, tenantID, paymentID, orderID, subRepo := openPaymentPostgresFixture(t)
	transactionID := int64(880013)
	if _, err := db.Exec(`UPDATE payments SET status = 'success', paymob_transaction_id = $1 WHERE id = $2`, transactionID, paymentID); err != nil {
		t.Fatalf("prepare successful payment: %v", err)
	}
	subRepo.infoByTenant[tenantID] = &domain.TenantSubscriptionInfo{SubscriptionStatus: domain.StatusActive}
	if _, err := db.Exec(`UPDATE tenants SET subscription_status = 'active' WHERE id = $1`, tenantID); err != nil {
		t.Fatalf("activate tenant for refund recovery fixture: %v", err)
	}
	provider := &postgresRacePaymobClient{validHMAC: true}
	uc.paymobClient = provider

	constraint := fmt.Sprintf("closure_reject_refund_%d", tenantID)
	if _, err := db.Exec(`ALTER TABLE payments ADD CONSTRAINT ` + constraint + ` CHECK (id <> ` + fmt.Sprint(paymentID) + ` OR status <> 'refunded')`); err != nil {
		t.Fatalf("install refund commit failure constraint: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Exec(`ALTER TABLE payments DROP CONSTRAINT IF EXISTS ` + constraint) })

	if _, err := uc.RefundPayment(context.Background(), paymentID, 77, nil, "", ""); err == nil {
		t.Fatal("refund request should surface the failed local commit")
	}
	if provider.refundCalls.Load() != 1 {
		t.Fatalf("provider refund calls = %d, want exactly one successful provider call", provider.refundCalls.Load())
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusSuccess, domain.StatusActive)
	var refundState string
	if err := db.QueryRow(`SELECT refund_state FROM payments WHERE id = $1`, paymentID).Scan(&refundState); err != nil {
		t.Fatalf("read durable refund state after local failure: %v", err)
	}
	if refundState != "pending" {
		t.Fatalf("refund state = %q after provider success/local failure, want pending", refundState)
	}
	if _, err := uc.RefundPayment(context.Background(), paymentID, 77, nil, "", ""); err == nil {
		t.Fatal("retry should not reissue provider refund while reconciliation is pending")
	}
	if provider.refundCalls.Load() != 1 {
		t.Fatalf("provider refund retry was duplicated: calls=%d", provider.refundCalls.Load())
	}
	if _, err := db.Exec(`ALTER TABLE payments DROP CONSTRAINT ` + constraint); err != nil {
		t.Fatalf("remove refund commit failure constraint: %v", err)
	}

	refundCallback := webhookPayloadForTest(orderID, 50000, true, false, false)
	refundCallback.Obj.ID = transactionID
	refundCallback.Obj.IsRefunded = true
	if err := handlePostgresWebhook(t, uc, refundCallback); err != nil {
		t.Fatalf("reconcile provider refund callback after local commit failure: %v", err)
	}
	assertPaymentAndTenantState(t, db, paymentID, tenantID, domain.PaymentStatusRefunded, domain.StatusExpired)
	assertPaymentRefundState(t, db, paymentID, "refunded")
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
		} else if !strings.Contains(err.Error(), "payment is no longer refundable") {
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
	assertPaymentRefundState(t, db, paymentID, "refunded")
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
		err = paymentRepo.UpdatePaymentOrder(ctx, tenantID, paymentID, orderID, "synthetic-auth-token")
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

func assertPaymentRefundState(t *testing.T, db *sql.DB, paymentID int64, want string) {
	t.Helper()
	var actual string
	if err := db.QueryRow(`SELECT refund_state FROM payments WHERE id = $1`, paymentID).Scan(&actual); err != nil {
		t.Fatalf("read payment refund state: %v", err)
	}
	if actual != want {
		t.Fatalf("payment refund state = %q, want %q", actual, want)
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
