package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/domain"
	"pos-saas/tests/testdb"
)

func TestSubscriptionExpiryWorkerPostgresScopesEachTenantByItsOwnExpiry(t *testing.T) {
	if os.Getenv(testdb.EnvVar) == "" {
		t.Skipf("set %s to an isolated disposable PostgreSQL database", testdb.EnvVar)
	}
	db, err := sql.Open("postgres", testdb.URL(t))
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close disposable PostgreSQL connection: %v", err)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}

	runID := fmt.Sprintf("%d", time.Now().UnixNano())
	expiredTenantID := insertExpiryWorkerTenant(t, ctx, db, runID+"-expired", time.Now().Add(-time.Hour))
	futureTenantID := insertExpiryWorkerTenant(t, ctx, db, runID+"-future", time.Now().Add(24*time.Hour))
	if err := NewSubscriptionExpiryWorker(db).runOnce(ctx); err != nil {
		t.Fatalf("run subscription expiry worker: %v", err)
	}

	assertExpiryWorkerTenant(t, ctx, db, expiredTenantID, "expired", 1)
	assertExpiryWorkerTenant(t, ctx, db, futureTenantID, "active", 0)
}

func insertExpiryWorkerTenant(t *testing.T, ctx context.Context, db *sql.DB, key string, expiresAt time.Time) int64 {
	t.Helper()
	var tenantID int64
	err := db.QueryRowContext(ctx, `
		INSERT INTO tenants (name, slug, domain, subscription_status, subscription_expires_at)
		VALUES ($1, $1, $2, 'active', $3)
		RETURNING id
	`, "closure-worker-"+key, "closure-worker-"+key+".example.test", expiresAt).Scan(&tenantID)
	if err != nil {
		t.Fatalf("insert synthetic subscription tenant: %v", err)
	}
	return tenantID
}

func assertExpiryWorkerTenant(t *testing.T, ctx context.Context, db *sql.DB, tenantID int64, wantStatus string, wantHistory int) {
	t.Helper()
	var status string
	if err := db.QueryRowContext(ctx, `SELECT subscription_status FROM tenants WHERE id = $1`, tenantID).Scan(&status); err != nil {
		t.Fatalf("read tenant subscription state: %v", err)
	}
	if status != wantStatus {
		t.Fatalf("tenant %d status=%q, want %q", tenantID, status, wantStatus)
	}
	var historyCount int
	if err := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM tenant_subscription_history WHERE tenant_id = $1 AND action = $2
	`, tenantID, domain.ActionSubscriptionExpired).Scan(&historyCount); err != nil {
		t.Fatalf("count tenant subscription history: %v", err)
	}
	if historyCount != wantHistory {
		t.Fatalf("tenant %d history rows=%d, want %d", tenantID, historyCount, wantHistory)
	}
}
