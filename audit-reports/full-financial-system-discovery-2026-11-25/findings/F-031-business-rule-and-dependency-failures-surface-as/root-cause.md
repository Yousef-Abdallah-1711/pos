# F-031 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): error mapping is done by passing err text to the client (see F-025 sites) and returning 500 for constraint violations (e.g. main.go:1302; hr handlers rely on DB check constraints)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/cmd/api/main.go | 1302 | backend/cmd/api/main.go | exact | yes |

