# Advanced Reporting & Analytics (Draft v0)

Status: Approved V1 MVP Scope (Sales first implementation in progress)

Owner: Dashboard + Backend team

Primary Apps:
- `frontend/apps/dashboard` (`:3002`)
- `backend` (`:8080`)

Related Docs:
- `docs-organized/PRODUCTION-DEPLOYMENT-GUIDE.md`
- `docs-organized/04-MODULES/README.md`

## 1. Feature Goal (Business)

Create an advanced reporting and analytics feature for restaurant owners that turns orders, sales, and inventory data into actionable insights.

Target business outcomes (from product brief):
- Increase profitability (menu/product profitability analysis, cost control)
- Improve operations (waste reduction, staffing optimization)
- Speed up decisions (near real-time dashboards + trends)
- Fast ROI (reporting investment recovered quickly through optimizations)

## 1.1 Approved V1 MVP Decisions (Captured February 23, 2026)

Release scope:
- `V1 only` (fast launch, validate with real customers)

Mandatory V1 reports:
- Sales Analytics (daily/weekly/monthly, top items, sales per hour)
- Profitability per Item (Top 10 products)
- Inventory Turnover

Access control:
- RBAC `REPORTS` module with `READ` for viewing
- Export actions will use existing permission hierarchy (`WRITE`) in this codebase for V1 implementation

Important RBAC note (repo reality):
- The current codebase already seeds `REPORTS` in `modules` and maps it to `MODULE_IDS.REPORTS = 8`
- We will NOT renumber to `15` because module IDs come from `BIGSERIAL` seed order and renumbering would break existing permissions/data
- A safe migration will ensure the `REPORTS` module exists/metadata is correct by `name`

Profitability:
- Included in V1
- Product cost source confirmed conceptually by product table cost field (repo source uses `products.cost`, not `cost_price`)

Exports in V1:
- CSV + PDF
- JSON remains available through API responses

Realtime mode in V1:
- Polling every 30 seconds

Localization:
- Arabic + RTL support required in V1

Out of scope for V1:
- Scheduled reports
- Redis requirement
- Forecasting
- SSE/WebSocket live updates

## 2. What We Can Confirm In Current Codebase

### Existing building blocks already available

- Dashboard app and layout:
  - `frontend/apps/dashboard/src/components/layout/DashboardLayout.tsx`
  - `frontend/apps/dashboard/src/components/layout/Sidebar.tsx`
- Dashboard API client + auth/tenant headers:
  - `frontend/apps/dashboard/src/lib/api.ts`
- RBAC utility (already includes `REPORTS` module ID):
  - `frontend/apps/dashboard/src/lib/rbac.ts`
- Existing export utility pattern (PDF/CSV style approach reference):
  - `frontend/apps/dashboard/src/lib/pdf-export.ts`
- Orders / products / categories / drivers / HR pages (data shape references):
  - `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/page.tsx`
  - `frontend/apps/dashboard/src/app/[locale]/dashboard/products/page.tsx`
  - `frontend/apps/dashboard/src/app/[locale]/dashboard/categories/page.tsx`
  - `frontend/apps/dashboard/src/app/[locale]/dashboard/drivers/page.tsx`
  - `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/employees/page.tsx`
- Backend route registration pattern and wrappers:
  - `backend/cmd/api/main.go`
- Backend RBAC + permission middleware:
  - `backend/internal/middleware/permission.go`
- Backend order data repositories:
  - `backend/internal/repository/order_repo.go`
  - `backend/internal/handler/http/admin_order_handler.go`
- Existing notifications endpoint patterns (alerts integration path):
  - `backend/migrations/050_create_notifications_system.sql`
  - `frontend/apps/dashboard/src/hooks/useNotifications.ts`

### Important constraints / gaps found

- `REPORTS` RBAC module already exists in DB modules migration and dashboard constants:
  - `backend/migrations/100_create_modules_table.sql`
  - `frontend/apps/dashboard/src/lib/rbac.ts`
- No chart library currently installed in dashboard (`recharts` / `chart.js` not present):
  - `frontend/apps/dashboard/package.json`
- `order_items` has pricing fields but no persisted COGS/cost snapshot:
  - `backend/migrations/061_create_order_items_table.sql`
- Product cost exists in repo source as `products.cost` (nullable), not `cost_price`
- `order_items` does not persist a cost snapshot yet, so profitability is currently estimated from current product cost (acceptable for V1 with clear labeling)
- `inventory` ledger table exists but does not carry tenant/restaurant IDs directly; analytics must scope via `products`

## 2.1 Phase 1 Completion (V1 Foundations - Profitability + Inventory)

Completed audit decisions (implemented / documented):
- Profitability V1 uses `products.cost` as current cost basis (`COGS = SUM(order_items.quantity * products.cost)`)
- Inventory turnover V1 uses estimated COGS from sales plus average inventory value from:
  - opening stock (latest `inventory.quantity_after` before range, fallback to current stock)
  - current stock (`products.quantity_in_stock`)
- Inventory turnover output is normalized to monthly ratio:
  - `monthly_turnover_ratio = (estimated_cogs / avg_inventory_value) * (30 / days_in_range)`
- Slow-moving item heuristic in V1:
  - current stock > 0 and monthly turnover < 3.0, OR current stock > 0 and sold quantity = 0

Exact SQL design implemented in active backend handler (`backend/internal/handler/http/report_handler.go`):
- Sales analytics aggregates from `orders` + `order_items`
- Profitability aggregates from `order_items` + `orders` + `products`
- Inventory turnover aggregates from:
  - `products` (cost + current stock)
  - `order_items` + `orders` (estimated COGS / sold quantity)
  - `inventory` (opening quantity + waste quantity)

Indexes added for V1 reporting performance:
- `backend/migrations/121_add_reporting_indexes.sql`
  - orders: tenant/restaurant/date(+status)
  - order_items: tenant/restaurant/product/date
  - inventory: product/date(+reason)

## 3. Scope Recommendation (Initial Delivery)

### Recommended V1 (2-week target, realistic)

Focus on high-value reporting with current data + minimal schema risk:
- Sales Overview (daily/weekly/monthly)
- Orders Trends (count, revenue, average ticket)
- Sales by Hour (peak analysis)
- Sales by Product / Category (revenue + quantity)
- Basic customer analytics (repeat rate if customer identity is reliable)
- Inventory movement indicators (using stock + sales volume, not full COGS if unavailable)
- CSV export for reporting tables
- Dashboard Reports Overview page + drill-down pages
- RBAC-protected reports routes (`READ` / `EXPORT`)

### V1.5 (short follow-up)

- PDF summaries
- Scheduled emailed reports (daily/weekly)
- Alerts (threshold-based: low margin / high waste / food cost)
- Cached aggregated endpoints

### V2 (heavier / optional)

- Forecasting (seasonality-aware)
- Drag-and-drop custom dashboards/widgets
- SSE/WebSocket live ticker
- Accounting integrations (QuickBooks/Xero/Fawry/Daftra)
- BI export API for Power BI/Looker

## 4. Questions We Need You To Answer Before Finalizing The Plan

Status:
- Core scope/product answers received on February 23, 2026 (see section 1.1)
- Remaining questions below are optional refinements for later phases unless implementation hits a blocker

Answer these first so implementation scope, timelines, and schema changes are correct.

### A. Product & Scope (Must Decide)

1. Which release do you want first?
- `V1 only` (core reports + charts + CSV export)
- `V1 + alerts`
- `V1 + scheduled reports`
- `Everything in one phase` (not recommended)

2. For the first release, which report pages are mandatory?
- Sales overview
- Sales by hour
- Product/category performance
- Profitability
- Inventory turnover
- Customer analytics
- Labor efficiency
- Forecasting

3. Should this feature be available to:
- owner only
- owner + manager
- all roles with RBAC `REPORTS` permission

4. Do you want reports under:
- `dashboard/reports` only (recommended)
- dashboard home widgets + reports pages
- separate top-level analytics app/page

### B. Profitability / Cost Data (Critical)

5. Do you already store product cost (`cost_price`) anywhere in production DB?
- If yes: exact column/table name?
- If no: are you OK adding `products.cost_price` and optionally a per-order item cost snapshot?

6. For profitability, should COGS be:
- current product cost (simpler but historically inaccurate)
- cost snapshot captured at order time (recommended for accuracy)

7. Do you need VAT/tax separated in profit calculations from day 1?
- If yes: confirm tax logic and whether tax is included in `orders.total`.

### C. Customer Analytics (Data Reliability)

8. How should “unique customer” be identified?
- phone only
- email only
- email OR phone
- dedicated customer table ID (if available in all order paths)

9. Is guest checkout common?
- If yes, repeat customer metrics will need fuzzy matching rules.

### D. Inventory / Waste / Operations

10. Do you want inventory turnover based on:
- product stock value only
- inventory ledger transactions (`inventory` table)
- estimated COGS from sales (recommended)

11. Do you currently maintain purchase receipts / supplier costs in the system?
- If no, “waste/cost” analytics will be estimated and should be labeled clearly.

12. For labor efficiency, which source should be used for labor hours?
- attendance module only
- scheduled shifts (if exists)
- manual input

### E. Forecasting

13. Forecasting in first version: yes or no?
- If yes, choose:
  - simple moving average + seasonality factor
  - day-of-week seasonality baseline
  - external ML service later (V2)

14. Do you need special event/holiday forecasting (e.g., Ramadan) in V1?

### F. Exports / Integrations

15. Which export formats are required in first release?
- CSV (recommended V1)
- PDF
- Excel (`.xlsx`)
- JSON

16. “Export to Accounting” in first release: do you need:
- generic CSV only
- provider-specific mapping (QuickBooks/Xero/Fawry/Daftra)
- both

17. Scheduled reports delivery channel:
- email only
- in-app notifications only
- both

### G. Real-time / Refresh

18. For dashboard updates, what is acceptable in V1?
- manual refresh button only
- polling every 30–60 seconds (recommended)
- SSE/WebSocket (V2)

19. How “real-time” do you need it?
- near real-time (1–5 min)
- 30–60 sec
- instant

### H. UX / Dashboard Customization

20. Custom drag-and-drop widgets in first release?
- yes
- no (ship fixed layout first, recommended)

21. Save report filters per user?
- yes (recommended)
- no

22. Need Arabic labels/RTL support in reports V1?
- yes (recommended if dashboard already bilingual)
- no (EN first)

### I. Operational / Infra

23. Can we add Redis as a required dependency for production report caching?
- yes
- no (use DB/materialized views + in-memory cache first)

24. Are background workers allowed to run in the current backend service process?
- yes
- separate worker service preferred

25. Are there SLA/response-time targets for reporting endpoints?
- e.g., `<500ms cached`, `<3s uncached`

## 5. Proposed Technical Direction (Pending Answers)

### Backend

- Add a reports module (handler/usecase/repository) under existing architecture:
  - Handler parses filters and permissions
  - UseCase validates/filter policy and orchestrates aggregation
  - Repository executes SQL aggregation queries/views
- Start with DB aggregation queries and optional SQL views/materialized views
- Add export endpoints (CSV first)
- Add cache strategy for heavy report endpoints
- Add RBAC protection using existing `REPORTS` module (ID 8)

### Frontend Dashboard

- Add new reports pages under dashboard:
  - overview page
  - report detail pages
- Add reusable reporting components:
  - KPI cards
  - chart wrappers
  - filter bar
  - export buttons
  - data tables
- Add dedicated reports API hook/client
- Integrate in sidebar navigation

## 6. File Map (Likely Files To Edit / Create)

This is the key path map you requested. Final list may adjust after your answers.

### A. Backend - Existing Files Likely To Edit

- `backend/cmd/api/main.go`
  - Register report endpoints and route wrappers
- `backend/internal/middleware/permission.go`
  - Likely no core changes, but reference for wrapper usage
- `backend/internal/handler/http/admin_order_handler.go`
  - Reference for pagination/filter patterns
- `backend/internal/repository/order_repo.go`
  - Reference for order query patterns and joins
- `backend/internal/repository/product_repo.go`
  - May need cost field integration
- `backend/internal/repository/employee_repo.go`
  - Labor metrics input
- `backend/internal/repository/driver_repo.go`
  - Delivery efficiency metrics input
- `backend/internal/service/webhook_retry_worker.go`
  - No direct edits expected for V1, reference only if scheduled/export jobs share worker patterns

### B. Backend - New Files Likely To Create

- `backend/internal/domain/report.go`
- `backend/internal/domain/report_filters.go`
- `backend/internal/repository/report_repository.go`
- `backend/internal/repository/report_repository_impl.go`
- `backend/internal/usecase/report_usecase.go`
- `backend/internal/handler/http/report_handler.go`
- `backend/internal/service/report_export_service.go`
- `backend/internal/service/report_cache.go` (or reuse existing cache patterns)
- `backend/internal/service/report_schedule_worker.go` (V1.5+ if scheduled reports)
- `backend/internal/handler/http/report_handler_test.go`
- `backend/internal/usecase/report_usecase_test.go`
- `backend/internal/repository/report_repository_test.go` (optional if DB integration exists)

### C. Backend - Migrations / SQL Assets Likely To Create

If profitability and cost analytics are in scope:
- `backend/migrations/###_ensure_product_cost_column.sql` (only if an environment is missing `products.cost`)
- `backend/migrations/###_add_cost_snapshot_to_order_items.sql` (recommended for accurate historical margins; V1 uses estimated current cost)

Report support tables/views:
- `backend/migrations/###_create_reports_table.sql` (scheduled reports / saved configs / export history)
- `backend/migrations/###_create_report_jobs_table.sql` (if background scheduling is included)
- `backend/migrations/###_create_sales_analytics_views.sql`
- `backend/migrations/###_create_profitability_views.sql`
- `backend/migrations/###_create_inventory_analytics_views.sql`

### D. Frontend Dashboard - Existing Files Likely To Edit

- `frontend/apps/dashboard/src/components/layout/Sidebar.tsx`
  - Add Reports navigation
- `frontend/apps/dashboard/src/lib/rbac.ts`
  - `REPORTS` already exists; may add helper labels/guards only
- `frontend/apps/dashboard/src/lib/api.ts`
  - Usually no major changes; use existing client
- `frontend/apps/dashboard/src/lib/translations.ts`
  - Reference only (translator helper)
- `frontend/apps/dashboard/src/i18n/messages/en.json`
- `frontend/apps/dashboard/src/i18n/messages/ar.json`
- `frontend/apps/dashboard/src/components/layout/DashboardLayout.tsx`
  - Possibly page-specific full-width allowances for reports pages if needed

### E. Frontend Dashboard - New Files Likely To Create

Pages:
- `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/overview/page.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/page.tsx` (redirect or overview)
- `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/[type]/page.tsx`

API / hooks / types:
- `frontend/apps/dashboard/src/lib/api/reportsApi.ts`
- `frontend/apps/dashboard/src/hooks/useReports.ts`
- `frontend/apps/dashboard/src/hooks/useReportExport.ts`
- `frontend/apps/dashboard/src/types/reports.ts`

Components:
- `frontend/apps/dashboard/src/components/reports/ReportsFilterBar.tsx`
- `frontend/apps/dashboard/src/components/reports/KpiCard.tsx`
- `frontend/apps/dashboard/src/components/reports/ReportChart.tsx`
- `frontend/apps/dashboard/src/components/reports/ReportTable.tsx`
- `frontend/apps/dashboard/src/components/reports/ExportButton.tsx`
- `frontend/apps/dashboard/src/components/reports/DateRangePicker.tsx` (or reuse existing if found)
- `frontend/apps/dashboard/src/components/reports/SavedViewSelector.tsx` (V1.5/V2)
- `frontend/apps/dashboard/src/components/reports/AlertsPanel.tsx` (V1.5+)

Tests:
- `frontend/apps/dashboard/src/lib/api/reportsApi.test.ts`
- `frontend/apps/dashboard/src/hooks/useReports.test.tsx`
- `frontend/apps/dashboard/src/components/reports/*.test.tsx` (selective)

### F. Dependencies Likely To Add

Charts (choose one):
- `recharts` (recommended for V1 simplicity)
- or `chart.js` + `react-chartjs-2`

Optional later:
- `xlsx` (Excel export)
- drag-and-drop library for custom widgets (V2)

## 7. Data Model Risks / Decisions To Finalize

These are the main blockers for reliable profitability analytics:

1. Historical margin accuracy requires cost snapshot at order time (not yet persisted on `order_items`)
2. `products.cost` is nullable, so missing costs reduce profitability accuracy unless data quality is enforced
3. Customer identity consistency for repeat-rate metrics
4. Labor-hours source reliability for labor efficiency KPIs

If these are not finalized, we should still ship V1 with:
- sales/orders/category/product performance
- trend analytics
- operational dashboards
- export
- and clearly mark profitability as “estimated / pending cost data”

## 8. Suggested Delivery Strategy (Recommended)

Recommended approach:
- Ship V1 reporting (sales/orders/category/product + trends + CSV export + filters)
- Add profitability after cost model is confirmed
- Add forecasting + scheduled reports after data quality and caching are stable

Reason:
- Fastest business value with lowest schema risk
- Avoids blocking everything on cost accounting decisions

## 9. What I Need From You Next (Reply Checklist)

Please reply with answers (copy/paste and fill):

1. Release scope (`V1 only` / `V1+alerts` / `V1+scheduled` / `all`)
2. Mandatory reports for first release
3. Who can access reports (roles / RBAC policy)
4. Profitability in V1? (`yes/no`)
5. Do we have `cost_price` now? (where?)
6. If not, approve adding `products.cost_price`? (`yes/no`)
7. Export formats for V1 (`CSV/PDF/Excel/JSON`)
8. Real-time mode for V1 (`manual/polling/SSE`)
9. Arabic/RTL in V1 (`yes/no`)
10. Scheduled reports in V1 (`yes/no`)
11. Redis required for V1 (`yes/no`)
12. Forecasting in V1 (`yes/no`)
