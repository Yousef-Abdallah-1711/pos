# F-076 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Super Admin Users list links every row to /dashboard/users/{id}, a route that does not exist (404).

## 2. Current behavior (verified)
Navigation to /dashboard/users/{id} returns the 404 page (7 consecutive rows in the control matrix logged a 404 resource error); only dashboard/tenants/[id] exists as a detail route

## 3. Required behavior
A user detail page (or no link)

## 4. Reproduction conditions
Headed Chrome, Super Admin app (port 13114), Users page: click View on a user row

## 5. Verified evidence
- evidence/network/RUN-004/PH-13/AUD-R4-SUPERADMIN-CONTROL-MATRIX.json
- tests/runtime/ph13_r4_superadmin_ui.mjs

## 6. Root cause and confidence
CONFIRMED — super-admin/src/app/dashboard/users/page.tsx:132 renders href /dashboard/users/${user.id}; the directory app/dashboard/users contains only page.tsx (no [id] route)

## 7. Affected paths
- Backend/database: `frontend/apps/super-admin/src/app/dashboard/users/page.tsx`:132, `frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx`:132, `frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx`
- Frontend: `frontend/apps/super-admin`
- API contracts: unchanged
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0
- Operational: Super Admin > Users workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: either add frontend/apps/super-admin/src/app/dashboard/users/[id]/page.tsx (user profile, memberships, status, last login, using the existing super-admin users API) or remove the View link at frontend/apps/super-admin/src/app/dashboard/users/page.tsx:132 until the page exists; add a link-check test that visits every internal href rendered on the page.

## 10. Alternative designs
A: build the detail page (recommended if support staff need it). B: remove the link.

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
- Regression in modules sharing the same source files (frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx, frontend/apps/super-admin/src/app/dashboard/users/page.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Every link rendered on the Users page resolves with 200 (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Every link on the Users page returns 200.
- The detail page shows the same email, status and role data as the list row.
