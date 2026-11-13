//go:build integration

package repository

import (
	"database/sql"
	"net/url"
	"os"
	"testing"

	"pos-saas/tests/testdb"
)

func accountingWriterConnection(t *testing.T) *sql.DB {
	t.Helper()
	raw := os.Getenv("POS_ACCOUNTING_DATABASE_URL")
	if raw == "" {
		t.Fatal("isolated accounting writer connection must be configured")
	}
	target, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	primary, err := url.Parse(testdb.URL(t))
	if err != nil {
		t.Fatal(err)
	}
	if target.Host != primary.Host || target.Path != primary.Path || target.User == nil || target.User.Username() != "pos_accounting" {
		t.Fatal("accounting probe must use pos_accounting in the same disposable database")
	}
	db, err := sql.Open("postgres", raw)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestAccountingWriterCannotMutateCoreOrBypassJournalGuards(t *testing.T) {
	admin, fixture := openAtomicPostingFixture(t)
	defer cleanupAtomicPostingFixture(t, admin, fixture)
	writer := accountingWriterConnection(t)
	var privileged bool
	if err := writer.QueryRow(`SELECT rolsuper OR rolcreatedb OR rolcreaterole OR rolbypassrls FROM pg_roles WHERE rolname=current_user`).Scan(&privileged); err != nil || privileged {
		t.Fatalf("privileged accounting identity: %v", err)
	}
	for _, table := range []string{"users", "tenants", "restaurants", "products", "orders", "schema_migrations", "order_stock_reservations"} {
		for _, operation := range []string{"INSERT", "UPDATE", "DELETE"} {
			var granted bool
			if err := writer.QueryRow(`SELECT has_table_privilege(current_user,$1,$2)`, table, operation).Scan(&granted); err != nil || granted {
				t.Fatalf("accounting writer has %s on %s: %v", operation, table, err)
			}
		}
	}
	var ddl bool
	if err := writer.QueryRow(`SELECT has_schema_privilege(current_user,'public','CREATE')`).Scan(&ddl); err != nil || ddl {
		t.Fatalf("accounting writer has DDL capability: %v", err)
	}
	entry := createCFR002DraftSQL(t, writer, fixture, 1, 1, true)
	repo := NewJournalEntryRepository(writer)
	if err := repo.PostEntry(fixture.tenantID, fixture.restaurantID, entry, 1); err != nil {
		t.Fatalf("authorized accounting post failed: %v", err)
	}
	if _, err := writer.Exec(`UPDATE general_ledger SET debit_amount=99 WHERE tenant_id=$1 AND restaurant_id=$2 AND journal_entry_id=$3`, fixture.tenantID, fixture.restaurantID, entry); err == nil {
		t.Fatal("accounting writer mutated posted ledger history")
	}
	if err := repo.VoidEntry(fixture.tenantID, fixture.restaurantID, entry, 1, "isolated role regression"); err != nil {
		t.Fatalf("authorized reversal failed: %v", err)
	}
	assertJournalLedgerEffects(t, admin, entry, 2)
}
