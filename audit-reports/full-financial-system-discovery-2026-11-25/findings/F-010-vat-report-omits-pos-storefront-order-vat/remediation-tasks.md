# F-010 — Implementation phases and tasks (all NOT STARTED)

Finding: VAT report omits POS/storefront order VAT · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-010-01
- **Finding:** F-010
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-010
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

### REM-F-010-02
- **Finding:** F-010
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-010 at code level
- **Existing source paths:** backend/internal/repository/tax_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph20_evidence_gaps.mjs).
- **Dependencies:** REM-F-010-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-010-03
- **Finding:** F-010
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-010
- **Existing source paths:** backend/internal/repository/tax_repo.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed: derive the VAT report from the GL (credits to 2120 / debits to 1150 within the period, net of reversals) rather than from invoice tables (backend/internal/repository/tax_repo.go:260-291), and keep an invoice-level drill-down. This includes POS/online order VAT and excludes draft documents once F-008 is fixed.
- **Dependencies:** REM-F-010-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Output VAT in the report equals GL 2120 credit movement for the window (30,880 vs 1,820 in the audit dataset before the fix). | Input VAT equals GL 1150 debit movement. | Net VAT equals output minus input.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-010-04
- **Finding:** F-010
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-010
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (Output VAT = GL 2120) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-010-05
- **Finding:** F-010
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-010
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Existing records created under the defect keep wrong values. Produce a read-only exception report (query by th); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-010-06
- **Finding:** F-010
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-010 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** VAT vs GL reconciliation
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-010-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Output VAT = GL 2120; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-010-07
- **Finding:** F-010
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-010
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-010 no longer reproduce (expected: Output VAT = GL 2120); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
