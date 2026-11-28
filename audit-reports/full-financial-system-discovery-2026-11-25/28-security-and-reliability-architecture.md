# 28 - Security and Reliability Architecture

Baseline SHA `133260db66b1143c254d79476daa7650b41a33f0`. Read-only discovery: source, migrations, compose/CI/ops files, and SELECT/catalog queries against the dev DB (`pos-postgres`, database `pos_saas`). No code, git, docker or data was modified; no mutating API call was made.

## 0. Evidence labels and the dev-DB caveat

Labels: CONFIRMED-BY-SOURCE (S), CONFIRMED-BY-READONLY-RUNTIME (RT), DOCUMENTED-INTENT (DOC), INFERRED (INF), UNKNOWN (UNK), REQUIRES-EXECUTION (REX).

**Dev-DB caveat (RT, important).** The dev stack connects as superuser `postgres` (`docker-compose.yml:57,86` `DB_USER: postgres`; `rolsuper=t, rolbypassrls=t`, queried from `pg_roles`). Superusers bypass RLS, table grants and `NOT NULL`-style role checks. Therefore **no runtime observation on this stack proves anything about production isolation, RLS enforcement, grant sufficiency or role restriction**. Catalog facts (policies, grants, triggers) are valid as to what is installed; whether the API works under `pos_app`/`pos_accounting` here is UNK/REX. The dev DB also lacks the full runtime-role provisioning: `pos_app` has explicit grants on only 7 tables (RT: `information_schema.role_table_grants`), because `backend/scripts/db/create_runtime_roles.sql` (which does `GRANT ... ON ALL TABLES ... TO pos_app`) was evidently not run here; the role `pos_auth_lookup` does not exist (RT: `pg_roles`) although migration 192 is recorded in `schema_migrations`. Dev DB holds 2 tenants, 0 posted journals, 0 fiscal_periods, 1 super_admin_audit_log row, so there is little data to inspect.

Also relevant: dev compose sets no `RBAC_ENFORCEMENT`, so dev runs in **audit mode** (S: `internal/middleware/permission.go:37-45` - empty value + non-production = audit). Dev runtime therefore never denies on RBAC.

---

## 1. Authentication

| Aspect | Intended (docs/config) | Actual (source) | DB-enforced (catalog) | Unknown / requires execution |
|---|---|---|---|---|
| Credential | bcrypt passwords; HttpOnly host-only session cookie `poos_session` (DOC `docs/security/auth-distributed-controls.md`) | bcrypt compare in `internal/usecase/auth_usecase.go:264`; browser sessions stored server-side in Redis/KV (`internal/service/auth_session.go`); idle 30 min, absolute 8 h (`auth_session.go:17,19`) (S) | `users.password_hash`; nothing at DB level | Bcrypt cost/parameters in prod UNK |
| JWT | >=32-char secrets; super-admin secret independent (DOC; `cmd/api/signing_secrets.go`) | HS256 only (`internal/pkg/jwt/token.go:70`), claims `user_id,tenant_id,restaurant_id,role,exp,iat`; **no `jti`, no revocation list, no per-request user-status recheck** (S `token.go:12-20,74-105`). `JWT_SECRET` default `"your-secret-key-change-this"` (`internal/config/config.go:84`); length/default rejected only when `Server.Env=="production"` (`config.go:91`) or env in {production,staging} for signing secrets (`signing_secrets.go:9`) | n/a | Whether prod secrets meet policy: UNK (secrets not inspected) |
| Login origin / CSRF | Origin allow-list protects session-creating POSTs (DOC `cors.go` comment) | `RequireTrustedOrigin` wraps register/login/login-confirm/refresh/logout and super-admin login (`cmd/api/main.go:556-562,1909`). `SessionAuthMiddleware` requires trusted Origin on unsafe methods **only on the cookie path** (`internal/middleware/auth.go:40-42`). **Bearer-header path has no Origin check** (S, correct for non-cookie auth). Double-submit `CSRFMiddleware` (`internal/middleware/csrf.go`) is **defined but not referenced by any route** (S: grep shows only definition) | n/a | Missing `Origin` header on same-site non-browser clients is rejected (`IsTrustedOrigin("")` false) - REX to confirm for native/POS clients |
| Multi-tenant login | One-time tenant-selection proof, 5 min TTL, SHA-256 key claimed with `SET NX` in Redis (DOC; S `auth_usecase.go:28-34,243-247`) | Login scans **all users with the email across tenants** (`auth_repo.go:251` `FindUserByEmailAllTenants`), keeps password-matching accounts, issues selection token if >1 (S). In-memory consumption fallback exists (`auth_usecase.go:43-76`) for non-prod | Legacy: `users` table has tenant_id; unique (tenant_id,email) assumed | Redis-down behaviour in prod: DOC says generic 503; REX |
| Token exchange | Super-admin one-time exchange "no longer supported" | `AuthHandler.TokenExchange` returns **410 Gone** (`auth_handler.go:411-414`); not routed (`main.go:561-563` comment) (S). Super-admin handoff is server-to-server `POST /super-admin/tenants/{id}/dashboard-login` and impersonation sessions (`main.go:1931,1938`) | `impersonation_sessions` table (no RLS) | Dashboard-login token lifetime and audit coverage: UNK |
| Identity/membership schema (migrations 191/192) | Canonical identity model; "runtime reads remain disabled until separately reviewed" (DOC 191 header; `auth_lookup_db.go` comment "deliberately opt-in") | Pool `connectAuthLookupDatabase` returns nil unless `AUTH_LOOKUP_DB_*` set (S `auth_lookup_db.go:35-40`); **login handler still uses legacy `users` path** (S) | RT: `auth_identities`, `identity_email_challenges`, `tenant_memberships`, `membership_restaurant_grants`, `membership_role_grants` all `ENABLE`+`FORCE` RLS with single policy `*_migrator_only` (`TO pos_migrator USING true`) (RT `pg_policies`) - so `pos_app` sees **zero rows** by design. `auth_private` schema exists with 2 SECURITY DEFINER lookup functions (mig. 192). Tables are empty (`auth_identities` count 0) (RT) | `pos_auth_lookup` role absent on dev (RT). Behaviour with the role present: REX |
| Refresh / revocation | Logout revokes sessions; refresh rotates (DOC `SERVIO_PRODUCTION_READINESS.md:45`) | `AuthSessionStore.Rotate/Revoke` delete the Redis key (`auth_session.go:78-130,227-232`); `POST /auth/refresh` requires existing session (`main.go:560`). **Raw Bearer JWT cannot be revoked before `exp`** (default `JWT_EXPIRY` 24h, `config.go:86`); the cookie session wraps a JWT whose expiry is clamped to the session absolute lifetime (`token.go:42-48`) | n/a | Revocation of sessions on password change / user disable / role change: no code found that sweeps sessions (INF) - REX |
| Account lockout | - | No per-account lockout for staff login (grep `lockout|failed_attempts` hits only super-admin code `super_admin_repo.go:62`) (S). Only per-IP rate limits (login 20/min, register 10/min, refresh 20/min, forgot-password 5/min) (`main.go:556-565`) | n/a | Brute-force exposure behind shared NAT: INF |
| Stub-mode login | none intended | When DB connect fails, API continues in "stub mode" (`main.go:99-105`); `FindUserByEmailAllTenants` then returns a mock **owner** user with password `password123`, tenant 1 (S `auth_repo.go:253-270`). Production refuses nil DB (`startup_guards.go:25-27`), **but `IsProductionMode()` is true only for `production|prod`** (`internal/config/runtime_env.go:20-23`); ops/aws runs `ENV: staging` (`ops/aws/compose.json:130`) so the guard does not run on staging | n/a | Whether staging API can start without DB: staging entrypoint waits for DB (`compose.json:192`), so exposure is conditional (INF); REX to confirm |

## 2. Authorization and RBAC

**Model (S).** Modules are integer IDs in `modules` (RT: 1 PRODUCTS, 2 CATEGORIES, 3 HR, 4 NOTIFICATIONS, 5 SETTINGS, 6 USERS, 7 ROLES, 8 REPORTS, 9 THEMES, 10 ORDERS, 11 LEAVES, 12 SALARIES, 13 KDS, 14 TABLES, 15 ACCOUNTING, 16 SUPPLIERS, 17 FINANCIAL_REPORTS, 18 PERFORMANCE, 19 TRAINING, 20 CONTRACTS). Levels are a linear order NONE<READ<EXPORT<WRITE<DELETE<ADMIN (`internal/domain/role.go:119-127`). So **WRITE implies EXPORT and READ; there is no separate "post", "approve", "reverse", "close period" or "export-only-financial" level** - the same WRITE/EXPORT ladder is used for everything.

**Evaluator (S, `permission.go:111-148`).** Highest of: DB `role_permissions` per user role, built-in KDS/TABLES fallbacks, and **JWT claim `role` in {owner, admin} => ADMIN on every module**. The claim is taken from the token, not re-read from DB per request (INF consequence: demoting an owner in DB does not take effect until the token/session expires).

**Enforce vs audit (S `permission.go:36-47`, `main.go:714-722`).** `RBAC_ENFORCEMENT` empty => enforce in production, **audit (log only, request proceeds) otherwise**; `audit` refused in production. Declared enforce in `docker-compose.prod.yml:25`, `ops/aws/compose.json:142`, `docker-compose.pos-test.yml:59`, `docker-compose.phase2a-e2e.yml:121`, CI `financial-acceptance.yml:124`. **Dev compose does not set it => audit** (RT-adjacent: `docker-compose.yml`). If the evaluator is nil or the lookup errors, enforce mode returns 503, audit mode lets the request through (S `permission.go:75-93`).

**Wrapper inventory (S: scripted extraction of `mux.Handle*` in `cmd/api/main.go`, `pos_routes.go`, `health_routes.go`, `translation_routes.go`; 376 registrations).** `wrapWithPermission` 261, `wrapAuthPermission` 14, `wrapKDSWithPermission` 8, POS `guards.requirePermission` 1 (=284 with an explicit module/level); 66 with no tenant wrapper (public, auth, super-admin which use `wrapSuperAdmin`, health, mock); `wrapProtectedNoStatus` 8, `wrapProtected` 9, `wrapAuthOnly` 3, bare `authenticateRequest` 1.

### 2.1 Financial routes and their permission wrapper

| Area | Routes (examples) | Module:level | Notes |
|---|---|---|---|
| Accounting master/journals | `POST /accounting/journal-entries`, `/{id}/post`, `/{id}/void`, `PUT /{id}`, accounts CRUD, bank accounts, customers | 15:WRITE (reads 15:READ, deletes 15:DELETE) (`main.go:1104-1167`) | **post and void share 15:WRITE with create/edit**; a seeded MANAGER role gets WRITE on module 15 (mig. 155:66) so managers can post and void journals |
| Paid-order event recovery | `GET /accounting/paid-order-events/failed`, `POST /{id}/redrive` | 15:READ / 15:WRITE (`main.go:1102-1103`) | redrive audited in `accounting_event_redrive_audit` (mig. 188) |
| AR/AP | customer-invoices create/send/void, supplier-invoices create/approve/void, payments received/made, bank reconciliation, recurring journals run, budgets, exchange rates, e-invoice submit | 15:WRITE / 15:READ | approve/void of supplier invoice is plain WRITE |
| Import | `POST /accounting/csv|excel/*/import` (chart of accounts, journal entries) | 15:WRITE (`main.go:1171-1179`) | journal import = bulk financial write at WRITE level |
| Export | `GET /accounting/csv|excel/*/export` | 15:EXPORT (`main.go:1170-1178`) | |
| Tax | `POST /accounting/tax/periods/{id}/close`, periods CRUD, vat-report | 17:WRITE / 17:READ (`main.go:1492-1498`) | closing a **tax** period is 17:WRITE; there is **no route that creates, closes or reopens `fiscal_periods`** (grep of `routes.manifest`: no "fiscal") |
| Reports | `/reports/financial/*` (11) | 17:READ; `/reports/overview|sales|profitability|inventory` 8:READ; `POST /reports/export` 8:EXPORT; HR reports 8:READ, social-insurance 8:EXPORT (`main.go:1467-1490`) | `GET /audit/export` is 17:EXPORT (`main.go:1491`) |
| Orders | `/admin/orders/*`: status, **payment-status** (10:WRITE), cancel item (10:WRITE), `DELETE /admin/orders/{id}` (10:DELETE), **`GET /admin/orders/export` (10:READ)** (`main.go:1432-1440`) | | **Order export needs only READ**, below the EXPORT level used elsewhere |
| POS | `POST /pos/orders` | 10:WRITE via `guards.requirePermission` (`pos_routes.go:18`) | |
| Payroll | `/hr/salaries*` incl. `POST /{id}/mark-paid` | 12:WRITE (`main.go:1091-1096`) | payroll accounting transaction on the accounting pool |
| Suppliers | `/suppliers*` | 16:* | |
| Inventory | `POST /inventory/adjust` | 1:WRITE (`main.go:785`) | module 1 = PRODUCTS, not a stock/finance module |
| KDS | bump/recall 13:WRITE, analytics 13:ADMIN | via `wrapKDSWithPermission` | |

### 2.2 Routes with NO permission wrapper that touch money or billing

| Route | Wrapper | Evidence | Assessment |
|---|---|---|---|
| `POST /subscription/initiate`, `GET /subscription/payments`, `GET /subscription/current`, `POST /subscription/cancel` | `wrapProtectedNoStatus` (auth + tenant only, **no RBAC**, no tenant-status gate by design) | `main.go:764-767` | SaaS billing (Paymob). **Any authenticated tenant user, whatever role, can start a payment or cancel the subscription** (S). Role intent UNK |
| `POST /public/orders`, `/quote`, `/validate`, `GET /public/orders/{id}[...]`, `DELETE /public/orders/{id}` | unauthenticated; `RateLimitByIP` (in-memory, **not** the Redis store) + `tenantContext` (X-Tenant-ID/X-Restaurant-ID verified against `restaurants`) (+ tenant status for POST) | `main.go:598-611` | Public by design; creates orders with money amounts. Access to an order is by id/number/token; see 6 and 10 |
| `POST /webhook/paymob` | HMAC (`PAYMOB_HMAC_SECRET`) + `RateLimitByIP(60)` | `main.go:656-660`; `paymob_client_live.go:186-203` constant-time compare | **`MockPaymobClient.VerifyHMAC` always returns true** (`paymob_client_mock.go:52-54`); safe only because production refuses mock mode (`paymob_mode.go:25-31`). Staging (`ENV=staging`) is not production, so mock mode is *allowed* there if `PAYMOB_MODE` is empty (INF) |
| `GET/POST /mock/paymob/*` | none | `main.go:667-668` | registered only when `PAYMOB_MODE` resolves to mock (`main.go:517-521`) |
| `POST /super-admin/payments/{id}/refund` | `wrapSuperAdmin` (separate super-admin JWT, IP whitelist, audit) | `main.go:1922` | refund of **tenant SaaS payments**, not restaurant sales |
| `/users/{userId}/settings*`, `/profile` | `wrapProtected` (no RBAC) | `main.go:1197-1203` | not financial; listed for completeness; IDOR protection is by RLS + handler (`user_settings` RLS) |

**Count: 0 tenant-ledger financial routes (accounting, payroll, reports, orders admin, suppliers) lack a permission wrapper (S); 4 tenant-authenticated SaaS-billing routes lack RBAC; 9 public order routes and 3 payment-provider routes are intentionally unauthenticated.**

## 3. Tenant and restaurant resolution

| Item | Intended | Actual (S) | DB-enforced | Unknown |
|---|---|---|---|---|
| Authenticated tenant | From JWT; client header ignored (DOC `tenant.go:36-45`) | `tenant.go:62-75`: tenant = `claims.TenantID`; `X-Tenant-ID` differing is logged and ignored. Restaurant = claim, else `X-Restaurant-ID` verified by `SQLRestaurantVerifier` against `restaurants(id,tenant_id)` (`tenant.go:75-84,111-129`); wired at `main.go:593` | `cfr002_guard_restaurant_scope_write` blocks changing `restaurants.tenant_id` (RT trigger) | Whether every handler uses `middleware.GetTenantID` rather than a body/query field: INF, needs per-handler review (REX: cross-tenant test suite) |
| Anonymous (storefront) | headers must be consistent | verifier checks the pair (S `tenant.go:86-100`) | n/a | Slug-to-tenant resolution for `/public/restaurants/{slug}/*` is in handlers (not reviewed) UNK |
| JWT with restaurant_id=0 | owner-level multi-restaurant | header-selected restaurant verified against the token's tenant | | |

## 4. Tenant isolation: app filters vs RLS

**App-level (S, INF).** Repositories pass `tenant_id`/`restaurant_id` into SQL `WHERE` clauses (e.g. `journal_entry_repo.go:442-447`, `auto_posting_repo.go:119-128`). Exhaustiveness across all ~100 tables is **not proven** (REX: grep/AST audit + two-tenant integration run).

**RLS (RT, catalog `pg_class.relrowsecurity/relforcerowsecurity`, `pg_policies`).** Of 102 public tables, **20 have RLS enabled and FORCEd** (20/20 forced). Mechanism: policies compare `tenant_id` with `NULLIF(current_setting('app.tenant_id', true),'')::bigint`, set transaction-locally by `database.WithTenantTx` (S `internal/pkg/database/tenant_tx.go:11-21,101`). `WithTenantTx` is used in only 12 source files (S grep): accounting_advanced, report_handler, auth_audit, notification, payment, rbac_role, role_permission, user_role, user_settings repos/handlers and payment_usecase.

| # | Table | Policy | Roles | Cmd |
|---|---|---|---|---|
| 1 | audit_export_logs | tenant_restaurant_scope | pos_app | ALL |
| 2 | auth_identities | migrator_only | pos_migrator | ALL (true) |
| 3 | e_invoice_submissions | tenant_restaurant_scope | pos_accounting, pos_app | ALL |
| 4 | email_change_requests | tenant_scope (via users) | pos_app | ALL |
| 5 | identity_email_challenges | migrator_only | pos_migrator | ALL |
| 6 | invoice_delivery_logs | tenant_restaurant_scope | pos_app | ALL |
| 7 | membership_restaurant_grants | migrator_only | pos_migrator | ALL |
| 8 | membership_role_grants | migrator_only | pos_migrator | ALL |
| 9 | notification_preferences | tenant_scope (via users) | pos_app | ALL |
| 10 | notifications | user_tenant_restaurant_scope | pos_app | ALL |
| 11 | password_change_history | tenant_scope (via users) | pos_app | ALL |
| 12 | payments | tenant_isolation (pos_app) + migrator_access | pos_app; pos_migrator | ALL |
| 13 | permission_audit_logs | tenant_scope | pos_accounting, pos_app | **SELECT only** |
| 14 | report_export_audit_logs | tenant_restaurant_scope | pos_app | ALL |
| 15 | role_permissions | restaurant_scope + migrator_access | pos_app; pos_migrator | ALL |
| 16 | tenant_memberships | migrator_only | pos_migrator | ALL |
| 17 | user_auth_audit_logs | user_auth_audit_tenant_scope | pos_app | ALL |
| 18 | user_roles | restaurant_scope + migrator_access | pos_app; pos_migrator | ALL |
| 19 | user_settings | tenant_scope (via users) | pos_app | ALL |
| 20 | user_settings_audit_log | user_settings_audit_tenant_scope | pos_app | ALL |

Consequences (S/RT): (a) The superuser/owner bypass RLS unless FORCE is set; FORCE is set, but **superusers and BYPASSRLS roles still bypass** (the dev backend does). (b) Table owner is `postgres` for all 102 tables in dev (RT `pg_tables`); production intent is owner `pos_migrator` (`create_runtime_roles.sql`, `transfer_public_ownership.sql`) (DOC). (c) RLS policies for `pos_accounting` exist only on `e_invoice_submissions` and `permission_audit_logs`; **the accounting pool operates on non-RLS tables with no DB tenant guard**.

**Tables with `tenant_id` but NO RLS: 57 of the 70 tables having a `tenant_id` column** (RT `information_schema.columns` x `pg_class`). Only 13 of the 70 are RLS-protected. The 57: account_budgets, accounting_event_outbox, accounting_inventory_movement_events, attendance, auto_posting_rules, bank_accounts_ledger, bank_reconciliation_lines, bank_reconciliation_sessions, categories, chart_of_accounts, customers_accounting, drivers, employee_roles, employees, employment_contracts, exchange_rates, fiscal_periods, fixed_assets, general_ledger, homepage_sections, hr_audit_logs, insurance_tracking, inventory_movements, invoice_customer_lines, invoice_supplier_lines, invoices_customer, invoices_supplier, journal_entries, journal_entry_lines, kds_queue, kds_settings, leaves, order_item_cancellations, order_items, order_stock_reservations, orders, payments_made, payments_received, performance_reviews, products, recurring_journal_templates, resignations, restaurant_tables, restaurant_themes, restaurants, roles, salaries, supplier_payment_schedules, suppliers, table_layouts, tax_periods, tenant_subscription_history, themes, themes_v2, training_records, users, website_settings. Additionally `customers_accounting`, `journal_*`, `general_ledger`, `orders` etc. are all in this set - **every core financial table (orders, order_items, journals, GL, invoices, payments_received/made, salaries, fiscal_periods, chart_of_accounts) has app-level isolation only.** `webhook_retry_queue`, `impersonation_sessions`, `super_admin_audit_logs`, `order_status_history`, `product_audit_log` have no tenant RLS either (webhook queue has no `tenant_id` column at all, RT `\d`).

SECURITY DEFINER functions callable by `pos_app` that bypass RLS by design: `resolve_payment_tenant_by_id`, `resolve_payment_tenant_by_paymob_order`, `get_payment_metrics` (RT `pg_proc.prosecdef`, mig. 174:61-63) - used for webhooks before tenant is known (S/DOC).

Whether RLS + app code work end-to-end under `pos_app`: CI has a `POS_RLS_APP_URL` job (`financial-acceptance.yml:99-102`) (DOC/CI config); results not re-run here (REX).

## 5. Database roles

| Role | Attributes (RT) | Grants (RT, dev DB) | Intended grants (S) |
|---|---|---|---|
| `postgres` | super, bypassrls, login | owner of all 102 public tables | bootstrap only; dev app uses it |
| `pos_migrator` | LOGIN, no super/createdb/createrole/replication/bypassrls | 20 table grants (5 identity tables + RLS-policy tables) | owns schema & runs migrations (`create_runtime_roles.sql:20-27`); **owns DB, so can disable triggers/alter policies** |
| `pos_app` | same restricted attributes; no role memberships (`pg_auth_members` empty) | explicit grants on 7 tables: outbox I/S/U, redrive_audit I/S, inventory_movement_events I/S, payments I/S/U, role_permissions DML, user_roles DML, webhook_retry_queue I/S/U | **all public tables DML via default privileges and `GRANT ... ON ALL TABLES`**, then REVOKE INSERT/UPDATE/DELETE on journals+GL (mig. 182:47), REVOKE on reservation/cancellation ledgers, outbox DELETE, `schema_migrations` (`create_runtime_roles.sql:35-70`) |
| `pos_accounting` | restricted; `LOGIN` only if password supplied (`create_accounting_role.sql`) | 106 table-grant rows: DML (incl. **DELETE**) on 24 accounting tables, SELECT on orders, order_items, products, restaurants, tenants, employees, salaries (col UPDATE), permission_audit_logs, etc. (RT) | "never core business-table writes or schema CREATE" (DOC `RUNBOOK.md:221-224`) |
| `pos_auth_lookup` | **role does not exist on dev** (RT) | none | read-only EXECUTE on two `auth_private` functions (mig. 192) |

Prod startup guards (S, `cmd/api/startup_guards.go:60-72`): runtime role must be exactly `pos_app` and not superuser/BYPASSRLS; also rejects the seeded super-admin default password hash (`startup_guards.go:11-12`). **Applies only when `ENV` is `production`/`prod`**; `ENV=staging` skips them (S `runtime_env.go:20-23`, `main.go:109-113`).

`pos_app` can still `DELETE` and `UPDATE` ordinary business tables (orders, payments_received, invoices, salaries...) in a correctly provisioned DB - no DB-level immutability there (see 9).

## 6. Accounting DB pool separation and fallback

- `connectAccountingDatabase` (S `cmd/api/accounting_db.go:27-58`): if env is development/dev/local **or `test`** and neither `ACCOUNTING_DB_USER` nor `ACCOUNTING_DB_PASSWORD` is set, **the accounting repos silently reuse the primary pool** (`return primary, nil`). In every other env (including staging and production) a dedicated `pos_accounting` login is required; it must not be super/createdb/createrole/bypassrls (`accounting_db.go:49-54`). Pool size 5/2 (`accounting_db.go:20-21`). Failure => `log.Fatalf` (`main.go:115-118`).
- **Config discrepancy (S):** `docker-compose.prod.yml` backend environment (`:12-38,170-175`) contains **no `ACCOUNTING_DB_USER/PASSWORD`**, whereas the guard makes them mandatory when `ENV=production`; `ops/aws/compose.json:153-154` does supply them (staging). Consequence (INF): the committed `docker-compose.prod.yml` stack would exit at startup unless the variables are injected another way (`--env-file`/override) - REX/UNK whether a different mechanism exists. Also `docker-compose.prod.yml` has no `AUTH_LOOKUP_DB_*` (expected, feature disabled).
- Outbox table grants are on `pos_app`, not `pos_accounting` (RT): the paid-order worker drives `orderUC` (`main.go:514`); outbox bookkeeping and journal writes therefore use **different pools/transactions** (INF) - see 12.
- Dev fallback means all dev/test runs exercise accounting as superuser; **role-restricted behaviour is only covered by CI jobs** (`financial-acceptance.yml`, `docker-compose.phase2a-e2e.yml:124-131`) (DOC).

## 7. Financial permissions: who may post / reverse / close / export

| Action | Route | Required | Seeded roles holding it (S mig. 155, 159; RT counts) |
|---|---|---|---|
| Create/edit draft journal | POST/PUT journal-entries | 15:WRITE | MANAGER (WRITE), ADMIN/OWNER/admin-claim (ADMIN) |
| **Post journal** | `.../post` | 15:WRITE | same |
| **Void (reverse) posted journal** | `.../void` | 15:WRITE | same |
| Redrive failed paid-order event | `.../redrive` | 15:WRITE | same |
| Close tax period | `tax/periods/{id}/close` | 17:WRITE | module 17: MANAGER only READ by seed (mig. 155:79-82) => only ADMIN/OWNER close tax periods |
| Close/open fiscal period | none (no API) | n/a | direct SQL only (UNK procedure) |
| Export ledger/COA/journals | `.../export` | 15:EXPORT | external-accountant role gets 15:READ only (mig. 159:78) so it **cannot** export accounting CSV/Excel but can export financial reports (17:EXPORT) and 8:EXPORT |
| Export audit log | `GET /audit/export` | 17:EXPORT | |
| Export orders | `GET /admin/orders/export` | 10:READ | any role with orders READ |
| Mark salary paid | `mark-paid` | 12:WRITE | |
Dev DB role_permissions (RT): module 15: ADMIN x2, WRITE x1, READ x1, NONE x3; module 17: ADMIN x2, EXPORT x1, READ x1, NONE x3. Segregation of duties (maker != checker) is **not enforced** anywhere found: the same WRITE holder creates, posts and voids (S routes above; `PostEntryTx` records `posted_by` only). Whether policy requires maker-checker: UNK.

## 8. Audit logging

| Log | Written by | Covers | RLS | Gap |
|---|---|---|---|---|
| `user_auth_audit_logs` | `auth_audit_repository.go`, `auth_handler.go:272` | login/logout/session events | yes | - |
| `permission_audit_logs` | `permission_audit_log_repository.go` | role/permission changes | yes (SELECT policy only; triggers `audit_role_changes` write it) | - |
| `audit_export_logs`, `report_export_audit_logs` | report/export handlers | exports | yes | order export (`/admin/orders/export`) coverage UNK |
| `super_admin_audit_logs` | `super_admin_audit_repo.go` | super-admin actions, impersonation | **no** | table global by design |
| `hr_audit_logs` | `create_hr_audit_log`, employee trigger, `salary_repo.go` | HR changes | no | |
| `accounting_event_redrive_audit` | `accounting_event_admin_repo.go` (mig. 188) | redrives, I/S only for pos_app | no | |
| `order_status_history`, `product_audit_log` | order/product repos | status, product changes | no | payment-status changes: covered by `order_status_history` (INF) |
| Journal posting/voiding | **no audit table**; only columns `posted_by/posted_at/voided_by/void_reason` enforced by trigger (S/RT `cfr002_guard_journal_header_write`) | | | no record of *who viewed/exported ledger*; no before/after for draft edits; no log of fiscal-period or COA edits (INF: `chart_of_accounts` has scope/normal-balance triggers, no audit trigger) |
Unused code: `internal/security/{audit_log,rls,gdpr_compliance,encryption,backup_isolation}.go` are imported only by `internal/handlers/*`, which `cmd/api/main.go` does **not** import (S: import-graph grep). These controls are therefore **not part of the running API** (INF - verify build tags/other binaries).

## 9. Historical immutability (triggers/guards)

RT: all listed triggers are enabled (`tgenabled='O'`). All guard functions are `SECURITY DEFINER`, `search_path=pg_catalog,public`.

| Trigger (table) | Blocks (S, function source dumped via `pg_get_functiondef`) | Does not block / bypass |
|---|---|---|
| `cfr002_journal_header_write` (journal_entries) | INSERT unless status=draft; DELETE unless draft; status transitions only draft->{posted,void}, posted->void with unchanged accounting columns and void audit fields; restaurant must belong to tenant | `fiscal_period`/`entry_date` freeze relies on equality check only when voiding; **no closed-fiscal-period check** |
| `cfr002_journal_line_write` (journal_entry_lines) | any INSERT/UPDATE/DELETE once header is not draft; cross-tenant/restaurant account use; moving lines | |
| `cfr002_general_ledger_insert` / `cfr002_general_ledger_immutable` (general_ledger) | INSERT only for posted/void journal, one effect per line, exact mirror on void; **all UPDATE/DELETE rejected** | **No TRUNCATE trigger** (RT: no statement-level triggers); table owner/superuser can `TRUNCATE`, `ALTER TABLE ... DISABLE TRIGGER`, `SET session_replication_role=replica` (INF - PostgreSQL semantics; REX on a restricted stack). `pos_accounting` has no TRUNCATE (RT grants) |
| `cfr002_journal_integrity_*` (3 deferred constraint triggers) | commit-time debit=credit/ledger consistency (mig. 180) | |
| `cfr003_lock_account_after_...` | serialises concurrent ledger inserts per account (lock) | |
| `cfr003_normal_balance_history_guard` (chart_of_accounts, mig. 190) | UPDATE of `normal_balance` once any GL row exists for the account | `account_type`, `code`, `is_active`, COA deletion are not covered here (INF); DELETE of accounts with history relies on FK `general_ledger.account_id -> chart_of_accounts(id)` (RT: NO ACTION) |
| `cfr002_account_scope_write`, `cfr002_restaurant_scope_write` | changing tenant/restaurant of COA/restaurant | |
| `trg_journal_balance` (journal_entry_lines) | balance validation | |
| **None** on `payments_received`, `payments_made`, `invoices_customer/supplier`, `orders` (only `order_stock_transition`), `salaries` (period setter), `fiscal_periods`, `accounting_event_outbox` (RT: trigger list) | | Posted invoices/payments/paid orders/salaries are mutable at DB level by any role with UPDATE/DELETE (`pos_accounting`, `pos_app` in provisioned DB); immutability there is app-logic only (INF) |

Hard-delete exposure: `salaries->tenants`, `payments->tenants`, `recurring_journal_templates->tenants` are `ON DELETE CASCADE`; `journal_entry_lines->journal_entries` CASCADE (guarded by line trigger) (RT `pg_constraint`). `general_ledger->tenants` is NO ACTION, so `DELETE /super-admin/tenants/{id}` (`main.go:1977`) on a tenant with GL history is blocked by FK (INF) - REX to confirm behaviour/soft-delete.

## 10. Transaction atomicity (db.Begin in financial paths)

S: `Begin/BeginTx` sites - `journal_entry_repo.go` (4: create/post/void/import), `invoice_customer_repo.go` (3), `invoice_supplier_repo.go` (4), `payment_acct_repo.go` (2), `accounting_advanced_repo.go` (3), `salary_repo.go` (1, payroll + posting via accounting pool, `main.go:421`), `order_creation_transaction.go`, `order_repo.go:552,614` (payment-status transition + `EnqueuePaidOrderTx` in same tx, `order_repo.go:595`), `inventory_movement_repo.go`, `payment_usecase.go`, `subscription_repo_impl.go` (4). Posting a journal appends GL rows in the same tx (`journal_entry_repo.go:358-425`); deferred triggers verify at commit.
Gaps/Unknowns: (a) outbox state (pos_app pool) and journal (pos_accounting pool) are **separate transactions** (INF). (b) Customer/supplier invoice + journal creation atomicity: uses the shared tx passed to `journal_entry_repo` `...Tx` methods (S) - rollback behaviour on failure REX. (c) Money type in Go paths is `float64` (`journal_entry_repo.go:800-806`, `postingLine`), DB columns NUMERIC(15,2) (RT trigger source) - rounding semantics belong to doc 04.

## 11. Concurrency control and idempotency

| Mechanism | Evidence |
|---|---|
| Row locks `FOR UPDATE` | journal post/void (`journal_entry_repo.go:446`), order payment-status (`order_repo.go:563,624`), document numbers (`financial_document_number.go`), salary, impersonation, kds, driver, theme, product update, subscription, auth_repo (S grep) |
| Advisory locks | none found for finance (grep `pg_advisory` = no matches in repo/handler files listed) |
| Trigger-level locks | `FOR UPDATE` inside cfr002 functions on journal header/account/restaurant rows (RT function source) - serialises ledger inserts per account |
| Outbox claim | `FOR UPDATE SKIP LOCKED` + lease + `lease_generation` fencing (S `accounting_event_repo.go:147-153,190,211`; mig. 183) |
| Unique idempotency keys (RT unique indexes) | `accounting_event_outbox`: `event_id` UNIQUE and `(event_type,tenant,restaurant,order_id)` UNIQUE; `accounting_inventory_movement_events` PK `(tenant,restaurant,order,product)`; `payments`: partial unique `paymob_order_id WHERE status='success'`; `orders.order_number` UNIQUE (global, not per-tenant); `orders.public_token_hash` partial UNIQUE; per-tenant unique numbers on invoices, payments_made/received, journal `entry_number`; `salaries (employee_id,period_start,period_end)` UNIQUE; `fiscal_periods (tenant,restaurant,period_name)` UNIQUE |
| **Missing at DB level** | no unique constraint on `journal_entries (tenant,restaurant,reference_type,reference_id,entry_type)`. Duplicate protection for auto-posted journals is a **check-then-insert** (`auto_posting_repo.go:118-139` `FindJournalEntryByReference`) (S); safety under concurrency depends on the outbox single-claim (INF). No idempotency key on `POST /public/orders`, journal create, payments_received/made create (only `order_cancellation_handler.go:36` requires a UUID `Idempotency-Key`; CORS allows the header, `cors.go:25`) |
| Payment webhook | HMAC check, then "idempotent skip" when payment already success/refunded (`payment_usecase.go` ~L205,209), unique success index, retry queue `webhook_retry_queue` enqueue on 5xx path (`payment_handler.go:~205-220`) (S/RT). This webhook covers **SaaS subscription payments**, not restaurant order payments |
Race outcomes (double post, double void, concurrent redrive, concurrent order payment) are REX.

## 12. Outbox durability, worker retries, recovery, redrive

(S/RT) Table `accounting_event_outbox` (mig. 181): statuses pending/processing/retryable/completed/failed; `max_attempts` 8; backoff `1 min * 2^min(attempts-1,6)` capped (`accounting_event_repo.go:202`); leases fenced by `lease_generation` (mig. 183); expired lease replays unless attempts exhausted, then `failed` (`:136-140`). Backfill of all already-paid orders at migration time (mig. 181:44-48); `ReconcilePaidOrders` re-enqueues paid orders missing events each 15 s (`accounting_event_worker.go:45-60`, interval line 20-22); partial index for recovery scan (mig. 182_index_paid_orders). Redrive of failed events via API 15:WRITE with audit row (mig. 188; `pos_app` I/S only on the audit table). Worker starts whenever `db != nil` (`main.go:513-516`), in-process (no separate service). **Event type constrained to `'paid_order'` only** (mig. 181 CHECK) - refunds, cancellations, voids, payroll, purchases are not outbox-driven (INF): losing such postings leaves no replay record. Dev DB outbox is empty (RT `GROUP BY status` returns no rows). Throughput/stuck-event alarms: UNK. Worker behaviour with multiple API replicas: SKIP LOCKED design says safe (S), REX to prove.

## 13. Rate limiting, CORS/origin, secrets, Redis, WebSocket, backups

**Rate limiting (S).** Auth routes: Redis-backed atomic per-IP window, fail-closed 503 on store error (`rate_limit.go:39-96`; production requires Redis `main.go:549-553`). Limits: register 10, login 20, confirm 20, refresh 20, forgot-password 5, check-subdomain 60, super-admin login 10 per minute. Public orders and webhook use **in-memory** `RateLimitByIP` (per process; no cross-replica limit) (`main.go:598-611,659`). KDS: per-user 60/min, bump 5/s. Client IP derivation honours `X-Forwarded-For` only via `TRUSTED_PROXY_CIDRS` (`super_admin_auth.go:222-300`) (S); prod value UNK. No rate limit on the authenticated accounting/export endpoints (S).

**CORS/origin (S `cors.go`).** Allow-list from `CORS_ORIGINS` plus wildcard single-label subdomains and hard-coded `https://accept.paymob.com`; credentials allowed for matched origins; dev allows any localhost port only when `IsDevelopmentMode()`. Prod/staging lists: `docker-compose.prod.yml:26` (env), `ops/aws/compose.json:150`.

**Secrets handling.** Prod compose uses `${VAR:?}` required vars (DOC `docker-compose.prod.yml`), Redis password via Docker secret file (`docker-entrypoint.sh` builds `REDIS_URL`), host secret files mode 0600 enforced by backup script (`scripts/production/backup-postgres.sh:7-10`). Dev compose embeds `JWT_SECRET: dev-only-jwt-secret-change-me` and postgres/postgres (`docker-compose.yml:57-59,90`) (S). A hard-coded legacy bcrypt hash of the seeded super admin exists in source (`startup_guards.go:11-12`) as a deny-marker. Secret rotation procedure: UNK.

**Redis (S/DOC).** Session store, one-time selection proofs, auth rate-limit counters, cache (`service/cache.go`); prod Redis: AOF everysec, `maxmemory 48mb noeviction`, password required (`ops/aws/compose.json:57-98`). With `noeviction` and 48 MB, a full Redis makes writes fail -> session/rate-limit 503 (INF). Redis state is **not** in the DB backup (RUNBOOK:259).

**WebSocket auth (KDS) (S `internal/handler/ws/kds_ws_handler.go:94-125`).** Cookie-only (never URL tokens), trusted-Origin check, session resolved from the store, then standard `wrapKDSWithPermission` (auth+tenant+status+13:READ) at route (`main.go:1457`); per-room connection limiter. A second WS handler `collaboration_websocket_handler.go:41` has its own CheckOrigin (not reviewed; whether routed: UNK).

**Backups (S/DOC).** Prod script `scripts/production/backup-postgres.sh`: `pg_dump -Fc` piped to GPG AES256 symmetric, listing verification, sha256 sidecar, plus encrypted uploads tarball; `restore-test-postgres.sh` exists. Staging `ops/aws/backup` + deployer takes a pre-migration dump and restores to a disposable DB to verify (RUNBOOK:154-158). Backups live on the same EBS disk; **off-host encrypted copy, retention, uploads/env/TLS/Redis restore are explicitly open production gates** (RUNBOOK:256-259; `SERVIO_PRODUCTION_READINESS.md:46`). No point-in-time recovery/WAL archiving configured in any compose file (S grep: none). RPO/RTO targets: UNK. Compose volumes: `postgres_data`, `redis_data`, `uploads_data` named volumes (`docker-compose.prod.yml:252-261`).

## 14. Security controls: claimed vs implemented

| # | Control | Claimed (source of claim) | Implemented (source) | Status |
|---|---|---|---|---|
| 1 | Row-level tenant isolation | "row-level tenant isolation" (`PROJECT_KNOWLEDGE_BASE.md:58,103`) | RLS on 20 tables, only 13 with `tenant_id`; 57 tenant tables incl. all core financial ones have none; app-level filters | PARTIAL (RT) |
| 2 | Runtime role has no BYPASSRLS | prod guard (`startup_guards.go:60-72`) | enforced for `ENV=production` only; dev = superuser; staging skips guard | PARTIAL (S/RT) |
| 3 | Accounting writer least-privilege | RUNBOOK:221-224 | `pos_accounting` restricted & verified at startup (`accounting_db.go:49-54`); has DELETE on journals/GL (blocked by triggers); dev silently falls back to primary | IMPLEMENTED w/ dev fallback (S/RT) |
| 4 | Journal/GL immutability | CFR-002 closure; mig. 180/190 | triggers enabled, source reviewed; no TRUNCATE guard; owner can disable | IMPLEMENTED, bypassable by owner/superuser (RT/INF) |
| 5 | Closed fiscal period blocks posting | accounting policy (UNK source) | app check only in `PostEntryTx` (`journal_entry_repo.go:386`), fails open if no period row (`:822-825`); no void check; no DB trigger; no API to close periods | PARTIAL (S) |
| 6 | RBAC enforced | `RBAC_ENFORCEMENT` | enforced only when `enforce`/production; dev default audit | IMPLEMENTED in prod config (DOC), not in dev (S) |
| 7 | Separate post/approve privilege | typical accounting practice | none (post=void=create=WRITE) | NOT IMPLEMENTED (S) |
| 8 | CSRF protection | `csrf.go` present | double-submit middleware unused; Origin allow-list on cookie path instead | DIFFERENT MECHANISM (S) |
| 9 | Session revocation | logout revokes (READINESS:45) | cookie sessions revocable; bearer JWT not; no sweep on password/role change | PARTIAL (S/INF) |
| 10 | Login brute-force protection | rate limits | per-IP only, no per-account lockout | PARTIAL (S) |
| 11 | Distributed rate limiting | auth-distributed-controls.md | Redis for auth routes; in-memory for public orders/webhook | PARTIAL (S) |
| 12 | Webhook authenticity | HMAC | constant-time HMAC in live mode; mock accepts all; mock only outside `production` | IMPLEMENTED (S) |
| 13 | Outbox durable accounting | CFR-004 | lease+fencing+retry+reconcile+redrive audit; only `paid_order` events | IMPLEMENTED for paid orders only (S/RT) |
| 14 | Idempotent journaling | CFR | check-then-insert, no unique key | PARTIAL (S/RT) |
| 15 | Audit trail of financial actions | audit tables | auth/permission/export/redrive/HR logged; journal post/void only via columns; no audit trigger on COA/fiscal periods | PARTIAL (RT/S) |
| 16 | Security package (RLS manager, encryption, GDPR, audit_log) | present in repo | not imported by API binary | NOT WIRED (S import graph) |
| 17 | Backups verified & restorable | scripts + runbook | encrypted dump + restore test exist; off-host/PITR not done | PARTIAL (DOC) |
| 18 | WebSocket auth | cookie only | implemented | IMPLEMENTED (S) |
| 19 | Mock payment provider blocked in prod | `paymob_mode.go` | production only; staging allowed | PARTIAL (S) |
| 20 | Prod secrets strength | startup validation | JWT default value rejected only in production; super-admin secret independence for prod+staging | PARTIAL (S) |
| 21 | Auth identity boundary (mig. 191/192) | designed | schema + policies installed, empty, runtime path disabled | NOT ACTIVE (S/RT) |
| 22 | Export gating | EXPORT level | accounting exports 15:EXPORT; order export only READ; WRITE implies EXPORT | PARTIAL (S) |

**Controls reviewed: 22 in the table above (plus 38 sub-items across sections 1-13).**

## 15. Reliability summary

- Health: `/api/v1/health`, `/health/ready` (`health_routes.go:17,33`), metrics behind `METRICS_TOKEN` (`health_routes.go:52`); DB health monitor every 15 s (`main.go:104`).
- Single-process workers (accounting, subscription expiry) with no leader election beyond SKIP LOCKED (S).
- Deployment: signed releases, pre-migration backups, rollback preserving financial guards (RUNBOOK:206-240) (DOC); migrations applied atomically per file; `180/187/188/190` have rollback SQL under `backend/scripts/` (S).
- Not verified here (REX): graceful degradation when Redis/DB unavailable; worker recovery after crash; restore RTO; behaviour under `pos_app`/`pos_accounting` on the dev DB.
