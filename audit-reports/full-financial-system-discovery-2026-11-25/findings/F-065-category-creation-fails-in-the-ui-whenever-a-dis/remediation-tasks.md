# F-065 — Implementation phases and tasks (all NOT STARTED)

Finding: Category creation fails in the UI whenever a display order is typed: the form sends display_order as a JSON string and the API answers 400 "Invalid request body" · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-065-01
- **Finding:** F-065
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-065 at code level
- **Existing source paths:** frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx, frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx, frontend/apps/dashboard/src/components/products/ProductForm.tsx
- **Proposed paths:** frontend/apps/dashboard/src/lib/formNumbers.ts (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph06_ui_category.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-065-02
- **Finding:** F-065
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-065
- **Existing source paths:** frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx, frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx, frontend/apps/dashboard/src/components/products/ProductForm.tsx
- **Proposed paths:** frontend/apps/dashboard/src/lib/formNumbers.ts (new)
- **Required changes:** Proposed: convert numeric form fields before submit. In frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx (:27 state, :129 payload) convert display_order with the same toInt helper the product form uses (components/products/ProductForm.tsx:100), move toInt to a shared lib, and show the server field error message instead of a generic toast; add a component test that submits a category with a display order and asserts the request body type.
- **Dependencies:** REM-F-065-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Saving a category with display order 7 returns 201 and stores 7. | A server validation error is shown next to the field.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-065-03
- **Finding:** F-065
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-065 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** frontend/apps/dashboard/src/lib/formNumbers.ts (new)
- **Required changes:** Create category with display order through the real UI; contract test posting the form payload
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-065-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Category saved with display_order 7; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-065-04
- **Finding:** F-065
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-065
- **Existing source paths:** deployment pipeline
- **Proposed paths:** frontend/apps/dashboard/src/lib/formNumbers.ts (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-065 no longer reproduce (expected: Category saved with display_order 7); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
