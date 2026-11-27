# F-025 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): handlers pass database errors to clients, e.g. backend/cmd/api/main.go:1302 (fmt.Sprintf("Failed to create user: %v", err)) and respondError(..., err.Error()) at admin_theme_handler_v2.go:80,120,154,361; 128 response sites in internal/handler/http use err.Error()


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/cmd/api/main.go | 1302 | backend/cmd/api/main.go | exact | yes |
| backend/internal/handler/http/admin_theme_handler_v2.go | 80 | backend/internal/handler/http/admin_theme_handler_v2.go | exact | yes |

