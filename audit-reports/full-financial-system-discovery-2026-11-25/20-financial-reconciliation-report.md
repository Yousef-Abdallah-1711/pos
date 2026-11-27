# 20 — Financial reconciliation report (tenant C, diagnostic COA; end of run)

Sources: `evidence/runtime/24-final-reconciliation-tenantC.txt`, `07-...`, `16-financial-oracle-results.csv`, `44-transaction-level-reconciliation.csv`. Amounts in EGP. Isolated DB snapshot: `evidence/database/isolated-audit-final-2026-11-25.dump`.

## Ledger integrity
- GL total debits = credits = 310,000.00 (difference 0.00). Journal-line sum equals GL sum. Tenant C holds 28 journals: 21 posted, 1 voided (capital), 6 manual drafts created by boundary tests (not in ledger); detail in evidence/runtime/24. (Whole DB: 30 journals incl. 2 for tenant A.)
- Every individual journal balanced (no unbalanced entry accepted; BND-01 rejected).
- This proves arithmetic integrity only; **balanced journals are not proof of correct accounting** (see below).

## Operational vs ledger (tenant C)
| Item | Operational | Ledger | Difference | Cause |
|---|---|---|---|---|
| Paid order subtotals vs Food Sales 4100 | 30,000 (3 paid orders) + invoice 4,500 | 34,500 | 0 | consistent |
| Paid order tax + invoice tax vs VAT payable | 3,000 + 700 | 3,700 | 0 | consistent (but wrong vs UI/14%: F-004, F-007) |
| COGS for 3 paid orders | 12,000 | 12,000 | 0 | consistent |
| Unpaid orders (8) | 88,000 total, stock consumed | not in ledger | 88,000 | card/wallet/duplicates unpaid (F-005, F-017) |
| Customer invoice balance vs AR 1120 | 0.00 | -1,800 | -1,800 | concurrent over-collection (F-006) |
| Supplier invoice balance vs AP 2101 | 800 | 800 | 0 | consistent |
| Stock × cost vs Inventory 1130 | 156,000 | 8,000 | 148,000 | opening stock without journal (F-011) |
| Card Payment Clearing 1113 | – | 11,000 | 11,000 | never cleared (F-026) |
| Owner Capital 3100 | 100,000 posted | 0 | -100,000 | voided in FIN-001 test (F-013) |

## Residuals
- **Ledger-level (journal debit−credit, TB debit−credit): signed 0.00, absolute 0.00.**
- **Balance-sheet equation (S13 snapshot): signed −2,000.00; absolute 2,000.00** (F-002).
- **Documented discrepancies are not summed into one number:** they are different kinds (policy-dependent, defect, test-induced) and summing could let offsets cancel. Per-item signed/absolute values are in 16 and 44. Largest single discrepancy: inventory valuation vs GL 180,000 (S13 snapshot) / 148,000 (end).
- Original historical 5,000-piaster residual: **not revalidated** (original fixture not available; no replacement claimed).

## Conclusion
Transaction-level agreement holds for the happy-path events that post (supplier invoice, supplier payment, customer receipt, manual journals, paid-order sales/COGS) but fails for tax consistency, payroll accrual, card/wallet capture, concurrency on customer settlement, inventory valuation, and statement-level reports. Aggregate agreement does not prove each transaction correct; the per-transaction table shows where it was checked.
