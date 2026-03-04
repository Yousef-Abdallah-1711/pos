package usecase

import (
	"math"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/repository"
)

func TestFinancialReportUseCase_GenerateProfitLossWithComparison_ManualCalculationMatches(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	uc := NewFinancialReportUseCase(db, nil)
	currentFrom := time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)
	currentTo := time.Date(2026, time.March, 31, 0, 0, 0, 0, time.UTC)
	compareFrom := time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)
	compareTo := time.Date(2026, time.February, 28, 0, 0, 0, 0, time.UTC)

	columns := []string{"account_code", "account_name", "account_type", "debit_total", "credit_total"}
	mock.ExpectQuery(`FROM journal_entry_lines jel`).
		WithArgs(1, 2, currentFrom, currentTo).
		WillReturnRows(
			sqlmock.NewRows(columns).
				AddRow("4100", "Food Sales", "revenue", 0.0, 1000.0).
				AddRow("5110", "Food COGS", "expense", 350.0, 0.0),
		)
	mock.ExpectQuery(`FROM journal_entry_lines jel`).
		WithArgs(1, 2, compareFrom, compareTo).
		WillReturnRows(
			sqlmock.NewRows(columns).
				AddRow("4100", "Food Sales", "revenue", 0.0, 800.0).
				AddRow("5110", "Food COGS", "expense", 300.0, 0.0),
		)

	report, err := uc.GenerateProfitLossWithComparison(1, 2, currentFrom, currentTo, compareFrom, compareTo)
	if err != nil {
		t.Fatalf("GenerateProfitLossWithComparison returned error: %v", err)
	}

	if report.DateFrom != "2026-03-01" || report.DateTo != "2026-03-31" {
		t.Fatalf("unexpected current range: %s to %s", report.DateFrom, report.DateTo)
	}
	if report.ComparisonDateFrom != "2026-02-01" || report.ComparisonDateTo != "2026-02-28" {
		t.Fatalf("unexpected comparison range: %s to %s", report.ComparisonDateFrom, report.ComparisonDateTo)
	}
	if !almostEqual(report.TotalRevenue, 1000.0) {
		t.Fatalf("expected total revenue 1000.0, got %f", report.TotalRevenue)
	}
	if !almostEqual(report.TotalExpenses, 350.0) {
		t.Fatalf("expected total expenses 350.0, got %f", report.TotalExpenses)
	}
	if !almostEqual(report.NetProfit, 650.0) {
		t.Fatalf("expected net profit 650.0, got %f", report.NetProfit)
	}
	if !almostEqual(report.ComparisonTotalRevenue, 800.0) {
		t.Fatalf("expected comparison revenue 800.0, got %f", report.ComparisonTotalRevenue)
	}
	if !almostEqual(report.ComparisonTotalExpenses, 300.0) {
		t.Fatalf("expected comparison expenses 300.0, got %f", report.ComparisonTotalExpenses)
	}
	if !almostEqual(report.ComparisonNetProfit, 500.0) {
		t.Fatalf("expected comparison net profit 500.0, got %f", report.ComparisonNetProfit)
	}
	if !almostEqual(report.RevenueChangePercent, 25.0) {
		t.Fatalf("expected revenue change 25.0%%, got %f", report.RevenueChangePercent)
	}
	if !almostEqual(report.ExpenseChangePercent, 16.6666667) {
		t.Fatalf("expected expense change ~16.67%%, got %f", report.ExpenseChangePercent)
	}
	if !almostEqual(report.NetProfitChangePercent, 30.0) {
		t.Fatalf("expected net profit change 30.0%%, got %f", report.NetProfitChangePercent)
	}

	if len(report.Revenue) != 1 || len(report.Expenses) != 1 {
		t.Fatalf("expected 1 revenue and 1 expense line, got %d/%d", len(report.Revenue), len(report.Expenses))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestFinancialReportUseCase_GenerateTrialBalance_ManualCalculationMatches(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	journalRepo := repository.NewJournalEntryRepository(db)
	generalLedgerRepo := repository.NewGeneralLedgerRepository(db, journalRepo)
	uc := NewFinancialReportUseCase(db, generalLedgerRepo)

	asOf := time.Date(2026, time.March, 31, 0, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`WITH latest_balances AS`).
		WithArgs(1, 2, asOf).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "account_code", "account_name", "account_type", "normal_balance", "net_balance",
			}).
				AddRow(1, "1110", "Cash on Hand", "asset", "debit", 200.0).
				AddRow(2, "4100", "Food Sales", "revenue", "credit", 200.0),
		)

	report, err := uc.GenerateTrialBalance(1, 2, asOf)
	if err != nil {
		t.Fatalf("GenerateTrialBalance returned error: %v", err)
	}

	if report.AsOfDate != "2026-03-31" {
		t.Fatalf("expected as_of_date 2026-03-31, got %s", report.AsOfDate)
	}
	if len(report.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(report.Rows))
	}
	if !almostEqual(report.TotalDebits, 200.0) {
		t.Fatalf("expected total debits 200.0, got %f", report.TotalDebits)
	}
	if !almostEqual(report.TotalCredits, 200.0) {
		t.Fatalf("expected total credits 200.0, got %f", report.TotalCredits)
	}
	if !report.IsBalanced {
		t.Fatalf("expected trial balance to be balanced")
	}
	if !almostEqual(report.Rows[0].DebitBalance, 200.0) || !almostEqual(report.Rows[0].CreditBalance, 0.0) {
		t.Fatalf("unexpected first row balances: debit=%f credit=%f", report.Rows[0].DebitBalance, report.Rows[0].CreditBalance)
	}
	if !almostEqual(report.Rows[1].DebitBalance, 0.0) || !almostEqual(report.Rows[1].CreditBalance, 200.0) {
		t.Fatalf("unexpected second row balances: debit=%f credit=%f", report.Rows[1].DebitBalance, report.Rows[1].CreditBalance)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.0001
}
