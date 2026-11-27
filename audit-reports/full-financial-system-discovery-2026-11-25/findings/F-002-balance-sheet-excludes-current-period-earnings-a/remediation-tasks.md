# F-002 — Implementation phases and tasks (all NOT STARTED)

Finding: Balance Sheet excludes current-period earnings and is UNRECONCILED by exactly the net result · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-002-01
- **Finding:** F-002
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-002
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

### REM-F-002-02
- **Finding:** F-002
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-002 at code level
- **Existing source paths:** backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-002-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-002-03
- **Finding:** F-002
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-002
- **Existing source paths:** backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: in the balance-sheet query (backend/internal/usecase/financial_report_usecase.go:285-306) add a computed equity line 'Current period earnings' = SUM(revenue) - SUM(expense) up to the as-of date (and prior-year retained earnings when a closed fiscal year exists), include it in total_equity and in the equation check; keep account-level rows unchanged.
- **Dependencies:** REM-F-002-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Clean-tenant month: equation_balanced = true and equation_difference = 0.00 (currently -27,710.00). | Trial balance and balance sheet agree for every as-of date in the dataset.
- **Risks:** medium
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-002-04
- **Finding:** F-002
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-002
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

### REM-F-002-05
- **Finding:** F-002
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-002 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** BS residual = 0 under profit and loss
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-002-06
- **Finding:** F-002
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-002
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
