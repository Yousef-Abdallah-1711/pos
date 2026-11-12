# POOS Brand Guidelines

POOS is an all-in-one restaurant management system for independent restaurants and multi-location groups. The public-facing brand promise is one platform for front- and back-of-house work. Product proof must come from verified POOS screens.

## Identity

The identity pairs a compact, continuous P-shaped route with a terminal node. The P anchors the wordmark; the route extension and node suggest connected restaurant work without using kitchen utensils or borrowed brand symbols. Use the supplied vector files as the source of truth.

## Color

| Token | Value | Use |
| --- | --- | --- |
| Primary violet | `#6D28D9` | Primary actions, symbol, active links |
| Secondary purple | `#8B5CF6` | Small route endpoint accent |
| Lavender | `#C4B5FD` | Dark-surface symbol and supporting accent |
| Light background | `#F7F5FF` | Main light canvas |
| Dark ink | `#171327` | Wordmark and primary text |
| Muted ink | `#575066` | Supporting text on light backgrounds |
| Deep aubergine | `#211833` | Dark editorial surfaces |
| Divider | `#D9D1EB` | Borders and ledger rules |

The existing landing CSS tokens are canonical for semantic states, borders, surfaces, and elevation. Keep success, warning, danger, and focus states distinct from brand accents.

## Typography

- Arabic headings: Alexandria.
- Arabic body and interface: IBM Plex Sans Arabic.
- English body and interface: Manrope.
- Logo wordmark: Manrope ExtraBold, with tight tracking. Preserve the SVG wordmark artwork instead of retyping it in another font.

## Use

- Prefer `poos-logo.svg` on light surfaces and `poos-logo-dark.svg` on dark surfaces.
- Use the symbol at 16, 32, or 48 CSS pixels only when the context already identifies POOS; use 24px or larger in product UI where possible.
- Use monochrome artwork for single-color printing and constrained placements.
- Keep a clear space of at least 8 units around the 48-unit symbol, or one quarter of the mark width around the full logo.
- Do not stretch, rotate, add effects, redraw the path, or place the light version on a low-contrast surface.
- The route node may not be separated from the P path; it is part of the symbol.

## Public product claims

Use “All-in-one restaurant management system” and “One platform for front- and back-of-house work.” Do not add customer counts, testimonials, pricing, integrations, or capability claims without verified product evidence and approval.
