# F-044 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 SQL): no trigger or code path updates kds_queue when an order is cancelled (pg_trigger on orders lists only order_stock_transition/order_stock_items); RUN-004 shows 3 cancelled orders still active in the queue (AUD-R4-KDS-after-read.json)


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/kds_repo.go`.
- Propagation: affected modules KDS; Orders; related findings F-028; observed consequence: Orders 34, 40, 42 are status cancelled but kds_queue rows stay status new and are listed in /kds/queue and on the board (ACTIVE ORDERS 25 / NEW 24 include them) | RUN-003 headed Chrome: 3 cancelled orders visible on KDS board (ACTIVE 25 / N
- Cause versus symptom: the symptom is what users or reports show (Orders 34, 40, 42 are status cancelled but kds_queue rows stay status new and are listed in /kds/queue and on the board (ACTIVE ORDERS 25 / ); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN005-DESTR-ORDER-CANCEL; scripts: tests/runtime/ph08_kds_cancelled.mjs, tests/runtime/ph04_r5_destructive_a.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/kds_repo.go | - | backend/internal/repository/kds_repo.go | exact | yes |
| backend/migrations/124_create_kds_queue.sql | - | backend/migrations/124_create_kds_queue.sql | exact | yes |
| backend/internal/repository/kds_repo.go | 307 | backend/internal/repository/kds_repo.go | exact | yes |

