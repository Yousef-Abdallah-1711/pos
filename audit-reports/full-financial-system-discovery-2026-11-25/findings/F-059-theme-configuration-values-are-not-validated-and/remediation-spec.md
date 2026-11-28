# F-059 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Theme configuration values are not validated and are served publicly (invalid colors and url(javascript:) accepted).

## 2. Current behavior (verified)
201, activate 200, public endpoint returns the raw values; theme name with HTML also stored; exploitability in the storefront CSS pipeline not tested (PROBABLE CSS injection) | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): invalid colors and url(javascript:) stored and served by public theme endpoint (REPRODUCED). 

## 3. Required behavior
Colors validated (hex/rgb/hsl allow-list); unsafe values rejected

## 4. Reproduction conditions
POST /admin/themes with colors.primary "not-a-color" and colors.text "url(javascript:alert(1))", activate, GET /public/restaurants/{slug}/theme

## 5. Verified evidence
- evidence/database/RUN-003/PH-08/AUD-THEME-DRIVER.json
- evidence/database/RUN-003/PH-08/AUD-THEME2.json
- tests/runtime/ph08_theme_driver.mjs
- tests/runtime/ph08_theme2.mjs

## 6. Root cause and confidence
CONFIRMED — the theme handler/service persist config JSON without colour or value validation (no validate call in the theme service; url(javascript:) was accepted and served)

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/admin_theme_handler_v2.go`:293
- Frontend: `frontend/apps/restaurant-website/app/[locale]`
- API contracts: POST/PUT /api/v1/admin/themes: 400 with field errors
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0 (integrity/security of public pages)
- Operational: Theme builder; Storefront; public theme API workflows behave as in section 2.
- Security: Yes — see category: Security / validation
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: schema-validate theme JSON. In the theme save path (backend/internal/handler/http/admin_theme_handler_v2.go:293 and the theme service) validate every colour against a hex/rgb(a)/hsl pattern, fonts against an allow-list, URLs against http/https and the upload domain only (reject javascript:, data: and url() with other schemes), sizes within ranges, and reject unknown keys; sanitise again when rendering on the storefront.

## 10. Alternative designs
A: strict schema on write + defensive render (recommended). B: render-side escaping only (bad data remains).

## 11. Required policy decisions
None.

## 12. Required data migrations
No; existing themes are scanned and reported.

## 13. Historical data implications
Existing themes with invalid values (evidence) are listed and fixed by the tenant.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
Keep RBAC module/level checks unchanged; add negative tests for roles without permission.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none and shared-file findings F-025

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/admin_theme_handler_v2.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Invalid/unsafe color or CSS rejected 400; public theme contains only validated values (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A colour "not-a-color" returns 400.
- A value containing url(javascript:...) returns 400 and never reaches the public theme endpoint.
- Valid themes round-trip unchanged.
