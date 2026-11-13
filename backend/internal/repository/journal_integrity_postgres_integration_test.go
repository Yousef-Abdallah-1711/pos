//go:build integration

package repository

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCFR002PostgresRejectsUnbalancedDirectSQLPosting(t *testing.T) {
	for _, test := range []struct {
		name   string
		debit  float64
		credit float64
		lines  bool
	}{
		{name: "debit heavy", debit: 100, credit: 90, lines: true},
		{name: "credit heavy", debit: 90, credit: 100, lines: true},
		{name: "zero lines"},
	} {
		t.Run(test.name, func(t *testing.T) {
			db, fixture := openAtomicPostingFixture(t)
			defer cleanupAtomicPostingFixture(t, db, fixture)

			entryID := createCFR002DraftSQL(t, db, fixture, test.debit, test.credit, test.lines)
			if _, err := db.Exec(`UPDATE journal_entries SET status = 'posted', posted_at = CURRENT_TIMESTAMP, posted_by = 1 WHERE id = $1`, entryID); err == nil {
				t.Fatal("direct SQL posted an invalid or empty journal")
			}
			var status string
			if err := db.QueryRow(`SELECT status FROM journal_entries WHERE id = $1`, entryID).Scan(&status); err != nil {
				t.Fatalf("read journal after rejected posting: %v", err)
			}
			if status != "draft" {
				t.Fatalf("journal status = %q after rejected posting, want draft", status)
			}
		})
	}
}

func TestCFR002PostgresProtectsPostedJournalAndLedgerWithDirectSQL(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	entryID := createCFR002DraftSQL(t, db, fixture, 100, 100, true)
	if _, err := db.Exec(`INSERT INTO journal_entry_lines (journal_entry_id, tenant_id, account_id, account_code, debit_amount, line_order) VALUES ($1, $2, $3, 'third-line', 50, 3)`, entryID, fixture.tenantID, fixture.debitAccountID); err != nil {
		t.Fatalf("add third balanced line: %v", err)
	}
	if _, err := db.Exec(`UPDATE journal_entry_lines SET debit_amount = 50 WHERE journal_entry_id = $1 AND line_order = 1`, entryID); err != nil {
		t.Fatalf("balance three-line journal: %v", err)
	}
	lineIDs := cfr002LineIDs(t, db, entryID)
	postCFR002DirectSQL(t, db, fixture, entryID, lineIDs)

	mutations := []struct {
		name  string
		query string
		args  []any
	}{
		{name: "line update", query: `UPDATE journal_entry_lines SET debit_amount = 101 WHERE id = $1`, args: []any{lineIDs[0]}},
		{name: "line delete", query: `DELETE FROM journal_entry_lines WHERE id = $1`, args: []any{lineIDs[0]}},
		{name: "line insert", query: `INSERT INTO journal_entry_lines (journal_entry_id, tenant_id, account_id, account_code, debit_amount, line_order) VALUES ($1, $2, $3, 'blocked', 1, 3)`, args: []any{entryID, fixture.tenantID, fixture.debitAccountID}},
		{name: "posted header amount", query: `UPDATE journal_entries SET total_debit = total_debit + 1 WHERE id = $1`, args: []any{entryID}},
		{name: "posted to draft", query: `UPDATE journal_entries SET status = 'draft' WHERE id = $1`, args: []any{entryID}},
		{name: "ledger update", query: `UPDATE general_ledger SET debit_amount = debit_amount + 1 WHERE journal_entry_id = $1`, args: []any{entryID}},
		{name: "ledger delete", query: `DELETE FROM general_ledger WHERE journal_entry_id = $1`, args: []any{entryID}},
		{name: "duplicate ledger effect", query: `INSERT INTO general_ledger (tenant_id, restaurant_id, account_id, journal_entry_id, journal_line_id, entry_date, fiscal_period, debit_amount, credit_amount, running_balance) SELECT tenant_id, restaurant_id, account_id, journal_entry_id, journal_line_id, entry_date, fiscal_period, debit_amount, credit_amount, running_balance FROM general_ledger WHERE journal_entry_id = $1 LIMIT 1`, args: []any{entryID}},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			if _, err := db.Exec(mutation.query, mutation.args...); err == nil {
				t.Fatalf("%s unexpectedly changed posted accounting history", mutation.name)
			}
		})
	}
	if _, err := db.Exec(`UPDATE chart_of_accounts SET restaurant_id = restaurant_id + 100000 WHERE id = $1`, fixture.debitAccountID); err == nil {
		t.Fatal("account scope move unexpectedly changed posted accounting scope")
	}
	otherTenant := openAtomicPostingTenantFixture(t, db)
	if _, err := db.Exec(`UPDATE restaurants SET tenant_id = $1 WHERE id = $2`, otherTenant.tenantID, fixture.restaurantID); err == nil {
		t.Fatal("restaurant tenant move unexpectedly changed posted accounting scope")
	}
	if _, err := db.Exec(`UPDATE general_ledger SET debit_amount = debit_amount + 1 WHERE journal_entry_id = $1`, entryID); err == nil {
		t.Fatal("second ledger mutation unexpectedly succeeded")
	}
}

func TestCFR002PostgresConcurrentPostingSerializesAndVoidAppendsInverseEffects(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	journalRepo := NewJournalEntryRepository(db)
	entry := fixture.entry()
	entry.ReferenceID = nil
	entry.ReferenceType = ""
	entry.ReferenceNumber = "CFR002-concurrent"
	entryID, err := journalRepo.CreateJournalEntry(entry)
	if err != nil {
		t.Fatalf("create draft journal: %v", err)
	}

	const attempts = 6
	start := make(chan struct{})
	results := make(chan error, attempts)
	var workers sync.WaitGroup
	for range attempts {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			results <- journalRepo.PostEntry(fixture.tenantID, fixture.restaurantID, entryID, 1)
		}()
	}
	close(start)
	workers.Wait()
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("successful concurrent postings = %d, want exactly one", successes)
	}
	assertJournalLedgerEffects(t, db, entryID, 1)

	if err := journalRepo.VoidEntry(fixture.tenantID, fixture.restaurantID, entryID, 2, "CFR-002 integration reversal"); err != nil {
		t.Fatalf("void posted journal: %v", err)
	}
	assertJournalLedgerEffects(t, db, entryID, 2)
	if err := journalRepo.VoidEntry(fixture.tenantID, fixture.restaurantID, entryID, 2, "duplicate"); err == nil {
		t.Fatal("duplicate void unexpectedly succeeded")
	}
	var status string
	if err := db.QueryRow(`SELECT status FROM journal_entries WHERE id = $1`, entryID).Scan(&status); err != nil {
		t.Fatalf("read voided journal status: %v", err)
	}
	if status != "void" {
		t.Fatalf("journal status = %q, want void", status)
	}
}

func TestCFR002PostgresRejectsCrossTenantJournalLine(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)
	otherTenant := openAtomicPostingTenantFixture(t, db)

	entryID := createCFR002DraftSQL(t, db, fixture, 100, 100, true)
	_, err := db.Exec(`
		INSERT INTO journal_entry_lines (journal_entry_id, tenant_id, account_id, account_code, debit_amount, line_order)
		VALUES ($1, $2, $3, 'cross-tenant', 1, 3)
	`, entryID, fixture.tenantID, otherTenant.debitAccountID)
	if err == nil {
		t.Fatal("cross-tenant account was accepted on a journal line")
	}
}

func TestCFR002RestrictedRuntimeRoleCannotWriteAccountingHistory(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	for _, table := range []string{"journal_entries", "journal_entry_lines", "general_ledger"} {
		for _, privilege := range []string{"INSERT", "UPDATE", "DELETE"} {
			var granted bool
			if err := db.QueryRow(`SELECT has_table_privilege('pos_app', $1, $2)`, table, privilege).Scan(&granted); err != nil {
				t.Fatalf("check pos_app %s privilege on %s: %v", privilege, table, err)
			}
			if granted {
				t.Fatalf("restricted pos_app unexpectedly has %s on %s", privilege, table)
			}
		}
	}
}

func createCFR002DraftSQL(t *testing.T, db *sql.DB, fixture *atomicPostingFixture, debit, credit float64, includeLines bool) int {
	t.Helper()
	var entryID int
	entryNumber := fmt.Sprintf("C2%016x", uint64(time.Now().UnixNano()))
	date := time.Now().UTC().Format("2006-01-02")
	period := time.Now().UTC().Format("2006-01")
	err := db.QueryRow(`
		INSERT INTO journal_entries (tenant_id, restaurant_id, entry_number, entry_date, fiscal_period, description, entry_type, status)
		VALUES ($1, $2, $3, $4, $5, 'CFR-002 direct SQL fixture', 'manual', 'draft')
		RETURNING id
	`, fixture.tenantID, fixture.restaurantID, entryNumber, date, period).Scan(&entryID)
	if err != nil {
		t.Fatalf("create draft journal by direct SQL: %v", err)
	}
	if includeLines {
		if _, err := db.Exec(`INSERT INTO journal_entry_lines (journal_entry_id, tenant_id, account_id, account_code, debit_amount, line_order) VALUES ($1, $2, $3, 'debit', $4, 1)`, entryID, fixture.tenantID, fixture.debitAccountID, debit); err != nil {
			t.Fatalf("insert debit line: %v", err)
		}
		if _, err := db.Exec(`INSERT INTO journal_entry_lines (journal_entry_id, tenant_id, account_id, account_code, credit_amount, line_order) VALUES ($1, $2, $3, 'credit', $4, 2)`, entryID, fixture.tenantID, fixture.creditAccountID, credit); err != nil {
			t.Fatalf("insert credit line: %v", err)
		}
	}
	return entryID
}

func cfr002LineIDs(t *testing.T, db *sql.DB, entryID int) []int {
	t.Helper()
	rows, err := db.Query(`SELECT id FROM journal_entry_lines WHERE journal_entry_id = $1 ORDER BY line_order`, entryID)
	if err != nil {
		t.Fatalf("query journal line IDs: %v", err)
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("scan journal line ID: %v", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate journal line IDs: %v", err)
	}
	return ids
}

func postCFR002DirectSQL(t *testing.T, db *sql.DB, fixture *atomicPostingFixture, entryID int, lineIDs []int) {
	t.Helper()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin direct SQL posting: %v", err)
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`UPDATE journal_entries SET status = 'posted', posted_at = CURRENT_TIMESTAMP, posted_by = 1 WHERE id = $1`, entryID); err != nil {
		t.Fatalf("transition balanced journal to posted: %v", err)
	}
	for index, lineID := range lineIDs {
		var accountID int
		var debit, credit float64
		if err := tx.QueryRow(`SELECT account_id, debit_amount, credit_amount FROM journal_entry_lines WHERE id = $1`, lineID).Scan(&accountID, &debit, &credit); err != nil {
			t.Fatalf("read line for ledger insert: %v", err)
		}
		if _, err := tx.Exec(`
			INSERT INTO general_ledger (tenant_id, restaurant_id, account_id, journal_entry_id, journal_line_id, entry_date, fiscal_period, debit_amount, credit_amount, running_balance)
			SELECT $1, $2, $3, $4, $5, entry_date, fiscal_period, $6, $7, 0 FROM journal_entries WHERE id = $4
		`, fixture.tenantID, fixture.restaurantID, accountID, entryID, lineID, debit, credit); err != nil {
			t.Fatalf("insert ledger effect %d: %v", index, err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit balanced direct SQL posting: %v", err)
	}
}

func assertJournalLedgerEffects(t *testing.T, db *sql.DB, entryID, expectedPerLine int) {
	t.Helper()
	var badLines int
	if err := db.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT jel.id FROM journal_entry_lines jel
			LEFT JOIN general_ledger gl ON gl.journal_entry_id = jel.journal_entry_id AND gl.journal_line_id = jel.id
			WHERE jel.journal_entry_id = $1 GROUP BY jel.id HAVING COUNT(gl.id) <> $2
		) invalid
	`, entryID, expectedPerLine).Scan(&badLines); err != nil {
		t.Fatalf("count invalid line ledger effects: %v", err)
	}
	if badLines != 0 {
		t.Fatalf("journal %d has %d lines with invalid ledger effect counts", entryID, badLines)
	}
}
