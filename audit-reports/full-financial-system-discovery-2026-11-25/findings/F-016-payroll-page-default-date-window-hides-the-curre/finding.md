# F-016 — Payroll page default date window hides the current month's salary records

- **Severity:** Low  ·  **Category:** UX defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** HR payroll
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Current records visible

Observed: 0 records until To-date changed | RUN-003 (AUD-EVD-F016): default window start_date=2026-10-27,end_date=2026-11-26 returns total 0 and UI says 'Create your first salary record', while salaries row id 1 (tenant 4, period 2026-11-01..2026-11-30, paid) exists (STILL REPRODUCIBLE). 

## Monetary / operational impact
n/a

## Where it was reproduced
Open payroll page with record period ending after today

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
