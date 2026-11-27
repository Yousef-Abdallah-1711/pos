# F-017 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source search): the only consumer of Idempotency-Key is backend/internal/handler/http/order_cancellation_handler.go:36 (item cancel); order creation handlers (public_order_handler.go CreateOrder, registerPOSOrderRoute backend/cmd/api/pos_routes.go:16) ignore the header although CORS advertises it; REL-001 re-run created 2 orders


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/order_cancellation_handler.go`:36.
- Propagation: affected modules POS; Orders; Inventory; related findings none recorded; observed consequence: 5 and 3 orders created; stock decremented each | RUN-003 REL-001: same order POST twice with identical Idempotency-Key creates 2 orders (REPRODUCED). 
- Cause versus symptom: the symptom is what users or reports show (5 and 3 orders created; stock decremented each | RUN-003 REL-001: same order POST twice with identical Idempotency-Key creates 2 orders (REP); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-REL-001; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/order_cancellation_handler.go | 36 | backend/internal/handler/http/order_cancellation_handler.go | exact | yes |
| backend/cmd/api/pos_routes.go | 16 | backend/cmd/api/pos_routes.go | exact | yes |
| backend/internal/handler/http/public_order_handler.go | 39 | backend/internal/handler/http/public_order_handler.go | exact | yes |
| public_order_handler.go | - | backend/internal/handler/http/public_order_handler.go | prefix | yes |

