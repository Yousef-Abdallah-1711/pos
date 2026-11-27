# F-042 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Order stats "total revenue" includes cancelled and unpaid orders and "delivered today" counts every order created today.

## 2. Current behavior (verified)
total_revenue 529,100 = SUM of ALL 41 orders incl 3 cancelled (11,440) and 21 unpaid; delivered_today 30 while delivered_count 0 (counts any order created today) | RUN-003: order stats total_revenue 529,100 = all 41 orders incl 3 cancelled; delivered_today 30 vs delivered_count 0 (CONFIRMED order_repo.go:859-870). 

## 3. Required behavior
Revenue excludes cancelled (and unpaid per chosen policy); delivered_today counts delivered orders

## 4. Reproduction conditions
GET /admin/orders/stats on tenant with cancelled, unpaid and non-delivered orders (RUN-003 batch capture)

## 5. Verified evidence
- evidence/database/RUN-003/PH-12/AUD-MET-ORD-KPI.txt
- evidence/network/RUN-003/PH-12/AUD-MET-BATCH.json
- tests/runtime/ph12_batch_apis.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/repository/order_repo.go:859-870 GetOrderStats sums total_amount with no status/payment filter and computes delivered_today as COUNT(created_at >= CURRENT_DATE)

## 7. Affected paths
- Backend/database: `backend/internal/repository/order_repo.go`:859-950, `backend/internal/repository/order_repo.go`:859-870
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: GET /api/v1/admin/orders/stats: adds booked_revenue and collected_revenue (additive); total_revenue follows POL-01
- Accounting logic / posting rules: none named
- Reports affected: Orders dashboard; Order stats API

## 8. Impact
- Financial: 529,100 shown vs 319,660 paid tax-inclusive vs 305,300 GL revenue
- Operational: Orders dashboard; Order stats API workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: make order statistics follow the approved revenue basis. In backend/internal/repository/order_repo.go GetOrderStats (:859-870) filter total revenue by status <> cancelled and payment_status = paid (or the basis approved in POL-01) and expose both values (booked vs collected); compute delivered_today from status = delivered and a delivered_at timestamp (add the column if missing) instead of COUNT(created_at >= CURRENT_DATE); use the restaurant time zone for the day boundary (see F-070).

## 10. Alternative designs
A: filter by basis and add delivered_at (recommended). B: filter only (delivered_today stays wrong).

## 11. Required policy decisions
POL-01 (revenue recognition basis) in 56.

## 12. Required data migrations
Possibly orders.delivered_at (nullable, back-filled from status history if available).

## 13. Historical data implications
None stored; numbers change on deploy (communicate to owners).

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
none and shared-file findings F-049

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/repository/order_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Stats with cancelled/unpaid/delivered orders must equal independent counts (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- For a dataset with cancelled and unpaid orders the stats revenue equals the approved basis computed independently in SQL.
- delivered_today counts only orders delivered today.
- The dashboard tiles use the same endpoint and show the same numbers.
