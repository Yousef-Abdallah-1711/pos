# ORDERLY social media launch kit

Production exports, editable SVG masters, brand rules, and Egyptian Arabic copy for the user-approved **ORDERLY — أوردرلي** collateral identity. The live product materials still use POOS; this package does not change the application or authorize posting to social accounts.

## Start here

- [Brand identity and rationale](brand-guidelines/BRAND_IDENTITY.md)
- [Colors](brand-guidelines/BRAND_COLORS.md)
- [Typography](brand-guidelines/BRAND_TYPOGRAPHY.md)
- [Platform copy](brand-guidelines/SOCIAL_MEDIA_COPY.md)
- [Brand overview preview](previews/brand-overview.png)
- [Social mockups preview](previews/social-media-mockups.png)

## Export inventory

PNG files are ready to upload. Every PNG has a sibling SVG master with the same base name; SVGs are editable vector compositions with embedded fonts and, for social artwork, embedded backgrounds. Logo and wordmark PNGs preserve transparency; profile avatars and platform graphics use their illustrated backgrounds. Rerun `python source/generate_assets.py` after editing master layout/content, then `python source/render_exports.py` to regenerate dimension-checked PNGs.

| Platform / use | PNG export | Dimensions |
|---|---|---:|
| Facebook cover | `facebook/cover.png` | 1640 × 624 |
| Facebook profile | `facebook/profile.png` | 1080 × 1080 |
| LinkedIn company cover | `linkedin/cover.png` | 1128 × 191 |
| LinkedIn profile | `linkedin/profile.png` | 400 × 400 |
| Instagram profile | `instagram/profile.png` | 1080 × 1080 |
| Instagram launch post | `instagram/launch-post.png` | 1080 × 1350 |
| Instagram story | `instagram/story.png` | 1080 × 1920 |
| Instagram highlight template | `instagram/highlight-template.png` | 1080 × 1080 |
| YouTube profile | `youtube/profile.png` | 800 × 800 |
| YouTube channel banner | `youtube/banner.png` | 2560 × 1440 |
| YouTube thumbnail template | `youtube/thumbnail-template.png` | 1280 × 720 |
| Brand overview | `previews/brand-overview.png` | 1200 × 720 |
| Social mockups | `previews/social-media-mockups.png` | 1600 × 1000 |

## Logo exports

Each listed logo has SVG and PNG versions. The Arabic and English wordmarks are included separately as well as the combined primary marks.

| Variant | SVG / PNG stem | Dimensions |
|---|---|---:|
| Primary horizontal English | `logo/primary/orderly-primary` | 820 × 180 |
| Primary horizontal Arabic | `logo/primary/orderly-primary-arabic` | 840 × 190 |
| English wordmark | `logo/english/orderly-wordmark-en` | 600 × 140 |
| Arabic wordmark | `logo/arabic/orderly-wordmark-ar` | 540 × 160 |
| Symbol only | `logo/symbol/orderly-symbol` | 160 × 160 |
| Light-surface logo | `logo/dark/orderly-logo-dark` | 820 × 180 |
| Dark-surface logo | `logo/light/orderly-logo-light` | 820 × 180 |
| Monochrome black | `logo/monochrome/orderly-logo-black` | 820 × 180 |
| Monochrome white | `logo/monochrome/orderly-logo-white` | 820 × 180 |
| Favicon | `logo/favicon/orderly-favicon` | 64 × 64 |
| App icon | `logo/favicon/orderly-app-icon` | 512 × 512 |

## Sources and production

Product positioning and the palette/type direction were taken from `frontend/apps/landing/PRODUCT.md`, the landing page, and its design tokens. Background illustrations were generated for this package; no product screenshots or customer claims were invented. Font subsets and their OFL licenses are in `fonts/`. The renderer checks each exported PNG against the SVG's declared dimensions.
