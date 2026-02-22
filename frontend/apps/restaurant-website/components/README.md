# Restaurant Website Components

## Folder Structure

- `layout/`
  - Shared page chrome for runtime and theme preview integration.
  - `ThemeHeader.tsx` and `ThemeFooter.tsx` are the single entry points.
  - `theme-layout-normalizers.ts` maps API/theme data into shared UI config.
- `providers/`
  - App-level providers (`query`, `preferences`, `web-vitals`).
- Root components (`*.tsx`)
  - `dynamic-home-page.tsx` renders homepage body sections using `@pos-saas/ui-themes` `SectionRenderer`.
  - `featured-products-section.tsx` appends live product data/cards after theme-configured sections.
  - Feature-specific UI (cart, checkout, product cards, tracking, etc.).
  - Legacy compatibility wrappers (`header-ssr.tsx`, `footer-ssr.tsx`).

## Rules

- Do not duplicate header/footer markup in app code.
- Use `@pos-saas/ui-themes` `Header` and `Footer` for all theme layouts.
- Use `@pos-saas/ui-themes` `SectionRenderer` + `frontend/layers/ui-themes/src/sections/*` for homepage theme sections.
- Put API-to-UI mapping logic only in `layout/theme-layout-normalizers.ts`.
- Keep route files importing wrappers (`header-ssr.tsx`, `footer-ssr.tsx`) unless doing a full migration.
