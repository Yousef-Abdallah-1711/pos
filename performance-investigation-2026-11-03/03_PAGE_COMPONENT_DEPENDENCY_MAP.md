# Page-to-Component Dependency Map (Static)

The exhaustive page-file list is `02_COMPLETE_FRONTEND_PAGE_INVENTORY.csv` (110 App Router `page.tsx` files: Dashboard 79, Landing 10, Restaurant Website 12, Super Admin 9). File counts are not guaranteed unique browser URLs: locale and dynamic segment templates represent multiple concrete URLs. Every page row is marked STATIC_ONLY; no authenticated page was rendered and no route was traced through every hook to every backend repository.

## Global/shared render dependencies observed

- **Dashboard:** root and locale layouts apply metadata/fonts and locale providers. `src/app/[locale]/layout.tsx` wraps the route tree in NextIntl, Suspense, auth, theme, WebVitals, and performance-monitoring providers. `src/app/[locale]/dashboard/layout.tsx` wraps all dashboard routes in `DashboardAuthProvider`. The latter delays route children behind client auth/persistence hydration and redirects anonymous users. `DashboardLayout` adds navigation/preferences and route preloading.
- **Landing:** root theme provider and locale `LanguageProvider`; custom locale routing. Pages do not share a React Query provider.
- **Restaurant Website:** `app/layout.tsx` includes WebVitals, SubdomainProvider, ThemeProvider, QueryProvider, and toaster globally. Host/tenant middleware may fetch tenant status before route handling. Server-rendered homepage fetches theme/products/homepage data.
- **Super Admin:** root `SuperAdminAuthProvider` wraps login and protected routes; dashboard routes have their own auth/layout loading state.

## Trace confidence

`Page → layout/provider` is partially statically mapped from source. `Page → component → hook → HTTP request → handler → service/repository → tables/Redis` remains untraced for most pages. This is intentional evidence accounting: endpoint consumers are not inferred from route names. `05_PAGE_TO_API_DEPENDENCY_MATRIX.csv` contains per-page status and follow-up fields; `06_FRONTEND_BACKEND_DEPENDENCY_GRAPH.md` gives the verified common paths and gaps.
