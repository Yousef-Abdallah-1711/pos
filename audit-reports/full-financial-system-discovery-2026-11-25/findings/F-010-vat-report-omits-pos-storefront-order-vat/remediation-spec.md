# F-010 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
VAT report omits POS/storefront order VAT.

## 2. Current behavior (verified)
output_vat 0 vs GL 1000 | RUN-003 (AUD-EVD-F010, tenant 4): VAT report API output_vat=700 vs SQL sum(orders.tax_amount, paid)=14,000 (+700 customer-invoice VAT) and GL VAT Payable (acct 87) credit balance 14,700; residual 14,000 = POS order VAT missing from report (STILL REPRODUCIBLE). 

## 3. Required behavior
Output VAT = GL 2120

## 4. Reproduction conditions
Open VAT report after POS sale

## 5. Verified evidence
- reports vat json | evidence/network/RUN-003/PH-20/AUD-EVD-gaps.json
- tests/runtime/ph20_evidence_gaps.mjs
- evidence/screenshots/RUN-003/PH-20/AUD-EVD-F010-vat-report.png

## 6. Root cause and confidence
CONFIRMED — backend/internal/repository/tax_repo.go:260-291 GetVATSummary sums tax_amount of invoices_customer / invoices_supplier only (status <> void, drafts included); order (POS/online) VAT never enters the report

## 7. Affected paths
- Backend/database: `backend/internal/repository/tax_repo.go`:260-291
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/restaurant-website/app/[locale]`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: unchanged unless the design below says otherwise (see 08-api-inventory.csv)
- Accounting logic / posting rules: accounts 2120
- Reports affected: VAT report

## 8. Impact
- Financial: 1000 here
- Operational: VAT report workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: derive the VAT report from the GL (credits to 2120 / debits to 1150 within the period, net of reversals) rather than from invoice tables (backend/internal/repository/tax_repo.go:260-291), and keep an invoice-level drill-down. This includes POS/online order VAT and excludes draft documents once F-008 is fixed.

## 10. Alternative designs
Alternative A: enforce in the database (constraint/trigger/lock) so every writer is covered; Alternative B: enforce in the application layer only (smaller change, but direct writers and future code paths stay unprotected). Prefer A for monetary invariants, B for presentation/validation defects.

## 11. Required policy decisions
POL-02.

## 12. Required data migrations
No schema change expected; confirm during implementation.

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
none

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/repository/tax_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
VAT vs GL reconciliation (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Output VAT in the report equals GL 2120 credit movement for the window (30,880 vs 1,820 in the audit dataset before the fix).
- Input VAT equals GL 1150 debit movement.
- Net VAT equals output minus input.
