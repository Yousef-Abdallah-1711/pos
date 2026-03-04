package repository

import (
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestGeneralLedgerRepository_GetLedgerByAccountCursor_ReturnsNextCursor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewGeneralLedgerRepository(db, nil)
	start := time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.February, 28, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`FROM general_ledger gl`).
		WithArgs(1, 2, 1101, start, end, 3).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "account_id", "account_code", "account_name", "journal_entry_id", "journal_line_id",
			"entry_date", "entry_number", "description", "reference_type", "reference_number",
			"debit_amount", "credit_amount", "running_balance", "line_description",
		}).
			AddRow(101, 1101, "1110", "Cash", 9001, 11, end, "JE-1", "Line 1", "", "", 100.0, 0.0, 100.0, "").
			AddRow(100, 1101, "1110", "Cash", 9000, 10, end.AddDate(0, 0, -1), "JE-0", "Line 0", "", "", 0.0, 20.0, 80.0, "").
			AddRow(99, 1101, "1110", "Cash", 8999, 9, end.AddDate(0, 0, -2), "JE-X", "Line X", "", "", 0.0, 10.0, 70.0, ""))

	lines, nextCursor, err := repo.GetLedgerByAccountCursor(1, 2, 1101, &start, &end, 2, nil)
	if err != nil {
		t.Fatalf("GetLedgerByAccountCursor returned error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if nextCursor == nil || *nextCursor != 100 {
		t.Fatalf("expected next cursor 100, got %v", nextCursor)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestGeneralLedgerRepository_GetLedgerPeriodSummaries_ReturnsAggregates(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}
	defer db.Close()

	repo := NewGeneralLedgerRepository(db, nil)

	mock.ExpectQuery(`WITH filtered AS`).
		WithArgs(1, 2, 1101).
		WillReturnRows(sqlmock.NewRows([]string{
			"fiscal_period", "entry_count", "debit_total", "credit_total", "net_change", "running_balance",
		}).
			AddRow("2026-02", 20, 2500.0, 2000.0, 500.0, 1500.0).
			AddRow("2026-01", 10, 1200.0, 900.0, 300.0, 1000.0))

	rows, err := repo.GetLedgerPeriodSummaries(1, 2, 1101, nil, nil)
	if err != nil {
		t.Fatalf("GetLedgerPeriodSummaries returned error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 summary rows, got %d", len(rows))
	}
	if rows[0].FiscalPeriod != "2026-02" || rows[0].EntryCount != 20 {
		t.Fatalf("unexpected first summary row: %+v", rows[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
