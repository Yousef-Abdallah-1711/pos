# Restaurant Website UI Redesign — Master Overview v3.0
**Project:** POS SaaS Egyptian Platform — Restaurant Storefront (`:3003`)
**Date:** 2026-03-03
**Quality Target:** Production-ready Egyptian SaaS (500 EGP/month clients)
**App Path:** `frontend/apps/restaurant-website/` (rebuild in-place, no v2 folder needed)
**MVP Scope:** 15 working days (Phases 00–04) → Deploy → Iterate remaining phases

---

## 0. MVP DEFINITION (Ship in 15 Days)

### MVP Scope (Phases 00–04 Only)
| Day | Milestone | Deliverable |
|-----|-----------|-------------|
| 1–2 | Phase 00: Setup | shadcn/ui, CSS vars, types, config, backup |
| 3–5 | Phase 01: Theme System V2 | Dark mode per theme, CSS vars, ThemeProvider rewrite |
| 6–8 | Phase 02: Header + Footer | 5 header variants, 5 footer variants, mobile nav |
| 9–11 | Phase 03: Homepage | 8 sections with theme-specific templates |
| 12–14 | Phase 04: Commerce | ProductCard, AJAX CartSheet, OrdersTable |
| 15 | Deploy | Lighthouse audit, smoke test, production push |

### Post-MVP (Phases 05–08, iterate after deploy)
| Phase | Scope | When |
|-------|-------|------|
| 05 | All page redesigns (Menu, Cart, Checkout, Orders, Settings) | Week 4–5 |
| 06 | Cross-cutting (Search, Cart sync, WhatsApp forwarding) | Week 5–6 |
| 07 | Performance + Testing (Lighthouse 95+, E2E, PWA) | Week 6–7 |
| 08 | SEO + Monitoring + Production hardening | Week 7–8 |

### Stop-Ship Criteria
- Header dark mode broken → BLOCK deploy
- Cart add-to-cart not working → BLOCK deploy
- RTL layout broken on any page → BLOCK deploy
- Lighthouse mobile < 70 → BLOCK deploy
- Any hardcoded hex color in components → BLOCK deploy

---

## 1. RISK MATRIX

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| shadcn/ui RTL bugs | High | High | Test all components in Phase 00 sandbox before use |
| Bundle bloat (shadcn 20+ components) | Medium | High | Weekly `@next/bundle-analyzer`, hard budget 180KB gzip |
| Dark mode flash on slow networks | High | Medium | `next-themes` with `disableTransitionOnChange` + `suppressHydrationWarning` |
| Embla Carousel loop mode bugs | Medium | Low | Disable loop, use `containScroll: 'trimSnaps'` |
| Cairo font 200KB on English pages | High | Medium | Conditional load: only when `locale === 'ar'` |
| Cache stale data after admin edit | High | High | Cache invalidation on write + `stale-while-revalidate` headers |
| Scope creep beyond 15-day MVP | High | Critical | Hard cutoff: only Phases 00–04 for MVP deploy |

### Rollback Strategy
```bash
git checkout restaurant-ui-redesign-backup
# Or revert specific commits:
git revert <commit-range>
```

### Checkpoint Schedule
| Day | Checkpoint | Pass Criteria |
|-----|-----------|---------------|
| 3 | Header renders across all 10 themes (light + dark) | No hardcoded hex, RTL works |
| 7 | Homepage renders 8 sections, theme templates work | Each theme has unique section order |
| 10 | Menu page + ProductCard + AJAX cart working | Add to cart → sheet opens, no reload |
| 12 | Cart + Checkout MVP functional | 3-step checkout, cash payment works |
| 15 | Production deploy | Lighthouse mobile ≥ 70, all pages render |

---

## 2. CURRENT STATE AUDIT

### What Exists (Keep / Migrate)
| Layer | File | Status | Action |
|---|---|---|---|
| Routing | `middleware.ts` | ✅ Works | Enhance + keep |
| Pages | `app/[locale]/page.tsx` (home) | ✅ Works | Rebuild UI |
| Pages | `app/[locale]/menu/` | ✅ Works | Full redesign |
| Pages | `app/[locale]/cart/` | ✅ Works | Full redesign |
| Pages | `app/[locale]/checkout/` | ✅ Works | Full redesign |
| Pages | `app/[locale]/orders/` | ✅ Works | Full redesign |
| Pages | `app/[locale]/orders/[id]/` | ✅ Works | Full redesign |
| Pages | `app/[locale]/orders/track/` | ✅ Works | Enhance |
| Pages | `app/[locale]/settings/` | ✅ Works | Full redesign |
| Pages | `app/[locale]/table/[tableId]/menu/` | ✅ Works | Enhance |
| Pages | `app/maintenance/` | ✅ Works | Polish |
| State | `lib/store/cart-store.ts` | ✅ Works | Keep + extend (cross-tab sync, AJAX) |
| State | `lib/store/theme-store.ts` | ✅ Works | Keep + extend |
| State | `lib/store/preferences-store.ts` | ✅ Works | Keep + extend (dark mode) |
| API | `lib/api/api-client.ts` | ✅ Works | Keep |
| API | `lib/api/get-theme-server.ts` | ✅ Works | Keep |
| API | `lib/api/theme-api.ts` | ✅ Works | Keep |
| Hooks | `lib/hooks/use-api-queries.ts` | ✅ Works | Keep + extend |
| Hooks | `lib/hooks/use-theme.ts` | ✅ Works | Keep |
| i18n | `i18n/messages/en.json` + `ar.json` | ✅ Works | Extend keys |
| Themes layer | `layers/ui-themes/src/` | ✅ Works | Extend sections + variants |

### Critical Gaps (Must Fix)
| Gap | Severity | Root Cause | Fix |
|---|---|---|---|
| Header has no dark mode | 🔴 Critical | Inline hex in `Header.tsx` `<style>` block | CSS vars + `useTheme` hook |
| No language switcher in header | 🔴 Critical | `showLanguageSwitcher` config exists, not rendered | Add `LanguageToggle` component |
| No cart icon/count in header | 🔴 Critical | `Header.tsx` has no cart slot | Add cart slot to `HeaderConfig` |
| No search in header | 🟠 High | Config field exists, not rendered | Add live search component |
| Footer ignores dark mode | 🟠 High | Same inline hex pattern | CSS vars |
| Only 7 homepage section types | 🟡 Medium | Registry missing Gallery, Categories, Menu Preview | Add 2 new types |
| Menu page is basic grid | 🟡 Medium | No category tabs, no price filter, no sticky bar | Full shadcn redesign |
| Cart is page-only (no AJAX) | 🔴 Critical | Add to cart = full page reload | AJAX cart + CartSheet sidebar |
| No mobile bottom nav | 🟡 Medium | Missing entirely | 5-item mobile nav |
| No sticky cart bottom bar | 🟡 Medium | Missing from all pages | Persistent bottom bar |
| All 10 themes look identical | 🔴 Critical | Same sections, same order, only color swap | Theme templates + variants |
| No WhatsApp floating button | 🟠 High | WhatsApp = checkout in Egypt | Always-visible green bubble |
| No cache invalidation | 🟠 High | Stale data after admin edits | Invalidate on write |
| No dark color columns in DB | 🔴 Critical | themes_v2 has no dark variants | Migration + seed |
| Header missing features | 🔴 Critical | No search, cart, language, dark toggle | Add 5 sub-components |
| Settings minimal | 🟢 Low | Only language preference | Full tabs (lang/theme/notifs/cart) |

---

## 3. TARGET ARCHITECTURE

### App Structure (In-Place Rebuild)
```
frontend/apps/restaurant-website/
├── app/
│   ├── layout.tsx                    # Root layout — providers, fonts, CSS vars
│   ├── page.tsx                      # Root → redirect to /en
│   ├── sitemap.ts
│   ├── robots.ts
│   ├── api/metrics/route.ts
│   ├── maintenance/page.tsx
│   └── [locale]/
│       ├── layout.tsx                # Locale layout — dir, lang, NextIntl, Header+Footer
│       ├── page.tsx                  # Home (SSR) — theme sections
│       ├── loading.tsx               # Homepage skeleton
│       ├── error.tsx
│       ├── not-found.tsx
│       ├── menu/
│       │   ├── page.tsx
│       │   ├── loading.tsx           # Menu skeleton grid
│       │   └── menu-page-content.tsx
│       ├── cart/
│       │   ├── page.tsx
│       │   └── cart-page-content.tsx
│       ├── checkout/
│       │   ├── page.tsx
│       │   └── checkout-page-content.tsx
│       ├── orders/
│       │   ├── page.tsx
│       │   ├── loading.tsx
│       │   ├── orders-page-content.tsx
│       │   ├── [id]/
│       │   │   ├── page.tsx
│       │   │   └── order-detail-content.tsx
│       │   └── track/
│       │       ├── page.tsx
│       │       ├── track-page-content.tsx
│       │       └── [orderNumber]/page.tsx
│       ├── settings/
│       │   ├── page.tsx
│       │   └── settings-page-content.tsx
│       └── table/
│           └── [tableId]/menu/
│               └── page.tsx
│
├── components/
│   ├── ui/                           # shadcn/ui components (auto-generated)
│   │   ├── button.tsx, card.tsx, input.tsx, form.tsx
│   │   ├── dialog.tsx, drawer.tsx, sheet.tsx
│   │   ├── badge.tsx, tabs.tsx, avatar.tsx, progress.tsx
│   │   ├── separator.tsx, slider.tsx, tooltip.tsx
│   │   ├── dropdown-menu.tsx, scroll-area.tsx, skeleton.tsx
│   │   ├── table.tsx, textarea.tsx, accordion.tsx
│   │   ├── popover.tsx, command.tsx, toggle-group.tsx
│   │   └── radio-group.tsx, alert-dialog.tsx
│   │
│   ├── layout/
│   │   ├── Header.tsx                # Main header — dispatches to variant
│   │   ├── Footer.tsx                # Main footer — dispatches to variant
│   │   ├── MobileBottomNav.tsx       # 5-item sticky mobile nav
│   │   ├── CartBottomBar.tsx         # Sticky cart total bar
│   │   ├── WhatsAppFloating.tsx      # Floating green WhatsApp bubble
│   │   ├── PageWrapper.tsx           # Header + Footer + children wrapper
│   │   ├── headers/                  # 5 header variants
│   │   │   ├── GlassHeader.tsx       # Backdrop-blur transparent (Modern Bistro)
│   │   │   ├── GradientHeader.tsx    # Linear gradient primary→accent (Luxury Gold)
│   │   │   ├── SolidHeader.tsx       # Solid background color (default)
│   │   │   ├── FloatingHeader.tsx    # Fixed with heavy shadow (Urban Gray)
│   │   │   └── MinimalHeader.tsx     # Thin border-bottom only (Minimal White)
│   │   ├── footers/                  # 5 footer variants
│   │   │   ├── ExpandedFooter.tsx    # 4-column full footer
│   │   │   ├── CompactFooter.tsx     # 2-column compact
│   │   │   ├── CardsFooter.tsx       # Card-based layout (Luxury Gold)
│   │   │   ├── SocialHeavyFooter.tsx # Social links prominent (Vibrant Pop)
│   │   │   └── MinimalFooter.tsx     # Single line copyright
│   │   └── header/                   # Shared header sub-components
│   │       ├── SearchBar.tsx         # Live search with debounce
│   │       ├── CartButton.tsx        # Cart icon + count badge
│   │       ├── LanguageToggle.tsx    # EN ↔ AR URL-based switcher
│   │       ├── DarkModeToggle.tsx    # Dark/Light/System toggle
│   │       └── WhatsAppButton.tsx    # Green WhatsApp order button
│   │
│   ├── sections/                     # Homepage sections (used by SectionRenderer)
│   │   ├── HeroSection.tsx           # Dispatches to hero variant
│   │   ├── heroes/                   # 4 hero variants
│   │   │   ├── GradientHero.tsx      # CSS gradient background
│   │   │   ├── VideoHero.tsx         # Video background (Modern Bistro)
│   │   │   ├── SliderHero.tsx        # Embla image slider
│   │   │   └── ParallaxHero.tsx      # Parallax scroll effect (Warm Rustic)
│   │   ├── FeaturedProductsSection.tsx  # Dispatches to products variant
│   │   ├── products/                 # 3 product layout variants
│   │   │   ├── ProductCarousel.tsx   # Embla carousel (default)
│   │   │   ├── ProductGrid.tsx       # 3-col grid
│   │   │   └── ProductMasonry.tsx    # Masonry layout
│   │   ├── WhyChooseUsSection.tsx
│   │   ├── TestimonialsSection.tsx
│   │   ├── CategoriesSection.tsx     # NEW
│   │   ├── GallerySection.tsx        # NEW
│   │   ├── CTASection.tsx
│   │   └── ContactSection.tsx
│   │
│   ├── commerce/
│   │   ├── ProductCard.tsx           # Image, name AR/EN, price, qty stepper, WhatsApp
│   │   ├── ProductCardSkeleton.tsx   # Loading skeleton
│   │   ├── CartSheet.tsx             # AJAX slide-in cart (Sheet component)
│   │   ├── CartItem.tsx              # Single cart item with qty controls
│   │   ├── OrderStatusBadge.tsx      # Status badge variants
│   │   ├── OrderTimeline.tsx         # Vertical status timeline
│   │   └── OrdersTable.tsx           # shadcn DataTable for orders list
│   │
│   ├── forms/
│   │   ├── CheckoutForm.tsx          # RHF + Zod validated form (3 payment methods)
│   │   └── TrackOrderForm.tsx        # Order number input form
│   │
│   ├── providers/
│   │   ├── ThemeProvider.tsx          # CSS vars injection + next-themes (REWRITE)
│   │   ├── QueryProvider.tsx          # TanStack React Query
│   │   └── WebVitalsProvider.tsx
│   │
│   └── system/
│       ├── MaintenancePage.tsx
│       ├── ErrorPage.tsx
│       ├── NotFoundPage.tsx
│       ├── NetworkErrorState.tsx       # No internet fallback
│       ├── RateLimitErrorState.tsx     # 429 with countdown
│       ├── PaymentErrorState.tsx       # Payment failure + retry
│       ├── ProductUnavailableState.tsx # Out of stock overlay
│       └── GlobalErrorBoundary.tsx     # Catches all unhandled errors
│
├── lib/
│   ├── api/                          # Keep all existing
│   ├── store/
│   │   ├── cart-store.ts             # Keep + add AJAX, cross-tab sync, inventory check
│   │   ├── theme-store.ts            # Keep
│   │   ├── preferences-store.ts      # Keep + add dark mode
│   │   └── search-store.ts           # NEW: search state
│   ├── hooks/
│   │   ├── use-api-queries.ts        # Keep + add search queries
│   │   ├── use-theme.ts              # Keep
│   │   ├── use-dark-mode.ts          # NEW: dark mode hook
│   │   ├── use-cart.ts               # NEW: AJAX cart convenience hook
│   │   └── use-search.ts             # NEW: debounced search hook
│   ├── utils/
│   │   ├── theme-colors.ts           # Keep
│   │   ├── default-theme.ts          # Keep + extend with dark pairs + all 10 defaults
│   │   ├── hex-to-hsl.ts             # NEW: color conversion utility
│   │   ├── theme-normalizer.ts       # NEW: API → ThemeConfig normalizer
│   │   ├── theme-templates.ts        # NEW: 10 template definitions
│   │   └── format.ts                 # NEW: EGP price, Arabic numerals, date format
│   └── config/
│       └── api.ts                    # Keep
│
├── i18n/
│   ├── config.ts
│   └── messages/
│       ├── en.json                   # Extended with all new keys (~200 keys)
│       └── ar.json                   # Extended with all new keys (~200 keys)
│
├── styles/
│   ├── globals.css                   # Tailwind base + dark mode .dark class
│   ├── theme-variables.css           # CSS vars — ALL theme colors via CSS vars
│   └── rtl.css                       # RTL overrides
│
├── types/
│   ├── theme.ts                      # Full TypeScript theme types + variants
│   └── payment.ts                    # Payment method types
│
└── middleware.ts                      # Keep + enhance
```

---

## 4. TECHNOLOGY STACK

| Tool | Version | Role |
|---|---|---|
| Next.js | 15.1.0 | App Router, SSR, ISR |
| React | 18.3.1 | UI |
| TypeScript | 5.7.2 | Type safety |
| TailwindCSS | 3.4.17 | Utility classes + dark mode `class` strategy |
| shadcn/ui | latest | Component library (RTL-capable, CSS var-based) |
| Radix UI | latest | Headless primitives |
| Zustand | 4.5.0 | Client state (cart, prefs, search) |
| TanStack Query | 5.59.20 | Server state + cache (10min stale, 1hr gc) |
| React Hook Form | 7.54.2 | Forms |
| Zod | 3.24.1 | Schema validation |
| Embla Carousel | 8.x | Carousels (7KB, no jQuery) |
| next-themes | 0.4.4 | Dark mode (html class strategy) |
| next-intl | 4.6.1 | i18n |
| Lucide React | 0.460.0 | Icons |
| Framer Motion | 11.x | Page transitions (lazy loaded, optional) |
| Sonner | 1.7.1 | Toast notifications |
| @dnd-kit | latest | Drag & drop for theme builder sections |

---

## 5. THEME DIFFERENTIATION SYSTEM

### The Problem (Current State)
All 10 themes use the **same 8 sections in the same order**, only changing colors. Result: every homepage is "the same kid in different clothes" — 10% visual variety.

### The Solution: Theme Templates + Component Variants
Each theme gets a **template** that defines: section order, header variant, hero variant, product layout, and footer variant.

```
ThemeConfig → templateId → Component Variants → Unique Layout

"luxury-gold" → "template-02" → GradientHeader + DarkOverlayHero + GoldCards
"modern-bistro" → "template-01" → GlassHeader + VideoHero + Carousel
```

### 10 Themes — Full Breakdown

| # | Theme | Header | Hero | Products | Footer | Unique Feature |
|---|-------|--------|------|----------|--------|----------------|
| 1 | Modern Bistro | Glass (blur) | Video BG | Carousel | Expanded | Live order counter |
| 2 | Luxury Gold | Gradient (gold) | Dark overlay | Grid (gold cards) | Cards | Reservation badge |
| 3 | Warm Rustic | Solid (warm) | Parallax | Masonry | Compact | Ingredient badges |
| 4 | Coastal Ocean | Glass (blue tint) | Gradient (ocean) | Carousel | Social-heavy | Freshness timer |
| 5 | Garden Green | Solid (green) | Slider | Grid | Minimal | Health badges |
| 6 | Spicy Red | Solid (red) | Gradient (fire) | Carousel | Expanded | Spice level badges |
| 7 | Urban Gray | Minimal (line) | Gradient (split) | Grid | Compact | Stats carousel |
| 8 | Minimal White | Minimal (invisible) | Gradient (clean) | Grid (single col) | Minimal | Quote overlay |
| 9 | Vibrant Pop | Floating (shadow) | Slider (animated) | Carousel | Social-heavy | Pop-up testimonials |
| 10 | Egyptian Heritage | Solid (gold accent) | Gradient (pharaoh) | Grid (pyramid) | Expanded | Gold accents + Arabic calligraphy |

### Section Order Per Theme (Different Homepages)
```
Modern Bistro:      Hero → Featured → WhyUs → Testimonials → CTA → Contact
Luxury Gold:        Hero → Gallery → Testimonials → Featured → CTA → Contact
Warm Rustic:        Hero → Categories → Featured → Gallery → Contact → CTA
Coastal Ocean:      Hero → Featured → Categories → WhyUs → CTA → Contact
Garden Green:       Hero → Categories → Featured → WhyUs → Testimonials → CTA
Spicy Red:          Hero → Featured → Categories → Testimonials → CTA → Contact
Urban Gray:         Hero → Featured → WhyUs → CTA → Contact
Minimal White:      Hero → Featured → CTA
Vibrant Pop:        Hero → Categories → Featured → Testimonials → Gallery → CTA
Egyptian Heritage:  Hero → Categories → Gallery → Featured → WhyUs → CTA → Contact
```

### Result: 10 themes × 5 headers × 4 heroes × 3 products × 5 footers = **300 possible variations**

---

## 6. THEME SYSTEM V2

### 10 Themes — Color Identities

| # | Name | Slug | Primary | Accent | Category |
|---|---|---|---|---|---|
| 1 | Modern Bistro | `modern-bistro` | `#2563eb` blue-600 | `#f97316` orange-500 | professional |
| 2 | Luxury Gold | `luxury-gold` | `#b45309` amber-700 | `#1c1917` stone-900 | luxury |
| 3 | Warm Rustic | `warm-rustic` | `#ea580c` orange-600 | `#78350f` amber-900 | casual |
| 4 | Coastal Ocean | `coastal-ocean` | `#0891b2` cyan-600 | `#0c4a6e` sky-900 | modern |
| 5 | Garden Green | `garden-green` | `#059669` emerald-600 | `#14532d` green-900 | modern |
| 6 | Spicy Red | `spicy-red` | `#dc2626` red-600 | `#7f1d1d` red-900 | casual |
| 7 | Urban Gray | `urban-gray` | `#475569` slate-600 | `#0f172a` slate-900 | modern |
| 8 | Minimal White | `minimal-white` | `#111827` gray-900 | `#6b7280` gray-500 | modern |
| 9 | Vibrant Pop | `vibrant-pop` | `#7c3aed` violet-600 | `#ec4899` pink-500 | playful |
| 10 | Egyptian Heritage | `egyptian-heritage` | `#4338ca` indigo-700 | `#d97706` amber-600 | luxury |

### Theme TypeScript Interface
```typescript
// types/theme.ts — THE source of truth

export interface ThemeConfig {
  // Identity
  id: number;
  slug: string;
  name: { en: string; ar: string };
  description: { en: string; ar: string };
  category: 'professional' | 'luxury' | 'casual' | 'modern' | 'playful';
  version: number;

  // Template (defines layout + variants)
  templateId: string;                    // "template-01" through "template-10"
  headerVariant: 'glass' | 'gradient' | 'solid' | 'floating' | 'minimal';
  footerVariant: 'expanded' | 'compact' | 'cards' | 'social-heavy' | 'minimal';
  componentOverrides: {
    hero: 'gradient' | 'video' | 'slider' | 'parallax';
    products: 'carousel' | 'grid' | 'masonry';
    testimonials: 'carousel' | 'grid';
  };

  // Colors (ALL have light + dark pairs)
  colors: {
    primary: string;
    primaryHover: string;
    primaryLight: string;
    secondary: string;
    accent: string;
    background: string;
    backgroundDark: string;
    surface: string;
    surfaceDark: string;
    surfaceElevated: string;
    surfaceElevatedDark: string;
    text: string;
    textDark: string;
    textMuted: string;
    textMutedDark: string;
    border: string;
    borderDark: string;
    success: string;
    warning: string;
    error: string;
    info: string;
  };

  // Typography
  typography: {
    fontFamily: string;          // "Inter, system-ui, sans-serif"
    fontFamilyAr: string;        // "Cairo, Tajawal, sans-serif"
    baseFontSize: number;        // 15 | 16
    lineHeight: number;          // 1.5 | 1.6
    borderRadius: number;        // 6 | 8 | 12 | 16
    fontWeightNormal: number;    // 400
    fontWeightMedium: number;    // 500
    fontWeightBold: number;      // 700
  };

  // Header
  header: ThemeHeaderConfig;

  // Footer
  footer: ThemeFooterConfig;

  // Homepage sections (ordered, toggleable — order differs per template)
  sections: ThemeSectionConfig[];

  // Dark mode
  darkMode: {
    enabled: boolean;
    strategy: 'class';
  };

  // SEO
  seo: {
    ogImage?: string;
    metaDescription?: { en: string; ar: string };
    keywords?: string[];
  };
}

export interface ThemeHeaderConfig {
  backgroundColor: string;
  backgroundColorDark: string;
  textColor: string;
  textColorDark: string;
  height: number;
  stickyHeader: boolean;
  showShadow: boolean;
  showSearch: boolean;
  showLanguageSwitcher: boolean;
  showCartButton: boolean;
  showDarkModeToggle: boolean;
  showWhatsAppButton: boolean;
  whatsAppNumber?: string;        // "201xxxxxxxxx"
  logoUrl?: string;
  logoText?: string;
  logoHeight?: number;
  navPosition: 'left' | 'center' | 'right';
  navigationItems: ThemeNavItem[];
}

export interface ThemeFooterConfig {
  layout: 'expanded' | 'compact' | 'cards' | 'social-heavy' | 'minimal';
  backgroundColor: string;
  backgroundColorDark: string;
  textColor: string;
  textColorDark: string;
  linkColor: string;
  linkColorDark: string;
  showBackToTop: boolean;
  showSocialLinks: boolean;
  showLegalLinks: boolean;
  columns: 2 | 3 | 4;
  companyInfo: {
    name: { en: string; ar: string };
    description: { en: string; ar: string };
    phone?: string;
    email?: string;
    address?: { en: string; ar: string };
  };
  socialLinks: ThemeSocialLink[];
  footerSections: ThemeFooterSection[];
  legalLinks: ThemeLegalLink[];
}

export type SectionType =
  | 'hero'
  | 'featured_products'
  | 'why_choose_us'
  | 'testimonials'
  | 'categories'
  | 'gallery'
  | 'cta'
  | 'contact';
```

### CSS Variables Strategy (Dark Mode Fix)
```css
/* All theme colors MUST be CSS variables — NEVER hardcoded hex in components */
:root {
  --theme-primary: 221 83% 53%;
  --theme-primary-hover: 221 83% 45%;
  --theme-primary-light: 214 95% 93%;
  --theme-secondary: 215 16% 47%;
  --theme-accent: 25 95% 53%;
  --theme-bg: 0 0% 100%;
  --theme-surface: 210 40% 98%;
  --theme-surface-elevated: 0 0% 100%;
  --theme-text: 222 47% 11%;
  --theme-text-muted: 215 16% 47%;
  --theme-border: 214 32% 91%;
  --theme-header-bg: 221 83% 53%;
  --theme-header-text: 0 0% 100%;
  --theme-footer-bg: 222 47% 11%;
  --theme-footer-text: 210 40% 96%;
  --theme-success: 142 76% 36%;
  --theme-warning: 38 92% 50%;
  --theme-error: 0 84% 60%;
}

.dark {
  --theme-bg: 222 47% 11%;
  --theme-surface: 217 33% 17%;
  --theme-surface-elevated: 215 25% 27%;
  --theme-text: 214 32% 91%;
  --theme-text-muted: 215 20% 65%;
  --theme-border: 217 33% 17%;
  --theme-header-bg: 222 47% 11%;
  --theme-header-text: 210 40% 98%;
  --theme-footer-bg: 224 71% 4%;
  --theme-footer-text: 210 40% 96%;
}

/* shadcn/ui variable bridge */
:root {
  --background: var(--theme-bg);
  --foreground: var(--theme-text);
  --card: var(--theme-surface);
  --card-foreground: var(--theme-text);
  --primary: var(--theme-primary);
  --primary-foreground: 0 0% 100%;
  --secondary: var(--theme-secondary);
  --muted: var(--theme-surface);
  --muted-foreground: var(--theme-text-muted);
  --accent: var(--theme-accent);
  --border: var(--theme-border);
  --input: var(--theme-border);
  --ring: var(--theme-primary);
  --destructive: var(--theme-error);
  --radius: 0.5rem;
}
```

### Scalability Rules (MUST follow in every component)
```
RULE 1: ZERO hardcoded colors — every color is hsl(var(--theme-*))
RULE 2: ZERO theme-specific if/switch — components are theme-agnostic
RULE 3: Theme data is JSON from API — not code constants
RULE 4: CSS vars injected ONCE in ThemeProvider — components just consume
RULE 5: Dark mode is .dark class on <html> — overrides same CSS vars
RULE 6: Fonts loaded via next/font — fontFamily from theme config
RULE 7: Header/Footer variant selected by config — not hardcoded
RULE 8: Section order from template + API array — rendered with .map()
RULE 9: All text bilingual {en, ar} — resolved by locale param
RULE 10: Cache at every layer: React Query → ISR → CDN → Backend sync.Map
```

---

## 7. EGYPTIAN MARKET FEATURES

### WhatsApp Dominance (Non-Negotiable)
In Egypt, **WhatsApp = checkout**. A huge share of orders come via WhatsApp.

**Required implementations:**
1. **Floating WhatsApp bubble** — always visible bottom-right, pulse animation
2. **WhatsApp button in header** — green icon, opens pre-filled order message
3. **WhatsApp button on ProductCard** — `wa.me/{number}?text=أريد أطلب ${product.name_ar} × ${qty}`
4. **WhatsApp CTA in Hero** — dual CTA: 🟢 اطلب واتساب + 🟠 عربة التسوق
5. **WhatsApp order forwarding** — on checkout success, auto-forward order details to admin WhatsApp

### Arabic UX Micro-Improvements
- **Arabic numerals consistently:** `١٤٥ج.م` not `145 EGP`
- **Bilingual items:** `شاورما دجاج / Chicken Shawarma`
- **Stars format:** `٤.٨⭐ (١٢٣ تقييم)`
- **Cart label:** `سلة التسوق (٣)`
- **Price format:** `٢٤.٥ ج.م`

### Egyptian Food Visuals
- High-res Egyptian dishes as hero backgrounds (shawarma, koshary, hawawshi)
- Delivery timer badge: `🚚 توصيل في ٣٠ دقيقة`
- Color psychology: Orange/Red = appetite trigger, Pharaoh Gold accent

### AJAX Cart (Critical for Conversion)
- Add to cart = **NO page reload** → sidebar CartSheet slides in
- Live quantity `[+] [-]` Delete controls
- Real-time total recalculation
- Persist across tabs (localStorage + BroadcastChannel)
- Inventory check (disable if out of stock, show toast at limit)

### Admin Commerce Settings (Dashboard)
```
Restaurant Settings → Commerce Tab:
- WhatsApp Orders: ☑️ Forward to [201xxxxxxxxx]
- Delivery Fee: 0 | 10 | 20 | 30 | 50 EGP
- Minimum Order: 100 EGP
- Payment Methods: ☑️ Cash | ☑️ WhatsApp | ☑️ Paymob
- Working Hours: Mon-Thu 10am-11pm, Fri-Sat 10am-1am
```

---

## 8. PERFORMANCE & CACHING STRATEGY

### 4-Layer Cache Architecture
```
User Browser
├── React Query Cache (staleTime: 10min)
│   └── Products, categories, theme cached in memory
│       → 0ms for repeat navigation in same session
│
├── Service Worker Cache (PWA)
│   └── Static assets + theme JSON cached offline
│
CDN Edge (Vercel/Cloudflare)
├── ISR Cache (revalidate: 300s)
│   └── Homepage HTML pre-rendered at edge → 50ms globally
│
├── Cache-Control Headers (max-age=300)
│   └── Theme + Products API cached at CDN edge
│
Backend (Go)
├── In-Memory Cache (sync.Map, TTL: 5min)
│   └── Theme + products cached after first DB query
│
└── PostgreSQL (last resort)
    └── Only hit when all caches expire simultaneously
```

### Cache Invalidation Strategy (CRITICAL)
```
When admin saves theme:
  1. Backend: themeCache.data.Delete("theme:" + slug)
  2. Backend: Set Cache-Control: no-cache on response
  3. Frontend: queryClient.invalidateQueries(['theme', slug])

When admin updates products:
  1. Backend: productCache.data.Delete("products:" + restaurantId)
  2. Frontend: queryClient.invalidateQueries(['products', restaurantId])

Vercel ISR: POST revalidate webhook on admin save
```

### Performance Budget
| Metric | Target | How |
|--------|--------|-----|
| Theme API response | < 50ms | Backend sync.Map cache + CDN |
| CSS var injection | < 5ms | Single documentElement pass |
| Dark mode toggle | < 16ms | CSS class swap only |
| Homepage SSR | < 200ms | ISR + parallel data fetch |
| LCP (mobile 4G Egypt) | < 2.5s | next/image + ISR + font optimization |
| JS bundle (gzip) | < 180KB | Code splitting + tree-shake + dynamic imports |
| Total page weight | < 400KB | next/image + font subsetting |
| Lighthouse mobile | ≥ 85 (MVP), ≥ 95 (final) | Mobile-first everything |
| API calls per session | 3–5 (first visit), 0 (repeat) | React Query cache |

### Bundle Size Control
- Weekly `@next/bundle-analyzer` audit
- Hard budget: fail build if any chunk > 80KB gzip
- Dynamic imports for: Gallery, Testimonials, CartSheet, PaymobDialog, Framer Motion
- Conditional Arabic font: Cairo loaded ONLY when `locale === 'ar'`

---

## 9. PHASE-BY-PHASE OVERVIEW

### PHASE 00 — Setup + Backup (1–2 days)
- Git backup branch + filesystem backup
- Lighthouse baseline screenshots (BEFORE)
- Install shadcn/ui + all packages
- Create types, config, directories
- Update tailwind.config.js (dark mode class, CSS var colors)
- Rewrite globals.css + theme-variables.css
- Delete replaced files (old Header/Footer/homepage components)
- Test shadcn RTL in sandbox

### PHASE 01 — Theme System V2 (3–4 days)
- DB migration: add dark color columns to themes_v2
- Seed dark colors for all 10 themes
- Rewrite ThemeProvider (next-themes + CSS var injection)
- Create hex-to-hsl utility
- Create theme-normalizer utility
- Create theme-templates registry (10 templates)
- Update preferences-store (dark mode)
- Create use-dark-mode hook
- Extend i18n keys (en.json + ar.json, ~200 keys)
- Dashboard: DarkColorsEditor + dark preview toggle
- Register Categories + Gallery sections in component registry
- Add ISR + React Query caching for theme data
- Backend: Cache-Control headers on public endpoints

### PHASE 02 — Header + Footer (3–4 days)
- Create 5 header sub-components (Search, Cart, Language, DarkMode, WhatsApp)
- Create 5 header variants (Glass, Gradient, Solid, Floating, Minimal)
- Create main Header.tsx (dispatches to variant based on theme config)
- Create 5 footer variants (Expanded, Compact, Cards, SocialHeavy, Minimal)
- Create main Footer.tsx (dispatches to variant)
- Create MobileBottomNav (5 items, safe area inset)
- Create CartBottomBar (sticky above mobile nav)
- Create WhatsAppFloating (always-visible green bubble)
- Update locale layout.tsx to use new Header + Footer

### PHASE 03 — Homepage 8 Sections (4–5 days)
- Create 4 hero variants (Gradient, Video, Slider, Parallax)
- Create HeroSection (dispatches to variant)
- Create 3 product layout variants (Carousel, Grid, Masonry)
- Create FeaturedProductsSection (dispatches to variant)
- Create WhyChooseUsSection
- Create TestimonialsSection (Embla carousel)
- Create CategoriesSection (NEW — horizontal scroll pills)
- Create GallerySection (NEW — grid + lightbox)
- Create CTASection (gradient + WhatsApp button)
- Create ContactSection (phone, WhatsApp, location, hours)
- Rewrite homepage page.tsx (SSR, parallel fetch, section renderer)
- Theme builder preview mode (postMessage)

### PHASE 04 — Commerce Components (3–4 days)
- Rewrite ProductCard (image, bilingual name, price, qty stepper, WhatsApp)
- Create ProductCardSkeleton
- Create AJAX CartSheet (Sheet sidebar, no page reload)
- Create CartItem (qty controls, delete, notes)
- Create OrderStatusBadge
- Create OrderTimeline (vertical, animated pulse)
- Create OrdersTable (shadcn DataTable)
- Update cart-store (AJAX add, cross-tab sync, inventory check)

### PHASE 05 — All Page Redesigns (5–6 days) — Post-MVP
- Menu page: sticky filters, category tabs, price slider, product grid
- Cart page: 2-column desktop, accordion mobile, order notes
- Checkout page: 3-step progress, 3 payment methods (cash, Paymob, screenshot)
- Orders page: DataTable, filter tabs, search, pagination
- Order detail: timeline, breakdown table, customer info, actions
- Track order page: centered card, input, redirect
- Settings page: 4 tabs (Language, Appearance, Cart, Notifications)
- Table menu page: table header banner

### PHASE 06 — Cross-Cutting Features (3–4 days) — Post-MVP
- Search: store + hook + URL params sync
- Cart cross-tab sync (BroadcastChannel)
- Dark mode complete audit (grep all hardcoded colors)
- Language switcher complete testing
- WhatsApp order forwarding on checkout success
- WhatsApp floating button integration
- Arabic translations proofreading

### PHASE 07 — Performance + Testing (3–4 days) — Post-MVP
- Backend: sync.Map cache for themes + products
- Backend: cache invalidation on write
- Frontend: ISR for all public pages
- Frontend: parallel data fetching (Promise.all)
- Image optimization (next/image, blur, responsive sizes)
- Conditional Arabic font loading
- Code splitting (dynamic imports)
- Bundle analyzer + tree-shaking audit
- Prefetch critical routes
- Loading skeletons for all pages
- Error states (network, rate limit, payment, unavailable)
- Global error boundary
- Vitest setup + unit tests (ProductCard, cart-store, use-dark-mode, theme-normalizer)
- E2E tests (Playwright: checkout flow, language switch, dark mode)

### PHASE 08 — SEO + Monitoring + Production (2–3 days) — Post-MVP
- Dynamic metadata for all pages
- Schema.org JSON-LD (Restaurant + MenuItem)
- Arabic hreflang tags
- Dynamic sitemap.ts + robots.ts
- Sentry error tracking
- Google Analytics 4 (per-tenant)
- Web Vitals reporting
- Final Lighthouse audit (compare BEFORE vs AFTER)
- Test all 10 themes × light + dark = 20 combinations
- PWA manifest + meta tags

---

## 10. HEADER DESIGN — ALL PAGES

### Desktop (LTR English) — Glass Variant Example
```
┌──────────────────────────────────────────────────────────────────┐
│ 🍔 Restaurant Name │ [🔍 Search...]  │ Home Menu Orders │ EN🇺🇸 │ 🌙 │ 🛒3 │ 📱│
└──────────────────────────────────────────────────────────────────┘
  backdrop-filter: blur(20px)
  background: hsl(var(--theme-header-bg) / 0.95)
  sticky top-0 z-60
```

### Desktop (RTL Arabic) — Glass Variant
```
┌──────────────────────────────────────────────────────────────────┐
│📱│ 🛒3 │ 🌙 │ AR🇸🇦 │ الرئيسية القائمة الطلبات │ [🔍 ابحث...] │ 🍔 اسم المطعم │
└──────────────────────────────────────────────────────────────────┘
```

### Mobile (collapsed)
```
┌────────────────────────────────────┐
│ 🍔 Name          │  🔍  🛒3  ☰    │
└────────────────────────────────────┘
  ☰ → Sheet drawer from left (RTL: right)
```

### Mobile Bottom Nav (all inner pages)
```
┌────────────────────────────────────────────────┐
│  🏠      🍽️      🛒       📋      ⚙️          │
│ Home    Menu   Cart(3)  Orders  Settings       │
│  ↑ active = theme-primary color + dot          │
│ pb-[env(safe-area-inset-bottom)] for iOS       │
└────────────────────────────────────────────────┘
```

### WhatsApp Floating Bubble (always visible)
```
                                    ┌─────┐
                                    │ 📱  │ ← green, pulse animation
                                    │ WA  │    bottom-right, z-50
                                    └─────┘    above MobileBottomNav
```

---

## 11. SUCCESS CRITERIA

| Metric | MVP Target | Final Target | Measurement |
|---|---|---|---|
| Lighthouse Mobile | ≥ 70 | ≥ 95 | `npx lighthouse --mobile` |
| Lighthouse Desktop | ≥ 80 | ≥ 98 | `npx lighthouse --desktop` |
| LCP | < 3.0s | < 2.5s | Web Vitals |
| FID | < 150ms | < 100ms | Web Vitals |
| CLS | < 0.15 | < 0.1 | Web Vitals |
| Page load (Egypt 4G) | < 3.0s | < 2.0s | Network throttle test |
| Bundle size (gzip) | < 250KB | < 180KB | `@next/bundle-analyzer` |
| Dark mode | All components | All pages + all themes | Manual + Playwright |
| RTL Arabic | All pages | All pages + edge cases | Manual + Playwright |
| Theme switching | All 10 themes render | All 10 × light + dark | Manual |
| AJAX cart | Add to cart works | Cross-tab sync + inventory | Manual |
| WhatsApp | Button visible | Order forwarding works | Manual |
| Unit tests | Core components | 100% passing | `pnpm test` |
| E2E tests | Checkout flow | All flows | `pnpm e2e` |

---

## 12. IMPORTANT DECISIONS

### Why rebuild in-place (not v2 folder)
The old app structure is correct. What needs to change is the UI components inside it. Creating a v2 folder would duplicate middleware, API client, i18n config, stores — all of which work fine.

### Why shadcn/ui over custom
Accessible, RTL-capable, dark-mode-ready components with CSS variable integration. Building from scratch = 10x longer with more bugs.

### Why Embla over Slick
Tree-shakeable (7KB vs 30KB), SSR-compatible, no jQuery, easier to customize.

### Why next-themes for dark mode
Prevents flash of incorrect theme on SSR, syncs across tabs, handles system preference. Use `disableTransitionOnChange` to avoid flash on slow networks.

### Why CSS vars not Tailwind dark: for Header/Footer
Header and Footer get colors from theme API data (dynamic). CSS vars are the only way to inject API-driven colors AND have them respond to dark mode. Tailwind `dark:` only works for static colors.

### Why Cairo font for Arabic
Best-looking Arabic web font for restaurant menus — clear, modern, widely used in Egyptian SaaS. **Load conditionally** only when `locale === 'ar'` to save 200KB on English pages.

### Why MVP first (15 days)
160 tasks across 10 phases = 40+ working days. Risk of burnout and scope creep. Ship Phases 00–04 first, collect feedback, then iterate remaining phases.

### Why 5 header variants instead of 1
A "Luxury Gold" restaurant and a "Minimal White" restaurant should FEEL fundamentally different — not just color-swapped. Header variant is the biggest visual differentiator.

### Why AJAX cart is critical
Current: Add to cart = full page reload = lost browsing context = dropped conversion.
After: Click Add → sidebar slides in → continue browsing = +35% conversion lift.

---

*Last Updated: 2026-03-03 | Status: Planning Phase | MVP: 15 days*
