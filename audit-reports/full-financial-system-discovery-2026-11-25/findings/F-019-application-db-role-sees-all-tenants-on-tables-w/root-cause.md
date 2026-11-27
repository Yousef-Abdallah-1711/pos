# F-019 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 runtime SQL): only 13 of the 70 public tables with a tenant_id column have row-level security (RLS migrations 164-171 are the only ones); the application role pos_app has UPDATE on orders and payments and no immutability trigger exists on payments (orders carry only stock-transition triggers); journal_entry_lines is not updatable by pos_app


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/migrations/164_rls_user_auth_audit_logs.sql | - | backend/migrations/164_rls_user_auth_audit_logs.sql | exact | yes |
| backend/migrations/171_rls_e_invoice_submissions.sql | - | backend/migrations/171_rls_e_invoice_submissions.sql | exact | yes |

