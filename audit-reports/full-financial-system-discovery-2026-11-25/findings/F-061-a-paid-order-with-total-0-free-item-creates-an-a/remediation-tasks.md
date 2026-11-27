# F-061 — Implementation phases and tasks (all NOT STARTED)

Finding: A paid order with total 0 (free item) creates an accounting event that can never post ("at least two journal lines are required") · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-061-01
- **Finding:** F-061
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-061 at code level
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/auto_posting_engine.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph16_masterdata.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-061-02
- **Finding:** F-061
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-061
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/auto_posting_engine.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: complete zero-value events. In backend/internal/usecase/auto_posting_engine.go (the two-line check at :535-537 and PostSalesEntry/PostCOGSEntry) return a distinct "nothing to post" result when the total is zero and mark the outbox event completed with a note, instead of building an empty journal that can never post and retries forever; keep rejecting genuinely unbalanced journals.
- **Dependencies:** REM-F-061-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Paying a zero-total order completes its outbox event with status completed and no journal. | The outbox shows no retryable entry for it. | A non-zero order still posts exactly one journal.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-061-03
- **Finding:** F-061
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-061
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Residual for this finding equals 0.00 EGP.
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-061-04
- **Finding:** F-061
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-061
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Stuck zero-total events (retryable:1 in tenant 4) are listed and closed.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-061-05
- **Finding:** F-061
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-061 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Zero-total paid order leaves the outbox completed with no journal
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-061-06
- **Finding:** F-061
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-061
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
