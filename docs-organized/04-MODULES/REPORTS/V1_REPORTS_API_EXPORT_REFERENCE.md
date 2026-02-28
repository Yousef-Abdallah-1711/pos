# V1 Reports API & Export Reference (Dashboard / Backend)

Status: V1 core reports + exports implemented (February 23, 2026)

Applies to:
- Backend: `backend/internal/handler/http/report_handler.go`
- Dashboard: `frontend/apps/dashboard/src/lib/api/reportsApi.ts`

## 1. Implemented V1 Endpoints

All endpoints are tenant/restaurant-scoped via auth + tenant middleware and return the standard shape:

```json
{ "success": true, "data": { ... } }
```

### GET `/api/v1/reports/overview`

Query params:
- `date_from=YYYY-MM-DD` (optional)
- `date_to=YYYY-MM-DD` (optional)

Example:

```http
GET /api/v1/reports/overview?date_from=2026-02-01&date_to=2026-02-23
Authorization: Bearer <jwt>
```

Response shape (summary):

```json
{
  "success": true,
  "data": {
    "filters": {
      "date_from": "2026-02-01",
      "date_to": "2026-02-23"
    },
    "sales": {
      "summary": {
        "total_revenue": 937.2,
        "total_orders": 2,
        "avg_ticket": 468.6
      },
      "top_item": {
        "product_id": 3,
        "name": "Grilled Chicken",
        "quantity": 5,
        "revenue": 74.95
      }
    },
    "profitability": {
      "summary": {
        "total_revenue": 937.2,
        "total_cogs": 0,
        "gross_profit": 937.2,
        "gross_margin_pct": 100,
        "food_cost_pct": 0
      }
    },
    "inventory": {
      "summary": {
        "tracked_products_count": 11,
        "average_monthly_turnover": 0
      }
    },
    "generated_at": "2026-02-23T17:00:00Z"
  }
}
```

### GET `/api/v1/reports/sales`

Query params:
- `date_from=YYYY-MM-DD` (optional)
- `date_to=YYYY-MM-DD` (optional)
- `granularity=day|week|month` (optional; default `day`)
- `top_limit` (optional; capped at `50`)

Example:

```http
GET /api/v1/reports/sales?date_from=2026-02-01&date_to=2026-02-23&granularity=day&top_limit=10
Authorization: Bearer <jwt>
```

Response includes:
- `summary`
- `trend`
- `hourly`
- `top_items`
- `top_categories`
- `generated_at`

### GET `/api/v1/reports/profitability`

Query params:
- `date_from=YYYY-MM-DD` (optional)
- `date_to=YYYY-MM-DD` (optional)
- `top_limit` (optional; capped at `50`)

Example:

```http
GET /api/v1/reports/profitability?date_from=2026-02-01&date_to=2026-02-23&top_limit=10
Authorization: Bearer <jwt>
```

V1 profitability notes:
- Uses `products.cost` as current cost basis (`COGS`)
- Historical cost snapshots are **not** implemented in V1
- Response includes `cost_field: "products.cost"` in filters

Response includes:
- `summary` (`gross_margin_pct`, `food_cost_pct`, missing-cost counters)
- `top_items`
- `lowest_margin_items`
- `generated_at`

### GET `/api/v1/reports/inventory`

Query params:
- `date_from=YYYY-MM-DD` (optional)
- `date_to=YYYY-MM-DD` (optional)
- `top_limit` (optional; capped at `50`)

Example:

```http
GET /api/v1/reports/inventory?date_from=2026-02-01&date_to=2026-02-23&top_limit=10
Authorization: Bearer <jwt>
```

V1 inventory notes:
- Estimated COGS uses `products.cost`
- Turnover is normalized to monthly ratio based on selected date range
- Slow-moving heuristic is threshold-based (`<3x/month` or stock with zero sales)

Response includes:
- `summary`
- `items`
- `slow_movers`
- `top_turnover`
- `generated_at`

## 2. Export Endpoint (Implemented in V1)

### POST `/api/v1/reports/export`

Supported report types (V1):
- `sales`
- `profitability`
- `inventory`

Supported formats (V1):
- `csv`
- `pdf`

RBAC:
- `REPORTS` + `EXPORT` permission (repo module id is `8`)

### Request Payload Schema

```json
{
  "report_type": "sales",
  "format": "csv",
  "date_from": "2026-02-01",
  "date_to": "2026-02-23",
  "granularity": "day",
  "top_limit": 10
}
```

Notes:
- `granularity` is used by `sales`
- `top_limit` is capped at `50`
- `date_from` / `date_to` default to V1 server defaults if omitted

### Response (CSV)

Headers:
- `Content-Type: text/csv; charset=utf-8`
- `Content-Disposition: attachment; filename="...csv"`

Compatibility:
- UTF-8 BOM included (Excel-friendly)
- CRLF line endings enabled

### Response (PDF)

Headers:
- `Content-Type: application/pdf`
- `Content-Disposition: attachment; filename="...pdf"`

Payload:
- Summary PDF (`%PDF` prefix)

## 3. V1 Accounting Export Schema (Generic CSV Mapping)

V1 decision:
- Use the existing report CSV export endpoint as the accounting-compatible base
- No provider-specific adapters (QuickBooks/Xero/Fawry/Daftra) in V1

### Generic Accounting CSV Columns (Recommended V1 Mapping)

Use these columns when exporting/importing into accounting workflows:

| Column | Source | Notes |
|---|---|---|
| `Date` | report bucket/date | Daily/period summary row date |
| `Revenue` | sales/profitability totals | Gross sales revenue |
| `Tax` | `0` (placeholder in V1 unless explicit tax split available) | V1 tax split not modeled consistently |
| `COGS` | profitability/inventory estimated COGS | Uses `products.cost` in V1 |
| `VAT` | `0` (placeholder in V1 unless explicit VAT split available) | Add in future tax-aware version |
| `Orders` | sales summary/trend | Optional but useful for reconciliation |
| `AverageTicket` | sales summary | Optional |
| `RestaurantID` | tenant context | For downstream mapping |
| `TenantID` | tenant context | For downstream mapping |

Important V1 caveat:
- `Tax`/`VAT` columns are placeholders unless order tax is modeled/reconciled consistently for the tenant.

## 4. Large Export Behavior (V1 Decision)

Decision (documented):
- **V1 uses synchronous exports only** via `POST /api/v1/reports/export`
- **Async export flow** (`GET /api/v1/reports/export/{id}` and background generation) is deferred to V1.5/V2

Current safeguards:
- `top_limit` capped at `50`
- Date range guardrail enforced (`max 366 days`)
- DB query timeout applied in reports queries (`5s`)

Pending validation (Phase 5/10):
- Large export load test with staging-scale data
- Sync vs async revisit if timeouts or payload sizes become problematic

## 5. Non-V1 / Deferred Endpoints (Documented)

Deferred (not part of approved V1 MVP):
- `GET /api/v1/reports/sales-by-hour` (sales hourly is included inside `/reports/sales`)
- `GET /api/v1/reports/products`
- `GET /api/v1/reports/categories`
- `GET /api/v1/reports/customers`
- `GET /api/v1/reports/labor-efficiency`
- `GET /api/v1/reports/export/{id}` (async export retrieval)

## 6. Security / Tenant Isolation Notes

- Tenant/restaurant context comes from authenticated JWT + tenant middleware
- Endpoint filters do not override tenant scope
- Cross-tenant leakage checks were completed in:
  - `docs-organized/04-MODULES/REPORTS/V1_REPORTS_VERIFICATION_RESULTS_2026-02-23.md`

