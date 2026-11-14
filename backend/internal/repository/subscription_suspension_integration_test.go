package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/tests/testdb"
)

func TestSuspensionPreservesPaidAndExpiredEntitlement(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	for _, status := range []string{"active", "expired"} {
		t.Run(status, func(t *testing.T) {
			db, fixture := openDriverScopeFixture(t)
			end := time.Now().UTC().Add(48 * time.Hour).Truncate(time.Second)
			if _, err := db.Exec(`UPDATE tenants SET is_active=true,is_disabled=false,subscription_status=$1,subscription_expires_at=$2 WHERE id=$3`, status, end, fixture.tenantA); err != nil {
				t.Fatal(err)
			}
			repo := NewSubscriptionRepository(db)
			ctx := context.Background()
			id := int64(fixture.tenantA)
			params := &domain.DisableTenantParams{TenantID: id, PerformedBy: 1, Reason: "isolated regression"}
			for i := 0; i < 2; i++ {
				if err := repo.DisableTenant(ctx, params); err != nil {
					t.Fatal(err)
				}
			}
			if allowed, _, err := repo.CheckTenantAccess(ctx, id); err != nil || allowed {
				t.Fatalf("disabled access=%v err=%v", allowed, err)
			}
			if err := repo.EnableTenant(ctx, &domain.EnableTenantParams{TenantID: id, PerformedBy: 1}); err != nil {
				t.Fatal(err)
			}
			var gotStatus string
			var gotEnd time.Time
			if err := db.QueryRow(`SELECT subscription_status,subscription_expires_at FROM tenants WHERE id=$1`, id).Scan(&gotStatus, &gotEnd); err != nil {
				t.Fatal(err)
			}
			if gotStatus != status || !gotEnd.Equal(end) {
				t.Fatalf("billing entitlement changed: status=%s end=%s", gotStatus, gotEnd)
			}
			if allowed, _, err := repo.CheckTenantAccess(ctx, id); err != nil || allowed != (status == "active") {
				t.Fatalf("restored access=%v status=%s err=%v", allowed, status, err)
			}
		})
	}
}

func TestLegacyMissingTrialDatesRemainAnchoredToCreation(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	for _, age := range []int{1, 16} {
		if _, err := db.Exec(`UPDATE tenants SET is_active=true,is_disabled=false,subscription_status='trial',trial_days=15,trial_start_date=NULL,trial_end_date=NULL,created_at=NOW()-$1*INTERVAL '1 day' WHERE id=$2`, age, fixture.tenantA); err != nil {
			t.Fatal(err)
		}
		repo := NewSubscriptionRepository(db)
		allowed, reason, err := repo.CheckTenantAccess(context.Background(), int64(fixture.tenantA))
		if err != nil || allowed != (age == 1) {
			t.Fatalf("age=%d allowed=%v reason=%s err=%v", age, allowed, reason, err)
		}
		info, err := repo.GetTenantSubscriptionInfo(context.Background(), int64(fixture.tenantA))
		if err != nil {
			t.Fatal(err)
		}
		if info.CanAccess != allowed || info.TrialStartDate == nil || info.TrialEndDate == nil {
			t.Fatal("public status and business entitlement disagree")
		}
		if !info.TrialEndDate.Equal(info.TrialStartDate.AddDate(0, 0, 15)) {
			t.Fatal("legacy trial was extended from now")
		}
	}
}
