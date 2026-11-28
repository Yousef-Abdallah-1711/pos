# F-064 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Tables API accepts a duplicate table number and returns 500 for an over-long name.

## 2. Current behavior (verified)
duplicate -> 201 (two tables numbered 1, ids 1 and 2); long name -> 500 "Failed to create table"; capacity (1-50), negative number and delete-with-active-orders (409) are validated correctly | RUN-003: duplicate table number accepted (201); 1,800-char name -> 500 (REPRODUCED). 

## 3. Required behavior
409 duplicate table number; 400 name length

## 4. Reproduction conditions
POST /tables {table_number:1,name:dup,capacity:4} when table 1 exists; POST /tables with a 1,800-character name

## 5. Verified evidence
- evidence/database/RUN-003/PH-16/AUD-TABLES.json
- tests/runtime/ph16_tables.mjs

## 6. Root cause and confidence
PROBABLE — PROBABLE: table create validation covers capacity and number sign but not uniqueness (no unique constraint on tenant, restaurant, table_number) or string length

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/table_handler.go`
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`
- API contracts: unchanged unless the design below says otherwise (see 08-api-inventory.csv)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0 (dine-in orders can attach to an ambiguous table)
- Operational: Tables; POS dine-in workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Unique constraint (tenant_id, restaurant_id, table_number) with 409 mapping; length validation

## 10. Alternative designs
Alternative A: enforce in the database (constraint/trigger/lock) so every writer is covered; Alternative B: enforce in the application layer only (smaller change, but direct writers and future code paths stay unprotected). Prefer A for monetary invariants, B for presentation/validation defects.

## 11. Required policy decisions
No accounting policy decision is required; product decisions are limited to the behaviour described in Required behavior.

## 12. Required data migrations
Yes — a numbered migration after the current highest in backend/migrations (data-correction + constraint/trigger change). Existing rows must be reconciled first.

## 13. Historical data implications
No financial history is affected; existing data needs no correction unless noted in the evidence.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-031, F-062

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/table_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Duplicate number returns 409; over-long name returns 400 (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Re-running the reproduction steps yields the Expected result: 409 duplicate table number; 400 name length
- The independent oracle/SQL comparison for the evidence dataset shows zero residual for this finding.
- All tests in regression-test-spec.md pass on the isolated stack with real headed Chrome for UI-facing checks.
