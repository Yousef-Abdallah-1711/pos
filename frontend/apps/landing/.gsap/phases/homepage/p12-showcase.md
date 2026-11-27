> Archived generator scaffold. For the implemented homepage, use `.gsap/current-state.md`, `DESIGN.md`, and `docs/landing-handoff.md`. This scaffold contains unverified guesses.

## Current P12 implementation

Status: implemented, browser QA in progress. `FeatureShowcase.tsx` pins `.cv-showcase-stage` after its separate introduction. The stage uses a measured header offset and `100svh`; horizontal travel is `track.scrollWidth - viewport.clientWidth`. Progress selects the nearest actual slide offset. At ≤900px, <760px viewport height, or reduced motion, the pin is absent and the viewport scrolls natively. Verification covers first/second/final/reverse positions, Arabic RTL, themes, and 104 Chromium viewport cases; final production checks remain.

# Phase P12 - Showcase

## Scope

- Page: homepage
- Section: Showcase
- Recipe Direction: parallax showcase or marquee
- Status: Planned

## Objective

Make the Showcase section feel intentional, modern, and aligned with the page hierarchy.

## Tasks

- [ ] Inspect the current code and selectors for the Showcase section before editing.
- [ ] Preserve any good existing Showcase motion and remove duplicate or noisy effects first.
- [ ] Align the Showcase motion with the inferred product type: Marketing Site, Docs / Knowledge Base, SaaS App.
- [ ] Document the motion goal, hierarchy, and fallback behavior for Showcase in the page artifact.
- [ ] Implement the Showcase animation with transform/opacity-first properties.
- [ ] Verify reduced-motion behavior for Showcase.
- [ ] Verify mobile downgrade behavior for Showcase.
- [ ] Update the Showcase status and notes in .gsap artifacts.

## Validation

- Reduced Motion:
- Mobile Downgrade:
- Notes:
