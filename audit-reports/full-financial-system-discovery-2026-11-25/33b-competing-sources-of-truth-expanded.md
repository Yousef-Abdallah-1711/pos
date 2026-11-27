# 33b — Competing sources of financial truth (expanded)

Extends the 12 sources in [27](27-financial-data-lineage.md). Built from [12](12-report-formula-catalog.md), [30](30-report-metric-inventory.csv) and [31](31-independent-finding-verification.csv). Evidence: CONFIRMED-BY-SOURCE unless noted; nothing runtime-verified. "Legit" = difference that accounting explains; "Inconsistency" = not explained by policy.

| Concept | Operational source | Accounting (GL) source | Reporting / dashboard / export | Intended authority (proposed, needs approval) | Divergence | Type | Reconciliation test |
|---|---|---|---|---|---|---|---|
| Gross sales | `orders.total_amount` (incl. tax/delivery) | not posted as such | Sales Analytics (all non-cancelled, any payment status); Order stats (ALL orders); KPI (paid only) | order table, paid-only, defined once | 3 filters on the same table | Inconsistency | SCN-06 + unpaid/cancelled orders; compare MET-SALES-001, MET-ORD-001, MET-KPI-001 |
| Net sales | order_items.total_price (excl. tax/delivery) | revenue accounts net of tax | Profitability revenue; P&L revenue | GL | tax/delivery inclusion differs | Legit + inconsistency (unlabelled) | SCN-06/12 |
| Recognised revenue | payment_status paid (outbox) | revenue journal dated order creation date (per 11) | P&L | GL after policy D-REV-01 | date basis = creation, trigger = paid | Inconsistency (policy open) | order created day 1, paid day 5 across a period boundary |
| Collected revenue | no payment records for orders (N21); payment_status only | cash debit in sales journal | KPI cash_position (bank cache); Predictive CF (AR/AP only) | GL cash+bank | no tender/split records | Inconsistency | SCN-08 split tender |
| Receivables | invoices_customer.balance_due | GL 1120 (invoices post at creation incl. drafts) | AR aging (current balance_due, as_of ignored for payments) | GL control account | aging uses subledger now-balance | Inconsistency | SCN-11 backdated as_of |
| Supplier liabilities | invoices_supplier.balance_due; suppliers.current_balance (maintenance unverified, VER-023) | GL 2101 (only after approval; supplier payment posting is dead code VER-007) | AP aging (includes drafts) | GL | draft included; payments may not post | Inconsistency | SCN-04 approve then pay |
| Cash balance | none for orders | GL 1110 | KPI cash_position reads bank_accounts cache, not GL | GL | cash on hand excluded | Inconsistency | S13 compare GL cash vs KPI |
| Bank balance | bank_accounts.current_balance (cache) | GL via gl_account_id | KPI cash_position | GL reconciled by bank-reconciliation sessions | cache vs GL | Unknown until run | SCN-20 |
| Inventory quantity | products.quantity_in_stock (decremented at order confirm; reservations mig 180) | none | Inventory valuation; Turnover (reads legacy inventory table, VER-025) | stock movements table | stock decrement at confirm vs GL COGS at payment | Inconsistency | REL-008, SCN-14 |
| Inventory value | qty × products.cost (current) | GL 1130 (supplier invoices go to expense 5110, not 1130; VER-006) | Inventory valuation (as_of ignored) | GL after costing policy INV-01 | valuation will not tie to GL | Inconsistency | S13 T2/T5 |
| COGS | order_items.qty × products.cost (reports) | COGS journal at payment (cost at that time) | Profitability (current cost), P&L (GL) | GL | cost timing, gross vs net-of-cancelled qty | Legit (timing) + inconsistency | SCN-06 change cost between sale and report |
| Gross profit | item revenue − COGS (Profitability) | none (P&L has no gross profit subtotal) | Profitability; Menu profitability (fails: cost_price missing) | GL once subtotal exists | missing subtotal in P&L | Inconsistency | S13 gross 6000 |
| Net profit | n/a | P&L = revenue − expenses (GL) | KPI net profit (orders revenue − journal expenses); BS equity (excludes it) | GL P&L | 3 different values possible | Inconsistency | S13 expected −2000 |
| VAT | orders.tax_amount (10% server vs 14% POS UI) | GL 2120 output / 1150 input | VAT report (invoice tables only; POS/storefront VAT excluded) | GL | order VAT never reaches VAT report | Inconsistency | SCN-06 plus invoice; compare |
| Discounts | order discount never assigned (N17); invoices.discount_amount | 4150 contra account | KPI "net of discount" | GL | invoice trigger taxes gross (VER-027) | Inconsistency | SCN-12 |
| Refunds | status 'refunded' value only; blocked from paid (N2) | none | none | n/a | no flow | Missing capability | SCN-13 capability test |
| Payroll cost | salaries.net_salary/overtime | payroll journal (sync) | HR cost analysis (net pay, paid only); HR overtime report estimates from attendance | GL | net vs gross/employer cost; overtime two sources | Legit + inconsistency | S13 T7/T8 |
| Payroll liabilities | salaries status | GL payroll payable | none dedicated | GL | no report | Gap | S13 payable 2000 |
| Operating expenses | none (no expense module, N8) | manual journals / supplier invoices | P&L; KPI breakdown | GL | KPI excludes void entries, P&L includes | Inconsistency | FIN-008 |
| Equity | n/a | GL equity accounts | BS equity excludes current-period earnings (VER-001) | GL + closing entries | BS unbalanced by net result | Inconsistency | S13 BS residual |
| Retained earnings | n/a | no closing entries (N10) | none | GL after CLOSE-01 | never moves profit | Missing | SCN-23 capability |
| Current-period earnings | n/a | derivable from revenue/expense | not shown on BS | policy EARN-01 | | Missing | S13 |

## Legitimate differences (do not flag as defects once labelled)
Operational vs accounting timing (order date vs payment date); gross (tax-inclusive) vs net sales; current cost vs historical cost for estimates; sub-ledger vs control account while unposted drafts exist — **provided** each report states its basis (none currently does).

## Observation
Revenue has at least six definitions (MET-PNL-001, MET-SALES-001, MET-PROF-001, MET-KPI-001, MET-ORD-001, dashboard tile). None is labelled with its basis in the catalogue; UI labels were not inspected.
