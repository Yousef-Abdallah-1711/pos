# Application architecture map

| Application | Root | Verified entry and ownership |
| --- | --- | --- |
| Landing | `frontend/apps/landing` | `src/app/[locale]/page.tsx`, `src/components/RootDocument.tsx`; platform marketing and auth entry, Coral marketing tokens. |
| Tenant dashboard | `frontend/apps/dashboard` | `src/app/[locale]/layout.tsx` → dashboard layout → `src/components/layout/DashboardLayout.tsx`; personal preferences plus site theme builder. |
| POS / kitchen | `frontend/apps/dashboard` | `src/app/[locale]/dashboard/pos` and `.../kds`; fullscreen routes have distinct shell treatment. |
| Public storefront | `frontend/apps/restaurant-website` | `app/layout.tsx`, `app/[locale]/layout.tsx`, `components/providers/theme-provider.tsx`, `lib/api/get-theme-server.ts`; separate site theme scope. |
| Platform admin | `frontend/apps/super-admin` | `src/app/layout.tsx`, `src/app/dashboard/layout.tsx`; independent administrative app. |
| Backend | `backend` | `cmd/api/main.go`, `internal/handler/http`, `internal/repository`, `migrations`; authentication, tenant/RLS and persistence. |
| Shared packages | `frontend/packages` | Includes `form-foundation` and `i18n`; inspect package consumers and version boundaries before changes. |

Actual permissions, routes and authenticated behavior require verification per application; framework declarations and package versions are available in each `package.json`. Do not assume storefront theme inheritance for the dashboard or tenant inheritance for super-admin.

At the inspected checkout, all four app manifests declare Next 16.4.6 and React 19.3.1. Dashboard and storefront declare `next-intl` 4.15.1 and `next-themes` ^0.4.4; landing declares `next-themes` ^0.4.4; the super-admin manifest did not declare those two dependencies directly. Dependency declarations do not prove runtime behavior or license terms.

## Theme, identity and localization boundaries

| Surface | Source evidence | Boundary for the future change |
| --- | --- | --- |
| Landing and public authentication | `frontend/apps/landing/src/components/RootDocument.tsx`, `src/app/[locale]/layout.tsx`, `src/components/marketing/MarketingShell.tsx` | Platform marketing palette and locale/theme selection. Login/register are entry points to backend sessions; changing tenant dashboard tokens must not depend on an authenticated landing session. |
| Tenant dashboard, POS and KDS | `frontend/apps/dashboard/src/app/[locale]/layout.tsx`, `src/components/layout/DashboardLayout.tsx`, `src/stores/authStore.ts`, `src/stores/preferencesStore.ts` | Inter wrapper, locale/RTL, cookie-session user and tenant context, personal settings and class-based mode; POS/KDS use fullscreen variants. See current lifecycle for unresolved cache behavior. |
| Storefront and website theme builder | `frontend/apps/restaurant-website/app/layout.tsx`, `app/[locale]/layout.tsx`, `components/providers/theme-provider.tsx`, `lib/subdomain-context.tsx` | Root loads Inter and Cairo, derives initial locale from request headers, nests subdomain, theme and query providers, and uses separate site theme variables. Locale layout uses `next-intl`. Public tenant resolution and site branding are separate from private workspace appearance. |
| Super-admin | `frontend/apps/super-admin/src/app/layout.tsx`, `src/components/providers/PreferencesProvider.tsx`, `SuperAdminAuthProvider.tsx` | Separate super-admin auth and `servio-super-admin-preferences` browser storage. Provider updates root `lang`, `dir`, dark class and color scheme. Do not assume a restaurant tenant theme should control this platform scope. |
| Shared packages and layers | `frontend/packages/{form-foundation,i18n,theme,ui,ui-components}`, `frontend/layers/` | Identify actual imports and ownership before replacing shared tokens: packages may serve multiple applications with different scope and release risk. |
| Backend and database | `backend/cmd/api/main.go`, `backend/internal/handler/http`, `backend/internal/repository`, `backend/internal/pkg/database/tenant_tx.go`, `backend/migrations/` | Authenticated claims and tenant-scoped transactions enforce data scope; UI state alone is never authoritative. Existing user and storefront theme data require separate migration decisions. |

The table identifies verified source entry points, not an exhaustive dependency graph. Authentication behavior, response caching, host-based tenant resolution and cross-app redirects require deployment-specific testing.
