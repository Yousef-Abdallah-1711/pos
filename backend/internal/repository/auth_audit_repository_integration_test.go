//go:build integration

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"pos-saas/internal/pkg/database"
)

// Runs the real migrations in an isolated schema on a disposable database. The admin URL must
// own the schema; the app URL must connect as pos_app with NOSUPERUSER and NOBYPASSRLS.
func TestAuthAuditRLSWithRestrictedRuntimeRole(t *testing.T) {
	admin := openAuthAuditTestDB(t, "POS_RLS_ADMIN_URL")
	appURL := os.Getenv("POS_RLS_APP_URL")
	if appURL == "" {
		t.Skip("POS_RLS_APP_URL is not set")
	}
	schema := fmt.Sprintf("auth_audit_rls_%d", time.Now().UnixNano())
	if !regexp.MustCompile(`^auth_audit_rls_[0-9]+$`).MatchString(schema) {
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
	createAuthAuditFixture(t, ctx, admin)
	grantAuthAuditFixture(t, ctx, admin, schema)
	applyAuthAuditMigrations(t, ctx, admin)

	appURL = withSearchPath(t, appURL, schema)
	app := openAuthAuditURL(t, appURL)
	app.SetMaxOpenConns(1) // reuse the connection to detect transaction-local context leakage
	if _, err := app.ExecContext(ctx, `SET search_path TO `+schema); err != nil {
		t.Fatalf("set runtime search path: %v", err)
	}
	assertRestrictedRuntimeRole(t, ctx, app)

	seedAuthAuditFixture(t, ctx, admin)
	repository := NewAuthAuditRepository(app)
	if err := repository.RecordAuthEvent(ctx, 1, 101, "login", "success", "192.0.2.1", "tenant-a"); err != nil {
		t.Fatalf("tenant A should record its own event: %v", err)
	}
	if err := repository.RecordAuthEvent(ctx, 1, 202, "login", "success", "192.0.2.1", "wrong-user-tenant"); err == nil {
		t.Fatal("tenant A must not record an event for tenant B's user")
	}
	assertTenantAuditRows(t, ctx, app, 1, 1)
	assertTenantAuditRows(t, ctx, app, 2, 1)
	assertUnscopedAuditTableIsEmpty(t, ctx, app)
	assertTenantCannotMutateOtherAuditRows(t, ctx, app)
	assertTenantCannotInsertOtherTenantAuditRow(t, ctx, app)
	assertInvalidAuditTenantFailsClosed(t, ctx, app)
	if err := repository.RecordAuthEvent(ctx, 0, 0, "logout", "success", "", ""); err == nil {
		t.Fatal("tenantless audit events must be rejected before reaching the database")
	}
}

func openAuthAuditTestDB(t *testing.T, env string) *sql.DB {
	t.Helper()
	raw := strings.TrimSpace(os.Getenv(env))
	if raw == "" {
		t.Skipf("%s is not set", env)
	}
	return openAuthAuditURL(t, raw)
}

func openAuthAuditURL(t *testing.T, raw string) *sql.DB {
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

func withSearchPath(t *testing.T, raw, schema string) string {
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

func createAuthAuditFixture(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`CREATE TABLE tenants (id BIGINT PRIMARY KEY)`,
		`CREATE TABLE users (id BIGINT PRIMARY KEY, tenant_id BIGINT NOT NULL REFERENCES tenants(id))`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("create fixture: %v", err)
		}
	}
}

func grantAuthAuditFixture(t *testing.T, ctx context.Context, db *sql.DB, schema string) {
	t.Helper()
	for _, statement := range []string{
		`GRANT USAGE ON SCHEMA ` + schema + ` TO pos_app`,
		`GRANT SELECT ON tenants, users TO pos_app`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("grant runtime fixture access: %v", err)
		}
	}
}

func applyAuthAuditMigrations(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, name := range []string{"163_add_user_auth_audit_logs.sql", "164_rls_user_auth_audit_logs.sql"} {
		migration, err := os.ReadFile(filepath.Join("..", "..", "migrations", name))
		if err != nil {
			t.Fatalf("read migration %s: %v", name, err)
		}
		if _, err := db.ExecContext(ctx, string(migration)); err != nil {
			t.Fatalf("apply migration %s: %v", name, err)
		}
	}
	if _, err := db.ExecContext(ctx, `GRANT SELECT, INSERT, UPDATE, DELETE ON user_auth_audit_logs TO pos_app`); err != nil {
		t.Fatalf("grant runtime audit access: %v", err)
	}
	if _, err := db.ExecContext(ctx, `GRANT USAGE, SELECT ON SEQUENCE user_auth_audit_logs_id_seq TO pos_app`); err != nil {
		t.Fatalf("grant runtime audit sequence: %v", err)
	}
}

func seedAuthAuditFixture(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`INSERT INTO tenants (id) VALUES (1), (2)`,
		`INSERT INTO users (id, tenant_id) VALUES (101, 1), (202, 2)`,
		`INSERT INTO user_auth_audit_logs (tenant_id, user_id, action, status, user_agent) VALUES (2, 202, 'login', 'success', 'tenant-b')`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("seed auth audit fixture: %v", err)
		}
	}
}

func assertRestrictedRuntimeRole(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	var role string
	var bypass bool
	if err := db.QueryRowContext(ctx, `SELECT current_user, rolsuper OR rolbypassrls FROM pg_roles WHERE rolname = current_user`).Scan(&role, &bypass); err != nil {
		t.Fatalf("inspect runtime role: %v", err)
	}
	if role != "pos_app" || bypass {
		t.Fatalf("runtime role must be restricted pos_app, got role=%q bypass=%v", role, bypass)
	}
}

func assertTenantAuditRows(t *testing.T, ctx context.Context, db *sql.DB, tenantID int64, want int) {
	t.Helper()
	var got int
	err := database.WithTenantTx(ctx, db, tenantID, 0, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `SELECT count(*) FROM user_auth_audit_logs`).Scan(&got)
	})
	if err != nil {
		t.Fatalf("count audit rows for tenant %d: %v", tenantID, err)
	}
	if got != want {
		t.Fatalf("tenant %d sees %d audit rows, want %d", tenantID, got, want)
	}
}

func assertUnscopedAuditTableIsEmpty(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	var got int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM user_auth_audit_logs`).Scan(&got); err != nil {
		t.Fatalf("unscoped audit read should fail closed: %v", err)
	}
	if got != 0 {
		t.Fatalf("unscoped pooled connection saw %d rows; tenant context leaked", got)
	}
}

func assertTenantCannotMutateOtherAuditRows(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	err := database.WithTenantTx(ctx, db, 1, 0, func(tx *sql.Tx) error {
		updated, err := tx.ExecContext(ctx, `UPDATE user_auth_audit_logs SET user_agent = 'changed' WHERE tenant_id = 2`)
		if err != nil {
			return err
		}
		updatedCount, err := updated.RowsAffected()
		if err != nil {
			return err
		}
		if updatedCount != 0 {
			return fmt.Errorf("tenant A updated %d tenant B audit rows", updatedCount)
		}
		result, err := tx.ExecContext(ctx, `DELETE FROM user_auth_audit_logs WHERE tenant_id = 2`)
		if err != nil {
			return err
		}
		deleted, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if deleted != 0 {
			return fmt.Errorf("tenant A deleted %d tenant B audit rows", deleted)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("cross-tenant update/delete must be isolated: %v", err)
	}
}

func assertTenantCannotInsertOtherTenantAuditRow(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	err := database.WithTenantTx(ctx, db, 1, 0, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO user_auth_audit_logs (tenant_id, user_id, action, status) VALUES (2, NULL, 'login', 'success')`)
		return err
	})
	if err == nil {
		t.Fatal("cross-tenant audit insert must be rejected by RLS")
	}
}

func assertInvalidAuditTenantFailsClosed(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	var count int
	err := database.WithTenantTx(ctx, db, 1, 0, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `SELECT set_config('app.tenant_id', 'invalid', true)`); err != nil {
			return err
		}
		return tx.QueryRowContext(ctx, `SELECT count(*) FROM user_auth_audit_logs`).Scan(&count)
	})
	if err == nil {
		t.Fatal("invalid tenant context must fail closed")
	}
}
