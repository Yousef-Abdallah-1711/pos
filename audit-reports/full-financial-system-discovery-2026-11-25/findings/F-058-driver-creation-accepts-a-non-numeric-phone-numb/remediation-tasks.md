# F-058 — Implementation phases and tasks (all NOT STARTED)

Finding: Driver creation accepts a non-numeric phone number · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-058-01
- **Finding:** F-058
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-058 at code level
- **Existing source paths:** backend/internal/handler/http/admin_driver_handler.go, backend/internal/handler/http/admin_driver_handler.go, backend/internal/handler/http/admin_driver_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph08_theme_driver.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-058-02
- **Finding:** F-058
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared server-side input validation helpers to cover F-058 (core component delivered once by the P3 task of F-062)
- **Existing source paths:** backend/internal/handler/http/admin_driver_handler.go, backend/internal/handler/http/admin_driver_handler.go, backend/internal/handler/http/admin_driver_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-058-specific cases. Design: Proposed: validate driver phone numbers in backend/internal/handler/http/admin_driver_handler.go (:146-150) with the E.164 validator used by registration, plus email format and licence number length; return 400 with field errors.
- **Dependencies:** F-062 P3 task; REM-F-058-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** Phone "abc" returns 400. | A valid E.164 number is stored normalised.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-058-03
- **Finding:** F-058
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-058 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Malformed phone rejected for drivers, employees, customers
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-058-04
- **Finding:** F-058
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-058
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
