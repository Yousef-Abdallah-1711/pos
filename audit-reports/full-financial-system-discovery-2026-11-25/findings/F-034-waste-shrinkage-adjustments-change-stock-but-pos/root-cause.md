# F-034 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/handler/http/inventory_adjustment_handler.go (Adjust :29-105) updates stock and writes an inventory movement but contains no journal/posting call (0 references); waste therefore posts nothing, and turnover reads the legacy table (F-046)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/inventory_adjustment_handler.go | 29-105 | backend/internal/handler/http/inventory_adjustment_handler.go | exact | yes |

