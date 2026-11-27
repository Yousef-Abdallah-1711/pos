# 42 — Complete operating month reconciliation

**STATUS: PARTIAL (RUN-003 compressed operating month executed on a fresh tenant; see the final section — the text below is the earlier batch status).** A multi-day restaurant month (daily purchases, sales, payroll, expenses, bank movements, cancellations, refunds) was not run. Reasons: time/context budget of this session; and capability blockers found that prevent a faithful month: no employee creation from the UI (F-014), no card/wallet payment capture (F-005), no refund workflow (N2), no bank transfers (N9), no expense module (N8), storefront blocked (F-022), no fiscal-period lifecycle (F-013).

What exists instead: a single-day deterministic control run (S13 variant) plus ad-hoc business events on diagnostic tenant C, fully reconciled in [44-transaction-level-reconciliation.csv](44-transaction-level-reconciliation.csv) and [20-financial-reconciliation-report.md](20-financial-reconciliation-report.md). That is **not** a substitute for a month-long reconciliation and must not be counted as one.

---
## RUN-003 compressed operating-month simulation (fresh tenant F, 2026-11-26) — PARTIAL (payroll and dated sales not possible)

Script `tests/runtime/ph14_month.mjs` (API-DIAGNOSTIC, independent integer-cents ledger); evidence `evidence/database/RUN-003/PH-14/AUD-MONTH.json`. Fresh tenant "Audit F" (registered by API; chart of accounts copied from tenant 4 by SQL fixture because of F-001).

| Event group (dates 2026-11-01..11-30) | Count | Result vs independent ledger |
|---|---|---|
| Capital contribution 100,000 | 1 | exact |
| Supplier invoices (no discount, 14% VAT) + 3 payments incl. one partial | 3 + 3 | exact (AP 7,779.40, input VAT 2,899.40) |
| Customer invoices (14% VAT) + 3 receipts incl. partial | 2 + 3 | exact (AR 1,850.00, output VAT 1,160.00 with POS VAT) |
| Rent/expense journals | 4 | exact (15,000) |
| Owner withdrawal via custom account 3150 | 1 | exact (no drawings account in seeded COA) |
| POS cash sales (10 orders, 25 units, 10% VAT, COGS 40/unit) | 10 | exact: revenue 2,500, VAT 250, COGS 1,000 |
| Payroll | 0 | NOT EXECUTABLE: employee creation is blocked (F-014) |

Independent ledger vs system GL: **0 differences on every account**. Trial balance balanced. P&L equals oracle (revenue 9,000; expenses 36,710; net loss -27,710).

Defects reproduced on the clean dataset:
- **F-002** balance sheet equation difference -27,710.00 = exactly the net loss (current-period earnings absent from equity); `equation_balanced=false`, `integrity_status=UNRECONCILED`.
- **F-003** cash flow reported 69,290.00 vs actual cash change 74,480.00; residual 5,190.00 = -1,850.00 (AR) +1,000.00 (inventory) -2,899.40 (input VAT) +7,779.40 (AP) +1,160.00 (VAT payable) — the working-capital terms the formula omits.

Limits (not hidden): orders cannot be back-dated so all sales land on the run date; payroll, bank reconciliation, depreciation, period close (F-013) and inventory opening/purchase accounting (F-011) were not part of the month. A full month with every module therefore remains blocked by F-014, F-013 and F-011.
