# WS01 — Database Least Privilege Report

**Status:** VERIFIED COMPLETE
**Date:** 2026-11-02
**Scope:** Local `pos-dev` development database and API startup only. No production system was contacted.

## Initial state

The development API was connecting as PostgreSQL's superuser. The migration runner and API shared the same database identity. Existing dev data and the `pos-dev` Docker volumes were treated as persistent owner data.

## Changes implemented

- Added `backend/scripts/db/transfer_public_ownership.sql` to transfer eligible application objects in `public` to the migration owner while leaving extension-owned objects untouched.
- Updated `backend/docker-entrypoint.dev.sh` so migrations use `MIGRATE_DB_USER`/`MIGRATE_DB_PASSWORD`, then the API starts with its separate runtime identity. The API role is explicitly denied access to `schema_migrations`.
- Configured `docker-compose.local.yml` to use local environment credentials for `pos_app` and `pos_migrator`. Generated values were stored only in ignored `.env.local`; existing values were preserved.
- The database and schema, 182 relations, and 31 functions are owned by `pos_migrator`. The API connects as `pos_app`.

## Recovery evidence

Before changing live development ownership, a custom-format backup was created outside the repository:

- File: `C:\Users\yosea\AppData\Local\POS\phase2a-ws01\pos_saas-before-ws01.dump`
- Size: 651,670 bytes
- SHA-256: `619DA9C698D6C5712B689F00C29322A110209ED446422BC43D4E6D63AAF8E24B`
- `pg_restore --list` completed successfully.
- The backup was restored into an isolated PostgreSQL 15.19 container. The restored database contained 91 application tables, 108 migration records, 1 tenant, 1 restaurant and 3 users.
- Ownership transfer and reverse transfer were both rehearsed on the isolated copy. The rehearsal restored the database/schema and eligible public objects to `postgres` ownership.

The backup is a recovery artifact, not a command to overwrite the running database. To inspect or restore it into a **new isolated** database, use the matching PostgreSQL 15 client and an isolated target:

```powershell
pg_restore --list 'C:\Users\yosea\AppData\Local\POS\phase2a-ws01\pos_saas-before-ws01.dump'
createdb -h <isolated-postgres-host> -p <isolated-port> -U <restore-owner> pos_saas_recovery
pg_restore --exit-on-error --no-owner --dbname='postgresql://<restore-owner>:<password>@<isolated-postgres-host>:<isolated-port>/pos_saas_recovery' 'C:\Users\yosea\AppData\Local\POS\phase2a-ws01\pos_saas-before-ws01.dump'
```

For an ownership rollback rehearsal, apply `backend/scripts/db/transfer_public_ownership.sql` with its `pos_migrator` setting changed to the original owner (`postgres`) against the isolated recovery DB. Do not run the restore or reverse-ownership operation against `pos-dev` without a separately reviewed recovery plan.

## Verification

### Isolated database

- Fresh database migrations completed using `pos_migrator`.
- Existing restored database started under production and development entrypoints with API/migration role separation.
- `pos_app` is `NOSUPERUSER` and `NOBYPASSRLS`.
- Application DML succeeded in a rolled-back transaction.
- Table/schema DDL, `TRUNCATE`, role creation/escalation and migration-history access were denied to `pos_app`.
- API health and public menu read returned HTTP 200.

### Existing development database

- Live catalog verified `pos_app` and `pos_migrator` are non-superusers and `NOBYPASSRLS`.
- `pos_migrator` owns the database/schema and application relations/functions; only `pos_app` has an API database connection.
- `pos_app` cannot read `schema_migrations`.
- Counts remained 1 tenant, 1 restaurant, 3 users, 91 tables and 108 migration records.
- Existing `pos-dev` volumes were preserved. Only POS services were restarted.
- API health and the landing, dashboard, storefront and Super Admin routes returned HTTP 200.

### Regression gates

- `go build ./...`: PASS.
- `go vet ./...`: PASS.
- `go test -count=1 ./...`: PASS.
- Disposable `scripts/test-stack/security-smoke.sh`: 73 passed, 0 failed.
- These Go checks ran in network-isolated containers with read-only source mounts. The security smoke ran against the disposable `pos-test` stack, not the dev database.

## Files changed

- `backend/scripts/db/transfer_public_ownership.sql` (new)
- `backend/docker-entrypoint.dev.sh`
- `docker-compose.local.yml`
- Ignored local `.env.local` (role credentials appended; secret values intentionally omitted here)
- `docs-organized/phase-2a/execution/PHASE_2A_EXECUTION_TRACKER.md`

## Security impact and limits

The API runtime no longer has PostgreSQL superuser or RLS-bypass capability, and migration-history access is separated. No RLS policies were enabled. This does not establish tenant isolation for every application query; that is WS02. The owner must keep the private backup artifact in a safe location until it is no longer needed.

## Closure

All WS01 exit criteria in the execution brief passed for the current local development environment. No commit, push, branch switch or deployment was performed.
