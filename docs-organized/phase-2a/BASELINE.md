# Phase 2A Baseline
Recorded at session start (2026-11-01). Branch `master`, HEAD `899dce68f81e4a1f8efc3b77e5a6957b3e96c3ba` (same as the Phase 1 audit).

## Pre-existing uncommitted work (owner's; must be preserved)
Tracked modifications: `.gitignore`, `README.md`, `backend/docker-entrypoint.sh` (CRLF→LF), `docker-compose.test.yml`, `docker-compose.yml` (now the DEV stack), `frontend/apps/{dashboard,landing,restaurant-website,super-admin}/next.config.js` (standalone output / type-gate opt-out).
Staged deletions: `Dockerfile.dashboard`, `Dockerfile.landing`, `frontend/apps/dashboard/Dockerfile`, `frontend/apps/landing/Dockerfile`.
Untracked: `.env.prod.example`, `.gitattributes`, `backend/Dockerfile.dev`, `backend/docker-entrypoint.dev.sh`, `docker-compose.prod.yml`, `frontend/.dockerignore`, `frontend/Dockerfile.dev`, `frontend/Dockerfile.prod`, `audit-reports/` (Phase 1 output).
Raw `git status --short` snapshot: `.baseline-status.txt`. Phase 1 diff snapshot: `audit-reports/phase-1-system-discovery/evidence/baseline-pre-audit.diff`.

## Rules in force
No commits/pushes/branch switches/resets/`git clean`; no sub-agents; no production data; Docker volumes of the dev stack (`pos-dev_*`) are not deleted. Phase 2A edits are made on top of the above and are listed in REFACTORING_LOG.md / SECURITY_REMEDIATION.md so they can be separated from the owner's changes.
