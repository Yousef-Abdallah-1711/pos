# 07 — Security and tenant tests

The report query includes tenant and restaurant predicates for chart accounts and general ledger. A closure follow-up used isolated PostgreSQL 15, actual cookie-session authentication, RBAC enforcement, two owners, one cashier, forged tenant/restaurant identifiers, response content checks, ledger ownership checks and export audit assertions. All recorded scenarios passed. Details are in `closure/01-AUTHENTICATED-TENANT-TEST-MATRIX.md`; raw output is in `closure/evidence/authenticated-tenant-results.json`.

**Result: PASS for tested report/export endpoints; broader product scope not established.** This evidence does not prove all endpoints, background jobs, caches, exports, super-admin flows or RLS behavior are isolated. Do not generalize the endpoint-scoped PASS to the entire product.
