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
	if _, err := admin.ExecContext(ctx, `GRANT USAGE ON SCHEMA `+schema+` TO pos_app; GRANT SELECT ON tenants, restaurants, users, invoices_customer, permission_audit_logs TO pos_app; GRANT SELECT, INSERT, UPDATE, DELETE ON audit_export_logs, invoice_delivery_logs TO pos_app; GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA `+schema+` TO pos_app`); err != nil {
		t.Fatalf("grant runtime fixture access: %v", err)
	}
	applyAuditExportMigration(t, ctx, admin)
	applyInvoiceDeliveryMigration(t, ctx, admin)
	applyPermissionAuditMigration(t, ctx, admin)

	appURL = withSearchPath(t, appURL, schema)
	app := openAuthAuditURL(t, appURL)
	app.SetMaxOpenConns(1)
	if _, err := app.ExecContext(ctx, `SET search_path TO `+schema); err != nil {
		t.Fatalf("set runtime search path: %v", err)
	}
	assertRestrictedRuntimeRole(t, ctx, app)
	seedAuditExportFixture(t, ctx, admin)
	repository := NewAccountingAdvancedRepository(app, nil)
	assertPermissionAuditExportScope(t, ctx, app, repository)
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
	assertInvoiceDeliveryRLS(t, ctx, app)
	rollbackAuditExportRLS(t, ctx, admin)
	rollbackInvoiceDeliveryRLS(t, ctx, admin)
	rollbackPermissionAuditRLS(t, ctx, admin)
	assertAuditExportRollback(t, ctx, app)
	assertInvoiceDeliveryRollback(t, ctx, app)
	assertPermissionAuditRollback(t, ctx, app)
}

func createAuditExportFixture(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`CREATE TABLE tenants (id INT PRIMARY KEY)`,
		`CREATE TABLE restaurants (id INT PRIMARY KEY, tenant_id INT NOT NULL REFERENCES tenants(id))`,
		`CREATE TABLE users (id INT PRIMARY KEY, tenant_id INT NOT NULL REFERENCES tenants(id), restaurant_id INT NULL REFERENCES restaurants(id))`,
		`ALTER TABLE users ADD COLUMN name TEXT, ADD COLUMN email TEXT`,
		`CREATE TABLE invoices_customer (id INT PRIMARY KEY, tenant_id INT NOT NULL REFERENCES tenants(id), restaurant_id INT NOT NULL REFERENCES restaurants(id))`,
		`CREATE TABLE audit_export_logs (id SERIAL PRIMARY KEY, tenant_id INT NOT NULL REFERENCES tenants(id), restaurant_id INT NOT NULL REFERENCES restaurants(id), export_type VARCHAR(40) NOT NULL, export_format VARCHAR(20) NOT NULL, params_json JSONB NOT NULL DEFAULT '{}'::jsonb, row_count INT NOT NULL DEFAULT 0, generated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, generated_by INT NULL REFERENCES users(id) ON DELETE SET NULL)`,
		`CREATE TABLE invoice_delivery_logs (id SERIAL PRIMARY KEY, tenant_id INT NOT NULL REFERENCES tenants(id), restaurant_id INT NOT NULL REFERENCES restaurants(id), invoice_id INT NOT NULL REFERENCES invoices_customer(id), channel VARCHAR(20) NOT NULL DEFAULT 'email', recipient_email VARCHAR(255), status VARCHAR(20) NOT NULL DEFAULT 'sent', sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, payload JSONB NOT NULL DEFAULT '{}'::jsonb, created_by INT NULL REFERENCES users(id) ON DELETE SET NULL)`,
		`CREATE TABLE permission_audit_logs (id BIGSERIAL PRIMARY KEY, tenant_id INT NOT NULL, admin_id INT NOT NULL, user_id INT NULL, role_id INT NULL, module_id INT NULL, action TEXT NOT NULL, details JSONB, ip_address INET, user_agent TEXT, created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("create audit export fixture: %v", err)
		}
	}
}

func applyPermissionAuditMigration(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	contents, err := os.ReadFile("../../migrations/170_rls_permission_audit_logs.sql")
	if err != nil {
		t.Fatalf("read permission audit RLS migration: %v", err)
	}
	if _, err := db.ExecContext(ctx, string(contents)); err != nil {
		t.Fatalf("apply permission audit RLS migration: %v", err)
	}
}

func applyInvoiceDeliveryMigration(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	contents, err := os.ReadFile("../../migrations/169_rls_invoice_delivery_logs.sql")
	if err != nil {
		t.Fatalf("read invoice delivery RLS migration: %v", err)
	}
	if _, err := db.ExecContext(ctx, string(contents)); err != nil {
		t.Fatalf("apply invoice delivery RLS migration: %v", err)
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
		`INSERT INTO invoices_customer (id, tenant_id, restaurant_id) VALUES (1001, 1, 11), (1002, 1, 12), (2001, 2, 21)`,
		`INSERT INTO audit_export_logs (tenant_id, restaurant_id, export_type, export_format, row_count, generated_by) VALUES (2, 21, 'orders', 'csv', 4, 202)`,
		`INSERT INTO invoice_delivery_logs (tenant_id, restaurant_id, invoice_id, created_by) VALUES (2, 21, 2001, 202)`,
		`INSERT INTO permission_audit_logs (tenant_id, admin_id, action, details) VALUES (1, 101, 'UPDATE_ROLE', '{}'::jsonb), (2, 202, 'UPDATE_ROLE', '{}'::jsonb)`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("seed audit export fixture: %v", err)
		}
	}
}

func assertPermissionAuditExportScope(t *testing.T, ctx context.Context, db *sql.DB, repository *AccountingAdvancedRepository) {
	t.Helper()
	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(time.Hour)
	rows, err := repository.GetAuditExportRows(ctx, 1, 11, from, to, "rbac")
	if err != nil {
		t.Fatalf("read scoped permission audit export rows: %v", err)
	}
	if len(rows) != 1 || rows[0]["entity_type"] != "rbac_audit" || rows[0]["reference"] != "UPDATE_ROLE" {
		t.Fatalf("tenant A audit export returned unexpected rows: %#v", rows)
	}
	var unscopedCount int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM permission_audit_logs`).Scan(&unscopedCount); err != nil || unscopedCount != 0 {
		t.Fatalf("unscoped runtime connection saw %d permission audit rows (err=%v)", unscopedCount, err)
	}
}

func rollbackPermissionAuditRLS(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`DROP POLICY permission_audit_logs_tenant_scope ON permission_audit_logs`,
		`ALTER TABLE permission_audit_logs NO FORCE ROW LEVEL SECURITY`,
		`ALTER TABLE permission_audit_logs DISABLE ROW LEVEL SECURITY`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("rollback permission audit RLS: %v", err)
		}
	}
}

func assertPermissionAuditRollback(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	var enabled, forced bool
	err := db.QueryRowContext(ctx, `SELECT relrowsecurity, relforcerowsecurity FROM pg_class WHERE oid = 'permission_audit_logs'::regclass`).Scan(&enabled, &forced)
	if err != nil || enabled || forced {
		t.Fatalf("permission audit rollback left RLS enabled=%v forced=%v (err=%v)", enabled, forced, err)
	}
	var count int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM permission_audit_logs`).Scan(&count); err != nil || count != 2 {
		t.Fatalf("expected two permission audit rows after rollback, got %d (err=%v)", count, err)
	}
}

func assertInvoiceDeliveryRLS(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	insertDelivery := func(tenantID, restaurantID, invoiceID, userID int64) error {
		return database.WithTenantTx(ctx, db, tenantID, restaurantID, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, `INSERT INTO invoice_delivery_logs (tenant_id, restaurant_id, invoice_id, created_by) VALUES ($1, $2, $3, $4)`, tenantID, restaurantID, invoiceID, userID)
			return err
		})
	}
	if err := insertDelivery(1, 11, 1001, 101); err != nil {
		t.Fatalf("same-restaurant actor should insert invoice delivery log: %v", err)
	}
	for _, invalid := range [][4]int64{{1, 11, 1001, 103}, {1, 11, 1001, 202}, {1, 11, 1002, 101}} {
		if err := insertDelivery(invalid[0], invalid[1], invalid[2], invalid[3]); err == nil {
			t.Fatalf("expected RLS to reject invoice delivery scope tenant=%d restaurant=%d invoice=%d actor=%d", invalid[0], invalid[1], invalid[2], invalid[3])
		}
	}
	if err := insertDelivery(1, 12, 1002, 103); err != nil {
		t.Fatalf("sibling actor should write within its own restaurant: %v", err)
	}
	assertInvoiceDeliveryCount(t, ctx, db, 1, 11, 1)
	assertInvoiceDeliveryCount(t, ctx, db, 1, 12, 1)
	assertInvoiceDeliveryCount(t, ctx, db, 2, 21, 1)
	var unscopedCount int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM invoice_delivery_logs`).Scan(&unscopedCount); err != nil || unscopedCount != 0 {
		t.Fatalf("unscoped runtime connection saw %d invoice delivery rows (err=%v)", unscopedCount, err)
	}
	for _, statement := range []string{
		`UPDATE invoice_delivery_logs SET status = 'failed' WHERE tenant_id = 2`,
		`DELETE FROM invoice_delivery_logs WHERE tenant_id = 2`,
	} {
		if err := database.WithTenantTx(ctx, db, 1, 11, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, statement)
			return err
		}); err != nil {
			t.Fatalf("foreign invoice delivery mutation should be filtered: %v", err)
		}
	}
	assertInvoiceDeliveryCount(t, ctx, db, 2, 21, 1)
}

func assertInvoiceDeliveryCount(t *testing.T, ctx context.Context, db *sql.DB, tenantID, restaurantID int64, want int) {
	t.Helper()
	var got int
	err := database.WithTenantTx(ctx, db, tenantID, restaurantID, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `SELECT count(*) FROM invoice_delivery_logs`).Scan(&got)
	})
	if err != nil || got != want {
		t.Fatalf("tenant %d restaurant %d saw %d invoice delivery rows, want %d (err=%v)", tenantID, restaurantID, got, want, err)
	}
}

func rollbackInvoiceDeliveryRLS(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`DROP POLICY invoice_delivery_logs_tenant_restaurant_scope ON invoice_delivery_logs`,
		`ALTER TABLE invoice_delivery_logs NO FORCE ROW LEVEL SECURITY`,
		`ALTER TABLE invoice_delivery_logs DISABLE ROW LEVEL SECURITY`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("rollback invoice delivery RLS: %v", err)
		}
	}
}

func assertInvoiceDeliveryRollback(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	var enabled, forced bool
	err := db.QueryRowContext(ctx, `SELECT relrowsecurity, relforcerowsecurity FROM pg_class WHERE oid = 'invoice_delivery_logs'::regclass`).Scan(&enabled, &forced)
	if err != nil || enabled || forced {
		t.Fatalf("invoice delivery rollback left RLS enabled=%v forced=%v (err=%v)", enabled, forced, err)
	}
	var count int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM invoice_delivery_logs`).Scan(&count); err != nil || count != 3 {
		t.Fatalf("expected three invoice delivery rows after rollback, got %d (err=%v)", count, err)
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
