> Archived generator scaffold. For the implemented homepage, use `.gsap/current-state.md`, `DESIGN.md`, and `docs/landing-handoff.md`. This scaffold contains unverified guesses.

## Current audit, 2026-11-27

- The old `.cv-problem-grid` showed three identical text cards and no product evidence. It has been removed; the asymmetric bento uses the existing four-way capture matrix and local illustrative states.
- In Chrome at 1440×900, the old showcase pinned the whole section at 78px. Its slide ran from 441px to 1083px, below the 900px viewport. The heading and 124px section padding consumed usable height. The new pin target is only `.cv-showcase-stage`; it measured 90px to 889px at the same viewport, with one spacer and a fully contained first slide.
- The old active counter rounded global ScrollTrigger progress. The new counter finds the nearest measured slide offset. The last slide and reverse travel are included in browser checks.
- Remaining release checks: clean production build, full route smoke, CI, and PR merge. The host `form-foundation/node_modules` links are damaged and block host typecheck/lint/auth routes; use the clean container build for final integration evidence.

# Animation Audit Report - landing

## Workflow State

- Last Audit Mode: gsap-refactor
- Resume Rule: read this file before repeating an audit

## Summary

- Date:
- Scope:
- Overall Health:

## Must Fix

- None recorded

## Should Fix

- None recorded

## Nice To Have

- None recorded

## Fixes Applied

- None recorded

## Refactor Snapshot - homepage

- Existing Motion Stack: None detected
- Matched Files: multi-tenant-login.integration.test.tsx, page.tsx, page.tsx, page.tsx, layout.tsx, page.tsx, page.tsx, page.tsx
- Planned Phases: P01 Navigation, P02 Feature Grid, P03 Stats, P04 Timeline, P05 Form, P06 Hero, P07 Testimonials, P08 CTA, P09 Pricing, P10 FAQ, P11 Team, P12 Showcase
- Findings:
  - GSAP does not appear to be implemented yet on the matched files.
  - Hero and repeated-card motion should not carry equal visual weight.
- Recommended Improvements:
  - Use a staged hero reveal with headline, supporting copy, and CTA arriving in deliberate sequence.
  - Use one card system for the grid and keep the stagger controlled so the section feels premium, not noisy.
  - Use count-up animation only when the stats enter view and only once.
  - Consider a subtle navbar polish on scroll rather than a dramatic transformation.
  - Use restrained pricing-card emphasis so the CTA hierarchy stays clear without making the section feel gimmicky.
  - Use state-based form polish like focus, validation, and success transitions instead of theatrical entrance effects.
  - Reserve marquee, parallax, or layered depth for a showcase section instead of distributing spectacle everywhere.
  - Animate charts and dashboards for clarity first, spectacle second.
  - Use fast, low-friction transitions for search and filter changes so the UI stays responsive.
  - Define reduced-motion and mobile downgrade rules before implementation starts.
- Open Questions:
  - Should the workflow install GSAP now, or only prepare the spec and phased plan?
  - This page does not appear to use GSAP yet. Should the refactor introduce GSAP or stay within the current motion stack?
  - Should charts and metrics animate softly for readability, or should they feel more dramatic and presentational?
  - Should the hero motion feel premium and cinematic, or quiet and product-focused?
  - Do you want smooth scrolling like Lenis, or should native scrolling stay untouched?
  - Should forms use only subtle focus polish, or do you want stronger field and success-state motion?
- Constraints:
  - Data-heavy screens need clarity-first motion and low distraction.
  - Forms and auth flows need utility-first motion, not theatrical timing.
  - Documentation surfaces should prioritize readability over spectacle.
