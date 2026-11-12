# Database and tenancy architecture audit

## Live development observation

Read-only SQL against `pos-postgres`, database `pos_saas`, on 2026-11-02 reported PostgreSQL 15.19, 91 public base tables, 0 tables with row-level security enabled, and 108 rows in `schema_migrations`. No credentials are included here. The current dev compose defines API database user `postgres`; role query on the live database returned only `postgres` among `postgres`, `pos_app`, `pos_migrator`, with `postgres` superuser and BYPASSRLS. This does not describe production.

## Ownership model

The checked-in matrix classifies 91 tables heuristically: 45 restaurant-owned; 8 tenant-owned; 12 tenant-owned through foreign keys; 4 restaurant-owned audit, 7 tenant-owned derived audit, 2 tenant-owned audit, 3 integration-flavoured, 3 platform-admin, 5 reference, one global tenant registry, and one global system queue (`webhook_retry_queue`). The retry queue is deliberately cross-tenant worker state: its schema has no tenant/restaurant key or ownership FK, and the retry worker claims callbacks across tenants. Do not add tenant RLS without first changing the queue contract to persist and enforce ownership. Keep tenant application handlers away from queue reads, and review stored provider-payload retention and access. The matrix remains an inventory rather than a complete verified policy map. See `docs-organized/phase-2a/DATABASE_TENANCY_MATRIX.md` for table-level keys and parent FK paths.

| Area | Current state | Risk |
|---|---|---|
| Tenant/restaurant keys | Many business tables have tenant_id and/or restaurant_id; derived tables rely on FK parent chains | A key column alone does not ensure every query scopes both tenant and restaurant. |
| Query isolation | Middleware derives tenant from validated identity and restaurant verifier; repositories and handler SQL are mixed | Object-ID access and internal call paths need systematic tenant/restaurant negative tests. |
| Runtime privilege | Dev API uses `postgres`; disposable stack proves `pos_app` can run DML with restricted role and migrations use `pos_migrator` | Dev/prod role rollout and ownership migration remain incomplete. |
| RLS | No live application tables have RLS enabled; `WithTenantTx` helper uses transaction-local settings; isolated pilot passed per prior report | Direct `*sql.DB` paths are not migrated; enabling RLS now could break legitimate reads/writes and is not safe. |
| Migrations | Dev startup runs migrations; live database has 108 migration records | Migrations are state-changing at startup; audit did not invoke them. Duplicate legacy numbering and historical migration shape require review. |

## Required safe path

Inventory and review all ownership classes, including derived rows; establish tenant+restaurant scoped repository contracts; move modules to `WithTenantTx`; run two-tenant/two-restaurant negative integration tests; rehearse DDL ownership transfer on a disposable database copy; then enable RLS per table family with `FORCE ROW LEVEL SECURITY` and least-privilege runtime role. Never enable policies globally before all access paths are transaction-scoped.

## Verification limits

This audit did not mutate `pos_saas`, apply migrations, or inspect real customer data. The existing 73-check smoke is disposable-stack evidence and does not establish every table/query’s isolation. The disposable role setup is in `scripts/db/create_runtime_roles.sql`; `docker-compose.pos-test.yml` differs materially from development compose.
