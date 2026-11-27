# F-044 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 SQL): no trigger or code path updates kds_queue when an order is cancelled (pg_trigger on orders lists only order_stock_transition/order_stock_items); RUN-004 shows 3 cancelled orders still active in the queue (AUD-R4-KDS-after-read.json)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/kds_repo.go | - | backend/internal/repository/kds_repo.go | exact | yes |
| backend/migrations/124_create_kds_queue.sql | - | backend/migrations/124_create_kds_queue.sql | exact | yes |
| backend/internal/repository/kds_repo.go | 307 | backend/internal/repository/kds_repo.go | exact | yes |

