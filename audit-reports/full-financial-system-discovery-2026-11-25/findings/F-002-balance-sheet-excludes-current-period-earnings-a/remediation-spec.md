# F-002 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Balance Sheet excludes current-period earnings and is UNRECONCILED by exactly the net result.

## 2. Current behavior (verified)
equation_difference -2000.00 (net loss); equity shows 100000 only | RUN-003 clean-tenant month (AUD-MONTH): balance sheet equation_balanced=false, equation_difference -27,710.00 = net loss exactly (REPRODUCED). 

## 3. Required behavior
A = L + E

## 4. Reproduction conditions
S13 loss scenario then open Balance Sheet

## 5. Verified evidence
- evidence/reports C-s13a-reports_financial_balance_sheet.json | evidence/database/RUN-003/PH-14/AUD-MONTH.json | tests/runtime/ph14_month.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/usecase/financial_report_usecase.go:285-306 balance-sheet query selects only account_type IN (asset, liability, equity); revenue/expense accounts never contribute, so current-period earnings are not added to equity (BS difference 82,700 measured earlier)

## 7. Affected paths
- Backend/database: `backend/internal/usecase/financial_report_usecase.go`:285-306
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`
- API contracts: unchanged unless the design below says otherwise (see 08-api-inventory.csv)
- Accounting logic / posting rules: none named
- Reports affected: Financial reports

## 8. Impact
- Financial: 2000 per period (signed residual -2000)
- Operational: Financial reports workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: in the balance-sheet query (backend/internal/usecase/financial_report_usecase.go:285-306) add a computed equity line 'Current period earnings' = SUM(revenue) - SUM(expense) up to the as-of date (and prior-year retained earnings when a closed fiscal year exists), include it in total_equity and in the equation check; keep account-level rows unchanged.

## 10. Alternative designs
Alternative A: fix at the source layer named in root-cause.md; Alternative B: add a compensating check in the calling layer (weaker, not recommended unless the source layer is vendor/owner-controlled).

## 11. Required policy decisions
Retained-earnings and year-end closing treatment (POL-05).

## 12. Required data migrations
No schema change expected; confirm during implementation.

## 13. Historical data implications
Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.

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
none and shared-file findings F-003, F-011, F-012

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/usecase/financial_report_usecase.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
BS residual = 0 under profit and loss (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Clean-tenant month: equation_balanced = true and equation_difference = 0.00 (currently -27,710.00).
- Trial balance and balance sheet agree for every as-of date in the dataset.
