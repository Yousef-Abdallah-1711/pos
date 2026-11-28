# F-079 — Implementation phases and tasks (all NOT STARTED)

Finding: Leave requests cannot be created through the UI: the form sends ISO date-times and the API accepts only YYYY-MM-DD (400) · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-079-01
- **Finding:** F-079
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-079 at code level
- **Existing source paths:** frontend/apps/dashboard/src/components/hr/LeaveForm.tsx, backend/internal/handler/http/hr_leave_handler.go, frontend/apps/dashboard/src/components/hr/LeaveForm.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph06_r4_form_matrix.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-079-02
- **Finding:** F-079
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-079
- **Existing source paths:** frontend/apps/dashboard/src/components/hr/LeaveForm.tsx, backend/internal/handler/http/hr_leave_handler.go, frontend/apps/dashboard/src/components/hr/LeaveForm.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: in frontend/apps/dashboard/src/components/hr/LeaveForm.tsx:114-115 send data.start_date and data.end_date exactly as the date input returns them (YYYY-MM-DD) instead of new Date(...).toISOString(); additionally make backend/internal/handler/http/hr_leave_handler.go:186 accept both YYYY-MM-DD and RFC3339 (normalise to a date in the tenant time zone, see F-070) so other clients are not affected; add a contract test that posts the exact form payload.
- **Dependencies:** REM-F-079-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Submitting the leave form with valid dates returns 201 and the request appears in the list with the same dates. | The API accepts YYYY-MM-DD and RFC3339 and rejects other formats with a field error. | A form-contract test fails if the payload format changes.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-079-03
- **Finding:** F-079
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-079
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/hr
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** API contract: POST/PUT /api/v1/hr/leaves: accepts both date formats (additive). Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** Submitting the leave form with valid dates returns 201 and the request appears in the list with the same dates.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-079-04
- **Finding:** F-079
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-079 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Submitting the leave form with valid dates returns 201 and lists the request
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-079-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 201 and the request appears in the list; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-079-05
- **Finding:** F-079
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-079
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-079 no longer reproduce (expected: 201 and the request appears in the list); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
