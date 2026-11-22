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

	profitLoss := generateCFR010ProfitLoss(t, db, fixture, from, to)
	trialBalance := generateCFR010TrialBalance(t, db, fixture, to)
	balanceSheet := generateCFR010BalanceSheet(t, db, fixture, to)
	if int64(profitLoss.TotalRevenue*100) != got.LedgerRevenueCents ||
		int64(profitLoss.TotalExpenses*100) != got.LedgerExpenseCents ||
		int64(trialBalance.TotalDebits*100) != got.TrialBalanceDebitsCents ||
		int64(trialBalance.TotalCredits*100) != got.TrialBalanceCreditsCents ||
		int64(balanceSheet.TotalAssets*100) != got.BalanceSheetAssetsCents ||
		int64(balanceSheet.TotalLiabilities*100) != got.BalanceSheetLiabilitiesCents ||
		int64(balanceSheet.TotalEquity*100) != got.BalanceSheetEquityCents {
		t.Fatalf("independent SQL oracles do not match reports: diagnostic=%+v P&L=%+v TB=%+v BS=%+v", got, profitLoss, trialBalance, balanceSheet)
	}
}

func TestCFR010CrossReportDiagnosticDatesVoidReversalAtVoidDate(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	defer db.Close()

	fixture := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, fixture, "2025-06-01", "synthetic cross-period revenue", cfr010Line{"cash", 50, 0}, cfr010Line{"revenue", 0, 50})
	var journalID int
	if err := db.QueryRow(`
		SELECT id FROM journal_entries
		WHERE tenant_id = $1 AND restaurant_id = $2 AND description = $3
	`, fixture.tenantID, fixture.restaurantID, "synthetic cross-period revenue").Scan(&journalID); err != nil {
		t.Fatalf("find synthetic posted journal: %v", err)
	}
	journalRepository := repository.NewJournalEntryRepository(db)
	if err := journalRepository.VoidEntry(fixture.tenantID, fixture.restaurantID, journalID, 1, "synthetic cross-period reversal"); err != nil {
		t.Fatalf("void synthetic posted journal: %v", err)
	}
	var voidDate time.Time
	if err := db.QueryRow(`SELECT voided_at::date FROM journal_entries WHERE id = $1`, journalID).Scan(&voidDate); err != nil {
		t.Fatalf("read synthetic reversal date: %v", err)
	}
	got, err := repository.NewFinancialIntegrityDiagnosticRepository(db).DiagnoseCrossReportIntegrity(
		t.Context(), fixture.tenantID, fixture.restaurantID, voidDate, voidDate,
	)
	if err != nil {
		t.Fatalf("diagnose reversal-date P&L scope: %v", err)
	}
	if got.CrossPeriodVoidCount != 1 || got.ProfitLossStatus != repository.CrossReportReconciled ||
		got.Status != repository.CrossReportPolicyPending || got.JournalRevenueCents != -5000 ||
		got.LedgerRevenueCents != -5000 || len(got.ProfitLossAccountMismatches) != 0 {
		t.Fatalf("dated cross-period reversal should reconcile while policy remains pending: %+v", got)
	}
}
