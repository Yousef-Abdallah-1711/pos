# F-039 — Implementation phases and tasks (all NOT STARTED)

Finding: Journal entry attachments accepted without schema, type or size validation · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-039-01
- **Finding:** F-039
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-039 at code level
- **Existing source paths:** backend/internal/domain/accounting.go, backend/internal/domain/accounting.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph19_regressions.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-039-02
- **Finding:** F-039
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-039
- **Existing source paths:** backend/internal/domain/accounting.go, backend/internal/domain/accounting.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: validate journal attachments. In backend/internal/domain/accounting.go (Attachments at :95 and :139) replace the raw json.RawMessage passthrough with a typed array of {name, key, size, mime}; enforce max 10 items, max 10 MB per file, an allow-list of MIME types and storage keys issued by the upload service (no arbitrary URLs); reject everything else with 400.
- **Dependencies:** REM-F-039-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** A 1 MB inline blob is rejected with 400. | A valid attachment descriptor is stored and returned unchanged. | javascript: or data: URLs are rejected.
- **Risks:** medium
- **Status:** NOT STARTED


## P6 Security and authorization

### REM-F-039-03
- **Finding:** F-039
- **Phase:** P6 Security and authorization
- **Objective:** Add negative authorization tests for F-039
- **Existing source paths:** backend/internal/domain/accounting.go, backend/internal/domain/accounting.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Cross-role and cross-tenant tests; confirm no new information exposure.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** security tests
- **Required tests:** API security battery (ph17_security.mjs pattern)
- **Acceptance criteria:** Negative authorization tests pass: roles without the permission and other tenants receive 401/403/404 on the affected routes for F-039; no cross-tenant data and no raw database text in responses
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-039-04
- **Finding:** F-039
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-039 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** API tests per variant => 400, no row created
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-039-* test in regression-test-spec.md passes on the isolated stack, including the expected result: 400/422 for non-array, unsafe URL scheme or oversize attachments; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-039-05
- **Finding:** F-039
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-039
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-039 no longer reproduce (expected: 400/422 for non-array, unsafe URL scheme or oversize attachments); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
