# F-009 — Menu Profitability and Break-even endpoints return HTTP 500

- **Severity:** Medium  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Financial reports advanced
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Reports render

Observed: HTTP 500 for both (cost_price column absent) | RUN-003: GET /reports/financial/menu-profitability and /break-even still 500. ROOT CAUSE CONFIRMED: repository/accounting_advanced_repo.go:976 and :1045 select p.cost_price, a column that does not exist on products (the column is cost) -> SQL error mapped to generic 500 (accounting_advanced_handler.go:782). 

## Severity rationale
Rated **Medium** (category: Software defect; evidence: REPRODUCED). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Financial reports advanced. n/a

## Technical impact
CONFIRMED (source + SQL error reproduced): backend/internal/repository/accounting_advanced_repo.go:976 and :1045 select p.cost_price, a column that does not exist on products (the column is cost); error mapped to generic 500 at accounting_advanced_handler.go:782

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
Open advanced insights and Refresh

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
