# F-060 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Theme history endpoint returns 501 and the active theme can be deleted.

## 2. Current behavior (verified)
history -> 501 Not Implemented; DELETE active theme -> 204 leaving the tenant without an active theme | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): history 501; active theme deleted with 204 (REPRODUCED). 

## 3. Required behavior
History list; deletion of the active theme refused or falls back to default

## 4. Reproduction conditions
GET /admin/themes/{id}/history; DELETE /admin/themes/{active id}

## 5. Verified evidence
- evidence/database/RUN-003/PH-08/AUD-THEME-DRIVER.json
- evidence/database/RUN-003/PH-08/AUD-THEME2.json
- tests/runtime/ph08_theme_driver.mjs
- tests/runtime/ph08_theme2.mjs

## 6. Root cause and confidence
CONFIRMED — GetThemeHistory returns 501 Not Implemented (admin_theme_handler.go:68)

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/admin_theme_handler.go`:68
- Frontend: none
- API contracts: GET /api/v1/admin/themes/{id}/history (implemented); POST .../restore (new)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0
- Operational: Theme builder workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: implement theme history and protect the active theme. Replace the 501 stub at backend/internal/handler/http/admin_theme_handler.go:68 with a version list (theme_versions table: id, theme_id, config, saved_by, saved_at) and a restore action; in DeleteTheme (admin_theme_handler_v2.go:371) refuse to delete the active theme (409) or require choosing a replacement.

## 10. Alternative designs
A: versions + guard (recommended). B: guard only.

## 11. Required policy decisions
None.

## 12. Required data migrations
Yes: theme_versions table.

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
- Regression in modules sharing the same source files (backend/internal/handler/http/admin_theme_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
History returns versions; deleting the active theme is rejected (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Each save creates a version; restoring version N reproduces its config.
- Deleting the active theme returns 409.
