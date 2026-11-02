# Dashboard pages inventory (frontend/apps/dashboard, port 3002)

## Dashboard

Source: static scan of `frontend/apps/dashboard/src/app` on 2026-11-01. Counts are grep-based and APPROXIMATE (JSX tags in the page file only, not child components). Machine-readable copy: `evidence/routes-dashboard.json`.

- Route files: **79 page.tsx + 8 layout.tsx = 87**; 0 route.ts (no Next API routes); 0 middleware (no `src/middleware.ts`, no root `middleware.ts`) so there is no server-side auth or locale routing.
- URL prefix `/{locale}` where locale in en|ar (`generateStaticParams`, `[locale]/layout.tsx:28`); unknown locale silently falls back to `en`.
- All `/{locale}/dashboard/*` pages are wrapped client-side by `DashboardAuthProvider` + `DashboardLayout` (`dashboard/layout.tsx`); gating below is page-level only.
- Stray app-root dirs outside `src/`: `components/` (14 files: component-manager, theme-editor), `lib/hooks/use-theme.ts` (imports unresolved `@/lib/api-client`), `i18n/request.ts` (re-export), `tests/visual-regression/` (Playwright). Not routes; included in tsconfig `include` (see FE-DASH-009/010).

### Layouts

| File | Purpose |
|---|---|
| `src/app/layout.tsx` | Root html (lang hard-coded en), Vercel analytics toggle |
| `src/app/[locale]/layout.tsx` | NextIntlClientProvider, AuthProvider, ThemeProvider, Toaster, dir/lang on inner div |
| `src/app/[locale]/dashboard/layout.tsx` | DashboardAuthProvider + DashboardLayout shell (sidebar, tenant-blocked overlay) |
| `src/app/[locale]/dashboard/kds/layout.tsx` | Viewport only (no zoom) |
| `src/app/[locale]/dashboard/pos/layout.tsx` | Viewport only (no zoom) |
| `src/app/[locale]/dashboard/settings/layout.tsx` | Settings sidebar, custom i18n, RTL border swap |
| `src/app/[locale]/dashboard/theme-builder/layout.tsx` | Pass-through |
| `src/app/[locale]/dashboard/theme-builder/editor/layout.tsx` | Editor layout |

### Pages

Legend: Gate = `MODULE_IDS.X` referenced in page file (none = no page-level gate); i18n: NI=next-intl `useTranslations`, CT=custom `createTranslator` (lib/translations.ts); RTL = count of isRTL/dir/rtl refs; F/B/T/D/I = forms/buttons/tables/dialogs/inputs (approx); BL = contains client-side arithmetic (`reduce`, `toFixed`, `*100`, TAX_RATE); raw = count of raw `fetch(` calls.

| Route | Module | Gate | Direct API calls | Data hooks / api modules | i18n | RTL | F/B/T/D/I | BL | raw |
|---|---|---|---|---|---|---|---|---|---|
| `/auth/callback` | auth | none | - | - | - | 0 | 0/0/0/0/0 | - | 0 |
| `/page.tsx` | shell | none | - | - | - | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}` | shell | none | - | - | - | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/auth/callback` | auth | none | - | - | - | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/auth/super-admin-login` | auth | none | POST /auth/token-exchange | - | - | 0 | 0/1/0/0/0 | - | 0 |
| `/{locale}/dashboard` | overview | CATEGORIES,HR,KDS,ORDERS,PRODUCTS,REPORTS,TABLES (READ,WRITE) | GET /admin/drivers<br>GET /admin/orders<br>GET /categories<br>GET /hr/employees<br>GET /notifications/stats<br>GET /products<br>GET /tables/stats<br>GET /tenant/status | useRbac | CT | 0 | 0/6/0/0/0 | Y | 0 |
| `/{locale}/dashboard/accounting/advanced` | accounting | none | - | accountingAdvancedApi | - | 0 | 0/11/0/0/21 | - | 0 |
| `/{locale}/dashboard/accounting/bank-accounts` | accounting | ACCOUNTING (DELETE,READ,WRITE) | - | accountingApi, bankApi, useRbac | CT | 0 | 0/5/0/1/0 | - | 0 |
| `/{locale}/dashboard/accounting/chart-of-accounts` | accounting | ACCOUNTING (DELETE,READ,WRITE) | - | accountingApi, useRbac | CT | 0 | 0/2/0/1/0 | - | 0 |
| `/{locale}/dashboard/accounting/customer-invoices` | accounting | ACCOUNTING (READ,WRITE) | - | accountingApi, bankApi, customersApi, useRbac | CT | 0 | 0/3/1/2/0 | Y | 0 |
| `/{locale}/dashboard/accounting/customers` | accounting | ACCOUNTING (DELETE,READ,WRITE) | - | accountingApi, customersApi, useRbac | CT | 0 | 0/3/1/2/0 | - | 0 |
| `/{locale}/dashboard/accounting/general-ledger` | accounting | ACCOUNTING (READ) | - | accountingApi, useRbac | CT | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/dashboard/accounting/journal-entries` | accounting | ACCOUNTING (READ,WRITE) | - | accountingApi, useRbac | CT | 0 | 0/2/0/1/0 | - | 0 |
| `/{locale}/dashboard/accounting/journal-entries/[id]` | accounting | ACCOUNTING (READ,WRITE) | - | accountingApi, useRbac | CT | 0 | 0/3/1/1/0 | - | 0 |
| `/{locale}/dashboard/accounting/supplier-invoices` | accounting | ACCOUNTING (READ,WRITE) | - | accountingApi, bankApi, suppliersApi, useRbac | CT | 0 | 0/3/1/2/0 | Y | 0 |
| `/{locale}/dashboard/accounting/suppliers` | accounting | SUPPLIERS (DELETE,READ,WRITE) | - | accountingApi, suppliersApi, useRbac | CT | 0 | 0/4/0/1/0 | - | 0 |
| `/{locale}/dashboard/categories` | categories | none | DELETE /categories/{}<br>GET /categories | - | CT | 0 | 0/3/1/0/0 | - | 0 |
| `/{locale}/dashboard/categories/[id]/edit` | categories | none | GET /categories/{}<br>PUT /categories/{} | - | CT | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/dashboard/categories/new` | categories | none | POST /categories | - | CT | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/dashboard/drivers` | drivers | none | DELETE /admin/drivers/{}<br>GET /admin/drivers<br>PUT /admin/drivers/{}/status | - | CT | 17 | 0/5/1/0/0 | Y | 0 |
| `/{locale}/dashboard/drivers/[id]` | drivers | none | GET /admin/drivers/{}<br>PUT /admin/drivers/{} | - | CT | 8 | 1/3/0/0/13 | Y | 0 |
| `/{locale}/dashboard/drivers/new` | drivers | none | POST /admin/drivers | - | CT | 6 | 1/3/0/0/4 | - | 0 |
| `/{locale}/dashboard/hr/attendance` | hr | HR (READ,WRITE) | DELETE /hr/attendance/{}<br>GET /hr/attendance<br>GET /hr/employees<br>POST /hr/attendance/clock-in<br>PUT /hr/attendance/{} | useRbac | CT | 0 | 0/3/0/1/2 | Y | 0 |
| `/{locale}/dashboard/hr/contracts` | hr | none | - | useContracts | NI | 0 | 0/1/0/1/0 | - | 0 |
| `/{locale}/dashboard/hr/employees` | hr | HR (DELETE,READ,WRITE) | DELETE /hr/employees/{}<br>GET /hr/employees<br>POST /hr/employees<br>PUT /hr/employees/{} | useRbac | CT | 0 | 0/3/0/1/0 | - | 0 |
| `/{locale}/dashboard/hr/insurance` | hr | none | - | useInsurance | NI | 0 | 0/1/0/1/0 | - | 0 |
| `/{locale}/dashboard/hr/leaves` | hr | LEAVES (READ,WRITE) | DELETE /hr/leaves/{}<br>GET /hr/employees<br>GET /hr/leaves<br>POST /hr/leaves<br>POST /hr/leaves/{}/approve<br>POST /hr/leaves/{}/reject<br>PUT /hr/leaves/{} | useRbac | CT | 0 | 0/5/0/2/3 | Y | 0 |
| `/{locale}/dashboard/hr/payroll` | hr | SALARIES (READ,WRITE) | DELETE /hr/salaries/{}<br>GET /hr/employees<br>GET /hr/salaries<br>POST /hr/salaries<br>POST /hr/salaries/{}/mark-paid<br>PUT /hr/salaries/{} | useRbac | CT | 0 | 0/3/0/1/2 | Y | 0 |
| `/{locale}/dashboard/hr/performance` | hr | none | - | usePerformanceReviews | NI | 0 | 0/1/0/1/0 | - | 0 |
| `/{locale}/dashboard/hr/resignations` | hr | none | - | useResignations | NI | 0 | 0/3/0/1/0 | - | 0 |
| `/{locale}/dashboard/hr/training` | hr | none | - | useTraining | NI | 0 | 0/1/0/1/0 | - | 0 |
| `/{locale}/dashboard/inventory` | inventory | none | GET /products<br>POST /inventory/adjust | - | CT | 0 | 2/3/1/1/2 | - | 0 |
| `/{locale}/dashboard/kds` | kds | KDS (ADMIN,DELETE,READ,WRITE) | - | useKDS, useKDSAnalytics, useKDSAudio, useKDSWebSocket, useRbac | CT | 10 | 0/2/0/0/0 | - | 0 |
| `/{locale}/dashboard/kds/settings` | kds | KDS (ADMIN) | GET /kds/settings<br>PUT /kds/settings | useKDS, useRbac | CT | 14 | 0/2/0/0/6 | - | 0 |
| `/{locale}/dashboard/notifications` | notifications | NOTIFICATIONS (DELETE,READ) | - | useNotifications, useRbac | CT | 5 | 0/3/0/0/4 | Y | 0 |
| `/{locale}/dashboard/orders` | orders | ORDERS (DELETE,READ) | DELETE /admin/orders/{}<br>GET /admin/orders | useRbac | CT | 14 | 0/7/1/0/4 | Y | 0 |
| `/{locale}/dashboard/orders/[id]` | orders | ORDERS (READ,WRITE) | GET /admin/orders/{}<br>PUT /admin/orders/{}/status | useRbac | CT | 0 | 0/6/0/0/1 | Y | 0 |
| `/{locale}/dashboard/pos` | pos | ORDERS (WRITE) | - | posApi, usePOSAudio, usePOSBackgroundSync, usePOSBarcodeScanner, usePOSData, usePOSKeyboard, usePOSOffline, usePOSPerformance, usePOSPrint, usePOSProductIndex, useRbac | CT | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/dashboard/pos/settings` | pos | ORDERS (WRITE) | - | useRbac | CT | 14 | 0/5/0/0/4 | - | 0 |
| `/{locale}/dashboard/products` | products | PRODUCTS (DELETE,READ,WRITE) | DELETE /products/{}<br>GET /categories<br>GET /products | useRbac | CT | 0 | 0/2/0/0/0 | - | 0 |
| `/{locale}/dashboard/products/[id]/edit` | products | PRODUCTS (WRITE) | GET /categories<br>GET /products/{}<br>PUT /products/{} | useRbac | CT | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/dashboard/products/new` | products | PRODUCTS (WRITE) | GET /categories<br>POST /products | useRbac | CT | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/dashboard/rbac/roles` | rbac | none | - | useRbac | CT | 0 | 1/6/0/1/2 | - | 0 |
| `/{locale}/dashboard/reports` | reports | none | - | - | - | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/dashboard/reports/financial` | reports | none | - | financialReportsApi, useFinancialReports | NI+CT | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/dashboard/reports/financial/advanced` | reports | none | - | accountingAdvancedApi | CT | 5 | 0/3/0/0/2 | - | 0 |
| `/{locale}/dashboard/reports/financial/balance-sheet` | reports | none | - | financialReportsApi, useFinancialReports | NI+CT | 0 | 0/1/0/0/1 | - | 0 |
| `/{locale}/dashboard/reports/financial/cash-flow` | reports | none | - | financialReportsApi, useFinancialReports | NI+CT | 0 | 0/1/0/0/2 | - | 0 |
| `/{locale}/dashboard/reports/financial/inventory-valuation` | reports | none | - | financialReportsApi, useFinancialReports | NI+CT | 0 | 0/1/0/0/1 | - | 0 |
| `/{locale}/dashboard/reports/financial/payables-aging` | reports | none | - | financialReportsApi, useFinancialReports | NI+CT | 0 | 0/1/0/0/1 | - | 0 |
| `/{locale}/dashboard/reports/financial/profit-loss` | reports | none | - | financialReportsApi, useFinancialReports | NI+CT | 0 | 0/1/0/0/4 | - | 0 |
| `/{locale}/dashboard/reports/financial/receivables-aging` | reports | none | - | financialReportsApi, useFinancialReports | NI+CT | 0 | 0/1/0/0/1 | - | 0 |
| `/{locale}/dashboard/reports/financial/vat` | reports | none | - | financialReportsApi, useFinancialReports | NI+CT | 0 | 0/1/0/0/2 | - | 0 |
| `/{locale}/dashboard/reports/hr` | reports | REPORTS (READ) | - | useHRReports, useRbac | NI+CT | 5 | 0/1/1/0/3 | - | 0 |
| `/{locale}/dashboard/reports/inventory` | reports | REPORTS (READ) | - | reportsApi, useRbac, useReports | CT | 7 | 0/1/0/0/2 | Y | 0 |
| `/{locale}/dashboard/reports/overview` | reports | REPORTS (READ) | - | reportsApi, useRbac, useReports | CT | 5 | 0/1/0/0/2 | Y | 0 |
| `/{locale}/dashboard/reports/profitability` | reports | REPORTS (READ) | - | reportsApi, useRbac, useReports | CT | 7 | 0/1/0/0/2 | Y | 0 |
| `/{locale}/dashboard/reports/sales` | reports | REPORTS (READ) | - | reportsApi, useRbac, useReports | CT | 7 | 0/1/0/0/3 | Y | 0 |
| `/{locale}/dashboard/settings` | settings | none | - | - | - | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/dashboard/settings/access-control` | settings | none | - | - | CT | 0 | 0/0/0/0/0 | - | 0 |
| `/{locale}/dashboard/settings/account` | settings | none | FETCH /users/{}/profile<br>FETCH /users/{}/settings/change-password | - | CT | 0 | 0/1/0/0/6 | - | 2 |
| `/{locale}/dashboard/settings/appearance` | settings | none | FETCH /users/{}/settings/theme<br>FETCH /users/{}/settings/theme-colors | themeApi | CT | 0 | 0/5/0/0/6 | - | 2 |
| `/{locale}/dashboard/settings/billing` | settings | none | - | usePayment | CT | 4 | 0/4/0/1/0 | - | 0 |
| `/{locale}/dashboard/settings/restaurant` | settings | none | GET /tenant/contact-info<br>PUT /tenant/contact-info | - | CT | 12 | 0/1/0/0/3 | - | 0 |
| `/{locale}/dashboard/settings/roles` | settings | SETTINGS (WRITE) | FETCH /rbac/roles<br>FETCH /rbac/roles/{} | useRbac | - | 0 | 0/3/1/0/0 | - | 2 |
| `/{locale}/dashboard/settings/roles/[id]/edit` | settings | ACCOUNTING,CATEGORIES,CONTRACTS,FINANCIAL_REPORTS,HR,KDS,LEAVES,NOTIFICATIONS,ORDERS,PERFORMANCE,PRODUCTS,REPORTS,ROLES,SALARIES,SETTINGS,SUPPLIERS,TABLES,THEMES,TRAINING,USERS (ADMIN,DELETE,READ,WRITE) | FETCH /rbac/roles/{} | useRbac | - | 0 | 1/2/0/0/2 | - | 2 |
| `/{locale}/dashboard/settings/roles/new` | settings | ACCOUNTING,CATEGORIES,CONTRACTS,FINANCIAL_REPORTS,HR,KDS,LEAVES,NOTIFICATIONS,ORDERS,PERFORMANCE,PRODUCTS,REPORTS,ROLES,SALARIES,SETTINGS,SUPPLIERS,TABLES,THEMES,TRAINING,USERS (ADMIN,DELETE,READ,WRITE) | FETCH /rbac/roles | useRbac | - | 0 | 1/2/0/0/2 | - | 1 |
| `/{locale}/dashboard/settings/users` | settings | SETTINGS (WRITE) | FETCH /rbac/roles<br>FETCH /rbac/users/{}/roles<br>FETCH /rbac/users/{}/roles/{}<br>FETCH /users<br>FETCH /users/{} | useRbac | - | 0 | 2/19/1/1/9 | - | 9 |
| `/{locale}/dashboard/tables` | tables | TABLES (DELETE,READ,WRITE) | - | tablesApi, useRbac, useTables | CT | 10 | 0/6/0/0/0 | - | 0 |
| `/{locale}/dashboard/tables/[id]` | tables | TABLES (READ,WRITE) | - | tablesApi, useRbac | CT | 4 | 0/5/0/0/0 | - | 0 |
| `/{locale}/dashboard/tables/[id]/edit` | tables | TABLES (READ,WRITE) | - | tablesApi, useRbac | CT | 4 | 0/1/0/0/0 | - | 0 |
| `/{locale}/dashboard/tables/[id]/qr` | tables | TABLES (READ) | - | tablesApi, useRbac | CT | 3 | 0/0/0/0/0 | - | 0 |
| `/{locale}/dashboard/tables/layout` | tables | TABLES (READ,WRITE) | - | tablesApi, useRbac, useTables | CT | 6 | 1/5/0/1/4 | - | 0 |
| `/{locale}/dashboard/theme-builder` | theme-builder | none | - | themeApi | CT | 1 | 0/11/0/0/3 | - | 0 |
| `/{locale}/dashboard/theme-builder/editor/[id]` | theme-builder | none | - | themeApi, useThemeStore | CT | 9 | 0/5/0/0/0 | - | 0 |
| `/{locale}/dashboard/theme-builder/import` | theme-builder | none | - | themeApi | CT | 0 | 0/4/0/0/2 | - | 0 |
| `/{locale}/dashboard/theme-builder/presets` | theme-builder | none | - | themeApi | CT | 0 | 0/1/0/0/0 | - | 0 |
| `/{locale}/dashboard/theme-builder/qr` | theme-builder | none | - | themeApi | CT | 3 | 0/0/0/0/0 | - | 0 |
| `/{locale}/test-url` | shell | none | - | - | - | 0 | 0/0/0/0/0 | - | 0 |

### Per-module page tally

| Module | Pages | Pages with no page-level gate |
|---|---|---|
| accounting | 10 | 1 |
| auth | 3 | 3 |
| categories | 3 | 3 |
| drivers | 3 | 3 |
| hr | 9 | 5 |
| inventory | 1 | 1 |
| kds | 2 | 0 |
| notifications | 1 | 0 |
| orders | 2 | 0 |
| overview | 1 | 0 |
| pos | 2 | 0 |
| products | 3 | 0 |
| rbac | 1 | 1 |
| reports | 15 | 10 |
| settings | 10 | 6 |
| shell | 3 | 3 |
| tables | 5 | 0 |
| theme-builder | 5 | 5 |

Notes
- Pages where the gate is absent but a child component gates (e.g. `components/hr/*`, `components/accounting/*`) are not detected by this grep; treat 'none' as 'no gate in page file'. Sidebar (`components/layout/Sidebar.tsx:44-62`) hides links by MODULE_IDS but direct URL access still renders the shell until the API returns 403.
- Pages calling axios directly: only `kds/settings/page.tsx` imports axios (7 refs) in addition to `lib/api.ts` and `lib/api/themeApi.ts`.
- Business logic in pages: dashboard overview (`recentPaidRevenue`), drivers (completion % = completed/total*100), accounting invoices (outstanding sums), inventory (low-stock filter), profitability chart rounding; POS tax lives in `stores/posStore.ts:106`.
