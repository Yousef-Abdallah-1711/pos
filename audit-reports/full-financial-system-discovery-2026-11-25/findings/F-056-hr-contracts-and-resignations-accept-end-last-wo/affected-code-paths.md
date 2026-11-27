# F-056 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/contract_handler.go | - | exact |
| backend/internal/handler/http/resignation_handler.go | - | exact |
| backend/internal/handler/http/contract_handler.go | 165 | exact |
| backend/internal/repository/resignation_repo.go | 117 | exact |
| backend/internal/handler/http/contract_handler.go | 165-170 | prefix |
| backend/internal/repository/resignation_repo.go | 117-143 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr | exists |

## Database objects involved
`attendance`

## API contracts involved
- contracts/resignations create and update: 400 for inverted ranges

## Proposed (future) files
- backend/migrations/<next>_hr_date_range_checks.sql (new)
