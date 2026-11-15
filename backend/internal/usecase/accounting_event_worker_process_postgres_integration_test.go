//go:build integration

package usecase

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"pos-saas/internal/repository"
	"pos-saas/internal/service"
	"pos-saas/tests/testdb"
)

const accountingWorkerTestMode = "SERVIO_ACCOUNTING_WORKER_TEST_MODE"

func TestPaidOrderAccountingWorkerSurvivesProcessInterruption(t *testing.T) {
	db := openAccountingEventTestDB(t)
	tenantID, restaurantID, orderID, productID := createInventoryWorkerOrder(t, db)
	orderRepo := repository.NewOrderRepository(db)
	if err := orderRepo.UpdatePaymentStatus(tenantID, restaurantID, orderID, "paid"); err != nil {
		t.Fatalf("persist paid order and accounting event: %v", err)
	}
	parkOtherAccountingEvents(t, db, tenantID, restaurantID, orderID)

	first, firstOutput := startAccountingWorkerProcess(t, "post-and-wait")
	posted := readAccountingWorkerSignal(t, firstOutput, "POSTED")
	if posted != fmt.Sprintf("POSTED paid-order:%d:%d:%d", tenantID, restaurantID, orderID) {
		stopAccountingWorkerProcess(first)
		t.Fatalf("first worker posted %q, want target event", posted)
	}
	assertInterruptedAccountingEvent(t, db, tenantID, restaurantID, orderID)
	parkOtherAccountingEvents(t, db, tenantID, restaurantID, orderID)
	firstPID := first.Process.Pid
	stopAccountingWorkerProcess(first)

	expireAccountingEventLease(t, db, tenantID, restaurantID, orderID)
	second, secondOutput := startAccountingWorkerProcess(t, "complete")
	completed := readAccountingWorkerSignal(t, secondOutput, "COMPLETED")
	if completed != "COMPLETED" {
		stopAccountingWorkerProcess(second)
		t.Fatalf("restarted worker signal=%q, want COMPLETED", completed)
	}
	if err := second.Wait(); err != nil {
		t.Fatalf("wait for restarted worker process: %v", err)
	}
	assertCompletedInventoryPosting(t, db, tenantID, restaurantID, orderID, productID)
	t.Logf("worker_recovery original_pid=%d restarted_pid=%d event=paid-order:%d:%d:%d result=completed",
		firstPID, second.Process.Pid, tenantID, restaurantID, orderID)
}

func TestAccountingWorkerProcessHelper(t *testing.T) {
	mode := os.Getenv(accountingWorkerTestMode)
	if mode == "" {
		t.Skip("subprocess helper")
	}
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL in worker process: %v", err)
	}
	defer db.Close()
	processor := newSubprocessAccountingProcessor(db, mode)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service.NewAccountingEventWorker(processor).Start(ctx)
	select {
	case signal := <-processor.signal:
		fmt.Fprintln(os.Stdout, signal)
		if signal == "IDLE" || strings.HasPrefix(signal, "ERROR ") {
			t.Fatal(signal)
		}
		if mode == "post-and-wait" {
			<-ctx.Done()
		}
	case <-time.After(60 * time.Second):
		t.Fatal("worker process did not reach its synchronized lifecycle point")
	}
}

type subprocessAccountingProcessor struct {
	db        *sql.DB
	mode      string
	signal    chan string
	processor *OrderUseCase
}

func newSubprocessAccountingProcessor(db *sql.DB, mode string) *subprocessAccountingProcessor {
	return &subprocessAccountingProcessor{
		db: db, mode: mode, signal: make(chan string, 1), processor: newAccountingEventOrderProcessor(db),
	}
}

func newAccountingEventOrderProcessor(db *sql.DB) *OrderUseCase {
	orderRepo := repository.NewOrderRepository(db)
	orderProcessor := NewOrderUseCase(orderRepo, nil)
	orderProcessor.SetAutoPostingEngine(newPostgresAutoPostingEngine(db))
	orderProcessor.SetAccountingEventRepository(repository.NewAccountingEventRepository(db))
	return orderProcessor
}

func (p *subprocessAccountingProcessor) ReconcileAccountingEvents(ctx context.Context) error {
	// Reconciliation is isolated so the process test cannot claim unrelated suite fixtures.
	return nil
}

func (p *subprocessAccountingProcessor) ProcessNextAccountingEvent(ctx context.Context) (bool, error) {
	if p.mode == "post-and-wait" {
		return p.postAndWait(ctx)
	}
	processed, err := p.processor.ProcessNextAccountingEvent(ctx)
	if err != nil {
		p.signal <- "ERROR " + err.Error()
	} else if processed {
		p.signal <- "COMPLETED"
	} else {
		p.signal <- "IDLE"
	}
	return processed, err
}

func (p *subprocessAccountingProcessor) postAndWait(ctx context.Context) (bool, error) {
	event, err := repository.NewAccountingEventRepository(p.db).ClaimNext(ctx)
	if err != nil || event == nil {
		return event != nil, err
	}
	order, err := p.processor.GetOrder(event.TenantID, event.RestaurantID, event.OrderID)
	if err != nil {
		p.signal <- "ERROR " + err.Error()
		return true, err
	}
	if err := p.processor.AutoPostPaidOrder(order, intPointer(event.PostedBy)); err != nil {
		p.signal <- "ERROR " + err.Error()
		return true, err
	}
	p.signal <- "POSTED " + event.EventID
	<-ctx.Done()
	return false, ctx.Err()
}

func createInventoryWorkerOrder(t *testing.T, db *sql.DB) (int64, int64, int64, int64) {
	t.Helper()
	tenantID, restaurantID, orderID := createPaidOrderEventFixture(t, db)
	seedPaidOrderPostingAccounts(t, db, tenantID, restaurantID)
	productID := createTrackedWorkerProduct(t, db, tenantID, restaurantID)
	if _, err := db.Exec(`
		INSERT INTO order_items (order_id, tenant_id, restaurant_id, product_id, product_name,
		                         quantity, unit_price, total_price)
		VALUES ($1, $2, $3, $4, 'Synthetic tracked item', 2, 50, 100)`,
		orderID, tenantID, restaurantID, productID); err != nil {
		t.Fatalf("create tracked order item: %v", err)
	}
	return tenantID, restaurantID, orderID, productID
}

func createTrackedWorkerProduct(t *testing.T, db *sql.DB, tenantID, restaurantID int64) int64 {
	t.Helper()
	var productID int64
	name := fmt.Sprintf("worker-restart-product-%d", time.Now().UnixNano())
	err := db.QueryRow(`
		INSERT INTO products (tenant_id, restaurant_id, name, name_en, price, cost, track_inventory)
		VALUES ($1, $2, $3, $3, 50, 5, TRUE) RETURNING id`, tenantID, restaurantID, name).Scan(&productID)
	if err != nil {
		t.Fatalf("create tracked worker product: %v", err)
	}
	return productID
}

func startAccountingWorkerProcess(t *testing.T, mode string) (*exec.Cmd, <-chan string) {
	t.Helper()
	command := exec.Command(os.Args[0], "-test.run=^TestAccountingWorkerProcessHelper$")
	command.Env = append(os.Environ(), accountingWorkerTestMode+"="+mode)
	command.Stderr = os.Stderr
	output, err := command.StdoutPipe()
	if err != nil {
		t.Fatalf("open accounting worker process output: %v", err)
	}
	if err := command.Start(); err != nil {
		t.Fatalf("start accounting worker process: %v", err)
	}
	signals := make(chan string, 64)
	go scanAccountingWorkerSignal(output, signals)
	return command, signals
}

func scanAccountingWorkerSignal(output io.Reader, signals chan<- string) {
	defer close(signals)
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		signals <- line
	}
}

func readAccountingWorkerSignal(t *testing.T, signals <-chan string, want string) string {
	t.Helper()
	var output []string
	deadline := time.After(60 * time.Second)
	for {
		select {
		case signal, ok := <-signals:
			if !ok {
				t.Fatalf("worker process exited before %s signal; output=%s", want, strings.Join(output, " | "))
			}
			output = append(output, signal)
			if signal == "COMPLETED" || signal == "IDLE" || strings.HasPrefix(signal, "POSTED ") || strings.HasPrefix(signal, "ERROR ") {
				return signal
			}
		case <-deadline:
			t.Fatalf("timed out waiting for worker process %s signal; output=%s", want, strings.Join(output, " | "))
			return ""
		}
	}
}

func stopAccountingWorkerProcess(command *exec.Cmd) {
	if command.Process == nil {
		return
	}
	_ = command.Process.Kill()
	_ = command.Wait()
}

func assertInterruptedAccountingEvent(t *testing.T, db *sql.DB, tenantID, restaurantID, orderID int64) {
	t.Helper()
	var status string
	var attempts int
	var journalCount, movementCount, movementQuantity int
	err := db.QueryRow(`
		SELECT e.status, e.attempts,
		       (SELECT COUNT(*) FROM journal_entries j WHERE j.tenant_id=e.tenant_id
		        AND j.restaurant_id=e.restaurant_id AND j.reference_type='order' AND j.reference_id=e.order_id
	        AND j.entry_type IN ('auto_sales','auto_cogs')),
	       (SELECT COUNT(*) FROM inventory_movements m WHERE m.tenant_id=e.tenant_id
	        AND m.restaurant_id=e.restaurant_id AND m.order_id=e.order_id AND m.movement_type='sale'),
	       (SELECT COALESCE(SUM(m.quantity),0)::int FROM inventory_movements m WHERE m.tenant_id=e.tenant_id
	        AND m.restaurant_id=e.restaurant_id AND m.order_id=e.order_id AND m.movement_type='sale')
		FROM accounting_event_outbox e WHERE e.tenant_id=$1 AND e.restaurant_id=$2 AND e.order_id=$3`,
		tenantID, restaurantID, orderID).Scan(&status, &attempts, &journalCount, &movementCount, &movementQuantity)
	if err != nil {
		t.Fatalf("read event after worker interruption: %v", err)
	}
	if status != repository.AccountingEventProcessing || attempts != 1 || journalCount != 2 || movementCount != 1 || movementQuantity != -2 {
		t.Fatalf("interrupted worker state=%s attempts=%d journals=%d movements=%d quantity=%d, want processing/1/2/1/-2 after financial commit", status, attempts, journalCount, movementCount, movementQuantity)
	}
}

func expireAccountingEventLease(t *testing.T, db *sql.DB, tenantID, restaurantID, orderID int64) {
	t.Helper()
	_, err := db.Exec(`UPDATE accounting_event_outbox SET lease_expires_at=CURRENT_TIMESTAMP-INTERVAL '1 second'
		WHERE tenant_id=$1 AND restaurant_id=$2 AND order_id=$3`, tenantID, restaurantID, orderID)
	if err != nil {
		t.Fatalf("expire interrupted worker lease: %v", err)
	}
}

func assertCompletedInventoryPosting(t *testing.T, db *sql.DB, tenantID, restaurantID, orderID, productID int64) {
	t.Helper()
	var status string
	var attempts, salesJournals, cogsJournals, ledgerRows, movements, movementQuantity int
	var sourceRevenueCents, glRevenueCents, sourceCOGSCents, glCOGSCents int64
	err := db.QueryRow(`
		SELECT e.status, e.attempts,
		 (SELECT COUNT(*) FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='order' AND reference_id=$3 AND entry_type='auto_sales' AND status='posted'),
		 (SELECT COUNT(*) FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='order' AND reference_id=$3 AND entry_type='auto_cogs' AND status='posted'),
		 (SELECT COUNT(*) FROM general_ledger WHERE tenant_id=$1 AND restaurant_id=$2 AND journal_entry_id IN (SELECT id FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='order' AND reference_id=$3)),
		 (SELECT COUNT(*) FROM inventory_movements WHERE tenant_id=$1 AND restaurant_id=$2 AND order_id=$3 AND product_id=$4 AND movement_type='sale'),
		 (SELECT COALESCE(SUM(quantity),0)::int FROM inventory_movements WHERE tenant_id=$1 AND restaurant_id=$2 AND order_id=$3 AND product_id=$4 AND movement_type='sale'),
		 (SELECT ROUND(o.total_amount*100)::bigint FROM orders o WHERE o.id=$3 AND o.tenant_id=$1 AND o.restaurant_id=$2),
		 (SELECT ROUND(COALESCE(SUM(gl.credit_amount-gl.debit_amount),0)*100)::bigint FROM general_ledger gl JOIN journal_entries j ON j.id=gl.journal_entry_id AND j.tenant_id=gl.tenant_id AND j.restaurant_id=gl.restaurant_id JOIN chart_of_accounts a ON a.id=gl.account_id AND a.tenant_id=gl.tenant_id AND a.restaurant_id=gl.restaurant_id WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND j.reference_type='order' AND j.reference_id=$3 AND j.entry_type='auto_sales' AND a.account_type='revenue'),
		 (SELECT ROUND(COALESCE(SUM(oi.quantity*p.cost),0)*100)::bigint FROM order_items oi JOIN products p ON p.id=oi.product_id AND p.tenant_id=oi.tenant_id AND p.restaurant_id=oi.restaurant_id WHERE oi.order_id=$3 AND oi.tenant_id=$1 AND oi.restaurant_id=$2),
		 (SELECT ROUND(COALESCE(SUM(gl.debit_amount-gl.credit_amount),0)*100)::bigint FROM general_ledger gl JOIN journal_entries j ON j.id=gl.journal_entry_id AND j.tenant_id=gl.tenant_id AND j.restaurant_id=gl.restaurant_id JOIN chart_of_accounts a ON a.id=gl.account_id AND a.tenant_id=gl.tenant_id AND a.restaurant_id=gl.restaurant_id WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND j.reference_type='order' AND j.reference_id=$3 AND j.entry_type='auto_cogs' AND a.account_type='expense')
		FROM accounting_event_outbox e WHERE e.tenant_id=$1 AND e.restaurant_id=$2 AND e.order_id=$3`,
		tenantID, restaurantID, orderID, productID).Scan(&status, &attempts, &salesJournals, &cogsJournals, &ledgerRows, &movements, &movementQuantity,
		&sourceRevenueCents, &glRevenueCents, &sourceCOGSCents, &glCOGSCents)
	if err != nil {
		t.Fatalf("reconcile recovered order journals and inventory: %v", err)
	}
	if status != repository.AccountingEventCompleted || attempts != 2 || salesJournals != 1 || cogsJournals != 1 || ledgerRows != 4 || movements != 1 || movementQuantity != -2 || sourceRevenueCents != glRevenueCents || sourceCOGSCents != glCOGSCents {
		t.Fatalf("recovered state=%s attempts=%d sales/cogs journals=%d/%d GL=%d movements=%d qty=%d source/GL revenue=%d/%d cents COGS=%d/%d cents", status, attempts, salesJournals, cogsJournals, ledgerRows, movements, movementQuantity, sourceRevenueCents, glRevenueCents, sourceCOGSCents, glCOGSCents)
	}
	t.Logf("inventory_worker_reconciliation source_revenue_cents=%d gl_revenue_cents=%d difference_cents=%d source_cogs_cents=%d gl_cogs_cents=%d difference_cents=%d movement_count=%d movement_quantity=%d", sourceRevenueCents, glRevenueCents, glRevenueCents-sourceRevenueCents, sourceCOGSCents, glCOGSCents, glCOGSCents-sourceCOGSCents, movements, movementQuantity)
}
