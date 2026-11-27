# F-055 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/subscription_handler.go | 137 | exact |
| backend/internal/handler/http/subscription_handler.go | 128-138 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/super-admin | exists |

## Database objects involved
`tenants`

## API contracts involved
- POST /api/v1/super-admin/tenants/{id}/extend-trial: reason required; days capped

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
