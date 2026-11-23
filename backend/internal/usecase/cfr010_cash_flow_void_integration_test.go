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

// This tests the existing report formula's movement dates, not approval of
// that formula as a cash flow methodology. Cash and tax are checked separately.
func TestCFR010CashFlowVoidMovementPeriods(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, f, "2025-01-01", "cash sale before void", cfr010Line{"cash", 55, 0}, cfr010Line{"payable", 0, 5}, cfr010Line{"revenue", 0, 50})
	var id int
	if err := db.QueryRow(`SELECT id FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND description=$3`, f.tenantID, f.restaurantID, "cash sale before void").Scan(&id); err != nil {
		t.Fatal(err)
	}
	uc := NewFinancialReportUseCase(db, nil)
	before, err := uc.GenerateCashFlow(f.tenantID, f.restaurantID, cfr010Date("2025-01-01"), cfr010Date("2025-12-31"))
	if err != nil || before.OperatingCashFlow != 50 {
		t.Fatalf("pre-void existing operating formula = %.2f, want 50: %v", before.OperatingCashFlow, err)
	}
	if err := repository.NewJournalEntryRepository(db).VoidEntry(f.tenantID, f.restaurantID, id, 1, "synthetic reversal"); err != nil {
		t.Fatal(err)
	}
	if err := repository.NewJournalEntryRepository(db).VoidEntry(f.tenantID, f.restaurantID, id, 1, "duplicate"); err == nil {
		t.Fatal("duplicate void unexpectedly succeeded")
	}
	var voidDate time.Time
	if err := db.QueryRow(`SELECT voided_at::date FROM journal_entries WHERE id=$1 AND tenant_id=$2 AND restaurant_id=$3`, id, f.tenantID, f.restaurantID).Scan(&voidDate); err != nil {
		t.Fatal(err)
	}
	var originalCash, reversedCash int64
	err = db.QueryRow(`SELECT (COALESCE(SUM(CASE WHEN gl.entry_date < $4 THEN gl.debit_amount-gl.credit_amount ELSE 0 END),0)*100)::bigint,
		(COALESCE(SUM(CASE WHEN gl.entry_date = $4 THEN gl.debit_amount-gl.credit_amount ELSE 0 END),0)*100)::bigint
		FROM general_ledger gl WHERE gl.tenant_id=$1 AND gl.restaurant_id=$2 AND gl.account_id=$3`, f.tenantID, f.restaurantID, f.accounts["cash"], voidDate).Scan(&originalCash, &reversedCash)
	if err != nil || originalCash != 5500 || reversedCash != -5500 {
		t.Fatalf("independent cash GL original/reversal = %d/%d cents, want 5500/-5500: %v", originalCash, reversedCash, err)
	}
	for _, tc := range []struct {
		name     string
		from, to time.Time
		want     float64
	}{
		{"original", cfr010Date("2025-01-01"), cfr010Date("2025-12-31"), 50},
		{"reversal", voidDate, voidDate, -50},
		{"combined", cfr010Date("2025-01-01"), voidDate, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			report, err := uc.GenerateCashFlow(f.tenantID, f.restaurantID, tc.from, tc.to)
			if err != nil || report.OperatingCashFlow != tc.want {
				t.Fatalf("existing operating formula = %.2f, want %.2f: %v", report.OperatingCashFlow, tc.want, err)
			}
		})
	}
	other := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, other, "2025-01-01", "other tenant", cfr010Line{"cash", 90, 0}, cfr010Line{"revenue", 0, 90})
	sibling := createCFR010RestaurantFixture(t, db, f.tenantID)
	postCFR010Entry(t, db, sibling, "2025-01-01", "other restaurant", cfr010Line{"cash", 80, 0}, cfr010Line{"capital", 0, 80})
	isolated, err := uc.GenerateCashFlow(f.tenantID, f.restaurantID, cfr010Date("2025-01-01"), cfr010Date("2025-12-31"))
	if err != nil || isolated.OperatingCashFlow != 50 || isolated.FinancingCashFlow != 0 {
		t.Fatalf("cash flow isolation = %+v, want operating 50 and financing 0: %v", isolated, err)
	}
}

func TestCFR010CashFlowSameDayVoid(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := createCFR010Fixture(t, db)
	var date time.Time
	if err := db.QueryRow(`SELECT CURRENT_DATE`).Scan(&date); err != nil {
		t.Fatal(err)
	}
	postCFR010Entry(t, db, f, date.Format("2006-01-02"), "same-day void", cfr010Line{"cash", 25, 0}, cfr010Line{"revenue", 0, 25})
	var id int
	if err := db.QueryRow(`SELECT id FROM journal_entries WHERE tenant_id=$1 AND restaurant_id=$2 AND description=$3`, f.tenantID, f.restaurantID, "same-day void").Scan(&id); err != nil {
		t.Fatal(err)
	}
	if err := repository.NewJournalEntryRepository(db).VoidEntry(f.tenantID, f.restaurantID, id, 1, "same day"); err != nil {
		t.Fatal(err)
	}
	report, err := NewFinancialReportUseCase(db, nil).GenerateCashFlow(f.tenantID, f.restaurantID, date, date)
	if err != nil || report.NetCashFlow != 0 {
		t.Fatalf("same-day original and reversal = %.2f, want 0: %v", report.NetCashFlow, err)
	}
}
