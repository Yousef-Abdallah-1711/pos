# F-046 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Inventory turnover report uses current stock as opening stock and reads waste from the legacy inventory table, ignoring inventory movements.

## 2. Current behavior (verified)
opening_stock_qty == current_stock_qty (33 and 21) so average inventory equals current stock; waste_qty 0 although inventory_movements has waste 1 for products 2 and 4; Audit Burger original stock 50 | RUN-003: turnover opening_stock_qty == current stock (33/21); waste_qty 0 vs movements waste 1 (REPRODUCED). 

## 3. Required behavior
Opening quantity from stock history; waste_qty reflects recorded waste

## 4. Reproduction conditions
GET /reports/inventory after recording a waste adjustment and sales (tenant 4)

## 5. Verified evidence
- evidence/network/RUN-003/PH-12/AUD-MET-BATCH2.json
- evidence/database/RUN-003/PH-12/AUD-MET-BATCH2-compare.json
- evidence/database/RUN-003/PH-06/OBS-INV-001-controlled.json
- tests/runtime/ph12_batch_apis2.mjs

## 6. Root cause and confidence
CONFIRMED — report_handler.go:557-645 / 2480-2609 reads legacy inventory.quantity_after else current stock and SUM(ABS(change)) from the legacy inventory table, while the live ledger is inventory_movements

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/report_handler.go`:557-645, `backend/internal/handler/http/report_handler.go`:2480-2609
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/dashboard/src/app/[locale]/dashboard/products`
- API contracts: response shape unchanged
- Accounting logic / posting rules: none named
- Reports affected: Inventory turnover report; Overview report

## 8. Impact
- Financial: Turnover ratio misstated (1.42 vs 1.01 using true average 41.5 qty)
- Operational: Inventory turnover report; Overview report workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: compute turnover from the inventory movement ledger. In backend/internal/handler/http/report_handler.go (:557-645 and :2480-2609) take the opening quantity from inventory_movements at the period start (not current stock or the legacy inventory table) and waste from movements of type waste/adjustment_out; average inventory = (opening value + closing value)/2 at product cost; document the cost basis.

## 10. Alternative designs
A: movement-ledger based (recommended). B: keep legacy table and label (stale source).

## 11. Required policy decisions
POL-03 (inventory accounting) in 56.

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
none and shared-file findings F-012, F-067

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/report_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Turnover report with opening stock, sales and waste equals independent calculation (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- For a fixture with known opening stock, purchases, sales and waste the turnover equals an independent calculation to 2 decimals.
- Waste in the report equals the sum of waste movements (see F-034).
