# Frontend and HTTP Caching Audit

## Current layers found in source

| Layer | Observed implementation | Scope / TTL / invalidation | Evidence confidence |
|---|---|---|---|
| Next server fetch cache | Restaurant Website homepage fetches declare revalidation values (homepage 300s, products 600s, theme 3600s in source); metadata also fetches theme | Server fetch URL/options; tenant-specific values depend on slug URL; actual Next cache-key behavior and invalidation not runtime verified | Static |
| Client query cache | Restaurant QueryProvider default stale 60s, cache/gc 5m; catalog hooks stale 5m/gc 30m; order list 1m, order detail 30s; tracking refetch 10/15s | Query keys vary by feature; catalog includes slug, order keys are ID-oriented; check tenant transition semantics | Static |
| Dashboard query cache | Theme hook uses `@tanstack/react-query`; no provider found in dashboard `src` search | Cache may not be active as intended; verify app wrapper before changes | Static, incomplete |
| Zustand/local storage | App stores include persisted client state; Super Admin auth token uses localStorage | User/session scope; inspect each store and migration behavior before cache decisions | Static |
| Redis | Optional `CacheService`, fallback to in-process map. Active consumers include tenant status (30s) and financial reports (5m) | Keys scoped by tenant; financial key also restaurant/path/raw query | Static + isolated Redis metadata |
| HTTP middleware cache | Cache middleware can assign API no-store/static immutable and ETag behavior, but caching stack is not wired in `main.go` | Not evidence of active headers | Static |
| Browser/CDN/service worker | Service worker files exist for Dashboard; global registration/production edge cache not verified | No cache-hit evidence | Static gap |

## Security and consistency notes

- Query key scope is a security boundary when data can change tenant/restaurant. Restaurant catalog keys include slug; order keys are ID-based. This is a candidate for verification, not a confirmed cross-tenant leak.
- Do not lengthen caching for stock, orders, user permissions, subscription access, or financial data without a freshness and isolation model.
- Financial report cache uses tenant+restaurant+path+raw query and a 5-minute TTL; no write invalidation was found in the handler. Query parameter order can create duplicate entries (inference from using raw query string).
- Redis errors on cache reads degrade to misses; Set errors are commonly ignored. On Redis outage, fallback cache is instance-local and may diverge across replicas.
- No ETag/304 or browser cache hit/miss was captured. Landing's measured page had zero `/api/` resource entries, which says nothing about server-side fetches.
