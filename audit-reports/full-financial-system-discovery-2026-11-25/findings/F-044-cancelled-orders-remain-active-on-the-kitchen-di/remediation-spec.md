# F-044 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Cancelled orders remain active on the Kitchen Display queue.

## 2. Current behavior (verified)
Orders 34, 40, 42 are status cancelled but kds_queue rows stay status new and are listed in /kds/queue and on the board (ACTIVE ORDERS 25 / NEW 24 include them) | RUN-003 headed Chrome: 3 cancelled orders visible on KDS board (ACTIVE 25 / NEW 24).  RUN-005 real-UI reproduction: cancelling unpaid order 5796 through the Orders page (confirm dialog, DELETE 200, status cancelled, stock restored 30 -> 32) left its kds_queue row in status "new".

## 3. Required behavior
Cancelled order is removed from the active kitchen queue

## 4. Reproduction conditions
Cancel a pending order (DELETE /admin/orders/{id}) and open the KDS board (headed Chrome)

## 5. Verified evidence
- evidence/database/RUN-003/PH-12/AUD-MET-ORD-KPI.txt
- evidence/network/RUN-003/PH-08/AUD-KDS-001.json
- evidence/screenshots/RUN-003/PH-08/AUD-KDS-001-board.png
- tests/runtime/ph08_kds_cancelled.mjs
- evidence/network/RUN-005/PH-04/AUD-R5-DESTRUCTIVE-A.run1.json
- tests/runtime/ph04_r5_destructive_a.mjs

## 6. Root cause and confidence
CONFIRMED — no trigger or code path updates kds_queue when an order is cancelled (pg_trigger on orders lists only order_stock_transition/order_stock_items); RUN-004 shows 3 cancelled orders still active in the queue (AUD-R4-KDS-after-read.json)

## 7. Affected paths
- Backend/database: `backend/internal/repository/kds_repo.go`, `backend/migrations/124_create_kds_queue.sql`, `backend/internal/repository/kds_repo.go`:307
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/kds`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: KDS queue/stats exclude cancelled (behaviour change)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: Operational: food prepared for cancelled orders (COGS without sale)
- Operational: KDS; Orders workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: propagate cancellation to the kitchen. When an order is cancelled (usecase/order_usecase.go cancellation path and repository/order_repo.go status update) update the matching kds_queue row to status cancelled inside the same transaction (or add a trigger on orders.status that does so); exclude cancelled rows from /kds/queue, /kds/stats active_orders and analytics; broadcast the change over the KDS channel so open screens remove the card.

## 10. Alternative designs
A: update in the cancellation transaction + broadcast (recommended). B: DB trigger only (covers every writer, no broadcast). C: filter at read time by joining orders (hides the data inconsistency).

## 11. Required policy decisions
None.

## 12. Required data migrations
Yes if the trigger option is used; plus a one-off update of existing active rows whose order is cancelled (3 in the evidence).

## 13. Historical data implications
Existing stale rows are listed and set to cancelled after review.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-028 and shared-file findings F-070

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/repository/kds_repo.go, backend/migrations/124_create_kds_queue.sql).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Cancel order and assert KDS queue/board and counts (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Cancelling an order removes it from the KDS board within 2 s and from active_orders.
- kds_queue has no active row whose order is cancelled (SQL assertion over the fixture).
- Analytics exclude cancelled orders.
