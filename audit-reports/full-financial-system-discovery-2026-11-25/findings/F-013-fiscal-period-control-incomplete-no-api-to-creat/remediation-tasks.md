# F-013 — Implementation phases and tasks (all NOT STARTED)

Finding: Fiscal-period control incomplete: no API to create/close periods; void and draft creation allowed in a closed period · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-013-01
- **Finding:** F-013
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-013
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-013-02
- **Finding:** F-013
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-013 at code level
- **Existing source paths:** backend/internal/repository/journal_entry_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-013-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-013-03
- **Finding:** F-013
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-013
- **Existing source paths:** backend/internal/repository/journal_entry_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: add fiscal-period endpoints (list/create/close/reopen) and enforce an open period in ensureFiscalPeriodOpenTx (backend/internal/repository/journal_entry_repo.go:813-842): reject posting when no period exists for the entry date (instead of returning nil), apply the same check to voids and draft creation, and require an accountant-level permission for close/reopen with an audit log.
- **Dependencies:** REM-F-013-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** POST journal dated 2019-01-15 or 2030-06-01 without a period returns 4xx. | After closing 2026-11, posting, voiding and creating drafts dated in November are rejected; reopening permits them. | GET /accounting/fiscal-periods returns the periods (currently 404).
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-013-04
- **Finding:** F-013
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-013
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

### REM-F-013-05
- **Finding:** F-013
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-013
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-013-06
- **Finding:** F-013
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-013 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Closed-period matrix tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-013-07
- **Finding:** F-013
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-013
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
