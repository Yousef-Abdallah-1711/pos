# FIN-020: Multi-currency / FX and consolidated financial reporting policy is not established

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P3. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-020 — Multi-currency / FX and consolidated financial reporting policy is not established

- **Severity/status:** P3 — SUSPECTED policy and implementation gap.
- **Problem / impact:** Combining restaurants/currencies without transaction, functional and presentation currency rules can make tenant-level consolidated reports materially wrong.
- **Current / expected:** Accounts and payment flows have currency fields, but conversion rates, dates, realized/unrealized FX and consolidation eliminations were not verified. Expected immutable original currency plus policy-defined base-currency valuation/reporting.
- **Root cause:** Not established; review candidates `domain/accounting.go`, payment/invoice repositories and report usecase.
- **Reproduction/evidence:** No FX test or consolidated report fixture.
- **Baseline scope:** Completeness unverified on main/current.
- **Recommended solution / components:** Keep FX as explicit policy/config with stored rate source/time and realized/unrealized gains/losses; scope consolidation.
- **Migration:** Potential rate/functional currency metadata after policy approval.
- **Dependencies/policy:** Functional currency, rate source/time, translation and consolidation rules; accountant approval.
- **Acceptance:** Independent multi-currency fixtures tie transaction-currency journal, base-currency ledger, FX gain/loss and consolidated statements.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
