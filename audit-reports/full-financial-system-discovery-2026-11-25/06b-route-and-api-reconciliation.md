# 06b - Route and API reconciliation (RUN-005, Phase L / T-P03-002)

Built by `tests/build_route_api_reconciliation.py` from the current source tree (read-only). Detail: [06b-route-and-api-reconciliation.csv](06b-route-and-api-reconciliation.csv).

## Pages and route handlers
- Page files: **117** (dashboard 79, landing 17, restaurant-website 12, super-admin 9); route handlers: **2** (super-admin /favicon.ico, restaurant-website /api/metrics). Pages + handlers = 119, which equals the 119 discovered entries in 06: the 119-vs-117 difference is exactly the 2 route handlers (T-P03-002 resolved; the earlier "117 loaded" counted pages only).
- Dynamic page patterns: 15 (need a valid id to render; invalid-id behaviour is tested separately).
- Dashboard pages not in the 53 control inventory (static): 4: /[locale]/dashboard/categories/new, /[locale]/dashboard/drivers/new, /[locale]/dashboard/settings/roles/new, /[locale]/dashboard/theme-builder/import.
- Dashboard pages without interaction results in 53b (static): 4.

## API registrations
- Distinct `METHOD /api/v1/...` literals registered in backend sources: **354**; rows in 08: **355** (after normalising path parameters).
- In source but not in 08 (0): none
- In 08 but not found in source literals (1): GET /uploads/{}
- Note: authentication-only sweeps (EX-API-*) are NOT functional API coverage; functional coverage is the set of endpoints exercised by the RUN-003/004 diagnostics and the UI matrices.
