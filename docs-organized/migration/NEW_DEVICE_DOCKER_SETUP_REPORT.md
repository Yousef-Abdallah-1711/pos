# New Device Docker Setup Report

**Result: DEVELOPMENT ENVIRONMENT FULLY RESTORED**
**Date:** 2026-11-02 (Africa/Cairo)
**Scope:** Local development only; no production deployment or source-code edits.

## A. Repository Verification

- Repository: `https://github.com/Yousef-Abdallah-1711/pos`
- Clone path: `H:\Projects\pos`
- `origin` fetch/push URLs point to the authoritative repository.
- Source history is present: 86 commits.
- Phase 1 and Phase 2A documentation, backend, frontend, migrations, and Compose files are present.
- The migration handoff directory did not exist in this clone, so it was created for these reports.

## B. Current Branch and Commit

- Branch: `main`, tracking `origin/main`.
- HEAD: `1815273924b8ade42f51863262b997115776322b` (`1815273`).
- Working tree has no source changes. The three setup documents in this directory are new and untracked; no commit or push was made.
- Some handoff documents describe an older checkout (`master` at `899dce6`). The current checkout and SHA above are authoritative for this setup.

## C. New Device Environment

| Item | Result |
|---|---|
| OS | Windows 10 Pro, build 19045, x64; WSL 2 active |
| CPU / RAM | AMD Ryzen 7 5700 / 31.9 GB visible RAM |
| Free disk | 138.2 GB on `H:` at inspection |
| Git | 2.54.0.windows.1 |
| Docker client / server | 29.7.2 / 29.7.2; Linux engine |
| Docker Compose | v5.5.0 |
| Host Node / pnpm | Node 24.16.0 / pnpm 11.28.9 |
| Go on host | Not installed; container toolchain used |
| Go in backend image | go1.24.13 linux/amd64 |
| Frontend container | Node 20.20.2, pnpm 9.0.0, Next.js 15.5.9 |

## D. Docker Build Results

- Development images built successfully: `pos-backend:dev`, the four `pos-dev-*` app images, and `pos-dev-frontend-install`.
- The frontend installer completed successfully with the frozen lockfile (941 packages, 10 workspace projects).
- Separate verification-only images built successfully from `frontend/Dockerfile.prod` with `NEXT_IGNORE_BUILD_ERRORS=false` and local URLs: landing, dashboard, restaurant website, and super admin. These images were not started as production services.
- Dockerfile versions match the workspace requirements; no dependency upgrades were made.

## Runtime Service Matrix

| Service | Container | Image | State / health evidence | Container port → host port | Network |
|---|---|---|---|---|---|
| PostgreSQL | `pos-postgres` | `postgres:15-alpine` | Running, healthy | 5432 → 127.0.0.1:15432 | `pos-dev_pos-network` |
| Redis | `pos-redis` | `redis:7-alpine` | Running, healthy; PING/PONG | 6379 → 127.0.0.1:16379 | `pos-dev_pos-network` |
| Go API | `pos-backend` | `pos-backend:dev` | Running; API health 200, DB healthy | 8080 → 127.0.0.1:8080 | `pos-dev_pos-network` |
| Landing | `pos-landing` | `pos-dev-landing` | Running; Next.js Ready | 3001 → 127.0.0.1:13001 | Shares backend network namespace |
| Dashboard | `pos-dashboard` | `pos-dev-dashboard` | Running; Next.js Ready | 3002 → 127.0.0.1:13002 | Shares backend network namespace |
| Restaurant Website | `pos-restaurant-website` | `pos-dev-restaurant-website` | Running; Next.js Ready | 3003 → 127.0.0.1:13003 | Shares backend network namespace |
| Super Admin | `pos-super-admin` | `pos-dev-super-admin` | Running; Next.js Ready | 3004 → 127.0.0.1:13004 | Shares backend network namespace |
| Frontend installer | one-shot `pos-dev-frontend-install-1` | `pos-dev-frontend-install` | Exited 0 as expected after dependency install | — | `pos-dev_pos-network` |
## E. PostgreSQL Status

- `pos-postgres`: running and healthy; PostgreSQL 15.19.
- Container port 5432 is published to `127.0.0.1:15432` to avoid an existing unrelated service on host port 5432.
- Persistent volume: `pos-dev_postgres_data` mounted at `/var/lib/postgresql/data`.

## F. Redis Status

- `pos-redis`: running and healthy; Redis 7.4.9.
- `PING` returned `PONG`; container port 6379 is published to `127.0.0.1:16379`.
- Persistent volume: `pos-dev_redis_data` mounted at `/data`.

## G. Database Migration Status

- The Go 1.24.13 container applied all 108 SQL files successfully, including migration 161 and legacy migration 999.
- The backend entrypoint reran the migration runner and skipped already-applied migrations.
- Database health is confirmed by `/api/v1/health` (`db_healthy=true`).

## H. Development Seed Status

- Database contains one development tenant, one restaurant, seeded themes, and no products/menu items.
- Local-only tenant slug/domain were aligned to `restaurant-1` / `restaurant-1.localhost` so the storefront middleware and restaurant slug resolve to the same fixture.
- Migration 109 creates a known default super-admin credential. That password was replaced in this local database with a newly generated random development password. The value is stored only in ignored `.env.local`; it is not included here. Super-admin login with that local credential returned HTTP 200.
- `backend/scripts/seed-test-users.sql` was not run because it references unrelated hard-coded test users and tenant ID 4.

## I. Go API Status

- `pos-backend`: running, hot reload active, listening on container port 8080 and published to `127.0.0.1:8080`.
- `/api/v1/health` returned HTTP 200, `status=ok`, `db_healthy=true`.
- `RBAC_ENFORCEMENT=enforce`; payment mode is `mock`.

## J. Landing Status

- `pos-landing`: running; Next.js reported Ready.
- `/en`, `/auth/login`, and `/auth/register` return HTTP 200 after locale routing.
- Known route defect: bare `/login` and `/register` redirect to `/en/login` and `/en/register`, which return 404. Localized implemented routes are `/en/auth/login` and `/en/auth/register`.

## K. Dashboard Status

- `pos-dashboard`: running; Next.js reported Ready; `/en` and `/en/dashboard` returned HTTP 200.
- Protected API requests still reject unauthenticated access (HTTP 401). The HTML response alone does not establish client-side route protection.

## L. Restaurant Website Status

- `pos-restaurant-website`: running; Next.js reported Ready.
- `http://restaurant-1.localhost:13003/en` returned HTTP 200.
- Tenant-status and public menu API requests for `restaurant-1` returned HTTP 200; the menu is empty because there are no product fixtures.

## M. Super Admin Status

- `pos-super-admin`: running; Next.js reported Ready; `/login` returned HTTP 200.
- `/api/v1/super-admin/tenants` rejected an unauthenticated request with HTTP 401.
- A local-only fresh super-admin password is saved in `.env.local`; its value is intentionally omitted here.

## N. Docker Networking

- Docker network: `pos-dev_pos-network` (bridge).
- Frontend containers use the backend container's network namespace as configured by the repository, so their `localhost:8080` API URL reaches the backend.
- Verified an HTTP health request from a frontend container to `http://localhost:8080/api/v1/health`.
- CORS preflight from `http://localhost:13001` returned HTTP 200 and an `Access-Control-Allow-Origin` header for that origin.
- Unrelated existing containers occupied host ports 3001, 5432, and 6379. They were not stopped or modified.

## O. API Connectivity

| Check | Result |
|---|---|
| API health and database | HTTP 200; healthy |
| Invalid login | HTTP 401 |
| Protected categories without bearer token | HTTP 401 |
| Super-admin API without token | HTTP 401 |
| Tenant status for `restaurant-1` | HTTP 200 |
| Public menu for `restaurant-1` | HTTP 200, empty product list |
| Public order status without tracking token | HTTP 404 (same not-found response used for unauthorized/missing order) |
| Fresh local super-admin login | HTTP 200; token returned, not logged in this report |

## P. Smoke Test Results

- All four frontend dev servers reported Ready; their local routes returned HTTP 200.
- API-to-PostgreSQL health, Redis PING, frontend-to-API connectivity, CORS, invalid-login rejection, tenant status, menu route, protected API rejection, and unauthenticated public-order rejection were exercised.
- No real payment or production credential was used.

## Q. Backend Build and Test Results

Executed in the repository's Go 1.24.13 container:

- `go build ./...` — PASS.
- `go vet ./...` — PASS.
- `go test -count=1 ./...` — PASS; all packages passed, matching the documented 680-test default baseline.
- Experimental suites were not run; `backend/tests/EXPERIMENTAL.md` says they are quarantined and expected to fail until their subsystems are repaired.

## R. Frontend Verification Results

App-code typecheck and separate production-mode Docker build (`NEXT_IGNORE_BUILD_ERRORS=false`) passed for all four applications.

| App | Unit test result |
|---|---|
| Landing | 4 tests passed; 3 suites fail to resolve the `@/stores/authStore` alias |
| Dashboard | 530 passed / 639 total; 109 failed across 15 suites, plus one test error |
| Restaurant website | 42 passed / 50 total; 8 failed across 3 suites, including test alias/store/color expectations |
| Super admin | 3 / 3 passed |

`typecheck:test` fails in all four apps on test-only TypeScript issues such as stale mock casts, writes to readonly `NODE_ENV`, stale component props, and Playwright/test typing. These failures were not changed or suppressed.

## S. Known Errors

- Landing bare `/login` and `/register` paths produce 404 after locale redirect; use the localized `/auth/*` routes.
- Dashboard Vitest has 109 failures; landing Vitest has alias-resolution failures; restaurant Jest has 8 failures. Super-admin's small suite passes.
- Dashboard build emits a critical dynamic-dependency warning from `packages/ui-components/registry/component-registry.ts`. Restaurant and super-admin builds emit React hook dependency warnings. All four production builds still succeed.
- Next.js logged transient Watchpack `ENOMEM` initial-scan warnings for dashboard and restaurant during startup. They later served routes successfully, and these warnings did not recur in the last 10 minutes checked.
- A successful super-admin login returned HTTP 200, but the API logged an audit insert failure (`invalid input syntax for type json`). Authentication succeeded; the audit record did not. This is an application defect to fix separately.
- Browserslist data in the frontend image is stale; it did not block builds.
- PostgreSQL also logged two failed diagnostic queries from setup checks (a query for a nonexistent `tenants.status` column and a psql variable attempt). Neither changed data.
- Database-level RLS is not enabled, and the dev API still uses the bootstrap PostgreSQL role; see `docs-organized/phase-2a/DATABASE_SECURITY.md`.

## T. Remaining Blockers

There is no blocker to running the development services. Phase 2A itself remains incomplete. The dashboard/landing/storefront unit-test failures, test-only type errors, super-admin audit insert error, and lack of database-level tenant enforcement remain follow-up engineering work.

## U. Local Changes Made

- Created ignored `.env.local` with new local-only database, JWT, and super-admin secrets; no values are included in reports.
- Created ignored `docker-compose.local.yml` with Redis, loopback-only port remaps, local credentials, mock payments, and RBAC enforcement. The paths were added to this checkout's `.git/info/exclude`.
- Created the setup tracker, this report, and the continuation handoff under `docs-organized/migration/`.
- Changed only local database seed state: aligned the test tenant slug/domain for storefront resolution and replaced the known seeded super-admin password with a fresh local value.
- No application source code, migrations, tracked Compose files, or production configuration were changed. No branch switch, commit, push, deployment, volume deletion, or image deletion occurred.

## V. Commands to Restart Everything

Run from `H:\Projects\pos` in PowerShell:

```powershell
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml up -d --build
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml ps
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml logs -f backend
```

Open these local URLs:

- API: `http://localhost:8080/api/v1/health`
- Landing: `http://localhost:13001/en`
- Dashboard: `http://localhost:13002/en`
- Restaurant website: `http://restaurant-1.localhost:13003/en`
- Super admin: `http://localhost:13004/login`

To restart or stop without removing the named volumes:

```powershell
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml restart
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml stop
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml start
```

## W. Phase 2A Continuation Readiness

The development stack is ready to continue Phase 2A work. See `NEW_DEVICE_CONTINUATION.md` for remaining security, test, and architecture work. Do not treat setup completion as Phase 2A completion.

**DEVELOPMENT ENVIRONMENT FULLY RESTORED**
