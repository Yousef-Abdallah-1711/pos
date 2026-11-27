# F-014 — Implementation phases and tasks (all NOT STARTED)

Finding: Employees cannot be created from the UI (HTTP 500; UI shows generic toast, raw DB error exposed in API response body) · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-014-01
- **Finding:** F-014
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-014 at code level
- **Existing source paths:** frontend/apps/dashboard/src/components/hr/EmployeeForm.tsx, backend/internal/handler/http/hr_employee_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph10_emp_create_headed.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-014-02
- **Finding:** F-014
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-014
- **Existing source paths:** frontend/apps/dashboard/src/components/hr/EmployeeForm.tsx, backend/internal/handler/http/hr_employee_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: align the EmployeeForm payload with the backend contract (frontend/apps/dashboard/src/components/hr/EmployeeForm.tsx:76-77 vs backend/internal/handler/http/hr_employee_handler.go:139-150) and make the handler return 400 with the offending field instead of a 500 with raw database text; add a contract test that posts the form's actual payload.
- **Dependencies:** REM-F-014-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Creating an employee through the real UI in headed Chrome returns 201 and the employee appears in the list. | Invalid payloads return 400 naming the field; no raw DB error text in any response.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-014-03
- **Finding:** F-014
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-014 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Employee create UI/API contract test
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-014-04
- **Finding:** F-014
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-014
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
