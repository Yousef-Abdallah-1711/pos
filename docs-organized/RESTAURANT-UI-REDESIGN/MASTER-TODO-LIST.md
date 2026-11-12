# Restaurant Website UI Redesign â€” Master Todo List v3.0
**All 9 Phases Â· ~180 Tasks Â· Production-Grade Â· MVP-First**
**Updated:** 2026-03-03 v3.0
**MVP:** Phases 00â€“04 (15 working days) â†’ Deploy â†’ Iterate Phases 05â€“08
**Quality:** Production-ready, 1000+ themes, scalable, Egyptian market optimized

> **Legend:** ðŸ”´ Blocking Â· ðŸŸ  High Â· ðŸŸ¡ Medium Â· ðŸŸ¢ Nice-to-have
> **Format:** `- [ ] Task` / `Path:` / `Why:` / `Test:` / `Time:`

---

## CHECKPOINT SCHEDULE
| Day | Checkpoint | Pass Criteria |
|-----|-----------|---------------|
| 3 | Header renders all 10 themes (light + dark) | No hardcoded hex, RTL flips correctly |
| 7 | Homepage renders 8 sections, templates work | Each theme has unique section order |
| 8 | Mobile Safari iOS test | `safe-area-inset-bottom` renders correctly, no viewport overflow, touch targets â‰¥ 48px |
| 10 | Menu page + ProductCard + AJAX cart working | Add â†’ sheet opens, no page reload |
| 12 | Cart + Checkout MVP functional | Cash payment completes, order created |
| 15 | Production deploy | Lighthouse mobile â‰¥ 70, all pages render |

---

## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
## PHASE 00 â€” BACKUP + SETUP (Day 1â€“2)
## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
> **Goal:** Safe backup, install foundations, test dependencies, establish conventions
> **Checkpoint Day 2:** All packages installed, types compile, tailwind configured

### 00.PRE â€” Backup & Baseline (ðŸ”´ MUST RUN FIRST)

- [x] **ðŸ”´ Create git backup branch + working branch**
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

- [x] **ðŸ”´ Take Lighthouse baseline screenshots**
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

- [x] **ðŸ”´ Document BEFORE-STATE.md**
  Path: `docs-organized/RESTAURANT-UI-REDESIGN/BEFORE-STATE.md`
  Record: git SHA, Lighthouse scores, page count, component count, bundle size
  Time: 10min

### 00.A â€” Package Installation

- [x] **Initialize shadcn/ui**
  Path: `frontend/apps/restaurant-website/`
  ```bash
  cd frontend/apps/restaurant-website
  pnpm dlx shadcn@latest init
  ```
  Choose: TypeScript YES, slate base, CSS variables YES, @/ aliases
  Creates: `components/ui/`, `components.json`, `lib/utils.ts`
  Time: 15min

- [x] **Install all shadcn components**
  Path: `frontend/apps/restaurant-website/`
  ```bash
  npx shadcn@latest add button card input label form textarea select \
    dialog drawer sheet badge tabs avatar separator scroll-area skeleton \
    progress tooltip dropdown-menu table slider accordion alert-dialog \
    popover command toggle-group radio-group
  ```
  Time: 15min

- [x] **Install additional packages**
  Path: `frontend/apps/restaurant-website/`
  ```bash
  pnpm add embla-carousel-react embla-carousel-autoplay
  pnpm add next-themes
  pnpm add framer-motion
  pnpm add sonner
  pnpm add -D @next/bundle-analyzer
  ```
  Time: 10min

- [x] **Install Cairo Arabic font**
  Path: `frontend/apps/restaurant-website/app/layout.tsx`
  ```typescript
  import { Inter, Cairo } from 'next/font/google'
  const inter = Inter({ subsets: ['latin'], variable: '--font-inter', display: 'swap' })
  const cairo = Cairo({ subsets: ['arabic', 'latin'], variable: '--font-cairo', display: 'swap' })
  ```
  Time: 10min

### 00.B â€” Dependency Sandbox Test (Risk Mitigation)

- [x] **ðŸ”´ Test shadcn/ui RTL with Arabic text**
  Path: create temp `frontend/apps/restaurant-website/__sandbox__/rtl-test.tsx`
  Test: Render `<Tabs>`, `<Sheet>`, `<DropdownMenu>`, `<Dialog>` in RTL mode
  Fix any RTL quirks discovered (document in BEFORE-STATE.md)
  Delete sandbox after testing
  Time: 30min

- [x] **Test Embla Carousel in RTL**
  Path: temp sandbox component
  Test: direction="rtl" prop, scroll-snap, autoplay
  Verify: no loop mode bugs, containScroll works
  Time: 15min

- [x] **Test next-themes hydration** _(completed via Phase 01 provider rewrite + `suppressHydrationWarning` + `disableTransitionOnChange`)_
  Path: temp sandbox
  Test: SSR â†’ client hydration, no flash of wrong theme
  Verify: `suppressHydrationWarning` + `disableTransitionOnChange` work
  Time: 15min

### 00.C â€” Configuration Files

- [x] **Update `tailwind.config.js`**
  Path: `frontend/apps/restaurant-website/tailwind.config.js`
  Changes:
  - `darkMode: ['class']` â€” manual toggle via `html.dark`
  - `fontFamily.sans` â†’ `var(--font-inter)`
  - `fontFamily.arabic` â†’ `var(--font-cairo)`
  - `colors.theme.*` â†’ all CSS var-based colors (primary, surface, text, border, etc.)
  - `minHeight['screen-svh']` â†’ `100svh`
  Time: 30min

- [x] **Rewrite `styles/globals.css`**
  Path: `frontend/apps/restaurant-website/styles/globals.css`
  Content: Tailwind directives + `:root` CSS vars + `.dark` overrides + shadcn bridge vars
  Time: 45min

- [x] **Create `styles/theme-variables.css`**
  Path: `frontend/apps/restaurant-website/styles/theme-variables.css`
  Content: All 25+ CSS variables with `.dark` overrides (see MASTER-OVERVIEW CSS section)
  Time: 30min

- [x] **Update `next.config.js`**
  Path: `frontend/apps/restaurant-website/next.config.js`
  Add: image domains, `optimizePackageImports` for lucide/framer-motion/embla, security headers
  Time: 20min

### 00.D â€” TypeScript Types

- [x] **Create `types/theme.ts`**
  Path: `frontend/apps/restaurant-website/types/theme.ts`
  Export: `ThemeConfig`, `ThemeHeaderConfig`, `ThemeFooterConfig`, `ThemeSectionConfig`, `SectionType`, `ThemeNavItem`, `ThemeSocialLink`, `ThemeFooterSection`, `ThemeLegalLink`
  Include: `templateId`, `headerVariant`, `footerVariant`, `componentOverrides` fields
  Test: `tsc --noEmit` passes
  Time: 45min

- [x] **Create `types/payment.ts`**
  Path: `frontend/apps/restaurant-website/types/payment.ts`
  Export: `PaymentMethodType` (`cash` | `paymob_online` | `screenshot_transfer`), `PaymentMethod`, `PaymobCheckoutResponse`, `PaymentScreenshot`
  Time: 30min

### 00.E â€” Create Directories + Delete Old Files

- [x] **Create all new directories**
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

- [x] **Delete replaced files** _(completed by migrating to `Header.tsx`/`Footer.tsx` and `home-page-content.tsx`)_
  ```
  components/layout/header-ssr.tsx       â†’ replaced by new Header.tsx
  components/layout/footer-ssr.tsx       â†’ replaced by new Footer.tsx
  components/layout/ThemeHeader.tsx       â†’ merged into Header.tsx
  components/layout/ThemeFooter.tsx       â†’ merged into Footer.tsx
  components/layout/theme-layout-normalizers.ts â†’ logic in Header/Footer
  components/homepage/dynamic-home-page.tsx     â†’ rebuilt in page.tsx
  components/homepage/featured-products-section.tsx â†’ new section component
  components/homepage/error-boundary.tsx         â†’ GlobalErrorBoundary
  ```
  Time: 5min

---

## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
## PHASE 01 â€” THEME SYSTEM V2 (Day 3â€“5)
## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
> **Goal:** Per-theme dark mode, CSS vars, theme templates, caching. Scalable to 1000 themes.
> **Checkpoint Day 5:** Toggle dark mode â†’ header/footer/cards all change per theme.

### 01.A â€” Database: Dark Mode Colors

- [x] **?? Create migration: add dark color columns**
  Path: `backend/migrations/161_add_theme_dark_colors.sql`
  Add columns: `primary_dark`, `secondary_dark`, `accent_dark`, `background_dark`, `text_dark`, `border_dark`, `shadow_dark`, `surface_color`, `surface_dark`, `header_bg_dark`, `header_text_dark`, `footer_bg_dark`, `footer_text_dark`, `dark_mode_enabled`
  Test: `SELECT primary_dark FROM themes_v2 LIMIT 1;`
  Time: 30min

- [x] **?? Create migration: add template + variant columns**
  Path: `backend/migrations/162_add_theme_template_variants.sql`
  Add columns: `template_id VARCHAR(20) DEFAULT 'template-01'`, `header_variant VARCHAR(20) DEFAULT 'solid'`, `footer_variant VARCHAR(20) DEFAULT 'expanded'`, `hero_variant VARCHAR(20) DEFAULT 'gradient'`, `products_variant VARCHAR(20) DEFAULT 'carousel'`
  Time: 20min

- [x] **?? Seed dark colors + templates for all 10 themes**
  Path: `backend/migrations/163_seed_theme_dark_colors_and_templates.sql`
  Each theme gets: unique dark palette + template assignment + variant assignments
  Time: 45min

- [x] **?? Update backend theme handler — return dark colors + variants**
  Path: `backend/internal/handler/http/admin_theme_handler_v2.go`
  Changes: Include dark color fields + template/variant fields in GET response + accept in PUT/POST
  Time: 45min

### 01.B â€” ThemeProvider Rewrite (Core)

- [x] **?? Rewrite ThemeProvider**
  Path: `frontend/apps/restaurant-website/components/providers/theme-provider.tsx`
  Architecture:
  - Wrap with `NextThemesProvider` (`attribute="class"`, `defaultTheme="system"`, `disableTransitionOnChange`)
  - `ThemeCSSInjector` component: reads `themeConfig.colors` + `resolvedTheme` â†’ injects CSS vars
  - Uses `hexToHSL()` to convert hex â†’ HSL for CSS vars
  - Re-runs ONLY when `themeConfig.slug`, `themeConfig.version`, or `isDark` changes
  - Single batch `style.setProperty()` loop â€” no layout thrashing
  Test: Toggle dark â†’ all colors update. Switch theme â†’ all colors update.
  Time: 2h

- [x] **Create `lib/utils/hex-to-hsl.ts`**
  Path: `frontend/apps/restaurant-website/lib/utils/hex-to-hsl.ts`
  Export: `hexToHSL(hex: string): string` â€” converts `#2563eb` â†’ `221 83% 53%`
  Export: `memoHexToHSL` â€” memoized version
  Test: `hexToHSL('#2563eb')` â†’ `'221 83% 53%'`
  Time: 30min

- [x] **Create `lib/utils/theme-normalizer.ts`**
  Path: `frontend/apps/restaurant-website/lib/utils/theme-normalizer.ts`
  Export: `normalizeThemeConfig(apiData: unknown, locale: 'en' | 'ar'): ThemeConfig`
  Handles: snake_case â†’ camelCase, missing fields â†’ defaults, bilingual extraction, nested structures
  Test: Unit test with 5 different API response shapes
  Time: 1.5h

- [x] **Update `lib/utils/default-theme.ts`**
  Path: `frontend/apps/restaurant-website/lib/utils/default-theme.ts`
  Add: ALL dark mode color pairs, header/footer dark colors, semantic colors, template defaults
  Export: `DEFAULT_THEME: ThemeConfig`
  Time: 1h

### 01.C â€” Theme Templates Registry

- [x] **?? Create `lib/utils/theme-templates.ts`**
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

### 01.D â€” Preferences Store + Dark Mode Hook

- [x] **Update `lib/store/preferences-store.ts`**
  Path: `frontend/apps/restaurant-website/lib/store/preferences-store.ts`
  Add: `darkMode: 'light' | 'dark' | 'system'`, `setDarkMode()`, sync with next-themes
  Time: 20min

- [x] **Create `lib/hooks/use-dark-mode.ts`**
  Path: `frontend/apps/restaurant-website/lib/hooks/use-dark-mode.ts`
  Returns: `{ mode, setMode, isDark, toggle }` â€” wraps next-themes + preferences store
  Time: 30min

### 01.E â€” i18n Keys (ALL new features)

- [x] **Extend `i18n/messages/en.json`**
  Path: `frontend/apps/restaurant-website/i18n/messages/en.json`
  Add ~200 keys: header, theme, language, cart, product, search, menu, checkout, orders, order_status, order_detail, track, settings, sections, errors, payment, whatsapp, common
  Time: 1h

- [x] **Extend `i18n/messages/ar.json`**
  Path: `frontend/apps/restaurant-website/i18n/messages/ar.json`
  Full Arabic translations for ALL keys. Critical:
  - `cart.title`: Ø³Ù„Ø© Ø§Ù„ØªØ³ÙˆÙ‚
  - `cart.empty`: Ø³Ù„ØªÙƒ ÙØ§Ø±ØºØ©
  - `product.add_to_cart`: Ø£Ø¶Ù
  - `checkout.place_order`: ØªØ£ÙƒÙŠØ¯ Ø§Ù„Ø·Ù„Ø¨
  - `payment.cash`: Ø§Ù„Ø¯ÙØ¹ Ø¹Ù†Ø¯ Ø§Ù„Ø§Ø³ØªÙ„Ø§Ù…
  - `whatsapp.order`: Ø§Ø·Ù„Ø¨ Ø¹Ø¨Ø± ÙˆØ§ØªØ³Ø§Ø¨
  Time: 1.5h

### 01.F â€” Format Utilities

- [x] **Create `lib/utils/format.ts`**
  Path: `frontend/apps/restaurant-website/lib/utils/format.ts`
  Export:
  - `formatPrice(amount: number, locale: string): string` â†’ `Ù¡Ù¤Ù¥ Ø¬.Ù…` or `145 EGP`
  - `formatArabicNumber(num: number): string` â†’ Arabic numerals
  - `formatDate(date: string, locale: string): string` â†’ localized date
  - `formatRating(rating: number, count: number, locale: string): string` â†’ `Ù¤.Ù¨â­ (Ù¡Ù¢Ù£ ØªÙ‚ÙŠÙŠÙ…)`
  Time: 45min

### 01.G â€” Section Registry Update

- [x] **Register categories + gallery in component registry**
  Path: `frontend/layers/ui-themes/src/registry/componentRegistry.ts`
  Add: `categories` (aliases: `categories-section`, `browse_menu`) + `gallery` (aliases: `gallery-section`, `food_gallery`)
  Keep backward compatibility with old type aliases
  Time: 30min

### 01.H â€” Backend Cache Invalidation (ðŸ”´ MUST COMPLETE BEFORE FRONTEND CACHING)

- [x] **?? Add backend sync.Map cache + invalidation on admin writes**
  Path: `backend/internal/handler/http/public_homepage_handler.go`
  Strategy: `sync.Map` keyed by `tenantID:themeID` â†’ cached JSON response
  On admin write (create/update/delete theme): delete cache key + return `Cache-Control: no-cache`
  Also add invalidation hooks to: `admin_theme_handler_v2.go`, `admin_menu_handler.go`
  Test: Update theme in admin â†’ public API returns new data within 1s (no stale cache)
  Time: 1h

- [x] **Add Cache-Control headers to public API endpoints**
  Path: `backend/internal/handler/http/admin_theme_handler_v2.go`
  Public: `Cache-Control: public, max-age=300, stale-while-revalidate=600`
  Admin: `Cache-Control: no-cache, private`
  Also add to: `public_menu_handler.go`, `public_order_handler.go`
  Time: 30min

- [x] **Add React Query cache defaults**
  Path: `frontend/apps/restaurant-website/components/providers/core/query-provider.tsx`
  Set: `staleTime: 10min`, `gcTime: 1hr`, `refetchOnWindowFocus: false`, `retry: 2`
  Depends on: Backend cache invalidation must be working first
  Time: 20min

---

## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
## PHASE 02 â€” HEADER + FOOTER (Day 6â€“8)
## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
> **Goal:** 5 header variants, 5 footer variants, mobile nav, WhatsApp bubble
> **Checkpoint Day 8:** Header with search/cart/lang/dark works across all themes + RTL

### 02.A â€” Header Sub-Components

- [x] **Create `components/layout/header/DarkModeToggle.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/header/DarkModeToggle.tsx`
  Cycles: light â†’ dark â†’ system. Icons: Sun/Moon/Monitor. Size: 48px touch target.
  Uses: `useDarkMode()` hook, Lucide icons, shadcn `<Tooltip>`
  Test: Click 3x â†’ cycles all modes, `html.dark` class updates
  Time: 30min

- [x] **Create `components/layout/header/LanguageToggle.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/header/LanguageToggle.tsx`
  Desktop: `[ðŸ‡ºðŸ‡¸ EN â–¾]` dropdown. Mobile: compact `[EN]`/`[AR]`.
  On switch: replace `/en/` â†” `/ar/` in pathname, `router.push()`, save to prefs
  Uses: shadcn `<DropdownMenu>`, `useRouter()`, `usePathname()`
  Test: ENâ†’AR: URL changes, RTL applies, Arabic font loads
  Time: 45min

- [x] **Create `components/layout/header/CartButton.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/header/CartButton.tsx`
  ShoppingCart icon + animated count badge. Click â†’ opens CartSheet.
  Badge: `bg-theme-error text-white`, scale bounce on change. Size: 48px touch.
  Uses: Zustand cart store, shadcn `<Badge>`, Lucide ShoppingCart
  Test: Add item â†’ badge appears, remove all â†’ badge disappears
  Time: 30min

- [x] **Create `components/layout/header/SearchBar.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/header/SearchBar.tsx`
  Desktop: expanding input with results `<Popover>`. Mobile: icon â†’ full-width overlay.
  300ms debounce â†’ filter cached products (no API call). Max 5 results per group.
  Results: Products (image + name AR/EN + price) | Categories. Keyboard navigation.
  Uses: shadcn `<Popover>` + `<Command>`, React Query cached products
  Test: Type "burger" â†’ results < 300ms, click â†’ navigate to menu
  Time: 2.5h

- [x] **Create `components/layout/header/WhatsAppButton.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/header/WhatsAppButton.tsx`
  Desktop: icon + "Order via WhatsApp". Mobile: green icon only.
  Click â†’ `wa.me/{number}?text=Ù…Ø±Ø­Ø¨Ø§ØŒ Ø£Ø±ÙŠØ¯ Ø£Ù† Ø£Ø·Ù„Ø¨`
  Hidden if `!theme.header.showWhatsAppButton`
  Time: 20min

### 02.B â€” Header Variants (5 Types)

- [x] **ðŸ”´ Create `components/layout/headers/GlassHeader.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/headers/GlassHeader.tsx`
  Style: `backdrop-filter: blur(20px)`, `bg-[hsl(var(--theme-header-bg)/0.1)]`, transparent
  Used by: Modern Bistro, Coastal Ocean
  Time: 1.5h

- [x] **ðŸ”´ Create `components/layout/headers/GradientHeader.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/headers/GradientHeader.tsx`
  Style: `background: linear-gradient(135deg, primary â†’ accent)`, opaque
  Used by: Luxury Gold
  Time: 1h

- [x] **ðŸ”´ Create `components/layout/headers/SolidHeader.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/headers/SolidHeader.tsx`
  Style: `bg-[hsl(var(--theme-header-bg))]`, solid color, default variant
  Used by: Warm Rustic, Garden Green, Spicy Red, Egyptian Heritage
  Time: 1h

- [x] **Create `components/layout/headers/FloatingHeader.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/headers/FloatingHeader.tsx`
  Style: fixed, heavy `shadow-xl`, rounded bottom corners
  Used by: Urban Gray, Vibrant Pop
  Time: 45min

- [x] **Create `components/layout/headers/MinimalHeader.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/headers/MinimalHeader.tsx`
  Style: thin `border-b` only, transparent background, minimal presence
  Used by: Minimal White
  Time: 45min

### 02.C â€” Main Header (Dispatcher)

- [x] **ðŸ”´ Create `components/layout/Header.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/Header.tsx`
  Props: `themeConfig`, `locale`, `restaurantSlug`
  Logic: reads `themeConfig.headerVariant` â†’ renders matching header variant
  All variants share: same sub-components (Search, Cart, Language, DarkMode, WhatsApp)
  All variants: sticky, z-60, RTL-aware (`flex-row-reverse` when Arabic), mobile sheet drawer
  ALL colors from CSS vars â€” responds to dark mode toggle instantly
  Time: 1.5h

### 02.D â€” Footer Variants (5 Types)

- [x] **ðŸ”´ Create `components/layout/footers/ExpandedFooter.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/footers/ExpandedFooter.tsx`
  4-column: Company info | Quick links | Hours | Contact. Social row. Legal row. Back to top.
  Used by: Modern Bistro, Spicy Red, Egyptian Heritage
  Time: 1.5h

- [x] **Create `components/layout/footers/CompactFooter.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/footers/CompactFooter.tsx`
  2-row: Name + contact | Copyright + legal. Compact.
  Used by: Warm Rustic, Urban Gray
  Time: 45min

- [x] **Create `components/layout/footers/CardsFooter.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/footers/CardsFooter.tsx`
  3 card blocks: About | Contact | Menu. Card-based layout.
  Used by: Luxury Gold
  Time: 1h

- [x] **Create `components/layout/footers/SocialHeavyFooter.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/footers/SocialHeavyFooter.tsx`
  Social links prominent, large icons, minimal text.
  Used by: Coastal Ocean, Vibrant Pop
  Time: 45min

- [x] **Create `components/layout/footers/MinimalFooter.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/footers/MinimalFooter.tsx`
  Single line: `Â© 2026 Restaurant Name | Privacy | Terms`
  Used by: Minimal White, Garden Green
  Time: 30min

### 02.E â€” Main Footer (Dispatcher)

- [x] **Create `components/layout/Footer.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/Footer.tsx`
  Props: `themeConfig`, `locale`
  Logic: reads `themeConfig.footer.layout` â†’ renders matching footer variant
  ALL colors from CSS vars. RTL: `dir` prop from locale.
  Time: 45min

### 02.F â€” Mobile Navigation

- [x] **ðŸ”´ Create `components/layout/MobileBottomNav.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/MobileBottomNav.tsx`
  5 items: Home / Menu / Cart(count) / Orders / Settings
  Active state: current route â†’ theme primary color + dot indicator
  Cart badge: live count from cart store. `md:hidden`.
  `pb-[env(safe-area-inset-bottom)]` for iOS. Height: 64px + safe area.
  Time: 1h

- [x] **Create `components/layout/CartBottomBar.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/CartBottomBar.tsx`
  Sticky above MobileBottomNav. Shows: cart count + total + "Checkout â†’" button.
  Hidden when cart empty, slides up with animation when items added.
  Shown on: Home, Menu pages. NOT on Cart/Checkout pages.
  Time: 1h

- [x] **ðŸ”´ Create `components/layout/WhatsAppFloating.tsx`**
  Path: `frontend/apps/restaurant-website/components/layout/WhatsAppFloating.tsx`
  Fixed bottom-right, green circle, WhatsApp icon, subtle pulse animation.
  Always visible (z-50, above MobileBottomNav). Click â†’ `wa.me/{number}`.
  Hidden if no whatsAppNumber in theme config.
  Time: 30min

### 02.G â€” Layout Integration

- [x] **ðŸ”´ Update `app/[locale]/layout.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/layout.tsx`
  Changes:
  - Fetch theme server-side (SSR)
  - Render: `<Header>` + `{children}` + `<Footer>` + `<MobileBottomNav>` + `<CartBottomBar>` + `<WhatsAppFloating>`
  - Preview mode: if `?preview=true` â†’ hide Header/Footer/Nav
  - Conditional Cairo font: only load when `locale === 'ar'`
  Time: 1h

---

## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
## PHASE 03 â€” HOMEPAGE 8 SECTIONS (Day 9â€“11)
## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
> **Goal:** 8 dynamic sections with theme-specific templates, all themes + dark mode + RTL
> **Checkpoint Day 11:** Each theme renders unique homepage with different section order

### 03.A â€” Hero Variants (4 Types)

- [ ] **ðŸ”´ Create `components/sections/heroes/GradientHero.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/heroes/GradientHero.tsx`
  Full viewport `min-h-[100svh]`. CSS gradient using theme primary â†’ accent.
  Dark overlay for text contrast. Content: headline AR, subtitle EN, CTA button, delivery badge.
  Dual CTA: ðŸŸ¢ Ø§Ø·Ù„Ø¨ ÙˆØ§ØªØ³Ø§Ø¨ + ðŸŸ  Order Now
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

### 03.B â€” Hero Dispatcher

- [ ] **ðŸ”´ Create `components/sections/HeroSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/HeroSection.tsx`
  Props: `config`, `locale`, `themeConfig`
  Logic: reads `themeConfig.componentOverrides.hero` â†’ renders matching hero variant
  Fallback: GradientHero if variant not found
  Time: 30min

### 03.C â€” Featured Products Variants (3 Types)

- [ ] **ðŸ”´ Create `components/sections/products/ProductCarousel.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/products/ProductCarousel.tsx`
  Embla Carousel: 1 mobile / 2 tablet / 3-4 desktop. Autoplay 4s, pause on hover.
  Dot pagination. "View Full Menu â†’" link.
  Time: 1.5h

- [ ] **Create `components/sections/products/ProductGrid.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/products/ProductGrid.tsx`
  3-col desktop, 2-col tablet, 1-col mobile. Static grid with 6-8 products.
  Time: 1h

- [ ] **Create `components/sections/products/ProductMasonry.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/products/ProductMasonry.tsx`
  Masonry layout with varying card heights. Used by: Warm Rustic.
  Time: 1h

### 03.D â€” Featured Products Dispatcher

- [ ] **Create `components/sections/FeaturedProductsSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/FeaturedProductsSection.tsx`
  Reads `themeConfig.componentOverrides.products` â†’ renders matching variant
  Real products from `useProducts(restaurantSlug)`. Loading: skeleton cards.
  Time: 30min

### 03.E â€” Other Sections

- [ ] **Create `components/sections/WhyChooseUsSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/WhyChooseUsSection.tsx`
  4 cards grid (2Ã—2 mobile, 4Ã—1 desktop). Lucide icons. Hover lift. Dark mode via CSS vars.
  Time: 1.5h

- [ ] **Create `components/sections/TestimonialsSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/TestimonialsSection.tsx`
  Embla Carousel. Avatar + Name + Stars (`â­â­â­â­â­`) + Review text AR/EN.
  Auto-play with pause on hover. Arabic star format: `Ù¤.Ù¨â­`
  Time: 1.5h

- [ ] **Create `components/sections/CategoriesSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/CategoriesSection.tsx`
  Horizontal scroll pills. Emoji/icon + Category name AR/EN + item count badge.
  Click â†’ navigate to `/[locale]/menu?category={id}`. Scroll-snap on mobile.
  Real categories from `useCategories(restaurantSlug)`.
  Time: 1.5h

- [ ] **Create `components/sections/GallerySection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/GallerySection.tsx`
  3Ã—2 image grid. Lightbox on click (shadcn `<Dialog>`). Hover: `scale(1.05)`.
  Images from theme config. Dark mode: dark overlay.
  Time: 2h

- [ ] **Create `components/sections/CTASection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/CTASection.tsx`
  Full width gradient (primary â†’ accent). Headline + sub + dual CTA:
  ðŸŸ  "Order Now" button + ðŸŸ¢ "Order via WhatsApp" button.
  Pulse animation on main button. Dark mode: gradient shifts darker.
  Time: 1h

- [ ] **Create `components/sections/ContactSection.tsx`**
  Path: `frontend/apps/restaurant-website/components/sections/ContactSection.tsx`
  4 cards: Phone (tel: link) | WhatsApp (wa.me/) | Location (Google Maps) | Hours.
  Each: Icon + Label AR/EN + Value + Action button. Mobile: 2Ã—2 grid.
  Time: 1h

### 03.F â€” Dynamic Homepage (SSR + Template Rendering)

- [ ] **ðŸ”´ Rewrite `app/[locale]/page.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/page.tsx`
  Server Component: `Promise.all([fetchTheme, fetchProducts, fetchCategories])` â€” parallel fetch
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

### 03.G â€” Theme Builder Preview Mode

- [ ] **Add preview mode to layout**
  Path: `frontend/apps/restaurant-website/app/[locale]/layout.tsx`
  If `?preview=true`: hide Header, Footer, MobileBottomNav, CartBottomBar, WhatsAppFloating
  Listen for `postMessage` â†’ update CSS vars in < 50ms
  Time: 1h

---

## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
## PHASE 04 â€” COMMERCE COMPONENTS (Day 12â€“14)
## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
> **Goal:** ProductCard with AJAX cart, CartSheet sidebar, OrdersTable
> **Checkpoint Day 14:** Add to cart â†’ sheet opens â†’ checkout works

### 04.A â€” Product Card

- [x] **🔴 Create `components/commerce/ProductCard.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/ProductCard.tsx`
  Image (next/image, blur placeholder) | Name AR (bold) + EN (muted) | Rating `Ù¤.Ù£â­`
  Description (2 lines truncated) | Price badge `Ù¤Ù¥ Ø¬.Ù…` | Qty stepper `[-][1][+]`
  Full-width "ðŸ›’ Ø£Ø¶Ù / Add" button | Optional WhatsApp button
  States: default, hover (lift), out-of-stock (grayscale + badge), added (checkmark toast)
  Dark mode: bg-theme-surface, border-theme-border
  AJAX: click Add â†’ `cartStore.addItem()` â†’ toast â†’ CartSheet opens (NO page reload)
  Time: 2h

- [x] **Create `components/commerce/ProductCardSkeleton.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/ProductCardSkeleton.tsx`
  Skeleton: image rectangle + 3 text lines + button bar (shadcn Skeleton)
  Time: 20min

### 04.B â€” AJAX Cart

- [x] **🔴 Update `lib/store/cart-store.ts`**
  Path: `frontend/apps/restaurant-website/lib/store/cart-store.ts`
  Add:
  - AJAX `addItem()` â€” instant Zustand update, NO page reload
  - Cross-tab sync via `StorageEvent` listener (BroadcastChannel)
  - Namespace per restaurant: key = `pos-cart-{restaurantSlug}`
  - `totalPrice` and `totalItems` computed getters
  - Inventory limit check â€” disable add if `qty >= quantity_in_stock`
  - Warning toast when at inventory limit
  Test: 2 tabs â†’ add in tab 1 â†’ tab 2 updates instantly
  Time: 1.5h

- [x] **Create `lib/hooks/use-cart.ts`**
  Path: `frontend/apps/restaurant-website/lib/hooks/use-cart.ts`
  Convenience hook: `{ items, totalItems, totalPrice, addItem, removeItem, updateQty, clearCart }`
  Wraps cart-store with AJAX behavior and inventory checks
  Time: 30min

- [x] **🔴 Create `components/commerce/CartSheet.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/CartSheet.tsx`
  shadcn `<Sheet>` sliding from right (RTL: left). NO full page reload.
  Content: header with count + close | ScrollArea cart items | separator | order notes textarea | order summary (items + delivery + total) | gradient Checkout button | "Continue Shopping" link
  Empty state: big cart icon + "Ø³Ù„ØªÙƒ ÙØ§Ø±ØºØ©" + "Browse Menu" button
  Time: 2h

- [x] **Create `components/commerce/CartItem.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/CartItem.tsx`
  Layout: `[img 64px] Name AR | Name EN | Qty [-][2][+] | Price | [ðŸ—‘ï¸]`
  Optional note text. Swipe-to-delete on mobile (future).
  Time: 1h

### 04.C â€” Order Components

- [x] **Create `components/commerce/OrderStatusBadge.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/OrderStatusBadge.tsx`
  Status â†’ color: pendingâ†’yellow, confirmedâ†’blue, preparingâ†’blue, out_for_deliveryâ†’orange, deliveredâ†’green, cancelledâ†’red
  Arabic labels: Ù‚ÙŠØ¯ Ø§Ù„Ø§Ù†ØªØ¸Ø§Ø±, ØªÙ… Ø§Ù„ØªØ£ÙƒÙŠØ¯, Ø¬Ø§Ø±ÙŠ Ø§Ù„ØªØ­Ø¶ÙŠØ±, ÙÙŠ Ø§Ù„Ø·Ø±ÙŠÙ‚, ØªÙ… Ø§Ù„ØªÙˆØµÙŠÙ„, Ù…Ù„ØºÙŠ
  Uses: shadcn `<Badge>` with dynamic variant
  Time: 30min

- [x] **Create `components/commerce/OrderTimeline.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/OrderTimeline.tsx`
  Vertical timeline: âœ… completed | â—‰ current (pulse ring in theme-primary) | â—‹ future (gray)
  Each step: Circle icon + Label AR/EN + Timestamp
  RTL: dots on RIGHT side. Dark mode: auto via CSS vars.
  Time: 1.5h

- [x] **Create `components/commerce/OrdersTable.tsx`**
  Path: `frontend/apps/restaurant-website/components/commerce/OrdersTable.tsx`
  shadcn `<Table>` with columns: #, Date, Items, Status badge, Total, Actions
  Responsive: cards layout on mobile (hide table). Row click â†’ order detail.
  Actions: Track (active) / Reorder (delivered) / New Order (cancelled)
  Time: 1.5h

### 04.D â€” MVP Deploy Checklist (Day 15)

- [ ] **ðŸ”´ Run Lighthouse mobile audit**
  Target: â‰¥ 70 mobile
  Time: 15min

- [x] **ðŸ”´ Test all 10 themes (light + dark)**
  Visit homepage with each theme â†’ check header variant, section order, footer variant
  Both light + dark mode for each = 20 combinations
  Time: 1.5h

- [x] **🔴 Test RTL (Arabic) on all pages**
  Switch to Arabic â†’ verify header flips, text right-aligned, sheet from correct side
  Time: 30min

- [x] **ðŸ”´ Test AJAX cart flow**
  Menu â†’ add product â†’ CartSheet opens â†’ qty change â†’ checkout button works
  Time: 15min

- [ ] **Deploy MVP**
  Push to production branch. Smoke test live.
  Time: 30min

---

## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
## PHASE 05 â€” ALL PAGE REDESIGNS (Post-MVP)
## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
> **Goal:** Full redesign of all 7 internal pages with shadcn/ui + 3 payment methods
> **Estimated Time:** 5â€“6 days

### 05.A â€” Menu Page

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

### 05.B â€” Cart Page

- [ ] **Rewrite `app/[locale]/cart/cart-page-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/cart/cart-page-content.tsx`
  Desktop: 2-column (items left, sticky summary right). Mobile: single column, summary accordion.
  Cart items with qty controls. Order notes textarea. Summary card (items + delivery + total).
  Checkout button. "Continue Shopping" link. Empty state.
  Time: 1.5h

### 05.C â€” Checkout Page (3 Payment Methods)

- [ ] **Create `components/forms/CheckoutForm.tsx`**
  Path: `frontend/apps/restaurant-website/components/forms/CheckoutForm.tsx`
  RHF + Zod validated form. Fields: name, phone (Egyptian regex), address, notes, payment method.
  3 payment methods (scalable registry): Cash | Paymob Online | Screenshot Transfer
  Zod: phone `/^(01[0-2,5][0-9]{8})$/`, address min 10 chars
  Cash flow: POST /orders â†’ clear cart â†’ redirect to order detail + confetti
  Paymob flow: POST /orders â†’ iframe dialog â†’ poll status â†’ redirect
  Screenshot flow: POST /orders â†’ upload dialog â†’ redirect with "Pending Verification" badge
  Paymob fallback: if API down â†’ auto-switch to cash + toast
  Time: 3.5h

- [ ] **Rewrite `app/[locale]/checkout/checkout-page-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/checkout/checkout-page-content.tsx`
  Progress bar (4 steps: Cart â†’ Info â†’ Payment â†’ Done)
  2-column desktop (form left, sticky summary right). Mobile: accordion summary.
  Guard: redirect to /cart if cart empty.
  Time: 1.5h

### 05.D â€” Orders History Page

- [ ] **Rewrite `app/[locale]/orders/orders-page-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/orders/orders-page-content.tsx`
  Filter tabs: All / Active / Delivered / Cancelled + count badges
  Search by order number. OrdersTable (desktop) / Card list (mobile).
  Pagination: shadcn Pagination. Empty state per tab.
  Time: 2h

### 05.E â€” Order Detail Page

- [ ] **Rewrite `app/[locale]/orders/[id]/order-detail-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/orders/[id]/order-detail-content.tsx`
  Header: back button + order # + status badge
  OrderTimeline (vertical, animated). Items breakdown table. Customer info card.
  Actions: Track (active) / Reorder + Rate (delivered) / New Order (cancelled)
  Screenshot pending: "Upload Receipt" button
  Polling: `refetchInterval: 15000` for active orders
  Time: 2h

### 05.F â€” Track Order Page

- [ ] **Create `components/forms/TrackOrderForm.tsx`**
  Path: `frontend/apps/restaurant-website/components/forms/TrackOrderForm.tsx`
  Centered card: order number input + "Track" button. Not found â†’ error + "Order Now" link.
  Time: 45min

- [ ] **Rewrite `app/[locale]/orders/track/track-page-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/orders/track/track-page-content.tsx`
  Centered search card. Found â†’ redirect to `/orders/{id}`.
  Time: 30min

### 05.G â€” Settings Page

- [ ] **Rewrite `app/[locale]/settings/settings-page-content.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/settings/settings-page-content.tsx`
  4 tabs: Language (radio cards with flags) | Appearance (Light/Dark/System radio cards) | Cart (item count, clear, saved address) | Notifications (coming soon badge)
  Time: 2h

### 05.H â€” Table Menu Page

- [ ] **Enhance `app/[locale]/table/[tableId]/menu/page.tsx`**
  Path: `frontend/apps/restaurant-website/app/[locale]/table/[tableId]/menu/page.tsx`
  Table header banner: `ðŸ½ï¸ Ø·Ø§ÙˆÙ„Ø© Ø±Ù‚Ù… Ù¥ Â· Table #5` | Change Table | End Session
  Same menu layout. Pass `tableId` to cart store.
  Time: 1h

---

## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
## PHASE 06 â€” CROSS-CUTTING FEATURES (Post-MVP)
## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
> **Goal:** Search, cart sync, WhatsApp forwarding, dark mode audit
> **Estimated Time:** 3â€“4 days

### 06.A â€” Search Feature

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
  Read `?search=` and `?category=` on mount. On change â†’ `router.replace()`.
  Time: 30min

### 06.B â€” Cart Cross-Tab Sync

- [ ] **Add BroadcastChannel sync to cart-store**
  Path: `frontend/apps/restaurant-website/lib/store/cart-store.ts`
  `BroadcastChannel` for instant cross-tab sync (fallback: `StorageEvent`)
  Time: 45min

- [ ] **ðŸŸ  Cross-browser test: cart cross-tab sync on Chrome + Firefox + Safari**
  Test matrix: Open 2 tabs â†’ add item in tab 1 â†’ verify tab 2 updates
  Chrome: `BroadcastChannel` (native support)
  Firefox: `BroadcastChannel` (native support)
  Safari: Verify `StorageEvent` fallback works (Safari BroadcastChannel support varies)
  Edge cases: tab opened before cart init, rapid add/remove, offline tab reconnect
  Time: 30min

### 06.C â€” Dark Mode Complete Audit

- [ ] **Audit ALL components for hardcoded colors**
  `grep -r "background.*#[0-9a-f]" components/` â†’ fix each to CSS var
  `grep -r "color.*#[0-9a-f]" components/` â†’ fix each
  Every color must be: `bg-theme-*` or `text-theme-*` or `hsl(var(--theme-*))`
  Time: 1h

- [ ] **Screenshot test: dark mode on ALL pages**
  All 8 pages Ã— light + dark. All: Header, Footer, Cards, Forms, Modals.
  Time: 1h

### 06.D â€” Language Switcher Complete

- [ ] **Test language switch on all pages**
  ENâ†’AR: URL changes, RTL applies, Arabic font loads, all text translated
  Cart preserved across language switch.
  Time: 1h

- [ ] **Proofread Arabic translations**
  Compare ar.json vs en.json â€” fill any gaps.
  Time: 1h

### 06.E â€” WhatsApp Integration Complete

- [ ] **Add WhatsApp order forwarding on checkout success**
  Path: `frontend/apps/restaurant-website/app/[locale]/checkout/checkout-page-content.tsx`
  On order success: auto-open `wa.me/{adminNumber}?text={formatted order details}`
  Format: Order #, Customer name, phone, address, items Ã— qty = price, total, payment method
  Time: 1h

- [ ] **Add WhatsApp button to ProductCard**
  Path: `frontend/apps/restaurant-website/components/commerce/ProductCard.tsx`
  `wa.me/{number}?text=Ø£Ø±ÙŠØ¯ Ø£Ø·Ù„Ø¨ ${product.name_ar} Ã— ${qty}`
  Only shown if WhatsApp configured in theme
  Time: 30min

- [ ] **Add WhatsApp button to CTASection**
  Path: `frontend/apps/restaurant-website/components/sections/CTASection.tsx`
  Green WhatsApp button alongside "Order Now"
  Time: 15min

---

## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
## PHASE 07 â€” PERFORMANCE + TESTING (Post-MVP)
## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
> **Goal:** Lighthouse 95+ mobile, backend caching, E2E tests, PWA
> **Estimated Time:** 3â€“4 days

### 07.A â€” Backend Caching

- [ ] **ðŸ”´ Add in-memory theme cache (sync.Map)**
  Path: `backend/internal/handler/http/admin_theme_handler_v2.go`
  `sync.Map` with 5min TTL. Check cache â†’ return cached â†’ or query DB + store.
  Set `X-Cache: HIT/MISS` header. Impact: DB queries drop 95%.
  Time: 1.5h

- [ ] **ðŸ”´ Add in-memory product cache**
  Path: `backend/internal/handler/http/public_menu_handler.go`
  Same pattern, 10min TTL (products change less often).
  Time: 1h

- [ ] **ðŸ”´ Add cache invalidation on write**
  Path: `backend/internal/handler/http/admin_theme_handler_v2.go`
  When admin saves theme: `themeCache.Delete("theme:" + slug)`
  When admin updates products: `productCache.Delete("products:" + restaurantId)`
  Time: 45min

### 07.B â€” Frontend Performance

- [ ] **ISR for all public pages**
  Path: `app/[locale]/page.tsx`, `app/[locale]/menu/page.tsx`
  `export const revalidate = 300` (5 min). Orders/track: `revalidate = 0` (always fresh).
  Time: 30min

- [ ] **Parallel data fetching**
  Path: `app/[locale]/page.tsx`, `app/[locale]/menu/page.tsx`
  `Promise.all([fetchTheme, fetchProducts, fetchCategories])` â€” no waterfall.
  Time: 30min

- [ ] **Optimize images with next/image**
  All product images: `fill`, `sizes` prop, `quality={75}`, `placeholder="blur"`, `loading="lazy"` (below fold), `priority` (first 4).
  Time: 1h

- [ ] **Conditional Arabic font loading**
  Path: `frontend/apps/restaurant-website/app/[locale]/layout.tsx`
  Load Cairo ONLY when `locale === 'ar'`. English pages save 200KB.
  Time: 30min

- [ ] **Code splitting â€” dynamic imports**
  Dynamic: GallerySection, TestimonialsSection, CartSheet, PaymobDialog, Framer Motion
  Time: 45min

- [ ] **Bundle analyzer audit**
  `ANALYZE=true pnpm build`. Fix: tree-shake lucide imports, remove unused shadcn components.
  Target: no chunk > 80KB gzip.
  Time: 1h

- [ ] **Prefetch critical routes**
  Homepage â†’ prefetch `/menu`. Menu â†’ prefetch `/cart`. Cart â†’ prefetch `/checkout`.
  Time: 15min

- [ ] **Loading skeletons for all pages**
  Path: `app/[locale]/menu/loading.tsx`, `app/[locale]/orders/loading.tsx`, `app/[locale]/cart/loading.tsx`
  Time: 45min

- [ ] **Preconnect to API domain**
  Path: `app/layout.tsx`
  `<link rel="preconnect" href="...">` + `<link rel="dns-prefetch" href="...">`
  Time: 10min

### 07.C â€” Error States

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
  Gray overlay on ProductCard + "ØºÙŠØ± Ù…ØªØ§Ø­" badge + disabled buttons
  Time: 20min

### 07.D â€” Testing

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
  Tests: 5 different API shapes â†’ valid ThemeConfig
  Time: 45min

- [ ] **E2E: Full checkout flow (Playwright)**
  Path: `frontend/apps/restaurant-website/e2e/checkout.spec.ts`
  Flow: home â†’ menu â†’ add product â†’ CartSheet â†’ checkout â†’ fill form â†’ cash â†’ order detail
  Time: 2h

- [ ] **E2E: Language switch**
  Path: `frontend/apps/restaurant-website/e2e/language.spec.ts`
  Flow: /en/menu â†’ toggle AR â†’ verify /ar/menu + RTL + Arabic text
  Time: 1h

- [ ] **E2E: Dark mode**
  Path: `frontend/apps/restaurant-website/e2e/dark-mode.spec.ts`
  Flow: light â†’ toggle dark â†’ verify html.dark + header colors â†’ reload â†’ persists
  Time: 45min

### 07.E â€” PWA

- [ ] **Create `public/manifest.json`**
  Path: `frontend/apps/restaurant-website/public/manifest.json`
  name, short_name, icons (192+512), start_url, display:standalone, theme_color
  Time: 30min

- [ ] **Add PWA meta tags to layout**
  Path: `frontend/apps/restaurant-website/app/layout.tsx`
  `<link rel="manifest">`, `<meta name="theme-color">`, apple-touch-icon
  Time: 20min

---

## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
## PHASE 08 â€” SEO + MONITORING + PRODUCTION (Post-MVP)
## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
> **Goal:** Egyptian Google rankings, error tracking, analytics, production hardening
> **Estimated Time:** 2â€“3 days

### 08.A â€” SEO

- [ ] **Add dynamic metadata to all pages**
  Path: each `page.tsx` in `app/[locale]/`
  Homepage: restaurant name + "Ø§Ø·Ù„Ø¨ Ø§ÙˆÙ†Ù„Ø§ÙŠÙ†" / "Order Online"
  Menu: "Ø§Ù„Ù‚Ø§Ø¦Ù…Ø©" / "Menu" + restaurant name
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

### 08.B â€” Monitoring

- [ ] **Add Sentry error tracking**
  Path: `frontend/apps/restaurant-website/sentry.client.config.ts` + `sentry.server.config.ts`
  DSN from `NEXT_PUBLIC_SENTRY_DSN` env var. Free tier: 5K errors/month.
  Time: 1h

- [ ] **Add Google Analytics 4 (optional, per-tenant)**
  Path: `frontend/apps/restaurant-website/app/layout.tsx`
  Only load if `NEXT_PUBLIC_GA_ID` is set. Track: page views, add-to-cart, checkout, orders.
  Time: 30min

### 08.C â€” Final Production Checklist

- [ ] **Final Lighthouse audit (AFTER)**
  Compare with BEFORE screenshots. Target: Mobile â‰¥ 95, Desktop â‰¥ 98.
  Time: 45min

- [ ] **Test all 10 themes Ã— light + dark = 20 combinations**
  Visit homepage â†’ check header variant, section order, footer variant, colors.
  Time: 1.5h

- [ ] **Test on Egypt 4G throttling**
  Chrome DevTools â†’ Fast 4G. Goal: LCP < 2.5s, FID < 100ms, CLS < 0.1.
  Time: 30min

- [ ] **Verify theme builder preview mode**
  Dashboard editor â†’ iframe loads â†’ change color â†’ instant update â†’ dark toggle
  Time: 30min

---

## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
## DASHBOARD TASKS (Parallel with Phases 01â€“02)
## â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•
> These tasks update the Dashboard app to support the new theme features

### Dashboard: Dark Mode Editor

- [ ] **Create `DarkColorsEditor.tsx`**
  Path: `frontend/apps/dashboard/src/components/editors/DarkColorsEditor.tsx`
  Side-by-side light vs dark color pickers. "Auto-Generate Dark Palette" button.
  Live preview toggle. ContrastChecker integration. Saves to themes_v2 dark columns.
  Time: 2.5h

- [ ] **Update `themeBuilderStore.ts` â€” dark color state**
  Path: `frontend/apps/dashboard/src/stores/themeBuilderStore.ts`
  Add: `colorsDark`, `setDarkColor()`, `autoGenerateDarkPalette()`
  Time: 45min

- [ ] **Update `themeApi.ts` â€” dark colors + variants**
  Path: `frontend/apps/dashboard/src/lib/api/themeApi.ts`
  Add: dark color fields + template/variant fields in normalize/denormalize
  Time: 30min

- [ ] **Add dark mode preview toggle to theme editor**
  Path: `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx`
  [â˜€ï¸ Light | ðŸŒ™ Dark] toggle in preview toolbar
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
| Homepage load (Egypt 4G) | 3â€“4s | 1.2s |
| Menu page load | 2â€“3s | 0.8s |
| Navigation (page to page) | 500ms+ | < 100ms |
| Product image per item | ~500KB | ~50KB |
| JS bundle (gzip) | ~450KB | < 180KB |
| Arabic page extra font | +200KB always | +0KB (conditional) |
| LCP (mobile) | 4.5s | < 2.0s |
| Lighthouse mobile | ~60 | â‰¥ 95 |
| API calls per user session | 15â€“20 | 3â€“5 first visit, 0 repeat |
| Theme visual variety | 10% (color swap) | 85% (templates + variants) |

---

*Last Updated: 2026-03-03 | Status: Planning Phase | MVP: 15 working days*








