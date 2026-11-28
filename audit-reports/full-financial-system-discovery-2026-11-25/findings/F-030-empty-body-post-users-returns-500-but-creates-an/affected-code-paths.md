# F-030 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/cmd/api/main.go | 1261 | exact |
| backend/cmd/api/main.go | 1300 | exact |
| backend/cmd/api/main.go | 1261-1335 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
No frontend surface is involved or none could be derived.

## Database objects involved
`users`

## API contracts involved
- POST /api/v1/users: 400 field errors (was 500)

## Proposed (future) files
- backend/internal/handler/http/user_admin_handler.go (new)
