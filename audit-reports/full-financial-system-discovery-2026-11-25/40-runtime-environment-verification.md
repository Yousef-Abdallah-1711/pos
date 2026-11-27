# 40 — Runtime environment verification (executed 2026-11-25 22:07–23:07Z)

Baseline: HEAD `133260db66b1143c254d79476daa7650b41a33f0`, working tree unchanged (6 tooling-modified files + untracked audit files; 2 stashes intact). Pre-change capture: `evidence/runtime/00-pre-change-baseline.txt`, `01-dev-db-fingerprint.txt`.

## What was stopped / started
Stopped (not removed): `pos-landing`, `pos-dashboard`, `pos-restaurant-website`, `pos-super-admin`, `pos-backend`. Left running: `pos-postgres` (dev DB, untouched). Started project `pos-p2a-ws05-e2e` from `docker-compose.phase2a-e2e.yml` unmodified: postgres (tmpfs), redis, init-roles, api, 4 frontends. Ports API 18082, apps 13111–13114 (no conflict with dev). Disk was 93% used (18 GB free) — image pulls/build succeeded.

## Acceptance gates (evidence `04-env-gate-raw.txt`, `08`, `09`, `22`)
| Gate | Result | Evidence |
|---|---|---|
| Fresh isolated DB `pos_phase2a_e2e` | PASS | 145 migrations; separate container; tmpfs |
| App role not superuser / no BYPASSRLS | PASS | `pos_app|f|f` |
| Accounting pool separate | PASS (connection observed) | pg_stat_activity: `pos_accounting` 1, `pos_app` 2 (fallback not active) |
| RBAC enforced, not only logged | PASS | Employee-role user got 403 on 12 of 18 probes (09) |
| RLS enforcement | PARTIAL | 20 RLS tables forced; **57/70 tenant tables have no RLS**; pos_app without tenant context sees all tenants on those tables (F-019) |
| Tenant isolation at API level | PASS | owner A/B got 404 on all 9 cross-tenant reads of tenant C (08) |
| Background worker uses intended roles | PASS (inferred from pool use + successful posting) | outbox posted via accounting pool |
| Migrations successful | PASS | 145/145 |
| No production integrations reachable | NOT VERIFIED beyond using local compose only; no real payment credentials used | — |
| Production-equivalence | NOT CLAIMED | `ENV=phase2a-e2e` is not `production`: production-only guards (VER-020/021) not exercised |

Not executed: accounting-pool failure and DB interruption (tmpfs would lose data; not altered), staging/production env tests, stub-login reachability.
