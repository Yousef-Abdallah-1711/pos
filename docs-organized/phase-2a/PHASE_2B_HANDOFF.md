# Phase 2B handoff

State: Phase 2A is **INCOMPLETE**. Backend is green and security-contained at the API level; the frontend is partly recovered; refactoring and database-level isolation are not done. Nothing is committed (working tree holds the owner's earlier Docker/Next work plus all Phase 2A changes).

## How to reproduce the verification
```
# Backend (Go is not installed on the host; use the toolchain image)
MSYS_NO_PATHCONV=1 docker run --rm -v <repo>\backend:/src -v pos2a_gomod:/go/pkg/mod -v pos2a_gobuild:/root/.cache/go-build -w /src golang:1.24-alpine sh -c "go build ./... && go vet ./... && go test -count=1 ./..."
# Disposable stack + live security smoke (73 checks)
docker compose -f docker-compose.pos-test.yml up -d --build
bash scripts/test-stack/security-smoke.sh
docker compose -f docker-compose.pos-test.yml down -v        # disposable (tmpfs); safe
# RLS pilot (needs the stack above)
docker run --rm --network pos-test_pos-test-net -e POS_RLS_ADMIN_URL=postgres://postgres:test-only-password@pos-test-postgres:5432/pos_test_disposable?sslmode=disable -e POS_RLS_APP_URL=postgres://pos_app:test-only-app-password@pos-test-postgres:5432/pos_test_disposable?sslmode=disable -v <repo>\backend:/src -w /src golang:1.24-alpine go test -tags integration ./internal/pkg/database -run RLSPilot
# Frontend type checks (dev stack volumes hold node_modules)
docker compose run --rm --no-deps frontend-install sh -c "cd apps/dashboard && npx tsc --noEmit -p tsconfig.json"
```
Experimental suites: `go test -tags experimental ./tests/...` (see `backend/tests/EXPERIMENTAL.md`; expected to fail until their subsystems are rewritten or deleted).

## Recommended order for 2B
1. **Frontend types/builds are done for application code (all four apps: 0 errors, production builds pass without the bypass).** Remaining type work is in tests only (`typecheck:test`: landing 23, super-admin 3, restaurant-website 24, dashboard 158) and the two parallel theme models (`ThemeComponent` numeric id vs `ComponentConfig` string id), currently bridged by one documented cast in `theme-builder/components/LivePreview.tsx`. Pick one model.
2. Make Vitest work everywhere: landing needs `jsdom`, `@vitejs/plugin-react` and testing-library as devDependencies (lockfile change) plus a config with the `@/` alias; dashboard has 15 failing files (109 of 639 tests: `enhanced-error-boundary` 40 — errors escape `render()` although the boundary has `getDerivedStateFromError`, cause not found; `error-logger` 14 severity-detection assertions; theme-builder/editor suites; `error-recovery-suggestions` 5; `rbac` hierarchy tests were updated for the EXPORT level) — triage into stale tests vs product bugs; run restaurant-website Jest. Decide what to do with the dashboard's unreferenced `components/`, `lib/hooks/use-theme.ts`, `i18n/` (excluded for now).
3. **Resource isolation tests (C006)** with two-tenant, two-restaurant fixtures for products, categories, drivers, users, roles, orders, inventory, settings; fix `driver_repo` (18/23 methods unscoped) and the theme repository (still id-only below the handler guard).
4. **Impersonation round trip (C009)** with a super-admin fixture; webhook success-path/out-of-order/concurrency tests with a real transaction (C008).
5. **Apply the restricted DB role to dev/prod compose** with a rehearsed ownership transfer on a copy of real data; then start RLS table families behind `WithTenantTx` (suggested first families: themes, notifications, categories/products, then orders) — each family: move repositories into tenant transactions, add policy, add negative tests, only then enable.
6. **Route extraction from `main.go`** (per-module route files + `Deps` struct), module by module with the existing green suite as the net; delete `legacy`/dead packages once the owner agrees.
7. Extend the test stack: Redis, the four frontends, storefront subdomain tenant, second restaurant per tenant, role fixtures; add the cross-application smoke (landing → login/tenant-select → dashboard; storefront order with token; super-admin).
8. Hardening items still open: server timeouts/graceful shutdown/body limits, typed errors + request IDs + structured logs, CI with build/test gates, backups, trial dates at registration, order atomicity/inventory, accounting invariants, webhook retry queue review, default RBAC mode `enforce` everywhere.

## Decisions waiting for the owner
- Delete or keep: `internal/auth`, `internal/handlers`, `internal/scaling`, `internal/security` RLS manager, sharding tooling, ~31.6k lines tagged `legacy`, dashboard root `components/`/`lib/`.
- Whether `enforce` becomes the default RBAC mode outside production.
- The public-order track page now asks for a phone number when the order was placed on another device — confirm the UX and add translations.
- Migration 161 (public order token) and the changed payloads (see SECURITY_REMEDIATION.md) need to ship with a coordinated storefront deploy.

## Files added or changed this session (beyond session 1)
Backend: `internal/middleware/{internal_token.go, permission.go (rewritten)}`, `cmd/api/startup_guards.go`, `internal/domain/order_public_access.go`, `migrations/161_*`, order repo/usecase/handler (public access), `usecase/payment_usecase.go` (guards), `handler/ws/kds_hub.go` (race), `sharding/router.go` (hash), `repository/kds_repo.go` (metrics), `pkg/database/tenant_tx.go`, `scripts/db/create_runtime_roles.sql`, `docker-entrypoint.sh`, test fixes/new tests (see TEST_RESULTS.md). Frontend: `restaurant-website/lib/{order-access.ts, api-client.ts}`, track page, `dashboard/src/{lib/form-errors.ts, types/qrcode.d.ts, components/products/product-types.ts}` and ~30 type fixes, tsconfig/vitest changes in all apps. Infra: `docker-compose.pos-test.yml`, `scripts/test-stack/*` (smoke, init-roles, and four one-off codemods prefixed `_` kept for reference).
