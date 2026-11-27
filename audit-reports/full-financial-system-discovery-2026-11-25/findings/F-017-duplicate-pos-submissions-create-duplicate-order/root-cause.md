# F-017 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source search): the only consumer of Idempotency-Key is backend/internal/handler/http/order_cancellation_handler.go:36 (item cancel); order creation handlers (public_order_handler.go CreateOrder, registerPOSOrderRoute backend/cmd/api/pos_routes.go:16) ignore the header although CORS advertises it; REL-001 re-run created 2 orders


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/order_cancellation_handler.go | 36 | backend/internal/handler/http/order_cancellation_handler.go | exact | yes |
| backend/cmd/api/pos_routes.go | 16 | backend/cmd/api/pos_routes.go | exact | yes |
| backend/internal/handler/http/public_order_handler.go | 39 | backend/internal/handler/http/public_order_handler.go | exact | yes |
| public_order_handler.go | - | backend/internal/handler/http/public_order_handler.go | prefix | yes |

