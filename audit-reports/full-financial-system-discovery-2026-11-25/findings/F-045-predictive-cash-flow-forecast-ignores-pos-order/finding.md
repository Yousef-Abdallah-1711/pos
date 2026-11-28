# F-045 — Predictive cash flow forecast ignores POS/order receipts and is built only from customer receipts and supplier payments

- **Severity:** Medium  ·  **Category:** Software defect (report definition)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Financial hub predictive cash flow
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: History reflects actual cash movements (POS cash/card receipts included)

Observed: 2026-11 history cash_in 13,000 / cash_out 32,400 -> forecast Dec -18,430 while 319,660 of paid POS orders and GL cash +276,270 exist | RUN-003: PCF Nov cash_in 13,000/cash_out 32,400 vs paid POS orders 319,660; forecast negative (CONFIRMED accounting_advanced_repo.go:1105-1175). 

## Severity rationale
Rated **Medium** (category: Software defect (report definition); evidence: REPRODUCED). Basis recorded in the register: Forecast understates inflows by 319,660 in this dataset. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Financial hub predictive cash flow. Forecast understates inflows by 319,660 in this dataset

## Technical impact
CONFIRMED (source): accounting_advanced_repo.go:1105-1175 GetPredictiveCashFlow aggregates only payments_received and payments_made

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
Forecast understates inflows by 319,660 in this dataset

## Where it was reproduced
GET /reports/financial/predictive-cash-flow on tenant 4

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
