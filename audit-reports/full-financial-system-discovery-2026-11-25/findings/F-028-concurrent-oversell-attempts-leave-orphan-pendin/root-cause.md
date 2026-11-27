# F-028 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/usecase/order_usecase.go:168 CreateOrderWithItems commits the pending order and items first; POS confirmation/stock reservation runs afterwards (:201-205) and its failure returns an error (HTTP 400) while the draft order stays persisted ("complete draft is preserved for audit" comment), leaving orphan pending orders that stats and sales reports count (F-042)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/order_usecase.go | 166-215 | backend/internal/usecase/order_usecase.go | exact | yes |
| backend/internal/usecase/order_usecase.go | 168 | backend/internal/usecase/order_usecase.go | exact | yes |

