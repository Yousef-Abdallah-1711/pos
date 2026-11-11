//go:build integration

package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/tests/testdb"
)

func TestGenerateBalanceSheetPostgresIntegrityStatus(t *testing.T) {
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	defer db.Close()

	asOf := time.Date(2026, time.November, 10, 0, 0, 0, 0, time.UTC)
	journalID := insertBalanceSheetFixture(t, db, asOf)
	defer deleteBalanceSheetFixture(t, db, journalID)
	report, err := NewFinancialReportUseCase(db, nil).GenerateBalanceSheet(1, 1, asOf)
	if err != nil {
		t.Fatalf("GenerateBalanceSheet returned error: %v", err)
	}
	if report.IntegrityStatus != "UNRECONCILED" || report.EquationBalanced {
		t.Fatalf("status/balanced = %s/%t, want UNRECONCILED/false", report.IntegrityStatus, report.EquationBalanced)
	}
	if report.TotalAssets != 100 || report.TotalEquity != 0 || report.EquationDifference != 100 {
		t.Fatalf("report changed source figures: assets=%v equity=%v difference=%v", report.TotalAssets, report.TotalEquity, report.EquationDifference)
	}
}

func insertBalanceSheetFixture(t *testing.T, db *sql.DB, entryDate time.Time) int64 {
	t.Helper()
	ctx := context.Background()
	entryNumber := fmt.Sprintf("P00-%d", time.Now().UnixNano()%1_000_000_000_000)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin fixture transaction: %v", err)
	}
	defer func() { _ = tx.Rollback() }()
	var journalID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO journal_entries (tenant_id, restaurant_id, entry_number, entry_date, description, entry_type, status, posted_at)
		VALUES (1, 1, $1, $2, 'Phase 00 integrity fixture', 'manual', 'posted', NOW()) RETURNING id
	`, entryNumber, entryDate).Scan(&journalID)
	if err != nil {
		t.Fatalf("insert fixture journal: %v", err)
	}
	fixtureLines := []struct {
		code, description      string
		debit, credit, balance float64
	}{
		{code: "1110", description: "Cash on Hand", debit: 100, balance: 100},
		{code: "4100", description: "Food Sales", credit: 100, balance: -100},
	}
	for index, fixtureLine := range fixtureLines {
		var accountID, lineID int64
		if err := tx.QueryRowContext(ctx, `SELECT id FROM chart_of_accounts WHERE tenant_id = $1 AND restaurant_id = $2 AND account_code = $3`, 1, 1, fixtureLine.code).Scan(&accountID); err != nil {
			t.Fatalf("find seeded account %s: %v", fixtureLine.code, err)
		}
		err = tx.QueryRowContext(ctx, `
			INSERT INTO journal_entry_lines (journal_entry_id, tenant_id, account_id, account_code, debit_amount, credit_amount, description, line_order)
			VALUES ($1, 1, $2, $3, $4, $5, $6, $7) RETURNING id
		`, journalID, accountID, fixtureLine.code, fixtureLine.debit, fixtureLine.credit, fixtureLine.description, index+1).Scan(&lineID)
		if err != nil {
			t.Fatalf("insert fixture journal line %s: %v", fixtureLine.code, err)
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO general_ledger (tenant_id, restaurant_id, account_id, journal_entry_id, journal_line_id, entry_date, fiscal_period, debit_amount, credit_amount, running_balance)
			VALUES (1, 1, $1, $2, $3, $4, '2026-11', $5, $6, $7)
		`, accountID, journalID, lineID, entryDate, fixtureLine.debit, fixtureLine.credit, fixtureLine.balance)
		if err != nil {
			t.Fatalf("insert fixture ledger row %s: %v", fixtureLine.code, err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit fixture: %v", err)
	}
	return journalID
}

func deleteBalanceSheetFixture(t *testing.T, db *sql.DB, journalID int64) {
	t.Helper()
	if _, err := db.Exec(`DELETE FROM general_ledger WHERE journal_entry_id = $1`, journalID); err != nil {
		t.Fatalf("delete fixture ledger row: %v", err)
	}
	if _, err := db.Exec(`DELETE FROM journal_entries WHERE id = $1`, journalID); err != nil {
		t.Fatalf("delete fixture journal: %v", err)
	}
}
