# F-007 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Customer invoice tax computed on gross before discount and UI/server totals differ.

## 2. Current behavior (verified)
Stored tax 700 total 5200 | RUN-003 (AUD-CINV-001): UI Tax 1,008 / Total 8,208 (matches discount-first oracle) vs stored tax 1,120 / total 8,320 (VAT on pre-discount subtotal): +112 tax, +112 AR (STILL REPRODUCIBLE).  RUN-005 real-UI reconfirmation: same invoice: UI preview tax 630.00 / total 5,130.00, stored tax 700.00 / total 5,200.00 (tax on the pre-discount 5,000), journal Dr 1120 5,200 / Cr 4100 4,500 / Cr 2120 700.

## 3. Required behavior
Tax 630 total 5130 (as UI shows)

## 4. Reproduction conditions
Invoice 5000 less 10% discount VAT 14%

## 5. Verified evidence
- network CI-C-001
- reconciliation TX-011 | evidence/screenshots/RUN-003/PH-09/AUD-CINV-001-form-before-submit.png
- evidence/screenshots/RUN-003/PH-09/AUD-CINV-001-after-submit.png
- evidence/network/RUN-003/PH-09/AUD-CINV-001.json
- evidence/database/RUN-003/PH-09/AUD-CINV-001-before-after.txt
- tests/runtime/ph09_custinv_draft_headed.mjs
- evidence/network/RUN-005/PH-09/AUD-R5-INVOICE-FORM-RESULT.json

## 6. Root cause and confidence
CONFIRMED — DB trigger taxes gross | RUN-003 SOURCE: DB trigger recalculate_invoice_customer_line_total (backend/migrations/141_create_invoices_customer.sql:79-83) computes tax_amount = ROUND(qty*unit_price*tax_rate/100,2) on the PRE-discount amount, then line_total = gross - discount + tax; the supplier twin (142_create_invoices_supplier.sql:75-79) is identical (see F-041).

## 7. Affected paths
- Backend/database: `backend/migrations/141_create_invoices_customer.sql`:80-82, `backend/migrations/141_create_invoices_customer.sql`:79-83, `backend/migrations/142_create_invoices_supplier.sql`:75-79
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices`
- API contracts: unchanged unless the design below says otherwise (see 08-api-inventory.csv)
- Accounting logic / posting rules: accounts 1120, 2120, 4100
- Reports affected: Customer invoices

## 8. Impact
- Financial: 70 per invoice
- Operational: Customer invoices workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: change the line trigger recalculate_invoice_customer_line_total (backend/migrations/141_create_invoices_customer.sql:79-83) so that tax_amount = ROUND((quantity*unit_price - discount_amount) * tax_rate/100, 2) with discount_amount computed first; deliver as a new migration that replaces the function (never edit the applied migration) and recalculates draft invoices only; posted invoices are corrected through reversing journals approved by the accountant. Make the frontend InvoiceForm use the same formula or display the server preview.

## 10. Alternative designs
A: fix the trigger (recommended, covers every writer). B: fix only the application layer (leaves direct SQL and imports wrong). C: switch to tax-inclusive pricing (policy change).

## 11. Required policy decisions
POL-02 in 56.

## 12. Required data migrations
Yes — a numbered migration after the current highest in backend/migrations (data-correction + constraint/trigger change). Existing rows must be reconciled first.

## 13. Historical data implications
Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-041 and shared-file findings F-041

## 21. Risks
- Regression in modules sharing the same source files (backend/migrations/141_create_invoices_customer.sql, backend/migrations/142_create_invoices_supplier.sql).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Discount+tax invoice tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Line 5,000 with 10% discount and 14% VAT stores tax 630.00 and total 5,130.00; UI preview equals stored values.
- Existing posted invoices are unchanged; an exception report lists invoices whose tax differs from the discount-first calculation.
- Customer and supplier triggers produce identical results for identical inputs (see F-041).
