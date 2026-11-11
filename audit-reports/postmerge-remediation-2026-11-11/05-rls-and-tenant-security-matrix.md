# PostgreSQL RLS and tenant security matrix

## Catalog evidence

Collected 2026-11-11 from the fresh, 126-migration PostgreSQL 17.7 disposable database `pos_remediation_20261111` at the tested source candidate. The catalog has 92 public ordinary tables. RLS is both enabled and forced on 15 tables; the other 77 tables do not have database RLS. `pos_app` and `pos_migrator` are both `NOSUPERUSER`, `NOCREATEDB`, `NOCREATEROLE`, and `NOBYPASSRLS`.

The complete per-table catalog extract is [05-table-security-matrix.csv](05-table-security-matrix.csv). It records table columns, direct scope signals, foreign-key parent paths, RLS and FORCE RLS state, policies, runtime role privileges, and a reminder that catalog evidence cannot certify API route/repository enforcement. Scope classes are schema-derived triage labels, not a claim that each table's business ownership or authorization has been fully reviewed.

## RLS-protected tables

The 15 protected tables are `audit_export_logs`, `e_invoice_submissions`, `email_change_requests`, `invoice_delivery_logs`, `notification_preferences`, `notifications`, `password_change_history`, `payments`, `permission_audit_logs`, `report_export_audit_logs`, `role_permissions`, `user_auth_audit_logs`, `user_roles`, `user_settings`, and `user_settings_audit_log`.

Existing policies use transaction-local `app.tenant_id`, `app.restaurant_id`, and in the notifications policy `app.user_id`; `user_roles` and `role_permissions` also have explicit `pos_migrator` policies. `permission_audit_logs` has a SELECT-only runtime policy. The table-level CSV captures the exact current policy definitions and the runtime grants.

## Restricted-role tests run

On the disposable database, existing integration tests passed:

- `TestRLSPilot_TenantIsolationWithRestrictedRole`: tenant A/B reads, cross-tenant insert rejection, update/delete filtering, no-context fail-closed behavior, connection reuse, and runtime-role attributes.
- `TestAuthAuditAndUserSettingsRLSWithRestrictedRuntimeRole`: auth audit, user settings, notification preferences, tenant/user ownership checks, unscoped denial, mutation isolation, and notification-preference rollback.
- Payment PostgreSQL integration tests: 14 tests passed, including restricted `pos_app` payment transactions, unscoped fail-closed reads, foreign-tenant denial, duplicate/out-of-order callbacks, concurrent void/success and refunds, invalid signatures, provider-success/local-commit recovery, durable refund recovery, and replay idempotency.
- Repository PostgreSQL isolation tests: 27 tests passed across accounting parent references and atomic writes, categories, drivers, HR employee references, inventory movement references, orders/items, products, RBAC assignments, themes, and user scope.

These are representative data-path tests, not proof of RLS or authorization for all 92 tables or every HTTP route.

## Remaining security work

RLS remains absent from core tables including `tenants`, `restaurants`, `users`, `orders`, `order_items`, `products`, `categories`, `inventory`, and `employees`. The repository contains table-scoped integration tests and API scoping, but a complete route-by-role/tenant matrix and comprehensive per-table runtime-role review have not been completed. The CSV includes rows marked for scope review, including child tables whose ownership must be inherited through parent references, global candidates, operations tables, and theme references with no declared FK.

No new RLS policies or migrations were added. Do not enable RLS globally based on column names alone; a policy plan still needs application call-path review, parent/child integrity tests, runtime grant compatibility, and per-migration rollback evidence.

**Status: PARTIAL.** Restricted-role RLS behavior passed for the exercised policies and repository scopes. Full tenant isolation/security acceptance remains open.
