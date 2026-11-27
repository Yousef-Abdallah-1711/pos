# F-073 — Bank reconciliation lines accept a journal_entry_id that belongs to another tenant (only a foreign key is checked)

- **Severity:** Medium  ·  **Category:** Security (tenant isolation integrity)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Accounting > Bank reconciliation
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-BANKRECON
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 400/404 because the journal is not visible to the tenant

Observed: HTTP 200; the line is stored in the tenant-4 session referencing tenant 2 journal 33; a nonexistent id returns 400 only because of the foreign key, with the raw constraint text in the error body

## Monetary / operational impact
0 direct; cross-tenant reference integrity and error-text leakage (F-025 class)

## Where it was reproduced
Tenant 4 owner adds a reconciliation line with journal_entry_id 33, a journal of tenant 2 (API-DIAGNOSTIC)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
