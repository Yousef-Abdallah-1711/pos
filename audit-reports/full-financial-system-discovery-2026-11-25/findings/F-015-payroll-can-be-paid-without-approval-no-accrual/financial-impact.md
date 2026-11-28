# F-015 — Financial impact

- **Quantified impact (register):** payroll payable 2000 absent
- **Accounts / statements named in the evidence:** 1111, 5210
- **Impact classes:** accounting, capability

## What is wrong in the books or reports
Paid with is_approved=false; journal only at payment as direct expense | RUN-003 SQL re-verification (AUD-EVD-F015): salaries id 1 and 2 are status=paid with is_approved=false; journals JE-2026-00007/00049 are Dr 5210 5,000 / Cr 1111 5,000 posted at payment only; no accrual journal exists (STILL REPRODUCIBLE; UI re-run pending T-P10-003). 

## Which reports and decisions are affected
Financial reports that read the affected data: HR payroll; Accounting

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

