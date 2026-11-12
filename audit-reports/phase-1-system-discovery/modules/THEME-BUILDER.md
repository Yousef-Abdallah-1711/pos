# Theme builder (dashboard) - static audit

Status: **Partial**. Create/list/edit/activate/duplicate/delete/history-read and presets implemented; rollback is a stub; import/export/component APIs target nonexistent routes; page-level gating absent.

## Pages (5)
`theme-builder` (671 lines; 10 api refs; list + activate/delete/duplicate), `theme-builder/editor/[id]` (605; editor with sidebar, ComponentLibrary, ThemeHistory, EditorPreview), `import` (269; creates via createTheme from pasted JSON), `presets` (115), `qr` (98; website QR). None has MODULE_IDS gating (THEMES module 9 exists and is used only in Sidebar `canViewThemes`).

## State
`stores/themeBuilderStore.ts` (zustand persist 'theme-builder-storage'; test 714 lines), `themeBuilderHooks.ts`, `themeBuilderUtils.ts`; shared types `types/theme*.ts`; `lib/validation/theme-schemas.ts`; `lib/themeLoader.ts`, `services/googleFontsService.ts`. Preview uses `@pos-saas/ui-themes` (workspace wiring issue FE-DASH-011).

## Contract (themeApi.ts uses its OWN axios instance, no tenant/impersonation headers)
Match: GET/POST `/admin/themes`, GET/PUT/DELETE `/admin/themes/{id}`, POST `/{id}/activate`, `/{id}/duplicate`, GET `/{id}/history`, GET `/admin/theme-presets` (main.go:1462-1472, `wrapProtected`, no RBAC).
Orphans (functions defined, no UI caller found): POST `/admin/themes/{id}/export`; POST `/admin/themes/import` (405: method mismatch); GET/POST `/admin/themes/{id}/components`, PUT/DELETE `/components/{cid}`, POST `/components/reorder`; GET `/public/component-library`. Handlers exist only in `backend/internal/routes/theme_routes_v2.go` under `//go:build legacy` (FE-DASH-013).
Stub: `ThemeHistory.tsx:65-82` shows success toast without calling any API (FE-DASH-012).
Sections/components are saved as part of the theme PUT payload (`updateTheme`, themeApi.ts:538-575) - whole-theme overwrite, no version/etag (last write wins; not verified server-side).

## Stray duplicates
App-root `components/theme-editor/*` and `components/component-manager/*` are an older editor, unreferenced from src, with type errors (section-editor `background_image`, theme-editor `updateTheme` from unresolved `@/lib/api-client`) (FE-DASH-010). `src/components/editors/*` (5 test files) and `src/components/theme/*` are the live editors.

## Findings
FE-DASH-010, 011, 012, 013, 014 (ungated, backend auth-only), FE-DASH-008 (second HTTP client), FE-DASH-007.

## Unverified
Publish pipeline to restaurant-website, image upload for logo (`logo-uploader.tsx` in dead dir), XSS safety of custom CSS/HTML fields.
