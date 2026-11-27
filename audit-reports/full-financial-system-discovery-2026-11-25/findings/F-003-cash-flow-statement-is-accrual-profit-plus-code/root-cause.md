# F-003 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + RUN-003 decomposition): backend/internal/usecase/financial_report_usecase.go:366-399 GenerateCashFlow is an indirect-method fragment: operating = P&L (revenue/expense credit-debit), investing = 12xx debit-credit (sign inverted), financing = equity or 22xx credit-debit; never reads cash/bank accounts and omits working-capital movements (see AUD-CF-TRACE.txt)


## Origin and propagation
- Origin (first verified source location): `backend/internal/usecase/financial_report_usecase.go`:366-399.
- Propagation: affected modules Financial reports; related findings none recorded; observed consequence: Operating = net profit (-2000); net cash flow 98000 vs GL cash+bank 93000 | RUN-003 (AUD-MET-FIN): cash-flow net 82,710.01 vs GL cash+bank 115,010.01, residual -32,300.00 (STILL REPRODUCIBLE).  | RUN-003 earlier run: cash-flow 82,710.01 vs 
- Cause versus symptom: the symptom is what users or reports show (Operating = net profit (-2000); net cash flow 98000 vs GL cash+bank 93000 | RUN-003 (AUD-MET-FIN): cash-flow net 82,710.01 vs GL cash+bank 1); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-MET-FIN, EX-RUN003-AUD-CF-TRACE, EX-RUN003-AUD-ACC-004, EX-RUN003-AUD-MONTH-001, EX-RUN004-MET-CF002; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/financial_report_usecase.go | 366-399 | backend/internal/usecase/financial_report_usecase.go | exact | yes |

