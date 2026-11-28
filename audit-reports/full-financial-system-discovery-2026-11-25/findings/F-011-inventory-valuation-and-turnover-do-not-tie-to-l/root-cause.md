# F-011 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/usecase/financial_report_usecase.go:488-526 GenerateInventoryValuation = SUM(quantity_in_stock x cost) for tracked products; it never reads GL 1130, and purchases/opening stock never post to 1130 (supplier invoice lines post to 5110, OBS-PUR-001)


## Origin and propagation
- Origin (first verified source location): `backend/internal/usecase/financial_report_usecase.go`:488-526.
- Propagation: affected modules Inventory; Products; Reports; related findings none recorded; observed consequence: 196000 (report) vs 16000 (GL) then 156000 vs 8000 | RUN-003 (AUD-EVD-F011): inventory valuation report 192,000 (48 x 4,000) = SQL sum(quantity_in_stock*cost); GL Inventory (acct 1130) balance -36,000; residual 228,000; opening stock has no 
- Cause versus symptom: the symptom is what users or reports show (196000 (report) vs 16000 (GL) then 156000 vs 8000 | RUN-003 (AUD-EVD-F011): inventory valuation report 192,000 (48 x 4,000) = SQL sum(quanti); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-OBS-PUR-001; scripts: tests/runtime/ph20_evidence_gaps.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/financial_report_usecase.go | 488-526 | backend/internal/usecase/financial_report_usecase.go | exact | yes |

