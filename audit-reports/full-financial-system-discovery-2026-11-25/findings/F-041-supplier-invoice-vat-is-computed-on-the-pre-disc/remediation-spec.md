# F-041 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Supplier invoice VAT is computed on the pre-discount amount, overstating input VAT and accounts payable.

## 2. Current behavior (verified)
API tax_amount 1,400, total 10,400; JE posted Dr 5110 9,000 / Dr 1150 1,400 / Cr 2101 10,400: input VAT +140 and AP +140 overstated | RUN-003 supplier flow (API-DIAGNOSTIC): invoice 10x1000 disc 10% VAT 14% => tax 1,400 total 10,400 vs oracle 1,260/10,260; JE Dr 5110 9,000 Dr 1150 1,400 Cr 2101 10,400 (CONFIRMED trigger 142:75-79). 

## 3. Required behavior
Net 9,000 + VAT 1,260 = 10,260 (discount-first oracle)

## 4. Reproduction conditions
Create supplier invoice line qty 10 x 1,000 with discount 10% and tax 14% (API-DIAGNOSTIC on isolated tenant 4), then approve

## 5. Verified evidence
- evidence/database/RUN-003/PH-06/AUD-SUP-FLOW.json
- tests/runtime/ph06_supplier_flow.mjs

## 6. Root cause and confidence
CONFIRMED — trigger recalculate_invoice_supplier_line_total backend/migrations/142_create_invoices_supplier.sql:75-79 computes tax on quantity*unit_price before subtracting the discount (same formula as customer trigger, F-007)

## 7. Affected paths
- Backend/database: `backend/migrations/142_create_invoices_supplier.sql`:72-81, `backend/migrations/141_create_invoices_customer.sql`:76-85, `backend/migrations/142_create_invoices_supplier.sql`:75-79
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`
- API contracts: unchanged unless the design below says otherwise (see 08-api-inventory.csv)
- Accounting logic / posting rules: accounts 1150, 2101, 5110
- Reports affected: Supplier invoices; Accounting (Input VAT 1150, AP 2101); VAT reports

## 8. Impact
- Financial: 140 per 10,000 gross at 10% discount and 14% VAT (VAT x discount)
- Operational: Supplier invoices; Accounting (Input VAT 1150, AP 2101); VAT reports workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: same trigger correction for backend/migrations/142_create_invoices_supplier.sql:75-79 as F-007, delivered in the same new migration; add a regression test comparing both triggers with the oracle (gross 10,000, discount 10%, tax 14% -> tax 1,260, total 10,260).

## 10. Alternative designs
Alternative A: enforce in the database (constraint/trigger/lock) so every writer is covered; Alternative B: enforce in the application layer only (smaller change, but direct writers and future code paths stay unprotected). Prefer A for monetary invariants, B for presentation/validation defects.

## 11. Required policy decisions
Decide and record in 56-accounting-policy-decision-package.md before implementation: revenue-recognition basis and tax base (discount-first vs gross)

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
F-007 and shared-file findings F-007

## 21. Risks
- Regression in modules sharing the same source files (backend/migrations/141_create_invoices_customer.sql, backend/migrations/142_create_invoices_supplier.sql).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Line with discount and tax on customer and supplier invoices must equal discount-first oracle (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Supplier invoice 10 x 1,000, discount 10%, tax 14% stores tax 1,260.00 total 10,260.00; JE debits 5110 9,000 and 1150 1,260; AP credit 10,260.
- Input VAT in the VAT report equals GL 1150 for the dataset.
