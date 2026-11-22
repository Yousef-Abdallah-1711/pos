# CFR-013: Voiding a journal writes a reversal dated today without adequate closed-period handling

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-013 — Voiding a journal writes a reversal dated today without adequate closed-period handling

- **Severity/status:** P2 — SOURCE-CONFIRMED.
- **Problem / impact:** Void can alter a closed/current period unexpectedly; original entry's status and GL representation may be confusing and audit trail may not link explicit separate reversal document.
- **Current / expected:** `VoidEntryTx` uses `time.Now()` then appends reversed GL lines; reviewed path has no explicit closed-period validation for the reversal date. Expected: linked immutable reversal journal in an allowed period, preserving original and controlled reason/approval.
- **Root cause:** `backend/internal/repository/journal_entry_repo.go:423-493` (void date 465-477).
- **Reproduction/evidence:** Source review; closed-period void not executed in PostgreSQL.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Replace implicit GL-only reversal with explicit reversal event/journal; use configured open period and uniqueness, authorization, audit.
- **Migration:** Reversal-link metadata or reversal entry as normal journal; preflight existing voided rows.
- **Dependencies/policy:** Backdate and closed-period correction policy.
- **Acceptance:** Closed-period void rejected or linked reversal in approved open period; no original mutation/deletion; duplicate void prevented; reports show both entries consistently.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
