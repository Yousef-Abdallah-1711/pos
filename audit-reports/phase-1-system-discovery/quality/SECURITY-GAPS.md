# Security Gaps (Phase 1, read-only audit)

Scope: authentication, RBAC, tenant isolation, payments, public endpoints, secrets, XSS/SQLi/upload review. Passive code inspection only; nothing was executed. Severity: P0 critical security/integrity, P1 high, P2 medium, P3 low, P4 informational. Evidence is `file:line` relative to the repo root (backend files are under `backend/`). IDs are stable and mirrored in `evidence/issues-security.json`. Financial/HR correctness items (FIN-*) live in `quality/FINANCIAL-AND-HR-CORRECTNESS.md`.

## Summary table

| ID | Sev | Title |
|---|---|---|
| SEC-001 | P0 | Any non-empty X-Impersonation-Token header bypasses authentication; tenant then chosen by request headers |
| SEC-002 | P0 | POST /auth/login/confirm issues a JWT with only email + tenant_id (no password or proof of prior login) |
| SEC-003 | P0 | Backend never enforces module permissions: WithRequiredPermission only annotates context |
| SEC-004 | P0 | Mock Paymob mode is the default; mock VerifyHMAC always true and webhook is public, allowing forged payments to activate subscriptions |
| SEC-005 | P1 | RBAC administration endpoints lack authorization; revoke is cross-tenant |
| SEC-006 | P1 | Unauthenticated order lookup/cancel by sequential id using caller-chosen tenant headers (IDOR, PII exposure) |
| SEC-007 | P1 | Unauthenticated POST /translations/{language}/upload overwrites global in-memory translations; unsynchronized map write |
| SEC-008 | P1 | Subscription/disabled status is not enforced server side on data routes |
| SEC-009 | P1 | No rate limiting or lockout on tenant login, login/confirm, register, forgot-password, token-exchange |
| SEC-010 | P1 | Seed migration creates a platform super admin with a documented password and no IP whitelist |
| SEC-012 | P1 | DB-down stub mode in production: server stays up and login accepts any email with a fixed test password as tenant owner |
| SEC-015 | P1 | Webhook replay/refund-callback can re-activate a refunded subscription; transaction type/flags ignored |
| SEC-026 | P1 | Authenticated tokens lacking restaurant_id can override tenant and restaurant through request headers |
| SEC-011 | P2 | JWT secret handling: insecure default, narrow production check, super-admin secret derived from tenant secret |
| SEC-013 | P2 | Credentials and token fragments written to logs |
| SEC-014 | P2 | No tenant session lifecycle: 24h non-revocable JWT, no refresh/logout, forgot-password is a stub, tokens in localStorage |
| SEC-016 | P2 | Paymob auth token and billing/card metadata stored in plaintext; raw webhook bodies queued |
| SEC-018 | P2 | Tenant-impacting routes with no permission annotation (users, categories, drivers, notifications, themes, contact info, subscription actions) |
| SEC-019 | P2 | Module ID mismatch between routes and migrations (HR/leaves/salaries) |
| SEC-020 | P2 | Unsanitized dangerouslySetInnerHTML for tenant-authored custom section HTML on public sites |
| SEC-017 | P3 | Subscription period/refund consistency issues |
| SEC-021 | P3 | Image upload trusts client Content-Type; no magic-byte/extension validation; serving origin unverified |
| SEC-022 | P3 | CSP permits unsafe-inline/unsafe-eval; HSTS conditional on r.TLS; health/metrics and test route public |
| SEC-023 | P3 | In-memory rate limiter and unbounded client map; mock handler builds callback URL from Host header |

## Positive controls verified
- Paymob HMAC-SHA512 verification: correct field order, constant-time compare, fails closed on empty secret in live mode (`backend/internal/service/paymob_client_live.go:186-204`).
- Super-admin routes uniformly wrapped with rate limit, separate JWT, IP whitelist, active/lock checks (`main.go:1761-1766`).
- Dashboard one-time token: crypto/rand 32 bytes, 30 s TTL, single use (`repository/super_admin_repo.go:624-678`).
- SQL: all values parameterized; dynamic fragments come from whitelisted constants (`handler/http/report_handler.go:2140-2160`, `hr_report_handler.go:60-77`, `payment_acct_repo.go:563-575`). No injectable string-built SQL found.
- Order prices and totals computed server-side (`usecase/order_usecase.go:127-182`); payment amount validated against plan price (`payment_usecase.go:217-232`).
- No real secrets found committed: tracked env files are `*.example` only; `docker-compose.yml:90` and `docker-compose.test.yml:39` carry development/test JWT placeholders; docs contain placeholders (`sk_live_xxxxxxxxxxxxx`). `.env.example` has `STRIPE_SECRET_KEY=sk_test_your_stripe_key` placeholder. Migration 109 and `backend/scripts/seed-test-users.sql:9` document default passwords (SEC-010).
- Uploads: Go's `mime/multipart` strips directory components from filenames, mitigating path traversal in `product_usecase.go:293`.

## Findings

### SEC-001 [P0] Any non-empty X-Impersonation-Token header bypasses authentication; tenant then chosen by request headers

- Category: auth
- Status: OPEN
- Files: `backend/internal/middleware/auth.go`, `backend/internal/middleware/tenant.go`, `backend/cmd/api/main.go`
- Evidence:
  - middleware/auth.go:18-31 accepts header and sets claims{UserID:1} without validation (comment says frontend validated it)
  - middleware/tenant.go:39-63 falls back to X-Tenant-ID/X-Restaurant-ID when claims lack tenant/restaurant
  - main.go:600-628,631-641 wrap all protected routes with AuthMiddleware; middleware.ImpersonationMiddleware (impersonation.go:167) is never mounted
- Root cause: Impersonation was implemented as a frontend-validated flow; backend middleware stub trusts header and legacy header-based tenant selection was kept for backward compatibility.
- Recommended fix: Remove the header shortcut. Validate impersonation tokens server side (hash lookup via ImpersonationRepository.ValidateToken), derive tenant/user from the session, mount ImpersonationMiddleware before TenantContext, restrict header tenant fallback to public routes only, audit-log impersonated mutations.
- Required regression tests:
  - Garbage X-Impersonation-Token + X-Tenant-ID returns 401 on every protected route (route-table test)
  - Valid session token yields only its target tenant/user
  - Expired/ended session rejected
  - Token without restaurant_id cannot select tenant via headers

### SEC-002 [P0] POST /auth/login/confirm issues a JWT with only email + tenant_id (no password or proof of prior login)

- Category: auth
- Status: OPEN
- Files: `backend/internal/usecase/auth_usecase.go`, `backend/internal/handler/http/auth_handler.go`, `backend/cmd/api/main.go`
- Evidence:
  - auth_usecase.go:202-257 looks up user by email and tenant and calls GenerateToken without credential check
  - main.go:444 route is public
  - Login verifies password only against users[0] (auth_usecase.go:131-143)
- Root cause: Two-step multi-tenant login was split without carrying proof of authentication between steps.
- Recommended fix: Make /login return a short-lived signed selection token (or session id) bound to email and the verified user ids; require it in /login/confirm; verify password per selected user; rate limit both endpoints.
- Required regression tests:
  - login/confirm without prior login returns 401
  - selection token for tenant A cannot select tenant B not in its list
  - expired selection token rejected
  - per-tenant password differences handled

### SEC-003 [P0] Backend never enforces module permissions: WithRequiredPermission only annotates context

- Category: rbac
- Status: OPEN
- Files: `backend/internal/middleware/permission.go`, `backend/cmd/api/main.go`, `backend/internal/auth/rbac_manager.go`
- Evidence:
  - permission.go:15-22 sets required_module_id/required_permission_level; grep finds no reader
  - main.go:624 permission check commented out
  - 193 wrapWithPermission usages in main.go all effectively auth-only
  - internal/auth/rbac_manager.go HasPermission not wired
- Root cause: Permission middleware was never completed (Phase 2 use cases commented out) while routes were annotated.
- Recommended fix: Implement PermissionMiddleware reading role_permissions for claims user (cached per tenant/user, invalidate on role change), compare with ordered level map, return 403; fix module IDs (SEC-019); deny by default for unannotated tenant routes.
- Required regression tests:
  - Table-driven: each seeded role x each route expects 200/403 per matrix
  - EMPLOYEE cannot POST /accounting/journal-entries
  - DRIVER cannot GET /hr/salaries
  - EXTERNAL_ACCOUNTANT cannot write
  - Cache invalidation after permission change

### SEC-004 [P0] Mock Paymob mode is the default; mock VerifyHMAC always true and webhook is public, allowing forged payments to activate subscriptions

- Category: payments
- Status: OPEN
- Files: `backend/cmd/api/main.go`, `backend/internal/service/paymob_client_mock.go`, `backend/internal/handler/http/payment_handler.go`, `docker-compose.prod.yml`, `.env.prod.example`
- Evidence:
  - main.go:197-211 any PAYMOB_MODE except 'live' selects mock; 427-431,583-588 register mock endpoints
  - paymob_client_mock.go:51-54 VerifyHMAC returns true
  - docker-compose.prod.yml:74 PAYMOB_MODE defaults to empty; .env.prod.example:24 empty
  - payment_usecase.go:185-232 only checks amount/currency vs public plan price; initiate returns paymob_order_id to the caller
- Root cause: Fail-open default and no production guard; mock client shares production wiring.
- Recommended fix: Fail startup when ENV=production and PAYMOB_MODE!=live or any of API key/HMAC secret/iframe/integration id missing; never register mock routes in production; make mock verify check a configured test secret; default mode must be 'disabled'.
- Required regression tests:
  - Startup test: production + mock aborts
  - Forged webhook in live mode returns 401
  - Mock routes 404 when ENV=production
  - Missing PAYMOB_HMAC_SECRET aborts live startup

### SEC-005 [P1] RBAC administration endpoints lack authorization; revoke is cross-tenant

- Category: rbac
- Status: OPEN
- Files: `backend/internal/handler/http/rbac_permission_handler.go`, `backend/internal/handler/http/rbac_user_role_handler.go`, `backend/internal/handler/http/rbac_role_handler.go`, `backend/internal/repository/role_permission_repository.go`, `backend/cmd/api/main.go`
- Evidence:
  - main.go:925-947 all /rbac/* use wrapAuthOnly
  - rbac_permission_handler.go:31-96 AssignPermission only requires claims; any user may grant ADMIN to any role
  - rbac_user_role_handler.go:34-89 any user may assign any role to any user in tenant including self
  - rbac_permission_handler.go:107-130 + role_permission_repository.go:139-145 Revoke deletes by role_id,module_id with no tenant filter
- Root cause: Authorization was deferred to the (missing) permission middleware.
- Recommended fix: Require ROLES module ADMIN (and prevent self-escalation / last-admin removal); add tenant_id to Revoke query and verify role ownership; audit log changes (permission_audit_logs table exists, migration 104).
- Required regression tests:
  - EMPLOYEE gets 403 on assign/revoke/create role
  - Cross-tenant revoke returns 404 and row survives
  - User cannot grant a level above own
  - Audit log row written

### SEC-006 [P1] Unauthenticated order lookup/cancel by sequential id using caller-chosen tenant headers (IDOR, PII exposure)

- Category: tenant-isolation
- Status: OPEN
- Files: `backend/internal/handler/http/public_order_handler.go`, `backend/cmd/api/main.go`, `backend/internal/middleware/tenant.go`
- Evidence:
  - main.go:472-484 public order routes need only TenantContextMiddleware (headers)
  - public_order_handler.go:116-157 GetOrder returns full order incl. customer fields by numeric id
  - public_order_handler.go:375-430 DELETE cancels any pending order by id without ownership proof
  - order numbers use unix-time modulo (order_usecase.go:63)
- Root cause: Public tracking designed around guessable identifiers without a per-order secret.
- Recommended fix: Issue an unguessable tracking token at order creation; require token (or phone+order number) for lookup/cancel; return a reduced public DTO; resolve tenant from slug/host not headers; rate limit.
- Required regression tests:
  - Order id enumeration without token returns 404
  - Cancel requires token and allowed status
  - Public DTO excludes email/phone/address
  - Rate limit enforced

### SEC-007 [P1] Unauthenticated POST /translations/{language}/upload overwrites global in-memory translations; unsynchronized map write

- Category: integrity
- Status: OPEN
- Files: `backend/internal/handler/http/translation_handler.go`, `backend/cmd/api/main.go`
- Evidence:
  - main.go:491 public route
  - translation_handler.go:170-205 reads upload with no auth, no size limit, writes h.translations[language] without mutex
- Root cause: Admin utility exposed on public router.
- Recommended fix: Remove from public router or require super-admin auth; add mutex and request size limit; validate schema; prefer static files.
- Required regression tests:
  - Unauthenticated upload returns 401
  - Concurrent upload/read race test (-race)
  - Oversize body rejected

### SEC-008 [P1] Subscription/disabled status is not enforced server side on data routes

- Category: billing
- Status: OPEN
- Files: `backend/cmd/api/main.go`, `backend/internal/middleware/tenant_status.go`, `backend/internal/middleware/status_check_config.go`
- Evidence:
  - main.go:610-616 wrapProtectedWithStatus defined; only use at main.go:681 (GET /tenant/status)
  - frontend blocks via TenantStatusProvider only
- Root cause: Status middleware built but not composed into wrapProtected/wrapWithPermission.
- Recommended fix: Compose TenantStatusMiddleware into all tenant wrappers except payment/status/auth paths per DefaultStatusCheckConfig; keep 30 s cache but invalidate on webhook/disable.
- Required regression tests:
  - Expired tenant gets 403 on /products
  - Expired tenant can still call /subscription/initiate and /tenant/status
  - Webhook activation unblocks within cache TTL or immediately

### SEC-009 [P1] No rate limiting or lockout on tenant login, login/confirm, register, forgot-password, token-exchange

- Category: auth
- Status: OPEN
- Files: `backend/cmd/api/main.go`, `backend/internal/middleware/rate_limit.go`
- Evidence:
  - main.go:441-447 plain HandleFunc; RateLimitByIP used only at 464-467 and 575-578
  - auth_usecase.go:40 min password length 6; bcrypt DefaultCost
- Root cause: Limiter exists but not applied to auth routes; limiter is per-process in memory.
- Recommended fix: Apply per-IP and per-email limits with backoff/lockout (shared store when scaled), raise min length to 8-10, consider bcrypt cost 12, uniform error timing.
- Required regression tests:
  - N failed logins trigger 429
  - Limiter keyed by email+IP
  - Register throttled

### SEC-010 [P1] Seed migration creates a platform super admin with a documented password and no IP whitelist

- Category: secrets
- Status: OPEN
- Files: `backend/migrations/109_seed_super_admin.sql`
- Evidence:
  - 109:5 and 109:10 comment states the password and hash; 109:22-28 ip_whitelist NULL means all IPs allowed; email superadmin@pos-saas.com
- Root cause: Convenience seed intended for setup shipped as a migration that runs in all environments.
- Recommended fix: Remove seeded credential; provision via one-time CLI with random password forced to change; require IP whitelist/2FA for super admin; rotate if any environment used it.
- Required regression tests:
  - Fresh DB has no super admin until provisioned
  - Login with documented password fails
  - Super admin without whitelist rejected in production

### SEC-012 [P1] DB-down stub mode in production: server stays up and login accepts any email with a fixed test password as tenant owner

- Category: auth
- Status: OPEN
- Files: `backend/cmd/api/main.go`, `backend/internal/repository/auth_repo.go`
- Evidence:
  - main.go:86-94 continues without DB; main.go:43-71,1061-1068,1476+ stub RBAC/users
  - auth_repo.go:250-268 synthetic owner user tenant 1 id 7 with test password
- Root cause: Development convenience retained in production binary.
- Recommended fix: Fatal on DB failure when ENV=production; compile stub mode behind a build tag or ENV=development guard; remove synthetic user.
- Required regression tests:
  - Production start without DB exits non-zero
  - Stub login unavailable unless ENV=development

### SEC-015 [P1] Webhook replay/refund-callback can re-activate a refunded subscription; transaction type/flags ignored

- Category: payments
- Status: OPEN
- Files: `backend/internal/usecase/payment_usecase.go`
- Evidence:
  - payment_usecase.go:202-205 idempotent skip only for status 'success'
  - 189-252 no use of is_refunded, is_voided, is_capture, has_parent_transaction, type, integration_id
  - updatePaymentStatusTx (750-760) updates by id with no status precondition
- Root cause: State machine lacks terminal-state guard; webhook treated as success when success flag true.
- Recommended fix: Reject payloads with is_refunded/is_voided/has_parent_transaction; verify integration_id and transaction type; make refunded/success terminal (UPDATE ... WHERE status IN allowed); store paymob transaction id UNIQUE and dedupe on it; handle refund callbacks explicitly.
- Required regression tests:
  - Replay of success webhook after refund leaves tenant unchanged
  - Refund callback payload does not activate
  - Concurrent duplicate webhooks activate once
  - Wrong integration_id rejected

### SEC-026 [P1] Authenticated tokens lacking restaurant_id can override tenant and restaurant through request headers

- Category: tenant-isolation
- Status: OPEN
- Files: `backend/internal/middleware/tenant.go`, `backend/internal/pkg/jwt/token.go`
- Evidence:
  - tenant.go:31 requires tenantID>0 && restaurantID>0 from claims, else falls through to headers at 40-63
  - token.go:14 restaurant_id omitempty; users created without restaurant get 0 (main.go:1167-1172)
- Root cause: Backward-compatibility header fallback not restricted to unauthenticated public flows.
- Recommended fix: When claims exist, always use claim tenant; reject header mismatch; resolve default restaurant server-side.
- Required regression tests:
  - Token with tenant A and restaurant 0 plus X-Tenant-ID B returns 403/ignored header

### SEC-011 [P2] JWT secret handling: insecure default, narrow production check, super-admin secret derived from tenant secret

- Category: auth
- Status: OPEN
- Files: `backend/internal/config/config.go`, `backend/internal/config/runtime_env.go`, `backend/cmd/api/main.go`, `docker-compose.yml`
- Evidence:
  - config.go:64 default 'your-secret-key-change-this'
  - config.go:69-73 enforces only when ENV == 'production' (runtime_env.go:20-23 also treats 'prod')
  - main.go:273 super admin secret = JWT_SECRET + '_super_admin'
  - pkg/jwt/token.go: no iss/aud; algorithm family check accepts HS384/512
  - docker-compose.yml:90 dev secret committed (development file)
- Root cause: Single-secret design with defaults.
- Recommended fix: Require JWT_SECRET (>=32 bytes random) and independent SUPER_ADMIN_JWT_SECRET in every non-test env; use IsProductionMode(); pin HS256 explicitly via WithValidMethods; add iss/aud/jti; support key rotation.
- Required regression tests:
  - Startup fails with default/short/equal secrets (prod and staging)
  - HS512-signed token rejected
  - Tenant token never validates on super-admin routes and vice versa

### SEC-013 [P2] Credentials and token fragments written to logs

- Category: logging
- Status: OPEN
- Files: `backend/internal/usecase/auth_usecase.go`, `backend/internal/handler/http/auth_handler.go`, `backend/internal/repository/super_admin_repo.go`
- Evidence:
  - auth_usecase.go:107,133,137 password length, first 30 chars of bcrypt hash, first 3 chars of password
  - auth_usecase.go:161,247,317,337-357 token prefixes and user details
  - super_admin_repo.go:654,675 first 16 chars of one-time token
  - jwt_manager.go logs claims (dead code)
- Root cause: Debug prints left in production paths (fmt.Println).
- Recommended fix: Remove debug prints; structured logger with redaction; never log secrets or token material.
- Required regression tests:
  - Log-capture test asserting no password/hash/token substrings

### SEC-014 [P2] No tenant session lifecycle: 24h non-revocable JWT, no refresh/logout, forgot-password is a stub, tokens in localStorage

- Category: session
- Status: OPEN
- Files: `backend/internal/usecase/auth_usecase.go`, `backend/internal/config/config.go`, `frontend/apps/dashboard/src/components/TenantSelectorModal.tsx`, `frontend/apps/landing/src/stores/authStore.ts`
- Evidence:
  - config.go:65 default expiry 24h
  - auth_usecase.go:288-305 ForgotPassword does nothing but responds 'Password reset email sent' (auth_handler.go:125)
  - internal/auth/jwt_manager.go and internal/handlers/auth_api.go (refresh/logout) not wired
  - TenantSelectorModal.tsx:65, landing authStore.ts:31 localStorage; dashboard authStore.ts:43 non-HttpOnly cookie
  - password change does not invalidate existing tokens (user_settings_handler.go:306-316)
- Root cause: Auth subsystem partially implemented; legacy managers left unwired.
- Recommended fix: Short access token (15m) + rotating refresh (HttpOnly Secure SameSite cookie), server-side revocation list or token version in users table, real reset flow with single-use hashed tokens and email, invalidate on password change.
- Required regression tests:
  - Logout invalidates token
  - Password change revokes old tokens
  - Reset token single use/expiry
  - Refresh rotation reuse detection

### SEC-016 [P2] Paymob auth token and billing/card metadata stored in plaintext; raw webhook bodies queued

- Category: data-protection
- Status: OPEN
- Files: `backend/migrations/115_create_payments_table.sql`, `backend/internal/repository/payment_repo_impl.go`, `backend/migrations/117_create_webhook_retry_queue.sql`
- Evidence:
  - 115:16 paymob_auth_token TEXT; payment_repo_impl.go:39,122
  - retry queue persists raw body incl. source_data.pan masked value and HMAC
- Root cause: Convenience persistence of provider token.
- Recommended fix: Do not persist auth token (short-lived, re-authenticate on demand); define retention for queue payloads; encrypt sensitive columns if kept.
- Required regression tests:
  - Schema test: no auth token column populated
  - Retention job purges processed queue rows

### SEC-018 [P2] Tenant-impacting routes with no permission annotation (users, categories, drivers, notifications, themes, contact info, subscription actions)

- Category: rbac
- Status: OPEN
- Files: `backend/cmd/api/main.go`
- Evidence:
  - main.go:671-675,725-753,1062-1289,1397-1403,1413-1458,1462-1472,645-648 use wrapProtected only
  - POST/DELETE /users lets any user create/delete users; POST /subscription/cancel any user cancels
- Root cause: Coverage gaps in route registration.
- Recommended fix: Assign modules (USERS=6, CATEGORIES=2, THEMES=9, SETTINGS=5, NOTIFICATIONS=4, new BILLING/DRIVERS module) and enforce; owner-only for billing.
- Required regression tests:
  - Route inventory test: every non-public tenant route has a permission annotation

### SEC-019 [P2] Module ID mismatch between routes and migrations (HR/leaves/salaries)

- Category: rbac
- Status: OPEN
- Files: `backend/cmd/api/main.go`, `backend/migrations/100_create_modules_table.sql`, `backend/migrations/105_add_rbac_modules.sql`, `frontend/apps/dashboard/src/lib/rbac.ts`
- Evidence:
  - main.go:757-976 HR routes use module 2 (comment 'HR'); migration 100 insert order gives 2=CATEGORIES, 3=HR; LEAVES=11, SALARIES=12 unused
  - rbac.ts:23 HR=3
- Root cause: Hard-coded numeric IDs duplicated in three places.
- Recommended fix: Reference modules by name constants resolved at startup from the modules table; add contract test.
- Required regression tests:
  - Test that route module ids equal modules.id for the named module
  - Payroll routes require SALARIES grant

### SEC-020 [P2] Unsanitized dangerouslySetInnerHTML for tenant-authored custom section HTML on public sites

- Category: xss
- Status: OPEN
- Files: `frontend/layers/ui-themes/src/sections/custom/CustomSection.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/components/POSReceipt.tsx`
- Evidence:
  - CustomSection.tsx:33 renders html_content/content raw
  - POSReceipt.tsx:12-46 custom DOM sanitizer (blocklist) over markup whose fields are escaped upstream (usePOSPrint.ts:197,244-247)
- Root cause: Page-builder content treated as trusted.
- Recommended fix: Sanitize with DOMPurify (allowlist) on save and render; consider CSP without unsafe-inline for the site; replace blocklist sanitizer with allowlist.
- Required regression tests:
  - Script/onerror payload stripped
  - Snapshot tests for allowed markup

### SEC-017 [P3] Subscription period/refund consistency issues

- Category: billing
- Status: OPEN
- Files: `backend/internal/usecase/payment_usecase.go`
- Evidence:
  - 283 expires_at = paidAt + period (early renewal loses remaining time)
  - 571-615 Paymob refund executed before DB transaction
  - 301-325 history rows written after commit, errors ignored
  - 99 order item name fixed to 'Monthly' for annual plan (paymob_client_live.go:99)
- Root cause: Simplified billing logic.
- Recommended fix: Extend from max(now, current expiry); persist refund intent first (saga/outbox); write history inside the transaction; use plan name in order item.
- Required regression tests:
  - Renew 10 days early extends from expiry
  - Refund DB failure reconciles
  - History row atomic with payment update

### SEC-021 [P3] Image upload trusts client Content-Type; no magic-byte/extension validation; serving origin unverified

- Category: upload
- Status: OPEN
- Files: `backend/internal/usecase/product_usecase.go`, `backend/internal/handler/http/product_handler.go`
- Evidence:
  - product_usecase.go:275-312 allowlist on header Content-Type, filename used in stored path (stdlib mime/multipart applies filepath.Base so traversal is mitigated)
  - no http.DetectContentType, no extension allowlist, no per-tenant quota
  - no FileServer in Go code; /uploads served elsewhere (unverified)
- Root cause: Header-only validation.
- Recommended fix: Sniff content, generate server-side filenames/extensions, re-encode images, serve from separate cookieless origin with nosniff and Content-Disposition.
- Required regression tests:
  - HTML file with image/png header rejected
  - Path-like filename sanitized
  - Size limit enforced

### SEC-022 [P3] CSP permits unsafe-inline/unsafe-eval; HSTS conditional on r.TLS; health/metrics and test route public

- Category: headers
- Status: OPEN
- Files: `backend/internal/middleware/cache.go`, `backend/cmd/api/main.go`
- Evidence:
  - cache.go:228-231 CSP, 238-240 HSTS only if r.TLS != nil
  - main.go:547-560 /health/metrics and /public/test-routing/{slug} unauthenticated
- Root cause: API headers copied from web defaults.
- Recommended fix: Set HSTS at proxy or honor X-Forwarded-Proto from trusted proxy; API CSP default-src 'none'; remove test route; restrict metrics.
- Required regression tests:
  - Header snapshot test
  - Metrics requires auth or internal network

### SEC-023 [P3] In-memory rate limiter and unbounded client map; mock handler builds callback URL from Host header

- Category: availability
- Status: OPEN
- Files: `backend/internal/middleware/rate_limit.go`, `backend/internal/handler/http/mock_paymob_handler.go`
- Evidence:
  - rate_limit.go:253-304 map never pruned, per process
  - mock_paymob_handler.go:243-251 uses r.Host
- Root cause: Simple implementation.
- Recommended fix: Use shared limiter (Redis) with expiry sweep; use configured base URL in mock handler.
- Required regression tests:
  - Limiter memory bounded test
  - Host header spoof does not redirect callback

## Not verified / out of scope
- KDS WebSocket authentication (`internal/handler/ws`), `/uploads` serving, RLS policies (`internal/security/rls_policies.sql` not part of migration list), frontend restaurant-website middleware tenant resolution, mock redirect sanitising, Paymob live payload shapes, dead-code packages `internal/auth`, `internal/handlers`, `internal/handler/*.go` (blockchain, white-label etc.) not routed from main.go.
- No dynamic testing, no dependency vulnerability scan (go.sum/package-lock not scanned).
