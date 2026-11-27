# F-047 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): repository/accounting_advanced_repo.go:396-420 RunRecurringJournalTemplate checks only template existence/is_active; it does not compare runDate with next_run_date/last_run_date and takes no row lock


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/accounting_advanced_repo.go`:396-470.
- Propagation: affected modules Accounting - recurring journals; related findings none recorded; observed consequence: Run 1 -> JE id 98 (rent +700); run 2 -> JE id 99 (rent +700 again): rent expense 1,400 for one period; next_run_date already advanced to 2026-12-26 | RUN-003: recurring run twice -> journals 98 and 99, rent +700 each (CONFIRMED accounting_a
- Cause versus symptom: the symptom is what users or reports show (Run 1 -> JE id 98 (rent +700); run 2 -> JE id 99 (rent +700 again): rent expense 1,400 for one period; next_run_date already advanced to 202); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-ACC-009; scripts: tests/runtime/ph11_accounting_flow2.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 396-470 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| repository/accounting_advanced_repo.go | 396-420 | backend/internal/repository/accounting_advanced_repo.go | prefix | yes |

