# F-062 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Product API performs no input validation (negative/zero price, negative cost or stock, empty name, discount above price, more than 2 decimals) and supplier/customer/category APIs accept duplicates and malformed contacts.

## 2. Current behavior (verified)
all of the above return 201 and are stored (product ids 8-16, suppliers 4-5, customers 4-5, category id 4 duplicate); overflow and over-length return 500 with raw pq text; the SKU unique violation returns 500; selling the negative-price product is blocked later by the order path ("invalid product price"), price 0 is sellable (F-061), 10.12345 is rounded per unit at order time | RUN-003 master-data battery: negative/zero price, negative cost/stock, empty name, discount above price, 5-decimal price, duplicate/malformed supplier/customer accepted (CONFIRMED product_usecase.go:294-306). 

## 3. Required behavior
400 with field errors; sellable products must have price > 0 and 2-decimal money

## 4. Reproduction conditions
POST /products with price -5, price 0, cost -1, cost 50 on price 10, empty name_en, quantity_in_stock -4 (tracked), discount_price 15 on price 10, price 10.12345; POST suppliers/customers with invalid email, duplicates, negative credit limit; duplicate category name

## 5. Verified evidence
- evidence/database/RUN-003/PH-16/AUD-MASTERDATA.json
- evidence/database/RUN-003/PH-16/AUD-MD-NEGPRICE.json
- tests/runtime/ph16_masterdata.mjs
- tests/runtime/ph16_negprice.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/usecase/product_usecase.go:294-306 validateProductRequest has every check commented out and returns nil; supplier/customer handlers validate only required names

## 7. Affected paths
- Backend/database: `backend/internal/usecase/product_usecase.go`:294-306, `backend/internal/handler/http/product_handler.go`:53, `backend/internal/handler/http/supplier_handler.go`:97
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/products`
- API contracts: POST/PUT /api/v1/products, categories, customers, suppliers: 400/409 with field errors
- Accounting logic / posting rules: none named
- Reports affected: Products; Categories; Suppliers; Customers (accounting)

## 8. Impact
- Financial: Negative margin and negative-stock products distort COGS and valuation; free items trigger F-061
- Operational: Products; Categories; Suppliers; Customers (accounting) workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: restore and extend input validation. In backend/internal/usecase/product_usecase.go validateProductRequest (:294-306, all checks commented out) re-enable and extend: name required and <= 100 chars, price > 0 for sellable products and <= a sane maximum with 2 decimals, cost >= 0, stock >= 0, discount <= price, SKU unique per restaurant (return 409, not 500); apply equivalent validation to categories (unique name), customers/suppliers (email/phone format, credit limit >= 0). Add table-driven tests for every rule.

## 10. Alternative designs
A: handler/usecase validation + DB CHECK constraints for the numeric rules (recommended). B: usecase validation only.

## 11. Required policy decisions
Whether a zero price (free item) is allowed and whether cost may exceed price are product decisions.

## 12. Required data migrations
CHECK constraints after listing violating rows.

## 13. Historical data implications
Existing invalid products (negative or zero price, empty name) are listed for correction.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-025, F-031, F-061 and shared-file findings F-040

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/product_handler.go, backend/internal/handler/http/supplier_handler.go, backend/internal/usecase/product_usecase.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Each invalid input returns 400 with the field name; duplicates return 409 (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Each boundary from file 54 (negative price, zero price, cost above price allowed or not per decision, negative stock, empty or over-long name, duplicate SKU) returns the documented 4xx and creates nothing.
- Valid payloads behave as before.
