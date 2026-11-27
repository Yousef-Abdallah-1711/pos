# SERVIO landing app: current implementation and handoff

Updated 2026-11-27. Repository: `C:/Users/Yousef/Desktop/Projects/POSS`. App: `frontend/apps/landing`. Continue in **this worktree**. The landing implementation and evidence were prepared for integration through a pull request. Unrelated working files in other apps and local browser/audit artifacts were excluded; preserve them. No deployment was performed.

Repository remote: `https://github.com/Yousef-Abdallah-1711/pos`. Landing branch: `feature/servio-bento-showcase-20261127`, head `2781efa7`, [PR #137](https://github.com/Yousef-Abdallah-1711/pos/pull/137). At this update the PR is open. Its earlier image CI failed Lighthouse performance; a font-loading improvement was pushed and its new CI is pending. Verify the PR and remote main before claiming a merge. After integration, replace this paragraph with the verified merge commit and current branch. Do not interpret the commit SHA as the final remote `main` SHA.

## First steps for the next agent

1. Read `CODING_STANDARDS.md` in full. It is the repository's primary code policy.
2. Read `frontend/apps/landing/AGENTS.md`; for Next.js API/convention changes, read the relevant installed guide under `frontend/apps/landing/node_modules/next/dist/docs/`.
3. Read `frontend/apps/landing/DESIGN.md`, `.gsap/current-state.md`, and `docs/landing-interactive-experience/product-capability-register.md` before altering layout, motion, or product claims.
4. Run `git status --short` and protect existing user work. The local landing preview is available at `http://localhost:3001/en` and `/ar` when the dev or production server is running.

## What is live

The localized `/en` and `/ar` homepages are the real Next.js conversion page. Coral supplies the palette, type, photographic hero, auth UI, header, and footer; the stronger section rhythm and interaction system now run in the app. The audience is Egyptian restaurant and café owners. The story follows an online request through staff confirmation, kitchen preparation, and fulfillment, with counter POS and management tools also shown. The primary CTA goes to actual restaurant registration.

The old static concept HTML pages (`coral.html`, `hybrid.html`, `ultraviolet.html`, `vantage.html`) are not the live homepage. Do not use their preview URLs as a source of truth. The current source is `src/components/marketing/ConversionLanding.tsx`, `src/content/conversionCopy.ts`, `src/styles/marketing-coral.css`, `src/styles/marketing-interactive.css`, and `DESIGN.md`.

## Current homepage, in order

| # | Section / anchor | Component or source | What visitors can do |
| --- | --- | --- | --- |
| 1 | Hero | `ConversionLanding.tsx`, `HeroCaptureMotion.tsx` | See the real dashboard capture over a credited café photo; register or jump to order story. Desktop capture settles with scroll. |
| 2 | Operational bento | `OperationalBento.tsx`, `conversionCopy.ts` | Explore three local illustrative steps in each asymmetric product card, open the actual localized capture, or continue to the order story. No order or payment mutation occurs. |
| 3 | Order story / `#how-orders-move` | `OrderStory.tsx`, `orderStoryCopy.ts` | Select seven stages or scroll through them; see actual screen areas with clearly labeled sample order and separate order/kitchen states. CSS sticky stage on desktop. |
| 4 | Storefront / `#online-ordering` | `ConversionLanding.tsx`, `WorkflowDemo.tsx` | Select Browse, Review, Submit, Confirm to see an explanatory customer-order sequence and open the feature details. |
| 5 | Kitchen / `#kitchen` | `ConversionLanding.tsx`, `WorkflowDemo.tsx` | Select New, In progress, Ready, Completed ticket states and open the kitchen feature. |
| 6 | Layered product statement | `LayeredTypography.tsx` | See the POS screen and restrained desktop text/capture parallax. |
| 7 | Feature showcase / `#capabilities` | `FeatureShowcase.tsx`, `src/content/features.ts` | Explore seven linked real-capture slides. A dedicated desktop stage pins beneath the measured header; the introduction scrolls normally. Narrow, short-height, and reduced-motion viewports use native scroll snap and previous/next controls. |
| 8 | Ecosystem | `BusinessEcosystem.tsx` | Select seven workspace modules and open the corresponding feature route. Lines show grouping, not invented automatic data transfer. |
| 9 | Management | `ManagementExperience.tsx` | Switch between overview, inventory, and team captures; sample numbers are labeled as sample data. |
| 10 | Getting started / `#getting-started` | `StackedSetupCards.tsx` | Read account, menu, and team setup steps and start registration. Cards stack with restrained desktop scroll treatment. |
| 11 | Pricing bridge | `LandingPricing.tsx` | Open the real pricing route or register. No unsupported price or payment promise is hardcoded here. |
| 12 | FAQ | `ConversionLanding.tsx` | Expand native `details` answers. |
| 13 | Final CTA | `ConversionLanding.tsx` | Register a restaurant account through the real localized form. |

## Section implementation map

All section routes are `/{locale}`. Shared Coral tokens, typography, responsive breakpoints and RTL selectors are in `frontend/apps/landing/src/styles/marketing-coral.css`; interactive layout and motion fallbacks are in `src/styles/marketing-interactive.css`. `src/content/conversionCopy.ts` supplies EN/AR copy except the order story (`src/content/orderStoryCopy.ts`), feature showcase (`src/content/features.ts` and `src/i18n/messages/{en,ar}/features.json`), and other component-local labels. Product captures are selected by `src/components/marketing/{ProductCapture,CaptureMedia}.tsx` from `public/images/product/servio/matrix/`; hero photography is `public/images/marketing/cafe-interior.jpg`. The matrix has English/Arabic and light/dark variants with synthetic data and a full-capture link. The script `src/components/marketing/desktopMotion.ts` lazily loads GSAP for eligible media queries and cleans up on changes.

| Section and anchor | Component / selectors | Motion and fallback |
| --- | --- | --- |
| Hero / `#cv-hero-title` | `ConversionLanding.tsx`, `HeroCaptureMotion.tsx`; `.cv-hero`, `.cv-hero-capture` | Desktop capture motion; static on narrow or reduced-motion displays. |
| Problems / `#cv-problem-title` | `OperationalBento.tsx`; `.cv-problem`, `.cv-bento-*` | Local card timelines and illustrative step controls; no pin. Mobile/reduced motion keeps manual controls. The actual POS screenshot remains static. |
| Order journey / `#how-orders-move` | `OrderStory.tsx`; `.cv-story-*` | Desktop CSS sticky stage with ScrollTrigger state synchronization; narrow layout stacks stage and steps. |
| Online / `#online-ordering` | `ConversionLanding.tsx`, `WorkflowDemo.tsx`; `.cv-online`, `.cv-workflow-demo--online` | Selectable illustrative states; no API mutation. |
| Kitchen / `#kitchen` | `ConversionLanding.tsx`, `WorkflowDemo.tsx`; `.cv-kitchen`, `.cv-workflow-demo--kitchen` | Selectable illustrative preparation states; no API mutation. |
| Layered statement | `LayeredTypography.tsx`; `.cv-layered-*` | Scoped desktop parallax; static fallback. |
| Showcase / `#capabilities` | `FeatureShowcase.tsx`; `.cv-showcase-*`, `.cv-feature-slide-*` | Desktop stage-only ScrollTrigger pin starts below the measured header; track travel is measured overflow. Width ≤900px, height <760px, and reduced motion use native scroll snap. RTL direction is handled explicitly. |
| Ecosystem | `BusinessEcosystem.tsx`; `.cv-ecosystem-*` | Selectable module nodes; desktop line animation and stacked mobile layout. |
| Management | `ManagementExperience.tsx`; `.cv-management-*` | Selectable real capture views and labeled sample metrics; reduced motion avoids nonessential effects. |
| Setup / `#getting-started` | `StackedSetupCards.tsx`; `.cv-setup-*` | Desktop sticky cards with scoped scroll effect; static stack on narrow screens. |
| Pricing | `LandingPricing.tsx`; `.cv-pricing-*` | Links to pricing and registration; responsive split stacks. |
| FAQ | `ConversionLanding.tsx`; `.cv-faq-*` | Native `details` disclosure; no animation dependency. |
| Final CTA | `ConversionLanding.tsx`, `TrialLink.tsx`; `.cv-close` | Localized registration link with funnel attribution. |

The shared `Navbar.tsx`/`Footer.tsx` and `src/styles/marketing-navigation.css` provide the header, locale and theme controls, mobile menu and footer. RTL and light/dark use the root `lang`, `dir`, and theme class. Test the complete page using `scripts/verify-operational-landing.mjs` with `SERVIO_QA_URL`, and the shared routes and forms with the other scripts and `tests/registration-e2e/`. Motion details and current verification gaps are documented in `.gsap/current-state.md` and `docs/landing-interactive-experience/operational-bento-showcase-report.md`.

English and Arabic content have equal priority. The route sets `lang` and `dir`; `next-themes` stores the choice under `servio-marketing-theme`. The hero and footer retain espresso treatment in both modes. Arabic copy is real Arabic, not auto-translation at runtime.

## Source map

| Concern | Paths |
| --- | --- |
| Homepage route and metadata | `frontend/apps/landing/src/app/[locale]/page.tsx`, `src/app/[locale]/layout.tsx`, `src/components/RootDocument.tsx` |
| Section composition | `src/components/marketing/ConversionLanding.tsx` |
| Interactive sections | `src/components/marketing/{HeroCaptureMotion,OperationalBento,OrderStory,WorkflowDemo,LayeredTypography,FeatureShowcase,BusinessEcosystem,ManagementExperience,StackedSetupCards,LandingPricing}.tsx` |
| Main homepage copy | `src/content/conversionCopy.ts`, `src/content/orderStoryCopy.ts` |
| Feature definitions and detail copy | `src/content/features.ts`, `src/i18n/messages/{en,ar}/features.json`, and related marketing JSON |
| Shared header and footer | `src/components/Navbar.tsx`, `src/components/Footer.tsx`, `src/components/navigation/`, `src/content/navigation.ts` |
| CTA and funnel | `src/components/marketing/TrialLink.tsx`, `src/lib/funnel.ts`, `src/lib/config.ts` |
| Real product captures | `src/components/marketing/{ProductCapture,CaptureMedia}.tsx`, `src/content/captures.json`, `public/images/product/servio/matrix/` |
| Café hero asset | `public/images/marketing/cafe-interior.jpg` |
| Theme and shared marketing shell | `src/components/RootDocument.tsx`, `src/app/marketing.css`, `src/styles/marketing-base.css`, `src/styles/marketing-navigation.css`, `src/styles/marketing-responsive.css` |
| Coral and interactive styling | `src/styles/marketing-coral.css`, `src/styles/marketing-interactive.css`, `src/styles/marketing-conversion.css`, `src/styles/marketing-captures.css`, `src/styles/marketing-workflows.css` |
| Auth screen copy and behavior | `src/app/[locale]/auth/{login,register}/page.tsx`, `src/i18n/messages/{en,ar}/legacy.json`, `src/lib/api.ts`, `src/lib/auth-session.ts`, `src/lib/config.ts`, shared form foundation package |
| Auth handoff and integration QA | `frontend/apps/dashboard/src/lib/{api,authSession}.ts`, `frontend/apps/landing/playwright.registration.config.ts`, `frontend/apps/landing/tests/registration-e2e/`, `docs/landing-interactive-experience/compose-auth-qa.override.yml` |
| Browser and performance evidence | `frontend/apps/landing/scripts/{verify-interactive-landing,verify-operational-landing,smoke-marketing-routes}.mjs`, `docs/landing-interactive-experience/artifacts/{acceptance,operational-revision}/`, `artifacts/route-smoke.json`, `artifacts/lighthouse-*-final.json` |
| Design and motion decisions | `frontend/apps/landing/DESIGN.md`, `.gsap/current-state.md` |

CSS class ownership: `cv-*` is homepage composition and the interactive system; `sv-*` is the shared marketing shell, capture, header, and footer; `auth-*` styles the account screens. The bento classes are `.cv-bento-*` in `marketing-coral.css`; obsolete `.cv-problem-grid` styles were removed. The footer is a real grid (`.sv-footer-top`, `.sv-footer-groups`) with four navigation groups and `.sv-footer-signature` oversized brand text. Keep the feature showcase viewport's `contain: layout paint`; removing it caused Chrome to expose the off-screen track as page width. The dedicated `.cv-showcase-stage` is the desktop GSAP pin target; its introduction is outside the pin. The desktop order-story panel is `.cv-story-stage-slot` with CSS `position: sticky`; pinning it with GSAP previously distorted the grid.

## Routes and destination ownership

| URL | Main source / purpose |
| --- | --- |
| `/` | `src/app/(entry)/page.tsx` redirects to default locale. |
| `/{locale}` | `src/app/[locale]/page.tsx` renders `MarketingShell` and `ConversionLanding`. |
| `/{locale}/features` | `src/app/[locale]/features/page.tsx`: feature directory. |
| `/{locale}/features/{slug}` | `src/app/[locale]/features/[slug]/page.tsx`: details, limitations, captures. Slugs: `pos`, `kitchen-display`, `online-ordering`, `website-builder`, `inventory`, `hr`, `analytics`. |
| `/{locale}/how-it-works` | `src/app/[locale]/how-it-works/page.tsx`: detailed workflow. |
| `/{locale}/pricing` | `src/app/[locale]/pricing/page.tsx`: offer and billing details. |
| `/{locale}/demo` | `src/app/[locale]/demo/page.tsx`: product explorer. |
| `/{locale}/contact`, `/privacy` | Their route `page.tsx` files: contact/demo and privacy. |
| `/{locale}/auth/login`, `/auth/register` | Real login and registration forms. The main CTA targets registration with a billing query; preserve it. |
| `/{locale}/auth/super-admin-login` | Separate admin auth route. |

`{locale}` is `en` or `ar`. `src/middleware.ts` redirects legacy paths; `src/app/[locale]/loading.tsx` and `error.tsx` handle route states. `src/app/{sitemap.ts,robots.ts,opengraph-image.tsx,icon.svg}` provide metadata assets. Header/footer labels come from `src/i18n/messages/{en,ar}/marketing.json`; navigation destinations come from `src/content/navigation.ts`.

## Product truth and capture handling

- `backend/internal/domain/order.go` defines order statuses `pending → confirmed → preparing → ready → out_for_delivery → delivered` for delivery paths. `frontend/apps/dashboard/src/stores/kdsStore.ts` defines separate kitchen ticket states `new → in_progress → ready → completed` (with recall handling). Online orders require staff confirmation before kitchen preparation. The storefront path uses cash on delivery; do not claim online card payment or automatic kitchen acceptance.
- Dashboard code exposes POS, orders, kitchen display, inventory, team/HR, overview/reports, and theme-builder routes. A disposable local backend was brought up for registration/login testing, but those product modules remain source-backed capabilities, not newly verified production transactions.
- `docs/landing-interactive-experience/product-capability-register.md` holds source paths and claim confidence. The prior financial audit in `audit-reports/full-financial-system-discovery-2026-11-25/` found accounting gaps; do not claim fully reliable accounting or financial outcomes.
- `ProductCapture` selects actual captured application screens by locale and theme from the matrix. Synthetic data and redacted QA identifiers are disclosed in captions. The order/demo overlays are explanatory local UI state, not live backend operations. Preserve provenance when replacing imagery.

## Motion and fallback rules

`gsap` and `@gsap/react` are direct landing dependencies in `frontend/apps/landing/package.json` and `frontend/pnpm-lock.yaml`. Client islands register `ScrollTrigger`, scope contexts to their root, and revert them on cleanup. The bento has short local timelines plus keyboard/touch step buttons; no backend call. The showcase measures header height and `track.scrollWidth - viewport.clientWidth`, pins only the stage at widths >900px and heights ≥760px, and tracks the nearest actual slide offset. At widths ≤900px, heights <760px, or `prefers-reduced-motion: reduce`, it becomes a touch-friendly native horizontal scroll container. The order story becomes a static stage above selectable steps on narrow screens, and setup cards stop sticking. No essential content is hidden behind animation.

## Commands and verification

Run from `frontend/`:

```powershell
npx --yes pnpm@10.34.9 --filter @pos-saas/landing dev
npx --yes pnpm@10.34.9 --filter @pos-saas/landing typecheck
npx --yes pnpm@10.34.9 --filter @pos-saas/landing lint
npx --yes pnpm@10.34.9 --filter @pos-saas/landing test:run
$env:NEXT_PUBLIC_LANDING_URL='https://localhost:3001'; npx --yes pnpm@10.34.9 --filter @pos-saas/landing build
npx --yes pnpm@10.34.9 --filter @pos-saas/landing start
```

From `frontend/apps/landing/` against port 3001:

```powershell
node scripts/verify-interactive-landing.mjs
node scripts/smoke-marketing-routes.mjs
```

For the disposable registration integration run, start Postgres and backend from the repository root with the two compose files below, then start the dashboard and landing apps in separate terminals. The override supplies current SQL mounts, the `api` command, local CORS, and loopback port 8080. Run migrations before starting the backend. The stack is named `pos-test`; do not use the unrelated existing `pos-postgres` container.

```powershell
docker compose -f docker-compose.pos-test.yml -f docs/landing-interactive-experience/compose-auth-qa.override.yml up -d postgres
docker compose -f docker-compose.pos-test.yml -f docs/landing-interactive-experience/compose-auth-qa.override.yml build backend
docker compose -f docker-compose.pos-test.yml -f docs/landing-interactive-experience/compose-auth-qa.override.yml run --rm backend migrate
docker compose -f docker-compose.pos-test.yml -f docs/landing-interactive-experience/compose-auth-qa.override.yml up -d backend
npx pnpm --filter @pos-saas/dashboard dev
```

With landing on port 3001 and dashboard on 3002, set `REG_E2E_LANDING_URL=http://localhost:3001`, `REG_E2E_DASHBOARD_URL=http://localhost:3002`, `REG_E2E_API_URL=http://localhost:8080/api/v1`, and `REG_E2E_BROWSER_CHANNEL=chrome`; run `npx pnpm --filter @pos-saas/landing exec playwright test --config playwright.registration.config.ts tests/registration-e2e --reporter=line`. After testing, `docker compose -f docker-compose.pos-test.yml -f docs/landing-interactive-experience/compose-auth-qa.override.yml down -v` removes the disposable services and database. Check `docker ps -a --filter name=pos-test` afterward; a detached one-off `backend migrate` container remained in this run and had to be removed separately before its Docker network could be removed. The unrelated `pos-postgres` container was left running. The registration suite writes accounts to its temporary database.

The production build, typecheck, and lint passed. Lint had 11 existing warnings in auth/API code and zero errors. Vitest passed **17 files / 92 tests**; jsdom prints its known unsupported navigation messages during otherwise passing auth tests. The production Chrome matrix passed **12/12** EN/AR × light/dark × 1440/768/390 cases with full-page PNGs and interactions. Firefox and WebKit passed **8/8 each** for EN/AR × light/dark × 390/1440 layout checks. The route/width sweep passed **44/44**: 30 localized marketing/auth routes at 390px and 14 localized homepage widths from 320px to 1920px. Evidence: `docs/landing-interactive-experience/artifacts/acceptance/results.json`, its twelve PNGs, and `docs/landing-interactive-experience/artifacts/route-smoke.json`. A manual reduced-motion desktop check found zero pin spacers and a native horizontally scrollable showcase. The local machine has Node 24.15.0 while the workspace declares ≥24.22.3; build and tests still completed.

For integration, a clean production build also passed using `frontend/Dockerfile.prod` with Node 24.22.3 and pnpm 10.34.9. The current Windows `frontend/packages/form-foundation/node_modules` has damaged generated junctions, so a subsequent host build failed to resolve its packages. A pinned host reinstall could not replace those junctions. Use the clean container build or repair the generated local links before relying on another host build.

Real registration E2E passed against the disposable `pos-test` Postgres/backend and local dashboard: tenant owner creation, HttpOnly session cookie, dashboard handoff, confirmed logout and denial, login again, role-management access, duplicate email/subdomain validation, and overlong phone rejection. The other registration suites passed session restoration, outage and DNS recovery, Arabic RTL, keyboard/theme controls, field announcements, and resilience checks. Across focused runs all **18/18 registration scenarios** passed. `frontend/apps/dashboard/src/lib/api.ts` now preserves the active locale on a 401 redirect via `getDashboardLoginUrl`. The local test stack used `docker-compose.pos-test.yml` plus `docs/landing-interactive-experience/compose-auth-qa.override.yml`; it is disposable, uses mock payments and a temporary database, and does not prove production delivery or billing.

Lighthouse 13.6.0 lab measurements: English mobile performance/accessibility **90/100**, LCP **3.31 s**, CLS **0**; Arabic mobile **78/100**, LCP **4.01 s**, CLS **0**; English desktop **100/100**, LCP **0.75 s**, CLS **0.004**. Full metrics and raw report paths are in `docs/landing-interactive-experience/performance-and-accessibility.md`. The mobile LCP target of 2.5 s remains unmet. The final scoped auth-outage notice CSS change was rebuilt and checked at 1200 px after the Lighthouse run.

## Remaining limits and next work

- Registration and login were verified only against the disposable local stack. Production API, payment processing, tenant DNS/provisioning, and full POS-to-kitchen transactions were not exercised.
- `src/components/marketing/WorkflowDemo.tsx` and several new components contain their own EN/AR UI copy; keep the pair aligned when editing. `conversionCopy.ts` and `orderStoryCopy.ts` own the longer homepage narrative.
- Cross-browser checks were layout/route checks, not full visual-diff or assistive technology sessions. A real screen-reader pass and on-device iOS/Android test remain release checks.
- The `.gsap` generator scaffold guessed unrelated sections. `.gsap/current-state.md` is authoritative; all scaffold files are marked archived.
- Asset and acceptance folders already contained user work. No broad cleanup was attempted. The older `docs/landing-interactive-experience/artifacts/baseline-en-light-1440.png` actually captured a dark saved-theme state; its filename is misleading and it should not be used as a light-mode comparison.

## Updating this implementation

Change section order and homepage composition in `ConversionLanding.tsx`; update EN and AR copy together in `conversionCopy.ts` and `orderStoryCopy.ts`. Change spacing, palette, dark mode, RTL, and breakpoints in the Coral/interactive CSS, keeping the feature containment. Change header/footer destinations and labels through `navigation.ts`, both `marketing.json` files, `Navbar.tsx`, and `Footer.tsx`. Change product capture provenance through the matrix, `captures.json`, and capture components. Change registration or billing links through `TrialLink.tsx` and `src/lib/funnel.ts`. For auth behavior, trace the localized forms, session state, API client, and dashboard callback together. Re-run the production browser matrix and route smoke after shared CSS or shell changes.
