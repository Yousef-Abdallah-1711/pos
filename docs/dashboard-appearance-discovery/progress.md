# Discovery progress

| Area | Status | Evidence and remaining work |
| --- | --- | --- |
| Source inventory | IN PROGRESS | Raw routes, colors, typography and SERVIO PNG CSVs generated; 63 static and three synthetic-ID dashboard URLs visited as an English owner. Classification, other roles, remaining dynamic routes, active asset usage and full application coverage remain. |
| Shell and sidebar | IN PROGRESS | Source behavior documented; authenticated desktop/mobile, RTL, permission and keyboard observation remains. |
| Appearance lifecycle | IN PROGRESS | Personal settings and implicit storefront sync traced; tenant switching, caching and failure behavior remain unverified. |
| Database/API | IN PROGRESS | User settings, RLS and site-theme endpoints identified; full schemas, effective role mapping and runtime isolation need review. |
| Future proposal | IN PROGRESS | Initial ownership, API direction, token subset, UX and test plan drafted; complete token registry, validations and migration decisions remain. |
| Visual baseline | IN PROGRESS | Twenty-four authenticated synthetic captures cover EN/AR, light/dark, collapsed desktop, Arabic mobile drawer, POS/KDS, additional routes, personal appearance and storefront theme settings, plus one staff account with no assigned roles. Additional restricted users, populated states, other route variants and physical devices remain. See `baseline/README.md`. |
| Implementation | NOT STARTED | No dashboard, database or storefront appearance behavior has been changed. |

Landing closure is being tracked in PR #137 and `docs/landing-interactive-experience/operational-bento-showcase-report.md`. At discovery time the PR has not been merged. The last local mobile Lighthouse sample missed the LCP target and the actual POS screenshot remains static with a separate illustrative cart; do not label the landing project fully verified. The 120-case browser matrix passed before the font-loading change and the new CI run is in progress.

The previous PR #137 image gate failed its Lighthouse performance threshold (EN 88 and AR 86 versus 90). A font preload change on landing head `2781efa7` improved single local comparisons in both languages but did not reach the 2.5-second LCP goal; its new CI run is pending. See `docs/landing-interactive-experience/performance-and-accessibility.md` for methodology and limitations. Do not merge by bypassing the failed gate.

After the font change, eight additional Chromium smoke cases covered EN/AR × light/dark × 390/1440-pixel widths on the new production build: each returned 200, retained 13 sections and expected direction/font family, and showed no page errors or document overflow. These focused checks do not replace the full three-engine matrix or the pending CI performance gate.

Additional local Arabic mobile Lighthouse sample on the latest production build: performance 45, accessibility 100, LCP 6.35 s, CLS 0, TBT 2.59 s. This single run occurred with other containers active and is not a controlled comparison or a physical-device result. English and Arabic LCP both remain above the 2.5 s target in the measured samples.

PR #137's `auth-browser-acceptance` job is skipped by the explicit PR-number/workflow condition in `.github/workflows/production-ci.yml:170`; it is not a successful fresh authentication integration run for this PR. Route smoke checks validate page rendering, not end-to-end registration or tenant provisioning.
