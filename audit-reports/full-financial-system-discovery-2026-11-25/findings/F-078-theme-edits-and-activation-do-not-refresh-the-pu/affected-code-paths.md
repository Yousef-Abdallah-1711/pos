# F-078 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/theme_public_cache.go | 45 | exact |
| backend/internal/handler/http/public_homepage_handler.go | 117 | exact |
| frontend/apps/restaurant-website/app/[locale]/page.tsx | 75 | exact |
| backend/internal/handler/http/public_homepage_handler.go | 117-144 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |
| frontend/apps/restaurant-website/app/[locale] | exists |

## Database objects involved
`restaurants`

## API contracts involved
- unchanged

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
