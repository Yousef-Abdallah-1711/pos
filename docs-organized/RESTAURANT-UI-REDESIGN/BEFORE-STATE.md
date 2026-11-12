# Restaurant Website UI Redesign - BEFORE State

Date: 2026-03-04
Branch: `restaurant-ui-redesign`
Git SHA: `899dce68f81e4a1f8efc3b77e5a6957b3e96c3ba`

## Lighthouse Baseline

Source reports:
- `docs-organized/RESTAURANT-UI-REDESIGN/lighthouse-BEFORE-mobile.html`
- `docs-organized/RESTAURANT-UI-REDESIGN/lighthouse-BEFORE-desktop.html`

Scores:
- Mobile: Performance `45`, Accessibility `100`, Best Practices `100`, SEO `54`
- Desktop: Performance `45`, Accessibility `100`, Best Practices `100`, SEO `54`

## Inventory Snapshot

- Route page files (`app/**/page.tsx`): `12`
- Component files (`components/**/*.ts(x)`): `19`
- Existing Next chunks footprint (`.next/static/chunks`): `12.85 MB` (local dev build artifact)

## Notes

- Baseline was captured before Phase 00 refactor tasks.
- A dedicated backup branch was created: `restaurant-ui-redesign-backup`.
- Phase 00 sandbox checks were executed for RTL shadcn primitives and Embla RTL configuration, then sandbox files were removed.
