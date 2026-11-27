# F-066 — Implementation phases and tasks (all NOT STARTED)

Finding: Dashboard home "Products" tile shows the page length (1) instead of the product total (12) · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-066-01
- **Finding:** F-066
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-066 at code level
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-066-02
- **Finding:** F-066
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-066
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: read the product total from the API contract. In frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx (extractTotalFromCollection :173-178, call :329) read pagination.total (fall back to total, then items.length only when no total exists) and add a unit test with the real /products response shape; consider a lightweight GET /products/count endpoint so the home page does not fetch rows.
- **Dependencies:** REM-F-066-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** With 12 non-deleted products the tile shows 12. | The unit test fails if the extractor reads the wrong key.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-066-03
- **Finding:** F-066
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-066
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/products
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** API contract: optional GET /api/v1/products/count (new). Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** With 12 non-deleted products the tile shows 12. (expected behaviour: Products tile = 12 (non-deleted menu items))
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-066-04
- **Finding:** F-066
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-066 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Dashboard tile equals the API total when there is more than one product
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-066-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Products tile = 12 (non-deleted menu items); evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-066-05
- **Finding:** F-066
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-066
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-066 no longer reproduce (expected: Products tile = 12 (non-deleted menu items)); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
