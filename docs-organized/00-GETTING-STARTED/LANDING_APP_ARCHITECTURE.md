# Landing App Architecture (Port 3001)

**Last Updated:** February 21, 2026

The landing app serves as the public-facing marketing site and authentication gateway for the POS SaaS platform.

---

## i18n Architecture

The landing app uses **URL-based i18n routing** — all pages live under `/en/...` or `/ar/...`. This differs from earlier versions that used localStorage/context switching.

### How It Works

```
User visits http://localhost:3001/
  ↓ middleware.ts
  → Redirect to /en (default locale)

User visits http://localhost:3001/ar/auth/login
  ↓ middleware.ts
  → Locale "ar" recognized — serve normally
  ↓ [locale]/layout.tsx (server)
  → Pass locale="ar" to LanguageProvider
  ↓ LanguageProvider (client)
  → Set document.documentElement.dir="rtl"
  → Provide translations["ar"] via context
```

### Middleware (`src/middleware.ts`)

```typescript
const SUPPORTED_LOCALES = ['en', 'ar']
const DEFAULT_LOCALE = 'en'

// Skips: /_next/*, /api/*, /favicon.ico, files with extension
// Passes through: paths starting with a supported locale
// Redirects: everything else to /en/...
```

### File Structure

```
src/
├── middleware.ts                    Locale detection + redirect
├── app/
│   ├── layout.tsx                  Root layout (ThemeProvider, html/body)
│   ├── page.tsx                    redirect('/en')
│   ├── [locale]/
│   │   ├── layout.tsx              Locale layout (LanguageProvider + Toaster)
│   │   ├── page.tsx                Landing page (all sections)
│   │   └── auth/
│   │       ├── login/page.tsx      Login (uses locale from URL)
│   │       ├── register/page.tsx   Register (uses locale from URL)
│   │       └── super-admin-login/  Super admin login
│   ├── auth/                       Legacy redirects → /en/auth/...
│   ├── login/page.tsx              Legacy redirect → /en/auth/login
│   └── register/page.tsx           Legacy redirect → /en/auth/register
├── components/
│   ├── Navbar.tsx                  Sticky glassmorphism, language toggle
│   ├── Hero.tsx                    Dark POS mockup, stats, CTA
│   ├── Features.tsx                6 feature cards with gradient icons
│   ├── HowItWorks.tsx              3-step process section (Sign Up → Set Up → Sell)
│   ├── Pricing.tsx                 Single plan: 500 EGP/month
│   ├── Testimonials.tsx            3 testimonials with star ratings
│   ├── CTA.tsx                     Red gradient call-to-action
│   ├── Footer.tsx                  4-column, all links translated
│   └── TenantSelectorModal.tsx     Multi-tenant login selection
├── contexts/
│   └── LanguageContext.tsx         URL-based locale, router.push for toggle
├── lib/
│   ├── translations.ts             EN + AR inline translations
│   └── config.ts                   Dashboard URL config
└── stores/
    └── authStore.ts                Zustand auth (email, token)
```

---

## Language Toggle

The language toggle swaps the first URL segment:

```typescript
// LanguageContext.tsx
const setLanguage = (lang: Language) => {
  const segments = pathname.split('/')
  segments[1] = lang  // swap /en/ ↔ /ar/
  router.push(segments.join('/'))
}
```

Clicking "العربية" on `/en/auth/login` → navigates to `/ar/auth/login`.

---

## Pricing Plan

Single plan (as of Feb 2026):

| | Details |
|---|---|
| **Trial** | 30 days free, no credit card required |
| **Price** | 500 EGP/month after trial |
| **Features** | All features included (no tiers) |

Features included: Full POS System, Kitchen Display System, Online Ordering Website, Analytics & Reports, Multi-Restaurant Management, Mobile Apps, Unlimited Staff, Priority Support, Free Updates.

---

## Auth Flow

### Login
1. User enters email + password at `/en/auth/login`
2. `POST /api/v1/auth/login`
3. If single tenant → receive JWT → redirect to Dashboard auth callback
4. If multiple tenants → show `TenantSelectorModal`
5. User selects tenant → `POST /api/v1/auth/login/confirm` → receive JWT
6. `window.location.href = ${dashboardUrl}/auth/callback?token=...&user=...`

### Register
1. User fills form at `/en/auth/register`
2. `POST /api/v1/auth/register` — creates Tenant + Restaurant + User (owner role)
3. JWT returned → redirect to Dashboard auth callback

---

## RTL Support

Arabic RTL is applied client-side via `useEffect` in `LanguageProvider`:

```typescript
useEffect(() => {
  document.documentElement.dir = dir   // 'rtl' or 'ltr'
  document.documentElement.lang = locale
}, [locale, dir])
```

The root `<html>` tag has `suppressHydrationWarning` to prevent hydration mismatch since `dir` is set after server render.

Components use:
- `dir={dir}` prop on containers
- `text-${dir === 'rtl' ? 'right' : 'left'}` for text alignment
- `${dir === 'rtl' ? 'mr-2' : 'ml-2'}` for conditional spacing

---

## Translations

All translations are in `src/lib/translations.ts` as a typed TypeScript object:

```typescript
const translations = {
  en: {
    nav: { ... },
    hero: { title, subtitle, badge, cta1, cta2, stats: { restaurants, orders, uptime } },
    features: { title, subtitle, items: [...] },
    howItWorks: { title, subtitle, step1, step2, step3 },
    pricing: { title, subtitle, trialBadge, price: '500', currency: 'EGP', features: [...] },
    testimonials: { ... },
    cta: { ... },
    footer: { ... },
    auth: { login: {...}, register: {...}, tenantSelector: {...} },
  },
  ar: { /* full Arabic parity */ }
}
```

---

## Build Notes

- `output: 'standalone'` is set in `next.config.js` for Docker packaging
- On Windows, standalone build shows EPERM symlink warnings — this is a Windows limitation, not a code error, and does not affect dev or production functionality
- `generateStaticParams()` in `[locale]/layout.tsx` pre-generates `/en` and `/ar` at build time
