# F-057 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Delivered cash-on-delivery orders stay unpaid and no settlement step records the cash collected by the driver.

## 2. Current behavior (verified)
status delivered with payment_status pending and 0 journals; revenue stays unrecognised until someone manually marks paid; no driver cash reconciliation exists | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): delivered COD order remained payment_status pending with 0 journals (REPRODUCED). 

## 3. Required behavior
A COD settlement action (driver cash handover) marks the order paid and posts revenue/COGS, or delivery completion prompts for payment

## 4. Reproduction conditions
Create delivery order (cash), assign driver, move status preparing -> ready -> out_for_delivery -> delivered, then read payment_status and journals

## 5. Verified evidence
- evidence/database/RUN-003/PH-08/AUD-THEME-DRIVER.json
- evidence/database/RUN-003/PH-08/AUD-THEME2.json
- tests/runtime/ph08_theme_driver.mjs
- tests/runtime/ph08_theme2.mjs

## 6. Root cause and confidence
CONFIRMED — CompleteDelivery only changes the order status to delivered (order_usecase.go:433-447) and records no collection or payment

## 7. Affected paths
- Backend/database: `backend/internal/usecase/order_usecase.go`:433-447, `backend/internal/handler/http/admin_driver_handler.go`, `backend/internal/usecase/order_usecase.go`:433
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: POST /api/v1/admin/orders/{id}/complete-delivery: accepts collected_amount (additive); POST /api/v1/admin/drivers/{id}/settlements (new)
- Accounting logic / posting rules: none named
- Reports affected: Delivery; Drivers; Orders; Accounting

## 8. Impact
- Financial: All unsettled COD sales are missing from GL revenue/VAT/cash
- Operational: Delivery; Drivers; Orders; Accounting workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: cash-on-delivery settlement. Extend OrderUseCase.CompleteDelivery (backend/internal/usecase/order_usecase.go:433-447) so completing a COD delivery records the amount the driver collected (field, default order total) and creates a driver cash liability (Dr driver cash receivable / Cr revenue-and-VAT through the normal sales posting); add a driver settlement operation that records the cash handed in (Dr cash / Cr driver receivable) and reports variances per driver and day.

## 10. Alternative designs
A: collected amount at delivery + driver settlement document (recommended). B: mark paid at delivery without tracking the driver (cash risk invisible).

## 11. Required policy decisions
Proposed POL-18 (driver cash handling and shortage treatment) to be added to 56.

## 12. Required data migrations
Yes: driver settlement tables and a driver receivable account mapping.

## 13. Historical data implications
Existing delivered-unpaid COD orders are listed for reconciliation with drivers.

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
none and shared-file findings F-028, F-049, F-058

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/admin_driver_handler.go, backend/internal/usecase/order_usecase.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Delivered COD order can be settled once; cash and revenue reconcile (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Completing a COD order records the collected amount and posts the sales journal once.
- Settling a driver for a day moves exactly the handed-in cash and shows any shortage.
- Unsettled driver cash is visible in a report.
