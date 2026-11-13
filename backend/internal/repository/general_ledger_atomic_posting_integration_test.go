//go:build integration

package repository

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"pos-saas/internal/domain"
	"pos-saas/tests/testdb"
)

func TestCreateAndPostAutoEntryPostgresSerializesDuplicateEvents(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	const concurrentAttempts = 8
	start := make(chan struct{})
	type result struct {
		id  int
		err error
	}
	results := make(chan result, concurrentAttempts)
	var workers sync.WaitGroup
	for range concurrentAttempts {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			id, err := fixture.ledger.CreateAndPostAutoEntry(fixture.entry(), 0)
			results <- result{id: id, err: err}
		}()
	}
	close(start)
	workers.Wait()
	close(results)

	var journalID int
	for result := range results {
		if result.err != nil {
			t.Fatalf("CreateAndPostAutoEntry returned an error: %v", result.err)
		}
		if journalID == 0 {
			journalID = result.id
		}
		if result.id != journalID {
			t.Fatalf("duplicate source event returned journal IDs %d and %d", journalID, result.id)
		}
	}

	assertAtomicPostingCounts(t, db, fixture, 1, 2)
	var status string
	if err := db.QueryRow(`SELECT status FROM journal_entries WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, journalID, fixture.tenantID, fixture.restaurantID).Scan(&status); err != nil {
		t.Fatalf("read posted journal status: %v", err)
	}
	if status != domain.JournalStatusPosted {
		t.Fatalf("journal status = %q, want posted", status)
	}
}

func TestCreateAndPostAutoEntryPostgresKeepsTenantAndRestaurantScopes(t *testing.T) {
	db, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, db, fixture)

	otherTenant := openAtomicPostingTenantFixture(t, db)
	defer cleanupAtomicPostingFixtureRows(t, db, otherTenant)
	siblingRestaurant := openAtomicPostingRestaurantFixture(t, db, fixture.tenantID)

	for _, scopedFixture := range []*atomicPostingFixture{fixture, otherTenant, siblingRestaurant} {
		if _, err := scopedFixture.ledger.CreateAndPostAutoEntry(scopedFixture.entry(), 0); err != nil {
			t.Fatalf("CreateAndPostAutoEntry for tenant %d restaurant %d: %v", scopedFixture.tenantID, scopedFixture.restaurantID, err)
		}
		assertAtomicPostingCounts(t, db, scopedFixture, 1, 2)
	}
}

func TestCreateAndPostAutoEntryPostgresRollsBackEveryFailureBoundary(t *testing.T) {
	stages := []struct {
		name   string
		table  string
		timing string
	}{
		{name: "journal line insert", table: "journal_entry_lines", timing: "BEFORE INSERT"},
		{name: "general ledger insert", table: "general_ledger", timing: "BEFORE INSERT"},
		{name: "journal posted status update", table: "journal_entries", timing: "BEFORE UPDATE"},
	}

	for _, stage := range stages {
		t.Run(stage.name, func(t *testing.T) {
			db, fixture := openAtomicPostingFixture(t)
			defer cleanupAtomicPostingFixture(t, db, fixture)
			triggerName := fmt.Sprintf("cfr004_fail_%d", fixture.tenantID)
			functionName := triggerName + "_fn"
			triggerSQL := fmt.Sprintf(`
				CREATE FUNCTION %s() RETURNS trigger LANGUAGE plpgsql AS $$
				BEGIN
					IF NEW.tenant_id = %d THEN
						RAISE EXCEPTION 'intentional CFR-004 integration failure';
					END IF;
					RETURN NEW;
				END
				$$;
				CREATE TRIGGER %s %s ON %s FOR EACH ROW EXECUTE FUNCTION %s();
			`, functionName, fixture.tenantID, triggerName, stage.timing, stage.table, functionName)
			if _, err := db.Exec(triggerSQL); err != nil {
				t.Fatalf("install failure-injection trigger: %v", err)
			}
			defer func() {
				if _, err := db.Exec(fmt.Sprintf(`DROP TRIGGER IF EXISTS %s ON %s`, triggerName, stage.table)); err != nil {
					t.Errorf("drop failure-injection trigger: %v", err)
				}
				if _, err := db.Exec(fmt.Sprintf(`DROP FUNCTION IF EXISTS %s()`, functionName)); err != nil {
					t.Errorf("drop failure-injection function: %v", err)
				}
			}()

			if _, err := fixture.ledger.CreateAndPostAutoEntry(fixture.entry(), 0); err == nil {
				t.Fatal("CreateAndPostAutoEntry succeeded despite injected database failure")
			}
			assertAtomicPostingCounts(t, db, fixture, 0, 0)
			for _, accountID := range []int{fixture.debitAccountID, fixture.creditAccountID} {
				var balance float64
				if err := db.QueryRow(`SELECT current_balance FROM chart_of_accounts WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`, accountID, fixture.tenantID, fixture.restaurantID).Scan(&balance); err != nil {
					t.Fatalf("read account balance after rollback: %v", err)
				}
				if balance != 0 {
					t.Fatalf("account %d balance = %v after rollback, want 0", accountID, balance)
				}
			}
		})
	}
}

type atomicPostingFixture struct {
	tenantID        int
	restaurantID    int
	debitAccountID  int
	creditAccountID int
	ledger          *GeneralLedgerRepository
}

func openAtomicPostingFixture(t *testing.T) (*sql.DB, *atomicPostingFixture) {
	t.Helper()
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	db.SetMaxOpenConns(20)
	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}

	var fixture atomicPostingFixture
	name := fmt.Sprintf("cfr004-%d", time.Now().UnixNano())
	if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1, $2, $3) RETURNING id`, name, name, name+".example.test").Scan(&fixture.tenantID); err != nil {
		db.Close()
		t.Fatalf("create isolated tenant fixture: %v", err)
	}
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, 'CFR-004 test restaurant', 'cfr004-test', 'test') RETURNING id`, fixture.tenantID).Scan(&fixture.restaurantID); err != nil {
		cleanupAtomicPostingFixture(t, db, &fixture)
		t.Fatalf("create isolated restaurant fixture: %v", err)
	}
	if err := db.QueryRow(`
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
		VALUES ($1, $2, 'CFR004-D', 'CFR-004 debit fixture', 'asset', 'debit') RETURNING id
	`, fixture.tenantID, fixture.restaurantID).Scan(&fixture.debitAccountID); err != nil {
		cleanupAtomicPostingFixture(t, db, &fixture)
		t.Fatalf("create debit account fixture: %v", err)
	}
	if err := db.QueryRow(`
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
		VALUES ($1, $2, 'CFR004-C', 'CFR-004 credit fixture', 'revenue', 'credit') RETURNING id
	`, fixture.tenantID, fixture.restaurantID).Scan(&fixture.creditAccountID); err != nil {
		cleanupAtomicPostingFixture(t, db, &fixture)
		t.Fatalf("create credit account fixture: %v", err)
	}
	journalRepo := NewJournalEntryRepository(db)
	fixture.ledger = NewGeneralLedgerRepository(db, journalRepo)
	return db, &fixture
}

func openAtomicPostingTenantFixture(t *testing.T, db *sql.DB) *atomicPostingFixture {
	t.Helper()
	fixture := newAtomicPostingLocationFixture(t, db, 0)
	return fixture
}

func openAtomicPostingRestaurantFixture(t *testing.T, db *sql.DB, tenantID int) *atomicPostingFixture {
	t.Helper()
	return newAtomicPostingLocationFixture(t, db, tenantID)
}

func newAtomicPostingLocationFixture(t *testing.T, db *sql.DB, tenantID int) *atomicPostingFixture {
	t.Helper()
	if tenantID == 0 {
		name := fmt.Sprintf("cfr004-scope-%d", time.Now().UnixNano())
		if err := db.QueryRow(`INSERT INTO tenants (name, slug, domain) VALUES ($1, $2, $3) RETURNING id`, name, name, name+".example.test").Scan(&tenantID); err != nil {
			t.Fatalf("create isolated CFR-004 tenant fixture: %v", err)
		}
	}
	fixture := &atomicPostingFixture{tenantID: tenantID}
	slug := fmt.Sprintf("cfr004-scope-%d", time.Now().UnixNano())
	if err := db.QueryRow(`INSERT INTO restaurants (tenant_id, name, slug, address) VALUES ($1, 'CFR-004 scope restaurant', $2, 'test') RETURNING id`, tenantID, slug).Scan(&fixture.restaurantID); err != nil {
		t.Fatalf("create isolated CFR-004 restaurant fixture: %v", err)
	}
	if err := db.QueryRow(`
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
		VALUES ($1, $2, 'CFR004-D', 'CFR-004 debit fixture', 'asset', 'debit') RETURNING id
	`, tenantID, fixture.restaurantID).Scan(&fixture.debitAccountID); err != nil {
		t.Fatalf("create isolated CFR-004 debit account fixture: %v", err)
	}
	if err := db.QueryRow(`
		INSERT INTO chart_of_accounts (tenant_id, restaurant_id, account_code, account_name, account_type, normal_balance)
		VALUES ($1, $2, 'CFR004-C', 'CFR-004 credit fixture', 'revenue', 'credit') RETURNING id
	`, tenantID, fixture.restaurantID).Scan(&fixture.creditAccountID); err != nil {
		t.Fatalf("create isolated CFR-004 credit account fixture: %v", err)
	}
	fixture.ledger = NewGeneralLedgerRepository(db, NewJournalEntryRepository(db))
	return fixture
}

func (fixture *atomicPostingFixture) entry() *domain.JournalEntry {
	referenceID := 4004
	entryDate := time.Date(2026, time.November, 13, 0, 0, 0, 0, time.UTC)
	return &domain.JournalEntry{
		TenantID:        fixture.tenantID,
		RestaurantID:    fixture.restaurantID,
		EntryDate:       entryDate,
		Description:     "CFR-004 isolated atomic posting fixture",
		ReferenceType:   "order",
		ReferenceID:     &referenceID,
		ReferenceNumber: "CFR004-4004",
		EntryType:       "auto_sales",
		Status:          domain.JournalStatusDraft,
		Lines: []domain.JournalEntryLine{
			{AccountID: fixture.debitAccountID, DebitAmount: 100, LineOrder: 1},
			{AccountID: fixture.creditAccountID, CreditAmount: 100, LineOrder: 2},
		},
	}
}

func assertAtomicPostingCounts(t *testing.T, db *sql.DB, fixture *atomicPostingFixture, journalCount, ledgerCount int) {
	t.Helper()
	var gotJournals, gotLedger int
	if err := db.QueryRow(`SELECT COUNT(*) FROM journal_entries WHERE tenant_id = $1 AND restaurant_id = $2 AND reference_type = 'order' AND reference_id = 4004 AND entry_type = 'auto_sales'`, fixture.tenantID, fixture.restaurantID).Scan(&gotJournals); err != nil {
		t.Fatalf("count auto-posting journals: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM general_ledger WHERE tenant_id = $1 AND restaurant_id = $2 AND journal_entry_id IN (SELECT id FROM journal_entries WHERE tenant_id = $1 AND restaurant_id = $2 AND reference_type = 'order' AND reference_id = 4004 AND entry_type = 'auto_sales')`, fixture.tenantID, fixture.restaurantID).Scan(&gotLedger); err != nil {
		t.Fatalf("count auto-posting ledger rows: %v", err)
	}
	if gotJournals != journalCount || gotLedger != ledgerCount {
		t.Fatalf("journal/ledger rows = %d/%d, want %d/%d", gotJournals, gotLedger, journalCount, ledgerCount)
	}
}

func cleanupAtomicPostingFixture(t *testing.T, db *sql.DB, fixture *atomicPostingFixture) {
	t.Helper()
	cleanupAtomicPostingFixtureRows(t, db, fixture)
	if err := db.Close(); err != nil {
		t.Errorf("close isolated PostgreSQL connection: %v", err)
	}
}

func cleanupAtomicPostingFixtureRows(t *testing.T, db *sql.DB, fixture *atomicPostingFixture) {
	t.Helper()
	// Financial integration fixtures live only in a dedicated disposable database.
	// CFR-002 makes ledger history immutable, so deleting test postings would weaken
	// the behavior under test; discard the entire database after the suite instead.
}
