# F-056 — Implementation phases and tasks (all NOT STARTED)

Finding: HR contracts and resignations accept end/last-working dates earlier than their start/resignation dates · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-056-01
- **Finding:** F-056
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-056 at code level
- **Existing source paths:** backend/internal/handler/http/contract_handler.go, backend/internal/handler/http/resignation_handler.go, backend/internal/handler/http/contract_handler.go, backend/internal/repository/resignation_repo.go
- **Proposed paths:** backend/migrations/<next>_hr_date_range_checks.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph10_hr_modules.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-056-02
- **Finding:** F-056
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared server-side input validation helpers to cover F-056 (core component delivered once by the P3 task of F-062)
- **Existing source paths:** backend/internal/handler/http/contract_handler.go, backend/internal/handler/http/resignation_handler.go, backend/internal/handler/http/contract_handler.go, backend/internal/repository/resignation_repo.go
- **Proposed paths:** backend/migrations/<next>_hr_date_range_checks.sql (new)
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-056-specific cases. Design: Proposed: validate ordered date ranges. In backend/internal/handler/http/contract_handler.go (create :165-170, update :221-226) and the resignation create/update (backend/internal/repository/resignation_repo.go:117-143 and its handler) require end_date >= start_date and last_working_date >= resignation_date, returning 400; add CHECK constraints so direct writers are covered.
- **Dependencies:** F-062 P3 task; REM-F-056-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** A contract ending before its start returns 400; a valid one returns 201. | A resignation with last working day before the resignation date returns 400. | The CHECK constraints reject the same data written directly.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-056-03
- **Finding:** F-056
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-056 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_hr_date_range_checks.sql (new)
- **Required changes:** Invalid date ranges rejected with 400 for contracts and resignations
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-056-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 400 with a clear message; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-056-04
- **Finding:** F-056
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-056
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_hr_date_range_checks.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-056 no longer reproduce (expected: 400 with a clear message); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
