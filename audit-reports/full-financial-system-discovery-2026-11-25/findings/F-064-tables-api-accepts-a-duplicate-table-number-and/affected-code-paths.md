# F-064 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/table_handler.go | - | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx | exists |

## Database objects involved
`orders`

## API contracts involved
See 08-api-inventory.csv for the endpoints in the affected modules; no contract change is proposed beyond the remediation spec.

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
