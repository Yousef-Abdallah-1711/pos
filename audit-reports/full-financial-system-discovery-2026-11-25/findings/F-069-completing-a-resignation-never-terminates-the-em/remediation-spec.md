# F-069 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Completing a resignation never terminates the employee: turnover stays 0 and the person stays active and payable until the employee is deleted.

## 2. Current behavior (verified)
Employee remains employment_status=active, is_active=true, termination_date NULL; turnover departures 0, rate 0. Only DELETE /hr/employees/{id} sets termination_date (to CURRENT_DATE, not the last working date); then turnover = 1/1 = 100% RUN-005 real-UI reproduction: Complete on resignation 9 returned 200 with settlement_amount 0; the employee stayed active and no termination date was set.

## 3. Required behavior
Employee becomes terminated with termination_date = last working date; turnover departures = 1

## 4. Reproduction conditions
Tenant 4: create employee (backend-native keys), resignation 2026-11-20 / last day 2026-11-25, acknowledge, complete (all 2xx); read the employee row and the turnover report

## 5. Verified evidence
- evidence/network/RUN-004/PH-12/AUD-R4-MET2-CAPTURE.json
- evidence/database/RUN-004/PH-12/AUD-R4-HR.json
- evidence/network/RUN-005/PH-04/AUD-R5-DESTRUCTIVE-B.run4-leave-resign.json

## 6. Root cause and confidence
CONFIRMED — resignation_repo.go:201-215 Complete updates only resignations; employee_repo.go:370-380 is the only termination_date writer; hr_report_handler.go:276-285 counts departures only from employees.termination_date

## 7. Affected paths
- Backend/database: `backend/internal/repository/resignation_repo.go`:201, `backend/internal/repository/employee_repo.go`:370, `backend/internal/handler/http/hr_report_handler.go`:276, `backend/internal/repository/resignation_repo.go`:201-215, `backend/internal/repository/employee_repo.go`:370-380, `backend/internal/handler/http/hr_report_handler.go`:276-285
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: POST /api/v1/hr/resignations/{id}/complete: also updates the employee (documented side effect)
- Accounting logic / posting rules: none named
- Reports affected: HR resignations, HR turnover report, employees, payroll eligibility

## 8. Impact
- Financial: Payroll and attendance continue for a resigned employee; settlement_amount is stored on the resignation with no payroll or GL effect
- Operational: HR resignations, HR turnover report, employees, payroll eligibility workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: terminate the employee when a resignation completes. In backend/internal/repository/resignation_repo.go Complete (:201-215) update employees in the same transaction: employment_status = terminated, termination_date = last_working_date, is_active = false from that date; make DeleteEmployee (employee_repo.go:370-380) use the contract/resignation date when one exists; exclude terminated employees from future payroll and attendance creation; the turnover report then counts the departure.

## 10. Alternative designs
A: complete resignation terminates employee (recommended). B: keep manual termination and warn when a completed resignation exists.

## 11. Required policy decisions
Notice period and final-settlement rules are HR policy (proposed POL-19).

## 12. Required data migrations
No schema change; data-correction for completed resignations whose employee is still active.

## 13. Historical data implications
Employees with completed resignations still active are listed for HR review.

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
F-014, F-056 and shared-file findings F-056

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/hr_report_handler.go, backend/internal/repository/employee_repo.go, backend/internal/repository/resignation_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Complete resignation -> employee terminated, departure counted, excluded from active payroll (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Completing a resignation with last working day 2026-11-25 sets termination_date 2026-11-25 and status terminated.
- The turnover report shows departures 1 and the matching rate for the period.
- No salary or attendance can be created for the employee after the termination date.
