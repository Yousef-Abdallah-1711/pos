# F-011 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/usecase/financial_report_usecase.go:488-526 GenerateInventoryValuation = SUM(quantity_in_stock x cost) for tracked products; it never reads GL 1130, and purchases/opening stock never post to 1130 (supplier invoice lines post to 5110, OBS-PUR-001)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/financial_report_usecase.go | 488-526 | backend/internal/usecase/financial_report_usecase.go | exact | yes |

