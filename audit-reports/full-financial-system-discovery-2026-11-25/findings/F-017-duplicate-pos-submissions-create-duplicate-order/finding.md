# F-017 — Duplicate POS submissions create duplicate orders; Idempotency-Key is ignored on order creation

- **Severity:** Medium  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** POS; Orders; Inventory
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-REL-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: One order

Observed: 5 and 3 orders created; stock decremented each | RUN-003 REL-001: same order POST twice with identical Idempotency-Key creates 2 orders (REPRODUCED). 

## Severity rationale
Rated **Medium** (category: Software defect; evidence: REPRODUCED). Basis recorded in the register: duplicate revenue risk. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: POS; Orders; Inventory. duplicate revenue risk

## Technical impact
CONFIRMED (source search): the only consumer of Idempotency-Key is backend/internal/handler/http/order_cancellation_handler.go:36 (item cancel); order creation handlers (public_order_handler.go CreateOrder, registerPOSOrderRoute backend/cmd/api/pos_routes.go:16) ignore the header although CORS advertises it; REL-001 re-run created 2 orders

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
duplicate revenue risk

## Where it was reproduced
5 identical parallel POST and 3 with same key

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
