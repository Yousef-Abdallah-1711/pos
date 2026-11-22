//go:build integration

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestReconcileJournalGL_PostedAndVoidedEntriesHaveExpectedEffects(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	entryID, err := fixture.ledger.CreateAndPostAutoEntry(fixture.entry(), 0)
	if err != nil {
		t.Fatalf("create synthetic posted journal: %v", err)
	}
	repo := NewJournalEntryRepository(db)
	got, err := repo.ReconcileJournalGL(context.Background(), fixture.tenantID, fixture.restaurantID)
	if err != nil {
		t.Fatalf("reconcile matched journal: %v", err)
	}
	if got.JournalCount != 1 || got.JournalLinesChecked != 2 || got.GLRowCount != 2 || got.Status != "reconciled" ||
		got.PostedLinesMissingGL != 0 || got.DuplicateOrInvalidGL != 0 || got.GLScopeMismatch != 0 ||
		got.GLAccountMismatch != 0 || got.GLAmountMismatch != 0 || got.GLDateMismatch != 0 ||
		got.InvalidLedgerLinks != 0 || got.NonPostedJournalsWithGL != 0 || got.UnbalancedJournals != 0 {
		t.Fatalf("matched posted journal reconciliation = %+v", got)
	}

	if err := repo.VoidEntry(fixture.tenantID, fixture.restaurantID, entryID, 1, "synthetic reconciliation reversal"); err != nil {
		t.Fatalf("void synthetic journal: %v", err)
	}
	got, err = repo.ReconcileJournalGL(context.Background(), fixture.tenantID, fixture.restaurantID)
	if err != nil {
		t.Fatalf("reconcile reversed journal: %v", err)
	}
	if got.Status != "reconciled" || got.GLRowCount != 4 || got.DuplicateOrInvalidGL != 0 ||
		got.PostedLinesMissingGL != 0 || got.GLScopeMismatch != 0 || got.GLAccountMismatch != 0 ||
		got.GLAmountMismatch != 0 || got.GLDateMismatch != 0 || got.InvalidLedgerLinks != 0 ||
		got.NonPostedJournalsWithGL != 0 || got.UnbalancedJournals != 0 {
		t.Fatalf("valid reversed journal reported as invalid: %+v", got)
	}
}

func TestReconcileJournalGL_DetectsIncorrectReversalDateAndPeriod(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)
	entryID, err := fixture.ledger.CreateAndPostAutoEntry(fixture.entry(), 0)
	if err != nil {
		t.Fatalf("post synthetic journal: %v", err)
	}
	if err := NewJournalEntryRepository(db).VoidEntry(fixture.tenantID, fixture.restaurantID, entryID, 1, "synthetic test reversal"); err != nil {
		t.Fatalf("void synthetic journal: %v", err)
	}
	if _, err := db.Exec(`ALTER TABLE general_ledger DISABLE TRIGGER ALL`); err != nil {
		t.Fatalf("disable ledger trigger in disposable test database: %v", err)
	}
	defer func() {
		if _, err := db.Exec(`ALTER TABLE general_ledger ENABLE TRIGGER ALL`); err != nil {
			t.Errorf("re-enable ledger triggers: %v", err)
		}
	}()

	var ledgerID int
	var originalDate time.Time
	var originalPeriod string
	if err := db.QueryRow(`
		SELECT gl.id, gl.entry_date, gl.fiscal_period
		FROM general_ledger gl
		JOIN journal_entry_lines jel ON jel.id = gl.journal_line_id
		JOIN journal_entries je ON je.id = gl.journal_entry_id
		WHERE je.id = $1 AND je.status = 'void'
		  AND gl.debit_amount = jel.credit_amount AND gl.credit_amount = jel.debit_amount
		  AND gl.entry_date = je.voided_at::date
		LIMIT 1
	`, entryID).Scan(&ledgerID, &originalDate, &originalPeriod); err != nil {
		t.Fatalf("find synthetic reversal effect: %v", err)
	}
	repo := NewJournalEntryRepository(db)
	if _, err := db.Exec(`UPDATE general_ledger SET entry_date = entry_date + 1 WHERE id = $1`, ledgerID); err != nil {
		t.Fatalf("change reversal date: %v", err)
	}
	got := reconcileFixture(t, repo, fixture)
	if got.GLDateMismatch != 1 || got.DuplicateOrInvalidGL == 0 || got.Status != "discrepancies" {
		t.Fatalf("incorrect reversal date result = %+v", got)
	}
	if _, err := db.Exec(`UPDATE general_ledger SET entry_date = $1 WHERE id = $2`, originalDate, ledgerID); err != nil {
		t.Fatalf("restore reversal date: %v", err)
	}
	if _, err := db.Exec(`UPDATE general_ledger SET fiscal_period = '2099-01' WHERE id = $1`, ledgerID); err != nil {
		t.Fatalf("change reversal period: %v", err)
	}
	got = reconcileFixture(t, repo, fixture)
	if got.GLDateMismatch != 1 || got.DuplicateOrInvalidGL == 0 || got.Status != "discrepancies" {
		t.Fatalf("incorrect reversal period result = %+v", got)
	}
	if _, err := db.Exec(`UPDATE general_ledger SET fiscal_period = $1 WHERE id = $2`, originalPeriod, ledgerID); err != nil {
		t.Fatalf("restore reversal period: %v", err)
	}
}

func TestReconcileJournalGL_EmptyScopeIsExplicit(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)
	got := reconcileFixture(t, NewJournalEntryRepository(db), fixture)
	if got.Status != "empty" || got.JournalCount != 0 || got.JournalLinesChecked != 0 || got.GLRowCount != 0 {
		t.Fatalf("empty scope result = %+v", got)
	}
}

func TestReconcileJournalGL_DetectsPostedAndVoidedHeadersWithoutLines(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)
	if _, err := db.Exec(`ALTER TABLE journal_entries DISABLE TRIGGER cfr002_journal_integrity_header`); err != nil {
		t.Fatalf("disable deferred integrity trigger for corrupt-state fixture: %v", err)
	}
	defer func() {
		if _, err := db.Exec(`ALTER TABLE journal_entries ENABLE TRIGGER cfr002_journal_integrity_header`); err != nil {
			t.Errorf("re-enable deferred integrity trigger: %v", err)
		}
	}()

	for _, status := range []string{"posted", "void"} {
		var entryID int
		if err := db.QueryRow(`
			INSERT INTO journal_entries (tenant_id, restaurant_id, entry_number, entry_date, fiscal_period,
			  description, entry_type, status)
			VALUES ($1, $2, $3, '2026-11-13', '2026-11', 'line-less diagnostic fixture', 'manual', 'draft')
			RETURNING id
		`, fixture.tenantID, fixture.restaurantID, "EMPTY-"+status).Scan(&entryID); err != nil {
			t.Fatalf("create draft for %s fixture: %v", status, err)
		}
		if _, err := db.Exec(`
			UPDATE journal_entries SET status = 'posted', posted_at = CURRENT_TIMESTAMP, posted_by = 1 WHERE id = $1
		`, entryID); err != nil {
			t.Fatalf("promote line-less fixture to posted: %v", err)
		}
		if status == "void" {
			if _, err := db.Exec(`
				UPDATE journal_entries SET status = 'void', voided_at = CURRENT_TIMESTAMP,
				  voided_by = 1, void_reason = 'synthetic corrupt state' WHERE id = $1
			`, entryID); err != nil {
				t.Fatalf("promote line-less fixture to void: %v", err)
			}
		}
	}

	got := reconcileFixture(t, NewJournalEntryRepository(db), fixture)
	if got.Status != "discrepancies" || got.JournalCount != 2 || got.JournalLinesChecked != 0 ||
		got.PostedJournalsWithoutLines != 2 || got.GLRowCount != 0 {
		t.Fatalf("line-less posted/void journal headers were not detected: %+v", got)
	}
}

func TestReconcileJournalGL_IsTenantAndRestaurantScoped(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)
	other := openAtomicPostingRestaurantFixture(t, db, fixture.tenantID)
	defer cleanupAtomicPostingFixtureRows(t, db, other)

	if _, err := fixture.ledger.CreateAndPostAutoEntry(fixture.entry(), 0); err != nil {
		t.Fatalf("post first synthetic journal: %v", err)
	}
	if _, err := other.ledger.CreateAndPostAutoEntry(other.entry(), 0); err != nil {
		t.Fatalf("post sibling synthetic journal: %v", err)
	}

	got, err := NewJournalEntryRepository(db).ReconcileJournalGL(context.Background(), fixture.tenantID, fixture.restaurantID)
	if err != nil {
		t.Fatalf("reconcile scoped journals: %v", err)
	}
	if got.JournalCount != 1 || got.JournalLinesChecked != 2 || got.GLRowCount != 2 ||
		got.Status != "reconciled" || got.PostedLinesMissingGL != 0 || got.DuplicateOrInvalidGL != 0 {
		t.Fatalf("scoped reconciliation included unexpected data or exceptions: %+v", got)
	}
}

func TestReconcileJournalGL_ConcurrentReadersSeeConsistentResults(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)
	if _, err := fixture.ledger.CreateAndPostAutoEntry(fixture.entry(), 0); err != nil {
		t.Fatalf("post synthetic journal: %v", err)
	}

	const readers = 6
	var workers sync.WaitGroup
	errors := make(chan error, readers)
	start := make(chan struct{})
	for range readers {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			for range 10 {
				got, err := NewJournalEntryRepository(db).ReconcileJournalGL(context.Background(), fixture.tenantID, fixture.restaurantID)
				if err != nil {
					errors <- err
					return
				}
				if got.Status != "reconciled" || got.JournalLinesChecked != got.JournalCount*2 || got.GLRowCount != got.JournalCount*2 {
					errors <- fmt.Errorf("unexpected concurrent reconciliation result: %+v", got)
					return
				}
			}
		}()
	}
	workers.Add(1)
	go func() {
		defer workers.Done()
		<-start
		entry := fixture.entry()
		*entry.ReferenceID = 4005
		entry.ReferenceNumber = "CFR004-4005"
		entry.Description = "CFR-004 concurrent reconciliation fixture"
		if _, err := fixture.ledger.CreateAndPostAutoEntry(entry, 0); err != nil {
			errors <- fmt.Errorf("post concurrent synthetic journal: %w", err)
		}
	}()
	close(start)
	workers.Wait()
	close(errors)
	for err := range errors {
		t.Errorf("concurrent reconciliation reader: %v", err)
	}
}

func TestReconcileJournalGL_DetectsIndividualLedgerAnomalies(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)
	entryID, err := fixture.ledger.CreateAndPostAutoEntry(fixture.entry(), 0)
	if err != nil {
		t.Fatalf("post synthetic journal: %v", err)
	}
	lineIDs := cfr002LineIDs(t, db, entryID)
	if len(lineIDs) != 2 {
		t.Fatalf("fixture journal has %d lines, want 2", len(lineIDs))
	}
	originalRows := make(map[int]ledgerFixtureRow, len(lineIDs))
	for _, lineID := range lineIDs {
		var row ledgerFixtureRow
		if err := db.QueryRow(`
			SELECT id, account_id, debit_amount::text, credit_amount::text, entry_date, fiscal_period, running_balance::text
			FROM general_ledger WHERE journal_entry_id = $1 AND journal_line_id = $2
		`, entryID, lineID).Scan(&row.id, &row.accountID, &row.debit, &row.credit, &row.date, &row.period, &row.runningBalance); err != nil {
			t.Fatalf("save original ledger row: %v", err)
		}
		row.lineID = lineID
		originalRows[lineID] = row
	}
	if _, err := db.Exec(`ALTER TABLE general_ledger DISABLE TRIGGER ALL`); err != nil {
		t.Fatalf("disable integrity triggers in disposable test database: %v", err)
	}
	defer restoreReconciliationFixture(t, db, entryID, originalRows)
	if _, err := db.Exec(`ALTER TABLE journal_entry_lines DISABLE TRIGGER ALL`); err != nil {
		t.Fatalf("disable line triggers in disposable test database: %v", err)
	}
	if _, err := db.Exec(`ALTER TABLE journal_entries DISABLE TRIGGER ALL`); err != nil {
		t.Fatalf("disable header triggers in disposable test database: %v", err)
	}
	repo := NewJournalEntryRepository(db)
	first := originalRows[lineIDs[0]]
	if _, err := db.Exec(`DELETE FROM general_ledger WHERE id = $1`, first.id); err != nil {
		t.Fatalf("remove one effect for missing-row case: %v", err)
	}
	got := reconcileFixture(t, repo, fixture)
	if got.PostedLinesMissingGL != 1 || got.Status != "discrepancies" {
		t.Fatalf("missing GL row result = %+v", got)
	}
	first = restoreOneLedgerRow(t, db, entryID, first)
	originalRows[lineIDs[0]] = first

	var duplicateID int
	if err := db.QueryRow(`
		INSERT INTO general_ledger (tenant_id, restaurant_id, account_id, journal_entry_id, journal_line_id, entry_date, fiscal_period, debit_amount, credit_amount, running_balance)
		SELECT tenant_id, restaurant_id, account_id, journal_entry_id, journal_line_id, entry_date, fiscal_period, debit_amount, credit_amount, running_balance
		FROM general_ledger WHERE id = $1 RETURNING id
	`, first.id).Scan(&duplicateID); err != nil {
		t.Fatalf("insert duplicate synthetic ledger row: %v", err)
	}
	got = reconcileFixture(t, repo, fixture)
	if got.DuplicateOrInvalidGL != 1 || got.Status != "discrepancies" {
		t.Fatalf("duplicate GL row result = %+v", got)
	}
	if _, err := db.Exec(`DELETE FROM general_ledger WHERE id = $1`, duplicateID); err != nil {
		t.Fatalf("remove duplicate synthetic ledger row: %v", err)
	}

	if _, err := db.Exec(`UPDATE general_ledger SET debit_amount = debit_amount + 1 WHERE id = $1`, first.id); err != nil {
		t.Fatalf("change amount for mismatch case: %v", err)
	}
	got = reconcileFixture(t, repo, fixture)
	if got.GLAmountMismatch != 1 {
		t.Fatalf("amount mismatch result = %+v", got)
	}
	if _, err := db.Exec(`UPDATE general_ledger SET debit_amount = $1 WHERE id = $2`, first.debit, first.id); err != nil {
		t.Fatalf("restore amount after mismatch case: %v", err)
	}

	if _, err := db.Exec(`UPDATE general_ledger SET account_id = $1 WHERE id = $2`, fixture.creditAccountID, first.id); err != nil {
		t.Fatalf("change account for mismatch case: %v", err)
	}
	got = reconcileFixture(t, repo, fixture)
	if got.GLAccountMismatch != 1 || got.InvalidLedgerLinks != 1 {
		t.Fatalf("account mismatch result = %+v", got)
	}
	if _, err := db.Exec(`UPDATE general_ledger SET account_id = $1 WHERE id = $2`, first.accountID, first.id); err != nil {
		t.Fatalf("restore account after mismatch case: %v", err)
	}

	if _, err := db.Exec(`UPDATE general_ledger SET entry_date = entry_date + 1 WHERE id = $1`, first.id); err != nil {
		t.Fatalf("change date for mismatch case: %v", err)
	}
	got = reconcileFixture(t, repo, fixture)
	if got.GLDateMismatch != 1 {
		t.Fatalf("date mismatch result = %+v", got)
	}
	if _, err := db.Exec(`UPDATE general_ledger SET entry_date = $1 WHERE id = $2`, first.date, first.id); err != nil {
		t.Fatalf("restore date after mismatch case: %v", err)
	}

	if _, err := db.Exec(`UPDATE general_ledger SET fiscal_period = '2099-01' WHERE id = $1`, first.id); err != nil {
		t.Fatalf("change fiscal period for mismatch case: %v", err)
	}
	got = reconcileFixture(t, repo, fixture)
	if got.GLDateMismatch != 1 {
		t.Fatalf("period mismatch result = %+v", got)
	}
	if _, err := db.Exec(`UPDATE general_ledger SET fiscal_period = $1 WHERE id = $2`, first.period, first.id); err != nil {
		t.Fatalf("restore fiscal period after mismatch case: %v", err)
	}

	if _, err := db.Exec(`UPDATE journal_entries SET status = 'draft' WHERE id = $1`, entryID); err != nil {
		t.Fatalf("change journal state for non-posted case: %v", err)
	}
	got = reconcileFixture(t, repo, fixture)
	if got.NonPostedJournalsWithGL != 1 {
		t.Fatalf("non-posted GL result = %+v", got)
	}
	if _, err := db.Exec(`UPDATE journal_entries SET status = 'posted' WHERE id = $1`, entryID); err != nil {
		t.Fatalf("restore journal state: %v", err)
	}

	otherTenant := openAtomicPostingTenantFixture(t, db)
	if _, err := db.Exec(`UPDATE general_ledger SET tenant_id = $1 WHERE id = $2`, otherTenant.tenantID, first.id); err != nil {
		t.Fatalf("move GL row to another synthetic tenant: %v", err)
	}
	got = reconcileFixture(t, repo, fixture)
	if got.GLScopeMismatch != 1 || got.InvalidLedgerLinks != 1 {
		t.Fatalf("cross-tenant GL result = %+v", got)
	}
	if _, err := db.Exec(`UPDATE general_ledger SET tenant_id = $1 WHERE id = $2`, fixture.tenantID, first.id); err != nil {
		t.Fatalf("restore GL tenant scope: %v", err)
	}

	if _, err := db.Exec(`UPDATE general_ledger SET journal_line_id = 2147483647 WHERE id = $1`, first.id); err != nil {
		t.Fatalf("orphan GL line reference: %v", err)
	}
	got = reconcileFixture(t, repo, fixture)
	if got.PostedLinesMissingGL != 1 || got.InvalidLedgerLinks != 1 {
		t.Fatalf("orphan GL line result = %+v", got)
	}
	if _, err := db.Exec(`UPDATE general_ledger SET journal_line_id = $1 WHERE id = $2`, first.lineID, first.id); err != nil {
		t.Fatalf("restore GL line reference: %v", err)
	}

	if _, err := db.Exec(`UPDATE journal_entry_lines SET debit_amount = debit_amount + 1 WHERE id = $1`, lineIDs[0]); err != nil {
		t.Fatalf("alter source line for unbalanced journal case: %v", err)
	}
	got = reconcileFixture(t, repo, fixture)
	if got.UnbalancedJournals != 1 {
		t.Fatalf("unbalanced journal result = %+v", got)
	}

	if _, err := db.Exec(`UPDATE journal_entry_lines SET debit_amount = debit_amount - 1 WHERE id = $1`, lineIDs[0]); err != nil {
		t.Fatalf("restore balanced source journal: %v", err)
	}
	if _, err := db.Exec(`
		UPDATE journal_entry_lines
		SET debit_amount = debit_amount + CASE WHEN debit_amount > 0 THEN 1 ELSE 0 END,
		    credit_amount = credit_amount + CASE WHEN credit_amount > 0 THEN 1 ELSE 0 END
		WHERE journal_entry_id = $1
	`, entryID); err != nil {
		t.Fatalf("apply offsetting source-line changes: %v", err)
	}
	if _, err := db.Exec(`UPDATE journal_entries SET total_debit = total_debit + 1, total_credit = total_credit + 1 WHERE id = $1`, entryID); err != nil {
		t.Fatalf("keep offsetting source totals balanced: %v", err)
	}
	got = reconcileFixture(t, repo, fixture)
	if got.GLAmountMismatch != 2 || got.UnbalancedJournals != 0 {
		t.Fatalf("offsetting line mismatch was hidden or misclassified: %+v", got)
	}

	if got.GLRowCount != 2 || got.JournalCount != 1 {
		t.Fatalf("fixture counts changed unexpectedly: %+v", got)
	}
}

type ledgerFixtureRow struct {
	id             int
	lineID         int
	accountID      int
	debit          string
	credit         string
	date           time.Time
	period         string
	runningBalance string
}

func reconcileFixture(t *testing.T, repo *JournalEntryRepository, fixture *atomicPostingFixture) JournalGLReconciliation {
	t.Helper()
	got, err := repo.ReconcileJournalGL(context.Background(), fixture.tenantID, fixture.restaurantID)
	if err != nil {
		t.Fatalf("reconcile synthetic fixture: %v", err)
	}
	return got
}

func restoreOneLedgerRow(t *testing.T, db *sql.DB, entryID int, row ledgerFixtureRow) ledgerFixtureRow {
	t.Helper()
	if err := db.QueryRow(`
		INSERT INTO general_ledger (tenant_id, restaurant_id, account_id, journal_entry_id, journal_line_id, entry_date, fiscal_period, debit_amount, credit_amount, running_balance)
		SELECT tenant_id, restaurant_id, $1, id, $2, $3, $4, $5, $6, $7
		FROM journal_entries WHERE id = $8 RETURNING id
	`, row.accountID, row.lineID, row.date, row.period, row.debit, row.credit, row.runningBalance, entryID).Scan(&row.id); err != nil {
		t.Fatalf("restore missing synthetic GL row: %v", err)
	}
	return row
}

func restoreReconciliationFixture(t *testing.T, db *sql.DB, entryID int, rows map[int]ledgerFixtureRow) {
	t.Helper()
	if _, err := db.Exec(`DELETE FROM general_ledger WHERE journal_entry_id = $1`, entryID); err != nil {
		t.Errorf("remove anomaly rows for journal %d: %v", entryID, err)
	}
	for _, row := range rows {
		if _, err := db.Exec(`
			INSERT INTO general_ledger (tenant_id, restaurant_id, account_id, journal_entry_id, journal_line_id, entry_date, fiscal_period, debit_amount, credit_amount, running_balance)
			SELECT tenant_id, restaurant_id, $1, id, $2, $3, $4, $5, $6, $7
			FROM journal_entries WHERE id = $8
		`, row.accountID, row.lineID, row.date, row.period, row.debit, row.credit, row.runningBalance, entryID); err != nil {
			t.Errorf("restore original GL effect for journal line %d: %v", row.lineID, err)
		}
		if _, err := db.Exec(`
			UPDATE journal_entry_lines SET debit_amount = $1, credit_amount = $2 WHERE id = $3
		`, row.debit, row.credit, row.lineID); err != nil {
			t.Errorf("restore source journal line %d: %v", row.lineID, err)
		}
	}
	if _, err := db.Exec(`
		UPDATE journal_entries je
		SET total_debit = (SELECT COALESCE(SUM(debit_amount), 0) FROM journal_entry_lines WHERE journal_entry_id = je.id),
		    total_credit = (SELECT COALESCE(SUM(credit_amount), 0) FROM journal_entry_lines WHERE journal_entry_id = je.id)
		WHERE je.id = $1
	`, entryID); err != nil {
		t.Errorf("restore journal header totals: %v", err)
	}
	if _, err := db.Exec(`UPDATE journal_entries SET status = 'posted' WHERE id = $1`, entryID); err != nil {
		t.Errorf("restore posted journal state: %v", err)
	}
	if _, err := db.Exec(`ALTER TABLE general_ledger ENABLE TRIGGER ALL`); err != nil {
		t.Errorf("re-enable ledger integrity triggers: %v", err)
	}
	if _, err := db.Exec(`ALTER TABLE journal_entry_lines ENABLE TRIGGER ALL`); err != nil {
		t.Errorf("re-enable journal-line integrity triggers: %v", err)
	}
	if _, err := db.Exec(`ALTER TABLE journal_entries ENABLE TRIGGER ALL`); err != nil {
		t.Errorf("re-enable journal-header integrity triggers: %v", err)
	}
}
