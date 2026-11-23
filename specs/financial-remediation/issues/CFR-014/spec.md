# CFR-014: Inventory valuation, COGS timing and refunds/adjustments are not proven to reconcile

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-014 — Inventory valuation, COGS timing and refunds/adjustments are not proven to reconcile

- **Severity/status:** P2 — SUSPECTED / partially source-confirmed; method and end-to-end arithmetic BLOCKED.
- **Problem / impact:** Incorrect cost basis or missing return/wastage reversal makes inventory assets and COGS/profit inaccurate.
- **Current / expected:** Inventory movement and COGS auto-post paths exist (`144_create_inventory_movements.sql`, `auto_posting_engine.go`); examined order path calculates item cost using stored unit cost. Costing method and correspondence among stock movement, cost layer and GL were not proven.
- **Root cause:** Not fully established; candidate paths `backend/internal/usecase/auto_posting_engine.go:155-240`, `backend/internal/repository/inventory_movement_repo.go`, migration 144.
- **Reproduction/evidence:** Historical source skim only; no current PostgreSQL or inventory fixture. Local-only adjustment handler is not origin/main.
- **Baseline scope:** Main has partial inventory/COGS path; completeness defect is unverified on both. Local handler is current-branch-only and unverified as integrated.
- **Recommended solution / components:** Approve FIFO/weighted average/specific cost and recognition timing; one atomic stock movement and journal event; explicit returns, comps, waste and negative-stock treatment.
- **Migration:** Cost layer/valuation schema only after method chosen; reconcile opening stock before rollout.
- **Dependencies/policy:** Inventory valuation method and count/cutoff approval.
- **Acceptance:** Independently calculated receipts/sales/returns/waste/adjustments reconcile quantities, valuation, COGS, GL and inventory report; concurrency and rollback injection.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
