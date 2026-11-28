# F-050 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): handler/http/admin_order_handler.go:410 returns 500 "Failed to update order status" for every UpdateOrderStatus error, including the paid-order guard


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/admin_order_handler.go`:357-410.
- Propagation: affected modules Orders API; related findings F-049; observed consequence: 500 "Failed to update order status" (DELETE returns 400 with the explicit reason) | RUN-003: PUT status cancelled on paid order -> 500.  | RUN-003: PUT status cancelled on paid order -> 500. Same class: unsigned Paymob webhook -> 500 and re
- Cause versus symptom: the symptom is what users or reports show (500 "Failed to update order status" (DELETE returns 400 with the explicit reason) | RUN-003: PUT status cancelled on paid order -> 500.  | R); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-HC10, EX-RUN003-AUD-SAAS-006, EX-RUN003-AUD-OT-003, EX-RUN005-DESTR-ORDER-CANCEL; scripts: tests/runtime/ph19_paid_cancel.mjs, tests/runtime/ph19_regs2.mjs, tests/runtime/ph07_order_types.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/admin_order_handler.go | 357-410 | backend/internal/handler/http/admin_order_handler.go | exact | yes |
| handler/http/admin_order_handler.go | 410 | backend/internal/handler/http/admin_order_handler.go | prefix | yes |

