package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
)

func TestTenantAccessUsesActivationAndEntitlement(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	future := time.Now().Add(time.Hour)
	cases := []struct {
		name   string
		active bool
		status string
		end    time.Time
		reason string
	}{
		{"legacy disabled", false, "trial", future, "tenant_disabled"},
		{"expired", true, "expired", future, "subscription_expired"},
		{"trial expired", true, "trial", past, "trial_expired"},
		{"trial valid", true, "trial", future, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			mock.ExpectQuery("SELECT is_active, is_disabled, subscription_status, trial_end_date, subscription_expires_at").WithArgs(int64(42)).WillReturnRows(sqlmock.NewRows([]string{"is_active", "is_disabled", "subscription_status", "trial_end_date", "subscription_expires_at"}).AddRow(tc.active, false, tc.status, tc.end, nil))
			allowed, reason, err := NewSubscriptionRepository(db).CheckTenantAccess(context.Background(), 42)
			if err != nil || reason != tc.reason || allowed != (tc.reason == "") {
				t.Fatalf("allowed=%v reason=%q err=%v", allowed, reason, err)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
