# 15 — Financial transaction execution ledger

Authoritative per-transaction table: [44-transaction-level-reconciliation.csv](44-transaction-level-reconciliation.csv) (17 transaction groups; tenants A and C). Oracle comparison: [16-financial-oracle-results.csv](16-financial-oracle-results.csv).

Counts (isolated DB at end): orders 12 (3 paid via UI/API marking; 9 unpaid incl. duplicates); journals 30; GL rows 56; payments_received 6; payments_made 5; supplier invoices 1; customer invoices 1; salaries 1 paid; outbox events 4 (all completed or recovered).

Method labels: **UI** = executed in browser; **API-diagnostic** = direct HTTP (employee fixture, concurrency, boundary, RBAC); **SQL-fixture** = COA/rules seed for tenants A (recovery) and C, closed fiscal period insert. Operational and diagnostic results are kept separate in 44. Capital and rent were entered as manual journals because no capital/expense workflow exists (N4, N8); this is a diagnostic, not a PASS of those capabilities.
