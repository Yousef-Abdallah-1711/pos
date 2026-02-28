# V1 Reports Verification Results (February 23, 2026)

Scope verified:
- `Sales`
- `Profitability`
- `Inventory`
- `Overview`
- Exports (`CSV`, `PDF`)

Environment:
- Backend: `http://localhost:8080`
- Dashboard: `http://localhost:3002`
- Tenant context used: `tenant_id=8`, `restaurant_id=8`, `user_id=7`

## 1. Manual API Tests (Real Tenant Data)

Verified endpoints (authenticated with local JWT):
- `GET /api/v1/reports/sales`
- `GET /api/v1/reports/profitability`
- `GET /api/v1/reports/inventory`
- `GET /api/v1/reports/overview`

Result:
- All returned `200`
- All returned `{ success: true, data: ... }`
- All included `generated_at`

Observed sample data (local DB):
- Sales (`2026-02-01` → `2026-02-23`): `937.2 EGP`, `2 orders`
- Shorter range (`2026-02-22` → `2026-02-23`): `0 EGP`, `0 orders`

Conclusion:
- Filters are affecting the result set as expected.

## 2. Tenant Scope / Cross-Tenant Leakage Checks

Checks performed:
1. Header override probe:
   - Sent valid tenant-8 JWT with conflicting headers `X-Tenant-ID=999`, `X-Restaurant-ID=999`
   - Result matched tenant-8 response
   - Confirms JWT tenant context takes precedence over headers

2. Cross-tenant token probe:
   - Sent JWT for `tenant_id=999`, `restaurant_id=999`
   - Sales summary returned zeros (no data), not tenant-8 data

Conclusion:
- No cross-tenant leakage observed in V1 report endpoints.

## 3. Export API Verification (Backend)

Verified endpoint:
- `POST /api/v1/reports/export`

Tested combinations:
- `sales` + `csv`
- `sales` + `pdf`
- `profitability` + `csv`
- `profitability` + `pdf`
- `inventory` + `csv`
- `inventory` + `pdf`

Checks:
- `200 OK`
- Correct `Content-Type`
- Correct `Content-Disposition` filename
- CSV payload non-empty
- PDF payload starts with `%PDF`

CSV compatibility notes:
- UTF-8 BOM added (Excel compatibility)
- CRLF line endings enabled (`csv.Writer.UseCRLF = true`)
- CSV parsed successfully in validation script

Note:
- Actual Excel/Google Sheets UI application-open test was not run inside this CLI environment, but generated CSV format is spreadsheet-compatible and parses correctly.

## 4. UI Verification (Dashboard / Playwright)

Routes verified:
- `/en/dashboard/reports/overview`
- `/en/dashboard/reports/sales`
- `/en/dashboard/reports/profitability`
- `/en/dashboard/reports/inventory`

Results:
- All routes returned `200`
- Reports loaded with authenticated local session injected via Playwright (`auth-storage`)

UI filter behavior (Sales page):
- Changed `date_to` in the date filter
- Observed follow-up request:
  - from `...date_to=2026-02-23`
  - to `...date_to=2026-02-22`

UI export behavior (Playwright downloads):
- `Sales`: CSV + PDF download succeeded
- `Profitability`: CSV + PDF download succeeded
- `Inventory`: CSV + PDF download succeeded

## 5. Bugs Found & Fixed During Verification

Fixed backend SQL bug:
- Reports queries used `GROUP BY` aliases (`product_name`, `category_name`) causing PostgreSQL errors in real execution.
- Patched to group by the underlying expressions.

Added export quality improvements:
- UTF-8 BOM + CRLF for CSV export (better Excel compatibility)

Added server-side safety:
- Date range guardrail (`max 366 days`) in report date parsing

## 6. Remaining Verification/Performance Work (Phase 5+)

Still pending:
- Live cold vs warm endpoint timing capture (API-level baseline with backend running)
- Cache-hit behavior confirmation from backend timing logs (`cache_hit=true`)
- Target response-time validation for common filters (30d/90d/365d)
- Large export load test (sync-only V1 export; async flow deferred)

Completed after this verification pass:
- DB query timeout enforcement (reports queries use `QueryContext`/`QueryRowContext`)
- Shared reports UI component refactor (`ReportChart` / `ReportTable`) across V1 report pages
- Sync vs async export decision documentation (V1 sync-only; async deferred)
