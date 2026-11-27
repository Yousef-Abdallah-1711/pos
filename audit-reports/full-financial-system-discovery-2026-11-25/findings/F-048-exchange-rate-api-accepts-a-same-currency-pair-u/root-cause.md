# F-048 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): the handler only requires non-empty codes and a positive rate (accounting_advanced_handler.go:475-482); there is no base != quote check


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/accounting_advanced_handler.go`.
- Propagation: affected modules Accounting - exchange rates; related findings none recorded; observed consequence: 200 {id:2}: pair stored | RUN-003: USD/USD pair accepted (id 2). 
- Cause versus symptom: the symptom is what users or reports show (200 {id:2}: pair stored | RUN-003: USD/USD pair accepted (id 2). ); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-ACC-011; scripts: tests/runtime/ph11_accounting_flow2.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/accounting_advanced_handler.go | - | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |
| backend/internal/handler/http/accounting_advanced_handler.go | 481 | backend/internal/handler/http/accounting_advanced_handler.go | exact | yes |
| accounting_advanced_handler.go | 475-482 | backend/internal/handler/http/accounting_advanced_handler.go | prefix | yes |

