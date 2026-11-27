# F-023 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
POS displays $ for EGP amounts.

## 2. Current behavior (verified)
Prices shown as $10000.00 | RUN-003 (AUD-EVD-F023): POS shows '$10,000.00' for an EGP product; no EGP label on page (STILL REPRODUCIBLE). 

## 3. Required behavior
EGP symbol

## 4. Reproduction conditions
Open POS

## 5. Verified evidence
- screenshots pos-A-cart | evidence/network/RUN-003/PH-20/AUD-EVD-gaps.json
- tests/runtime/ph20_evidence_gaps.mjs
- evidence/screenshots/RUN-003/PH-20/AUD-EVD-F023-pos.png

## 6. Root cause and confidence
CONFIRMED — frontend/apps/dashboard/src/lib/rtl.ts:81-91 formatCurrencyWithLocale defaults currency to USD (Intl.NumberFormat style currency), so any caller that omits the tenant currency prints $; the storefront passes EGP and shows "EGP" (AUD-STORE-004) while the POS shows $

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/lib/rtl.ts`:81-91
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`
- API contracts: GET tenant settings gains currency (additive) if not present
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: n/a
- Operational: POS workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: stop defaulting to USD. In frontend/apps/dashboard/src/lib/rtl.ts:81-91 (formatCurrencyWithLocale) require an explicit currency argument (no default) and make every POS/cart/receipt caller pass the tenant currency from one source (a tenant settings hook, today missing: no currency setting exists); until a currency setting exists use the single functional currency decided in POL-11 as a named constant, not an implicit default.

## 10. Alternative designs
A: required argument + tenant currency source (recommended). B: change the default to EGP (hides the same defect for other currencies). C: remove the symbol (loses information).

## 11. Required policy decisions
POL-11 (currency policy) in 56.

## 12. Required data migrations
Possibly a tenants.currency column (default EGP).

## 13. Historical data implications
None (display only).

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
- Regression in modules sharing the same source files (frontend/apps/dashboard/src/lib/rtl.ts).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Currency display test (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- The POS, receipts and orders list show the same currency code/symbol as the storefront and reports for the same tenant.
- A lint/test fails when formatCurrencyWithLocale is called without a currency.
- Switching the tenant currency setting changes all displays.
