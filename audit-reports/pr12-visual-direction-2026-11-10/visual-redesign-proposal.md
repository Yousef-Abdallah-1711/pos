# POOS product-led landing page visual direction

**Repository:** `Yousef-Abdallah-1711/pos`  
**PR:** [#12](https://github.com/Yousef-Abdallah-1711/pos/pull/12)  
**Reviewed source head:** `1a709a0aec0bee164085f6ea92d4a520116909bd`  
**Proposal date:** 2026-11-10  
**Scope:** evidence-backed design proposal only. No application source, database, running service, or PR metadata was changed.

## Executive decision

**VISUAL AUDIT: PARTIAL — proposal and concept boards complete; authenticated product inspection and screenshot capture blocked.**  
**Recommended direction:** Direction 1, Premium Product Showcase, conditional on obtaining authentic, current product captures.  
**Implementation readiness:** **CONDITIONAL GO** to select a visual direction; **NO-GO** for implementation that depends on product proof until an isolated API and disposable account are available and the screens are recaptured and verified.

This direction uses the product interface as the primary visual evidence and keeps supporting copy concise. It relies less on exact order-state claims than a workflow-led concept. Today, however, the only landing-page dashboard image cannot be independently traced back to its original capture. Do not ship that image as verified proof until its source is recovered or the screen is recaptured.

## 1. Current experience: verified critique

The PR #12 audit recorded a coherent brand foundation: a purple palette, Alexandria and IBM Plex Sans Arabic, Arabic RTL and English LTR, a localized marketing page, and responsive layouts. The redesigned page and auth UI passed the recorded local layout checks at 390, 767, 1024, 1440, and 1920 CSS-pixel widths, with no observed horizontal overflow. The bilingual foundation should be retained.

The current hero centers its product proof on a small Access Control crop. That is too narrow to explain an all-in-one restaurant management system or the approved differentiator, “One platform for front and back of house.” The page audit also found long text-led sections with repeated visual rhythm, so the product story lacks clear visual variety. These are visual-hierarchy findings from the audited source and screenshots, not claims about user conversion.

The prior technical audit found issues that affect marketing confidence and page delivery: the generated Open Graph image URL resolves through locale middleware to a 404; `/robots.txt` and `/sitemap.xml` return landing HTML rather than their expected formats; canonical and hreflang metadata are absent. These are verified findings in the linked technical audit, not addressed by this design proposal.

Prior visual evidence: [English desktop 1440](../pr12-post-implementation-audit-2026-11-10/evidence/home-en-css-1440.png), [English mobile 390](../pr12-post-implementation-audit-2026-11-10/evidence/home-en-css-390.png), [Arabic desktop 1440](../pr12-post-implementation-audit-2026-11-10/evidence/home-ar-css-1440.png), and [mobile navigation](../pr12-post-implementation-audit-2026-11-10/evidence/mobile-menu-en.png). These are landing/auth screenshots, not authenticated dashboard product captures.

## 2. Screenshot authenticity and capture inventory

### Existing hero asset

| Asset | Classification | Evidence and limits |
|---|---|---|
| `frontend/apps/landing/public/images/poos-access-control-capture.png` | **UNVERIFIED** | SHA-256 `29B65CD22F92E8C2DB367500F6690C34993B9F8B6784B4813A26B09FC3271AA6`; 970×450. The provenance note says it is a crop from a real sample-workspace Dashboard → Settings → Access Control screen, but the original full-size capture is not in the repository. The source cannot be reconstructed or independently compared. Commit `2361e51a306a57c2bc03739dd8a296a6f61ae186` introduced the file. No matching original was located among the retained Phase 1 screenshots. Treat as unverified until independently recaptured. |
| Three new visual direction SVGs | **MOCKUP / CONCEPT** | They are labeled wireframes with explicit screenshot slots. They contain no simulated product interface or metrics and must not be represented as actual POOS screens. |
| Historical Phase 1 screenshots | **UNVERIFIED for marketing use** | Local screenshots exist for landing/login/register, storefront/menu, and super-admin at 320/768/1440 sizes. Per-image runtime, date, data, and authentication provenance is not sufficient to certify them as current marketing proof. |
| Newly verified authenticated product screenshots in this run | **0** | No product captures were made. |

### Product-screen inventory

Source presence below confirms that a route/component exists in this checkout; it does **not** prove successful runtime behavior, current visual state, data authenticity, or advertising readiness.

| Product area | Source route/component to inspect | Runtime / advertising status |
|---|---|---|
| Dashboard overview | `frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx` | Source present; live API unavailable; capture and behavior unverified. |
| POS and order creation | `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `.../pos/components/POSTerminal.tsx` | Source present; payment/order completion unverified. |
| Kitchen Display System | `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/page.tsx`, `.../kds/components/KDSGrid.tsx`, `KDSOrderCard.tsx` | Source present; real order-state sequence unverified. |
| Products and categories | `frontend/apps/dashboard/src/app/[locale]/dashboard/products/page.tsx`, `.../categories/page.tsx`, `frontend/apps/dashboard/src/components/products/ProductList.tsx` | Source present; current disposable-data runtime unverified. |
| Inventory | `frontend/apps/dashboard/src/app/[locale]/dashboard/inventory/page.tsx` | Source present; stock-state semantics and API data unverified. |
| Staff | `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/employees/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/users/page.tsx` | Source present; staff lifecycle and tenant data unverified. |
| Permissions | `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/access-control/page.tsx`, `frontend/apps/dashboard/src/components/settings/RBACSettings.tsx` | Source present; live role/permission behavior unverified. Existing crop is unverified. |
| Restaurant storefront | Dashboard theme builder under `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/`; public storefront screenshots exist in Phase 1 evidence | A historical UI capture exists, but exact route/data/date provenance and current ordering behavior are insufficient for marketing certification. |
| Reports | `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/page.tsx`, `.../reports/sales/page.tsx` | Source present; report values and live data are unverified. Do not display numeric proof. |

**Capture library status:** no verified product captures are available from this run. Once isolated infrastructure is available, retain each unmodified original in a dated evidence directory and record route, viewport, UTC capture time, environment, disposable-data fixture, authentication role, SHA-256, and any derivative crop transformation. Redact secrets and personal data. Mark each feature “ready to advertise” only after its interaction and displayed data are verified end-to-end.

## 3. Application inspection and blockers

- The local landing and dashboard services observed in prior work were started from separate dirty worktrees (`H:\pos-registration-validation-remediation` and `H:\pos-unified-phase2a-2c`). They were preserved and not stopped or modified.
- No API responded on `localhost:8080`; browser inspection of the dashboard also recorded a refused connection to `localhost:18082`.
- There was no isolated API plus disposable account suitable for authenticated inspection. No customer or persistent database data was touched.
- Consequently, this run could not verify dashboard, POS, KDS, inventory, staff, role assignment, reports, payment, or storefront behavior against live disposable data. No real product screenshot was captured.

## 4. Reference research

The review used official product pages as references for communication principles, not as templates to copy:

- [Linear](https://linear.app/homepage) — lead with the interface and a focused product promise; keep supporting copy short.
- [Stripe](https://stripe.com/) — vary product demonstration and explanatory compositions across the page instead of repeating one section pattern.
- [Square Restaurants](https://squareup.com/us/en/point-of-sale/restaurants) — organize the story around restaurant jobs and capabilities.
- [Toast products](https://pos.toasttab.com/products) — name operational product areas clearly so a restaurant operator can scan the offering.
- [Vercel](https://vercel.com/) — use strong typographic hierarchy and restrained high-contrast product framing.

Apply the underlying lessons in POOS’s own purple and bilingual system. Do not reproduce another company’s identity, art, screenshots, copy, or distinctive layout.

## 5. Header and hero proposal

### Header

Use a light, softly tinted navigation bar with a compact, distinctive POOS wordmark and one small ledger-line motif that suggests front-of-house and back-of-house entries joining the same shift. Keep this motif as a brand cue, not as a complex navigation metaphor. Use a restrained border or shadow after scroll instead of a permanently heavy floating pill.

Desktop order in English: logo; Product, How it works, For restaurant groups; language switch; Log in; one prominent “Create your restaurant workspace” action. Mirror semantic order for Arabic RTL rather than mechanically moving each visual element. Keep focus order logical in both locales. The mobile menu should open as a clearly labeled panel, keep the registration action visible, close on Escape, return focus to its trigger, and expose a 48×48px minimum touch target.

### Hero

Use the approved positioning without expanding it into unverified claims:

- **Arabic headline:** `نظام متكامل لإدارة المطاعم`
- **Arabic support:** `منصة واحدة تجمع العمل أمام الضيوف وخلف الكواليس.`
- **Arabic primary action:** `أنشئ مساحة مطعمك`
- **Arabic secondary action:** `استكشف المنتج`
- **English headline:** `All-in-one restaurant management`
- **English support:** `One platform for the work in front of the house and behind it.`
- **English primary action:** `Create your restaurant workspace`
- **English secondary action:** `Explore the product`

Arabic copy is meaning-aligned rather than literal. Keep its line breaks authored for RTL and Alexandria, not generated by mirroring English. Avoid “real-time,” sales totals, order counts, stock alerts, workflow status, customer quotes, savings, uptime, or growth promises until their accuracy and public use are approved.

Composition: Arabic copy on the right and a large, legible verified dashboard capture on the left; English copy on the left and capture on the right. Frame the screenshot as a quiet browser window with a narrow title bar and subtle purple ambient field. Do not add floating metric or notification cards until their exact data and behavior are verified. On mobile, order content as eyebrow/value proposition, headline, support, primary CTA, secondary CTA, then a single large crop with horizontal scrolling avoided. Preserve readable labels; do not shrink a full desktop dashboard to an illegible thumbnail.

## 6. Three design directions

The boards are actual SVG visual concepts. Each visibly labels its product screenshot area as a slot and states that it is not product evidence.

| Direction | Visual | Composition | Benefits | Trade-offs |
|---|---|---|---|---|
| **1 — Premium Product Showcase** | [Open concept board](evidence/direction-1-premium-showcase.svg) | Hero built around one large dashboard image; alternating capability sections; subtle browser framing; minimal floating UI. | Strongest product-led introduction; simple story; least dependent on unverified workflow transitions. | Requires one high-quality, authentic overview capture; a single screen may not prove the full suite by itself. |
| **2 — Interactive Restaurant Workflow** | [Open concept board](evidence/direction-2-workflow-story.svg) | Follow an order through connected front-of-house and kitchen scenes, with a vertical mobile timeline. | Makes cross-team coordination concrete if backend states and transitions are validated. | Highest verification burden; exact labels/transitions are not yet confirmed. Animation and interactivity can obscure the story or hurt performance. |
| **3 — Modern SaaS Editorial** | [Open concept board](evidence/direction-3-editorial-bento.svg) | Strong editorial type, varied split panels and bento compositions, alternating light and dark purple surfaces. | Distinct page rhythm; accommodates multiple verified product captures and supports a broad product map. | More content-heavy; risks repeating generic cards unless every panel has a distinct communication purpose. |

## 7. Section-by-section plan

1. **Hero / product overview:** one verified dashboard screenshot and the approved positioning. No fake “live” overlay.
2. **Front-of-house: POS and orders:** one POS screen plus a short explanation of the verified operator task. Do not assert payment completion until exercised end-to-end.
3. **Kitchen:** KDS capture and workflow description only after the actual status model and transitions are confirmed. The proposed sequence `Order Created → Kitchen Queue → Preparing → Ready → Completed` is a hypothesis to validate, not approved copy.
4. **Products and inventory:** product/inventory screen pair or one deliberate crop. Any stock warning must be visible in the captured app and tied to disposable data.
5. **People and permissions:** staff/role view with a verified role and permission state; avoid implying a permission boundary until access tests pass.
6. **Restaurant website and online ordering:** pair storefront and mobile capture only after the route, menu data, cart, checkout, and order handoff are verified.
7. **Reports:** use a real, explicitly labeled sample workspace; no public numeric claim. Omit this section if values cannot be safely explained.
8. **Audience close:** speak to both independent restaurants and multi-location groups without claiming enterprise readiness or unsupported multi-location features.

Use no five-card row repeated for each feature. Alternate image/text split, wide product panel, compact capability index, and a workflow diagram only where it explains a validated relationship.

## 8. Workflow and floating interface specification

Treat the connected workflow as a **conditional concept**. Before using it, inspect backend and frontend state enums, event handlers, cancellation/refund paths, and actual role permissions. Capture each state from one disposable scenario and record its source. Use a static, accessible sequence by default; motion should be a short, one-time reveal with `prefers-reduced-motion` support, no autoplay loop, and no information conveyed by color alone. On mobile, use a vertical timeline with labels and a visible current/next state only if the data model supports it.

Floating UI cards are excluded from the first approved direction until verified. If later approved, they must be crops or faithful values from the same verified sample session, be visually labeled as illustrative sample workspace data where appropriate, avoid covering the screenshot’s controls, and remain static for reduced motion. Never imply live customer activity from seeded demo data.

## 9. Mobile, Arabic, accessibility, and performance

- Validate at 390px, 768px, and a wide desktop viewport. The existing audit checked 390, 767, 1024, 1440, and 1920 CSS px; 768px should be explicitly added in the next design acceptance pass.
- Build mobile sections intentionally: readable screen crop, no tiny full-dashboard miniature, a vertical process diagram, and no decorative cards over the headline or key controls.
- Preserve `lang` and `dir` per locale. Check Arabic punctuation, numerals, line wrapping, button width, mixed Latin product labels, and screenshot orientation in actual RTL layout.
- Keep Alexandria for Arabic display type and IBM Plex Sans Arabic for body copy, subject to rendered glyph and line-wrap review. Maintain visible focus, semantic heading order, keyboard menu operation, descriptive image alternative text, and contrast of at least 4.5:1 for normal text. Decorative concept labels must not masquerade as app controls.
- Keep a stable hero media aspect ratio and explicit dimensions to prevent layout shift. Use responsive image sizes and modern formats; lazy-load below-fold captures, but prioritize the hero image. Avoid large animated video, unbounded parallax, and loading every screenshot at mobile resolution when it is not visible.
- The prior unthrottled local sample recorded LCP/FCP around 156ms and early layout-shift sum 0.0158. It is one warm local observation, not a production or mobile performance acceptance result. Collect repeatable cold-cache samples at target devices/network profiles and report median and p75 CLS/LCP/INP; retain the existing development CLS improvement as a separate observation.

## 10. Provenance and feasibility

For every screenshot, retain the original lossless capture and a manifest containing SHA-256, capture date/time, URL/route, viewport, browser/build, environment, tenant fixture, user role, data origin (disposable seed or real), redactions, and whether the image is original or a named derivative crop. Keep transformations reproducible and never overwrite originals. The current provenance note is not sufficient independent evidence because the referenced source capture is absent.

Design implementation is technically feasible in the existing Next.js landing app (`frontend/apps/landing/src/app/[locale]/page.tsx` and its section components), but visual approval does not clear the technical audit defects: OG image route, robots/sitemap content types, canonical/hreflang, live API auth, and production security headers remain separate engineering acceptance items. The audit found local typecheck, 35 frontend unit tests, lint, and a production build with disposable local URL configuration passed; live auth was blocked. This proposal does not rerun or upgrade those results.

## 11. Recommendation and owner decision

Choose **Direction 1 — Premium Product Showcase** as the working design direction, because it gives the clearest product-led story with the fewest assumptions about unverified order workflow semantics. Before implementation, provide an isolated API with disposable roles/tenants and recover or recreate a verifiable Dashboard/POS/KDS/inventory/storefront sample session. Then capture the overview and the highest-value supporting screens, verify them against the running UI, and approve exact usage/crops. If product proof remains unavailable, keep the current implementation and do not substitute a wireframe for a screenshot.

### Exact blockers

1. No responding API at the inspected local endpoints (`localhost:8080` and dashboard browser API `localhost:18082`).
2. No disposable authenticated account/session for product route inspection.
3. Source original for the existing 970×450 access-control crop is not retained.
4. No newly verified product screenshots: **0**.
5. Product workflow transitions and displayed values have not been verified against the running backend.

## Evidence index

- [Prior implementation audit](../pr12-post-implementation-audit-2026-11-10/audit-report.md)
- [Prior machine-readable acceptance summary](../pr12-post-implementation-audit-2026-11-10/evidence/acceptance-results.json)
- [Direction 1 — Premium Product Showcase](evidence/direction-1-premium-showcase.svg)
- [Direction 2 — Interactive Restaurant Workflow](evidence/direction-2-workflow-story.svg)
- [Direction 3 — Modern SaaS Editorial](evidence/direction-3-editorial-bento.svg)

---

**No redesign was implemented, no PR was merged, and no existing service, database, dirty worktree, or untracked mock artifact was modified by this proposal.**
