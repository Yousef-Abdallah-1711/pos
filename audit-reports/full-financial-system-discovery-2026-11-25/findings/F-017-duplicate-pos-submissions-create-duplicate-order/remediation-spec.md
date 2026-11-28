# F-017 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Duplicate POS submissions create duplicate orders; Idempotency-Key is ignored on order creation.

## 2. Current behavior (verified)
5 and 3 orders created; stock decremented each | RUN-003 REL-001: same order POST twice with identical Idempotency-Key creates 2 orders (REPRODUCED). 

## 3. Required behavior
One order

## 4. Reproduction conditions
5 identical parallel POST and 3 with same key

## 5. Verified evidence
- runtime/11 | evidence/database/RUN-003/PH-18/AUD-REL-BATTERY.json | tests/runtime/ph18_reliability.mjs

## 6. Root cause and confidence
CONFIRMED — the only consumer of Idempotency-Key is backend/internal/handler/http/order_cancellation_handler.go:36 (item cancel); order creation handlers (public_order_handler.go CreateOrder, registerPOSOrderRoute backend/cmd/api/pos_routes.go:16) ignore the header although CORS advertises it; REL-001 re-run created 2 orders

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/order_cancellation_handler.go`:36, `backend/cmd/api/pos_routes.go`:16, `backend/internal/handler/http/public_order_handler.go`:39, `backend/internal/handler/http/public_order_handler.go`
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/products`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: POST /api/v1/pos/orders and POST /api/v1/public/.../orders: honour Idempotency-Key (additive header)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: duplicate revenue risk
- Operational: POS; Orders; Inventory workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: make order creation idempotent. (1) Add a table order_idempotency(tenant_id, key, order_id, request_hash, created_at) with a unique (tenant_id, key); (2) in the POS order handler registered at backend/cmd/api/pos_routes.go and public_order_handler.go CreateOrder, read the Idempotency-Key header (the only current consumer is order_cancellation_handler.go:36), insert the key inside the same transaction as the order, and on a repeated key return the original order with 200 and the same body; a key reused with a different payload returns 422; (3) the POS frontend generates one key per cart submission (UUID) and reuses it on retry/double-click.

## 10. Alternative designs
A: server-side idempotency table (recommended; protects every client). B: client-side double-submit guard only (does not protect retries or other clients). C: natural key on cart content within a time window (false positives for repeated identical orders).

## 11. Required policy decisions
None (technical control); retention period for keys (proposed 24 h) to be confirmed.

## 12. Required data migrations
Yes: idempotency table with TTL cleanup job.

## 13. Historical data implications
Existing duplicates (REL-001 evidence) are listed; cancel through the normal cancellation flow after owner review.

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
none and shared-file findings F-052

## 21. Risks
- Regression in modules sharing the same source files (backend/cmd/api/pos_routes.go, backend/internal/handler/http/order_cancellation_handler.go, backend/internal/handler/http/public_order_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Duplicate submit tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Posting the same cart twice with the same Idempotency-Key creates exactly one order and decrements stock once.
- A different payload with the same key returns 422 and creates nothing.
- Five concurrent identical submissions create one order.
