# F-045 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): accounting_advanced_repo.go:1105-1175 GetPredictiveCashFlow aggregates only payments_received and payments_made


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 1105-1175 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |

