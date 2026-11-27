# F-008 — Implementation phases and tasks (all NOT STARTED)

Finding: Customer invoices post to the ledger at creation while status is draft · Severity Medium · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-008-01
- **Finding:** F-008
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-008
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

### REM-F-008-02
- **Finding:** F-008
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-008 at code level
- **Existing source paths:** backend/internal/repository/invoice_customer_repo.go, backend/internal/repository/invoice_customer_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (reproduction steps).
- **Dependencies:** REM-F-008-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-008-03
- **Finding:** F-008
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-008
- **Existing source paths:** backend/internal/repository/invoice_customer_repo.go, backend/internal/repository/invoice_customer_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: post to the ledger only when an invoice leaves draft. In backend/internal/repository/invoice_customer_repo.go createInvoiceTx (:452-560) skip syncJournalEntryTx (:539) when status = draft; add the same call to the issue/send transition and to the approve transition; editing a draft must never touch the ledger. Existing draft invoices that already carry a journal are listed by an exception query and reversed by approved journals.
- **Dependencies:** REM-F-008-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Creating a draft invoice creates no journal entry and leaves AR (1120) and revenue (4100) unchanged. | Issuing the invoice creates exactly one journal; re-issuing does not create a second. | Deleting or voiding a draft leaves the ledger unchanged.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-008-04
- **Finding:** F-008
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-008
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

### REM-F-008-05
- **Finding:** F-008
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-008
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Draft invoices with posted journals (e.g. INV-2026-00002 / JE-2026-00050 in tenant 4) must be listed and reversed explicitly.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Accountant sign-off; no in-place mutation of posted journals.
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-008-06
- **Finding:** F-008
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-008 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Draft invoice accounting test
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** All REG tests pass; evidence stored under evidence/ for the release.
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-008-07
- **Finding:** F-008
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-008
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
