//go:build integration

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestFinancialAcceptanceLiveAPISupplierRollbackAndConcurrentPayments(t *testing.T) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("POS_FINANCIAL_ACCEPTANCE_API_URL")), "/")
	if baseURL == "" {
		t.Skip("requires the dedicated disposable financial acceptance API process")
	}
	db, ctx, cancel := openSupplierAcceptanceDB(t)
	defer cancel()
	defer db.Close()
	fixture := seedLiveFinancialOrderFixture(t, db)
	api := &supplierAcceptanceAPI{baseURL: baseURL, db: db, ctx: ctx, fixture: fixture, client: supplierAcceptanceOwnerClient(t, baseURL, fixture)}
	supplierID := createAcceptanceSupplier(t, api)
	invoiceID := createAcceptanceSupplierInvoice(t, api, supplierID)
	journalID := assertSupplierPostingRollback(t, api, invoiceID)
	assertSupplierPaymentRollback(t, api, supplierID, invoiceID)
	assertConcurrentSupplierPayments(t, api, supplierID, invoiceID, journalID)
}

func assertSupplierPostingRollback(t *testing.T, api *supplierAcceptanceAPI, invoiceID int) int64 {
	t.Helper()
	removeTrigger := installSupplierJournalFailure(t, api.db, invoiceID, "supplier_invoice")
	defer removeTrigger()
	path := fmt.Sprintf("%s/api/v1/accounting/supplier-invoices/%d/approve", api.baseURL, invoiceID)
	status, body := acceptanceJSONRequest(t, api.client, http.MethodPost, path, map[string]any{})
	removeTrigger()
	if status < 400 {
		t.Fatalf("forced supplier invoice posting failure returned %d: %s", status, body)
	}
	assertSupplierInvoiceState(t, api, invoiceID, "draft", 0, 0)
	if countSupplierInvoiceJournals(t, api, invoiceID) != 0 {
		t.Fatal("failed supplier invoice approval left a journal behind")
	}
	retryStatus, retryBody := acceptanceJSONRequest(t, api.client, http.MethodPost, path, map[string]any{})
	if retryStatus != http.StatusOK {
		t.Fatalf("supplier invoice retry after rollback returned %d: %s", retryStatus, retryBody)
	}
	status, body = acceptanceJSONRequest(t, api.client, http.MethodPost, path, map[string]any{})
	if status < 400 || countSupplierInvoiceJournals(t, api, invoiceID) != 1 {
		t.Fatalf("repeated supplier approval status=%d expected rejection and one journal: %s", status, body)
	}
	var journalID int64
	if err := api.db.QueryRowContext(api.ctx, `SELECT journal_entry_id FROM invoices_supplier WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, invoiceID, api.fixture.tenantID, api.fixture.restaurantID).Scan(&journalID); err != nil || journalID <= 0 {
		t.Fatalf("read supplier invoice journal after retry: journal_id=%d err=%v", journalID, err)
	}
	return journalID
}

func assertSupplierPaymentRollback(t *testing.T, api *supplierAcceptanceAPI, supplierID, invoiceID int) {
	t.Helper()
	removeTrigger := installSupplierJournalFailure(t, api.db, invoiceID, "payment_made")
	defer removeTrigger()
	status, body := supplierPaymentRequest(t.Context(), api.client, api.baseURL, supplierID, invoiceID, 220)
	removeTrigger()
	if status < 400 {
		t.Fatalf("forced supplier payment posting failure returned %d: %s", status, body)
	}
	assertSupplierInvoiceState(t, api, invoiceID, "approved", 0, 22000)
	var paymentCount, journalCount, lineCount, ledgerCount int
	err := api.db.QueryRowContext(api.ctx, `
		SELECT COUNT(pm.id), COUNT(je.id), COUNT(jel.id), COUNT(gl.id)
		FROM invoices_supplier inv
		LEFT JOIN payments_made pm ON pm.invoice_id=inv.id AND pm.tenant_id=inv.tenant_id AND pm.restaurant_id=inv.restaurant_id
		LEFT JOIN journal_entries je ON je.reference_type='payment_made' AND je.reference_id=pm.id AND je.tenant_id=inv.tenant_id AND je.restaurant_id=inv.restaurant_id
		LEFT JOIN journal_entry_lines jel ON jel.journal_entry_id=je.id AND jel.tenant_id=je.tenant_id AND jel.restaurant_id=je.restaurant_id
		LEFT JOIN general_ledger gl ON gl.journal_entry_id=je.id AND gl.tenant_id=je.tenant_id AND gl.restaurant_id=je.restaurant_id
		WHERE inv.id=$1 AND inv.tenant_id=$2 AND inv.restaurant_id=$3`, invoiceID, api.fixture.tenantID, api.fixture.restaurantID).Scan(&paymentCount, &journalCount, &lineCount, &ledgerCount)
	if err != nil || paymentCount+journalCount+lineCount+ledgerCount != 0 {
		t.Fatalf("failed payment rollback left payments/journals/lines/GL=%d/%d/%d/%d err=%v", paymentCount, journalCount, lineCount, ledgerCount, err)
	}
}

func assertConcurrentSupplierPayments(t *testing.T, api *supplierAcceptanceAPI, supplierID, invoiceID int, invoiceJournalID int64) {
	t.Helper()
	succeeded := submitConcurrentSupplierPayments(t, api, supplierID, invoiceID)
	paid, balance, payments := readSupplierPaymentTotals(t, api, invoiceID)
	if succeeded != 4 || paid != 20000 || balance != 2000 || payments != 4 {
		t.Fatalf("concurrent supplier payments expected success/source/balance/count=4/20000/2000/4; actual=%d/%d/%d/%d", succeeded, paid, balance, payments)
	}
	assertSupplierAPBalance(t, api, invoiceID, invoiceJournalID, balance)
	status, body := supplierPaymentRequest(t, api.client, api.baseURL, supplierID, invoiceID, 21)
	if status < 400 {
		t.Fatalf("supplier payment exceeding remaining 20.00 balance returned %d: %s", status, body)
	}
	assertSupplierInvoiceState(t, api, invoiceID, "partially_paid", 20000, 2000)
	if err := api.db.QueryRowContext(api.ctx, `SELECT COUNT(*) FROM payments_made WHERE invoice_id=$1 AND tenant_id=$2 AND restaurant_id=$3`, invoiceID, api.fixture.tenantID, api.fixture.restaurantID).Scan(&payments); err != nil || payments != 4 {
		t.Fatalf("overpayment rejection expected four persisted payments; actual=%d err=%v", payments, err)
	}
	t.Logf("financial_acceptance scenario=supplier_rollback_and_concurrent_payments expected_successes=4 actual_successes=%d expected_paid_cents=20000 actual_paid_cents=%d expected_balance_cents=2000 actual_balance_cents=%d expected_ap_net_credit_cents=2000", succeeded, paid, balance)
}

func submitConcurrentSupplierPayments(t *testing.T, api *supplierAcceptanceAPI, supplierID, invoiceID int) int {
	t.Helper()
	const requests = 5
	start := make(chan struct{})
	statuses := make(chan int, requests)
	errors := make(chan error, requests)
	var workers sync.WaitGroup
	for range requests {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			status, _, err := doSupplierPaymentRequest(api.ctx, api.client, api.baseURL, supplierID, invoiceID, 50)
			if err != nil {
				errors <- err
				return
			}
			statuses <- status
		}()
	}
	close(start)
	workers.Wait()
	close(statuses)
	close(errors)
	for err := range errors {
		t.Fatalf("concurrent synthetic payment request: %v", err)
	}
	succeeded := 0
	for status := range statuses {
		if status == http.StatusCreated {
			succeeded++
		} else if status < 400 {
			t.Fatalf("unexpected concurrent payment status %d", status)
		}
	}
	return succeeded
}

func readSupplierPaymentTotals(t *testing.T, api *supplierAcceptanceAPI, invoiceID int) (int64, int64, int) {
	t.Helper()
	var paid, balance int64
	var payments int
	err := api.db.QueryRowContext(api.ctx, `SELECT ROUND(amount_paid*100)::bigint, ROUND(balance_due*100)::bigint,
		(SELECT COUNT(*) FROM payments_made WHERE invoice_id=inv.id AND tenant_id=inv.tenant_id AND restaurant_id=inv.restaurant_id)
		FROM invoices_supplier inv WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, invoiceID, api.fixture.tenantID, api.fixture.restaurantID).Scan(&paid, &balance, &payments)
	if err != nil {
		t.Fatalf("read supplier payment source totals: %v", err)
	}
	return paid, balance, payments
}

func assertSupplierAPBalance(t *testing.T, api *supplierAcceptanceAPI, invoiceID int, invoiceJournalID int64, balance int64) {
	t.Helper()
	var apNetCredit int64
	err = api.db.QueryRowContext(api.ctx, `
		SELECT ROUND(COALESCE(SUM(gl.credit_amount-gl.debit_amount),0)*100)::bigint
		FROM general_ledger gl JOIN chart_of_accounts a ON a.id=gl.account_id AND a.tenant_id=gl.tenant_id AND a.restaurant_id=gl.restaurant_id
		JOIN journal_entries je ON je.id=gl.journal_entry_id AND je.tenant_id=gl.tenant_id AND je.restaurant_id=gl.restaurant_id
		WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND a.account_code='2101' AND je.status='posted'
		AND ((je.reference_type='supplier_invoice' AND je.reference_id=$3) OR
		 (je.reference_type='payment_made' AND je.reference_id IN (SELECT id FROM payments_made WHERE invoice_id=$3 AND tenant_id=$1 AND restaurant_id=$2)))`,
		api.fixture.tenantID, api.fixture.restaurantID, invoiceID).Scan(&apNetCredit)
	if err != nil || apNetCredit != balance || invoiceJournalID <= 0 {
		t.Fatalf("supplier AP reconciliation expected GL net credit=%d, actual=%d invoice_journal=%d err=%v", balance, apNetCredit, invoiceJournalID, err)
	}
}

func installSupplierJournalFailure(t *testing.T, db *sql.DB, id int, entryType string) func() {
	t.Helper()
	trigger := fmt.Sprintf("financial_acceptance_%d_%s_trigger", id, entryType)
	function := fmt.Sprintf("financial_acceptance_%d_%s_failure", id, entryType)
	_, _ = db.Exec(`DROP TRIGGER IF EXISTS ` + trigger + ` ON journal_entries`)
	_, _ = db.Exec(`DROP FUNCTION IF EXISTS ` + function + `()`)
	if _, err := db.Exec(`CREATE FUNCTION ` + function + `() RETURNS trigger LANGUAGE plpgsql AS $body$
		BEGIN RAISE EXCEPTION 'synthetic financial acceptance posting failure'; END $body$`); err != nil {
		t.Fatalf("install synthetic posting failure function: %v", err)
	}
	if _, err := db.Exec(`CREATE TRIGGER ` + trigger + ` BEFORE UPDATE OF status ON journal_entries
		FOR EACH ROW WHEN (NEW.status='posted' AND NEW.entry_type='` + entryType + `') EXECUTE FUNCTION ` + function + `()`); err != nil {
		t.Fatalf("install synthetic posting failure trigger: %v", err)
	}
	return func() {
		if _, err := db.Exec(`DROP TRIGGER IF EXISTS ` + trigger + ` ON journal_entries`); err != nil {
			t.Errorf("remove synthetic posting failure trigger: %v", err)
		}
		if _, err := db.Exec(`DROP FUNCTION IF EXISTS ` + function + `()`); err != nil {
			t.Errorf("remove synthetic posting failure function: %v", err)
		}
	}
}

func supplierPaymentRequest(t *testing.T, client *http.Client, baseURL string, supplierID, invoiceID int, amount float64) (int, []byte) {
	t.Helper()
	status, body, err := doSupplierPaymentRequest(t.Context(), client, baseURL, supplierID, invoiceID, amount)
	if err != nil {
		t.Fatalf("send supplier payment request: %v", err)
	}
	return status, body
}

func doSupplierPaymentRequest(ctx context.Context, client *http.Client, baseURL string, supplierID, invoiceID int, amount float64) (int, []byte, error) {
	payload, err := json.Marshal(map[string]any{"payment_date": time.Now().UTC().Format("2006-01-02"), "supplier_id": supplierID,
		"invoice_id": invoiceID, "amount": amount, "payment_method": "cash", "reference": "FIN-ACCEPTANCE-CONCURRENT"})
	if err != nil {
		return 0, nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/api/v1/accounting/payments/made", strings.NewReader(string(payload)))
	if err != nil {
		return 0, nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "http://localhost:3002")
	response, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<16))
	if err != nil {
		return 0, nil, err
	}
	return response.StatusCode, body, nil
}

func assertSupplierInvoiceState(t *testing.T, api *supplierAcceptanceAPI, invoiceID int, expectedState string, expectedPaid, expectedBalance int64) {
	t.Helper()
	var state string
	var paid, balance int64
	err := api.db.QueryRowContext(api.ctx, `SELECT status, ROUND(amount_paid*100)::bigint, ROUND(balance_due*100)::bigint
		FROM invoices_supplier WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, invoiceID, api.fixture.tenantID, api.fixture.restaurantID).Scan(&state, &paid, &balance)
	if err != nil || state != expectedState || paid != expectedPaid || balance != expectedBalance {
		t.Fatalf("supplier invoice expected state/paid/balance=%s/%d/%d; actual=%s/%d/%d err=%v", expectedState, expectedPaid, expectedBalance, state, paid, balance, err)
	}
}

func countSupplierInvoiceJournals(t *testing.T, api *supplierAcceptanceAPI, invoiceID int) int {
	t.Helper()
	var count int
	err := api.db.QueryRowContext(api.ctx, `SELECT COUNT(*) FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2
		AND reference_type='supplier_invoice' AND reference_id=$3`, api.fixture.tenantID, api.fixture.restaurantID, invoiceID).Scan(&count)
	if err != nil {
		t.Fatalf("count supplier invoice journals: %v", err)
	}
	return count
}
