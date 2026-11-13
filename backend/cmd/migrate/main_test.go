package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRunMigrationsFromDirLocksAndReleasesConnection(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "001_test.sql"), []byte("SELECT 1;"), 0o600); err != nil {
		t.Fatal(err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT pg_advisory_lock($1)")).
		WithArgs(migrationAdvisoryLockID).WillReturnRows(sqlmock.NewRows([]string{"pg_advisory_lock"}).AddRow(nil))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM schema_migrations WHERE version = $1")).
		WithArgs("001_test.sql").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("SELECT 1;")).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO schema_migrations (version) VALUES ($1)")).
		WithArgs("001_test.sql").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT pg_advisory_unlock($1)")).
		WithArgs(migrationAdvisoryLockID).WillReturnRows(sqlmock.NewRows([]string{"pg_advisory_unlock"}).AddRow(true))

	if err := runMigrationsFromDir(context.Background(), db, dir); err != nil {
		t.Fatalf("runMigrationsFromDir returned an error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestRunMigrationsFromDirReleasesLockAfterMigrationFailure(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "001_test.sql"), []byte("SELECT broken;"), 0o600); err != nil {
		t.Fatal(err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT pg_advisory_lock($1)")).
		WithArgs(migrationAdvisoryLockID).WillReturnRows(sqlmock.NewRows([]string{"pg_advisory_lock"}).AddRow(nil))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM schema_migrations WHERE version = $1")).
		WithArgs("001_test.sql").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("SELECT broken;")).WillReturnError(errors.New("migration failed"))
	mock.ExpectRollback()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT pg_advisory_unlock($1)")).
		WithArgs(migrationAdvisoryLockID).WillReturnRows(sqlmock.NewRows([]string{"pg_advisory_unlock"}).AddRow(true))

	if err := runMigrationsFromDir(context.Background(), db, dir); err == nil {
		t.Fatal("expected migration failure")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}
