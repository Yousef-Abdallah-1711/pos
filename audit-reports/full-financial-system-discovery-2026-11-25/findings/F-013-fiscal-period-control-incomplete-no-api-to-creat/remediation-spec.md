# F-013 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Fiscal-period control incomplete: no API to create/close periods; void and draft creation allowed in a closed period.

## 2. Current behavior (verified)
Post blocked; draft create allowed; void of posted entry allowed (capital 100000 reversed) | RUN-003 (AUD-ACC-FLOW JE-13/JE-14): journals dated 2019-01-15 and 2030-06-01 post successfully and auto-create fiscal_period 2019-01/2030-06; GET /accounting/fiscal-periods is 404 (no period API). 

## 3. Required behavior
Closed period blocks all postings including voids

## 4. Reproduction conditions
Insert closed period (SQL fixture) then post / create / void

## 5. Verified evidence
- runtime/14 | evidence/database/RUN-003/PH-11/AUD-ACC-FLOW.json | tests/runtime/ph11_accounting_flow.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/repository/journal_entry_repo.go:813-842 ensureFiscalPeriodOpenTx returns nil when no fiscal_periods row matches, so any date posts; only status=closed blocks; no API creates/closes periods (GET /accounting/fiscal-periods is 404)

## 7. Affected paths
- Backend/database: `backend/internal/repository/journal_entry_repo.go`:813-842
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: GET/POST /api/v1/accounting/fiscal-periods (new); POST /api/v1/accounting/fiscal-periods/{id}/close and /reopen (new)
- Accounting logic / posting rules: none named
- Reports affected: Accounting

## 8. Impact
- Financial: 100000 reversed in closed period
- Operational: Accounting workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: add fiscal-period endpoints (list/create/close/reopen) and enforce an open period in ensureFiscalPeriodOpenTx (backend/internal/repository/journal_entry_repo.go:813-842): reject posting when no period exists for the entry date (instead of returning nil), apply the same check to voids and draft creation, and require an accountant-level permission for close/reopen with an audit log.

## 10. Alternative designs
Alternative A: enforce in the database (constraint/trigger/lock) so every writer is covered; Alternative B: enforce in the application layer only (smaller change, but direct writers and future code paths stay unprotected). Prefer A for monetary invariants, B for presentation/validation defects.

## 11. Required policy decisions
POL-05.

## 12. Required data migrations
No schema change expected; confirm during implementation.

## 13. Historical data implications
Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.

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
none

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/repository/journal_entry_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Closed-period matrix tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- POST journal dated 2019-01-15 or 2030-06-01 without a period returns 4xx.
- After closing 2026-11, posting, voiding and creating drafts dated in November are rejected; reopening permits them.
- GET /accounting/fiscal-periods returns the periods (currently 404).
