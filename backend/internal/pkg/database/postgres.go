package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/config"
)

func Connect(cfg *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Phase 14.4: make pool tunable via env while preserving safe defaults.
	maxOpenConns := cfg.MaxOpenConns
	if maxOpenConns <= 0 {
		maxOpenConns = 25
	}
	maxIdleConns := cfg.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = 10
	}

	connMaxLifetime := 5 * time.Minute
	if cfg.ConnMaxLifetime != "" {
		parsed, err := time.ParseDuration(cfg.ConnMaxLifetime)
		if err == nil && parsed > 0 {
			connMaxLifetime = parsed
		}
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	log.Printf("Database connected successfully (max_open=%d max_idle=%d conn_max_lifetime=%s)", maxOpenConns, maxIdleConns, connMaxLifetime)
	return db, nil
}
