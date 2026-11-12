# Redis Architecture and Runtime Audit

## Source integration

- Single go-redis/v9 client constructor is in `backend/internal/service/cache.go:85–100`; `REDIS_URL` is optional. The client pings with a 2s timeout. Service exposes Get/Set/Delete and default TTL 5m (`:103–148`). Missing/unavailable Redis falls back to process-local memory during construction.
- Observed active call sites: tenant-status middleware key `tenant_access:{tenantID}` with 30s TTL (`middleware/tenant_status.go:16,56–96`); financial-report handler key `financial_reports:{tenant}:{restaurant}:{path}:{rawQuery}` and 5m TTL (`handler/http/financial_report_handler.go:207–256`).
- Separate theme cache is process-local, `theme:{id}`, 10m. No Redis pub/sub, queue, session, distributed lock or comprehensive rate-limit usage was found in source review.
- Runtime cache middleware exists but is not wired into the global chain.

## Safe isolated runtime observation

Verified exact project labels pointed to `H:\Projects\pos\docker-compose.phase2a-e2e.yml`, `com.pos.disposable=true`, WS05, `pos-p2a-ws05-e2e-20261103`; Redis was isolated with persistence disabled in Compose. Read-only Redis observation: Redis 7.4.9, zero keys, about 1.1 MB used, 114 hits/11 misses, zero evictions, 11 expirations, no LATENCY LATEST events, SLOWLOG length 1. Key values were not read. Empty keyspace and low fixture load mean this is infrastructure health metadata only, not cache effectiveness or production memory behavior.

## Risks / missing evidence

- Redis cache failures after client initialization need code-path review; callers may treat errors as misses and often ignore write failures.
- Process-local fallback permits per-instance inconsistency and stampede risk after restart/expiry.
- No sampled key classification or TTL distribution was produced because the isolated Redis DB had zero keys; no values were dumped.
- No cache hit/miss attribution by endpoint and no DB-versus-Redis request spans exist.
- Do not assume more cache is beneficial until tenant/restaurant/user scope, freshness, invalidation and query cost are measured.
