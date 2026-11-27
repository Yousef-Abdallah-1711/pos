# F-008 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/invoice_customer_repo.go:452-560 createInvoiceTx calls syncJournalEntryTx (:539) which builds (buildJournalEntry :561-636) and posts the AR/revenue/VAT journal for every invoice at creation regardless of status; draft invoices therefore hit AR, revenue and VAT


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/invoice_customer_repo.go`:452-636.
- Propagation: affected modules Customer invoices; related findings none recorded; observed consequence: Journal JE-14 posted for draft invoice (journal_entry_id returned) | RUN-003 SQL (AUD-EVD-F008) PARTIAL: invoices_customer id 2 (INV-2026-00001; now status paid) has journal JE-2026-00014 posted 2026-11-25 22:49:27 Dr 1120 5,200 / Cr 4100 4
- Cause versus symptom: the symptom is what users or reports show (Journal JE-14 posted for draft invoice (journal_entry_id returned) | RUN-003 SQL (AUD-EVD-F008) PARTIAL: invoices_customer id 2 (INV-2026-00); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-CINV-001, EX-RUN005-INVOICE-UI; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/invoice_customer_repo.go | 452-636 | backend/internal/repository/invoice_customer_repo.go | exact | yes |
| backend/internal/repository/invoice_customer_repo.go | 452-560 | backend/internal/repository/invoice_customer_repo.go | exact | yes |

