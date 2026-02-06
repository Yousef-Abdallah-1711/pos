package repository

import (
	"context"
	"pos-saas/internal/domain"
)

// SubscriptionRepository handles subscription-related database operations
type SubscriptionRepository interface {
	// GetTenantSubscriptionInfo retrieves subscription information for a tenant
	GetTenantSubscriptionInfo(ctx context.Context, tenantID int64) (*domain.TenantSubscriptionInfo, error)

	// GetTenantSubscriptionBySlug retrieves subscription info by tenant slug (for public endpoints)
	GetTenantSubscriptionBySlug(ctx context.Context, slug string) (*domain.TenantSubscriptionInfo, error)

	// UpdateSubscriptionStatus updates the subscription status of a tenant
	UpdateSubscriptionStatus(ctx context.Context, tenantID int64, status string) error

	// ExtendTrial extends the trial period for a tenant and logs the action
	ExtendTrial(ctx context.Context, params *domain.ExtendTrialParams) error

	// DisableTenant disables a tenant's access
	DisableTenant(ctx context.Context, params *domain.DisableTenantParams) error

	// EnableTenant re-enables a disabled tenant
	EnableTenant(ctx context.Context, params *domain.EnableTenantParams) error

	// ScheduleDeletion schedules a tenant for deletion
	ScheduleDeletion(ctx context.Context, params *domain.DeleteTenantParams) error

	// CancelDeletion cancels a scheduled deletion
	CancelDeletion(ctx context.Context, tenantID int64, performedBy int64) error

	// GetSubscriptionHistory retrieves paginated subscription change history
	GetSubscriptionHistory(ctx context.Context, tenantID int64, page int, pageSize int) ([]domain.SubscriptionHistory, int64, error)

	// CreateHistoryRecord creates a subscription history record
	CreateHistoryRecord(ctx context.Context, params *domain.CreateHistoryRecordParams) error

	// GetExpiredTrials retrieves all tenants with expired trials
	GetExpiredTrials(ctx context.Context) ([]domain.TenantSubscriptionInfo, error)

	// GetTenantsScheduledForDeletion retrieves all tenants scheduled for deletion
	GetTenantsScheduledForDeletion(ctx context.Context) ([]domain.TenantSubscriptionInfo, error)

	// CheckTenantAccess checks if a tenant can access the system
	CheckTenantAccess(ctx context.Context, tenantID int64) (canAccess bool, blockReason string, err error)
}
