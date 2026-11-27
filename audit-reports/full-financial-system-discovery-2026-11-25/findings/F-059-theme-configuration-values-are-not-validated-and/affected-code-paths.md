# F-059 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/admin_theme_handler_v2.go | 293 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/restaurant-website/app/[locale] | exists |

## Database objects involved
`restaurants`

## API contracts involved
- POST/PUT /api/v1/admin/themes: 400 with field errors

## Proposed (future) files
- backend/internal/service/theme_validator.go (new)
