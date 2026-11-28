# F-029 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/bank_account_repo.go:101-122 inserts opening_balance/current_balance into the bank ledger with no journal entry (the file only reads journal ids in queries at :224-270); the bank handler has no posting code; fresh account in RUN-003 left GL 1112 = 0


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/bank_account_repo.go`:101-122.
- Propagation: affected modules Bank accounts; Reports KPI; Accounting; related findings none recorded; observed consequence: No journal; account shows Current balance 50000 while its transaction list runs -10000; KPI cash_position 50000 vs GL 1111 -10000 | RUN-003 (AUD-ACC-FLOW BNK-01): new bank account opening_balance 12,000 linked to GL 1112 creates no journal 
- Cause versus symptom: the symptom is what users or reports show (No journal; account shows Current balance 50000 while its transaction list runs -10000; KPI cash_position 50000 vs GL 1111 -10000 | RUN-003 ); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-ACC-008; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/bank_account_repo.go | 101-122 | backend/internal/repository/bank_account_repo.go | exact | yes |
| backend/internal/handler/http/bank_account_handler.go | 97 | backend/internal/handler/http/bank_account_handler.go | exact | yes |

