# PostgreSQL Query and Index Inventory

This pass maps selected source hotspots and runtime metadata. It does not claim to inventory every repository query or establish production bottlenecks. Static call-site grep found roughly 375 direct `Query`/`QueryRow`/`Exec` versus 211 context-aware variants in non-test repository files; this is a call-site count, not unique SQL or runtime use.

| Query family | Source path | Static observations | Runtime plan |
|---|---|---|---|
| Admin order list + detail hydration | `backend/internal/handler/http/admin_order_handler.go:108–225`; `internal/repository/order_repo.go:254–390` | List is followed by per-order `GetOrder`; OFFSET pagination; request page size has no evident maximum clamp | EXPLAIN / per-query timing not captured |
| Admin order export | `admin_order_handler.go:~609–675` | Up to 1,000 list records followed by per-record fetch | Not executed |
| Product list | `product_handler.go:179–240`; `product_repo.go:~285–306` | Handler defaults page size 20 and caps 100; count + paged query; SQL sort field/order interpolated | Not executed; validate allowlist and plan |
| KDS repository | KDS repository call sites | Selected DB duration instrumentation exists | No endpoint-correlated plan sampled |
| Financial reports | report handler + repositories | Redis cache can skip report work; query plan varies by report/date scope | No hit-path / miss-path plan captured |

## Safe runtime metadata

Disposable E2E PostgreSQL is version 15.19, database `pos_phase2a_e2e`, about 18.4 MB, 91 public tables, low synthetic tuple counts. `pg_stat_statements` is absent; no useful historical query fingerprints. No EXPLAIN ANALYZE was run. The sparse database does not model production cardinalities and is insufficient to rank indexes.

## Index recommendation discipline

No specific new index is recommended without query plan and realistic row-count evidence. For the admin orders N+1 pattern, reduce round trips first and inspect existing tenant/restaurant/date indexes; adding indexes does not eliminate per-row round trips and adds write cost. For product sort, use allowlisted identifiers before plan review.
