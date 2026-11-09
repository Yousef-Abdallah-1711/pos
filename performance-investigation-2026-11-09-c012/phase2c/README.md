# Phase 2C browser measurements — 2026-11-09

These measurements used the labeled disposable stack `pos-p2a-ws05-e2e-20261109-c012-final`, synthetic tenant/restaurant `2/2`, Playwright 1.64.1, and the worktree at `8151e6d` plus its uncommitted closure changes. The stack runs Next.js development servers, so these measurements are diagnostics and do not establish production-build acceptance.

The page profiler used five hard navigations per route in one browser context per route. Each route's first row is labeled `first-in-context`; remaining rows are reloads in that context. The output contains no cookies, headers, response bodies, or raw query values. API query strings are represented by truncated SHA-256 fingerprints.

## Results

| Route | Result | Assessment |
|---|---|---|
| `/en/dashboard` | First sample: CLS 0.3189, LCP 952 ms. Four reloads: CLS 0.0152 each; LCP 732–848 ms. All five reached the visible `main h1` selector with HTTP 200. | Cold CLS exceeds the 0.1 target. Warm reload CLS is below 0.1. Cold dashboard acceptance remains open and needs investigation. |
| `/en/dashboard/pos` | Five ready samples with HTTP 200. LCP: 908, 824, 808, 1056, and 888 ms; CLS: 0.0014–0.0021. | Observed LCP stays below 1.1 s and under the 2.5 s budget. This is a development-stack observation, not final production acceptance. |
| `/en/dashboard/orders` | All five reached `main h1` with HTTP 200. First navigation TTFB was 7048 ms during the dev server's cold route compilation; four reloads were 284–325 ms. Two identical `/api/v1/admin/orders` GETs appeared per sample; four duplicate requests were counted per page. | Warm measurements are diagnostic only. The repeated GETs may be development Strict Mode behavior; production-build request count remains unverified. Existing Admin Orders SQL benchmark evidence is separate. |

An authenticated first-navigation attribution capture in `dashboard-cls-attribution.json` recorded CLS 0.2001 across two events. The largest event (0.1321) included the Dashboard sidebar content container changing from y=65/height=536 to y=89/height=631 and summary content reflowing. The next event (0.0681) included the dashboard summary grid at the bottom of the viewport reflowing from y=643. This narrows the cold CLS issue to layout changes during Dashboard initialization; it does not establish which source change should be made, so the CLS gate remains open.

## Artifacts

- `dashboard-cold-page-samples.csv` and `dashboard-api-waterfall.csv`
- `dashboard-cls-attribution.json`
- `pos-warm-page-samples.csv` and `pos-api-waterfall.csv`
- `admin-orders-page-samples.csv` and `admin-orders-api-waterfall.csv`
- `dashboard-route.json`, `pos-route.json`, and `orders-route.json`

The separately documented Admin Orders PostgreSQL benchmark at commit `e30ffe3` reports page-size 50 dropping from 102 to 3 SQL statements, with median execution 291.768 to 11.784 ms and p95 504.784 to 23.033 ms; canonical response and payload checks matched. This session did not rerun that benchmark, so it remains prior-candidate evidence.
