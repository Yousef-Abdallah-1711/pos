# F-030 — Implementation phases and tasks (all NOT STARTED)

Finding: Empty-body POST /users returns 500 but creates an active blank user; other empty-body POSTs return 500 · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-030-01
- **Finding:** F-030
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-030 at code level
- **Existing source paths:** backend/cmd/api/main.go, backend/cmd/api/main.go, backend/cmd/api/main.go
- **Proposed paths:** backend/internal/handler/http/user_admin_handler.go (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-030-02
- **Finding:** F-030
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-030
- **Existing source paths:** backend/cmd/api/main.go, backend/cmd/api/main.go, backend/cmd/api/main.go
- **Proposed paths:** backend/internal/handler/http/user_admin_handler.go (new)
- **Required changes:** Proposed: validate and make user creation atomic. Move the inline handler at backend/cmd/api/main.go:1261-1335 into the handler layer; validate name (required, length), email (format, unique per tenant) and password (minimum length/complexity) before hashing; create the user and assign the EMPLOYEE role in one transaction and fail the whole request if the role assignment fails (today the error is discarded with _ =); return 400 with field errors for empty bodies instead of 500.
- **Dependencies:** REM-F-030-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** POST /users with {} returns 400 with field errors and creates no row. | A failing role assignment rolls back the user. | Concurrent creation with the same email creates exactly one user.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-030-03
- **Finding:** F-030
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-030 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/handler/http/user_admin_handler.go (new)
- **Required changes:** User API validation tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-030-04
- **Finding:** F-030
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-030
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/handler/http/user_admin_handler.go (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
