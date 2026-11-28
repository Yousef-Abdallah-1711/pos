# F-075 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + runtime): CreateAccount only trims the strings and passes them on (backend/internal/handler/http/account_handler.go:124-170: no regexp, length or format check); repository errors are returned to the client as 400 text (F-025)


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/account_handler.go`:124.
- Propagation: affected modules Accounting > Chart of accounts; related findings F-003, F-025, F-031, F-062; observed consequence: POST /accounting/accounts returns 201 and stores account id 389 with code <b>x</b>'; -- (the response HTML-escapes it as unicode); a 300-character value returns 400 with the raw database length error; no maxlength on the inputs
- Cause versus symptom: the symptom is what users or reports show (POST /accounting/accounts returns 201 and stores account id 389 with code <b>x</b>'; -- (the response HTML-escapes it as unicode); a 300-cha); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-FORMMATRIX; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/account_handler.go | 124 | backend/internal/handler/http/account_handler.go | exact | yes |
| backend/internal/handler/http/account_handler.go | 124-170 | backend/internal/handler/http/account_handler.go | exact | yes |

