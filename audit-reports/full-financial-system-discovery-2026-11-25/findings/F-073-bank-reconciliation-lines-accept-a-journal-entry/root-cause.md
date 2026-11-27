# F-073 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): AddBankReconciliationLine verifies only that the session exists for the tenant, then inserts journal_entry_id as supplied (accounting_advanced_repo.go:113-166); the table has a plain foreign key and the table has no row-level security (F-019)


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/accounting_advanced_repo.go`:113.
- Propagation: affected modules Accounting > Bank reconciliation; related findings F-019, F-025, F-072; observed consequence: HTTP 200; the line is stored in the tenant-4 session referencing tenant 2 journal 33; a nonexistent id returns 400 only because of the foreign key, with the raw constraint text in the error body
- Cause versus symptom: the symptom is what users or reports show (HTTP 200; the line is stored in the tenant-4 session referencing tenant 2 journal 33; a nonexistent id returns 400 only because of the forei); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-BANKRECON; scripts: tests/runtime/ph11_r4_bank_recon.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 113 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |
| backend/internal/handler/http/accounting_advanced_handler.go | 58 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |
| accounting_advanced_repo.go | 113-166 | backend/internal/repository/accounting_advanced_repo.go | prefix | yes |

