# F-041 — Financial impact

- **Quantified impact (register):** 140 per 10,000 gross at 10% discount and 14% VAT (VAT x discount)
- **Accounts / statements named in the evidence:** 1150, 2101, 5110
- **Impact classes:** accounting, report, validation

## What is wrong in the books or reports
API tax_amount 1,400, total 10,400; JE posted Dr 5110 9,000 / Dr 1150 1,400 / Cr 2101 10,400: input VAT +140 and AP +140 overstated | RUN-003 supplier flow (API-DIAGNOSTIC): invoice 10x1000 disc 10% VAT 14% => tax 1,400 total 10,400 vs oracle 1,260/10,260; JE Dr 5110 9,000 Dr 1150 1,400 Cr 2101 10,400 (CONFIRMED trigger 142:75-79). 

## Which reports and decisions are affected
Financial reports that read the affected data: Supplier invoices; Accounting (Input VAT 1150, AP 2101); VAT reports

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

