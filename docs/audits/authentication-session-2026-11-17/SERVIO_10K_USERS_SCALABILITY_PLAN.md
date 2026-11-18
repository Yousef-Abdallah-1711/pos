# SERVIO 10,000 User and Concurrency Scalability Plan

**Audit date:** 2026-11-17 UTC. This report defines explicit load envelopes; it does not claim the platform can meet them. Staging is a small integration host and is not evidence of production capacity.

## Labels and evidence

- **VERIFIED CURRENT BEHAVIOR:** code/runtime fact observed at the stated revision/time.
- **OBSERVED DEFECT:** current limit or behavior that prevents a defensible scale claim.
- **PROPOSED DESIGN:** measurement or architecture target, not implemented.
- **UNVERIFIED ASSUMPTION:** workload or production topology not established by telemetry.
- **REQUIRED BUSINESS DECISION:** owner/product input required for a valid service objective.

## Count users and workload separately

“10,000 users” is ambiguous. The following test cases are independent; do not add their concurrency values or treat one as a production profile.

| Scenario | Explicit load envelope | Derived starting request rate | What it measures |
|---|---|---:|---|
| Registered users | 10,000 active account rows, loaded across tenants with realistic role/membership distribution | Not a request rate | Indexes, auth lookup, tenant membership lookup, audit retention and batch maintenance. |
| 10k concurrent authenticated sessions, idle | 10,000 valid cookie sessions; 0 background polling | 0 RPS baseline, then navigation/focus cohorts | Redis session footprint, session check p50/p95/p99, cookie correctness, idle expiry. |
| 10k concurrent sessions, light interactive | 10,000 sessions, assume 1 protected API request per user every 30 seconds | 10,000 / 30 = **333 RPS** average; a 3× 30-second burst is **1,000 RPS** | API/Redis/database request path and tenant enforcement. Assumption must be replaced by telemetry. |
| 10k concurrent sessions, active | Assume 1 protected API request/user/5 seconds for POS-heavy period | 10,000 / 5 = **2,000 RPS**; 2× burst = **4,000 RPS** | High interactive/ordering envelope; not a claim of observed behavior. |
| Login burst | 10,000 logins over 10 minutes, then separately over 60 seconds | **16.7 RPS** / **166.7 RPS** average respectively; add 2× synchronized burst | Bcrypt/Argon2 CPU, DB lookup, rate-limit correctness, session insert and failure behavior. Current bcrypt cost and account distribution must be measured. |
| 10k concurrent landing visitors | 10,000 open sessions; assume 1 navigation per visitor/minute | **167 page requests/sec**; 10 embedded static assets on first view implies up to 1,667 edge asset requests/sec, mostly CDN | CDN cache hit ratio, origin request collapse, static asset throughput, JS/hydration impact. |
| 10k landing visits in 5 min | 10,000 document requests over 300 seconds | **33.3 document RPS** average; test 3× synchronized burst = 100 RPS | Warm and cold cache, route generation, locale/SEO, API session island load. |
| Authenticated header island | For 10k dashboard users opening one landing tab over 5 min | **33.3 `/auth/session` RPS** average, independent of marketing document cache | Cross-origin CORS, Redis lookup and minimal DTO serialization. |
| POS session validation | Assume every business API call resolves the cookie/session (as current middleware does) | Included in the protected API request rates above | Redis round trips per request; avoid multiplying “session-check requests” by “business requests” in a way that hides actual Redis load. |
| Order burst | **UNVERIFIED ASSUMPTION:** 100 orders/sec for 5 minutes, plus 2× short burst = 200 order writes/sec | Independent of user count | Transaction/write contention, idempotency, inventory/menu checks, payment/webhook behavior, queue and DB pool. Owner to supply event profile. |
| KDS live updates | **UNVERIFIED ASSUMPTION:** 1,000 simultaneously connected KDS/browser sockets; separately test 10,000 | Heartbeat/poll rate depends on implementation; measure source and network | WebSocket fan-out, reconnect storm, Redis/pubsub topology and per-tenant broadcast isolation. |
| Super Admin | 100 simultaneous admins with 1 protected request/10 seconds | 10 RPS | Privileged rate limiting, expensive aggregate queries, IP policy and step-up. Not a normal staff profile. |

These rates are test envelopes, not capacity estimates. For workload measurement, derive actual (RPS = concurrent actors × requests per actor ÷ interval), document think time and response bytes, and compare arrival-rate tests with closed-loop user tests.

## Current measured staging snapshot

At 2026-11-17 17:51 UTC on `servio-staging-01`:

- Root filesystem: 48 GiB reported by `df`, 33 GiB used, 16 GiB available (69% used); inode availability was not captured in the original readout and must be measured before a release.
- Memory: 1.9 GiB total, 996 MiB used, 909 MiB available, 124 MiB free; 2 GiB swap, 584 MiB used.
- Containers healthy: eight. Current `docker stats --no-stream` aggregate working set about 346 MiB (proxy 26, Super Admin 40, storefront 58, dashboard 75, landing 89, backend 20, Redis 4, PostgreSQL 28 MiB). This is a single idle/low-load sample, not peak data.
- Per-container configured caps in `ops/aws/compose.json`: PostgreSQL 256 MiB, Redis 96 MiB (48 MiB maxmemory), Go API 256 MiB, each Next app 224 MiB; proxy memory cap was not in that file segment and was not verified there. Frontend Node heap is 160 MiB.
- `docker image ls` reports approximate image sizes: backend 159 MB, landing 485–490 MB, dashboard 498 MB, admin 467 MB, storefront 471 MB; five-image displayed size sum about **2.08 GB** before transfer/layer reuse. Docker's displayed per-tag size includes shared layers, so not equal to incremental disk needed.
- Docker host retains multiple historical SHA tags (many releases); deployment pull/update peak space is current layers + new unique layers + writable layers + DB/WAL/log/backup growth. Snapshot had 16 GiB filesystem space; a precise candidate image layer delta was not calculated.
- Live staging release receipt still reported SHA `96773ef22b818113327fb65ad044fdf7624c5656`; repository branch HEAD is `18ac7a4d8fcb918ad14a15930e4718fc30fca0e9`. Do not infer deployed code from current working tree.

**VERIFIED CURRENT BEHAVIOR:** This host can serve its current eight-container staging workload within its configured limits at the sampled moment. **OBSERVED DEFECT:** It does not demonstrate 10k concurrent user capacity; memory headroom is under 1 GiB and swap is already used, and the sample is idle. No load test/latency SLO or sustained traffic profile was available.

## Landing performance and cache model

**VERIFIED CURRENT BEHAVIOR:** Live anonymous landing `/en` and `/ar` responses included `Cache-Control: public, max-age=0, s-maxage=300, stale-while-revalidate=60`; they had staging `X-Robots-Tag: noindex,nofollow`. Current Navbar is rendered with guest actions, so this HTML is not personalized. The Landing profile feature is not implemented.

**PROPOSED DESIGN:** Preserve CDN caching for public route shell, locale content, fonts and hashed static assets; add a small session island that calls no-store `/auth/session`. Measure both fully personalized SSR and cached shell + island:

| Metric | Personalized SSR | Cached shell + session island |
|---|---|---|
| Landing HTML cache hit | Usually private/bypass or personalized cache key required | Shared public cache remains available |
| Auth/session store work | At least one lookup on each uncached personalized render | One lookup per app visit/session check interval |
| First-paint profile state | Ready if SSR cookie forwarding/session succeeds | Neutral placeholder until small API request returns |
| Leakage risk | High if cache directives or cache key mishandled | Lower; personalized data isolated to no-store API response |
| Latency dependency | API/session service on HTML critical path | Profile slot only depends on session API; page content is independent |

Measure cache-hit ratios and TTFB/LCP/INP with guest and authenticated synthetic users; confirm authenticated API/session data never appears in a shared cache entry. Compare initial blank/neutral header duration and layout shift.

## Backend and Redis sizing method

**VERIFIED CURRENT BEHAVIOR:** Cookie-auth routes resolve a session through Redis and then validate JWT; API-level `NoStoreAPIResponses` covers API endpoints. In staging, Redis has 48 MiB maxmemory with `noeviction`; backend pool config shows max 15 open and 5 idle DB connections; PostgreSQL is configured at max 40 connections, 64 MiB shared buffers and 2 MiB work_mem. Auth route rate limits are in-process maps.

**OBSERVED DEFECT:** No metrics were found in the reviewed evidence for session lookup latency, Redis memory per session, API request rate by route, DB pool waits, SQL lock time, cache misses, login compute, WebSocket fan-out or queue depth. Local process rate limiting is not horizontally consistent.

**PROPOSED DESIGN:** Benchmark with candidate build and production-equivalent Redis/Postgres topology. For 10k live sessions, measure Redis memory per record and connection utilization under a mix of 333/2,000 API RPS test cases. DB pool concurrency must be sized from observed query latency, not one connection per request. Watch saturation and queue time; increase pools only if PostgreSQL max connections and memory budget permit. Use DB pooler if needed after measuring. Keep session validation in Redis; avoid a database query for every session check unless role/version policy requires it. Use indexes for normalized email, membership lookup and session maintenance. Run expiry cleanup with bounded batches and avoid full-table delete scans.

### Required metrics

- Auth: login attempts/success/failure by category, bcrypt duration, proof-consume conflict, session create/resolve/revoke/refresh latency, 401/403/429/5xx.
- Redis: command p50/p95/p99, memory, fragmentation, evictions, rejected writes, connections, failover, key TTL distribution.
- PostgreSQL: pool open/idle/wait and wait duration, query p95, active connections, lock waits, dead tuples, WAL/checkpoint, CPU/IO, RLS query plans.
- API: request rate/latency/bytes per route and tenant bucket; no email/session/token labels; rate limit and Go heap/GC.
- Web apps/CDN: cache hit/miss, origin RPS, TTFB, LCP, INP, CLS, JS bytes, hydration failures and API session island latency.
- POS/KDS/storefront: order submit p95/error, idempotency conflict, queue delay, WebSocket active count/fan-out and reconnect rate.

## Failure and deployment scenarios

Test each independently at production-like traffic:

1. Redis primary failover while sessions are checked: API fails closed, UI displays retryable service issue, no fallback local auth; recovery rate is bounded.
2. PostgreSQL pool saturation while Redis remains healthy: login/tenant lookup surfaces 503 without credential ambiguity; business writes remain isolated.
3. CDN cold cache and cache-key miss storm: static file origin is protected; origin request collapse prevents synchronized 10k reload stampede.
4. API/frontend cold starts: health/readiness gates prevent traffic before DB/Redis dependencies are ready.
5. Rolling deploy with old/new app versions: cookie/session format remains compatible; migration follows expand-migrate-contract; rollback binary can read the same session format.
6. Session expiration synchronized burst: clients single-flight refresh and randomize bounded rechecks; no thundering herd.
7. 10k login arrival burst: rate limit policy permits legitimate cohorts; bcrypt compute does not starve POS orders.
8. Order storm plus KDS fan-out: order writes and notifications do not share an unbounded queue or cross tenant.

## Capacity test gates (proposed)

No hardware recommendation is made before these tests. First capture an agreed latency/error SLO and production region/topology. Then run staged arrival rates: 30, 100, 333, 1,000, 2,000 and 4,000 protected API RPS; 33, 100 and 167 landing document RPS; and separate login, order and KDS scenarios. Hold each tier long enough to reach steady state, then burst and recovery. Require:

- p95/p99 API and session endpoint within owner-approved SLO and no rising queue/latency trend;
- errors, authentication false-denials and tenant isolation failures within explicit limits (zero cross-tenant leakage);
- Redis and DB remain below agreed saturation thresholds with failover/recovery verified;
- origin and CDN behavior verified for both guest and authenticated sessions;
- resource utilization remains bounded, with no OOM, swap storm or pool exhaustion;
- results include exact build SHA, data distribution, region, topology, generator, scripts, warm/cold conditions and raw histograms.

**UNVERIFIED ASSUMPTION:** The scenario mixes above approximate moderate restaurant SaaS behavior. Replace them with measured request traces from representative synthetic accounts before sizing production.

**REQUIRED BUSINESS DECISION:** Define latency/error SLOs, geographic footprint, peak order and KDS counts, registered-vs-concurrent customer definition, retention/availability target, acceptable degraded mode during Redis/database outage, and budget ceiling.

See [current architecture](SERVIO_CURRENT_ARCHITECTURE_AUDIT.md), [security model](SERVIO_AUTH_SECURITY_THREAT_MODEL.md), and [acceptance matrix](SERVIO_AUTH_E2E_ACCEPTANCE_MATRIX.md).

## Implementation update — 2026-11-17

This report was reconciled against `origin/main` at `c052c817e114df4d060d8428e6e236f6011876f8`. The current implementation slice, P0/P1 traceability, test evidence, and outstanding work are recorded in [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md); the refreshed repository/CI baseline is in [BASELINE.md](BASELINE.md). Historical audit conclusions remain intact. This update does not mark unimplemented Redis replay/rate-limit work, the full E2E matrix, Super Admin/storefront acceptance, cache-boundary testing, or 10K measurements as complete.

## Release-specific load-test specification — proposed, not approved

This is a test plan only. The 10 RPS / 10 VU smoke run does not validate these targets. Do not execute the high-load profile until the owner authorizes a separate production-like environment with capacity, monitoring, rollback, and synthetic data confirmed.

### Proposed workload

1. Create and verify **10,000 active, idle staff sessions** in Redis with the approved 30-minute idle and eight-hour absolute policy. Record serialized session bytes, key TTL distribution, memory overhead, and expiry behavior. Keep the sessions valid using a separately documented idle cohort so the test does not silently extend TTLs.
2. Ramp to **333 protected API requests/second for 30 minutes** across a representative tenant/restaurant distribution. Include the session lookup and ordinary read mix; keep mutation traffic separately labeled.
3. Generate a **2,000 RPS protected-request burst for 60 seconds**, then measure the full return to steady latency and error rates.
4. Run a login cohort of **100 simultaneous synthetic login attempts**, ramped over 30 seconds, with valid and invalid credentials tracked separately. Repeat only after verifying rate-limit behavior and monitoring password-hash CPU use.
5. During a separate recovery scenario, interrupt Redis using the production-like failover procedure; confirm staff authentication fails closed, then confirm recovery after Redis is healthy without accepting revoked or expired sessions.

### Proposed SLOs for owner approval

| Measure | Proposed target |
|---|---:|
| Protected request latency at 333 RPS | p50 ≤ 100 ms, p95 ≤ 300 ms, p99 ≤ 750 ms |
| Protected request latency during 2,000 RPS burst | p95 ≤ 500 ms, p99 ≤ 1,000 ms |
| Valid login latency | p95 ≤ 1,000 ms, p99 ≤ 2,000 ms |
| Protected-route timeout/5xx rate at steady load | < 0.1% |
| Protected-route timeout/5xx rate during burst | < 0.5% |
| False 401/403 for valid sessions | 0 |
| Cross-tenant or cross-restaurant data disclosures | 0 |
| Redis failover and stable recovery | ≤ 60 seconds, fail-closed throughout outage |

These are proposed thresholds, not accepted product SLOs. Product and engineering owners must approve them before a load run can determine release readiness.

### Required instrumentation and stop conditions

- Record p50/p95/p99, request count, timeouts, 4xx/5xx, valid-session false denials, and recovery timelines by route class.
- Record Redis used/peak memory, bytes per session, active clients, command latency, TTL distribution, rejected writes, and evictions. Proposed steady cap: ≤70% of configured Redis memory and zero evictions.
- Record PostgreSQL active/idle/waiting connections, pool wait time, query latency, lock waits, and CPU. Proposed steady cap: ≤70% of configured connection capacity, no pool exhaustion.
- Record service/host CPU and memory, OOM events, swap, container restarts, network, and disk/WAL growth. Proposed sustained CPU <70% and memory <80%; any OOM, unbounded growth, tenant leak, session false denial, or exhausted pool stops the run.
- Include exact source SHA, immutable image digests, dataset generation seed, region/topology, runner/generator capacity, warm/cold cache state, request mix, raw histograms, and a signed run receipt.
- Use isolated synthetic identities only. Do not point load generators at production or the shared staging stack. No infrastructure upgrades are authorized by this plan.

### Approval status

High-load execution environment: **not identified or authorized**. Proposed SLOs: **pending product and engineering approval**. 10K readiness: **not validated**. The existing 10 RPS smoke remains a small functional load check only.
