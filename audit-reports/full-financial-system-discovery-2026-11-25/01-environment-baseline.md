# 01 — Environment Baseline (Phase 0)

Captured: 2026-11-25 (UTC ~19:40). All values below were observed by command, not assumed.

## Repository identity
| Item | Observed |
|---|---|
| Working directory | `C:/Users/Yousef/Desktop/Projects/POSS` |
| Local branch | `master` (tracks the OBSOLETE remote) |
| Local HEAD (full SHA) | `133260db66b1143c254d79476daa7650b41a33f0` |
| HEAD commit | 2026-11-25 14:01:19 +0200, "Merge pull request #134 … codex/financial-recovery-review-batch2-20261125" |
| Correct upstream | `https://github.com/Yousef-Abdallah-1711/pos`, default branch `main` |
| Upstream `main` SHA (fetched into local ref `pos-real/main`) | `133260db66b1143c254d79476daa7650b41a33f0` — **identical to local HEAD** |
| `origin` remote | **Still `https://github.com/yousefabdallah171/POSS.git` (obsolete repo)** — NOT changed (no config changes without approval). `git pull` would pull the wrong repo. The correct repo is reachable only via the explicit URL / ref `pos-real/main`. |
| Worktrees | 1 (main working tree only) |
| Stashes | `stash@{0}` "local edits before reset to origin/master 2026-11-25", `stash@{1}` "pre-pull local changes 2026-11-25" (owner's earlier local work, preserved) |

### Working-tree changes at baseline (none are product-logic changes)
Modified (tooling rewrote on container start — `jsx: preserve→react-jsx`, `.next/dev/types` include, `next-env.d.ts`):
`frontend/apps/{dashboard,landing,restaurant-website,super-admin}/tsconfig.json`, `frontend/apps/{restaurant-website,super-admin}/next-env.d.ts`.
Untracked: `writing-block.md` (the audit prompt), `audit-reports/full-financial-system-discovery-2026-11-25/` (this package).

## Runtime (development stack, compose project `pos-dev`)
| Container | Image | Notes |
|---|---|---|
| pos-postgres | postgres:15-alpine (PostgreSQL 15.19) | db `pos_saas` |
| pos-backend | `pos-backend:dev` | Go, hot-reload (`air`), `ENV=development` |
| pos-landing / pos-dashboard / pos-restaurant-website / pos-super-admin | `pos-dev-*` | Next.js dev servers, ports 3001–3004 |

Containers were built from the working tree at `133260d` (source is bind-mounted, so running code == checked-out code; frontend/backend binary build hashes are not separately embedded — UNVERIFIED beyond that).

### Database
- Applied migrations: **145** rows in `schema_migrations`; **145** `.sql` files in `backend/migrations` (counts agree; highest numbered file 192; `999_theme_migration_legacy.sql` also present).
- Roles present: `postgres` (superuser, **bypassrls=true**, login), `pos_app` (NOLOGIN), `pos_migrator` (NOLOGIN), `pos_accounting` (NOLOGIN, no password).
  `pos_app` / `pos_migrator` / `pos_accounting` were **created manually during local startup** because the dev compose does not provision them but migrations 164+ and 182 require them (see environment-concern table).

### Backend DB identity (critical for security conclusions)
`DB_USER=postgres` (superuser, BYPASSRLS), `ENV=development`, `ACCOUNTING_DB_USER` unset ⇒ accounting pool falls back to the operational pool. `RBAC_ENFORCEMENT` unset ⇒ audit (non-enforcing) mode in development.

> **Consequence:** nothing observed on this dev stack can support a claim of production-grade tenant isolation, RLS enforcement, role separation or RBAC enforcement. Those gates (Phase 9, gates 16/17/22) must be tested on the restricted-role disposable stack below.

## Isolated test environments provided by the repository
| Compose file | Project | Isolation | Frontends | Restricted roles | Accounting pool | Redis |
|---|---|---|---|---|---|---|
| `docker-compose.phase2a-e2e.yml` | `pos-p2a-ws05-e2e` | tmpfs Postgres, loopback ports (API 18082, apps 13111–13114, PG 15437) | YES (all 4) | `pos_app`/`pos_migrator`/`pos_accounting` w/ passwords | YES (`ACCOUNTING_DB_USER`) | YES |
| `docker-compose.pos-test.yml` | `pos-test` | tmpfs Postgres, API 18080, PG 15433, mock payments | no | `pos_app`/`pos_migrator` | no | no |
| `docker-compose.security-closure.yml`, `docker-compose.test.yml` | — | tmpfs PG | partial | varies | — | — |

**Selected plan:** use `docker-compose.phase2a-e2e.yml` for all write/transaction/security/concurrency tests (production-equivalent roles, disposable data, owner's dev data untouched). Not started yet: the Docker VM has 7.7 GiB RAM and the dev stack already uses ~5 GiB, so the dev stack must be stopped first (reversible, volumes kept). **Awaiting owner approval** — see progress.md.

## Environment concerns from prior reports — verification status
| Lead | Verified now? | Evidence |
|---|---|---|
| Migrations through 192 | VERIFIED | 145 applied == 145 files; highest file `192_auth_lookup_boundary.sql`; startup log "All migrations completed successfully" |
| Missing roles at local startup (`pos_app`, `pos_migrator`, `pos_accounting`) | VERIFIED (and worked around manually) | `pg_roles` now shows all three NOLOGIN; failures were reproduced at migrations 164, 174, 182 before creation |
| Dev accounting connection fallback | VERIFIED | `ACCOUNTING_DB_USER` empty in `pos-backend` env; backend started and posted nothing yet (runtime behaviour NOT YET TESTED) |
| Separate accounting pool | CODE PRESENT (`backend/cmd/api/accounting_db.go`), runtime-separated only on the e2e stack — NOT YET TESTED |
| Existing QA accounts/products | VERIFIED | 5 users (super admin + 4 QA), 6 products, tenant 1 "Test Tenant" — **owner dev data, must not be modified** |
| Local frontend config files modified | VERIFIED | see working-tree changes above |
| Browser Origin requirement for auth | VERIFIED | `POST /auth/login` without `Origin` → "Forbidden - invalid request origin"; with `Origin: http://localhost:3001` → 200 |
| `origin` still obsolete | VERIFIED | `git remote -v` |

## Reproduction commands
```
git rev-parse HEAD pos-real/main
git remote -v
docker ps --format "{{.Names}} | {{.Image}} | {{.Status}}"
docker exec pos-postgres psql -U postgres pos_saas -At -c "select rolname,rolsuper,rolbypassrls,rolcanlogin from pg_roles"
docker exec pos-backend printenv | grep -E "^(ENV|DB_USER|ACCOUNTING_DB|RBAC_ENFORCEMENT)"
```
