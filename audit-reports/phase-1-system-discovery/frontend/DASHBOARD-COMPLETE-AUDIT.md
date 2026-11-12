# Dashboard app audit (frontend/apps/dashboard) - Phase 1, static analysis

Date 2026-11-01. Method: READ-ONLY source review plus scripted regex scans (route inventory, API call extraction vs `backend/cmd/api/main.go`, import resolution, dead-import heuristic, i18n key diff). Nothing was run: no `tsc`/`next build` (no TypeScript binary in the checkout), no browser, no tests. Counts from grep are labelled approximate. Findings: `evidence/issues-dashboard.json` (FE-DASH-001..027). Per-route detail: `frontend/pages-dashboard.md`, `evidence/routes-dashboard.json`.

## 1. Summary
- Next.js 15 / React 18 app, port 3002, next-intl 4, zustand 4, axios, react-query (declared; usage not mapped), recharts 3, jsPDF/html2canvas, qrcode. 262 files under `src/` (outside `app/`), 87 route files (79 pages + 8 layouts), 0 route handlers, **0 middleware**. 46 test files (vitest; not executed).
- Auth and locale handling are 100% client-side; there is no server-side protection of any route.
- 283 API call sites: 274 match a backend route, 8 orphan + 1 method mismatch (1 reachable from UI: inventory adjust).
- Top risks: POS offline/payment integrity (FE-DASH-002/003), POS using the public order route (024), orphan inventory adjust (001), build gate (009), authorization gaps (014), token handling/logout/impersonation (005-007).

## 2. Architecture
- Providers: `[locale]/layout.tsx` (NextIntl, AuthProvider, next-themes, Toaster, web-vitals), `dashboard/layout.tsx` (`DashboardAuthProvider` -> `DashboardLayout`: sidebar, blocked-tenant overlay, trial banner, language switch).
- Locale: URL prefix en|ar; root `app/page.tsx` redirects to `/en`; `<html lang="en">` fixed; `dir` applied on an inner div; `preferencesStore` can also set documentElement dir.
- Two i18n helpers: next-intl (`useTranslations`, 14 pages) and `lib/translations.ts createTranslator(locale)(key, englishFallback)` (most pages). en.json = ar.json = 1719 keys; 36 literal keys used in code are missing (pos settings 21, bank accounts 8, kds 3, sidebar 3, layout 1).
- HTTP: `lib/api.ts` axios (6s timeout, GET de-dupe 100 ms, tenant/restaurant headers, impersonation header, 401 -> clear localStorage + redirect to landing login), `lib/api/themeApi.ts` second axios, ~21 raw `fetch` calls, `useCache.ts`/`apiCache.ts` generic fetch caches (unreferenced per heuristic).
- Stray dirs: app-root `components/` (14 files), `lib/hooks/use-theme.ts`, `i18n/request.ts`, `tests/` - duplicates/abandoned but inside tsconfig `include`.
- Security headers in `next.config.js`: CSP with unsafe-inline/eval, `X-Frame-Options DENY`, HSTS in prod, but `Cache-Control: public, max-age=3600` for all paths (FE-DASH-020). `output: standalone` except on Windows.

## 3. Module status
| Module | Status | Key evidence |
|---|---|---|
| POS | Partial | cart/hold/barcode/print/offline queue implemented; cash auto-paid, no tender; SW sync broken; public-route order create (FE-DASH-002,003,017,024) |
| KDS | Verified implemented (code) | WS + REST fallback + IndexedDB queue + snapshot; reliability/RBAC gaps (019) |
| Orders | Verified implemented | list/detail/status/cancel/assign; export = current page only; stats/history unused |
| Products & categories | Verified implemented | CRUD with multipart image; categories ungated (014) |
| Inventory | Partial / adjust = Stub | adjust endpoint missing (001) |
| Tables & QR | Verified implemented | CRUD, floor plan, QR; single-domain/guessable QR (025) |
| Drivers | Partial | CRUD + assign; no tracking/stats UI; ungated (014) |
| HR | Verified (core) / Partial (advanced) | 4 core pages + 5 thin advanced pages ungated |
| Accounting | Verified (core) / Partial (advanced) | 70+ endpoints wired; Excel, PDF, period close have no UI; client-computed invoice amounts (018) |
| Reports / financial | Verified implemented | server-computed statements; client PDFs via html2canvas; ungated financial pages |
| RBAC / settings / users | Verified implemented | raw fetch pages; cache not per user (015) |
| Theme builder | Partial | rollback stub (012), orphan APIs (013), ungated (014) |
| Billing | Verified implemented (UI) | Paymob iframe flow via `/subscription/*`; bypasses blocked overlay by design |
| Notifications | Verified implemented | REST + 5 min poll; `NotificationBell` unreferenced |
Detail in `modules/*.md`, `integrations/*.md`.

## 4. TypeScript build failure (reported ~705 errors) - static characterization
`tsc` could not be run here, so **the 705 figure is unverified and the per-category counts below are not tsc counts**. What static reading verifies:
1. **Unresolved imports (3)**: `@/lib/api-client` (`lib/hooks/use-theme.ts:3`; cascades to every consumer such as `components/theme-editor/theme-editor.tsx:19,61` -> `updateTheme`/`isSaving`/`theme` implicit any), `@/hooks/useAuth` (`TenantSelectorModal.tsx:5`), `@/pages/heavy-route` (`usePreloadRoute.ts:5`).
2. **Undeclared/untyped packages**: `@radix-ui/react-tabs` (`ui/tabs.tsx:4`), `qrcode` without `@types/qrcode` (`tables/page.tsx:15`, QRPreview), `@playwright/test` (tests/ is in tsconfig include), `@pos-saas/ui` imported but not a dependency (4 files), `@pos-saas/ui-themes` not in tsconfig paths and not in `transpilePackages` (FE-DASH-011).
3. **Stray-dir model mismatches**: `components/component-manager/section-editor.tsx:194-195` reads `formData.background_image`, absent from `ThemeSection` (`component-sidebar.tsx:7-19`).
4. **recharts v3 typings**: 15 `formatter={...}` props (5 files import recharts) annotate `(value: number)` (e.g. `reports/sales/page.tsx:320`) while v3 passes `ValueType | undefined`.
5. **Loose typing**: 318 `: any`/`as any` occurrences; HR/accounting pages with `[key: string]: any` interfaces and untyped callbacks (employees page) are the likely noImplicitAny sources (`strict: true`, `tsconfig.json`).
6. **Scope inflation**: `include: ["**/*.ts","**/*.tsx"]` drags in app-root dirs, tests and dead files (39 unreferenced, heuristic).
Gate bypass: `next.config.js:11` ignoreBuildErrors only when `NEXT_IGNORE_BUILD_ERRORS=true` (Docker prod builds, per comment). 0 `@ts-ignore`/`@ts-expect-error` in code.
Recommended triage order: (a) fix 3 unresolved imports and delete/relocate stray dirs, (b) add missing deps/types, (c) align workspace packages, (d) fix recharts formatter signatures, (e) type HR/accounting models, then run `tsc --noEmit` to get a real histogram.

## 5. Stores, auth flow, error handling
- Stores: `authStore` (persist `auth-storage`: user+token in localStorage), `preferencesStore` (`preferences-storage`), `tenantStatusStore` (memory), `posStore` (`pos-storage`, cart + held orders + customer PII), `kdsStore` (memory), `themeBuilderStore` (`theme-builder-storage`). None of them is cleared on logout.
- Login happens in the landing app; dashboard receives credentials via query params (`auth_token`+`auth_user`, `/auth/callback?token&user`, super-admin one-time `token` -> POST `/auth/token-exchange`, `impersonation_token` -> GET `/impersonation/validate`). `DashboardAuthProvider` waits fixed 100/200 ms timeouts for persistence (`:116-130`).
- 401: axios interceptor clears localStorage keys and redirects to `${landingUrl}/auth/login`; zustand in-memory state is not reset; skipped while impersonating. Raw fetch/themeApi bypass it.
- Impersonation: dummy bearer `'impersonated'` + sessionStorage `impersonation_token` -> only `api` client honours it (FE-DASH-007). Tenant blocked state: `tenantStatusStore.checkStatus()` once after auth; 403 sets `canAccess=false`; other errors fail open.
- Error normalization: none; each page extracts `error.response?.data?.error` or shows `err.message`. Duplicate fetch patterns: `useHasPermission` per-hook request, `fetchData` in many pages with manual loading state instead of react-query (declared dependency, not used in pages scanned), `useCache.ts` unused; dashboard overview fires ~9 parallel list calls just to compute counts (`dashboard/page.tsx`).

## 6. Positive observations (verified)
POS order payload excludes client prices; KDS reconnect/backoff and Strict-Mode handling; permission check fails closed (`hasPermission=false`) on error when no cache; i18n key parity en/ar; no `dangerouslySetInnerHTML` was searched (unverified); Next security headers present; tenant blocked overlay with billing bypass.

## 7. Unverified / needs runtime or backend confirmation
tsc error count and categories; whether `useHasPermission` pages render before 403 data loads; RTL behaviour of Radix portals; SW registration in production build; backend enforcement for POS public route, users CRUD, theme routes; react-query usage map; `dangerouslySetInnerHTML`/XSS review of theme builder; HR statutory math; visual regression suite (`tests/visual-regression`) and 46 unit tests not run.

## 8. Finding index
FE-DASH-001 P1 inventory adjust orphan | 002 P1 POS replay dupes/unpaid | 003 P1 pos-sw sync wrong origin | 009 P1 build/typecheck | 024 P1 POS public route trusts tenant headers | 004,005,006,007,008,010,011,012,014,015,016,017,018,019,025 P2 | 013,020,021,022,023,026 P3 (and 027 P2). No P0 verified from the dashboard source alone.
