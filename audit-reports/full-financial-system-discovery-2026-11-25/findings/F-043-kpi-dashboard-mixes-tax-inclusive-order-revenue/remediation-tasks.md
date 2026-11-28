# F-043 — Implementation phases and tasks (all NOT STARTED)

Finding: KPI dashboard mixes tax-inclusive order revenue with GL expenses, and cash position ignores GL cash accounts · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-043-01
- **Finding:** F-043
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-043
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

### REM-F-043-02
- **Finding:** F-043
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-043 at code level
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph12_batch_apis.mjs).
- **Dependencies:** REM-F-043-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-043-03
- **Finding:** F-043
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared shared revenue definitions to cover F-043 (core component delivered once by the P3 task of F-012)
- **Existing source paths:** backend/internal/repository/accounting_advanced_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-043-specific cases. Design: Proposed: derive the KPI dashboard from the GL. In backend/internal/repository/accounting_advanced_repo.go:838-880 replace revenue = SUM(orders.total_amount) with the same GL revenue query used by the P&L (financial_report_usecase.go:208), expenses from the GL expense accounts, net profit = revenue - expenses, and cash position = balances of the cash and bank GL accounts (1110/1120-family) instead of bank_accounts_ledger.current_balance; show a reconciliation line when the bank module differs from the GL (see F-029).
- **Dependencies:** F-012 P3 task; REM-F-043-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** KPI revenue, expenses and net profit equal the P&L for the same period. | KPI cash position equals the GL cash and bank balance. | A difference between the bank module and the GL is displayed, not hidden.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-043-04
- **Finding:** F-043
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-043
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (Revenue/expenses/net profit tie to P&L; cash position ties to GL cash and bank accounts) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-043-05
- **Finding:** F-043
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-043 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** KPI summary equals P&L and TB-derived cash on controlled data
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-043-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Revenue/expenses/net profit tie to P&L; cash position ties to GL cash and bank accounts; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-043-06
- **Finding:** F-043
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-043
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-043 no longer reproduce (expected: Revenue/expenses/net profit tie to P&L; cash position ties to GL cash and bank accounts); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
