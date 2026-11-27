# F-041 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): trigger recalculate_invoice_supplier_line_total backend/migrations/142_create_invoices_supplier.sql:75-79 computes tax on quantity*unit_price before subtracting the discount (same formula as customer trigger, F-007)


## Origin and propagation
- Origin (first verified source location): `backend/migrations/142_create_invoices_supplier.sql`:72-81.
- Propagation: affected modules Supplier invoices; Accounting (Input VAT 1150, AP 2101); VAT reports; related findings F-007; observed consequence: API tax_amount 1,400, total 10,400; JE posted Dr 5110 9,000 / Dr 1150 1,400 / Cr 2101 10,400: input VAT +140 and AP +140 overstated | RUN-003 supplier flow (API-DIAGNOSTIC): invoice 10x1000 disc 10% VAT 14% => tax 1,400 total 10,400 vs orac
- Cause versus symptom: the symptom is what users or reports show (API tax_amount 1,400, total 10,400; JE posted Dr 5110 9,000 / Dr 1150 1,400 / Cr 2101 10,400: input VAT +140 and AP +140 overstated | RUN-00); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-SUP-001; scripts: tests/runtime/ph06_supplier_flow.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/migrations/142_create_invoices_supplier.sql | 72-81 | backend/migrations/142_create_invoices_supplier.sql | exact | yes |
| backend/migrations/141_create_invoices_customer.sql | 76-85 | backend/migrations/141_create_invoices_customer.sql | exact | yes |
| backend/migrations/142_create_invoices_supplier.sql | 75-79 | backend/migrations/142_create_invoices_supplier.sql | exact | yes |

