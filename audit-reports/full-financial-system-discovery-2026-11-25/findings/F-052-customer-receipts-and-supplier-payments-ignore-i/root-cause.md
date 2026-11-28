# F-052 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source search): the only Idempotency-Key consumer in the backend is order_cancellation_handler.go (item cancel); payment_acct handler/repository/usecase contain no key handling although the CORS allow-list advertises the header


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/payment_acct_handler.go`.
- Propagation: affected modules Receipts; Supplier payments; Accounting (cash, AR/AP); related findings F-017; observed consequence: Both return 201; payments_received has 2 rows (sum 800); invoice amount_paid 800 / balance 200; AR credited twice | RUN-003 REL-025: identical Idempotency-Key replay creates 2 receipts (sum 800) (REPRODUCED). 
- Cause versus symptom: the symptom is what users or reports show (Both return 201; payments_received has 2 rows (sum 800); invoice amount_paid 800 / balance 200; AR credited twice | RUN-003 REL-025: identic); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-REL-025; scripts: tests/runtime/ph18_reliability.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/payment_acct_handler.go | - | backend/internal/handler/http/payment_acct_handler.go | exact | yes |
| backend/internal/repository/payment_acct_repo.go | - | backend/internal/repository/payment_acct_repo.go | exact | yes |
| backend/internal/handler/http/order_cancellation_handler.go | 36 | backend/internal/handler/http/order_cancellation_handler.go | exact | yes |
| order_cancellation_handler.go | - | backend/internal/handler/http/order_cancellation_handler.go | prefix | yes |

