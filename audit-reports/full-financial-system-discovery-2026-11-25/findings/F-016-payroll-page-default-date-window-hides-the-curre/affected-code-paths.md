# F-016 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx | 57-58 | exact |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx | 91 | exact |
| backend/internal/repository/accounting_advanced_repo.go | 1236-1243 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr | exists |

## Database objects involved
`salaries`

## API contracts involved
- GET /api/v1/hr/salaries: no change except overlap semantics for date filters

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
