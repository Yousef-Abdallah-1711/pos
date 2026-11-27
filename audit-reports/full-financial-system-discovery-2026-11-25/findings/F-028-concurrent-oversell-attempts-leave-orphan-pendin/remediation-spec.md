# F-028 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Concurrent oversell attempts leave orphan pending orders with items although the API returned 400.

## 2. Current behavior (verified)
5 accepted; 15 returned 400 yet 15 pending orders with order_items exist (single rejected request leaves none); orphans appear on KDS and in Sales Analytics

## 3. Required behavior
5 orders accepted; 15 rejected leave no order

## 4. Reproduction conditions
Set stock 5; 20 parallel POS orders of 1

## 5. Verified evidence
- evidence/runtime/33

## 6. Root cause and confidence
CONFIRMED — backend/internal/usecase/order_usecase.go:168 CreateOrderWithItems commits the pending order and items first; POS confirmation/stock reservation runs afterwards (:201-205) and its failure returns an error (HTTP 400) while the draft order stays persisted ("complete draft is preserved for audit" comment), leaving orphan pending orders that stats and sales reports count (F-042)

## 7. Affected paths
- Backend/database: `backend/internal/usecase/order_usecase.go`:166-215, `backend/internal/usecase/order_usecase.go`:168
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/kds`, `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`, `frontend/apps/dashboard/src/app/[locale]/dashboard/products`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: POST /api/v1/pos/orders: insufficient stock returns 409 with available quantities (was 400)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 15 x 11000 = 165000 gross phantom pending sales
- Operational: POS; Orders; Inventory; KDS workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: make order creation and stock reservation one transaction. In backend/internal/usecase/order_usecase.go CreateOrderWithItems (:168) move the stock reservation/confirmation (currently :201-205, after the commit) inside the same database transaction that inserts the order and items, locking the product rows (SELECT ... FOR UPDATE ordered by product id to avoid deadlocks); on insufficient stock roll back everything so no pending order remains; return 409 with the available quantity.

## 10. Alternative designs
A: single transaction with row locks (recommended). B: keep two steps and add a compensating cleanup job for orphan pending orders (leaves a window of inconsistent data). C: optimistic stock check (does not prevent oversell under concurrency).

## 11. Required policy decisions
None (technical); whether unpaid orders hold stock is POL-16 (proposed).

## 12. Required data migrations
No.

## 13. Historical data implications
Existing orphan pending orders with items are listed (REL-005/REL-009 evidence) and cancelled through the normal flow.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-042 and shared-file findings F-049, F-057

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/usecase/order_usecase.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Concurrent oversell tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- N concurrent orders for stock S < N create exactly S orders (or S units) and no pending order without a reservation; every rejected request returns 409/400 and leaves no rows.
- Sum of order quantities = stock decrease = ledger movement after the burst.
- No deadlock errors under 100 concurrent submissions of mixed multi-item carts.
