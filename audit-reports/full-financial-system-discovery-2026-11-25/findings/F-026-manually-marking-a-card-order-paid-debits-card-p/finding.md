# F-026 — Manually marking a card order paid debits Card Payment Clearing with no clearing/settlement flow

- **Severity:** Medium  ·  **Category:** Missing capability  ·  **Evidence status:** OBSERVED
- **Affected modules:** Accounting
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Clearing cleared on settlement

Observed: 1113 balance 11000 persists

## Monetary / operational impact
11000 stuck

## Where it was reproduced
Mark card order paid

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
