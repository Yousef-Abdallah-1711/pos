# CFR-009: Paid payroll remains editable/deletable and duplicate mark-paid guard is unproven

Status: IMPLEMENTATION MERGED OR OPEN — ACCEPTANCE PENDING; verify scope and PRs. Original severity: P1. This is a source-backed recovery draft, not acceptance.

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

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
