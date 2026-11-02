# Frontend <-> backend contracts - Dashboard section

Method: regex extraction of 363 `METHOD /api/v1/...` patterns from `backend/cmd/api/main.go` (Go 1.22 ServeMux, full paths, no router groups) and of 283 literal API call sites in `frontend/apps/dashboard/src` (axios `api.*`, `apiClient.*`, raw `fetch`, `*Api.*`). Paths with template variables were matched with `${}` -> `{}` wildcards. Machine-readable: `evidence/api-calls-dashboard.json` (file, line, method, path, status, backend line), `evidence/api-inventory.json` (backend). Calls whose URL is a variable (`api.post(url ...)` in `accounting/advanced/page.tsx:169`, `api.get(path ...)` in `accountingAdvancedApi.ts:201`) were not resolved.

## Result
- 283 call sites: **274 match**, **8 orphan**, **1 method mismatch**. All 9 are in `inventory/page.tsx` or `lib/api/themeApi.ts`.
- Base URL: `NEXT_PUBLIC_API_URL` default `http://localhost:8080/api/v1` (`lib/config.ts:15`); `themeApi.ts` re-derives `API_ROOT` itself (`:15-17`).

| # | Call (file:line) | Status | Backend evidence | UI reachable |
|---|---|---|---|---|
| 1 | POST /inventory/adjust (`inventory/page.tsx:107`) | ORPHAN | no /inventory route in main.go | **Yes** (FE-DASH-001) |
| 2 | POST /admin/themes/{id}/export (`themeApi.ts:625`) | ORPHAN | only in legacy-tagged theme_routes_v2.go | No caller |
| 3 | POST /admin/themes/import (`themeApi.ts:641`) | METHOD_MISMATCH | `/admin/themes/{id}` has GET/PUT/DELETE only | No caller |
| 4-7 | GET/POST /admin/themes/{id}/components, PUT/DELETE .../components/{cid}, POST .../components/reorder (`themeApi.ts:679-753`) | ORPHAN | not in main.go (main.go:495 commented) | No caller |
| 8 | GET /public/component-library (`themeApi.ts:768`) | ORPHAN | not in main.go | No caller |

## Backend routes with no dashboard caller (excluding public/super-admin/webhook/auth/health)
`GET /me`, `GET /products/alerts/low-stock`, `/hr/roles` CRUD (4), `GET /hr/employees/{id}`, `/hr/attendance/{id}`, `POST /hr/attendance`, `/hr/leaves/{id}`, `/hr/salaries/{id}`, `/hr/performance/{id}`, `/hr/performance/employee/{id}`, `/hr/resignations/{id}`, `GET /reports/hr/social-insurance`, `GET /audit/export`, `/admin/orders/{stats,export,{id}/history}`, `/admin/drivers/{id}/stats`, `/admin/drivers/{id}/location` (GET/POST), `GET /accounting/customer-invoices/{id}/pdf`, `PUT|DELETE /accounting/recurring-journals/{id}`, four `/accounting/excel/*`, `POST /accounting/tax/periods/{id}/close`, `PUT /users/{id}/settings/language`, `POST /notifications`. (`GET /kds/ws` is used via WebSocket.) Features with backend but no UI: order history/audit, driver tracking, VAT period close, Excel import/export, invoice PDF, low-stock alerts, HR roles.

## Auth / headers contract
- axios `api` adds `Authorization: Bearer <auth-storage token>`, `X-Tenant-ID`, `X-Restaurant-ID` (from user object) except URLs containing `/auth/`; under impersonation adds `X-Impersonation-Token` instead of bearer (`api.ts:149-192`). Backend `TenantContextMiddleware` prefers JWT claims, falls back to those headers (`middleware/tenant.go:17-66`); CORS allows these headers (`cors.go:25`).
- Raw-fetch pages send only `Authorization` (rely on JWT claims) - fine for normal login, broken for impersonation (FE-DASH-007).
- Module ids (frontend `MODULE_IDS`, `lib/rbac.ts:21-42`) vs backend `wrapWithPermission` ids observed: orders 10, KDS 13, reports 8, financial reports 17, products 1 - consistent where checked; HR/accounting ids not individually diffed.
- Authz mismatch: frontend treats drivers/categories/themes/users as module-gated in sidebar; backend uses `wrapProtected` for them (FE-DASH-014).

## Response-shape handling (fragility)
Clients accept several shapes: `data.data ?? data`, `data.products|categories|orders|items` (`posApi.ts:143-163`, `dashboard/page.tsx`, users page `Array.isArray(data)? data : data.data || data.users`). Pagination field names vary (`total`, `page_size`, `pageSize` in theme list `?pageSize=` vs `page_size` elsewhere). Error shape read as `error.response.data.error` (string) or `.error.code/.reason/.message` (tenant 403, `tenantStatusStore.ts:75-79`) - two backend error envelopes in use.

## Findings
FE-DASH-001, 007, 008, 013, 014, 024, 026.
