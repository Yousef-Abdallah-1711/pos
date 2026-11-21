//go:build integration

package usecase

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestCFR010CrossReportDiagnosticKeepsUnclosedEarningsPolicyPending(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	defer db.Close()
	fixture := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, fixture, "2025-01-01", "synthetic owner contribution", cfr010Line{"cash", 1000, 0}, cfr010Line{"capital", 0, 1000})
	postCFR010Entry(t, db, fixture, "2025-06-01", "synthetic unclosed earnings", cfr010Line{"cash", 50, 0}, cfr010Line{"revenue", 0, 50})
	otherTenant := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, otherTenant, "2025-06-01", "isolated other-tenant revenue", cfr010Line{"cash", 700, 0}, cfr010Line{"revenue", 0, 700})
	siblingRestaurant := createCFR010RestaurantFixture(t, db, fixture.tenantID)
	var siblingRevenueID int
	if err := db.QueryRow(`
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
		VALUES ($1, $2, '4100', 'Sibling sales', 'revenue', 'credit') RETURNING id
	`, siblingRestaurant.tenantID, siblingRestaurant.restaurantID).Scan(&siblingRevenueID); err != nil {
		t.Fatalf("create sibling revenue account: %v", err)
	}
	siblingRestaurant.accounts["revenue"] = siblingRevenueID
	siblingRestaurant.accountCodes["revenue"] = "4100"
	postCFR010Entry(t, db, siblingRestaurant, "2025-06-01", "isolated sibling revenue", cfr010Line{"cash", 900, 0}, cfr010Line{"revenue", 0, 900})

	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	got, err := repository.NewFinancialIntegrityDiagnosticRepository(db).DiagnoseCrossReportIntegrity(
		t.Context(), fixture.tenantID, fixture.restaurantID, from, to,
	)
	if err != nil {
		t.Fatalf("diagnose cross-report integrity: %v", err)
	}
	if got.Status != repository.CrossReportPolicyPending || got.ProfitLossStatus != repository.CrossReportReconciled ||
		got.TrialBalanceStatus != repository.CrossReportGLBalanced || got.BalanceSheetStatus != repository.CrossReportPolicyPending ||
		got.JournalRevenueCents != 5000 || got.LedgerRevenueCents != 5000 || got.BalanceSheetResidualCents != 5000 ||
		got.TrialBalanceDebitsCents != 105000 || got.TrialBalanceCreditsCents != 105000 ||
		len(got.ProfitLossAccountMismatches) != 0 {
		t.Fatalf("cross-report diagnostic = %+v", got)
	}
}
