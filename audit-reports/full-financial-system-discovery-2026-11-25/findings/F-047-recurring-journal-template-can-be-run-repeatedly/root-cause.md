# F-047 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): repository/accounting_advanced_repo.go:396-420 RunRecurringJournalTemplate checks only template existence/is_active; it does not compare runDate with next_run_date/last_run_date and takes no row lock


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 396-470 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| repository/accounting_advanced_repo.go | 396-420 | backend/internal/repository/accounting_advanced_repo.go | prefix | yes |

