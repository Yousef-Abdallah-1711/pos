# F-046 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/report_handler.go | 557-645 | exact |
| backend/internal/handler/http/report_handler.go | 2480-2609 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/products | exists |

## Database objects involved
`products`, `inventory_movements`

## API contracts involved
- response shape unchanged

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
