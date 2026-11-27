# F-056 — HR contracts and resignations accept end/last-working dates earlier than their start/resignation dates

- **Severity:** Low  ·  **Category:** Validation gap  ·  **Evidence status:** REPRODUCED
- **Affected modules:** HR contracts; HR resignations
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-CON-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 400 with a clear message

Observed: both return 201 and the invalid periods are stored (contract id 1, resignation id 1); duplicate contract number returns 400 with raw pq constraint text (F-031 class); attendance/leave/salary periods in the same module family are protected only by DB check constraints that surface as 500 | RUN-003: contract end 2026-10-01 < start 2026-11-01 stored (201); resignation last day before resignation date stored (201) (REPRODUCED). 

## Monetary / operational impact
0 directly; wrong notice/tenure calculations downstream

## Where it was reproduced
POST /hr/contracts with start_date 2026-11-01 and end_date 2026-10-01; POST /hr/resignations with resignation_date 2026-11-26 and last_working_date 2026-11-01

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
