package domain

import (
	"time"
)

// Subscription status constants
const (
	StatusTrial     = "trial"
	StatusActive    = "active"
	StatusExpired   = "expired"
	StatusSuspended = "suspended"
	StatusDisabled  = "disabled"
)

// Subscription action constants
const (
	ActionTrialStarted      = "TRIAL_STARTED"
	ActionTrialExtended     = "TRIAL_EXTENDED"
	ActionTrialExpired      = "TRIAL_EXPIRED"
	ActionTenantDisabled    = "TENANT_DISABLED"
	ActionTenantEnabled     = "TENANT_ENABLED"
	ActionDeletionScheduled = "DELETION_SCHEDULED"
	ActionDeletionCancelled = "DELETION_CANCELLED"
	ActionTenantDeleted     = "TENANT_DELETED"
	ActionStatusChanged     = "STATUS_CHANGED"
)

// SubscriptionHistory represents a subscription change event
type SubscriptionHistory struct {
	ID              int64                  `json:"id"`
	TenantID        int64                  `json:"tenant_id"`
	Action          string                 `json:"action"`
	PerformedBy     *int64                 `json:"performed_by,omitempty"`
	OldStatus       *string                `json:"old_status,omitempty"`
	NewStatus       *string                `json:"new_status,omitempty"`
	OldTrialEndDate *time.Time             `json:"old_trial_end_date,omitempty"`
	NewTrialEndDate *time.Time             `json:"new_trial_end_date,omitempty"`
	DaysAdded       *int                   `json:"days_added,omitempty"`
	Reason          *string                `json:"reason,omitempty"`
	IPAddress       *string                `json:"ip_address,omitempty"`
	UserAgent       *string                `json:"user_agent,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

// TenantSubscriptionInfo holds comprehensive subscription information for a tenant
type TenantSubscriptionInfo struct {
	IsActive              bool                     `json:"is_active"`
	TenantID              int64                    `json:"tenant_id"`
	TenantName            string                   `json:"tenant_name"`
	Slug                  string                   `json:"slug"`
	SubscriptionStatus    string                   `json:"subscription_status"`
	SubscriptionPlan      string                   `json:"subscription_plan,omitempty"`
	TrialStartDate        *time.Time               `json:"trial_start_date,omitempty"`
	TrialEndDate          *time.Time               `json:"trial_end_date,omitempty"`
	SubscriptionExpiresAt *time.Time               `json:"subscription_expires_at"`
	LastPaymentAt         *time.Time               `json:"last_payment_at"`
	TrialDays             int                      `json:"trial_days"`
	DaysRemaining         int                      `json:"days_remaining"`
	IsDisabled            bool                     `json:"is_disabled"`
	DisabledAt            *time.Time               `json:"disabled_at,omitempty"`
	DisabledReason        *string                  `json:"disabled_reason,omitempty"`
	CanAccess             bool                     `json:"can_access"`
	BlockReason           *string                  `json:"block_reason,omitempty"`
	Plan                  *SubscriptionPlanSummary `json:"plan,omitempty"`
	ContactEmail          string                   `json:"contact_email,omitempty"`
	ContactPhone          string                   `json:"contact_phone,omitempty"`
}

// CreateHistoryRecordParams holds parameters for creating a subscription history record
type CreateHistoryRecordParams struct {
	TenantID        int64
	Action          string
	PerformedBy     *int64
	OldStatus       *string
	NewStatus       *string
	OldTrialEndDate *time.Time
	NewTrialEndDate *time.Time
	DaysAdded       *int
	Reason          *string
	IPAddress       *string
	UserAgent       *string
	Metadata        map[string]interface{}
}

// ExtendTrialParams holds parameters for extending trial
type ExtendTrialParams struct {
	TenantID    int64
	Days        int
	PerformedBy int64
	Reason      *string
	IPAddress   string
	UserAgent   string
}

// DisableTenantParams holds parameters for disabling a tenant
type DisableTenantParams struct {
	TenantID    int64
	PerformedBy int64
	Reason      string
	IPAddress   string
	UserAgent   string
}

// EnableTenantParams holds parameters for enabling a tenant
type EnableTenantParams struct {
	TenantID    int64
	PerformedBy int64
	IPAddress   string
	UserAgent   string
}

// DeleteTenantParams holds parameters for deleting a tenant
type DeleteTenantParams struct {
	TenantID    int64
	PerformedBy int64
	Reason      *string
	IPAddress   string
	UserAgent   string
}

// ValidateSubscriptionStatus checks if a status string is valid
func ValidateSubscriptionStatus(status string) bool {
	switch status {
	case StatusTrial, StatusActive, StatusExpired, StatusSuspended, StatusDisabled:
		return true
	}
	return false
}

// ValidateSubscriptionAction checks if an action string is valid
func ValidateSubscriptionAction(action string) bool {
	switch action {
	case ActionTrialStarted, ActionTrialExtended, ActionTrialExpired,
		ActionTenantDisabled, ActionTenantEnabled, ActionDeletionScheduled,
		ActionDeletionCancelled, ActionTenantDeleted, ActionStatusChanged,
		ActionPaymentInitiated, ActionPaymentSuccess, ActionPaymentFailed,
		ActionSubscriptionActivated, ActionSubscriptionExpired,
		ActionPaymentRefunded, ActionSubscriptionCancelled:
		return true
	}
	return false
}
