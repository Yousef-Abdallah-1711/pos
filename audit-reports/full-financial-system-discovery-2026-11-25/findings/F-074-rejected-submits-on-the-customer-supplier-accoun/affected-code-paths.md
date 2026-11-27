# F-074 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | 129 | exact |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx | 95 | exact |
| frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx | 69 | exact |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx | 18 | exact |
| frontend/apps/dashboard/src/components/hr-advanced/ContractForm.tsx | 27 | exact |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx | 95-118 | exact |
| frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx | 69-79 | prefix |
| frontend/apps/dashboard/src/components/accounting/SupplierForm.tsx | 75-85 | basename |
| frontend/apps/dashboard/src/components/accounting/AccountForm.tsx | 81-117 | basename |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx | 18-32 | prefix |
| frontend/apps/dashboard/src/components/hr-advanced/ContractForm.tsx | 27-29 | prefix |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |

## Database objects involved
`suppliers`

## API contracts involved
- error bodies gain field_errors (additive, see F-025)

## Proposed (future) files
- frontend/apps/dashboard/src/hooks/useFormSubmit.ts (new)
