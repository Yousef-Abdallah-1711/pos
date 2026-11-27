# F-025 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/cmd/api/main.go | 1302 | exact |
| backend/internal/handler/http/admin_theme_handler_v2.go | 80 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
No frontend surface is involved or none could be derived.

## Database objects involved
None identified from the finding text.

## API contracts involved
- error bodies gain {code, request_id} additively; message text becomes generic

## Proposed (future) files
- backend/internal/handler/http/errors.go (new)
- backend/internal/handler/http/errors_scan_test.go (new)
