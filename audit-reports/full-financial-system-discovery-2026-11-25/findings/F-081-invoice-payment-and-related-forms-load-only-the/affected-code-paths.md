# F-081 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/handler/http/customer_acct_handler.go | 29 | exact |
| backend/internal/handler/http/supplier_handler.go | 29 | exact |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx | 120 | exact |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices/page.tsx | 107 | exact |
| frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | 107 | basename-ambiguous(117) |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |

## Database objects involved
`suppliers`

## API contracts involved
- unchanged (search and page_size parameters already supported or additive)

## Proposed (future) files
- frontend/apps/dashboard/src/components/accounting/EntityCombobox.tsx (new)
