# F-012 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): three independent queries define revenue differently: Sales = SUM(o.total_amount) tax-inclusive (report_handler.go:2177), Profitability = SUM(oi.total_price) net of tax for non-cancelled orders (report_handler.go:2367), P&L = GL revenue-type accounts (financial_report_usecase.go:208); no shared definition or label; numeric bridge in evidence/network/RUN-003/PH-12/AUD-REV-TRACE.txt


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/report_handler.go`:2177.
- Propagation: affected modules Reports; Dashboard; related findings none recorded; observed consequence: Sales 121000; Profitability 110000; KPI 22000; GL revenue 20000; KPI net profit 6000 vs P&L 4000; KPI cash position 0 vs GL cash
- Cause versus symptom: the symptom is what users or reports show (Sales 121000; Profitability 110000; KPI 22000; GL revenue 20000; KPI net profit 6000 vs P&L 4000; KPI cash position 0 vs GL cash); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/report_handler.go | 2177 | backend/internal/handler/http/report_handler.go | exact | yes |
| backend/internal/handler/http/report_handler.go | 2367 | backend/internal/handler/http/report_handler.go | exact | yes |
| backend/internal/usecase/financial_report_usecase.go | 208 | backend/internal/usecase/financial_report_usecase.go | exact | yes |

