//go:build integration

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFinancialAcceptanceLiveAPIPayrollLifecycle(t *testing.T) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("POS_FINANCIAL_ACCEPTANCE_API_URL")), "/")
	if baseURL == "" {
		t.Skip("requires the dedicated disposable financial acceptance API process")
	}
	db, ctx, cancel := openFinancialAcceptanceDatabase(t)
	defer cancel()
	defer db.Close()
	fixture := seedLiveFinancialOrderFixture(t, db)
	payrollUser := seedLiveFinancialIdentity(t, db, fixture, "payroll-manager", "ACCEPTANCE-PAYROLL", []liveFinancialPermission{{moduleID: 12, level: "DELETE"}})
	client := loginLiveFinancialIdentity(t, baseURL, payrollUser)
	employeeID := seedLiveFinancialEmployee(t, db, fixture)
	seedLivePayrollPostingRule(t, db, fixture)
	salaryID := createLiveFinancialSalary(t, baseURL, client, employeeID)
	var grossCents, netCents int64
	var month, year int
	if err := db.QueryRowContext(ctx, `SELECT ROUND(gross_salary*100)::bigint, ROUND(net_salary*100)::bigint, month, year FROM salaries WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, salaryID, fixture.tenantID, fixture.restaurantID).Scan(&grossCents, &netCents, &month, &year); err != nil || grossCents != 150000 || netCents != 150000 || month != int(time.Now().UTC().Month()) || year != time.Now().UTC().Year() {
		t.Fatalf("synthetic source salary gross/net=%d/%d cents period=%d/%d, expected 150000/150000 and current period (err=%v)", grossCents, netCents, month, year, err)
	}
	markPath := financialAcceptancePath(baseURL, "/api/v1/hr/salaries/%d/mark-paid", salaryID)
	status, body := acceptanceJSONRequest(t, client, http.MethodPost, markPath, map[string]any{"payment_reference": "FA-PAYROLL-1500"})
	if status != http.StatusOK {
		t.Fatalf("authorized payroll mark-paid returned %d: %s", status, body)
	}
	duplicateStatus, duplicateBody := acceptanceJSONRequest(t, client, http.MethodPost, markPath, map[string]any{"payment_reference": "FA-PAYROLL-DUPLICATE"})
	if duplicateStatus != http.StatusOK {
		t.Fatalf("duplicate payroll mark-paid returned %d: %s", duplicateStatus, duplicateBody)
	}
	journalID, journalCount := readLivePayrollJournal(t, db, ctx, fixture, salaryID)
	expenseCents, cashCreditCents, ledgerRows := readLivePayrollLedger(t, db, ctx, fixture, journalID)
	if expenseCents != grossCents || cashCreditCents != netCents || grossCents != netCents || ledgerRows != 2 {
		t.Fatalf("payroll source/GL gross=%d net=%d expense=%d cash-credit=%d rows=%d", grossCents, netCents, expenseCents, cashCreditCents, ledgerRows)
	}
	assertPaidPayrollIsImmutable(t, baseURL, client, salaryID)
	assertForeignPayrollAccessDenied(t, db, baseURL, fixture, salaryID)
	assertLivePayrollPaymentIdentity(t, db, ctx, fixture, salaryID)
	t.Logf("financial_acceptance scenario=authenticated_payroll_api tenant_id=%d restaurant_id=%d salary_id=%d journal_id=%d source_gross_cents=%d source_net_cents=%d expense_gl_cents=%d cash_gl_cents=%d journal_count=%d gl_rows=%d source_gl_difference_cents=0 duplicate_mark_paid_http_status=%d immutable_edit_delete=true foreign_tenant_denied=true", fixture.tenantID, fixture.restaurantID, salaryID, journalID, grossCents, netCents, expenseCents, cashCreditCents, journalCount, ledgerRows, duplicateStatus)
}

func seedLiveFinancialEmployee(t *testing.T, db *sql.DB, fixture liveFinancialFixture) int64 {
	t.Helper()
	key := fmt.Sprintf("FA-EMP-%d", time.Now().UnixNano())
	var id int64
	if err := db.QueryRow(`INSERT INTO employees (tenant_id,restaurant_id,employee_code,first_name,last_name,email,hire_date,employment_type,position) VALUES ($1,$2,$3,'Synthetic','Employee',$4,CURRENT_DATE,'full_time','Acceptance fixture') RETURNING id`, fixture.tenantID, fixture.restaurantID, key, key+"@example.test").Scan(&id); err != nil {
		t.Fatalf("seed synthetic payroll employee: %v", err)
	}
	return id
}

func seedLivePayrollPostingRule(t *testing.T, db *sql.DB, fixture liveFinancialFixture) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO chart_of_accounts (tenant_id,restaurant_id,account_code,account_name,account_type,normal_balance) VALUES ($1,$2,'5210','Synthetic Payroll Expense','expense','debit')`, fixture.tenantID, fixture.restaurantID); err != nil {
		t.Fatalf("seed synthetic payroll expense account: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO auto_posting_rules (tenant_id,restaurant_id,rule_code,rule_name,is_active,default_accounts) VALUES ($1,$2,'SALARY_PAID','Synthetic payroll acceptance',TRUE,'{"bank":"1110","cash":"1110","salary_expense":"5210"}'::jsonb)`, fixture.tenantID, fixture.restaurantID); err != nil {
		t.Fatalf("seed synthetic payroll posting rule: %v", err)
	}
}

func createLiveFinancialSalary(t *testing.T, baseURL string, client *http.Client, employeeID int64) int64 {
	t.Helper()
	start := time.Now().UTC().Format("2006-01-02")
	end := time.Now().UTC().AddDate(0, 1, 0).AddDate(0, 0, -1).Format("2006-01-02")
	status, body := acceptanceJSONRequest(t, client, http.MethodPost, baseURL+"/api/v1/hr/salaries", map[string]any{
		"employee_id": employeeID, "pay_period_start": start, "pay_period_end": end,
		"base_salary": 1500, "currency": "EGP", "payment_method": "bank_transfer", "status": "pending",
	})
	var response struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	if status != http.StatusCreated || json.Unmarshal(body, &response) != nil || response.Data.ID <= 0 {
		t.Fatalf("authenticated salary creation returned %d or no id: %s", status, body)
	}
	return response.Data.ID
}

func readLivePayrollJournal(t *testing.T, db *sql.DB, ctx context.Context, fixture liveFinancialFixture, salaryID int64) (int64, int) {
	t.Helper()
	var count int
	var journalID int64
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*), COALESCE(MIN(id),0) FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND reference_type='salary' AND reference_id=$3 AND entry_type='auto_payroll' AND status='posted'`, fixture.tenantID, fixture.restaurantID, salaryID).Scan(&count, &journalID); err != nil || count != 1 || journalID <= 0 {
		t.Fatalf("payroll posted journal count/id=%d/%d err=%v, want one journal", count, journalID, err)
	}
	return journalID, count
}

func readLivePayrollLedger(t *testing.T, db *sql.DB, ctx context.Context, fixture liveFinancialFixture, journalID int64) (int64, int64, int) {
	t.Helper()
	var expenseCents, cashCreditCents int64
	var rows int
	if err := db.QueryRowContext(ctx, `SELECT COALESCE(ROUND(SUM(gl.debit_amount-gl.credit_amount) FILTER (WHERE a.account_code='5210')*100),0)::bigint, COALESCE(ROUND(SUM(gl.credit_amount-gl.debit_amount) FILTER (WHERE a.account_code='1110')*100),0)::bigint, COUNT(*) FROM general_ledger gl JOIN chart_of_accounts a ON a.id=gl.account_id AND a.tenant_id=gl.tenant_id AND a.restaurant_id=gl.restaurant_id WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND gl.journal_entry_id=$3`, fixture.tenantID, fixture.restaurantID, journalID).Scan(&expenseCents, &cashCreditCents, &rows); err != nil {
		t.Fatalf("independently reconcile payroll GL: %v", err)
	}
	return expenseCents, cashCreditCents, rows
}

func assertPaidPayrollIsImmutable(t *testing.T, baseURL string, client *http.Client, salaryID int64) {
	t.Helper()
	path := financialAcceptancePath(baseURL, "/api/v1/hr/salaries/%d", salaryID)
	status, body := acceptanceJSONRequest(t, client, http.MethodPut, path, map[string]any{"base_salary": 2000})
	if status != http.StatusConflict {
		t.Fatalf("paid salary edit returned %d, want 409: %s", status, body)
	}
	status, body = acceptanceJSONRequest(t, client, http.MethodDelete, path, nil)
	if status != http.StatusConflict {
		t.Fatalf("paid salary delete returned %d, want 409: %s", status, body)
	}
}

func assertForeignPayrollAccessDenied(t *testing.T, db *sql.DB, baseURL string, fixture liveFinancialFixture, salaryID int64) {
	t.Helper()
	runID := fmt.Sprintf("%d", time.Now().UnixNano())
	tenantName := "financial-payroll-foreign-" + runID
	var tenantID, restaurantID int64
	if err := db.QueryRow(`INSERT INTO tenants (name,slug,domain) VALUES ($1,$2,$3) RETURNING id`, tenantName, tenantName, tenantName+".example.test").Scan(&tenantID); err != nil {
		t.Fatalf("seed foreign payroll tenant: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id,name,slug,address) VALUES ($1,$2,$2,'synthetic') RETURNING id`, tenantID, tenantName).Scan(&restaurantID); err != nil {
		t.Fatalf("seed foreign payroll restaurant: %v", err)
	}
	foreign := seedLiveFinancialIdentity(t, db, liveFinancialFixture{tenantID: tenantID, restaurantID: restaurantID}, "owner", "ACCEPTANCE-FOREIGN-OWNER", nil)
	foreignClient := loginLiveFinancialIdentity(t, baseURL, foreign)
	path := financialAcceptancePath(baseURL, "/api/v1/hr/salaries/%d/mark-paid", salaryID)
	status, body := acceptanceJSONRequest(t, foreignClient, http.MethodPost, path, map[string]any{"payment_reference": "FA-FOREIGN-PAYROLL"})
	if status != http.StatusNotFound {
		t.Fatalf("foreign tenant payroll payment returned %d, want 404: %s", status, body)
	}
	var paid bool
	if err := db.QueryRow(`SELECT is_paid FROM salaries WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, salaryID, fixture.tenantID, fixture.restaurantID).Scan(&paid); err != nil || !paid {
		t.Fatalf("foreign payroll attempt changed salary paid=%v err=%v", paid, err)
	}
}

func assertLivePayrollPaymentIdentity(t *testing.T, db *sql.DB, ctx context.Context, fixture liveFinancialFixture, salaryID int64) {
	t.Helper()
	var paid bool
	var paymentReference string
	if err := db.QueryRowContext(ctx, `SELECT is_paid,payment_reference FROM salaries WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, salaryID, fixture.tenantID, fixture.restaurantID).Scan(&paid, &paymentReference); err != nil || !paid || paymentReference != "FA-PAYROLL-1500" {
		t.Fatalf("persisted paid salary is_paid/reference=%v/%s err=%v", paid, paymentReference, err)
	}
}
