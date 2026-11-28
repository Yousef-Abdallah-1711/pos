# F-065 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Category creation fails in the UI whenever a display order is typed: the form sends display_order as a JSON string and the API answers 400 "Invalid request body".

## 2. Current behavior (verified)
Name only -> POST /categories 201 with display_order 0; with display order -> POST body {"display_order":"7"} -> 400 {"error":"Invalid request body"}; no field-level message is shown to the user | RUN-003 real-UI reproduction: display order typed -> string payload -> 400 Invalid request body; name only -> 201 (CONFIRMED CategoryFormPage.tsx:27,129 vs ProductForm.tsx:100). 

## 3. Required behavior
Category saved with display_order 7

## 4. Reproduction conditions
Headed Chrome: /en/dashboard/categories/new, fill Name and Display order 7, click Save; repeat with Name only

## 5. Verified evidence
- evidence/network/RUN-003/PH-06/AUD-UIFORM-CATEGORY.json
- evidence/screenshots/RUN-003/PH-06/AUD-UIFORM-category-name+display-order.png
- tests/runtime/ph06_ui_category.mjs
- evidence/network/RUN-003/PH-06/AUD-UIFORMS.json

## 6. Root cause and confidence
CONFIRMED — frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx:27,129 keeps display_order as the raw input string while the product forms convert it with toInt (components/products/ProductForm.tsx:100); the backend decodes strictly into an int

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx`:27, `frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx`:129, `frontend/apps/dashboard/src/components/products/ProductForm.tsx`:100
- Frontend: `frontend/apps/restaurant-website/app/[locale]`, `frontend/apps/dashboard/src/app/[locale]/dashboard/products`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: unchanged
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0 (blocks menu setup; products then fall back to no category and break the public menu, see F-022)
- Operational: Categories (dashboard); Products/Storefront dependency workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: convert numeric form fields before submit. In frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx (:27 state, :129 payload) convert display_order with the same toInt helper the product form uses (components/products/ProductForm.tsx:100), move toInt to a shared lib, and show the server field error message instead of a generic toast; add a component test that submits a category with a display order and asserts the request body type.

## 10. Alternative designs
A: shared converter + field errors (recommended). B: make the backend accept numeric strings (loosens the API contract).

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
F-014, F-022 and shared-file findings F-074

## 21. Risks
- Regression in modules sharing the same source files (frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx, frontend/apps/dashboard/src/components/products/ProductForm.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Create category with display order through the real UI; contract test posting the form payload (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Saving a category with display order 7 returns 201 and stores 7.
- A server validation error is shown next to the field.
