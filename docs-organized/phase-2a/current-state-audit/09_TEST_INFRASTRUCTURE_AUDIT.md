# Test infrastructure audit

## Backend

Current prior-checkpoint verification: Go 1.24 toolchain; `go build ./...` PASS; `go vet ./...` PASS; `go test -count=1 ./...` PASS with 680 tests, 0 failures (13 packages with tests and 10 with no test files per `TEST_RESULTS.md`). The checks were run in a disposable toolchain/container, not against live dev DB. `-race` was not run. Six experimental package groups behind `-tags experimental` are quarantined and excluded from passing totals; reasons and required environments are in `backend/tests/EXPERIMENTAL.md`. RLS pilot integration passed against isolated DB as historical current-checkpoint evidence.

## Frontend current results

| App | Test command/environment | Result |
|---|---|---|
| Dashboard | Vitest in prior checkpoint | 530/639 pass, 109 fail; 15/41 test files fail. Root-cause-by-test exhaustive grouping remains undone. |
| Landing | `docker exec pos-landing pnpm --filter @pos-saas/landing test:run` | 4 pass; 3 suites fail resolving `@/stores/authStore` in component/integration/login test files. |
| Restaurant website | `docker exec pos-restaurant-website pnpm --filter @pos-saas/restaurant-website test -- --runInBand` | 42 pass, 8 fail; 5 suites pass, 3 fail. Failures include stale theme-store API/persistence expectations, shorthand hex/contrast expectation mismatch, and alias mapping to missing path. |
| Super-admin | Prior checkpoint | 3/3 passed. |

Host-side Vitest/Jest binaries are unavailable; running the actual commands in existing app containers worked. No DB writes or external provider calls were part of those test runs.

## Integration, fixtures and E2E

`docker-compose.pos-test.yml` is a disposable API+PostgreSQL stack with tmpfs DB, mock Paymob, explicit `RBAC_ENFORCEMENT=enforce`, and separate `pos_app`/`pos_migrator`; smoke script asserts disposable labels. The 73-check smoke passed in the previous checkpoint. It is not a full-stack stack: no Redis or four frontends, only a partial deterministic fixture set; missing fixtures include second restaurant per tenant, products/inventory/drivers, broad roles, subscription variants, storefront subdomain and super-admin. No comprehensive browser E2E was run. Existing Playwright config/specs do not amount to a verified CI E2E pipeline.

## CI and coverage

Only two workflows currently exist; they run Python URL/runtime config guardrails and one workflow installs frontend dependencies, but neither executes Go build/vet/tests, frontend typecheck/test/build, migration validation, security smoke, dependency scanning or secret scanning. No measured per-module coverage percentages are available; do not infer coverage from 680 test count.

## Exact gaps

1. Fix/rewrite landing alias and setup so all current tests execute; keep app/test TypeScript scopes separately visible.
2. Resolve storefront test contract, implementation defects and Jest alias config; independently typecheck test sources.
3. Triage dashboard's 109 failures into broken assertions, setup/mocks and product defects; do not blanket-disable tests.
4. Define stable disposable fixtures and run full API+DB+Redis+all apps in isolated CI.
5. Add cross-tenant negative tests, payment concurrency/transaction scenarios, and browser critical journeys.
6. Add CI gates and produce measured coverage by module/layer.
