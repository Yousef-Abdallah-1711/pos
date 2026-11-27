# GSAP homepage implementation state

Updated 2026-11-27. This file supersedes the auto-generated `.gsap` scaffold, which guessed unrelated sections such as testimonials, team, and a form. The approved architecture and claim boundaries are in `DESIGN.md`, `docs/landing-handoff.md`, and the user's attached 43-section brief. Package manager: **pnpm 10.34.9**, workspace root `frontend/`.

## Runtime and ownership

- `gsap` and `@gsap/react` are direct dependencies of `@pos-saas/landing`; `ScrollTrigger` is imported from `gsap/ScrollTrigger` and registered in each client motion island.
- Each island uses `useGSAP` with an element scope, `gsap.matchMedia()` for desktop motion, and cleanup through `media.revert()` and trigger/tween cleanup. Do not animate shared marketing routes by selecting document-wide classes.
- `HeroCaptureMotion.tsx`: dashboard capture perspective settles on desktop scroll. Mobile and reduced motion retain the static capture.
- `OrderStory.tsx`: seven-stage narrative and sample status update. The stage sticks via CSS; GSAP updates stage selection on desktop scroll. Avoid `ScrollTrigger.pin` inside its grid; it previously expanded the grid and page height.
- `WorkflowDemo.tsx`: local storefront and kitchen step selection with restrained status animation. It is an illustrative UI, not a backend simulator.
- `LayeredTypography.tsx`: restrained background text/capture parallax on desktop.
- `OperationalBento.tsx`: three asymmetric product cards with actual locale/theme POS, overview, and kitchen captures. Three-step local explanatory timelines run once as cards enter on desktop; keyboard and touch buttons select steps at every size. Reduced motion and mobile keep the same static content and manual controls. This section has no pin and no backend mutation.
- `FeatureShowcase.tsx`: seven real capture panels. The editorial introduction is outside the pin, while `.cv-showcase-stage` fills the space below the measured header. Desktop ScrollTrigger pins only that stage; its travel uses track overflow and progress uses actual slide offsets. Widths ≤900px, heights <760px, and reduced motion use native horizontal scrolling. `.cv-showcase-viewport { contain: layout paint }` prevents the off-screen track from expanding document width in Chrome.
- `BusinessEcosystem.tsx`: SVG line draw and selectable module detail.
- `ManagementExperience.tsx`: sample number reveal and selectable real product captures.
- `StackedSetupCards.tsx`: desktop scroll scale; static stacked cards on mobile and reduced motion.

## Verification

See `docs/landing-interactive-experience/artifacts/acceptance/results.json` for the 12-case Chrome locale/theme/viewport matrix and Firefox/WebKit checks, and `route-smoke.json` for shared routes and extra widths. Verify results after any layout or motion changes. `prefers-reduced-motion: reduce` was manually inspected: no pin spacer, native horizontal feature scrolling, identity track transform.

## Safe continuation

Read `CODING_STANDARDS.md`, app `AGENTS.md`, and installed Next guides. Preserve existing worktree changes. Keep actual product screen captions and synthetic sample disclosures. Test Arabic RTL, both themes, keyboard controls, the footer, and a production build before claiming completion. This document is the authoritative phase record; generated phase notes under `.gsap/phases/homepage/` are archived scaffolding and should not guide implementation.
