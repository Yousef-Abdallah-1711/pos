# F-004 — POS shows 14% tax but server stores 10%

- **Severity:** High  ·  **Category:** Software defect (rate itself policy TAX-01)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** POS; Orders; Accounting
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-POS-001, EX-RUN003-AUD-POS-002
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Displayed total equals stored total

Observed: UI total 11400 (tax 1400); stored order tax 1000 total 11000 | RUN-003 (AUD-POS-001, order 38): POS label Tax (14%) vs stored tax 6,000 on subtotal 60,000 (10%); UI displayed value not legible in capture (overlay) - rate mismatch STILL REPRODUCIBLE by label vs stored; earlier evidence holds values. 

## Monetary / operational impact
400 per 10000 order (customer shown more than recorded)

## Where it was reproduced
POS add item price 10000 then Place Order

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
