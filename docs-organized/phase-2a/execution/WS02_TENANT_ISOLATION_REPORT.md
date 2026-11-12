# WS02 — Tenant & Restaurant Isolation Report

**Status:** PARTIAL — driver repository hardening implemented; full workstream acceptance remains unverified.
**Date:** 2026-11-02
**Database safety:** mutation tests use SQL mocks. PostgreSQL syntax rehearsal used only the isolated `pos-ws01-postgres` restored copy (`pos_saas_restore2`) and ended in rollback. The `pos-dev` database was not used for mutation tests.

## Initial state and scope reviewed

The current-state resource matrix marked products, categories, orders, inventory, themes, users/roles, settings, subscriptions/payments, accounting/HR, drivers and storefront paths as partial or open. It explicitly required a two-tenant/two-restaurants-per-tenant fixture and database side-effect assertions. That integration fixture has not yet been created.

The driver repository exposed 23 methods. Before the change, driver and assignment methods commonly selected or mutated by ID alone. Several current handlers fetched by ID and checked ownership in Go, then issued a second ID-only update/delete. Assignment creation separately checked the order and driver before writing, leaving the ownership boundary and assignment state split across calls.

## Changes implemented

- Recounted and reviewed all 23 `DriverRepository` methods. Each now takes tenant/restaurant scope, applies those predicates in SQL, or (for driver creation) requires the selected restaurant to belong to the supplied tenant.
- Scoped driver get-by-ID and get-by-email lookups, updates, deletes, status/availability/count/location changes, statistics, and delivery metrics.
- Scoped driver location reads through the owning driver and excluded location records linked to an out-of-scope order. Location writes verify driver ownership and, when provided, order ownership.
- Scoped assignment reads and history through both linked driver and order ownership. Assignment status and rating updates require the assignment's driver and order to match the same tenant and restaurant.
- Made order assignment transactional. The repository locks the in-scope order and driver, checks for an existing assignment, inserts the assignment, and updates driver availability/count in the same transaction. Concurrent assignment attempts for the same order serialize on the order row.
- Changed driver admin and admin-order callers to pass tenant/restaurant IDs from middleware claims into the repository. Cross-scope records resolve as not found at these repository boundaries.

## Driver method inventory

| Method group | Scope enforced |
|---|---|
| Create | Restaurant must match tenant in `restaurants` |
| Get by ID/email, tenant/restaurant list, available list | Tenant and restaurant predicates |
| Update, delete, status, availability, active-order count, location update, stats, delivery metrics | Tenant and restaurant predicates in the write/read query |
| Location history read/write | Driver ownership plus tenant/restaurant; optional order must share the same scope |
| Assignment create/read/list/status/rating/complete | Linked driver and order must share tenant/restaurant; assignment create and state update are transactional |

No driver repository method is declared global by design. Unused repository methods were also given scoped signatures instead of being left as latent ID-only APIs.

## Tests and evidence

- Added `driver_repo_test.go` cases for foreign-scope driver lookup, update, delete and order assignment. These assert the scope parameters are part of the SQL and that a missing scoped row produces no successful write/assignment.
- Targeted `go test ./internal/repository ./internal/handler/http`: PASS.
- PostgreSQL 15.19 `PREPARE` syntax checks for scoped driver insert/status update, driver location write/read, assignment ownership lock, assignment rating update and delivery metrics passed against the isolated restored DB; the SQL session rolled back.
- Full backend gate in a temporary Go 1.24 container with network disabled, source read-only, and the existing module cache read-only:
  - `gofmt -d` on touched Go files: clean.
  - `go build ./...`: PASS.
  - `go vet ./...`: PASS.
  - `go test -count=1 ./...`: PASS across all packages; no-test packages were reported separately.

## Files changed

- `backend/internal/repository/driver_repo.go`
- `backend/internal/repository/driver_repo_test.go` (new)
- `backend/internal/handler/http/admin_driver_handler.go`
- `backend/internal/handler/http/admin_order_handler.go`
- `docs-organized/phase-2a/execution/PHASE_2A_EXECUTION_TRACKER.md`

## Remaining work and limitations

- SQL-mock tests verify query construction and call behavior; they do not provide actual two-tenant/two-restaurant database side-effect evidence.
- PostgreSQL PREPARE validates SQL parsing and table/column references, not runtime isolation behavior under fixtures or transaction concurrency.
- The requested deterministic fixture (Tenant A/B, two restaurants each, owners/managers/cashiers and representative products/orders/inventory/drivers/themes/roles) and HTTP cross-tenant/cross-restaurant read/update/delete matrix are still required.
- Other resource families in `08_RESOURCE_ISOLATION_MATRIX.md` were not yet remediated in this workstream: products, categories, orders beyond the driver-assignment path, inventory, themes, users/roles, settings, subscriptions, accounting/HR and public storefront paths.
- Authentication-context and user-controlled-header negative HTTP tests are still needed.
- The existing 73/73 disposable security smoke belongs to WS01 and was run before this WS02 change; it is not represented as a WS02 regression result.
- No live Paymob provider was called. No RLS policy was enabled.

## Closure decision

The driver isolation remediation is implemented and the backend regression gates pass. WS02 remains **PARTIAL / IMPLEMENTED BUT UNVERIFIED** until the isolated integration fixture and cross-resource side-effect matrix pass. Phase 2A cannot be declared complete on this evidence.

## Closure execution update — 2026-11-02

### New isolated evidence

Added `docker-compose.security-closure.yml`. It runs PostgreSQL 15 under the unique Compose project `pos-security-closure-20261102`, exposes only `127.0.0.1:15434`, stores data in container tmpfs, and is marked disposable. The complete current migration directory was applied to this new database. The existing `pos-dev`, `pos-test`, and other Docker projects were not used for mutations.

Added `backend/internal/repository/driver_repo_integration_test.go`. The tests create synthetic Tenant A and B, Restaurants A1/A2/B1/B2 and one driver per restaurant. They verify:

- Tenant A cannot read, change status on, or delete Tenant B's driver.
- Restaurant A1 cannot read, change status on, or delete Restaurant A2's driver.
- In both write cases, a direct PostgreSQL read confirms the target row remains present and `active`.
- A driver cannot be inserted for Tenant A with a Restaurant B ID, and a direct row count confirms no row was created.

The isolated integration tests skip unless `POS_TEST_DATABASE_URL` is explicitly supplied. This prevents them from defaulting to a developer's local test database.

### Driver method recount

Current source has **24** interface methods and 24 concrete methods: 16 driver, location and statistics operations plus 8 assignment operations. This corrects the older 23-method report. Each is scoped by tenant/restaurant in its signature and SQL or validates the tenant↔restaurant relationship. The current relevant callers are the admin driver and admin order handlers. There is no intentional global operation in this repository.

### Verification

- Targeted PostgreSQL tests: 3 test cases pass, repeated three times (9 test executions).
- The tests verify repository behavior and base-row side effects, not HTTP authorization. No product/category/order/inventory/theme/user/role/settings/subscription/accounting/HR/storefront ownership matrix has been completed.
- Final Go `build`, `vet`, and full test regression passed; the closure report records commands and aggregate results.

### Current closure decision

WS02 remains **PARTIAL**. The concrete driver repository boundary now has PostgreSQL side-effect evidence, but there is no full synthetic account/resource fixture, HTTP forgery matrix, bulk mutation matrix, background-job coverage, or cross-resource ownership matrix. The WS02 exit gate is not met.

## Cross-user settings HTTP side-effect test — 2026-11-03

Added `TestUserSettingsHTTPRejectsCrossUserMutationWithoutDatabaseSideEffects`. It authenticates a Tenant A user, sends a valid language update to a different same-tenant user's settings path, and verifies the response is 403, the target language remains unchanged in PostgreSQL, and no target audit row was created. The test passed against the isolated `pos-ws02-closure-20261103` PostgreSQL database.

The focused WS02 packages (`internal/repository`, `internal/handler/http`, `internal/usecase`, and `internal/service`) passed against that database. A subsequent full Go 1.24 regression also passed: `go build ./...`, `go vet ./...`, and `go test -count=1 ./...`. The full test command used only the isolated WS02 `closure_test` database identity. Remaining matrix gaps are still open; this test closes only the user-settings HTTP mutation case.

### HR test fixture packaging correction — 2026-11-03

- Made the new HR child-resource ownership test self-contained so it does not depend on helper code in another uncommitted test file.
- The focused PostgreSQL test passed after the correction. The full Go 1.24.13 build/vet/test gate then passed in the separate Projects checkout after safely integrating the two verified HR commits; the WS02 closure database was the only test database used.
- PR #1 now contains both verified HR commits. WS02 remains partial for the matrix's other listed cases.
