# Functional E2E

## Exercised with real browser and API calls

- Landing login form handled invalid credentials (API 401).
- Registration created two disposable users/tenants (API 201); authenticated session calls succeeded.
- Dashboard loaded authenticated tenant data.
- Category creation returned API 201.
- Product creation returned API 201 and appeared in POS.
- POS added the product to the cart and placed a cash order (API 201). The UI then updated payment status (API 200).
- KDS displayed the created order. Bumping it updated the order from New to In Progress.
- KDS showed `Disconnected` while it loaded orders and allowed a status transition. Live socket/reconnect behavior remains unverified.
- Relevant captured request statuses are summarized in `findings.json`.

## Not covered

Refunds/provider payments, inventory deduction and negative-stock rules, modifier flows, restaurant branch configuration, employee/role workflows, storefront ordering, accounting journals/reconciliation, exports, concurrency, super-admin impersonation, receipt/print output, and cross-tenant CRUD matrices were NOT RUN. This is not a full functional E2E pass.


- Restaurant website host for the disposable tenant resolved to a maintenance page ('We Are Currently Offline / Scheduled Maintenance'), so public menu and ordering acceptance is BLOCKED by the tenant state.
- Super-admin login UI rendered, but its session request to localhost:18080 failed with connection refused after the temporary API forward was stopped; super-admin auth remains NOT RUN.

