# Phase 3C CI Safety Status

**Date:** 2026-11-15 (Africa/Cairo)
**Repository:** `Yousef-Abdallah-1711/pos`
**Latest verified `origin/main`:** `3ed336c5f00d25596da3186176256a629d0d474c`
**Safety implementation:** branch `codex/ci-safety-lockdown`, pending PR and required CI

## Remediation in this change

- API, UI, browser, and synthetic trial acceptance now require explicit CI targets on exact `.servio.test` hosts. Public hosts, unknown hosts, HTTP, credentials, query strings, fragments, and unsafe ports fail validation. Redirect following is disabled for API requests; browser routing aborts requests outside the allowlist.
- The GitHub Actions browser stack rewrites compiled frontend origins, backend origins, the seeded storefront domain, and Caddy virtual hosts to `.servio.test`. It binds its proxy only to loopback and checks the actual disposable PostgreSQL database, backend environment, mocked payment mode, Redis container, and internal-only Docker networks before registering synthetic users.
- Browser package installation now happens on the GitHub runner before synthetic fixture credentials are created. The browser container receives only synthetic credentials and a read-only dependency mount; its network has no external route.
- Automatic GHCR image publication, signed release publication on main pushes, and workflow-run staging verification were removed. The remaining staging workflow is manual and performs no staging request or write.
- CI artifacts include the exact commit SHA, run ID, isolated environment label, database name, and API target without credentials.

## Historical external-write incident

Existing CI evidence for run [43728374754](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43728374754) records that the acceptance script reached the public API target. The sanitized incident record reports two synthetic chart-of-account records created and a balanced 1.00 journal posted and then reversed. The retained summary does not identify the two account IDs or establish whether account cleanup occurred. The two account records therefore remain an **unresolved external cleanup/reconciliation item**. No additional request was made to the external environment, and this change intentionally does not delete, modify, or reverse those records.

## Verification in this worktree

- Node target tests: 4 passed, including rejection of the historical public API URL by browser routing.
- Python target tests: 3 passed.
- CI isolation guard tests: 4 passed, including a source scan that rejects public production API literals in every mutating acceptance entrypoint.
- Changed JavaScript syntax and Python compilation checks passed; `git diff --check` passed.
- The disposable Docker stack and full GitHub Actions pipeline have not yet run on this safety branch. CI acceptance remains pending; no green aggregate is claimed.

## Release and financial status

This remediation addresses the CI target boundary; it does not close the broader financial acceptance gaps. Worker interruption/replay, positive AR aging, parsed PDF monetary reconciliation, authenticated English/Arabic desktop/mobile financial browser acceptance, CFR-011 product/accounting disposition, and qualified accounting approvals remain unaccepted or policy-blocked as recorded in the Phase 3 matrix. Historical supplier, Balance Sheet, and financial CI results remain historical evidence and are not promoted to current acceptance by this change.

**Financial release remains NO-GO.** Independent review and full required CI are outstanding. No merge, deployment, or external cleanup was performed.
