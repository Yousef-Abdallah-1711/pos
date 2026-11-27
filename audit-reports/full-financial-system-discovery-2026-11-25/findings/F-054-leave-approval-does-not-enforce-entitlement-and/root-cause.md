# F-054 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): total_days is taken from the client when non-zero (hr_leave_handler.go:207-210; update :328-329) and no allowance or date-consistency check exists


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/hr_leave_handler.go | 132-200 | backend/internal/handler/http/hr_leave_handler.go | exact | yes |
| backend/internal/handler/http/hr_leave_handler.go | 207 | backend/internal/handler/http/hr_leave_handler.go | exact | yes |
| backend/internal/handler/http/hr_leave_handler.go | 328 | backend/internal/handler/http/hr_leave_handler.go | exact | yes |
| hr_leave_handler.go | 207-210 | backend/internal/handler/http/hr_leave_handler.go | prefix | yes |

