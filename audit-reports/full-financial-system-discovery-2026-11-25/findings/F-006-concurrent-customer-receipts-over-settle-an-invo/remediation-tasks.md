# F-006 — Implementation phases and tasks (all NOT STARTED)

Finding: Concurrent customer receipts over-settle an invoice · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-006-01
- **Finding:** F-006
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-006
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/internal/repository/payment_acct_repo_concurrency_integration_test.go (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-006-02
- **Finding:** F-006
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-006 at code level
- **Existing source paths:** backend/internal/repository/payment_acct_repo.go, backend/internal/repository/payment_acct_repo.go, backend/internal/repository/payment_acct_repo.go
- **Proposed paths:** backend/internal/repository/payment_acct_repo_concurrency_integration_test.go (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-006-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-006-03
- **Finding:** F-006
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-006
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/internal/repository/payment_acct_repo_concurrency_integration_test.go (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-006-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-006-04
- **Finding:** F-006
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-006
- **Existing source paths:** backend/internal/repository/payment_acct_repo.go, backend/internal/repository/payment_acct_repo.go, backend/internal/repository/payment_acct_repo.go
- **Proposed paths:** backend/internal/repository/payment_acct_repo_concurrency_integration_test.go (new)
- **Required changes:** Proposed: in PaymentAccountingRepository.CreatePaymentReceived (backend/internal/repository/payment_acct_repo.go:29) lock the invoice row with SELECT ... FOR UPDATE inside the same transaction (mirror payment_acct_repo.go:499-503 for suppliers), compare amount with the locked balance_due and reject with a 409 'payment amount exceeds invoice balance'; make applyCustomerInvoicePaymentTx (:569-582) conditional (UPDATE ... WHERE balance_due >= amount) and check RowsAffected; remove the silent GREATEST clamp. Add a CHECK (amount_paid <= total_amount) as a safety net after the data pre-check.
- **Dependencies:** REM-F-006-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** 5 parallel receipts of 1,000 on a 2,000 balance: exactly 2 return 201, 3 return 409; invoice amount_paid = 2,000 and balance 0. | GL 1120 for the invoice equals invoice total minus receipts (no negative AR). | Sequential overpayment still returns 400/409 and creates no receipt.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-006-05
- **Finding:** F-006
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-006
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/internal/repository/payment_acct_repo_concurrency_integration_test.go (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Residual for this finding equals 0.00 EGP.
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-006-06
- **Finding:** F-006
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-006
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/internal/repository/payment_acct_repo_concurrency_integration_test.go (new)
- **Required changes:** Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-006-07
- **Finding:** F-006
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-006 (API, real Chrome, reconciliation, concurrency)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/internal/repository/payment_acct_repo_concurrency_integration_test.go (new)
- **Required changes:** Concurrent settlement tests both sides
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-006-08
- **Finding:** F-006
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-006
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/internal/repository/payment_acct_repo_concurrency_integration_test.go (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** Reproduction no longer fails on staging; rollback rehearsed.
- **Risks:** low
- **Status:** NOT STARTED
