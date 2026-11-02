# 06 - Multi-Tenancy and Sharding (Phase 1, read-only)

Auditor scope: tenant isolation, sharding/scaling. Static analysis only; DB container was not running and no request was sent to a live API, so every "exploitable" statement below is derived from code reading and should be confirmed by Phase 2 dynamic tests.

## 1. Isolation model (what actually exists)

- Hierarchy: tenant -> restaurant -> data, every operational table has (tenant_id, restaurant_id) denormalised columns (see 05).
- **Enforcement is application-layer only.** No RLS in migrations (DB-005); app DB user is the `postgres` superuser.
- Tenant source: `middleware/tenant.go:19-69` `TenantContextMiddleware` - JWT claims first, **headers as fallback** (DB-001). `AuthMiddleware` runs first (main.go:604-605). Some RBAC routes read claims directly (`wrapAuthOnly`, main.go:597-599) which is safer.
- Public routes: menu/homepage by slug (no tenant header), but **public order routes use the header fallback unauthenticated** (DB-002).
- WebSocket KDS: tenant from JWT claims, room per tenant+restaurant (kds_hub.go:171) - OK.
- Caches: financial reports key includes tenant and restaurant (financial_report_handler.go:72 `financialReportCacheKey(tenantID, restaurantID, path, query)`); tenant_access cache keyed by tenant (tenant_status.go:57); theme caches `theme:%d` keyed by id only (theme_service.go:200) and `theme:list:%d:...` by tenant. Legacy services keyed by restaurantID only (service/legacy, not imported by cmd/api - grep).
- Exports: orders export (admin_order_handler.go:610-613), report export, audit export (accounting_advanced_handler.go:943-952) all take tenant/restaurant from middleware context and pass to repos with tenant+restaurant - SCOPED (subject to DB-001 header override).
- Uploads: path includes tenant id (product_usecase.go:299); serving route not found (DB-021).

## 2. Repository scoping verdicts

Method: script (`evidence/database-evidence/repository-tenant-scan.txt`) flags functions containing SQL whose body never mentions tenant_id/tenantID or restaurant_id/restaurantID; flagged ones were then read manually. "SCOPED" = every query read filters tenant and (where the table has it) restaurant; "PARTIAL" = some methods lack a predicate or rely on caller pre-check; "UNSCOPED" = tenant-owned data reachable by id only with no check anywhere on the path. Platform-level tables (plans, modules, super admin) are marked N/A.

| Repository (backend/internal/repository) | Verdict | Notes |
|---|---|---|
| account_repo, accounting_advanced_repo, attendance_repo, auto_posting_repo, bank_account_repo, category_repo, contract_repo, customer_acct_repo, employee_repo, general_ledger_repo, insurance_repo, inventory_movement_repo, kds_repo, leave_repo, payment_acct_repo, performance_review_repo, rbac_role_repository, resignation_repo, role_repo, salary_repo, supplier_repo, table_repo, tax_repo, training_repo | SCOPED | scan: 0 functions without tenant; restaurant also present |
| order_repo | PARTIAL | GetOrderByID/ByNumber/List/UpdateStatus/Cancel/Delete scoped (order_repo.go:136,238,268,479,499,574); **GetOrderItems :700 and GetOrderStatusHistory :762 query by order_id only** (DB-017) |
| product_repo | PARTIAL | UpdateProduct `WHERE id AND tenant_id` (:372) but no restaurant (:340); GetProductImages by product_id only (:606) |
| invoice_customer_repo / invoice_supplier_repo / journal_entry_repo | SCOPED (headers) | child-line helpers (getInvoiceLines, getLines, getPostingLinesTx) by parent id after parent scoped query |
| restaurant_repo | PARTIAL | GetBySlug global (:32) - needed for public routes but slug uniqueness not DB-enforced globally (DB-018) |
| user_repo | PARTIAL | DeleteUser has tenant (:382); UpdatePassword/UpdateProfile/VerifyPassword by user id (self-service; callers not individually traced) |
| user_role_repository | SCOPED (tenant) | all methods take tenantID; no restaurant predicate (acceptable: tenant-level roles) |
| role_permission_repository | PARTIAL | **Revoke / RevokeAll / GetByRoleAndModule by role id only** (:139-166) - DB-004 |
| driver_repo | PARTIAL (18 of 23 methods unscoped) | handlers pre-check ownership (admin_driver_handler.go:263,292,359,...); repo itself UNSCOPED for Update/Delete/UpdateStatus/assignments/stats (DB-017) |
| notification_repo | SCOPED (user) | MarkAsRead/Unread/Delete filter by user_id (:260,285,326); user belongs to tenant |
| homepage_section_repository, website_settings_repository | PARTIAL | restaurant-scoped, no tenant predicate; restaurant id comes from context (subject to DB-001) |
| restaurant_theme_repo | PARTIAL | UpdateTheme tenant-only; section ops by theme_id only |
| theme_repository (themes_v2) | **UNSCOPED** for Update/Delete/GetByID path (:300,:360) - DB-003 |
| payment_repo_impl, subscription_repo_impl | PARTIAL | tenant-level billing; UpdatePaymentOrder/Status by payment id; invoked from payment flow/webhooks (not traced end-to-end); no restaurant col on payments |
| user_settings_repo | PARTIAL | by user id only (self-service) |
| auth_repo | N/A/PARTIAL | login/registration paths are pre-tenant by nature; CheckEmailExists global |
| plan_repo, module_definition_repository, super_admin_repo, super_admin_audit_repo, impersonation_repo, permission_audit_log_repository | N/A (platform) | super-admin routes; authorization correctness is out of my scope (auth auditor) |
| blockchain_repository, payment_repo, subscription_repo | no SQL / interfaces | |
| repository/legacy/* | not wired into cmd/api (grep) | not assessed |

SQL in handlers/usecases/services: direct db.Query/Exec exists in handler/http/{hr_report,report,accounting_advanced,hr_attendance,hr_leave,hr_salary,payment_monitoring}_handler.go, usecase/{financial_report,payment}_usecase.go and service/{webhook_retry,subscription_expiry}_worker.go. Spot-checked: attendance/leave/salary deletes filter id+tenant_id+restaurant_id (hr_attendance_handler.go:331, hr_leave_handler.go:336, hr_salary_handler.go:371); customers_accounting lookup scoped (accounting_advanced_handler.go:1310-1312); hr_report builds a dynamic `where` (hr_report_handler.go:77-146, 23 tenant_id mentions - clause content not fully traced); report_handler.go:2459,2480 are inner joins of already-scoped queries. payment_usecase.go UPDATE payments/tenants by id (:751,791,823,850) and payment_monitoring_handler / webhook_retry_worker were not traced in detail (unverified). Unused package internal/handlers was not assessed.

**IDOR risk list (ordered):** DB-003 themes (any authed user, cross-tenant write), DB-002 public orders (unauth read/cancel), DB-001 header override, DB-004 RBAC revoke, DB-017 drivers/order items/payments (guarded by handler pre-checks today).

## 3. Sharding and scaling

| Question | Answer | Evidence |
|---|---|---|
| Implemented or experimental? | **Experimental / scaffolding.** Not imported by cmd/api. | DB-008 |
| Routing key | restaurant_id (string) via MD5 + a broken jump hash that is key-independent | DB-006 |
| Tenant->shard mapping | none (no directory/table); all restaurants of a tenant may land on different shards | DB-009 |
| Cross-shard transactions | none; no 2PC/saga | DB-009 |
| Main API uses shard connections? | No: repositories receive a single `*sql.DB` (main.go:125 area) | DB-008 |
| Migration tooling | `cmd/migrate-shard`, `internal/sharding/migration.go`: wrong tables, no scoping, destructive rollback | DB-007 |
| Scaling (autoscaler, LB, failover, multi-region) | in-memory simulation; cloud calls are comments; scaling_api routes use unsupported `:id` syntax and no auth | DB-008 |
| Docs vs implementation | IMPLEMENTATION_SUMMARY.md says COMPLETED/production-ready; DATABASE_SHARDING.md says In Progress; tests (tests/unit, tests/integration) require live shards and were not run | DB-008 |

Conclusion: the platform is a single-Postgres monolith. Claims of "1000+ restaurants across 4 shards" are unsupported by running code.

## 4. Findings (multi-tenancy and sharding)

### DB-001 [P0] Request headers X-Tenant-ID / X-Restaurant-ID override JWT tenant when JWT restaurant_id is 0

- **Category:** multi-tenancy
- **Source:** `backend/internal/middleware/tenant.go`, `backend/internal/pkg/jwt/token.go`, `backend/internal/repository/auth_repo.go`
- **Evidence:** tenant.go:26-37 uses JWT claims only if tenantID>0 AND restaurantID>0; otherwise tenant.go:39-63 falls through to headers and trusts them for BOTH tenant and restaurant, even when a valid JWT is present. token.go:14 has RestaurantID int omitempty, and users.restaurant_id is nullable (003_create_users.sql:7, ON DELETE SET NULL; auth_repo.go:230 scans *int), so any user with NULL restaurant_id (e.g. the seeded super admin, migration 109 inserts no restaurant_id) gets restaurant_id 0 in the token and therefore full header control of tenant+restaurant on every wrapProtected / wrapWithPermission route (main.go:604-639). Static analysis; not exercised dynamically.
- **Suspected root cause:** Backward-compat header fallback applied to authenticated requests; "claims present but incomplete" is treated like "no claims".
- **Recommended fix:** If claims != nil, ALWAYS use claims.TenantID; reject (403) when restaurant_id is 0 unless the route is explicitly tenant-level, and never read headers on authenticated routes. Restrict header fallback to a separate explicit public middleware that resolves tenant/restaurant from the slug server-side (see DB-002).
- **Required regression tests:**
  - Authenticated user with JWT restaurant_id=0 sends X-Tenant-ID of another tenant -> 403, no data
  - JWT tenant A + header tenant B -> always tenant A
  - Table-driven test across every wrapProtected route group

### DB-002 [P0] Unauthenticated public order endpoints trust client-supplied tenant/restaurant headers; sequential IDs allow order PII enumeration and cancellation

- **Category:** multi-tenancy
- **Source:** `backend/cmd/api/main.go`, `backend/internal/handler/http/public_order_handler.go`, `backend/internal/middleware/tenant.go`
- **Evidence:** main.go:472-484 register POST/GET/DELETE /api/v1/public/orders*, wrapped only in TenantContextMiddleware (no auth; no rate limiter on these lines). With no JWT, tenant.go:40-63 takes tenant/restaurant from headers. public_order_handler.go:117-139 GetOrder(tenantID,restaurantID,orderID) and :376-404 CancelOrder use that header-derived scope; orders.id is BIGSERIAL (060_create_orders_table.sql:7) and order_number is predictable (DB-010). Restaurant IDs are discoverable via public slug endpoints (main.go:449-457). I found no proof-of-ownership (customer phone/token) check in these handlers. Orders hold name/email/phone/address (060:13-22). Not exercised dynamically.
- **Suspected root cause:** Public endpoints reuse the internal header-based tenant middleware; orders identified by guessable sequential ids/numbers with no per-order secret.
- **Recommended fix:** Resolve restaurant from slug server-side; require an unguessable per-order tracking token for status/track/cancel; redact PII in public responses; add rate limiting; remove header trust.
- **Required regression tests:**
  - Unauthenticated GET /public/orders/{id} with correct tenant/restaurant headers but no tracking token -> 401/404
  - Cancel without token rejected
  - Enumeration over ids returns no PII
  - Rate-limit test

### DB-003 [P0] Theme v2 get/update/delete/activate have no tenant check (cross-tenant IDOR, any authenticated role)

- **Category:** multi-tenancy
- **Source:** `backend/internal/handler/http/admin_theme_handler_v2.go`, `backend/internal/service/theme_service.go`, `backend/internal/repository/theme_repository.go`, `backend/cmd/api/main.go`
- **Evidence:** Routes main.go:1462-1469 use wrapProtected only (no RBAC permission). admin_theme_handler_v2.go:322 UpdateTheme and :356 DeleteTheme pass only the path id to the service; theme_service.go:224-300 does repo.GetByID(ctx,id) then Update/Delete with no comparison of theme.TenantID to the caller tenant (theme.TenantID is read only for cache invalidation, :268,:298). theme_repository.go:300-330 UPDATE themes_v2 ... WHERE id=$10 and :355-360 DELETE FROM themes_v2 WHERE id=$1 carry no tenant_id. Cache key theme:%d (theme_service.go:200) is tenant-agnostic. Static analysis; not exercised dynamically.
- **Suspected root cause:** Handler/service signatures omit tenant context; repository methods keyed by id only.
- **Recommended fix:** Pass tenantID/restaurantID from context into service+repo; add AND tenant_id=$ AND restaurant_id=$ to every themes_v2 statement; add RBAC permission to theme routes; tenant-scope cache keys.
- **Required regression tests:**
  - Tenant B PUT/DELETE/activate/GET on tenant A theme id -> 404
  - Cache does not serve tenant A theme to tenant B
  - Cashier role cannot delete theme

### DB-004 [P1] RBAC RevokePermission deletes role permissions by role_id without tenant check; RBAC admin routes have no authorization-level gate

- **Category:** multi-tenancy
- **Source:** `backend/internal/handler/http/rbac_permission_handler.go`, `backend/internal/repository/role_permission_repository.go`, `backend/cmd/api/main.go`
- **Evidence:** rbac_permission_handler.go:94-117 reads roleId/moduleId from the path and calls role_permission_repository.go:139-145 DELETE FROM role_permissions WHERE role_id=$1 AND module_id=$2 (no tenant_id, no role ownership lookup), unlike AssignPermission (:72-84) which verifies the role via roleRepo.GetByID(tenant,restaurant,role). Routes main.go:925-936 use wrapAuthOnly: any authenticated user reaches role/permission/user-role management; handlers check tenant membership but grep for admin/role checks in rbac_*_handler.go found only the system-role delete guard (rbac_role_handler.go:328). Privilege-escalation reachability is for the auth auditor to confirm.
- **Suspected root cause:** Repository method omitted tenant param; RBAC handlers rely on authentication only.
- **Recommended fix:** Add tenant_id to Revoke/RevokeAll/GetByRoleAndModule; verify role ownership in handler; gate all /rbac mutations behind an explicit admin permission.
- **Required regression tests:**
  - Tenant B revoke on tenant A role -> 404
  - Cashier JWT cannot assign permissions/roles
  - Regression for RevokeAll

### DB-005 [P1] Row-Level Security does not exist in the deployed schema; RLS script and manager are dead code; documentation claims RLS

- **Category:** database
- **Source:** `backend/internal/security/rls_policies.sql`, `backend/internal/security/rls.go`, `backend/migrations`, `backend/tests/integration/security_rls_integration_test.go`, `docs-organized/05-INFRASTRUCTURE/MULTI_TENANCY/TENANT_ISOLATION_SECURITY.md`
- **Evidence:** grep for "ENABLE ROW LEVEL|CREATE POLICY" across backend/migrations returns 0 hits (DB was down; static). The only RLS SQL is internal/security/rls_policies.sql (340 lines; policies for restaurants/orders/customers/payments/order_items, :57-61), not in migrations/, not run by cmd/migrate, and it references a customers table that no migration creates (accounting uses customers_accounting). grep shows NewRLSManager/SetTenantContext are not used outside internal/security and tests. App connects as DB_USER=postgres (docker-compose.yml:86, config.go:54, .env.prod.example:13), a superuser that bypasses RLS even if enabled. TENANT_ISOLATION_SECURITY.md:10,18-38 presents RLS as implemented.
- **Suspected root cause:** RLS designed in a later phase but never integrated into migrations or the request path; docs written ahead of implementation.
- **Recommended fix:** Either implement properly (migration enabling RLS+FORCE on all tenant tables, non-superuser app role, per-request SET LOCAL app.current_tenant_id in a transaction) or remove claims from docs. Treat app-layer scoping as the only control until then.
- **Required regression tests:**
  - Integration test with non-superuser role: SELECT without tenant GUC returns 0 rows
  - Cross-tenant SELECT/UPDATE blocked at DB
  - CI check that every table with tenant_id has RLS enabled

### DB-006 [P2] Shard router jump-hash ignores the key: every restaurant maps to the same shard

- **Category:** sharding
- **Source:** `backend/internal/sharding/router.go`
- **Evidence:** router.go:67-83: the loop mutates hashValue (line 77) but never uses it; bucket depends only on numShards. Computed loop result (python reproduction of lines 71-80): numShards 2->1, 3->2, 4->3, 5->3, 8->5, 16->12 for ALL keys, then % numShards (line 59). Real jump hash needs (b+1)*2^31/((key>>33)+1). tests/unit/middleware/sharding_test.go:125 claims distribution across 4 shards; not run (needs live shards). Not wired into the main API (DB-008), so no production impact today.
- **Suspected root cause:** Incorrect port of the Lamping-Veach jump consistent hash; no key-distribution unit test without live DB.
- **Recommended fix:** Replace with a correct jump hash or a persisted tenant->shard directory table; add a pure unit test asserting uniform distribution over 100k keys and minimal movement when adding a shard.
- **Required regression tests:**
  - Distribution chi-square test over 100k keys, N=4
  - Consistency test: adding shard N+1 moves ~1/(N+1) keys
  - Determinism test

### DB-007 [P1] Shard migration tool would copy entire tables (not one restaurant), uses invalid SQL, and rollback deletes ALL target data

- **Category:** sharding
- **Source:** `backend/internal/sharding/migration.go`, `backend/cmd/migrate-shard/main.go`, `backend/internal/handlers/shard_migration_api.go`
- **Evidence:** migration.go:223-245: SELECT * FROM <table> ORDER BY id LIMIT/OFFSET with no restaurant/tenant filter; tables hardcoded ["orders","order_items","payments","customers","restaurants"] (:200) - omits products, journal_entries, employees etc.; "customers" does not exist. Offset carried across tables (:217). rows.Close deferred inside a loop (:244). Column list formatted with fmt.Sprintf("%s", columns) (:307, :437) which renders "[a b c]" - invalid SQL. Rollback (:373-389) runs DELETE FROM <table> on the TARGET shard with no WHERE; cmd/migrate-shard defaults --rollback-on-error=true (main.go:36), so any error wipes the target shard. Dual-write follower write is fire-and-forget (:455-461) with unsynchronised errorCount. Never executed (no shards available); CLI defaults to user postgres.
- **Suspected root cause:** Prototype written without a real schema; no tests with live data.
- **Recommended fix:** Do not use against real data. Rewrite as a per-tenant scoped, FK-ordered, transactional copy with checksum verification and a safe rollback (delete only rows tagged with the migration).
- **Required regression tests:**
  - Migrate one tenant between two throwaway shards; row counts + checksums per table equal
  - Other tenants on target untouched after rollback
  - Failure injection mid-copy

### DB-008 [P2] Sharding and scaling packages are not wired into the API (experimental); documentation claims production-ready

- **Category:** sharding
- **Source:** `backend/cmd/api/main.go`, `backend/internal/sharding`, `backend/internal/scaling`, `backend/internal/handlers/scaling_api.go`, `docs-organized/05-INFRASTRUCTURE/MULTI_TENANCY/IMPLEMENTATION_SUMMARY.md`
- **Evidence:** cmd/api/main.go imports (top ~60 lines) include no internal/sharding, internal/scaling or internal/handlers. grep: ShardingMiddleware/NewShardRouter appear only in internal/sharding, internal/middleware/sharding.go, internal/handlers/sharding_example.go and tests. Repositories take a single *sql.DB. scaling/autoscaler.go:226,267 "In production, call cloud provider API" (no-op), failover_manager.go:250 stub, load_balancer.go:239 stub GeoIP. scaling_api.go:34-35 registers "servers/:id/metrics" (not Go 1.22 pattern syntax) and has no auth. IMPLEMENTATION_SUMMARY.md:3,422 says COMPLETED / production-ready; DATABASE_SHARDING.md:3 says In Progress. ~3,970 lines unused.
- **Suspected root cause:** Scaffolding committed as "complete" without integration.
- **Recommended fix:** Mark packages experimental in docs, move behind a build tag or delete; do not advertise horizontal scalability.
- **Required regression tests:**
  - CI check that unused packages are excluded
  - If revived: end-to-end test routing two tenants to two shards through the real API

### DB-009 [P2] Shard key is restaurant_id but accounting/RBAC/subscription data is tenant-scoped; cross-shard transactions undefined

- **Category:** sharding
- **Source:** `backend/internal/middleware/sharding.go`, `backend/internal/sharding/router.go`
- **Evidence:** middleware/sharding.go:35-49 routes by RestaurantID only; users, tenants, roles, role_permissions, subscriptions and payments (115_create_payments_table.sql: tenant_id only) are tenant-level and several repos query across a tenant (subscription_repo_impl.go). No code coordinates 2PC/saga. GetAllConnections (router.go:134-160) opens unmanaged pools and leaks them on error. Global UNIQUEs (orders.order_number 060:10, tenants.slug) cannot be enforced across shards. DSN hardcodes sslmode=disable (router.go:81).
- **Suspected root cause:** Shard key chosen without data-model analysis.
- **Recommended fix:** Shard by tenant_id (all restaurants of a tenant co-located); keep a global control-plane DB for tenants/users/plans; replace global uniques with (tenant_id, ...) uniques.
- **Required regression tests:**
  - Tenant-level report across its restaurants runs on a single shard
  - Global uniqueness test for control-plane tables

### DB-017 [P2] Repository methods without tenant/restaurant predicates (defence-in-depth gaps); GetOrderItems/GetOrderStatusHistory ignore tenant

- **Category:** multi-tenancy
- **Source:** `backend/internal/repository/driver_repo.go`, `backend/internal/repository/order_repo.go`, `backend/internal/repository/homepage_section_repository.go`, `backend/internal/repository/restaurant_theme_repo.go`, `backend/internal/repository/theme_repository.go`, `backend/internal/repository/payment_repo_impl.go`, `backend/internal/repository/user_repo.go`
- **Evidence:** See table in 06 and evidence/database-evidence/repository-tenant-scan.txt. driver_repo.go: Update (:198 WHERE id=$16), Delete (:228), UpdateStatus (:249), UpdateAvailability, UpdateLocation, assignments, stats are by id only; safety depends on handler pre-checks comparing driver.TenantID (admin_driver_handler.go:263,292,359,...). Drivers routes use wrapProtected only (main.go:1397-1403, no RBAC). order_repo.go:693-700 GetOrderItems and :756-762 GetOrderStatusHistory accept tenantID but query WHERE order_id=$1 only (callers check the order first, e.g. public_order_handler.go:285). homepage_section DeleteSection/website_settings are restaurant-scoped but not tenant-scoped; restaurant_theme_repo UpdateSection/DeleteSection by theme_id only; payment_repo_impl.go UpdatePaymentOrder/UpdatePaymentStatus WHERE id=$n only; user_repo.go UpdatePassword/UpdateProfile by user id (callers not individually traced).
- **Suspected root cause:** Authorization done in handlers instead of enforced in the data layer.
- **Recommended fix:** Make tenant_id (and restaurant_id where applicable) mandatory parameters of every repository method and part of WHERE; CI lint scanning repo SQL for tenant predicates.
- **Required regression tests:**
  - Per-repository cross-tenant matrix test (tenant B cannot read/update/delete tenant A row by id)
  - Static SQL-lint test in CI

### DB-021 [P3] KDS websocket is JWT-scoped (good) but CheckOrigin always true and token in query string; upload serving unverified

- **Category:** multi-tenancy
- **Source:** `backend/internal/handler/ws/kds_ws_handler.go`, `backend/internal/handler/ws/kds_hub.go`, `backend/internal/usecase/product_usecase.go`
- **Evidence:** kds_ws_handler.go:98-114 validates the JWT from ?token= and derives tenantID/restaurantID from claims (not headers); room key tenant:%d:restaurant:%d (kds_hub.go:171) - SCOPED. CheckOrigin returns true (:83-87, TODO) enabling cross-site websocket hijack if a token is obtainable; token in URL leaks to logs. Uploads: product images are stored under ./uploads/products/<tenant>/<product>/ (product_usecase.go:299-304) but I found no route serving /uploads in cmd/api/main.go; serving (nginx/static) and ACLs not verified.
- **Suspected root cause:** Dev shortcuts.
- **Recommended fix:** Whitelist origins; use first-message auth or short-lived ws ticket; verify upload serving path/ACL.
- **Required regression tests:**
  - WS from disallowed Origin rejected
  - Tenant A token cannot receive tenant B room messages


## 5. Not verified

- No live DB: no pg_catalog counts, no EXPLAIN, no actual trigger/constraint list; actual `schema_migrations` state of any environment unknown.
- No dynamic requests: DB-001/002/003/004 exploitability is from code reading and must be confirmed in Phase 2.
- Sharding/scaling tests were not run (need live shards).
- Authorization (who may call RBAC/driver/theme routes) is cross-checked only for route wrappers; role semantics belong to the auth auditor.
- Upload serving path, payment webhook/worker SQL scoping, hr_report dynamic where clauses: not fully traced.
