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

func TestAuthLookupBoundaryIsReadOnlyAndReturnsOnlyEligibleMemberships(t *testing.T) {
	admin := openAuthAuditTestDB(t, "POS_RLS_ADMIN_URL")
	ctx := context.Background()
	admin.SetMaxOpenConns(1)
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	schema := "auth_lookup_test_" + suffix
	privateSchema := "auth_lookup_private_" + suffix
	role := "auth_lookup_role_" + suffix
	identifier := regexp.MustCompile(`^[a-z][a-z0-9_]+$`)
	for _, value := range []string{schema, privateSchema, role} {
		if !identifier.MatchString(value) {
			t.Fatalf("generated unsafe PostgreSQL identifier %q", value)
		}
	}
	if _, err := admin.ExecContext(ctx, `CREATE SCHEMA `+schema); err != nil {
		t.Fatalf("create disposable identity schema: %v", err)
	}
	t.Cleanup(func() {
		if _, err := admin.ExecContext(context.Background(), `DROP SCHEMA `+schema+` CASCADE`); err != nil {
			t.Errorf("drop disposable identity schema: %v", err)
		}
		if _, err := admin.ExecContext(context.Background(), `DROP SCHEMA IF EXISTS `+privateSchema+` CASCADE`); err != nil {
			t.Errorf("drop disposable auth function schema: %v", err)
		}
	})
	if _, err := admin.ExecContext(ctx, `SET search_path TO `+schema); err != nil {
		t.Fatalf("set disposable schema search path: %v", err)
	}
	createIdentitySchemaFixture(t, ctx, admin, schema)
	for _, statement := range []string{
		`ALTER TABLE ` + schema + `.tenants OWNER TO pos_migrator`,
		`ALTER TABLE ` + schema + `.users OWNER TO pos_migrator`,
		`ALTER TABLE ` + schema + `.restaurants OWNER TO pos_migrator`,
		`ALTER TABLE ` + schema + `.roles OWNER TO pos_migrator`,
		`ALTER SCHEMA ` + schema + ` OWNER TO pos_migrator`,
		`SET ROLE pos_migrator`,
	} {
		if _, err := admin.ExecContext(ctx, statement); err != nil {
			t.Fatalf("prepare migration-owned disposable schema: %v", err)
		}
	}

	identitySchemaMigration, err := os.ReadFile(filepath.Join("..", "..", "migrations", "191_auth_identity_membership_schema.sql"))
	if err != nil {
		t.Fatalf("read identity schema migration: %v", err)
	}
	if _, err := admin.ExecContext(ctx, string(identitySchemaMigration)); err != nil {
		t.Fatalf("apply identity schema migration as pos_migrator: %v", err)
	}
	lookupMigration, err := os.ReadFile(filepath.Join("..", "..", "migrations", "192_auth_lookup_boundary.sql"))
	if err != nil {
		t.Fatalf("read auth lookup migration: %v", err)
	}
	rewrittenLookupMigration := strings.NewReplacer(
		"public.", schema+".",
		"auth_private", privateSchema,
		"pos_auth_lookup", role,
	).Replace(string(lookupMigration))
	if _, err := admin.ExecContext(ctx, rewrittenLookupMigration); err != nil {
		t.Fatalf("apply auth lookup migration as pos_migrator: %v", err)
	}
	if _, err := admin.ExecContext(ctx, rewrittenLookupMigration); err != nil {
		t.Fatalf("reapply auth lookup migration without role provisioning: %v", err)
	}
	if _, err := admin.ExecContext(ctx, `RESET ROLE`); err != nil {
		t.Fatalf("reset migration role: %v", err)
	}
	if _, err := admin.ExecContext(ctx, `CREATE ROLE `+role+` NOLOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS NOINHERIT`); err != nil {
		t.Fatalf("create disposable lookup role after migration: %v", err)
	}
	t.Cleanup(func() {
		if _, err := admin.ExecContext(context.Background(), `DROP OWNED BY `+role); err != nil {
			t.Errorf("drop disposable lookup role grants: %v", err)
		}
		if _, err := admin.ExecContext(context.Background(), `DROP ROLE `+role); err != nil {
			t.Errorf("drop disposable lookup role: %v", err)
		}
	})
	for _, statement := range []string{
		`GRANT USAGE ON SCHEMA ` + schema + ` TO ` + role,
		`GRANT USAGE ON SCHEMA ` + privateSchema + ` TO ` + role,
		`GRANT EXECUTE ON FUNCTION ` + privateSchema + `.lookup_verified_identity_by_email(TEXT) TO ` + role,
		`GRANT EXECUTE ON FUNCTION ` + privateSchema + `.lookup_eligible_memberships(UUID) TO ` + role,
	} {
		if _, err := admin.ExecContext(ctx, statement); err != nil {
			t.Fatalf("grant restricted lookup role capability: %v", err)
		}
	}

	seedAndAssertIdentitySchema(t, ctx, admin)
	seedAdditionalLookupMemberships(t, ctx, admin, schema)
	assertLookupFunctionOwnership(t, ctx, admin, privateSchema)
	assertAuthLookupRoleResults(t, ctx, admin, schema, privateSchema, role)
	assertAuthLookupRollback(t, ctx, admin, schema, privateSchema)
}

func seedAdditionalLookupMemberships(t *testing.T, ctx context.Context, db *sql.DB, schema string) {
	t.Helper()
	if _, err := db.ExecContext(ctx, `SET ROLE pos_migrator`); err != nil {
		t.Fatalf("assume migration role for additional auth fixtures: %v", err)
	}
	defer func() {
		if _, err := db.ExecContext(context.Background(), `RESET ROLE`); err != nil {
			t.Errorf("reset role after additional auth fixtures: %v", err)
		}
	}()
	for _, statement := range []string{
		`INSERT INTO ` + schema + `.tenant_memberships (id, identity_id, tenant_id, legacy_user_id, state, accepted_at) VALUES ('00000000-0000-4000-8000-000000000041', '00000000-0000-4000-8000-000000000001', 2, 202, 'active', CURRENT_TIMESTAMP)`,
		`INSERT INTO ` + schema + `.membership_restaurant_grants (id, membership_id, tenant_id, restaurant_id, state) VALUES ('00000000-0000-4000-8000-000000000042', '00000000-0000-4000-8000-000000000041', 2, 22, 'active')`,
		`INSERT INTO ` + schema + `.membership_role_grants (id, membership_id, tenant_id, restaurant_id, role_id, state) VALUES ('00000000-0000-4000-8000-000000000043', '00000000-0000-4000-8000-000000000041', 2, 22, 222, 'active')`,
		`INSERT INTO ` + schema + `.tenants (id, name) VALUES (3, 'Pending Tenant')`,
		`INSERT INTO ` + schema + `.users (id, tenant_id) VALUES (303, 3)`,
		`INSERT INTO ` + schema + `.restaurants (id, tenant_id, name) VALUES (33, 3, 'Pending Restaurant')`,
		`INSERT INTO ` + schema + `.roles (id, tenant_id, restaurant_id, role_code) VALUES (333, 3, 33, 'PENDING')`,
		`INSERT INTO ` + schema + `.tenant_memberships (id, identity_id, tenant_id, legacy_user_id, state) VALUES ('00000000-0000-4000-8000-000000000032', '00000000-0000-4000-8000-000000000001', 3, 303, 'pending')`,
		`INSERT INTO ` + schema + `.membership_restaurant_grants (id, membership_id, tenant_id, restaurant_id, state) VALUES ('00000000-0000-4000-8000-000000000033', '00000000-0000-4000-8000-000000000032', 3, 33, 'active')`,
		`INSERT INTO ` + schema + `.membership_role_grants (id, membership_id, tenant_id, restaurant_id, role_id, state) VALUES ('00000000-0000-4000-8000-000000000034', '00000000-0000-4000-8000-000000000032', 3, 33, 333, 'active')`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("seed pending synthetic membership: %v", err)
		}
	}
}

func assertLookupFunctionOwnership(t *testing.T, ctx context.Context, db *sql.DB, schema string) {
	t.Helper()
	for _, function := range []string{
		"lookup_verified_identity_by_email(text)",
		"lookup_eligible_memberships(uuid)",
	} {
		var owner string
		var securityDefiner, pinnedSearchPath bool
		query := `
			SELECT owner.rolname, p.prosecdef,
			       p.proconfig @> ARRAY['search_path=pg_catalog']::text[]
			FROM pg_proc p
			JOIN pg_namespace n ON n.oid = p.pronamespace
			JOIN pg_roles owner ON owner.oid = p.proowner
			WHERE n.nspname = $1 AND p.oid = ($1 || '.' || $2)::regprocedure
		`
		if err := db.QueryRowContext(ctx, query, schema, function).Scan(&owner, &securityDefiner, &pinnedSearchPath); err != nil {
			t.Fatalf("inspect lookup function %s: %v", function, err)
		}
		if owner != "pos_migrator" || !securityDefiner || !pinnedSearchPath {
			t.Errorf("lookup function %s must be SECURITY DEFINER, migration-owned, and search-path pinned; got owner=%q definer=%v pinned=%v", function, owner, securityDefiner, pinnedSearchPath)
		}
	}
}

func assertAuthLookupRoleResults(t *testing.T, ctx context.Context, db *sql.DB, dataSchema, functionSchema, role string) {
	t.Helper()
	if _, err := db.ExecContext(ctx, `SET ROLE `+role); err != nil {
		t.Fatalf("assume disposable auth lookup role: %v", err)
	}
	defer func() {
		if _, err := db.ExecContext(context.Background(), `RESET ROLE`); err != nil {
			t.Errorf("reset disposable auth lookup role: %v", err)
		}
	}()

	var identityID, email, passwordHash string
	var version int64
	if err := db.QueryRowContext(ctx, `
		SELECT identity_id, normalized_email, password_hash, security_version
		FROM `+functionSchema+`.lookup_verified_identity_by_email($1)
	`, " OWNER@EXAMPLE.TEST ").Scan(&identityID, &email, &passwordHash, &version); err != nil {
		t.Fatalf("lookup verified synthetic identity: %v", err)
	}
	if identityID != "00000000-0000-4000-8000-000000000001" || email != "owner@example.test" || passwordHash != "synthetic-hash" || version != 1 {
		t.Fatalf("unexpected synthetic identity lookup: id=%q email=%q hash=%q version=%d", identityID, email, passwordHash, version)
	}

	rows, err := db.QueryContext(ctx, `
		SELECT membership_id, tenant_id, restaurant_id, role_id, role_code
		FROM `+functionSchema+`.lookup_eligible_memberships($1::uuid)
		ORDER BY tenant_id
	`, identityID)
	if err != nil {
		t.Fatalf("lookup eligible memberships: %v", err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var membershipID, roleCode string
		var tenantID, restaurantID, roleID int
		if err := rows.Scan(&membershipID, &tenantID, &restaurantID, &roleID, &roleCode); err != nil {
			t.Fatalf("read eligible membership: %v", err)
		}
		got = append(got, fmt.Sprintf("%d/%d/%d/%s", tenantID, restaurantID, roleID, roleCode))
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate eligible memberships: %v", err)
	}
	if len(got) != 2 || got[0] != "1/11/111/TEST" || got[1] != "2/22/222/TEST" {
		t.Fatalf("eligible memberships = %v, want both active explicit memberships and no pending tenant 3", got)
	}

	if _, err := db.ExecContext(ctx, `SELECT count(*) FROM `+dataSchema+`.auth_identities`); err == nil {
		t.Fatal("auth lookup role must not directly read identity tables")
	}
	if _, err := db.ExecContext(ctx, `UPDATE `+dataSchema+`.auth_identities SET status='revoked'`); err == nil {
		t.Fatal("auth lookup role must not write identity tables")
	}
}

func assertAuthLookupRollback(t *testing.T, ctx context.Context, db *sql.DB, dataSchema, schema string) {
	t.Helper()
	rollback, err := os.ReadFile(filepath.Join("..", "..", "migrations", "192_auth_lookup_boundary.rollback"))
	if err != nil {
		t.Fatalf("read auth lookup rollback: %v", err)
	}
	if _, err := db.ExecContext(ctx, `SET ROLE pos_migrator`); err != nil {
		t.Fatalf("assume migration role for auth lookup rollback: %v", err)
	}
	rewrittenRollback := strings.ReplaceAll(string(rollback), "auth_private", schema)
	if _, err := db.ExecContext(ctx, rewrittenRollback); err != nil {
		_, _ = db.ExecContext(ctx, `ROLLBACK`)
		t.Fatalf("rollback auth lookup functions without touching identity tables: %v", err)
	}
	if _, err := db.ExecContext(ctx, `RESET ROLE`); err != nil {
		t.Fatalf("reset role after auth lookup rollback: %v", err)
	}
	var schemaExists bool
	if err := db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname=$1)`, schema).Scan(&schemaExists); err != nil {
		t.Fatalf("inspect auth lookup rollback schema: %v", err)
	}
	if schemaExists {
		t.Fatal("auth lookup rollback left its private function schema behind")
	}
	var identityCount int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM `+dataSchema+`.auth_identities`).Scan(&identityCount); err != nil {
		t.Fatalf("confirm auth lookup rollback preserved identity schema data: %v", err)
	}
	if identityCount == 0 {
		t.Fatal("auth lookup rollback modified identity data")
	}
}
