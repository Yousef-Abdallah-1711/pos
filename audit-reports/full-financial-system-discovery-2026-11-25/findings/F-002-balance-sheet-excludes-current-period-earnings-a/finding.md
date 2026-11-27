# F-002 — Balance Sheet excludes current-period earnings and is UNRECONCILED by exactly the net result

- **Severity:** High  ·  **Category:** Software defect (policy EARN-01 related)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Financial reports
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-MONTH-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A = L + E

Observed: equation_difference -2000.00 (net loss); equity shows 100000 only | RUN-003 clean-tenant month (AUD-MONTH): balance sheet equation_balanced=false, equation_difference -27,710.00 = net loss exactly (REPRODUCED). 

## Monetary / operational impact
2000 per period (signed residual -2000)

## Where it was reproduced
S13 loss scenario then open Balance Sheet

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
