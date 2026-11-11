# PostgreSQL RLS and tenant security matrix

## Evidence carried forward

The read-only catalog result from the prior disposable database run found 92 public ordinary tables, with RLS enabled on 15 (16.3%). It specifically reported RLS disabled on these core tenant-bearing tables:

| Table | Preliminary ownership | Tenant key | RLS evidence | Risk / next step |
|---|---|---|---|---|
| users | Tenant/user | `tenant_id` | Disabled in prior catalog | Test all user-scoped repository and handler paths under `pos_app`; design a scoped policy. |
| tenants | Tenant root | Primary tenant identity | Disabled in prior catalog | Keep cross-tenant access limited to explicit platform administration. |
| restaurants | Tenant-owned | `tenant_id` | Disabled in prior catalog | Verify restaurant/tenant relationship constraints and runtime scoping. |
| orders | Tenant/restaurant-owned | `tenant_id`, `restaurant_id` | Disabled in prior catalog | Test reads, mutations, and nested items with cross-tenant IDs. |
| order_items | Parent/order-owned | Parent order relationship | Disabled in prior catalog | Enforce ownership through the order relationship and test item IDs across tenants. |
| products | Tenant/restaurant-owned | `tenant_id`, restaurant relationship | Disabled in prior catalog | Test parent-child category links and all CRUD operations. |
| categories | Tenant/restaurant-owned | `tenant_id`, restaurant relationship | Disabled in prior catalog | Test product/category cross-tenant references and mutations. |
| inventory | Tenant/restaurant-owned | `tenant_id`, restaurant relationship | Disabled in prior catalog | Test read/write isolation and order side effects. |
| employees | Tenant/restaurant-owned | `tenant_id`, restaurant relationship | Disabled in prior catalog | Test role, salary, and employee references across tenants. |

The source migrations contain 55 explicit `CREATE TABLE` statements, while the prior live catalog observed 92 public tables. Source-only parsing cannot reliably enumerate all deployed tables or their final columns, current policies, grants, views, dynamic SQL objects, or migration state.

## Status

**INCOMPLETE / BLOCKED.** Docker Desktop and PostgreSQL services are unavailable in this environment. Therefore no current catalog query, table-by-table ownership classification, restricted-role negative test, tenant/RLS integration test, or forward/rollback migration acceptance was run. The prior API evidence included one cross-user settings denial (403), own-user settings access (200), and unauthenticated denials (401); it is not a complete tenant matrix.

No RLS policies or database migrations were added. Do not enable RLS broadly until the full deployed schema, role behavior, policies, transaction-local tenant context, application query coverage, and rollback plan are verified on disposable PostgreSQL.
