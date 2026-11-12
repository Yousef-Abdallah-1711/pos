# WS05 — Isolated Test Infrastructure Report

**Status:** PARTIAL / NOT VERIFIED
**Date:** 2026-11-02

## Current state

The existing disposable test project contains `pos-test-backend` and `pos-test-postgres`; Docker reports the test PostgreSQL healthy. It does not currently provide the required Redis and four frontend services or the deterministic two-tenant/two-restaurants-per-tenant fixture and browser journeys.

The `pos-dev` Compose project was inspected read-only using the base and local Compose files plus `.env.local`. All seven services were up; PostgreSQL and Redis were healthy. API health, Landing `/en`, Dashboard `/en`, Restaurant Website `/en` (including a restaurant host header), and Super Admin `/login` returned HTTP 200. Direct DNS resolution of `restaurant-1.localhost` is unavailable in this shell, so the browser's custom-host resolution remains unverified. No service was started, stopped, recreated, or removed by this workstream.

## Safety boundary

No fixture seeding or destructive test was run against `pos-dev`. Existing test/dev volumes were not removed. WS05 changes were not made because the complete disposable topology and ownership fixture require further verified wiring and cleanup boundaries.

## Blockers

- Add Redis and all four frontends to the disposable project with distinct project/network/volume names.
- Add deterministic synthetic tenants, restaurants, roles and representative resources, including super-admin and storefront domains.
- Prove no connection path reaches `pos-dev` or production.
- Add repeatable startup, smoke/browser execution, and teardown evidence.

## Closure execution update — 2026-11-02

Added `docker-compose.security-closure.yml` for the minimum WS02/WS03 database-backed tests. Its unique project runs PostgreSQL 15 on loopback port 15434 with `/var/lib/postgresql/data` on tmpfs and disposable labels. The full migration directory was applied to this database. Tests used the cached Go 1.24 container, a read-only module cache and `GOPROXY=off`.

This project completed the driver ownership and payment concurrency/rollback tests. It does not replace the broader WS05 target: no isolated Redis, four frontend apps or browser journey were added. The project was torn down by its exact Compose project name after verification; the Compose file and repeat-run instructions remain.

## Environment recovery continuation — 2026-11-02

Docker Desktop was started and the Linux engine became available. The existing `pos-dev` project reports seven services up; PostgreSQL and Redis are healthy, and `http://localhost:8080/api/v1/health` returned HTTP 200. `http://localhost:13001/en`, `http://localhost:13002/en`, and `http://localhost:13004/login` returned HTTP 200. The storefront's default route redirects (HTTP 307); a host-routed storefront probe is still required. Frontend containers log repeated Watchpack `EIO` scans from the Windows bind-mounted workspace, so these HTTP checks do not establish full UI readiness.

Docker Desktop also brought unrelated Compose projects online during daemon startup. Their services were left untouched after startup. No development database or volume was reset. A separate tmpfs PostgreSQL container was used for the WS10 RLS pilot and stopped after the test. The full isolated seven-service topology and browser journeys remain incomplete.

## Fresh isolated checks — 2026-11-03

A new PostgreSQL 15 instance was launched under the unique Compose project `pos-p2a-rls-20261103`. It used the report's tmpfs-only security-closure service and a restricted `pos_app` role. The RLS/tenant-transaction integration test passed. The project was removed with `docker compose -p pos-p2a-rls-20261103 -f docker-compose.security-closure.yml down`; no `pos-dev` service was operated on, and all seven `pos-dev` containers remained up afterward.

This remains only a database pilot. No seven-service isolated stack, full fixtures, Redis integration, browser E2E, or full-stack teardown exercise was completed. WS05 remains **PARTIAL**.

## Isolated full-stack and browser smoke — 2026-11-03

- Added `docker-compose.phase2a-e2e.yml`, using Compose project `pos-p2a-ws05-e2e-20261103`, PostgreSQL tmpfs on loopback port 15435, Redis, API on 18080, and Landing, Dashboard, Restaurant Website, and Super Admin on ports 13101–13104. It has separate named volumes and labels each disposable service. The `init-roles` shell invocation was corrected after its first start failed; all migrations then completed and API health reported healthy.
- Added `scripts/phase2a-e2e/seed.ps1`. It refuses to seed unless Docker labels match the exact disposable project/database. It is idempotent and creates two independent owner accounts, restaurants, categories, and menu products. Running it twice passed; PostgreSQL showed one product per tenant with distinct tenant and restaurant IDs. The fixed test password is documented below and must only be used with this disposable stack.
- Playwright browser smoke signed in through the Landing login UI as Tenant A, followed the auth callback to `/en/dashboard`, and rendered the authenticated owner dashboard. The initial dashboard request caught incomplete fixture fields (`name_en`/`name_ar`); those were added to the seed, and reloading the dashboard then loaded without the partial-data message.
- Authenticated HTTP/browser-context checks using both account tokens confirmed Tenant A's product list did not include Tenant B's product, foreign product lookup returned 404, and attempted foreign product deletion did not change Tenant B's row. The delete route currently returns 500 for the scoped no-row result; this is a response-status defect, not an observed cross-tenant database side effect.
- The authenticated dashboard rendered Tenant A's seeded owner and products after the seed was corrected for required localized product names. The real POS browser route loaded, adding the seeded burger to the cart succeeded, and `Place Order` created a confirmed POS order for Tenant A in PostgreSQL. The configured stock count remains 20 after order submission, so no claim is made that POS order creation currently decrements stock.
- Inventory route verification exposed a React `Rendered fewer hooks than expected` client crash in the shared dashboard `Sidebar`: tenant/permission early returns preceded a `useMemo` hook. The early returns were moved below the hooks in `frontend/apps/dashboard/src/components/layout/Sidebar.tsx`. A dashboard container restart and route reload still showed the error in the running browser bundle, so this source fix is **not verified** yet.
- The public storefront route loaded at `http://phase2a-restaurant-a.localhost:13103/en` with the seeded restaurant title, and the Super Admin `/login` route loaded at `http://localhost:13104/login`. These are page-load smoke checks only; a storefront order and authenticated administration actions were not exercised. The frontend-development stack logs `ENOMEM` during first compilation and compiles large dashboard routes slowly on the Windows bind mount. Sign-in and POS order creation are proven browser journeys, but the inventory UI adjustment, administration workflow, cross-tenant order display through the UI, and deterministic automated Playwright spec remain open.

### Test-only credentials

Both accounts use `WS05-Disposable-Only-2026!`:

| Tenant | Email | Restaurant slug | Seed product |
| --- | --- | --- | --- |
| A | `phase2a-owner-a@example.test` | `phase2a-restaurant-a` | `WS05 E2E Burger A` |
| B | `phase2a-owner-b@example.test` | `phase2a-restaurant-b` | `WS05 E2E Burger B` |

These fixtures belong only to the tmpfs `pos_phase2a_e2e` database. Do not use the credentials against another environment.

**WS05 remains PARTIAL.** The topology, deterministic seed, auth, product isolation, and POS order journey are verified. The required committed/repeatable browser automation, verified inventory mutation, storefront and administration journeys, and teardown evidence remain open. The isolated Compose project is intentionally left running for continued testing; it has not been applied to `pos-dev`.

## Continued implementation — 2026-11-03

- Moved the Sidebar loading/blocked returns below its permission `useMemo`. After clearing only the disposable Dashboard `.next` cache and restarting the Dashboard service, `/en/dashboard/inventory` rendered without a client-side exception; the seeded product and stock controls were visible. The slow bind-mounted development compile took 28–102 seconds, and old cached bundles continued to show the exception until the cache was cleared.
- Added an authenticated `POST /api/v1/inventory/adjust` endpoint. It scopes and locks the product row by JWT tenant and restaurant, rejects zero/negative resulting stock, and commits the quantity update and `inventory_movements` row in one PostgreSQL transaction. A browser adjustment of +2 returned HTTP 200; direct PostgreSQL verification showed product 1 at quantity 22 and one movement row (baseline 20/0). The preceding POS order did not decrement inventory, as observed in the existing order flow.
- Added `frontend/apps/dashboard/tests/phase2a-e2e/phase2a.spec.ts` and the `test:e2e:phase2a` command for tenant isolation, POS order creation, persisted inventory adjustment, storefront ordering, and Super Admin login. The Super Admin happy path requires `E2E_SUPERADMIN_EMAIL` and `E2E_SUPERADMIN_PASSWORD`; these are intentionally not included in the owner fixture.
- The first repeatable Playwright runner attempt did not start tests because the disposable Dashboard container did not have Playwright Chromium installed. The suite was subsequently run from `mcr.microsoft.com/playwright:v1.64.1-noble` on the isolated stack's shared network namespace.
- The first actual run found that browser-origin API calls from the internal frontend ports were missing from the isolated API CORS allowlist. Added origins `http://localhost:3001` through `http://localhost:3004` to this disposable Compose project and recreated only the API and frontend containers. The tmpfs database and Redis were preserved.
- Final suite result: **4 passed, 1 skipped**. Passed tests cover product tenant isolation and foreign delete response, authenticated POS order creation, inventory page adjustment with PostgreSQL quantity verification, and public storefront order creation. The Super Admin login page rendered, but the authenticated action test skipped because no isolated Super Admin credentials are configured.
- The Dashboard `.next` cache was cleared only inside the disposable named volume. The inventory route then returned HTTP 200 and rendered without the Sidebar hook-order exception.

Repeat command from the repository root after starting/seeding this exact disposable project:

```powershell
docker run --rm --network=container:pos-p2a-ws05-e2e-20261103-api-1 --volumes-from pos-p2a-ws05-e2e-20261103-dashboard-1 -w /app/apps/dashboard -e WS05_DASHBOARD_URL=http://localhost:3002 -e WS05_API_URL=http://localhost:18080/api/v1 -e WS05_STOREFRONT_URL=http://phase2a-restaurant-a.localhost:3003 -e WS05_SUPER_ADMIN_URL=http://localhost:3004 -e WS05_MAP_LOCALHOST=127.0.0.1 mcr.microsoft.com/playwright:v1.64.1-noble sh -lc './node_modules/.bin/playwright test --config playwright.phase2a.config.ts'
```

WS05 remains **PARTIAL**. Repeatable browser verification now passes for tenant isolation, POS, inventory adjustment, and storefront ordering. An authenticated Super Admin workflow and exact-project teardown/reseed proof remain open. POS order creation does not decrement inventory in this flow; no stock decrement is claimed.

## Super Admin and browser regression continuation — 2026-11-03

- Added a disposable Super Admin user and `super_admins` record to the idempotent WS05 seed; credentials remain test-only and are created only after the script validates the labeled disposable database.
- Completed the Super Admin browser journey through login, dashboard, and tenant listing. The initial rerun exposed an ambiguous `Tenants` link locator because the sidebar and dashboard card both link there; the test now selects the exact sidebar link. The focused journey passed.
- The complete Playwright suite passed **5/5**: forged-header product isolation, POS order creation, persisted inventory adjustment, storefront order creation, and Super Admin tenant listing. POS page navigation now waits for DOM content and allows additional time for cold development compilation.
- Added `scripts/phase2a-e2e/stack.ps1` with explicit `Up`, `Seed`, `Test`, and `Down` actions. Teardown checks that every target container belongs to the exact project and Compose file and verifies the disposable labels on PostgreSQL and API before removing that project's containers, network, and named volumes.
- Idempotent seed action passed twice with unchanged Tenant A/B IDs. The exact-project teardown completed; a cold rebuild/start is currently being repeated and must finish with a healthy API and successful reseed before lifecycle closure is claimed.

## Cold-start Super Admin diagnosis and verified fix — 2026-11-03

The first full-project restart run showed the browser returning to `/login` after Super Admin login. API logs and a direct disposable-stack login established that credential verification succeeded; the resulting JWT had `exp: 0` because `SuperAdminClaims` declared an `exp` field that shadowed the embedded `RegisteredClaims.ExpiresAt`. Protected dashboard calls rejected that token. The duplicate JSON field is removed; expiry is now held as a non-serialized derived value, and tests cover valid and expired Super Admin tokens.

The same browser run showed login audit inserts failing because `lib/pq` bound `[]byte` as `bytea` for a JSONB column. The repository now passes encoded details as JSON text and passes nil details as SQL NULL. SQL-mock tests cover both cases. The API hot-reloaded these changes; the subsequent complete Playwright suite passed **5/5**, including Super Admin dashboard and tenant listing. The audit log insert also succeeded in the API log.

The isolated stack lifecycle was verified with guarded teardown of only `pos-p2a-ws05-e2e-20261103`, removal of its tmpfs database and named volumes, a cold Compose rebuild/start, healthy PostgreSQL/Redis/API, and successful reseed. The stack is left running with freshly seeded disposable accounts.

## Final regression after Super Admin and route extraction fixes — 2026-11-03

- Projects checkout backend: `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` all passed. This includes health readiness behavior, route manifest parity, Super Admin JWT expiry, and audit JSONB parameter tests.
- WS02 checkout backend: `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` all passed against only the isolated WS02 test database. The new cross-user settings HTTP negative test also passed under the restricted `pos_app` role.
- WS05 Playwright via `scripts/phase2a-e2e/stack.ps1 -Action Test`: **5 passed, 0 failed** after the cold-start fixes. The Super Admin journey verified the dashboard and both disposable tenant slugs.
- PostgreSQL currently contains three `LOGIN_SUCCESS` and two `VIEW_TENANT_LIST` audit records from disposable browser/API verification. API, PostgreSQL, Redis and all four frontends are healthy after lifecycle restart and reseeding.

## Regression rerun — 2026-11-03

- Re-ran the isolated WS05 Playwright suite: all 5 journeys passed (tenant isolation, POS order creation, inventory adjustment persisted in PostgreSQL, storefront order, and Super Admin login).
- Dashboard Vitest initially attempted to load the Playwright spec as a Vitest suite. Updated the Vitest exclusion to omit the `tests/` browser suite; rerun passed all 41 test files / 639 tests.
- Dashboard `pnpm typecheck` passed in the disposable dashboard container.
- The full-stack stack remains running; no teardown or development/production data changes were made.

## Redirect and inventory concurrency continuation — 2026-11-03

- Added a deterministic unauthenticated Dashboard redirect journey. The first container run reached the configured host-published Landing URL, which is not reachable as a published port from the Playwright container sharing the API network namespace. The test now routes that URL to the Landing service's internal port while still asserting the published redirect was requested and that the real login form rendered.
- Added an inventory regression covering a rejected adjustment with unchanged stock and six concurrent `+1` adjustments. The resulting product quantity equals the starting quantity plus all six writes, demonstrating the row-locked stock update does not lose concurrent changes.
- Read-only PostgreSQL verification on the disposable WS05 database after three full-suite runs showed 18 `WS05 concurrent adjustment` movement rows and zero `WS05 rejected adjustment` rows; the synthetic Burger A stock was 49, consistent with all accepted adjustments.
- Removed the inventory journey's assumption that the seeded product has database ID `1`; it now resolves the synthetic fixture by its stable localized name.
- `scripts/phase2a-e2e/stack.ps1 -Action Test`: **7 passed, 0 failed** on the labeled disposable WS05 project, including all five prior journeys, unauthenticated redirect, and concurrent/rejected inventory adjustments.
- `pnpm --filter @pos-saas/dashboard typecheck`: PASS. `pnpm --filter @pos-saas/dashboard test:run`: **41 files passed, 639 tests passed**.

WS05 remains **PARTIAL**. This batch did not verify cold startup/teardown again, invalid login/logout/expired-session behavior, KDS/table journeys, or complete movement-history assertions. The isolated stack remains running; no development or production data was changed.

## Authentication lifecycle continuation — 2026-11-03

- The Landing login form now declares `method="post"`, preventing a browser-native form submission from placing credentials in the query string before client hydration.
- Added repeatable invalid-login and logout/protected-route Playwright journeys, plus explicit port routing for the Landing and Dashboard services in the shared disposable network namespace.
- Focused invalid-login journey: **PASS**. It verified the login form method, stayed on the login route, showed a credential error, and found no password query parameter.
- Focused logout/protected-route journey: **FAIL**. The Dashboard auth callback page remained at `/auth/callback`; the browser reported failed Next.js CSS/JS chunk requests during the cross-app navigation, so authentication was not established and logout assertions could not run. This is not counted as a pass.
- The attempted callback still carries the access token and serialized user in its query string. Dashboard development-server request logs recorded those query parameters. This is a confirmed credential-exposure gap and requires a safer cross-app auth handoff; the login-form POST change does not close it.
- Earlier full-stack baseline remains **7/7 passed** before adding these auth lifecycle checks. No expanded full-suite rerun is claimed. The labeled WS05 disposable stack remains the only database touched.

## P0 authentication closure verification — 2026-11-04

- All persistence and browser checks used the positively identified `pos-p2a-ws05-e2e-20261103` disposable Compose project. PostgreSQL `pos_phase2a_e2e` is tmpfs-backed; development and production databases were not used.
- Migration 162 adds the Super Admin `ACCESS_DASHBOARD` audit action. The real PostgreSQL audit repository integration passed and verified persistence. Dashboard access now fails if its audit write fails.
- Real Redis integration passed using two independent clients: create/lookup, expiry, revocation visibility, atomic concurrent rotation (one winner), and fail-closed store outage. In-memory concurrent rotation also passed.
- Review found the handoff converted Super Admin access into the tenant owner's unrestricted session without preserving a distinct actor/effective principal or restricted permissions. The UI button is disabled and the API returns 503 without creating a tenant session. The browser test verifies this denial and confirms the Super Admin session remains valid. **The requested handoff is intentionally unavailable; P0 authentication is NO-GO.**
- Fresh Go `build`, `vet`, and uncached `test ./...` passed. Frontend tests passed: Dashboard 639/639, Super Admin 3/3, Landing 35/35, Restaurant Website 64/64. All four app typechecks and production builds passed. Dashboard generated 139 static pages; Super Admin 13 pages; Landing 17 pages; Restaurant Website 7 pages/routes.
- `scripts/phase2a-e2e/stack.ps1 -Action Test`: **10 passed, 0 failed**, covering login/session, logout denial, Dashboard redirect, tenant isolation, POS, persisted inventory adjustment, concurrent/rejected stock adjustments, storefront order, Super Admin login/listing, and the fail-closed tenant-access guard.
- The disposable stack is restored and running after the builds. No teardown was run.

WS05 remains **PARTIAL** because secure Super Admin support access requires actor attribution, restricted authorization, and auditable start/end handling. Other Phase 2A workstream blockers remain. No commit, push or new PR was made because both source worktrees contain extensive mixed dirty changes. Existing PR #1 remains open with failed `runtime-localhost-guardrail` and `url-domain-config-smoke` checks.

## Authentication security continuation — 2026-11-04

- Removed debug statements from API login and login-confirm handlers and their use-case code, including email, request URL, JWT-prefix, and user-identity diagnostics. Focused Go package tests passed: `internal/handler/http`, `internal/usecase`.
- This does not resolve token transfer through callback query strings or token storage in frontend localStorage. Current successful login still sends the JWT and user in the Dashboard callback URL. The previous invalid-login browser result remains valid; logout and protected-route browser verification remain failed/unverified as described above.
- No full Playwright suite was run in this continuation. WS05 remains **PARTIAL** and the P0 auth flow remains blocked on a coordinated HttpOnly session design and callback lifecycle repair.
