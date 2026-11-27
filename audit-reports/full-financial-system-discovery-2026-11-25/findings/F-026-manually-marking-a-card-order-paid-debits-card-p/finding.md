# F-026 — Manually marking a card order paid debits Card Payment Clearing with no clearing/settlement flow

- **Severity:** Medium  ·  **Category:** Missing capability  ·  **Evidence status:** REPRODUCED (RUN-005 API-DIAGNOSTIC)
- **Affected modules:** Accounting
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN005-OBSERVED-REPRO
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Clearing cleared on settlement

Observed: 1113 balance 11000 persists RUN-005: card order 5795 (110.00) marked paid -> journal Dr 1113 Card Payment Clearing 110.00 / Cr 4100 100.00 / Cr 2120 10.00; no settlement route exists.

## Severity rationale
Rated **Medium** (category: Missing capability; evidence: REPRODUCED (RUN-005 API-DIAGNOSTIC)). Basis recorded in the register: 11000 stuck. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Accounting. 11000 stuck

## Technical impact
CONFIRMED (RUN-004 source): card payments resolve the debit account to 1113 Card Payment Clearing (auto_posting_engine.go:615-616; domain/auto_posting.go:37; seeded in migration 156:35); no code or migration references 1113 anywhere else, so no settlement/clearing flow exists

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
11000 stuck

## Where it was reproduced
Mark card order paid

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
