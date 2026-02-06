package usecase

import (
	"context"
	"fmt"
	"log"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

// SuperAdminTenantUseCase handles tenant management for super admins
type SuperAdminTenantUseCase struct {
	superAdminRepo *repository.SuperAdminRepository
	auditRepo      *repository.SuperAdminAuditRepository
}

// NewSuperAdminTenantUseCase creates a new super admin tenant usecase
func NewSuperAdminTenantUseCase(
	superAdminRepo *repository.SuperAdminRepository,
	auditRepo *repository.SuperAdminAuditRepository,
) *SuperAdminTenantUseCase {
	return &SuperAdminTenantUseCase{
		superAdminRepo: superAdminRepo,
		auditRepo:      auditRepo,
	}
}

// GetDashboardStats retrieves statistics for the super admin dashboard
func (uc *SuperAdminTenantUseCase) GetDashboardStats(ctx context.Context, superAdminID int64, ip, userAgent string) (*domain.SuperAdminDashboardStats, error) {
	log.Printf("[SuperAdminTenantUseCase] GetDashboardStats: superAdminID=%d\n", superAdminID)

	stats, err := uc.superAdminRepo.GetDashboardStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	return stats, nil
}

// ListTenants retrieves all tenants with filtering and pagination
func (uc *SuperAdminTenantUseCase) ListTenants(
	ctx context.Context,
	superAdminID int64,
	filter domain.TenantFilter,
	pagination domain.Pagination,
	ip, userAgent string,
) (*domain.PaginatedResponse, error) {
	log.Printf("[SuperAdminTenantUseCase] ListTenants: superAdminID=%d, page=%d\n", superAdminID, pagination.Page)

	// Set defaults for pagination
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 20
	}

	tenants, total, err := uc.superAdminRepo.ListAllTenants(ctx, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	// Log the action
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionViewTenantList,
		Details: map[string]interface{}{
			"page":      pagination.Page,
			"page_size": pagination.PageSize,
			"total":     total,
		},
		IPAddress: ip,
		UserAgent: userAgent,
	})

	return domain.NewPaginatedResponse(tenants, total, pagination.Page, pagination.PageSize), nil
}

// GetTenantDetails retrieves detailed information about a specific tenant
func (uc *SuperAdminTenantUseCase) GetTenantDetails(
	ctx context.Context,
	superAdminID int64,
	tenantID int64,
	ip, userAgent string,
) (*domain.TenantWithStats, error) {
	log.Printf("[SuperAdminTenantUseCase] GetTenantDetails: superAdminID=%d, tenantID=%d\n", superAdminID, tenantID)

	tenant, err := uc.superAdminRepo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	// Log the action
	targetType := domain.AuditTargetTenant
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionViewTenant,
		TargetType:   &targetType,
		TargetID:     &tenantID,
		Details: map[string]interface{}{
			"tenant_name": tenant.Name,
			"tenant_slug": tenant.Slug,
		},
		IPAddress: ip,
		UserAgent: userAgent,
	})

	return tenant, nil
}

// GetTenantUsers retrieves all users for a specific tenant
func (uc *SuperAdminTenantUseCase) GetTenantUsers(
	ctx context.Context,
	superAdminID int64,
	tenantID int64,
	pagination domain.Pagination,
	ip, userAgent string,
) (*domain.PaginatedResponse, error) {
	log.Printf("[SuperAdminTenantUseCase] GetTenantUsers: superAdminID=%d, tenantID=%d\n", superAdminID, tenantID)

	// Set defaults for pagination
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 20
	}

	filter := domain.UserFilter{
		TenantID: &tenantID,
	}

	users, total, err := uc.superAdminRepo.ListAllUsers(ctx, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return domain.NewPaginatedResponse(users, total, pagination.Page, pagination.PageSize), nil
}

// ListAllUsers retrieves all users across all tenants
func (uc *SuperAdminTenantUseCase) ListAllUsers(
	ctx context.Context,
	superAdminID int64,
	filter domain.UserFilter,
	pagination domain.Pagination,
	ip, userAgent string,
) (*domain.PaginatedResponse, error) {
	log.Printf("[SuperAdminTenantUseCase] ListAllUsers: superAdminID=%d, page=%d\n", superAdminID, pagination.Page)

	// Set defaults for pagination
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 20
	}

	users, total, err := uc.superAdminRepo.ListAllUsers(ctx, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Log the action
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionViewUserList,
		Details: map[string]interface{}{
			"page":      pagination.Page,
			"page_size": pagination.PageSize,
			"total":     total,
		},
		IPAddress: ip,
		UserAgent: userAgent,
	})

	return domain.NewPaginatedResponse(users, total, pagination.Page, pagination.PageSize), nil
}

// GetUserDetails retrieves detailed information about a specific user
func (uc *SuperAdminTenantUseCase) GetUserDetails(
	ctx context.Context,
	superAdminID int64,
	userID int64,
	ip, userAgent string,
) (*domain.User, error) {
	log.Printf("[SuperAdminTenantUseCase] GetUserDetails: superAdminID=%d, userID=%d\n", superAdminID, userID)

	user, err := uc.superAdminRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Log the action
	targetType := domain.AuditTargetUser
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionViewUser,
		TargetType:   &targetType,
		TargetID:     &userID,
		Details: map[string]interface{}{
			"user_email":  user.Email,
			"user_name":   user.Name,
			"user_tenant": user.TenantID,
		},
		IPAddress: ip,
		UserAgent: userAgent,
	})

	// Clear sensitive data
	user.PasswordHash = ""

	return user, nil
}

// DisableTenant disables a tenant
func (uc *SuperAdminTenantUseCase) DisableTenant(
	ctx context.Context,
	superAdminID int64,
	tenantID int64,
	reason string,
	ip, userAgent string,
) error {
	log.Printf("[SuperAdminTenantUseCase] DisableTenant: superAdminID=%d, tenantID=%d\n", superAdminID, tenantID)

	// Get tenant first for logging
	tenant, err := uc.superAdminRepo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	// Update tenant status
	err = uc.superAdminRepo.UpdateTenantStatus(ctx, tenantID, false)
	if err != nil {
		return fmt.Errorf("failed to disable tenant: %w", err)
	}

	// Log the action
	targetType := domain.AuditTargetTenant
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionDisableTenant,
		TargetType:   &targetType,
		TargetID:     &tenantID,
		Details: map[string]interface{}{
			"tenant_name": tenant.Name,
			"tenant_slug": tenant.Slug,
			"reason":      reason,
		},
		IPAddress: ip,
		UserAgent: userAgent,
	})

	log.Printf("[SuperAdminTenantUseCase] Tenant disabled: %d (%s)\n", tenantID, tenant.Name)
	return nil
}

// EnableTenant enables a tenant
func (uc *SuperAdminTenantUseCase) EnableTenant(
	ctx context.Context,
	superAdminID int64,
	tenantID int64,
	ip, userAgent string,
) error {
	log.Printf("[SuperAdminTenantUseCase] EnableTenant: superAdminID=%d, tenantID=%d\n", superAdminID, tenantID)

	// Get tenant first for logging
	tenant, err := uc.superAdminRepo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	// Update tenant status
	err = uc.superAdminRepo.UpdateTenantStatus(ctx, tenantID, true)
	if err != nil {
		return fmt.Errorf("failed to enable tenant: %w", err)
	}

	// Log the action
	targetType := domain.AuditTargetTenant
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionEnableTenant,
		TargetType:   &targetType,
		TargetID:     &tenantID,
		Details: map[string]interface{}{
			"tenant_name": tenant.Name,
			"tenant_slug": tenant.Slug,
		},
		IPAddress: ip,
		UserAgent: userAgent,
	})

	log.Printf("[SuperAdminTenantUseCase] Tenant enabled: %d (%s)\n", tenantID, tenant.Name)
	return nil
}

// GetAuditLogs retrieves audit logs with filtering
func (uc *SuperAdminTenantUseCase) GetAuditLogs(
	ctx context.Context,
	superAdminID int64,
	filter domain.AuditLogFilter,
	pagination domain.Pagination,
	ip, userAgent string,
) (*domain.PaginatedResponse, error) {
	log.Printf("[SuperAdminTenantUseCase] GetAuditLogs: superAdminID=%d\n", superAdminID)

	// Set defaults for pagination
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 50
	}

	logs, total, err := uc.auditRepo.ListLogs(ctx, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}

	// Log viewing audit logs
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionViewAuditLogs,
		Details: map[string]interface{}{
			"page":      pagination.Page,
			"page_size": pagination.PageSize,
			"total":     total,
		},
		IPAddress: ip,
		UserAgent: userAgent,
	})

	return domain.NewPaginatedResponse(logs, total, pagination.Page, pagination.PageSize), nil
}

// GenerateDashboardLoginToken generates a one-time token for accessing the tenant dashboard
func (uc *SuperAdminTenantUseCase) GenerateDashboardLoginToken(
	ctx context.Context,
	superAdminID int64,
	tenantID int64,
	ip, userAgent string,
) (string, error) {
	log.Printf("[SuperAdminTenantUseCase] GenerateDashboardLoginToken: superAdminID=%d, tenantID=%d\n", superAdminID, tenantID)

	// Verify tenant exists
	tenant, err := uc.superAdminRepo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to verify tenant: %w", err)
	}

	// Get first admin/manager user from tenant for impersonation
	// For now, we'll use a generic token that the dashboard will resolve
	// The actual user selection will happen on the dashboard side
	token, err := uc.superAdminRepo.GenerateDashboardLoginToken(ctx, superAdminID, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Log the action
	targetType := domain.AuditTargetTenant
	_ = uc.auditRepo.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionAccessDashboard,
		TargetType:   &targetType,
		TargetID:     &tenantID,
		Details: map[string]interface{}{
			"tenant_name": tenant.Name,
			"tenant_slug": tenant.Slug,
		},
		IPAddress: ip,
		UserAgent: userAgent,
	})

	log.Printf("[SuperAdminTenantUseCase] Dashboard login token generated for tenant: %d\n", tenantID)
	return token, nil
}
