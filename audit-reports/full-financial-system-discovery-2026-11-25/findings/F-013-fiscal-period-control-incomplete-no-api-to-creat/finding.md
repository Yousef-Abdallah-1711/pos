# F-013 — Fiscal-period control incomplete: no API to create/close periods; void and draft creation allowed in a closed period

- **Severity:** High  ·  **Category:** Software defect / missing capability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Accounting
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-ACC-006
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Closed period blocks all postings including voids

Observed: Post blocked; draft create allowed; void of posted entry allowed (capital 100000 reversed) | RUN-003 (AUD-ACC-FLOW JE-13/JE-14): journals dated 2019-01-15 and 2030-06-01 post successfully and auto-create fiscal_period 2019-01/2030-06; GET /accounting/fiscal-periods is 404 (no period API). 

## Severity rationale
Rated **High** (category: Software defect / missing capability; evidence: REPRODUCED). Basis recorded in the register: 100000 reversed in closed period. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Accounting. 100000 reversed in closed period

## Technical impact
CONFIRMED (source) for the missing-period part: backend/internal/repository/journal_entry_repo.go:813-842 ensureFiscalPeriodOpenTx returns nil when no fiscal_periods row matches, so any date posts; only status=closed blocks; no API creates/closes periods (GET /accounting/fiscal-periods is 404)

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
100000 reversed in closed period

## Where it was reproduced
Insert closed period (SQL fixture) then post / create / void

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
