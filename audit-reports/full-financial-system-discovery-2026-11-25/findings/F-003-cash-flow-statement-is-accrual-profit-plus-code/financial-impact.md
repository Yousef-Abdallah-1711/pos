# F-003 — Financial impact

- **Quantified impact (register):** 5000 mismatch (this dataset)
- **Accounts / statements named in the evidence:** 1110, 1111, 2120
- **Impact classes:** accounting, report

## What is wrong in the books or reports
Operating = net profit (-2000); net cash flow 98000 vs GL cash+bank 93000 | RUN-003 (AUD-MET-FIN): cash-flow net 82,710.01 vs GL cash+bank 115,010.01, residual -32,300.00 (STILL REPRODUCIBLE).  | RUN-003 earlier run: cash-flow 82,710.01 vs GL cash+bank 115,010.01 (residual 32,300.00). Re-measured after more activity: report 172,710.01 vs GL 1110+1111 = 280,010.01, residual 107,300.00. ROOT CAUSE CONFIRMED (AUD-CF-TRACE): usecase/financial_report_usecase.go:366-399 GenerateCashFlow is an indirect-method fragment that never reads cash/bank accounts and omits every working-capital movement (AR -6,520; inventory +96,000; input VAT -2,800; card clearing -11,000; AP +800; VAT payable +30,820 = +107,300); investing = '12%' debit-credit has inverted sign for asset purchases; financing excludes 21xx/2120 liabilities. 

## Which reports and decisions are affected
Financial reports that read the affected data: Financial reports

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.
Measured residuals: 32,300.00 earlier, 107,300.00 on the busy tenant, 5,190.00 on the clean tenant; all equal the omitted working-capital terms listed in evidence AUD-CF-TRACE.txt and AUD-MONTH.json.
