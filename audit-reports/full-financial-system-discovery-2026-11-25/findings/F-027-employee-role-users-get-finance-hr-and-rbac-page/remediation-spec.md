# F-027 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Employee-role users get finance, HR and RBAC page shells in the UI (no UI guard) while the API returns 403.

## 2. Current behavior (verified)
Shells render (P&L, BS, CF, VAT, aging, payroll, RBAC roles) with blank data; only 24 pages gated

## 3. Required behavior
Page guard consistent with API permissions

## 4. Reproduction conditions
Crawl dashboard as Employee user: 50 of 79 routes render; 24 show Access denied; finance report pages render with 403 data calls

## 5. Verified evidence
- evidence/crawl-dash-cashier.json

## 6. Root cause and confidence
CONFIRMED — only 6 dashboard page files reference a permission hook; e.g. settings/roles/page.tsx has no guard reference, so the shell renders for any signed-in role while the API returns 403; other pages guard inconsistently

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx`
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: unchanged
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: n/a
- Operational: Frontend routing; RBAC workflows behave as in section 2.
- Security: Yes — see category: Security/usability
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: one route-guard mechanism. (1) Add a shared guard wrapper for dashboard pages that checks the same module and level the API enforces (the permission hook already exists and is used by 6 page files); (2) apply it to every page under accounting, hr, settings/roles, reports and billing, driven by a single route-to-permission table kept next to the navigation definition so the menu and the guard cannot diverge; (3) unauthorised users get an access-denied page, never the page shell.

## 10. Alternative designs
A: central route table + guard component (recommended). B: per-page checks (current state, drifts). C: rely on the API 403 only (shell leaks structure).

## 11. Required policy decisions
POL-15 (default role matrix) defines what each role may see.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
UI guards are advisory; API enforcement stays authoritative (see F-020/F-071).

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none

## 21. Risks
- Regression in modules sharing the same source files (frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx, frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Page-permission matrix tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- For the Employee role every finance/HR/RBAC route renders the access-denied page and issues no data requests.
- The menu hides exactly the routes the guard would block.
- A test enumerates all dashboard routes and fails when a route is missing from the permission table.
