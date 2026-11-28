# F-043 — KPI dashboard mixes tax-inclusive order revenue with GL expenses, and cash position ignores GL cash accounts

- **Severity:** Medium  ·  **Category:** Software defect (report definition)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Financial KPI dashboard
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Revenue/expenses/net profit tie to P&L; cash position ties to GL cash and bank accounts

Observed: KPI revenue 319,660 (paid orders incl tax) vs P&L revenue 305,300; net_profit 181,540 vs P&L 167,180 (+14,360); cash_position 50,000 (bank_accounts_ledger) vs GL cash+bank 276,270.01 | RUN-003: KPI net profit 181,540 vs P&L 167,180; cash_position 50,000 vs GL cash+bank 276,270.01 (CONFIRMED accounting_advanced_repo.go:838-880). 

## Severity rationale
Rated **Medium** (category: Software defect (report definition); evidence: REPRODUCED). Basis recorded in the register: 14,360 profit overstatement; 226,270.01 cash understatement in this dataset. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Financial KPI dashboard. 14,360 profit overstatement; 226,270.01 cash understatement in this dataset

## Technical impact
CONFIRMED (source): backend/internal/repository/accounting_advanced_repo.go:838-880 revenue=SUM(orders.total_amount paid, created_at), expenses=GL expense lines, cash_position=SUM(bank_accounts_ledger.current_balance)

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
14,360 profit overstatement; 226,270.01 cash understatement in this dataset

## Where it was reproduced
GET /reports/financial/kpis (tenant 4, window 2026-10-28..2026-11-26)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
