# F-012 — At least six revenue/profit definitions disagree and are unlabeled

- **Severity:** High  ·  **Category:** Documentation/design inconsistency  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Reports; Dashboard
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Reports state basis or reconcile

Observed: Sales 121000; Profitability 110000; KPI 22000; GL revenue 20000; KPI net profit 6000 vs P&L 4000; KPI cash position 0 vs GL cash

## Severity rationale
Rated **High** (category: Documentation/design inconsistency; evidence: REPRODUCED). Basis recorded in the register: up to 101000 spread. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Reports; Dashboard. up to 101000 spread

## Technical impact
CONFIRMED (RUN-004 source): three independent queries define revenue differently: Sales = SUM(o.total_amount) tax-inclusive (report_handler.go:2177), Profitability = SUM(oi.total_price) net of tax for non-cancelled orders (report_handler.go:2367), P&L = GL revenue-type accounts (financial_report_usecase.go:208); no shared definition or label; numeric bridge in evidence/network/RUN-003/PH-12/AUD-REV-TRACE.txt

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
up to 101000 spread

## Where it was reproduced
Open Sales Profitability KPI P&L after mixed paid/unpaid orders

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
