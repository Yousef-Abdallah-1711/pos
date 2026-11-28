# F-079 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + request capture): frontend/apps/dashboard/src/components/hr/LeaveForm.tsx:114-115 send new Date(...).toISOString(); backend/internal/handler/http/hr_leave_handler.go:186 parses the field as YYYY-MM-DD only


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/components/hr/LeaveForm.tsx`:114.
- Propagation: affected modules HR > Leave requests (dashboard form, API); related findings F-014, F-054, F-074; observed consequence: POST /hr/leaves returns 400 {"error":"Invalid start_date format (expected YYYY-MM-DD)"}; the page shows "Request failed with status code 400" for every valid submission; the request body carries start_date "2026-11-20T00:00:00.000Z"
- Cause versus symptom: the symptom is what users or reports show (POST /hr/leaves returns 400 {"error":"Invalid start_date format (expected YYYY-MM-DD)"}; the page shows "Request failed with status code 400); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-FORMMATRIX2; scripts: tests/runtime/ph06_r4_form_matrix.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/components/hr/LeaveForm.tsx | 114 | frontend/apps/dashboard/src/components/hr/LeaveForm.tsx | exact | yes |
| backend/internal/handler/http/hr_leave_handler.go | 186 | backend/internal/handler/http/hr_leave_handler.go | exact | yes |
| frontend/apps/dashboard/src/components/hr/LeaveForm.tsx | 114-115 | frontend/apps/dashboard/src/components/hr/LeaveForm.tsx | exact | yes |

