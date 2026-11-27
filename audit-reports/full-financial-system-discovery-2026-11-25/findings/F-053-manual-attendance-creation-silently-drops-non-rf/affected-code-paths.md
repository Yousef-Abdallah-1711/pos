# F-053 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/hr_attendance_handler.go | 188-260 | exact |
| backend/internal/handler/http/hr_attendance_handler.go | 235 | exact |
| backend/internal/handler/http/hr_attendance_handler.go | 235-243 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr | exists |

## Database objects involved
`attendance`

## API contracts involved
- POST /api/v1/hr/attendance: 400 field errors; response includes total_hours

## Proposed (future) files
- backend/migrations/<next>_attendance_clock_check.sql (new)
