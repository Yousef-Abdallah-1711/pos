# API Contracts Used by landing, restaurant-website, super-admin

Source of truth for existence: `backend/cmd/api/main.go` (line numbers below). Frontend prefix is `NEXT_PUBLIC_API_URL` (`/api/v1`) except super-admin which uses `/api/v1/super-admin`. Request/response shapes were taken from frontend call sites and handler signatures only; not executed.

## Landing
| Method path | Caller | main.go | Auth | Contract notes |
|---|---|---|---|---|
| POST /auth/login | `[locale]/auth/login/page.tsx` | 443 | none | body {email,password}; returns {user,token} or {multiple_tenants:true,tenants:[{tenant_id,tenant_name,user_id,roles,default_restaurant_id}]}; tenant_name is placeholder `"Tenant {id}"` (auth_usecase.go TODO) |
| POST /auth/login/confirm | `TenantSelectorModal.tsx` | 444 | none | body {email,tenant_id} only - P0 FE-OTH-001 |
| POST /auth/register | register page | 442 | none | body {name,email,password,restaurant_name,restaurant_slug,phone}; 201 {user,token}; 409 on duplicates (string match on error text, auth_handler.go:30) |
| GET /auth/check-subdomain?slug= | register page | 447 | none | {available,message,slug}; no rate limit seen |
| POST /auth/token-exchange | super-admin-login page | 445 | one-time token | returns {success,data:{user,token}} (frontend reads `data.data \|\| data`) |
| POST /auth/forgot-password | (none) | 446 | none | stub, no mail sent |
Hand-off contract to dashboard: `GET {dashboard}/auth/callback?token&user` (user JSON in URL) and `/en/dashboard?auth_token&auth_user`. Not a backend contract; dashboard owns it.
Not called but advertised: GET /subscription/plans (main.go:574).

## Restaurant website
| Method path | Caller | main.go | Notes |
|---|---|---|---|
| GET /public/tenant-status/{slug} | middleware.ts:68 | 566 (conditional on subscriptionHandler) | response {can_access, block_reason, is_maintenance, tenant_name, contact_email, contact_phone}; frontend reads `json.data \|\| json`; backend wraps in OKResponse |
| GET /public/restaurants/{slug}/info | subdomain-context.tsx:48 | 453 | frontend expects top-level `tenant_id`, `id` (comment in code); response format not verified |
| GET /public/restaurants/{slug}/homepage, /theme, /products | `[locale]/page.tsx` | 460, 461, 451 | ISR fetches |
| GET .../categories, .../categories/{id}/products, .../search?q=, .../products/{id} | use-api-queries.ts | 452, 454, 456, 455 | `?lang=en` always sent |
| GET /public/table/{id}/info | table menu page | 465 | RateLimitByIP(120); returns {success,data:{id,table_number,name,status,restaurant_slug}} |
| POST /public/orders | use-api-queries.ts:240 | 472 | headers X-Tenant-ID, X-Restaurant-ID from localStorage; body {customer_name,customer_email,customer_phone,delivery_address,payment_method:'cash',notes,items[{product_id,quantity,special_instructions}],table_id?,order_source} |
| POST /public/orders/validate | hook defined (usage not found in pages) | 473 | |
| GET /public/orders/number/{n}/lookup | track page | 475 | |
| GET /public/orders/{id}, /status, /track | order pages | 481, 476, 477 | sequential ids |
| DELETE /public/orders/{id} | hook `useCancelOrder` | 484 | no ownership proof - P0 FE-OTH-002 |
| GET /public/orders | `useOrders` (`/orders` page) | NOT REGISTERED | 405/404 - FE-OTH-015 |
| GET /public/themes/{slug}, /public/themes | `lib/api/theme-api.ts` | NOT REGISTERED | FE-OTH-015 |
| POST /api/metrics | web-vitals (same origin Next route) | n/a | |
Header contract: `TenantContextMiddleware` accepts JWT claims first, else X-Tenant-ID/X-Restaurant-ID integers (middleware/tenant.go:39-63). No server-side binding of header ids to the host slug.

## Super admin
All under `/api/v1/super-admin` (see SUPER-ADMIN-COMPLETE-AUDIT.md section 4 for the table). Response envelope: `ApiResponse<T> {success,data,error?}`; frontend unwraps `response.data.data!`. List endpoints return paginated structures normalized in `lib/api.ts` (`listTenants`, `listUsers`, `getAuditLogs`). Contract mismatches: `POST /impersonate/validate` (unused client method) vs real `GET /api/v1/impersonation/validate`; tenant disable/enable exist as PUT (SuperAdminHandler) and POST (SubscriptionHandler), client uses POST; `dashboard-login` returns {token, dashboard_url, redirect_url, expires_in_seconds:30}.
Dashboard-facing public endpoints used by the super-admin handoff: GET /api/v1/impersonation/validate, POST /api/v1/impersonation/end, GET /api/v1/impersonation/check (main.go:1811-1813, unauthenticated by design, token in request).

## External services referenced from these apps
Next/Google Fonts (Inter, Poppins), `images.unsplash.com`, Vercel Analytics/web-vitals. No direct third-party payment calls from these apps; Paymob is backend-only (`/subscription/initiate`, `/webhook/paymob`, mock iframe routes main.go:585-586, dev-only when mock handler non-nil).

## Recommended contract tests
Enumerate frontend endpoint literals and assert each matches a registered mux pattern (would have caught `/public/orders` list and `/public/themes*`); schema tests for tenant-status, info, table info and super-admin login envelopes.
