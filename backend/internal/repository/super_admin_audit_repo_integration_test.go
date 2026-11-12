package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/domain"
)

func TestSuperAdminAuditRepositoryAcceptsDashboardAccessInDisposablePostgres(t *testing.T) {
	const envVar = "POOS_PHASE2A_AUDIT_TEST_DATABASE_URL"
	rawURL := strings.TrimSpace(os.Getenv(envVar))
	if rawURL == "" {
		t.Skipf("set %s to the identified WS05 disposable PostgreSQL database", envVar)
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Hostname() != "postgres" || parsedURL.User == nil || parsedURL.User.Username() != "pos_app" || strings.TrimPrefix(parsedURL.Path, "/") != "pos_phase2a_e2e" {
		t.Fatalf("refusing database URL: %s must target pos_app@postgres/pos_phase2a_e2e", envVar)
	}
	db, err := sql.Open("postgres", rawURL)
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close disposable PostgreSQL: %v", err)
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}

	var superAdminID int64
	if err := db.QueryRowContext(ctx, `SELECT id FROM super_admins ORDER BY id LIMIT 1`).Scan(&superAdminID); err != nil {
		t.Fatalf("find synthetic super admin fixture: %v", err)
	}
	marker := fmt.Sprintf("ws05-dashboard-audit-%d", time.Now().UnixNano())
	t.Cleanup(func() {
		if _, err := db.ExecContext(context.Background(), `DELETE FROM super_admin_audit_logs WHERE user_agent = $1`, marker); err != nil {
			t.Errorf("remove synthetic dashboard audit event: %v", err)
		}
	})

	repository := NewSuperAdminAuditRepository(db)
	targetType := domain.AuditTargetTenant
	targetID := int64(2)
	if err := repository.LogAction(ctx, &domain.SuperAdminAuditLog{
		SuperAdminID: superAdminID,
		Action:       domain.AuditActionAccessDashboard,
		TargetType:   &targetType,
		TargetID:     &targetID,
		IPAddress:    "127.0.0.1",
		UserAgent:    marker,
	}); err != nil {
		t.Fatalf("persist dashboard access audit event: %v", err)
	}

	var action string
	if err := db.QueryRowContext(ctx, `SELECT action FROM super_admin_audit_logs WHERE user_agent = $1`, marker).Scan(&action); err != nil {
		t.Fatalf("verify dashboard access audit event: %v", err)
	}
	if action != domain.AuditActionAccessDashboard {
		t.Fatalf("dashboard audit action = %q", action)
	}
}
