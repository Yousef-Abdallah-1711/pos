# F-068 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED: accounting_advanced_repo.go:1421-1446 LEFT JOIN insurance_tracking it ON ... AND it.status active has no DISTINCT or per-employee limit


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 1437 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/handler/http/accounting_advanced_handler.go | 912 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |
| accounting_advanced_repo.go | 1421-1446 | backend/internal/repository/accounting_advanced_repo.go | prefix | yes |

