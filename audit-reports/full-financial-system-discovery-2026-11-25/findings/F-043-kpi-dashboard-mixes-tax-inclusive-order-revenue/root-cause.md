# F-043 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/accounting_advanced_repo.go:838-880 revenue=SUM(orders.total_amount paid, created_at), expenses=GL expense lines, cash_position=SUM(bank_accounts_ledger.current_balance)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 838-880 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |

