# F-054 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/hr_leave_handler.go | 132-200 | exact |
| backend/internal/handler/http/hr_leave_handler.go | 207 | exact |
| backend/internal/handler/http/hr_leave_handler.go | 328 | exact |
| backend/internal/handler/http/hr_leave_handler.go | 207-210 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr | exists |

## Database objects involved
`leaves`

## API contracts involved
- POST/PUT /api/v1/hr/leaves: total_days ignored on input (returned computed)

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
