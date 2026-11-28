# F-010 — Financial impact

- **Quantified impact (register):** 1000 here
- **Accounts / statements named in the evidence:** 2120
- **Impact classes:** accounting, report

## What is wrong in the books or reports
output_vat 0 vs GL 1000 | RUN-003 (AUD-EVD-F010, tenant 4): VAT report API output_vat=700 vs SQL sum(orders.tax_amount, paid)=14,000 (+700 customer-invoice VAT) and GL VAT Payable (acct 87) credit balance 14,700; residual 14,000 = POS order VAT missing from report (STILL REPRODUCIBLE). 

## Which reports and decisions are affected
Financial reports that read the affected data: VAT report

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

