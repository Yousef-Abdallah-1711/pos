# Restaurant Website Audit (frontend/apps/restaurant-website, port 3003)

Static read-only analysis; IDs refer to `evidence/issues-frontend-others.json`.

## 1. Stack
Next.js ^15.1 App Router (root dir `app/`, no `src/`), React 18, @tanstack/react-query, zustand (persist), axios, next-intl ^4.6 (used only for a 3-key inline message set in `app/[locale]/layout.tsx`; grep shows no `useTranslations` consumer), custom translator `lib/translations.ts` + `i18n/messages/{en,ar}.json` (a second `messages/` dir also exists). Declares `@pos-saas/ui-themes` (workspace) but also imports `@pos-saas/ui` (9 files, undeclared; resolved via tsconfig paths) - FE-OTH-026. Scripts: dev 3003, build, start, lint, jest (test/watch/coverage). 8 test files. `next.config.js` has no `transpilePackages`; `output: standalone` except Windows; images allow unsplash, `cdn.example.com`, localhost (FE-OTH-030).

## 2. Routes (details in pages-others.md)
12 pages, 1 route handler (`/api/metrics`), 2 layouts, error, not-found, 1 middleware. Pages: `/`, `/maintenance`, `/[locale]` {`/`, `menu`, `cart`, `checkout`, `orders`, `orders/[id]`, `orders/track`, `orders/track/[orderNumber]`, `settings`, `table/[tableId]/menu`}.

## 3. Middleware (`middleware.ts`)
Order of operations:
1. Legacy `/table/{id}/{locale}/menu` -> `/{locale}/table/{id}/menu` redirect.
2. `/maintenance*` passes (no headers set).
3. `restaurantSlug = host.split('.')[0]` (raw, includes port for bare hosts).
4. Non-locale first segment -> redirect to `/en/...` (except `_next`, `favicon.ico`, `public`, `api`).
5. Sets request headers `x-restaurant-slug`, `x-restaurant-locale`, and `x-table-id` for table URLs. `Headers.set` overwrites client-supplied values, so direct header injection on locale routes is not possible; the real risks are the unvalidated Host-derived value (FE-OTH-013) and the early-return paths that do not overwrite.
6. For "real subdomain" (`!= localhost, www, ''`) calls `GET {api}/public/tenant-status/{slug}`; `!can_access` -> redirect `/maintenance?reason&name&email&phone` (data from backend response); 404 -> `/maintenance`; fetch error or other status -> allow (fail open, FE-OTH-012). No caching (`lib/checkTenantStatus.ts` with 30 s cache exists but is unused).
Bug: for `localhost:3003` the slug is `localhost:3003`, treated as a real subdomain, tenant-status returns 404, redirect to maintenance (FE-OTH-011). Maintenance page renders query parameters verbatim (FE-OTH-018). Backend tenant-status handler is public and returns tenant name/contact info for any slug (by design, enumerable).

## 4. Tenant resolution at runtime
- Server: pages read `x-restaurant-slug` (`menu/page.tsx:24` defaults to `'demo'`).
- Client: `SubdomainProvider` (`lib/subdomain-context.tsx`) reads cookies `restaurant-slug`/`is-restaurant-site` (never set by middleware - dead) else `parseSubdomain(window.location.host)`, fetches `GET /public/restaurants/{slug}/info`, stores `tenant-id`/`restaurant-id` in localStorage; `lib/api-client.ts` sends them as `X-Tenant-ID`/`X-Restaurant-ID`. Two API client modules exist (FE-OTH-016).

## 5. Cart / checkout / order flow
- Cart: `lib/store/cart-store.ts` zustand persisted as `cart-store` (items, tableId, tableName); `addItem` blocks out-of-stock if `trackInventory`; line totals computed client side; persisted prices can go stale (server recomputes).
- Checkout (`checkout/checkout-page-content.tsx`): zod form, payment fixed `cash` (line 26), delivery fee 5.99 / free over 50 and `$` formatting computed client-side (FE-OTH-017). Inventory pre-check from persisted stock numbers. Submit -> `POST /public/orders` with `{customer_*, delivery_address, payment_method, notes, items:[{product_id, quantity, special_instructions}], table_id?, order_source: 'dine_in'|'website'}`; prices are not sent (good). Success shows order number and links to `/orders/track/{n}`; `clearCart()`.
- Tracking: `GET /public/orders/number/{n}/lookup` then `/public/orders/{id}/track`; polling 10-15 s via react-query. Cancel: `DELETE /public/orders/{id}`.
- Backend exposure: these public endpoints trust client tenant headers and identify orders by sequential id (FE-OTH-002, P0). `/orders` list calls `GET /public/orders`, not registered (FE-OTH-015).

## 6. Table QR menu
`/[locale]/table/[tableId]/menu`: server-side `GET /public/table/{id}/info` (`no-store`), invalid/maintenance statuses show notice (`lib/table-ordering.ts`), cart store gets table context, order source `dine_in`. Table id is global and enumerable; slug from table overrides host (FE-OTH-020). Test files cover table warning and cart table context.

## 7. Theme loading and caching
Server: `lib/api/get-theme-server.ts` fetches `/public/restaurants/{slug}/theme` with `next:{revalidate:3600}`; homepage ISR fetches homepage 300 s, theme 3600 s, products 600 s (`[locale]/page.tsx:149-156`). Pages also call `headers()`, which makes them dynamic regardless. Cache key is the full URL so tenants do not mix; theme edits take up to 1 h to appear. Client: zustand `theme-store.ts` persisted `theme-storage` (memory -> localStorage -> API), `theme-api.ts` calls unregistered `/public/themes*` (FE-OTH-015). Rendering uses `@pos-saas/ui-themes` (`SectionRenderer`, `Header`, `Footer`, `ThemeProvider`) with fallback `getDefaultThemeData`.

## 8. i18n / RTL
Locales en/ar from path. Root layout `<html lang dir>` is always en/ltr (FE-OTH-014); RTL via inner div in `[locale]/layout.tsx`. Translator is custom (`createTranslator`), next-intl effectively unused. Checkout/cart strings go through `t()`; some components keep English fallbacks.

## 9. UI inventory (approx)
~48 button usages, 2 `<form>`, ~11 inputs, 3 select/textarea, 1 table, 0 dialogs, 2 `useForm`; component folders: commerce (cart, order-list, product-card, status tracker), homepage, layout (SSR + theme header/footer), providers, system/MaintenancePage.

## 10. Env
`lib/config/api.ts`: `NEXT_PUBLIC_API_BASE_URL` or `NEXT_PUBLIC_API_URL`; throws in production when missing (`getApiBaseUrl`), dev fallback `http://localhost:8080/api/v1`; `NEXT_PUBLIC_API_TIMEOUT` (default 5000). `lib/api/api-client.ts:11` evaluates `getApiBaseUrl()` at import, so a missing var fails at module load; inside middleware the call is within try/catch so it silently fails open. Note `NEXT_PUBLIC_` values are build-time inlined; middleware (edge) reads them from the build.

## 11. Other
`/api/metrics` POST unauthenticated, in-memory (FE-OTH-019). `@vercel/analytics` and web-vitals included. `app/sitemap.ts` exists with a commented-out restaurants fetch of a non-existent `/public/restaurants`.

## 12. Priority findings
FE-OTH-002 (P0), 011, 012, 014, 015, 017 (P2), 013, 016, 018-020 (P3).
