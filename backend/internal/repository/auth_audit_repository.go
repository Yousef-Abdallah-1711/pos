package repository

import (
	"context"
	"database/sql"
	"errors"

	"pos-saas/internal/pkg/database"
)

// AuthAuditRepository writes successful browser authentication events.
type AuthAuditRepository struct {
	db *sql.DB
}

func NewAuthAuditRepository(db *sql.DB) *AuthAuditRepository {
	return &AuthAuditRepository{db: db}
}

func (r *AuthAuditRepository) RecordAuthEvent(ctx context.Context, tenantID, userID int64, action, status, ipAddress, userAgent string) error {
	if r == nil || r.db == nil {
		return errors.New("auth audit database is not configured")
	}
	if action != "login" && action != "logout" && action != "refresh" {
		return errors.New("unsupported auth audit action")
	}
	if status != "success" {
		return errors.New("unsupported auth audit status")
	}
	if tenantID <= 0 {
		return errors.New("tenant id is required for auth audit events")
	}
	return database.WithTenantTx(ctx, r.db, tenantID, 0, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
		INSERT INTO user_auth_audit_logs (tenant_id, user_id, action, status, ip_address, user_agent)
		VALUES ($1, NULLIF($2, 0), $3, $4, NULLIF($5, ''), NULLIF($6, ''))`,
			tenantID, userID, action, status, ipAddress, userAgent,
		)
		return err
	})
}
