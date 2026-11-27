# 12 - Report and formula catalogue

Baseline SHA 133260db66b1143c254d79476daa7650b41a33f0. All line refs are relative to `backend/internal/` unless stated. Routes: `backend/cmd/api/main.go:1116-1117, 1467-1497` and `backend/cmd/api/routes.manifest`. The dev DB has no transactional data, so no report was executed against data; formulas below are read from source. Where behaviour was only inferred it is marked UNVERIFIED.

## 0. Source-of-truth map (mixed sources are the main risk)

| Family | Reads | Date field | Basis |
|---|---|---|---|
| Ledger reports (P&L, Balance Sheet, Cash Flow, Trial Balance, Ledger) | `journal_entries`/`journal_entry_lines` (P&L, Cash Flow) or `general_ledger` (Balance Sheet, Trial Balance) joined to `chart_of_accounts` | `journal_entries.entry_date` (accounting date) / `general_ledger.entry_date` | Accrual, posted journals |
| Sub-ledger reports (AR/AP aging, VAT) | `invoices_customer` / `invoices_supplier` current columns | `due_date` (aging), `invoice_date` (VAT) | Operational tables, NOT GL |
| Inventory Valuation | `products.quantity_in_stock x products.cost` | none (as-of ignored) | Operational snapshot |
| Sales Analytics / Profitability / Inventory Turnover / Overview | `orders` + `order_items` + `products` (+ legacy `inventory`) | `orders.created_at` or `order_items.created_at` | Operational, payment_status NOT filtered |
| Advanced Insights (KPI, Break-even, Menu profitability) | MIX: `orders` (paid) for revenue + posted journal expense accounts + `bank_accounts_ledger.current_balance` | `orders.created_at`, `journal_entries.entry_date` | Mixed cash/accrual/operational |
| Predictive cash flow | `payments_received`, `payments_made` | `payment_date` | Cash, operational |
| HR reports | `employees`, `attendance`, `leaves`, `training_records`, `salaries`, `insurance_tracking` | various | Operational |

Hardcoded constants found: order tax 10% (`usecase/order_pricing.go:11,46`), delivery fee forced to 0 (`:52`), currency "EGP" (`:99`; frontend Financial Hub `reports/financial/page.tsx` Intl currency EGP), P&L comparison period = same length immediately before (`financial_report_usecase.go:126-131`), aging buckets 30/60/90 days (`:443-446`), turnover slow-moving threshold 3.0/month (`handler/http/report_handler.go:2553 slowThreshold := 3.0`), overtime multiplier 1.5 and 4.33 weeks/month (`handler/http/hr_report_handler.go:650-655`), predictive forecast weights 0.7/0.3 (`repository/accounting_advanced_repo.go:1198`), default COA codes in `domain/auto_posting.go:33-91`. CoA seeds VAT 14% on 1150/2120 (`chart_of_accounts.tax_rate`) which conflicts with the 10% order tax.

## 1. Ledger-based financial reports

### 1.1 Profit & Loss - `GET /api/v1/reports/financial/profit-loss` (perm module 17 READ)
- Handler: `handler/http/financial_report_handler.go:42-84`; use case `usecase/financial_report_usecase.go:120-250` (`generateProfitLossCore`).
- Params: `date_from`, `date_to` (inclusive, default last 30 days UTC; `parseFinancialDateRange` :276-296), optional `compare_date_from/to`.
- SQL (:183-211): movements CTE = lines of journals with `status='posted'` OR (`status='void'` AND `posted_at IS NOT NULL`) where `entry_date BETWEEN from AND to`, UNION ALL inverse of lines of voided journals where `voided_at::date BETWEEN from AND to`; join `chart_of_accounts` filtered to `account_type IN ('revenue','expense')`; group by account.
- Formulas: revenue line = credit - debit; expense line = debit - credit; `TotalRevenue = sum(revenue lines)`; `TotalExpenses = sum(expense lines)`; `NetProfit = TotalRevenue - TotalExpenses` (:248). Contra-revenue 4150 Sales Discounts (debit) reduces revenue automatically. COGS (5110) appears under "Expenses"; there is no Gross Profit / operating profit subtotal.
- Comparison: `percentChange = ((cur-prev)/prev)*100`, prev==0 -> 100 (or 0 if both 0) (:532-540). Default comparison window computed :126-131.
- Date basis: posting/accounting date `entry_date` (for auto sales = ORDER CREATION date, see matrix). Void reversal is dated `voided_at` (:197) so the original period keeps the revenue and the void period shows the reversal.
- Rounding: none in SQL; Go float64 sums, JSON floats. Cache: none (`respondFinancialReportWithoutCache`). Export: JSON only (no CSV/PDF export endpoint for the financial statements; `POST /reports/export` covers only sales/profitability/inventory).
- Flags: float64 accumulation; includes only posted journals so unpaid orders excluded (vs Sales Analytics which includes them) - the two reports will not reconcile by design.

### 1.2 Trial Balance - `GET /accounting/trial-balance` and `GET /reports/financial/trial-balance`
- Handlers: `handler/http/ledger_handler.go` (GetTrialBalance) and `financial_report_handler.go:183-200`; repo `repository/general_ledger_repo.go:411-470` (GetTrialBalance), use case `:252-274`.
- SQL: active, non-header accounts LEFT JOIN `general_ledger` with `entry_date <= asOf`; `net_balance = SUM(credit-debit)` if normal_balance credit else `SUM(debit-credit)`; row debit/credit column derived from sign and normal balance (mapTrialBalanceRow :542+). `IsBalanced = round2(totalDebits) == round2(totalCredits)` (:272).
- `round2` = `float64(int(v*100+0.5))/100` (:528-530) - wrong for negative values.
- Source: general_ledger (includes void reversal rows dated at void date). Inactive accounts with balances are silently excluded (`is_active = TRUE`) - UNVERIFIED impact on balance.

### 1.3 Balance Sheet - `GET /reports/financial/balance-sheet`
- Handler `financial_report_handler.go:86-103`; use case `financial_report_usecase.go:276-364`.
- SQL (:284-307): accounts `is_active AND NOT is_header AND account_type IN (asset, liability, equity)` LEFT JOIN `general_ledger` `entry_date <= asOf`; balance = credit-debit for credit-normal accounts else debit-credit; `balanceSheetContribution` flips sign when an account's normal_balance differs from the type default (:355-364) (so 1240 Accumulated Depreciation (credit) reduces assets; contra-equity would reduce equity).
- Totals: `TotalAssets`, `TotalLiabilities`, `TotalEquity`; check `EquationBalanced = round2(A) == round2(L+E)`; difference exposed; status RECONCILED / UNRECONCILED (:344-351).
- **Gap**: revenue/expense accounts are excluded and there is no closing-entry or "current year earnings" line (3300 exists in CoA but is unused; diagnostics only: `usecase/current_earnings_diagnostic.go:13 POLICY_PENDING`). Any period with profit/loss makes the equation unbalanced by exactly net earnings. Source: GL only. No cache.

### 1.4 Cash Flow - `GET /reports/financial/cash-flow`
- Handler `:105-122`; use case `:366-400`.
- Same movements CTE as P&L; formulas (:388-392):
  - `OperatingCashFlow = SUM(credit-debit)` over revenue+expense accounts (= net profit, accrual, NOT cash)
  - `InvestingCashFlow = SUM(debit-credit)` over accounts whose CODE starts with `12` (fixed assets; sign is debit-credit, i.e. an asset purchase is POSITIVE)
  - `FinancingCashFlow = SUM(credit-debit)` over equity-type accounts OR codes starting `22`
  - `NetCashFlow = Operating + Investing + Financing`.
- **Flags**: this is not a cash flow statement. It ignores cash/bank account movements, working-capital changes (AR, AP, inventory, VAT), and mis-signs investing; classification by hardcoded code prefix ("12%", "22%"). Source: posted+void journals (not GL). Contrast with Predictive Cash Flow (section 3.3), which uses payments tables.

### 1.5 Ledger (account statement) - `GET /accounting/ledger`
- `handler/http/ledger_handler.go` -> `repository/general_ledger_repo.go` (account entries with `running_balance` column stored at post time: computed from prior GL rows with `entry_date <= new entry_date` (journal_entry_repo.go:553-654 appendGeneralLedgerTx)). Stored running_balance on backdated entries does not rewrite later rows (documented in code comment) - UNVERIFIED how the UI orders/reads it.
- `chart_of_accounts.current_balance` is a cache recomputed from SUM(GL) after each posting (journal_entry_repo.go:~630-654); reports do not read it.

## 2. Sub-ledger financial reports

### 2.1 Receivables / Payables Aging - `GET /reports/financial/receivables-aging`, `/payables-aging`
- Use case `generateAging` `financial_report_usecase.go:428-486`; handlers `:124-161`.
- SQL: from `invoices_customer` (or `invoices_supplier`) where `status <> 'void' AND balance_due > 0`, grouped by customer/supplier; buckets by `due_date` vs `asOf`: current `due_date >= asOf`; 1-30 `due_date < asOf AND >= asOf-30d`; 31-60; 61-90; >90 (:442-446). `Total = sum of buckets`.
- **Flags**: uses CURRENT `balance_due`, ignoring the `asOf` for payments made after that date (not a true historic as-of); draft invoices (customer: journal posted at create; supplier: draft has no journal) are included in AP aging even though unapproved supplier invoices have no ledger payable; no reconciliation to GL 1120/2101; entity-name `Unknown` when customer missing; the dynamic SQL builds table/column names from constants only (no injection). No cache, no export.

### 2.2 VAT Report - `GET /accounting/tax/vat-report`; tax periods CRUD `/accounting/tax/periods`
- `handler/http/tax_handler.go:~176`; repo `repository/tax_repo.go:260-291 GetVATSummary`.
- `OutputVAT = SUM(invoices_customer.tax_amount)` and `InputVAT = SUM(invoices_supplier.tax_amount)` where `status <> 'void'` and `invoice_date BETWEEN from AND to`; `NetVAT = Output - Input`.
- **Flags**: (a) reads invoice tables, not GL 2120/1150; VAT charged on storefront/POS orders (posted to 2120 by auto_sales) is NOT included; (b) VAT rate is not a parameter anywhere - tax_amount is whatever the invoice line carries; orders use a hardcoded 10% (`order_pricing.go:46`) while CoA tax_rate is 14; (c) includes draft/unapproved supplier invoices; (d) `tax_periods.output_vat/input_vat/net_vat` columns are not populated by close (`CloseTaxPeriod` tax_repo.go:293-318 only flips status) - periods are not frozen; (e) frontend Financial Hub shows these numbers as KPI cards for the last 30 days (`frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/page.tsx:41-66`).

### 2.3 Inventory Valuation - `GET /reports/financial/inventory-valuation`
- `financial_report_usecase.go:488-526`; handler `:163-181` (the ONLY financial report cached: KV cache TTL 5 minutes, key tenant+restaurant+path+query, `financial_report_handler.go:21,177-181,202-241`).
- SQL: products with `track_inventory=true AND status <> 'deleted'`: `quantity_in_stock` x `cost` (both read as DOUBLE PRECISION); `TotalValue = sum(qty * cost)`.
- **Flags**: `asOf` is echoed but NOT used (always current stock); valuation method = current standard `products.cost` (not last purchase cost used by COGS (`inventory_movement_repo.go:158-206`) nor weighted average); products with NULL cost are valued 0 silently; stock quantity decremented at order confirmation while GL inventory (1130) is relieved at payment - valuation will not tie to GL 1130; 5-minute stale cache after stock changes.

## 3. Advanced Insights (`handler/http/accounting_advanced_handler.go`, repo `repository/accounting_advanced_repo.go`)
Date params: `parseDateRangeForReports` (`handler/http/accounting_advanced_handler.go:1170-1198`): default from today-30d, `date_to` inclusive converted to exclusive (+1 day), requires to > from.

### 3.1 Financial KPIs - `GET /reports/financial/kpis` (repo :837-964)
- `revenue = SUM(orders.total_amount)` where `payment_status='paid' AND created_at >= from AND < to` (operational table; gross incl. tax and delivery; net of discount).
- `expenses = SUM(debit-credit)` over lines of `journal_entries.status='posted'` on `account_type='expense'` accounts with `entry_date >= from AND < to`. (**voided entries are excluded here (unlike P&L)**.)
- `net_profit = revenue - expenses`; `cash_position = SUM(bank_accounts_ledger.current_balance)` for active accounts (cached balance column, as of NOW, ignores date range; excludes cash on hand 1110).
- `revenue_trend`: monthly `date_trunc('month', created_at)` sum of paid order totals; `expense_breakdown`: top 10 expense accounts by `COALESCE(account_subtype, account_name)`.
- **Flags (mixed sources)**: revenue from orders (including VAT, which is a liability) vs expenses from journals; revenue excludes AR invoices and the discount account; net_profit therefore differs from P&L NetProfit; revenue date = order creation not payment; cash_position unrelated to the period.

### 3.2 Menu Item Profitability - `GET /reports/financial/menu-profitability` (:966-1025) and Break-even - `GET /reports/financial/break-even` (:1027-1103)
- **Defect (confirmed)**: both queries reference `products.cost_price` (lines 976 and 1045), but the column is named `cost`. Running the query on the dev DB gives `ERROR: column p.cost_price does not exist` (`evidence/database/cost_price_probe.txt`; SQL `sql/12_cost_price_column_probe.sql`). Both endpoints therefore return a server error at runtime regardless of data. (Grep of `cost_price` across backend finds only these two uses.)
- Intended formulas: menu: per product `quantity = SUM(qty - cancelled)`, `revenue = SUM(total_price)`, `cogs = SUM((qty-cancelled)*cost)`, `gross_profit = revenue - cogs`, `margin% = gross_profit/revenue*100`, paid orders only, `LIMIT` default 20. Break-even: `variable_costs = SUM((qty-cancelled)*cost)`, `fixed_costs = MAX(0, total journal expenses - variable_costs)` (double counts COGS because journal expenses already include 5110), `contribution = revenue - variable`, `CM ratio = contribution/revenue`, `break_even_revenue = fixed / CMratio`, `gap = revenue - break_even`.

### 3.3 Predictive Cash Flow - `GET /reports/financial/predictive-cash-flow` (:1105-1215)
- History from `payments_received`/`payments_made` (status <> 'void') by `date_trunc('month', payment_date)` for N months (default 6); forecast K months (default 3): `projected = 0.7*lastMonthNet + 0.3*avgNet + slope*i`, `slope = (last-first)/(n-1)` (:1198). Excludes order sales receipts (POS/website cash), payroll and any non-subledger cash; labelled "cash flow" but only AR/AP settlements.

### 3.4 HR Cost Analysis - `GET /reports/hr/cost-analysis` (:1218-1315)
- `payroll_cost = SUM(salaries.net_salary)` for `status='paid' AND pay_period_end in [from,to)` (net pay, not employer cost; excludes gross deductions liabilities); `overtime_cost = SUM(overtime_amount)`; `insurance = SUM(insurance_tracking.monthly_premium) x ceil(days/30)` for active policies (estimate, ignores policy dates); `total_hr_cost = payroll + insurance`; `labor_pct = total_hr_cost / revenue(paid orders) * 100`; by department via employees.department.

## 4. Operational reports (`handler/http/report_handler.go`)
Routes: `GET /reports/overview|sales|profitability|inventory`, `POST /reports/export` (main.go:1467-1471; perm module 8; export audited in `report_export_audit_logs`, `logReportExportAudit` :1059). Params `date_from/date_to`, `parseReportDateRange` (:2656) -> half-open [from, to+1d). Caches: in-process `reportResponseCache`; TTL overview 30s, sales 30s, profitability 60s, inventory 60s (`:280-283`); key = tenant+restaurant+dates (+ granularity/limit).

### 4.1 Sales Analytics - `GET /reports/sales` (`GetSalesReport` :351-470)
- Summary (:2174): `total_revenue = SUM(orders.total_amount)`, `total_orders = COUNT(*)`, `avg_ticket = AVG(total_amount)` over orders with `status <> 'cancelled'` and `created_at` in range - **no payment_status filter**: pending, failed and unpaid delivered orders count as revenue; includes tax and delivery.
- Trend (:2195): bucket day/week (IYYY-"W"IW)/month by `date_trunc(created_at)`, same filter. Hourly (:2237): `EXTRACT(HOUR FROM created_at)` (DB timezone, not tenant timezone). Top items (:2283): `SUM(oi.quantity)`, `SUM(oi.total_price)` filtered by `oi.created_at` (item timestamp, not order) and `o.status <> 'cancelled'`; quantity is gross of `cancelled_quantity`, revenue is net -> qty and revenue inconsistent for partially cancelled items. Top categories (:2323) similarly.
- Export: CSV (`buildSalesCSV` :1425, with `csvSafeTextCell` formula-injection guard :1456) and PDF (`buildSalesPDFSummary` :1531; Arabic fonts `report_fonts`).

### 4.2 Profitability - `GET /reports/profitability` (:471-556)
- Summary (:2364): `revenue = SUM(oi.total_price)`, `COGS = SUM(oi.quantity * COALESCE(p.cost,0))`, `gross_profit = revenue - COGS`, `gross_margin% = GP/revenue*100`, `food_cost% = COGS/revenue*100`; flags `units_missing_cost`. Filter: `o.status <> 'cancelled'`, `oi.created_at` range - again no payment filter, and uses CURRENT product cost (not historical cost used in the GL COGS entry) and gross quantity.
- Items (:2395): same per product, ordered by gross profit; top N + lowest N. Revenue here excludes tax and delivery (item totals) whereas Sales Analytics revenue includes them: the two operational reports intentionally differ but are not labelled.
- Export CSV/PDF (:1460, :1609).

### 4.3 Inventory Turnover - `GET /reports/inventory` (:557-645; `queryInventoryTurnoverItems` :2480-2586)
- `opening_qty = latest legacy inventory.quantity_after before from` else current stock; `avg_inventory_qty = (opening + current)/2`; `avg_inventory_value = avg_qty * cost`; `estimated_cogs = SUM(oi.quantity * cost)` of non-cancelled orders in range; `turnover = estimated_cogs/avg_inventory_value`; `monthly_turnover = turnover * 30/daysInRange`; slow-moving if monthly < 3.0 or (stock>0 and sold=0); `waste_qty = SUM(ABS(quantity_change))` from legacy `inventory` where reason='waste'.
- **Flags**: reads the legacy `inventory` table that nothing writes (no INSERT into `inventory` in Go; `/inventory/adjust` writes `inventory_movements`), so opening = current stock and waste = 0 always; `total_current_inventory_value = qty * cost` (:2588-2609).
- Export CSV/PDF (:1502, :1676).

### 4.4 Overview - `GET /reports/overview` (:646-813)
Aggregates sales summary + top item + profitability summary + top profitable item + inventory summary under the same cache; inherits all flags above.

### 4.5 Order statistics - `GET /admin/orders/stats` (`repository/order_repo.go:859-950`)
`total_revenue = SUM(total_amount)` over ALL orders (including cancelled/unpaid) in `created_at <= endDate` inclusive; `delivered_today` counts orders created today regardless of status (label vs logic mismatch).

### 4.6 Dashboard home (frontend) 
`frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx:271-340` loads order lists, product/category/employee/driver counts, table stats - counts only; no money KPIs computed client-side (UNVERIFIED beyond the first 340 lines). Super-admin dashboard: `GET /super-admin/dashboard/stats` (`handler/http/super_admin_handler.go:157`) reads tenant/subscription tables (platform metrics, outside tenant books; formulas UNVERIFIED).

### 4.7 KDS analytics - `GET /kds/analytics` (`handler/http/kds_handler.go`; repo KDS) - operational timing metrics, not financial (UNVERIFIED formulas).

## 5. HR reports (`handler/http/hr_report_handler.go`)
| Report | Endpoint | Formula (source tables) |
|---|---|---|
| Headcount | `/reports/hr/headcount` (:25) | COUNT of employees (filters), grouped by department, employment_status, employment_type |
| Turnover | `/reports/hr/turnover` (:246) | `turnover% = departures / ((headcountStart + headcountEnd)/2) * 100` from `employees.hire_date/termination_date` |
| Absenteeism | `/reports/hr/absenteeism` (:397) | `rate% = absent attendance records / total attendance records * 100`, by weekday/employee |
| Overtime | `/reports/hr/overtime` (:571) | per employee `hourly_rate = base_salary / (working_hours_per_week * 4.33)`; `overtime_cost = overtime_hours * hourly_rate * 1.5` (estimate from `attendance.overtime_hours`, NOT from salaries.overtime_amount - can disagree with payroll and the P&L) |
| Leave balances | `/reports/hr/leave-balances` (:694) | approved leave days and annual leave used per employee for the year (`l.status='approved'`) |
| Training compliance | `/reports/hr/training-compliance` (:801) | `certified employees / active employees * 100`, plus per department and record completion rate |
| Cost analysis | `/reports/hr/cost-analysis` | see 3.4 (served by accounting handler) |
| Social insurance | `/reports/hr/social-insurance` (perm EXPORT) | CSV from `GetSocialInsuranceReportRows` (accounting_advanced_repo.go:1421-1492) |

## 6. Accounting exports
`GET /accounting/csv|excel/chart-of-accounts|journal-entries/export` (repo `GetChartOfAccountsCSVRows` :1620, `GetJournalEntriesCSVRows` :1695; audit logged in `audit_export_logs`, RLS), imports at `/accounting/csv|excel/.../import`. `GET /accounting/customer-invoices/{id}/pdf`, `.../send-email`, customer statement `/accounting/customers/{id}/statement`. Format: CSV/XLSX/PDF; no CSV/PDF for financial statements.

## 7. Findings (ranked)
| # | Finding | Evidence |
|---|---|---|
| F1 | Menu Profitability and Break-even endpoints always fail: they select non-existent `products.cost_price` | accounting_advanced_repo.go:976,1045; `evidence/database/cost_price_probe.txt` |
| F2 | "Cash Flow" is not cash-based: operating = net P&L accrual, investing/financing by hardcoded account-code prefix and inverted signs | financial_report_usecase.go:388-392 |
| F3 | Balance Sheet cannot balance when period profit exists: P&L accounts excluded, no close/current-earnings policy | financial_report_usecase.go:304,344-351; current_earnings_diagnostic.go:13 |
| F4 | VAT report ignores GL and order-level VAT; hardcoded 10% order tax vs 14% CoA rate; period close is cosmetic | tax_repo.go:260-318; order_pricing.go:46 |
| F5 | Operational revenue reports (Sales, Profitability, order stats) include unpaid/unconfirmed orders and use current product cost; Advanced KPIs use paid orders; P&L uses posted journals - three different "revenue" numbers | report_handler.go:2174-2190,2364-2380; accounting_advanced_repo.go:837-850 |
| F6 | Inventory Valuation ignores as-of date, caches 5 min, values at standard cost not tied to GL 1130; Turnover reads a dead legacy table | financial_report_usecase.go:488-526; report_handler.go:2513 |
| F7 | Aging uses current balance_due rather than as-of; includes drafts | financial_report_usecase.go:438-455 |
| F8 | `round2` truncates negatives incorrectly; all report math in float64 | financial_report_usecase.go:528 |
| F9 | Break-even double counts COGS (journal expenses already contain 5110) | accounting_advanced_repo.go:1027-1103 |
| F10 | Overtime report estimates cost (x1.5, 4.33 wk/month) independent of payroll actuals | hr_report_handler.go:646-656 |
