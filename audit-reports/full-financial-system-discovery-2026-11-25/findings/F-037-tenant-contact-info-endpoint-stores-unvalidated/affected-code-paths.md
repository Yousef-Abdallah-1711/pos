# F-037 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/cmd/api/main.go | 845-873 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings | exists |

## Database objects involved
`tenants`

## API contracts involved
- PUT contact-info: 400 field errors (was 200)

## Proposed (future) files
- backend/internal/handler/http/tenant_contact_handler.go (new)
