# New Device Continuation Handoff

**Current state:** Development environment is running on the authoritative `main` branch at `1815273924b8ade42f51863262b997115776322b`. This is environment restoration only; Phase 2A remains incomplete and Phase 2B has not started.

## Running the environment

From `H:\Projects\pos` in PowerShell:

```powershell
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml up -d
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml ps
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml logs -f backend
```

The machine-local ignored `.env.local` contains fresh development-only DB, JWT, and super-admin secrets. Do not copy its values into documentation or source control. The machine-local ignored `docker-compose.local.yml` supplies Redis, local port remaps, `RBAC_ENFORCEMENT=enforce`, and `PAYMOB_MODE=mock`.

Host URLs:

- Landing: `http://localhost:13001/en`
- Dashboard: `http://localhost:13002/en`
- Storefront: `http://restaurant-1.localhost:13003/en`
- Super admin: `http://localhost:13004/login`
- API health: `http://localhost:8080/api/v1/health`

The backend container listens on 8080; the apps listen internally on 3001–3004. PostgreSQL and Redis host ports are 15432 and 16379 because unrelated projects already use 5432 and 6379. To stop services without removing their persistent data, use Compose `stop`; do not remove named volumes.

## Current verification baseline

- Go 1.24.13 in Docker: `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` pass (documented baseline: 680 tests).
- App-code typecheck and production-mode frontend builds pass for all four apps.
- Unit tests: landing has 3 suites blocked by the `@/stores/authStore` alias (4 tests pass); dashboard has 109 failures of 639; restaurant website has 8 failures of 50; super admin passes 3/3.
- Test-only TypeScript checks still fail in all four apps.
- The test suites behind the `experimental` build tag were not run. Read `backend/tests/EXPERIMENTAL.md` before considering them.
- The API's invalid credentials and protected route checks return 401. Restaurant tenant status/menu return 200. The public order status route returns 404 without a valid tracking token.

Commands for backend checks:

```powershell
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml run --rm --no-deps --entrypoint go backend build ./...
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml run --rm --no-deps --entrypoint go backend vet ./...
docker compose --env-file .env.local -f docker-compose.yml -f docker-compose.local.yml run --rm --no-deps --entrypoint go backend test -count=1 ./...
```

Frontend app-code typecheck command pattern:

```powershell
docker exec pos-landing pnpm --filter @pos-saas/landing typecheck
docker exec pos-dashboard pnpm --filter @pos-saas/dashboard typecheck
docker exec pos-restaurant-website pnpm --filter @pos-saas/restaurant-website typecheck
docker exec pos-super-admin pnpm --filter @pos-saas/super-admin typecheck
```

## Phase 2A work still open

1. **Database-level isolation:** Dev and production Compose still use PostgreSQL's bootstrap account. The restricted `pos_app`/`pos_migrator` role setup exists for the disposable test stack, but the safe ownership-transfer procedure has not been rehearsed for dev/prod. RLS is not enabled on live tables. Move repositories into tenant-scoped transactions and enable policies table-family by table-family with cross-tenant negative tests. Start with themes, notifications, categories/products, then orders, following `docs-organized/phase-2a/DATABASE_SECURITY.md`.
2. **Resource isolation:** Add two-tenant/two-restaurant tests for products, categories, drivers, users, roles, orders, inventory, and settings. The handoff flags driver repository methods and theme repository scoping as open.
3. **Frontend test recovery:** Fix landing's Vitest alias/config; triage dashboard's 109 failing tests; repair the storefront Jest alias and stale theme-store/color expectations. Resolve test-only TypeScript errors rather than disabling checks. Keep app-code typechecks and production builds green.
4. **Security round trips:** Add a real super-admin impersonation round trip. Add webhook success-path, out-of-order, and concurrent-delivery tests against a real transaction. Fix the current successful-super-admin-login audit insert failure (`invalid input syntax for type json`).
5. **Backend hardening:** Review server timeouts/graceful shutdown/body limits, typed errors/request IDs, CI build/test gates, tenant trial dates, order atomicity/inventory behavior, accounting invariants, and webhook retry behavior.
6. **Architecture:** Extract route registration from `backend/cmd/api/main.go` module by module using the existing Go suite. Decide with the owner before removing legacy/dead packages or duplicate frontend theme structures.
7. **Storefront fixtures and smoke:** Add product/menu fixtures, cross-app login/tenant selection, public order token flow, and role fixtures to the disposable test stack. The local development DB currently has one test tenant and restaurant but no products.

## Security constraints and owner decisions

- Keep development isolated: no production database, customer data, payment credentials, SMTP, or live Paymob transactions. `PAYMOB_MODE` stays `mock` here.
- Keep RBAC enforcement enabled in this local Compose override. Whether `enforce` becomes the default outside production remains an owner decision.
- Do not enable RLS on live tables until each relevant repository uses tenant-scoped transactions and the migration/ownership upgrade is rehearsed on a copy.
- Do not ship migration 161/public-order payload changes without coordinated storefront/API deployment.
- Do not delete or rewrite legacy auth, handlers, scaling, security/RLS, sharding, or dashboard code without the owner's decision listed in `PHASE_2B_HANDOFF.md`.
- Do not switch to `restaurant-ui-redesign`, commit, push, deploy, or delete development volumes as part of continuation setup.
- The known seeded super-admin password was replaced locally. Use the local value in ignored `.env.local`; never restore the repository default.

## Known startup/UI issues

- Landing bare `/login` and `/register` redirect to missing locale routes and return 404. Use `/en/auth/login` and `/en/auth/register`; fix this separately as a narrow routing defect.
- Initial Watchpack scans logged transient ENOMEM warnings for dashboard and restaurant. All apps then served successfully; monitor hot reload during normal development.
- Dashboard frontend builds have a dynamic dependency warning from `component-registry.ts`; storefront/admin builds have React Hook dependency warnings.
- Database-level tenant enforcement is not enabled, and the development API uses the bootstrap PostgreSQL account. Do not treat the running stack as production-safe.

## Source-of-truth handoff documents

- `docs-organized/phase-2a/PHASE_2A_TASKS.md`
- `docs-organized/phase-2a/PHASE_2A_FINAL_REPORT.md`
- `docs-organized/phase-2a/PHASE_2B_HANDOFF.md`
- `docs-organized/phase-2a/DATABASE_SECURITY.md`
- `docs-organized/phase-2a/SECURITY_REMEDIATION.md`
- `docs-organized/phase-2a/TEST_RESULTS.md`
- `backend/tests/EXPERIMENTAL.md`

Some older handoff text describes an earlier branch/commit. Verify statements against the current `main` checkout at `1815273` before acting.
