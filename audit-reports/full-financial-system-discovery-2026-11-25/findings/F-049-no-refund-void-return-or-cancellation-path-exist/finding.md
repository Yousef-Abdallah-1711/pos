# F-049 — No refund, void, return or cancellation path exists for a paid order (missing financial capability)

- **Severity:** High  ·  **Category:** Missing capability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Orders; Payments; Accounting (revenue/COGS/VAT/inventory reversal)
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-HC10
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Controlled refund/void producing reversing journals, stock return and VAT reversal

Observed: DELETE -> 400 "cannot cancel paid or refunded order before a supported refund and accounting reversal workflow exists"; payment-status refunded/failed/pending from paid -> 400 invalid transition; voided -> "invalid payment status"; status=cancelled -> 500; paid revenue/COGS stay posted; cash 307,330.01 unchanged | RUN-003: paid order cannot be cancelled/refunded/voided through any endpoint; explicit guard message (REPRODUCED). 

## Monetary / operational impact
Any customer refund cannot be recorded; every refunded sale overstates revenue, VAT and cash

## Where it was reproduced
Pay an order, then try DELETE /admin/orders/{id}, PUT payment-status refunded/voided/failed/pending, PUT status cancelled, item cancel

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
