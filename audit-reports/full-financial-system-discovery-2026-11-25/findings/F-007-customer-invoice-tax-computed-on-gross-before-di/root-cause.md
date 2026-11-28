# F-007 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source, RUN-003): DB trigger taxes gross | RUN-003 SOURCE: DB trigger recalculate_invoice_customer_line_total (backend/migrations/141_create_invoices_customer.sql:79-83) computes tax_amount = ROUND(qty*unit_price*tax_rate/100,2) on the PRE-discount amount, then line_total = gross - discount + tax; the supplier twin (142_create_invoices_supplier.sql:75-79) is identical (see F-041).


## Origin and propagation
- Origin (first verified source location): `backend/migrations/141_create_invoices_customer.sql`:80-82.
- Propagation: affected modules Customer invoices; related findings F-041; observed consequence: Stored tax 700 total 5200 | RUN-003 (AUD-CINV-001): UI Tax 1,008 / Total 8,208 (matches discount-first oracle) vs stored tax 1,120 / total 8,320 (VAT on pre-discount subtotal): +112 tax, +112 AR (STILL REPRODUCIBLE).  RUN-005 real-UI reconf
- Cause versus symptom: the symptom is what users or reports show (Stored tax 700 total 5200 | RUN-003 (AUD-CINV-001): UI Tax 1,008 / Total 8,208 (matches discount-first oracle) vs stored tax 1,120 / total 8); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-CINV-001, EX-RUN005-INVOICE-UI; scripts: tests/runtime/ph09_custinv_draft_headed.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| migrations/141_create_invoices_customer.sql | 80-82 | backend/migrations/141_create_invoices_customer.sql | prefix | yes |
| backend/migrations/141_create_invoices_customer.sql | 79-83 | backend/migrations/141_create_invoices_customer.sql | exact | yes |
| 142_create_invoices_supplier.sql | 75-79 | backend/migrations/142_create_invoices_supplier.sql | basename | yes |

