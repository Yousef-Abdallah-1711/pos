# F-052 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/payment_acct_handler.go | - | exact |
| backend/internal/repository/payment_acct_repo.go | - | exact |
| backend/internal/handler/http/order_cancellation_handler.go | 36 | exact |
| backend/internal/handler/http/order_cancellation_handler.go | - | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |

## Database objects involved
`payments_received`

## API contracts involved
- All money-creating POSTs: header Idempotency-Key (UUID) required; documented in 08-api-inventory.csv

## Proposed (future) files
- backend/internal/middleware/idempotency.go (new)
- backend/migrations/<next>_create_idempotency_keys.sql (new)
