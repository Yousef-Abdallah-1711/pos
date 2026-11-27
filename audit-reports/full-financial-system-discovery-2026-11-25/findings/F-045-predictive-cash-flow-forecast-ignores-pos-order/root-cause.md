# F-045 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): accounting_advanced_repo.go:1105-1175 GetPredictiveCashFlow aggregates only payments_received and payments_made


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/accounting_advanced_repo.go`:1105-1175.
- Propagation: affected modules Financial hub predictive cash flow; related findings F-001, F-003; observed consequence: 2026-11 history cash_in 13,000 / cash_out 32,400 -> forecast Dec -18,430 while 319,660 of paid POS orders and GL cash +276,270 exist | RUN-003: PCF Nov cash_in 13,000/cash_out 32,400 vs paid POS orders 319,660; forecast negative (CONFIRMED 
- Cause versus symptom: the symptom is what users or reports show (2026-11 history cash_in 13,000 / cash_out 32,400 -> forecast Dec -18,430 while 319,660 of paid POS orders and GL cash +276,270 exist | RUN-0); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: tests/runtime/ph12_batch_apis.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 1105-1175 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |

