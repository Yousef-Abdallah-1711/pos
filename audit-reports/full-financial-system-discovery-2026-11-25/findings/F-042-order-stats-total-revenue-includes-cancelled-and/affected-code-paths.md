# F-042 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/order_repo.go | 859-950 | exact |
| backend/internal/repository/order_repo.go | 859-870 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`orders`

## API contracts involved
- GET /api/v1/admin/orders/stats: adds booked_revenue and collected_revenue (additive); total_revenue follows POL-01

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
