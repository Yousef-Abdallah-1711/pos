# F-066 — Dashboard home "Products" tile shows the page length (1) instead of the product total (12)

- **Severity:** Low  ·  **Category:** Software defect (UI aggregation)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Dashboard home
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-MET-DASH004
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Products tile = 12 (non-deleted menu items)

Observed: Tile shows 1. The dashboard requests page_size=1 and reads total at the payload root, but the API returns it under pagination.total, so the code falls back to items.length = 1

## Severity rationale
Rated **Low** (category: Software defect (UI aggregation); evidence: REPRODUCED). Basis recorded in the register: 0 (wrong operational count; misleads the owner about catalogue size). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Dashboard home. 0 (wrong operational count; misleads the owner about catalogue size)

## Technical impact
CONFIRMED (source + captured response): frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx:173-178 extractTotalFromCollection reads direct?.total (undefined) not pagination.total; request with page_size 1 at :314 and call at :329

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0 (wrong operational count; misleads the owner about catalogue size)

## Where it was reproduced
Headed Chrome, tenant 4 owner: open /en/dashboard; compare the Products tile with GET /products?page=1&page_size=1 (pagination.total=12) and SQL (13 rows, 1 with status=deleted)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
