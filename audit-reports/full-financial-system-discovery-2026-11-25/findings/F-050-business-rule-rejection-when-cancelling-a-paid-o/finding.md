# F-050 — Business-rule rejection when cancelling a paid order via PUT /admin/orders/{id}/status returns HTTP 500 instead of a 4xx

- **Severity:** Low  ·  **Category:** Software defect (error handling)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Orders API
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-HC10, EX-RUN003-AUD-SAAS-006, EX-RUN003-AUD-OT-003
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 409/400 with the same message as DELETE

Observed: 500 "Failed to update order status" (DELETE returns 400 with the explicit reason) | RUN-003: PUT status cancelled on paid order -> 500.  | RUN-003: PUT status cancelled on paid order -> 500. Same class: unsigned Paymob webhook -> 500 and refund of unknown payment -> 500 (AUD-SAAS-MATRIX PAY-02/PAY-03). 

## Monetary / operational impact
0

## Where it was reproduced
PUT /admin/orders/{paid id}/status {status:cancelled}

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
