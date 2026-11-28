# F-002 — Balance Sheet excludes current-period earnings and is UNRECONCILED by exactly the net result

- **Severity:** High  ·  **Category:** Software defect (policy EARN-01 related)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Financial reports
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-MONTH-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A = L + E

Observed: equation_difference -2000.00 (net loss); equity shows 100000 only | RUN-003 clean-tenant month (AUD-MONTH): balance sheet equation_balanced=false, equation_difference -27,710.00 = net loss exactly (REPRODUCED). 

## Severity rationale
Rated **High** (category: Software defect (policy EARN-01 related); evidence: REPRODUCED). Basis recorded in the register: 2000 per period (signed residual -2000). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Financial reports. 2000 per period (signed residual -2000)

## Technical impact
CONFIRMED (source, RUN-003 read): backend/internal/usecase/financial_report_usecase.go:285-306 balance-sheet query selects only account_type IN (asset, liability, equity); revenue/expense accounts never contribute, so current-period earnings are not added to equity (BS difference 82,700 measured earlier)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
2000 per period (signed residual -2000)

## Where it was reproduced
S13 loss scenario then open Balance Sheet

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
