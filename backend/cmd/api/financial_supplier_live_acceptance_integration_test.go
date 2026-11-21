//go:build integration

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/tests/testdb"
)

type supplierAcceptanceAPI struct {
	baseURL string
	db      *sql.DB
	ctx     context.Context
	client  *http.Client
	fixture liveFinancialFixture
}

type supplierAcceptanceResult struct {
	supplierID, invoiceID, paymentID                int
	invoiceJournalID, paymentJournalID              int64
	totalCents, invoiceBalanceCents                 int64
	expenseCents, inputVATCents, payableCreditCents int64
	payableSettlementCents, cashCreditCents         int64
	invoiceState                                    string
}

func TestFinancialAcceptanceLiveAPISupplierLifecycle(t *testing.T) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("POS_FINANCIAL_ACCEPTANCE_API_URL")), "/")
	if baseURL == "" {
		t.Skip("requires the dedicated disposable financial acceptance API process")
	}
	db, ctx, cancel := openSupplierAcceptanceDB(t)
	defer cancel()
	defer db.Close()
	fixture := seedLiveFinancialOrderFixture(t, db)
	api := &supplierAcceptanceAPI{baseURL: baseURL, db: db, ctx: ctx, fixture: fixture, client: supplierAcceptanceOwnerClient(t, baseURL, fixture)}
	result := runSupplierAcceptanceLifecycle(t, api)
	logSupplierAcceptanceResult(t, fixture, result)
}

func openSupplierAcceptanceDB(t *testing.T) (*sql.DB, context.Context, context.CancelFunc) {
	t.Helper()
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable acceptance database: %v", err)
	}
	db.SetMaxOpenConns(4)
	ctx, cancel := context.WithTimeout(t.Context(), 45*time.Second)
	if err := db.PingContext(ctx); err != nil {
		cancel()
		db.Close()
		t.Fatalf("connect acceptance database: %v", err)
	}
	return db, ctx, cancel
}

func supplierAcceptanceOwnerClient(t *testing.T, baseURL string, fixture liveFinancialFixture) *http.Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("create isolated cookie jar: %v", err)
	}
	client := &http.Client{Jar: jar, Timeout: 30 * time.Second}
	status, body := acceptanceJSONRequest(t, client, http.MethodPost, baseURL+"/api/v1/auth/login", map[string]string{
		"email": fixture.email, "password": fixture.password,
	})
	if status != http.StatusOK {
		t.Fatalf("synthetic owner login returned %d: %s", status, body)
	}
	return client
}

func runSupplierAcceptanceLifecycle(t *testing.T, api *supplierAcceptanceAPI) supplierAcceptanceResult {
	t.Helper()
	result := supplierAcceptanceResult{totalCents: 22000}
	result.supplierID = createAcceptanceSupplier(t, api)
	result.invoiceID = createAcceptanceSupplierInvoice(t, api, result.supplierID)
	result.invoiceJournalID = approveAcceptanceSupplierInvoice(t, api, result.invoiceID)
	result.paymentID = payAcceptanceSupplierInvoice(t, api, result.supplierID, result.invoiceID)
	reconcileAcceptanceSupplierLifecycle(t, api, &result)
	return result
}

func createAcceptanceSupplier(t *testing.T, api *supplierAcceptanceAPI) int {
	t.Helper()
	status, body := acceptanceJSONRequest(t, api.client, http.MethodPost, api.baseURL+"/api/v1/suppliers", map[string]any{
		"company_name": "Synthetic Financial Acceptance Supplier", "currency": "EGP",
	})
	var response struct {
		Data struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	if status != http.StatusCreated || json.Unmarshal(body, &response) != nil || response.Data.ID <= 0 {
		t.Fatalf("authenticated supplier create returned %d or no id: %s", status, body)
	}
	return response.Data.ID
}

func createAcceptanceSupplierInvoice(t *testing.T, api *supplierAcceptanceAPI, supplierID int) int {
	t.Helper()
	today := time.Now().UTC().Format("2006-01-02")
	status, body := acceptanceJSONRequest(t, api.client, http.MethodPost, api.baseURL+"/api/v1/accounting/supplier-invoices", map[string]any{
		"supplier_ref": "FIN-ACCEPTANCE-" + today, "invoice_date": today,
		"due_date":    time.Now().UTC().AddDate(0, 0, 30).Format("2006-01-02"),
		"supplier_id": supplierID, "currency": "EGP",
		"lines": []map[string]any{{"description": "Synthetic ingredients", "quantity": 2, "unit_price": 100, "tax_rate": 10}},
	})
	var response struct {
		Data struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	if status != http.StatusCreated || json.Unmarshal(body, &response) != nil || response.Data.ID <= 0 {
		t.Fatalf("authenticated supplier invoice create returned %d or no id: %s", status, body)
	}
	return response.Data.ID
}

func approveAcceptanceSupplierInvoice(t *testing.T, api *supplierAcceptanceAPI, invoiceID int) int64 {
	t.Helper()
	path := fmt.Sprintf("%s/api/v1/accounting/supplier-invoices/%d/approve", api.baseURL, invoiceID)
	status, body := acceptanceJSONRequest(t, api.client, http.MethodPost, path, map[string]any{})
	if status != http.StatusOK {
		t.Fatalf("authenticated supplier invoice approval returned %d: %s", status, body)
	}
	var journalID int64
	var totalCents, paidCents, balanceCents int64
	var invoiceState string
	err := api.db.QueryRowContext(api.ctx, `
		SELECT ROUND(total_amount * 100)::bigint, ROUND(amount_paid * 100)::bigint,
		       ROUND(balance_due * 100)::bigint, status, journal_entry_id
		FROM invoices_supplier WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`,
		invoiceID, api.fixture.tenantID, api.fixture.restaurantID).Scan(&totalCents, &paidCents, &balanceCents, &invoiceState, &journalID)
	if err != nil || totalCents != 22000 || paidCents != 0 || balanceCents != 22000 || invoiceState != "approved" || journalID <= 0 {
		t.Fatalf("approved invoice expected total/paid/balance=22000/0/22000 status=approved journal>0; actual=%d/%d/%d status=%s journal=%d err=%v",
			totalCents, paidCents, balanceCents, invoiceState, journalID, err)
	}
	return journalID
}

func payAcceptanceSupplierInvoice(t *testing.T, api *supplierAcceptanceAPI, supplierID, invoiceID int) int {
	t.Helper()
	status, body := acceptanceJSONRequest(t, api.client, http.MethodPost, api.baseURL+"/api/v1/accounting/payments/made", map[string]any{
		"payment_date": time.Now().UTC().Format("2006-01-02"), "supplier_id": supplierID,
		"invoice_id": invoiceID, "amount": 220, "payment_method": "cash", "reference": "FIN-ACCEPTANCE-SUPPLIER",
	})
	var response struct {
		Data struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	if status != http.StatusCreated || json.Unmarshal(body, &response) != nil || response.Data.ID <= 0 {
		t.Fatalf("authenticated supplier payment returned %d or no id: %s", status, body)
	}
	return response.Data.ID
}

func reconcileAcceptanceSupplierLifecycle(t *testing.T, api *supplierAcceptanceAPI, result *supplierAcceptanceResult) {
	t.Helper()
	readPaidSupplierInvoice(t, api, result)
	readSupplierJournalIDs(t, api, result)
	readSupplierGLTotals(t, api, result)
	if result.expenseCents != 20000 || result.inputVATCents != 2000 || result.payableCreditCents != result.totalCents ||
		result.payableSettlementCents != result.totalCents || result.cashCreditCents != result.totalCents {
		t.Fatalf("independent supplier totals subtotal/tax/total=20000/2000/%d; GL expense/inputVAT/AP-credit/AP-settlement/cash-credit=%d/%d/%d/%d/%d",
			result.totalCents, result.expenseCents, result.inputVATCents, result.payableCreditCents,
			result.payableSettlementCents, result.cashCreditCents)
	}
}

func readPaidSupplierInvoice(t *testing.T, api *supplierAcceptanceAPI, result *supplierAcceptanceResult) {
	t.Helper()
	var totalCents, paidCents int64
	err := api.db.QueryRowContext(api.ctx, `
		SELECT ROUND(total_amount * 100)::bigint, ROUND(amount_paid * 100)::bigint,
		       ROUND(balance_due * 100)::bigint, status
		FROM invoices_supplier WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`,
		result.invoiceID, api.fixture.tenantID, api.fixture.restaurantID).Scan(
		&totalCents, &paidCents, &result.invoiceBalanceCents, &result.invoiceState)
	if err != nil || totalCents != result.totalCents || paidCents != result.totalCents || result.invoiceBalanceCents != 0 || result.invoiceState != "paid" {
		t.Fatalf("paid invoice expected total/paid/balance=%d/%d/0 status=paid; actual=%d/%d/%d status=%s err=%v",
			result.totalCents, result.totalCents, totalCents, paidCents, result.invoiceBalanceCents, result.invoiceState, err)
	}
}

func readSupplierJournalIDs(t *testing.T, api *supplierAcceptanceAPI, result *supplierAcceptanceResult) {
	t.Helper()
	var invoiceCount, paymentCount int
	if err := api.db.QueryRowContext(api.ctx, `
		SELECT COUNT(*), COALESCE(MIN(id), 0) FROM journal_entries
		WHERE tenant_id = $1 AND restaurant_id = $2 AND reference_type = 'supplier_invoice'
		  AND reference_id = $3 AND entry_type = 'supplier_invoice' AND status = 'posted'`,
		api.fixture.tenantID, api.fixture.restaurantID, result.invoiceID).Scan(&invoiceCount, &result.invoiceJournalID); err != nil {
		t.Fatalf("read posted supplier invoice journal: %v", err)
	}
	err := api.db.QueryRowContext(api.ctx, `
		SELECT COUNT(*), COALESCE(MIN(je.id), 0) FROM payments_made pm
		JOIN journal_entries je ON je.id = pm.journal_entry_id AND je.tenant_id = pm.tenant_id AND je.restaurant_id = pm.restaurant_id
		WHERE pm.id = $1 AND pm.tenant_id = $2 AND pm.restaurant_id = $3 AND je.status = 'posted'`,
		result.paymentID, api.fixture.tenantID, api.fixture.restaurantID).Scan(&paymentCount, &result.paymentJournalID)
	if err != nil || invoiceCount != 1 || paymentCount != 1 || result.invoiceJournalID <= 0 || result.paymentJournalID <= 0 {
		t.Fatalf("supplier journals invoice=%d/%d payment=%d/%d; expected one posted journal each; err=%v",
			invoiceCount, result.invoiceJournalID, paymentCount, result.paymentJournalID, err)
	}
}

func readSupplierGLTotals(t *testing.T, api *supplierAcceptanceAPI, result *supplierAcceptanceResult) {
	t.Helper()
	err := api.db.QueryRowContext(api.ctx, `
		SELECT
		  COALESCE(ROUND(SUM(gl.debit_amount - gl.credit_amount) FILTER (WHERE a.account_code = '5110' AND gl.journal_entry_id = $3) * 100), 0)::bigint,
		  COALESCE(ROUND(SUM(gl.debit_amount - gl.credit_amount) FILTER (WHERE a.account_code = '1150' AND gl.journal_entry_id = $3) * 100), 0)::bigint,
		  COALESCE(ROUND(SUM(gl.credit_amount - gl.debit_amount) FILTER (WHERE a.account_code = '2101' AND gl.journal_entry_id = $3) * 100), 0)::bigint,
		  COALESCE(ROUND(SUM(gl.debit_amount - gl.credit_amount) FILTER (WHERE a.account_code = '2101' AND gl.journal_entry_id = $4) * 100), 0)::bigint,
		  COALESCE(ROUND(SUM(gl.credit_amount - gl.debit_amount) FILTER (WHERE a.account_code = '1110' AND gl.journal_entry_id = $4) * 100), 0)::bigint
		FROM general_ledger gl JOIN chart_of_accounts a
		  ON a.id = gl.account_id AND a.tenant_id = gl.tenant_id AND a.restaurant_id = gl.restaurant_id
		WHERE gl.tenant_id = $1 AND gl.restaurant_id = $2 AND gl.journal_entry_id IN ($3, $4)`,
		api.fixture.tenantID, api.fixture.restaurantID, result.invoiceJournalID, result.paymentJournalID).Scan(
		&result.expenseCents, &result.inputVATCents, &result.payableCreditCents,
		&result.payableSettlementCents, &result.cashCreditCents)
	if err != nil {
		t.Fatalf("independently reconcile supplier GL totals: %v", err)
	}
}

func logSupplierAcceptanceResult(t *testing.T, fixture liveFinancialFixture, result supplierAcceptanceResult) {
	t.Logf("financial_acceptance scenario=authenticated_supplier_invoice_to_payment tenant_id=%d restaurant_id=%d supplier_id=%d invoice_id=%d source_subtotal_cents=20000 source_tax_cents=2000 source_total_cents=%d invoice_journal_id=%d payment_id=%d payment_journal_id=%d expense_debit_cents=%d input_vat_debit_cents=%d ap_invoice_credit_cents=%d ap_payment_debit_cents=%d cash_credit_cents=%d invoice_balance_cents=%d invoice_status=%s",
		fixture.tenantID, fixture.restaurantID, result.supplierID, result.invoiceID, result.totalCents,
		result.invoiceJournalID, result.paymentID, result.paymentJournalID, result.expenseCents,
		result.inputVATCents, result.payableCreditCents, result.payableSettlementCents,
		result.cashCreditCents, result.invoiceBalanceCents, result.invoiceState)
}
