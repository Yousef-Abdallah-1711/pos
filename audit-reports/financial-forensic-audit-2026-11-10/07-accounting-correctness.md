# Accounting correctness

**Status: PARTIAL — one balanced posting and selected APIs exercised; no independent accountant sign-off.**

## Runtime invariants

- Two-line debit 100 / credit 90 was rejected HTTP 400: journal entry is not balanced.
- One-line entry rejected HTTP 400: journal entry must contain at least two lines.
- Balanced EGP 100 manual journal persisted and posted.
- Trial balance returned 100 debit, 100 credit, is_balanced true.
- Balance Sheet returned assets 100, liabilities 0, equity 0, equation_balanced false, while P&L net profit was 100. Current-period profit/retained earnings treatment appears absent from that response; report-query review is needed to locate root cause.

## Source risks

- Auto-posting normalizeBalancedLines adjusts a final line to absorb differences above 0.009 (auto_posting_engine.go:633 onward), potentially reassigning material value rather than diagnosing the source.
- Order/accounting amounts use float64 while journal columns use NUMERIC(15,2); currency and rounding semantics need end-to-end review.
- Order tax is fixed at 10% (order_usecase.go:173); jurisdiction is unknown.
- Manual-journal UI omitted JSONB attachments and failed; API accepted an explicit empty array.

IAS 1 describes a complete financial statement set including position, profit/loss/OCI, changes in equity, cash flows and notes. IAS 2 covers inventory. IFRS 15 covers revenue recognition. These are reference points, not a conclusion that IFRS applies. Tax, fiscal year, currency, inventory costing, revenue recognition and close policy require local professional decisions.

## Root-cause closure — 2026-11-11

### FIN-001 (P1) — VERIFIED, four fixtures

`backend/internal/usecase/financial_report_usecase.go:261-336` computes each balance-sheet total only from ledger balances of active non-header asset, liability, and equity accounts. It has no current-period earnings or retained-earnings calculation. Independent posted journal fixture states produced:

| As of | Assets | Liabilities | Equity | P&L net profit | Equation |
|---|---:|---:|---:|---:|---|
| 2026-11-10 | 100 | 0 | 0 | 100 | Fails by 100 |
| 2026-11-11 | 60 | 0 | 0 | 60 | Fails by 60 |
| 2026-11-12 | 85 | 25 | 0 | 60 | Fails by 60 |
| 2026-11-13 | 135 | 25 | 50 | 60 | Fails by 60 |

The latter three include balanced revenue/expense, liability, asset, and equity activity; the persistent EGP 60 gap equals cumulative net profit. No signed-off close/retained-earnings policy was supplied, but the queried omission is confirmed.

### FIN-003 (P1) — normalizer behavior VERIFIED; actual auto-post persistence BLOCKED by separate defects

`backend/internal/usecase/auto_posting_engine.go:633-660` changes the final line for any rounded difference over 0.009. A runtime probe transformed Dr 100 / Cr 90 into Dr 100 / Cr 100 by increasing the final revenue credit by EGP 10. The actual `PostSupplierInvoiceEntry` call did not persist because nil attachments caused a PostgreSQL JSONB error. Using valid JSON and the normalized lines with `entry_type: manual` through the journal/ledger repository persisted a posted balanced journal with two ledger rows. Thus material alteration is verified; persistence through that conditional repository route is verified; present automatic supplier-method persistence is not demonstrated. Reverse imbalance Dr 100 / Cr 110 was rejected after normalization created a two-sided line that violated line constraints.

### Additional automatic posting blockers

`createAndPostEntry` (`auto_posting_engine.go:522-559`) omits `Attachments`; automatic supplier insert fails with invalid JSON. With valid attachments, `auto_supplier_invoice` (21 chars) exceeds migration 137 `journal_entries.entry_type VARCHAR(20)`. Both are evidenced in the isolated overlay. The order paid-posting call logs failures and suppresses them at `order_usecase.go:492-520`; no paid-order runtime trace was performed.

Required regression gates: current and prior-period earnings/retained-earnings fixtures reconcile to the trial balance; material imbalances reject before persistence; bounded approved rounding has a distinct auditable line; every automatic posting type fits the schema; JSONB fields have valid defaults; injected automatic-post failure is visible/recoverable and cannot leave an untracked draft.

The normalizer has no configured maximum adjustment: it uses the full rounded debit/credit difference. Callers include payroll and the shared `createAndPostEntry` paths for order sales, COGS, supplier invoice/payment and depreciation. Database NUMERIC and journal line checks impose downstream limits, but these are not accounting materiality controls.
