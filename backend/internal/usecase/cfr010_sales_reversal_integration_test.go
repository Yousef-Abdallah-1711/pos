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

// TestCFR010SalesAndReversalAccounting uses a fresh synthetic scope and actual
// journal posting. The residual remains policy-pending for open earnings.
func TestCFR010SalesAndReversalAccounting(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := createCFR010Fixture(t, db)
	date := cfr010Date("2026-04-30")
	postCFR010Entry(t, db, f, "2026-04-01", "sale 1", cfr010Line{"cash", 55, 0}, cfr010Line{"payable", 0, 5}, cfr010Line{"revenue", 0, 50})
	postCFR010Entry(t, db, f, "2026-04-02", "sale 2", cfr010Line{"cash", 22, 0}, cfr010Line{"payable", 0, 2}, cfr010Line{"revenue", 0, 20})
	assertSalesAndReversal(t, db, f, date, 77, 7, 70)
	postCFR010Entry(t, db, f, "2026-04-03", "partial sale reversal", cfr010Line{"revenue", 10, 0}, cfr010Line{"payable", 1, 0}, cfr010Line{"cash", 0, 11})
	assertSalesAndReversal(t, db, f, date, 66, 6, 60)
}

func TestCFR010VoidPreservesHistoricalProfitLoss(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, f, "2025-01-01", "historical sale before void", cfr010Line{"cash", 50, 0}, cfr010Line{"revenue", 0, 50})
	var id int
	if err := db.QueryRow(`SELECT id FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND description=$3`, f.tenantID, f.restaurantID, "historical sale before void").Scan(&id); err != nil {
		t.Fatal(err)
	}
	if err := repository.NewJournalEntryRepository(db).VoidEntry(f.tenantID, f.restaurantID, id, 1, "synthetic reversal"); err != nil {
		t.Fatal(err)
	}
	bs := generateCFR010BalanceSheet(t, db, f, cfr010Date("2025-12-31"))
	pl := generateCFR010ProfitLoss(t, db, f, cfr010Date("2025-01-01"), cfr010Date("2025-12-31"))
	if bs.TotalAssets != 50 || pl.NetProfit != 50 {
		t.Fatalf("historical asset/P&L = %.2f/%.2f, want 50/50 before the later void", bs.TotalAssets, pl.NetProfit)
	}
	today := cfr010Date(time.Now().Format("2006-01-02"))
	reversal := generateCFR010ProfitLoss(t, db, f, today, today)
	combined := generateCFR010ProfitLoss(t, db, f, cfr010Date("2025-01-01"), today)
	current := generateCFR010BalanceSheet(t, db, f, today)
	if reversal.NetProfit != -50 || combined.NetProfit != 0 || current.TotalAssets != 0 {
		t.Fatalf("void-date P&L/combined P&L/current assets = %.2f/%.2f/%.2f, want -50/0/0", reversal.NetProfit, combined.NetProfit, current.TotalAssets)
	}
	var originalCents, reversalCents int64
	err = db.QueryRow(`SELECT
		(COALESCE(SUM(CASE WHEN gl.entry_date < $3 THEN gl.credit_amount-gl.debit_amount ELSE 0 END),0)*100)::bigint,
		(COALESCE(SUM(CASE WHEN gl.entry_date = $3 THEN gl.credit_amount-gl.debit_amount ELSE 0 END),0)*100)::bigint
		FROM general_ledger gl JOIN chart_of_accounts a ON a.id=gl.account_id AND a.tenant_id=gl.tenant_id AND a.restaurant_id=gl.restaurant_id
		WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND a.account_type='revenue'`, f.tenantID, f.restaurantID, today).Scan(&originalCents, &reversalCents)
	if err != nil || originalCents != 5000 || reversalCents != -5000 {
		t.Fatalf("independent GL revenue original/reversal = %d/%d cents, want 5000/-5000: %v", originalCents, reversalCents, err)
	}
	other := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, other, "2025-01-01", "other tenant sale", cfr010Line{"cash", 90, 0}, cfr010Line{"revenue", 0, 90})
	if isolated := generateCFR010ProfitLoss(t, db, f, cfr010Date("2025-01-01"), cfr010Date("2025-12-31")); isolated.NetProfit != 50 {
		t.Fatalf("other tenant contaminated historical P&L: %.2f", isolated.NetProfit)
	}
}

func assertSalesAndReversal(t *testing.T, db *sql.DB, f *cfr010Fixture, asOf time.Time, assets, liabilities, net float64) {
	t.Helper()
	bs := generateCFR010BalanceSheet(t, db, f, asOf)
	pl := generateCFR010ProfitLoss(t, db, f, cfr010Date("2026-04-01"), asOf)
	tb := generateCFR010TrialBalance(t, db, f, asOf)
	if bs.TotalAssets != assets || bs.TotalLiabilities != liabilities || bs.TotalEquity != 0 || bs.EquationDifference != net || bs.EquationBalanced || pl.NetProfit != net || tb.TotalDebits != assets || tb.TotalCredits != assets || !tb.IsBalanced {
		t.Fatalf("report mismatch: balance=%+v P&L=%+v trial=%+v; expected assets/liabilities/net=%.2f/%.2f/%.2f", bs, pl, tb, assets, liabilities, net)
	}
	var actualAssets, actualLiabilities, actualNet int64
	err := db.QueryRow(`SELECT
		(COALESCE(SUM(CASE WHEN a.account_type='asset' THEN gl.debit_amount-gl.credit_amount ELSE 0 END),0)*100)::bigint,
		(COALESCE(SUM(CASE WHEN a.account_type='liability' THEN gl.credit_amount-gl.debit_amount ELSE 0 END),0)*100)::bigint,
		(COALESCE(SUM(CASE WHEN a.account_type='revenue' THEN gl.credit_amount-gl.debit_amount WHEN a.account_type='expense' THEN gl.credit_amount-gl.debit_amount ELSE 0 END),0)*100)::bigint
		FROM general_ledger gl JOIN chart_of_accounts a ON a.id=gl.account_id AND a.tenant_id=gl.tenant_id AND a.restaurant_id=gl.restaurant_id
		WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND gl.entry_date<=$3`, f.tenantID, f.restaurantID, asOf).Scan(&actualAssets, &actualLiabilities, &actualNet)
	if err != nil {
		t.Fatal(err)
	}
	if actualAssets != int64(assets*100) || actualLiabilities != int64(liabilities*100) || actualNet != int64(net*100) {
		t.Fatalf("independent GL oracle cents = %d/%d/%d, want %.0f/%.0f/%.0f", actualAssets, actualLiabilities, actualNet, assets*100, liabilities*100, net*100)
	}
	t.Logf("assets/liabilities/equity/residual/P&L/trial debit/trial credit cents = %d/%d/0/%d/%d/%d/%d; POLICY_PENDING", actualAssets, actualLiabilities, actualNet, actualNet, actualAssets, actualAssets)
}
