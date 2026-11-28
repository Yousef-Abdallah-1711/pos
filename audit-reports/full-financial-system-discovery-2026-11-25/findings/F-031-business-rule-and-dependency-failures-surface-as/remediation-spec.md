# F-031 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Business-rule and dependency failures surface as wrong status codes and leak internals.

## 2. Current behavior (verified)
Paid-order cancel 500 (item cancel 409); receipt during outage 400 with raw role error pq role pos_accounting not permitted; DB down: /health 200 degraded and authenticated calls 403 TENANT_ACCESS_BLOCKED contact support | RUN-003 (AUD-ACC-FLOW ACC-02): duplicate account code returns raw DB text pq: duplicate key value violates unique constraint chart_of_accounts_tenant_id_restaurant_id_account_code_key (REPRODUCED). 

## 3. Required behavior
409/503 with generic message

## 4. Reproduction conditions
Cancel a paid order; receipt during accounting-role outage; DB outage

## 5. Verified evidence
- evidence/runtime/35 36 40 | evidence/database/RUN-003/PH-11/AUD-ACC-FLOW.json | tests/runtime/ph11_accounting_flow.mjs

## 6. Root cause and confidence
CONFIRMED — error mapping is done by passing err text to the client (see F-025 sites) and returning 500 for constraint violations (e.g. main.go:1302; hr handlers rely on DB check constraints)

## 7. Affected paths
- Backend/database: `backend/cmd/api/main.go`:1302
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: status codes change for the listed cases (documented in 08-api-inventory.csv)
- Accounting logic / posting rules: none named
- Reports affected: Orders; Receipts; Platform

## 8. Impact
- Financial: n/a
- Operational: Orders; Receipts; Platform workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: status-code mapping rules, delivered together with the F-025 error mapper. Rules: business-rule violations 409/422 (paid-order cancel, over-settlement, closed period), dependency outages 503 with Retry-After (accounting pool or Redis down), missing resources 404, validation 400; the audited cases (paid-order cancel 500, receipt during outage 400 with a role error, DB down returning 403 TENANT_ACCESS_BLOCKED) each get an explicit mapping and a regression test; /health reports degraded with the failing component.

## 10. Alternative designs
A: mapping table + tests (recommended). B: case-by-case patches.

## 11. Required policy decisions
None.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-025 and shared-file findings F-018, F-020, F-025, F-030, F-033, F-036, F-037, F-071

## 21. Risks
- Regression in modules sharing the same source files (backend/cmd/api/main.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Error-contract tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Each audited case returns the documented status and a generic message.
- During a database outage authenticated calls return 503, not 403.
- No response leaks pq or role names.
