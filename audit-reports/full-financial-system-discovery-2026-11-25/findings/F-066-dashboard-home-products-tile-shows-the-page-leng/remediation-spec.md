# F-066 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Dashboard home "Products" tile shows the page length (1) instead of the product total (12).

## 2. Current behavior (verified)
Tile shows 1. The dashboard requests page_size=1 and reads total at the payload root, but the API returns it under pagination.total, so the code falls back to items.length = 1

## 3. Required behavior
Products tile = 12 (non-deleted menu items)

## 4. Reproduction conditions
Headed Chrome, tenant 4 owner: open /en/dashboard; compare the Products tile with GET /products?page=1&page_size=1 (pagination.total=12) and SQL (13 rows, 1 with status=deleted)

## 5. Verified evidence
- evidence/network/RUN-004/PH-12/AUD-R4-MET-CAPTURE.json
- evidence/network/RUN-004/PH-12/AUD-R4-MET2-CAPTURE.json
- evidence/screenshots/RUN-004/PH-12/AUD-R4-MET-dash.png

## 6. Root cause and confidence
CONFIRMED — frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx:173-178 extractTotalFromCollection reads direct?.total (undefined) not pagination.total; request with page_size 1 at :314 and call at :329

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx`:173, `frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx`:329, `frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx`:389, `frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx`:173-178
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/products`
- API contracts: optional GET /api/v1/products/count (new)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0 (wrong operational count; misleads the owner about catalogue size)
- Operational: Dashboard home workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: read the product total from the API contract. In frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx (extractTotalFromCollection :173-178, call :329) read pagination.total (fall back to total, then items.length only when no total exists) and add a unit test with the real /products response shape; consider a lightweight GET /products/count endpoint so the home page does not fetch rows.

## 10. Alternative designs
A: fix the extractor (recommended, smallest). B: add a count endpoint (cleaner long term).

## 11. Required policy decisions
None.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

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
- Regression in modules sharing the same source files (frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Dashboard tile equals the API total when there is more than one product (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- With 12 non-deleted products the tile shows 12.
- The unit test fails if the extractor reads the wrong key.
