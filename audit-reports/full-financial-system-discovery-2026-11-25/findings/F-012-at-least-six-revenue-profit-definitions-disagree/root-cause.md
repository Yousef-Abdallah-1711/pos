# F-012 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): three independent queries define revenue differently: Sales = SUM(o.total_amount) tax-inclusive (report_handler.go:2177), Profitability = SUM(oi.total_price) net of tax for non-cancelled orders (report_handler.go:2367), P&L = GL revenue-type accounts (financial_report_usecase.go:208); no shared definition or label; numeric bridge in evidence/network/RUN-003/PH-12/AUD-REV-TRACE.txt


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/report_handler.go | 2177 | backend/internal/handler/http/report_handler.go | exact | yes |
| backend/internal/handler/http/report_handler.go | 2367 | backend/internal/handler/http/report_handler.go | exact | yes |
| backend/internal/usecase/financial_report_usecase.go | 208 | backend/internal/usecase/financial_report_usecase.go | exact | yes |

