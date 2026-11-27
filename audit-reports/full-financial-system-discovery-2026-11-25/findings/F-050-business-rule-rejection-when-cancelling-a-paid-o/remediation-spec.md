# F-050 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Business-rule rejection when cancelling a paid order via PUT /admin/orders/{id}/status returns HTTP 500 instead of a 4xx.

## 2. Current behavior (verified)
500 "Failed to update order status" (DELETE returns 400 with the explicit reason) | RUN-003: PUT status cancelled on paid order -> 500.  | RUN-003: PUT status cancelled on paid order -> 500. Same class: unsigned Paymob webhook -> 500 and refund of unknown payment -> 500 (AUD-SAAS-MATRIX PAY-02/PAY-03). 

## 3. Required behavior
409/400 with the same message as DELETE

## 4. Reproduction conditions
PUT /admin/orders/{paid id}/status {status:cancelled}

## 5. Verified evidence
- evidence/database/RUN-003/PH-19/AUD-REFUND-01-paid-order-cancel.json
- evidence/database/RUN-003/PH-19/AUD-REG-HC10-FIN019-CFR006.json
- tests/runtime/ph19_paid_cancel.mjs
- tests/runtime/ph19_regs2.mjs
- evidence/database/RUN-003/PH-07/AUD-ORDER-TYPES.json
- tests/runtime/ph07_order_types.mjs

## 6. Root cause and confidence
CONFIRMED — handler/http/admin_order_handler.go:410 returns 500 "Failed to update order status" for every UpdateOrderStatus error, including the paid-order guard

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/admin_order_handler.go`:357-410, `backend/internal/handler/http/admin_order_handler.go`:410
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: PUT /api/v1/admin/orders/{id}/status: 409 for guard violations (was 500)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0
- Operational: Orders API workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: map domain guard errors to client errors. Introduce sentinel errors (for example ErrPaidOrderCannotBeCancelled) in the usecase layer and translate them in handler/http/admin_order_handler.go:410 to 409 with a stable code; keep 500 only for unexpected errors; add table-driven tests for every guard.

## 10. Alternative designs
A: sentinel errors + mapper (recommended; shares the F-025 helper). B: string matching on messages (fragile).

## 11. Required policy decisions
None.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

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
F-049

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/admin_order_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Paid-order cancel through every endpoint returns the same 4xx (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Cancelling a paid order via PUT /admin/orders/{id}/status returns 409 with code ORDER_PAID_CANNOT_CANCEL.
- Unexpected failures still return 500 with a request id.
