# F-054 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): total_days is taken from the client when non-zero (hr_leave_handler.go:207-210; update :328-329) and no allowance or date-consistency check exists


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/hr_leave_handler.go`:132-200.
- Propagation: affected modules HR leaves; leave-balance report; related findings none recorded; observed consequence: 30-day leave approved; report shows used 30, remaining 0 (clamped) vs allowance 21; 3-day range stored total_days 40 | RUN-003: 30-day annual leave approved over allowance 21; total_days 40 stored for a 3-day range (REPRODUCED). 
- Cause versus symptom: the symptom is what users or reports show (30-day leave approved; report shows used 30, remaining 0 (clamped) vs allowance 21; 3-day range stored total_days 40 | RUN-003: 30-day annua); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-LEV-001; scripts: tests/runtime/ph10_hr_matrix.mjs, tests/runtime/ph10_hr_matrix2.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/hr_leave_handler.go | 132-200 | backend/internal/handler/http/hr_leave_handler.go | exact | yes |
| backend/internal/handler/http/hr_leave_handler.go | 207 | backend/internal/handler/http/hr_leave_handler.go | exact | yes |
| backend/internal/handler/http/hr_leave_handler.go | 328 | backend/internal/handler/http/hr_leave_handler.go | exact | yes |
| hr_leave_handler.go | 207-210 | backend/internal/handler/http/hr_leave_handler.go | prefix | yes |

