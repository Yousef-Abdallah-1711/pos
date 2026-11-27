# F-020 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/auth_repo.go | 380 | exact |
| backend/cmd/api/main.go | 877 | exact |
| backend/internal/domain/hr.go | 44 | exact |
| backend/internal/repository/auth_repo.go | - | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr | exists |

## Database objects involved
`orders`, `salaries`, `employees`, `roles`

## API contracts involved
- GET /api/v1/hr/employees: pay fields omitted unless SALARIES READ (breaking for clients that relied on them; coordinate with dashboard)

## Proposed (future) files
- backend/internal/handler/http/employee_dto.go (new)
