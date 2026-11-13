# Security and tenant isolation

**Status: PARTIAL / UNVERIFIED.** Accounting routes are wrapped in permission checks in main.go (module 15; reports module 17). Runtime positive path used synthetic owner. No separate low-privilege user or second tenant was used for BOLA checks.

Observed:
- Routes wrap handlers in permission checks.
- Owner created accounts and posted a journal in the disposable tenant.
- No unauthorized, cross-tenant, cross-restaurant, object-id swap, or export authorization negative case was executed.
- Tenant/restaurant ids exist on journal and ledger rows; this alone does not prove every query scopes by both.
- No RLS conclusion is made; policies were not exhaustively inspected.

Required: object-id swap across two synthetic tenants; read/write/export route checks with owner and least privilege; tenant predicates and trusted tenant resolution review; audit-log immutability review.

## Completion status — 2026-11-11

This continuation did not create a second synthetic tenant or execute cross-tenant financial object/export requests. Low-privilege role enforcement was not exercised. Tenant isolation and financial RBAC therefore remain **NOT TESTED** for the requested adversarial matrix; earlier route access with a synthetic owner does not establish isolation. Do not interpret this report as a pass.
