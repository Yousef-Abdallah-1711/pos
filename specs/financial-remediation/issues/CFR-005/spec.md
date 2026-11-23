# CFR-005: Paid orders can suppress posting failures and have no proven cancellation/refund reversal

Status: PARTIALLY IMPLEMENTED — refund lifecycle and acceptance pending. Original severity: P1.

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

At main `3fa3702148771f39af5c9f7f25f377ecdd3f0a98`, `OrderRepository.UpdatePaymentStatusAs` locks the scoped order and commits paid state with a durable paid-order accounting intent. `CanTransitionOrderPaymentStatus` rejects a generic refunded transition and paid downgrade. PR [#39](https://github.com/Yousef-Abdallah-1711/pos/pull/39) merged at `d208489b4946b39eb832931fbc9ac929a2013069`. The repository has PostgreSQL tests for unbacked refunds, cancelled orders, tenant scope and concurrent paid/cancel transitions; historical focused passes are described in the [tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md), not freshly rerun here.

The historical tracker reports no restaurant-order provider refund lifecycle or linked accounting reversal. Subscription payment refunds do not demonstrate order refund accounting. That absence needs a complete current-path audit and provider contract; the existing guard prevents declaring an unsupported refund but does not perform one. Recognition, allocation, tax and closed-period decisions are pending.

## Definition of done

- Paid transition and intent are atomic; failures, replay, concurrency and cross-scope attempts are tested with synthetic orders and an independent journal/GL oracle.
- An authorized provider contract supports full and partial refund, timeout/unknown outcome, callback replay, concurrent requests and recovery after provider success with local failure.
- Refunds cannot exceed the remaining captured amount; each approved reversal links to its original and follows approved allocation, tax, recognition and period rules.
- Reconcile source payments and refunds to journal and GL, exercise API/operator flows, and retain current-SHA PostgreSQL, CI, review and main evidence. Do not mark complete while accounting decisions or provider contract are absent.
