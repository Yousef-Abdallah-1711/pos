# CFR-005: Paid orders can suppress posting failures and have no proven cancellation/refund reversal

Status: IMPLEMENTATION MERGED OR OPEN — ACCEPTANCE PENDING; verify scope and PRs. Original severity: P1. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-005 — Paid orders can suppress posting failures and have no proven cancellation/refund reversal

- **Severity/status:** P1 — SOURCE-CONFIRMED for best-effort call; refund/cancel reversal missing in examined path.
- **Problem / impact:** Paid orders can be absent from financial reports, and refunded/canceled sales may remain as revenue/receivable/cash; period totals can be wrong.
- **Current / expected:** Handler marks payment state and calls `AutoPostPaidOrderByID` without checking a returned error; no reverse call on paid→refunded/canceled was found. Entry date uses order creation date. Expected: policy-date accounting, failure surfaced/retried, and linked reversal/refund entry.
- **Root cause:** `backend/internal/handler/http/admin_order_handler.go` payment-status handler around lines 470-515; `backend/internal/usecase/order_usecase.go:459-474`; `backend/internal/usecase/auto_posting_engine.go:149`.
- **Reproduction/evidence:** Source-confirmed; current session did not run a database end-to-end transition. Historical order auto-posting test proves only happy-path call behavior, not reversal/failure recovery.
- **Baseline scope:** Present on main/current. Phase 00 local warnings/diagnostics do not reverse entries.
- **Recommended solution / components:** Canonical payment/refund events, atomic update/outbox, post-once business event, linked reversing journal, payment-to-GL reconciliation and recovery dashboard.
- **Migration:** Event identity/outbox/reversal-link columns may be needed after design.
- **Dependencies/policy:** Revenue recognition trigger, delivery/fulfilment timing, refund/chargeback fee allocation, backdate policy.
- **Acceptance:** Paid then partial/full refund yields independently calculated net GL; cancellation before/after capture behaves distinctly; failed post visibly retries; paid order vs sales GL reconciliation has zero unexplained difference.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
