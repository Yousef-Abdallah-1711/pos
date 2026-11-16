//go:build integration

package main

import (
	"context"
	"database/sql"
	"encoding/json"
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
	Accountant      financialBrowserIdentityFixture `json:"accountant"`
	Cashier         financialBrowserIdentityFixture `json:"cashier"`
	SecondTenant    financialBrowserIdentityFixture `json:"second_tenant"`
	PrimaryAccount  int64                           `json:"primary_account_id"`
	RevenueAccount  int64                           `json:"revenue_account_id"`
	PrimaryTenantID int64                           `json:"primary_tenant_id"`
	ForeignAccount  int64                           `json:"foreign_account_id"`
	AccountName     string                          `json:"account_name"`
	FormulaName     string                          `json:"formula_name"`
}

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
		{moduleID: 15, level: "WRITE"},
		{moduleID: 17, level: "READ"},
	})
	cashier := seedLiveFinancialIdentity(t, db, primary, "cashier", "BROWSER-CASHIER", nil)
	primaryAccountID := liveFinancialAccountID(t, db, primary, "1110")
	revenueAccountID := liveFinancialAccountID(t, db, primary, "4100")
	formulaName := `=HYPERLINK("https://example.invalid","Synthetic account")`
	if _, err := db.ExecContext(ctx, `
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
		VALUES ($1, $2, 'QA-FORMULA-1', $3, 'asset', 'debit')`, primary.tenantID, primary.restaurantID, formulaName); err != nil {
		t.Fatalf("seed spreadsheet-formula account fixture: %v", err)
	}

	foreignTenantID, foreignRestaurantID, foreignEmail, foreignPassword := seedLiveFinancialOtherTenant(t, db)
	foreignAccountID := int64(0)
	if err := db.QueryRowContext(ctx, `
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
		VALUES ($1, $2, 'QA-FOREIGN-1110', 'Foreign synthetic browser account', 'asset', 'debit') RETURNING id`,
		foreignTenantID, foreignRestaurantID).Scan(&foreignAccountID); err != nil {
		t.Fatalf("seed foreign tenant browser account: %v", err)
	}

	fixture := financialBrowserFixture{
		Accountant:      financialBrowserIdentityFixture{Email: accountant.email, Password: accountant.password, Role: accountant.role},
		Cashier:         financialBrowserIdentityFixture{Email: cashier.email, Password: cashier.password, Role: cashier.role},
		SecondTenant:    financialBrowserIdentityFixture{Email: foreignEmail, Password: foreignPassword, Role: "owner"},
		PrimaryAccount:  primaryAccountID,
		RevenueAccount:  revenueAccountID,
		PrimaryTenantID: primary.tenantID,
		ForeignAccount:  foreignAccountID,
		AccountName:     "Synthetic Cash",
		FormulaName:     formulaName,
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
