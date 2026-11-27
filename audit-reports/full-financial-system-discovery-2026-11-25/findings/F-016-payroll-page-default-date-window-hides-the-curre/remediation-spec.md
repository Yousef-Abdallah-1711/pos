# F-016 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Payroll page default date window hides the current month's salary records.

## 2. Current behavior (verified)
0 records until To-date changed | RUN-003 (AUD-EVD-F016): default window start_date=2026-10-27,end_date=2026-11-26 returns total 0 and UI says 'Create your first salary record', while salaries row id 1 (tenant 4, period 2026-11-01..2026-11-30, paid) exists (STILL REPRODUCIBLE). 

## 3. Required behavior
Current records visible

## 4. Reproduction conditions
Open payroll page with record period ending after today

## 5. Verified evidence
- explore | evidence/network/RUN-003/PH-20/AUD-EVD-gaps.json
- tests/runtime/ph20_evidence_gaps.mjs
- evidence/screenshots/RUN-003/PH-20/AUD-EVD-F016-payroll-default.png

## 6. Root cause and confidence
CONFIRMED — frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx:57-58,91 initialise the filter to the last 30 days (today-30d .. today) and send it to GET /hr/salaries; HR cost report filters by pay_period_end < date_to (accounting_advanced_repo.go:1236-1243), so a salary for the whole current month ending in the future is hidden

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx`:57-58, `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx`:91, `backend/internal/repository/accounting_advanced_repo.go`:1236-1243
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: GET /api/v1/hr/salaries: no change except overlap semantics for date filters
- Accounting logic / posting rules: none named
- Reports affected: HR payroll

## 8. Impact
- Financial: n/a
- Operational: HR payroll workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: default the payroll list filter to the current pay period, not the last 30 days. In frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx:57-58,91 initialise start_date to the first day of the current month and end_date to the last day, and filter by overlap with the pay period (pay_period_start <= end and pay_period_end >= start) in the list query; show the active window in the page header; the HR cost report filter (pay_period_end < date_to, accounting_advanced_repo.go) should use the same overlap rule.

## 10. Alternative designs
A: current-period default + overlap rule (recommended). B: remove the default filter (large lists). C: keep the 30-day window and add an empty-state hint (does not fix the report).

## 11. Required policy decisions
None (UX/default). Date-filter semantics documented in 12-report-formula-catalog.md.

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
none and shared-file findings F-009, F-043, F-045, F-047, F-068, F-072, F-073

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/repository/accounting_advanced_repo.go, frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
UI filter test (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A salary record for the current month is visible on first load without changing filters.
- The HR cost report includes a pay period that straddles date_to.
- The empty state never says "Create your first salary record" when records exist outside the window.
