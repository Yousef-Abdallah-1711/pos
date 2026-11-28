# F-003 — Implementation phases and tasks (all NOT STARTED)

Finding: Cash Flow statement is accrual profit plus code-prefix guesses and does not tie to cash · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-003-01
- **Finding:** F-003
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-003
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/usecase/financial_report_usecase_cashflow_test.go (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-003-02
- **Finding:** F-003
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-003 at code level
- **Existing source paths:** backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/usecase/financial_report_usecase_cashflow_test.go (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-003-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-003-03
- **Finding:** F-003
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-003
- **Existing source paths:** backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/usecase/financial_report_usecase_cashflow_test.go (new)
- **Required changes:** Proposed: rewrite GenerateCashFlow (backend/internal/usecase/financial_report_usecase.go:366-399) as an indirect-method statement derived from the GL: net income + non-cash items +/- change in working-capital accounts (1120, 1130, 1150, 2101, 2120, others by account range) = operating; asset account movements (12xx) with the correct sign = investing; equity and long-term liability movements (3xxx, 22xx) = financing; closing check: opening cash + net change = closing balance of cash accounts (1110, 1111, 1112, 1113 per policy POL-04). Return the reconciliation line so a mismatch is visible instead of silent.
- **Dependencies:** REM-F-003-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** For the clean-tenant month dataset (evidence PH-14) the report net change equals the GL cash change 74,480.00 exactly. | Posting a 20,000 asset purchase (Dr 1250 / Cr 1110) reduces net cash flow by 20,000 and shows investing -20,000. | The report includes closing-balance reconciliation with difference 0.00; a deliberately inconsistent dataset shows a non-zero difference.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-003-04
- **Finding:** F-003
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-003
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/usecase/financial_report_usecase_cashflow_test.go (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (Net cash change equals change in cash accounts) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-003-05
- **Finding:** F-003
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-003
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/internal/usecase/financial_report_usecase_cashflow_test.go (new)
- **Required changes:** Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Existing records created under the defect keep wrong values. Produce a read-only exception report (query by th); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-003-06
- **Finding:** F-003
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-003 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/usecase/financial_report_usecase_cashflow_test.go (new)
- **Required changes:** CF net = delta cash for scenarios
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-003-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Net cash change equals change in cash accounts; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-003-07
- **Finding:** F-003
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-003
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/usecase/financial_report_usecase_cashflow_test.go (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-003 no longer reproduce (expected: Net cash change equals change in cash accounts); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
