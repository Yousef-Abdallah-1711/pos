# F-034 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/handler/http/inventory_adjustment_handler.go (Adjust :29-105) updates stock and writes an inventory movement but contains no journal/posting call (0 references); waste therefore posts nothing, and turnover reads the legacy table (F-046)


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/inventory_adjustment_handler.go`:29-105.
- Propagation: affected modules Inventory; Reports; related findings F-046; observed consequence: Stock 48 to 47; no journal; waste_qty 0; reason count_shortage not accepted (allowed: purchase, sale, adjustment, waste, return)
- Cause versus symptom: the symptom is what users or reports show (Stock 48 to 47; no journal; waste_qty 0; reason count_shortage not accepted (allowed: purchase, sale, adjustment, waste, return)); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/inventory_adjustment_handler.go | 29-105 | backend/internal/handler/http/inventory_adjustment_handler.go | exact | yes |

