> Archived generator scaffold. For the implemented homepage, use `.gsap/current-state.md`, `DESIGN.md`, and `docs/landing-handoff.md`. This scaffold contains unverified guesses.

# Animation Spec - landing

## Workflow State

- Owner Workflow: gsap-refactor
- Source Of Truth: .gsap artifacts
- Resume Rule: Read these files before asking new questions

## Philosophy

- Personality: Subtle & Professional
- Density: Moderate
- Scroll Behavior: Scroll-triggered reveals

## Global Effects

- [ ] Smooth scroll
- [ ] Magnetic buttons
- [ ] Page transitions
- [ ] Parallax images

## Performance Rules

- Device Priority: Mobile-first
- Reduced Motion: Required
- Heavy Effects On Mobile: Disabled unless intentionally approved
- Preferred Animated Properties: transform, opacity

## Missing Decisions

- Brand references:
- Sections needing strongest emphasis:
- Elements to avoid over-animating:

## Project Intelligence

- Inferred Archetypes: Marketing Site, Docs / Knowledge Base, SaaS App
- Supporting Infrastructure:
  - State: zustand
  - UI Kit: @radix-ui/react-dropdown-menu, @radix-ui/react-icons, @radix-ui/react-label, @radix-ui/react-slot, @radix-ui/react-toast
  - Testing: @playwright/test, @testing-library/jest-dom, playwright, vitest
- Project Constraints:
  - Data-heavy screens need clarity-first motion and low distraction.
  - Forms and auth flows need utility-first motion, not theatrical timing.
  - Documentation surfaces should prioritize readability over spectacle.

## Interview Strategy

- Interview Mode: targeted-discovery
- Priority Categories: Technical Stack, Section Behavior, Creative Direction, Workflow Clarification
- Priority Questions:
  - Should the workflow install GSAP now, or only prepare the spec and phased plan?
  - This page does not appear to use GSAP yet. Should the refactor introduce GSAP or stay within the current motion stack?
  - Should charts and metrics animate softly for readability, or should they feel more dramatic and presentational?
  - Should the hero motion feel premium and cinematic, or quiet and product-focused?
  - Do you want smooth scrolling like Lenis, or should native scrolling stay untouched?
  - Should forms use only subtle focus polish, or do you want stronger field and success-state motion?

## Spec-Driven Rules

- Planning must happen before implementation.
- Work must be split into phases with one major section per phase.
- Each phase must have explicit reduced-motion and mobile downgrade notes.
- The next agent should be able to resume from artifacts alone.

## Discovery Snapshot

- Framework: next
- Package Manager: npm
- Packages:
  - next@16.4.6
  - react@19.3.1
  - tailwindcss@^3.4.17
  - typescript@^5.7.2
- Routes Detected: /(entry), /[locale], /[locale]/contact, /[locale]/demo, /[locale]/features, /[locale]/how-it-works, /[locale]/pricing, /[locale]/privacy, /[locale]/features/[slug], /[locale]/auth/login, /[locale]/auth/register, /[locale]/auth/super-admin-login

## Brand And Design Signals

- Colors:
  - #f4e8d9
  - #fffaf2
  - #f0ecf8
  - #fff
  - #171327
  - #575066
  - #f8f2e8
  - #2b211d
  - #3a2e27
  - #e17861
- Fonts:
  - var(--font-plex),sans-serif
  - var(--font-fraunces),Georgia,serif
  - var(--font-plex), var(--font-manrope), sans-serif
  - var(--font-manrope), var(--font-plex), sans-serif
  - var(--font-alexandria), var(--font-plex), sans-serif
  - var(--font-manrope),sans-serif
- Tone Hints:
  - ai
  - 3d
- CSS Variables:
  - --poos-paper
  - --poos-surface
  - --poos-surface-raised
  - --poos-ink
  - --poos-ink-soft
  - --poos-muted
  - --poos-rule
  - --sv-ink
  - --sv-muted
  - --sv-violet
  - --sv-paper
  - --sv-line
- Visual Direction:
  - gradient-rich
  - glass-friendly
  - rounded-ui
  - depth-layering
  - editorial-typography
- Blur Usage Signals: 2
- Backdrop Filter Signals: 0

## Questions To Resolve

- Should the workflow install GSAP now, or only prepare the spec and phased plan?
- This page does not appear to use GSAP yet. Should the refactor introduce GSAP or stay within the current motion stack?
- Should charts and metrics animate softly for readability, or should they feel more dramatic and presentational?
- Should the hero motion feel premium and cinematic, or quiet and product-focused?
- Do you want smooth scrolling like Lenis, or should native scrolling stay untouched?
- Should forms use only subtle focus polish, or do you want stronger field and success-state motion?

## Phase Strategy

- Active Page: homepage
- Rule: implement one major section per phase.
- Planned Phase Count: 12
- First Phase: P01 - Navigation
