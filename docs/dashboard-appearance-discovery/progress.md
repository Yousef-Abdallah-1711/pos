# Discovery progress

| Area | Status | Evidence and remaining work |
| --- | --- | --- |
| Source inventory | IN PROGRESS | Raw routes, colors, typography and SERVIO PNG CSVs generated; 63 static dashboard URLs visited as a synthetic English owner. Classification, other roles, dynamic routes, active asset usage and full application coverage remain. |
| Shell and sidebar | IN PROGRESS | Source behavior documented; authenticated desktop/mobile, RTL, permission and keyboard observation remains. |
| Appearance lifecycle | IN PROGRESS | Personal settings and implicit storefront sync traced; tenant switching, caching and failure behavior remain unverified. |
| Database/API | IN PROGRESS | User settings, RLS and site-theme endpoints identified; full schemas, effective role mapping and runtime isolation need review. |
| Future proposal | IN PROGRESS | Initial ownership, API direction, token subset, UX and test plan drafted; complete token registry, validations and migration decisions remain. |
| Visual baseline | IN PROGRESS | Sixteen authenticated synthetic captures cover EN/AR, light/dark, collapsed desktop, Arabic mobile drawer, empty POS/KDS, a synthetic POS cart, six additional routes and one staff account with no assigned roles. Additional restricted users, populated states, other route variants and physical devices remain. See `baseline/README.md`. |
| Implementation | NOT STARTED | No dashboard, database or storefront appearance behavior has been changed. |

Landing closure is being tracked in PR #137 and `docs/landing-interactive-experience/operational-bento-showcase-report.md`. At discovery time the PR has not been merged. The last local mobile Lighthouse sample missed the LCP target and the detailed POS screenshot interaction remains incomplete; do not label the landing project fully verified. CI and the final browser rerun are in progress.

Additional local Arabic mobile Lighthouse sample on the latest production build: performance 45, accessibility 100, LCP 6.35 s, CLS 0, TBT 2.59 s. This single run occurred with other containers active and is not a controlled comparison or a physical-device result. English and Arabic LCP both remain above the 2.5 s target in the measured samples.

PR #137's `auth-browser-acceptance` job is skipped by the explicit PR-number/workflow condition in `.github/workflows/production-ci.yml:170`; it is not a successful fresh authentication integration run for this PR. Route smoke checks validate page rendering, not end-to-end registration or tenant provisioning.
