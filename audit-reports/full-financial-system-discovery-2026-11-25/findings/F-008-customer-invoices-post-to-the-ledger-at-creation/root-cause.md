# F-008 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/invoice_customer_repo.go:452-560 createInvoiceTx calls syncJournalEntryTx (:539) which builds (buildJournalEntry :561-636) and posts the AR/revenue/VAT journal for every invoice at creation regardless of status; draft invoices therefore hit AR, revenue and VAT


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/invoice_customer_repo.go | 452-636 | backend/internal/repository/invoice_customer_repo.go | exact | yes |
| backend/internal/repository/invoice_customer_repo.go | 452-560 | backend/internal/repository/invoice_customer_repo.go | exact | yes |

