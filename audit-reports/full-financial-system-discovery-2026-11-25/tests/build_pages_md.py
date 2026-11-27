import json, re
pages = json.load(open('tests/_pages_raw.json'))
P = [  # (regex on url, purpose, financial relevance)
 (r'/dashboard/accounting/advanced$', 'Advanced accounting console: bank reconciliation, budgets, FX rates, supplier payment schedules, e-invoice submit, CSV import/export, recurring journals', 'HIGH - writes money; may post journals (recurring run, CSV import)'),
 (r'/accounting/bank-accounts$', 'Bank accounts CRUD + transactions list', 'HIGH - master data for cash/bank GL linkage'),
 (r'/accounting/chart-of-accounts$', 'Chart of accounts tree + AccountForm modal', 'HIGH - GL structure'),
 (r'/accounting/customer-invoices$', 'AR invoices list/create/edit/send/void + record payment received', 'HIGH - posts journal on create/void/payment'),
 (r'/accounting/customers$', 'AR customer accounts + statements', 'HIGH - receivable balances'),
 (r'/accounting/general-ledger$', 'Ledger by account + trial balance', 'HIGH - reads GL'),
 (r'/accounting/journal-entries/\[id\]$', 'Journal entry detail; post/void/edit draft', 'HIGH - posts/voids journal'),
 (r'/accounting/journal-entries$', 'Journal entry list + create form', 'HIGH - creates drafts'),
 (r'/accounting/supplier-invoices$', 'AP invoices list/create/edit/approve/void + record payment made', 'HIGH - posts journal on approve/void/payment'),
 (r'/accounting/suppliers$', 'Supplier master CRUD (module SUPPLIERS=16)', 'HIGH - AP master'),
 (r'/dashboard/hr/payroll$', 'Payroll: salary records list/create/edit/delete/mark-paid', 'HIGH - mark-paid posts payroll journal'),
 (r'/dashboard/hr/employees$', 'Employee master (may include pay terms)', 'MEDIUM - HR financial terms'),
 (r'/dashboard/hr/contracts$', 'Employment contracts (salary terms)', 'MEDIUM'),
 (r'/dashboard/hr/insurance$', 'Insurance tracking', 'MEDIUM'),
 (r'/dashboard/hr/resignations$', 'Resignations + settlement calculation', 'MEDIUM - settlement amounts (read)'),
 (r'/dashboard/hr/(attendance|leaves|performance|training)$', 'HR operational module', 'LOW - feeds payroll indirectly (attendance/leave)'),
 (r'/dashboard/inventory$', 'Stock adjustment UI', 'HIGH - stock change without cost/journal'),
 (r'/dashboard/kds(/settings)?$', 'Kitchen display (WS + bump/recall) / settings', 'LOW'),
 (r'/dashboard/orders/\[id\]$', 'Order detail, status update, driver assign', 'MEDIUM - totals read; no payment-status control'),
 (r'/dashboard/orders$', 'Order list, filter by payment status, cancel', 'MEDIUM'),
 (r'/dashboard/pos$', 'POS terminal: cart, payment method, create order, auto-mark paid (cash only)', 'HIGH - revenue trigger'),
 (r'/dashboard/pos/settings$', 'POS settings', 'LOW'),
 (r'/dashboard/products(/new|/\[id\]/edit)?$', 'Product list / create / edit (price, cost, stock)', 'HIGH - price/cost master'),
 (r'/dashboard/categories', 'Category management', 'NONE'),
 (r'/dashboard/reports/financial/(balance-sheet|cash-flow|profit-loss|inventory-valuation|payables-aging|receivables-aging|vat)$', 'Financial statement report', 'HIGH - reads money (module 17)'),
 (r'/dashboard/reports/financial/advanced$', 'KPIs, break-even, menu profitability, predictive cash flow, HR cost, audit export', 'HIGH - reads money'),
 (r'/dashboard/reports/financial$', 'Financial reports hub (VAT summary)', 'HIGH - reads money'),
 (r'/dashboard/reports/(sales|overview|profitability|inventory|hr)$', 'Operational report (module REPORTS=8)', 'MEDIUM - reads operational money (orders table)'),
 (r'/dashboard/reports$', 'Reports hub', 'LOW'),
 (r'/dashboard/settings/billing$', 'SaaS subscription billing (Paymob)', 'MEDIUM - platform billing, not restaurant ledger'),
 (r'/dashboard/settings/restaurant$', 'Restaurant name/email/phone (no tax or currency fields)', 'LOW'),
 (r'/dashboard/(rbac|settings/(roles|users|access-control))', 'RBAC / user admin', 'MEDIUM - controls who can post/approve'),
 (r'/dashboard/settings', 'Settings page', 'LOW'),
 (r'/dashboard/tables', 'Table management / QR', 'LOW'),
 (r'/dashboard/drivers', 'Delivery drivers', 'LOW'),
 (r'/dashboard/notifications$', 'Notifications', 'NONE'),
 (r'/dashboard/theme-builder', 'Storefront theme builder (calls theme-component endpoints that have no backend route - see section 4)', 'NONE'),
 (r'/dashboard$', 'Dashboard home / KPIs', 'LOW - reads order stats'),
 (r'/auth/', 'Auth callback / login', 'NONE'),
 (r'^/$|^/\[locale\]$', 'Root / locale redirect to dashboard or login', 'NONE'),
 (r'/test-url$', 'Dev test page (should not exist in production)', 'NONE'),
]
LANDING = {
 '/': 'Entry redirect', '/[locale]': 'Marketing home', '/[locale]/pricing': 'Pricing page', '/[locale]/contact': 'Contact', '/[locale]/demo': 'Demo request', '/[locale]/features': 'Features list', '/[locale]/features/[slug]': 'Feature detail', '/[locale]/how-it-works': 'How it works', '/[locale]/privacy': 'Privacy policy',
 '/[locale]/auth/login': 'Login (calls POST /auth/login)', '/[locale]/auth/register': 'Tenant registration (POST /auth/register, GET /auth/check-subdomain)', '/[locale]/auth/super-admin-login': 'Super-admin login redirect', '/auth/login': 'Entry login', '/auth/register': 'Entry register', '/auth/super-admin-login': 'Entry super-admin login', '/login': 'Entry login alias', '/register': 'Entry register alias'}
RW = {
 '/': 'Root redirect to locale', '/[locale]': 'Storefront home', '/[locale]/menu': 'Menu browse', '/[locale]/cart': 'Cart', '/[locale]/checkout': 'Checkout: quote + create public order (cash only)', '/[locale]/orders': 'Customer orders list', '/[locale]/orders/[id]': 'Order detail/status', '/[locale]/orders/track': 'Order tracking form', '/[locale]/orders/track/[orderNumber]': 'Track by number', '/[locale]/settings': 'Storefront settings', '/[locale]/table/[tableId]/menu': 'Dine-in QR menu/ordering', '/api/metrics': 'Next route handler: metrics', '/maintenance': 'Maintenance page (tenant blocked)'}
SA = {
 '/': 'Redirect', '/login': 'Super-admin login', '/dashboard': 'Platform stats', '/dashboard/tenants': 'Tenant list', '/dashboard/tenants/[id]': 'Tenant detail: subscription, extend trial, enable/disable, impersonate', '/dashboard/users': 'All users', '/dashboard/payments': 'SaaS payment metrics/webhook health/refund', '/dashboard/audit-logs': 'Super-admin audit log', '/dashboard/settings': 'IP whitelist', '/favicon.ico': 'Route handler'}
def purpose(app, url):
    if app == 'landing': return LANDING.get(url, ''), 'NONE'
    if app == 'restaurant-website':
        fr = 'HIGH - creates orders' if url in ('/[locale]/checkout',) else ('MEDIUM - shows totals' if url in ('/[locale]/cart', '/[locale]/orders', '/[locale]/orders/[id]') else 'LOW')
        return RW.get(url, ''), fr
    if app == 'super-admin':
        return SA.get(url, ''), ('MEDIUM - platform billing/refund' if 'payments' in url else 'LOW')
    for rx, p, f in P:
        if re.search(rx, url): return p, f
    return '', ''
def guard(x):
    parts = []
    if x['perms']:
        parts.append('client checkPermission: ' + ', '.join(x['perms']))
    if x['app'] == 'dashboard' and '/dashboard' in x['url']:
        parts.append('DashboardAuthProvider (layout.tsx) session gate')
    if x['app'] == 'super-admin' and x['url'] != '/login': parts.append('super-admin layout (dashboard/layout.tsx)')
    if x['app'] == 'dashboard' and '/dashboard' in x['url'] and not x['perms']: parts.append('NO client permission check in page closure (backend only)')
    if not parts: parts.append('public / none')
    return '; '.join(parts)
out = []
out.append('# 06 - Complete Page Inventory\n')
out.append('Baseline SHA 133260db66b1143c254d79476daa7650b41a33f0. Generated by `tests/analyze_pages.py` (static import-closure scan, depth<=4, named-import-aware) then enriched manually. Read-only.\n')
out.append('## Method and honest limits\n')
out.append('- Page list = every `page.tsx|page.ts|route.ts` under `frontend/apps/{dashboard,landing,super-admin}/src/app` and `frontend/apps/restaurant-website/app` (119 entries; layouts listed separately). ')
out.append('- Endpoints = string literals that match a registered backend route in `08-api-inventory.csv` (verified match). Method shown as `?` where the call style could not be inferred. Endpoints reached only through dynamic string building are UNVERIFIED/missed.')
out.append('- Shared permission hook calls (`/rbac/me/permissions`, `/rbac/check-permission`, `/rbac/*`, `/auth/*`) are omitted from per-page lists (every dashboard page pulls them via `useUserPermissions`); they appear only on RBAC/settings pages.')
out.append('- "Unmatched" = a literal that looks like an API path but matches no registered route (possible missing backend route or a Next.js path).')
out.append('- Locale: dashboard/landing/restaurant-website routes sit under `/[locale]` (en|ar). No Next route handlers exist in dashboard/landing; restaurant-website has `/api/metrics`; super-admin has `favicon.ico/route.ts`.\n')
for app in ['dashboard', 'landing', 'restaurant-website', 'super-admin']:
    ps = [p for p in pages if p['app'] == app]
    out.append(f'\n## {app} ({len(ps)} pages)\n')
    out.append('| URL pattern | File | Purpose | API endpoints called (verified vs 08) | Permission / guard | Financial relevance |')
    out.append('|---|---|---|---|---|---|')
    for x in ps:
        pu, fr = purpose(app, x['url'])
        eps = [e for e in x['endpoints'] if not re.search(r' /(rbac|auth)(/|$)', e)]
        if app == 'dashboard' and re.search(r'/(rbac|settings/(roles|users|access-control))', x['url']):
            eps = x['endpoints']
        eps_s = '<br>'.join(eps[:16]) + (f'<br>(+{len(eps)-16} more)' if len(eps) > 16 else '') if eps else '-'
        um = [u for u in x['unmatched'] if not u.endswith('/auth/') and not u.endswith('/auth/callback') and not u.endswith(' /orders')]
        if um: eps_s += '<br>UNMATCHED: ' + '<br>'.join(um[:6])
        out.append(f"| `{x['url']}` | `{x['file']}` | {pu} | {eps_s} | {guard(x)} | {fr} |")
out.append('''
## Layout / middleware files discovered
- dashboard: `src/app/layout.tsx`, `[locale]/layout.tsx`, `[locale]/dashboard/layout.tsx` (DashboardAuthProvider), `kds/layout.tsx`, `pos/layout.tsx`, `settings/layout.tsx`, `theme-builder/layout.tsx`, `theme-builder/editor/layout.tsx`; `src/middleware.ts` only sets trusted `x-dashboard-locale` header (no auth).
- landing: `(entry)/layout.tsx`, `[locale]/layout.tsx`, `[locale]/auth/layout.tsx`.
- restaurant-website: `app/layout.tsx`, `[locale]/layout.tsx`, `middleware.ts` (tenant/subdomain resolution; contents not audited in this pass - UNVERIFIED).
- super-admin: `src/app/layout.tsx`, `dashboard/layout.tsx`.
''')
out.append('''
## Verified frontend<->backend mismatches
1. Theme builder calls endpoints with no registered backend route: `/admin/themes/{id}/components`, `/admin/themes/{id}/components/{id}`, `/admin/themes/{id}/components/reorder`, `/admin/themes/import`, `/admin/themes/{id}/export`, `/public/component-library` (main.go:618 shows the component-library route commented out). Restaurant-website calls `/public/themes` and `/public/themes/{id}` - also unregistered. Not financial.
2. No frontend page calls `GET /accounting/paid-order-events/failed` or `POST /accounting/paid-order-events/{id}/redrive` (grep over frontend/apps found no reference): outbox recovery is API-only, no UI.
3. Only `lib/api/posApi.ts:266` calls `PUT /admin/orders/{id}/payment-status`; the order list/detail pages have no control to mark an order paid.

## Financial pages: modals, tabs, forms discovered
### Accounting (dashboard/src/components/accounting)
- Accounts: `AccountTree` + `AccountForm` (create/edit modal), delete confirm.
- Journal entries: `JournalEntryList`, `JournalEntryForm` (create/edit draft), detail page with Post and Void actions.
- General ledger: `LedgerView` (account/date filters), `TrialBalanceTable`.
- Customer invoices / supplier invoices: `InvoiceList`, `InvoiceForm` (mode=customer|supplier, line grid), `PaymentForm` (mode=received|made) opened from invoice rows, Void and Send/Approve row actions.
- Customers / Suppliers: `CustomerList`/`CustomerForm`, `SupplierList`/`SupplierForm`, customer statement view.
- Bank accounts: `BankAccountForm`, `BankTransactionList`.
- Advanced (`accounting/advanced/page.tsx`, single page, inline forms, no tabs): bank reconciliation session (+lines), budget upsert, exchange rate upsert, supplier payment schedule + send-reminder, e-invoice submit (invoice ID defaults to 1, page.tsx:75), recurring journal create/run, CSV import/export for chart of accounts and journal entries (file inputs accept .csv).
### Reports (dashboard/src/components/financial-reports, components/reports)
- `BalanceSheet`, `ProfitLossStatement`, `CashFlowStatement`, `AgingReport` (payables and receivables), `VATReport`, `InventoryValuationReport`, `ReportExportButton`; operational: `ReportTable`, `ReportChart`, `KpiCard`, `ExportButton`.
### HR / payroll
- `SalaryForm` (react-hook-form; modal opened from payroll page), `SalaryList` (mark-paid, edit, delete), `EmployeeForm`, `LeaveForm`, `AttendanceForm`; hr-advanced: `ContractForm`, `InsuranceForm`, `ResignationForm/Workflow/ExitChecklist`, `PerformanceReviewForm`, `TrainingForm`.
### POS (dashboard/pos/components)
- `POSProductGrid/VirtualGrid`, `POSCartPanel/CartItem`, `POSCustomerInput`, `POSPaymentSelector` (cash|card), `POSTableSelector`, `POSOrderSummary` (tax shown as 14%), `POSSubmitButton`, `POSReceipt`, offline queue (`usePOSOffline.ts`).
### Orders / KDS
- Orders list filters (payment_status), cancel, driver-assignment modal (`components/driver-assignment-modal.tsx`), order detail status control. KDS board bump/recall, KDS settings.
### Products / inventory
- `ProductForm`, `ProductFormPage`, `ProductList`, inventory adjust UI (`/dashboard/inventory`).
### Settings
- `settings/restaurant` (name/email/phone only), `settings/billing` (`PlanCard`, `PaymentModal`, `PaymentHistoryTable`), RBAC: `RBACSettings`, roles new/edit, users, access-control.
### Restaurant storefront
- `cart/page` (`CartSheet`, `CartItem`), `checkout/checkout-page-content.tsx` (react-hook-form + zod; `OrderQuoteSummary`), `orders` list/detail/track, table QR menu.
''')
open('06-complete-page-inventory.md', 'w', encoding='utf-8').write('\n'.join(out))
print(len(pages))
