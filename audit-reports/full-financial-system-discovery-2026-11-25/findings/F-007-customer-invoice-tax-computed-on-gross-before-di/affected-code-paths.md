# F-007 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/migrations/141_create_invoices_customer.sql | 80-82 | prefix |
| backend/migrations/141_create_invoices_customer.sql | 79-83 | exact |
| backend/migrations/142_create_invoices_supplier.sql | 75-79 | basename |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices | exists |

## Database objects involved
None identified from the finding text.

## API contracts involved
See 08-api-inventory.csv for the endpoints in the affected modules; no contract change is proposed beyond the remediation spec.

## Proposed (future) files
- backend/migrations/<next>_fix_invoice_line_tax_discount_order.sql (new)
