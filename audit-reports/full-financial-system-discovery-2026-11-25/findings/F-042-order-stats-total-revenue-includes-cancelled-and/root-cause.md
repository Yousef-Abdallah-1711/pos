# F-042 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/order_repo.go:859-870 GetOrderStats sums total_amount with no status/payment filter and computes delivered_today as COUNT(created_at >= CURRENT_DATE)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/order_repo.go | 859-950 | backend/internal/repository/order_repo.go | exact | yes |
| backend/internal/repository/order_repo.go | 859-870 | backend/internal/repository/order_repo.go | exact | yes |

