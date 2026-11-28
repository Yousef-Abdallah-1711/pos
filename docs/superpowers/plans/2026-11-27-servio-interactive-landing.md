# SERVIO Interactive Landing Experience Implementation Plan

> Execution method: inline continuation in the existing worktree, as explicitly directed in the supplied brief.

**Goal:** Extend the existing SERVIO Coral landing page into a localized, responsive, product-led experience with verified claims, accessible GSAP motion, repaired shared shell, and evidence-based handoff.

**Architecture:** Keep the localized Next.js route and marketing shell server-rendered. Add focused homepage sections and one isolated client motion controller that imports GSAP only after mount; all copy and product states remain in the server-rendered page so reduced motion and no-JavaScript use remain understandable. Preserve current shared routes, auth flows, tokens, and existing user work.

**Tech Stack:** Next.js 16.4.6, React 19.3.1, TypeScript, existing Coral CSS tokens, GSAP Core and ScrollTrigger (loaded dynamically), Vitest, Playwright.

**Spec:** User-provided SERVIO brief at C:/Users/Yousef/.codex/attachments/489b48e9-4856-48e6-96c2-5a57f886b5e9/Pasted text.txt; repository guidance in AGENTS.md, CODING_STANDARDS.md, frontend/apps/landing/AGENTS.md and DESIGN.md.

## Global Constraints

- Preserve Coral colors, Fraunces/Manrope/IBM Plex Sans Arabic, authentic captures, current route destinations, session controls, attribution and synthetic-data disclosure.
- Keep unsupported financial, payment, industry, and end-to-end claims out of marketing copy.
- Use native scroll and one pinned story; honor reduced motion, keyboard access, RTL shaping, and mobile simplification.
- Do not change package-management strategy, downgrade dependencies, modify unrelated apps, or commit, push, merge, deploy, or delete protected artifacts.
- Keep the current worktree and its uncommitted changes intact.
- Verify with current environment evidence; Node currently reports v24.15.0, pnpm v9.0.0, and port 3001 is already occupied.

## Review Focus

- Reduced-motion visitors must see complete content without pinned or horizontal traps; test the motion controller's reduced-motion branch and browser behavior.
- Locale/theme changes must not leave stale captures or direction-specific layout; test EN/AR and light/dark route states.
- Keyboard users must reach every story stage and feature; test tab/arrow/escape behavior and link destinations.
- Narrow viewports must not create document overflow; test 320, 360, 390, 768, 1024, 1440 and 1920 widths.
- API downtime must remain an explicit integration limit without disrupting page rendering; record console/network evidence and test guest route behavior.

## File Map

- Homepage composition: frontend/apps/landing/src/components/marketing/ConversionLanding.tsx.
- Localized homepage strings: frontend/apps/landing/src/content/conversionCopy.ts.
- New motion behavior: frontend/apps/landing/src/components/marketing/HomeMotion.tsx and its focused tests.
- New product-story sections: focused components under frontend/apps/landing/src/components/marketing/home/.
- Homepage and shared-shell layout: frontend/apps/landing/src/styles/marketing-coral.css, marketing-navigation.css, Footer.tsx, Navbar.tsx.
- Dependency metadata: frontend/apps/landing/package.json and frontend/pnpm-lock.yaml.
- Source of truth for motion implementation: frontend/apps/landing/.gsap/.
- Evidence and handoff: docs/landing-interactive-experience/ and docs/landing-handoff.md.

## Approaches

1. Replace the Coral page wholesale with a new marketing framework. Rejected because it risks existing account/navigation behavior and duplicates the approved design system.
2. Reuse the current page, retain the shared shell, add focused sections and progressive motion around existing server-rendered captures. Selected because it preserves route behavior and authentic product assets while meeting the requested story architecture.
3. Build a separate preview and migrate later. Rejected because the request explicitly targets the live application and would leave its broken or incomplete experience unresolved.

## Tasks

### Task 1: Baseline and capability evidence

Inspect actual landing, dashboard, storefront, backend paths, tests, capture provenance, ports and runtime. Save baseline browser evidence and a capability register with explicit confidence and advertise-safety status. No implementation claims until workflow evidence exists.

Acceptance: baseline report separates existing API refusal from page behavior; capability register cites source paths and test evidence; the running port 3001 service is reused.

### Task 2: Shared shell and footer

Use the Coral token file and responsive styles to repair footer grid structure and desktop/mobile navigation behavior. Add focused regression coverage for group destinations and structural layout, then verify the four locale/theme combinations in the browser.

Acceptance: CTA, brand, navigation groups and legal/disclosure row remain distinct; desktop/tablet/mobile and RTL layouts have no overlap or document overflow.

### Task 3: Motion foundation and hero/order story

Add GSAP to the landing workspace only if absent. Implement a scoped, lifecycle-safe motion controller with responsive matchMedia rules, one pinned order journey, restrained hero reveal, reduced-motion fallback and ScrollTrigger refresh only after relevant media/fonts settle. Use existing verified copy and captures.

Acceptance: all story text and order states are available without animation; one pinned region only; no timeline or trigger duplication after remount/locale changes.

### Task 4: Product demonstrations and showcase sections

Add the storefront and kitchen workflow presentations, layered typography transition, horizontal capability showcase, connected business-module presentation, and management visibility section. Use authentic captures and synthetic demonstration annotations; include only features supported by the evidence register.

Acceptance: each section has a distinct purpose, accessible non-scroll alternative, localized copy and real route destination where one exists; no unsupported money/accounting correctness claim.

### Task 5: Setup, pricing, FAQ and shell finish

Add stacking setup cards and localized pricing details only from the existing commercial offer source; preserve current FAQ, CTA, auth funnel tracking, oversized footer mark and synthetic data disclosure.

Acceptance: account CTAs preserve locale and billing behavior; pricing values agree with source data; no unsupported discount or trial claim.

### Task 6: Full acceptance and handoff

Run landing typecheck, lint, full Vitest suite, production build, route regression, browser visual matrix, keyboard/reduced-motion, responsive overflow, console/network and available performance/accessibility checks. Record unverified API integration and unavailable browser engines honestly. Update DESIGN.md and docs/landing-handoff.md.

Acceptance: evidence-backed status for every requirement, changed file inventory and preview URL; no deployment or commit.

