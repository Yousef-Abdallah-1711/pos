package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/config"
	"pos-saas/internal/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("✅ All migrations completed successfully")
}

func runMigrations(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	return runMigrationsFromDir(ctx, db, "migrations")
}

const migrationAdvisoryLockID int64 = 0x504f53504f53

// runMigrationsFromDir serializes migration runners on a dedicated PostgreSQL
// connection so the session advisory lock covers every status check and DDL run.
func runMigrationsFromDir(ctx context.Context, db *sql.DB, migrationsDir string) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire migration connection: %w", err)
	}
	defer conn.Close()

	var locked interface{}
	if err := conn.QueryRowContext(ctx, "SELECT pg_advisory_lock($1)", migrationAdvisoryLockID).Scan(&locked); err != nil {
		return fmt.Errorf("failed to acquire migration lock: %w", err)
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var unlocked bool
		if err := conn.QueryRowContext(unlockCtx, "SELECT pg_advisory_unlock($1)", migrationAdvisoryLockID).Scan(&unlocked); err != nil {
			log.Printf("failed to release migration lock: %v", err)
		}
	}()

	return runMigrationsOnConn(ctx, conn, migrationsDir)
}

func runMigrationsOnConn(ctx context.Context, conn *sql.Conn, migrationsDir string) error {
	// Create migrations table if it doesn't exist
	_, err := conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Get migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort files by name
	var sqlFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Execute each migration
	for _, filename := range sqlFiles {
		// Check if migration has already been applied
		var count int
		err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = $1", filename).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count > 0 {
			log.Printf("⏭️  Skipping (already applied): %s", filename)
			continue
		}

		log.Printf("Running migration: %s", filename)

		content, err := os.ReadFile(filepath.Join(migrationsDir, filename))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		if _, err := conn.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Mark migration as applied
		_, err = conn.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", filename)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		log.Printf("✅ Completed: %s", filename)
	}

	return nil
}
