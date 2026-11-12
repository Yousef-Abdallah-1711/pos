package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

// ImpersonationUseCase handles user impersonation logic
type ImpersonationUseCase struct {
	impersonationRepo *repository.ImpersonationRepository
	superAdminRepo    *repository.SuperAdminRepository
	auditRepo         *repository.SuperAdminAuditRepository
	tokenExpiry       time.Duration
}

// ErrImpersonationAudit marks failures that must prevent an impersonation flow from reporting success.
var ErrImpersonationAudit = errors.New("impersonation audit write failed")

// NewImpersonationUseCase creates a new impersonation usecase
func NewImpersonationUseCase(
	impersonationRepo *repository.ImpersonationRepository,
	superAdminRepo *repository.SuperAdminRepository,
	auditRepo *repository.SuperAdminAuditRepository,
) *ImpersonationUseCase {
	return &ImpersonationUseCase{
		impersonationRepo: impersonationRepo,
		superAdminRepo:    superAdminRepo,
		auditRepo:         auditRepo,
		tokenExpiry:       1 * time.Hour, // Default 1 hour expiry
	}
}

// StartImpersonation starts an impersonation session
func (uc *ImpersonationUseCase) StartImpersonation(
	ctx context.Context,
	superAdminID int64,
	targetUserID int64,
	ip, userAgent string,
) (*domain.ImpersonationResponse, error) {
	log.Printf("[ImpersonationUseCase] StartImpersonation: superAdminID=%d, targetUserID=%d\n", superAdminID, targetUserID)

	// Check if super admin has an active session - only one allowed
	existingSession, err := uc.impersonationRepo.GetActiveSession(ctx, superAdminID)
	if err != nil {
		log.Printf("[ImpersonationUseCase] Error checking existing session: %v\n", err)
		return nil, fmt.Errorf("failed to check existing sessions: %w", err)
	}

	if existingSession != nil {
		log.Printf("[ImpersonationUseCase] Super admin already has active session: %d\n", existingSession.ID)
		return nil, errors.New("you already have an active impersonation session. Please end it first")
	}

	// Get target user details
	targetUser, err := uc.superAdminRepo.GetUserByID(ctx, targetUserID)
	if err != nil {
		log.Printf("[ImpersonationUseCase] Target user not found: %d\n", targetUserID)
		return nil, errors.New("target user not found")
	}

	// Cannot impersonate another super admin
	if targetUser.IsSuperAdmin {
		log.Printf("[ImpersonationUseCase] Cannot impersonate super admin: %d\n", targetUserID)
		return nil, errors.New("cannot impersonate another super admin")
	}

	// Generate secure impersonation token
	token, err := GenerateSecureToken(64) // 128 hex characters
	if err != nil {
		log.Printf("[ImpersonationUseCase] Failed to generate token: %v\n", err)
		return nil, fmt.Errorf("failed to generate impersonation token: %w", err)
	}

	// Create session
	expiresAt := time.Now().Add(uc.tokenExpiry)
	session := &domain.ImpersonationSession{
		SuperAdminID:       superAdminID,
		TargetUserID:       targetUserID,
		TargetTenantID:     int64(targetUser.TenantID),
		ImpersonationToken: token,
		IPAddress:          ip,
		UserAgent:          userAgent,
		ExpiresAt:          expiresAt,
	}

	createdSession, err := uc.impersonationRepo.CreateSession(ctx, session)
	if err != nil {
		if errors.Is(err, repository.ErrActiveImpersonationSession) {
			return nil, errors.New("you already have an active impersonation session. Please end it first")
		}
		log.Printf("[ImpersonationUseCase] Failed to create session: %v\n", err)
		return nil, fmt.Errorf("failed to create impersonation session: %w", err)
	}

	// Log the action
	targetType := domain.AuditTargetImpersonation
	targetID := targetUserID
	if err := uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionImpersonateStart,
		TargetType:   &targetType,
		TargetID:     &targetID,
		Details: map[string]interface{}{
			"target_user_email": targetUser.Email,
			"target_tenant_id":  targetUser.TenantID,
			"session_id":        createdSession.ID,
			"expires_at":        expiresAt.Format(time.RFC3339),
		},
		IPAddress: ip,
		UserAgent: userAgent,
	}); err != nil {
		cleanupErr := uc.impersonationRepo.EndSession(ctx, createdSession.ID, domain.ImpersonationEndForced)
		if cleanupErr != nil {
			return nil, fmt.Errorf("%w: start audit: %v; failed to revoke unaudited session: %w", ErrImpersonationAudit, err, cleanupErr)
		}
		return nil, fmt.Errorf("%w: start audit: %v", ErrImpersonationAudit, err)
	}

	log.Printf("[ImpersonationUseCase] Impersonation started: sessionID=%d\n", createdSession.ID)

	// Build dashboard URL
	var dashboardURL string
	if targetUser.RestaurantID != nil {
		dashboardURL = fmt.Sprintf("/en/dashboard?impersonate=%s", token)
	} else {
		dashboardURL = fmt.Sprintf("/en/dashboard?impersonate=%s", token)
	}

	var restaurantID *int64
	if targetUser.RestaurantID != nil {
		rid := int64(*targetUser.RestaurantID)
		restaurantID = &rid
	}

	return &domain.ImpersonationResponse{
		Token:              token,
		TargetUser:         targetUser,
		TargetTenantID:     int64(targetUser.TenantID),
		TargetRestaurantID: restaurantID,
		ExpiresAt:          expiresAt,
		DashboardURL:       dashboardURL,
	}, nil
}

// ValidateImpersonation validates an impersonation token
func (uc *ImpersonationUseCase) ValidateImpersonation(ctx context.Context, token string) (*domain.ImpersonationSession, *domain.User, error) {
	log.Printf("[ImpersonationUseCase] ValidateImpersonation\n")

	session, err := uc.impersonationRepo.ValidateToken(ctx, token)
	if err != nil {
		log.Printf("[ImpersonationUseCase] Invalid token: %v\n", err)
		return nil, nil, err
	}

	// Get target user
	targetUser, err := uc.superAdminRepo.GetUserByID(ctx, session.TargetUserID)
	if err != nil {
		log.Printf("[ImpersonationUseCase] Failed to get target user: %v\n", err)
		return nil, nil, fmt.Errorf("failed to get target user: %w", err)
	}

	return session, targetUser, nil
}

// EndImpersonation ends an active impersonation session
func (uc *ImpersonationUseCase) EndImpersonation(ctx context.Context, superAdminID int64, ip, userAgent string) error {
	log.Printf("[ImpersonationUseCase] EndImpersonation: superAdminID=%d\n", superAdminID)

	// Get active session
	session, err := uc.impersonationRepo.GetActiveSession(ctx, superAdminID)
	if err != nil {
		log.Printf("[ImpersonationUseCase] Error getting active session: %v\n", err)
		return fmt.Errorf("failed to get active session: %w", err)
	}

	if session == nil {
		log.Printf("[ImpersonationUseCase] No active session found\n")
		return errors.New("no active impersonation session")
	}

	// End the session
	err = uc.impersonationRepo.EndSession(ctx, session.ID, domain.ImpersonationEndManual)
	if err != nil {
		log.Printf("[ImpersonationUseCase] Failed to end session: %v\n", err)
		return fmt.Errorf("failed to end session: %w", err)
	}

	// Log the action
	targetType := domain.AuditTargetImpersonation
	sessionID := session.ID
	if err := uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionImpersonateEnd,
		TargetType:   &targetType,
		TargetID:     &sessionID,
		Details: map[string]interface{}{
			"target_user_id":   session.TargetUserID,
			"target_tenant_id": session.TargetTenantID,
			"duration_minutes": time.Since(session.StartedAt).Minutes(),
		},
		IPAddress: ip,
		UserAgent: userAgent,
	}); err != nil {
		return fmt.Errorf("%w: end audit: %v", ErrImpersonationAudit, err)
	}

	log.Printf("[ImpersonationUseCase] Impersonation ended: sessionID=%d\n", session.ID)
	return nil
}

// EndImpersonationByToken ends an impersonation session by token
func (uc *ImpersonationUseCase) EndImpersonationByToken(ctx context.Context, token, ip, userAgent string) error {
	log.Printf("[ImpersonationUseCase] EndImpersonationByToken\n")

	session, err := uc.impersonationRepo.ValidateToken(ctx, token)
	if err != nil {
		log.Printf("[ImpersonationUseCase] Invalid token: %v\n", err)
		return errors.New("invalid impersonation token")
	}

	// End the session
	err = uc.impersonationRepo.EndSession(ctx, session.ID, domain.ImpersonationEndManual)
	if err != nil {
		log.Printf("[ImpersonationUseCase] Failed to end session: %v\n", err)
		return fmt.Errorf("failed to end session: %w", err)
	}

	// Log the action
	targetType := domain.AuditTargetImpersonation
	sessionID := session.ID
	if err := uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: session.SuperAdminID,
		Action:       domain.AuditActionImpersonateEnd,
		TargetType:   &targetType,
		TargetID:     &sessionID,
		Details: map[string]interface{}{
			"target_user_id":   session.TargetUserID,
			"target_tenant_id": session.TargetTenantID,
			"duration_minutes": time.Since(session.StartedAt).Minutes(),
			"ended_by":         "token",
		},
		IPAddress: ip,
		UserAgent: userAgent,
	}); err != nil {
		return fmt.Errorf("%w: end audit: %v", ErrImpersonationAudit, err)
	}

	log.Printf("[ImpersonationUseCase] Impersonation ended by token: sessionID=%d\n", session.ID)
	return nil
}

// GetImpersonationStatus returns the current impersonation status for a super admin
func (uc *ImpersonationUseCase) GetImpersonationStatus(ctx context.Context, superAdminID int64) (*domain.ImpersonationSession, error) {
	log.Printf("[ImpersonationUseCase] GetImpersonationStatus: superAdminID=%d\n", superAdminID)

	session, err := uc.impersonationRepo.GetActiveSession(ctx, superAdminID)
	if err != nil {
		log.Printf("[ImpersonationUseCase] Error getting session: %v\n", err)
		return nil, fmt.Errorf("failed to get impersonation status: %w", err)
	}

	return session, nil
}

// ForceEndAllSessions ends all active sessions for a super admin (used during logout)
func (uc *ImpersonationUseCase) ForceEndAllSessions(ctx context.Context, superAdminID int64) error {
	log.Printf("[ImpersonationUseCase] ForceEndAllSessions: superAdminID=%d\n", superAdminID)

	return uc.impersonationRepo.EndAllSessionsForSuperAdmin(ctx, superAdminID, domain.ImpersonationEndForced)
}

// CleanupExpiredSessions cleans up expired sessions (to be called periodically)
func (uc *ImpersonationUseCase) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	log.Printf("[ImpersonationUseCase] CleanupExpiredSessions\n")

	count, err := uc.impersonationRepo.ExpireOldSessions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup sessions: %w", err)
	}

	log.Printf("[ImpersonationUseCase] Cleaned up %d expired sessions\n", count)
	return count, nil
}
