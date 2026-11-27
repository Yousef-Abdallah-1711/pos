# F-003 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + RUN-003 decomposition): backend/internal/usecase/financial_report_usecase.go:366-399 GenerateCashFlow is an indirect-method fragment: operating = P&L (revenue/expense credit-debit), investing = 12xx debit-credit (sign inverted), financing = equity or 22xx credit-debit; never reads cash/bank accounts and omits working-capital movements (see AUD-CF-TRACE.txt)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/usecase/financial_report_usecase.go | 366-399 | backend/internal/usecase/financial_report_usecase.go | exact | yes |

