# F-079 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + request capture): frontend/apps/dashboard/src/components/hr/LeaveForm.tsx:114-115 send new Date(...).toISOString(); backend/internal/handler/http/hr_leave_handler.go:186 parses the field as YYYY-MM-DD only


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/components/hr/LeaveForm.tsx | 114 | frontend/apps/dashboard/src/components/hr/LeaveForm.tsx | exact | yes |
| backend/internal/handler/http/hr_leave_handler.go | 186 | backend/internal/handler/http/hr_leave_handler.go | exact | yes |
| frontend/apps/dashboard/src/components/hr/LeaveForm.tsx | 114-115 | frontend/apps/dashboard/src/components/hr/LeaveForm.tsx | exact | yes |

