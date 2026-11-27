# F-021 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/supplier_repo.go | 44-189 | exact |
| backend/internal/repository/payment_acct_repo.go | 302 | exact |
| backend/internal/repository/supplier_repo.go | 172 | exact |
| backend/internal/handler/http/supplier_handler.go | 73 | exact |
| backend/internal/repository/payment_acct_repo.go | 302-306 | exact |
| backend/internal/repository/supplier_repo.go | 172-200 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers | exists |

## Database objects involved
`suppliers`

## API contracts involved
- supplier responses keep current_balance (value now derived)

## Proposed (future) files
- backend/migrations/<next>_supplier_balance_view.sql (new)
