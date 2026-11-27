# F-046 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source per inventory 30 rows TURN-001/008): report_handler.go:557-645 / 2480-2609 reads legacy inventory.quantity_after else current stock and SUM(ABS(change)) from the legacy inventory table, while the live ledger is inventory_movements


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/report_handler.go | 557-645 | backend/internal/handler/http/report_handler.go | exact | yes |
| backend/internal/handler/http/report_handler.go | 2480-2609 | backend/internal/handler/http/report_handler.go | exact | yes |

