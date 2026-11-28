# F-003 — Cash Flow statement is accrual profit plus code-prefix guesses and does not tie to cash

- **Severity:** High  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Financial reports
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-MET-FIN, EX-RUN003-AUD-CF-TRACE, EX-RUN003-AUD-ACC-004, EX-RUN003-AUD-MONTH-001, EX-RUN004-MET-CF002
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Net cash change equals change in cash accounts

Observed: Operating = net profit (-2000); net cash flow 98000 vs GL cash+bank 93000 | RUN-003 (AUD-MET-FIN): cash-flow net 82,710.01 vs GL cash+bank 115,010.01, residual -32,300.00 (STILL REPRODUCIBLE).  | RUN-003 earlier run: cash-flow 82,710.01 vs GL cash+bank 115,010.01 (residual 32,300.00). Re-measured after more activity: report 172,710.01 vs GL 1110+1111 = 280,010.01, residual 107,300.00. ROOT CAUSE CONFIRMED (AUD-CF-TRACE): usecase/financial_report_usecase.go:366-399 GenerateCashFlow is an indirect-method fragment that never reads cash/bank accounts and omits every working-capital movement (AR -6,520; inventory +96,000; input VAT -2,800; card clearing -11,000; AP +800; VAT payable +30,820 = +107,300); investing = '12%' debit-credit has inverted sign for asset purchases; financing excludes 21xx/2120 liabilities. 

## Severity rationale
Rated **High** (category: Software defect; evidence: REPRODUCED). Basis recorded in the register: 5000 mismatch (this dataset). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Financial reports. 5000 mismatch (this dataset)

## Technical impact
CONFIRMED (source + RUN-003 decomposition): backend/internal/usecase/financial_report_usecase.go:366-399 GenerateCashFlow is an indirect-method fragment: operating = P&L (revenue/expense credit-debit), investing = 12xx debit-credit (sign inverted), financing = equity or 22xx credit-debit; never reads cash/bank accounts and omits working-capital movements (see AUD-CF-TRACE.txt)

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
5000 mismatch (this dataset)

## Where it was reproduced
Open Cash Flow after S13

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
