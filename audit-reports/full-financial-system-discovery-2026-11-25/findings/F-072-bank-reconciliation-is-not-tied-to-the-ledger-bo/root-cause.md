# F-072 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): refreshReconciliationSummaryTx sums bank_reconciliation_lines.book_amount supplied by the client (accounting_advanced_repo.go:2126-2149); AddBankReconciliationLine inserts without checking that the journal touches the bank GL account or that amounts equal the journal line (:113-166); CompleteBankReconciliationSession writes last_reconciled_* from the summed client amounts and sets completed_at regardless of status (:169-230)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 113 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/repository/accounting_advanced_repo.go | 169 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/repository/accounting_advanced_repo.go | 2126 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/handler/http/accounting_advanced_handler.go | 58 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |
| accounting_advanced_repo.go | 2126-2149 | backend/internal/repository/accounting_advanced_repo.go | prefix | yes |

