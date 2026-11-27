# F-062 — Product API performs no input validation (negative/zero price, negative cost or stock, empty name, discount above price, more than 2 decimals) and supplier/customer/category APIs accept duplicates and malformed contacts

- **Severity:** Medium  ·  **Category:** Validation gap  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Products; Categories; Suppliers; Customers (accounting)
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-MD-001, EX-RUN004-FORMMATRIX, EX-RUN004-FORMMATRIX2
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 400 with field errors; sellable products must have price > 0 and 2-decimal money

Observed: all of the above return 201 and are stored (product ids 8-16, suppliers 4-5, customers 4-5, category id 4 duplicate); overflow and over-length return 500 with raw pq text; the SKU unique violation returns 500; selling the negative-price product is blocked later by the order path ("invalid product price"), price 0 is sellable (F-061), 10.12345 is rounded per unit at order time | RUN-003 master-data battery: negative/zero price, negative cost/stock, empty name, discount above price, 5-decimal price, duplicate/malformed supplier/customer accepted (CONFIRMED product_usecase.go:294-306). 

## Severity rationale
Rated **Medium** (category: Validation gap; evidence: REPRODUCED). Basis recorded in the register: Negative margin and negative-stock products distort COGS and valuation; free items trigger F-061. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Products; Categories; Suppliers; Customers (accounting). Negative margin and negative-stock products distort COGS and valuation; free items trigger F-061

## Technical impact
CONFIRMED (source): backend/internal/usecase/product_usecase.go:294-306 validateProductRequest has every check commented out and returns nil; supplier/customer handlers validate only required names

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
Negative margin and negative-stock products distort COGS and valuation; free items trigger F-061

## Where it was reproduced
POST /products with price -5, price 0, cost -1, cost 50 on price 10, empty name_en, quantity_in_stock -4 (tracked), discount_price 15 on price 10, price 10.12345; POST suppliers/customers with invalid email, duplicates, negative credit limit; duplicate category name

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
