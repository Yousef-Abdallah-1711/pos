# F-056 — HR contracts and resignations accept end/last-working dates earlier than their start/resignation dates

- **Severity:** Low  ·  **Category:** Validation gap  ·  **Evidence status:** REPRODUCED
- **Affected modules:** HR contracts; HR resignations
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-CON-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 400 with a clear message

Observed: both return 201 and the invalid periods are stored (contract id 1, resignation id 1); duplicate contract number returns 400 with raw pq constraint text (F-031 class); attendance/leave/salary periods in the same module family are protected only by DB check constraints that surface as 500 | RUN-003: contract end 2026-10-01 < start 2026-11-01 stored (201); resignation last day before resignation date stored (201) (REPRODUCED). 

## Severity rationale
Rated **Low** (category: Validation gap; evidence: REPRODUCED). Basis recorded in the register: 0 directly; wrong notice/tenure calculations downstream. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: HR contracts; HR resignations. 0 directly; wrong notice/tenure calculations downstream

## Technical impact
CONFIRMED (RUN-004 source): contract create/update validate only the date format (contract_handler.go:165-170, 221-226); no end >= start rule; resignation create has the same gap (resignation_repo.go:117-143)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0 directly; wrong notice/tenure calculations downstream

## Where it was reproduced
POST /hr/contracts with start_date 2026-11-01 and end_date 2026-10-01; POST /hr/resignations with resignation_date 2026-11-26 and last_working_date 2026-11-01

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
