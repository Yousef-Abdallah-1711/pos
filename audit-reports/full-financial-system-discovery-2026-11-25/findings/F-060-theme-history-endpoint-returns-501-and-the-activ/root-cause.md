# F-060 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): GetThemeHistory returns 501 Not Implemented (admin_theme_handler.go:68)


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/admin_theme_handler.go`:68.
- Propagation: affected modules Theme builder; related findings none recorded; observed consequence: history -> 501 Not Implemented; DELETE active theme -> 204 leaving the tenant without an active theme | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): history 501; active theme deleted with 204 (REPRODUCED). 
- Cause versus symptom: the symptom is what users or reports show (history -> 501 Not Implemented; DELETE active theme -> 204 leaving the tenant without an active theme | RUN-003 (AUD-THEME-DRIVER / AUD-THEM); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-THM-001, EX-RUN004-THEMEBUILDER-UI, EX-RUN005-DESTR-THEME; scripts: tests/runtime/ph08_theme_driver.mjs, tests/runtime/ph08_theme2.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/admin_theme_handler.go | 68 | backend/internal/handler/http/admin_theme_handler.go | exact | yes |

