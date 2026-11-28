# F-041 — Implementation phases and tasks (all NOT STARTED)

Finding: Supplier invoice VAT is computed on the pre-discount amount, overstating input VAT and accounts payable · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-041-01
- **Finding:** F-041
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-041
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

### REM-F-041-02
- **Finding:** F-041
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-041 at code level
- **Existing source paths:** backend/migrations/142_create_invoices_supplier.sql, backend/migrations/141_create_invoices_customer.sql, backend/migrations/142_create_invoices_supplier.sql
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph06_supplier_flow.mjs).
- **Dependencies:** REM-F-041-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-041-03
- **Finding:** F-041
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-041
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-041-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-041-04
- **Finding:** F-041
- **Phase:** P3 Backend implementation
- **Objective:** Apply the shared discount-first tax calculation to cover F-041 (core component delivered once by the P3 task of F-007)
- **Existing source paths:** backend/migrations/142_create_invoices_supplier.sql, backend/migrations/141_create_invoices_customer.sql, backend/migrations/142_create_invoices_supplier.sql
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** SHARED IMPLEMENTATION - do not re-implement the component; add only the F-041-specific cases. Design: Proposed: same trigger correction for backend/migrations/142_create_invoices_supplier.sql:75-79 as F-007, delivered in the same new migration; add a regression test comparing both triggers with the oracle (gross 10,000, discount 10%, tax 14% -> tax 1,260, total 10,260).
- **Dependencies:** F-007 P3 task; REM-F-041-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change (finding-specific cases)
- **Required tests:** unit + integration
- **Acceptance criteria:** Supplier invoice 10 x 1,000, discount 10%, tax 14% stores tax 1,260.00 total 10,260.00; JE debits 5110 9,000 and 1150 1,260; AP credit 10,260. | Input VAT in the VAT report equals GL 1150 for the dataset.
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-041-05
- **Finding:** F-041
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-041
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (Net 9,000 + VAT 1,260 = 10,260 (discount-first oracle)) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-041-06
- **Finding:** F-041
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-041
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

### REM-F-041-07
- **Finding:** F-041
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-041 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Line with discount and tax on customer and supplier invoices must equal discount-first oracle
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-041-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Net 9,000 + VAT 1,260 = 10,260 (discount-first oracle); evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-041-08
- **Finding:** F-041
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-041
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-041 no longer reproduce (expected: Net 9,000 + VAT 1,260 = 10,260 (discount-first oracle)); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
