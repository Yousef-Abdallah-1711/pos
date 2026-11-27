# F-011 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Inventory valuation and turnover do not tie to ledger; opening stock has no journal.

## 2. Current behavior (verified)
196000 (report) vs 16000 (GL) then 156000 vs 8000 | RUN-003 (AUD-EVD-F011): inventory valuation report 192,000 (48 x 4,000) = SQL sum(quantity_in_stock*cost); GL Inventory (acct 1130) balance -36,000; residual 228,000; opening stock has no journal (STILL REPRODUCIBLE). 

## 3. Required behavior
Valuation = GL inventory

## 4. Reproduction conditions
Create product with stock 50 cost 4000 then valuation

## 5. Verified evidence
- reports inventory json | evidence/network/RUN-003/PH-20/AUD-EVD-gaps.json
- tests/runtime/ph20_evidence_gaps.mjs
- evidence/screenshots/RUN-003/PH-20/AUD-EVD-F011-inventory-valuation.png

## 6. Root cause and confidence
CONFIRMED — backend/internal/usecase/financial_report_usecase.go:488-526 GenerateInventoryValuation = SUM(quantity_in_stock x cost) for tracked products; it never reads GL 1130, and purchases/opening stock never post to 1130 (supplier invoice lines post to 5110, OBS-PUR-001)

## 7. Affected paths
- Backend/database: `backend/internal/usecase/financial_report_usecase.go`:488-526
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/dashboard/src/app/[locale]/dashboard/products`
- API contracts: unchanged unless the design below says otherwise (see 08-api-inventory.csv)
- Accounting logic / posting rules: accounts 1130, 5110
- Reports affected: Inventory; Products; Reports

## 8. Impact
- Financial: 180000 / 148000
- Operational: Inventory; Products; Reports workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: choose perpetual inventory (POL-03): opening stock creates Dr 1130 / Cr opening-balance equity journal; supplier invoice lines for inventory items debit 1130 (not 5110) and increase product stock through an inventory movement; COGS relieves 1130 at weighted-average cost; the valuation report reads GL 1130 and reconciles to quantity x cost with the difference displayed.

## 10. Alternative designs
Alternative A: enforce in the database (constraint/trigger/lock) so every writer is covered; Alternative B: enforce in the application layer only (smaller change, but direct writers and future code paths stay unprotected). Prefer A for monetary invariants, B for presentation/validation defects.

## 11. Required policy decisions
POL-03.

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
none and shared-file findings F-002, F-003, F-012

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/usecase/financial_report_usecase.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Stock vs GL reconciliation (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- After opening stock 48 x 4,000 and purchases, GL 1130 equals SUM(quantity_in_stock x average cost) within 0.01.
- Inventory valuation report shows GL balance and a reconciliation difference of 0.00.
