//go:build integration

package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestPaidOrderOutboxRetriesAfterFailureAndDuplicateDelivery(t *testing.T) {
	db := openAccountingEventTestDB(t)
	tenantID, restaurantID, orderID := createPaidOrderEventFixture(t, db)
	events := repository.NewAccountingEventRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	if err := orderRepo.UpdatePaymentStatus(tenantID, restaurantID, orderID, "paid"); err != nil {
		t.Fatalf("mark order paid and enqueue intent: %v", err)
	}
	if err := orderRepo.UpdatePaymentStatus(tenantID, restaurantID, orderID, "paid"); err != nil {
		t.Fatalf("duplicate paid delivery: %v", err)
	}
	parkOtherAccountingEvents(t, db, tenantID, restaurantID, orderID)
	assertAccountingEventState(t, db, tenantID, restaurantID, orderID, repository.AccountingEventPending, 0, 1)

	posting := newPostgresAutoPostingEngine(db)
	processor := NewOrderUseCase(orderRepo, nil)
	processor.SetAutoPostingEngine(posting)
	processor.SetAccountingEventRepository(events)
	processed, err := processor.ProcessNextAccountingEvent(context.Background())
	if !processed || err == nil {
		t.Fatalf("first worker attempt = (%v, %v), want a persisted retryable failure", processed, err)
	}
	assertAccountingEventState(t, db, tenantID, restaurantID, orderID, repository.AccountingEventRetryable, 1, 1)
	seedPaidOrderPostingAccounts(t, db, tenantID, restaurantID)
	if _, err := db.Exec(`UPDATE accounting_event_outbox SET next_attempt_at = CURRENT_TIMESTAMP WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3`, tenantID, restaurantID, orderID); err != nil {
		t.Fatalf("make retry due: %v", err)
	}

	// A new use-case instance models a process restart; the intent remains in PostgreSQL.
	restarted := NewOrderUseCase(repository.NewOrderRepository(db), nil)
	restarted.SetAutoPostingEngine(newPostgresAutoPostingEngine(db))
	restarted.SetAccountingEventRepository(repository.NewAccountingEventRepository(db))
	processed, err = restarted.ProcessNextAccountingEvent(context.Background())
	if !processed || err != nil {
		t.Fatalf("restarted worker attempt = (%v, %v), want successful replay", processed, err)
	}
	assertAccountingEventState(t, db, tenantID, restaurantID, orderID, repository.AccountingEventCompleted, 2, 1)
	reconciliation, err := events.ReconcilePaidOrders(context.Background())
	if err != nil {
		t.Fatalf("reconcile completed paid order: %v", err)
	}
	if reconciliation.MissingIntents != 0 || reconciliation.CompletedWithoutSalesJournal != 0 {
		t.Fatalf("reconciliation after posting=%+v, want no missing intent or journal", reconciliation)
	}
	var journals, ledgerRows int
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entries WHERE tenant_id = $1 AND restaurant_id = $2 AND reference_type = 'order' AND reference_id = $3 AND entry_type = 'auto_sales' AND status = 'posted'`, tenantID, restaurantID, orderID).Scan(&journals); err != nil {
		t.Fatalf("count posted sales journals: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM general_ledger WHERE tenant_id = $1 AND restaurant_id = $2 AND journal_entry_id IN (SELECT id FROM journal_entries WHERE tenant_id = $1 AND restaurant_id = $2 AND reference_type = 'order' AND reference_id = $3 AND entry_type = 'auto_sales')`, tenantID, restaurantID, orderID).Scan(&ledgerRows); err != nil {
		t.Fatalf("count posted sales GL rows: %v", err)
	}
	if journals != 1 || ledgerRows != 2 {
		t.Fatalf("sales accounting effects journals/GL=%d/%d, want exactly 1/2 after retry", journals, ledgerRows)
	}
}

func TestAccountingEventConcurrentClaimsAndExpiredLeaseRecovery(t *testing.T) {
	db := openAccountingEventTestDB(t)
	tenantID, restaurantID, orderID := createPaidOrderEventFixture(t, db)
	if err := repository.NewOrderRepository(db).UpdatePaymentStatus(tenantID, restaurantID, orderID, "paid"); err != nil {
		t.Fatalf("enqueue accounting intent: %v", err)
	}
	parkOtherAccountingEvents(t, db, tenantID, restaurantID, orderID)
	events := repository.NewAccountingEventRepository(db)

	const claims = 8
	var group sync.WaitGroup
	var mu sync.Mutex
	claimed := 0
	var firstClaim *repository.AccountingEvent
	start := make(chan struct{})
	for range claims {
		group.Add(1)
		go func() {
			defer group.Done()
			<-start
			event, err := events.ClaimNext(context.Background())
			if err != nil {
				t.Errorf("claim accounting event: %v", err)
				return
			}
			if event != nil {
				mu.Lock()
				claimed++
				firstClaim = event
				mu.Unlock()
			}
		}()
	}
	close(start)
	group.Wait()
	if claimed != 1 {
		t.Fatalf("concurrent successful claims=%d, want one", claimed)
	}
	assertAccountingEventState(t, db, tenantID, restaurantID, orderID, repository.AccountingEventProcessing, 1, 1)

	// Simulate a process dying after its claim but before it can acknowledge the event.
	if _, err := db.Exec(`UPDATE accounting_event_outbox SET lease_expires_at = CURRENT_TIMESTAMP - INTERVAL '1 second' WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3`, tenantID, restaurantID, orderID); err != nil {
		t.Fatalf("expire simulated worker lease: %v", err)
	}
	recovered, err := repository.NewAccountingEventRepository(db).ClaimNext(context.Background())
	if err != nil || recovered == nil {
		t.Fatalf("claim after simulated restart: event=%v err=%v", recovered, err)
	}
	if recovered.Attempts != 2 || recovered.LeaseGeneration != firstClaim.LeaseGeneration+1 || recovered.EventID != fmt.Sprintf("paid-order:%d:%d:%d", tenantID, restaurantID, orderID) {
		t.Fatalf("recovered event=%+v, want stable identity and second attempt", recovered)
	}
	if err := events.Complete(context.Background(), firstClaim); err == nil {
		t.Fatal("expired worker was allowed to complete a newer worker's lease")
	}
	if err := events.Fail(context.Background(), firstClaim, errors.New("stale worker")); err == nil {
		t.Fatal("expired worker was allowed to fail a newer worker's lease")
	}
	if err := events.Fail(context.Background(), recovered, errors.New("simulated restart before posting")); err != nil {
		t.Fatalf("requeue recovered event: %v", err)
	}
	assertAccountingEventState(t, db, tenantID, restaurantID, orderID, repository.AccountingEventRetryable, 2, 1)
}

func TestPaidOrderReconciliationRestoresMissingIntent(t *testing.T) {
	db := openAccountingEventTestDB(t)
	tenantID, restaurantID, orderID := createPaidOrderEventFixture(t, db)
	if _, err := db.Exec(`UPDATE orders SET payment_status = 'paid' WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, orderID, tenantID, restaurantID); err != nil {
		t.Fatalf("set pre-existing paid state: %v", err)
	}
	result, err := repository.NewAccountingEventRepository(db).ReconcilePaidOrders(context.Background())
	if err != nil {
		t.Fatalf("reconcile missing intent: %v", err)
	}
	if result.MissingIntents != 0 || result.Pending < 1 {
		t.Fatalf("reconciliation after repair=%+v, want at least one pending and no missing intent", result)
	}
	assertAccountingEventState(t, db, tenantID, restaurantID, orderID, repository.AccountingEventPending, 0, 1)
}

func TestOrderSaleInventoryMovementRecoveryIsIdempotent(t *testing.T) {
	db := openAccountingEventTestDB(t)
	tenantID, restaurantID, orderID := createPaidOrderEventFixture(t, db)
	var legacyProductID, newProductID int64
	for _, target := range []*int64{&legacyProductID, &newProductID} {
		if err := db.QueryRow(`
			INSERT INTO products (tenant_id, restaurant_id, name, price, cost)
			VALUES ($1, $2, $3, 10, 5) RETURNING id`, tenantID, restaurantID,
			fmt.Sprintf("recovery-product-%d", time.Now().UnixNano())).Scan(target); err != nil {
			t.Fatalf("create movement product: %v", err)
		}
	}
	legacyMovement := orderSaleMovement(tenantID, restaurantID, orderID, legacyProductID)
	if _, err := repository.NewInventoryMovementRepository(db).CreateMovement(&legacyMovement); err != nil {
		t.Fatalf("create pre-recovery legacy movement: %v", err)
	}
	movementRepo := repository.NewInventoryMovementRepository(db)
	if _, err := movementRepo.CreateOrderSaleMovementIfAbsent(&legacyMovement); err != nil {
		t.Fatalf("recognize legacy movement: %v", err)
	}
	newMovement := orderSaleMovement(tenantID, restaurantID, orderID, newProductID)
	if _, err := db.Exec(`
		CREATE FUNCTION cfr004_sale_movement_failure() RETURNS trigger LANGUAGE plpgsql AS $$
		BEGIN RAISE EXCEPTION 'synthetic inventory movement failure'; END $$;
		CREATE TRIGGER cfr004_sale_movement_failure BEFORE INSERT ON inventory_movements
		FOR EACH ROW EXECUTE FUNCTION cfr004_sale_movement_failure()`); err != nil {
		t.Fatalf("install movement failure trigger: %v", err)
	}
	if _, err := movementRepo.CreateOrderSaleMovementIfAbsent(&newMovement); err == nil {
		t.Fatal("movement write succeeded despite injected failure")
	}
	var markersAfterFailure int
	if err := db.QueryRow(`
		SELECT COUNT(*) FROM accounting_inventory_movement_events
		WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3 AND product_id = $4`,
		tenantID, restaurantID, orderID, newProductID).Scan(&markersAfterFailure); err != nil {
		t.Fatalf("check idempotency marker rollback: %v", err)
	}
	if markersAfterFailure != 0 {
		t.Fatalf("idempotency markers after failed movement=%d, want 0", markersAfterFailure)
	}
	if _, err := db.Exec(`DROP TRIGGER cfr004_sale_movement_failure ON inventory_movements; DROP FUNCTION cfr004_sale_movement_failure()`); err != nil {
		t.Fatalf("remove movement failure trigger: %v", err)
	}
	const deliveries = 8
	var group sync.WaitGroup
	for range deliveries {
		group.Add(1)
		go func() {
			defer group.Done()
			if _, err := movementRepo.CreateOrderSaleMovementIfAbsent(&newMovement); err != nil {
				t.Errorf("idempotent movement delivery: %v", err)
			}
		}()
	}
	group.Wait()
	var movementCount, markerCount int
	if err := db.QueryRow(`
		SELECT COUNT(*) FROM inventory_movements
		WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3
		  AND product_id = $4 AND movement_type = 'sale'`, tenantID, restaurantID, orderID, newProductID).Scan(&movementCount); err != nil {
		t.Fatalf("count new sale movements: %v", err)
	}
	if err := db.QueryRow(`
		SELECT COUNT(*) FROM accounting_inventory_movement_events
		WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3`, tenantID, restaurantID, orderID).Scan(&markerCount); err != nil {
		t.Fatalf("count movement idempotency keys: %v", err)
	}
	if movementCount != 1 || markerCount != 2 {
		t.Fatalf("new movement/key counts=%d/%d, want one movement and keys for both products", movementCount, markerCount)
	}
}

func orderSaleMovement(tenantID, restaurantID, orderID, productID int64) domain.InventoryMovement {
	orderReference := int(orderID)
	return domain.InventoryMovement{
		TenantID: int(tenantID), RestaurantID: int(restaurantID), ProductID: productID,
		OrderID: &orderID, ReferenceType: domain.AutoPostingReferenceTypeOrder,
		ReferenceID: &orderReference, MovementType: domain.InventoryMovementTypeSale,
		Quantity: -2, UnitCost: 5, TotalCost: 10,
	}
}

func TestAccountingEventStopsAfterConfiguredRetryLimit(t *testing.T) {
	db := openAccountingEventTestDB(t)
	tenantID, restaurantID, orderID := createPaidOrderEventFixture(t, db)
	if err := repository.NewOrderRepository(db).UpdatePaymentStatus(tenantID, restaurantID, orderID, "paid"); err != nil {
		t.Fatalf("enqueue accounting intent: %v", err)
	}
	parkOtherAccountingEvents(t, db, tenantID, restaurantID, orderID)
	if _, err := db.Exec(`UPDATE accounting_event_outbox SET max_attempts = 1 WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3`, tenantID, restaurantID, orderID); err != nil {
		t.Fatalf("set bounded attempt limit: %v", err)
	}
	processor := NewOrderUseCase(repository.NewOrderRepository(db), nil)
	processor.SetAccountingEventRepository(repository.NewAccountingEventRepository(db))
	processor.SetAutoPostingEngine(newPostgresAutoPostingEngine(db))
	processed, err := processor.ProcessNextAccountingEvent(context.Background())
	if !processed || err == nil {
		t.Fatalf("final failed attempt=(%v,%v), want durable failure", processed, err)
	}
	assertAccountingEventState(t, db, tenantID, restaurantID, orderID, repository.AccountingEventFailed, 1, 1)
}

func TestPaidStatusAndAccountingIntentRollbackTogether(t *testing.T) {
	db := openAccountingEventTestDB(t)
	tenantID, restaurantID, orderID := createPaidOrderEventFixture(t, db)
	_, err := db.Exec(`
		CREATE FUNCTION cfr004_outbox_insert_failure() RETURNS trigger LANGUAGE plpgsql AS $$
		BEGIN RAISE EXCEPTION 'synthetic outbox failure'; END $$;
		CREATE TRIGGER cfr004_outbox_insert_failure BEFORE INSERT ON accounting_event_outbox
		FOR EACH ROW EXECUTE FUNCTION cfr004_outbox_insert_failure()`)
	if err != nil {
		t.Fatalf("install outbox failure trigger: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.Exec(`DROP TRIGGER IF EXISTS cfr004_outbox_insert_failure ON accounting_event_outbox`)
		_, _ = db.Exec(`DROP FUNCTION IF EXISTS cfr004_outbox_insert_failure()`)
	})
	if err := repository.NewOrderRepository(db).UpdatePaymentStatus(tenantID, restaurantID, orderID, "paid"); err == nil {
		t.Fatal("payment update succeeded despite outbox insert failure")
	}
	var status string
	if err := db.QueryRow(`SELECT payment_status FROM orders WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, orderID, tenantID, restaurantID).Scan(&status); err != nil {
		t.Fatalf("read payment status after rollback: %v", err)
	}
	if status != "pending" {
		t.Fatalf("payment status after failed outbox insert=%q, want pending", status)
	}
	assertAccountingEventState(t, db, tenantID, restaurantID, orderID, "", 0, 0)
}

func TestPaidOrderAccountingIntentPreservesTenantAndRestaurantScope(t *testing.T) {
	db := openAccountingEventTestDB(t)
	tenantID, restaurantID, orderID := createPaidOrderEventFixture(t, db)
	otherTenantID, otherRestaurantID, otherOrderID := createPaidOrderEventFixture(t, db)
	orderRepo := repository.NewOrderRepository(db)
	if err := orderRepo.UpdatePaymentStatus(otherTenantID, restaurantID, orderID, "paid"); err == nil {
		t.Fatal("foreign tenant payment update unexpectedly succeeded")
	}
	if err := orderRepo.UpdatePaymentStatus(tenantID, otherRestaurantID, orderID, "paid"); err == nil {
		t.Fatal("foreign restaurant payment update unexpectedly succeeded")
	}
	var foreignIntentCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM accounting_event_outbox WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3`,
		tenantID, restaurantID, orderID).Scan(&foreignIntentCount); err != nil {
		t.Fatalf("count intents after denied scope changes: %v", err)
	}
	if foreignIntentCount != 0 {
		t.Fatalf("intent count after denied scope changes=%d, want 0", foreignIntentCount)
	}
	if err := orderRepo.UpdatePaymentStatus(tenantID, restaurantID, orderID, "paid"); err != nil {
		t.Fatalf("mark owned order paid: %v", err)
	}
	if err := orderRepo.UpdatePaymentStatus(otherTenantID, otherRestaurantID, otherOrderID, "paid"); err != nil {
		t.Fatalf("mark other tenant's own order paid: %v", err)
	}
	var scopedEvents int
	if err := db.QueryRow(`
		SELECT COUNT(*) FROM accounting_event_outbox
		WHERE (tenant_id = $1 AND restaurant_id = $2 AND order_id = $3)
		   OR (tenant_id = $4 AND restaurant_id = $5 AND order_id = $6)`,
		tenantID, restaurantID, orderID, otherTenantID, otherRestaurantID, otherOrderID).Scan(&scopedEvents); err != nil {
		t.Fatalf("count tenant-scoped intents: %v", err)
	}
	if scopedEvents != 2 {
		t.Fatalf("scoped intent count=%d, want 2 independent tenant events", scopedEvents)
	}
}

func openAccountingEventTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	db.SetMaxOpenConns(20)
	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func createPaidOrderEventFixture(t *testing.T, db *sql.DB) (int64, int64, int64) {
	t.Helper()
	name := fmt.Sprintf("cfr004-event-%d", time.Now().UnixNano())
	var tenantID, restaurantID, orderID int64
	if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1, $2, $3) RETURNING id`, name, name, name+".example.test").Scan(&tenantID); err != nil {
		t.Fatalf("create event tenant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, 'CFR-004 event test', $2, 'test') RETURNING id`, tenantID, name).Scan(&restaurantID); err != nil {
		t.Fatalf("create event restaurant: %v", err)
	}
	if err := db.QueryRow(`
		INSERT INTO orders (tenant_id, restaurant_id, order_number, customer_name, customer_phone, subtotal, total_amount, payment_method, payment_status)
		VALUES ($1, $2, $3, 'Guest', 'N/A', 100, 100, 'cash', 'pending') RETURNING id`, tenantID, restaurantID, name).Scan(&orderID); err != nil {
		t.Fatalf("create event order: %v", err)
	}
	return tenantID, restaurantID, orderID
}

func newPostgresAutoPostingEngine(db *sql.DB) *AutoPostingEngine {
	journalRepo := repository.NewJournalEntryRepository(db)
	return NewAutoPostingEngine(repository.NewAutoPostingRepository(db), repository.NewAccountRepository(db),
		repository.NewInventoryMovementRepository(db), repository.NewGeneralLedgerRepository(db, journalRepo), nil)
}

func seedPaidOrderPostingAccounts(t *testing.T, db *sql.DB, tenantID, restaurantID int64) {
	t.Helper()
	accounts := []struct{ code, name, kind, normal string }{
		{"1110", "Cash", "asset", "debit"},
		{"4100", "Sales", "revenue", "credit"},
		{"2120", "Tax payable", "liability", "credit"},
		{"5110", "Food COGS", "expense", "debit"},
		{"1130", "Inventory", "asset", "debit"},
	}
	for _, account := range accounts {
		if _, err := db.Exec(`
			INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
			VALUES ($1, $2, $3, $4, $5, $6)`, tenantID, restaurantID, account.code, account.name, account.kind, account.normal); err != nil {
			t.Fatalf("create posting account %s: %v", account.code, err)
		}
	}
}

func assertAccountingEventState(t *testing.T, db *sql.DB, tenantID, restaurantID, orderID int64, wantStatus string, wantAttempts, wantCount int) {
	t.Helper()
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM accounting_event_outbox WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3`, tenantID, restaurantID, orderID).Scan(&count); err != nil {
		t.Fatalf("count accounting intents: %v", err)
	}
	if count != wantCount {
		t.Fatalf("accounting intent count=%d, want %d", count, wantCount)
	}
	if wantCount == 0 {
		return
	}
	var status string
	var attempts int
	if err := db.QueryRow(`SELECT status, attempts FROM accounting_event_outbox WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3`, tenantID, restaurantID, orderID).Scan(&status, &attempts); err != nil {
		t.Fatalf("read accounting intent: %v", err)
	}
	if status != wantStatus || attempts != wantAttempts {
		t.Fatalf("accounting intent status/attempts=%s/%d, want %s/%d", status, attempts, wantStatus, wantAttempts)
	}
}

func parkOtherAccountingEvents(t *testing.T, db *sql.DB, tenantID, restaurantID, orderID int64) {
	t.Helper()
	_, err := db.Exec(`
		UPDATE accounting_event_outbox
		SET next_attempt_at = CURRENT_TIMESTAMP + INTERVAL '1 day'
		WHERE status IN ('pending', 'retryable')
		  AND NOT (tenant_id = $1 AND restaurant_id = $2 AND order_id = $3)`, tenantID, restaurantID, orderID)
	if err != nil {
		t.Fatalf("isolate target event from unrelated disposable-fixture events: %v", err)
	}
}
