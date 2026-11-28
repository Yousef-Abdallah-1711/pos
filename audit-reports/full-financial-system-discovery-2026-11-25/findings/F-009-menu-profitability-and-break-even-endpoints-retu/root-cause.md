# F-009 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + SQL error reproduced): backend/internal/repository/accounting_advanced_repo.go:976 and :1045 select p.cost_price, a column that does not exist on products (the column is cost); error mapped to generic 500 at accounting_advanced_handler.go:782


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/accounting_advanced_repo.go`:976.
- Propagation: affected modules Financial reports advanced; related findings none recorded; observed consequence: HTTP 500 for both (cost_price column absent) | RUN-003: GET /reports/financial/menu-profitability and /break-even still 500. ROOT CAUSE CONFIRMED: repository/accounting_advanced_repo.go:976 and :1045 select p.cost_price, a column that does 
- Cause versus symptom: the symptom is what users or reports show (HTTP 500 for both (cost_price column absent) | RUN-003: GET /reports/financial/menu-profitability and /break-even still 500. ROOT CAUSE CONF); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 976 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/repository/accounting_advanced_repo.go | 1045 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/handler/http/accounting_advanced_handler.go | 782 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |

