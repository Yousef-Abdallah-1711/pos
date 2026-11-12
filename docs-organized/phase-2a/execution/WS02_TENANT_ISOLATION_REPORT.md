# WS02 — Tenant & Restaurant Isolation Report

**Status:** PARTIAL — verified isolation now covers additional repository families, one authenticated HTTP mutation path, a transactional bulk path, and the tenant-wide expiry worker; full workstream acceptance remains unverified.
**Date:** 2026-11-03 (latest execution update; history below remains dated)
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

## Final security closure continuation — 2026-11-03

### Implemented during this continuation

- Continued the existing `codex/ws02-final-security-closure` branch in the isolated `H:\ws02` worktree. The unrelated, dirty `H:\Projects\pos` checkout was left untouched. The external standards path named by the owner instructions was absent; the repository `CODING_STANDARDS.md` was used as required by the task handoff.
- Kept every PostgreSQL mutation on the separate `pos-ws02-closure-20261103` project, loopback port `15434`, using its `pos_security_closure_test` database. `pos-dev` was not used. Repository isolation tests ran as restricted `pos_app`; tests requiring a temporary DDL constraint ran as schema owner `pos_migrator`.
- Extended repository-level, real-PostgreSQL negative coverage for categories, products, inventory movement foreign references, orders and order items/history, RBAC assignments, themes, users, and tables. The denial tests assert the targeted rows or related row counts remain unchanged.
- Added an authenticated HTTP test using a signed Tenant A/Restaurant A1 JWT while sending forged Tenant B/Restaurant B1 headers and foreign product IDs. Both foreign-tenant and same-tenant sibling-restaurant deletion attempts are rejected; PostgreSQL confirms each product remains active. Existing middleware tests also cover invalid, conflicting, missing, and anonymous tenant/restaurant contexts.
- Added a mixed-scope bulk table-position test. A valid row followed by a foreign restaurant or foreign tenant table causes the transaction to roll back; all involved positions remain unchanged.
- Added PostgreSQL coverage for the intentionally platform-wide subscription expiry sweep. One expired tenant is marked expired and gets one history row; a separate future-expiry tenant remains active with no expiry event.
- Updated order-item SQL-mock expectations to match the new tenant/restaurant-scoped query and made order test fixture numbers unique while staying inside the schema's 50-character limit.
- Added [WS02_RESOURCE_COVERAGE_MATRIX.csv](WS02_RESOURCE_COVERAGE_MATRIX.csv). It records known route families, handler/service/repository/table mappings, ownership rules, auth, evidence, and gaps. Rows remain explicitly `PARTIAL` or `OPEN` where coverage is incomplete.
- Updated this repository's canonical tracker and final-report files: `PHASE_2A_EXECUTION_TRACKER.md` and `PHASE_2A_EXECUTION_FINAL_REPORT.md`. The alternate filenames `PHASE_2A_REMAINING_ENGINEERING_TRACKER.md` and `PHASE_2A_ENGINEERING_CLOSURE_FINAL.md` named in the request do not exist in this checkout.

### Latest verification

The final post-change Go gate passed in the Go 1.24.13 container against the isolated PostgreSQL database:

```text
go build ./...                 PASS
go vet ./...                   PASS
go test -count=1 ./...         PASS; 13 test-bearing packages, 0 failed
go test -json -count=1 ./...   PASS; 715 passing test/subtest events, 0 failures
```

The database-backed resource suite and authenticated HTTP forged-header test also passed separately with `POS_TEST_DATABASE_URL` set to `postgres://pos_app@host.docker.internal:15434/pos_security_closure_test?sslmode=disable`. This limited run included category, driver, inventory movement, order/item/history, RBAC, product, theme, user, table bulk rollback, and product HTTP context-forgery coverage. The subscription expiry worker test passed using `pos_migrator` because it inserts and verifies tenant history in the migrated database. The full suite, run as `pos_migrator`, includes that worker test.

The final focused restricted-role command was:

```sh
POS_TEST_DATABASE_URL=postgres://pos_app@host.docker.internal:15434/pos_security_closure_test?sslmode=disable \
  go test ./internal/repository -run 'Test(CategoryRepositoryPostgres|DriverRepositoryPostgres|InventoryMovementRepositoryPostgres|OrderRepositoryPostgres|OrderItemRepositoryPostgres|RBACRepositoriesPostgres|ProductRepositoryPostgres|ThemeRepositoryPostgres|UserRepositoryPostgres|TableRepositoryPostgres)' -count=1
POS_TEST_DATABASE_URL=postgres://pos_app@host.docker.internal:15434/pos_security_closure_test?sslmode=disable \
  go test ./internal/handler/http -run '^TestProductHTTPForgedScopeHeadersCannotDeleteForeignProducts$' -count=1
```

### Coverage decision and remaining blockers

**WS02 remains PARTIAL / NOT VERIFIED COMPLETE.** This execution adds real PostgreSQL and HTTP evidence, but acceptance requires coverage across every applicable resource and route, not representative samples. The CSV matrix documents the gaps. In particular, direct HTTP forged-context side effects are verified for product deletion only; HR and accounting table families lack the requested PostgreSQL negative matrix; public storefront/menu configuration, restaurant settings, and subscriptions lack full HTTP object-ID coverage; bulk table positions are covered but no separate bulk assignment/order endpoint has been verified; and other worker types have not been inventoried/tested. These gaps prevent the workstream from being marked complete.

No RLS policy was enabled. No live payment provider was called. WS03 payment behavior and WS06 CI remain outside this task's scope.

## WS02 continuation — 2026-11-03

### Additional confirmed gaps closed

- HR attendance creation, clock-in, leave creation, and salary creation now select the employee only when its tenant and restaurant match the requested scope. Clock-in conflict updates additionally require the existing attendance row's tenant and restaurant to match the incoming row. This prevents a wrong-scope request from inserting through an otherwise valid employee foreign key or updating a legacy conflicting row in another restaurant.
- Accounting account create/update now accept a parent account only when that parent is in the same tenant and restaurant. Journal entry lines already resolve each account using tenant and restaurant predicates; new PostgreSQL evidence confirms a foreign sibling account line aborts the entry and all its lines.
- Added authenticated HTTP coverage for public menu slug scoping, public restaurant theme settings, current-subscription tenant ownership under forged headers, and mixed-scope table-position bulk rollback. The HTTP tests use the isolated PostgreSQL fixture and assert response ownership and/or unchanged database state.
- Reviewed registered worker and bulk routes in `backend/cmd/api/main.go`: the only applicable explicit bulk mutation route is `PUT /api/v1/tables/positions`, now covered at repository and authenticated HTTP levels. The registered subscription expiry worker is covered with PostgreSQL side effects. The other registered worker is the payment webhook retry worker and remains excluded from WS02.
- Updated `WS02_RESOURCE_COVERAGE_MATRIX.csv` with these verified cases and concrete remaining gaps.

### Verification in this continuation

- Focused PostgreSQL repository tests passed using `pos_app` against the isolated `pos-ws02-closure-20261103` database: HR foreign-employee writes and clock-in conflict, account parent ownership, and atomic journal-line ownership.
- Focused PostgreSQL-backed HTTP tests passed using `pos_app`: public menu and theme settings ownership, subscription tenant ownership with forged headers, and mixed-scope bulk table-position rollback.
- Full backend regression passed in Go 1.24.13 with the isolated test database: `go build ./...`, `go vet ./...`, and `go test -count=1 ./...`. All 13 test-bearing packages passed; the remaining packages have no test files.
- No production or developer database was used. No payment behavior or GitHub Actions files were changed.

### Current closure decision

WS02 remains **PARTIAL / NOT ACCEPTANCE-COMPLETE**. The focused gaps above are closed with real PostgreSQL evidence, but the matrix still identifies uncovered HR tables and routes, Accounting invoices/suppliers/customers/bank accounts/payments/tax/recurring routes, settings mutation paths, subscription cancellation/platform administration evidence, and menu/modifier administration. These remaining acceptance cases prevent a complete WS02 sign-off. Branch HEAD still has no commit divergence from `origin/main`; the WS02 implementation remains uncommitted on the requested topic branch, so no PR is ready to open yet.

### WS02 Accounting and HR reference follow-on — 2026-11-03

- Employee creation now verifies the employee's restaurant belongs to the tenant and any manager is in that same tenant and restaurant. Employee updates reject a manager outside the employee's scope and report the scoped no-row result as an error. PostgreSQL tests confirm sibling-manager create/update attempts leave no new employee and do not change the target employee.
- Account creation now verifies restaurant/tenant consistency. Supplier/customer payable and receivable accounts, invoice parties, optional customer invoice orders, and invoice-line accounts are checked against matching tenant and restaurant scope in their SQL write paths. Supplier/customer invoice reads join party rows on matching scope to prevent exposing legacy mismatched relationships.
- Focused PostgreSQL repository tests passed under `pos_app` against the isolated `pos-ws02-closure-20261103` database for employee manager references, attendance/leave/salary employee references, account hierarchy and restaurant ownership, supplier/customer account references, invoice parties and line accounts, and journal-line atomic rollback.
- Route inventory found no registered admin menu/modifier endpoints; currently registered menu routes are public reads. The matrix records this source finding and retains child modifier repository tests as an open item.

WS02 remains **PARTIAL**. HR table families beyond employees/attendance/leaves/salaries and Accounting bank accounts, accounting payments, tax periods, recurring templates, settings mutation routes, subscription mutations (without changing WS03 payment behavior), and complete authenticated HTTP side-effect matrices remain uncovered. No production or development database was used.

### Continuation regression rerun — 2026-11-03

- Corrected the order-usecase SQL mock to expect the tenant/restaurant-scoped write of the public tracking-token hash before the dine-in table occupancy write. The focused `TestOrderUseCase_CreateOrderAutoOccupiesTable` now passes.
- Against the existing isolated `pos-ws02-closure-20261103` PostgreSQL service and `pos_security_closure_test`, Go 1.24.13 `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` all passed using the schema-owner `pos_migrator` identity. The owner identity is required by the existing payment integration test that installs a rollback-only database constraint; this did not change or exercise payment behavior beyond the existing full regression suite.
- Re-ran the HR employee/reference and Accounting account, journal, party, invoice, and invoice-line PostgreSQL test cases with restricted runtime role `pos_app`: PASS. Re-ran public storefront/settings, subscription ownership, and product forged-context HTTP tests plus the subscription expiry worker PostgreSQL test under `pos_app`: PASS.
- An exploratory attempt to run the whole Go suite as `pos_app` correctly failed where an existing payment integration test attempted to alter `tenants`; that role intentionally lacks ownership/DDL privileges. This is a test-role limitation, not a runtime-role regression.

WS02 remains **PARTIAL / NOT ACCEPTANCE-COMPLETE**. These reruns add no new coverage for the remaining matrix rows; uncovered HR/Accounting tables, settings mutation paths, subscription cancellation/platform administration, menu modifier repositories, and full route-level ownership side effects remain open. The source work and evidence remain uncommitted on `codex/ws02-final-security-closure`; no PR was prepared because the acceptance matrix is still incomplete.

## Foreign product delete status correction — 2026-11-03

Changed the product delete handler to map a scoped `sql.ErrNoRows` result to HTTP 404. The authenticated PostgreSQL test `TestProductHTTPForgedScopeHeadersCannotDeleteForeignProducts` passed against the isolated `pos-ws02-closure-20261103` database as restricted `pos_app`. It covers foreign-tenant and sibling-restaurant product IDs with forged tenant/restaurant headers and confirms both target product rows remain unchanged. This closes the previously recorded foreign-delete 500 response defect; WS02 remains partial for the uncovered matrix rows.

## Additional isolated regression — 2026-11-03

Using the same `pos-ws02-closure-20261103` tmpfs PostgreSQL database as its disposable `closure_test` owner role, and a source read-only Go 1.24 container, ran `go test -count=1 ./internal/repository ./internal/handler/http ./internal/usecase ./internal/service` with `POS_TEST_DATABASE_URL=postgres://closure_test@postgres:5432/pos_security_closure_test?sslmode=disable`. All four packages passed. This re-executes the present PostgreSQL repository, HTTP isolation, transaction/usecase, and worker tests, including the forged-header product-delete test. It does not replace the separate restricted-`pos_app` results above or close the explicitly open coverage-matrix cases.

## Additional WS05-driven Super Admin security fixes — 2026-11-03

Fixed duplicate `exp` JSON tags in Super Admin JWT claims that produced tokens with an unusable zero expiration and caused protected dashboard requests to return to login. Added valid/expired token tests. Fixed Super Admin audit JSONB inserts by binding JSON text (and SQL NULL for absent details); the repository tests cover both parameter forms. These changes are outside the WS02 worktree and do not change payment behavior.

## Final focused regression — 2026-11-03

The latest complete backend gate passed in both working checkouts (`go build ./...`, `go vet ./...`, `go test -count=1 ./...`). The WS02 gate used only the isolated `pos-ws02-closure-20261103` PostgreSQL database; the new same-tenant cross-user settings mutation denial was additionally run under restricted `pos_app` and confirmed unchanged settings/audit state. The Projects checkout added the Super Admin token expiration and JSONB audit regressions plus health-route extraction; its full Go gate passed, and the WS05 browser journey passed after those fixes.

### HR child-resource ownership regression — 2026-11-03

- Employment contracts, training records, insurance tracking, and resignations now require referenced employees to match the requested tenant and restaurant on create and update.
- `TestHRChildRepositoriesRejectForeignEmployeesOnCreateAndUpdate` passed against the isolated `pos_security_closure_test` database for sibling-restaurant and foreign-tenant references across all four resource families. Denied creates leave no child rows; denied updates leave the existing in-scope child rows unchanged.
- The complete backend gate passed in a disposable Go 1.24.13 container using that isolated database: `go build ./...`, `go vet ./...`, and `go test -count=1 ./...`.
- WS02 remains **PARTIAL** because other documented HR/Accounting/Settings/Subscription HTTP and database cases remain open.

### HR child-resource test packaging correction — 2026-11-03

The HR child-resource regression now owns its employee fixture helper, so its committed test is independent of other uncommitted test files. The focused PostgreSQL test passed after this correction. The full Go build/vet/test gate then passed from the Projects checkout after these two HR commits were integrated. WS02 remains **PARTIAL** for the remaining acceptance cases in `WS02_RESOURCE_COVERAGE_MATRIX.csv`.

### Subscription administration ownership and concurrency — 2026-11-03

- Added isolated PostgreSQL-backed HTTP tests proving a tenant JWT cannot call platform trial administration even with forged tenant/restaurant headers; both tenants and the target audit rows remain unchanged.
- Added a valid super-admin trial extension test proving only the selected tenant changes, the tenant history identifies the actor/reason, and a platform audit row is recorded. This exposed that audit IP values included the request port and that the audit action name violated the table's allowed-action constraint. IP values are now normalized to PostgreSQL-compatible addresses, and the constrained audit action records `UPDATE_SETTINGS` with `EXTEND_TENANT_TRIAL` preserved as the operation detail.
- Trial extension now locks the tenant row while calculating/updating the end date. A 12-worker PostgreSQL concurrency test confirms all extensions and history rows are retained.
- Focused subscription HTTP and concurrency tests passed against the isolated `pos_security_closure_test` database. Full Go build, vet, and tests passed in Go 1.24.13 against the same disposable database.
- Subscription cancellation and the remaining platform-admin mutation routes remain open; no payment provider behavior was changed.

### Accounting tax-period ownership follow-on — 2026-11-03

- `TaxRepository.CreateTaxPeriod` now inserts only when the supplied restaurant belongs to the supplied tenant. The prior insert relied on independent foreign keys and could persist a mismatched tenant/restaurant pair.
- Added PostgreSQL tests proving mismatched-scope creation leaves no period row and foreign-tenant/sibling-restaurant periods cannot be read, updated, closed, or deleted through another restaurant scope; direct row checks confirmed the target records remained unchanged.
- Against the isolated `pos-ws02-closure-20261103` database using restricted runtime role `pos_app`, `go test -count=1 ./internal/repository -run '^TestTaxPeriodRepository'` passed. The focused Accounting regression `go test -count=1 ./internal/repository -run 'Test(Accounting|Journal|TaxPeriod|AccountRepositoryPostgres)'` also passed.
- WS02 remains **PARTIAL**. Tax-period repository coverage is now recorded in the matrix; bank-account and accounting-payment PostgreSQL coverage, recurring templates, authenticated Accounting HTTP routes, and other listed resource gaps remain open.
- Final Go regression in the WS02 checkout passed `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` with `POS_TEST_DATABASE_URL` pointed at the isolated closure database as table-owner role `pos_migrator`. An initial full-suite attempt as restricted role `pos_app` failed only because an existing payment rollback integration test installs a temporary constraint and lacks table ownership; the focused security tests separately passed under `pos_app`, and the full suite passed under the required owner role. No payment code was changed.

### Accounting recurring-template ownership follow-on — 2026-11-03

- `AccountingAdvancedRepository.CreateRecurringJournalTemplate` now inserts only when the supplied restaurant belongs to the supplied tenant. The prior insert permitted mismatched tenant/restaurant identifiers because the columns had independent foreign keys.
- PostgreSQL tests confirm a mismatched create leaves no template row, and a foreign-tenant or sibling-restaurant template cannot be read, updated, deleted, or executed through another scope. They also verify the target template is unchanged and no recurring journal entry is created.
- Focused tests passed against the isolated `pos-ws02-closure-20261103` database under restricted `pos_app`: `go test -count=1 ./internal/repository -run '^TestRecurringTemplateRepository'` and `go test -count=1 ./internal/repository -run 'Test(Accounting|Journal|TaxPeriod|RecurringTemplate|AccountRepositoryPostgres)'`.
- WS02 remains **PARTIAL**. The matrix now records tax-period and recurring-template repository evidence; bank accounts, accounting payments, full authenticated Accounting HTTP coverage, and other resource families remain open.
