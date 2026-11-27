# F-042 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/order_repo.go:859-870 GetOrderStats sums total_amount with no status/payment filter and computes delivered_today as COUNT(created_at >= CURRENT_DATE)


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/order_repo.go`:859-950.
- Propagation: affected modules Orders dashboard; Order stats API; related findings none recorded; observed consequence: total_revenue 529,100 = SUM of ALL 41 orders incl 3 cancelled (11,440) and 21 unpaid; delivered_today 30 while delivered_count 0 (counts any order created today) | RUN-003: order stats total_revenue 529,100 = all 41 orders incl 3 cancelled;
- Cause versus symptom: the symptom is what users or reports show (total_revenue 529,100 = SUM of ALL 41 orders incl 3 cancelled (11,440) and 21 unpaid; delivered_today 30 while delivered_count 0 (counts any); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: tests/runtime/ph12_batch_apis.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/order_repo.go | 859-950 | backend/internal/repository/order_repo.go | exact | yes |
| backend/internal/repository/order_repo.go | 859-870 | backend/internal/repository/order_repo.go | exact | yes |

