# F-040 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/usecase/auto_posting_engine.go:173 (inside PostCOGSEntry :143) skips every order item whose product has track_inventory=false (`if !product.TrackInventory { continue }`), so no COGS/inventory journal or sale movement is created; profitability report computes cogs from products.cost for all items. CreateProduct defaults track_inventory=false when omitted.


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/auto_posting_engine.go | 143-176 | backend/internal/usecase/auto_posting_engine.go | exact | yes |
| backend/internal/handler/http/product_handler.go | 53 | backend/internal/handler/http/product_handler.go | exact | yes |
| backend/migrations/180_order_stock_reservations.sql | 40-70 | backend/migrations/180_order_stock_reservations.sql | exact | yes |
| backend/internal/usecase/auto_posting_engine.go | 173 | backend/internal/usecase/auto_posting_engine.go | exact | yes |

