# F-024 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | 623 | exact |
| frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | 637 | exact |
| frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx | 620-650 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
No frontend surface is involved or none could be derived.

## Database objects involved
None identified from the finding text.

## API contracts involved
- unchanged (values are rounded)

## Proposed (future) files
- frontend/apps/dashboard/src/lib/money.ts (new)
