# F-066 — Dashboard home "Products" tile shows the page length (1) instead of the product total (12)

- **Severity:** Low  ·  **Category:** Software defect (UI aggregation)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Dashboard home
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-MET-DASH004
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Products tile = 12 (non-deleted menu items)

Observed: Tile shows 1. The dashboard requests page_size=1 and reads total at the payload root, but the API returns it under pagination.total, so the code falls back to items.length = 1

## Monetary / operational impact
0 (wrong operational count; misleads the owner about catalogue size)

## Where it was reproduced
Headed Chrome, tenant 4 owner: open /en/dashboard; compare the Products tile with GET /products?page=1&page_size=1 (pagination.total=12) and SQL (13 rows, 1 with status=deleted)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
