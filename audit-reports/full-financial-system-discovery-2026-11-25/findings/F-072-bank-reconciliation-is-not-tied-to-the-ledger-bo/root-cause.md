# F-072 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): refreshReconciliationSummaryTx sums bank_reconciliation_lines.book_amount supplied by the client (accounting_advanced_repo.go:2126-2149); AddBankReconciliationLine inserts without checking that the journal touches the bank GL account or that amounts equal the journal line (:113-166); CompleteBankReconciliationSession writes last_reconciled_* from the summed client amounts and sets completed_at regardless of status (:169-230)


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/accounting_advanced_repo.go`:113.
- Propagation: affected modules Accounting > Bank reconciliation, bank accounts; related findings F-029, F-043; observed consequence: Session status completed, difference 0.00, bank_accounts_ledger.last_reconciled_balance = 777,777.77 and last_reconciled_date set, while the ledger balance is 12,000.00 and the journal has 0 bank-account lines. A session with difference 123
- Cause versus symptom: the symptom is what users or reports show (Session status completed, difference 0.00, bank_accounts_ledger.last_reconciled_balance = 777,777.77 and last_reconciled_date set, while the); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-BANKRECON; scripts: tests/runtime/ph11_r4_bank_recon.mjs, tests/runtime/ph11_r4_bank_recon2.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 113 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/repository/accounting_advanced_repo.go | 169 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/repository/accounting_advanced_repo.go | 2126 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/handler/http/accounting_advanced_handler.go | 58 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |
| accounting_advanced_repo.go | 2126-2149 | backend/internal/repository/accounting_advanced_repo.go | prefix | yes |

