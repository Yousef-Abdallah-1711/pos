# SERVIO operational bento and feature showcase revision

Updated 2026-11-27. Source of truth: the active `frontend/apps/landing` app in this worktree. This report covers the two supplemental landing directives. The full site handoff remains `docs/landing-handoff.md`.

## Decision and design review

| Option | Clarity and product evidence | Mobile, RTL, themes, accessibility | Decision |
| --- | --- | --- | --- |
| A. Asymmetric interactive bento | One dominant POS area and two smaller situations give the actual screens an editorial hierarchy. | Grid can stack without absolute positioning; each state is a real button. | Selected. |
| B. Operational timeline | Explains sequence well but duplicates the seven-stage order story immediately below. | Easy to stack, but a second timeline would lengthen the page. | Rejected for repetition. |
| C. Layered collage | Could suggest a busy service day, but floating tickets would rely on fabricated UI fragments. | More fragile for Arabic, small screens, focus order, and reduced motion. | Rejected for authenticity and robustness. |

The old `.cv-problem-grid` and its three identical text cards were removed. `OperationalBento.tsx` now uses compact EN/AR copy, a large POS card, and smaller online-request and next-shift cards. The actual POS, dashboard, and kitchen screenshots come from `public/images/product/servio/matrix/` via `ProductCapture` and `CaptureMedia`; those components choose EN/AR and light/dark files and keep synthetic-data captions plus a full-capture link. The animated three-step labels are explicitly illustrative. They change only component state and never call the backend. The dashboard capture depicts the operations workspace, not a live incoming request; its annotation explains the verified staff-review workflow. The bridge still targets `#how-orders-move`. **Remaining requirement:** the approved brief called for visible product-level selection, order and status changes; the current static captures with changing explanatory labels do not demonstrate those UI changes. Obtain verified product fragments or synthetic capture states before treating that interaction as complete.

Three visual passes were made: (1) replace the large text-only cards with a bento and product captures; (2) remove an oversized empty area in the dominant card and reduce the supporting cards; (3) refine mobile/RTL/theme layouts and responsive capture crops. The desktop section measured about 1,495px before the compacting pass and about 1,308px after it at 1440×900. The heading stays within two lines on wide English desktop.

## Showcase root cause and correction

Chrome measured the prior 1440×900 state with a 76.7px header. The whole showcase section was pinned at about 78px, while the first slide began at 441px and ended at 1,083px. The heading and large section padding were inside the pinned element, so the capture extended 183px beyond the viewport as horizontal motion started.

Architecture A was selected: the introduction scrolls normally and a dedicated `.cv-showcase-stage` pins below the measured header. At the same 1440×900 size, the new stage measured about 90–889px and the first slide about 156–889px. The seven-slide track uses `track.scrollWidth - viewport.clientWidth` for travel. The active counter finds the nearest measured slide offset, including RTL, rather than rounding the overall scroll percentage. The stage keeps its progress buttons visible, and its capture remains in a side-by-side layout with the text. ScrollTrigger and its context revert on viewport/media changes or unmount.

At widths ≤900px, viewport heights <760px, or reduced-motion preference, the stage does not pin. The viewport provides native horizontal scrolling with snap, and the previous/next buttons remain available. Keyboard focus on an off-screen feature link reveals its slide. The showcase retains `contain: layout paint` on its viewport to prevent document-width overflow. No other homepage section or route was redesigned by this correction.

## Files changed

| Concern | Path |
| --- | --- |
| Homepage composition | `frontend/apps/landing/src/components/marketing/ConversionLanding.tsx` |
| New operational section | `frontend/apps/landing/src/components/marketing/OperationalBento.tsx` |
| Showcase stage and measured progress | `frontend/apps/landing/src/components/marketing/FeatureShowcase.tsx` |
| Motion media-query option | `frontend/apps/landing/src/components/marketing/desktopMotion.ts` |
| EN/AR problem copy and steps | `frontend/apps/landing/src/content/conversionCopy.ts` |
| Bento and responsive Coral styling | `frontend/apps/landing/src/styles/marketing-coral.css` |
| Stage, track, short-height fallback | `frontend/apps/landing/src/styles/marketing-interactive.css` |
| Browser regression script | `frontend/apps/landing/scripts/verify-operational-landing.mjs` |
| Design and handoff | `frontend/apps/landing/DESIGN.md`, `docs/landing-handoff.md`, `docs/landing-interactive-experience/progress.md`, `.gsap/` |

## Browser evidence and status

- Baseline images: `artifacts/problems-before.png` and `artifacts/showcase-before.png`.
- New layout images: `artifacts/operational-revision/` contains locale/theme mobile bento, the loaded desktop POS card, and desktop first/final showcase screenshots. An earlier stale dev-bundle screenshot was discarded.
- Chromium: 104/104 cases passed across 26 viewports in EN/AR and light/dark, including all specified mobile sizes, tablet sizes, large desktop, 1280×600/1440×650/1920×700, and two landscape sizes. Assertions covered page overflow, section/card count, locale/theme capture choice, bento bounds, single pin, first and final slide fit, final progress, second-slide progress, and reverse progress. JSON: `artifacts/operational-revision/results.json` (the Chromium run).
- Firefox and WebKit: 16/16 representative mobile and 1440×900 cases passed across EN/AR and light/dark after allowing asynchronous theme hydration and GSAP initialization to settle. Initial premature checks in WebKit exposed the need to wait for these states; the wait was added to the regression script. JSON: `artifacts/operational-revision/cross-browser/results.json`.
- Reduced motion: four Chromium EN/AR × light/dark cases passed with no pin spacer, native horizontal overflow, and keyboard focus revealing the final slide. JSON: `artifacts/operational-revision/smoke/results.json`.
- Host development route smoke: ten homepage/feature/demo routes returned 200. Six shared routes returned 500 because the pre-existing `frontend/packages/form-foundation/node_modules` links cannot resolve React, Zod, Lucide, or libphonenumber in this Windows checkout. This host dependency failure is separate from the landing changes; a clean production container build and route smoke are required before release.

## Verification still required for release

The clean Docker build, container route smoke, final lint/tests, CI checks, and pull-request merge must be recorded here before a GO verdict. Physical iOS/Android and screen-reader sessions remain unavailable in this environment. The product captures use synthetic sample data; public marketing components do not prove production order, payment, or tenant behavior. The prior landing report recorded mobile LCP above its 2.5s target; a new performance pass is still needed after these changes.

## Verification update (2026-11-27)

The repository's generated Windows `form-foundation/node_modules` links remained unreliable after a local repair; use the clean container for reproducible verification. The clean Node 24.22.3/pnpm 10.34.9 production build completed with exit 0, including TypeScript and all 46 static page outputs. The production-container browser script passed 120 viewport/locale/theme cases (104 Chromium, eight Firefox, eight WebKit), four reduced-motion/focus cases and 32 localized route checks. This first container run preceded the final showcase measurement-cache and compact-card CSS edits; the matrix is being rerun on the final source. The local unit suite passed 17 files and 93 tests. Lint originally failed because the React Hooks rule override applied to file extensions outside Next's plugin registration; limiting the override to the same extensions restored lint to zero errors and 11 existing warnings. Physical devices and assistive technologies remain untested, and no new Lighthouse measurement has yet been collected.

The final-source production-container rerun passed 120 browser cases (104 Chromium and eight each Firefox/WebKit), four reduced-motion/focus cases and 32 localized route checks. The final clean build and typecheck passed; lint passed with 11 warnings and zero errors. Unit tests passed 93/93 in the container with `NODE_ENV=test` and production public URL build variables cleared. Running the suite with the image's production environment caused React test-runtime failures; retaining the image's public URL variables caused two URL assertion failures, so the test-specific environment is necessary. These browser checks are emulated and do not prove physical-device behavior.

Current verdict: **CONDITIONAL GO for the two revised sections; CI and new performance measurements remain pending.**
