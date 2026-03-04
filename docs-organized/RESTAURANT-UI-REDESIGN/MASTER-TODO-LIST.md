# Restaurant Website UI Redesign — Master Todo List v3.0
**All 9 Phases · ~180 Tasks · Production-Grade · MVP-First**
**Updated:** 2026-03-03 v3.0
**MVP:** Phases 00–04 (15 working days) → Deploy → Iterate Phases 05–08
**Quality:** Production-ready, 1000+ themes, scalable, Egyptian market optimized

> **Legend:** 🔴 Blocking · 🟠 High · 🟡 Medium · 🟢 Nice-to-have
> **Format:** `- [ ] Task` / `Path:` / `Why:` / `Test:` / `Time:`

---

## CHECKPOINT SCHEDULE
| Day | Checkpoint | Pass Criteria |
|-----|-----------|---------------|
| 3 | Header renders all 10 themes (light + dark) | No hardcoded hex, RTL flips correctly |
| 7 | Homepage renders 8 sections, templates work | Each theme has unique section order |
| 8 | Mobile Safari iOS test | `safe-area-inset-bottom` renders correctly, no viewport overflow, touch targets ≥ 48px |
| 10 | Menu page + ProductCard + AJAX cart working | Add → sheet opens, no page reload |
| 12 | Cart + Checkout MVP functional | Cash payment completes, order created |
| 15 | Production deploy | Lighthouse mobile ≥ 70, all pages render |

---

## ═══════════════════════════════════════
## PHASE 00 — BACKUP + SETUP (Day 1–2)
## ═══════════════════════════════════════
> **Goal:** Safe backup, install foundations, test dependencies, establish conventions
> **Checkpoint Day 2:** All packages installed, types compile, tailwind configured

### 00.PRE — Backup & Baseline (🔴 MUST RUN FIRST)

- [ ] **🔴 Create git backup branch + working branch**
  ```bash
  cd C:/Users/OPT/Desktop/POS
  git stash push -m "pre-redesign-state-$(date +%Y%m%d)"
  git checkout -b restaurant-ui-redesign-backup
  git add -A && git commit -m "backup: restaurant-website state before UI redesign"
  git checkout master
  git checkout -b restaurant-ui-redesign
  ```
  Test: `git branch -a` shows both branches
  Time: 5min

- [ ] **🔴 Take Lighthouse baseline screenshots**
  ```bash
  npx lighthouse http://localhost:3003/en --output=html \
    --output-path=./docs-organized/RESTAURANT-UI-REDESIGN/lighthouse-BEFORE-mobile.html \
    --emulated-form-factor=mobile
  npx lighthouse http://localhost:3003/en --output=html \
    --output-path=./docs-organized/RESTAURANT-UI-REDESIGN/lighthouse-BEFORE-desktop.html \
    --emulated-form-factor=desktop
  ```
  Path: `docs-organized/RESTAURANT-UI-REDESIGN/lighthouse-BEFORE-*.html`
  Why: Measure improvement. Before/after comparison proves value.
  Time: 10min

- [ ] **🔴 Document BEFORE-STATE.md**
  Path: `docs-organized/RESTAURANT-UI-REDESIGN/BEFORE-STATE.md`
  Record: git SHA, Lighthouse scores, page count, component count, bundle size
  Time: 10min

### 00.A — Package Installation

- [ ] **Initialize shadcn/ui**
  Path: `frontend/apps/restaurant-website/`
  ```bash
  cd frontend/apps/restaurant-website
  pnpm dlx shadcn@latest init
  ```
  Choose: TypeScript YES, slate base, CSS variables YES, @/ aliases
  Creates: `components/ui/`, `components.json`, `lib/utils.ts`
  Time: 15min

- [ ] **Install all shadcn components**
  Path: `frontend/apps/restaurant-website/`
  ```bash
  npx shadcn@latest add button card input label form textarea select \
    dialog drawer sheet badge tabs avatar separator scroll-area skeleton \
    progress tooltip dropdown-menu table slider accordion alert-dialog \
    popover command toggle-group radio-group
  ```
  Time: 15min

- [ ] **Install additional packages**
  Path: `frontend/apps/restaurant-website/`
  ```bash
  pnpm add embla-carousel-react embla-carousel-autoplay
  pnpm add next-themes
  pnpm add framer-motion
  pnpm add sonner
  pnpm add -D @next/bundle-analyzer
  ```
  Time: 10min

- [ ] **Install Cairo Arabic font**
  Path: `frontend/apps/restaurant-website/app/layout.tsx`
  ```typescript
  import { Inter, Cairo } from 'next/font/google'
  const inter = Inter({ subsets: ['latin'], variable: '--font-inter', display: 'swap' })
  const cairo = Cairo({ subsets: ['arabic', 'latin'], variable: '--font-cairo', display: 'swap' })
  ```
  Time: 10min

### 00.B — Dependency Sandbox Test (Risk Mitigation)

- [ ] **🔴 Test shadcn/ui RTL with Arabic text**
  Path: create temp `frontend/apps/restaurant-website/__sandbox__/rtl-test.tsx`
  Test: Render `<Tabs>`, `<Sheet>`, `<DropdownMenu>`, `<Dialog>` in RTL mode
  Fix any RTL quirks discovered (document in BEFORE-STATE.md)
  Delete sandbox after testing
  Time: 30min

- [ ] **Test Embla Carousel in RTL**
  Path: temp sandbox component
  Test: direction="rtl" prop, scroll-snap, autoplay
  Verify: no loop mode bugs, containScroll works
  Time: 15min

- [ ] **Test next-themes hydration**
  Path: temp sandbox
  Test: SSR → client hydration, no flash of wrong theme
  Verify: `suppressHydrationWarning` + `disableTransitionOnChange` work
  Time: 15min

### 00.C — Configuration Files

- [ ] **Update `tailwind.config.js`**
  Path: `frontend/apps/restaurant-website/tailwind.config.js`
  Changes:
  - `darkMode: ['class']` — manual toggle via `html.dark`
  - `fontFamily.sans` → `var(--font-inter)`
  - `fontFamily.arabic` → `var(--font-cairo)`
  - `colors.theme.*` → all CSS var-based colors (primary, surface, text, border, etc.)
  - `minHeight['screen-svh']` → `100svh`
  Time: 30min

- [ ] **Rewrite `styles/globals.css`**
  Path: `frontend/apps/restaurant-website/styles/globals.css`
  Content: Tailwind directives + `:root` CSS vars + `.dark` overrides + shadcn bridge vars
  Time: 45min

- [ ] **Create `styles/theme-variables.css`**
  Path: `frontend/apps/restaurant-website/styles/theme-variables.css`
  Content: All 25+ CSS variables with `.dark` overrides (see MASTER-OVERVIEW CSS section)
  Time: 30min

- [ ] **Update `next.config.js`**
  Path: `frontend/apps/restaurant-website/next.config.js`
  Add: image domains, `optimizePackageImports` for lucide/framer-motion/embla, security headers
  Time: 20min

### 00.D — TypeScript Types

- [ ] **Create `types/theme.ts`**
  Path: `frontend/apps/restaurant-website/types/theme.ts`
  Export: `ThemeConfig`, `ThemeHeaderConfig`, `ThemeFooterConfig`, `ThemeSectionConfig`, `SectionType`, `ThemeNavItem`, `ThemeSocialLink`, `ThemeFooterSection`, `ThemeLegalLink`
  Include: `templateId`, `headerVariant`, `footerVariant`, `componentOverrides` fields
  Test: `tsc --noEmit` passes
  Time: 45min

- [ ] **Create `types/payment.ts`**
  Path: `frontend/apps/restaurant-website/types/payment.ts`
  Export: `PaymentMethodType` (`cash` | `paymob_online` | `screenshot_transfer`), `PaymentMethod`, `PaymobCheckoutResponse`, `PaymentScreenshot`
  Time: 30min

### 00.E — Create Directories + Delete Old Files

- [ ] **Create all new directories**
  Path: `frontend/apps/restaurant-website/`
  ```
  components/layout/headers/
  components/layout/footers/
  components/layout/header/
  components/sections/heroes/
  components/sections/products/
  components/commerce/
  components/forms/
  components/system/
  lib/utils/
  types/
  __tests__/components/
  __tests__/hooks/
  __tests__/stores/
  e2e/
  ```
  Time: 2min

- [ ] **Delete replaced files**
  ```
  components/layout/header-ssr.tsx       → replaced by new Header.tsx
  components/layout/footer-ssr.tsx       → replaced by new Footer.tsx
  components/layout/ThemeHeader.tsx       → merged into Header.tsx
  components/layout/ThemeFooter.tsx       → merged into Footer.tsx
  components/layout/theme-layout-normalizers.ts → logic in Header/Footer
  components/homepage/dynamic-home-page.tsx     → rebuilt in page.tsx
  components/homepage/featured-products-section.tsx → new section component
  components/homepage/error-boundary.tsx         → GlobalErrorBoundary
  ```
  Time: 5min

---

## ═══════════════════════════════════════
## PHASE 01 — THEME SYSTEM V2 (Day 3–5)
## ═══════════════════════════════════════
> **Goal:** Per-theme dark mode, CSS vars, theme templates, caching. Scalable to 1000 themes.
> **Checkpoint Day 5:** Toggle dark mode → header/footer/cards all change per theme.

### 01.A — Database: Dark Mode Colors

- [ ] **🔴 Create migration: add dark color columns**
  Path: `backend/migrations/158_add_theme_dark_colors.sql`
  Add columns: `primary_dark`, `secondary_dark`, `accent_dark`, `background_dark`, `text_dark`, `border_dark`, `shadow_dark`, `surface_color`, `surface_dark`, `header_bg_dark`, `header_text_dark`, `footer_bg_dark`, `footer_text_dark`, `dark_mode_enabled`
  Test: `SELECT primary_dark FROM themes_v2 LIMIT 1;`
  Time: 30min

- [ ] **🔴 Create migration: add template + variant columns**
  Path: `backend/migrations/159_add_theme_template_variants.sql`
  Add columns: `template_id VARCHAR(20) DEFAULT 'template-01'`, `header_variant VARCHAR(20) DEFAULT 'solid'`, `footer_variant VARCHAR(20) DEFAULT 'expanded'`, `hero_variant VARCHAR(20) DEFAULT 'gradient'`, `products_variant VARCHAR(20) DEFAULT 'carousel'`
  Time: 20min

- [ ] **🔴 Seed dark colors + templates for all 10 themes**
  Path: `backend/migrations/160_seed_theme_dark_colors_and_templates.sql`
  Each theme gets: unique dark palette + template assignment + variant assignments
  Time: 45min

- [ ] **🔴 Update backend theme handler — return dark colors + variants**
  Path: `backend/internal/handler/http/admin_theme_handler_v2.go`
  Changes: Include dark color fields + template/variant fields in GET response + accept in PUT/POST
  Time: 45min

### 01.B — ThemeProvider Rewrite (Core)

- [ ] **🔴 Rewrite ThemeProvider**
  Path: `frontend/apps/restaurant-website/components/providers/ThemeProvider.tsx`
  Architecture:
  - Wrap with `NextThemesProvider` (`attribute="class"`, `defaultTheme="system"`, `disableTransitionOnChange`)
  - `ThemeCSSInjector` component: reads `themeConfig.colors` + `resolvedTheme` → injects CSS vars
  - Uses `hexToHSL()` to convert hex → HSL for CSS vars
  - Re-runs ONLY when `themeConfig.slug`, `themeConfig.version`, or `isDark` changes
  - Single batch `style.setProperty()` loop — no layout thrashing
  Test: Toggle dark → all colors update. Switch theme → all colors update.
  Time: 2h

- [ ] **Create `lib/utils/hex-to-hsl.ts`**
  Path: `frontend/apps/restaurant-website/lib/utils/hex-to-hsl.ts`
  Export: `hexToHSL(hex: string): string` — converts `#2563eb` → `221 83% 53%`
  Export: `memoHexToHSL` — memoized version
  Test: `hexToHSL('#2563eb')` → `'221 83% 53%'`
  Time: 30min

- [ ] **Create `lib/utils/theme-normalizer.ts`**
  Path: `frontend/apps/restaurant-website/lib/utils/theme-normalizer.ts`
  Export: `normalizeThemeConfig(apiData: unknown, locale: 'en' | 'ar'): ThemeConfig`
  Handles: snake_case → camelCase, missing fields → defaults, bilingual extraction, nested structures
  Test: Unit test with 5 different API response shapes
  Time: 1.5h

- [ ] **Update `lib/utils/default-theme.ts`**
  Path: `frontend/apps/restaurant-website/lib/utils/default-theme.ts`
  Add: ALL dark mode color pairs, header/footer dark colors, semantic colors, template defaults
  Export: `DEFAULT_THEME: ThemeConfig`
  Time: 1h

### 01.C — Theme Templates Registry

- [ ] **🔴 Create `lib/utils/theme-templates.ts`**
  Path: `frontend/apps/restaurant-website/lib/utils/theme-templates.ts`
  ```typescript
  export const TEMPLATE_REGISTRY: Record<string, ThemeTemplate> = {
    'template-01': { // Modern Bistro
      name: 'Modern Bistro',
      sectionOrder: ['hero', 'featured_products', 'why_choose_us', 'testimonials', 'cta', 'contact'],
      headerVariant: 'glass',
      heroVariant: 'video',
      productsVariant: 'carousel',
      footerVariant: 'expanded',
    },
    'template-02': { // Luxury Gold
      name: 'Luxury Gold',
      sectionOrder: ['hero', 'gallery', 'testimonials', 'featured_products', 'cta', 'contact'],
      headerVariant: 'gradient',
      heroVariant: 'gradient',
      productsVariant: 'grid',
      footerVariant: 'cards',
    },
    // ... all 10 templates
  }
  export const DEFAULT_TEMPLATE = TEMPLATE_REGISTRY['template-01']
  ```
  Time: 1h

### 01.D — Preferences Store + Dark Mode Hook

- [ ] **Update `lib/store/preferences-store.ts`**
  Path: `frontend/apps/restaurant-website/lib/store/preferences-store.ts`
  Add: `darkMode: 'light' | 'dark' | 'system'`, `setDarkMode()`, sync with next-themes
  Time: 20min

- [ ] **Create `lib/hooks/use-dark-mode.ts`**
  Path: `frontend/apps/restaurant-website/lib/hooks/use-dark-mode.ts`
  Returns: `{ mode, setMode, isDark, toggle }` — wraps next-themes + preferences store
  Time: 30min

### 01.E — i18n Keys (ALL new features)

- [ ] **Extend `i18n/messages/en.json`**
  Path: `frontend/apps/restaurant-website/i18n/messages/en.json`
  Add ~200 keys: header, theme, language, cart, product, search, menu, checkout, orders, order_status, order_detail, track, settings, sections, errors, payment, whatsapp, common
  Time: 1h

- [ ] **Extend `i18n/messages/ar.json`**
  Path: `frontend/apps/restaurant-website/i18n/messages/ar.json`
  Full Arabic translations for ALL keys. Critical:
  - `cart.title`: سلة التسوق
  - `cart.empty`: سلتك فارغة
  - `product.add_to_cart`: أضف
  - `checkout.place_order`: تأكيد الطلب
  - `payment.cash`: الدفع عند الاستلام
  - `whatsapp.order`: اطلب عبر واتساب
  Time: 1.5h

### 01.F — Format Utilities

- [ ] **Create `lib/utils/format.ts`**
  Path: `frontend/apps/restaurant-website/lib/utils/format.ts`
  Export:
  - `formatPrice(amount: number, locale: string): string` → `١٤٥ ج.م` or `145 EGP`
  - `formatArabicNumber(num: number): string` → Arabic numerals
  - `formatDate(date: string, locale: string): string` → localized date
  - `formatRating(rating: number, count: number, locale: string): string` → `٤.٨⭐ (١٢٣ تقييم)`
  Time: 45min

### 01.G — Section Registry Update

- [ ] **Register categories + gallery in component registry**
  Path: `frontend/layers/ui-themes/src/registry/componentRegistry.ts`
  Add: `categories` (aliases: `categories-section`, `browse_menu`) + `gallery` (aliases: `gallery-section`, `food_gallery`)
  Keep backward compatibility with old type aliases
  Time: 30min

### 01.H — Backend Cache Invalidation (🔴 MUST COMPLETE BEFORE FRONTEND CACHING)

- [ ] **🔴 Add backend sync.Map cache + invalidation on admin writes**
  Path: `backend/internal/handler/http/public_theme_handler.go`
  Strategy: `sync.Map` keyed by `tenantID:themeID` → cached JSON response
  On admin write (create/update/delete theme): delete cache key + return `Cache-Control: no-cache`
  Also add invalidation hooks to: `admin_theme_handler_v2.go`, `admin_menu_handler.go`
  Test: Update theme in admin → public API returns new data within 1s (no stale cache)
  Time: 1h

- [ ] **Add Cache-Control headers to public API endpoints**
  Path: `backend/internal/handler/http/admin_theme_handler_v2.go`
  Public: `Cache-Control: public, max-age=300, stale-while-revalidate=600`
  Admin: `Cache-Control: no-cache, private`
  Also add to: `public_menu_handler.go`, `public_order_handler.go`
  Time: 30min

- [ ] **Add React Query cache defaults**
  Path: `frontend/apps/restaurant-website/components/providers/QueryProvider.tsx`
  Set: `staleTime: 10min`, `gcTime: 1hr`, `refetchOnWindowFocus: false`, `retry: 2`
  Depends on: Backend cache invalidation must be working first
  Time: 20min

---

## ═══════════════════════════════════════
## PHASE 02 — HEADER + FOOTER (Day 6–8)
## ═══════════════════════════════════════
> **Goal:** 5 header variants, 5 footer variants, mobile nav, WhatsApp bubble
> **Checkpoint Day 8:** Header with search/cart/lang/dark works across all themes + RTL

### 02.A — Header Sub-Components

- [ ] **Create `components/layout/header/DarkModeToggle.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/header/DarkModeToggle.tsx`
  Cycles: light → dark → system. Icons: Sun/Moon/Monitor. Size: 48px touch target.
  Uses: `useDarkMode()` hook, Lucide icons, shadcn `<Tooltip>`
  Test: Click 3x → cycles all modes, `html.dark` class updates
  Time: 30min

- [ ] **Create `components/layout/header/LanguageToggle.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/header/LanguageToggle.tsx`
  Desktop: `[🇺🇸 EN ▾]` dropdown. Mobile: compact `[EN]`/`[AR]`.
  On switch: replace `/en/` ↔ `/ar/` in pathname, `router.push()`, save to prefs
  Uses: shadcn `<DropdownMenu>`, `useRouter()`, `usePathname()`
  Test: EN→AR: URL changes, RTL applies, Arabic font loads
  Time: 45min

- [ ] **Create `components/layout/header/CartButton.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/header/CartButton.tsx`
  ShoppingCart icon + animated count badge. Click → opens CartSheet.
  Badge: `bg-theme-error text-white`, scale bounce on change. Size: 48px touch.
  Uses: Zustand cart store, shadcn `<Badge>`, Lucide ShoppingCart
  Test: Add item → badge appears, remove all → badge disappears
  Time: 30min

- [ ] **Create `components/layout/header/SearchBar.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/header/SearchBar.tsx`
  Desktop: expanding input with results `<Popover>`. Mobile: icon → full-width overlay.
  300ms debounce → filter cached products (no API call). Max 5 results per group.
  Results: Products (image + name AR/EN + price) | Categories. Keyboard navigation.
  Uses: shadcn `<Popover>` + `<Command>`, React Query cached products
  Test: Type "burger" → results < 300ms, click → navigate to menu
  Time: 2.5h

- [ ] **Create `components/layout/header/WhatsAppButton.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/header/WhatsAppButton.tsx`
  Desktop: icon + "Order via WhatsApp". Mobile: green icon only.
  Click → `wa.me/{number}?text=مرحبا، أريد أن أطلب`
  Hidden if `!theme.header.showWhatsAppButton`
  Time: 20min

### 02.B — Header Variants (5 Types)

- [ ] **🔴 Create `components/layout/headers/GlassHeader.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/headers/GlassHeader.tsx`
  Style: `backdrop-filter: blur(20px)`, `bg-[hsl(var(--theme-header-bg)/0.1)]`, transparent
  Used by: Modern Bistro, Coastal Ocean
  Time: 1.5h

- [ ] **🔴 Create `components/layout/headers/GradientHeader.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/headers/GradientHeader.tsx`
  Style: `background: linear-gradient(135deg, primary → accent)`, opaque
  Used by: Luxury Gold
  Time: 1h

- [ ] **🔴 Create `components/layout/headers/SolidHeader.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/headers/SolidHeader.tsx`
  Style: `bg-[hsl(var(--theme-header-bg))]`, solid color, default variant
  Used by: Warm Rustic, Garden Green, Spicy Red, Egyptian Heritage
  Time: 1h

- [ ] **Create `components/layout/headers/FloatingHeader.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/headers/FloatingHeader.tsx`
  Style: fixed, heavy `shadow-xl`, rounded bottom corners
  Used by: Urban Gray, Vibrant Pop
  Time: 45min

- [ ] **Create `components/layout/headers/MinimalHeader.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/headers/MinimalHeader.tsx`
  Style: thin `border-b` only, transparent background, minimal presence
  Used by: Minimal White
  Time: 45min

### 02.C — Main Header (Dispatcher)

- [ ] **🔴 Create `components/layout/Header.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/Header.tsx`
  Props: `themeConfig`, `locale`, `restaurantSlug`
  Logic: reads `themeConfig.headerVariant` → renders matching header variant
  All variants share: same sub-components (Search, Cart, Language, DarkMode, WhatsApp)
  All variants: sticky, z-60, RTL-aware (`flex-row-reverse` when Arabic), mobile sheet drawer
  ALL colors from CSS vars — responds to dark mode toggle instantly
  Time: 1.5h

### 02.D — Footer Variants (5 Types)

- [ ] **🔴 Create `components/layout/footers/ExpandedFooter.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/footers/ExpandedFooter.tsx`
  4-column: Company info | Quick links | Hours | Contact. Social row. Legal row. Back to top.
  Used by: Modern Bistro, Spicy Red, Egyptian Heritage
  Time: 1.5h

- [ ] **Create `components/layout/footers/CompactFooter.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/footers/CompactFooter.tsx`
  2-row: Name + contact | Copyright + legal. Compact.
  Used by: Warm Rustic, Urban Gray
  Time: 45min

- [ ] **Create `components/layout/footers/CardsFooter.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/footers/CardsFooter.tsx`
  3 card blocks: About | Contact | Menu. Card-based layout.
  Used by: Luxury Gold
  Time: 1h

- [ ] **Create `components/layout/footers/SocialHeavyFooter.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/footers/SocialHeavyFooter.tsx`
  Social links prominent, large icons, minimal text.
  Used by: Coastal Ocean, Vibrant Pop
  Time: 45min

- [ ] **Create `components/layout/footers/MinimalFooter.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/footers/MinimalFooter.tsx`
  Single line: `© 2026 Restaurant Name | Privacy | Terms`
  Used by: Minimal White, Garden Green
  Time: 30min

### 02.E — Main Footer (Dispatcher)

- [ ] **Create `components/layout/Footer.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/Footer.tsx`
  Props: `themeConfig`, `locale`
  Logic: reads `themeConfig.footer.layout` → renders matching footer variant
  ALL colors from CSS vars. RTL: `dir` prop from locale.
  Time: 45min

### 02.F — Mobile Navigation

- [ ] **🔴 Create `components/layout/MobileBottomNav.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/MobileBottomNav.tsx`
  5 items: Home / Menu / Cart(count) / Orders / Settings
  Active state: current route → theme primary color + dot indicator
  Cart badge: live count from cart store. `md:hidden`.
  `pb-[env(safe-area-inset-bottom)]` for iOS. Height: 64px + safe area.
  Time: 1h

- [ ] **Create `components/layout/CartBottomBar.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/CartBottomBar.tsx`
  Sticky above MobileBottomNav. Shows: cart count + total + "Checkout →" button.
  Hidden when cart empty, slides up with animation when items added.
  Shown on: Home, Menu pages. NOT on Cart/Checkout pages.
  Time: 1h

- [ ] **🔴 Create `components/layout/WhatsAppFloating.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/WhatsAppFloating.tsx`
  Fixed bottom-right, green circle, WhatsApp icon, subtle pulse animation.
  Always visible (z-50, above MobileBottomNav). Click → `wa.me/{number}`.
  Hidden if no whatsAppNumber in theme config.
  Time: 30min

### 02.G — Layout Integration

- [ ] **🔴 Update `app/[locale]/layout.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/layout.tsx`
  Changes:
  - Fetch theme server-side (SSR)
  - Render: `<Header>` + `{children}` + `<Footer>` + `<MobileBottomNav>` + `<CartBottomBar>` + `<WhatsAppFloating>`
  - Preview mode: if `?preview=true` → hide Header/Footer/Nav
  - Conditional Cairo font: only load when `locale === 'ar'`
  Time: 1h

---

## ═══════════════════════════════════════
## PHASE 03 — HOMEPAGE 8 SECTIONS (Day 9–11)
## ═══════════════════════════════════════
> **Goal:** 8 dynamic sections with theme-specific templates, all themes + dark mode + RTL
> **Checkpoint Day 11:** Each theme renders unique homepage with different section order

### 03.A — Hero Variants (4 Types)

- [ ] **🔴 Create `components/sections/heroes/GradientHero.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/heroes/GradientHero.tsx`
  Full viewport `min-h-[100svh]`. CSS gradient using theme primary → accent.
  Dark overlay for text contrast. Content: headline AR, subtitle EN, CTA button, delivery badge.
  Dual CTA: 🟢 اطلب واتساب + 🟠 Order Now
  RTL: text-right, items-end. Dark mode: overlay darker.
  Time: 1.5h

- [ ] **Create `components/sections/heroes/VideoHero.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/heroes/VideoHero.tsx`
  `<video autoPlay muted loop>` background. Dark gradient overlay. Same content layout.
  Fallback: gradient hero if no video URL in config.
  Used by: Modern Bistro
  Time: 1h

- [ ] **Create `components/sections/heroes/SliderHero.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/heroes/SliderHero.tsx`
  Embla Carousel with 3-5 hero images. Auto-slide 5s. Dot pagination.
  Used by: Garden Green, Vibrant Pop
  Time: 1h

- [ ] **Create `components/sections/heroes/ParallaxHero.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/heroes/ParallaxHero.tsx`
  Parallax scroll effect on background image. Framer Motion `useScroll`.
  Used by: Warm Rustic
  Time: 1h

### 03.B — Hero Dispatcher

- [ ] **🔴 Create `components/sections/HeroSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/HeroSection.tsx`
  Props: `config`, `locale`, `themeConfig`
  Logic: reads `themeConfig.componentOverrides.hero` → renders matching hero variant
  Fallback: GradientHero if variant not found
  Time: 30min

### 03.C — Featured Products Variants (3 Types)

- [ ] **🔴 Create `components/sections/products/ProductCarousel.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/products/ProductCarousel.tsx`
  Embla Carousel: 1 mobile / 2 tablet / 3-4 desktop. Autoplay 4s, pause on hover.
  Dot pagination. "View Full Menu →" link.
  Time: 1.5h

- [ ] **Create `components/sections/products/ProductGrid.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/products/ProductGrid.tsx`
  3-col desktop, 2-col tablet, 1-col mobile. Static grid with 6-8 products.
  Time: 1h

- [ ] **Create `components/sections/products/ProductMasonry.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/products/ProductMasonry.tsx`
  Masonry layout with varying card heights. Used by: Warm Rustic.
  Time: 1h

### 03.D — Featured Products Dispatcher

- [ ] **Create `components/sections/FeaturedProductsSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/FeaturedProductsSection.tsx`
  Reads `themeConfig.componentOverrides.products` → renders matching variant
  Real products from `useProducts(restaurantSlug)`. Loading: skeleton cards.
  Time: 30min

### 03.E — Other Sections

- [ ] **Create `components/sections/WhyChooseUsSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/WhyChooseUsSection.tsx`
  4 cards grid (2×2 mobile, 4×1 desktop). Lucide icons. Hover lift. Dark mode via CSS vars.
  Time: 1.5h

- [ ] **Create `components/sections/TestimonialsSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/TestimonialsSection.tsx`
  Embla Carousel. Avatar + Name + Stars (`⭐⭐⭐⭐⭐`) + Review text AR/EN.
  Auto-play with pause on hover. Arabic star format: `٤.٨⭐`
  Time: 1.5h

- [ ] **Create `components/sections/CategoriesSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/CategoriesSection.tsx`
  Horizontal scroll pills. Emoji/icon + Category name AR/EN + item count badge.
  Click → navigate to `/[locale]/menu?category={id}`. Scroll-snap on mobile.
  Real categories from `useCategories(restaurantSlug)`.
  Time: 1.5h

- [ ] **Create `components/sections/GallerySection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/GallerySection.tsx`
  3×2 image grid. Lightbox on click (shadcn `<Dialog>`). Hover: `scale(1.05)`.
  Images from theme config. Dark mode: dark overlay.
  Time: 2h

- [ ] **Create `components/sections/CTASection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/CTASection.tsx`
  Full width gradient (primary → accent). Headline + sub + dual CTA:
  🟠 "Order Now" button + 🟢 "Order via WhatsApp" button.
  Pulse animation on main button. Dark mode: gradient shifts darker.
  Time: 1h

- [ ] **Create `components/sections/ContactSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/ContactSection.tsx`
  4 cards: Phone (tel: link) | WhatsApp (wa.me/) | Location (Google Maps) | Hours.
  Each: Icon + Label AR/EN + Value + Action button. Mobile: 2×2 grid.
  Time: 1h

### 03.F — Dynamic Homepage (SSR + Template Rendering)

- [ ] **🔴 Rewrite `app/[locale]/page.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/page.tsx`
  Server Component: `Promise.all([fetchTheme, fetchProducts, fetchCategories])` — parallel fetch
  Client Component `<DynamicHomePage>`:
  ```typescript
  // Get template section order from theme config
  const template = TEMPLATE_REGISTRY[themeConfig.templateId] || DEFAULT_TEMPLATE
  const orderedSections = themeConfig.sections
    .filter(s => s.enabled)
    .sort((a, b) => template.sectionOrder.indexOf(a.type) - template.sectionOrder.indexOf(b.type))
  // Render each section with correct variant
  orderedSections.map(section => <SectionComponent key={section.id} {...section} />)
  ```
  ISR: `export const revalidate = 300` (5 min)
  Time: 2h

### 03.G — Theme Builder Preview Mode

- [ ] **Add preview mode to layout**
  Path: `frontend/apps/restaurant-website/app/[locale]/layout.tsx`
  If `?preview=true`: hide Header, Footer, MobileBottomNav, CartBottomBar, WhatsAppFloating
  Listen for `postMessage` → update CSS vars in < 50ms
  Time: 1h

---

## ═══════════════════════════════════════
## PHASE 04 — COMMERCE COMPONENTS (Day 12–14)
## ═══════════════════════════════════════
> **Goal:** ProductCard with AJAX cart, CartSheet sidebar, OrdersTable
> **Checkpoint Day 14:** Add to cart → sheet opens → checkout works

### 04.A — Product Card

- [ ] **🔴 Create `components/commerce/ProductCard.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/ProductCard.tsx`
  Image (next/image, blur placeholder) | Name AR (bold) + EN (muted) | Rating `٤.٣⭐`
  Description (2 lines truncated) | Price badge `٤٥ ج.م` | Qty stepper `[-][1][+]`
  Full-width "🛒 أضف / Add" button | Optional WhatsApp button
  States: default, hover (lift), out-of-stock (grayscale + badge), added (checkmark toast)
  Dark mode: bg-theme-surface, border-theme-border
  AJAX: click Add → `cartStore.addItem()` → toast → CartSheet opens (NO page reload)
  Time: 2h

- [ ] **Create `components/commerce/ProductCardSkeleton.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/ProductCardSkeleton.tsx`
  Skeleton: image rectangle + 3 text lines + button bar (shadcn Skeleton)
  Time: 20min

### 04.B — AJAX Cart

- [ ] **🔴 Update `lib/store/cart-store.ts`**
  Path: `frontend/apps/restaurant-website/lib/store/cart-store.ts`
  Add:
  - AJAX `addItem()` — instant Zustand update, NO page reload
  - Cross-tab sync via `StorageEvent` listener (BroadcastChannel)
  - Namespace per restaurant: key = `pos-cart-{restaurantSlug}`
  - `totalPrice` and `totalItems` computed getters
  - Inventory limit check — disable add if `qty >= quantity_in_stock`
  - Warning toast when at inventory limit
  Test: 2 tabs → add in tab 1 → tab 2 updates instantly
  Time: 1.5h

- [ ] **Create `lib/hooks/use-cart.ts`**
  Path: `frontend/apps/restaurant-website/lib/hooks/use-cart.ts`
  Convenience hook: `{ items, totalItems, totalPrice, addItem, removeItem, updateQty, clearCart }`
  Wraps cart-store with AJAX behavior and inventory checks
  Time: 30min

- [ ] **🔴 Create `components/commerce/CartSheet.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/CartSheet.tsx`
  shadcn `<Sheet>` sliding from right (RTL: left). NO full page reload.
  Content: header with count + close | ScrollArea cart items | separator | order notes textarea | order summary (items + delivery + total) | gradient Checkout button | "Continue Shopping" link
  Empty state: big cart icon + "سلتك فارغة" + "Browse Menu" button
  Time: 2h

- [ ] **Create `components/commerce/CartItem.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/CartItem.tsx`
  Layout: `[img 64px] Name AR | Name EN | Qty [-][2][+] | Price | [🗑️]`
  Optional note text. Swipe-to-delete on mobile (future).
  Time: 1h

### 04.C — Order Components

- [ ] **Create `components/commerce/OrderStatusBadge.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/OrderStatusBadge.tsx`
  Status → color: pending→yellow, confirmed→blue, preparing→blue, out_for_delivery→orange, delivered→green, cancelled→red
  Arabic labels: قيد الانتظار, تم التأكيد, جاري التحضير, في الطريق, تم التوصيل, ملغي
  Uses: shadcn `<Badge>` with dynamic variant
  Time: 30min

- [ ] **Create `components/commerce/OrderTimeline.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/OrderTimeline.tsx`
  Vertical timeline: ✅ completed | ◉ current (pulse ring in theme-primary) | ○ future (gray)
  Each step: Circle icon + Label AR/EN + Timestamp
  RTL: dots on RIGHT side. Dark mode: auto via CSS vars.
  Time: 1.5h

- [ ] **Create `components/commerce/OrdersTable.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/OrdersTable.tsx`
  shadcn `<Table>` with columns: #, Date, Items, Status badge, Total, Actions
  Responsive: cards layout on mobile (hide table). Row click → order detail.
  Actions: Track (active) / Reorder (delivered) / New Order (cancelled)
  Time: 1.5h

### 04.D — MVP Deploy Checklist (Day 15)

- [ ] **🔴 Run Lighthouse mobile audit**
  Target: ≥ 70 mobile
  Time: 15min

- [ ] **🔴 Test all 10 themes (light + dark)**
  Visit homepage with each theme → check header variant, section order, footer variant
  Both light + dark mode for each = 20 combinations
  Time: 1.5h

- [ ] **🔴 Test RTL (Arabic) on all pages**
  Switch to Arabic → verify header flips, text right-aligned, sheet from correct side
  Time: 30min

- [ ] **🔴 Test AJAX cart flow**
  Menu → add product → CartSheet opens → qty change → checkout button works
  Time: 15min

- [ ] **Deploy MVP**
  Push to production branch. Smoke test live.
  Time: 30min

---

## ═══════════════════════════════════════
## PHASE 05 — ALL PAGE REDESIGNS (Post-MVP)
## ═══════════════════════════════════════
> **Goal:** Full redesign of all 7 internal pages with shadcn/ui + 3 payment methods
> **Estimated Time:** 5–6 days

### 05.A — Menu Page

- [ ] **Rewrite `app/[locale]/menu/menu-page-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/menu/menu-page-content.tsx`
  Sticky category filter tabs (shadcn `<Tabs>`) + price range slider (`<Slider>`)
  Search input (client-side filter, no API). Sort dropdown. "Clear filters" button.
  Product grid: 3 cols desktop / 2 tablet / 1 mobile. Load more button.
  URL params sync: `?category=5&search=burger&sort=price_asc`
  Loading: grid of `<ProductCardSkeleton>`. Empty state with illustration.
  CartBottomBar visible at bottom.
  Time: 3h

- [ ] **Create `app/[locale]/menu/loading.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/menu/loading.tsx`
  9 ProductCardSkeletons in grid layout
  Time: 15min

### 05.B — Cart Page

- [ ] **Rewrite `app/[locale]/cart/cart-page-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/cart/cart-page-content.tsx`
  Desktop: 2-column (items left, sticky summary right). Mobile: single column, summary accordion.
  Cart items with qty controls. Order notes textarea. Summary card (items + delivery + total).
  Checkout button. "Continue Shopping" link. Empty state.
  Time: 1.5h

### 05.C — Checkout Page (3 Payment Methods)

- [ ] **Create `components/forms/CheckoutForm.tsx`**
  Path: `frontend/apps/restaurant-website/components/forms/CheckoutForm.tsx`
  RHF + Zod validated form. Fields: name, phone (Egyptian regex), address, notes, payment method.
  3 payment methods (scalable registry): Cash | Paymob Online | Screenshot Transfer
  Zod: phone `/^(01[0-2,5][0-9]{8})$/`, address min 10 chars
  Cash flow: POST /orders → clear cart → redirect to order detail + confetti
  Paymob flow: POST /orders → iframe dialog → poll status → redirect
  Screenshot flow: POST /orders → upload dialog → redirect with "Pending Verification" badge
  Paymob fallback: if API down → auto-switch to cash + toast
  Time: 3.5h

- [ ] **Rewrite `app/[locale]/checkout/checkout-page-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/checkout/checkout-page-content.tsx`
  Progress bar (4 steps: Cart → Info → Payment → Done)
  2-column desktop (form left, sticky summary right). Mobile: accordion summary.
  Guard: redirect to /cart if cart empty.
  Time: 1.5h

### 05.D — Orders History Page

- [ ] **Rewrite `app/[locale]/orders/orders-page-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/orders/orders-page-content.tsx`
  Filter tabs: All / Active / Delivered / Cancelled + count badges
  Search by order number. OrdersTable (desktop) / Card list (mobile).
  Pagination: shadcn Pagination. Empty state per tab.
  Time: 2h

### 05.E — Order Detail Page

- [ ] **Rewrite `app/[locale]/orders/[id]/order-detail-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/orders/[id]/order-detail-content.tsx`
  Header: back button + order # + status badge
  OrderTimeline (vertical, animated). Items breakdown table. Customer info card.
  Actions: Track (active) / Reorder + Rate (delivered) / New Order (cancelled)
  Screenshot pending: "Upload Receipt" button
  Polling: `refetchInterval: 15000` for active orders
  Time: 2h

### 05.F — Track Order Page

- [ ] **Create `components/forms/TrackOrderForm.tsx`**
  Path: `frontend/apps/restaurant-website/components/forms/TrackOrderForm.tsx`
  Centered card: order number input + "Track" button. Not found → error + "Order Now" link.
  Time: 45min

- [ ] **Rewrite `app/[locale]/orders/track/track-page-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/orders/track/track-page-content.tsx`
  Centered search card. Found → redirect to `/orders/{id}`.
  Time: 30min

### 05.G — Settings Page

- [ ] **Rewrite `app/[locale]/settings/settings-page-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/settings/settings-page-content.tsx`
  4 tabs: Language (radio cards with flags) | Appearance (Light/Dark/System radio cards) | Cart (item count, clear, saved address) | Notifications (coming soon badge)
  Time: 2h

### 05.H — Table Menu Page

- [ ] **Enhance `app/[locale]/table/[tableId]/menu/page.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/table/[tableId]/menu/page.tsx`
  Table header banner: `🍽️ طاولة رقم ٥ · Table #5` | Change Table | End Session
  Same menu layout. Pass `tableId` to cart store.
  Time: 1h

---

## ═══════════════════════════════════════
## PHASE 06 — CROSS-CUTTING FEATURES (Post-MVP)
## ═══════════════════════════════════════
> **Goal:** Search, cart sync, WhatsApp forwarding, dark mode audit
> **Estimated Time:** 3–4 days

### 06.A — Search Feature

- [ ] **Create `lib/store/search-store.ts`**
  Path: `frontend/apps/restaurant-website/lib/store/search-store.ts`
  State: `query`, `results: { products, categories }`, `isLoading`
  Time: 20min

- [ ] **Create `lib/hooks/use-search.ts`**
  Path: `frontend/apps/restaurant-website/lib/hooks/use-search.ts`
  300ms debounce, client-side filter of cached products. Returns: `{ results, isLoading, query, setQuery, clearSearch }`
  Time: 45min

- [ ] **Wire search URL params in menu page**
  Path: `frontend/apps/restaurant-website/app/[locale]/menu/menu-page-content.tsx`
  Read `?search=` and `?category=` on mount. On change → `router.replace()`.
  Time: 30min

### 06.B — Cart Cross-Tab Sync

- [ ] **Add BroadcastChannel sync to cart-store**
  Path: `frontend/apps/restaurant-website/lib/store/cart-store.ts`
  `BroadcastChannel` for instant cross-tab sync (fallback: `StorageEvent`)
  Time: 45min

- [ ] **🟠 Cross-browser test: cart cross-tab sync on Chrome + Firefox + Safari**
  Test matrix: Open 2 tabs → add item in tab 1 → verify tab 2 updates
  Chrome: `BroadcastChannel` (native support)
  Firefox: `BroadcastChannel` (native support)
  Safari: Verify `StorageEvent` fallback works (Safari BroadcastChannel support varies)
  Edge cases: tab opened before cart init, rapid add/remove, offline tab reconnect
  Time: 30min

### 06.C — Dark Mode Complete Audit

- [ ] **Audit ALL components for hardcoded colors**
  `grep -r "background.*#[0-9a-f]" components/` → fix each to CSS var
  `grep -r "color.*#[0-9a-f]" components/` → fix each
  Every color must be: `bg-theme-*` or `text-theme-*` or `hsl(var(--theme-*))`
  Time: 1h

- [ ] **Screenshot test: dark mode on ALL pages**
  All 8 pages × light + dark. All: Header, Footer, Cards, Forms, Modals.
  Time: 1h

### 06.D — Language Switcher Complete

- [ ] **Test language switch on all pages**
  EN→AR: URL changes, RTL applies, Arabic font loads, all text translated
  Cart preserved across language switch.
  Time: 1h

- [ ] **Proofread Arabic translations**
  Compare ar.json vs en.json — fill any gaps.
  Time: 1h

### 06.E — WhatsApp Integration Complete

- [ ] **Add WhatsApp order forwarding on checkout success**
  Path: `frontend/apps/restaurant-website/app/[locale]/checkout/checkout-page-content.tsx`
  On order success: auto-open `wa.me/{adminNumber}?text={formatted order details}`
  Format: Order #, Customer name, phone, address, items × qty = price, total, payment method
  Time: 1h

- [ ] **Add WhatsApp button to ProductCard**
  Path: `frontend/apps/restaurant-website/components/commerce/ProductCard.tsx`
  `wa.me/{number}?text=أريد أطلب ${product.name_ar} × ${qty}`
  Only shown if WhatsApp configured in theme
  Time: 30min

- [ ] **Add WhatsApp button to CTASection**
  Path: `frontend/apps/restaurant-website/components/sections/CTASection.tsx`
  Green WhatsApp button alongside "Order Now"
  Time: 15min

---

## ═══════════════════════════════════════
## PHASE 07 — PERFORMANCE + TESTING (Post-MVP)
## ═══════════════════════════════════════
> **Goal:** Lighthouse 95+ mobile, backend caching, E2E tests, PWA
> **Estimated Time:** 3–4 days

### 07.A — Backend Caching

- [ ] **🔴 Add in-memory theme cache (sync.Map)**
  Path: `backend/internal/handler/http/admin_theme_handler_v2.go`
  `sync.Map` with 5min TTL. Check cache → return cached → or query DB + store.
  Set `X-Cache: HIT/MISS` header. Impact: DB queries drop 95%.
  Time: 1.5h

- [ ] **🔴 Add in-memory product cache**
  Path: `backend/internal/handler/http/public_menu_handler.go`
  Same pattern, 10min TTL (products change less often).
  Time: 1h

- [ ] **🔴 Add cache invalidation on write**
  Path: `backend/internal/handler/http/admin_theme_handler_v2.go`
  When admin saves theme: `themeCache.Delete("theme:" + slug)`
  When admin updates products: `productCache.Delete("products:" + restaurantId)`
  Time: 45min

### 07.B — Frontend Performance

- [ ] **ISR for all public pages**
  Path: `app/[locale]/page.tsx`, `app/[locale]/menu/page.tsx`
  `export const revalidate = 300` (5 min). Orders/track: `revalidate = 0` (always fresh).
  Time: 30min

- [ ] **Parallel data fetching**
  Path: `app/[locale]/page.tsx`, `app/[locale]/menu/page.tsx`
  `Promise.all([fetchTheme, fetchProducts, fetchCategories])` — no waterfall.
  Time: 30min

- [ ] **Optimize images with next/image**
  All product images: `fill`, `sizes` prop, `quality={75}`, `placeholder="blur"`, `loading="lazy"` (below fold), `priority` (first 4).
  Time: 1h

- [ ] **Conditional Arabic font loading**
  Path: `frontend/apps/restaurant-website/app/[locale]/layout.tsx`
  Load Cairo ONLY when `locale === 'ar'`. English pages save 200KB.
  Time: 30min

- [ ] **Code splitting — dynamic imports**
  Dynamic: GallerySection, TestimonialsSection, CartSheet, PaymobDialog, Framer Motion
  Time: 45min

- [ ] **Bundle analyzer audit**
  `ANALYZE=true pnpm build`. Fix: tree-shake lucide imports, remove unused shadcn components.
  Target: no chunk > 80KB gzip.
  Time: 1h

- [ ] **Prefetch critical routes**
  Homepage → prefetch `/menu`. Menu → prefetch `/cart`. Cart → prefetch `/checkout`.
  Time: 15min

- [ ] **Loading skeletons for all pages**
  Path: `app/[locale]/menu/loading.tsx`, `app/[locale]/orders/loading.tsx`, `app/[locale]/cart/loading.tsx`
  Time: 45min

- [ ] **Preconnect to API domain**
  Path: `app/layout.tsx`
  `<link rel="preconnect" href="...">` + `<link rel="dns-prefetch" href="...">`
  Time: 10min

### 07.C — Error States

- [ ] **Create `components/system/GlobalErrorBoundary.tsx`**
  Path: `frontend/apps/restaurant-website/components/system/GlobalErrorBoundary.tsx`
  Class component. Catches unhandled errors. Sends to Sentry if configured.
  Renders: ErrorFallbackUI with "Retry" button.
  Time: 45min

- [ ] **Create `components/system/NetworkErrorState.tsx`**
  Path: `frontend/apps/restaurant-website/components/system/NetworkErrorState.tsx`
  WifiOff icon + "No Internet" + "Retry" + "Cart is saved offline" note
  Time: 30min

- [ ] **Create `components/system/PaymentErrorState.tsx`**
  Path: `frontend/apps/restaurant-website/components/system/PaymentErrorState.tsx`
  XCircle icon + "Payment Failed" + [Try Again] + [Pay Cash Instead]
  Time: 30min

- [ ] **Create `components/system/ProductUnavailableState.tsx`**
  Path: `frontend/apps/restaurant-website/components/system/ProductUnavailableState.tsx`
  Gray overlay on ProductCard + "غير متاح" badge + disabled buttons
  Time: 20min

### 07.D — Testing

- [ ] **Setup Vitest**
  Path: `frontend/apps/restaurant-website/vitest.config.ts`
  jsdom env, path aliases, setup file
  Time: 30min

- [ ] **Unit test: ProductCard**
  Path: `frontend/apps/restaurant-website/__tests__/components/ProductCard.test.tsx`
  Tests: renders name in locale, out-of-stock badge, add-to-cart AJAX, qty stepper
  Time: 1h

- [ ] **Unit test: cart-store**
  Path: `frontend/apps/restaurant-website/__tests__/stores/cart-store.test.ts`
  Tests: addItem, removeItem, updateQty, clearCart, totalPrice, cross-tab sync
  Time: 45min

- [ ] **Unit test: theme-normalizer**
  Path: `frontend/apps/restaurant-website/__tests__/hooks/theme-normalizer.test.ts`
  Tests: 5 different API shapes → valid ThemeConfig
  Time: 45min

- [ ] **E2E: Full checkout flow (Playwright)**
  Path: `frontend/apps/restaurant-website/e2e/checkout.spec.ts`
  Flow: home → menu → add product → CartSheet → checkout → fill form → cash → order detail
  Time: 2h

- [ ] **E2E: Language switch**
  Path: `frontend/apps/restaurant-website/e2e/language.spec.ts`
  Flow: /en/menu → toggle AR → verify /ar/menu + RTL + Arabic text
  Time: 1h

- [ ] **E2E: Dark mode**
  Path: `frontend/apps/restaurant-website/e2e/dark-mode.spec.ts`
  Flow: light → toggle dark → verify html.dark + header colors → reload → persists
  Time: 45min

### 07.E — PWA

- [ ] **Create `public/manifest.json`**
  Path: `frontend/apps/restaurant-website/public/manifest.json`
  name, short_name, icons (192+512), start_url, display:standalone, theme_color
  Time: 30min

- [ ] **Add PWA meta tags to layout**
  Path: `frontend/apps/restaurant-website/app/layout.tsx`
  `<link rel="manifest">`, `<meta name="theme-color">`, apple-touch-icon
  Time: 20min

---

## ═══════════════════════════════════════
## PHASE 08 — SEO + MONITORING + PRODUCTION (Post-MVP)
## ═══════════════════════════════════════
> **Goal:** Egyptian Google rankings, error tracking, analytics, production hardening
> **Estimated Time:** 2–3 days

### 08.A — SEO

- [ ] **Add dynamic metadata to all pages**
  Path: each `page.tsx` in `app/[locale]/`
  Homepage: restaurant name + "اطلب اونلاين" / "Order Online"
  Menu: "القائمة" / "Menu" + restaurant name
  Orders: `noindex: true` (private)
  Include: `alternates.languages` for en/ar
  Time: 1.5h

- [ ] **Add Schema.org JSON-LD**
  Path: `frontend/apps/restaurant-website/app/[locale]/layout.tsx`
  Restaurant schema + MenuItem schema for rich results.
  Time: 1h

- [ ] **Add hreflang tags**
  Path: `frontend/apps/restaurant-website/app/[locale]/layout.tsx`
  `ar` + `en` + `x-default`
  Time: 20min

- [ ] **Create `app/sitemap.ts`**
  Path: `frontend/apps/restaurant-website/app/sitemap.ts`
  Dynamic: `/{locale}`, `/{locale}/menu`, `/{locale}/orders/track`
  Time: 30min

- [ ] **Create `app/robots.ts`**
  Path: `frontend/apps/restaurant-website/app/robots.ts`
  Allow: /, /menu, /orders/track. Disallow: /orders, /checkout, /settings.
  Time: 15min

### 08.B — Monitoring

- [ ] **Add Sentry error tracking**
  Path: `frontend/apps/restaurant-website/sentry.client.config.ts` + `sentry.server.config.ts`
  DSN from `NEXT_PUBLIC_SENTRY_DSN` env var. Free tier: 5K errors/month.
  Time: 1h

- [ ] **Add Google Analytics 4 (optional, per-tenant)**
  Path: `frontend/apps/restaurant-website/app/layout.tsx`
  Only load if `NEXT_PUBLIC_GA_ID` is set. Track: page views, add-to-cart, checkout, orders.
  Time: 30min

### 08.C — Final Production Checklist

- [ ] **Final Lighthouse audit (AFTER)**
  Compare with BEFORE screenshots. Target: Mobile ≥ 95, Desktop ≥ 98.
  Time: 45min

- [ ] **Test all 10 themes × light + dark = 20 combinations**
  Visit homepage → check header variant, section order, footer variant, colors.
  Time: 1.5h

- [ ] **Test on Egypt 4G throttling**
  Chrome DevTools → Fast 4G. Goal: LCP < 2.5s, FID < 100ms, CLS < 0.1.
  Time: 30min

- [ ] **Verify theme builder preview mode**
  Dashboard editor → iframe loads → change color → instant update → dark toggle
  Time: 30min

---

## ═══════════════════════════════════════
## DASHBOARD TASKS (Parallel with Phases 01–02)
## ═══════════════════════════════════════
> These tasks update the Dashboard app to support the new theme features

### Dashboard: Dark Mode Editor

- [ ] **Create `DarkColorsEditor.tsx`**
  Path: `frontend/apps/dashboard/src/components/editors/DarkColorsEditor.tsx`
  Side-by-side light vs dark color pickers. "Auto-Generate Dark Palette" button.
  Live preview toggle. ContrastChecker integration. Saves to themes_v2 dark columns.
  Time: 2.5h

- [ ] **Update `themeBuilderStore.ts` — dark color state**
  Path: `frontend/apps/dashboard/src/stores/themeBuilderStore.ts`
  Add: `colorsDark`, `setDarkColor()`, `autoGenerateDarkPalette()`
  Time: 45min

- [ ] **Update `themeApi.ts` — dark colors + variants**
  Path: `frontend/apps/dashboard/src/lib/api/themeApi.ts`
  Add: dark color fields + template/variant fields in normalize/denormalize
  Time: 30min

- [ ] **Add dark mode preview toggle to theme editor**
  Path: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx`
  [☀️ Light | 🌙 Dark] toggle in preview toolbar
  Time: 45min

### Dashboard: Template Selector

- [ ] **Add template selector to theme editor**
  Path: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx`
  Tabs or dropdown: 10 templates with preview thumbnail.
  On select: auto-updates section order + variant selections.
  Time: 1.5h

### Dashboard: Section DnD Upgrade

- [ ] **Upgrade DnD to @dnd-kit**
  Path: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/components/ComponentBuilder.tsx`
  ```bash
  pnpm add @dnd-kit/core @dnd-kit/sortable @dnd-kit/utilities
  ```
  Replace native HTML5 DnD: smooth animations, visual drop indicators, keyboard accessibility, touch support, drag overlay ghost.
  Time: 2h

- [ ] **Add "Add Section" dropdown with all 9 types**
  Path: same ComponentBuilder.tsx
  Types: Hero, Featured Products, Categories, Why Choose Us, Gallery, Testimonials, CTA, Contact, Custom
  Each: icon + name + short description
  Time: 45min

### Dashboard: New Section Editors

- [ ] **Create `CategoriesEditor.tsx`**
  Path: `frontend/apps/dashboard/src/components/editors/CategoriesEditor.tsx`
  Fields: title_en/ar, subtitle_en/ar, layout (scroll|grid), max_items, show_count
  Time: 1h

- [ ] **Create `GalleryEditor.tsx`**
  Path: `frontend/apps/dashboard/src/components/editors/GalleryEditor.tsx`
  Fields: title_en/ar, images[] (upload/URL), layout (grid|masonry), columns, lightbox toggle
  Time: 1h

---

## PERFORMANCE SUMMARY

| Optimization | Before | After |
|---|---|---|
| Backend DB queries (10K restaurants) | ~350K QPS | ~7 QPS |
| Theme API response time | 200ms | 5ms (cache hit) |
| Homepage load (Egypt 4G) | 3–4s | 1.2s |
| Menu page load | 2–3s | 0.8s |
| Navigation (page to page) | 500ms+ | < 100ms |
| Product image per item | ~500KB | ~50KB |
| JS bundle (gzip) | ~450KB | < 180KB |
| Arabic page extra font | +200KB always | +0KB (conditional) |
| LCP (mobile) | 4.5s | < 2.0s |
| Lighthouse mobile | ~60 | ≥ 95 |
| API calls per user session | 15–20 | 3–5 first visit, 0 repeat |
| Theme visual variety | 10% (color swap) | 85% (templates + variants) |

---

*Last Updated: 2026-03-03 | Status: Planning Phase | MVP: 15 working days*
