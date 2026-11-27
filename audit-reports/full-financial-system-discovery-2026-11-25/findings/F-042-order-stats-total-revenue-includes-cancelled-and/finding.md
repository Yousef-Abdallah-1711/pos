# F-042 — Order stats "total revenue" includes cancelled and unpaid orders and "delivered today" counts every order created today

- **Severity:** Medium  ·  **Category:** Software defect (report definition)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Orders dashboard; Order stats API
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Revenue excludes cancelled (and unpaid per chosen policy); delivered_today counts delivered orders

Observed: total_revenue 529,100 = SUM of ALL 41 orders incl 3 cancelled (11,440) and 21 unpaid; delivered_today 30 while delivered_count 0 (counts any order created today) | RUN-003: order stats total_revenue 529,100 = all 41 orders incl 3 cancelled; delivered_today 30 vs delivered_count 0 (CONFIRMED order_repo.go:859-870). 

## Monetary / operational impact
529,100 shown vs 319,660 paid tax-inclusive vs 305,300 GL revenue

## Where it was reproduced
GET /admin/orders/stats on tenant with cancelled, unpaid and non-delivered orders (RUN-003 batch capture)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
