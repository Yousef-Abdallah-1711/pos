# F-060 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/admin_theme_handler.go | 68 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
No frontend surface is involved or none could be derived.

## Database objects involved
None identified from the finding text.

## API contracts involved
- GET /api/v1/admin/themes/{id}/history (implemented); POST .../restore (new)

## Proposed (future) files
- backend/migrations/<next>_theme_versions.sql (new)
