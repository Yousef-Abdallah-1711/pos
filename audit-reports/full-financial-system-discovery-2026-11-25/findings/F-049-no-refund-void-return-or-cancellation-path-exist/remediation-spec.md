# F-049 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
No refund, void, return or cancellation path exists for a paid order (missing financial capability).

## 2. Current behavior (verified)
DELETE -> 400 "cannot cancel paid or refunded order before a supported refund and accounting reversal workflow exists"; payment-status refunded/failed/pending from paid -> 400 invalid transition; voided -> "invalid payment status"; status=cancelled -> 500; paid revenue/COGS stay posted; cash 307,330.01 unchanged | RUN-003: paid order cannot be cancelled/refunded/voided through any endpoint; explicit guard message (REPRODUCED). 

## 3. Required behavior
Controlled refund/void producing reversing journals, stock return and VAT reversal

## 4. Reproduction conditions
Pay an order, then try DELETE /admin/orders/{id}, PUT payment-status refunded/voided/failed/pending, PUT status cancelled, item cancel

## 5. Verified evidence
- evidence/database/RUN-003/PH-19/AUD-REFUND-01-paid-order-cancel.json
- evidence/database/RUN-003/PH-19/AUD-REG-HC10-FIN019-CFR006.json
- tests/runtime/ph19_paid_cancel.mjs
- tests/runtime/ph19_regs2.mjs | evidence/database/RUN-003/PH-07/AUD-ORDER-TYPES.json

## 6. Root cause and confidence
CONFIRMED — usecase/order_usecase.go:341 and repository/order_repo.go:634 reject cancelling paid/refunded orders with "before a supported refund and accounting reversal workflow exists"; domain/order.go:219 CanTransitionOrderPaymentStatus forbids refunded and undoing paid; no refund use case exists

## 7. Affected paths
- Backend/database: `backend/internal/usecase/order_usecase.go`:341, `backend/internal/repository/order_repo.go`:634, `backend/internal/domain/order.go`:215-240, `backend/migrations/175_durable_payment_refund_state.sql`, `backend/internal/domain/order.go`:219
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/dashboard/src/app/[locale]/dashboard/products`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: POST /api/v1/admin/orders/{id}/refunds (new); GET /api/v1/admin/orders/{id}/refunds (new)
- Accounting logic / posting rules: none named
- Reports affected: Orders; Payments; Accounting (revenue/COGS/VAT/inventory reversal)

## 8. Impact
- Financial: Any customer refund cannot be recorded; every refunded sale overstates revenue, VAT and cash
- Operational: Orders; Payments; Accounting (revenue/COGS/VAT/inventory reversal) workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: design and build a refund/void/return workflow for paid orders (the largest missing capability). Remove the blanket rejection at usecase/order_usecase.go:341 and repository/order_repo.go:634 only behind the new workflow: (1) a refund request document (order, items/quantities, reason, tender, approver) with partial and full refunds; (2) approval gated by a new RBAC action; (3) one transaction that records the refund, returns stock for tracked products at their original cost, and posts reversing journals (Dr revenue and VAT payable, Cr cash/bank/card clearing; Dr inventory / Cr COGS for returned goods); (4) idempotency key on the refund; (5) the order moves to refunded or partially_refunded, reports exclude refunded amounts per POL-06, and KDS/stock are updated.

## 10. Alternative designs
A: refund document with approval and reversing journals (recommended). B: cancel-and-recreate (loses audit trail). C: manual correcting journals only (stock and orders stay wrong).

## 11. Required policy decisions
POL-06 (refund / void / return policy) in 56; blocked until approved.

## 12. Required data migrations
Yes: refunds and refund_items tables, status values, refunded totals on orders.

## 13. Historical data implications
None (new capability).

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Mandatory Idempotency-Key on refund creation, stored with the refund (see F-017).

## 16. Authorization requirements
New permission for refund approval, separate from order write.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-005, F-026 and shared-file findings F-028, F-042, F-057

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/domain/order.go, backend/internal/repository/order_repo.go, backend/internal/usecase/order_usecase.go, backend/migrations/175_durable_payment_refund_state.sql).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Full and partial refund per tender type reconcile revenue, VAT, COGS, inventory, cash/clearing (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A full refund of a paid cash order returns the cash, reverses revenue/VAT/COGS exactly (sum of journals for the order = 0) and restores stock.
- A partial refund of one item posts proportional reversals.
- The same refund request twice posts once; concurrent refunds cannot exceed the paid amount.
- The P&L, balance sheet and sales reports reflect the refund per POL-06.
