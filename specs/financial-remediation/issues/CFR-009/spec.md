# CFR-009: Paid payroll remains editable/deletable and duplicate mark-paid guard is unproven

Status: PARTIALLY IMPLEMENTED — correction workflow and current acceptance pending. Original severity: P1.

## Original issue text

### CFR-009 — Paid payroll remains editable/deletable and duplicate mark-paid guard is unproven

- **Severity/status:** P1 — SOURCE-CONFIRMED for missing status checks on update/delete; DB reproduction BLOCKED.
- **Problem / impact:** Payroll history can be rewritten after payment/posting, leaving journal, bank and payroll records inconsistent; a hard delete can orphan its journal.
- **Current / expected:** `UpdateSalary` updates fields without paid-state refusal; `DeleteSalary` invokes repository deletion; `MarkAsPaid` behavior has no proven concurrency-safe one-time precondition. Expected: paid record immutable; correction through approval and linked reversal; unpaid-only deletion; mark-paid idempotent.
- **Root cause:** `backend/internal/handler/http/hr_salary_handler.go:237-382,388-438`; `backend/internal/repository/salary_repo.go:216-227`; migration `024_create_salaries_table.sql`.
- **Reproduction/evidence:** Source review; no live salary fixture. Existing sqlmock payroll auto-post test covers a happy path only.
- **Baseline scope:** Present on main/current.
- **Recommended solution / components:** Salary state machine, guarded atomic update, audit trail, reversal/correction workflow, approval and statutory deduction policy.
- **Migration:** Audit/reversal link or state columns only as designed; preserve payroll history.
- **Dependencies/policy:** Payroll statutory rules/roles and correction approval by owner/accountant.
- **Acceptance:** Update/delete paid returns conflict; duplicate concurrent mark-paid produces one payment/journal; correction preserves original and nets correctly; payroll report equals GL.
- **Checklist:** OPEN.



## Current evidence and open questions

At main `3fa3702148771f39af5c9f7f25f377ecdd3f0a98`, salary repository update and delete predicates reject paid rows. `MarkAsPaidTx` locks the salary, preserves repeated-payment metadata, and writes an audit record; the HTTP handler uses the same transaction for payroll posting. PR [#40](https://github.com/Yousef-Abdallah-1711/pos/pull/40) merged as `0d8fea4272bde8f9c9da62d486a02474276a8815`. Seven PostgreSQL integration tests cover unpaid mutations, paid immutability, concurrent payment, failure rollback, disabled posting and scope. The [historical tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md) reports those seven passed on PostgreSQL 17, but that is not a current-SHA rerun. A paid-salary correction and reversal policy has not been established by this trace.

## Definition of done

- Paid edit/delete return conflict, unpaid edits remain available, and concurrent repeated payment creates one transition, journal and audit trail in scoped PostgreSQL fixtures.
- Failures at salary, audit, journal and GL writes roll back all effects; disabled posting cannot mark paid.
- Authorized correction preserves original payment and uses an approved linked reversal; payroll reports independently reconcile to GL. Obtain payroll/statutory decisions before implementing policy-dependent corrections.
- Retain exact-SHA PostgreSQL, API, role isolation, CI, review and main evidence before closure.
