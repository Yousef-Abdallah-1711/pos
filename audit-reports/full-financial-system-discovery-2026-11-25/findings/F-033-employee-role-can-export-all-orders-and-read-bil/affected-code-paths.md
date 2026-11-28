# F-033 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/auth_repo.go | 380 | exact |
| backend/cmd/api/main.go | 877 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr | exists |
| frontend/apps/super-admin | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`orders`, `attendance`, `leaves`, `employees`, `roles`

## API contracts involved
- several GET endpoints start returning 403 for the Employee role (breaking, coordinated with the dashboard)

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
