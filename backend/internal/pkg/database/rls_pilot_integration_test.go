//go:build integration

package database_test

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"pos-saas/internal/pkg/database"
)

// Integration pilot for staged RLS (Phase 2A C012). Runs only with `-tags integration` against a
// DISPOSABLE database:
//   POS_RLS_ADMIN_URL  superuser/owner URL (creates the pilot table and policy)
//   POS_RLS_APP_URL    restricted runtime role URL (pos_app: NOSUPERUSER NOBYPASSRLS)
// Both database names must contain "test".

func openDB(t *testing.T, env string) *sql.DB {
	t.Helper()
	url := os.Getenv(env)
	if url == "" {
		t.Skipf("%s not set", env)
	}
	if !strings.Contains(strings.ToLower(url), "test") {
		t.Fatalf("refusing to run: %s does not point at a test database", env)
	}
	db, err := sql.Open("postgres", url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestRLSPilot_TenantIsolationWithRestrictedRole(t *testing.T) {
	admin := openDB(t, "POS_RLS_ADMIN_URL")
	app := openDB(t, "POS_RLS_APP_URL")
	app.SetMaxOpenConns(1) // force connection reuse: proves context cannot leak through the pool
	ctx := context.Background()

	for _, stmt := range []string{
		`DROP TABLE IF EXISTS rls_pilot_items`,
		`CREATE TABLE rls_pilot_items (id serial PRIMARY KEY, tenant_id bigint NOT NULL, name text NOT NULL)`,
		`ALTER TABLE rls_pilot_items ENABLE ROW LEVEL SECURITY`,
		`ALTER TABLE rls_pilot_items FORCE ROW LEVEL SECURITY`,
		`CREATE POLICY tenant_isolation ON rls_pilot_items
		   USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint)
		   WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::bigint)`,
		`GRANT SELECT, INSERT, UPDATE, DELETE ON rls_pilot_items TO pos_app`,
		`GRANT USAGE, SELECT ON SEQUENCE rls_pilot_items_id_seq TO pos_app`,
	} {
		if _, err := admin.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("setup %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() { _, _ = admin.Exec(`DROP TABLE IF EXISTS rls_pilot_items`) })

	insert := func(tenant int64, name string) error {
		return database.WithTenantTx(ctx, app, tenant, 0, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, `INSERT INTO rls_pilot_items (tenant_id, name) VALUES ($1, $2)`, tenant, name)
			return err
		})
	}
	names := func(tenant int64) []string {
		var out []string
		err := database.WithTenantTx(ctx, app, tenant, 0, func(tx *sql.Tx) error {
			rows, err := tx.QueryContext(ctx, `SELECT name FROM rls_pilot_items ORDER BY id`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var n string
				if err := rows.Scan(&n); err != nil {
					return err
				}
				out = append(out, n)
			}
			return rows.Err()
		})
		if err != nil {
			t.Fatalf("read as tenant %d: %v", tenant, err)
		}
		return out
	}

	if err := insert(1, "a-1"); err != nil {
		t.Fatal(err)
	}
	if err := insert(1, "a-2"); err != nil {
		t.Fatal(err)
	}
	if err := insert(2, "b-1"); err != nil {
		t.Fatal(err)
	}

	if got := names(1); len(got) != 2 || got[0] != "a-1" || got[1] != "a-2" {
		t.Fatalf("tenant 1 must see only its rows, got %v", got)
	}
	if got := names(2); len(got) != 1 || got[0] != "b-1" {
		t.Fatalf("tenant 2 must see only its rows, got %v", got)
	}

	// Without a tenant context (plain connection, as any code that forgets WithTenantTx would use) nothing is visible.
	var n int
	if err := app.QueryRowContext(ctx, `SELECT count(*) FROM rls_pilot_items`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("no tenant context must fail closed (0 rows), saw %d (context leaked through the pool?)", n)
	}

	// Writing a row for another tenant is rejected by WITH CHECK.
	err := database.WithTenantTx(ctx, app, 1, 0, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO rls_pilot_items (tenant_id, name) VALUES (2, 'smuggled')`)
		return err
	})
	if err == nil {
		t.Fatal("cross-tenant insert must violate the policy")
	}

	// Updates and deletes cannot reach another tenant's rows.
	var affected int64
	_ = database.WithTenantTx(ctx, app, 1, 0, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `UPDATE rls_pilot_items SET name = 'hijacked' WHERE tenant_id = 2`)
		if err != nil {
			return err
		}
		affected, _ = res.RowsAffected()
		return nil
	})
	if affected != 0 {
		t.Fatalf("tenant 1 updated %d of tenant 2's rows", affected)
	}
	_ = database.WithTenantTx(ctx, app, 1, 0, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM rls_pilot_items WHERE tenant_id = 2`)
		if err == nil {
			affected, _ = res.RowsAffected()
		}
		return err
	})
	if affected != 0 {
		t.Fatalf("tenant 1 deleted %d of tenant 2's rows", affected)
	}
	if got := names(2); len(got) != 1 || got[0] != "b-1" {
		t.Fatalf("tenant 2's data must be intact, got %v", got)
	}

	// A superuser/BYPASSRLS role would bypass the policy: that is why the runtime role must not be one.
	var bypass bool
	if err := app.QueryRowContext(ctx, `SELECT rolsuper OR rolbypassrls FROM pg_roles WHERE rolname = current_user`).Scan(&bypass); err != nil {
		t.Fatal(err)
	}
	if bypass {
		t.Fatal("the runtime role must not be superuser or BYPASSRLS")
	}
}
