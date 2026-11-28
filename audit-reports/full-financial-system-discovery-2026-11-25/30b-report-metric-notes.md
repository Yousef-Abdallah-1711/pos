# 30b — Report metric inventory notes

File: [30-report-metric-inventory.csv](30-report-metric-inventory.csv) — **112 metric rows** across 22 report groups (PNL 6, TB 5, BS 5, CF 4, GL 3, Aging 11, VAT 4, Inventory valuation 3, KPI 6, Menu profitability 5, Break-even 5, Predictive CF 3, HR cost 6, Sales analytics 8, Profitability 7, Inventory turnover 9, Overview 1, Order stats 3, Dashboard home 4, KDS 3, HR reports 7, Super Admin 4).

## How it was built
From the formulas already catalogued in [12](12-report-formula-catalog.md) (which cites file:line), re-read this session for the structure; per-report defaults (API, source, date basis, cache, export) + per-metric formula/source/conflicts. Endpoints were **not** counted as metrics.

## Limitations (honest)
1. **KDS (3 rows) and Super Admin (4 rows) are placeholders**: their formulas were never read (evidence label UNKNOWN). The labels are indicative, not verified display labels. Real count for those reports is UNKNOWN. Excluding them, the verified-formula count is **105**.
2. Display labels for other reports are taken from the formula/field names in the backend responses, **not from the rendered UI text or the i18n files** (en/ar). A UI-label pass during execution is required.
3. Dashboard home (4 rows) was read only to line 340 of the page.
4. HR reports: six endpoints listed with one row each; per-field output metrics (e.g. per-department breakdowns) are not individually enumerated.
5. Charts/tables/exports: exports exist only for Sales/Profitability/Inventory (CSV+PDF, `POST /reports/export`), HR social insurance (CSV) and accounting CSV/Excel exports; financial statements (P&L, TB, BS, CF) have **no export** (JSON only). Chart components were not separately enumerated.
6. Evidence labels are CONFIRMED-BY-SOURCE unless stated; no metric has been checked at runtime.
