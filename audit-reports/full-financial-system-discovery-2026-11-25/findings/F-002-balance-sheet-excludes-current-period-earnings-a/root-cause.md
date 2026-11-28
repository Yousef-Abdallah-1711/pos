# F-002 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source, RUN-003 read): backend/internal/usecase/financial_report_usecase.go:285-306 balance-sheet query selects only account_type IN (asset, liability, equity); revenue/expense accounts never contribute, so current-period earnings are not added to equity (BS difference 82,700 measured earlier)


## Origin and propagation
- Origin (first verified source location): `backend/internal/usecase/financial_report_usecase.go`:285-306.
- Propagation: affected modules Financial reports; related findings none recorded; observed consequence: equation_difference -2000.00 (net loss); equity shows 100000 only | RUN-003 clean-tenant month (AUD-MONTH): balance sheet equation_balanced=false, equation_difference -27,710.00 = net loss exactly (REPRODUCED). 
- Cause versus symptom: the symptom is what users or reports show (equation_difference -2000.00 (net loss); equity shows 100000 only | RUN-003 clean-tenant month (AUD-MONTH): balance sheet equation_balanced=); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-MONTH-001; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/financial_report_usecase.go | 285-306 | backend/internal/usecase/financial_report_usecase.go | exact | yes |

