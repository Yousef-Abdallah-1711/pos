# F-039 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source + RUN-003 runtime): CreateJournalEntryRequest.Attachments and JournalEntry.Attachments are json.RawMessage passthroughs with no size or schema limit (domain/accounting.go:95,139); a 1 MB attachment was stored (HC-05)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/domain/accounting.go | 95 | backend/internal/domain/accounting.go | exact | yes |
| backend/internal/domain/accounting.go | 139 | backend/internal/domain/accounting.go | exact | yes |

