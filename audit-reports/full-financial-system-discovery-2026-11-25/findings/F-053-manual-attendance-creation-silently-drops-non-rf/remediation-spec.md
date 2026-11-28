# F-053 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Manual attendance creation silently drops non-RFC3339 clock times and accepts the record with zero hours, so overtime and absence metrics miss the shift.

## 2. Current behavior (verified)
"09:00" form -> 201 with clock_in/clock_out NULL and regular_hours 0 / overtime_hours 0; reversed order on the same form also 201; RFC3339 form computes 8.00 + 2.50; reversed RFC3339 -> 500 with raw constraint text valid_clock_times; overtime report for November shows 0 hours | RUN-003: attendance with 09:00/19:30 stored NULL times and 0 hours (201); RFC3339 computes 8+2.5 (REPRODUCED). 

## 3. Required behavior
Reject invalid time format; compute 8 regular + 2.5 overtime hours

## 4. Reproduction conditions
POST /hr/attendance with clock_in "09:00", clock_out "19:30" (10.5 h shift); repeat with RFC3339 datetimes

## 5. Verified evidence
- evidence/database/RUN-003/PH-10/AUD-HR-MATRIX.json
- evidence/database/RUN-003/PH-10/AUD-HR-MATRIX2.json
- tests/runtime/ph10_hr_matrix.mjs
- tests/runtime/ph10_hr_matrix2.mjs

## 6. Root cause and confidence
CONFIRMED — clock_in/clock_out are parsed only as RFC3339 and a failed parse is silently ignored (hr_attendance_handler.go:235-243), so "09:00" is dropped and the row stores NULL times

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/hr_attendance_handler.go`:188-260, `backend/internal/handler/http/hr_attendance_handler.go`:235, `backend/internal/handler/http/hr_attendance_handler.go`:235-243
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: POST /api/v1/hr/attendance: 400 field errors; response includes total_hours
- Accounting logic / posting rules: none named
- Reports affected: HR attendance; HR overtime/absenteeism reports; payroll overtime input

## 8. Impact
- Financial: Overtime cost understated whenever times are entered in the simple form
- Operational: HR attendance; HR overtime/absenteeism reports; payroll overtime input workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: strict clock-time handling in backend/internal/handler/http/hr_attendance_handler.go (CreateAttendance, parse block at :235-243): accept RFC3339 and HH:MM (combined with attendance_date in the tenant time zone), return 400 with a field error when a supplied value cannot be parsed (never ignore the error), require clock_out > clock_in, compute and return total hours and overtime server-side; add a CHECK clock_out > clock_in on the table.

## 10. Alternative designs
A: strict parse + 400 + computed hours (recommended). B: keep lenient parsing and warn (data stays wrong).

## 11. Required policy decisions
Overtime rules (thresholds) are HR policy; not changed here.

## 12. Required data migrations
CHECK constraint after listing existing rows with NULL times and zero hours.

## 13. Historical data implications
Rows with NULL times and 0 hours (evidence) are listed for HR correction.

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
F-031

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/hr_attendance_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Simple and RFC3339 inputs; reversed times; hours equal independent calculation (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- POST with clock_in "09:00" and clock_out "19:30" stores both times and 10.5 hours.
- POST with clock_in "abc" returns 400 and stores nothing.
- clock_out before clock_in returns 400.
