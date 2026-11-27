# F-008 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/internal/repository/invoice_customer_repo.go | 452-636 | exact |
| backend/internal/repository/invoice_customer_repo.go | 452-560 | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices | exists |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting | exists |

## Database objects involved
`invoices_customer`

## API contracts involved
- POST /api/v1/accounting/customer-invoices: response journal_entry_id is null for drafts
- POST /api/v1/accounting/customer-invoices/{id}/issue: posts the journal (confirm exact route in 08-api-inventory.csv)

## Proposed (future) files
- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.
