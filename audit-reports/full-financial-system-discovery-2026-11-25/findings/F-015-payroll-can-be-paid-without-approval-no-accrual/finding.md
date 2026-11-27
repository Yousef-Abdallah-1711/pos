# F-015 — Payroll can be paid without approval; no accrual; no partial payment

- **Severity:** Medium  ·  **Category:** Missing capability / policy  ·  **Evidence status:** REPRODUCED
- **Affected modules:** HR payroll; Accounting
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-SAL-002
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Approval gate; accrual per PAY-01

Observed: Paid with is_approved=false; journal only at payment as direct expense | RUN-003 SQL re-verification (AUD-EVD-F015): salaries id 1 and 2 are status=paid with is_approved=false; journals JE-2026-00007/00049 are Dr 5210 5,000 / Cr 1111 5,000 posted at payment only; no accrual journal exists (STILL REPRODUCIBLE; UI re-run pending T-P10-003). 

## Monetary / operational impact
payroll payable 2000 absent

## Where it was reproduced
Create salary then Mark Paid

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
