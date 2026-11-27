# F-040 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
No COGS or inventory movement is posted for products with track_inventory=false even when a cost is set; Profitability report still deducts cost.

## 2. Current behavior (verified)
Order 39 (untracked): JE-2026-00064 posts only Dr 1110 330 / Cr 4100 300 / Cr 2120 30 - no COGS; tracked order 41: JE-2026-00066 Dr 5110 120 / Cr 1130 120. Profitability report shows cogs 120 for both products => report COGS exceeds GL COGS by 120 | RUN-003 controlled experiment: untracked product order 39 posts no COGS while Profitability report shows cogs 120; tracked twin order 41 posts COGS 120 (CONFIRMED, auto_posting_engine.go:173). 

## 3. Required behavior
Cost-bearing product sold produces COGS consistent with the Profitability report (or the report excludes untracked cost); policy for untracked items documented

## 4. Reproduction conditions
Create product with cost 40 and track_inventory=false (API default when the flag is omitted), POS order qty 3, mark paid; compare with identical tracked product

## 5. Verified evidence
- evidence/database/RUN-003/PH-06/AUD-INV-COGS-UNTRACKED.txt
- evidence/database/RUN-003/PH-06/OBS-INV-001-controlled-A-untracked.json
- evidence/database/RUN-003/PH-06/OBS-INV-001-controlled.json
- tests/runtime/ph06_inventory_controlled.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/usecase/auto_posting_engine.go:173 (inside PostCOGSEntry :143) skips every order item whose product has track_inventory=false (`if !product.TrackInventory { continue }`), so no COGS/inventory journal or sale movement is created; profitability report computes cogs from products.cost for all items. CreateProduct defaults track_inventory=false when omitted.

## 7. Affected paths
- Backend/database: `backend/internal/usecase/auto_posting_engine.go`:143-176, `backend/internal/handler/http/product_handler.go`:53, `backend/migrations/180_order_stock_reservations.sql`:40-70, `backend/internal/usecase/auto_posting_engine.go`:173
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/dashboard/src/app/[locale]/dashboard/products`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: unchanged unless the design below says otherwise (see 08-api-inventory.csv)
- Accounting logic / posting rules: accounts 1110, 1130, 2120, 4100, 5110
- Reports affected: Products; POS/Orders; Accounting; Profitability report

## 8. Impact
- Financial: 120 in test (cost x quantity of every untracked sale)
- Operational: Products; POS/Orders; Accounting; Profitability report workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: remove the track_inventory guard at backend/internal/usecase/auto_posting_engine.go:173 for COGS (post cost-of-goods for any sold item with a cost, relieving 1130 only when tracked, otherwise a configurable 'cost of non-stocked items' account), default track_inventory to true in the product form or make its effect explicit, and align the profitability report with the GL.

## 10. Alternative designs
Alternative A: enforce in the database (constraint/trigger/lock) so every writer is covered; Alternative B: enforce in the application layer only (smaller change, but direct writers and future code paths stay unprotected). Prefer A for monetary invariants, B for presentation/validation defects.

## 11. Required policy decisions
POL-03.

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
F-011 and shared-file findings F-015, F-026, F-061, F-062

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/product_handler.go, backend/internal/usecase/auto_posting_engine.go, backend/migrations/180_order_stock_reservations.sql).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Sell untracked and tracked cost-bearing products; assert GL COGS equals Profitability COGS (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Order of an untracked product with cost 40 x 3 posts COGS 120 to the configured account; Profitability COGS equals GL COGS for the period.
