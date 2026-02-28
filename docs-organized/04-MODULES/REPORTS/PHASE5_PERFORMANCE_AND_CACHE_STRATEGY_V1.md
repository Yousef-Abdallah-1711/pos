# Phase 5 - Performance & Cache Strategy (V1)

Status: Complete for V1 local baseline (February 23, 2026) - staging-scale re-profiling recommended before wider rollout

This document captures the Phase 5 implementation, local runtime verification, and remaining staging-scale profiling tasks for V1 reports.

## Implemented in this Phase 5 Start

- Shared frontend reports components created:
  - `frontend/apps/dashboard/src/components/reports/KpiCard.tsx`
  - `frontend/apps/dashboard/src/components/reports/ReportChart.tsx`
  - `frontend/apps/dashboard/src/components/reports/ReportTable.tsx`
  - `frontend/apps/dashboard/src/components/reports/ExportButton.tsx`
  - `frontend/apps/dashboard/src/components/reports/LastUpdatedIndicator.tsx`
- Shared polling/data-fetch hook created:
  - `frontend/apps/dashboard/src/hooks/useReports.ts`
- Shared reports wrappers adopted across all V1 detail pages:
  - `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/sales/page.tsx`
  - `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/profitability/page.tsx`
  - `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/inventory/page.tsx`
- Reports sidebar grouped with child pages:
  - Overview / Sales / Profitability / Inventory
- Server-side date-range guardrail added:
  - max range = `366` days
- DB query timeout added in reports handler:
  - `5s` timeout using `QueryContext` / `QueryRowContext`
- In-process backend cache enabled for report GET endpoints (TTL-based)
- Endpoint timing logs added for report GET endpoints (`cache_hit`, `duration_ms`, tenant/restaurant)
- Local SQL profiling completed and documented:
  - `docs-organized/04-MODULES/REPORTS/PHASE5_SQL_PROFILING_RESULTS_2026-02-23.md`
- Cold vs warm timing capture helper added and executed successfully (30d/90d/365d local):
  - `scripts/capture_reports_cache_timings.py`

## V1 Cache Strategy (Implemented + Documented)

Redis is out of scope for V1, so cache strategy is intentionally lightweight:

### Frontend (Dashboard)

- Polling interval: `30s`
- Per-page in-memory React state only (no persistence)
- Display both:
  - server `generated_at`
  - client `last updated` (dashboard refresh time)

### Backend (Implemented Lightweight V1 Cache)

In-process TTL cache is enabled for:
- `GET /api/v1/reports/overview` -> TTL `30s`
- `GET /api/v1/reports/sales` -> TTL `30s`
- `GET /api/v1/reports/profitability` -> TTL `60s`
- `GET /api/v1/reports/inventory` -> TTL `60s`

Cache key dimensions:
- `tenant_id`
- `restaurant_id`
- `endpoint`
- `date_from`
- `date_to`
- `granularity`
- `top_limit`

Cache invalidation (V1 simple rule, implemented):
- TTL-based only (no active invalidation)
- UI should always show “last updated” timestamp to avoid stale-data confusion

## SQL Profiling Plan (Next Performance Pass)

Local baseline has now been captured (SQL + API cold/warm timings). Next pass should be on staging with higher data volume.
Use `scripts/capture_reports_cache_timings.py` to capture API cold/warm timings before and after cache hits.

Run against local/staging with representative tenant data.

### Endpoints to profile first

1. `sales`
2. `profitability`
3. `inventory`
4. `overview` (composite endpoint)

### Measurements to capture

- Response time (cold/warm app)
- Query execution plan (`EXPLAIN ANALYZE`)
- Rows scanned vs rows returned
- Index usage
- Date range impact (`30d`, `90d`, `365d`)

### SQL profiling examples (psql)

Use the exact SQL from `backend/internal/handler/http/report_handler.go` and run:

```sql
EXPLAIN (ANALYZE, BUFFERS)
SELECT ...
```

Track:
- sequential scans on `orders`, `order_items`, `inventory`
- join selectivity on `tenant_id`, `restaurant_id`
- whether reporting indexes from `121_add_reporting_indexes.sql` are used

## Query Timeout / Guardrails Plan

Implemented:
- max date range guardrail (`366` days)
- Per-request DB query timeout using `QueryContext/QueryRowContext` (`5s`)

Current V1 decision after local timing verification:
- Keep `5s` DB query timeout for report GET endpoints
- Keep current endpoint TTLs (`30s` overview/sales, `60s` profitability/inventory)
- Revisit only if staging profiling shows timeouts or heavy scans

## Performance Baseline Checklist (V1 Local Baseline)

- [x] Capture response times for all report endpoints (`30d`, `90d`, `365d`) - local timing capture completed
- [x] Capture `EXPLAIN ANALYZE` for top heavy queries (local baseline)
- [x] Document bottlenecks and index coverage (local baseline + staging template/checklist documented in `PHASE5_SQL_PROFILING_RESULTS_2026-02-23.md`)
- [x] Decide if in-process cache is needed before V1 release (decision: yes, lightweight in-process TTL cache)
- [x] Add endpoint performance logging (report GET endpoints)
- [x] Add DB query timeouts

## Remaining Before Wider Rollout (Non-blocking for local V1 completion)

- Run staging-scale timing capture (`30d/90d/365d`) with representative tenant volume
- Fill staging bottlenecks and index-coverage findings in `PHASE5_SQL_PROFILING_RESULTS_2026-02-23.md`
- Re-tune TTLs/timeouts only if staging timings justify changes
