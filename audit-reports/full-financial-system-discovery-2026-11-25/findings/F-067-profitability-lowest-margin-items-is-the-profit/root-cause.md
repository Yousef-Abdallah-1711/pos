# F-067 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED: report_handler.go:2432 ORDER BY gross_profit DESC, then :2470-2476 builds the lowest list by iterating the same slice backwards; no margin sort exists


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/report_handler.go | 2432 | backend/internal/handler/http/report_handler.go | exact | yes |
| backend/internal/handler/http/report_handler.go | 2470 | backend/internal/handler/http/report_handler.go | exact | yes |
| backend/internal/handler/http/report_handler.go | 212 | backend/internal/handler/http/report_handler.go | exact | yes |

