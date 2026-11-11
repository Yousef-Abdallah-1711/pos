# Clean build and regression

- `corepack pnpm install --frozen-lockfile` from `frontend/`: PASS; lockfile unchanged.
- Landing production build: PASS.
- Dashboard production build: PASS.
- Restaurant website production build: PASS (existing effect-dependency warning).
- Super Admin production build: PASS (existing effect-dependency warnings).
- Dashboard Vitest: 655 tests passed across 47 files.
- Landing Vitest: 53 tests passed across 8 files. Existing jsdom navigation warnings are emitted by redirect-related tests.
- Super Admin Vitest: 16 tests passed across 3 files.
- Go unit tests: PASS (`go test -mod=readonly ./...`, Go 1.27.4).
- Go vet: PASS (`go vet -mod=readonly ./...`).
- Go build: PASS (`go build -mod=readonly ./...`).
- Disposable PostgreSQL fresh migration: PASS; all 126 migration files applied through `cmd/migrate` to a new test database and produced 92 public tables.
- Migration rollback/reapply: PASS for migration `174_rls_payments.sql`; its documented rollback removed policies/functions and disabled RLS, then the forward migration restored forced RLS and both policies. This is a targeted rollback check, not a full-suite rollback.
- Restricted-role RLS integration: PASS (`TestRLSPilot_TenantIsolationWithRestrictedRole` and `TestAuthAuditAndUserSettingsRLSWithRestrictedRuntimeRole`).
- Payment PostgreSQL integration: PASS (14 tests; callback replay/order, void/refund, recovery, concurrency, restricted runtime role).
- Webhook retry PostgreSQL integration: PASS (4 tests; claim/complete, backoff, stale lease recovery, exhaustion).
- Repository PostgreSQL isolation integration: PASS (27 tests; accounting, category, driver, HR, inventory, order, product, RBAC, theme, and user scope).
- Docker Compose configuration validation: PASS.
- `git diff --check`: PASS.
- Phase 2A Playwright: PASS (16 passed, 0 failed) on the isolated seven-service stack after configuring its host-published ports and tenant-subdomain CORS origin. The final test verifies authenticated KDS WebSocket order delivery following POS order creation.
- Restaurant website API config regression: PASS (5 tests); restaurant website TypeScript check: PASS; production build: PASS after adding the server-only API base URL.
- Registration Playwright: PASS (12 passed, 0 failed) against the disposable API, database, and frontend services.
- Dashboard application and test TypeScript checks: PASS after frozen-lockfile workspace installation; dashboard config and table-list unit tests: PASS (6 tests); dashboard production build: PASS after the mobile header change.

Blocked checks:

- Full migration rollback across all migrations: NOT RUN; the repository has no complete set of executable down migrations. One documented payment-RLS rollback was exercised.
- Production Docker images: NOT RUN. The disposable development stack and browser suite were verified; this is not production-image acceptance.
- Full per-table RLS policy acceptance and route-by-role HTTP matrix: NOT RUN; see the bounded results and uncovered scope in report 05.
- Performance acceptance (Lighthouse/Core Web Vitals/load testing) remains NOT RUN.

Restaurant website has no unit-test script in the workspace package configuration.
