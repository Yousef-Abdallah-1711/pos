package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"pos-saas/internal/config"
	"pos-saas/internal/pkg/database"
)

func accountingDatabaseConfig(base config.DatabaseConfig) (config.DatabaseConfig, error) {
	user := strings.TrimSpace(os.Getenv("ACCOUNTING_DB_USER"))
	password := os.Getenv("ACCOUNTING_DB_PASSWORD")
	if user != "pos_accounting" || password == "" {
		return base, fmt.Errorf("dedicated ACCOUNTING_DB_USER and ACCOUNTING_DB_PASSWORD are required")
	}
	base.User = user
	base.Password = password
	base.MaxOpenConns = 5
	base.MaxIdleConns = 2
	return base, nil
}

func connectAccountingDatabase(primary *sql.DB, base config.DatabaseConfig) (*sql.DB, error) {
	allowFallback := config.IsDevelopmentMode() || config.CurrentEnv() == "test"
	if primary == nil {
		if allowFallback {
			return nil, nil
		}
		return nil, fmt.Errorf("primary database is required for accounting")
	}
	if allowFallback && os.Getenv("ACCOUNTING_DB_USER") == "" && os.Getenv("ACCOUNTING_DB_PASSWORD") == "" {
		return primary, nil
	}
	cfg, err := accountingDatabaseConfig(base)
	if err != nil {
		return nil, err
	}
	db, err := database.Connect(&cfg)
	if err != nil {
		return nil, fmt.Errorf("accounting database connection failed: %w", err)
	}
	var role string
	var privileged bool
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.QueryRowContext(ctx, `SELECT current_user,rolsuper OR rolcreatedb OR rolcreaterole OR rolbypassrls FROM pg_roles WHERE rolname=current_user`).Scan(&role, &privileged)
	if err != nil || role != "pos_accounting" || privileged {
		db.Close()
		return nil, fmt.Errorf("accounting database must use the restricted pos_accounting role")
	}
	return db, nil
}
