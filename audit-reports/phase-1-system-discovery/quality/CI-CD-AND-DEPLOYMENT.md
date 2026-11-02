# CI/CD and Deployment (docker files are uncommitted; treated as current state)

## GitHub Actions (.github/workflows)
- runtime-localhost-guardrail.yml: runs scripts/check_runtime_localhost.py on PR and push to main/master.
- url-domain-config-smoke.yml: sets up Python 3.11, Go (from go.mod), pnpm 9, Node 20, `pnpm --dir frontend install --frozen-lockfile`, then scripts/smoke_check_url_domain_config.py.
- Not present: go build/vet/test, golangci-lint, vitest/jest, tsc, next build, docker build, security scan (govulncheck, pnpm audit), deploy job.

## scripts/
capture_reports_cache_timings.py, check_runtime_localhost.py, smoke_check_url_domain_config.py, verify_deployment_endpoints.py; plus 2 tracked __pycache__ .pyc files. backend/scripts: PowerShell local start scripts, benchmark.ts, baseline/optimized metrics json, SQL seeds.

## Docker
- backend/Dockerfile: golang:1.24-alpine builder, alpine:3.20 runtime. backend/Dockerfile.dev uses air hot reload. frontend/Dockerfile.prod: multi-stage, node:20-alpine, corepack pnpm@9.0.0, build arg APP. frontend/Dockerfile.dev node:20-alpine.
- docker-compose.yml (dev): postgres:15-alpine with POSTGRES_PASSWORD/DB_PASSWORD hard-coded `postgres` (lines 58, 87).
- docker-compose.prod.yml: DB_PASSWORD required (`:?set DB_PASSWORD`), healthchecks and restart: always on services, volumes declared. No Redis service, no reverse proxy/TLS service found, no backup service.
- Old Dockerfile.dashboard/Dockerfile.landing and per-app Dockerfiles are staged deleted.
- config.go defaults DB_SSLMODE to disable and DB password to postgres; production env only guards JWT_SECRET.

## Backup / restore / rollback evidence
- Rollback: backend/scripts/rollback_phase6_accounting_hr.sql (one phase only). No down migrations for other migrations (migrations/ holds forward .sql only; not exhaustively verified).
- No pg_dump/backup script, schedule or restore doc found in scripts/ or compose files (grep for backup in docker-compose.prod.yml returned nothing). docs-organized/PRODUCTION-DEPLOYMENT-GUIDE.md and 05-INFRASTRUCTURE/DEPLOYMENT were not read; may contain procedures (unverified).
- Migrations are applied via docker-entrypoint.sh / cmd/migrate; no migration checksum/lock verification reviewed.
