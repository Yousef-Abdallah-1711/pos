package domain

import (
	"testing"
	"time"
)

func TestTenantAccessEntitlement(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)
	cases := []struct {
		name   string
		state  TenantAccessState
		reason string
	}{
		{"inactive", TenantAccessState{Status: "active"}, "tenant_disabled"},
		{"disabled", TenantAccessState{Active: true, Disabled: true, Status: "active"}, "tenant_disabled"},
		{"expired", TenantAccessState{Active: true, Status: "expired"}, "subscription_expired"},
		{"unknown", TenantAccessState{Active: true, Status: "unknown"}, "subscription_expired"},
		{"missing trial", TenantAccessState{Active: true, Status: "trial"}, "trial_expired"},
		{"expired trial", TenantAccessState{Active: true, Status: "trial", TrialEnd: &past}, "trial_expired"},
		{"valid trial", TenantAccessState{Active: true, Status: "trial", TrialEnd: &future}, ""},
		{"expired paid", TenantAccessState{Active: true, Status: "active", SubscriptionEnd: &past}, "subscription_expired"},
		{"valid paid", TenantAccessState{Active: true, Status: "active", SubscriptionEnd: &future}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.state.BlockReason(now); got != tc.reason {
				t.Fatalf("got %q want %q", got, tc.reason)
			}
		})
	}
}
