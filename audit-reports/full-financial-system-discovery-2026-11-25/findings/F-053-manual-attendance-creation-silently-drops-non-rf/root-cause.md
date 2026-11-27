# F-053 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): clock_in/clock_out are parsed only as RFC3339 and a failed parse is silently ignored (hr_attendance_handler.go:235-243), so "09:00" is dropped and the row stores NULL times


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/hr_attendance_handler.go`:188-260.
- Propagation: affected modules HR attendance; HR overtime/absenteeism reports; payroll overtime input; related findings F-031; observed consequence: "09:00" form -> 201 with clock_in/clock_out NULL and regular_hours 0 / overtime_hours 0; reversed order on the same form also 201; RFC3339 form computes 8.00 + 2.50; reversed RFC3339 -> 500 with raw constraint text valid_clock_times; overti
- Cause versus symptom: the symptom is what users or reports show ("09:00" form -> 201 with clock_in/clock_out NULL and regular_hours 0 / overtime_hours 0; reversed order on the same form also 201; RFC3339 f); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-ATT-001; scripts: tests/runtime/ph10_hr_matrix.mjs, tests/runtime/ph10_hr_matrix2.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/hr_attendance_handler.go | 188-260 | backend/internal/handler/http/hr_attendance_handler.go | exact | yes |
| backend/internal/handler/http/hr_attendance_handler.go | 235 | backend/internal/handler/http/hr_attendance_handler.go | exact | yes |
| hr_attendance_handler.go | 235-243 | backend/internal/handler/http/hr_attendance_handler.go | prefix | yes |

