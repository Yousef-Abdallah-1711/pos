# Phase 3C CI Safety Status

**Date:** 2026-11-15 (Africa/Cairo)
**Repository:** `Yousef-Abdallah-1711/pos`
**Latest verified `origin/main`:** `3ed336c5f00d25596da3186176256a629d0d474c`
**Safety implementation:** branch `codex/ci-safety-lockdown`, PR #63 open; required CI is still pending

## Remediation in this change

- API, UI, browser, and synthetic trial acceptance now require explicit CI targets on exact `.servio.test` hosts. Public hosts, unknown hosts, HTTP, credentials, query strings, fragments, and unsafe ports fail validation. Redirect following is disabled for API requests; browser routing aborts requests outside the allowlist.
- The GitHub Actions browser stack rewrites compiled frontend origins, backend origins, the seeded storefront domain, and Caddy virtual hosts to `.servio.test`. It publishes no application ports to the runner host, seeds users through backend-local loopback, and checks the actual disposable PostgreSQL database, backend environment, mocked payment mode, Redis container, and internal-only Docker networks before browser testing.
- Caddy's API virtual host is explicitly rewritten to `api.servio.test`, matching the guarded browser/API origin. The generic public-to-test suffix rewrite had left it at `api.pos.servio.test`, causing TLS failures in API-backed browser scenarios.
- Upload-browser logout is bounded and browser shutdown is guaranteed even when logout fails, preventing a teardown exception from holding the whole smoke job open.
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
- Full required CI run [43736821949](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43736821949) passed backend, all frontend suites, security/config, and image builds, but failed during browser-stack setup because the runner host could not reach a loopback-published proxy port. The registration request was refused before reaching the API; no synthetic registration occurred in that run. The proxy port publication has since been removed, and full CI is pending a rerun. No green aggregate is claimed.
- Full required CI run [43737700405](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43737700405) passed backend, all frontend suites, security/config, image builds, and target/isolation guardrails. Anonymous page smoke passed at 1440, 768, and 390 pixels with no measured horizontal overflow. Authenticated tenant/admin login and storefront checks failed because of the Caddy API hostname mismatch; the upload smoke then errored on logout and did not close Chromium, causing the 60-minute job timeout. The run's exact-head artifact was downloaded and inspected. The API route and logout cleanup are fixed in the current worktree; focused tests pass, but the fix has not yet been verified by full CI.

## Release and financial status

This remediation addresses the CI target boundary; it does not close the broader financial acceptance gaps. Worker interruption/replay, positive AR aging, parsed PDF monetary reconciliation, authenticated English/Arabic desktop/mobile financial browser acceptance, CFR-011 product/accounting disposition, and qualified accounting approvals remain unaccepted or policy-blocked as recorded in the Phase 3 matrix. Historical supplier, Balance Sheet, and financial CI results remain historical evidence and are not promoted to current acceptance by this change.

**Financial release remains NO-GO.** Independent review and full required CI are outstanding. No merge, deployment, or external cleanup was performed.
