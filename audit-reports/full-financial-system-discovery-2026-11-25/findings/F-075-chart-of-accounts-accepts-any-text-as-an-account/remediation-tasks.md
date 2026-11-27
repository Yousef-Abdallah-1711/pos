# F-075 — Implementation phases and tasks (all NOT STARTED)

Finding: Chart of accounts accepts any text as an account code or name (HTML and SQL characters, non-numeric codes) with no format or length rule in the UI or API · Severity Low · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-075-01
- **Finding:** F-075
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-075 at code level
- **Existing source paths:** backend/internal/handler/http/account_handler.go, backend/internal/handler/http/account_handler.go
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

### REM-F-075-02
- **Finding:** F-075
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-075
- **Existing source paths:** backend/internal/handler/http/account_handler.go, backend/internal/handler/http/account_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed (root cause PROBABLE: confirm the handler path first): validate account_code against the chart pattern (digits only, 3-10 characters, prefix must match an existing parent when parent_id is set), trim names and limit them to the column size (return 400 not 500), and set maxlength on the account form inputs; escape account fields in CSV/PDF exports; reject codes that would change report classification prefixes unexpectedly.
- **Dependencies:** REM-F-075-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** A code containing letters or HTML characters returns 400 with a field error. | A 300-character name returns 400, not a database error. | Valid numeric codes are accepted and appear in reports.
- **Risks:** medium
- **Status:** NOT STARTED


## P4 Frontend and API contract

### REM-F-075-03
- **Finding:** F-075
- **Phase:** P4 Frontend and API contract
- **Objective:** Align UI and API contract for F-075
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/accounting
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** API contract: POST/PUT /api/v1/accounting/accounts: 400 field errors. Update the affected screens to the corrected contract (EN and AR); keep response fields additive.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** frontend change + contract test
- **Required tests:** headed Chrome check of the affected screens (EN and AR)
- **Acceptance criteria:** UI value equals API value equals stored value for the reproduction dataset.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-075-04
- **Finding:** F-075
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-075 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Invalid code or name returns 400 with a field message; valid code accepted
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-075-05
- **Finding:** F-075
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-075
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
