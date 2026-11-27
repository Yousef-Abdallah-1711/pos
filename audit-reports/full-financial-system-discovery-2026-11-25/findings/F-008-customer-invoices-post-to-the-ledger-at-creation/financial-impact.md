# F-008 — Financial impact

- **Quantified impact (register):** 5200 AR+revenue recognised at draft
- **Accounts / statements named in the evidence:** 1120, 2120, 4100
- **Impact classes:** accounting

## What is wrong in the books or reports
Journal JE-14 posted for draft invoice (journal_entry_id returned) | RUN-003 SQL (AUD-EVD-F008) PARTIAL: invoices_customer id 2 (INV-2026-00001; now status paid) has journal JE-2026-00014 posted 2026-11-25 22:49:27 Dr 1120 5,200 / Cr 4100 4,500 / Cr 2120 700, i.e. posted at creation; the draft state itself could not be re-observed because the invoice has since been paid - fresh draft re-run scheduled in T-P09-001.  | RUN-003 FULL REPRODUCTION (AUD-CINV-001, headed Chrome): new invoice id 3 status=draft; journal JE-2026-00050 posted at creation Dr 1120 8,320 / Cr 2120 1,120 / Cr 4100 7,200 (STILL REPRODUCIBLE). 

## Which reports and decisions are affected
Financial reports that read the affected data: Customer invoices

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

