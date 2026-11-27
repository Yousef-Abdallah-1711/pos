# F-009 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Menu Profitability and Break-even endpoints return HTTP 500.

## 2. Current behavior (verified)
HTTP 500 for both (cost_price column absent) | RUN-003: GET /reports/financial/menu-profitability and /break-even still 500. ROOT CAUSE CONFIRMED: repository/accounting_advanced_repo.go:976 and :1045 select p.cost_price, a column that does not exist on products (the column is cost) -> SQL error mapped to generic 500 (accounting_advanced_handler.go:782). 

## 3. Required behavior
Reports render

## 4. Reproduction conditions
Open advanced insights and Refresh

## 5. Verified evidence
- evidence/reports C-s13c-advanced.json | evidence/network/RUN-003/PH-12/AUD-MET-BATCH.json | evidence/database/RUN-003/PH-12/AUD-F009-ROOT-CAUSE.txt | tests/runtime/ph12_batch_apis.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/repository/accounting_advanced_repo.go:976 and :1045 select p.cost_price, a column that does not exist on products (the column is cost); error mapped to generic 500 at accounting_advanced_handler.go:782

## 7. Affected paths
- Backend/database: `backend/internal/repository/accounting_advanced_repo.go`:976, `backend/internal/repository/accounting_advanced_repo.go`:1045, `backend/internal/handler/http/accounting_advanced_handler.go`:782
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`
- API contracts: unchanged (response shape of both endpoints)
- Accounting logic / posting rules: none named
- Reports affected: Financial reports advanced

## 8. Impact
- Financial: n/a
- Operational: Financial reports advanced workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: change the two SQL statements in backend/internal/repository/accounting_advanced_repo.go at :976 and :1045 from p.cost_price to p.cost (the actual products column, see migrations and report_handler.go:2367 which already uses p.cost); map a missing-column/SQL error to a logged 500 with a generic message instead of leaking it; add a repository integration test that runs both queries against the migrated schema so a column rename fails CI.

## 10. Alternative designs
A: rename in SQL (recommended). B: add a generated column cost_price on products (adds a second cost field - rejected, creates a new competing source).

## 11. Required policy decisions
None for the column fix; cost basis (current vs historical cost) stays under POL-03 (inventory accounting).

## 12. Required data migrations
No.

## 13. Historical data implications
No stored data is affected; reports become available.

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
none and shared-file findings F-016, F-043, F-045, F-047, F-048, F-063, F-068, F-072, F-073

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Endpoint smoke tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- GET /reports/financial/menu-profitability and /break-even return 200 with revenue and cost computed from products.cost for a fixture with 3 products.
- The totals equal an independent SQL calculation over order_items for the same date range.
- A schema-contract test fails if a referenced column does not exist.
