# F-080 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx | 61 | exact |
| backend/internal/handler/http/customer_acct_handler.go | 141 | exact |
| backend/internal/repository/journal_entry_repo.go | 708 | exact |
| frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx | 61-62 | basename |

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
- POST/PUT /api/v1/accounting/customers and /api/v1/suppliers: 400 for invalid account ids (additive validation)

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
