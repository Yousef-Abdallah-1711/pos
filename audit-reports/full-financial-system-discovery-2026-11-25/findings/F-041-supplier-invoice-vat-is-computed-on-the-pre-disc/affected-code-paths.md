# F-041 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/migrations/142_create_invoices_supplier.sql | 72-81 | exact |
| backend/migrations/141_create_invoices_customer.sql | 76-85 | exact |
| backend/migrations/142_create_invoices_supplier.sql | 75-79 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/reports | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
See 08-api-inventory.csv for the endpoints in the affected modules; no contract change is proposed beyond the remediation spec.

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
