# F-059 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source absence + RUN-003 runtime): the theme handler/service persist config JSON without colour or value validation (no validate call in the theme service; url(javascript:) was accepted and served)


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/admin_theme_handler_v2.go`:293.
- Propagation: affected modules Theme builder; Storefront; public theme API; related findings none recorded; observed consequence: 201, activate 200, public endpoint returns the raw values; theme name with HTML also stored; exploitability in the storefront CSS pipeline not tested (PROBABLE CSS injection) | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): invalid colors and url
- Cause versus symptom: the symptom is what users or reports show (201, activate 200, public endpoint returns the raw values; theme name with HTML also stored; exploitability in the storefront CSS pipeline n); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-THM-001; scripts: tests/runtime/ph08_theme_driver.mjs, tests/runtime/ph08_theme2.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/admin_theme_handler_v2.go | 293 | backend/internal/handler/http/admin_theme_handler_v2.go | exact | yes |

