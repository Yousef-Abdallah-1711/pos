# F-061 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/usecase/auto_posting_engine.go:535-537 rejects any posting with fewer than two lines; PostSalesEntry builds no lines for a zero total instead of returning nil


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/auto_posting_engine.go | 50-142 | backend/internal/usecase/auto_posting_engine.go | exact | yes |
| backend/internal/usecase/auto_posting_engine.go | 535-537 | backend/internal/usecase/auto_posting_engine.go | exact | yes |

