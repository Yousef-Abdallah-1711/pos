# Phase 2A task list — state at the end of session 2

Statuses: TODO | IN_PROGRESS | BLOCKED | FAILED | VERIFIED_COMPLETE. A task is VERIFIED_COMPLETE only when implementation **and** verification exist.
**Phase 2A: INCOMPLETE** (see PHASE_2A_FINAL_REPORT.md). Session-1 reconciliation: TASK_RECONCILIATION.md.

| ID | Title | Pri | Status | Remaining blockers / evidence |
|---|---|---|---|---|
| C001 | Recover repository state | P0 | VERIFIED_COMPLETE | CONTINUATION_BASELINE.md |
| C002 | Reconcile session-1 tasks | P0 | VERIFIED_COMPLETE | TASK_RECONCILIATION.md |
| C003 | Complete route-level RBAC | P0 | **VERIFIED_COMPLETE (API level)** | 274 of 363 registrations carry a module permission (was 200); the 20 authenticated-only routes are self-service/billing/status and reviewed (ROUTE_AUTHORIZATION_MATRIX.md). Smoke: employee without grants denied on categories/users/drivers/themes/notifications/rbac/contact-info, privilege-escalation attempt denied, owner allowed. Not covered: per-role matrix (read-only vs write vs admin roles) beyond owner/employee-without-roles; KDS WebSocket does not check module 13 |
| C004 | Enforce RBAC in isolated tests | P0 | VERIFIED_COMPLETE for the test stack | `RBAC_ENFORCEMENT=enforce` set explicitly; audit mode never counts as a pass (assertions use 403). **Open:** default outside production is still `audit` (log-only) — decide whether to make `enforce` the default everywhere |
| C005 | Public order object-level authorization | P0 | **VERIFIED_COMPLETE (backend); storefront wired but untested in a browser** | migration 161 (hash column), 256-bit token returned once, hash stored, phone proof for reads only, identical 404, legacy orders unreachable. 15 smoke checks + 3 unit tests. Storefront stores token in localStorage and sends `X-Order-Token`; track page gained a phone field (hard-coded EN/AR strings, not i18n keys). **Not exercised:** real order creation through `/public/orders` (needs product fixtures), POS flow, browser |
| C006 | Theme and resource isolation | P0 | IN_PROGRESS | Themes done (handler layer, 3 unit + 6 smoke). Products/categories/drivers/users/roles/orders/inventory/settings: not tested for cross-tenant access with two-tenant fixtures; Phase 1 found `driver_repo` 18/23 methods unscoped and theme repo filters by id only |
| C007 | Remove unsafe auth defaults | P0 | **VERIFIED_COMPLETE for listed items** | production refuses stub mode + seeded super-admin password; auth endpoints rate limited (5–20/min/IP); translations upload needs a super-admin token; metrics need `METRICS_TOKEN` (404 in production when unset); auth log lines removed. Not done: password-reset flow review, `fmt.Printf` auth debug output elsewhere, per-account lockout, shared (multi-replica) rate limiter |
| C008 | Paymob webhook security | P0 | IN_PROGRESS | 8 new tests: invalid/missing signature, unknown order, failure-after-success, replay on refunded, refund/void callbacks, duplicate failures, currency mismatch + guards added. **Not covered:** the success-activation transaction (needs a real DB tx), out-of-order failure→success, concurrent duplicate deliveries, retry queue |
| C009 | Impersonation integration test | P0 | TODO | Needs a super-admin fixture and the full flow (login → impersonate → act → expiry → revoke → audit). Middleware-level behaviour is tested (6 unit tests) |
| C010 | Restricted DB runtime role | P0 | **VERIFIED_COMPLETE for the isolated stack** | DATABASE_SECURITY.md; 15 smoke checks. **Not applied** to dev/prod compose |
| C011 | Tenancy classification | P1 | **VERIFIED_COMPLETE (generated, heuristic)** | DATABASE_TENANCY_MATRIX.md: 91 tables; 1 needs manual review |
| C012 | Staged RLS | P1 | IN_PROGRESS | Tenant-transaction helper and generic pilot PASS. Migration 164 plus actual `user_auth_audit_logs` policy/repository slice passed disposable PostgreSQL 16 integration as restricted `pos_app` at `fbb247b099b0ce86c947907538b36ab0a2779b29`. No live database migration was applied; remaining runtime repositories still need scoped transaction conversion and per-family policies/tests. Gate stays PARTIAL. |
| C013 | Repair Go test failures | P0 | **VERIFIED_COMPLETE (default suite)** | 16 failing tests + 9 failing packages → 0 failing; 680 tests pass. Root causes: stale KDS test fixtures (new `order_source` column), stale SQL regexes, COUNT-first pagination queries, tautological theme tests with inverted logic, one contradictory slug test, unwired-subsystem suites. **Quarantined, not passing:** 6 packages under `-tags experimental` (tests/EXPERIMENTAL.md) |
| C014 | Complete backend build verification | P0 | VERIFIED_COMPLETE | `go build ./...`, `go vet ./...`, and `cmd/{api,migrate,migrate-shard,seed-themes,hash}` all OK. `scripts/*.go` standalone mains tagged `ignore`. Race detector not run (needs cgo) |
| C015 | Frontend test configuration | P1 | IN_PROGRESS | dashboard: jest-dom added to setup, Playwright specs + dead dirs excluded; landing: still no usable config (dependency/lockfile change needed); restaurant-website: Jest not run |
| C016 | Landing TypeScript | P1 | **VERIFIED_COMPLETE (app code)** | `typecheck` = 0; tests: 23 errors remain |
| C017 | Dashboard TypeScript | P1 | **VERIFIED_COMPLETE (app code)** | 705 combined → 163 app-code (after test/app split) → **0**. Test files still have 158 type errors (`typecheck:test`). Two parallel theme type models remain bridged by one documented cast |
| C018 | Verify all four frontends | P1 | IN_PROGRESS | App typecheck 0/0/0/0 and production build PASS for all four without the bypass. Tests: super-admin 3/3, dashboard 530/639, landing unrunnable, restaurant-website not run; lint not run separately |
| C019 | Complete deterministic fixtures | P1 | IN_PROGRESS | Two tenants/restaurants/owners, employee, orders (direct insert), theme, per run. Missing: second restaurant per tenant, products/inventory/tables/drivers/employees/accounting data, subscriptions, storefront subdomain tenant, super-admin fixture, other roles |
| C020 | Four frontends in the test stack | P1 | TODO | `docker-compose.pos-test.yml` has API + Postgres only (no Redis, no frontends) |
| C021 | Integration smoke tests | P0 | IN_PROGRESS | API/Postgres smoke 73/73. Landing → dashboard, storefront and super-admin flows not run in the isolated stack |
| C022 | Finalize backend architecture | P2 | IN_PROGRESS | BACKEND_ARCHITECTURE.md (current + target) |
| C023 | Refactor backend composition/routes | P2 | TODO | Wrapper helpers and `tenantContext`/`requirePermission`/`statusGuard` centralised in `main.go`; routes not extracted |
| C024 | Refactor backend business modules | P2 | TODO | — |
| C025 | Finalize frontend architecture | P2 | IN_PROGRESS | FRONTEND_ARCHITECTURE.md |
| C026 | Refactor frontend shared packages | P2 | TODO | — |
| C027 | Refactor frontend feature boundaries | P2 | TODO (products component props typed as a side effect) | — |
| C028 | Error-handling foundations | P2 | TODO | Only sensitive auth log lines removed; typed errors/request IDs not started |
| C029 | Full backend regression | P0 | VERIFIED_COMPLETE for what exists | TEST_RESULTS.md |
| C030 | Full frontend regression | P0 | IN_PROGRESS | Typecheck + builds green for all four; tests still red (dashboard 109 failing, landing unrunnable, storefront Jest not run) |
| C031 | Cross-application smoke | P0 | TODO | — |
| C032 | Architecture/compatibility verification | P0 | IN_PROGRESS | Route set compared by script (363 registrations before and after: none removed); API contract diffs for the changed public-order and login-confirm payloads are documented in SECURITY_REMEDIATION.md |

Count: **13 of 32** session-2 tasks are verified complete (C001, C002, C003, C004, C005, C007, C010, C011, C013, C014, C016, C017, C029), several with the caveats written in their rows. 11 are in progress, 0 are blocked and 8 have not been started or only touched. Session-1 tasks 2A-010, 2A-011 and 2A-012 remain verified complete.
