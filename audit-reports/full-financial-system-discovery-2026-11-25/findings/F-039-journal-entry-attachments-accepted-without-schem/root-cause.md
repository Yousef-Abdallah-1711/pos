# F-039 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source + RUN-003 runtime): CreateJournalEntryRequest.Attachments and JournalEntry.Attachments are json.RawMessage passthroughs with no size or schema limit (domain/accounting.go:95,139); a 1 MB attachment was stored (HC-05)


## Origin and propagation
- Origin (first verified source location): `backend/internal/domain/accounting.go`:95.
- Propagation: affected modules Accounting / journal entries; related findings none recorded; observed consequence: All five variants returned 201 and were stored verbatim (journals 62-66, up to 1,048,605 bytes); historical issue HC-05 STILL REPRODUCIBLE
- Cause versus symptom: the symptom is what users or reports show (All five variants returned 201 and were stored verbatim (journals 62-66, up to 1,048,605 bytes); historical issue HC-05 STILL REPRODUCIBLE); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-REG-HC05; scripts: tests/runtime/ph19_regressions.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/domain/accounting.go | 95 | backend/internal/domain/accounting.go | exact | yes |
| backend/internal/domain/accounting.go | 139 | backend/internal/domain/accounting.go | exact | yes |

