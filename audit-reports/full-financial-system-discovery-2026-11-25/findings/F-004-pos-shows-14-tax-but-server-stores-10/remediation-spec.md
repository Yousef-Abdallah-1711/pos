# F-004 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
POS shows 14% tax but server stores 10%.

## 2. Current behavior (verified)
UI total 11400 (tax 1400); stored order tax 1000 total 11000 | RUN-003 (AUD-POS-001, order 38): POS label Tax (14%) vs stored tax 6,000 on subtotal 60,000 (10%); UI displayed value not legible in capture (overlay) - rate mismatch STILL REPRODUCIBLE by label vs stored; earlier evidence holds values. 

## 3. Required behavior
Displayed total equals stored total

## 4. Reproduction conditions
POS add item price 10000 then Place Order

## 5. Verified evidence
- evidence/runtime/06
- network pos-POS-C-001.json | evidence/screenshots/RUN-003/PH-07/AUD-POS-001-before-submit.png | evidence/network/RUN-003/PH-07/AUD-POS-001.json | evidence/database/RUN-003/PH-07/AUD-POS-001-before-after.txt | tests/runtime/ph07_pos_cash_headed.mjs | evidence/screenshots/RUN-003/PH-07/AUD-POS-002-cart-totals-overlay-hidden.png | evidence/network/RUN-003/PH-07/AUD-POS-002.json | tests/runtime/ph07_pos_ui_totals.mjs

## 6. Root cause and confidence
CONFIRMED — frontend/apps/dashboard/src/stores/posStore.ts:107 hard-codes TAX_RATE = 0.14 for the cart; backend/internal/usecase/order_pricing.go:46 computes the stored tax as 10% (subtotal+5)/10 and ignores any client tax

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/stores/posStore.ts`:107, `backend/internal/usecase/order_pricing.go`:40-52, `backend/internal/usecase/order_pricing.go`:46
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: GET /api/v1/pos/pricing-config (new, read-only); POST /api/v1/pos/orders: unchanged contract; response keeps tax_amount and total_amount
- Accounting logic / posting rules: none named
- Reports affected: POS; Orders; Accounting

## 8. Impact
- Financial: 400 per 10000 order (customer shown more than recorded)
- Operational: POS; Orders; Accounting workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: make the server the only tax authority. (1) Add a read endpoint (proposed GET /api/v1/pos/pricing-config) returning the effective tax rate and currency from tenant settings; (2) change frontend/apps/dashboard/src/stores/posStore.ts:107 to read that value instead of the hard-coded TAX_RATE = 0.14 and label it with the returned rate; (3) change backend/internal/usecase/order_pricing.go:46 to compute tax from the same tenant setting (default 14% only if policy POL-02 approves) instead of the fixed (subtotal+5)/10 expression; (4) keep rejecting client-supplied tax (already ignored) and return the server-computed preview from POST /pos/orders/preview (proposed) so the UI total equals the stored total.

## 10. Alternative designs
A: tenant-level tax setting (recommended; also fixes F-037/no tax setting). B: hard-code the same 14% on both sides (fast, but the rate stays unchangeable and unverified). C: compute tax only on the server and show no tax line in the cart until submit (worse UX).

## 11. Required policy decisions
POL-02 (tax base and display) in 56, plus POL-11 for the currency symbol; no implementation before approval.

## 12. Required data migrations
Possibly: a tenant tax-rate setting column/table (proposed) with a default equal to the approved policy; backfill existing tenants to the policy default.

## 13. Historical data implications
Existing orders keep the stored 10% tax. Produce a read-only exception report of orders whose stored rate differs from the approved rate; any correction is a reversing/adjusting journal approved by the accountant.

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
- Regression in modules sharing the same source files (backend/internal/usecase/order_pricing.go, frontend/apps/dashboard/src/stores/posStore.ts).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
UI vs server total equality tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- For subtotal 60,000 with the approved rate r, the POS cart tax, the order response tax_amount and orders.tax_amount are all equal to ROUND(60,000 x r, 2).
- Changing the tenant tax rate changes the POS preview and the stored tax without a code change.
- A mismatch between the UI preview and the stored total is impossible: contract test posts the cart and compares totals.
