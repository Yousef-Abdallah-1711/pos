# Runtime and Docker status

Observation time: 2026-11-02, read-only inspection. Docker also contains unrelated CRM, WebAudit, Supabase, school and other project services; these were not treated as POS components or modified.

## POOS services

| Container | Status | Health / URL |
|---|---|---|
| `pos-postgres` | Up | Healthy; host bind `127.0.0.1:15432` in `docker-compose.local.yml` |
| `pos-redis` | Up | Healthy; `127.0.0.1:16379` |
| `pos-backend` | Up | No Docker healthcheck configured; `GET http://localhost:8080/api/v1/health` returned 200 |
| `pos-landing` | Up | `http://localhost:13001` returned 200 |
| `pos-dashboard` | Up | `http://localhost:13002` returned 200 |
| `pos-restaurant-website` | Up | `http://localhost:13003` returned 200 |
| `pos-super-admin` | Up | `http://localhost:13004` returned 200 |

The four frontend ports are published through the backend service's network namespace. Normal browser route and API-flow behavior was not established by these GET checks. Container health is not configured on the Go API/apps; their `Up` state is not a functional health assertion.

## Runtime configuration and storage

The Compose project is `pos-dev`. Dev Postgres uses persistent named volume `postgres_data`; Redis uses `redis_data`; backend source is bind-mounted, and Go module/build caches are named volumes. Frontend source is bind-mounted with named `node_modules` and `.next` volumes. Local override binds ports to loopback, sets mock payment mode and RBAC enforcement. Base compose contains literal local-only DB/JWT defaults and configures API DB user `postgres`; do not reuse for production.

Live read-only DB evidence: PostgreSQL 15.19, 91 public base tables, 0 RLS-enabled tables, 108 migration rows. Active local DB/API connections were not enumerated by the first filtered query because the driver application name did not include `pos`; source and compose establish configured role, and role catalog showed `postgres` only among expected runtime/migrator roles.

No container was stopped/restarted, no volume altered, no migration invoked, and no service credentials included in this report.
