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
