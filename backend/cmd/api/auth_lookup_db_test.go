package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"pos-saas/internal/config"
)

func TestAuthLookupPoolRequiresDedicatedIdentity(t *testing.T) {
	base := config.DatabaseConfig{Host: "localhost", User: "pos_app", Password: "synthetic-core", DBName: "servio_test", MaxOpenConns: 20}
	t.Setenv("AUTH_LOOKUP_DB_PASSWORD", "synthetic-lookup")
	for _, user := range []string{"", "pos_app", "pos_migrator", "pos_accounting", "postgres"} {
		t.Setenv("AUTH_LOOKUP_DB_USER", user)
		if _, err := authLookupDatabaseConfig(base); err == nil {
			t.Fatalf("accepted authentication lookup identity %q", user)
		}
	}
	t.Setenv("AUTH_LOOKUP_DB_USER", "pos_auth_lookup")
	cfg, err := authLookupDatabaseConfig(base)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.User != "pos_auth_lookup" || cfg.Password != "synthetic-lookup" || cfg.MaxOpenConns != 5 || cfg.MaxIdleConns != 2 {
		t.Fatal("dedicated bounded auth lookup pool not configured")
	}
	if base.User != "pos_app" || base.Password != "synthetic-core" {
		t.Fatal("core database identity was mutated")
	}
	t.Setenv("AUTH_LOOKUP_DB_PASSWORD", "")
	if _, err = authLookupDatabaseConfig(base); err == nil {
		t.Fatal("missing auth lookup secret accepted")
	}
}

func TestAuthLookupPoolIsDisabledWithoutCredentialsAndNeverFallsBack(t *testing.T) {
	t.Setenv("AUTH_LOOKUP_DB_USER", "")
	t.Setenv("AUTH_LOOKUP_DB_PASSWORD", "")
	db, err := connectAuthLookupDatabase(&sql.DB{}, config.DatabaseConfig{})
	if err != nil || db != nil {
		t.Fatalf("unconfigured lookup pool must remain disabled: db=%v err=%v", db, err)
	}
	t.Setenv("AUTH_LOOKUP_DB_USER", "pos_app")
	t.Setenv("AUTH_LOOKUP_DB_PASSWORD", "synthetic-core")
	if _, err := connectAuthLookupDatabase(&sql.DB{}, config.DatabaseConfig{}); err == nil {
		t.Fatal("auth lookup pool accepted a shared application identity")
	}
	t.Setenv("AUTH_LOOKUP_DB_USER", "pos_auth_lookup")
	if _, err := connectAuthLookupDatabase(nil, config.DatabaseConfig{}); err == nil {
		t.Fatal("auth lookup pool allowed a missing primary database")
	}
}

func TestValidateAuthLookupDatabaseRequiresReadOnlyUnprivilegedPool(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT r\.rolname`).WillReturnRows(sqlmock.NewRows([]string{
		"role", "privileged", "database_create", "schema_usage", "schema_create", "default_read_only",
		"transaction_read_only", "identity_function", "memberships_function", "inherited_role",
	}).AddRow("pos_auth_lookup", false, false, true, false, "on", "on", true, true, false))
	mock.ExpectQuery(`SELECT count\(\*\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT count\(\*\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT count\(\*\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT count\(\*\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	if err := validateAuthLookupDatabase(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestValidateAuthLookupDatabaseRejectsSharedOrWritableRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectQuery(`SELECT r\.rolname`).WillReturnRows(sqlmock.NewRows([]string{
		"role", "privileged", "database_create", "schema_usage", "schema_create", "default_read_only",
		"transaction_read_only", "identity_function", "memberships_function", "inherited_role",
	}).AddRow("pos_app", false, false, true, false, "off", "off", false, false, false))
	if err := validateAuthLookupDatabase(context.Background(), db); err == nil {
		t.Fatal("shared writable database role passed auth lookup validation")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
