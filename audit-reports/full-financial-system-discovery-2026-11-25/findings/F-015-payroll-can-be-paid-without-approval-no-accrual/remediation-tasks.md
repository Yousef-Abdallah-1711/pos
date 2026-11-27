# F-015 — Implementation phases and tasks (all NOT STARTED)

Finding: Payroll can be paid without approval; no accrual; no partial payment · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-015-01
- **Finding:** F-015
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-015
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/migrations/<next>_salary_payments_and_accrual.sql (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-015-02
- **Finding:** F-015
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-015 at code level
- **Existing source paths:** backend/internal/repository/salary_repo.go, backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/auto_posting_engine.go
- **Proposed paths:** backend/migrations/<next>_salary_payments_and_accrual.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-015-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-015-03
- **Finding:** F-015
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-015
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/migrations/<next>_salary_payments_and_accrual.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-015-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-015-04
- **Finding:** F-015
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-015
- **Existing source paths:** backend/internal/repository/salary_repo.go, backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/auto_posting_engine.go
- **Proposed paths:** backend/migrations/<next>_salary_payments_and_accrual.sql (new)
- **Required changes:** Proposed: require approval and book payroll in two steps. (1) In backend/internal/repository/salary_repo.go MarkAsPaidTx (:344-380) add the precondition is_approved = true and reject otherwise with 409; (2) introduce an accrual journal at approval (Dr 5210 salaries expense / Cr 2110 accrued salaries) and have payment post Dr 2110 / Cr cash-or-bank (usecase/auto_posting_engine.go postPayrollEntry :276), so expense is recognised in the correct period; (3) allow partial payment with a paid_amount ledger and status partially_paid; (4) UI: disable "Mark paid" until approved.
- **Dependencies:** REM-F-015-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Marking an unapproved salary paid returns 409 and writes nothing. | Approval creates one accrual journal; payment creates one settlement journal; the 2110 balance returns to zero after full payment. | Two concurrent mark-paid calls produce exactly one payment journal.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-015-05
- **Finding:** F-015
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-015
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/migrations/<next>_salary_payments_and_accrual.sql (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Residual for this finding equals 0.00 EGP.
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-015-06
- **Finding:** F-015
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-015
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/migrations/<next>_salary_payments_and_accrual.sql (new)
- **Required changes:** Rows already paid without approval (salaries 1 and 2 in the evidence) are listed in an exception report for management sign-off; no automatic change.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-015-07
- **Finding:** F-015
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-015 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_salary_payments_and_accrual.sql (new)
- **Required changes:** Payroll approval and accrual tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-015-08
- **Finding:** F-015
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-015
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_salary_payments_and_accrual.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
