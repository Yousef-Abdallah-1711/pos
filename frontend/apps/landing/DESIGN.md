# POOS landing and authentication design system

## Visual idea

The Shared Shift Ledger uses the restaurant close-out record as its organizing form. Thin violet rules connect work areas and onboarding steps; aligned entries replace repeated feature cards. The page keeps useful navigation, product evidence, and account creation visible without inventing commercial proof.

## Foundations

- Arabic-first public experience, with URL locale controlling page `lang` and `dir`.
- Alexandria for headings; IBM Plex Sans Arabic for interface and body copy; Inter as the Latin body fallback.
- Violet `#6D28D9` is the primary action and rule color; lavender `#C4B5FD` marks continuity; paper `#F7F5FF`, ink `#171327`, and white define surfaces.
- Muted body text `#5C566F`, success `#17643B`, warning `#8A5100`, danger `#A3233D`, and violet focus outlines supply semantic states.
- Logical CSS properties support both RTL and LTR. Mobile uses a single-column ledger and a vertical onboarding line.

## Components and interaction

- `BrandMark` supplies the inline wordmark and symbol; `icon.svg` provides the favicon.
- `poos-button`, `poos-input`, and `poos-skip-link` set shared controls and keyboard entry.
- Landing sections use ledger rows and native disclosure elements. Registration and login keep the existing API/session behavior. Multi-tenant selection uses a native modal dialog with Escape/cancel behavior.
- Reduced-motion preferences are respected; visible focus styling is global.

## Evidence and claims

Only the dashboard access-control capture is used as product evidence, with its source and limitations recorded beside the image. No testimonials, customer logos, metrics, approved prices, or trial terms are presented.
