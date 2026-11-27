> Archived generator scaffold. For the implemented homepage, use `.gsap/current-state.md`, `DESIGN.md`, and `docs/landing-handoff.md`. This scaffold contains unverified guesses.

# Page Animation Spec - homepage

## Page

- Page Route: /homepage
- Source File: C:/Users/Yousef/Desktop/Projects/POSS/frontend/apps/landing/src/__tests__/integration/multi-tenant-login.integration.test.tsx
- Page Status: Needs phased refactor

## Resume State

- Next Agent Action: open phase P01 and implement only that section.
- Blocking Questions: Should the workflow install GSAP now, or only prepare the spec and phased plan? | This page does not appear to use GSAP yet. Should the refactor introduce GSAP or stay within the current motion stack? | Should charts and metrics animate softly for readability, or should they feel more dramatic and presentational? | Should the hero motion feel premium and cinematic, or quiet and product-focused? | Do you want smooth scrolling like Lenis, or should native scrolling stay untouched? | Should forms use only subtle focus polish, or do you want stronger field and success-state motion?
- Discovery Confidence: Medium unless the agent verifies source files manually.
- Active Phase: P01

## Sections

### Hero

- Type:
- Elements:
- Status: Not started

### Content Grid / Cards

- Type:
- Elements:
- Status: Not started

### Stats / Numbers

- Type:
- Elements:
- Status: Not started

### CTA / Footer

- Type:
- Elements:
- Status: Not started

## Mobile Rules

- Simplifications:
- Disabled Effects:

## Reduced Motion

- Fallback Behavior:

## Recommended Motion Directions

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

## Scenario Notes

- Inferred Product Type: Marketing Site, Docs / Knowledge Base, SaaS App
- Motion Constraints:
  - Data-heavy screens need clarity-first motion and low distraction.
  - Forms and auth flows need utility-first motion, not theatrical timing.
  - Documentation surfaces should prioritize readability over spectacle.

## Phase Status

| Phase | Section | Recipe | Status | Notes |
|---|---|---|---|---|
| P01 | Navigation | navbar polish on scroll | Planned | Make the Navigation section feel intentional, modern, and aligned with the page hierarchy. |
| P02 | Feature Grid | staggered card reveal + hover depth | Planned | Make the Feature Grid section feel intentional, modern, and aligned with the page hierarchy. |
| P03 | Stats | count-up on enter | Planned | Make the Stats section feel intentional, modern, and aligned with the page hierarchy. |
| P04 | Timeline | step-by-step reveal | Planned | Make the Timeline section feel intentional, modern, and aligned with the page hierarchy. |
| P05 | Form | field-focus polish + success transitions | Planned | Make the Form section feel intentional, modern, and aligned with the page hierarchy. |
| P06 | Hero | hero-text-reveal + layered media parallax | Planned | Make the Hero section feel intentional, modern, and aligned with the page hierarchy. |
| P07 | Testimonials | quote reveal + soft carousel polish | Planned | Make the Testimonials section feel intentional, modern, and aligned with the page hierarchy. |
| P08 | CTA | cta emphasis + magnetic hover | Planned | Make the CTA section feel intentional, modern, and aligned with the page hierarchy. |
| P09 | Pricing | pricing-card emphasis + trust hierarchy | Planned | Make the Pricing section feel intentional, modern, and aligned with the page hierarchy. |
| P10 | FAQ | accordion timing polish | Planned | Make the FAQ section feel intentional, modern, and aligned with the page hierarchy. |
| P11 | Team | profile stagger + soft hover depth | Planned | Make the Team section feel intentional, modern, and aligned with the page hierarchy. |
| P12 | Showcase | parallax showcase or marquee | Planned | Make the Showcase section feel intentional, modern, and aligned with the page hierarchy. |

## Discovery Snapshot

- Matched Files: multi-tenant-login.integration.test.tsx, page.tsx, page.tsx, page.tsx, layout.tsx, page.tsx, page.tsx, page.tsx
- Sections Detected: Navigation, Feature Grid, Stats, Timeline, Form, Hero, Testimonials, CTA, Pricing, FAQ, Team, Showcase
- Repeated Components: LoginPage, RootDocument, Navbar, Link, FormErrorSummary, FormField, PasswordField, Loader2, Footer, AlertCircle
- Structure Patterns: Modal interactions, Auth surface, Repeated collection render, Search and filtering, Data visualization
- Selector Samples: poos-skip-link, auth-page, auth-layout, auth-story, auth-story__rule, auth-form-panel, auth-form, auth-form__back
