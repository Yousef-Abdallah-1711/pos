package main

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"pos-saas/tests/testdb"
)

func TestExplicitWrapperReceiptFailureRollsBackSchemaAndKeepsConnection(t *testing.T) {
	testWrappedMigrationReceipt(t, "180_enforce_journal_integrity.sql")
}

func TestOutboxWrapperReceiptFailureRollsBackSchemaAndKeepsConnection(t *testing.T) {
	testWrappedMigrationReceipt(t, "181_durable_accounting_event_outbox.sql")
}

func testWrappedMigrationReceipt(t *testing.T, filename string) {
	t.Helper()
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	// Session-local tables cannot alter the real migration receipt or application schema.
	if _, err = conn.ExecContext(ctx, `CREATE TEMP TABLE schema_migrations(version TEXT PRIMARY KEY)`); err != nil {
		t.Fatal(err)
	}
	if _, err = conn.ExecContext(ctx, `INSERT INTO schema_migrations(version) VALUES($1)`, filename); err != nil {
		t.Fatal(err)
	}
	content := "-- reviewed wrapper\nBEGIN;\nCREATE TEMP TABLE migration_receipt_probe(marker INT);\nCOMMIT;\n"
	if err = applyMigration(ctx, conn, filename, content); err == nil {
		t.Fatal("duplicate receipt unexpectedly committed")
	}
	var absent bool
	if err = conn.QueryRowContext(ctx, `SELECT to_regclass('pg_temp.migration_receipt_probe') IS NULL`).Scan(&absent); err != nil {
		t.Fatalf("receipt failure broke the migration connection: %v", err)
	}
	if !absent {
		t.Fatal("schema change survived failed receipt insertion")
	}
	var count int
	if err = conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_migrations WHERE version=$1`, filename).Scan(&count); err != nil || count != 1 {
		t.Fatalf("prior receipt changed: count=%d err=%v", count, err)
	}
}
