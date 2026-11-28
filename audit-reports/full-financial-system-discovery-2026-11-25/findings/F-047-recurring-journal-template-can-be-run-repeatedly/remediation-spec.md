# F-047 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Recurring journal template can be run repeatedly for the same period and posts a duplicate journal each time.

## 2. Current behavior (verified)
Run 1 -> JE id 98 (rent +700); run 2 -> JE id 99 (rent +700 again): rent expense 1,400 for one period; next_run_date already advanced to 2026-12-26 | RUN-003: recurring run twice -> journals 98 and 99, rent +700 each (CONFIRMED accounting_advanced_repo.go:396-420). 

## 3. Required behavior
Second run for the same period rejected or returns the first journal

## 4. Reproduction conditions
Create template (auto_post) and POST /accounting/recurring-journals/{id}/run twice on 2026-11-26

## 5. Verified evidence
- evidence/database/RUN-003/PH-11/AUD-ACC-FLOW2.json
- tests/runtime/ph11_accounting_flow2.mjs

## 6. Root cause and confidence
CONFIRMED — repository/accounting_advanced_repo.go:396-420 RunRecurringJournalTemplate checks only template existence/is_active; it does not compare runDate with next_run_date/last_run_date and takes no row lock

## 7. Affected paths
- Backend/database: `backend/internal/repository/accounting_advanced_repo.go`:396-470, `backend/internal/repository/accounting_advanced_repo.go`:396-420
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: POST /api/v1/accounting/recurring-journals/{id}/run: 409 on duplicate period; optional force+reason
- Accounting logic / posting rules: none named
- Reports affected: Accounting - recurring journals

## 8. Impact
- Financial: 700 duplicated per extra run in test
- Operational: Accounting - recurring journals workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: make recurring journals idempotent per period. In backend/internal/repository/accounting_advanced_repo.go RunRecurringJournalTemplate (:396-420) lock the template row (SELECT ... FOR UPDATE), reject a run when runDate < next_run_date unless an explicit force flag with a reason is supplied, store a unique key (template_id, period) on the generated journal and advance next_run_date/last_run_date in the same transaction.

## 10. Alternative designs
A: lock + unique (template, period) + next_run_date check (recommended). B: unique key only (concurrency handled, but manual re-runs still blocked silently).

## 11. Required policy decisions
POL-13 (recurring-journal idempotency) in 56.

## 12. Required data migrations
Yes: unique index on (template_id, period); existing duplicates listed first.

## 13. Historical data implications
Duplicate recurring journals in the evidence (2 runs, 2 journals) are reversed by approved journals.

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
none and shared-file findings F-009, F-016, F-043, F-045, F-068, F-072, F-073

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/repository/accounting_advanced_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Double-click / concurrent run creates exactly one journal per period (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Running the same template twice for the same period creates one journal; the second call returns 409 with the existing journal id.
- Five concurrent runs create one journal.
- A forced run records the reason and creates a second journal only when explicitly requested.
