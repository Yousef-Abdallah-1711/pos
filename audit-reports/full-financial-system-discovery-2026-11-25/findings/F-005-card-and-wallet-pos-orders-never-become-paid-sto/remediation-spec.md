# F-005 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Card and wallet POS orders never become paid; stock is consumed anyway.

## 2. Current behavior (verified)
Order stays pending; no outbox; no journal; stock decremented (8 unpaid orders = 88000 not in ledger)

## 3. Required behavior
Payment captured or at least recorded; revenue recognised per policy

## 4. Reproduction conditions
POS order with Card or Wallet

## 5. Verified evidence
- runtime/10
- 24

## 6. Root cause and confidence
CONFIRMED — frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx:302-308 calls updateOrderPaymentStatus(paid) only when paymentMethod === cash; card and wallet orders stay payment_status pending, and no other flow settles them (provider settlement belongs to subscription code)

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`:302-308
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/products`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: POST /api/v1/admin/orders/{id}/payments (new); PUT /api/v1/admin/orders/{id}/payment-status: unchanged, but rejected for card/wallet without a payment record (proposed)
- Accounting logic / posting rules: none named
- Reports affected: POS; Orders; Accounting; Inventory

## 8. Impact
- Financial: 88000 gross unrecorded in this run
- Operational: POS; Orders; Accounting; Inventory workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: define the tender lifecycle for non-cash POS orders. (1) In pos/page.tsx:302-308 keep cash settling immediately, but for card/wallet call a new "record payment" action that stores the provider reference and marks the order paid only after a confirmation step (manual confirmation by the cashier for terminal card payments, or a webhook for online wallets); (2) in the backend add POST /api/v1/admin/orders/{id}/payments (proposed) accepting method, amount and reference, validated against the order balance in one transaction (see F-006 locking); (3) release or reserve stock according to policy POL-05: stock is consumed at order creation today even if payment never arrives, so define an expiry that releases stock for unpaid orders (this reservation rule is not yet a numbered decision: proposed as POL-16).

## 10. Alternative designs
A: explicit payment-recording endpoint plus a pending-payment expiry job (recommended). B: force card/wallet orders to be created only after provider confirmation (simplest accounting, needs provider integration). C: treat card orders as paid at creation (inflates revenue when a terminal declines).

## 11. Required policy decisions
POL-07 (non-cash tender settlement) in 56, plus proposed new decision POL-16 (stock reservation for unpaid orders) which must be added to 56 before implementation.

## 12. Required data migrations
Yes: order payment records table (order_id, method, amount, reference, recorded_by, recorded_at) with a unique (order_id, reference).

## 13. Historical data implications
Existing pending card/wallet orders (e.g. 8 orders, 88,000 in the evidence) need a one-off review: confirm or cancel each; do not auto-pay them.

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
- Regression in modules sharing the same source files (frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Card/wallet/split order tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A card order that is confirmed shows payment_status paid, produces one outbox event and one sales journal debiting Card Payment Clearing (1113).
- A card order that is never confirmed does not produce revenue and releases its stock after the configured expiry.
- Recording the same payment twice (same reference) is idempotent.
