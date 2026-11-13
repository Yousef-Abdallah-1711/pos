# 12 — Financial security and tenant isolation (PROPOSED)

Derive tenant identity from verified auth/session context, not caller-provided IDs. Every read/write/export/job checks tenant and restaurant at object level. Test roles for journal create/post/void, account changes, bank details, close/reopen, audits and exports. Super Admin access is separate, reasoned, time-bound and audited. Suspension blocks new writes/provider operations but preserves historical evidence.

RLS may be defense in depth after connection-pool reset and privileged-role analysis; it does not replace application checks. Cache keys/jobs include tenant+restaurant. Negative tests: Tenant A guesses B journal/order/export; sibling restaurant requests record; read-only posts; cashier exports bank report; suspended tenant replays callback; worker consumes cross-tenant message. Expect deny and audit. Prior audit tenant/RBAC matrix is NOT TESTED; wrappers alone do not prove isolation.
