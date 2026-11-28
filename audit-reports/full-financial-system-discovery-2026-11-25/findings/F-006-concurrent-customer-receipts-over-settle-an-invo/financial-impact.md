# F-006 — Financial impact

- **Quantified impact (register):** 1800
- **Accounts / statements named in the evidence:** 1120
- **Impact classes:** accounting, concurrency

## What is wrong in the books or reports
5 accepted; amount_paid 7000 on 5200 invoice; AR GL -1800 | RUN-003 TRANSACTION-LEVEL TRACE (AUD-AR-RESIDUAL): GL 1120 = Dr 5,200 (INV-2026-00001) + Dr 8,320 (INV-2026-00002 draft, F-008) - Cr 7,000 (6 receipts: 2,000 + 5x1,000) = 6,520; aging 8,320; residual 1,800 = receipts 7,000 minus invoice 5,200 = over-settlement of invoice 2 (amount_paid 7,000 > total 5,200, balance clamped to 0). ROOT CAUSE CONFIRMED (source): backend/internal/repository/payment_acct_repo.go:445-449 customer invoice balance SELECT has no FOR UPDATE (supplier twin at :499-503 has FOR UPDATE OF inv), and applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0) clamp that hides the overpayment.  | RUN-003 FRESH-INVOICE RE-REPRODUCTION (AUD-CINV-FLOW): invoice 3,000, 5 concurrent receipts of 1,000 -> all 5 accepted (201), amount_paid 6,000, AR credit -3,000; single sequential overpayment 5,000 is rejected (guard exists but is not serialized). RUN-003 TRANSACTION-LEVEL TRACE (AUD-AR-RESIDUAL): GL 1120 = Dr 5,200 (INV-2026-00001) + Dr 8,320 (INV-2026-00002 draft, F-008) - Cr 7,000 (6 receipts: 2,000 + 5x1,000) = 6,520; aging 8,320; residual 1,800 = receipts 7,000 minus invoice 5,200 = over-settlement of invoice 2 (amount_paid 7,000 > total 5,200, balance clamped to 0). ROOT CAUSE CONFIRMED (source): backend/internal/repository/payment_acct_repo.go:445-449 customer invoice balance SELECT has no FOR UPDATE (supplier twin at :499-503 has FOR UPDATE OF inv), and applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0) clamp that hides the overpayment. 

## Which reports and decisions are affected
Financial reports that read the affected data: Customer invoices; Receipts

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.
Evidence of over-collection: tenant 4 invoice 2 (7,000 on 5,200 -> AR residual 1,800) and fresh invoice 4 (6,000 on 3,000).
