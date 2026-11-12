# Browser Network Waterfalls

## What was actually observed

- Isolated Chromium context `pos-performance-audit` on the verified WS05 E2E stack, unauthenticated, viewport 1397×1007.
- Landing `/en` first isolated-context navigation: responseStart 10,259 ms; DOMContentLoaded 11,279 ms; load 12,100 ms; 12 resource entries; 2,271,411 transferred bytes as reported by Resource Timing; no API pathname resource entries. Development mode and initial compilation/HMR are likely contributors; do not treat as production baseline.
- Repeat Landing `/en`: responseStart 2,277 ms; DOMContentLoaded 2,378 ms; load 3,075 ms; 12 resources; 2,222,676 bytes; no API resource entries.
- Dashboard direct navigation timed out at 60 seconds waiting for DOMContentLoaded. Afterward the browser was at Landing `/en/auth/login`. The dashboard runtime sample is BLOCKED and must not be treated as a 60-second page-load measurement.

No HAR, request bodies, response bodies, credentials, cookies, or tokens were captured. The Playwright MCP did not expose a trace/HAR export operation in the available interface; snapshots and console records are retained under `16_BROWSER_TRACES/`. They are not Playwright trace archives.

## Request count limitations

Resource Timing counts the Landing navigation document and resources observed in that page context; it is not a network proxy capture and can omit server-side API fetches, cross-origin timing detail, cached transfer sizes, and requests prior to navigation. Full API waterfall, duplicates, compressed-vs-decoded bytes, and endpoint payload contracts remain unmeasured.
