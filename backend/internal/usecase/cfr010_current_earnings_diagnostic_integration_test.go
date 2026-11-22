//go:build integration

package usecase

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"pos-saas/internal/repository"
	"pos-saas/tests/testdb"
)

func TestCFR010CurrentEarningsDiagnosticReportsExactCandidateAndPolicyRisk(t *testing.T) {
	db := openCFR010EarningsDiagnosticDB(t)
	fixture := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, fixture, "2025-06-01", "synthetic earnings revenue", cfr010Line{"cash", 50, 0}, cfr010Line{"revenue", 0, 50})
	postCFR010Entry(t, db, fixture, "2025-06-30", "synthetic earnings expense", cfr010Line{"expense", 15, 0}, cfr010Line{"cash", 0, 15})
	postCFR010Entry(t, db, fixture, "2025-06-30", "unclassified retained-earnings movement", cfr010Line{"cash", 10, 0}, cfr010Line{"retained", 0, 10})

	got := diagnoseCurrentEarnings(t, db, fixture, "2025-06-01", "2025-06-30")
	if got.Status != CurrentEarningsPolicyPending || got.FiscalPeriodState != "SOURCE_UNAVAILABLE" ||
		got.RevenueCents != 5000 || got.ExpenseCents != 1500 || got.CandidateNetEarningsCents != 3500 ||
		got.RetainedEarningsMovementCents != 1000 || got.RetainedEarningsJournalCount != 1 ||
		!got.PotentialDoubleCountingRisk || got.CloseJournalIdentification != "UNAVAILABLE" {
		t.Fatalf("current-earnings diagnostic = %+v", got)
	}
}

func TestCFR010CurrentEarningsDiagnosticReportsKnownFiveThousandCentScenario(t *testing.T) {
	db := openCFR010EarningsDiagnosticDB(t)
	fixture := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, fixture, "2025-06-15", "synthetic known earnings", cfr010Line{"cash", 50, 0}, cfr010Line{"revenue", 0, 50})

	got := diagnoseCurrentEarnings(t, db, fixture, "2025-06-01", "2025-06-30")
	if got.Status != CurrentEarningsPolicyPending || got.RevenueCents != 5000 || got.ExpenseCents != 0 ||
		got.CandidateNetEarningsCents != 5000 || got.RetainedEarningsMovementCents != 0 || got.PotentialDoubleCountingRisk {
		t.Fatalf("known five-thousand-cent candidate = %+v", got)
	}
}

func TestCFR010CurrentEarningsDiagnosticIncludesPeriodEndAndReversalEffects(t *testing.T) {
	db := openCFR010EarningsDiagnosticDB(t)
	fixture := createCFR010Fixture(t, db)
	postCFR010Entry(t, db, fixture, "2025-12-31", "synthetic boundary revenue", cfr010Line{"cash", 25, 0}, cfr010Line{"revenue", 0, 25})
	postCFR010Entry(t, db, fixture, "2026-01-01", "synthetic next-period loss", cfr010Line{"expense", 40, 0}, cfr010Line{"cash", 0, 40})
	var revenueJournalID int
	if err := db.QueryRow(`
		SELECT id FROM journal_entries
		WHERE tenant_id = $1 AND restaurant_id = $2 AND description = 'synthetic boundary revenue'
	`, fixture.tenantID, fixture.restaurantID).Scan(&revenueJournalID); err != nil {
		t.Fatalf("read synthetic boundary journal: %v", err)
	}
	if err := repository.NewJournalEntryRepository(db).VoidEntry(
		fixture.tenantID, fixture.restaurantID, revenueJournalID, 1, "synthetic earnings diagnostic reversal",
	); err != nil {
		t.Fatalf("void synthetic boundary journal: %v", err)
	}

	priorPeriod := diagnoseCurrentEarnings(t, db, fixture, "2025-01-01", "2025-12-31")
	if priorPeriod.RevenueCents != 2500 || priorPeriod.ExpenseCents != 0 || priorPeriod.CandidateNetEarningsCents != 2500 || priorPeriod.CrossPeriodVoidCount != 1 {
		t.Fatalf("pre-reversal-period result = %+v", priorPeriod)
	}
	combinedPeriods := diagnoseCurrentEarnings(t, db, fixture, "2025-01-01", "2026-12-31")
	if combinedPeriods.RevenueCents != 0 || combinedPeriods.ExpenseCents != 4000 || combinedPeriods.CandidateNetEarningsCents != -4000 {
		t.Fatalf("combined-period reversal and loss result = %+v", combinedPeriods)
	}
}

func TestCFR010CurrentEarningsDiagnosticReportsClosedEmptyPeriod(t *testing.T) {
	db := openCFR010EarningsDiagnosticDB(t)
	fixture := createCFR010Fixture(t, db)
	if _, err := db.Exec(`
		INSERT INTO fiscal_periods (tenant_id, restaurant_id, period_name, period_type, start_date, end_date, status)
		VALUES ($1, $2, '2027', 'year', '2027-01-01', '2027-12-31', 'closed')
	`, fixture.tenantID, fixture.restaurantID); err != nil {
		t.Fatalf("create closed synthetic fiscal period: %v", err)
	}
	got := diagnoseCurrentEarnings(t, db, fixture, "2027-01-01", "2027-12-31")
	if got.Status != CurrentEarningsNoActivity || got.FiscalPeriodState != "CLOSED_PERIOD_OBSERVED" ||
		got.RevenueCents != 0 || got.ExpenseCents != 0 || got.CandidateNetEarningsCents != 0 {
		t.Fatalf("closed empty-period result = %+v", got)
	}
}

func openCFR010EarningsDiagnosticDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func diagnoseCurrentEarnings(t *testing.T, db *sql.DB, fixture *cfr010Fixture, start, end string) CurrentEarningsDiagnostic {
	t.Helper()
	from, to := cfr010Date(start), cfr010Date(end)
	got, err := NewFinancialReportUseCase(db, nil).DiagnosePotentialCurrentEarnings(t.Context(), fixture.tenantID, fixture.restaurantID, from, to)
	if err != nil {
		t.Fatalf("diagnose current earnings: %v", err)
	}
	return got
}
