# Phase 5 SQL Profiling Results (V1 Reports) - 2026-02-23

## Scope

Local profiling for the heaviest V1 report queries using:

- `EXPLAIN (ANALYZE, BUFFERS)` in PostgreSQL (Docker local)
- Tenant: `tenant_id=8`, `restaurant_id=8`
- Date range: `2026-01-24` to `2026-02-24`

This profiling was run after:

- Phase 1/2/4 reporting endpoints/pages implementation
- reporting indexes migration (`121_add_reporting_indexes.sql`)
- server-side date range guardrail (`max 366 days`)
- DB query timeout implementation in reports handler (`5s`)

## Local Dataset Context (Important)

Local database volume is small, so execution times are not representative of staging/production scale.

- `orders` (tenant 8 / restaurant 8): `6`
- `order_items` (tenant 8 / restaurant 8): `4`
- tracked `products` (tenant 8 / restaurant 8): `11`
- `inventory` rows (all tenants): `0`

Because of this, PostgreSQL will prefer some sequential scans that may change at higher row counts.

## Query Profiles (EXPLAIN ANALYZE)

### 1) Sales Trend Query (`/reports/sales`)

Raw output:
- `tmp/reports-sql-profiling/sales_trend_explain.txt`

Result:
- Planning Time: `3.108 ms`
- Execution Time: `0.278 ms`

Notes:
- Uses `Seq Scan on orders` (expected at tiny row count)
- `GroupAggregate` + sort cost is trivial at current volume

### 2) Top Items Query (`/reports/sales`)

Raw output:
- `tmp/reports-sql-profiling/top_items_explain.txt`

Result:
- Planning Time: `11.403 ms`
- Execution Time: `0.432 ms`

Notes:
- Joins `order_items -> orders -> products`
- Uses `Index Scan using products_pkey` for product join
- `orders` / `order_items` still use `Seq Scan` due tiny dataset
- This query is one of the key candidates for re-profiling at staging scale

### 3) Profitability Items Query (`/reports/profitability`)

Raw output:
- `tmp/reports-sql-profiling/profitability_items_explain.txt`

Result:
- Planning Time: `9.630 ms`
- Execution Time: `0.846 ms`

Notes:
- CTE + grouping on product/cost dimensions
- Similar join shape to top items query
- `products_pkey` index used
- Execution still sub-1ms on local dataset

### 4) Inventory Turnover Query (`/reports/inventory`)

Raw output:
- `tmp/reports-sql-profiling/inventory_turnover_explain.txt`

Result:
- Planning Time: `16.548 ms`
- Execution Time: `1.579 ms`

Notes:
- Heaviest V1 query (expected)
- Uses:
  - `idx_products_restaurant_tenant`
  - `idx_inventory_product_created`
- Includes LATERAL lookup for opening stock snapshot (`inventory` history)
- Current `inventory` table has `0` rows locally, so this must be re-profiled with real stock movement data

## Summary Table

| Query | Planning | Execution | Local Assessment |
|---|---:|---:|---|
| Sales Trend | 3.108 ms | 0.278 ms | Very fast |
| Top Items | 11.403 ms | 0.432 ms | Very fast (re-profile at scale) |
| Profitability Items | 9.630 ms | 0.846 ms | Very fast (re-profile at scale) |
| Inventory Turnover | 16.548 ms | 1.579 ms | Fast locally, heaviest query |

## API Sanity Check After Timeout Patch

Authenticated probes (tenant 8 / restaurant 8) returned `200` for:

- `/api/v1/reports/sales`
- `/api/v1/reports/profitability`
- `/api/v1/reports/inventory`
- `/api/v1/reports/overview`

## Frontend Route Sanity Check After `useReports` Refactor

Browser checks (headless Playwright) returned `200` and no page/console errors for:

- `/en/dashboard/reports/overview`
- `/en/dashboard/reports/sales`
- `/en/dashboard/reports/profitability`
- `/en/dashboard/reports/inventory`

This also fixed the reports overview flashing loop caused by `useReports` callback dependency churn.

## What Was Implemented In This Pass (Performance/Guardrails)

- DB query timeout for reports queries: `5s`
  - `backend/internal/handler/http/report_handler.go`
  - now uses `QueryContext` / `QueryRowContext`
- Export audit insert timeout also bounded (`ExecContext` with timeout)
- Date range guardrail already active (`max 366 days`)
- Last updated indicator now shown on all report pages (reduces stale-data confusion)

## Remaining Phase 5 Performance Work

- Re-profile on staging with realistic data volume (orders, order_items, inventory history)
- Capture cold vs warm endpoint timings (API-level baseline)
- Validate endpoint timing logs + cache-hit behavior with live backend runs
- Re-tune TTLs/timeouts only if staging timings require it

Status update:
- Endpoint timing logs are now implemented in `backend/internal/handler/http/report_handler.go`
- In-process TTL cache is now enabled for report GET endpoints
- Live timing capture completed successfully on February 23, 2026 after rebuilding and starting the Docker backend (`docker compose up -d --build backend`)
- Cache-hit behavior is confirmed by both:
  - stable `generated_at` values on repeat requests within TTL
  - backend timing logs showing `cache_hit=false` then `cache_hit=true`
- A prior failed attempt (`ConnectionRefusedError [WinError 10061]`) occurred before Docker/backend were available and is now resolved

## Cold vs Warm Timing Capture (Local / Staging)

Use the timing capture script to measure cold vs warm responses and infer cache hits from `generated_at` stability:

```bash
python scripts/capture_reports_cache_timings.py \
  --api-base http://localhost:8080/api/v1 \
  --jwt-secret your-super-secret-jwt-key-change-in-production-2024 \
  --tenant-id 8 \
  --restaurant-id 8 \
  --user-id 7 \
  --email owner@example.com \
  --repeat 2
```

For staging, prefer an existing real owner token:

```bash
python scripts/capture_reports_cache_timings.py \
  --api-base https://api-staging.yourdomain.com/api/v1 \
  --token "$REPORTS_TEST_TOKEN" \
  --date-from 2026-01-01 \
  --date-to 2026-01-31
```

Expected cache-hit signal:
- Request 1: cold (cache miss)
- Request 2 (same URL/filters within TTL): `generated_at` should remain unchanged (likely cache hit)

## Local Runtime Timing Results (Phase 5 Completion - 2026-02-23)

Environment:
- Backend: Docker Compose local (`pos-backend`, `pos-postgres`)
- API Base: `http://localhost:8080/api/v1`
- Tenant/Restaurant: `8 / 8`
- Repeat Count: `2` per endpoint/window
- DB Query Timeout: `5s`

### API Endpoint Timing Baseline (Cold vs Warm) - Local

| Endpoint | Filter Window | Cold ms | Warm ms | Cache Hit Verified | Notes |
|---|---|---:|---:|---|---|
| `/reports/overview` | 30d | 59.1 | 5.3 | Yes | `generated_at` unchanged; backend log `cache_hit=true` |
| `/reports/sales` | 30d | 16.2 | 17.6 | Yes | Warm request not faster in client wall-time due local jitter; backend log shows `cache_hit=true duration_ms=0` |
| `/reports/profitability` | 30d | 10.4 | 5.3 | Yes | |
| `/reports/inventory` | 30d | 9.3 | 6.2 | Yes | |
| `/reports/overview` | 90d | 48.2 | 4.9 | Yes | |
| `/reports/sales` | 90d | 18.0 | 7.1 | Yes | |
| `/reports/profitability` | 90d | 10.0 | 6.7 | Yes | |
| `/reports/inventory` | 90d | 10.1 | 5.9 | Yes | |
| `/reports/overview` | 365d | 62.8 | 18.0 | Yes | |
| `/reports/sales` | 365d | 20.9 | 4.8 | Yes | |
| `/reports/profitability` | 365d | 11.4 | 4.8 | Yes | |
| `/reports/inventory` | 365d | 13.6 | 4.7 | Yes | |

### Cache-Hit Log Evidence (Backend)

Sample backend log lines captured via:

```bash
docker logs pos-backend --since 3m 2>&1 | findstr /I "ReportHandler cache_hit"
```

Observed pattern (for all 4 endpoints across repeated runs):
- First request: `cache_hit=false`
- Second request (same filters): `cache_hit=true`

Example lines:
- `[ReportHandler] endpoint=overview ... cache_hit=false ...`
- `[ReportHandler] endpoint=overview ... cache_hit=true ...`
- `[ReportHandler] endpoint=sales ... cache_hit=true duration_ms=0 ...`

## Staging-Scale Bottlenecks & Index Coverage (Template)

Complete this section after running staging profiling with representative tenant volume.
Local timing verification is now complete; the table below remains staging-specific and should be filled before broader rollout.

### Staging Dataset Context

- Tenant/Restaurant profiled: `TODO`
- Orders in date range (`30d` / `90d` / `365d`): `TODO`
- Order items in date range: `TODO`
- Inventory movement rows in date range: `TODO`
- Tracked products count: `TODO`

### API Endpoint Timing Baseline (Cold vs Warm)

| Endpoint | Filter Window | Cold ms | Warm ms | Cache Hit Verified | Notes |
|---|---|---:|---:|---|---|
| `/reports/overview` | 30d | TODO | TODO | TODO | |
| `/reports/sales` | 30d | TODO | TODO | TODO | |
| `/reports/profitability` | 30d | TODO | TODO | TODO | |
| `/reports/inventory` | 30d | TODO | TODO | TODO | |
| `/reports/overview` | 90d | TODO | TODO | TODO | |
| `/reports/sales` | 90d | TODO | TODO | TODO | |
| `/reports/profitability` | 90d | TODO | TODO | TODO | |
| `/reports/inventory` | 90d | TODO | TODO | TODO | |
| `/reports/overview` | 365d | TODO | TODO | TODO | |
| `/reports/sales` | 365d | TODO | TODO | TODO | |
| `/reports/profitability` | 365d | TODO | TODO | TODO | |
| `/reports/inventory` | 365d | TODO | TODO | TODO | |

### SQL Bottlenecks (Staging Findings)

| Query | Symptom | Root Cause | Current Index Usage | Action |
|---|---|---|---|---|
| Sales Trend | TODO | TODO | TODO | TODO |
| Top Items | TODO | TODO | TODO | TODO |
| Profitability Items | TODO | TODO | TODO | TODO |
| Inventory Turnover | TODO | TODO | TODO | TODO |

### Index Coverage Review

Verify the following indexes are used (or explain why not) in staging `EXPLAIN ANALYZE`:

- `121_add_reporting_indexes.sql` indexes on:
  - `orders`
  - `order_items`
  - `inventory`
  - `products`
- Existing PK/FK indexes (`products_pkey`, etc.)

Checklist:
- [ ] `orders` reporting filters avoid large sequential scans at staging volume
- [ ] `order_items` joins are index-assisted for report queries
- [ ] `inventory_turnover` query uses `idx_inventory_product_created` (or better alternative)
- [ ] `products` tenant/restaurant index used in inventory query path
- [ ] Any remaining sequential scans are justified by row counts/selectivity

### TTL / Timeout Tuning Decision (Staging)

Current defaults:
- Overview/Sales TTL: `30s`
- Profitability/Inventory TTL: `60s`
- DB query timeout: `5s`

Record decisions:
- [x] Keep current TTLs/timeouts (local 30d/90d/365d timings all well within V1 targets; no timeout pressure observed)
- [ ] Adjust TTLs (document per-endpoint)
- [ ] Adjust DB query timeout (document reason)
- [ ] Add endpoint-specific timeout targets

## How To Re-run

1. Ensure backend + Postgres containers are up:

```bash
docker compose up -d
```

2. Run the `EXPLAIN ANALYZE` commands (examples already used) and save outputs under:

- `tmp/reports-sql-profiling/`

3. Compare new timings against this file and update if dataset size changes.
