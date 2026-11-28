# F-015 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/salary_repo.go | 344-380 | exact |
| backend/internal/usecase/auto_posting_engine.go | 261-345 | exact |
| backend/internal/usecase/auto_posting_engine.go | - | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr | exists |

## Database objects involved
`salaries`

## API contracts involved
- PUT /api/v1/hr/salaries/{id}/pay: 409 when not approved; accepts amount for partial payment (additive)

## Proposed (future) files
- backend/migrations/<next>_salary_payments_and_accrual.sql (new)
