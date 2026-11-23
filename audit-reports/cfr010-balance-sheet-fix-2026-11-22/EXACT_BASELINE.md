# Exact Baseline

Date: 2026-11-22

- Repository: `Yousef-Abdallah-1711/pos`
- Remote `origin/main` and exact isolated-worktree base: `c9a57c1db92ce24acdd3ced0c55650a3ad031d90`
- Work branch: `codex/cfr010-bs-rootcause-20261122`
- Remediation worktree: `H:\Projects\pos-cfr010-bs-rootcause-20261122`
- The primary checkout at `H:\Projects\pos` had extensive unrelated changes and was left untouched.
- A fresh disposable PostgreSQL 16.4 database was migrated from the exact baseline. The task-local API, dashboard, Redis, and browser fixture use synthetic data only.

The user-supplied candidate string was not treated as a Git SHA. No production database, financial record, or deployment was changed.
