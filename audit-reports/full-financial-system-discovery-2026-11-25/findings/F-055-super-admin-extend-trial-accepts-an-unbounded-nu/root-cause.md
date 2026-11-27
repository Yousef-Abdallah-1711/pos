# F-055 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): ExtendTenantTrial validates only days > 0 (subscription_handler.go:128-138); there is no upper bound


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/subscription_handler.go | 137 | backend/internal/handler/http/subscription_handler.go | exact | yes |
| subscription_handler.go | 128-138 | backend/internal/handler/http/subscription_handler.go | prefix | yes |

