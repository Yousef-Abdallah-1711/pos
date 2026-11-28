# 47 — Environment restoration report (2026-11-25 23:07Z)

Actions: stopped the 7 isolated containers (`docker stop`, no `down`, no `-v`, no prune); started `pos-backend`, then `pos-landing/dashboard/restaurant-website/super-admin`. Isolated project containers remain (Exited) with named volumes; **the isolated Postgres data was tmpfs and is gone**; it is preserved as `evidence/database/isolated-audit-final-2026-11-25.dump` (+ `isolated-audit-financial-rows.sql`).

Verification (`evidence/runtime/25-dev-restoration-check.txt`):
- Dev URLs: 3001/3002/3003 → 307, 3004 → 200, API 8080 → healthy; backend env unchanged (`ENV=development`, `DB_USER=postgres`).
- Dev DB counts identical to the pre-change fingerprint (tenants 2, users 6, restaurants 2, orders 0, journals 0, GL 0, products 6, COA 62, migrations 145); owner tenant 2 untouched.
- 18 `pos-dev_*` volumes present; `pos-postgres` never stopped.
- Git: HEAD `133260db…`, same 6 tooling-modified files + untracked audit files, 2 stashes intact; no remote/branch change.
- No failures in restoration. Left behind: stopped `pos-p2a-ws05-e2e-*` containers and `p2a_*` volumes (reusable; remove only with owner approval).

## Batch 2 restoration (2026-11-26 01:43Z)
Docker Desktop had been closed by a host reboot (02:42 local); I started it, which auto-started the dev containers; dev DB re-verified intact before and after. For Batch 2 the dev app containers were stopped again (`pos-postgres` never stopped); the isolated stack ran on the persistent volume `pos-p2a-ws05-e2e_p2a_audit_pgdata` (audit-only compose copy `tests/runtime/docker-compose.audit-persistent.yml`). At the end: final dump `evidence/database/isolated-audit-batch2-final-2026-11-26.dump` (tenants 4, orders 34, journals 54, GL 115); isolated containers stopped (volumes kept, including the persistent audit DB volume); dev containers restarted and verified (`evidence/runtime/45`): URLs 3001-3004 + API healthy, dev DB counts identical (tenants 2, users 6, orders 0, journals 0, COA 62, migrations 145), git HEAD and the 6 tooling-modified files, 2 stashes unchanged. Correction to earlier text: the super-admin app has 9 page routes (not 17).
