# F-059 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source absence + RUN-003 runtime): the theme handler/service persist config JSON without colour or value validation (no validate call in the theme service; url(javascript:) was accepted and served)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/admin_theme_handler_v2.go | 293 | backend/internal/handler/http/admin_theme_handler_v2.go | exact | yes |

