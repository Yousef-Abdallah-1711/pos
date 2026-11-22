//go:build integration

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestAuthIdentityMembershipSchemaIsDefaultOffAndTenantConsistent(t *testing.T) {
	admin := openAuthAuditTestDB(t, "POS_RLS_ADMIN_URL")
	appURL := os.Getenv("POS_RLS_APP_URL")
	if appURL == "" {
		t.Skip("POS_RLS_APP_URL is not set")
	}
	ctx := context.Background()
	admin.SetMaxOpenConns(1)
	schema := fmt.Sprintf("auth_identity_schema_%d", time.Now().UnixNano())
	if !regexp.MustCompile(`^auth_identity_schema_[0-9]+$`).MatchString(schema) {
		t.Fatal("generated unsafe schema name")
	}
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
	createIdentitySchemaFixture(t, ctx, admin, schema)

	migration, err := os.ReadFile(filepath.Join("..", "..", "migrations", "191_auth_identity_membership_schema.sql"))
	if err != nil {
		t.Fatalf("read identity schema migration: %v", err)
	}
	if _, err := admin.ExecContext(ctx, string(migration)); err != nil {
		t.Fatalf("apply identity schema migration: %v", err)
	}

	assertIdentitySchemaPrivileges(t, ctx, admin, schema)
	assertIdentityRoleReprovisioningPreservesIsolation(t, ctx, admin, schema)
	assertIdentitySchemaRLS(t, ctx, admin)
	seedAndAssertIdentitySchema(t, ctx, admin)
	assertIdentitySchemaRuntimeRoleDenied(t, ctx, appURL, schema)
	assertIdentitySchemaRollback(t, ctx, admin, schema)
}

func createIdentitySchemaFixture(t *testing.T, ctx context.Context, db *sql.DB, schema string) {
	t.Helper()
	// This fixture is deliberately confined to a unique disposable schema.
	for _, statement := range []string{
		`GRANT USAGE ON SCHEMA ` + schema + ` TO pos_app, pos_migrator`,
		`CREATE TABLE tenants (id INTEGER PRIMARY KEY)`,
		`CREATE TABLE users (id INTEGER PRIMARY KEY, tenant_id INTEGER NOT NULL REFERENCES tenants(id))`,
		`CREATE TABLE restaurants (id INTEGER PRIMARY KEY, tenant_id INTEGER NOT NULL REFERENCES tenants(id))`,
		`CREATE TABLE roles (id INTEGER PRIMARY KEY, tenant_id INTEGER NOT NULL REFERENCES tenants(id), restaurant_id INTEGER NOT NULL REFERENCES restaurants(id))`,
		`INSERT INTO tenants (id) VALUES (1), (2)`,
		`INSERT INTO users (id, tenant_id) VALUES (101, 1), (202, 2)`,
		`INSERT INTO restaurants (id, tenant_id) VALUES (11, 1), (22, 2)`,
		`INSERT INTO roles (id, tenant_id, restaurant_id) VALUES (111, 1, 11), (222, 2, 22)`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("prepare schema fixture: %v", err)
		}
	}
}

func assertIdentitySchemaPrivileges(t *testing.T, ctx context.Context, db *sql.DB, schema string) {
	t.Helper()
	for _, table := range []string{
		"auth_identities",
		"identity_email_challenges",
		"tenant_memberships",
		"membership_restaurant_grants",
		"membership_role_grants",
	} {
		var appCanSelect, appCanInsert, appCanUpdate, appCanDelete, migratorCanSelect bool
		qualifiedName := schema + "." + table
		if err := db.QueryRowContext(ctx, `
			SELECT has_table_privilege('pos_app', $1, 'SELECT'),
			       has_table_privilege('pos_app', $1, 'INSERT'),
			       has_table_privilege('pos_app', $1, 'UPDATE'),
			       has_table_privilege('pos_app', $1, 'DELETE'),
			       has_table_privilege('pos_migrator', $1, 'SELECT')
		`, qualifiedName).Scan(&appCanSelect, &appCanInsert, &appCanUpdate, &appCanDelete, &migratorCanSelect); err != nil {
			t.Fatalf("inspect %s privileges: %v", table, err)
		}
		if appCanSelect || appCanInsert || appCanUpdate || appCanDelete {
			t.Errorf("pos_app must not read or write %s before the auth boundary is implemented", table)
		}
		if !migratorCanSelect {
			t.Errorf("pos_migrator must retain reviewed migration access to %s", table)
		}
	}
}

func assertIdentityRoleReprovisioningPreservesIsolation(t *testing.T, ctx context.Context, db *sql.DB, schema string) {
	t.Helper()
	roleScript, err := os.ReadFile(filepath.Join("..", "..", "scripts", "db", "create_runtime_roles.sql"))
	if err != nil {
		t.Fatalf("read runtime role provisioning script: %v", err)
	}
	grantAll := strings.Index(string(roleScript), "GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO pos_app;")
	resetInclude := strings.Index(string(roleScript), `\ir revoke_phase1_auth_runtime_privileges.sql`)
	if grantAll < 0 || resetInclude < grantAll {
		t.Fatal("runtime role provisioning must reapply identity-table denials after broad grants")
	}

	for _, table := range []string{
		"auth_identities",
		"identity_email_challenges",
		"tenant_memberships",
		"membership_restaurant_grants",
		"membership_role_grants",
	} {
		qualifiedName := schema + "." + table
		if _, err := db.ExecContext(ctx, `GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE `+qualifiedName+` TO pos_app`); err != nil {
			t.Fatalf("simulate role re-provisioning grant on %s: %v", table, err)
		}
	}
	var appCanSelect bool
	if err := db.QueryRowContext(ctx, `SELECT has_table_privilege('pos_app', $1, 'SELECT')`, schema+".auth_identities").Scan(&appCanSelect); err != nil {
		t.Fatalf("inspect simulated re-provisioning grant: %v", err)
	}
	if !appCanSelect {
		t.Fatal("test did not simulate the broad runtime-role re-provisioning grant")
	}

	reset, err := os.ReadFile(filepath.Join("..", "..", "scripts", "db", "revoke_phase1_auth_runtime_privileges.sql"))
	if err != nil {
		t.Fatalf("read identity runtime privilege reset: %v", err)
	}
	if _, err := db.ExecContext(ctx, string(reset)); err != nil {
		t.Fatalf("reapply identity runtime privilege boundary: %v", err)
	}
	assertIdentitySchemaPrivileges(t, ctx, db, schema)
}

func assertIdentitySchemaRLS(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	rows, err := db.QueryContext(ctx, `
		SELECT c.relname, c.relrowsecurity, c.relforcerowsecurity
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = current_schema()
		  AND c.relname IN ('auth_identities', 'identity_email_challenges', 'tenant_memberships', 'membership_restaurant_grants', 'membership_role_grants')
	`)
	if err != nil {
		t.Fatalf("inspect row-level security: %v", err)
	}
	defer rows.Close()
	seen := map[string]bool{}
	for rows.Next() {
		var name string
		var enabled, forced bool
		if err := rows.Scan(&name, &enabled, &forced); err != nil {
			t.Fatalf("read row-level security state: %v", err)
		}
		seen[name] = enabled && forced
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate row-level security state: %v", err)
	}
	for _, table := range []string{
		"auth_identities",
		"identity_email_challenges",
		"tenant_memberships",
		"membership_restaurant_grants",
		"membership_role_grants",
	} {
		if !seen[table] {
			t.Errorf("%s must enable and force row-level security", table)
		}
	}
}

func seedAndAssertIdentitySchema(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	if _, err := db.ExecContext(ctx, `SET ROLE pos_migrator`); err != nil {
		t.Fatalf("assume migration role: %v", err)
	}
	defer func() { _, _ = db.ExecContext(context.Background(), `RESET ROLE`) }()

	const identityID = "00000000-0000-4000-8000-000000000001"
	const membershipID = "00000000-0000-4000-8000-000000000002"
	const restaurantGrantID = "00000000-0000-4000-8000-000000000003"
	const roleGrantID = "00000000-0000-4000-8000-000000000004"
	const challengeID = "00000000-0000-4000-8000-000000000005"
	if _, err := db.ExecContext(ctx, `INSERT INTO auth_identities (id, normalized_email, email_verified_at, password_hash, status) VALUES ($1, 'owner@example.test', CURRENT_TIMESTAMP, 'synthetic-hash', 'active')`, identityID); err != nil {
		t.Fatalf("insert verified identity: %v", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO auth_identities (id, normalized_email, password_hash) VALUES ('00000000-0000-4000-8000-000000000011', 'unverified@example.test', 'synthetic-hash'), ('00000000-0000-4000-8000-000000000012', 'unverified@example.test', 'synthetic-hash')`); err != nil {
		t.Fatalf("unverified duplicates must remain representable: %v", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO auth_identities (id, normalized_email, email_verified_at, password_hash, status) VALUES ('00000000-0000-4000-8000-000000000013', 'owner@example.test', CURRENT_TIMESTAMP, 'synthetic-hash', 'active')`); err == nil {
		t.Fatal("a verified email must identify at most one global identity")
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO auth_identities (id, normalized_email, password_hash, status) VALUES ('00000000-0000-4000-8000-000000000014', 'pending@example.test', 'synthetic-hash', 'active')`); err == nil {
		t.Fatal("an unverified identity must not become active")
	}

	if _, err := db.ExecContext(ctx, `INSERT INTO tenant_memberships (id, identity_id, tenant_id, legacy_user_id, state, accepted_at) VALUES ($1, $2, 1, 101, 'active', CURRENT_TIMESTAMP)`, membershipID, identityID); err != nil {
		t.Fatalf("insert tenant membership: %v", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO tenant_memberships (id, identity_id, tenant_id, legacy_user_id) VALUES ('00000000-0000-4000-8000-000000000021', $1, 1, 202)`, identityID); err == nil {
		t.Fatal("membership must not bind a legacy user from a different tenant")
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO tenant_memberships (id, identity_id, tenant_id) VALUES ('00000000-0000-4000-8000-000000000022', $1, 1)`, identityID); err == nil {
		t.Fatal("one identity must not have duplicate memberships in one tenant")
	}

	if _, err := db.ExecContext(ctx, `INSERT INTO membership_restaurant_grants (id, membership_id, tenant_id, restaurant_id, state) VALUES ($1, $2, 1, 11, 'active')`, restaurantGrantID, membershipID); err != nil {
		t.Fatalf("insert explicit restaurant grant: %v", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO membership_restaurant_grants (id, membership_id, tenant_id, restaurant_id) VALUES ('00000000-0000-4000-8000-000000000031', $1, 1, 22)`, membershipID); err == nil {
		t.Fatal("restaurant grants must reject a restaurant in a different tenant")
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO membership_role_grants (id, membership_id, tenant_id, restaurant_id, role_id, state) VALUES ($1, $2, 1, 11, 111, 'active')`, roleGrantID, membershipID); err != nil {
		t.Fatalf("insert scoped role grant: %v", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO membership_role_grants (id, membership_id, tenant_id, restaurant_id, role_id) VALUES ('00000000-0000-4000-8000-000000000041', $1, 1, 11, 222)`, membershipID); err == nil {
		t.Fatal("role grants must reject a role from a different tenant and restaurant")
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO membership_role_grants (id, membership_id, tenant_id, restaurant_id, role_id) VALUES ('00000000-0000-4000-8000-000000000042', $1, 1, 22, 222)`, membershipID); err == nil {
		t.Fatal("role grants must reject a restaurant outside the membership's explicit grant")
	}

	if _, err := db.ExecContext(ctx, `INSERT INTO identity_email_challenges (id, identity_id, email, purpose, token_hash, expires_at) VALUES ($1, $2, 'owner@example.test', 'verify_email', decode(repeat('ab', 32), 'hex'), CURRENT_TIMESTAMP + INTERVAL '10 minutes')`, challengeID, identityID); err != nil {
		t.Fatalf("insert synthetic verification challenge: %v", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO identity_email_challenges (id, identity_id, email, purpose, token_hash, expires_at) VALUES ('00000000-0000-4000-8000-000000000051', $1, 'different@example.test', 'verify_email', decode(repeat('cd', 32), 'hex'), CURRENT_TIMESTAMP + INTERVAL '10 minutes')`, identityID); err == nil {
		t.Fatal("email challenges must bind to the identity's exact normalized email")
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO identity_email_challenges (id, identity_id, email, purpose, token_hash, expires_at, attempt_count) VALUES ('00000000-0000-4000-8000-000000000052', $1, 'owner@example.test', 'verify_email', decode(repeat('ef', 32), 'hex'), CURRENT_TIMESTAMP + INTERVAL '10 minutes', 11)`, identityID); err == nil {
		t.Fatal("challenge attempt count must remain within the database limit")
	}
}

func assertIdentitySchemaRuntimeRoleDenied(t *testing.T, ctx context.Context, rawURL, schema string) {
	t.Helper()
	appURL := withSearchPath(t, rawURL, schema)
	app := openAuthAuditURL(t, appURL)
	app.SetMaxOpenConns(1)
	assertRestrictedRuntimeRole(t, ctx, app)
	if _, err := app.ExecContext(ctx, `SET search_path TO `+schema); err != nil {
		t.Fatalf("set runtime schema search path: %v", err)
	}
	var count int
	if err := app.QueryRowContext(ctx, `SELECT count(*) FROM auth_identities`).Scan(&count); err == nil {
		t.Fatal("pos_app must not read global identity data before a least-privilege auth boundary exists")
	}
	if _, err := app.ExecContext(ctx, `INSERT INTO auth_identities (id, normalized_email, password_hash) VALUES ('00000000-0000-4000-8000-000000000099', 'runtime@example.test', 'synthetic-hash')`); err == nil {
		t.Fatal("pos_app must not write global identity data")
	}
}

func assertIdentitySchemaRollback(t *testing.T, ctx context.Context, db *sql.DB, schema string) {
	t.Helper()
	rollback, err := os.ReadFile(filepath.Join("..", "..", "migrations", "191_auth_identity_membership_schema.rollback"))
	if err != nil {
		t.Fatalf("read identity schema rollback: %v", err)
	}
	var preservedIdentities int
	func() {
		conn, err := db.Conn(ctx)
		if err != nil {
			t.Fatalf("acquire connection for protected rollback: %v", err)
		}
		defer conn.Close()
		if _, err := conn.ExecContext(ctx, string(rollback)); err == nil {
			t.Fatal("rollback must refuse to discard populated identity data")
		}
		if _, err := conn.ExecContext(ctx, `ROLLBACK`); err != nil {
			t.Fatalf("clear failed protected rollback transaction: %v", err)
		}
		if err := conn.QueryRowContext(ctx, `SELECT count(*) FROM auth_identities`).Scan(&preservedIdentities); err != nil {
			t.Fatalf("confirm protected rollback preserved identity data: %v", err)
		}
	}()
	if preservedIdentities == 0 {
		t.Fatal("protected rollback removed identity data after its guard failed")
	}
	for _, table := range []string{"membership_role_grants", "membership_restaurant_grants", "identity_email_challenges", "tenant_memberships", "auth_identities"} {
		if _, err := db.ExecContext(ctx, `DELETE FROM `+table); err != nil {
			t.Fatalf("clear synthetic %s rows before rollback: %v", table, err)
		}
	}
	if _, err := db.ExecContext(ctx, string(rollback)); err != nil {
		t.Fatalf("rollback empty identity schema: %v", err)
	}
	for _, object := range []string{"auth_identities", "identity_email_challenges", "tenant_memberships", "membership_restaurant_grants", "membership_role_grants", "idx_users_phase1_tenant_id", "idx_restaurants_phase1_tenant_id", "idx_roles_phase1_scope_id"} {
		var exists bool
		if err := db.QueryRowContext(ctx, `SELECT to_regclass($1) IS NOT NULL`, schema+"."+object).Scan(&exists); err != nil {
			t.Fatalf("inspect rolled-back object %s: %v", object, err)
		}
		if exists {
			t.Errorf("rollback left identity schema object %s behind", object)
		}
	}
}
