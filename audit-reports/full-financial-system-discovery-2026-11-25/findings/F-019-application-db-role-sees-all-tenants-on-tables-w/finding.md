# F-019 — Application DB role sees all tenants on tables without RLS and can alter paid orders/payment amounts

- **Severity:** Medium  ·  **Category:** Security / defense in depth  ·  **Evidence status:** REPRODUCED (runtime catalog)
- **Affected modules:** Platform
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** none recorded in 34b (evidence is in the register/evidence files)
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Cross-tenant rows hidden; paid records immutable

Observed: All 4 tenants' orders/journals/GL/COA/users visible; UPDATE orders/payments_received succeeded (rolled back); ledger tables protected by grants

## Severity rationale
Rated **Medium** (category: Security / defense in depth; evidence: REPRODUCED (runtime catalog)). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Platform. n/a

## Technical impact
CONFIRMED (RUN-004 runtime SQL): only 13 of the 70 public tables with a tenant_id column have row-level security (RLS migrations 164-171 are the only ones); the application role pos_app has UPDATE on orders and payments and no immutability trigger exists on payments (orders carry only stock-transition triggers); journal_entry_lines is not updatable by pos_app

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: yes, see category. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
Connect as pos_app without tenant context

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
