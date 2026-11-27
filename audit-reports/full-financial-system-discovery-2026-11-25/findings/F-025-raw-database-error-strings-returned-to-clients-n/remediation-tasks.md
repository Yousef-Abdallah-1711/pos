# F-025 — Implementation phases and tasks (all NOT STARTED)

Finding: Raw database error strings returned to clients (numeric overflow; constraint names) · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-025-01
- **Finding:** F-025
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-025 at code level
- **Existing source paths:** backend/cmd/api/main.go, backend/internal/handler/http/admin_theme_handler_v2.go
- **Proposed paths:** backend/internal/handler/http/errors.go (new); backend/internal/handler/http/errors_scan_test.go (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-025-02
- **Finding:** F-025
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-025
- **Existing source paths:** backend/cmd/api/main.go, backend/internal/handler/http/admin_theme_handler_v2.go
- **Proposed paths:** backend/internal/handler/http/errors.go (new); backend/internal/handler/http/errors_scan_test.go (new)
- **Required changes:** Proposed: introduce one error-mapping layer. (1) A shared helper (proposed backend/internal/handler/http/errors.go) maps known sentinel errors and pq error codes (23505 unique, 23503 foreign key, 23514 check, 22003 numeric overflow, 22001 value too long) to 400/409/422 with a stable code and a generic message; unknown errors log the detail with a request id and return 500 "internal error" plus the request id; (2) replace direct err.Error() / fmt.Sprintf(... %v, err) responses (e.g. backend/cmd/api/main.go:1302, admin_theme_handler_v2.go:80,120,154,361 and the other ~128 sites found by grep) with the helper; (3) add a test that scans handlers for err.Error() in response writers.
- **Dependencies:** REM-F-025-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** No API response body contains pq:, SQLSTATE, constraint names or table names for any induced failure in the fixture battery (duplicate account code, over-long name, numeric overflow). | Duplicate keys return 409 with a stable error code; check violations return 422. | Every 500 response carries a request id that appears in the server log.
- **Risks:** medium
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-025-03
- **Finding:** F-025
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-025
- **Existing source paths:** backend/cmd/api/main.go, backend/internal/handler/http/admin_theme_handler_v2.go
- **Proposed paths:** backend/internal/handler/http/errors.go (new); backend/internal/handler/http/errors_scan_test.go (new)
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Denied roles receive 401/403; no cross-tenant data.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-025-04
- **Finding:** F-025
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-025 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/handler/http/errors.go (new); backend/internal/handler/http/errors_scan_test.go (new)
- **Required changes:** Error-contract tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-025-05
- **Finding:** F-025
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-025
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/handler/http/errors.go (new); backend/internal/handler/http/errors_scan_test.go (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
