# F-002 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source, RUN-003 read): backend/internal/usecase/financial_report_usecase.go:285-306 balance-sheet query selects only account_type IN (asset, liability, equity); revenue/expense accounts never contribute, so current-period earnings are not added to equity (BS difference 82,700 measured earlier)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/financial_report_usecase.go | 285-306 | backend/internal/usecase/financial_report_usecase.go | exact | yes |

