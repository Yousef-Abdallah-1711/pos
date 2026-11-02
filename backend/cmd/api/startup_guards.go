package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	// Seeded by migration 109 and documented in the repository; must never survive in production.
	seededSuperAdminEmail    = "superadmin@pos-saas.com"
	seededSuperAdminPassword = "SuperAdmin@123"
)

// isSeededDefaultPassword reports whether hash is the bcrypt hash of the documented seed password.
func isSeededDefaultPassword(hash string) bool {
	if hash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(seededSuperAdminPassword)) == nil
}

// productionStartupGuards refuses to start a production API that is unsafe by default:
// no database (stub-mode authentication) or the seeded super-admin password still in place.
func productionStartupGuards(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return errors.New("production requires a database connection; stub mode is not allowed")
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var hash string
	err := db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE email = $1`, seededSuperAdminEmail).Scan(&hash)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("cannot verify seeded super-admin credentials: %w", err)
	}
	if isSeededDefaultPassword(hash) {
		return fmt.Errorf("the seeded super-admin account %q still uses its documented default password; rotate it before starting in production (UPDATE users SET password_hash = <bcrypt hash> WHERE email = '%s')", seededSuperAdminEmail, seededSuperAdminEmail)
	}
	return nil
}
