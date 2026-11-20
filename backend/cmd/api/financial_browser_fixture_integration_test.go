//go:build integration

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/tests/testdb"
)

type financialBrowserIdentityFixture struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type financialBrowserFixture struct {
	Accountant           financialBrowserIdentityFixture `json:"accountant"`
	Cashier              financialBrowserIdentityFixture `json:"cashier"`
	Payroll              financialBrowserIdentityFixture `json:"payroll"`
	SecondTenant         financialBrowserIdentityFixture `json:"second_tenant"`
	PrimaryAccount       int64                           `json:"primary_account_id"`
	RevenueAccount       int64                           `json:"revenue_account_id"`
	PrimaryTenantID      int64                           `json:"primary_tenant_id"`
	ForeignAccount       int64                           `json:"foreign_account_id"`
	AccountName          string                          `json:"account_name"`
	FormulaName          string                          `json:"formula_name"`
	CustomerName         string                          `json:"customer_name"`
	SupplierName         string                          `json:"supplier_name"`
	EmployeeName         string                          `json:"employee_name"`
	ExpectedBalanceSheet financialBrowserBalanceSheet    `json:"expected_balance_sheet"`
}

type financialBrowserBalanceSheet struct {
	AssetsCents      int64 `json:"assets_cents"`
	LiabilitiesCents int64 `json:"liabilities_cents"`
	EquityCents      int64 `json:"equity_cents"`
}

const financialBrowserJournalAmountCents int64 = 5000

// TestFinancialAcceptanceSeedBrowserFixtures runs only in the isolated CI
// financial acceptance workflow and writes credentials to a private local file.
func TestFinancialAcceptanceSeedBrowserFixtures(t *testing.T) {
	fixturePath := strings.TrimSpace(os.Getenv("POS_FINANCIAL_BROWSER_FIXTURE_FILE"))
	if fixturePath == "" || strings.TrimSpace(os.Getenv("POS_FINANCIAL_ACCEPTANCE_API_URL")) == "" {
		t.Skip("requires the isolated financial browser acceptance workflow")
	}
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable browser acceptance database: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(8)
	ctx, cancel := context.WithTimeout(t.Context(), 45*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("connect browser acceptance database: %v", err)
	}

	primary := seedLiveFinancialOrderFixture(t, db)
	accountant := seedLiveFinancialIdentity(t, db, primary, "accountant", "BROWSER-ACCOUNTANT", []liveFinancialPermission{
		{moduleID: 3, level: "READ"},
		{moduleID: 15, level: "WRITE"},
		{moduleID: 16, level: "READ"},
		{moduleID: 17, level: "READ"},
		{moduleID: 12, level: "READ"},
	})
	cashier := seedLiveFinancialIdentity(t, db, primary, "cashier", "BROWSER-CASHIER", nil)
	// This persona mirrors the dedicated payroll API acceptance role and only
	// receives the approved HR employee read and Salaries delete levels.
	payroll := seedLiveFinancialIdentity(t, db, primary, "payroll-manager", "BROWSER-PAYROLL", []liveFinancialPermission{
		{moduleID: 3, level: "READ"},
		{moduleID: 12, level: "DELETE"},
	})
	primaryAccountID := liveFinancialAccountID(t, db, primary, "1110")
	revenueAccountID := liveFinancialAccountID(t, db, primary, "4100")
	receivableAccountID := liveFinancialAccountID(t, db, primary, "1120")
	payableAccountID := liveFinancialAccountID(t, db, primary, "2101")
	cashAccountID := primaryAccountID
	customerName := "Synthetic Browser AR Customer"
	customerCode := strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("FB-AR-%d", time.Now().UnixNano()%1_000_000_000_000), " ", ""))
	if _, err := db.ExecContext(ctx, `INSERT INTO customers_accounting (tenant_id, restaurant_id, customer_code, customer_name, receivable_account_id, is_active) VALUES ($1,$2,$3,$4,$5,TRUE)`, primary.tenantID, primary.restaurantID, customerCode, customerName, receivableAccountID); err != nil {
		t.Fatalf("seed browser AR customer: %v", err)
	}
	supplierName := "Synthetic Browser AP Supplier"
	supplierCode := fmt.Sprintf("FB-AP-%d", time.Now().UnixNano()%1_000_000_000_000)
	if _, err := db.ExecContext(ctx, `INSERT INTO suppliers (tenant_id, restaurant_id, supplier_code, company_name, payable_account_id, is_active) VALUES ($1,$2,$3,$4,$5,TRUE)`, primary.tenantID, primary.restaurantID, supplierCode, supplierName, payableAccountID); err != nil {
		t.Fatalf("seed browser AP supplier: %v", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO bank_accounts_ledger (tenant_id, restaurant_id, account_name, bank_name, account_number, currency, account_type, opening_date, gl_account_id, is_active, is_primary) VALUES ($1,$2,'Synthetic Browser Cash','Synthetic QA Bank',$3,'EGP','petty_cash',CURRENT_DATE,$4,TRUE,FALSE)`, primary.tenantID, primary.restaurantID, fmt.Sprintf("FB-%d", time.Now().UnixNano()), cashAccountID); err != nil {
		t.Fatalf("seed browser payment account: %v", err)
	}
	seedLiveFinancialEmployee(t, db, primary)
	seedLivePayrollPostingRule(t, db, primary)
	formulaName := `=HYPERLINK("https://example.invalid","Synthetic account")`
	if _, err := db.ExecContext(ctx, `
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
		VALUES ($1, $2, 'QA-FORMULA-1', $3, 'asset', 'debit')`, primary.tenantID, primary.restaurantID, formulaName); err != nil {
		t.Fatalf("seed spreadsheet-formula account fixture: %v", err)
	}

	foreignTenantID, foreignRestaurantID, foreignEmail, foreignPassword := seedLiveFinancialOtherTenant(t, db)
	secondTenant := seedLiveFinancialIdentity(t, db, liveFinancialFixture{
		tenantID: foreignTenantID, restaurantID: foreignRestaurantID,
		email: foreignEmail, password: foreignPassword,
	}, "owner", "BROWSER-SECOND-TENANT", []liveFinancialPermission{
		{moduleID: 15, level: "READ"},
	})
	foreignAccountID := int64(0)
	if err := db.QueryRowContext(ctx, `
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
		VALUES ($1, $2, 'QA-FOREIGN-1110', 'Foreign synthetic browser account', 'asset', 'debit') RETURNING id`,
		foreignTenantID, foreignRestaurantID).Scan(&foreignAccountID); err != nil {
		t.Fatalf("seed foreign tenant browser account: %v", err)
	}
	expectedBalanceSheet := readFinancialBrowserBalanceSheet(t, db, primary, primaryAccountID)

	fixture := financialBrowserFixture{
		Accountant:           financialBrowserIdentityFixture{Email: accountant.email, Password: accountant.password, Role: accountant.role},
		Cashier:              financialBrowserIdentityFixture{Email: cashier.email, Password: cashier.password, Role: cashier.role},
		Payroll:              financialBrowserIdentityFixture{Email: payroll.email, Password: payroll.password, Role: payroll.role},
		SecondTenant:         financialBrowserIdentityFixture{Email: secondTenant.email, Password: secondTenant.password, Role: secondTenant.role},
		PrimaryAccount:       primaryAccountID,
		RevenueAccount:       revenueAccountID,
		PrimaryTenantID:      primary.tenantID,
		ForeignAccount:       foreignAccountID,
		AccountName:          "Synthetic Cash",
		FormulaName:          formulaName,
		CustomerName:         customerName,
		SupplierName:         supplierName,
		EmployeeName:         "Synthetic Employee",
		ExpectedBalanceSheet: expectedBalanceSheet,
	}
	encoded, err := json.MarshalIndent(fixture, "", "  ")
	if err != nil {
		t.Fatalf("encode private browser fixture: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(fixturePath), 0o700); err != nil {
		t.Fatalf("create private fixture directory: %v", err)
	}
	if err := os.WriteFile(fixturePath, encoded, 0o600); err != nil {
		t.Fatalf("write private browser fixture: %v", err)
	}
	t.Log("created isolated accountant, cashier, second-tenant identities and account fixtures")
}

func readFinancialBrowserBalanceSheet(t *testing.T, db *sql.DB, fixture liveFinancialFixture, cashAccountID int64) financialBrowserBalanceSheet {
	t.Helper()
	var expected financialBrowserBalanceSheet
	err := db.QueryRowContext(t.Context(), `
		WITH balances AS (
			SELECT a.account_type, a.normal_balance, SUM(gl.debit_amount-gl.credit_amount) AS net_debit
			FROM general_ledger gl
			JOIN chart_of_accounts a ON a.id=gl.account_id AND a.tenant_id=gl.tenant_id AND a.restaurant_id=gl.restaurant_id
			WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND gl.entry_date<=CURRENT_DATE
			GROUP BY a.account_type,a.normal_balance
		)
		SELECT
			ROUND(COALESCE(SUM(CASE WHEN account_type='asset' THEN CASE WHEN normal_balance='debit' THEN net_debit ELSE -net_debit END ELSE 0 END),0)*100)::bigint,
			ROUND(COALESCE(SUM(CASE WHEN account_type='liability' THEN CASE WHEN normal_balance='credit' THEN -net_debit ELSE net_debit END ELSE 0 END),0)*100)::bigint,
			ROUND(COALESCE(SUM(CASE WHEN account_type='equity' THEN CASE WHEN normal_balance='credit' THEN -net_debit ELSE net_debit END ELSE 0 END),0)*100)::bigint
		FROM balances`, fixture.tenantID, fixture.restaurantID).Scan(&expected.AssetsCents, &expected.LiabilitiesCents, &expected.EquityCents)
	if err != nil {
		t.Fatalf("calculate independent browser Balance Sheet oracle: %v", err)
	}
	var accountType, normalBalance string
	if err := db.QueryRowContext(t.Context(), `SELECT account_type,normal_balance FROM chart_of_accounts WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, cashAccountID, fixture.tenantID, fixture.restaurantID).Scan(&accountType, &normalBalance); err != nil {
		t.Fatalf("verify browser journal debit account: %v", err)
	}
	if accountType != "asset" || normalBalance != "debit" {
		t.Fatalf("browser journal's planned debit must increase a debit-normal asset; account type=%q normal balance=%q", accountType, normalBalance)
	}
	// The browser journey posts this exact cash debit before opening the report.
	expected.AssetsCents += financialBrowserJournalAmountCents
	return expected
}
