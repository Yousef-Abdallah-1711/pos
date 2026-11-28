# F-031 — Business-rule and dependency failures surface as wrong status codes and leak internals

- **Severity:** Medium  ·  **Category:** Software defect (API contract)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Orders; Receipts; Platform
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-ACC-001, EX-RUN003-AUD-ATT-001, EX-RUN003-AUD-LEV-001, EX-RUN003-AUD-CON-001, EX-RUN003-AUD-MD-001, EX-RUN003-AUD-UIFORM-001, EX-RUN003-AUD-UIFORM-003, EX-RUN004-FORMMATRIX, EX-RUN004-FORMMATRIX2
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 409/503 with generic message

Observed: Paid-order cancel 500 (item cancel 409); receipt during outage 400 with raw role error pq role pos_accounting not permitted; DB down: /health 200 degraded and authenticated calls 403 TENANT_ACCESS_BLOCKED contact support | RUN-003 (AUD-ACC-FLOW ACC-02): duplicate account code returns raw DB text pq: duplicate key value violates unique constraint chart_of_accounts_tenant_id_restaurant_id_account_code_key (REPRODUCED). 

## Severity rationale
Rated **Medium** (category: Software defect (API contract); evidence: REPRODUCED). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Orders; Receipts; Platform. n/a

## Technical impact
CONFIRMED (RUN-004 source): error mapping is done by passing err text to the client (see F-025 sites) and returning 500 for constraint violations (e.g. main.go:1302; hr handlers rely on DB check constraints)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
Cancel a paid order; receipt during accounting-role outage; DB outage

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
