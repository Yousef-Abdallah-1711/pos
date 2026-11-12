package repository

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"pos-saas/internal/domain"
	"pos-saas/tests/testdb"
)

func TestSubscriptionTrialExtensionsSerializeConcurrentUpdates(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db, fixture := openDriverScopeFixture(t)
	baseEnd := time.Now().UTC().Truncate(time.Second)
	const initialDays = 3
	if _, err := db.Exec(`
		UPDATE tenants
		SET subscription_status = 'expired', trial_end_date = $1, trial_days = $2
		WHERE id = $3
	`, baseEnd, initialDays, fixture.tenantA); err != nil {
		t.Fatalf("seed trial concurrency state: %v", err)
	}

	const extensions = 12
	start := make(chan struct{})
	results := make(chan error, extensions)
	var workers sync.WaitGroup
	repo := NewSubscriptionRepository(db)
	for index := 0; index < extensions; index++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			results <- repo.ExtendTrial(ctx, &domain.ExtendTrialParams{
				TenantID: int64(fixture.tenantA), Days: 1, PerformedBy: 1,
				IPAddress: "127.0.0.1", UserAgent: "isolated-concurrency-test",
			})
		}()
	}
	close(start)
	workers.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatalf("extend trial concurrently: %v", err)
		}
	}

	var gotEnd time.Time
	var gotDays int
	var gotStatus string
	err := db.QueryRow(`SELECT trial_end_date, trial_days, subscription_status FROM tenants WHERE id = $1`, fixture.tenantA).Scan(&gotEnd, &gotDays, &gotStatus)
	if err != nil {
		t.Fatalf("read trial state after concurrent extensions: %v", err)
	}
	wantEnd := baseEnd.AddDate(0, 0, extensions)
	if !gotEnd.Equal(wantEnd) || gotDays != initialDays+extensions || gotStatus != "trial" {
		t.Fatalf("concurrent extensions lost an update: end=%s days=%d status=%s; want end=%s days=%d status=trial",
			gotEnd, gotDays, gotStatus, wantEnd, initialDays+extensions)
	}

	var historyCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM tenant_subscription_history WHERE tenant_id = $1 AND action = $2`, fixture.tenantA, domain.ActionTrialExtended).Scan(&historyCount); err != nil {
		t.Fatalf("count trial extension history: %v", err)
	}
	if historyCount != extensions {
		t.Fatalf("concurrent extensions recorded %d history rows, want %d (tenant %d)", historyCount, extensions, fixture.tenantA)
	}
}
