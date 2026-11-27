# F-068 — Social-insurance CSV export lists an employee once per active insurance policy, repeating base salary and deductions

- **Severity:** Medium  ·  **Category:** Software defect (export)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** HR reports > Social insurance export
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-MET-HRR007
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: One row per employee salary record

Observed: 1 salary row (SQL) produced 2 data rows after the second active policy; both repeat base_salary 5,000.00 and the deductions, so any total over the file double-counts

## Monetary / operational impact
Proportional: each extra active policy repeats the employee base salary and deductions (here 5,000.00 base salary listed twice)

## Where it was reproduced
Tenant 4, month 2026-11: export with one active policy, add a second active policy for the same employee (API-DIAGNOSTIC), export again

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
