# F-058 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/admin_driver_handler.go | 129 | exact |
| backend/internal/handler/http/admin_driver_handler.go | 146 | exact |
| backend/internal/handler/http/admin_driver_handler.go | 146-150 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
No frontend surface is involved or none could be derived.

## Database objects involved
`employees`

## API contracts involved
- POST /api/v1/admin/drivers: 400 for invalid phone

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
