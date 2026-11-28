# F-069 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/resignation_repo.go | 201 | exact |
| backend/internal/repository/employee_repo.go | 370 | exact |
| backend/internal/handler/http/hr_report_handler.go | 276 | exact |
| backend/internal/repository/resignation_repo.go | 201-215 | prefix |
| backend/internal/repository/employee_repo.go | 370-380 | prefix |
| backend/internal/handler/http/hr_report_handler.go | 276-285 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr | exists |

## Database objects involved
`attendance`, `employees`

## API contracts involved
- POST /api/v1/hr/resignations/{id}/complete: also updates the employee (documented side effect)

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
