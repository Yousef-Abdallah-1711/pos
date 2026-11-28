# F-015 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Payroll can be paid without approval; no accrual; no partial payment.

## 2. Current behavior (verified)
Paid with is_approved=false; journal only at payment as direct expense | RUN-003 SQL re-verification (AUD-EVD-F015): salaries id 1 and 2 are status=paid with is_approved=false; journals JE-2026-00007/00049 are Dr 5210 5,000 / Cr 1111 5,000 posted at payment only; no accrual journal exists (STILL REPRODUCIBLE; UI re-run pending T-P10-003). 

## 3. Required behavior
Approval gate; accrual per PAY-01

## 4. Reproduction conditions
Create salary then Mark Paid

## 5. Verified evidence
- runtime TX-007 TX-008 | evidence/database/RUN-003/PH-20/AUD-EVD-F008-F015-sql.txt

## 6. Root cause and confidence
CONFIRMED — backend/internal/repository/salary_repo.go:344-380 MarkAsPaidTx updates status=paid with no is_approved precondition, and payroll posting happens once at payment (usecase/auto_posting_engine.go postPayrollEntry :276), so no accrual journal exists

## 7. Affected paths
- Backend/database: `backend/internal/repository/salary_repo.go`:344-380, `backend/internal/usecase/auto_posting_engine.go`:261-345, `backend/internal/usecase/auto_posting_engine.go`
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: PUT /api/v1/hr/salaries/{id}/pay: 409 when not approved; accepts amount for partial payment (additive)
- Accounting logic / posting rules: accounts 1111, 5210
- Reports affected: HR payroll; Accounting

## 8. Impact
- Financial: payroll payable 2000 absent
- Operational: HR payroll; Accounting workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: require approval and book payroll in two steps. (1) In backend/internal/repository/salary_repo.go MarkAsPaidTx (:344-380) add the precondition is_approved = true and reject otherwise with 409; (2) introduce an accrual journal at approval (Dr 5210 salaries expense / Cr 2110 accrued salaries) and have payment post Dr 2110 / Cr cash-or-bank (usecase/auto_posting_engine.go postPayrollEntry :276), so expense is recognised in the correct period; (3) allow partial payment with a paid_amount ledger and status partially_paid; (4) UI: disable "Mark paid" until approved.

## 10. Alternative designs
A: approval gate + accrual + partial payments (recommended). B: approval gate only (keeps direct-expense booking). C: leave as is and document (rejected: unapproved payroll can be paid).

## 11. Required policy decisions
POL-08 (payroll accrual and approval) in 56.

## 12. Required data migrations
Yes: salary payments table or paid_amount column plus status value; backfill existing paid rows as fully paid.

## 13. Historical data implications
Rows already paid without approval (salaries 1 and 2 in the evidence) are listed in an exception report for management sign-off; no automatic change.

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
none and shared-file findings F-026, F-040, F-061

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/repository/salary_repo.go, backend/internal/usecase/auto_posting_engine.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Payroll approval and accrual tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Marking an unapproved salary paid returns 409 and writes nothing.
- Approval creates one accrual journal; payment creates one settlement journal; the 2110 balance returns to zero after full payment.
- Two concurrent mark-paid calls produce exactly one payment journal.
