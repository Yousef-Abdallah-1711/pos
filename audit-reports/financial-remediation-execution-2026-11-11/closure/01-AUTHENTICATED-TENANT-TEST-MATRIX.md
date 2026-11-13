# Authenticated Tenant Test Matrix

## Environment and evidence

Run `291a83d6`, 2026-11-11 16:10 UTC. Real HTTP API, actual cookie-session login, RBAC enforcement enabled, synthetic users/restaurants/ledger fixtures, isolated disposable PostgreSQL 15 tmpfs database. No credentials or tokens are retained. Raw result: [`authenticated-tenant-results.json`](evidence/authenticated-tenant-results.json); reproducible runner: [`run-authenticated-tenant-tests.ps1`](evidence/run-authenticated-tenant-tests.ps1).

| Scenario | Authenticated role / requested scope | Expected | Actual | Result |
|---|---|---|---|---|
| Owner A actual login | Owner A / tenant 14, restaurant 13 | Session established | HTTP 200 | PASS |
| Owner B actual login | Owner B / tenant 15, restaurant 14 | Session established | HTTP 200 | PASS |
| Owner A own report | Owner A / tenant 14 | Assets and equity 110; no B marker | HTTP 200; 110; tenant 14 | PASS |
| Owner B own report | Owner B / tenant 15 | Assets and equity 220; no A marker | HTTP 200; 220; tenant 15 | PASS |
| Owner A forged tenant and restaurant | Owner A session; attempted tenant 15 / restaurant 14 | Caller scope only | HTTP 200; assets 110, tenant 14 | PASS |
| Owner B forged tenant and restaurant | Owner B session; attempted tenant 14 / restaurant 13 | Caller scope only | HTTP 200; assets 220, tenant 15 | PASS |
| Missing authentication | No session | Deny | HTTP 401 | PASS |
| Invalid authentication | Invalid session | Deny | HTTP 401 | PASS |
| Cashier report | Cashier / tenant 16 | Deny | HTTP 403 | PASS |
| Owner A export with Tenant B headers | Owner A session; attempted tenant 15 | CSV contains A 123.45, not B 987.65 | HTTP 200; A present, B absent; tenant 14 | PASS |
| Cashier export | Cashier session | Deny | HTTP 403 | PASS |
| Fixture and audit database assertions | A/B journal and ledger ownership | Balanced scoped rows; export audit under A | Balanced journals; two ledger rows per tenant; A export audit under tenant 14 | PASS |

The successful `200` forged-header cases are accepted only because response contents and effective tenant scope were independently checked. These results establish the tested endpoint behavior, not universal tenant isolation across the product.
