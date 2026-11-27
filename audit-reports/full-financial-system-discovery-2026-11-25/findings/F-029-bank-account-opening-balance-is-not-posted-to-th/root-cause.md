# F-029 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/bank_account_repo.go:101-122 inserts opening_balance/current_balance into the bank ledger with no journal entry (the file only reads journal ids in queries at :224-270); the bank handler has no posting code; fresh account in RUN-003 left GL 1112 = 0


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/bank_account_repo.go | 101-122 | backend/internal/repository/bank_account_repo.go | exact | yes |
| backend/internal/handler/http/bank_account_handler.go | 97 | backend/internal/handler/http/bank_account_handler.go | exact | yes |

