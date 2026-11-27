# F-054 — Leave approval does not enforce entitlement and total_days is not validated against the requested dates

- **Severity:** Medium  ·  **Category:** Software defect (business rule)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** HR leaves; leave-balance report
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-LEV-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Reject or flag leave beyond remaining entitlement; total_days derived from dates

Observed: 30-day leave approved; report shows used 30, remaining 0 (clamped) vs allowance 21; 3-day range stored total_days 40 | RUN-003: 30-day annual leave approved over allowance 21; total_days 40 stored for a 3-day range (REPRODUCED). 

## Monetary / operational impact
Unpaid-leave exposure / payroll deduction not triggered

## Where it was reproduced
Create annual leave 2026-12-01..12-30 total_days 30 for an employee with allowance 21 and approve; create a 3-day range with total_days 40

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
