# F-082 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
The Profit and Loss report overflows horizontally at tablet width (768 px): content is 803 px wide and the page scrolls sideways.

## 2. Current behavior (verified)
document.documentElement.scrollWidth is 803 against a 768 px viewport; the same report fits at 1024 px, 390 px and 360 px and in Arabic at 1440 px

## 3. Required behavior
No horizontal page scroll; wide tables scroll inside their own container

## 4. Reproduction conditions
Headed Chrome, viewport 768x1024, English, tenant 4 owner: open /en/dashboard/reports/financial/profit-loss and compare document scrollWidth with the viewport width

## 5. Verified evidence
- evidence/network/RUN-005/PH-04/AUD-R5-LANG-VIEWPORT.json
- evidence/screenshots/RUN-005/PH-04/AUD-R5-LANG-en-768-reports_financial_profit-loss.png
- tests/runtime/ph04_r5_lang_viewport.mjs

## 6. Root cause and confidence
PROBABLE — a fixed-width element or table in the report page without an overflow wrapper at the md breakpoint

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx`
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`
- API contracts: unchanged
- Accounting logic / posting rules: none named
- Reports affected: Reports > Financial > Profit and Loss

## 8. Impact
- Financial: 0
- Operational: Reports > Financial > Profit and Loss workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed (root cause PROBABLE: first isolate the 803 px element with the browser inspector at 768 px): wrap the report table or summary grid in frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx with an overflow-x-auto container, or let the grid collapse to one column at the md breakpoint; add a viewport test at 768, 1024 and 1440 px asserting scrollWidth <= innerWidth.

## 10. Alternative designs
A: overflow wrapper (recommended, smallest). B: responsive grid collapse.

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
- Regression in modules sharing the same source files (frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
At 768, 1024 and 1440 px the page has no horizontal scroll (scrollWidth <= innerWidth) (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- At 768 px the page has no horizontal scroll (document scrollWidth <= 768).
- The table remains fully readable by scrolling inside its container.
- No regression at 360, 390, 1024 and 1440 px, English and Arabic.
