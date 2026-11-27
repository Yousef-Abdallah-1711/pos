# F-061 — A paid order with total 0 (free item) creates an accounting event that can never post ("at least two journal lines are required")

- **Severity:** Medium  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** POS/Orders; Accounting event outbox
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-MD-002
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Zero-value sale completes without a journal (or posts a zero-line memo) and the event is completed

Observed: order 73 total 0.00 paid: outbox status retryable attempts 1/8 last_error "at least two journal lines are required"; it will exhaust retries and appear as a failed event needing manual redrive; no journal exists | RUN-003: zero-total paid order 73 -> outbox retryable, error at least two journal lines are required (CONFIRMED auto_posting_engine.go:535). 

## Monetary / operational impact
0 value, but permanent failed-event noise hides real failures and needs manual redrive

## Where it was reproduced
Sell a product with price 0 (product creation allows it), mark the order paid, read accounting_event_outbox

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
