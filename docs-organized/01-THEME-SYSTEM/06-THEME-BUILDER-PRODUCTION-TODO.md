# Theme Builder Production TODO

Last updated: 2026-02-21

## Goal
Stabilize Theme Builder so each restaurant can run exactly 2 reliable themes with distinct header/footer/homepage content, while shared app pages remain common.

## Completed In This Pass
- [x] Normalize dashboard theme API mapping (snake_case/camelCase/config parsing).
  - File: `frontend/apps/dashboard/src/lib/api/themeApi.ts`
- [x] Persist full theme config on save (header, footer, components, custom CSS).
  - File: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx`
  - File: `frontend/apps/dashboard/src/types/theme.ts`
- [x] Fix Theme Builder list date rendering (`Invalid Date` issue).
  - File: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/page.tsx`
- [x] Fix preview website target to support `demo.localhost:3003` domain flow.
  - File: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/page.tsx`
- [x] Fix component registry mapping and ID consistency.
  - File: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/components/ComponentLibrary.tsx`
  - File: `frontend/apps/dashboard/src/hooks/useThemeStore.ts`
- [x] Fix component ordering fallback (`displayOrder` and `order`).
  - File: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/components/EditorPreview.tsx`
  - File: `frontend/apps/restaurant-website/components/dynamic-home-page.tsx`
- [x] Enable header/footer runtime compatibility for both camelCase and snake_case configs.
  - File: `frontend/apps/restaurant-website/components/header-ssr.tsx`
  - File: `frontend/apps/restaurant-website/components/footer-ssr.tsx`
- [x] Fix public theme response extraction to support nested API payload variants (`data.data`, `data.config`, `theme.config`, JSON string config).
  - File: `frontend/apps/restaurant-website/lib/api/get-theme-server.ts`
  - File: `frontend/apps/restaurant-website/app/[locale]/page.tsx`
- [x] Remove incorrect per-page theme re-extraction and use normalized theme object directly across website pages.
  - File: `frontend/apps/restaurant-website/app/[locale]/cart/page.tsx`
  - File: `frontend/apps/restaurant-website/app/[locale]/menu/page.tsx`
  - File: `frontend/apps/restaurant-website/app/[locale]/checkout/page.tsx`
  - File: `frontend/apps/restaurant-website/app/[locale]/orders/page.tsx`
  - File: `frontend/apps/restaurant-website/app/[locale]/settings/page.tsx`
  - File: `frontend/apps/restaurant-website/app/[locale]/orders/[id]/page.tsx`
  - File: `frontend/apps/restaurant-website/app/[locale]/orders/track/page.tsx`
  - File: `frontend/apps/restaurant-website/app/[locale]/orders/track/[orderNumber]/page.tsx`
- [x] Replace static-like storefront header/footer with theme-config-driven runtime renderer (uses HeaderEditor/FooterEditor fields).
  - Header fields now respected at runtime: `logoUrl`, `logoText`, `logoHeight`, `showLogo`, `navigationItems`, `navPosition`, `navAlignment`, `backgroundColor`, `textColor`, `height`, `padding`, `stickyHeader`, `showShadow|shadow`, `hideNavOnMobile`.
  - Footer fields now respected at runtime: `companyName`, `companyDescription`, `address`, `phone`, `email`, `copyrightText`, `socialLinks`, `footerSections`, `legalLinks`, `backgroundColor`, `textColor`, `linkColor`, `padding`, `showLinks`, `showLegal`, `showBackToTop`, `columns`, `layout`.
  - File: `frontend/apps/restaurant-website/components/header-ssr.tsx`
  - File: `frontend/apps/restaurant-website/components/footer-ssr.tsx`
- [x] Apply custom CSS on restaurant homepage runtime.
  - File: `frontend/apps/restaurant-website/components/dynamic-home-page.tsx`
- [x] Add color sync from dashboard appearance settings to active website theme.
  - File: `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/appearance/page.tsx`
- [x] Add Theme Builder bulk actions: migrate configs and keep only 2 themes.
  - File: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/page.tsx`
- [x] Remove inline React style usage from active Theme Builder editors and related theme runtime surfaces.
  - Updated editors: `ContactEditor.tsx`, `TestimonialsEditor.tsx`, `ProductGridEditor.tsx`, `HeaderEditor.tsx`, `FooterEditor.tsx`, `HeroEditor.tsx`, `WhyUsEditor.tsx`, `TypographyPreview.tsx`, `FontPicker.tsx`, `ColorPalettes.tsx`, `ContrastChecker.tsx`
  - Related theme layer update: `frontend/layers/ui-themes/src/components/ThemeProvider.tsx`
  - Validation:
    - `rg -n "style=\\{\\{" frontend -S -g "!**/*.disabled"` -> 0
    - `rg -n "style=\\{" frontend -S -g "!**/*.disabled"` -> 0

## Database/Data Migration Applied (Demo Restaurant)
- [x] Restaurant slug `demo` (`restaurant_id=8`) now has only 2 themes.
- [x] Existing theme JSON config normalized to include:
  - `header`
  - `footer`
  - `customCss`
  - normalized `components` with `displayOrder` and `config`
- [x] Theme identity domain updated to `demo.localhost:3003`.

## Remaining Hardening Tasks
- [ ] Add automated integration test: dashboard save -> public `/theme` response reflects header/footer/customCss.
  - Suggested files: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/**/__tests__/*`
  - Suggested files: `frontend/apps/restaurant-website/lib/api/get-theme-server.ts` tests
- [ ] Add E2E test for `/en/`, `/en/menu`, `/en/cart`, `/en/orders` verifying same theme header/footer.
  - Suggested files: `playwright.config.ts` and new tests under `frontend/apps/restaurant-website`
- [ ] Add validation to prevent deleting active theme unless another theme is activated first.
  - File: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/page.tsx`
- [ ] Add API guardrail to reject partial theme updates that drop required keys.
  - File: `backend/internal/handler/http/admin_theme_handler_v2.go`
- [ ] Add backend migration script for all restaurants (not only demo) to normalize legacy theme configs.
  - Suggested path: `backend/migrations/` or `backend/cmd/`
- [x] Update root docs for actual dev ports and subdomain flow (`3003` for restaurant app).
  - File: `README.md`
- [ ] Update remaining theme-system docs to reflect the same port/subdomain standards.
  - File: `docs-organized/01-THEME-SYSTEM/*.md`
- [ ] Add integration tests for header/footer parity between Theme Builder preview and storefront runtime.
  - Suggested files: `frontend/apps/restaurant-website/components/__tests__/header-ssr.test.tsx`
  - Suggested files: `frontend/apps/restaurant-website/components/__tests__/footer-ssr.test.tsx`

## Verification Checklist
- [x] Open storefront home with host header (`Host: demo.localhost:3003`) and confirm active theme header/footer/colors.
  - Verified result: dark header color `#1f2937` and title `My Restaurant` returned in HTML.
- [x] Open storefront cart route with host header and confirm same active theme header/footer.
  - Verified result: theme footer description + legal links returned in HTML.
- [ ] In dashboard Theme Builder, edit header nav item and save; refresh website and verify update.
- [ ] In dashboard Theme Builder, edit footer links and save; refresh website and verify update.
- [ ] In dashboard Theme Builder, change primary color and save; refresh website and verify branding.
- [ ] In dashboard Theme Builder, switch active theme between the two themes and verify full change.
