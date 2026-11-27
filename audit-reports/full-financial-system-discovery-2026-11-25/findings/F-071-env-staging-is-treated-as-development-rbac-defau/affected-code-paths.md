# F-071 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/config/runtime_env.go | 20 | exact |
| backend/internal/middleware/permission.go | 36 | exact |
| backend/internal/config/paymob_mode.go | 29 | exact |
| backend/cmd/api/health_routes.go | 52 | exact |
| backend/cmd/api/main.go | 705 | exact |
| backend/cmd/api/signing_secrets.go | 9 | exact |
| backend/internal/config/runtime_env.go | 20-23 | basename |
| backend/internal/middleware/permission.go | 36-42 | prefix |
| backend/cmd/api/main.go | 148-156 | basename-ambiguous(5) |
| backend/internal/handler/http/auth_handler.go | 223 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr | exists |

## Database objects involved
`salaries`, `users`

## API contracts involved
- /health/metrics requires a token outside development

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
