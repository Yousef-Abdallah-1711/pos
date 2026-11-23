//go:build integration

package usecase

import (
	"database/sql"
	"fmt"
	"math"
	"testing"

	_ "github.com/lib/pq"
	"pos-saas/tests/testdb"
)

func TestCFR010CalendarBoundaryAndSmallAmounts(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, f, "2025-12-31", "calendar year end", cfr010Line{"cash", 10, 0}, cfr010Line{"capital", 0, 10})
	for i := 0; i < 20; i++ {
		postCFR010Entry(t, db, f, "2026-01-01", fmt.Sprintf("minor unit %d", i), cfr010Line{"cash", 0.01, 0}, cfr010Line{"capital", 0, 0.01})
	}
	for _, tc := range []struct {
		date      string
		wantCents int64
	}{{"2025-12-31", 1000}, {"2026-01-01", 1020}} {
		asOf := cfr010Date(tc.date)
		report := generateCFR010BalanceSheet(t, db, f, asOf)
		var debit, credit int64
		err := db.QueryRow(`SELECT (COALESCE(SUM(debit_amount),0)*100)::bigint, (COALESCE(SUM(credit_amount),0)*100)::bigint FROM general_ledger WHERE tenant_id=$1 AND restaurant_id=$2 AND entry_date <= $3`, f.tenantID, f.restaurantID, asOf).Scan(&debit, &credit)
		if err != nil || debit != tc.wantCents || credit != tc.wantCents || int64(math.Round(report.TotalAssets*100)) != tc.wantCents || int64(math.Round(report.TotalEquity*100)) != tc.wantCents || !report.EquationBalanced {
			t.Fatalf("%s report assets/equity=%.2f/%.2f GL cents=%d/%d, want %d: %v", tc.date, report.TotalAssets, report.TotalEquity, debit, credit, tc.wantCents, err)
		}
		t.Logf("%s report assets/equity and GL debit/credit = %d cents", tc.date, tc.wantCents)
	}
}
