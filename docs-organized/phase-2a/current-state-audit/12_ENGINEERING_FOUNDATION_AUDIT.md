# Engineering foundation audit

| Area | Current state | Status |
|---|---|---|
| Dev topology | Compose project `pos-dev`: PostgreSQL 15, Go API, Redis, four Next apps; frontends share backend network namespace | Implemented; live services responding. |
| DB/runtime separation | Dev compose API uses bootstrap `postgres`; isolated stack has restricted role | Critical configuration drift. |
| Payment mode | Local override selects mock; production config guards require explicit live prerequisites | Code/config exists; production deployment settings not observed. |
| Builds | Backend and all frontend builds previously passed without TS bypass; frontend dev stack currently active | Prior evidence; current production builds intentionally not rerun to avoid mounted `.next` interference. |
| Errors | Many handler paths still return raw error text; unified typed error mapping incomplete | Partial. |
| Logs/request correlation | Logging exists, some sensitive auth logs removed; no complete structured request ID/redaction proof | Partial. |
| Health/readiness | `/api/v1/health` returns 200; app URLs return 200; DB/Redis health checks pass | Live smoke verified, readiness semantics not exhaustively checked. |
| Observability | `internal/observability/perf_metrics.go`, health/metrics endpoints and error tracker exist | Partial; metrics auth/config and alerting not shown as deployed. |
| Migration lifecycle | Dev startup applies migrations; 108 records live; audit made no migration call | Implemented but drift/rollback strategy not verified. |
| Config safety | `.env.local` uncommitted; local override requires values; base dev compose has unsafe local literals; secrets not displayed | Mixed, production separation not fully verified. |
| Shutdown/workers | Background worker and WebSocket code exists | Lifecycle, readiness, drain behavior not verified. |

### Classification of remaining work

Phase 2A foundation: clear module ownership, stable repeatable checks, documented environment matrix and least-privilege strategy. Later hardening: TLS/database network restrictions, full metrics/error collection and alerting, load testing, backup/restore drills, production secret manager and graceful deployment strategy. Critical dev superuser issue should be contained before shared or production-like data is used.
