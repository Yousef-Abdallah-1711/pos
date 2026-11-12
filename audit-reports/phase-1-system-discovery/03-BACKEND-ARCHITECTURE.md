# 03 - Backend Architecture (Phase 1, read-only)

Scope: `backend/` (Go 1.24 per go.mod, stdlib `net/http` ServeMux with Go 1.22 method patterns). Method: static reading and grep only; Go is not installed, nothing was compiled, run or tested. Every claim cites `file:line`. Where something could not be verified it is labelled so. Companion files: `07-COMPLETE-API-INVENTORY.md`, `evidence/api-inventory.json`, `evidence/issues-backend.json`.

## 1. Executive summary

- 363 route registrations, all in `cmd/api/main.go` (1,894 lines). 294 protected, 45 public, 24 super-admin.
- The request pipeline has **three critical gaps** that the README presents as working: an authentication bypass header (BE-001), RBAC that is declared per route but never evaluated (BE-002), and RBAC-admin endpoints without any caller check (BE-003). Tenant isolation relies solely on `WHERE tenant_id` in application code (no Postgres RLS is applied) and on a middleware that falls back to client headers (BE-005).
- The layering claimed in README.md:96/206 (`Handler -> UseCase -> Repository`) holds for roughly a third of endpoints (BE-012).
- About a third of all non-test Go lines are `//go:build legacy` or unwired (BE-014).
- Positives observed (verified in source): parameterised SQL (`$n` placeholders everywhere sampled), bcrypt for passwords (main.go:1157), JWT algorithm pinned to HMAC (pkg/jwt/token.go:60), constant-time HMAC compare in the live Paymob client (paymob_client_live.go:203), production JWT-secret length validation (config.go:70), most repositories filter by tenant_id and restaurant_id, CORS is allow-list based and production-safe by default (middleware/cors.go:47-76), webhook idempotency and amount/currency checks (payment_usecase.go:202-232).

## 2. Entry point and composition (cmd/api/main.go)

| Lines | What happens |
|---|---|
| 1-68 | In-memory `rbacStorage` seeded in `init()` for stub mode |
| 71-93 | `config.Load`; DB connect failure is only a warning -> `db=nil` stub mode (BE-015) |
| 96-110 | KDS WebSocket hub goroutine |
| 112-345 | Manual constructor wiring of ~45 repositories, ~14 use cases, ~47 handlers |
| 396-430 | Webhook retry + subscription-expiry workers started with `context.Background()` (BE-009) |
| 438-573 | Public routes (auth, menu, homepage, public orders, translations, health, mock paymob) |
| 595-646 | Wrapper closures: `wrapAuthOnly`, `wrapProtected`, `wrapProtectedWithStatus`, `wrapWithPermission`, `wrapKDSWithPermission` |
| 650-1291 | Protected routes, including 40 inline handlers with raw SQL and 25 db==nil stubs |
| 1294-1700 | Orders, tables, KDS, reports, advanced HR/accounting, drivers, notifications, themes, RBAC stubs |
| 1700-1850 | Super-admin routes (`wrapSuperAdmin` closure is defined twice) |
| 1851-1856 | Global middleware: CORS -> Gzip -> SecurityHeaders -> RequestMetrics -> `http.ListenAndServe` (no timeouts) |

There is no dependency-injection container or route table; adding a module means editing this file. `cmd/` also contains `hash`, `migrate`, `migrate-shard`, `seed-themes` (not analysed in depth).

## 3. Real layering vs README

| Metric (from main.go wiring and grep) | Value |
|---|---|
| Endpoints whose handler talks only to a use case | 87 |
| Endpoints using use case + repo(s) | 40 |
| Endpoints whose handler talks directly to a repository | 123 |
| Endpoints whose handler receives `*sql.DB` | 65 |
| Inline closures in main.go (handler = SQL in composition root) | 40 |
| Files in `internal/handler/http` containing direct `Query/Exec/Begin` | 29 (accounting_advanced, admin_order, auth, contract, financial_report, hr_attendance, hr_leave, hr_report, hr_salary, impersonation, insurance, journal_entry, kds, ledger, mock_paymob, notification, payment, payment_monitoring, product, public_menu, rbac_permission, report, subscription, super_admin, table, tax, training, theme_v2, utils) |
| Use cases holding `*sql.DB` | payment_usecase.go, financial_report_usecase.go |
| Use-case files vs repository files | 18 vs ~50 |
| `context`-free DB calls vs `*Context` calls (non-legacy, non-test) | 468 vs 237 |
| Transactions (`Begin`/`BeginTx`) | 30 sites, none in order_repo.go; accounting (journal/invoices/payments) and payments do use them |

Verdict: it is a handler-centric layered monolith with use cases only for the older/core flows (auth, products, orders, tables, KDS, payments, subscriptions, notifications, super-admin) and handler->repo for HR/accounting/RBAC. Business logic in handlers: report_handler.go (2,630 lines, 73 funcs, PDF/Excel generation + SQL), hr_report_handler.go (938, 29 SQL calls), accounting_advanced_handler.go (1,346), main.go inline user/role handlers.

## 4. Middleware and security pipeline (as implemented)

| Wrapper | Chain (outer -> inner) | Notes |
|---|---|---|
| `wrapProtected` (76 routes) | AuthMiddleware -> TenantContextMiddleware | impersonation header bypass (BE-001); header tenant fallback (BE-005) |
| `wrapWithPermission(m,l)` (193) | AuthMiddleware -> TenantContextMiddleware -> WithRequiredPermission | permission never read (BE-002) |
| `wrapKDSWithPermission` (7) | Auth -> Tenant -> WithRequiredPermission(13) -> RateLimitByUser(60/min) -> [RateLimitByUserWindow(5/s)] | same RBAC gap |
| `wrapAuthOnly` (17 routes, RBAC admin) | AuthMiddleware | no admin check (BE-003) |
| `wrapProtectedWithStatus` (1) | Auth -> Tenant -> TenantStatus | only `/tenant/status` (BE-008) |
| `wrapSuperAdmin` (24) | RateLimitSuperAdmin(30/min) -> SuperAdminAuth -> ValidateIPWhitelist -> RequireSuperAdmin | separate JWT secret derived from main secret (BE-026) |
| Global | CORS -> Gzip -> SecurityHeaders -> RequestMetrics | no panic recovery, request id, body limit, timeouts |

Dead middleware: `ImpersonationMiddleware` (impersonation.go:33), `sharding.go`, `cache.go`/`cache_integration.go` helpers (only `SecurityHeadersMiddleware` from cache.go is wired; not exhaustively verified).

## 5. Dead / legacy / duplicated code

| Area | Evidence | Lines |
|---|---|---|
| `//go:build legacy` files | 80 files | 31,571 (of 96,070 non-test) |
| internal/service/legacy | 27 services (blockchain, elasticsearch, CDN, churn prediction ...) | 11,563 |
| internal/handlers (shard/scale/compliance APIs) | all legacy | 3,629 |
| internal/repository/legacy | 14 files | 3,687 |
| internal/models.disabled | 14 files (also tagged legacy) | 2,979 |
| internal/handler/*.go (gin-based: white_label, reporting, blockchain, search ...) | legacy | ~4,900 |
| internal/routes (chi theme routes v2) | legacy | 354 |
| internal/auth (JWT/RBAC/API-key managers, auth middleware) | not referenced by cmd/api; only `tenant_validator.go` used | ~1,900 of 2,182 |
| internal/security, scaling, sharding | not referenced by cmd/api (sharding used by cmd/migrate-shard only) | 2,161 / 2,209 / 1,434 |
| Tests tagged legacy (do not run by default) | admin_theme_handler*_test, employee_handler_test, payment_handler_test, rbac_integration_test, role_handler_test, service/legacy/theme_service_test | 7+1 files |
| Tracked binaries/junk | backend/api.exe, api-local.exe, pos-api.exe, api, app, main, pos-api, pos-backend, tmpclaude-* | git ls-files |

Duplicates: two JWT libraries (v4 live-unused, v5 live), two RBAC implementations, `respondError` x3, `min` x2, two theme stacks (v1 handler+theme_service live, v2 legacy). `go.mod` still requires gin and chi only for legacy code.

## 6. Transactions, context, goroutines, errors

- **Transactions:** present in journal_entry_repo (4), invoice_customer (3), invoice_supplier (4), payment_acct (2), accounting_advanced (2), subscription_repo_impl (4), payment_usecase (3), table_repo, kds_repo, contract_repo, auth_repo. Absent where it matters most: order creation (BE-010). Auto-posting is outside any transaction (BE-011).
- **Context cancellation:** BE-016. Handlers that do pass `r.Context()` (inline main.go handlers, payment, super-admin) are the exception.
- **Goroutines:** KDS hub (`go kdsHub.Run()`, main.go:104), per-connection ReadPump/WritePump (ws/kds_ws_handler.go:148-149, bounded by a per-restaurant connection limiter), retry worker, expiry worker (both `context.Background()`), DB health monitor (main.go:1862+), in-memory cache cleanup goroutines. No unbounded per-request goroutine spawning was found in non-legacy code (grep `go func|go x(`), but workers have no shutdown path (BE-009).
- **Error handling:** decentralised (BE-018). Domain error classification by string matching on `err.Error()` (e.g. public_order_handler.go:78-102 `classifyCreateOrderError`, "not found"/"cannot" checks in handlers) rather than typed errors.
- **Validation:** hand-rolled per handler (BE-019).
- **Auth/validation duplication:** `TenantValidator` (auth/tenant_validator.go) is used by 2 of ~47 handlers; the rest read `middleware.GetTenantID` directly.

## 7. Recommended target architecture (future; no change made now)

1. `cmd/api` = config + wiring only; routes declared in per-module `routes.go` returning `[]Route{Method, Path, Handler, Policy}` so a single table test can assert policy (auth class, module, level) for every route.
2. One security middleware stack: Recover -> RequestID -> Auth(JWT only) -> Tenant(from claims) -> TenantStatus -> `Authorize(module, level)` (enforcing) -> RateLimit(policy) -> BodyLimit -> Handler.
3. Handlers decode/validate/encode; use cases accept `ctx` + a `Principal{UserID, TenantID, RestaurantID, Perms}` and own transactions via a `UnitOfWork`; repositories require tenantID in every method.
4. Typed domain errors mapped once to HTTP; structured logging with redaction.
5. Remove legacy/dead trees; consider Postgres RLS (`SET LOCAL app.tenant_id`) as defence in depth.
6. Public API on slug + unguessable order tokens, no client-chosen tenant headers.

Refactor sequence with lowest risk first: (a) fix BE-001/002/003/007 behind tests, (b) TenantStatus + header fallback removal, (c) transactions for orders, (d) move inline SQL out of main.go, (e) context propagation, (f) split god files.

## 8. Test coverage observed (not executed)

- 62 `_test.go` files under backend/. Active (non-legacy) handler tests exist for admin orders, public orders, tables, KDS (+ws), reports (export/security), ledger validation, financial report cache, and accounting auto-posting. Use-case tests exist for payment, KDS, table, financial report, auto-posting engine, phase-6 accounting cycle. Repo tests: order_repo, general ledger, journal edge cases.
- Not covered by any active test: AuthMiddleware, TenantContextMiddleware, RBAC enforcement, user management, drivers, categories, HR handlers (employee test is legacy-tagged), subscription handlers, payment webhook handler (test is legacy-tagged), RBAC handlers (legacy-tagged), translation upload, impersonation.
- `cmd/api/main_phase6_rbac_routes_test.go:22` asserts route wrapping by regex over main.go source, which passes while enforcement is absent (BE-002).
- `backend/tests/` holds separate unit/integration/load/security suites (tenant_isolation_test.go, security_rls_integration_test.go, etc.); whether they compile or pass was not verified.

## 9. Could not verify

- Any runtime behaviour (no Go toolchain, no DB, no network). All severities are from static reading; BE-001/005/006/007 should be reproduced against a running stack in Phase 2.
- Whether production JWTs ever have `restaurant_id = 0` (BE-005) and whether the frontends rely on `X-Impersonation-Token` or header tenant selection (BE-001/005/006); belongs to the frontend audit.
- Production deployment env (PAYMOB_MODE, TRUSTED_PROXY_CIDRS, reverse-proxy rate limits).
- Exhaustive tenant-scoping of all ~700 SQL call sites (sampled only).
- Whether the `tests/` suites and the 62 test files compile/pass.
- README claim "157+ migrations": `backend/migrations` has 123 files numbered up to 160 plus `999_theme_migration_legacy.sql`; `internal/db/migrations` has 4 more (belongs to the DB audit).

## 10. Findings

Severity scale: P0 = exploitable now / data breach / auth failure; P1 = serious, config- or condition-dependent; P2 = significant quality/reliability risk; P3 = hygiene; P4 = informational.

### BE-001 [P0] Authentication bypass via X-Impersonation-Token header
- **File(s):** `backend/internal/middleware/auth.go`, `backend/internal/middleware/tenant.go`, `backend/internal/middleware/cors.go`
- **Line range:** auth.go:18-31; tenant.go:39-63; cors.go:25
- **Problem:** AuthMiddleware treats any non-empty X-Impersonation-Token header as authenticated: it builds `jwt.Claims{UserID: 1}` without validating the token and calls next. TenantContextMiddleware then accepts X-Tenant-ID / X-Restaurant-ID from the same request. CORS explicitly allow-lists both headers.
- **Evidence:** auth.go:19-31 ("If impersonation token is present, allow the request ... Temporary value"); tenant.go:39-63 header fallback when claims lack tenant+restaurant; cmd/api/main.go:595-646 every wrapper starts with middleware.AuthMiddleware. The real impersonation validator (middleware.ImpersonationMiddleware, impersonation.go:33) is defined but never wired (grep: no caller). Static reading only; not exercised at runtime.
- **Runtime impact:** Unauthenticated caller can reach every handler registered through wrapProtected/wrapWithPermission/KDS/wrapAuthOnly, acting as user id 1 inside an arbitrary tenant/restaurant chosen by headers. Handlers reading middleware.GetTenantID (all inline main.go handlers, HR, accounting, orders, products, drivers, themes) are affected; handlers using auth.TenantValidator (products, notifications) would reject because claims.TenantID is 0, which is the only mitigation observed.
- **Recommended future architecture:** Single AuthMiddleware that only trusts a signed JWT (or a signed, server-validated impersonation JWT whose claims carry the target tenant/user and an `imp` marker). Remove header-based tenant selection for authenticated routes.
- **Refactoring risk:** Medium: the super-admin "view as tenant" frontend flow may rely on the header; must be migrated to a real token exchange (the /auth/token-exchange route already exists) before removal.
- **Required regression tests:** Request with only X-Impersonation-Token + X-Tenant-ID/X-Restaurant-ID returns 401 on every protected route; valid impersonation token yields claims of the target user with an impersonation flag; audit-log entry created.

### BE-002 [P0] RBAC "wrapWithPermission(module, level)" is decorative: permission is never checked
- **File(s):** `backend/internal/middleware/permission.go`, `backend/cmd/api/main.go`, `backend/internal/handler/http/product_handler.go`, `backend/internal/handler/http/notification_handler.go`, `backend/cmd/api/main_phase6_rbac_routes_test.go`
- **Line range:** permission.go:16-24; main.go:621-629,631-640; product_handler.go:33-48; notification_handler.go:32-45
- **Problem:** WithRequiredPermission only stores moduleID and level in the request context. A repo-wide grep finds no reader of RequiredModuleContextKey / RequiredPermissionContextKey outside permission.go. checkProductPermission and checkNotificationPermission contain the comment "If permission middleware is available, check RBAC" followed by no code. 193 wrapWithPermission routes + 7 KDS routes therefore only require a valid login.
- **Evidence:** grep for RequiredModuleContextKey|RequiredPermissionContextKey: only permission.go:11-20 hits; product_handler.go:44-48 empty step 2; internal/auth/rbac_manager.go (HasPermission, line 312) is not referenced from cmd/api. The only RBAC test (cmd/api/main_phase6_rbac_routes_test.go:22) regex-matches main.go text and never issues a request.
- **Runtime impact:** Any authenticated tenant user (cashier, waiter, driver) can read/write HR salaries, journal entries, financial reports, products, orders, KDS settings (module "ADMIN" level) etc. RBAC UI and RBAC tables give a false sense of control. README claims "[RBAC]" in the middleware chain.
- **Recommended future architecture:** A real `RequirePermission(moduleID, level)` middleware that loads effective permissions (user_roles x role_permissions, cached per user/tenant), denies by default, and is covered by a table-driven route test that issues HTTP requests with a low-privilege token.
- **Refactoring risk:** High: turning enforcement on will break any UI flow whose users lack seeded permissions (fallbackPermissionLevelForRole in rbac_permission_handler.go:281 shows ad-hoc defaults). Needs role seed audit + staged rollout (log-only mode first).
- **Required regression tests:** For each of the 200 wrapped routes: token without permission -> 403, with READ only on a WRITE route -> 403, owner -> 200; cross-tenant role must not grant access.

### BE-003 [P0] RBAC administration endpoints have no caller-authorization: any tenant user can grant themselves permissions/roles
- **File(s):** `backend/cmd/api/main.go`, `backend/internal/handler/http/rbac_permission_handler.go`, `backend/internal/handler/http/rbac_user_role_handler.go`, `backend/internal/handler/http/rbac_role_handler.go`
- **Line range:** main.go:924-947; rbac_permission_handler.go:35-92; rbac_user_role_handler.go:34-55; rbac_role_handler.go:144,220
- **Problem:** POST /rbac/permissions/assign, POST /rbac/users/{userId}/roles, POST|PUT|DELETE /rbac/roles, users-by-email variants are registered with wrapAuthOnly. Handlers only verify claims != nil and that the target belongs to claims.TenantID; none checks that the caller is owner/admin.
- **Evidence:** rbac_permission_handler.go:36-91 (claims nil check then rolePermissionRepo.Assign); rbac_user_role_handler.go:35-50 (tenant match then userRoleRepo.Assign). No occurrence of an admin/owner check in rbac_user_role_handler.go (grep "admin|owner" returns nothing); rbac_role_handler.go matches are only IsSystemRole fields.
- **Runtime impact:** Privilege escalation inside a tenant: a low-privilege user can assign an Admin role to self (POST /rbac/users/{self}/roles) or add ADMIN permissions to their role. Combined with BE-002 this is currently masked only because nothing is enforced; fixing BE-002 alone would leave this hole.
- **Recommended future architecture:** Authorization at the use-case layer: RBAC write operations require a dedicated `rbac:ADMIN` capability (or owner flag) and forbid self-elevation / last-owner removal; every mutation audited (permission_audit_log table exists).
- **Refactoring risk:** Low-medium.
- **Required regression tests:** Non-admin token -> 403 on every /rbac write route; admin cannot remove last owner; self-assignment of higher role -> 403; audit row written.

### BE-004 [P1] Sensitive mutations protected only by login (no RBAC wrapper): users, drivers, themes, categories, notifications, tenant contact info
- **File(s):** `backend/cmd/api/main.go`
- **Line range:** main.go:702-756 (contact-info), 1062-1290 (users), 1341-1352 and surrounding driver/theme/category routes (see inventory flag NO-RBAC)
- **Problem:** 22 mutating routes use wrapProtected: POST/PUT/DELETE /users, POST/PUT/DELETE /categories, /admin/drivers*, /admin/themes*, POST /notifications, PUT /tenant/contact-info. POST /users creates a login with a caller-supplied password and DELETE /users/{id} hard-deletes any user of the tenant (repository/user_repo.go:378-383) with no owner/self protection seen.
- **Evidence:** main.go:1133-1204 POST /users inline; 1264-1291 DELETE inline; 07-COMPLETE-API-INVENTORY.md rows flagged NO-RBAC(BE-004).
- **Runtime impact:** Any employee can create accounts, delete the owner, change tenant billing contact e-mail/phone, or restyle the public website.
- **Recommended future architecture:** Every mutation routed through a use case that takes an authorization context; user management moved out of main.go into a UserUseCase.
- **Refactoring risk:** Low.
- **Required regression tests:** Per-route negative tests (non-privileged token -> 403); delete-last-owner and delete-self guards.

### BE-005 [P1] Tenant/restaurant can be taken from request headers even on authenticated routes
- **File(s):** `backend/internal/middleware/tenant.go`, `backend/internal/pkg/jwt/token.go`
- **Line range:** tenant.go:24-63
- **Problem:** TenantContextMiddleware uses JWT claims only when BOTH tenantID>0 and restaurantID>0. Otherwise it falls through to X-Tenant-ID/X-Restaurant-ID. A valid JWT whose RestaurantID is 0 (GenerateToken accepts *int restaurantID, auth_usecase.go ExchangeSuperAdminToken passes user.RestaurantID which can be nil) therefore lets the caller choose any tenant via headers; the JWT tenant is discarded rather than compared.
- **Evidence:** tenant.go:31 (&&), 39-63; README says "Extracts tenant_id + restaurant_id from JWT claims" (README.md TenantContextMiddleware row). Whether production JWTs with restaurant_id=0 exist was not verified (needs DB/data).
- **Runtime impact:** Cross-tenant read/write for any user lacking a restaurant (e.g. owner mid-onboarding, tenant-level admins) and for BE-001 callers.
- **Recommended future architecture:** For authenticated routes: tenant = claims.TenantID, always; restaurant from claims or a validated membership lookup; headers ignored. For public routes use slug -> tenant resolution server-side (public/restaurants/{slug}).
- **Refactoring risk:** Medium: frontends may send headers today; confirm in Phase 1 frontend audit.
- **Required regression tests:** JWT(tenant A, restaurant 0) + X-Tenant-ID: B -> 403/ignored; JWT(A,R1) + X-Restaurant-ID: R2 -> 403.

### BE-006 [P1] Unauthenticated order endpoints trust client-supplied tenant/restaurant headers and expose full orders by sequential id (IDOR) incl. cancel
- **File(s):** `backend/cmd/api/main.go`, `backend/internal/handler/http/public_order_handler.go`, `backend/internal/usecase/order_usecase.go`
- **Line range:** main.go:455-482; public_order_handler.go:37-111,115-157,370-420; order_usecase.go:218-230
- **Problem:** GET /public/orders/{id}, /number/{n}/lookup, /{id}/status, /{id}/track and DELETE /public/orders/{id} are wrapped only by TenantContextMiddleware (headers). No customer secret/token is required; ids are integers. GetOrder returns the whole order object (customer name/phone/address). CreateOrder lets the caller choose order_source; value "pos" auto-confirms the order into the KDS queue. publicOrderHandler holds a restaurantRepo that is never used to verify that the restaurant belongs to the tenant.
- **Evidence:** main.go:469-482; public_order_handler.go:116-157 returns `"data": order`; :385-420 CancelOrder; order_usecase.go:218-229 auto-confirm for pos source.
- **Runtime impact:** PII disclosure and order cancellation/spam against any restaurant whose numeric tenant/restaurant ids are guessed (small integers); kitchen disruption via forged POS tickets; no per-IP rate limit on these routes (only table-info has one).
- **Recommended future architecture:** Public API addressed by restaurant slug (server-resolved tenant); order tracking via unguessable per-order token (UUID) returned at creation; trimmed public DTO; order_source forced to "website"/"qr" for unauthenticated callers; rate limit + captcha for create.
- **Refactoring risk:** Medium: restaurant-website frontend currently sends headers and numeric ids.
- **Required regression tests:** Anonymous GET of another customer order -> 404/403; create with order_source=pos -> forced to website; cancel requires token; rate limit returns 429.

### BE-007 [P1] Paymob webhook signature check is a no-op unless PAYMOB_MODE=live; production compose defaults to mock
- **File(s):** `backend/internal/service/paymob_client_mock.go`, `backend/cmd/api/main.go`, `docker-compose.prod.yml`, `.env.prod.example`, `backend/internal/handler/http/mock_paymob_handler.go`
- **Line range:** paymob_client_mock.go:51-54; main.go:197-210,425-432; docker-compose.prod.yml:74; .env.prod.example:24; mock_paymob_handler.go:243-252
- **Problem:** MockPaymobClient.VerifyHMAC returns true. The mock client is chosen for any PAYMOB_MODE other than "live", including empty; docker-compose.prod.yml passes `${PAYMOB_MODE:-}` and .env.prod.example leaves it blank. config.Load validates JWT secret for production but nothing validates PAYMOB_MODE. The mock routes (/mock/paymob/iframe|process) are registered in the same default state, and ProcessPayment posts to a webhook URL built from the request Host header (SSRF-style, mock_paymob_handler.go:251).
- **Evidence:** See cited lines. HandleWebhook (payment_usecase.go:176-188) relies solely on VerifyHMAC for authenticity; it then checks stored amount == plan price, so an attacker needs a valid paymob_order_id (integer) of a pending payment and no signature.
- **Runtime impact:** If production runs with default env, an unauthenticated POST to /api/v1/webhook/paymob can mark pending subscription payments as successful (free activation). Severity is configuration-dependent; the repo default is the unsafe one.
- **Recommended future architecture:** Mock client allowed only when ENV is development; startup fails in production unless PAYMOB_MODE=live and all Paymob secrets set; webhook handler additionally idempotency-keyed by transaction id.
- **Refactoring risk:** Low.
- **Required regression tests:** Production-mode boot test: PAYMOB_MODE empty -> process exits; webhook with bad hmac -> 401 in every mode other than explicit dev mock.

### BE-008 [P1] Tenant subscription/suspension status is enforced on exactly one route
- **File(s):** `backend/cmd/api/main.go`, `backend/internal/middleware/tenant_status.go`
- **Line range:** main.go:610-618,681; tenant_status.go:36-60
- **Problem:** wrapProtectedWithStatus (TenantStatusMiddleware) is referenced once, for GET /tenant/status. All other 283 tenant routes use wrappers without it, so expired-trial, disabled or unpaid tenants keep full API access. README lists [TenantStatus] in the chain.
- **Evidence:** grep wrapProtectedWithStatus in main.go: lines 610 (definition) and 681 (only use). tenant_status.go itself has a cache and skip-list, i.e. it was designed to be global.
- **Runtime impact:** Revenue leakage and ineffective super-admin disable-tenant action for API clients (frontend may block UI, API does not).
- **Recommended future architecture:** Global tenant-status middleware after Tenant context, with allow-list for billing/auth/health routes.
- **Refactoring risk:** Medium: must keep billing/payment routes reachable for blocked tenants.
- **Required regression tests:** Disabled tenant -> 403 on /products etc.; allowed on /subscription/*; expired trial boundary tests.

### BE-009 [P1] HTTP server has no timeouts, no body-size limits, no graceful shutdown; workers use context.Background()
- **File(s):** `backend/cmd/api/main.go`
- **Line range:** main.go:416,420,1856
- **Problem:** http.ListenAndServe(":"+port, handler) is used directly: no ReadHeaderTimeout/ReadTimeout/WriteTimeout/IdleTimeout, no signal handling or Shutdown. No http.MaxBytesReader anywhere (95 json.NewDecoder/io.ReadAll(r.Body) sites). Retry and expiry workers are started with context.Background() and are never cancelled; kdsHub.Shutdown is deferred but unreachable on SIGTERM because main never returns normally.
- **Evidence:** main.go:1856; grep MaxBytesReader = 0 hits outside legacy; main.go:416,420; grep -c "json.NewDecoder(r.Body)|io.ReadAll(r.Body)" = 95.
- **Runtime impact:** Slowloris / large-body memory exhaustion; in-flight requests and DB transactions are killed on deploy; workers may double-run during rolling restarts.
- **Recommended future architecture:** http.Server with timeouts, signal.NotifyContext, Shutdown(ctx), worker contexts derived from the root context, a body-limit middleware.
- **Refactoring risk:** Low.
- **Required regression tests:** Shutdown test (in-flight request completes); oversized body -> 413; slow-header client closed.

### BE-010 [P1] Order creation is non-atomic and order numbers are globally unique but derived from the clock
- **File(s):** `backend/internal/usecase/order_usecase.go`, `backend/internal/repository/order_repo.go`, `backend/internal/domain/order.go`, `backend/migrations/060_create_orders_table.sql`
- **Line range:** order_usecase.go:64,190-203; order_repo.go:24,630; domain/order.go:240-244; 060_create_orders_table.sql:11
- **Problem:** (a) CreateOrder inserts the order row and then each item with separate autocommit statements (order_repo.go has no Begin; grep shows 0 transactions in it); a failure leaves an order without items and without stock effects. Stock is only checked, never decremented here. (b) order_number = "ORD-YYYY-" + (unix_seconds % 1,000,000 + 1), column is `UNIQUE` across ALL tenants: two orders in the same second (any tenant) collide, the insert fails and the customer sees 500; numbers repeat every ~11.6 days. (c) Tax is hard-coded at 10% (order_usecase.go:173) regardless of restaurant settings.
- **Evidence:** order_usecase.go:64 `GenerateOrderNumber(tenantID, time.Now().Unix()%1000000)`; domain/order.go:242-243 ignores tenantID; migration 060 line 11 `order_number VARCHAR(50) NOT NULL UNIQUE`.
- **Runtime impact:** Intermittent failed checkouts under concurrency, orphaned orders, wrong tax for non-10% restaurants. Not reproduced (static).
- **Recommended future architecture:** Order aggregate persisted in one DB transaction; per-restaurant sequence (DB sequence/counter table) for numbers; tax from restaurant config.
- **Refactoring risk:** Medium: touches core money flow; needs data migration for unique index (tenant_id, restaurant_id, order_number).
- **Required regression tests:** Concurrent create test (N goroutines same second) -> no duplicate/500; failure injection on item insert -> no order row; tax per restaurant.

### BE-011 [P2] Accounting auto-posting runs outside the business transaction and swallows failures
- **File(s):** `backend/internal/usecase/order_usecase.go`, `backend/internal/handler/http/public_order_handler.go`
- **Line range:** order_usecase.go:460-474; public_order_handler.go:69-72
- **Problem:** AutoPostPaidOrder posts the sales entry and the COGS entry as two independent operations after the order is saved; errors are only logged ([AccountingAutoPost] ... failed). No retry, outbox or reconciliation job is registered in main.go.
- **Evidence:** order_usecase.go:468-473.
- **Runtime impact:** Ledger can silently diverge from orders (paid order with no journal entry, or sales without COGS).
- **Recommended future architecture:** Transactional outbox: order paid event row in same tx, worker posts idempotently with retries.
- **Refactoring risk:** Medium.
- **Required regression tests:** Failure injection in PostSalesEntry -> retried and eventually consistent; idempotency on re-post.

### BE-012 [P2] README claims Handler->UseCase->Repository; reality is mixed, with SQL in handlers, usecases and main.go
- **File(s):** `README.md`, `backend/cmd/api/main.go`, `backend/internal/handler/http/*.go`
- **Line range:** README.md:96,206; main.go:702-1291 (inline SQL handlers)
- **Problem:** Counting the 363 registrations by how main.go wires the handler: 87 handler->usecase only, 40 usecase+repo, 123 handler->repository directly, 65 handler receives *sql.DB, 40 inline closures in main.go containing SQL (hr/roles, users, tenant contact-info, tenant status fallback). 29 files in handler/http contain direct db calls; usecase/payment_usecase.go and usecase/financial_report_usecase.go also hold *sql.DB. 18 usecase files exist versus ~50 repository files. Business rules therefore live in handlers (e.g. report_handler.go, hr_report_handler.go) and in main.go.
- **Evidence:** See 07 inventory "Layering" table; grep of handler/http/*.go for .Query|.Exec|.QueryRow gives the 29 files listed in the architecture section.
- **Runtime impact:** Rules cannot be unit-tested or reused; authorization and tenant checks are re-implemented per handler (BE-002/005); behavior drift between stubs and real handlers.
- **Recommended future architecture:** Handlers = decode/validate/encode only; use cases own authorization, transactions and rules; repositories own SQL; main.go only composes. See section 7 of this doc.
- **Refactoring risk:** High effort, incremental: start with users, hr/roles, reports.
- **Required regression tests:** Characterization tests per endpoint before moving code.

### BE-013 [P2] God files
- **File(s):** `backend/internal/handler/http/report_handler.go`, `backend/internal/repository/accounting_advanced_repo.go`, `backend/cmd/api/main.go`, `backend/internal/handler/http/accounting_advanced_handler.go`, `backend/internal/usecase/payment_usecase.go`
- **Line range:** see table
- **Problem:** Non-test Go files over 800 lines (wc -l): report_handler.go 2630; accounting_advanced_repo.go 2133; cmd/api/main.go 1894 (composition root + 40 inline handlers + 15 stub handlers + routing); accounting_advanced_handler.go 1346; payment_usecase.go 1110; handler/white_label_handler.go 1086 (legacy tag); handlers/scaling_api.go 1041 (legacy tag); hr_report_handler.go 938; order_repo.go 904; journal_entry_repo.go 898; kds_repo.go 862; validation/theme_validator.go 835 (no production importer); subscription_repo_impl.go 829.
- **Evidence:** wc -l output on backend/**/*.go excluding _test.
- **Runtime impact:** Review/merge-conflict hot spots; AccountingAdvancedHandler registers 33 endpoints through one struct.
- **Recommended future architecture:** Split by bounded context (reporting, accounting-advanced, payments), one route file per module.
- **Refactoring risk:** Low-medium.
- **Required regression tests:** Existing route-table test generalised.

### BE-014 [P2] ~33% of non-test Go lines are legacy/disabled; auth, RLS, scaling and sharding packages are not wired
- **File(s):** `backend/internal/auth`, `backend/internal/security`, `backend/internal/scaling`, `backend/internal/sharding`, `backend/internal/handlers`, `backend/internal/handler`, `backend/internal/service/legacy`, `backend/internal/repository/legacy`, `backend/internal/models.disabled`, `backend/internal/routes`, `backend/.gitignore`
- **Line range:** n/a
- **Problem:** 80 files carry `//go:build legacy` (31,571 lines of 96,070 non-test lines): internal/handlers/*, internal/routes/*, internal/service/legacy/* (11.5k), internal/repository/legacy/*, internal/models.disabled/* (14 files, 2,979 lines, also tagged), and the top-level internal/handler/*.go (gin-based). Untagged but unreferenced from cmd/api: internal/auth jwt_manager/rbac_manager/api_key_manager/auth_middleware (~1.9k lines; only tenant_validator.go is used), internal/security (RLS, encryption, GDPR; rls_policies.sql is not applied: no ROW LEVEL SECURITY in backend/migrations and no set_config/current_tenant in Go), internal/scaling, internal/sharding, middleware/sharding.go, middleware.ImpersonationMiddleware. go.mod pulls gin, chi and two JWT majors (v4 in dead auth, v5 in live pkg/jwt) mainly for dead code. Compiled binaries are tracked in git (backend/api.exe, api-local.exe, pos-api.exe and ELF files api, app, main, pos-api, pos-backend, plus tmpclaude-* files).
- **Evidence:** grep "go:build legacy" (80 files); grep NewJWTManager/NewRBACManager/NewAPIKeyManager outside internal/auth = only tests/unit; `git ls-files backend` lists binaries.
- **Runtime impact:** Reader confusion (two RBAC stacks, two JWT stacks), supply-chain surface, false confidence that RLS/sharding/RBAC exist. Database-level tenant isolation is absent; isolation is purely `WHERE tenant_id` in app code.
- **Recommended future architecture:** Delete legacy/dead packages after tagging a release; keep one auth package; add Postgres RLS as defense-in-depth if desired.
- **Refactoring risk:** Low (build-tagged code is not compiled by default).
- **Required regression tests:** Build default + `-tags legacy` removed; go vet clean.

### BE-015 [P2] Server boots in "stub mode" when the DB is unreachable and serves fake successful responses; inline handlers dereference nil db
- **File(s):** `backend/cmd/api/main.go`
- **Line range:** main.go:41-68,84-93,1060-1131,1470-1560
- **Problem:** If database.Connect fails the process only logs a warning and continues with db=nil. 25 routes then register in-memory stubs (users list with a hard-coded personal e-mail, user create/update/delete that return success without persisting, notifications, RBAC with a mock Admin role for user id 7). Routes not stubbed (hr/roles inline handlers, tenant/contact-info, many repos) call db methods on a nil *sql.DB and panic per request (net/http recovers the goroutine; no recovery middleware or error response).
- **Evidence:** main.go:89 "Continuing without database (stub mode)"; 1062-1066 hard-coded email; 41-68 rbacStorage seeded in init().
- **Runtime impact:** A transient DB outage at boot in production silently yields a server that reports success for writes it does not persist and grants an in-memory Admin role; load balancer health (/health) still returns HTTP 200 with status "degraded".
- **Recommended future architecture:** Fail fast (log.Fatal) if DB is required; dev stubs behind an explicit DEV_STUB=1; recovery + request-id middleware.
- **Refactoring risk:** Low.
- **Required regression tests:** Boot with bad DSN in production mode -> exit non-zero.

### BE-016 [P2] Database calls mostly ignore request cancellation
- **File(s):** `backend/internal/repository`, `backend/internal/usecase`, `backend/internal/handler/http`
- **Line range:** repo-wide
- **Problem:** Counting call sites outside legacy and tests: 468 context-free db.Query/QueryRow/Exec/Begin/Prepare calls vs 237 *Context variants. Order, accounting, journal, invoice, product, table, notification, KDS and most HR repositories (e.g. order_repo.go, journal_entry_repo.go) and use-case signatures (OrderUseCase.CreateOrder(tenantID, restaurantID, req)) take no context. Client disconnects and timeouts never cancel queries.
- **Evidence:** grep counts; order_usecase.go:53-57 signature.
- **Runtime impact:** Slow queries pile up under load; no per-request deadline; harder graceful shutdown.
- **Recommended future architecture:** context.Context as first parameter through handler->usecase->repo; per-request timeout middleware.
- **Refactoring risk:** Medium (large mechanical change).
- **Required regression tests:** Race/timeout tests with cancelled contexts.

### BE-017 [P2] In-memory rate limiter: unbounded map, per-instance, proxy-IP dependent; almost no routes use it
- **File(s):** `backend/internal/middleware/rate_limit.go`, `backend/internal/middleware/super_admin_auth.go`, `backend/cmd/api/main.go`
- **Line range:** rate_limit.go:46-100; super_admin_auth.go:205-270; main.go:466,577,634-636
- **Problem:** rateLimitWithKey stores per-key counters in a map that is never pruned (entries are replaced only when the same key returns). Limits are per process (no shared store). GetClientIP trusts X-Forwarded-For only from loopback unless TRUSTED_PROXY_CIDRS is set; behind a docker/nginx proxy all clients may share one IP bucket (deployment not verified). Only 4 non-super-admin usages exist (table info, webhook, KDS per-user).
- **Evidence:** cited lines.
- **Runtime impact:** Memory growth under IP churn; either ineffective limits or self-DoS behind a proxy; see BE-020 for missing coverage.
- **Recommended future architecture:** Shared limiter (Redis) with TTL, applied via one policy table per route class.
- **Refactoring risk:** Low.
- **Required regression tests:** Limiter unit tests with clock injection; proxy header tests.

### BE-018 [P2] Error handling and logging are decentralised; internal errors and secrets reach clients/logs
- **File(s):** `backend/internal/handler/http`, `backend/cmd/api/main.go`, `backend/internal/usecase/auth_usecase.go`, `backend/internal/handler/http/auth_handler.go`
- **Line range:** auth_usecase.go:309-380; auth_handler.go:131-166; utils.go; main.go:1140-1290
- **Problem:** Three different response helpers (respondError x735, ErrorResponseWithMessage/UnauthorizedResponse in utils.go, plain http.Error x37) and several envelopes ({"error"}, {"success":false,"error"}, plain text). 121 handler sites return err.Error() to the client (main.go user handlers: "Failed to create user: %v"). 166 fmt.Print* calls remain in non-legacy code; token exchange prints the one-time token prefix, user email, role (auth_usecase.go:316-377; auth_handler.go:148). No request-id, no panic-recovery middleware, no structured logger. `min` and `respondError` are duplicated in several packages.
- **Evidence:** grep counts; auth_usecase.go:321-324,364-380.
- **Runtime impact:** Information disclosure (SQL/driver errors), inconsistent client error handling, secrets in logs.
- **Recommended future architecture:** One error type (domain errors -> HTTP mapping in a single middleware/helper), slog with request-id, redaction.
- **Refactoring risk:** Low-medium (frontend may parse current shapes).
- **Required regression tests:** Contract tests for error envelope.

### BE-019 [P2] No real validation layer; struct validate tags are inert and validation is duplicated
- **File(s):** `backend/internal/domain/order.go`, `backend/go.mod`, `backend/internal/validation/theme_validator.go`
- **Line range:** order.go:107-136; go.mod:32
- **Problem:** Domain structs carry `validate:"required"` tags but go-playground/validator is only an indirect dependency and is not imported by any production file; internal/validation (theme_validator.go, 835 lines) is imported only by a legacy-tagged test. Each handler hand-writes checks (e.g. main.go inline handlers accept empty email/password for POST /users: only JSON decoding is checked, main.go:1133-1171).
- **Evidence:** grep imports; main.go:1146-1171.
- **Runtime impact:** Empty/invalid inputs reach repositories (e.g. users created with empty password hash input, negative quantities handled only in some paths).
- **Recommended future architecture:** DTO validation at the handler boundary via one library, domain invariants in use cases.
- **Refactoring risk:** Low.
- **Required regression tests:** Table-driven invalid payload tests per endpoint.

### BE-020 [P2] No rate limiting on authentication and public mutation endpoints
- **File(s):** `backend/cmd/api/main.go`
- **Line range:** main.go:441-447,470-482
- **Problem:** POST /auth/register, /auth/login, /auth/login/confirm, /auth/token-exchange, /auth/forgot-password and POST /public/orders have no limiter (only RateLimitByIP on table info and the webhook).
- **Evidence:** main.go:441-447 plain mux.HandleFunc.
- **Runtime impact:** Credential stuffing and order spam are unthrottled in the app (a reverse proxy may mitigate; not verified).
- **Recommended future architecture:** Per-IP + per-account limiter on auth, captcha on public create.
- **Refactoring risk:** Low.
- **Required regression tests:** 429 behavior tests.

### BE-021 [P3] Forgot-password is a stub that reports success
- **File(s):** `backend/internal/usecase/auth_usecase.go`, `backend/internal/handler/http/auth_handler.go`
- **Line range:** auth_usecase.go:289-306; auth_handler.go:112-128
- **Problem:** ForgotPassword looks the user up and returns nil with a TODO; the handler answers "Password reset email sent".
- **Evidence:** cited lines.
- **Runtime impact:** Users believe a reset e-mail was sent; no recovery path exists.
- **Recommended future architecture:** Token table + mailer + rate limit.
- **Refactoring risk:** Low.
- **Required regression tests:** E2E reset flow test.

### BE-022 [P1] Unauthenticated endpoint mutates global translations for all tenants
- **File(s):** `backend/internal/handler/http/translation_handler.go`, `backend/cmd/api/main.go`
- **Line range:** translation_handler.go:163-207; main.go:491
- **Problem:** POST /api/v1/translations/{language}/upload has no auth; it parses an uploaded JSON file and replaces h.translations[language] in a process-wide map without a lock or size limit; GetTranslations serves it to every client.
- **Evidence:** translation_handler.go:196-204 `h.translations[language] = translations`; main.go:491 registered with plain HandleFunc.
- **Runtime impact:** Anyone can deface/poison UI text for all tenants (stored XSS if any frontend renders it as HTML), cause data races (map write without mutex) and memory exhaustion.
- **Recommended future architecture:** Remove the route or require super-admin; load translations from files at build time.
- **Refactoring risk:** Low (verify no tool depends on it).
- **Required regression tests:** Unauthenticated POST -> 401; concurrent GET/POST race test.

### BE-023 [P3] KDS WebSocket: origin check disabled, JWT in query string, no permission check
- **File(s):** `backend/internal/handler/ws/kds_ws_handler.go`
- **Line range:** kds_ws_handler.go:83-86,98-106
- **Problem:** CheckOrigin returns true (TODO), token is read from ?token=, only tenant/restaurant presence is checked.
- **Evidence:** cited lines.
- **Runtime impact:** Cross-site WebSocket hijack risk (mitigated by needing a token), token leakage in access logs, any role can watch kitchen orders.
- **Recommended future architecture:** Origin allow-list, ticket-based ws auth, RBAC module 13 check.
- **Refactoring risk:** Low.
- **Required regression tests:** Ws tests with bad origin.

### BE-024 [P3] Some repositories are not tenant-scoped in SQL and rely on handler pre-checks
- **File(s):** `backend/internal/repository/driver_repo.go`, `backend/internal/repository/homepage_section_repository.go`
- **Line range:** driver_repo.go:89,228,249,256,307,314; homepage_section_repository.go:91
- **Problem:** driver GetByID/Delete/UpdateStatus/UpdateAvailability use `WHERE id = $1` only; AdminDriverHandler fetches then compares TenantID/RestaurantID (admin_driver_handler.go:257-264, 354-365) before calling Delete/UpdateStatus. Correct today but fragile (new call sites bypass it). Other repos sampled (orders, employees, journal, KDS, categories) do filter by tenant_id/restaurant_id; a full sweep of every query was not performed.
- **Evidence:** cited lines.
- **Runtime impact:** Latent cross-tenant bug on the next caller.
- **Recommended future architecture:** All repo methods take tenantID and include it in WHERE.
- **Refactoring risk:** Low.
- **Required regression tests:** Cross-tenant repo tests.

### BE-025 [P3] Duplicated auth/JWT/response utilities
- **File(s):** `backend/internal/auth`, `backend/internal/pkg/jwt`, `backend/internal/middleware/super_admin_auth.go`, `backend/internal/handler/http/auth_handler.go`
- **Line range:** n/a
- **Problem:** Two JWT implementations (auth/jwt_manager.go on golang-jwt v4, pkg/jwt/token.go on v5), a third for super-admin tokens, two RBAC models (auth/rbac_manager.go vs repository RBAC), TenantValidator re-checks claims that middleware already set (and logs 5 lines per call), `min()` defined in two packages, respondError defined in handler/http and middleware.
- **Evidence:** cited files.
- **Runtime impact:** Drift risk; log noise (TenantValidator logs on every product/notification request).
- **Recommended future architecture:** Single auth package.
- **Refactoring risk:** Low.
- **Required regression tests:** n/a

### BE-026 [P3] Super-admin token secret is derived from the tenant JWT secret; public impersonation endpoints take tokens in query strings; duplicate disable/enable routes
- **File(s):** `backend/cmd/api/main.go`, `backend/internal/handler/http/impersonation_handler.go`
- **Line range:** main.go:273,1792-1793,1811-1813,1835-1836; impersonation_handler.go:198-240
- **Problem:** Super-admin JWT secret = cfg.JWT.Secret + "_super_admin" (one leaked secret compromises both realms). /impersonation/validate|end|check are unauthenticated and read ?token=. Tenant disable/enable exist twice: PUT via SuperAdminHandler and POST via SubscriptionHandler (different code paths).
- **Evidence:** cited lines.
- **Runtime impact:** Weaker isolation of super-admin realm; token leakage via URLs; inconsistent audit trail.
- **Recommended future architecture:** Independent secret, POST body tokens, single disable/enable use case.
- **Refactoring risk:** Low.
- **Required regression tests:** n/a

### BE-027 [P3] Public diagnostics: health returns raw DB errors, metrics snapshot, debug route
- **File(s):** `backend/cmd/api/main.go`
- **Line range:** main.go:512-556
- **Problem:** /health and /health/ready return `db_status: err.Error()` (driver/DSN-level text); /health/metrics returns observability.Snapshot() unauthenticated; /public/test-routing/{slug} echoes the slug (written with fmt.Sprintf into JSON, so a slug containing a quote breaks the JSON).
- **Evidence:** cited lines.
- **Runtime impact:** Information disclosure; minor.
- **Recommended future architecture:** Return status only; metrics on internal port.
- **Refactoring risk:** Low.
- **Required regression tests:** n/a
