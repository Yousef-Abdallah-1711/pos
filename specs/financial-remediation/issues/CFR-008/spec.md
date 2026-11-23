# CFR-008: Order tax hardcoded; invoice tax computed before discount

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-008 — Order tax hardcoded; invoice tax computed before discount

- **Severity/status:** P2 — SOURCE-CONFIRMED; jurisdictional expected value requires approval.
- **Problem / impact:** Tax can be under/over-collected or discounts misapplied; amounts may not match local policy.
- **Current / expected:** Order creation sets 10% in Go and doesn't set discount from a configured policy; customer invoice trigger calculates tax on gross line amount, then subtracts discount. Expected: tenant/jurisdiction-configured tax and approved inclusive/exclusive tax base/discount ordering.
- **Root cause:** `backend/internal/usecase/order_usecase.go:147-182`; `backend/migrations/141_create_invoices_customer.sql:78-82` (supplier calculation also review `invoice_supplier_repo.go:449-476`).
- **Reproduction/evidence:** Source formula confirms order 100 → tax 10 absent policy, and gross-based invoice tax formula. No approved jurisdiction/rate or DB test is available; numerical result is formula evidence, not legal assessment.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Tax-rule/config service, effective-dated jurisdiction and item tax category, deterministic allocation and approved calculation order; no embedded legal rate.
- **Migration:** Tax config/effective-date schema may be required; do not rewrite historical tax blindly.
- **Dependencies/policy:** Accountant/legal approval of tax registrations, exemptions, inclusive/exclusive, discount and rounding rules.
- **Acceptance:** Approved jurisdiction fixtures match independently calculated invoice/order tax to the minor unit; changing tenant policy cannot leak across tenants; historic records retain rate/basis.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
