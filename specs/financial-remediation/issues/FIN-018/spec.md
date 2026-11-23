# FIN-018: No complete fiscal close and retained-earnings roll-forward established

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-018 — No complete fiscal close and retained-earnings roll-forward established

- **Severity/status:** P2 — SOURCE-CONFIRMED report gap / workflow incompleteness.
- **Problem / impact:** Open periods lack reliable current earnings in equity; closed-period corrections and year-end carryforward may be inconsistent.
- **Current / expected:** Fiscal periods and posting checks exist, but a full close/reopen/earnings transfer/audit workflow was not established. Expected role-gated, repeatable close with TB/report lock and linked corrections.
- **Root cause:** Missing/partial lifecycle across `backend/migrations/153_create_fiscal_periods.sql`, `financial_report_usecase.go`, journal repository.
- **Reproduction/evidence:** Static inventory; no fiscal close test or direct DB run.
- **Baseline scope:** Main/current completeness not established.
- **Recommended solution / components:** Design close state machine, lock, retained earnings policy and reopen/adjustment event. Coordinate CFR-010/013.
- **Migration:** Possible close run/approval and lock metadata after design.
- **Dependencies/policy:** Fiscal year, adjustment periods, retained earnings, reopening authority require accountant/owner approval.
- **Acceptance:** Close twice is idempotent; posting blocked in closed period; current earnings roll forward exactly; reopen/reversal audited; all statements reconcile.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
