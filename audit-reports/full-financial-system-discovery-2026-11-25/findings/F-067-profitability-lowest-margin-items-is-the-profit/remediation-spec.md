# F-067 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Profitability "Lowest Margin Items" is the profit-amount ranking read backwards, so a 100%-margin item ranks second-lowest.

## 2. Current behavior (verified)
Order is MD probe 0%, MD probe 100%, Probe A 60%, Probe B 60%, Burger 60%: the list is the gross-PROFIT descending list iterated backwards

## 3. Required behavior
Items ordered by gross_margin_pct ascending (0%, 60%, 60%, 60%, 100%)

## 4. Reproduction conditions
Headed Chrome /en/dashboard/reports/profitability, tenant 4, 2026-10-29..2026-11-27; compare lowest_margin_items with margin percentages from SQL

## 5. Verified evidence
- evidence/network/RUN-004/PH-12/AUD-R4-MET-CAPTURE.json
- evidence/network/RUN-004/PH-12/AUD-R4-MET2-CAPTURE.json
- evidence/database/RUN-004/PH-12/AUD-R4-ORACLE-PROF-CF.json
- evidence/screenshots/RUN-004/PH-12/AUD-R4-MET2-prof.png

## 6. Root cause and confidence
CONFIRMED — report_handler.go:2432 ORDER BY gross_profit DESC, then :2470-2476 builds the lowest list by iterating the same slice backwards; no margin sort exists

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/report_handler.go`:2432, `backend/internal/handler/http/report_handler.go`:2470, `backend/internal/handler/http/report_handler.go`:212
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`
- API contracts: response shape unchanged; ordering of lowest_margin_items changes
- Accounting logic / posting rules: none named
- Reports affected: Reports > Profitability (API, UI, CSV/PDF)

## 8. Impact
- Financial: 0 direct (decision quality: menu pricing review points at the wrong items)
- Operational: Reports > Profitability (API, UI, CSV/PDF) workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: sort the lowest-margin list by margin. In backend/internal/handler/http/report_handler.go replace the backwards iteration at :2470-2476 with a query (or an explicit sort) by gross_margin_pct ASC then revenue DESC, keeping the top list ordered by gross profit (:2432); label both lists with their sort key in the UI and exports.

## 10. Alternative designs
A: separate sort for the lowest list (recommended). B: rename the section "Lowest profit" (changes the business meaning).

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
F-007, F-040 and shared-file findings F-012, F-046

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/report_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Fixture with items of different margin and profit; assert both lists (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- For items with margins 0%, 60%, 60%, 60%, 100% the lowest list is ordered 0%, 60%, 60%, 60%, 100%.
- The top list stays ordered by gross profit.
- CSV and PDF exports match the screen.
