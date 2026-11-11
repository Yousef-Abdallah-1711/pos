# Final handoff

- Base: `origin/main` SHA `a6a2c6700a5fe883256deb1e5a3893108bd9865c`.
- Branch: `codex/postmerge-p0p1-remediation-2026-11-11`.
- Implemented: unique landing keys and priority image loading, server-side dashboard settings redirect, narrow dashboard flex shrink fix, cookie-session KDS WebSocket eligibility with regression tests, configurable API host-port routing, server-only restaurant SSR API URL, fixture tenant-subdomain CORS allowlisting, and the Phase 2A seed phone correction.
- Verified: four production builds (restaurant website rebuilt after the latest changes), restaurant API config tests (5), restaurant typecheck/build, frontend suites (dashboard 655, landing 53, super-admin 16), Go tests/vet/build, 126 fresh PostgreSQL migrations, targeted migration-174 rollback/reapply, restricted-role RLS tests, 14 payment recovery/idempotency/concurrency tests, 4 webhook retry database tests, 27 repository tenant/parent-reference tests, Phase 2A Playwright (15 passed, 0 failed; CLS 0.047), and registration Playwright (12 passed, 0 failed). Compose validation and diff whitespace checks also pass.
- PR #14 is open against `main`; the two hosted smoke checks (`runtime-localhost-guardrail` and `url-domain-config-smoke`) passed. No independent approval is recorded.
- Blocked/open: full RLS/tenant table and HTTP matrix, live KDS WebSocket event delivery, production Docker images, full migration rollback, registration/session refresh/expiry flows, and Lighthouse/load-performance metrics. RLS is enabled and forced on only 15/92 public tables; no blanket policy was added.
- The isolated Compose stack uses tmpfs PostgreSQL and its own project/network. Existing persistent databases and unrelated Docker projects were not modified.
- Required repository action after push: wait for CI and obtain eligible independent review. Do not merge until required gates pass.
