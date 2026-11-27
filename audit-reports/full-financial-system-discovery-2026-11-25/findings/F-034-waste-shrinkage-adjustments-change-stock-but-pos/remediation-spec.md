# F-034 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Waste/shrinkage adjustments change stock but post nothing and the turnover report shows waste 0.

## 2. Current behavior (verified)
Stock 48 to 47; no journal; waste_qty 0; reason count_shortage not accepted (allowed: purchase, sale, adjustment, waste, return)

## 3. Required behavior
Expense posted and waste reported

## 4. Reproduction conditions
POST inventory adjust reason waste -1

## 5. Verified evidence
- evidence/runtime/39

## 6. Root cause and confidence
CONFIRMED — backend/internal/handler/http/inventory_adjustment_handler.go (Adjust :29-105) updates stock and writes an inventory movement but contains no journal/posting call (0 references); waste therefore posts nothing, and turnover reads the legacy table (F-046)

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/inventory_adjustment_handler.go`:29-105
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/dashboard/src/app/[locale]/dashboard/products`
- API contracts: POST /api/v1/inventory/adjustments: response includes journal_entry_id (additive)
- Accounting logic / posting rules: none named
- Reports affected: Inventory; Reports

## 8. Impact
- Financial: n/a
- Operational: Inventory; Reports workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: post inventory adjustments. In backend/internal/handler/http/inventory_adjustment_handler.go Adjust (:29-105) call the auto-posting engine inside the same transaction: waste/shrinkage Dr waste expense (or COGS adjustment) / Cr inventory 1130 at the product cost; count gains the reverse; include the amounts in the turnover report waste column (today 0) by reading the movement type.

## 10. Alternative designs
A: post at adjustment time (recommended). B: post at period close from a stock valuation difference (delays recognition). C: leave unposted (current).

## 11. Required policy decisions
POL-03 (inventory accounting) in 56.

## 12. Required data migrations
No schema change; correction journals for historical adjustments behind approval.

## 13. Historical data implications
Past adjustments are listed with cost impact; accountant approves a catch-up journal.

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
F-046

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/inventory_adjustment_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Adjustment/waste scenario tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Adjusting 3 units of a product costing 40 posts Dr waste 120 / Cr 1130 120 and stock and ledger agree.
- The turnover report shows waste 120 for the period.
- Replaying the same adjustment request (idempotency key) posts once.
