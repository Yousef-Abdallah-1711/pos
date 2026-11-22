//go:build integration

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/tests/testdb"
)

func TestDiagnoseOpeningBalances_ReportsPolicyBoundaryAndExactObservation(t *testing.T) {
	db := openOpeningBalanceDiagnosticDB(t)
	fixture := createOpeningBalanceDiagnosticFixture(t, db)
	defer cleanupOpeningBalanceDiagnosticFixture(t, db, fixture)
	insertOpeningBalanceDiagnosticAccount(t, db, fixture, "1000.25", "2025-01-01", &fixture.accountID)

	got, err := NewBankAccountRepository(db).DiagnoseOpeningBalances(context.Background(), fixture.tenantID, fixture.restaurantID)
	if err != nil {
		t.Fatalf("diagnose opening balance: %v", err)
	}
	if got.Status != OpeningBalancePolicyPending || got.Authoritative || got.SourceModel != "bank_accounts_ledger" || len(got.Accounts) != 1 {
		t.Fatalf("opening diagnostic status/source/accounts = %s/%s/%t/%d", got.Status, got.SourceModel, got.Authoritative, len(got.Accounts))
	}
	observation := got.Accounts[0]
	if observation.OpeningBalanceCents != 100025 || observation.GLAccountScope != "scoped" ||
		observation.GLMovementOnOrBeforeDateCents != 0 || observation.GLRowsOnOpeningDate != 0 ||
		observation.SourceReferenceAvailable || observation.ApprovalMetadataAvailable {
		t.Fatalf("opening observation = %+v", observation)
	}
	if len(got.PolicyGaps) != 3 {
		t.Fatalf("policy gaps = %v, want source, approval, and mapping gaps", got.PolicyGaps)
	}
}

func TestDiagnoseOpeningBalances_EmptyScopeDoesNotInferOpeningSource(t *testing.T) {
	db := openOpeningBalanceDiagnosticDB(t)
	fixture := createOpeningBalanceDiagnosticFixture(t, db)
	defer cleanupOpeningBalanceDiagnosticFixture(t, db, fixture)

	got, err := NewBankAccountRepository(db).DiagnoseOpeningBalances(context.Background(), fixture.tenantID, fixture.restaurantID)
	if err != nil {
		t.Fatalf("diagnose empty opening scope: %v", err)
	}
	if got.Status != OpeningBalanceSourceUnavailable || got.Authoritative || len(got.Accounts) != 0 {
		t.Fatalf("empty opening diagnostic = %+v", got)
	}
}

func TestDiagnoseOpeningBalances_RejectsMissingScope(t *testing.T) {
	db := openOpeningBalanceDiagnosticDB(t)
	_, err := NewBankAccountRepository(db).DiagnoseOpeningBalances(context.Background(), 0, 1)
	if err == nil {
		t.Fatal("diagnose opening balances accepted a missing tenant")
	}
}

type openingBalanceDiagnosticFixture struct {
	tenantID, restaurantID, accountID int
	name                              string
}

func openOpeningBalanceDiagnosticDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func createOpeningBalanceDiagnosticFixture(t *testing.T, db *sql.DB) openingBalanceDiagnosticFixture {
	t.Helper()
	fixture := openingBalanceDiagnosticFixture{name: fmt.Sprintf("open-diagnostic-%d", time.Now().UnixNano())}
	err := db.QueryRow(`
		INSERT INTO tenants (name, slug, domain) VALUES ($1, $1, $2) RETURNING id
	`, fixture.name, fixture.name+".example.test").Scan(&fixture.tenantID)
	if err != nil {
		t.Fatalf("create diagnostic tenant: %v", err)
	}
	err = db.QueryRow(`
		INSERT INTO restaurants (tenant_id, name, slug, address)
		VALUES ($1, $2, $2, 'test') RETURNING id
	`, fixture.tenantID, fixture.name).Scan(&fixture.restaurantID)
	if err != nil {
		t.Fatalf("create diagnostic restaurant: %v", err)
	}
	err = db.QueryRow(`
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
		VALUES ($1, $2, '1101', 'Diagnostic cash', 'asset', 'debit') RETURNING id
	`, fixture.tenantID, fixture.restaurantID).Scan(&fixture.accountID)
	if err != nil {
		t.Fatalf("create diagnostic GL account: %v", err)
	}
	return fixture
}

func insertOpeningBalanceDiagnosticAccount(t *testing.T, db *sql.DB, fixture openingBalanceDiagnosticFixture, amount, date string, accountID *int) {
	t.Helper()
	accountNumber := fmt.Sprintf("d%d", time.Now().UnixNano())
	_, err := db.Exec(`
		INSERT INTO bank_accounts_ledger (
			tenant_id, restaurant_id, account_name, bank_name, account_number, currency,
			account_type, opening_balance, current_balance, opening_date, gl_account_id
		) VALUES ($1, $2, 'Synthetic opening account', 'Synthetic bank', $3, 'EGP',
		         'checking', $4, 0, $5, $6)
	`, fixture.tenantID, fixture.restaurantID, accountNumber, amount, date, accountID)
	if err != nil {
		t.Fatalf("create scoped bank opening metadata: %v", err)
	}
}

func cleanupOpeningBalanceDiagnosticFixture(t *testing.T, db *sql.DB, fixture openingBalanceDiagnosticFixture) {
	t.Helper()
	if _, err := db.Exec(`DELETE FROM bank_accounts_ledger WHERE tenant_id = $1 AND restaurant_id = $2`, fixture.tenantID, fixture.restaurantID); err != nil {
		t.Errorf("delete diagnostic bank account: %v", err)
	}
	if _, err := db.Exec(`DELETE FROM chart_of_accounts WHERE tenant_id = $1 AND restaurant_id = $2`, fixture.tenantID, fixture.restaurantID); err != nil {
		t.Errorf("delete diagnostic chart of accounts: %v", err)
	}
	if _, err := db.Exec(`DELETE FROM restaurants WHERE id = $1 AND tenant_id = $2`, fixture.restaurantID, fixture.tenantID); err != nil {
		t.Errorf("delete diagnostic restaurant: %v", err)
	}
	if _, err := db.Exec(`DELETE FROM tenants WHERE id = $1`, fixture.tenantID); err != nil {
		t.Errorf("delete diagnostic tenant: %v", err)
	}
}
