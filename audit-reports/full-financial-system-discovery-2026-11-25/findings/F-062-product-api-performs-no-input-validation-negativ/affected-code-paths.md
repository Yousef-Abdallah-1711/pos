# F-062 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/usecase/product_usecase.go | 294-306 | exact |
| backend/internal/handler/http/product_handler.go | 53 | exact |
| backend/internal/handler/http/supplier_handler.go | 97 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/products | exists |

## Database objects involved
`products`, `suppliers`

## API contracts involved
- POST/PUT /api/v1/products, categories, customers, suppliers: 400/409 with field errors

## Proposed (future) files
- backend/migrations/<next>_product_check_constraints.sql (new)
