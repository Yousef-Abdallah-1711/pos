# Existing marketing capture integration

`frontend/apps/landing/src/content/captures.json` defines product capture dimensions; `ProductCapture.tsx` delegates to `CaptureMedia.tsx`, which selects files under `public/images/product/servio/matrix/` by locale and resolved theme for known capture names. Unknown/non-matrix captures use their fallback and disclosure. `screenshot-matrix.json` and `provenance.json` contain existing capture provenance. The current CSV covers 37 image assets under `public/images/product` plus the legacy capture identified below; it does not establish that every product screenshot elsewhere under `public` has been found or prove active usage for each asset. Audit feature routes and all product-image references before considering the inventory complete.

Hero, bento, order story, workflow demos, feature showcase, ecosystem and management components consume these captures in different frames and may crop/scale them. Replacing files requires validating intrinsic dimensions, locale/theme source selection, disclosure, full-image links and animation geometry. Existing screenshots represent the old product appearance; they are not proof of a Coral-rebranded dashboard.

| Source selection | Verified usage path |
| --- | --- |
| `dashboard` | `src/components/marketing/ConversionLanding.tsx` hero, `ManagementExperience.tsx`, and staff-confirmation/delivery stages in `src/content/orderStoryCopy.ts`. |
| `pos` | `src/components/marketing/LayeredTypography.tsx`; also inspect `OperationalBento.tsx` and the feature definitions. |
| `checkout` | Online `WorkflowDemo.tsx` and browse/submit stages in `orderStoryCopy.ts`. |
| `kitchen` | Kitchen `WorkflowDemo.tsx` and kitchen stages in `orderStoryCopy.ts`. |
| `inventory`, `staff` | Selectable views in `ManagementExperience.tsx`; additional feature routes use their mapped capture definitions. |
| Feature-defined names | `src/content/features.ts` reads localized feature definitions; `FeatureShowcase.tsx` and `ProductExplorer.tsx` pass the selected name to `ProductCapture`. Verify all locale JSON mappings when changing an asset. |

`ProductCapture.tsx` uses `CaptureMedia.tsx` for selection and full-capture disclosure. `captures.json` lists 1440×900 for dashboard/POS/kitchen/inventory/theme/staff and 390×2124 for checkout; those are registry dimensions, not a guarantee every rendered container uses that size. The CSV reports intrinsic dimensions per file. Inspect `marketing-coral.css`, `marketing-interactive.css` and each consuming component before changing aspect ratios or replacing a screenshot.

The two `roles-*.png` assets under `public/images/product` are referenced by `src/components/Hero.tsx`; the inspected source search found that component imported by its unit test and did not find an active page import. Treat active usage as unverified until route-level analysis confirms it. `public/images/poos-access-control-capture.png` exists outside the scanned product directory and had no source reference in the search; it remains historical/unverified, not approved for deletion. The café photograph is an intentionally non-product marketing asset.

As a broader filename check on this checkout, `rg --files frontend/apps/landing/public -g '*.png' -g '*.webp' -g '*.jpg' -g '*.jpeg'` returned 56 paths: the 38 inventoried product-related images, one café photo, and 17 raster logo/icon files under `public/brand/png`. That establishes the raster filename coverage for those extensions in this local checkout; it does not prove actual usage, provenance or that SVGs/remote assets cannot depict product interfaces.
