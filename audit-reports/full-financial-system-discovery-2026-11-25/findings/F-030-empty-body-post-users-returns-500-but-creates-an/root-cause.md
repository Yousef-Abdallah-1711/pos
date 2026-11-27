# F-030 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): POST /users handler decodes name/email/phone/password with no validation, hashes the (possibly empty) password and calls CreateUserInTenant, then assigns the EMPLOYEE role with the error discarded (_ =) in separate non-atomic steps (backend/cmd/api/main.go:1261-1335)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/cmd/api/main.go | 1261 | backend/cmd/api/main.go | exact | yes |
| backend/cmd/api/main.go | 1300 | backend/cmd/api/main.go | exact | yes |
| backend/cmd/api/main.go | 1261-1335 | backend/cmd/api/main.go | exact | yes |

