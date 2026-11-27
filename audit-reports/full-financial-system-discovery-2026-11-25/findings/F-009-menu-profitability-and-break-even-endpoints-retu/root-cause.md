# F-009 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + SQL error reproduced): backend/internal/repository/accounting_advanced_repo.go:976 and :1045 select p.cost_price, a column that does not exist on products (the column is cost); error mapped to generic 500 at accounting_advanced_handler.go:782


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 976 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/repository/accounting_advanced_repo.go | 1045 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/handler/http/accounting_advanced_handler.go | 782 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |

