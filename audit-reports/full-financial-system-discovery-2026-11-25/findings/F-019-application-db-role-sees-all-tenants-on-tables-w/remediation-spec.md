# F-019 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Application DB role sees all tenants on tables without RLS and can alter paid orders/payment amounts.

## 2. Current behavior (verified)
All 4 tenants' orders/journals/GL/COA/users visible; UPDATE orders/payments_received succeeded (rolled back); ledger tables protected by grants

## 3. Required behavior
Cross-tenant rows hidden; paid records immutable

## 4. Reproduction conditions
Connect as pos_app without tenant context

## 5. Verified evidence
- runtime/22

## 6. Root cause and confidence
CONFIRMED — only 13 of the 70 public tables with a tenant_id column have row-level security (RLS migrations 164-171 are the only ones); the application role pos_app has UPDATE on orders and payments and no immutability trigger exists on payments (orders carry only stock-transition triggers); journal_entry_lines is not updatable by pos_app

## 7. Affected paths
- Backend/database: `backend/migrations/164_rls_user_auth_audit_logs.sql`, `backend/migrations/171_rls_e_invoice_submissions.sql`
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: unchanged
- Accounting logic / posting rules: none named
- Reports affected: Platform

## 8. Impact
- Financial: n/a
- Operational: Platform workflows behave as in section 2.
- Security: Yes — see category: Security / defense in depth
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: extend row-level security and add immutability where money is stored. (1) New migrations following the pattern of migrations 164-171 enable and force RLS (tenant_id = current_setting(app.tenant_id)) on orders, order_items, payments, invoices_customer/supplier, payments_received/made, chart_of_accounts, employees, salaries and the other tenant tables (57 of 70 are unprotected), with the application setting the tenant GUC per request; (2) revoke UPDATE/DELETE on payments and paid orders from pos_app and route corrections through reversing records; (3) add a trigger that rejects changes to total_amount/payment fields once an order is paid (same style as cfr002 triggers on the ledger).

## 10. Alternative designs
A: RLS + revoke + immutability triggers (recommended, defence in depth). B: RLS only (does not stop an authorised-but-buggy writer from altering paid amounts). C: application-layer checks only (current state).

## 11. Required policy decisions
Security architecture decision: approve table list and role grants before implementation.

## 12. Required data migrations
Yes: several migrations; staged rollout (enable on staging first, then production) with a rollback script.

## 13. Historical data implications
No data change.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
Tighten database roles: pos_app loses UPDATE/DELETE on payment tables; add negative tests that attempt the forbidden writes.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none

## 21. Risks
- Regression in modules sharing the same source files (backend/migrations/164_rls_user_auth_audit_logs.sql, backend/migrations/171_rls_e_invoice_submissions.sql).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Restricted-role isolation suite (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- With app.tenant_id set to tenant A, SELECT on every tenant table returns only A rows; without the setting it returns none.
- UPDATE of orders.total_amount or payments.amount for a paid order by pos_app fails; the reversing-record path works.
- The full existing test suite passes with RLS enabled (no missing tenant GUC).
