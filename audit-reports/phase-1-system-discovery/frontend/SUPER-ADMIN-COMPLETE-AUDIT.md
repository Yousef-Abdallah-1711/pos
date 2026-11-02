# Super Admin App Audit (frontend/apps/super-admin, port 3004)

Static read-only analysis; IDs refer to `evidence/issues-frontend-others.json`.

## 1. Stack
Next.js ^15.1 App Router under `src/app`, all pages are client components. axios, zustand (persist), react-hook-form/zod declared (no `useForm` found in source), sonner, lucide, Radix dialog/dropdown, @tanstack/react-query declared but not used by pages (they call `superAdminAPI` in `useEffect`). Declares `@pos-saas/ui-components` and `@pos-saas/ui-themes` and transpiles them (`next.config.js:21-24`) but imports neither (FE-OTH-026). Scripts: dev 3004, build, start, lint, vitest (1 test: `src/lib/config.test.ts`). No `middleware.ts`.

## 2. Routes (9 pages)
`/` (client redirect), `/login`, `/dashboard`, `/dashboard/tenants`, `/dashboard/tenants/[id]`, `/dashboard/users`, `/dashboard/payments`, `/dashboard/audit-logs`, `/dashboard/settings`; `favicon.ico/route.ts`; layouts root + dashboard. `/dashboard/users/{id}` is linked (`users/page.tsx:130`) but missing (FE-OTH-023).

## 3. Auth (separate JWT, IP whitelist)
- Separate API namespace `/api/v1/super-admin/*` via its own base URL `NEXT_PUBLIC_SUPER_ADMIN_API_URL` (`lib/config.ts`; throws in production if unset; same for `NEXT_PUBLIC_DASHBOARD_URL`).
- Login: `POST /super-admin/auth/login` (public route, `main.go:1770`) -> token stored in `localStorage['super-admin-token']` and in persisted zustand `super-admin-auth-store` (FE-OTH-021). Axios interceptor adds Bearer; 401 clears both keys and hard-redirects `/login`.
- Guard: `SuperAdminAuthProvider` (client) redirects unauthenticated users from `/dashboard*`; no server middleware (FE-OTH-022).
- Logout: nav link removes both keys only; backend `/super-admin/auth/logout` is never called.
- Backend chain (`main.go:1761-1766`, executed in order): RateLimitSuperAdmin(30/min) -> SuperAdminAuthMiddleware(own token service) -> ValidateIPWhitelist -> RequireSuperAdmin. Limiter is an unsynchronized per-route map (FE-OTH-006). Whitelist empty = allow all; `X-Forwarded-For` trusted only from `TRUSTED_PROXY_CIDRS` (FE-OTH-025). Failed-login lockout exists in repo (`super_admin_repo.go:190`, 3 attempts -> 30 min). The login route itself has no rate limit wrapper.

## 4. API calls vs backend
All calls in `lib/api.ts` resolve against `/api/v1/super-admin`:
| Frontend | Backend (main.go) | Status |
|---|---|---|
| POST /auth/login, /auth/logout | 1770, 1773 | exists (logout unused) |
| GET /dashboard/stats | 1774 | exists |
| GET /payment-metrics, /webhook-health | 1776-1777 | exists only if paymentMonitoringHandler != nil |
| POST /payments/{id}/refund | 1782 | exists only if paymentHandler != nil |
| GET /tenants, /tenants/{id}, /tenants/{id}/users | 1788-1790 | exists |
| POST /tenants/{id}/dashboard-login | 1791 | exists |
| GET /tenants/{id}/subscription, GET .../subscription-history | 1833, 1838 | exists |
| POST /tenants/{id}/extend-trial, /disable, /enable | 1834-1836 | exists (POST variants from SubscriptionHandler) |
| DELETE /tenants/{id}?confirm=true | 1837 | exists |
| GET /users, /users/{id} | 1796-1797 | exists |
| POST /impersonate/{userID}, /impersonate/end, GET /impersonate/status | 1800-1802 | exists |
| GET /audit-logs | 1805 | exists |
| PUT /settings/ip-whitelist | 1808 | exists |
| POST /impersonate/validate (api.ts, unused) | not registered (would hit {userID}); real route is public GET /api/v1/impersonation/validate | wrong contract, dead code |
PUT `/tenants/{id}/disable|enable` (SuperAdminHandler) are registered but unused (FE-OTH-024).

## 5. Impersonation flow
`ImpersonationButton` -> `ImpersonationModal`: `POST /impersonate/{userID}` returns token; modal builds `<dashboard>/en/dashboard?impersonation_token=<token>` (token in URL, FE-OTH-005), offers copy/open in new tab; `endImpersonation` -> `POST /impersonate/end`. Dashboard validates through public `/api/v1/impersonation/validate|check|end` (`main.go:1811-1813`). `ImpersonationBanner` component is never rendered; `getImpersonationStatus` has no caller. The alternative "Login to dashboard" button (`tenants/[id]/page.tsx:144`) calls `POST /tenants/{id}/dashboard-login` and opens `redirect_url` (30 s one-time token in an in-memory server store; owner/admin user JWT is minted via `/auth/token-exchange`). Audit logging occurs backend-side (not verified in detail).

## 6. Tenant actions
Modals: Disable (reason), Enable, Extend trial (days, reason), Delete (reason, `?confirm=true`), refund last payment, view subscription history table (`SubscriptionHistoryTable`), `TrialDaysCounter` (progress = remaining/trial_days; with NULL trial_end for new tenants shows no counter, see FE-OTH-003). Actions are gated only by the backend; UI has confirm modals.

## 7. State, i18n
zustand store only; no i18n/RTL (English hardcoded, `<html lang="en">`) - acceptable for an internal tool but note for Arabic-speaking operators. Light design: Tailwind, no dark mode.

## 8. UI inventory (approx)
~48 buttons, 4 `<form>`, 7 inputs, 5 selects/textareas, 5 tables, ~98 dialog/modal references (8 modal components).

## 9. Env
`NEXT_PUBLIC_SUPER_ADMIN_API_URL`, `NEXT_PUBLIC_DASHBOARD_URL` (throw in production when missing). README mentions `NEXT_PUBLIC_API_URL`, which the code does not read.

## 10. Priority findings
FE-OTH-005, 006 (P1), 021 (P2), 022-025 (P3), 026.
