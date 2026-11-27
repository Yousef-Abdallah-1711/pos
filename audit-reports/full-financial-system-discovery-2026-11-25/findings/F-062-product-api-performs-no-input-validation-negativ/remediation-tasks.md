# F-062 — Implementation phases and tasks (all NOT STARTED)

Finding: Product API performs no input validation (negative/zero price, negative cost or stock, empty name, discount above price, more than 2 decimals) and supplier/customer/category APIs accept duplicates and malformed contacts · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-062-01
- **Finding:** F-062
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-062 at code level
- **Existing source paths:** backend/internal/usecase/product_usecase.go, backend/internal/handler/http/product_handler.go, backend/internal/handler/http/supplier_handler.go
- **Proposed paths:** backend/migrations/<next>_product_check_constraints.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph16_masterdata.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-062-02
- **Finding:** F-062
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-062
- **Existing source paths:** backend/internal/usecase/product_usecase.go, backend/internal/handler/http/product_handler.go, backend/internal/handler/http/supplier_handler.go
- **Proposed paths:** backend/migrations/<next>_product_check_constraints.sql (new)
- **Required changes:** Proposed: restore and extend input validation. In backend/internal/usecase/product_usecase.go validateProductRequest (:294-306, all checks commented out) re-enable and extend: name required and <= 100 chars, price > 0 for sellable products and <= a sane maximum with 2 decimals, cost >= 0, stock >= 0, discount <= price, SKU unique per restaurant (return 409, not 500); apply equivalent validation to categories (unique name), customers/suppliers (email/phone format, credit limit >= 0). Add table-driven tests for every rule.
- **Dependencies:** REM-F-062-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Each boundary from file 54 (negative price, zero price, cost above price allowed or not per decision, negative stock, empty or over-long name, duplicate SKU) returns the documented 4xx and creates nothing. | Valid payloads behave as before.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-062-03
- **Finding:** F-062
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-062
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting
- **Proposed paths:** backend/migrations/<next>_product_check_constraints.sql (new)
- **Required changes:** API contract: POST/PUT /api/v1/products, categories, customers, suppliers: 400/409 with field errors. Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** UI value equals API value equals stored value for the reproduction dataset.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-062-04
- **Finding:** F-062
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-062 (API, real Chrome, reconciliation, concurrency)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_product_check_constraints.sql (new)
- **Required changes:** Each invalid input returns 400 with the field name; duplicates return 409
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-062-05
- **Finding:** F-062
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-062
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_product_check_constraints.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
