# F-007 — Implementation phases and tasks (all NOT STARTED)

Finding: Customer invoice tax computed on gross before discount and UI/server totals differ · Severity High · Root-cause confidence CONFIRMED

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P0 Requirements and policy

### REM-F-007-01
- **Finding:** F-007
- **Phase:** P0 Requirements and policy
- **Objective:** Record the policy decision for F-007
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md
- **Proposed paths:** backend/migrations/<next>_fix_invoice_line_tax_discount_order.sql (new)
- **Required changes:** Owner/accountant selects one option; recorded with date and approver. No code until approved.
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** signed decision entry
- **Required tests:** n/a
- **Acceptance criteria:** Decision recorded; remediation spec section 11 updated.
- **Risks:** low
- **Status:** NOT STARTED


## P1 Failing test first

### REM-F-007-02
- **Finding:** F-007
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-007 at code level
- **Existing source paths:** backend/migrations/141_create_invoices_customer.sql, backend/migrations/141_create_invoices_customer.sql, backend/migrations/142_create_invoices_supplier.sql
- **Proposed paths:** backend/migrations/<next>_fix_invoice_line_tax_discount_order.sql (new)
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph09_custinv_draft_headed.mjs).
- **Dependencies:** REM-F-007-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P2 Database and data integrity

### REM-F-007-03
- **Finding:** F-007
- **Phase:** P2 Database and data integrity
- **Objective:** Add migration/constraint for F-007
- **Existing source paths:** backend/migrations (next free number)
- **Proposed paths:** backend/migrations/<next>_fix_invoice_line_tax_discount_order.sql (new)
- **Required changes:** Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.
- **Dependencies:** REM-F-007-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** migration + rollback note
- **Required tests:** migration applies on a copy of the isolated dump; RLS tables unchanged
- **Acceptance criteria:** Migration applies cleanly on the isolated dump and on an empty database.
- **Risks:** high
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-007-04
- **Finding:** F-007
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-007
- **Existing source paths:** backend/migrations/141_create_invoices_customer.sql, backend/migrations/141_create_invoices_customer.sql, backend/migrations/142_create_invoices_supplier.sql
- **Proposed paths:** backend/migrations/<next>_fix_invoice_line_tax_discount_order.sql (new)
- **Required changes:** Proposed: change the line trigger recalculate_invoice_customer_line_total (backend/migrations/141_create_invoices_customer.sql:79-83) so that tax_amount = ROUND((quantity*unit_price - discount_amount) * tax_rate/100, 2) with discount_amount computed first; deliver as a new migration that replaces the function (never edit the applied migration) and recalculates draft invoices only; posted invoices are corrected through reversing journals approved by the accountant. Make the frontend InvoiceForm use the same formula or display the server preview.
- **Dependencies:** REM-F-007-02
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** Line 5,000 with 10% discount and 14% VAT stores tax 630.00 and total 5,130.00; UI preview equals stored values. | Existing posted invoices are unchanged; an exception report lists invoices whose tax differs from the discount-first calculation. | Customer and supplier triggers produce identical results for identical inputs (see F-041).
- **Risks:** medium
- **Status:** NOT STARTED


## P5 Posting and reporting integration

### REM-F-007-05
- **Finding:** F-007
- **Phase:** P5 Posting and reporting integration
- **Objective:** Verify posting rules and reports for F-007
- **Existing source paths:** backend/internal/usecase/auto_posting_engine.go, backend/internal/usecase/financial_report_usecase.go
- **Proposed paths:** backend/migrations/<next>_fix_invoice_line_tax_discount_order.sql (new)
- **Required changes:** Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** reconciliation evidence
- **Required tests:** oracle reconciliation (Decimal / integer minor units)
- **Acceptance criteria:** Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (Tax 630 total 5130 (as UI shows)) and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged
- **Risks:** low
- **Status:** NOT STARTED


## P7 Historical data remediation

### REM-F-007-06
- **Finding:** F-007
- **Phase:** P7 Historical data remediation
- **Objective:** Produce exception report and approved correction plan for F-007
- **Existing source paths:** audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)
- **Proposed paths:** backend/migrations/<next>_fix_invoice_line_tax_discount_order.sql (new)
- **Required changes:** Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.
- **Dependencies:** P3 task
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** read-only exception report
- **Required tests:** report reconciles before/after
- **Acceptance criteria:** Read-only exception report lists the affected historical records (Existing records created under the defect keep wrong values. Produce a read-only exception report (query by th); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals
- **Risks:** low
- **Status:** NOT STARTED


## P8 Verification

### REM-F-007-07
- **Finding:** F-007
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-007 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** backend/migrations/<next>_fix_invoice_line_tax_discount_order.sql (new)
- **Required changes:** Discount+tax invoice tests
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-007-* test in regression-test-spec.md passes on the isolated stack, including the expected result: Tax 630 total 5130 (as UI shows); evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-007-08
- **Finding:** F-007
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-007
- **Existing source paths:** deployment pipeline
- **Proposed paths:** backend/migrations/<next>_fix_invoice_line_tax_discount_order.sql (new)
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-007 no longer reproduce (expected: Tax 630 total 5130 (as UI shows)); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
