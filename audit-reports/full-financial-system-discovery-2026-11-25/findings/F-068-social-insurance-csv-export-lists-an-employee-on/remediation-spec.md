# F-068 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Social-insurance CSV export lists an employee once per active insurance policy, repeating base salary and deductions.

## 2. Current behavior (verified)
1 salary row (SQL) produced 2 data rows after the second active policy; both repeat base_salary 5,000.00 and the deductions, so any total over the file double-counts

## 3. Required behavior
One row per employee salary record

## 4. Reproduction conditions
Tenant 4, month 2026-11: export with one active policy, add a second active policy for the same employee (API-DIAGNOSTIC), export again

## 5. Verified evidence
- evidence/network/RUN-004/PH-12/AUD-R4-HRR007-before.csv
- evidence/network/RUN-004/PH-12/AUD-R4-HRR007-after-second-policy.csv
- evidence/database/RUN-004/PH-12/AUD-R4-HR.json

## 6. Root cause and confidence
CONFIRMED — accounting_advanced_repo.go:1421-1446 LEFT JOIN insurance_tracking it ON ... AND it.status active has no DISTINCT or per-employee limit

## 7. Affected paths
- Backend/database: `backend/internal/repository/accounting_advanced_repo.go`:1437, `backend/internal/handler/http/accounting_advanced_handler.go`:912, `backend/internal/repository/accounting_advanced_repo.go`:1421-1446
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: CSV columns unchanged
- Accounting logic / posting rules: none named
- Reports affected: HR reports > Social insurance export

## 8. Impact
- Financial: Proportional: each extra active policy repeats the employee base salary and deductions (here 5,000.00 base salary listed twice)
- Operational: HR reports > Social insurance export workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: one row per salary record in the social-insurance export. In backend/internal/repository/accounting_advanced_repo.go GetSocialInsuranceReportRows (:1421-1446) replace the LEFT JOIN insurance_tracking at :1437 with a LATERAL subquery that returns one active policy per employee (latest start) or aggregates policy numbers with string_agg, so salary and deductions appear once per employee; add a test with two active policies.

## 10. Alternative designs
A: LATERAL one-row-per-employee (recommended). B: DISTINCT ON (employee) in the outer query.

## 11. Required policy decisions
Which policy is reported when several are active is an HR/legal decision (proposed: list all numbers in one cell).

## 12. Required data migrations
No.

## 13. Historical data implications
Previously exported files may contain duplicates; list affected months.

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
none and shared-file findings F-009, F-016, F-043, F-045, F-047, F-048, F-063, F-072, F-073

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Two active policies -> one row (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- With two active policies for one employee the CSV has exactly one data row for that employee.
- Totals of base_salary and deductions over the file equal the salaries table.
