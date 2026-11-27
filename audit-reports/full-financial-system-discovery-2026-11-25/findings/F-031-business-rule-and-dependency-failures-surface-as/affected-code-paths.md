# F-031 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/cmd/api/main.go | 1302 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
- status codes change for the listed cases (documented in 08-api-inventory.csv)

## Proposed (future) files
- backend/internal/handler/http/errors.go (shared with F-025)
