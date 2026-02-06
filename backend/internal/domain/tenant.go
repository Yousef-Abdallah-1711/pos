package domain

import "time"

type Tenant struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Domain           string    `json:"domain"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	LogoURL          string    `json:"logo_url"`
	SubscriptionPlan string    `json:"subscription_plan"`
	MaxRestaurants   int       `json:"max_restaurants"`
	MaxUsers         int       `json:"max_users"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Subscription and trial fields
	SubscriptionStatus  string     `json:"subscription_status"`
	TrialStartDate      *time.Time `json:"trial_start_date,omitempty"`
	TrialEndDate        *time.Time `json:"trial_end_date,omitempty"`
	TrialDays           int        `json:"trial_days"`
	IsDisabled          bool       `json:"is_disabled"`
	DisabledAt          *time.Time `json:"disabled_at,omitempty"`
	DisabledBy          *int64     `json:"disabled_by,omitempty"`
	DisabledReason      *string    `json:"disabled_reason,omitempty"`
	EnabledAt           *time.Time `json:"enabled_at,omitempty"`
	GracePeriodDays     int        `json:"grace_period_days"`
	DeletionScheduledAt *time.Time `json:"deletion_scheduled_at,omitempty"`
}

// IsTrialActive returns true if the trial is currently active
func (t *Tenant) IsTrialActive() bool {
	if t.TrialEndDate == nil {
		return false
	}
	return time.Now().Before(*t.TrialEndDate)
}

// IsTrialExpired returns true if the trial has expired
func (t *Tenant) IsTrialExpired() bool {
	if t.TrialEndDate == nil {
		return true
	}
	return time.Now().After(*t.TrialEndDate)
}

// DaysRemainingInTrial calculates the number of remaining trial days
func (t *Tenant) DaysRemainingInTrial() int {
	if t.TrialEndDate == nil {
		return 0
	}
	if time.Now().After(*t.TrialEndDate) {
		return 0
	}
	days := t.TrialEndDate.Sub(time.Now()).Hours() / 24
	return int(days) + 1 // Round up to include current day
}

// CanAccess returns true if the tenant can access the system
func (t *Tenant) CanAccess() bool {
	// Cannot access if manually disabled
	if t.IsDisabled {
		return false
	}

	// Cannot access if trial expired and not a paid subscription
	if t.SubscriptionStatus == "trial" && t.IsTrialExpired() {
		return false
	}

	// Can access for all other statuses (active, expired with payment plan, etc)
	return true
}

// GetBlockReason returns the reason why tenant cannot access, or empty string if can access
func (t *Tenant) GetBlockReason() string {
	if t.IsDisabled {
		return "tenant_disabled"
	}

	if t.SubscriptionStatus == "trial" && t.IsTrialExpired() {
		return "trial_expired"
	}

	return ""
}
