# Browser and test results

## Reproducible successful acceptance

[CI run 43579188642](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43579188642) passed the complete verified workflow for PR head `60691410a2245c1f194efc613dd77c3f39c1a705`, including all five production image builds. Evidence: [summary](evidence/candidate-verification.json); full `chromium-smoke-results` and `backend-test-results` artifacts on the run.

| Check | Passed | Failed | Scope |
|---|---:|---:|---|
| Backend integration/unit tests | 1,242 | 0 | Two environment-specific tests skipped |
| Landing unit/component tests | 64 | 0 | Eleven files |
| Dashboard unit/component tests | 665 | 0 | Fifty-two files |
| Marketing browser assertions | 1,638 | 0 | All public routes at all six widths in EN/AR, metadata, axe, navigation, annual mobile pricing |
| New-owner trial assertions | 23 | 0 | Two accounts; annual/monthly persistence, 14-day entitlement, dashboard, category creation, isolation, logout |
| Existing landing visual samples | 60 | 0 | Ten widths × three routes × two locales |
| Existing mobile-menu checks | 4 | 0 | Open and Escape close |

The two skipped backend tests are the opt-in disposable Admin Orders benchmark and the WS05-specific Super Admin audit database test. None of the new commercial/trial checks were skipped. Existing authentication/navigation/library lint warnings remain; lint, typecheck and test typecheck passed for all frontend applications.

## Final release gate

`frontend/apps/landing/scripts/marketing-acceptance.mjs` now covers every public route in both languages at 360, 390, 768, 1024, 1440 and 1920px: 168 route/viewport samples, plus interaction journeys and annual mobile pricing captures. It validates raw server language/direction, single h1, canonical URLs, unique descriptions, hreflang, WebSite structured data, screenshot decoding, document overflow, axe WCAG A/AA rules, browser/console errors, and failed HTTP responses.

`trial-registration.mjs` requires explicit synthetic-write enablement and landing/API URLs. It creates unique workspaces, verifies the subscription through the authenticated API, performs a real first menu action through the dashboard, checks workspace separation and revokes sessions. No credentials are written to reports.

`capture-landing-acceptance.mjs` retains the pre-existing 60-sample regression suite with updated navigation and hidden/lazy-image handling. `lighthouse-acceptance.mjs` uses pinned Lighthouse 13.6.0 and rejects performance below 90, accessibility below 95, or CLS at/above 0.1.

`ops/aws/run-ci-browser-checks.sh` runs all suites and returns failure if any suite fails. Nothing is converted to a warning. Final focus/logout assertion counts and the exact approved main/release SHA are recorded in the [PR delivery record](https://github.com/Yousef-Abdallah-1711/pos/pull/43) after the final required workflow; do not infer AWS acceptance from a PR image build.

The browser method is repository Playwright tooling in the existing pinned CI browser image; Playwright MCP was not available.

## Cross-browser release gate

The same public-route, viewport, axe, keyboard-focus, metadata and pricing checks run in Firefox as well as Chromium in the pinned Playwright image. Both suites must pass; Firefox results are uploaded separately under `marketing-firefox/`. [Official Docker reproduction guidance](https://playwright.dev/docs/docker). The final PR delivery record includes both actual outcomes.
