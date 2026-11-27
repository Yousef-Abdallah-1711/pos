# F-031 — Business-rule and dependency failures surface as wrong status codes and leak internals

- **Severity:** Medium  ·  **Category:** Software defect (API contract)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Orders; Receipts; Platform
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-ACC-001, EX-RUN003-AUD-ATT-001, EX-RUN003-AUD-LEV-001, EX-RUN003-AUD-CON-001, EX-RUN003-AUD-MD-001, EX-RUN003-AUD-UIFORM-001, EX-RUN003-AUD-UIFORM-003, EX-RUN004-FORMMATRIX, EX-RUN004-FORMMATRIX2
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 409/503 with generic message

Observed: Paid-order cancel 500 (item cancel 409); receipt during outage 400 with raw role error pq role pos_accounting not permitted; DB down: /health 200 degraded and authenticated calls 403 TENANT_ACCESS_BLOCKED contact support | RUN-003 (AUD-ACC-FLOW ACC-02): duplicate account code returns raw DB text pq: duplicate key value violates unique constraint chart_of_accounts_tenant_id_restaurant_id_account_code_key (REPRODUCED). 

## Monetary / operational impact
n/a

## Where it was reproduced
Cancel a paid order; receipt during accounting-role outage; DB outage

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
