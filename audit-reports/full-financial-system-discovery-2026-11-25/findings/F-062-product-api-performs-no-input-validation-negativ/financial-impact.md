# F-062 — Financial impact

- **Quantified impact (register):** Negative margin and negative-stock products distort COGS and valuation; free items trigger F-061
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, concurrency, ui, validation

## What is wrong in the books or reports
all of the above return 201 and are stored (product ids 8-16, suppliers 4-5, customers 4-5, category id 4 duplicate); overflow and over-length return 500 with raw pq text; the SKU unique violation returns 500; selling the negative-price product is blocked later by the order path ("invalid product price"), price 0 is sellable (F-061), 10.12345 is rounded per unit at order time | RUN-003 master-data battery: negative/zero price, negative cost/stock, empty name, discount above price, 5-decimal price, duplicate/malformed supplier/customer accepted (CONFIRMED product_usecase.go:294-306). 

## Which reports and decisions are affected
Financial reports that read the affected data: Products; Categories; Suppliers; Customers (accounting)

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

