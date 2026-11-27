# F-050 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): handler/http/admin_order_handler.go:410 returns 500 "Failed to update order status" for every UpdateOrderStatus error, including the paid-order guard


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/admin_order_handler.go | 357-410 | backend/internal/handler/http/admin_order_handler.go | exact | yes |
| handler/http/admin_order_handler.go | 410 | backend/internal/handler/http/admin_order_handler.go | prefix | yes |

