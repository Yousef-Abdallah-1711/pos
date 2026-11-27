# F-024 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): InvoiceForm computes discountAmount, taxAmount and lineTotal with raw floating-point multiplication and sends them unrounded (frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx:620-650)


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx`:623.
- Propagation: affected modules Invoices; related findings none recorded; observed consequence: tax_amount 630.0000000000001 (server stores rounded) RUN-005 (real invoice form, 1 x 5,000, 10% discount, 14% VAT): the POST payload carries "tax_amount":630.0000000000001 (raw float artefact) while the UI shows Tax EGP 630.00.
- Cause versus symptom: the symptom is what users or reports show (tax_amount 630.0000000000001 (server stores rounded) RUN-005 (real invoice form, 1 x 5,000, 10% discount, 14% VAT): the POST payload carries); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN005-OBSERVED-REPRO, EX-RUN005-INVOICE-UI; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | 623 | frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | exact | yes |
| frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | 637 | frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | exact | yes |
| frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | 620-650 | frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | exact | yes |

