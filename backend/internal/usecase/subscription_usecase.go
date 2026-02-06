package usecase

import (
	"context"
	"fmt"
	"log"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

// SubscriptionUseCase handles subscription-related business logic
type SubscriptionUseCase struct {
	subscriptionRepo repository.SubscriptionRepository
	auditRepo        *repository.SuperAdminAuditRepository
}

// NewSubscriptionUseCase creates a new subscription usecase
func NewSubscriptionUseCase(
	subscriptionRepo repository.SubscriptionRepository,
	auditRepo *repository.SuperAdminAuditRepository,
) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		subscriptionRepo: subscriptionRepo,
		auditRepo:        auditRepo,
	}
}

// GetTenantSubscription retrieves subscription info for a tenant
func (uc *SubscriptionUseCase) GetTenantSubscription(ctx context.Context, tenantID int64) (*domain.TenantSubscriptionInfo, error) {
	log.Printf("[SubscriptionUseCase] GetTenantSubscription - tenant_id: %d\n", tenantID)

	info, err := uc.subscriptionRepo.GetTenantSubscriptionInfo(ctx, tenantID)
	if err != nil {
		log.Printf("[SubscriptionUseCase] ERROR: Failed to get subscription info - tenant_id: %d, error: %v\n", tenantID, err)
		return nil, fmt.Errorf("failed to get subscription info: %w", err)
	}

	log.Printf("[SubscriptionUseCase] GetTenantSubscription success - status: %s, can_access: %v\n", info.SubscriptionStatus, info.CanAccess)
	return info, nil
}

// GetTenantSubscriptionBySlug retrieves subscription info by tenant slug
func (uc *SubscriptionUseCase) GetTenantSubscriptionBySlug(ctx context.Context, slug string) (*domain.TenantSubscriptionInfo, error) {
	log.Printf("[SubscriptionUseCase] GetTenantSubscriptionBySlug - slug: %s\n", slug)

	info, err := uc.subscriptionRepo.GetTenantSubscriptionBySlug(ctx, slug)
	if err != nil {
		log.Printf("[SubscriptionUseCase] ERROR: Failed to get subscription by slug - slug: %s, error: %v\n", slug, err)
		return nil, fmt.Errorf("failed to get subscription info: %w", err)
	}

	return info, nil
}

// ExtendTrial extends the trial period for a tenant
func (uc *SubscriptionUseCase) ExtendTrial(
	ctx context.Context,
	tenantID int64,
	days int,
	performedBy int64,
	reason *string,
	ipAddress string,
	userAgent string,
) error {
	log.Printf("[SubscriptionUseCase] ExtendTrial - tenant_id: %d, days: %d\n", tenantID, days)

	if days <= 0 {
		log.Printf("[SubscriptionUseCase] ERROR: Invalid days - days: %d\n", days)
		return fmt.Errorf("days must be greater than 0")
	}

	params := &domain.ExtendTrialParams{
		TenantID:    tenantID,
		Days:        days,
		PerformedBy: performedBy,
		Reason:      reason,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}

	if err := uc.subscriptionRepo.ExtendTrial(ctx, params); err != nil {
		log.Printf("[SubscriptionUseCase] ERROR: Failed to extend trial - error: %v\n", err)
		return fmt.Errorf("failed to extend trial: %w", err)
	}

	// Log the action in audit log
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: performedBy,
		Action:       "EXTEND_TENANT_TRIAL",
		TargetType:   stringPtr("TENANT"),
		TargetID:     &tenantID,
		Details: map[string]interface{}{
			"days_added": days,
			"reason":     reason,
		},
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})

	log.Printf("[SubscriptionUseCase] ExtendTrial success\n")
	return nil
}

// DisableTenant disables a tenant's access
func (uc *SubscriptionUseCase) DisableTenant(
	ctx context.Context,
	tenantID int64,
	performedBy int64,
	reason string,
	ipAddress string,
	userAgent string,
) error {
	log.Printf("[SubscriptionUseCase] DisableTenant - tenant_id: %d\n", tenantID)

	if reason == "" {
		log.Printf("[SubscriptionUseCase] ERROR: Reason is required\n")
		return fmt.Errorf("reason is required")
	}

	params := &domain.DisableTenantParams{
		TenantID:    tenantID,
		PerformedBy: performedBy,
		Reason:      reason,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}

	if err := uc.subscriptionRepo.DisableTenant(ctx, params); err != nil {
		log.Printf("[SubscriptionUseCase] ERROR: Failed to disable tenant - error: %v\n", err)
		return fmt.Errorf("failed to disable tenant: %w", err)
	}

	// Log the action in audit log
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: performedBy,
		Action:       "DISABLE_TENANT",
		TargetType:   stringPtr("TENANT"),
		TargetID:     &tenantID,
		Details: map[string]interface{}{
			"reason": reason,
		},
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})

	log.Printf("[SubscriptionUseCase] DisableTenant success\n")
	return nil
}

// EnableTenant re-enables a disabled tenant
func (uc *SubscriptionUseCase) EnableTenant(
	ctx context.Context,
	tenantID int64,
	performedBy int64,
	ipAddress string,
	userAgent string,
) error {
	log.Printf("[SubscriptionUseCase] EnableTenant - tenant_id: %d\n", tenantID)

	params := &domain.EnableTenantParams{
		TenantID:    tenantID,
		PerformedBy: performedBy,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}

	if err := uc.subscriptionRepo.EnableTenant(ctx, params); err != nil {
		log.Printf("[SubscriptionUseCase] ERROR: Failed to enable tenant - error: %v\n", err)
		return fmt.Errorf("failed to enable tenant: %w", err)
	}

	// Log the action in audit log
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: performedBy,
		Action:       "ENABLE_TENANT",
		TargetType:   stringPtr("TENANT"),
		TargetID:     &tenantID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	})

	log.Printf("[SubscriptionUseCase] EnableTenant success\n")
	return nil
}

// DeleteTenant permanently deletes a tenant
func (uc *SubscriptionUseCase) DeleteTenant(
	ctx context.Context,
	tenantID int64,
	performedBy int64,
	reason *string,
	ipAddress string,
	userAgent string,
) error {
	log.Printf("[SubscriptionUseCase] DeleteTenant - tenant_id: %d\n", tenantID)

	params := &domain.DeleteTenantParams{
		TenantID:    tenantID,
		PerformedBy: performedBy,
		Reason:      reason,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}

	if err := uc.subscriptionRepo.ScheduleDeletion(ctx, params); err != nil {
		log.Printf("[SubscriptionUseCase] ERROR: Failed to delete tenant - error: %v\n", err)
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	// Log the action in audit log
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: performedBy,
		Action:       "DELETE_TENANT",
		TargetType:   stringPtr("TENANT"),
		TargetID:     &tenantID,
		Details: map[string]interface{}{
			"reason": reason,
		},
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})

	log.Printf("[SubscriptionUseCase] DeleteTenant success\n")
	return nil
}

// GetSubscriptionHistory retrieves subscription change history
func (uc *SubscriptionUseCase) GetSubscriptionHistory(
	ctx context.Context,
	tenantID int64,
	page int,
	pageSize int,
) ([]domain.SubscriptionHistory, int64, error) {
	log.Printf("[SubscriptionUseCase] GetSubscriptionHistory - tenant_id: %d, page: %d\n", tenantID, page)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	history, total, err := uc.subscriptionRepo.GetSubscriptionHistory(ctx, tenantID, page, pageSize)
	if err != nil {
		log.Printf("[SubscriptionUseCase] ERROR: Failed to get history - error: %v\n", err)
		return nil, 0, fmt.Errorf("failed to get history: %w", err)
	}

	log.Printf("[SubscriptionUseCase] GetSubscriptionHistory success - count: %d, total: %d\n", len(history), total)
	return history, total, nil
}

// CheckTenantAccess checks if a tenant can access the system
func (uc *SubscriptionUseCase) CheckTenantAccess(ctx context.Context, tenantID int64) (bool, string, error) {
	log.Printf("[SubscriptionUseCase] CheckTenantAccess - tenant_id: %d\n", tenantID)

	canAccess, blockReason, err := uc.subscriptionRepo.CheckTenantAccess(ctx, tenantID)
	if err != nil {
		log.Printf("[SubscriptionUseCase] ERROR: Failed to check access - error: %v\n", err)
		return false, "", fmt.Errorf("failed to check access: %w", err)
	}

	log.Printf("[SubscriptionUseCase] CheckTenantAccess result - canAccess: %v, blockReason: %s\n", canAccess, blockReason)
	return canAccess, blockReason, nil
}

// Helper function to convert string to pointer
func stringPtr(s string) *string {
	return &s
}
