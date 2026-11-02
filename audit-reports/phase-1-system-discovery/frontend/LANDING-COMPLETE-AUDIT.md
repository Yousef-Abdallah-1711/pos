# Landing App Audit (frontend/apps/landing, port 3001)

Method: static read-only analysis. Findings reference `issues-frontend-others.json` (FE-OTH-nnn). Nothing was executed.

## 1. Stack and entry points
- Next.js ^15.1 App Router, React 18, Tailwind 3, zustand (persist), axios, react-hook-form + zod, sonner, next-themes, lucide (`package.json`).
- Scripts: `dev` (3001), `build`, `start`, `lint`, `test` (vitest), `test:run`, `test:config-resolver`. No workspace deps declared; `next.config.js:6` lists `@pos-saas/ui` and non-existent `@pos-saas/types` in `transpilePackages`; `output: 'standalone'` except on Windows.
- Entry: `src/app/layout.tsx` (html lang="en", ThemeProvider, Inter), `src/app/[locale]/layout.tsx` (LanguageProvider + Toaster), `src/middleware.ts`.

## 2. Routes (see pages-others.md)
Effective pages: `/[locale]`, `/[locale]/auth/login`, `/[locale]/auth/register`, `/[locale]/auth/super-admin-login`. Six more files are redirect stubs.

### Duplicate routes: are /login, /register vs /auth/login, /auth/register redundant?
Yes, and partly broken (FE-OTH-009). `src/middleware.ts:19-35` redirects any path without `en|ar` prefix to `/en<path>`. So `/auth/login` -> `/en/auth/login` (works, stub never runs), but `/login` -> `/en/login` which does not exist (404). `src/app/login`, `src/app/register`, `src/app/auth/*`, `src/app/page.tsx` are dead code. The stubs also hardcode `/en`, ignoring locale. `lib/api.ts:40` and backend-generated links use `/auth/login`, which work only via the middleware.

## 3. Auth mechanism
1. Login: `POST {API}/auth/login` (axios `lib/api.ts`). If `multiple_tenants` -> `TenantSelectorModal` -> `fetch POST /auth/login/confirm` with `{email, tenant_id}` only (modal line 46-60; backend gap = FE-OTH-001, P0).
2. On success: `useAuthStore.setAuth` (zustand persist `auth-storage` + `localStorage.token/user`), then `window.location.href = <dashboard>/auth/callback?token=<JWT>&user=<json>` (login page:76, register:131, modal:74) -> FE-OTH-005.
3. Register: `GET /auth/check-subdomain?slug=` (debounced), `POST /auth/register` {name,email,password,restaurant_name,restaurant_slug,phone}; client password min 6 mirrors backend.
4. Super-admin handoff: `/[locale]/auth/super-admin-login?token=` -> `POST /auth/token-exchange` -> redirect `<dashboard>/en/dashboard?auth_token=..&auth_user=..` (hardcoded `/en`). Backend's `redirect_url` points at the dashboard's own page of the same name (`super_admin_handler.go:313-320`), so this landing copy is redundant (FE-OTH-030).
5. `lib/api.ts` request interceptor skips token for `/auth/` URLs; 401 interceptor clears storage and redirects to `<landing>/auth/login` (not used for auth endpoints).

## 4. API calls vs backend (main.go)
| Call | main.go | Status |
|---|---|---|
| POST /auth/login | :443 | exists |
| POST /auth/login/confirm | :444 | exists (insecure, FE-OTH-001) |
| POST /auth/register | :442 | exists |
| GET /auth/check-subdomain | :447 | exists |
| POST /auth/token-exchange | :445 | exists (one-time token store is in-memory) |
| POST /auth/forgot-password | :446 | exists but unused by UI; backend is a stub (FE-OTH-008) |
No call to `/subscription/plans` (pricing is hardcoded).

## 5. State, i18n/RTL
- State: zustand `authStore.ts`; LanguageContext; no react-query.
- i18n: custom `lib/translations.ts` (437 lines, en+ar), `LanguageProvider` takes locale from route; `setLanguage` swaps first path segment. RTL via `dir` on page wrappers and a `useEffect` that sets `document.documentElement.dir/lang` (client-side, root layout is static `lang="en"`) -> FE-OTH-014. Locale default `en`, no Accept-Language detection. Some strings hardcoded English (e.g. "Pricing" label in Pricing.tsx:19, "Privacy"/"Terms" in Footer).

## 6. UI inventory (approx grep)
2 `<form>`, ~18 button usages, ~10 inputs, 4 `useForm`, 1 modal (TenantSelectorModal, custom div not Radix dialog), 0 tables. Local `components/ui/{button,input,label}.tsx` duplicate `packages/ui`.

## 7. Env dependencies
`src/lib/config.ts`: `NEXT_PUBLIC_API_URL`, `NEXT_PUBLIC_DASHBOARD_URL` throw in production when unset (`NODE_ENV==='production'`); `NEXT_PUBLIC_LANDING_URL` falls back to browser origin; `NEXT_PUBLIC_RESTAURANT_DOMAIN` throws. Because they are `NEXT_PUBLIC_*`, they are inlined at build; `Dockerfile.prod` passes them as build ARGs. `getConfig()` throwing at module load of `lib/api.ts:4` means a missing var crashes any page importing api. Dev fallbacks: localhost 3001/3002/8080/3003. `.env.production.example` documents all four.

## 8. Marketing claims vs implementation (FE-OTH-004, -010)
| Claim | Evidence | Verdict |
|---|---|---|
| 30-day free trial | translations.ts:12,63,78,102,119 (+AR) | Conflicts with DB default 15 (migrations 110/112) and dashboard "15-day"; and new tenants get no trial dates at all (FE-OTH-003) |
| 500 EGP/month | translations.ts:81-102; DB plan monthly 50000 cents (migration 118), env default | Matches, but hardcoded; annual plan 5000 EGP exists and is not shown |
| Mobile Apps (iOS & Android) | translations.ts:52-55, 95 | No mobile project in repo |
| Unlimited restaurants/staff | translations.ts:88-89 | Registration sets max_restaurants=1, max_users=5 (auth_repo.go:62-63); enforcement not found by grep |
| Trusted by 500+ restaurants, +30% revenue testimonials | translations.ts:18,107 | Unverifiable static copy |
| No credit card required / cancel anytime | pricing | Consistent: trial needs no payment; backend CancelSubscription at period end |
| Footer: demo, about, blog, help, docs, api, privacy, terms; #about, #contact | Footer.tsx:35-67, Navbar.tsx:36-37 | All dead (no pages/anchors) |

## 9. Tests/build
Tests: 4 files (login test, TenantSelectorModal, multi-tenant integration, config). `login.test.tsx:274` expects `/auth/forgot-password` which does not match the localized route. No middleware tests.

## 10. Priority findings
FE-OTH-001 (P0), 005 (P1), 004, 008, 010 (P2), 009 (P3).
