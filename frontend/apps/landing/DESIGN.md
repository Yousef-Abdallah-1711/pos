---
name: Servio Coral Marketing and Authentication
description: A warm editorial presentation of connected restaurant service for Egypt.
colors:
  paper: '#F8F2E8'
  paper-2: '#EFE4D5'
  paper-3: '#E6D6C3'
  ink: '#2B201C'
  ink-2: '#5A4941'
  muted: '#6E5A4E'
  rule: '#D8C8B6'
  accent: '#A13224'
  accent-hover: '#8D291C'
  dark-paper: '#211B18'
  dark-surface: '#2D2520'
  dark-ink: '#F2E9DC'
  dark-accent: '#E17861'
typography:
  english-display: 'Fraunces, Georgia, serif'
  english-body: 'Manrope, Arial, sans-serif'
  arabic: 'IBM Plex Sans Arabic, sans-serif'
---

# Servio Coral design system

## Direction

The approved landing combines the Coral color and type system with the Ultraviolet concept's section rhythm. The hero is centered over a subdued café photograph and brings a large real dashboard capture into the first screen. This direction supersedes the earlier violet ledger concept for the landing and authentication pages.

The page speaks to owners of restaurants and cafés in Egypt. It follows a concrete service day: counter and online orders, staff confirmation, kitchen preparation, fulfillment, and the owner's view of daily work. Each claim must match a working product path. Online ordering means a guest storefront with cash on delivery order submission and staff confirmation before kitchen preparation. Do not imply card payment, automatic kitchen acceptance, or unsupported metrics.

## Palette and modes

- Light mode uses warm paper (`#F8F2E8`), sand sections (`#EFE4D5`, `#E6D6C3`), dark brown text (`#2B201C`), and terracotta actions (`#A13224`).
- Dark mode uses espresso surfaces (`#211B18`, `#2D2520`), warm ivory text (`#F2E9DC`), and a lighter coral action (`#E17861`).
- The photo hero and footer use deep espresso, with a legible ivory foreground in both modes.
- Focus is a visible two pixel outline with offset. Contrast must remain readable in both modes.
- Source tokens live in `src/styles/marketing-coral.css`. The `cv-*` classes belong to the homepage; shared `sv-*` and `auth-*` selectors carry the system into navigation, footer, and account entry.

## Typography and layout

Fraunces gives English editorial headings a warm character; Manrope keeps navigation and body text clear. IBM Plex Sans Arabic handles Arabic display and text. Never apply Latin letter spacing to Arabic. English and Arabic have equal content priority and the page `lang` and `dir` follow the route locale.

The content container is at most 1200 pixels wide. Large headings and asymmetric section pairings create rhythm, with fine rules and restrained cards for steps and capabilities. On narrow screens, pairings and cards stack, calls to action span the available width, and product captures retain their natural proportion. Use logical spacing and alignment so RTL is first class.

## Product evidence and calls to action

Use the real capture matrix under `public/images/product/servio/matrix/`, selecting the right locale and color mode. Captions identify the captures as actual product screens with synthetic data. The café photo at `public/images/marketing/cafe-interior.jpg` is credited to Julen Rey Azcona. Calls to action go to the real localized registration route and preserve funnel tracking. Avoid fabricated logos, testimonials, counts, and commercial promises.

The order story, online storefront, kitchen workflow, capabilities, getting started steps, FAQ, and final account action are the homepage's core sections. Navigation and footer stay functional across locales and modes; registration and login forms keep their real validation and authentication behavior.

## Interactive homepage system

The current homepage has thirteen sections: photographic hero, service-day problems, seven-step order story, storefront demo, kitchen demo, layered product typography, seven-feature showcase, ecosystem map, management view, stacked setup cards, pricing bridge, FAQ, and final registration action. `ConversionLanding.tsx` owns this order; individual interactive components live beside it in `src/components/marketing/`.

`src/styles/marketing-interactive.css` extends the Coral tokens. It contains the desktop sticky story stage, horizontal feature panels, localized demo controls, the ecosystem map, management samples, and responsive fallbacks. Keep the showcase viewport's `contain: layout paint`: without it, full-page Chrome capture and horizontal scrolling expose the entire off-screen feature track. Desktop feature motion uses a scoped GSAP `ScrollTrigger` pin; mobile and reduced-motion desktop use native horizontal scrolling with snap. The story stage uses CSS sticky positioning and GSAP only to synchronize narrative state. The hero capture, layered text, ecosystem lines, setup cards, and sample metrics have small, scoped GSAP effects. Motion must never gate access to content or links.

The story combines three actual screen areas (`checkout`, `dashboard`, `kitchen`) with a visibly labeled synthetic order overlay. Its order states and separate kitchen ticket states follow `backend/internal/domain/order.go` and the dashboard KDS store. The storefront and kitchen controls demonstrate explanatory state locally; they do not submit orders or change the backend. Management sample figures are not performance claims. Feature links lead to the existing detail routes for product evidence and limitations.

The oversized SERVIO footer signature follows the real navigation grid. Keep the footer groups as actual grid children; `display: contents` caused fragile alignment. The shared header, account links, theme control, and locale switch remain in the existing shell. Preserve the localized registration route and `TrialLink` tracking when editing calls to action.
