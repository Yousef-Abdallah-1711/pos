# F-047 — Recurring journal template can be run repeatedly for the same period and posts a duplicate journal each time

- **Severity:** Medium  ·  **Category:** Software defect (idempotency)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Accounting - recurring journals
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-ACC-009
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Second run for the same period rejected or returns the first journal

Observed: Run 1 -> JE id 98 (rent +700); run 2 -> JE id 99 (rent +700 again): rent expense 1,400 for one period; next_run_date already advanced to 2026-12-26 | RUN-003: recurring run twice -> journals 98 and 99, rent +700 each (CONFIRMED accounting_advanced_repo.go:396-420). 

## Monetary / operational impact
700 duplicated per extra run in test

## Where it was reproduced
Create template (auto_post) and POST /accounting/recurring-journals/{id}/run twice on 2026-11-26

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
