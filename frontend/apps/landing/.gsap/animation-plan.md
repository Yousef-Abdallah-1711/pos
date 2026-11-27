> Archived generator scaffold. For the implemented homepage, use `.gsap/current-state.md`, `DESIGN.md`, and `docs/landing-handoff.md`. This scaffold contains unverified guesses.

# Animation Plan - landing

## Workflow State

- Current Mode: gsap-refactor
- Resume State: Discovery complete, phased execution pending
- Last Updated By:
- Suggested Next Command: Implement P01 only, then update all .gsap artifacts.
- Active Phase: P01

## Target Scope

- Pages:
- Primary Goal:
- Constraints:

## Implementation Order

1. Read .gsap artifacts
2. Inspect current code
3. Resolve missing decisions
4. Implement section by section
5. Verify mobile and reduced motion
6. Update artifacts after changes

## Recipes By Section

| Page | Section | Recipe | Status | Notes |
|---|---|---|---|---|
| homepage | Navigation | navbar polish on scroll | Planned | Make the Navigation section feel intentional, modern, and aligned with the page hierarchy. |
| homepage | Feature Grid | staggered card reveal + hover depth | Planned | Make the Feature Grid section feel intentional, modern, and aligned with the page hierarchy. |
| homepage | Stats | count-up on enter | Planned | Make the Stats section feel intentional, modern, and aligned with the page hierarchy. |
| homepage | Timeline | step-by-step reveal | Planned | Make the Timeline section feel intentional, modern, and aligned with the page hierarchy. |
| homepage | Form | field-focus polish + success transitions | Planned | Make the Form section feel intentional, modern, and aligned with the page hierarchy. |
| homepage | Hero | hero-text-reveal + layered media parallax | Planned | Make the Hero section feel intentional, modern, and aligned with the page hierarchy. |
| homepage | Testimonials | quote reveal + soft carousel polish | Planned | Make the Testimonials section feel intentional, modern, and aligned with the page hierarchy. |
| homepage | CTA | cta emphasis + magnetic hover | Planned | Make the CTA section feel intentional, modern, and aligned with the page hierarchy. |
| homepage | Pricing | pricing-card emphasis + trust hierarchy | Planned | Make the Pricing section feel intentional, modern, and aligned with the page hierarchy. |
| homepage | FAQ | accordion timing polish | Planned | Make the FAQ section feel intentional, modern, and aligned with the page hierarchy. |
| homepage | Team | profile stagger + soft hover depth | Planned | Make the Team section feel intentional, modern, and aligned with the page hierarchy. |
| homepage | Showcase | parallax showcase or marquee | Planned | Make the Showcase section feel intentional, modern, and aligned with the page hierarchy. |

## Phases

| Phase | Page | Section | Objective | Recipe | Status |
|---|---|---|---|---|---|
| P01 | homepage | Navigation | Make the Navigation section feel intentional, modern, and aligned with the page hierarchy. | navbar polish on scroll | Planned |
| P02 | homepage | Feature Grid | Make the Feature Grid section feel intentional, modern, and aligned with the page hierarchy. | staggered card reveal + hover depth | Planned |
| P03 | homepage | Stats | Make the Stats section feel intentional, modern, and aligned with the page hierarchy. | count-up on enter | Planned |
| P04 | homepage | Timeline | Make the Timeline section feel intentional, modern, and aligned with the page hierarchy. | step-by-step reveal | Planned |
| P05 | homepage | Form | Make the Form section feel intentional, modern, and aligned with the page hierarchy. | field-focus polish + success transitions | Planned |
| P06 | homepage | Hero | Make the Hero section feel intentional, modern, and aligned with the page hierarchy. | hero-text-reveal + layered media parallax | Planned |
| P07 | homepage | Testimonials | Make the Testimonials section feel intentional, modern, and aligned with the page hierarchy. | quote reveal + soft carousel polish | Planned |
| P08 | homepage | CTA | Make the CTA section feel intentional, modern, and aligned with the page hierarchy. | cta emphasis + magnetic hover | Planned |
| P09 | homepage | Pricing | Make the Pricing section feel intentional, modern, and aligned with the page hierarchy. | pricing-card emphasis + trust hierarchy | Planned |
| P10 | homepage | FAQ | Make the FAQ section feel intentional, modern, and aligned with the page hierarchy. | accordion timing polish | Planned |
| P11 | homepage | Team | Make the Team section feel intentional, modern, and aligned with the page hierarchy. | profile stagger + soft hover depth | Planned |
| P12 | homepage | Showcase | Make the Showcase section feel intentional, modern, and aligned with the page hierarchy. | parallax showcase or marquee | Planned |

## Validation Checklist

- [ ] Reduced motion covered
- [ ] Mobile-heavy effects reviewed
- [ ] Existing code inspected before changes
- [ ] Page artifact updated after implementation
- [ ] Discovery questions resolved or consciously deferred
- [ ] Tasks file created or refreshed
- [ ] Phase files created or refreshed

## Workflow Snapshot - homepage

- Mode: gsap-refactor
- Framework: next
- Package Manager: npm
- Existing Motion Stack: No motion libraries detected
- Matched Files: multi-tenant-login.integration.test.tsx, page.tsx, page.tsx, page.tsx, layout.tsx, page.tsx, page.tsx, page.tsx
- Inferred Product Type: Marketing Site, Docs / Knowledge Base, SaaS App

## Implementation Plan - homepage

- Read .gsap artifacts before editing code.
- Work in one major section phase at a time.
- Finish hero or top-priority story beat before supporting sections.
- Add spectacle only where the story earns it.
- Update phase, tasks, and page artifacts after each section.

## Detected Project Signals - homepage

- Sections: Navigation, Feature Grid, Stats, Timeline, Form, Hero, Testimonials, CTA, Pricing, FAQ, Team, Showcase
- Structure Patterns: Modal interactions, Auth surface, Repeated collection render, Search and filtering, Data visualization
- Selector Samples: poos-skip-link, auth-page, auth-layout, auth-story, auth-story__rule, auth-form-panel, auth-form, auth-form__back, auth-form__intro, auth-form__fields
- Infrastructure: State, UI Kit, Testing
- Recommendations:
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
