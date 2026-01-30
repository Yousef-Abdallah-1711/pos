package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"pos-saas/internal/domain"
)

// ImpersonationRepository handles impersonation session database operations
type ImpersonationRepository struct {
	db *sql.DB
}

// NewImpersonationRepository creates a new impersonation repository
func NewImpersonationRepository(db *sql.DB) *ImpersonationRepository {
	return &ImpersonationRepository{db: db}
}

// CreateSession creates a new impersonation session
func (r *ImpersonationRepository) CreateSession(ctx context.Context, session *domain.ImpersonationSession) (*domain.ImpersonationSession, error) {
	log.Printf("[ImpersonationRepository] CreateSession: superAdminID=%d, targetUserID=%d\n", session.SuperAdminID, session.TargetUserID)

	// Hash the token for storage
	tokenHash := domain.HashToken(session.ImpersonationToken)

	var result domain.ImpersonationSession
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO impersonation_sessions
			(super_admin_id, target_user_id, target_tenant_id, impersonation_token, token_hash,
			 ip_address, user_agent, started_at, expires_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), $8, true)
		RETURNING id, super_admin_id, target_user_id, target_tenant_id, token_hash,
		          ip_address, user_agent, started_at, expires_at, is_active
	`, session.SuperAdminID, session.TargetUserID, session.TargetTenantID,
		session.ImpersonationToken, tokenHash, session.IPAddress, session.UserAgent,
		session.ExpiresAt,
	).Scan(
		&result.ID, &result.SuperAdminID, &result.TargetUserID, &result.TargetTenantID,
		&result.TokenHash, &result.IPAddress, &result.UserAgent, &result.StartedAt,
		&result.ExpiresAt, &result.IsActive,
	)

	if err != nil {
		log.Printf("[ImpersonationRepository] ERROR: Failed to create session: %v\n", err)
		return nil, fmt.Errorf("failed to create impersonation session: %w", err)
	}

	// Copy the original token (not stored in DB, but needed for response)
	result.ImpersonationToken = session.ImpersonationToken

	log.Printf("[ImpersonationRepository] Session created: id=%d\n", result.ID)
	return &result, nil
}

// GetActiveSession retrieves the active impersonation session for a super admin
func (r *ImpersonationRepository) GetActiveSession(ctx context.Context, superAdminID int64) (*domain.ImpersonationSession, error) {
	log.Printf("[ImpersonationRepository] GetActiveSession: superAdminID=%d\n", superAdminID)

	var session domain.ImpersonationSession
	var endedAt sql.NullTime
	var endReason sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT id, super_admin_id, target_user_id, target_tenant_id, token_hash,
		       ip_address, user_agent, started_at, expires_at, ended_at, is_active, end_reason
		FROM impersonation_sessions
		WHERE super_admin_id = $1 AND is_active = true AND expires_at > NOW()
		ORDER BY started_at DESC
		LIMIT 1
	`, superAdminID).Scan(
		&session.ID, &session.SuperAdminID, &session.TargetUserID, &session.TargetTenantID,
		&session.TokenHash, &session.IPAddress, &session.UserAgent, &session.StartedAt,
		&session.ExpiresAt, &endedAt, &session.IsActive, &endReason,
	)

	if err == sql.ErrNoRows {
		log.Printf("[ImpersonationRepository] No active session for super admin %d\n", superAdminID)
		return nil, nil
	}

	if err != nil {
		log.Printf("[ImpersonationRepository] ERROR: Failed to query session: %v\n", err)
		return nil, fmt.Errorf("failed to fetch session: %w", err)
	}

	if endedAt.Valid {
		session.EndedAt = &endedAt.Time
	}
	if endReason.Valid {
		session.EndReason = &endReason.String
	}

	return &session, nil
}

// ValidateToken validates an impersonation token and returns the session
func (r *ImpersonationRepository) ValidateToken(ctx context.Context, token string) (*domain.ImpersonationSession, error) {
	log.Printf("[ImpersonationRepository] ValidateToken\n")

	tokenHash := domain.HashToken(token)

	var session domain.ImpersonationSession
	var endedAt sql.NullTime
	var endReason sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT id, super_admin_id, target_user_id, target_tenant_id, token_hash,
		       ip_address, user_agent, started_at, expires_at, ended_at, is_active, end_reason
		FROM impersonation_sessions
		WHERE token_hash = $1
	`, tokenHash).Scan(
		&session.ID, &session.SuperAdminID, &session.TargetUserID, &session.TargetTenantID,
		&session.TokenHash, &session.IPAddress, &session.UserAgent, &session.StartedAt,
		&session.ExpiresAt, &endedAt, &session.IsActive, &endReason,
	)

	if err == sql.ErrNoRows {
		log.Printf("[ImpersonationRepository] Token not found\n")
		return nil, fmt.Errorf("invalid token")
	}

	if err != nil {
		log.Printf("[ImpersonationRepository] ERROR: Failed to validate token: %v\n", err)
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	if endedAt.Valid {
		session.EndedAt = &endedAt.Time
	}
	if endReason.Valid {
		session.EndReason = &endReason.String
	}

	// Check if session is still valid
	if !session.IsActive {
		log.Printf("[ImpersonationRepository] Session is not active\n")
		return nil, fmt.Errorf("session is not active")
	}

	if time.Now().After(session.ExpiresAt) {
		log.Printf("[ImpersonationRepository] Session has expired\n")
		// Auto-expire the session
		_ = r.EndSession(ctx, session.ID, domain.ImpersonationEndExpired)
		return nil, fmt.Errorf("session has expired")
	}

	return &session, nil
}

// EndSession ends an impersonation session
func (r *ImpersonationRepository) EndSession(ctx context.Context, sessionID int64, reason string) error {
	log.Printf("[ImpersonationRepository] EndSession: id=%d, reason=%s\n", sessionID, reason)

	result, err := r.db.ExecContext(ctx, `
		UPDATE impersonation_sessions
		SET is_active = false, ended_at = NOW(), end_reason = $1
		WHERE id = $2 AND is_active = true
	`, reason, sessionID)

	if err != nil {
		log.Printf("[ImpersonationRepository] ERROR: Failed to end session: %v\n", err)
		return fmt.Errorf("failed to end session: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("[ImpersonationRepository] Session not found or already ended: %d\n", sessionID)
		return fmt.Errorf("session not found or already ended")
	}

	log.Printf("[ImpersonationRepository] Session ended: id=%d\n", sessionID)
	return nil
}

// EndAllSessionsForSuperAdmin ends all active sessions for a super admin
func (r *ImpersonationRepository) EndAllSessionsForSuperAdmin(ctx context.Context, superAdminID int64, reason string) error {
	log.Printf("[ImpersonationRepository] EndAllSessionsForSuperAdmin: superAdminID=%d\n", superAdminID)

	_, err := r.db.ExecContext(ctx, `
		UPDATE impersonation_sessions
		SET is_active = false, ended_at = NOW(), end_reason = $1
		WHERE super_admin_id = $2 AND is_active = true
	`, reason, superAdminID)

	if err != nil {
		log.Printf("[ImpersonationRepository] ERROR: Failed to end sessions: %v\n", err)
		return fmt.Errorf("failed to end sessions: %w", err)
	}

	return nil
}

// GetSessionByID retrieves a session by its ID
func (r *ImpersonationRepository) GetSessionByID(ctx context.Context, sessionID int64) (*domain.ImpersonationSession, error) {
	log.Printf("[ImpersonationRepository] GetSessionByID: %d\n", sessionID)

	var session domain.ImpersonationSession
	var endedAt sql.NullTime
	var endReason sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT id, super_admin_id, target_user_id, target_tenant_id, token_hash,
		       ip_address, user_agent, started_at, expires_at, ended_at, is_active, end_reason
		FROM impersonation_sessions
		WHERE id = $1
	`, sessionID).Scan(
		&session.ID, &session.SuperAdminID, &session.TargetUserID, &session.TargetTenantID,
		&session.TokenHash, &session.IPAddress, &session.UserAgent, &session.StartedAt,
		&session.ExpiresAt, &endedAt, &session.IsActive, &endReason,
	)

	if err == sql.ErrNoRows {
		log.Printf("[ImpersonationRepository] Session not found: %d\n", sessionID)
		return nil, fmt.Errorf("session not found")
	}

	if err != nil {
		log.Printf("[ImpersonationRepository] ERROR: Failed to query session: %v\n", err)
		return nil, fmt.Errorf("failed to fetch session: %w", err)
	}

	if endedAt.Valid {
		session.EndedAt = &endedAt.Time
	}
	if endReason.Valid {
		session.EndReason = &endReason.String
	}

	return &session, nil
}

// ListSessionsForSuperAdmin lists all sessions for a super admin
func (r *ImpersonationRepository) ListSessionsForSuperAdmin(ctx context.Context, superAdminID int64, pagination domain.Pagination) ([]domain.ImpersonationSession, int64, error) {
	log.Printf("[ImpersonationRepository] ListSessionsForSuperAdmin: superAdminID=%d\n", superAdminID)

	// Count total
	var total int64
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM impersonation_sessions WHERE super_admin_id = $1
	`, superAdminID).Scan(&total)
	if err != nil {
		log.Printf("[ImpersonationRepository] ERROR: Failed to count sessions: %v\n", err)
		return nil, 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	// Query with pagination
	offset := (pagination.Page - 1) * pagination.PageSize
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, super_admin_id, target_user_id, target_tenant_id, token_hash,
		       ip_address, user_agent, started_at, expires_at, ended_at, is_active, end_reason
		FROM impersonation_sessions
		WHERE super_admin_id = $1
		ORDER BY started_at DESC
		LIMIT $2 OFFSET $3
	`, superAdminID, pagination.PageSize, offset)

	if err != nil {
		log.Printf("[ImpersonationRepository] ERROR: Failed to query sessions: %v\n", err)
		return nil, 0, fmt.Errorf("failed to fetch sessions: %w", err)
	}
	defer rows.Close()

	var sessions []domain.ImpersonationSession
	for rows.Next() {
		var session domain.ImpersonationSession
		var endedAt sql.NullTime
		var endReason sql.NullString

		err := rows.Scan(
			&session.ID, &session.SuperAdminID, &session.TargetUserID, &session.TargetTenantID,
			&session.TokenHash, &session.IPAddress, &session.UserAgent, &session.StartedAt,
			&session.ExpiresAt, &endedAt, &session.IsActive, &endReason,
		)
		if err != nil {
			log.Printf("[ImpersonationRepository] ERROR: Failed to scan session: %v\n", err)
			return nil, 0, fmt.Errorf("failed to scan session: %w", err)
		}

		if endedAt.Valid {
			session.EndedAt = &endedAt.Time
		}
		if endReason.Valid {
			session.EndReason = &endReason.String
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[ImpersonationRepository] ERROR: Row iteration error: %v\n", err)
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return sessions, total, nil
}

// ExpireOldSessions expires sessions that have passed their expiry time
func (r *ImpersonationRepository) ExpireOldSessions(ctx context.Context) (int64, error) {
	log.Printf("[ImpersonationRepository] ExpireOldSessions\n")

	result, err := r.db.ExecContext(ctx, `
		UPDATE impersonation_sessions
		SET is_active = false, ended_at = NOW(), end_reason = $1
		WHERE is_active = true AND expires_at < NOW()
	`, domain.ImpersonationEndExpired)

	if err != nil {
		log.Printf("[ImpersonationRepository] ERROR: Failed to expire sessions: %v\n", err)
		return 0, fmt.Errorf("failed to expire sessions: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("[ImpersonationRepository] Expired %d sessions\n", rowsAffected)
	return rowsAffected, nil
}
