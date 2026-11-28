# F-013 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source) for the missing-period part: backend/internal/repository/journal_entry_repo.go:813-842 ensureFiscalPeriodOpenTx returns nil when no fiscal_periods row matches, so any date posts; only status=closed blocks; no API creates/closes periods (GET /accounting/fiscal-periods is 404)


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/journal_entry_repo.go`:813-842.
- Propagation: affected modules Accounting; related findings none recorded; observed consequence: Post blocked; draft create allowed; void of posted entry allowed (capital 100000 reversed) | RUN-003 (AUD-ACC-FLOW JE-13/JE-14): journals dated 2019-01-15 and 2030-06-01 post successfully and auto-create fiscal_period 2019-01/2030-06; GET /
- Cause versus symptom: the symptom is what users or reports show (Post blocked; draft create allowed; void of posted entry allowed (capital 100000 reversed) | RUN-003 (AUD-ACC-FLOW JE-13/JE-14): journals da); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-ACC-006; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/journal_entry_repo.go | 813-842 | backend/internal/repository/journal_entry_repo.go | exact | yes |

