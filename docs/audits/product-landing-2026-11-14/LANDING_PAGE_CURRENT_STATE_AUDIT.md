# Landing page current-state audit

Audit date: 14 November 2026 (Africa/Cairo). Source revision: `ea23d95a213f46fd8cdc59648b56e8474fc679f2` on `main`. Public environment: **staging**, API release receipt `f733efb5e7f5730fea5a56a65d66d9c339e6c040`, status `candidate`.

**Updated evidence boundary:** the user authorized synthetic staging provisioning and testing. Authenticated browser/API/SQL checks now cover core orders, kitchen, stock, theme configuration, HR records, manual accounting and tenant suspension. See [authenticated evidence and limitations](evidence/authenticated/README.md). Successful tests are scoped to candidate `f733efb`; they do not establish full production readiness. Synthetic data was created and retained as audit evidence; application source, deployment and infrastructure were not changed.

Evidence: [browser observations](evidence/browser-observations.json), [HTTP observations](evidence/http-observations.json), [release receipt](evidence/release-probes.json), [source declarations](evidence/source-index.json), [schema definitions](evidence/schema-index.json), [existing tests](evidence/test-index.json). Paths below are repository-relative and pinned to the audited SHA; source schema is supplemented by fresh installed-schema inspection, 138 migration receipts and container identity observations in the authenticated evidence directory.

## What was inspected

Chromium 159 through Playwright, anonymous clean pages: EN and AR homepages at 1440×1000, 768×1024 and 390×844; each FAQ expanded; mobile navigation opened, Escape tested; Tab sampled. Both localized login and registration forms inspected without submission. Nine core dashboard routes redirected to landing login. Super Admin displayed its login screen. Public storefront homepage, Arabic homepage, menu, empty cart and empty checkout inspected. This initial anonymous pass created no records; the authorized continuation below created synthetic data only.

Six full-page screenshots are in `evidence/{en,ar}-{desktop,tablet,mobile}.png`. Desktop EN and mobile AR were visually inspected. DOM/overflow observations cover all six. No Safari, Firefox, real device, screen reader or field performance data was available. This is not a WCAG certification or Core Web Vitals assessment.

## Current page and messaging

| Surface | Observed content | Assessment |
|---|---|---|
| Hero | “All-in-one restaurant management.”; front/back-house description | Product category clear; concrete workflow and scope absent. Broad scope claim needs qualification |
| Hero visual | Real-looking role-management capture explicitly labeled synthetic sample workspace | Honest sample label, but secondary admin task is a poor explanation of everyday restaurant value. Embedded navigation labels are incidental mentions, not feature explanations |
| Platform section | Two large front/back-house panels and connection mark | Repeats category language without actions, outputs or examples |
| Audience | Independent restaurants and multi-location groups | Group-management claim requires proof beyond tenant isolation |
| How it works | Create account → workspace → choose area | Explains registration, not running a service or customer order |
| FAQ | Platform identity, audience, registration data, public EN/AR | Covers basic signup questions; no pricing, rollout, hardware, offline, payments, support or migration answers |
| Final CTA | Create workspace | Clear action, little additional reason to act |
| Footer | Homepage anchors, login, account creation | No pricing/contact/support/legal destinations in rendered footer |

English/Arabic content has equivalent high-level coverage. Arabic is readable and RTL rendering works at tested widths. The hero screenshot remains English even on Arabic; it does not demonstrate Arabic dashboard capability. Much of the page allocates large vertical space to repeated general statements; the mobile image consumes substantial space before meaningful product education.

## Route, link and conversion inventory

Homepage links in both locales point to `#product`, `#how-it-works`, `#audience`, `#faq`, localized `/auth/login`, localized `/auth/register`, homepage logo and `/images/product/roles-1440x900.png`. All displayed anchor IDs exist in source/DOM; CTA forms were reached. Primary CTA appears in navigation, hero and final block; mobile navigation exposes account creation. No prices, free-trial duration, “no card” promise, ROI numbers, testimonials or ratings were observed. Do not invent them.

`/en/features`, `/en/pricing`, `/en/how-it-works` returned 404; the current “how it works” is an in-page anchor. The route tree also includes auth/redirect aliases and super-admin access routes. Those technical entry points are not feature pages and were not exercised with tokens. No additional marketing routes were found in the landing route tree or sitemap. Arbitrary public URL enumeration is outside this crawl.

Forms show labels, optional phone and password guidance. Registration asks for restaurant subdomain and previews a tenant address. The initial anonymous pass did not execute registration. In the authorized continuation, unique registration returned 201, redirected to the dashboard and created an active 15-day trial. Duplicate restaurant name produced an inline error. Multi-organization selection and full onboarding remain unverified; the new storefront hostname failed TLS. Account creation should not be described as an instant, frictionless trial until synthetic acceptance establishes that behavior.

No sales/contact CTA or accessible live demo was found. Add these only when an operational destination and response owner exist. Establish commercial terms and privacy/terms pages before acquisition. Do not add a fake demo button or made-up price.

## Responsive behavior and accessibility

| Check | Result | Limit |
|---|---|---|
| Load | 200 in all six homepage views; no pageerror events | Does not prove all requests or workflows succeed |
| Overflow | No horizontal document overflow in six views | Tested sizes only; not all breakpoints/zoom levels |
| Language direction | EN ltr, AR rtl in document and computed body | Storefront localization is weaker, see below |
| Mobile menu | Opens; Escape closes and returns focus in both languages | Full keyboard traversal and screen-reader modal behavior pending |
| FAQ | Four native details controls expand | Complete screen-reader reading sequence pending |
| Semantics | Single H1, ordered H2/H3, skip link and image alt | Link purpose, all focus styles and visual contrast need manual review |
| Automated a11y | axe-core 4.10.3 reported 0 violations across requested WCAG A/AA tags in six homepage states | Incomplete `aria-prohibited-attr` checks in all; 3 contrast nodes need review on each desktop. Zero automated violations is not compliance |
| Project touch target standard | 9 visible targets below 48px in each sampled view | Mobile footer links about 44.8px tall, logo 36px, screenshot link 40px; WCAG exceptions/spacing need separate evaluation |

See `followup-observations.json` for axe nodes and `browser-observations.json` for bounding boxes. Do not call the site inaccessible based only on stricter project targets; do not declare it accessible based only on axe.

## Technical SEO

- Localized titles/descriptions observed, canonical URLs correct for EN/AR, reciprocal EN/AR hreflang present in document and sitemap. No x-default observed; it is an optional opportunity, not a failed requirement.
- Open Graph/Twitter title, image, dimensions and description present. No JSON-LD observed. Add appropriate Organization/SoftwareApplication data only with truthful commercial information; FAQ markup does not guarantee a search enhancement.
- `/robots.txt` allows public crawl and excludes localized auth routes; sitemap lists EN/AR homepage only. Crawl permission and indexing are different: responses carry **X-Robots-Tag: noindex, nofollow** and **X-Servio-Environment: staging**.
- Staging noindex is intentional protection, not an instruction to remove it now. A future approved production site needs an explicit indexing release check, correct canonical host, status codes, sitemap and auth noindex policy.
- Static/cached page delivery observed through Next.js/Caddy, cache HIT and s-maxage. This does not establish optimal caching for authenticated content or adequate production capacity.
- Navigation timing samples saved; no throttled Lighthouse, field LCP/INP/CLS, percentile traffic data, bandwidth-normalized load budget or image transfer audit completed. No performance score or conversion uplift claimed.
- HTML language/dir and descriptive alt exist. No broken hero image was recorded. App screenshot must remain authentic and sanitized.

## Product surfaces reached and access limitations

Anonymous Dashboard/POS/KDS/theme builder/inventory/HR/accounting/reports/billing requests redirect to login. The later authorized continuation logged in normally, traversed 79 route/viewport combinations and exercised scoped workflows; no authentication bypass was attempted. Super Admin is internal and not part of a restaurant package. Storefront `/en/menu` lists one “SERVIO staging QA item” at 25 EGP; the later authorized continuation submitted one synthetic guest order. Arabic storefront still displays English navigation/footer and hero content in this deployed candidate. Main includes localization/control changes that cannot be assumed deployed. The initial empty-state observations are supplemented by a later actual guest checkout→staff confirmation→KDS fixture at 27.50 EGP, followed by cancellation restoring stock. The staff confirmation action used the authenticated API; customer and kitchen views were browser-verified.

## Coverage methodology

Ten core buyer modules are the denominator: overview, cashier POS, KDS, orders/tables/delivery, catalog/inventory, HR, accounting, reports/analytics, storefront/theme builder, access/settings. A module is “meaningfully explained” only if visible marketing prose names a specific task or workflow, states a benefit and supplies context beyond a navigation label embedded in an image. Current page meaningfully explains **0/10 operational modules (0%)**. This deliberately measures explanatory coverage, not product completeness. Generic front/back-house copy and visible screenshot menu labels do not qualify. Public bilingual support is explicitly explained but outside the ten operational modules. Registration is described but also outside that denominator. Missing explanatory coverage is 10/10; no claim that every underlying module is ready to market.

## Authenticated product evidence affecting marketing

The service story is now supported by real synthetic POS/KDS and online-order sequences. The page still does not explain them. Screenshot candidates now exist in `evidence/authenticated`, including POS cart, kitchen ticket, guest checkout/confirmation/detail, populated HR, journal/ledger and theme editor. These are audit captures: POS totals are demonstrably inconsistent, and editor mock/live scope plus Arabic omissions must be disclosed. Do not publish captures as proof of corrected behavior. The new tenant hostname failed TLS despite registration advertising its URL. That belongs in onboarding/conversion readiness, not a landing-only visual critique. Super Admin tenant detail displayed a history-fetch error and restricted support access; internal operator tooling must not be marketed as restaurant features.
