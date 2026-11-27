# F-075 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): CreateAccount only trims the strings and passes them on (backend/internal/handler/http/account_handler.go:124-170: no regexp, length or format check); repository errors are returned to the client as 400 text (F-025)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/account_handler.go | 124 | backend/internal/handler/http/account_handler.go | exact | yes |
| backend/internal/handler/http/account_handler.go | 124-170 | backend/internal/handler/http/account_handler.go | exact | yes |

