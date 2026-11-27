# F-073 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): AddBankReconciliationLine verifies only that the session exists for the tenant, then inserts journal_entry_id as supplied (accounting_advanced_repo.go:113-166); the table has a plain foreign key and the table has no row-level security (F-019)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 113 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/handler/http/accounting_advanced_handler.go | 58 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |
| accounting_advanced_repo.go | 113-166 | backend/internal/repository/accounting_advanced_repo.go | prefix | yes |

