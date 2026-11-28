# F-010 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/tax_repo.go:260-291 GetVATSummary sums tax_amount of invoices_customer / invoices_supplier only (status <> void, drafts included); order (POS/online) VAT never enters the report


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/tax_repo.go`:260-291.
- Propagation: affected modules VAT report; related findings none recorded; observed consequence: output_vat 0 vs GL 1000 | RUN-003 (AUD-EVD-F010, tenant 4): VAT report API output_vat=700 vs SQL sum(orders.tax_amount, paid)=14,000 (+700 customer-invoice VAT) and GL VAT Payable (acct 87) credit balance 14,700; residual 14,000 = POS order
- Cause versus symptom: the symptom is what users or reports show (output_vat 0 vs GL 1000 | RUN-003 (AUD-EVD-F010, tenant 4): VAT report API output_vat=700 vs SQL sum(orders.tax_amount, paid)=14,000 (+700 c); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: tests/runtime/ph20_evidence_gaps.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/tax_repo.go | 260-291 | backend/internal/repository/tax_repo.go | exact | yes |

