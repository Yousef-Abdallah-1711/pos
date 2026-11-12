# Phase 2B Production Build and Browser Baseline

## Environment

- Compose project: pos-p2b-perf-20261103. Four Next.js applications and the Go API were built/run in production mode. Dashboard production image build completed in 368.62 seconds, 2026-11-03 22:33:23–22:39:32 +02:00.
- Source fingerprint at Dashboard build start: A00BD5F7ED5E0B239291B6412A58439088098E31DC4C8E90787D0CD5F20259AC. This differs from initial checkpoint F3F900... because Phase 2A worktree files changed before the build. No application source was changed for this investigation; existing worktree edits remain preserved.
- Next.js resolved to 15.5.9. Dashboard shared first-load JS was 103 kB.
- Fixture totals: 3 tenants overall (2 synthetic Phase 2B tenants plus one bootstrap test tenant), 5 restaurants overall, 24 users, 6 categories, 62 products, 400 orders, 400 order items, 80 tables, 180 inventory movements. Small synthetic fixture; not representative production scale.

## Browser results

| App / route | Mode and sample | Result | Notes |
|---|---|---|---|
| Dashboard /en/dashboard | Production build; 5 hard reloads | FCP/LCP median 472 ms; meaningful content median 436 ms; 23 API requests; 73 resources | CLS 0.1059 each sample; a 3-column section disappeared around 647 ms; no long tasks |
| Anonymous Dashboard /en/dashboard | 3 runs, 8 s networkidle timeout | Redirected to Landing login and rendered quickly | Next RSC prefetch to /en/auth/forgot-password returned 404; not evidence of user-facing Dashboard hang |
| Dashboard static route sweep | 13 templates | All HTTP 200 | 5 samples each for Dashboard, POS, Orders, Inventory; one for other static routes |
| Dashboard dynamic routes | 6 fixture-resolved templates | All HTTP 200 | Order 397, product 1, category 1, table 2 detail/edit/QR |
| Restaurant storefront | /en, /en/menu, /en/cart, /en/checkout | Loaded | Menu showed 21 synthetic products and made 2 API requests; no mutation |
| Super Admin | Dashboard, tenants, tenant 2, users, payments, settings | HTTP 200 | Audit Logs returned HTTP 200 but crashed client-side |
| Landing | /en, login, registration | Loaded | Super-admin login showed Access Failed |

Hard reload Dashboard values use five samples; the median is reported. Per-route TTFB, DCL, load, INP, per-app resource bytes, cold startup, mobile/slow network, and complete frontend timing breakdowns were not captured consistently. These remain unmeasured. No raw credential-bearing HAR is included. See 02_AUTHENTICATED_PAGE_READINESS.csv for all 110 source templates.

## Source boundary

Final hash comparison against the recorded build input found two later worktree changes: backend/cmd/api/main.go and frontend/apps/landing/src/app/[locale]/auth/login/page.tsx. These are existing Phase 2A changes and are not mounted into the running containers. Runtime results describe the immutable images built at A00BD5F7...; static source statements are scoped to that build. Details: 14_RUNTIME_SOURCE_BOUNDARY.md and 13_SOURCE_STATE_AT_COMPLETION.json.
