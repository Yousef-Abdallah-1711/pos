# F-036 — Implementation phases and tasks (all NOT STARTED)

Finding: Mock Paymob endpoints reachable without authentication in non-production environments · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-036-01
- **Finding:** F-036
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-036 at code level
- **Existing source paths:** backend/cmd/api/main.go, backend/internal/handler/http/mock_paymob_handler.go
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

### REM-F-036-02
- **Finding:** F-036
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-036
- **Existing source paths:** backend/cmd/api/main.go, backend/internal/handler/http/mock_paymob_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: gate the mock Paymob routes registered at backend/cmd/api/main.go:667-668 behind an explicit flag (PAYMOB_MODE=mock AND ENABLE_MOCK_PAYMENT_ROUTES=true) and require the same authentication or a signed token as the real flow; never register them when PAYMOB_MODE is live; with F-071 also refuse mock mode when ENV is staging.
- **Dependencies:** REM-F-036-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** With PAYMOB_MODE=live the mock routes return 404. | In mock mode they require the flag and a valid token. | Staging refuses to start in mock mode.
- **Risks:** medium
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-036-03
- **Finding:** F-036
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-036
- **Existing source paths:** backend/cmd/api/main.go, backend/internal/handler/http/mock_paymob_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Denied roles receive 401/403; no cross-tenant data.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-036-04
- **Finding:** F-036
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-036 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Env guard tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-036-05
- **Finding:** F-036
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-036
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
