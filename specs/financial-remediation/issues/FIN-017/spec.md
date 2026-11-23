# FIN-017: Fixed-asset acquisition, depreciation and disposal lifecycle not verified end to end

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P3. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-017 — Fixed-asset acquisition, depreciation and disposal lifecycle not verified end to end

- **Severity/status:** P3 — BLOCKED for isolated DB fixture.
- **Problem / impact:** Asset register, depreciation expense/accumulated depreciation and disposal proceeds may diverge from GL.
- **Current / expected:** Schema exists at migration 145; complete journal lifecycle not verified. Expected approved useful life, method, convention, disposal and reversal behavior.
- **Root cause:** Not established; candidate migration 145 and auto-posting depreciation path.
- **Reproduction/evidence:** No lifecycle test run.
- **Baseline scope:** Feature exists partially on main/current; completeness BLOCKED.
- **Recommended solution / components:** Define asset classes/schedules and connect acquisition/depreciation/disposal to immutable journals.
- **Migration:** Only if lifecycle fields or event identity required.
- **Dependencies/policy:** Accountant approves capitalization threshold, depreciation method/life and disposal policy.
- **Acceptance:** Schedule independently ties book value/expense/GL at period boundaries; disposal and duplicate monthly run safe.
- **Checklist:** BLOCKED.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
