# F-038 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Registration form keeps a stale "choose an available subdomain" error after availability check succeeds.

## 2. Current behavior (verified)
No POST sent; red error persists >=8 s while the same field shows "Subdomain is available"; second click succeeds (clearErrors)

## 3. Required behavior
Error clears (or submit proceeds) once availability is confirmed

## 4. Reproduction conditions
Open /en/auth/register, fill fields, type the subdomain at ~90 ms/char, click Create account ~200 ms after typing (or with check-subdomain delayed 3 s)

## 5. Verified evidence
- evidence/network/RUN-003/PH-05/OBS-ONB-001-probes.json
- evidence/screenshots/RUN-003/PH-05/OBS-ONB-001-S1-click-200ms-after-typing.png
- evidence/screenshots/RUN-003/PH-05/OBS-ONB-001-S4-slow-network-3s-check.png
- tests/runtime/ph05_subdomain_probe.mjs

## 6. Root cause and confidence
CONFIRMED — in the current landing code (frontend/apps/landing/src/app/[locale]/auth/register/page.tsx) onSubmit sets the slug error with setError when availability is not available (lines 59-64) and nothing clears it later: the only useEffect (lines 20-24) handles billing and analytics, availability is read from useDebouncedAvailability (line 39) without an effect that calls clearErrors, and the schema resolver does not know availability, so the manual error persists until the field is edited or the form is submitted again. RUN-003 reproduced the persistence in headed Chrome (2 of 3 early-submit scenarios, error visible >= 8 s while the field showed "available"). The landing source moved under the comm

## 7. Affected paths
- Backend/database: `frontend/apps/landing/src/app/[locale]/auth/register/page.tsx`:34-63, `frontend/apps/landing/src/app/[locale]/auth/register/page.tsx`:59, `frontend/apps/landing/src/app/[locale]/auth/register/page.tsx`:39
- Frontend: `frontend/apps/landing/src/app/[locale]/auth/register/page.tsx`
- API contracts: unchanged
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: none
- Operational: Landing / registration workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed (root cause CONFIRMED in RUN-005 by source: nothing clears the manual slug error): in the landing register page (frontend/apps/landing/src/app/[locale]/auth/register/page.tsx, around lines 59-63) clear the slug error when availability becomes available (useEffect on availability.state) and either disable submit while the check is pending or queue the submit until the check resolves.

## 10. Alternative designs
A: clear on available + disable while checking (recommended). B: re-run validation on submit only.

## 11. Required policy decisions
None (UX). Note: the landing frontend is under separate development; coordinate before editing.

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
- Regression in modules sharing the same source files (frontend/apps/landing/src/app/[locale]/auth/register/page.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Component test: submit during check then check resolves => no stale error (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- After the availability check returns available the error disappears within 200 ms and submit succeeds on the first click.
- Submitting before the check completes does not show a stale error.
