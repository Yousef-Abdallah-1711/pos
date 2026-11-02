# Localization and RTL/LTR (landing, restaurant-website, super-admin; dashboard summarised)

Static analysis only. IDs refer to `evidence/issues-frontend-others.json`.

| App | Locales | Mechanism | Locale source | `<html lang/dir>` | RTL implementation |
|---|---|---|---|---|---|
| landing | en, ar | `lib/translations.ts` (typed object, 437 lines) via `LanguageContext` | URL `/{locale}`; middleware forces prefix, default `en`, no Accept-Language/cookie | Root layout static `lang="en"`; `useEffect` sets `documentElement.dir/lang` after hydration | `dir={dir}` on page wrappers (20 `dir` usages), conditional classes (`dir==='rtl' ? 'left-3' : 'right-3'`), logical utilities (7 `ms/me/ps/pe`) vs 3 physical |
| restaurant-website | en, ar | `i18n/messages/{en,ar}.json` (224 keys each, identical key sets) via `createTranslator`; a separate `messages/*.json` with a different top-level key set also exists; next-intl provider holds only 3 inline keys and no `useTranslations` consumer | URL `/{locale}` (middleware), default `en` | Root layout computes locale from `params.locale`, which is unavailable there -> always `lang="en" dir="ltr"` (FE-OTH-014); inner `<div lang dir>` in `[locale]/layout.tsx` | Mostly physical Tailwind utilities: 53 matches of `ml-/mr-/pl-/pr-/left-/right-/text-left|right` vs 1 logical -> likely mirrored-layout defects in RTL; 5 `dir=`/`rtl:` usages |
| super-admin | en only | none, strings hardcoded | n/a | `lang="en"` | n/a |
| dashboard | next-intl (`src/i18n`), `[locale]` routes | per package.json/config only | URL | not audited here | not audited here |

## Specific observations
1. Server-rendered Arabic pages in restaurant-website start with `dir="ltr"` on `<html>`: scrollbar side, `<body>` default alignment, and any element outside the inner div (toasts: `Toaster position="top-right"`, portals, modals) stay LTR. Landing's Toaster is locale-aware (`top-left` for ar, `[locale]/layout.tsx:21`) but the html attribute flips only after hydration (flash/hydration mismatch risk).
2. Landing `setLanguage` replaces the first path segment (works for all routes); restaurant-website switching is via URL too. Neither persists preference; restaurant has `preferences-store` (settings page) - not tied to routing.
3. Hardcoded English in localized apps: landing Pricing section label "Pricing" (`Pricing.tsx:19`), Footer "Privacy"/"Terms"; restaurant default theme labels ('Orders', 'Home' in `lib/utils/default-theme.ts`, `get-theme-server.ts:180`, `theme-layout-normalizers.ts:82`); `use-api-queries.ts` hardcodes `?lang=en` on categories/products/search (lines 101-165) so API content is always English regardless of locale (backend support for `lang` not verified here).
4. Currency/number/date formatting: checkout uses `$` and `toFixed(2)` regardless of locale (FE-OTH-017); landing prices are translated strings ('500' with 'EGP'/'جنيه مصري'), no Intl formatting. Arabic digits not used.
5. Deep-link and redirect stubs hardcode `/en` (landing stubs, super-admin-login handoff `?/en/dashboard`, `impersonation_token` URL `/en/dashboard`), so Arabic users are redirected to English dashboard routes.
6. Translation parity: restaurant-website en/ar have identical key sets (224). Landing parity enforced only by the TS `Translations` type (not run).
7. Marketing copy in AR includes the same unverified claims (trial 30 days, 500+ restaurants); AR price string '500 جنيه مصري/شهر'.
8. RTL layout of Arabic theme sections comes from `@pos-saas/ui-themes`; not audited for logical properties.

## Recommendations
Move `<html>` to the `[locale]` layout in restaurant-website; adopt logical Tailwind utilities (`ms/me/ps/pe/start/end`) or the `rtl:` variant; centralize locale-aware currency via `Intl.NumberFormat` driven by tenant settings; pass the active locale to the `lang` query param; keep one message source in restaurant-website and remove next-intl or use it fully; add an RTL visual regression suite.
