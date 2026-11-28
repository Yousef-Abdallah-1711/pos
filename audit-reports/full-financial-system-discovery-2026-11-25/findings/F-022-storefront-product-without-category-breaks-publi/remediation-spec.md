# F-022 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Storefront: product without category breaks public menu API (500); storefront UI still renders 0 items with a valid category.

## 2. Current behavior (verified)
Uncategorised product (category_id NULL, allowed by the product form) makes GET public menu/products return 500 for the whole tenant; after assigning a category the API returns the product but the storefront page still shows 0 items; new-tenant homepage shows Unable to Load Homepage | RUN-003: (1) REAL-UI reproduction: product form without a category sends category_id 0 -> POST /products 201 -> public menu/products API 500 for the whole tenant; (2) storefront UI zero items with default isolated CORS (environment), 3 items after the wildcard origin was allowed. 

## 3. Required behavior
Home and menu render

## 4. Reproduction conditions
Open storefront subdomain /en and /en/menu

## 5. Verified evidence
- runtime/18 | evidence/network/RUN-003/PH-06/AUD-UIFORM-F022-public-menu-500.txt | evidence/network/RUN-003/PH-06/AUD-UIFORMS.json | evidence/network/RUN-003/PH-08/AUD-STORE-002-public-api.txt | evidence/network/RUN-003/PH-08/AUD-STORE-003.json | evidence/network/RUN-003/PH-08/AUD-STORE-004.json | tests/runtime/ph06_ui_forms.mjs

## 6. Root cause and confidence
CONFIRMED — 63 / product listing with NULL category); UI zero-items cause CONFIRMED in RUN-003: with CORS_ORIGINS lacking a wildcard (the middleware supports http://*.localhost:13113, middleware/cors.go:66-92) browsers block the storefront subdomain origin; once the isolated API was recreated with that origin allowed, see AUD-STORE-004. The isolated compose default does not include the wildcard (environment configuration)

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/public_menu_handler.go`:63, `backend/internal/middleware/cors.go`:66-92, `audit-reports/full-financial-system-discovery-2026-11-25/tests/runtime/docker-compose.audit-persistent.yml`:140
- Frontend: `frontend/apps/restaurant-website/app/[locale]`, `frontend/apps/dashboard/src/app/[locale]/dashboard/products`
- API contracts: POST /api/v1/products: category_id required (400 otherwise); GET /api/v1/public/menu: uncategorised products appear under a synthetic section (additive)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: storefront orders blocked
- Operational: Storefront; Products workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: (1) make the product listing used by the public menu robust to a missing category: in the repository behind backend/internal/handler/http/public_menu_handler.go:63 use a LEFT JOIN on categories and group uncategorised products under a synthetic "Other" section, so one bad row can never return 500 for the whole tenant; (2) require a category when creating a product (backend validation plus the dashboard product form, which currently posts category_id 0); (3) for the storefront UI, document the required CORS origin pattern (http://*.localhost:<port> in non-production, the real tenant domains in production) in the deployment guide and add a startup warning when the storefront origin is missing from CORS_ORIGINS.

## 10. Alternative designs
A: LEFT JOIN + required category + deployment warning (recommended). B: LEFT JOIN only (hides the data-quality problem). C: wildcard CORS in production (rejected: widens the attack surface).

## 11. Required policy decisions
None (data-quality rule); naming of the fallback section is a product decision.

## 12. Required data migrations
Optional: NOT NULL on products.category_id after back-filling existing rows (list rows first; do not guess a category).

## 13. Historical data implications
Existing products with NULL or 0 category are listed for the owner to categorise (tenant 4 product 18 in the evidence).

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none

## 21. Risks
- Regression in modules sharing the same source files (audit-reports/full-financial-system-discovery-2026-11-25/tests/runtime/docker-compose.audit-persistent.yml, backend/internal/handler/http/public_menu_handler.go, backend/internal/middleware/cors.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Storefront smoke test after registration + product create (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A tenant with one uncategorised product still returns 200 from the public menu and products endpoints, listing that product under "Other".
- POST /products without a valid category_id returns 400 with a field error.
- Starting the API with a storefront origin absent from CORS_ORIGINS logs a warning that names the missing origin.
