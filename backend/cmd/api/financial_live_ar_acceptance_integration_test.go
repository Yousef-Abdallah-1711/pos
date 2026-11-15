//go:build integration

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFinancialAcceptanceLiveAPIReceivableLifecycle(t *testing.T) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("POS_FINANCIAL_ACCEPTANCE_API_URL")), "/")
	if baseURL == "" {
		t.Skip("requires the dedicated disposable financial acceptance API process")
	}
	db, ctx, cancel := openFinancialAcceptanceDatabase(t)
	defer cancel()
	defer db.Close()
	fixture := seedLiveFinancialOrderFixture(t, db)
	accountant := seedLiveFinancialIdentity(t, db, fixture, "accountant", "ACCEPTANCE-AR-ACCOUNTANT", []liveFinancialPermission{
		{moduleID: 15, level: "WRITE"}, {moduleID: 17, level: "READ"},
	})
	client := loginLiveFinancialIdentity(t, baseURL, accountant)
	customerID, receivableID := seedLiveReceivableCustomer(t, db, fixture)
	invoiceID := createLiveReceivableInvoice(t, baseURL, fixture, client, customerID)
	createLiveReceivablePayment(t, baseURL, fixture, client, customerID, invoiceID)
	assertLiveReceivableBalances(t, db, ctx, fixture, customerID, invoiceID, receivableID)
	assertLiveReceivablesAging(t, baseURL, client, db, ctx, fixture)
	assertLiveReceivableTenantIsolation(t, baseURL, db, fixture, invoiceID)
	t.Logf("financial_acceptance scenario=authenticated_ar_lifecycle tenant_id=%d restaurant_id=%d customer_id=%d invoice_id=%d invoice_cents=10000 payment_cents=4000 open_ar_cents=6000 subledger_difference_cents=0 receivable_gl_difference_cents=0 aging_difference_cents=0 foreign_invoice_status=404", fixture.tenantID, fixture.restaurantID, customerID, invoiceID)
}

func seedLiveReceivableCustomer(t *testing.T, db *sql.DB, fixture liveFinancialFixture) (int64, int64) {
	t.Helper()
	receivableID := liveFinancialAccountID(t, db, fixture, "1120")
	var customerID int64
	code := fmt.Sprintf("FA-AR-%d", time.Now().UnixNano())
	err := db.QueryRow(`INSERT INTO customers_accounting (tenant_id, restaurant_id, customer_code, customer_name, receivable_account_id, is_active) VALUES ($1,$2,$3,'Synthetic AR Customer',$4,TRUE) RETURNING id`, fixture.tenantID, fixture.restaurantID, code, receivableID).Scan(&customerID)
	if err != nil {
		t.Fatalf("seed synthetic receivable customer: %v", err)
	}
	return customerID, receivableID
}

func createLiveReceivableInvoice(t *testing.T, baseURL string, fixture liveFinancialFixture, client *http.Client, customerID int64) int64 {
	t.Helper()
	today := time.Now().UTC().Format("2006-01-02")
	due := time.Now().UTC().AddDate(0, 0, -35).Format("2006-01-02")
	status, body := acceptanceJSONRequest(t, client, http.MethodPost, baseURL+"/api/v1/accounting/customer-invoices", map[string]any{
		"invoice_date": today, "due_date": due, "customer_id": customerID, "currency": "EGP", "status": "sent",
		"lines": []map[string]any{{"description": "Synthetic service", "quantity": 1, "unit_price": 100, "line_total": 100, "line_order": 1}},
	})
	var response struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	if status != http.StatusCreated || json.Unmarshal(body, &response) != nil || response.Data.ID <= 0 {
		t.Fatalf("authenticated customer invoice creation returned %d or no invoice id: %s", status, body)
	}
	return response.Data.ID
}

func createLiveReceivablePayment(t *testing.T, baseURL string, fixture liveFinancialFixture, client *http.Client, customerID, invoiceID int64) {
	t.Helper()
	status, body := acceptanceJSONRequest(t, client, http.MethodPost, baseURL+"/api/v1/accounting/payments/received", map[string]any{
		"payment_date": time.Now().UTC().Format("2006-01-02"), "customer_id": customerID, "invoice_id": invoiceID,
		"amount": 40, "payment_method": "cash", "reference": fmt.Sprintf("FA-AR-PAY-%d", time.Now().UnixNano()),
	})
	if status != http.StatusCreated {
		t.Fatalf("authenticated partial receivable payment returned %d: %s", status, body)
	}
}

func assertLiveReceivableBalances(t *testing.T, db *sql.DB, ctx context.Context, fixture liveFinancialFixture, customerID, invoiceID, receivableID int64) {
	t.Helper()
	var invoiceBalance, customerBalance, receivableBalance float64
	err := db.QueryRowContext(ctx, `SELECT balance_due FROM invoices_customer WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, invoiceID, fixture.tenantID, fixture.restaurantID).Scan(&invoiceBalance)
	if err != nil {
		t.Fatalf("read receivable invoice balance: %v", err)
	}
	err = db.QueryRowContext(ctx, `SELECT current_balance FROM customers_accounting WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, customerID, fixture.tenantID, fixture.restaurantID).Scan(&customerBalance)
	if err != nil {
		t.Fatalf("read customer subledger balance: %v", err)
	}
	err = db.QueryRowContext(ctx, `SELECT COALESCE(SUM(debit_amount-credit_amount),0) FROM general_ledger WHERE tenant_id=$1 AND restaurant_id=$2 AND account_id=$3`, fixture.tenantID, fixture.restaurantID, receivableID).Scan(&receivableBalance)
	if err != nil {
		t.Fatalf("calculate receivable GL balance: %v", err)
	}
	for name, amount := range map[string]float64{"invoice": invoiceBalance, "customer subledger": customerBalance, "receivable GL": receivableBalance} {
		cents := int64(math.Round(amount * 100))
		if cents != 6000 {
			t.Fatalf("%s open balance=%d cents, expected 6000; difference=%d", name, cents, cents-6000)
		}
	}
}

func assertLiveReceivablesAging(t *testing.T, baseURL string, client *http.Client, db *sql.DB, ctx context.Context, fixture liveFinancialFixture) {
	t.Helper()
	date := time.Now().UTC().Format("2006-01-02")
	status, body := acceptanceJSONRequest(t, client, http.MethodGet, baseURL+"/api/v1/reports/financial/receivables-aging?as_of="+date, nil)
	var response struct {
		Data struct {
			Total float64 `json:"total"`
		} `json:"data"`
	}
	if status != http.StatusOK || json.Unmarshal(body, &response) != nil {
		t.Fatalf("authenticated receivables aging returned %d: %s", status, body)
	}
	var expected float64
	err := db.QueryRowContext(ctx, `SELECT COALESCE(SUM(balance_due),0) FROM invoices_customer WHERE tenant_id=$1 AND restaurant_id=$2 AND status <> 'void' AND balance_due > 0`, fixture.tenantID, fixture.restaurantID).Scan(&expected)
	if err != nil || int64(math.Round(response.Data.Total*100)) != int64(math.Round(expected*100)) || int64(math.Round(expected*100)) != 6000 {
		t.Fatalf("AR aging API/SQL/expected cents=%d/%d/6000 err=%v", int64(math.Round(response.Data.Total*100)), int64(math.Round(expected*100)), err)
	}
}

func assertLiveReceivableTenantIsolation(t *testing.T, baseURL string, db *sql.DB, fixture liveFinancialFixture, invoiceID int64) {
	t.Helper()
	foreignTenantID, foreignRestaurantID, email, password := seedLiveFinancialOtherTenant(t, db)
	client := loginLiveFinancialIdentity(t, baseURL, liveFinancialIdentity{email: email, password: password, role: "foreign-owner"})
	status, body := acceptanceJSONRequest(t, client, http.MethodGet, financialAcceptancePath(baseURL, "/api/v1/accounting/customer-invoices/%d", invoiceID), nil)
	if status != http.StatusNotFound {
		t.Fatalf("foreign tenant customer invoice read returned %d, want 404: %s", status, body)
	}
	if foreignTenantID == fixture.tenantID || foreignRestaurantID == fixture.restaurantID {
		t.Fatal("foreign receivable fixture unexpectedly shares tenant or restaurant scope")
	}
}
