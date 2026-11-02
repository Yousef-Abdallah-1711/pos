# Frontend architecture — current state and target (Phase 2A)

## Current state
- Four Next.js 15 apps (`landing` 3001, `dashboard` 3002, `restaurant-website` 3003, `super-admin` 3004) in a pnpm workspace. `packages/ui`, `ui-components`, `theme` and `layers/ui-themes`, `layers/shared-components`; no circular dependencies (Phase 1). `@pos-saas/ui` is imported but not declared; `packages/theme` and `layers/shared-components` have no consumers.
- Authentication is client-side in all apps (no server middleware except landing locale redirect and storefront subdomain resolution). Tokens still travel in URLs between apps (FE-OTH-005).
- Type-check split done this phase: each app's `tsconfig.json` now covers application code only; `tsconfig.test.json` adds tests with Vitest/Jest + jest-dom types. New scripts: `typecheck`, `typecheck:test` in all four apps.

## Facts measured this phase
| App | `typecheck` (app code) | `typecheck:test` (incl. tests) | Unit tests |
|---|---|---|---|
| landing | **0 errors** | 23 errors | cannot run: no Vitest config/`@/` alias; needs `jsdom`, `@vitejs/plugin-react`, testing-library as devDependencies (lockfile change) |
| super-admin | **0 errors** | 3 errors | 1 file, 3 tests **pass** |
| restaurant-website | **0 errors** | 24 errors | Jest; not run |
| dashboard | **0 errors** (was 163 app-code + 542 test = 705) | 158 errors | 41 files: 26 pass / 15 fail; 639 tests: 530 pass / 109 fail |

Root causes found in the dashboard (all verified): 542 of the 705 errors were test files type-checked by the app config (413 were missing jest-dom matcher types); 5 were in unreferenced duplicate directories at the app root (`components/`, `lib/hooks/use-theme.ts`, `i18n/`) that nothing in `src/` imports; 45 were untyped product component props; 23 were duplicate object keys; 16 were react-hook-form errors rendered as React nodes.

## Real defects found while fixing types (fixed)
1. Products pages used `useHasPermission(...)` (returns `{hasPermission, loading, error}`) as a boolean, so view/create/delete gating **never fired** (UI only; the API enforces permissions). Fixed on products list/new/edit.
2. `generatePaletteFromPrimary` called an undefined `hexToHex` (ReferenceError); now uses `rgbToHex`.
3. `PerformanceMonitor` metric cards looked up the key `'gray'` in a map without it (class `undefined`).
4. `product.price.toFixed` crashed for products without a price.
6. `useNotifications` called `apiCache.remove()`, which does not exist (the method is `clear`), so marking notifications read/deleting threw while invalidating the cache; fixed (6 sites).
7. Header/Footer/Hero/CTA/ProductGrid/Testimonials/WhyUs previews read `component.title.en` although titles can be plain strings (always showed the fallback); now go through `previewText()`.
5. Theme editor page built config objects with 31 dead duplicate property lines (the later key always won); removed with no behaviour change.

## Open defects found (not fixed)
- `NotificationBell` (unused) targets a removed low-stock data model; it was patched only to compile.
- The dashboard `TenantSelectorModal` (unused) imported a non-existent hook; it now uses `useAuthStore`.
- Two parallel theme type models (`ThemeComponent` with numeric ids vs `ComponentConfig` with string ids) cause the remaining theme-builder errors.
- Header/Footer editors treat `navigationItems[].id` as always present while the shared type makes it optional (items loaded without ids cannot be edited correctly).
- Two parallel theme models remain (cast documented in `LivePreview.tsx`).
- Dashboard `components/` and `lib/` at the app root are dead and are excluded from type-check and Vitest (documented, not deleted); `tests/visual-regression` Playwright specs were being run by Vitest and are now excluded.

## Target (decision record; not implemented)
Keep the four apps. (1) One `api-client` package replacing per-app axios/fetch wrappers with a normalised error type; (2) delete or adopt `packages/theme` and `layers/shared-components`; declare `@pos-saas/ui`; (3) move route components to thin composition + feature hooks, starting with products, theme-builder and HR where types showed the coupling; (4) server-side auth guards (Next middleware) instead of client redirects; (5) fix the Vitest setup per app so tests are a CI gate. Not done in Phase 2A: the type and test recovery had to come first and is not finished.
