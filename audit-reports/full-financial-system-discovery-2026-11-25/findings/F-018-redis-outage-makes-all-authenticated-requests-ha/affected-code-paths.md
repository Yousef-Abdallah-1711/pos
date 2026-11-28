# F-018 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/service/cache.go | 231-246 | exact |
| backend/cmd/api/main.go | 2003 | exact |
| backend/internal/middleware/auth.go | 60 | exact |
| backend/internal/service/auth_session.go | 166 | exact |
| backend/cmd/api/main.go | 2004 | exact |
| backend/internal/service/auth_session.go | 166-171 | exact |
| backend/internal/middleware/cache.go | 248 | basename-ambiguous(2) |
| backend/internal/middleware/cache.go | 232-237 | basename-ambiguous(2) |
| backend/internal/middleware/auth.go | 61-67 | basename |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |

## Database objects involved
`orders`, `products`

## API contracts involved
- /health response gains a redis component (additive)

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
