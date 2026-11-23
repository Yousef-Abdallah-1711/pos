# FIN-016: No complete gift-card/deposit/deferred-revenue lifecycle established

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P3. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-016 — No complete gift-card/deposit/deferred-revenue lifecycle established

- **Severity/status:** P3 — SUSPECTED (scope search did not establish a full lifecycle).
- **Problem / impact:** Gift-card sale/redemption/expiry or customer deposits may be treated as immediate sales or omitted from liability/revenue reporting.
- **Current / expected:** No complete financially integrated state machine was evidenced in inspected financial catalog; expected liability on issue and controlled revenue recognition on redemption/expiry per approved policy.
- **Root cause:** Product capability gap not yet proven by exhaustive route/schema search. Candidate areas: order/payment/accounting APIs and `backend/migrations/136_create_chart_of_accounts.sql` onward.
- **Reproduction/evidence:** No gift-card lifecycle tests located/run; absence is not a claim that no partial code exists.
- **Baseline scope:** Completeness unverified on main/current.
- **Recommended solution / components:** Confirm product scope, define stored value liability, redemption, breakage/expiry, refund and reconciliation; then implement.
- **Migration:** Likely balance/transaction ledger schema if required.
- **Dependencies/policy:** Owner/accountant/legal policy; local stored-value rules.
- **Acceptance:** Issue/redeem/refund/expire flows reconcile liability, cash and revenue; duplicates and cross-tenant access denied.
- **Checklist:** OPEN (scope confirmation).



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
