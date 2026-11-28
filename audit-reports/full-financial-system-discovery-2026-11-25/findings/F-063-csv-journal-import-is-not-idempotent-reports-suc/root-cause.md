# F-063 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): backend/internal/handler/http/accounting_advanced_handler.go:1067-1140 ImportJournalEntriesCSV groups rows by entry_key/entry_number, creates each journal sequentially and always answers success:true with counters; no unique key, no job/async model, no per-request deadline handling


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/accounting_advanced_handler.go`:1067-1140.
- Propagation: affected modules Accounting CSV import; journal entries; related findings F-047, F-052; observed consequence: identical re-upload creates a second journal (IMP-03); unbalanced, unknown account, bad date and wrong header all return HTTP 200 success:true with failed_entries 1 and created_entries 0; a 6,000-journal import ran ~122 s, the client socket
- Cause versus symptom: the symptom is what users or reports show (identical re-upload creates a second journal (IMP-03); unbalanced, unknown account, bad date and wrong header all return HTTP 200 success:tr); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-IMP-001, EX-RUN004-BENCH-50K; scripts: tests/runtime/ph11_import.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/accounting_advanced_handler.go | 1067-1140 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |

