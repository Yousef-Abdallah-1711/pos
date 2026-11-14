package domain

import "time"

// TenantAccessState applies the same entitlement policy to public and authenticated APIs.
type TenantAccessState struct {
	Active          bool
	Disabled        bool
	Status          string
	TrialEnd        *time.Time
	SubscriptionEnd *time.Time
}

func (s TenantAccessState) BlockReason(now time.Time) string {
	if !s.Active || s.Disabled {
		return "tenant_disabled"
	}
	switch s.Status {
	case "trial":
		if s.TrialEnd == nil || !now.Before(*s.TrialEnd) {
			return "trial_expired"
		}
	case "active":
		if s.SubscriptionEnd != nil && !now.Before(*s.SubscriptionEnd) {
			return "subscription_expired"
		}
	default:
		return "subscription_expired"
	}
	return ""
}
