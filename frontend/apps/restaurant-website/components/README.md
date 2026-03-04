# Restaurant Website Components

## Folder Structure

- `layout/`
  - Shared page chrome for runtime and theme preview integration.
  - `Header.tsx` and `Footer.tsx` normalize server theme data and render `@pos-saas/ui-themes` layout components.
  - `header-ssr.tsx` and `footer-ssr.tsx` are the route-facing wrappers used by storefront pages.
  - `layout-normalizers.ts` maps API/theme data into shared UI config.
- `homepage/`
  - Homepage theme-builder runtime components.
  - `home-page-content.tsx` renders section order from theme config.
  - Featured products are appended inside `home-page-content.tsx`.
  - Section error boundaries live in `home-page-content.tsx`.
- `commerce/`
  - Reusable storefront commerce UI used by menu/cart/orders pages.
  - Product cards, cart widget, order list, and order tracking UI.
- `system/`
  - App-level non-commerce UI (e.g. maintenance page).
- `providers/`
  - `theme-provider.tsx` (storefront app theme vars: `--theme-*`).
  - `core/` contains app-wide infrastructure providers (`query`, `web-vitals`).

## Rules

- Do not duplicate header/footer markup in app code.
- Use `@pos-saas/ui-themes` `Header` and `Footer` for all theme layouts.
- Use `@pos-saas/ui-themes` `SectionRenderer` + `frontend/layers/ui-themes/src/sections/*` for homepage theme sections.
- Put API-to-UI mapping logic only in `layout/layout-normalizers.ts`.
- Keep route files importing layout wrappers (`layout/header-ssr.tsx`, `layout/footer-ssr.tsx`) unless doing a full migration.

