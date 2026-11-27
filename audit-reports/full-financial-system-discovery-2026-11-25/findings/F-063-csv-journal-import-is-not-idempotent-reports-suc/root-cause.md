# F-063 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): backend/internal/handler/http/accounting_advanced_handler.go:1067-1140 ImportJournalEntriesCSV groups rows by entry_key/entry_number, creates each journal sequentially and always answers success:true with counters; no unique key, no job/async model, no per-request deadline handling


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/accounting_advanced_handler.go | 1067-1140 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |

