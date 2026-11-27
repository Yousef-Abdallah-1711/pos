# F-037 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED: inline handler cmd/api/main.go:845-873 decodes body and runs UPDATE tenants SET email=$1, phone=$2 with no validation or length limit


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/cmd/api/main.go | 845-873 | backend/cmd/api/main.go | exact | yes |

