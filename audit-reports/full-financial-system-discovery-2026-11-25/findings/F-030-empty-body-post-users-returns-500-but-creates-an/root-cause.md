# F-030 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): POST /users handler decodes name/email/phone/password with no validation, hashes the (possibly empty) password and calls CreateUserInTenant, then assigns the EMPLOYEE role with the error discarded (_ =) in separate non-atomic steps (backend/cmd/api/main.go:1261-1335)


## Origin and propagation
- Origin (first verified source location): `backend/cmd/api/main.go`:1261.
- Propagation: affected modules Users; Notifications; related findings none recorded; observed consequence: 500 returned yet user id 6 created (blank email and name; role assigned; 60-char hash of empty password); PUT /users/{id} and POST /notifications also 500
- Cause versus symptom: the symptom is what users or reports show (500 returned yet user id 6 created (blank email and name; role assigned; 60-char hash of empty password); PUT /users/{id} and POST /notifica); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/cmd/api/main.go | 1261 | backend/cmd/api/main.go | exact | yes |
| backend/cmd/api/main.go | 1300 | backend/cmd/api/main.go | exact | yes |
| backend/cmd/api/main.go | 1261-1335 | backend/cmd/api/main.go | exact | yes |

