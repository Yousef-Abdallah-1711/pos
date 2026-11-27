# 37 — Isolated Environment Execution Runbook (NOT EXECUTED)

Source: `docker-compose.phase2a-e2e.yml` (read in full, 218 lines). Nothing below has been run. Every step is status **NOT EXECUTED** and needs the owner's separate authorization before use.

## 1. Identity
| Item | Value (from file) |
|---|---|
| Compose project name | `pos-p2a-ws05-e2e` (set via `name:`; always pass `-p pos-p2a-ws05-e2e`) |
| Dev stack (must not be touched) | project `pos-dev`, ports 3001–3004 (+ API) |
| Network | own default network, subnet `10.253.255.0/24` (**check it does not overlap the `pos-dev` network subnet before starting**: `docker network inspect`) |
| Labels | `com.pos.disposable=true`, `com.pos.workstream=WS05` on postgres/redis/api |

## 2. Containers
| Service | Purpose | Notes |
|---|---|---|
| postgres (postgres:15-alpine) | DB `pos_phase2a_e2e`, superuser `postgres` | **tmpfs** `/var/lib/postgresql/data`: all data lost on container stop/restart |
| redis (redis:7-alpine) | cache/limits | no persistence (`--save ""`, `--appendonly no`) |
| init-roles | one-shot, runs `backend/scripts/db/create_runtime_roles.sql` | creates `pos_migrator`, `pos_app`, `pos_accounting` with passwords |
| api | Go API via `go run`; migrates as `pos_migrator`, then runs as `pos_app` | health: `/api/v1/health` on 18080, start_period 90 s |
| frontend-install | `pnpm install --frozen-lockfile` into named volumes | one-shot |
| landing / dashboard / restaurant-website / super-admin | Next dev servers | `network_mode: service:api` (share the API's network namespace) |

## 3. Ports (host loopback only, 127.0.0.1)
API 18082 → 18080 · Landing 13111 → 3001 · Dashboard 13112 → 3002 · Storefront 13113 → 3003 · Super Admin 13114 → 3004 · Postgres 15437 → 5432. No overlap with dev ports 3001–3004. Overridable through `P2A_*_HOST_PORT`. Redis is not published.

## 4. Database & roles
- DB name `pos_phase2a_e2e`; superuser password `ws05-local-only-postgres` (test-only value from the compose file).
- Runtime roles: `pos_app` (API; restricted), `pos_migrator` (migrations), `pos_accounting` (accounting pool via `ACCOUNTING_DB_USER/PASSWORD`).
- API `ENV=phase2a-e2e`. **This is not `production`/`prod`, so `IsProductionMode()` is false** (VER-020): production-only guards are NOT exercised by this stack. `RBAC_ENFORCEMENT=enforce` is set explicitly, `POS_TEST_SEED=true`.
- Source is bind-mounted (`./backend`, `./frontend`), the same working tree as the dev stack. Dev and isolated stacks therefore run identical code but separate containers, DB and `.next`/`node_modules` volumes.

## 5. Volume isolation and persistence
- DB and Redis are ephemeral. Evidence that must survive must be exported **before** any container restart: use `pg_dump` / `psql \copy` to `evidence/database/` and save API responses/screenshots to `evidence/`.
- Named volumes are prefixed `p2a_` (go caches, pnpm store, node_modules, .next) — safe to keep; they do not touch dev volumes.
- Backup/restore: `docker exec <pg> pg_dump -U postgres -Fc pos_phase2a_e2e > evidence/database/<stamp>.dump`; restore into a fresh tmpfs DB with `pg_restore` after migrations. A dump is the only way to keep state across a restart.

## 6. Resource constraint (owner decision)
Docker VM has 7.7 GiB; dev stack uses ~5 GiB (observed earlier). The compose has no memory limits. Starting both may exhaust memory. **Options for the owner:** (a) stop the dev stack (`docker compose -p pos-dev stop`, volumes kept) — needs approval, takes the 3001–3004 URLs down; (b) start only `postgres redis init-roles api` plus the dashboard + landing (and others later) to reduce RAM; (c) raise Docker Desktop memory. Recommended: (b) first, measure with `docker stats --no-stream`, then decide on (a).

## 7. Start-up sequence (NOT EXECUTED)
1. Pre-flight: `docker ps`, `docker network ls`, `docker network inspect` (subnet overlap), `docker system df`, confirm `git rev-parse HEAD` = `133260db…`, record in `evidence/logs/`.
2. `docker compose -p pos-p2a-ws05-e2e -f docker-compose.phase2a-e2e.yml config` (render only).
3. `… up -d postgres redis init-roles` → wait healthy / init-roles exit 0.
4. `… up -d api` → watch logs for "All migrations completed successfully" and health 200 (`curl http://127.0.0.1:18082/api/v1/health`). First build downloads Go modules; allow minutes.
5. `… up -d frontend-install` then the four apps; verify 13111–13114 return 200.
6. Verify roles: `psql … -c "select rolname,rolsuper,rolbypassrls from pg_roles"` expects `pos_app` without superuser/BYPASSRLS.
7. Verify `RBAC_ENFORCEMENT=enforce` in `docker exec <api> printenv`.
8. Verify accounting pool actually separate (log line from `accounting_db.go` / `pg_stat_activity` by usename).
9. Checkpoint 0: `pg_dump` the freshly migrated empty DB (golden baseline for resets).

## 8. Tenants, users, data (to be created through the UI during execution)
- Tenant A "Audit A": register at Landing 13111 (UI). Registration seeds **no COA** (VER-011) → accounting initialization step is part of the test (COA import/seed route or API; to be determined, **requires owner approval if it means a manual SQL seed**).
- Tenant B "Audit B": second registration for isolation tests.
- Roles/users per tenant: owner/admin, cashier, accountant, HR, restricted-viewer (created through the Users/Roles UI to also test RBAC).
- Super admin: `POS_TEST_SEED=true` may seed fixtures — contents NOT inspected yet (open item).
- Payments: no real Paymob credentials; mock handler `mock_paymob_handler.go` exists; customer orders have no card gateway (VER-010), so payments are recorded manually via the admin payment-status route.

## 9. Evidence handling
Per test: request/response JSON, screenshot, SQL result (`evidence/database/`), journal/TB export, log excerpt. File names `<test_id>_<n>.<ext>`; hash at end of phase (`sha256sum`).

## 10. Health checks (repeat before each phase)
API `/api/v1/health`; all four apps 200; `select count(*) from schema_migrations`; `select * from pg_roles` unchanged; outbox worker running (log line every 15 s); `docker stats` memory.

## 11. Teardown and restoring the original dev stack
- Isolated: `docker compose -p pos-p2a-ws05-e2e -f docker-compose.phase2a-e2e.yml down` (add `-v` only with approval; tmpfs DB is already ephemeral).
- Dev stack: if it was stopped, `docker compose -p pos-dev start` (or `up -d` from the dev compose) and verify 3001–3004 = 200, DB rows intact (tenant 1/2 counts), backend as before. Record before/after counts of tenants/users/orders to prove no change.

## 12. Changes that require additional owner approval
1. Stopping the `pos-dev` stack (RAM).
2. Starting any container / building images (several GB download and disk).
3. Any change to `docker-compose.phase2a-e2e.yml` or env (e.g. setting `ENV=production`, `staging`, or removing accounting vars) — needed for VER-019/020/021 tests; use a **copy** of the compose file in the scratchpad rather than editing the repo file.
4. Manual SQL seeding of COA/posting rules (if no UI/API route exists).
5. Fault injection (stopping postgres/api/redis in the isolated stack) — isolated only, but still a lifecycle action.
6. Seeding via `POS_TEST_SEED=true` content (inspect first).
7. Any use of real payment credentials — **not permitted**.

## 13. Known limitations of this environment
`ENV=phase2a-e2e` does not test production guards; tmpfs loses state on restart; dev servers (Next dev, `go run`) differ from production builds; Redis/rate limiting behaviour in a single instance only; clock/timezone same as host.
