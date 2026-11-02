# 09 - Authentication and RBAC (Phase 1, read-only audit)

Scope: backend (Go, `backend/`) and frontend (Next.js, `frontend/`). Passive code inspection only. Secret values are never reproduced. Paths below are relative to repo root. Findings use IDs from `quality/SECURITY-GAPS.md` and `evidence/issues-security.json`.

## 1. Authentication architecture (as wired in `backend/cmd/api/main.go`)

| Concern | Implementation | Evidence |
|---|---|---|
| Tenant JWT library | `golang-jwt/jwt/v5`, HS256 signing, `Claims{UserID,Email,TenantID,RestaurantID,Role}` | `backend/internal/pkg/jwt/token.go:10-56` |
| Algorithm pinning | Checks `token.Method.(*jwt.SigningMethodHMAC)` (HS256/384/512 all accepted, `none` rejected) | `token.go:59-64` |
| Expiry | `exp` set from `JWT_EXPIRY` (default `24h`); validated by v5 parser. No `iss`/`aud`/`jti`. | `config.go:65`, `token.go:48-51` |
| Secret handling | `JWT_SECRET` default `"your-secret-key-change-this"`; rejected only when `cfg.Server.Env == "production"` (exact string, `ENV` var) and len < 32. `prod`, `staging` etc. are not covered. | `config/config.go:64,69-73` vs `config/runtime_env.go:20-23` (treats `prod` as production) |
| Password hashing | bcrypt `DefaultCost`; register min length 6, change-password min 8 | `usecase/auth_usecase.go:40`, `user_settings_handler.go:277`, `main.go:1160` |
| Super-admin JWT | Separate `SuperAdminTokenService`, secret = `JWT_SECRET + "_super_admin"` (derived, not independent), expiry `"15m"`, also HS256 + HMAC family check | `main.go:273`, `pkg/jwt/super_admin_token.go:30,58,69-70` |
| Super-admin login | bcrypt compare, failed-attempt counter and lockout (`IsLocked`), IP whitelist (empty = allow all), 30 req/min limiter on protected routes | `usecase/super_admin_auth_usecase.go:87-105`, `middleware/super_admin_auth.go:133-166,168+` |
| Token exchange (super-admin -> tenant dashboard) | 32-byte `crypto/rand` hex token, in-memory map, 30 s TTL, single use; exchange mints a normal tenant JWT for first owner/admin user of the tenant | `repository/super_admin_repo.go:624-678`, `auth_usecase.go:309-389` |
| Impersonation | DB-backed sessions with hashed tokens (`impersonation_repo.go:27-28,100-128`); `ImpersonationMiddleware` exists but is NOT mounted anywhere in `main.go` (grep: only definition). The live path is a bypass (SEC-001). | `middleware/impersonation.go:167`, `middleware/auth.go:18-31` |
| Logout / refresh / revocation | Not implemented for tenant users. `internal/auth/*` (JWTManager with refresh, RBACManager) and `internal/handlers/auth_api.go` (`/auth/refresh`, `/auth/logout`) are not referenced by `main.go` (dead code). Only super-admin has `POST /super-admin/auth/logout`. | `main.go:441-447,1773`; grep `NewJWTManager` outside `internal/auth` = none |
| Password reset | `ForgotPassword` is a stub: looks up user, returns success, handler replies "Password reset email sent". No token, no email. | `auth_usecase.go:288-305`, `auth_handler.go:111-128` |
| Login rate limiting | None on `/auth/login`, `/auth/login/confirm`, `/auth/register`, `/auth/token-exchange`, `/auth/forgot-password`. `RateLimitByIP` is applied only to webhook and public table info. | `main.go:441-447,464-467,575-578` |
| CORS | Allowlist from `CORS_ORIGINS` (+ hard-coded `https://accept.paymob.com`); dev adds localhost. `Allow-Credentials: true` only for allowed origins. Allowed headers include `X-Tenant-ID`, `X-Restaurant-ID`, `X-Impersonation-Token`. | `middleware/cors.go:118-177` |
| CSRF | Bearer token in `Authorization` header, so classic CSRF is low relevance. Dashboard also reads a non-HttpOnly `auth_token` cookie (`frontend/apps/dashboard/src/stores/authStore.ts:43-55`) and stores JWTs in `localStorage` (`TenantSelectorModal.tsx:65`, `landing/src/stores/authStore.ts:31`, `super-admin/src/app/login/page.tsx:36`), so any XSS yields token theft. |
| Security headers | `nosniff`, `X-Frame-Options: DENY`, Referrer-Policy, Permissions-Policy, CSP with `'unsafe-inline' 'unsafe-eval'`; HSTS only when `r.TLS != nil` (absent behind a TLS-terminating proxy). | `middleware/cache.go:216-244` |
| Client IP | `X-Forwarded-For`/`X-Real-IP` trusted only from loopback or `TRUSTED_PROXY_CIDRS`. | `middleware/super_admin_auth.go:205-268` |

### Request pipeline for protected routes
`wrapProtected` = `AuthMiddleware` -> `TenantContextMiddleware` -> handler (`main.go:600-607`). `wrapWithPermission` adds `WithRequiredPermission` as innermost wrapper (`main.go:621-628`). `wrapProtectedWithStatus` (adds subscription check) is used only for `GET /api/v1/tenant/status` (`main.go:610-616,681`).

### Critical auth defects (details in SECURITY-GAPS.md)
1. `AuthMiddleware` accepts any non-empty `X-Impersonation-Token` header and builds claims `{UserID:1}` with no validation (`middleware/auth.go:18-31`). `TenantContextMiddleware` then falls back to caller-supplied `X-Tenant-ID`/`X-Restaurant-ID` headers (`tenant.go:39-63`). Net effect: unauthenticated access to every route wrapped by `wrapProtected`/`wrapWithPermission`/`wrapKDSWithPermission` for any tenant. (SEC-001, P0)
2. `POST /auth/login/confirm` issues a JWT given only `email` + `tenant_id`, no password (`auth_usecase.go:202-257`, route `main.go:444`). (SEC-002, P0)
3. `Login` verifies the password only against `users[0]`, then for single-tenant returns a token for `users[0]`, for multi-tenant returns the tenant list (which LoginConfirm then trusts). (`auth_usecase.go:131-143`)
4. `TenantContextMiddleware` falls back to headers whenever the JWT lacks a non-zero `restaurant_id` (field is `omitempty`), so an authenticated user whose token has no restaurant can pick any tenant/restaurant via headers (`tenant.go:31-63`; claims `token.go:14`). (SEC-026, P1)
5. Stub mode: if DB connection fails the server keeps running (`main.go:86-94`) and `FindUserByEmailAllTenants` returns a synthetic owner for any email with a fixed known test password (`repository/auth_repo.go:250-268`). (SEC-012, P1)

## 2. RBAC data model

Migrations: `100_create_modules_table.sql`, `101_create_permission_levels_table.sql`, `103_create_user_roles_table.sql`, `104_create_role_permissions_table.sql` (`102_...disabled`, `099_...disabled`), `105_add_rbac_modules.sql`, `106_backfill_rbac_defaults.sql`, `120_ensure_reports_module_seed.sql`, `122_add_export_permission_level.sql`, `126/127` (KDS), `132/135` (tables), `155_add_accounting_rbac_modules.sql`, `159_seed_external_accountant_role.sql`.

Permission levels (`domain/role.go:118-126`, migration 122): NONE=0, READ=1, EXPORT=2, WRITE=3, DELETE=4, ADMIN=5 (migration 101 originally had no EXPORT; 122 shifts values).

Modules (IDs follow insert order; confirmed by frontend `lib/rbac.ts:18-39`):
1 PRODUCTS, 2 CATEGORIES, 3 HR, 4 NOTIFICATIONS, 5 SETTINGS, 6 USERS, 7 ROLES, 8 REPORTS, 9 THEMES, 10 ORDERS, 11 LEAVES, 12 SALARIES, 13 KDS, 14 TABLES, 15 ACCOUNTING, 16 SUPPLIERS, 17 FINANCIAL_REPORTS, 18 PERFORMANCE, 19 TRAINING, 20 CONTRACTS. Full matrix in `modules/RBAC.md`.

## 3. Enforcement reality

`WithRequiredPermission` only stores `required_module_id` and `required_permission_level` in the request context (`middleware/permission.go:15-22`). Nothing reads those keys (grep over `backend/` finds only the definition; the real checker is commented out at `main.go:624`). `internal/auth/rbac_manager.go` (which has `HasPermission`) is not wired. Therefore all `wrapWithPermission(..., module, level)` annotations are decorative: any authenticated tenant user (any role, including DRIVER/EMPLOYEE) can call every accounting, HR, payroll, reports, orders, tables and product route. (SEC-003, P0)

Authorization elsewhere:
- RBAC admin endpoints (`/rbac/roles*`, `/rbac/permissions/*`, `/rbac/users/*/roles`) use `wrapAuthOnly` and handlers only check that claims exist; no admin check. Any user can create a role, grant ADMIN on any module, and assign it to themselves (`handler/http/rbac_permission_handler.go:31-96`, `rbac_user_role_handler.go:34-89`). `RevokePermission` deletes by `(role_id, module_id)` with no tenant filter (`rbac_permission_handler.go:107-130`, `repository/role_permission_repository.go:139-145`). (SEC-005, P1)
- `GET /rbac/me/permissions` and `GET /rbac/check-permission` feed the frontend gates.
- Only server-side role logic found: `AdminOrderHandler` uses `userRoleRepo` (driver-specific filtering, `admin_order_handler.go:728-737`).
- `claims.Role` (string on users table: owner/admin/manager/employee...) is carried in the JWT but never consulted for authorization.

## 4. Roles that exist
- System roles per tenant (migrations 021 and 106): ADMIN, MANAGER, EMPLOYEE (021 seeds with Arabic names), DRIVER (106). KDS/tables backfills reference KITCHEN_STAFF, KITCHEN_LEAD, CHEF, KITCHEN_VIEWER, WAITER, HOST, CASHIER, OWNER, ADMINISTRATOR (127, 135, fallback in `rbac_permission_handler.go:283-310`) but these are not seeded as roles by the migrations read.
- EXTERNAL_ACCOUNTANT (migration 159): READ-only accounting; EXPORT on REPORTS (8) and FINANCIAL_REPORTS (17).
- Custom roles via `POST /hr/roles` (inline SQL, `role_code` auto `ROLE_<ms>`) and `POST /rbac/roles`. Note two parallel role systems: `/hr/roles` writes `roles.permissions='{}'` with access flags, `/rbac/*` uses `role_permissions`.
- Platform level: super admin (`users.is_super_admin`, `super_admins` table).
- `users.role` string: `owner`, `super_admin`, others.

## 5. Frontend gates vs backend enforcement
- Dashboard has no Next middleware (`frontend/apps/dashboard` has no `middleware.ts`); gating is client-side only via `checkPermission(MODULE_IDS.X, LEVEL)` in pages and `components/rbac/PermissionBoundary.tsx`, fed by `/rbac/me/permissions` (with a 2.5 s timeout, `lib/rbac.ts:7`).
- Page files with no direct gate call (may be gated in child components; unverified): accounting/advanced; categories (3 pages); drivers (3); hr/contracts, insurance, performance, resignations, training; inventory; products (3); rbac/roles; reports/page and all reports/financial/*; settings/access-control, account, appearance, billing, restaurant, settings root; theme-builder (editor, import, presets, qr, root).
- Because the backend does not enforce, bypassing the UI (curl with any valid token) grants full access regardless of role.
- Module ID mismatch: backend HR routes use module `2` (comment says "HR") but migration 100 assigns 2 = CATEGORIES, 3 = HR (`main.go:757-976`); leaves and salaries use `2` instead of 11/12 (SEC-019).
- Subscription block (`TenantStatusProvider`, `useTenantStatus`) is a client redirect; API stays open for expired/disabled tenants (SEC-008).

## 6. Public (unauthenticated) endpoints worth noting
`/auth/*` (register, login, login/confirm, token-exchange, forgot-password, check-subdomain); `/public/restaurants/{slug}/*` (menu, homepage, theme, settings); `/public/orders*` with tenant headers (create, lookup by id/number, status, track, cancel); `/translations/*` including `POST /translations/{language}/upload` (overwrites global state); `/subscription/plans`; `POST /webhook/paymob`; `/mock/paymob/*` (when mock mode); `/health*` incl. `/health/metrics`; `/public/test-routing/{slug}`; `/public/tenant-status/{slug}`; `/impersonation/validate|end|check`; `/super-admin/auth/login`; `/kds/ws` (WebSocket, token handling in `kds_ws` handler not reviewed).

## 7. Positive observations
- Super-admin routes are wrapped consistently (rate limit, JWT, IP whitelist, active check) `main.go:1761-1766`.
- Payment HMAC verification is constant-time and fails closed on empty secret (live mode) `paymob_client_live.go:186-204`.
- Parameterized SQL everywhere reviewed; dynamic fragments are whitelisted constants (`report_handler.go:2140-2160`, `hr_report_handler.go:60-77`).
- Dashboard login token exchange is single-use with short TTL.
- Order pricing is server-side from product table (`order_usecase.go:127-160`).

## 8. Not verified
KDS WebSocket auth (`handler/ws`), `/uploads` static serving origin, RLS policies (`internal/security/rls_policies.sql` not applied by migrations list), frontend `restaurant-website` middleware tenant resolution, live behaviour (no execution performed), Paymob webhook field coverage versus live payloads.
