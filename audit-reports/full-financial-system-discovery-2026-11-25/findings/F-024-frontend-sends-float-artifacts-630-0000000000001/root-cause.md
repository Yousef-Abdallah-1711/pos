# F-024 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): InvoiceForm computes discountAmount, taxAmount and lineTotal with raw floating-point multiplication and sends them unrounded (frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx:620-650)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | 623 | frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | exact | yes |
| frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | 637 | frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | exact | yes |
| frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | 620-650 | frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | exact | yes |

