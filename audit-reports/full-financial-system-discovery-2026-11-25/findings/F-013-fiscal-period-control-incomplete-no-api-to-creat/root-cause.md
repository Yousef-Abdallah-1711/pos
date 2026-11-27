# F-013 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source) for the missing-period part: backend/internal/repository/journal_entry_repo.go:813-842 ensureFiscalPeriodOpenTx returns nil when no fiscal_periods row matches, so any date posts; only status=closed blocks; no API creates/closes periods (GET /accounting/fiscal-periods is 404)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/journal_entry_repo.go | 813-842 | backend/internal/repository/journal_entry_repo.go | exact | yes |

