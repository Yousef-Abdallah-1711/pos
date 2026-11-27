# F-049 — Financial impact

- **Quantified impact (register):** Any customer refund cannot be recorded; every refunded sale overstates revenue, VAT and cash
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, capability

## What is wrong in the books or reports
DELETE -> 400 "cannot cancel paid or refunded order before a supported refund and accounting reversal workflow exists"; payment-status refunded/failed/pending from paid -> 400 invalid transition; voided -> "invalid payment status"; status=cancelled -> 500; paid revenue/COGS stay posted; cash 307,330.01 unchanged | RUN-003: paid order cannot be cancelled/refunded/voided through any endpoint; explicit guard message (REPRODUCED). 

## Which reports and decisions are affected
Financial reports that read the affected data: Orders; Payments; Accounting (revenue/COGS/VAT/inventory reversal)

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

