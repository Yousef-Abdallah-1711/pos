# Distributed authentication controls

## Status

This note describes code in the `codex/auth-distributed-controls-20261117`
branch. It is not a deployment receipt. Production behavior must be confirmed
after the pull request is merged and released.

## One-time tenant selection proofs

`POST /api/v1/auth/login/confirm` validates the signed proof, confirms that the
selected tenant still has a matching account, then atomically claims the proof
identifier in Redis with `SET NX` and a TTL bounded by the proof expiration.
The Redis key contains a SHA-256 digest of the random proof identifier. The
stored value is a digest of the normalized identity and authorized tenant set.
The proof is claimed only after the identity and selected tenant checks pass.

Redis errors return a generic `503` and do not fall back to process-local
consumption. Production startup requires the Redis-backed session store. The
in-memory proof store remains available to isolated tests and non-production
construction paths.

## Distributed IP limits

The following current per-minute limits are preserved and now use an atomic
Redis counter with an expiration set on its first increment:

| Route scope | Limit |
| --- | ---: |
| Staff registration | 10 |
| Staff login | 20 |
| Tenant login confirmation | 20 |
| Session refresh | 20 |
| Forgot password | 5 |
| Subdomain availability | 60 |
| Super Admin login | 10 |

Redis keys contain a hash of the route scope and resolved client IP, not the
raw address. Store failures return a generic `503`; requests are not allowed
through without a limiter decision. Non-production uses a bounded 10,000-key
process-local store, which fails closed at capacity.

`GetClientIP` trusts forwarded headers only when the immediate peer is in
`TRUSTED_PROXY_CIDRS`. It walks `X-Forwarded-For` from the nearest hop toward
the client and returns the first untrusted address, or the leftmost address
when every forwarded hop is trusted. Configure only the exact
proxy CIDRs that overwrite or append the forwarded chain. Without a trusted
proxy entry, the immediate peer address is used and forwarded headers are
ignored.

## Metrics and operations

The internal-token-protected `GET /api/v1/health/metrics` response includes
fixed, low-cardinality `auth_security` counters for requests admitted by the
limiter, rate-limit rejections, limiter store errors, consumed selection
proofs, proof replay rejections, and proof store errors. It contains no user,
email, tenant, IP, token, or session labels. Scrape each API replica and sum
the counters in the existing monitoring system.

Alert on sustained increases in `rate_limit_store_errors` or
`selection_proof_store_errors`; these indicate Redis availability or
connectivity problems. The counters are process-local observations, while the
enforcement state is shared in Redis.

No database migration or cookie change is required. Public storefront and
order route limiters are unchanged by this slice.
