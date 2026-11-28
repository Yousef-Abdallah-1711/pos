# F-028 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/usecase/order_usecase.go:168 CreateOrderWithItems commits the pending order and items first; POS confirmation/stock reservation runs afterwards (:201-205) and its failure returns an error (HTTP 400) while the draft order stays persisted ("complete draft is preserved for audit" comment), leaving orphan pending orders that stats and sales reports count (F-042)


## Origin and propagation
- Origin (first verified source location): `backend/internal/usecase/order_usecase.go`:166-215.
- Propagation: affected modules POS; Orders; Inventory; KDS; related findings F-042; observed consequence: 5 accepted; 15 returned 400 yet 15 pending orders with order_items exist (single rejected request leaves none); orphans appear on KDS and in Sales Analytics
- Cause versus symptom: the symptom is what users or reports show (5 accepted; 15 returned 400 yet 15 pending orders with order_items exist (single rejected request leaves none); orphans appear on KDS and in); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/order_usecase.go | 166-215 | backend/internal/usecase/order_usecase.go | exact | yes |
| backend/internal/usecase/order_usecase.go | 168 | backend/internal/usecase/order_usecase.go | exact | yes |

