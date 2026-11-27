# F-004 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): frontend/apps/dashboard/src/stores/posStore.ts:107 hard-codes TAX_RATE = 0.14 for the cart; backend/internal/usecase/order_pricing.go:46 computes the stored tax as 10% (subtotal+5)/10 and ignores any client tax


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/stores/posStore.ts | 107 | frontend/apps/dashboard/src/stores/posStore.ts | exact | yes |
| backend/internal/usecase/order_pricing.go | 40-52 | backend/internal/usecase/order_pricing.go | exact | yes |
| backend/internal/usecase/order_pricing.go | 46 | backend/internal/usecase/order_pricing.go | exact | yes |

