# F-007 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source, RUN-003): DB trigger taxes gross | RUN-003 SOURCE: DB trigger recalculate_invoice_customer_line_total (backend/migrations/141_create_invoices_customer.sql:79-83) computes tax_amount = ROUND(qty*unit_price*tax_rate/100,2) on the PRE-discount amount, then line_total = gross - discount + tax; the supplier twin (142_create_invoices_supplier.sql:75-79) is identical (see F-041).


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| migrations/141_create_invoices_customer.sql | 80-82 | backend/migrations/141_create_invoices_customer.sql | prefix | yes |
| backend/migrations/141_create_invoices_customer.sql | 79-83 | backend/migrations/141_create_invoices_customer.sql | exact | yes |
| 142_create_invoices_supplier.sql | 75-79 | backend/migrations/142_create_invoices_supplier.sql | basename | yes |

