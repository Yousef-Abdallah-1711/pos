# F-049 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): usecase/order_usecase.go:341 and repository/order_repo.go:634 reject cancelling paid/refunded orders with "before a supported refund and accounting reversal workflow exists"; domain/order.go:219 CanTransitionOrderPaymentStatus forbids refunded and undoing paid; no refund use case exists


## Origin and propagation
- Origin (first verified source location): `backend/internal/usecase/order_usecase.go`:341.
- Propagation: affected modules Orders; Payments; Accounting (revenue/COGS/VAT/inventory reversal); related findings F-005, F-026; observed consequence: DELETE -> 400 "cannot cancel paid or refunded order before a supported refund and accounting reversal workflow exists"; payment-status refunded/failed/pending from paid -> 400 invalid transition; voided -> "invalid payment status"; status=c
- Cause versus symptom: the symptom is what users or reports show (DELETE -> 400 "cannot cancel paid or refunded order before a supported refund and accounting reversal workflow exists"; payment-status refun); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-HC10, EX-RUN005-DESTR-ORDER-CANCEL; scripts: tests/runtime/ph19_paid_cancel.mjs, tests/runtime/ph19_regs2.mjs | evidence/database/RUN-003/PH-07/AUD-ORDER-TYPES.json, tests/runtime/ph04_r5_destructive_a.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/order_usecase.go | 341 | backend/internal/usecase/order_usecase.go | exact | yes |
| backend/internal/repository/order_repo.go | 634 | backend/internal/repository/order_repo.go | exact | yes |
| backend/internal/domain/order.go | 215-240 | backend/internal/domain/order.go | exact | yes |
| backend/migrations/175_durable_payment_refund_state.sql | - | backend/migrations/175_durable_payment_refund_state.sql | exact | yes |
| domain/order.go | 219 | backend/internal/domain/order.go | prefix | yes |

