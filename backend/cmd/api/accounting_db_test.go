package main

import (
	"database/sql"
	"pos-saas/internal/config"
	"testing"
)

func TestAccountingPoolRequiresDedicatedIdentity(t *testing.T) {
	base := config.DatabaseConfig{Host: "localhost", User: "pos_app", Password: "synthetic-core", DBName: "servio_test", MaxOpenConns: 20}
	t.Setenv("ACCOUNTING_DB_PASSWORD", "synthetic-accounting")
	for _, user := range []string{"", "pos_app", "pos_migrator", "postgres"} {
		t.Setenv("ACCOUNTING_DB_USER", user)
		if _, err := accountingDatabaseConfig(base); err == nil {
			t.Fatalf("accepted accounting identity %q", user)
		}
	}
	t.Setenv("ACCOUNTING_DB_USER", "pos_accounting")
	cfg, err := accountingDatabaseConfig(base)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.User != "pos_accounting" || cfg.Password != "synthetic-accounting" || cfg.MaxOpenConns != 5 || cfg.MaxIdleConns != 2 {
		t.Fatal("dedicated bounded accounting pool not configured")
	}
	if base.User != "pos_app" || base.Password != "synthetic-core" {
		t.Fatal("core database identity was mutated")
	}
	t.Setenv("ACCOUNTING_DB_PASSWORD", "")
	if _, err = accountingDatabaseConfig(base); err == nil {
		t.Fatal("missing accounting secret accepted")
	}
}

func TestAccountingPoolCannotFallbackInStagingOrProduction(t *testing.T) {
	t.Setenv("ACCOUNTING_DB_USER", "")
	t.Setenv("ACCOUNTING_DB_PASSWORD", "")
	for _, environment := range []string{"staging", "production", "unknown"} {
		t.Setenv("ENV", environment)
		if _, err := connectAccountingDatabase(&sql.DB{}, config.DatabaseConfig{}); err == nil {
			t.Fatalf("shared accounting fallback allowed in %s", environment)
		}
		if _, err := connectAccountingDatabase(nil, config.DatabaseConfig{}); err == nil {
			t.Fatalf("missing database allowed in %s", environment)
		}
	}
}
