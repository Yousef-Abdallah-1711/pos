# F-022 — Implementation phases and tasks (all NOT STARTED)

Finding: Storefront: product without category breaks public menu API (500); storefront UI still renders 0 items with a valid category · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-022-01
- **Finding:** F-022
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-022 at code level
- **Existing source paths:** backend/internal/handler/http/public_menu_handler.go, backend/internal/middleware/cors.go, audit-reports/full-financial-system-discovery-2026-11-25/tests/runtime/docker-compose.audit-persistent.yml
- **Proposed paths:** backend/internal/repository/public_menu_repo_test.go (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-022-02
- **Finding:** F-022
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-022
- **Existing source paths:** backend/internal/handler/http/public_menu_handler.go, backend/internal/middleware/cors.go, audit-reports/full-financial-system-discovery-2026-11-25/tests/runtime/docker-compose.audit-persistent.yml
- **Proposed paths:** backend/internal/repository/public_menu_repo_test.go (new)
- **Required changes:** Proposed: (1) make the product listing used by the public menu robust to a missing category: in the repository behind backend/internal/handler/http/public_menu_handler.go:63 use a LEFT JOIN on categories and group uncategorised products under a synthetic "Other" section, so one bad row can never return 500 for the whole tenant; (2) require a category when creating a product (backend validation plus the dashboard product form, which currently posts category_id 0); (3) for the storefront UI, document the required CORS origin pattern (http://*.localhost:<port> in non-production, the real tenant domains in production) in the deployment guide and add a startup warning when the storefront origin is missing from CORS_ORIGINS.
- **Dependencies:** REM-F-022-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** A tenant with one uncategorised product still returns 200 from the public menu and products endpoints, listing that product under "Other". | POST /products without a valid category_id returns 400 with a field error. | Starting the API with a storefront origin absent from CORS_ORIGINS logs a warning that names the missing origin.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-022-03
- **Finding:** F-022
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-022
- **Existing source paths:** frontend/apps/restaurant-website/app/[locale], frontend/apps/dashboard/src/app/[locale]/dashboard/products
- **Proposed paths:** backend/internal/repository/public_menu_repo_test.go (new)
- **Required changes:** API contract: POST /api/v1/products: category_id required (400 otherwise); GET /api/v1/public/menu: uncategorised products appear under a synthetic section (additive). Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** A tenant with one uncategorised product still returns 200 from the public menu and products endpoints, listing that product under "Other".
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-022-04
- **Finding:** F-022
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-022 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/repository/public_menu_repo_test.go (new)
- **Required changes:** Storefront smoke test after registration + product create
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-022-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Home and menu render; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-022-05
- **Finding:** F-022
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-022
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/repository/public_menu_repo_test.go (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-022 no longer reproduce (expected: Home and menu render); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
