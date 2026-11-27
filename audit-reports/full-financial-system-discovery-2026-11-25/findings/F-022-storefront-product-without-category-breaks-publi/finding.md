# F-022 — Storefront: product without category breaks public menu API (500); storefront UI still renders 0 items with a valid category

- **Severity:** High  ·  **Category:** Software defect  ·  **Evidence status:** REPRODUCED (API root cause); UI cause UNKNOWN
- **Affected modules:** Storefront; Products
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-STORE-004, EX-RUN003-AUD-UIFORM-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: Home and menu render

Observed: Uncategorised product (category_id NULL, allowed by the product form) makes GET public menu/products return 500 for the whole tenant; after assigning a category the API returns the product but the storefront page still shows 0 items; new-tenant homepage shows Unable to Load Homepage | RUN-003: (1) REAL-UI reproduction: product form without a category sends category_id 0 -> POST /products 201 -> public menu/products API 500 for the whole tenant; (2) storefront UI zero items with default isolated CORS (environment), 3 items after the wildcard origin was allowed. 

## Monetary / operational impact
storefront orders blocked

## Where it was reproduced
Open storefront subdomain /en and /en/menu

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
