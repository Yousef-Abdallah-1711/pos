# F-040 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/usecase/auto_posting_engine.go:173 (inside PostCOGSEntry :143) skips every order item whose product has track_inventory=false (`if !product.TrackInventory { continue }`), so no COGS/inventory journal or sale movement is created; profitability report computes cogs from products.cost for all items. CreateProduct defaults track_inventory=false when omitted.


## Origin and propagation
- Origin (first verified source location): `backend/internal/usecase/auto_posting_engine.go`:143-176.
- Propagation: affected modules Products; POS/Orders; Accounting; Profitability report; related findings F-011; observed consequence: Order 39 (untracked): JE-2026-00064 posts only Dr 1110 330 / Cr 4100 300 / Cr 2120 30 - no COGS; tracked order 41: JE-2026-00066 Dr 5110 120 / Cr 1130 120. Profitability report shows cogs 120 for both products => report COGS exceeds GL COGS
- Cause versus symptom: the symptom is what users or reports show (Order 39 (untracked): JE-2026-00064 posts only Dr 1110 330 / Cr 4100 300 / Cr 2120 30 - no COGS; tracked order 41: JE-2026-00066 Dr 5110 120); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-OBS-INV-001, EX-RUN003-AUD-INV-COGS-UNTRACKED; scripts: tests/runtime/ph06_inventory_controlled.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/auto_posting_engine.go | 143-176 | backend/internal/usecase/auto_posting_engine.go | exact | yes |
| backend/internal/handler/http/product_handler.go | 53 | backend/internal/handler/http/product_handler.go | exact | yes |
| backend/migrations/180_order_stock_reservations.sql | 40-70 | backend/migrations/180_order_stock_reservations.sql | exact | yes |
| backend/internal/usecase/auto_posting_engine.go | 173 | backend/internal/usecase/auto_posting_engine.go | exact | yes |

