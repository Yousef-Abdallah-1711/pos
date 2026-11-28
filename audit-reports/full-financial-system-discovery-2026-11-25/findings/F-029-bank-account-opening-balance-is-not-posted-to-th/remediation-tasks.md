# F-029 — Implementation phases and tasks (all NOT STARTED)

Finding: Bank account opening balance is not posted to the ledger; bank module and KPI cash differ from GL · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-029-01
- **Finding:** F-029
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-029
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

### REM-F-029-02
- **Finding:** F-029
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-029 at code level
- **Existing source paths:** backend/internal/repository/bank_account_repo.go, backend/internal/handler/http/bank_account_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-029-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-029-03
- **Finding:** F-029
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-029
- **Existing source paths:** backend/internal/repository/bank_account_repo.go, backend/internal/handler/http/bank_account_handler.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: post the bank opening balance. In backend/internal/repository/bank_account_repo.go:101-122 create the account and a balanced opening journal in one transaction (Dr bank account GL / Cr opening balance equity, or Cr the account chosen by the accountant) dated at the opening date; make the KPI cash position read the GL (see F-043) so the bank module and the ledger agree.
- **Dependencies:** REM-F-029-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Creating a bank account with opening balance 50,000 produces one journal Dr bank / Cr equity of 50,000 and GL cash-and-bank equals the bank module balance. | Re-saving the account does not post again. | The balance sheet stays balanced.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-029-04
- **Finding:** F-029
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-029
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (Opening journal or clear disclosure; balances agree) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-029-05
- **Finding:** F-029
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-029
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Existing bank accounts with unposted opening balances are listed with the difference to the GL.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Existing bank accounts with unposted opening balances are listed with the difference to the GL.); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-029-06
- **Finding:** F-029
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-029 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Bank vs GL reconciliation
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-029-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Opening journal or clear disclosure; balances agree; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-029-07
- **Finding:** F-029
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-029
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-029 no longer reproduce (expected: Opening journal or clear disclosure; balances agree); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
