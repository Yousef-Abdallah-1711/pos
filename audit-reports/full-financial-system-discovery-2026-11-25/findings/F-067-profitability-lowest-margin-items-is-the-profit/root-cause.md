# F-067 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED: report_handler.go:2432 ORDER BY gross_profit DESC, then :2470-2476 builds the lowest list by iterating the same slice backwards; no margin sort exists


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/report_handler.go`:2432.
- Propagation: affected modules Reports > Profitability (API, UI, CSV/PDF); related findings F-007, F-040; observed consequence: Order is MD probe 0%, MD probe 100%, Probe A 60%, Probe B 60%, Burger 60%: the list is the gross-PROFIT descending list iterated backwards
- Cause versus symptom: the symptom is what users or reports show (Order is MD probe 0%, MD probe 100%, Probe A 60%, Probe B 60%, Burger 60%: the list is the gross-PROFIT descending list iterated backwards); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-MET-PROF007; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/report_handler.go | 2432 | backend/internal/handler/http/report_handler.go | exact | yes |
| backend/internal/handler/http/report_handler.go | 2470 | backend/internal/handler/http/report_handler.go | exact | yes |
| backend/internal/handler/http/report_handler.go | 212 | backend/internal/handler/http/report_handler.go | exact | yes |

