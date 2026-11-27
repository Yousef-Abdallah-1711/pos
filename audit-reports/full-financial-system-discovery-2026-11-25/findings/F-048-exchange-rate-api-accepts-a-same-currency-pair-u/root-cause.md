# F-048 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): the handler only requires non-empty codes and a positive rate (accounting_advanced_handler.go:475-482); there is no base != quote check


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/accounting_advanced_handler.go | - | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |
| backend/internal/handler/http/accounting_advanced_handler.go | 481 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |
| accounting_advanced_handler.go | 475-482 | backend/internal/handler/http/accounting_advanced_handler.go | prefix | yes |

