# Local performance measurement scripts

These scripts are for repeatable measurements against a disposable local stack loaded only with synthetic data. They do not mutate data. Confirm every configured route is a safe GET before running the API sampler; never point them at production or a shared development stack.

## API GET sampling

1. Copy `safe-get-routes.example.json` and add only verified read-only routes available in the disposable stack. It includes `GET /api/v1/health` as a concrete read-only example; it checks database health and should only be sampled against the disposable local stack.
2. Set `P2C_API_BASE_URL` to that stack's API origin. For authenticated endpoints, use a synthetic token and synthetic tenant/restaurant IDs only.
3. Run from the repository root:

```powershell
$env:P2C_API_BASE_URL = 'http://127.0.0.1:8080'
$env:P2C_BEARER_TOKEN = '<synthetic-token>'
$env:P2C_TENANT_ID = '<synthetic-tenant-id>'
$env:P2C_RESTAURANT_ID = '<synthetic-restaurant-id>'
$env:P2C_API_SAMPLES = '25'
$env:P2C_API_OUTPUT = 'performance-investigation-2026-11-03/phase2c/api-samples.csv'
node scripts/performance/measure-api.mjs --routes scripts/performance/safe-get-routes.example.json
```

The first request is a warm-up and excluded from latency percentiles. The CSV includes status counts, warm p50/p75/p95/p99, response byte count, and a best-effort record count. It does not contain bodies, headers, query strings, or credentials. Keep the token only in the process environment; do not put it in route JSON or output files.

## Authenticated browser page sampling

Copy `local-pages.example.json` and add routes that exist in your local app. Create a Playwright storage-state file from a synthetic account in the disposable stack and provide visible-content selectors that indicate the route is actually ready. The state file is sensitive even when synthetic; keep it outside version control and delete it during cleanup.

```powershell
$env:P2C_BASE_URL = 'http://127.0.0.1:3004'
$env:P2C_STORAGE_STATE = 'C:\path\outside\repo\synthetic-state.json'
$env:P2C_PAGE_SAMPLES = '5'
$env:P2C_PAGE_OUTPUT = 'performance-investigation-2026-11-03/phase2c/page-samples.csv'
node scripts/performance/profile-pages.mjs --routes scripts/performance/local-pages.example.json
```

`local-pages.example.json` is an array of `{ "path": "/local/path", "readySelector": "main h1" }`. The first hard navigation is labeled separately from repeat reloads in the same browser context. The CSV records browser navigation milestones, content readiness and its readiness status, layout shifts, requests, transferred resource bytes, long tasks, console errors, uncaught JavaScript errors, failed requests, and final path. The sanitized waterfall includes a truncated SHA-256 fingerprint of the query string so duplicate counts distinguish differently filtered requests without storing query values. A redirect or visible-content timeout is reported as data and must be classified; it is not automatically a page failure.

## Interpretation limits

- Run equivalent route sets, fixture sizes, browser versions, builds, and sample counts for comparisons.
- The API sampler measures end-to-end client latency only; it does not instrument SQL, Redis, handler, or serialization time.
- The page sampler measures browser-visible milestones; it does not isolate authentication, middleware, SSR/RSC, API, or hydration without trace-level follow-up.
- These scripts do not establish a performance budget. Record measured baselines and separately label proposed targets.
