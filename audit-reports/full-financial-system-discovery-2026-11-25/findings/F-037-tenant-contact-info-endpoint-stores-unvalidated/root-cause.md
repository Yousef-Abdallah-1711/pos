# F-037 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED: inline handler cmd/api/main.go:845-873 decodes body and runs UPDATE tenants SET email=$1, phone=$2 with no validation or length limit


## Origin and propagation
- Origin (first verified source location): `backend/cmd/api/main.go`:845-873.
- Propagation: affected modules Settings / tenant profile; related findings none recorded; observed consequence: PUT /tenant/contact-info 200; tenants.email/phone persisted as entered; empty values also accepted; UI shows no validation for phone and only native browser validation for email (bypassed by API)
- Cause versus symptom: the symptom is what users or reports show (PUT /tenant/contact-info 200; tenants.email/phone persisted as entered; empty values also accepted; UI shows no validation for phone and onl); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-SET-001, EX-RUN003-AUD-SET-002, EX-RUN003-AUD-SET-003; scripts: tests/runtime/ph05_settings_forms.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/cmd/api/main.go | 845-873 | backend/cmd/api/main.go | exact | yes |

