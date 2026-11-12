//go:build integration

package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/pkg/database"
)

func TestReportExportAuditUsesRestrictedTenantRLS(t *testing.T) {
	admin := openReportRLSDB(t, "POS_RLS_ADMIN_URL")
	appURL := strings.TrimSpace(os.Getenv("POS_RLS_APP_URL"))
	if appURL == "" {
		t.Skip("POS_RLS_APP_URL is not set")
	}
	schema := fmt.Sprintf("report_export_rls_%d", time.Now().UnixNano())
	if !regexp.MustCompile(`^report_export_rls_[0-9]+$`).MatchString(schema) {
		t.Fatal("generated unsafe schema name")
	}
	ctx := context.Background()
	admin.SetMaxOpenConns(1)
	if _, err := admin.ExecContext(ctx, `CREATE SCHEMA `+schema); err != nil {
		t.Fatalf("create isolated schema: %v", err)
	}
	t.Cleanup(func() {
		if _, err := admin.ExecContext(context.Background(), `DROP SCHEMA `+schema+` CASCADE`); err != nil {
			t.Errorf("drop isolated schema: %v", err)
		}
	})
	if _, err := admin.ExecContext(ctx, `SET search_path TO `+schema); err != nil {
		t.Fatalf("set admin search path: %v", err)
	}
	createReportRLSFixture(t, ctx, admin)
	if _, err := admin.ExecContext(ctx, `GRANT USAGE ON SCHEMA `+schema+` TO pos_app; GRANT SELECT ON tenants, restaurants, users TO pos_app; GRANT SELECT, INSERT, UPDATE, DELETE ON report_export_audit_logs TO pos_app; GRANT USAGE, SELECT ON SEQUENCE report_export_audit_logs_id_seq TO pos_app`); err != nil {
		t.Fatalf("grant runtime fixture access: %v", err)
	}
	applyReportRLSMigration(t, ctx, admin)

	appURL = reportURLWithSearchPath(t, appURL, schema)
	app := openReportRLSURL(t, appURL)
	app.SetMaxOpenConns(1)
	if _, err := app.ExecContext(ctx, `SET search_path TO `+schema); err != nil {
		t.Fatalf("set runtime search path: %v", err)
	}
	assertReportRLSAppRole(t, ctx, app)
	seedReportRLSFixture(t, ctx, admin)
	handler := NewReportHandler(app)
	if err := handler.logReportExportAudit(ctx, reportExportAuditRecord{TenantID: 1, RestaurantID: 11, UserID: 101, ReportType: "sales", Format: "csv", ExportedBytes: 42, Filters: map[string]interface{}{"range": "today"}}); err != nil {
		t.Fatalf("same-restaurant report export should be logged: %v", err)
	}
	if err := handler.logReportExportAudit(ctx, reportExportAuditRecord{TenantID: 1, RestaurantID: 11, UserID: 103, ReportType: "sales", Format: "csv"}); err == nil {
		t.Fatal("sibling-restaurant user must not log an export for another restaurant")
	}
	if err := handler.logReportExportAudit(ctx, reportExportAuditRecord{TenantID: 1, RestaurantID: 11, UserID: 202, ReportType: "sales", Format: "csv"}); err == nil {
		t.Fatal("foreign-tenant user must not log an export")
	}
	if err := handler.logReportExportAudit(ctx, reportExportAuditRecord{TenantID: 1, RestaurantID: 12, UserID: 103, ReportType: "sales", Format: "csv"}); err != nil {
		t.Fatalf("sibling user should log an export in its own restaurant: %v", err)
	}
	assertReportExportCount(t, ctx, app, 1, 11, 1)
	assertReportExportCount(t, ctx, app, 1, 12, 1)
	assertReportExportCount(t, ctx, app, 2, 21, 1)
	assertReportExportUnscopedEmpty(t, ctx, app)
	assertReportExportForeignMutationsFiltered(t, ctx, app)
	rollbackReportExportRLS(t, ctx, admin)
	assertReportExportRollback(t, ctx, app)
}

func openReportRLSDB(t *testing.T, env string) *sql.DB {
	t.Helper()
	raw := strings.TrimSpace(os.Getenv(env))
	if raw == "" {
		t.Skipf("%s is not set", env)
	}
	return openReportRLSURL(t, raw)
}

func openReportRLSURL(t *testing.T, raw string) *sql.DB {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse disposable PostgreSQL URL: %v", err)
	}
	if !strings.Contains(strings.ToLower(strings.Trim(parsed.Path, "/")), "test") {
		t.Fatalf("refusing to run against non-test database %q", parsed.Path)
	}
	db, err := sql.Open("postgres", raw)
	if err != nil {
		t.Fatalf("open disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.PingContext(context.Background()); err != nil {
		t.Fatalf("connect to disposable PostgreSQL: %v", err)
	}
	return db
}

func createReportRLSFixture(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`CREATE TABLE tenants (id BIGINT PRIMARY KEY)`,
		`CREATE TABLE restaurants (id BIGINT PRIMARY KEY, tenant_id BIGINT NOT NULL REFERENCES tenants(id))`,
		`CREATE TABLE users (id BIGINT PRIMARY KEY, tenant_id BIGINT NOT NULL REFERENCES tenants(id), restaurant_id BIGINT NULL REFERENCES restaurants(id))`,
		`CREATE TABLE report_export_audit_logs (id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE, restaurant_id BIGINT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE, user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE, report_type VARCHAR(50) NOT NULL, format VARCHAR(20) NOT NULL, file_name VARCHAR(255), exported_bytes INT NOT NULL DEFAULT 0, filters JSONB NOT NULL DEFAULT '{}'::jsonb, request_path TEXT, user_agent TEXT, ip_address VARCHAR(64), created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("create report export fixture: %v", err)
		}
	}
}

func applyReportRLSMigration(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	contents, err := os.ReadFile("../../../migrations/168_rls_report_export_audit_logs.sql")
	if err != nil {
		t.Fatalf("read report export RLS migration: %v", err)
	}
	if _, err := db.ExecContext(ctx, string(contents)); err != nil {
		t.Fatalf("apply report export RLS migration: %v", err)
	}
}

func seedReportRLSFixture(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`INSERT INTO tenants (id) VALUES (1), (2)`,
		`INSERT INTO restaurants (id, tenant_id) VALUES (11, 1), (12, 1), (21, 2)`,
		`INSERT INTO users (id, tenant_id, restaurant_id) VALUES (101, 1, 11), (103, 1, 12), (202, 2, 21)`,
		`INSERT INTO report_export_audit_logs (tenant_id, restaurant_id, user_id, report_type, format) VALUES (2, 21, 202, 'sales', 'csv')`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("seed report export fixture: %v", err)
		}
	}
}

func reportURLWithSearchPath(t *testing.T, raw, schema string) string {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse runtime PostgreSQL URL: %v", err)
	}
	query := parsed.Query()
	query.Set("options", "-csearch_path="+schema)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func assertReportRLSAppRole(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	var role string
	var bypass bool
	if err := db.QueryRowContext(ctx, `SELECT current_user, rolsuper OR rolbypassrls FROM pg_roles WHERE rolname = current_user`).Scan(&role, &bypass); err != nil || role != "pos_app" || bypass {
		t.Fatalf("runtime role must be restricted pos_app, got role=%q bypass=%v err=%v", role, bypass, err)
	}
}

func assertReportExportCount(t *testing.T, ctx context.Context, db *sql.DB, tenantID, restaurantID int64, want int) {
	t.Helper()
	var got int
	err := database.WithTenantTx(ctx, db, tenantID, restaurantID, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `SELECT count(*) FROM report_export_audit_logs`).Scan(&got)
	})
	if err != nil || got != want {
		t.Fatalf("tenant %d restaurant %d saw %d report export rows, want %d (err=%v)", tenantID, restaurantID, got, want, err)
	}
}

func assertReportExportUnscopedEmpty(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	var count int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM report_export_audit_logs`).Scan(&count); err != nil || count != 0 {
		t.Fatalf("unscoped runtime connection saw %d report export rows (err=%v)", count, err)
	}
}

func assertReportExportForeignMutationsFiltered(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`UPDATE report_export_audit_logs SET file_name = 'tampered.csv' WHERE tenant_id = 2`,
		`DELETE FROM report_export_audit_logs WHERE tenant_id = 2`,
	} {
		if err := database.WithTenantTx(ctx, db, 1, 11, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, statement)
			return err
		}); err != nil {
			t.Fatalf("foreign report export mutation should be filtered: %v", err)
		}
	}
	assertReportExportCount(t, ctx, db, 2, 21, 1)
}

func rollbackReportExportRLS(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`DROP POLICY report_export_audit_logs_tenant_restaurant_scope ON report_export_audit_logs`,
		`ALTER TABLE report_export_audit_logs NO FORCE ROW LEVEL SECURITY`,
		`ALTER TABLE report_export_audit_logs DISABLE ROW LEVEL SECURITY`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("rollback report export RLS: %v", err)
		}
	}
}

func assertReportExportRollback(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	var enabled, forced bool
	if err := db.QueryRowContext(ctx, `SELECT relrowsecurity, relforcerowsecurity FROM pg_class WHERE oid = 'report_export_audit_logs'::regclass`).Scan(&enabled, &forced); err != nil || enabled || forced {
		t.Fatalf("report export rollback left RLS enabled=%v forced=%v (err=%v)", enabled, forced, err)
	}
	var count int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM report_export_audit_logs`).Scan(&count); err != nil || count != 3 {
		t.Fatalf("expected three report export rows after rollback, got %d (err=%v)", count, err)
	}
}
