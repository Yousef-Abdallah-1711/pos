# Dashboard Navigation Timeout Root Cause

Three anonymous navigations to http://localhost:13202/en/dashboard used an 8-second Playwright waitUntil networkidle condition; all timed out. The browser redirected to Landing login and rendered quickly. This is a harness wait-condition failure, not a demonstrated visual Dashboard hang.

An RSC prefetch to /en/auth/forgot-password returned 404. The Landing login page links to this missing route. This is a reproducible functional issue and a source of continuing network activity relevant to networkidle; available evidence does not prove it alone accounts for the full timeout.

Classification:
- Functional failure: linked forgot-password route is unavailable (404).
- Harness timeout: networkidle failed 3/3 while the login UI rendered.
- Not established: authenticated Dashboard hang, redirect loop, tenant initialization failure, SSR stall, or API-latency cause.
- Authenticated Dashboard: rendered; five hard reloads showed median FCP/LCP 472 ms and meaningful content at 436 ms. No long tasks.
- Follow-up: wait for explicit login/dashboard content in automation and test background request settling separately; repair/verify forgot-password route.

Authenticated Dashboard home made 23 successful API requests (73 total resources): 8 permission checks, 2 each for orders, notifications, notification statistics, and tenant status, plus 7 other calls. Counts establish repeated request families, not whether requests are redundant or critical-path dependencies. CLS was 0.1059; a section disappeared around 647 ms. No credentials or raw HAR are retained.
