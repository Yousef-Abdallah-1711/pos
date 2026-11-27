# F-019 — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
| path | lines | resolution |
|---|---|---|
| backend/migrations/164_rls_user_auth_audit_logs.sql | - | exact |
| backend/migrations/171_rls_e_invoice_submissions.sql | - | exact |

## Frontend (existing, module-derived — verify the exact component before editing)
| path | status |
|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/orders | exists |

## Database objects involved
`orders`, `journal_entry_lines`, `payments_received`, `tenants`, `users`

## API contracts involved
- unchanged

## Proposed (future) files
- backend/migrations/<next>_rls_core_money_tables.sql (new)
- backend/internal/middleware/tenant_guc.go (new, proposed)
