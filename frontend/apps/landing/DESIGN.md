---
name: POOS Landing and Authentication
description: Arabic-first restaurant operations presented as a precise shared shift ledger.
colors:
  violet: "#6D28D9"
  violet-deep: "#4C1D95"
  purple: "#8B5CF6"
  lavender: "#C4B5FD"
  paper: "#F7F5FF"
  surface: "#F0ECF8"
  surface-raised: "#FFFFFF"
  ink: "#171327"
  ink-soft: "#373147"
  muted: "#575066"
  dark: "#211833"
  dark-soft: "#302447"
  rule: "#D9D1EB"
  success: "#17643B"
  warning: "#8A5100"
  danger: "#A3233D"
  focus: "#6D28D9"
typography:
  display:
    fontFamily: "Alexandria, IBM Plex Sans Arabic, sans-serif"
    fontSize: "clamp(2.7rem, 5.3vw, 5rem)"
    fontWeight: 700
    lineHeight: 1.12
    letterSpacing: "-0.035em"
  headline:
    fontFamily: "Alexandria, IBM Plex Sans Arabic, sans-serif"
    fontSize: "clamp(2.15rem, 4.3vw, 3.8rem)"
    fontWeight: 700
    lineHeight: 1.15
    letterSpacing: "-0.035em"
  title:
    fontFamily: "Alexandria, IBM Plex Sans Arabic, sans-serif"
    fontSize: "1.3rem"
    fontWeight: 700
    lineHeight: 1.35
  body:
    fontFamily: "IBM Plex Sans Arabic, Manrope, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.75
  label:
    fontFamily: "IBM Plex Sans Arabic, Manrope, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 700
    lineHeight: 1.4
rounded:
  xs: "0.35rem"
  sm: "0.5rem"
  md: "0.625rem"
  control: "0.75rem"
  panel: "1.25rem"
spacing:
  xs: "0.35rem"
  sm: "0.75rem"
  md: "1rem"
  lg: "1.5rem"
  xl: "2rem"
components:
  button-primary:
    backgroundColor: "{colors.violet}"
    textColor: "{colors.surface-raised}"
    typography: "{typography.label}"
    rounded: "{rounded.control}"
    padding: "0.8rem 1.2rem"
    height: "3.25rem"
  button-secondary:
    backgroundColor: "{colors.surface-raised}"
    textColor: "{colors.ink}"
    typography: "{typography.label}"
    rounded: "{rounded.control}"
    padding: "0.8rem 1.2rem"
    height: "3.25rem"
  input:
    backgroundColor: "{colors.surface-raised}"
    textColor: "{colors.ink}"
    typography: "{typography.body}"
    rounded: "{rounded.md}"
    padding: "0.75rem 0.875rem"
    height: "3rem"
  ledger-proof:
    backgroundColor: "{colors.surface-raised}"
    textColor: "{colors.ink-soft}"
    rounded: "{rounded.panel}"
---

# Design System: POOS Landing and Authentication

## Overview

**Creative North Star: “The Shared Shift Ledger”**

POOS turns the restaurant close-out record into a digital visual language: aligned entries, fine rules, and a continuous line joining the counter to the kitchen and close-out work. Violet marks action and continuity on soft violet paper, while dark ink keeps the experience grounded and businesslike. The result is a premium product presentation with clear working structure rather than a generic SaaS hero followed by repeated feature cards.

The landing page moves from a genuine POOS product capture and account-creation action into connected front- and back-of-house entries, onboarding steps, and practical answers. Authentication keeps the same visual notation while giving forms a calm, readable paper surface. The public experience is Arabic-first; locale controls both language and reading direction, and English remains a fully supported presentation. Public product claims stay within verified behavior and approved evidence.

**Key Characteristics:**
- Ledger rules and aligned entries connect operational areas and steps.
- Violet and lavender accents sit on paper, white, and deep ink surfaces.
- Product evidence is a real POOS screen, not a fabricated dashboard.
- Arabic-first RTL and English LTR layouts use logical spacing and direction-aware alignment.

## Colors

The palette pairs POOS violet with pale lavender, soft paper surfaces, dark ink, and restrained semantic status colors.

### Primary
- **POOS Violet** (`{colors.violet}`): Primary actions, active links, ledger rules, and focus outlines.
- **Deep Violet** (`{colors.violet-deep}`): Button hover, emphasized links, and high-contrast violet text.
- **Bright Purple** (`{colors.purple}`): Secondary brand accent; use sparingly beside the primary violet.
- **Soft Lavender** (`{colors.lavender}`): Continuity rules, numbered markers, and accents on dark surfaces.

### Neutral
- **Violet Paper** (`{colors.paper}`): Default page background and the surface behind the landing ledger.
- **Soft Surface** (`{colors.surface}`): Subtle tonal separation and muted ledger regions.
- **Raised White** (`{colors.surface-raised}`): Forms, evidence frames, and cards that sit above the paper.
- **Dark Ink** (`{colors.ink}`): Main text and deep footer/background areas.
- **Soft Ink** (`{colors.ink-soft}`): Secondary interface text such as evidence headings.
- **Muted Ink** (`{colors.muted}`): Supporting body copy and explanatory text.
- **Deep Aubergine** (`{colors.dark}`): Dark story panels, audience ledger, and final calls to action.
- **Aubergine Surface** (`{colors.dark-soft}`): Secondary dark tone for related deep surfaces.
- **Ledger Rule** (`{colors.rule}`): Dividers, card borders, and fine structural lines.
- **White** (`{colors.surface-raised}`): Text and clear content surfaces on dark regions.

### Named Rules
**The One Ledger Rule.** Use violet lines and aligned entries to connect related restaurant work; do not replace the page's organizing idea with a repeated feature-card grid.

## Typography

**Display Font:** Alexandria (with IBM Plex Sans Arabic and sans-serif fallbacks)
**Body Font:** IBM Plex Sans Arabic for Arabic; Manrope for English (with the other interface family as fallback)
**Label/Mono Font:** IBM Plex Sans Arabic for Arabic UI labels; Manrope for English UI labels

**Character:** Alexandria gives Arabic headings a distinctive, confident voice. IBM Plex Sans Arabic keeps Arabic interface and body copy legible; Manrope is the English body and interface face. The source applies Alexandria to headings in both locales, with Latin glyphs available in its loaded subsets.

### Hierarchy
- **Display** (700, `clamp(2.7rem, 5.3vw, 5rem)`, 1.12): Landing promise; keep the line measure short and expressive.
- **Headline** (700, `clamp(2.15rem, 4.3vw, 3.8rem)`, 1.15): Major section and closing-call headlines.
- **Title** (700, `1.3rem`, 1.35): Ledger entries, form headings, and compact section titles.
- **Body** (400, `1rem`, 1.75): Explanations, form guidance, and long-form interface copy; preserve comfortable line lengths around 35–40ch where constrained.
- **Label** (700, `0.875rem`, 1.4): Controls, captions, navigation, and compact evidence labels.

### Named Rules
**The Locale Pairing Rule.** Arabic body and interface text use IBM Plex Sans Arabic with Alexandria headings; English body and interface text use Manrope. Keep the route locale reflected in `lang` and `dir`.

## Layout

The landing uses a wide ledger spread: a two-column hero places the promise and action opposite authentic product evidence, followed by rule-bound entries for product areas and a three-step horizontal line. The shared page container is capped at 78rem with 1.25rem side gutters, tightening to 0.75rem on narrow screens. Section spacing scales with viewport size, generally from 4.5rem to 8rem.

At narrower widths, the hero, paired ledger entries, and section headings stack. Onboarding steps change from a horizontal row to a vertical line with numbered markers. The navigation becomes a compact menu below 68rem; at phone widths, controls and account actions use full-width or single-column layouts as space requires. Prefer logical properties such as `padding-inline`, `inset-inline`, and `margin-block` so RTL and LTR share the same layout rules.

## Elevation & Depth

Depth uses both tonal surfaces and restrained shadows. The paper-to-white shift and fine borders do most of the structural work; shadows lift the product evidence, sticky navigation state, auth frame, and modal where separation matters.

### Shadow Vocabulary
- **Evidence and Auth Card** (`0 22px 58px rgb(35 23 61 / 12%)`): Lift the product capture and desktop authentication frame from the paper.
- **Navigation Lift** (`0 8px 24px rgb(35 23 61 / 7%)`): Subtle sticky-header separation when active.
- **Tenant Dialog** (`0 1.5rem 4rem rgb(23 19 39 / 24%)`): Clearly separate the native tenant-selection dialog from its dimmed backdrop.

### Named Rules
**The Paper First Rule.** Use surface contrast and fine rules for ordinary grouping; reserve stronger shadows for overlays and prominent product/authentication surfaces.

## Shapes

Forms are gently rounded and consistent rather than pill-shaped. Controls use 0.625–0.75rem corners; framed evidence uses a larger 1.25rem radius, reduced on mobile. Ledger entries and section boundaries rely on straight one-pixel rules, with circles reserved for step markers and continuity nodes. Use pill ends only for compact status/sample chips or the short central join mark.

## Components

### Buttons
- **Character:** Confident, clear actions with compact, comfortable touch targets.
- **Primary:** Violet fill, white text, 0.75rem corners, 0.8rem by 1.2rem padding, and at least 3.25rem height.
- **Hover / Active / Focus:** Deep-violet hover; a slight one-pixel press on active; global 3px violet focus outline with 3px offset.
- **Secondary:** White fill with ink text and a fine ledger-rule border; hover adds a pale violet surface and deep-violet text.
- **Disabled:** Reduce opacity and show a non-action cursor; preserve readable text.

### Cards / Containers
- **Character:** Quiet paper frames around genuine content.
- **Evidence Frame:** White surface, 1px muted lavender border, 1.25rem desktop radius, and the evidence/auth shadow. The caption and heading are separated by thin rules.
- **Authentication Frame:** White form surface beside a deep aubergine story panel; rounded 1.5rem on desktop and edge-to-edge on mobile.
- **Internal Padding:** Follow the spacing rhythm; evidence heading/caption use about 0.75–1rem and 1–1.5rem form or story insets.

### Inputs / Fields
- **Style:** Full-width white field, 1px muted stroke, 0.625rem corners, 3rem minimum height, and 0.75rem by 0.875rem padding.
- **Focus:** Visible 3px violet outline with offset; do not depend on color shift alone.
- **Error / Disabled:** Danger border and supporting message for invalid values; disabled fields reduce opacity and stop accepting input.

### Navigation
- **Style:** Sticky paper header with a fine bottom rule, 5rem desktop row, ink links, and a violet account-start action. Hover uses a soft lavender wash and deep violet text.
- **Mobile:** Below 68rem, replace inline links with an accessible menu; keep language switching and login/account actions reachable. On phone widths the header row contracts to 4.25rem.

### Shared Shift Ledger
Use horizontal rules, paired front/back entries, concise labels, and a visible continuity mark to show that restaurant work remains connected. Step numbers sit on the connecting line; at phone widths that line turns vertical. Keep the story legible in both writing directions.

### Tenant Selection Dialog
Use a native modal dialog with a white surface, strong separation shadow, clear close action, and full-width tenant choices. Selected and hovered choices receive violet border and pale violet fill; roles remain small outlined labels. Keep keyboard cancellation and focus visibility intact.

## Do's and Don'ts

### Do:
- **Do** keep the Shared Shift Ledger as the organizing visual idea, using aligned entries and fine violet rules to connect work.
- **Do** use real POOS interface captures with their documented source; never present mock data as a live account.
- **Do** preserve Arabic-first RTL and English LTR through locale-aware `lang`, `dir`, and logical CSS properties.
- **Do** use the English and Arabic font pairing and color tokens defined above.
- **Do** honor visible keyboard focus and reduced-motion preferences.
- **Do** keep public claims within verified capabilities and approved evidence.

### Don't:
- **Don't** turn the landing into a centered generic SaaS hero followed by repeated feature cards.
- **Don't** invent customer counts, testimonials, customer logos, prices, trial terms, or performance claims.
- **Don't** use fabricated dashboard data as product proof.
- **Don't** hard-code left/right layout assumptions where a logical property supports both locales.
- **Don't** use the unused Tailwind blue scale as POOS brand color guidance.
