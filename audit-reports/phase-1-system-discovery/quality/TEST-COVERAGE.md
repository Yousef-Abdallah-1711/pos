# Test Coverage (static inventory; tests NOT run)

## Go: 62 `_test.go` files, 458 `func Test*`
| Area | Files | Notes |
|---|---|---|
| internal/handler/http | 20 | 8 files carry `//go:build` tags (admin_theme x3, employee, payment, rbac_integration, role) so default `go test` skips them; kds_flow_e2e, accounting/salary/order auto-posting integration |
| internal/repository | 9 | sqlmock-based; tenant isolation (13 tests), SQL injection (3), journal edge cases, kds, table |
| internal/usecase | 7 | auto_posting, financial_report, kds, table, phase6 accounting cycle |
| internal/domain | 3 | |
| internal/handler/ws | 2 | |
| middleware, config, validation, service, service/legacy | 1 each | |
| cmd/api | 2 | |
| backend/tests/unit (auth 2, middleware 1, scaling 5, security 1, sharding 2) | 11 | scaling/sharding code is not wired into main.go |
| backend/tests/integration | 2 | security_rls_integration_test.go references live DB |
| backend/tests/security | 1 | tenant_isolation_test.go |

Live-DB indicators (grep for sql.Open / DATABASE_URL / embedded): backend/tests/integration/security_rls_integration_test.go, backend/tests/unit/auth/rbac_test.go, backend/tests/unit/security/rls_test.go. Others use sqlmock (by grep). Data mutation was not determined beyond this.
Non-Go: backend/tests/load/*.js (3), paymob_api_tests.json + run_phase10_newman.ps1 (Newman against a live API), backward-compatibility suite (.ts).

## Frontend: 59 test files
dashboard 46 (vitest; 4 config variants), landing 4 (vitest), restaurant-website 8 (jest), super-admin 1 (vitest). packages/* and layers/*: none. Dashboard covers KDS, tables, POS, theme-builder editors, error-handling, posApi/reportsApi, stores, hooks. No tests found for accounting/HR pages or accountingApi/hrAdvancedApi.

## Playwright
Root playwright.config.ts: testDir ./e2e/tests (directory does not exist), baseURL http://localhost:3000, webServers `npm run dev --workspace=frontend` / `--workspace=backend` (root package.json declares no workspaces; apps run on 3001-3004; backend is Go). Only spec: frontend/apps/dashboard/tests/visual-regression/visual-regression.spec.ts, outside testDir.

## CI
No workflow runs go test, vitest, jest, lint, tsc or build.
