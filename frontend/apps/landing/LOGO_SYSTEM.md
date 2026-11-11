# POOS Logo System

## Files

All vector masters live in `public/brand/`:

- `poos-logo.svg`: primary horizontal logo for light backgrounds.
- `poos-logo-dark.svg`: primary horizontal logo for dark backgrounds.
- `poos-logo-mono.svg`: single-color horizontal version.
- `poos-logo-vertical.svg`: stacked light-surface lockup.
- `poos-logo-vertical-dark.svg`: stacked dark-surface lockup.
- `poos-logo-vertical-mono.svg`: stacked monochrome lockup.
- `poos-wordmark.svg`: wordmark-only lockup.
- `poos-symbol.svg`: standalone violet symbol.
- `poos-symbol-dark.svg`: lavender/white symbol for dark backgrounds.
- `poos-symbol-mono.svg`: single-color symbol.
- `poos-symbol-white.svg`: white symbol for saturated brand backgrounds.
- `poos-app-icon.svg`: square app icon.
- `poos-social-profile.svg`: large square social profile artwork.

`src/app/icon.svg` is the browser favicon and shares the same mark geometry.

## Construction and spacing

The 48-unit symbol uses a continuous P-shaped route with a rounded endpoint node. Keep at least 8 units of clear space around the symbol. At small sizes, the route is a 7-unit stroke so its counters remain legible. Use the full horizontal logo at 120 CSS pixels wide or larger; below that, use the standalone symbol with an accessible brand name supplied by the surrounding interface.

## Minimum size

- Symbol: 16px minimum for favicon use; 24px minimum in interface navigation.
- Horizontal logo: 120px minimum width.
- Wordmark: 96px minimum width.
- App and social icons: use the supplied square artwork without cropping or adding a second border.

## PNG exports

Run `pnpm --filter @pos-saas/landing export:brand-png` to regenerate the PNG derivatives in `public/brand/png/`. The script rasterizes the SVG masters with the repository's Playwright Chromium dependency and preserves transparent backgrounds. Keep SVG masters as the editable source; raster exports are delivery derivatives.
