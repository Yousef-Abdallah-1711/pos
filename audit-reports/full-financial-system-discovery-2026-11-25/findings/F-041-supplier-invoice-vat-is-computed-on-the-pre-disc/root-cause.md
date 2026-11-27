# F-041 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): trigger recalculate_invoice_supplier_line_total backend/migrations/142_create_invoices_supplier.sql:75-79 computes tax on quantity*unit_price before subtracting the discount (same formula as customer trigger, F-007)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/migrations/142_create_invoices_supplier.sql | 72-81 | backend/migrations/142_create_invoices_supplier.sql | exact | yes |
| backend/migrations/141_create_invoices_customer.sql | 76-85 | backend/migrations/141_create_invoices_customer.sql | exact | yes |
| backend/migrations/142_create_invoices_supplier.sql | 75-79 | backend/migrations/142_create_invoices_supplier.sql | exact | yes |

