# POSS Landing Page Redesign Implementation Plan

> **For agentic workers:** Implement natively in the current worktree. Preserve the user's existing edits.

**Goal:** Redesign the public landing page using the installed Hallmark workflow while preserving its content, routes, captures, and working interactions.

**Architecture:** Keep the existing Next.js locale route and component boundaries. Rework the marketing page's shared tokens, layout and component styling in place; retain the accessible locale-aware markup and product evidence. Use Hallmark's before/after records and screenshots to compare the rendered page.

**Tech Stack:** Next.js App Router, React, TypeScript, CSS, next/image, Lucide, Hallmark sloplint, MCP Playwright.

**Spec:** `frontend/apps/landing/DESIGN.md`, `frontend/apps/landing/BRAND_GUIDELINES.md`, `frontend/apps/landing/PRODUCT.md`, installed Hallmark skill.

## Global Constraints

- Keep the Shared Shift Ledger visual idea, token palette and genuine product captures.
- Preserve Arabic RTL and English LTR, existing routes, links, IDs, handlers, form fields and analytics hooks.
- Do not add unsupported claims, pricing, trial terms, testimonials or customer logos.
- Keep and do not overwrite unrelated files already changed in the shared worktree.
- Make no dependency or lockfile changes.

## Review Focus

- Arabic RTL layout and line wrapping — inspect `/ar` at phone and desktop widths.
- English hero fit and overflow — inspect `/en` at 375, 1280, 1440 and 1920 pixel widths.
- Landing actions, navigation and pricing selector — compare retained handlers and links in Hallmark records.
- Focus visibility, reduced motion and contrast — run Hallmark sloplint and inspect keyboard/focus rules.
- Claims and brand inconsistency — keep current page wording; report SERVIO/POOS and public offer approval discrepancies.

---

### Task 1: Redesign landing visual system

**Files:**
- Modify: `frontend/apps/landing/src/styles/marketing-base.css`
- Modify: `frontend/apps/landing/src/styles/marketing-hero.css`
- Modify: `frontend/apps/landing/src/styles/marketing-product.css`
- Modify: `frontend/apps/landing/src/styles/marketing-conversion.css`
- Modify: `frontend/apps/landing/src/styles/marketing-navigation.css`
- Modify: `frontend/apps/landing/src/styles/marketing-responsive.css`
- Modify: `frontend/apps/landing/src/components/marketing/HomeHero.tsx`

**Interfaces:**
- Consumes: Existing marketing class names, token values in `DESIGN.md`, `getSiteCopy(locale)`, `TrialLink`, `ProductCapture`.
- Produces: Responsive, token-led page styles with the shared-shift ledger layout and an above-the-fold product hero.

- [x] Update the existing CSS token layer to match the locked palette in light and dark themes.
- [x] Refine hero typography, content measure, stage composition and spacing; remove decorative status/uppercase stage treatment while retaining translated copy and navigation.
- [x] Refactor remaining marketing rules to use CSS tokens, logical properties, visible focus, motion preferences, mobile-safe widths and correct contrast.

### Task 2: Render and compare the complete landing experience

**Files:**
- Modify: affected marketing components only when required for the approved layout.
- Create/update: `frontend/apps/landing/.hallmark/after.json`, `after.txt`, screenshots and `compare.html`.
- Create: `frontend/apps/landing/.hallmark/finish-review.md`.

**Interfaces:**
- Consumes: Task 1 page styles and the existing dev server on port 3001.
- Produces: Hallmark-kept comparison, responsive rendered captures, and a finish review documenting issues and waivers.

- [x] Run Hallmark sloplint with `--render --kept` against the same files and `/en` URL as the baseline.
- [x] Inspect English and Arabic mobile/desktop screenshots and correct actionable visual/accessibility issues.
- [x] Create the Hallmark comparison and review notes; preserve all pre-existing user worktree changes.
