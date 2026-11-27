# F-004 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): frontend/apps/dashboard/src/stores/posStore.ts:107 hard-codes TAX_RATE = 0.14 for the cart; backend/internal/usecase/order_pricing.go:46 computes the stored tax as 10% (subtotal+5)/10 and ignores any client tax


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/stores/posStore.ts`:107.
- Propagation: affected modules POS; Orders; Accounting; related findings none recorded; observed consequence: UI total 11400 (tax 1400); stored order tax 1000 total 11000 | RUN-003 (AUD-POS-001, order 38): POS label Tax (14%) vs stored tax 6,000 on subtotal 60,000 (10%); UI displayed value not legible in capture (overlay) - rate mismatch STILL REPR
- Cause versus symptom: the symptom is what users or reports show (UI total 11400 (tax 1400); stored order tax 1000 total 11000 | RUN-003 (AUD-POS-001, order 38): POS label Tax (14%) vs stored tax 6,000 on s); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-POS-001, EX-RUN003-AUD-POS-002; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/stores/posStore.ts | 107 | frontend/apps/dashboard/src/stores/posStore.ts | exact | yes |
| backend/internal/usecase/order_pricing.go | 40-52 | backend/internal/usecase/order_pricing.go | exact | yes |
| backend/internal/usecase/order_pricing.go | 46 | backend/internal/usecase/order_pricing.go | exact | yes |

