# F-032 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Debug page /en/test-url is shipped in the dashboard app, publicly reachable and not behind the auth provider.

## 2. Current behavior (verified)
Page renders TEST PAGE LOADED with no auth protection and a hydration error

## 3. Required behavior
No test pages in production build

## 4. Reproduction conditions
Open /en/test-url

## 5. Verified evidence
- evidence/crawl-dash-owner.json

## 6. Root cause and confidence
CONFIRMED — frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx exists in the shipped dashboard app

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx`
- Frontend: none
- API contracts: unchanged
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: n/a
- Operational: Dashboard app workflows behave as in section 2.
- Security: Yes — see category: Security/hygiene
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: delete the debug route frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx from the shipped app (or move it behind a development-only guard that is compiled out in production builds), and add a build check that fails when a route named test/debug/dev exists in the production route manifest.

## 10. Alternative designs
A: delete (recommended). B: guard by NODE_ENV (risk of misconfiguration).

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
Keep RBAC module/level checks unchanged; add negative tests for roles without permission.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none

## 21. Risks
- Regression in modules sharing the same source files (frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Route inventory check (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- GET /en/test-url returns 404 in the production build.
- The CI route-manifest check passes.
