# Sidebar and shell evidence

`frontend/apps/dashboard/src/components/layout/Sidebar.tsx` defines the navigation items inline. It uses locale-aware labels, Lucide icons, pathname/search parameters, `useUserPermissions`, `MODULE_IDS` and `PERMISSION_LEVELS`; `useTenantStatus` can substitute blocked navigation. Groups include orders, products, HR, accounting, suppliers, financial reports, reports, theme builder and settings. Visibility checks do not substitute for server-side authorization; route-level enforcement requires independent review.

The sidebar uses fixed `w-64` expanded and `w-20` collapsed widths, gray surfaces and borders, primary-color logo accents, responsive transforms for a mobile drawer, and locale-dependent left/right alignment. `DashboardLayout.tsx` owns open/collapsed state and persists `dashboard.sidebar.collapsed` in localStorage. `src/app/globals.css` coordinates 16rem/5rem content padding and RTL on desktop. The shell has a sticky 64px header, notifications, locale switcher, mode toggle and user menu. POS/KDS fullscreen routes skip header and sidebar; review their distinct surfaces and controls separately. Portal styling and keyboard behavior remain unverified.

`isChildActive` matches the localized pathname and optional `tab` query; groups are expanded by label in local component state. The source returns early from `toggleExpanded` when collapsed and only renders nested links when expanded. Assess discoverability and keyboard access of nested destinations in collapsed mode during a real interaction audit. There are special visibility fallbacks for orders/tables, suppliers/accounting, financial reports/general reports, and advanced HR; record these as compatibility behavior before changing navigation. The loading/blocked branches and their duplicated source blocks also need review rather than assuming the normal sidebar covers every state.

## Sidebar groups and source-level visibility

These are navigation-display rules from `Sidebar.tsx:45`–`:320`, not a guarantee of backend authorization. Child routes may have additional checks.

| Entry | Destination or children | Visibility evidence |
| --- | --- | --- |
| Dashboard / Settings | `/dashboard`, `/dashboard/settings` | No module ID in the normal sidebar; shown unconditionally by its filter. |
| Cashier POS | `/dashboard/pos` | Orders WRITE. |
| Kitchen Display | `/dashboard/kds` | KDS READ. |
| Orders | Orders, tables, drivers | Parent visible for Orders READ or Tables READ; children with module IDs filtered separately. |
| Products | Products, categories, inventory | Products READ parent; categories use Categories READ, inventory uses Products READ. |
| Human Resources | Employees, attendance, leaves, payroll, performance, training, insurance, contracts, resignations | HR READ parent; Leaves and Salaries have individual checks; performance/training/contracts children also have an HR compatibility fallback. |
| Accounting | Chart of accounts, journals, ledger, bank accounts, advanced tools, customers, customer invoices and payment tab | Accounting READ. |
| Suppliers | Suppliers, supplier invoices and payment tab | Supplier READ or Accounting READ parent; individual child module checks still apply. |
| Financial Reports | Financial hub, P&L, balance sheet, cash flow, VAT, aging, valuation, advanced | Financial Reports READ or Reports READ compatibility fallback. |
| Reports | Overview, sales, profitability, inventory and HR reports | Reports READ. |
| Theme Builder | `/dashboard/theme-builder` | Themes READ; this is the storefront builder. |

The source marks a group active by matching child pathname and optional `tab` search parameter. Its navigation definitions use Lucide React icons and locale-prefixed links at render time. Blocked tenants render `BlockedSidebar`; loading and fullscreen POS/KDS bypass the ordinary final navigation state. Subscription badges and the blocked state require their own future capture matrix.

Relevant shared paths: `src/components/layout/{DashboardLayout,Sidebar,SidebarSubscriptionBadge,LanguageSwitcher}.tsx`, `src/components/blocked/BlockedSidebar.tsx`, `src/hooks/useRbac.ts`, `src/lib/rbac.ts`, `src/app/globals.css`, and `src/styles/rtl.css`. Color roles to migrate include shell surface, foreground, muted foreground, active/hover, icon, border, focus and badge states. Preserve semantic subscription and permission indications. Actual captures and detailed keyboard checks are still required.

## Specific source anchors and proposed migration roles

Paths below are relative to `frontend/apps/dashboard/`; line numbers describe the inspected revision and may move.

| Current dependency | Source | Current mechanism | Candidate role / review |
| --- | --- | --- | --- |
| Sidebar width, surface, divider, mobile transform | `src/components/layout/Sidebar.tsx:402` | `w-20` / `w-64`, white/gray Tailwind classes, side switches for RTL; `md` changes drawer behavior | Sidebar surface and border tokens; preserve width and focus behavior. |
| Logo icon / title | `src/components/layout/Sidebar.tsx:410` | `text-primary-600` / `dark:text-primary-400` and gray/white title | Brand icon and sidebar foreground; do not recolor semantic icons globally. |
| Group hover and nested navigation | `src/components/layout/Sidebar.tsx:433` | Gray foreground/hover utilities, conditional collapse and directional chevrons | Sidebar hover and foreground; check expanded state and keyboard operation. |
| Active child / main link | `src/components/layout/Sidebar.tsx:455` and `:475` | Primary utility classes for active, gray for inactive | Selection background/foreground, test both modes and contrast. |
| Skeleton and blocked branches | `src/components/layout/Sidebar.tsx:330`, `:360` | Separate branch with fixed gray Tailwind utilities / separate blocked sidebar | Map loading and blocked surfaces explicitly; do not assume the normal branch covers them. |
| Shell, header and profile menu | `src/components/layout/DashboardLayout.tsx:171`, `:192`, `:240`, `:258` | Gray/white Tailwind utilities and local theme toggle; profile popover rendered in shell | Workspace background, header, popover, borders and controls; confirm portal behavior separately. |
| Collapse preference and page offset | `src/components/layout/DashboardLayout.tsx:59`, `:66`; `src/app/globals.css:93` | Local storage plus layout data attributes and CSS padding | Preserve preference, responsive thresholds and RTL content offsets. |

Focus indicators, subscription badge colors, tooltips, mobile overlay accessibility, workspace switching, search and breadcrumbs require individual source and runtime verification; absence from this table is not evidence they are unsupported.

In the [synthetic staff baseline](baseline/dashboard-en-restricted-light.png), the account had no `user_roles` rows and the sidebar showed Dashboard and Settings. A Theme Builder shortcut was nevertheless visible within the main dashboard. In that same disposable browser session, direct authenticated GET requests to `/api/v1/admin/themes` and `/api/v1/products` returned HTTP 403. This narrow observation does not prove authorization of every route/action or explain all partially loaded widgets; compare shortcut visibility with server authorization in a future role matrix.

Separate synthetic sessions using the tenant's seeded [Manager](baseline/baseline-en-manager-light.png) and [Employee](baseline/baseline-en-employee-light.png) roles rendered different navigation groups. Manager showed POS, KDS, orders, products, HR, financial reports, reports, theme builder and settings; Employee showed KDS, orders, products, HR and settings. These are observations for these seeded role definitions, not a general security guarantee or proof that child actions are permitted.
