# F-046 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source per inventory 30 rows TURN-001/008): report_handler.go:557-645 / 2480-2609 reads legacy inventory.quantity_after else current stock and SUM(ABS(change)) from the legacy inventory table, while the live ledger is inventory_movements


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/report_handler.go`:557-645.
- Propagation: affected modules Inventory turnover report; Overview report; related findings none recorded; observed consequence: opening_stock_qty == current_stock_qty (33 and 21) so average inventory equals current stock; waste_qty 0 although inventory_movements has waste 1 for products 2 and 4; Audit Burger original stock 50 | RUN-003: turnover opening_stock_qty ==
- Cause versus symptom: the symptom is what users or reports show (opening_stock_qty == current_stock_qty (33 and 21) so average inventory equals current stock; waste_qty 0 although inventory_movements has w); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: tests/runtime/ph12_batch_apis2.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/report_handler.go | 557-645 | backend/internal/handler/http/report_handler.go | exact | yes |
| backend/internal/handler/http/report_handler.go | 2480-2609 | backend/internal/handler/http/report_handler.go | exact | yes |

