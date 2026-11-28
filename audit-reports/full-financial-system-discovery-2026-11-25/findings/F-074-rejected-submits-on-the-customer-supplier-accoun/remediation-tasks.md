# F-074 — Implementation phases and tasks (all NOT STARTED)

Finding: Rejected submits on the customer, supplier, account and contract forms raise an unhandled promise rejection; three show the raw database error in a toast and the contract form shows nothing · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-074-01
- **Finding:** F-074
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-074 at code level
- **Existing source paths:** frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx, frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx
- **Proposed paths:** frontend/apps/dashboard/src/hooks/useFormSubmit.ts (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph06_r4_form_matrix.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-074-02
- **Finding:** F-074
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-074
- **Existing source paths:** frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx, frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx
- **Proposed paths:** frontend/apps/dashboard/src/hooks/useFormSubmit.ts (new)
- **Required changes:** Proposed (root cause CONFIRMED in RUN-005: parent rethrows, children and the contracts page do not catch): add one shared hook (proposed frontend/apps/dashboard/src/hooks/useFormSubmit.ts) that wraps the submit promise in try/catch, maps a server {error, field_errors} body to react-hook-form setError calls and shows a generic toast otherwise; replace the ad-hoc handlers in the customer, supplier, chart-of-accounts and contract forms (and CategoryFormPage.tsx:129) with it; combine with the backend error mapper (F-025) so messages are safe to display.
- **Dependencies:** REM-F-074-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** A 400 on each listed form shows a readable inline or toast message and no framework error overlay. | The form remains editable and the entered values are kept. | A unit test simulates a 400 and asserts the message is rendered.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-074-03
- **Finding:** F-074
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-074
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers, frontend/apps/dashboard/src/app/[locale]/dashboard/accounting
- **Proposed paths:** frontend/apps/dashboard/src/hooks/useFormSubmit.ts (new)
- **Required changes:** API contract: error bodies gain field_errors (additive, see F-025). Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** A 400 on each listed form shows a readable inline or toast message and no framework error overlay.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-074-04
- **Finding:** F-074
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-074 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** frontend/apps/dashboard/src/hooks/useFormSubmit.ts (new)
- **Required changes:** For each listed form, a rejected submit shows a readable message and no framework error overlay
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-074-* test in regression-test-spec.md passes on the isolated stack, including the expected result: A field-level or toast message built from the server error, and the form stays usable; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-074-05
- **Finding:** F-074
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-074
- **Existing source paths:** deployment pipeline
- **Proposed paths:** frontend/apps/dashboard/src/hooks/useFormSubmit.ts (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-074 no longer reproduce (expected: A field-level or toast message built from the server error, and the form stays usable); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
