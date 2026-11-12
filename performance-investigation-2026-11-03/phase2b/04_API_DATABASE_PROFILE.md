# API Latency, Response Size, and Query Work

25 sequential authenticated GET samples per admin-order page size against the isolated synthetic database. Warm latency excludes the first sample. Response size is observed JSON bytes. Component spans were unavailable.

| Method / route | Context | Size | Status | Bytes | Records | Warm median ms | Warm p95 ms | First ms | n |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| GET /api/v1/admin/orders | synthetic tenant/restaurant | 1 | 200 | 417 | 1 | 12.9 | 15.1 | 349.9 | 25 |
| GET /api/v1/admin/orders | synthetic tenant/restaurant | 5 | 200 | 1,859 | 5 | 20.2 | 24.2 | — | 25 |
| GET /api/v1/admin/orders | synthetic tenant/restaurant | 25 | 200 | 9,055 | 25 | 50.1 | 61.2 | — | 25 |
| GET /api/v1/admin/orders | synthetic tenant/restaurant | 50 | 200 | 18,055 | 50 | 91.2 | 117.4 | — | 25 |
| GET /api/v1/super-admin/audit-logs | synthetic Super Admin | page 1 size 50 | 200 | unmeasured | 9 | unmeasured | unmeasured | — | 1 observed |

P75 and p99 are omitted because the retained aggregate contains median and p95 only. Handler/SQL/Redis timings, runtime query counts, error rates, and production-scale samples are unavailable. Source shows the admin list/count plus one GetOrder per returned row; each GetOrder reads order and items. This is a static N+1 candidate supported by rising end-to-end latency with page size, not a measured SQL query-count curve.

Synthetic PostgreSQL plans: order list uses idx_orders_tenant_restaurant_created (0.933 ms), count 0.295 ms, item lookup uses idx_order_items_order_id (0.165 ms). At 100 orders per restaurant, no index recommendation is justified. pg_stat_statements is unavailable.
