package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"pos-saas/internal/domain"
	"pos-saas/internal/pkg/jwt"
	"pos-saas/internal/repository"
)

// SuperAdminTokenService interface for generating super admin tokens
type SuperAdminTokenService interface {
	GenerateSuperAdminToken(userID int64, email string, superAdminID int64) (string, time.Time, error)
	ValidateSuperAdminToken(token string) (*jwt.SuperAdminClaims, error)
}

// SuperAdminAuthUseCase handles super admin authentication
type SuperAdminAuthUseCase struct {
	superAdminRepo *repository.SuperAdminRepository
	userRepo       *repository.UserRepository
	auditRepo      *repository.SuperAdminAuditRepository
	tokenService   SuperAdminTokenService
}

// NewSuperAdminAuthUseCase creates a new super admin auth usecase
func NewSuperAdminAuthUseCase(
	superAdminRepo *repository.SuperAdminRepository,
	userRepo *repository.UserRepository,
	auditRepo *repository.SuperAdminAuditRepository,
	tokenService SuperAdminTokenService,
) *SuperAdminAuthUseCase {
	return &SuperAdminAuthUseCase{
		superAdminRepo: superAdminRepo,
		userRepo:       userRepo,
		auditRepo:      auditRepo,
		tokenService:   tokenService,
	}
}

// Login authenticates a super admin
func (uc *SuperAdminAuthUseCase) Login(ctx context.Context, req *domain.SuperAdminLoginRequest, ip, userAgent string) (*domain.SuperAdminLoginResponse, error) {
	log.Printf("[SuperAdminAuthUseCase] Login attempt for email: %s from IP: %s\n", req.Email, ip)

	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Find user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("[SuperAdminAuthUseCase] User not found: %s\n", req.Email)
		uc.logFailedLogin(ctx, 0, ip, userAgent, "user_not_found")
		return nil, errors.New("invalid credentials")
	}

	// Check if user is a super admin
	if !user.IsSuperAdmin {
		log.Printf("[SuperAdminAuthUseCase] User is not a super admin: %s\n", req.Email)
		uc.logFailedLogin(ctx, 0, ip, userAgent, "not_super_admin")
		return nil, errors.New("invalid credentials")
	}

	// Get super admin record
	superAdmin, err := uc.superAdminRepo.GetByUserID(ctx, int64(user.ID))
	if err != nil {
		log.Printf("[SuperAdminAuthUseCase] Super admin record not found for user: %d\n", user.ID)
		uc.logFailedLogin(ctx, 0, ip, userAgent, "super_admin_record_not_found")
		return nil, errors.New("invalid credentials")
	}

	// Check if super admin is active
	if !superAdmin.IsActive {
		log.Printf("[SuperAdminAuthUseCase] Super admin account is inactive: %d\n", superAdmin.ID)
		uc.logFailedLogin(ctx, superAdmin.ID, ip, userAgent, "account_inactive")
		return nil, errors.New("account is inactive")
	}

	// Check if account is locked
	if superAdmin.IsLocked() {
		log.Printf("[SuperAdminAuthUseCase] Super admin account is locked: %d\n", superAdmin.ID)
		uc.logFailedLogin(ctx, superAdmin.ID, ip, userAgent, "account_locked")
		return nil, fmt.Errorf("account is locked until %s", superAdmin.LockedUntil.Format(time.RFC3339))
	}

	// Check IP whitelist
	if !superAdmin.IsIPAllowed(ip) {
		log.Printf("[SuperAdminAuthUseCase] IP not in whitelist: %s for super admin: %d\n", ip, superAdmin.ID)
		uc.logFailedLogin(ctx, superAdmin.ID, ip, userAgent, "ip_not_allowed")
		return nil, errors.New("access denied from this IP address")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Printf("[SuperAdminAuthUseCase] Invalid password for super admin: %d\n", superAdmin.ID)

		// Update failed login attempts
		_ = uc.superAdminRepo.UpdateLoginFailure(ctx, superAdmin.ID)
		uc.logFailedLogin(ctx, superAdmin.ID, ip, userAgent, "invalid_password")

		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, expiresAt, err := uc.tokenService.GenerateSuperAdminToken(int64(user.ID), user.Email, superAdmin.ID)
	if err != nil {
		log.Printf("[SuperAdminAuthUseCase] Failed to generate token: %v\n", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update login success
	err = uc.superAdminRepo.UpdateLoginSuccess(ctx, superAdmin.ID, ip)
	if err != nil {
		log.Printf("[SuperAdminAuthUseCase] Warning: Failed to update login success: %v\n", err)
	}

	// Log successful login
	uc.logSuccessfulLogin(ctx, superAdmin.ID, ip, userAgent)

	// Attach user to super admin for response
	user.PasswordHash = "" // Clear sensitive data
	superAdmin.User = user

	log.Printf("[SuperAdminAuthUseCase] Login successful for super admin: %d\n", superAdmin.ID)

	return &domain.SuperAdminLoginResponse{
		SuperAdmin: superAdmin,
		User:       user,
		Token:      token,
		ExpiresAt:  expiresAt,
	}, nil
}

// Logout logs out a super admin
func (uc *SuperAdminAuthUseCase) Logout(ctx context.Context, superAdminID int64, ip, userAgent string) error {
	log.Printf("[SuperAdminAuthUseCase] Logout for super admin: %d\n", superAdminID)

	// Log logout action
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionLogout,
		IPAddress:    ip,
		UserAgent:    userAgent,
	})

	return nil
}

// ValidateSession validates a super admin session
func (uc *SuperAdminAuthUseCase) ValidateSession(ctx context.Context, superAdminID int64, ip string) (*domain.SuperAdmin, error) {
	log.Printf("[SuperAdminAuthUseCase] ValidateSession for super admin: %d\n", superAdminID)

	superAdmin, err := uc.superAdminRepo.GetByID(ctx, superAdminID)
	if err != nil {
		return nil, errors.New("session invalid")
	}

	if !superAdmin.IsActive {
		return nil, errors.New("account is inactive")
	}

	if superAdmin.IsLocked() {
		return nil, errors.New("account is locked")
	}

	if !superAdmin.IsIPAllowed(ip) {
		return nil, errors.New("access denied from this IP address")
	}

	return superAdmin, nil
}

// UpdateIPWhitelist updates the IP whitelist for a super admin
func (uc *SuperAdminAuthUseCase) UpdateIPWhitelist(ctx context.Context, superAdminID int64, ipWhitelist []string, ip, userAgent string) error {
	log.Printf("[SuperAdminAuthUseCase] UpdateIPWhitelist for super admin: %d\n", superAdminID)

	err := uc.superAdminRepo.UpdateIPWhitelist(ctx, superAdminID, ipWhitelist)
	if err != nil {
		return err
	}

	// Log the action
	targetType := domain.AuditTargetSystem
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionUpdateIPWhitelist,
		TargetType:   &targetType,
		Details:      map[string]interface{}{"ip_whitelist": ipWhitelist},
		IPAddress:    ip,
		UserAgent:    userAgent,
	})

	return nil
}

// GenerateSecureToken generates a cryptographically secure token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Helper methods for audit logging
func (uc *SuperAdminAuthUseCase) logFailedLogin(ctx context.Context, superAdminID int64, ip, userAgent, reason string) {
	if superAdminID == 0 {
		return // Can't log without a super admin ID
	}
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionLoginFailed,
		Details:      map[string]interface{}{"reason": reason},
		IPAddress:    ip,
		UserAgent:    userAgent,
	})
}

func (uc *SuperAdminAuthUseCase) logSuccessfulLogin(ctx context.Context, superAdminID int64, ip, userAgent string) {
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionLoginSuccess,
		IPAddress:    ip,
		UserAgent:    userAgent,
	})
}
