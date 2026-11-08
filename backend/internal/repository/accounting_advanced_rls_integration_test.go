//go:build integration

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/pkg/database"
)

func TestAccountingAdvancedAuditExportRLSWithRestrictedRuntimeRole(t *testing.T) {
	admin := openAuthAuditTestDB(t, "POS_RLS_ADMIN_URL")
	appURL := os.Getenv("POS_RLS_APP_URL")
	if appURL == "" {
		t.Skip("POS_RLS_APP_URL is not set")
	}
	schema := fmt.Sprintf("audit_export_rls_%d", time.Now().UnixNano())
	if !regexp.MustCompile(`^audit_export_rls_[0-9]+$`).MatchString(schema) {
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
	createAuditExportFixture(t, ctx, admin)
	if _, err := admin.ExecContext(ctx, `GRANT USAGE ON SCHEMA `+schema+` TO pos_app; GRANT SELECT ON tenants, restaurants, users TO pos_app; GRANT SELECT, INSERT, UPDATE, DELETE ON audit_export_logs TO pos_app; GRANT USAGE, SELECT ON SEQUENCE audit_export_logs_id_seq TO pos_app`); err != nil {
		t.Fatalf("grant runtime fixture access: %v", err)
	}
	applyAuditExportMigration(t, ctx, admin)

	appURL = withSearchPath(t, appURL, schema)
	app := openAuthAuditURL(t, appURL)
	app.SetMaxOpenConns(1)
	if _, err := app.ExecContext(ctx, `SET search_path TO `+schema); err != nil {
		t.Fatalf("set runtime search path: %v", err)
	}
	assertRestrictedRuntimeRole(t, ctx, app)
	seedAuditExportFixture(t, ctx, admin)
	repository := NewAccountingAdvancedRepository(app, nil)
	if err := repository.LogAuditExport(ctx, 1, 11, 101, "orders", "csv", map[string]interface{}{"range": "today"}, 3); err != nil {
		t.Fatalf("same-restaurant actor should write export log: %v", err)
	}
	if err := repository.LogAuditExport(ctx, 1, 11, 103, "orders", "csv", nil, 1); err == nil {
		t.Fatal("sibling-restaurant actor must not write an export log")
	}
	if err := repository.LogAuditExport(ctx, 1, 11, 202, "orders", "csv", nil, 1); err == nil {
		t.Fatal("foreign-tenant actor must not write an export log")
	}
	if err := repository.LogAuditExport(ctx, 1, 12, 103, "orders", "csv", nil, 2); err != nil {
		t.Fatalf("sibling restaurant actor should write within its own restaurant: %v", err)
	}
	assertAuditExportCount(t, ctx, app, 1, 11, 1)
	assertAuditExportCount(t, ctx, app, 1, 12, 1)
	assertAuditExportCount(t, ctx, app, 2, 21, 1)
	assertAuditExportNoContextIsEmpty(t, ctx, app)
	assertAuditExportForeignMutationsAreFiltered(t, ctx, app)
	assertAuditExportForeignInsertIsRejected(t, ctx, app)
	rollbackAuditExportRLS(t, ctx, admin)
	assertAuditExportRollback(t, ctx, app)
}

func createAuditExportFixture(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`CREATE TABLE tenants (id INT PRIMARY KEY)`,
		`CREATE TABLE restaurants (id INT PRIMARY KEY, tenant_id INT NOT NULL REFERENCES tenants(id))`,
		`CREATE TABLE users (id INT PRIMARY KEY, tenant_id INT NOT NULL REFERENCES tenants(id), restaurant_id INT NULL REFERENCES restaurants(id))`,
		`CREATE TABLE audit_export_logs (id SERIAL PRIMARY KEY, tenant_id INT NOT NULL REFERENCES tenants(id), restaurant_id INT NOT NULL REFERENCES restaurants(id), export_type VARCHAR(40) NOT NULL, export_format VARCHAR(20) NOT NULL, params_json JSONB NOT NULL DEFAULT '{}'::jsonb, row_count INT NOT NULL DEFAULT 0, generated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, generated_by INT NULL REFERENCES users(id) ON DELETE SET NULL)`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("create audit export fixture: %v", err)
		}
	}
}

func applyAuditExportMigration(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	contents, err := os.ReadFile("../../migrations/167_rls_audit_export_logs.sql")
	if err != nil {
		t.Fatalf("read audit export RLS migration: %v", err)
	}
	if _, err := db.ExecContext(ctx, string(contents)); err != nil {
		t.Fatalf("apply audit export RLS migration: %v", err)
	}
}

func seedAuditExportFixture(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`INSERT INTO tenants (id) VALUES (1), (2)`,
		`INSERT INTO restaurants (id, tenant_id) VALUES (11, 1), (12, 1), (21, 2)`,
		`INSERT INTO users (id, tenant_id, restaurant_id) VALUES (101, 1, 11), (103, 1, 12), (202, 2, 21)`,
		`INSERT INTO audit_export_logs (tenant_id, restaurant_id, export_type, export_format, row_count, generated_by) VALUES (2, 21, 'orders', 'csv', 4, 202)`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("seed audit export fixture: %v", err)
		}
	}
}

func assertAuditExportCount(t *testing.T, ctx context.Context, db *sql.DB, tenantID, restaurantID int64, want int) {
	t.Helper()
	var got int
	err := database.WithTenantTx(ctx, db, tenantID, restaurantID, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `SELECT count(*) FROM audit_export_logs`).Scan(&got)
	})
	if err != nil || got != want {
		t.Fatalf("tenant %d restaurant %d saw %d audit export rows, want %d (err=%v)", tenantID, restaurantID, got, want, err)
	}
}

func assertAuditExportNoContextIsEmpty(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	var got int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM audit_export_logs`).Scan(&got); err != nil {
		t.Fatalf("query without RLS context: %v", err)
	}
	if got != 0 {
		t.Fatalf("unscoped runtime connection saw %d audit export rows", got)
	}
}

func assertAuditExportForeignMutationsAreFiltered(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`UPDATE audit_export_logs SET row_count = 99 WHERE tenant_id = 2`,
		`DELETE FROM audit_export_logs WHERE tenant_id = 2`,
	} {
		if err := database.WithTenantTx(ctx, db, 1, 11, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, statement)
			return err
		}); err != nil {
			t.Fatalf("foreign export mutation should be filtered: %v", err)
		}
	}
	assertAuditExportCount(t, ctx, db, 2, 21, 1)
}

func assertAuditExportForeignInsertIsRejected(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	err := database.WithTenantTx(ctx, db, 1, 11, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO audit_export_logs (tenant_id, restaurant_id, export_type, export_format, generated_by) VALUES (2, 21, 'orders', 'csv', 202)`)
		return err
	})
	if err == nil {
		t.Fatal("tenant A must not insert an export record for tenant B")
	}
}

func rollbackAuditExportRLS(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`DROP POLICY audit_export_logs_tenant_restaurant_scope ON audit_export_logs`,
		`ALTER TABLE audit_export_logs NO FORCE ROW LEVEL SECURITY`,
		`ALTER TABLE audit_export_logs DISABLE ROW LEVEL SECURITY`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("rollback audit export RLS: %v", err)
		}
	}
}

func assertAuditExportRollback(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	var enabled, forced bool
	err := db.QueryRowContext(ctx, `SELECT relrowsecurity, relforcerowsecurity FROM pg_class WHERE oid = 'audit_export_logs'::regclass`).Scan(&enabled, &forced)
	if err != nil || enabled || forced {
		t.Fatalf("audit export rollback left RLS enabled=%v forced=%v (err=%v)", enabled, forced, err)
	}
	var count int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM audit_export_logs`).Scan(&count); err != nil || count != 3 {
		t.Fatalf("expected three rows after rollback, got %d (err=%v)", count, err)
	}
}
