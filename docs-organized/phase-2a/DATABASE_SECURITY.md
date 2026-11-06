# Database security model (Phase 2A C010 / C012)

## Identities

| Role | Purpose | Attributes | Privileges |
|---|---|---|---|
| `postgres` (bootstrap) | One-time setup only (create roles, database ownership). **Not used by the API.** | superuser | everything |
| `pos_migrator` | Owns the database and `public` schema; runs `cmd/migrate` | `NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS` | DDL on its own schema |
| `pos_app` | The API's runtime connection (`DB_USER`) | `NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS` | `CONNECT`, `USAGE` on `public`, `SELECT/INSERT/UPDATE/DELETE` on tables, `USAGE/SELECT` on sequences, `EXECUTE` on functions. No DDL, no `TRUNCATE`, no access to `schema_migrations` |

Created by `backend/scripts/db/create_runtime_roles.sql` (idempotent; takes the database name and both passwords as psql variables). `backend/docker-entrypoint.sh` runs migrations with `MIGRATE_DB_USER/MIGRATE_DB_PASSWORD` when set, then revokes `schema_migrations` from the runtime role and starts the API with `DB_USER`. Default privileges (`ALTER DEFAULT PRIVILEGES FOR ROLE pos_migrator`) make every table the migrator creates in future usable by `pos_app` with DML only.

## Verified (disposable stack `docker-compose.pos-test.yml`, `scripts/test-stack/security-smoke.sh`)
- `pos_app` attribute flags all false; `pos_migrator` is not a superuser and cannot create roles.
- The API holds network connections as `pos_app` and none as `postgres`.
- As `pos_app`: reads work; `CREATE TABLE`, `DROP TABLE`, `ALTER TABLE`, `TRUNCATE`, `CREATE ROLE`, `DROP SCHEMA`, `CREATE FUNCTION`, reading/updating `schema_migrations`, and `COPY ... TO PROGRAM` are all denied.
- All 73 smoke checks (auth, tenancy, RBAC, themes, orders, subscription status, ...) pass with the API running as `pos_app`. This is smoke coverage, not exhaustive coverage of every table/endpoint.

## Not yet applied
- **Production/dev compose files still connect as the bootstrap user.** `docker-compose.yml` (dev) and `docker-compose.prod.yml` were not switched; migration 050 grants to the role `postgres`, so existing databases must also `REASSIGN OWNED BY <old owner> TO pos_migrator` first. Needs an upgrade procedure and a rehearsal on a copy.
- No TLS requirement on the database connection (`DB_SSLMODE=disable` in all compose files).

## Row-level security — staged, NOT enabled on any live table
- Classification of all 91 tables: `DATABASE_TENANCY_MATRIX.md` (45 restaurant-owned, 8 tenant-owned, 12 tenant-owned via FK, 4+7+2 audit-flavoured tenant tables, 3 integration, 3 platform, 5 reference, `tenants`, and 1 needing review: `webhook_retry_queue`).
- `internal/pkg/database.WithTenantTx` runs work in a transaction with `app.tenant_id` / `app.restaurant_id` set via `set_config(..., is_local => true)`, so a pooled connection cannot carry a tenant into the next request. Unit-tested with sqlmock (6 tests).
- `internal/pkg/database/rls_pilot_integration_test.go` (`-tags integration`, needs `POS_RLS_ADMIN_URL` and `POS_RLS_APP_URL`) proves the generic pilot policy against the restricted role and a **single-connection pool**. `internal/repository/auth_audit_repository_integration_test.go` runs the actual migrations 163/164 in an isolated schema and checks the auth-audit table, tenant/user ownership, writes, invalid/missing context, role restrictions, and connection reuse. Both tests passed against disposable PostgreSQL 16 on 2026-11-06; the auth-audit verification is recorded at commit `fbb247b099b0ce86c947907538b36ab0a2779b29`.
- **Why not enabled broadly on live tables:** the 363-route API reads and writes through `*sql.DB` directly (repositories take the pool, ~123 endpoints skip use cases, 40 inline SQL closures). With RLS on, any query that does not run inside `WithTenantTx` returns no rows or fails. The auth-audit table is the first migration-backed slice; enabling additional tables requires moving each module's data access to tenant transactions first, family by family, with policies and regression tests added in the same change. Order suggested in `PHASE_2B_HANDOFF.md`.
- The older `internal/security` RLS manager and `rls_policies.sql` (never applied, schema-incompatible) should be deleted or rewritten, not revived: its tests are quarantined in `backend/tests/EXPERIMENTAL.md`.
- Gate G04/G05 (tenant and restaurant isolation) stays **PARTIAL**: application-level isolation is tested; database-level enforcement is not enabled.
