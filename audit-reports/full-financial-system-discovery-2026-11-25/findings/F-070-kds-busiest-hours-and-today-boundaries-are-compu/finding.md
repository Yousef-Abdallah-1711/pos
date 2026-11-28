# F-070 — KDS busiest-hours and "today" boundaries are computed in UTC, not the restaurant local time

- **Severity:** Low  ·  **Category:** Software defect (time zone)  ·  **Evidence status:** REPRODUCED (KDS only; other reports unverified)
- **Affected modules:** KDS analytics and stats
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-MET-KDS003
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Hours and the "today" cut-off follow the restaurant time zone (Egypt UTC+2/+3)

Observed: API hours equal the UTC buckets exactly (1:7, 17:4, 19:4, 20:18, 22:11); the Cairo buckets differ (0:11, 3:7, 19:4, 21:4, 22:18); kds_queue.received_at is timestamp without time zone; no tenant time-zone setting was found

## Severity rationale
Rated **Low** (category: Software defect (time zone); evidence: REPRODUCED (KDS only; other reports unverified)). Basis recorded in the register: 0 (scheduling signal shifted 2-3 hours); other date-bounded reports may share the cut-off issue (CANDIDATE, not tested). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: KDS analytics and stats. 0 (scheduling signal shifted 2-3 hours); other date-bounded reports may share the cut-off issue (CANDIDATE, not tested)

## Technical impact
CONFIRMED for KDS: kds_repo.go:418-430 EXTRACT(HOUR FROM received_at) and CURRENT_DATE at :314-317; DB timezone is UTC

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0 (scheduling signal shifted 2-3 hours); other date-bounded reports may share the cut-off issue (CANDIDATE, not tested)

## Where it was reproduced
Tenant 4: compare /kds/analytics busiest_hours with SQL bucketed by UTC and by Africa/Cairo for the same kds_queue rows

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
