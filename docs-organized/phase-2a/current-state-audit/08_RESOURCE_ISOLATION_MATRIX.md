# Resource isolation matrix

The table records known evidence at module level. A `tenant_id` column alone is not accepted as proof. Complete per-method query review and two-tenant/two-restaurant mutation tests remain open.

| Resource | Tenant scope | Restaurant scope | Current evidence | Status / missing work |
|---|---|---|---|---|
| Themes | Middleware + handler/service checks; repositories have ID-based methods | Partial | 3 unit and 6 disposable smoke checks per current task records | PARTIAL: prove every read/update/delete and component/theme path. |
| Products | API permission wrapping; repository scope varies | Partial | Some tests and route smoke | OPEN: foreign tenant and same-tenant foreign restaurant read/write/delete. |
| Categories | API permission wrapping | Partial | Employee denied route smoke | OPEN: object-level cross-scope CRUD. |
| Orders | Public order token path plus protected API routes | Required | Public access tests, migration 161 and 15 smoke checks | PARTIAL: POS flow, restaurant scope, direct ID/object paths, actual public create browser flow. |
| Inventory | Middleware/RBAC; legacy inventory paths | Required | No complete negative matrix | OPEN: stock reads/writes and derived records across tenant/restaurant. |
| Drivers/assignments | Permission route coverage | Required | Historic report: driver repo 18/23 methods unscoped; not recounted in current session | OPEN / HIGH: inspect every repository method and parent ownership chain. |
| Users/roles/permissions | Auth and tenant middleware; RBAC endpoints | Tenant-level | Select employee denial + privilege escalation smoke | PARTIAL: per-role CRUD and foreign tenant user/role ID tests. |
| Settings | Auth-only profile/preferences paths | User/tenant | Route inventory, no comprehensive owner tests | OPEN: ensure URL userId is bound to authenticated principal and tenant. |
| Subscriptions/payments | Tenant identity + billing handlers | Tenant | Status/billing smoke and payment unit cases | PARTIAL: ID ownership, callback idempotency/transactions. |
| Accounting/HR | Permission routes and module queries | Tenant + restaurant | Go unit suite has domain/usecase tests; little end-to-end isolation evidence | OPEN: representative query and mutation negatives per table family. |
| Storefront/public menu | Slug-based tenant/restaurant resolver | Restaurant | All page URLs return HTTP 200; no real cross-tenant browser scenario | PARTIAL: slug collision, serialization, rate limits and tenant fixture coverage. |

Acceptance for closure: each listed resource has explicit ownership contract and negative tests using Tenant A/Tenant B plus two restaurants under each tenant; foreign resource IDs cannot read, mutate or delete; tests use isolated fixtures and prove DB is not the live development database. The live database currently has no RLS, so DB-level isolation is NOT VERIFIED/enforced.
