# Backend architecture audit

## Current structure

The executable API composition root is `backend/cmd/api/main.go`. It creates configuration, database connections, repositories, use cases, handlers, middleware and registers the HTTP routes. Current inventory: 363 registrations and 338 distinct method/path pairs. Other commands include `migrate`, `migrate-shard`, `seed-themes`, and `hash`.

| Layer / area | Observed implementation | Assessment |
|---|---|---|
| Entrypoint/composition | `cmd/api/main.go`, ~1,900 lines; inline constructors and route registration | Implemented, centralized, oversized; safe extraction by module is documented but not implemented. |
| Handlers | `internal/handler/http`, `internal/handlers`, plus legacy handler packages | Mixed generations; handlers range from use-case orchestration to direct repository/SQL work. |
| Use cases | `internal/usecase` includes auth, orders, payments, accounting, HR, KDS, themes | Real boundaries exist, but are bypassed for many routes. Existing architecture inventory estimates ~87 endpoints through use cases, ~123 direct handler→repository, 65 handlers accepting `*sql.DB`, and 40 inline SQL route closures. Treat counts as prior static inventory, not new execution metrics. |
| Repositories/domain | `internal/repository`, `internal/domain` | Broad module coverage; some repositories have ID-only or incompletely scoped methods. |
| Database | `database/sql` pool used across handlers/repos; tenant transaction helper exists but is not generally adopted | App-level scope is not universally a transaction/repository boundary; live RLS cannot safely be enabled. |
| Integrations | Paymob interface with live/mock clients, webhook retry worker; Redis cache and notification services | Implemented in part; payment success/concurrency/retry paths remain incompletely verified. |
| Real-time/background | KDS WebSocket hub, workers for subscriptions/webhook retry | Lifecycle and shutdown coverage is incomplete; functional E2E not demonstrated. |
| Experimental/dead | `internal/auth`, `internal/security`, `internal/scaling`, sharding code and `//go:build legacy` tree | Several are not wired to `cmd/api`; six experimental test families are explicitly quarantined. Do not count them as product functionality or passing tests. |

## Dependency direction and risks

The intended handler→use case→domain/repository flow exists in some modules, but the codebase is not uniformly layered. Inline SQL and direct `*sql.DB` access make it difficult to enforce transaction-scoped tenant context consistently. Multiple package generations (`handler`/`handlers`, `service`/`usecase`, legacy files) increase ambiguity and create a larger change surface. Existing records report about 31.6k of 96k non-test Go lines behind the `legacy` build tag; this is a static report figure, not a live metric.

## Reliability foundation

| Concern | Status | Evidence / gap |
|---|---|---|
| Panic recovery, timeouts, graceful shutdown | PARTIAL / NOT VERIFIED globally | Must be confirmed route/server-wide; no complete runtime lifecycle evidence in current reports. |
| Error mapping | PARTIAL | Existing architecture inventory reports many handlers returning raw `err.Error()`; common typed error mapper is a target, not implemented. |
| Request IDs and structured/redacted logs | PARTIAL | Logging exists; request ID, structured correlation and comprehensive sensitive-data redaction are not established end-to-end. |
| Context-aware DB/external calls | PARTIAL | Some calls accept contexts; no exhaustive audit proves every query/call has bounded timeout. |
| WebSocket and worker shutdown | NOT VERIFIED | No complete start/stop lifecycle or E2E evidence. |
| Build/static verification | VERIFIED (checkpoint evidence) | `go build ./...`, `go vet ./...`, `go test -count=1 ./...` passed, 680 tests. |

## Safe sequence (not implementation authorization)

1. Preserve route/security smoke as a baseline and add missing cross-tenant negatives.
2. Extract route registration into module files while retaining middleware order and machine-comparing all method/path pairs.
3. Move one module at a time to handler→use case→repository; select existing tested modules first.
4. Adopt tenant transaction context at use-case boundaries before any RLS rollout.
5. Establish typed domain errors, centralized HTTP mapping, request correlation and shutdown ownership.
6. Remove or rehabilitate legacy/unwired subsystems only after owner decisions and call-site proof.

Detailed historical current/target architecture context: `docs-organized/phase-2a/BACKEND_ARCHITECTURE.md`.
