# F-056 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): contract create/update validate only the date format (contract_handler.go:165-170, 221-226); no end >= start rule; resignation create has the same gap (resignation_repo.go:117-143)


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/contract_handler.go`.
- Propagation: affected modules HR contracts; HR resignations; related findings F-031, F-053, F-054; observed consequence: both return 201 and the invalid periods are stored (contract id 1, resignation id 1); duplicate contract number returns 400 with raw pq constraint text (F-031 class); attendance/leave/salary periods in the same module family are protected o
- Cause versus symptom: the symptom is what users or reports show (both return 201 and the invalid periods are stored (contract id 1, resignation id 1); duplicate contract number returns 400 with raw pq cons); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-CON-001; scripts: tests/runtime/ph10_hr_modules.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/contract_handler.go | - | backend/internal/handler/http/contract_handler.go | exact | yes |
| backend/internal/handler/http/resignation_handler.go | - | backend/internal/handler/http/resignation_handler.go | exact | yes |
| backend/internal/handler/http/contract_handler.go | 165 | backend/internal/handler/http/contract_handler.go | exact | yes |
| backend/internal/repository/resignation_repo.go | 117 | backend/internal/repository/resignation_repo.go | exact | yes |
| contract_handler.go | 165-170 | backend/internal/handler/http/contract_handler.go | prefix | yes |
| resignation_repo.go | 117-143 | backend/internal/repository/resignation_repo.go | prefix | yes |

