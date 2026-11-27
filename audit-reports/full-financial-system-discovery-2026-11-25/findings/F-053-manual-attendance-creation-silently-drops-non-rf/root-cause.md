# F-053 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): clock_in/clock_out are parsed only as RFC3339 and a failed parse is silently ignored (hr_attendance_handler.go:235-243), so "09:00" is dropped and the row stores NULL times


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/hr_attendance_handler.go | 188-260 | backend/internal/handler/http/hr_attendance_handler.go | exact | yes |
| backend/internal/handler/http/hr_attendance_handler.go | 235 | backend/internal/handler/http/hr_attendance_handler.go | exact | yes |
| hr_attendance_handler.go | 235-243 | backend/internal/handler/http/hr_attendance_handler.go | prefix | yes |

