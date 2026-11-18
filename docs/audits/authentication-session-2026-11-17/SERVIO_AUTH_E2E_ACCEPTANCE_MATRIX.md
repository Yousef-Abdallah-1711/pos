# SERVIO Authentication End-to-End Acceptance Matrix

## Current execution status — 2026-11-18

This section supersedes the historical 2026-11-17 statuses below. A scenario is PASS only when its stated acceptance checks were actually exercised. Partial subchecks remain visible but do not count as a full PASS.

**Exact candidate:** PR #82 head `d11e8a8912b93201e209e57644ab64621b12da42`; [SERVIO verified CI run 44337495674](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44337495674). The pull-request workflow built merge ref `f71f223cc09a31e94974723f303856ab3f459e50` against base `c052c817e114df4d060d8428e6e236f6011876f8`. Related exact-head runs: financial acceptance [44337495661](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44337495661), URL/domain smoke [44337495641](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44337495641), and runtime localhost guardrail [44337495844](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44337495844). **Environment:** GitHub-hosted `ubuntu-24.04`, isolated disposable Compose project `servio-auth-gha-44337495674`, PostgreSQL and Redis, synthetic WS05 fixtures; the workflow verified ownership and removed only this run's stack. Landing and Dashboard were built in production mode. **Browser:** Playwright Chromium; Firefox and actual screen-reader testing were not run. Sanitized browser, storefront, accessibility, and load evidence is in artifact `auth-browser-acceptance-44337495674` on the linked run.

**Current result:** PASS 25 · FAIL 0 · BLOCKED 2 · NOT RUN 3. On exact head `d11e8a8`, backend, auth-browser-acceptance, all four frontend jobs, financial acceptance, URL/domain smoke, and runtime-localhost guardrail passed. `security-and-config` failed the all-severity production dependency audit because of the High picomatch finding; the required rollup failed, and application/release image jobs were skipped. No finding was suppressed or waived.

### Executed commands and evidence

- **Real browser authentication:** `auth-lifecycle.spec.ts` — 11/11 Playwright tests passed. This includes guest/login, invalid credentials, nine real-login malicious redirect probes plus three malformed-locale paths, cross-site write rejection, the full registered unsafe-route Origin matrix, forged tenant headers, secure login-time tenant selection, logout/revocation, Redis outage/recovery, expired Redis session TTL/browser recovery, and explicit no-guest-flash/no-hydration-mismatch checks.
- **CSRF route matrix:** the exact source route manifest contained 163 unsafe routes after documented exclusions: 150 staff and 13 Super Admin. Every route was exercised with both missing and hostile Origin, for 326 expected 403 responses. Trusted login/refresh/logout were checked separately. Staff and Super Admin session cookies each received 401 when presented to the other portal's session endpoint.
- **Anonymous ordering:** isolated storefront browser test passed 1/1; a guest added an item and completed an order through the public customer capability flow.
- **Locale/theme/accessibility:** dark and light Landing audits each completed 5,497 checks with zero failures. EN and AR/RTL were exercised across the tested responsive viewports. Axe checks, ARIA structure, browser arrow-key behavior, Escape/focus restoration, no horizontal overflow, Arabic `lang/dir`, and a 44 px sign-out target passed. Actual screen-reader announcements were not tested.
- **Small authenticated-session load:** PR #85 harness, 10 RPS for 30 seconds, max 10 VUs. It produced 301 requests at 10.0327 requests/second, all 301 session checks passing, zero failed requests, and zero dropped iterations. Latency: p50 1.167 ms, p95 1.736 ms, p99 2.103 ms, max 2.634 ms. These measurements on a disposable hosted runner do not establish high-concurrency capacity or an SLO.
- **Additional isolated API evidence:** `evidence/final-api-lifecycle.json` records same-cookie session success after API restart, fail-closed Redis outage, recovery, logout revocation, and rejected old-cookie replay. `evidence/final-api-security.json` records missing/untrusted Origin refresh rejection and authorized tenant selection. `evidence/final-api-admin-boundary.json` records staff/Super Admin cookie separation and cross-portal 401 denials. The source-level concurrent refresh tests recorded one winner among 12 handler attempts and among 20 Redis rotations using two Redis clients.
- **Cross-tenant check:** Browser response was 404 and its body contained none of the target product, tenant, or restaurant identifiers. The disposable database query using the owner tenant/restaurant scope returned zero matching target rows. This demonstrates response non-disclosure and app-level scope evidence; it is not a claim that the products table has RLS enabled.
- **Local host recovery:** Current root volume is 48 GiB total, 4.2 GiB available (92% used); inode use 17%. RAM is 1.9 GiB with about 827 MiB available; swap is 2 GiB with about 899 MiB used. Eight healthy staging containers remain on this host. No local browser stack was started and no local Docker resources, caches, worktrees, KDS files, or SQL fixtures were deleted. Browser/load execution ran on GitHub-hosted disposable infrastructure.
- **Dependency gate:** The exact-head CI production audit reported one High picomatch finding through `apps__dashboard > next-intl > @parcel/watcher > micromatch > picomatch`. PR #86 removes the production watcher path; its production audit passed but its required development audit still found the High advisory through the watcher and Vite toolchain. No supported patched release is published; the upstream releases page still lists 4.0.7 as latest. The advisory and upstream candidates remain open/unfixed: [GitHub Advisory](https://github.com/advisories/GHSA-fc5q-39xc-vjpf), [upstream releases](https://github.com/micromatch/picomatch/releases), [closed upstream PR #235](https://github.com/micromatch/picomatch/pull/235), [upstream PR #238](https://github.com/micromatch/picomatch/pull/238), [upstream PR #246](https://github.com/micromatch/picomatch/pull/246), [upstream PR #261](https://github.com/micromatch/picomatch/pull/261). A development-only exception proposal is prepared at [picomatch-development-exception-proposal-2026-11-18.md](../../security/picomatch-development-exception-proposal-2026-11-18.md); it is not approved and has no assigned independent approver.
- **Staging:** no merge or deployment occurred. The running staging app SHA remains `96773ef22b818113327fb65ad044fdf7624c5656`; it is not this candidate.

### Current 30-scenario matrix

| # | Scenario | Status | Evidence and exact remaining gap |
|---:|---|---|---|
| 1 | Guest opens Landing | PASS | Production Chromium run for this PR head's merge ref returned 200 at `/en`; server HTML contained no synthetic account identifiers. EN/AR pages also passed dark/light audit checks. |
| 2 | Guest clicks Login | PASS | Click reached canonical localized `/en/auth/login`; login form rendered. |
| 3 | Successful login | PASS | Live API login returned 200; invalid credentials returned 401. Session cookie was HttpOnly/SameSite=Lax and auth state persisted no bearer token. |
| 4 | Dashboard opens | PASS | Successful login reached the authenticated Dashboard and `GET /api/v1/auth/session` returned 200. |
| 5 | Return to Landing | PASS | Same browser context opened Landing in a second tab; the session island resolved the signed-in profile and authorized restaurant. |
| 6 | Profile menu | PASS | Account panel showed display name and authorized restaurant, omitted phone/IDs/session/token, and exposed canonical Dashboard/settings links. Enter, Escape, and focus restoration were exercised. |
| 7 | Dashboard link | PASS | Account-panel links matched configured canonical Dashboard URLs in EN and AR. |
| 8 | Refresh preserves session | PASS | Dashboard and Landing were reloaded after login; both retained the authenticated session. The guest test also confirmed exactly one initial session request before login. |
| 9 | Dashboard in a new tab | PASS | A second Dashboard tab resolved the same browser session and retained it after reload. |
| 10 | Browser restart policy | PASS | User approved session-only cookie, 30-minute server-enforced idle timeout, and eight-hour absolute lifetime. Cookie has no persistent expiry; backend session and credential lifetimes are bounded by the approved idle/absolute policy. |
| 11 | Logout revokes session | PASS | Landing and Dashboard logout returned 204, browser cookie state cleared, and replay of each old cookie returned 401. |
| 12 | Landing returns to guest state | PASS | After Landing logout, the mobile account drawer showed the guest Login action. |
| 13 | Dashboard rejects revoked session | PASS | After Dashboard logout, old-cookie session replay returned 401 and a protected Dashboard route returned to Login. |
| 14 | Expired session and transient API outage | PASS | Live Redis outage/recovery passed; a second test expired the synthetic Redis session key and verified 401 plus browser recovery to canonical localized Login. |
| 15 | Concurrent refresh | PASS | Source integration tests recorded one winner among 12 concurrent handler refresh attempts and among 20 Redis rotations across two clients. Browser retry-storm behavior and independent multi-process rotation are not established. |
| 16 | Multiple restaurant memberships | NOT RUN | Login-time multi-tenant selection passed, but no fixture tests multiple restaurant memberships within one tenant. The current schema stores one `restaurant_id` per `(tenant_id,email)` membership; same-tenant multi-restaurant support needs a product decision and schema/login design. |
| 17 | Tenant switching | BLOCKED | In-session tenant switching remains disabled. No product decision confirms it is out of scope for release 1; adding it would require a membership-check contract, atomic session rotation/revocation, and cache invalidation behavior. Secure login-time tenant selection passed. |
| 18 | Cross-tenant denial | PASS | Forged tenant/restaurant headers produced 404; response body contained no target product/tenant/restaurant identifiers; owner-scoped `pos_app` query found zero target rows. This is app-level scope evidence, not proof of database RLS. |
| 19 | Super Admin boundary | PASS | Isolated API evidence: staff cookie could not read Super Admin session and Super Admin cookie could not read staff session (401 each); Super Admin login/session succeeded (200). See `evidence/final-api-admin-boundary.json`. |
| 20 | Anonymous storefront ordering | PASS | The real anonymous browser flow completed an order through the public customer capability route; staff session capability was not needed. |
| 21 | EN/AR navigation | PASS | EN→AR locale switch preserved the account state; document reported `lang=ar`, `dir=rtl`. Dark/light audits covered both locales. |
| 22 | Light/dark behavior | PASS | Authenticated theme toggle persisted after reload; dark and light audit suites each completed 5,497 checks with zero failures. |
| 23 | Mobile profile menu | PASS | 320/390/768/1024/1440 px were exercised; no horizontal overflow, accessible drawer labels/expanded state, Escape close/focus return, and a ≥44 px sign-out target passed. |
| 24 | Keyboard and screen reader | NOT RUN | Browser arrow-key behavior, ARIA snapshots, axe, focus restoration, and Escape passed. Actual screen-reader announcements were not exercised, so the combined scenario remains incomplete. |
| 25 | CSRF rejection | PASS | All 163 registered unsafe routes (150 staff, 13 Super Admin) rejected both missing and hostile Origin (326/326 expected 403); trusted login, refresh, and logout flows passed. |
| 26 | Malicious redirect rejection | PASS | Nine real-login query vectors covered protocol-relative, encoded, userinfo, hostile-host, host-suffix, alternate-port, backslash, and JavaScript schemes; three malformed locale paths stayed on the Landing origin. All logins ended on the configured Dashboard origin with no malicious query. |
| 27 | CDN cache isolation | NOT RUN | API session responses are private/no-store, but no candidate staging CDN was available for a cross-user warm-cache A/B test. |
| 28 | Backend restart with shared Redis | PASS | Isolated API evidence restarted the API while Redis remained shared; the same session cookie resolved to 200 afterward. See `evidence/final-api-lifecycle.json`. |
| 29 | Redis outage and recovery | PASS | Browser run observed session 503 with private cache control and no Set-Cookie, kept the user signed in, and recovered to 200 on the first explicit retry after Redis health returned. |
| 30 | High-concurrency validation | BLOCKED | The 10 RPS smoke passed, but 10k-session/333/2,000-RPS targets still need an owner-defined capacity/SLO and separately authorized production-like environment. No 10K-user readiness claim is made. |

**Release interpretation:** the production-mode auth journey, approved session policy, cross-tenant denial, Redis failure/expiry/recovery, anonymous storefront order, full unsafe-route CSRF matrix, redirect matrix, EN/AR responsive behavior, and a small disposable load smoke have execution evidence. Release is blocked by the red production dependency audit, unresolved tenant-switching and high-load decisions, untested same-tenant multi-restaurant membership, incomplete screen-reader/CDN checks, and lack of independent approval. No merge or staging deployment was performed from this candidate.

**Current candidate SHA:** `d11e8a8912b93201e209e57644ab64621b12da42` · **Last reported staging SHA:** `96773ef22b818113327fb65ad044fdf7624c5656`.

## Existing source test coverage inspected

- `backend/internal/handler/http/auth_http_session_integration_test.go`: creates a synthetic account in a disposable PostgreSQL DB, checks no JWT in login response, cookie attributes, session resolve, untrusted Origin logout denial, server revocation and 401 afterward. The test skips without its explicitly configured isolated DB.
- `backend/internal/handler/http/auth_refresh_test.go`: checks 12 concurrent refresh attempts against an in-memory/test store; one succeeds and old session is consumed. Separate Redis integration exists in session store tests.
- `backend/internal/service/auth_session_test.go`: create/resolve/expire/rotate and Redis-sharing behavior.
- `backend/internal/middleware/tenant_context_test.go`: verifies tenant/restaurant header handling and fail-closed behavior.
- Frontend tests include dashboard auth redirect/bootstrap helpers and Super Admin provider behavior. They do not prove live cross-origin behavior.

These definitions are **SOURCE TEST ONLY** until run on the candidate SHA. They do not establish staging acceptance, screenshot evidence, cache isolation or capacity results.

## Playwright execution design

Use isolated per-test browser contexts for guest, owner, manager, cashier, chef/KDS, customer and Super Admin. Provision records only in disposable tenant fixtures. Use API fixture setup with generated credentials passed via CI secret environment; never print them. Assert browser cookie names/flags/domain/path while redacting cookie values. Add request listeners that record method, URL host/path, status, cache headers, and request ID, but redact `Cookie`, `Set-Cookie`, Authorization, password, selection proof and order-token headers. For logout/replay, retain old state only within a disposable test context and assert backend rejection.

Store screenshots with names like `sha-role-locale-theme-width-height-utc.png`, plus a JSON sidecar containing SHA, origin, role label, locale, theme, viewport, test ID and timestamp. Approved visual baseline update requires human review; test failure must not auto-update baseline.

## Release acceptance gate proposal

Block release on any tenant leak, Admin boundary failure, old session accepted after logout, untrusted Origin write accepted, auth endpoint shared-cacheable, secret/token in browser storage/response, open redirect, keyboard-inaccessible account menu, or failed auth/session dependencies. Require all applicable scenarios 1–30 executed or explicitly waived with a named owner and rationale. A scenario marked NOT RUN is not a pass.

**UNVERIFIED ASSUMPTION:** Synthetic staff/admin/customer accounts can be provisioned safely in the staging environment with permission sets representative of production. **REQUIRED BUSINESS DECISION:** privacy-safe screenshot retention, cookie lifetime policy, customer identity model and availability SLOs for load scenarios.

See [security threat model](SERVIO_AUTH_SECURITY_THREAT_MODEL.md) and [implementation roadmap](SERVIO_IMPLEMENTATION_ROADMAP.md).

## Implementation update — 2026-11-17

This report was reconciled against `origin/main` at `c052c817e114df4d060d8428e6e236f6011876f8`. The current implementation slice, P0/P1 traceability, test evidence, and outstanding work are recorded in [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md); the refreshed repository/CI baseline is in [BASELINE.md](BASELINE.md). Historical audit conclusions remain intact. This update does not mark unimplemented Redis replay/rate-limit work, the full E2E matrix, Super Admin/storefront acceptance, cache-boundary testing, or 10K measurements as complete.

At the time of this 2026-11-17 update, the E2E matrix was **PASS 0 · FAIL 0 · BLOCKED 0 · NOT RUN 30** for browser-level scenarios. Selected unit/component tests were run for that branch and did not count as matrix executions.

## Candidate browser subset — 2026-11-17

The Landing session-island UI subset was subsequently executed on candidate commit `b7986af` with Playwright 1.64.1 / Chromium against a disposable local Next dev server. All four tests passed. The session and logout endpoints were fulfilled by Playwright route mocks with synthetic profile/cookie fixtures; the Go API, Redis, PostgreSQL, Dashboard, CDN, and real cross-origin revocation were not part of this run. Treat these as focused mock-browser checks, not PASS results for the end-to-end rows above.

| Executed browser check | Result | Scope and limit |
|---|---|---|
| Session profile, shared HTML, account links, new tab, confirmed logout | PASS (mock API) | Asserts private profile values are absent from Landing's initial HTML, account details and authorized links render after session discovery, a second tab discovers the fixture session, and mocked 204 logout changes both Landing views and expires the test cookie. No backend revocation was exercised. |
| Temporary outage and retry | PASS (mock API) | A mocked 503 keeps the header in its temporary-unavailable state; retry after the route changes to 401 reaches guest actions. |
| Responsive account navigation | PASS (mock API) | Checks 320, 390, 768, 1024, and 1440 px, horizontal overflow, 44 px minimum sign-out target, Escape/focus restoration, light theme, Arabic locale and RTL. This is not a screen-reader, arrow-key, contrast, or WCAG audit. |
| Logout failure | PASS (mock API) | A mocked 503 displays an unconfirmed-logout message and keeps the test session cookie and authenticated state. |

Reviewed UI captures: [authenticated desktop](evidence/landing-authenticated-desktop.png) and [authenticated mobile at 320 px](evidence/landing-authenticated-mobile-320.png). Passing-run Playwright traces were retained locally after replacing the synthetic cookie value with `[REDACTED]`; they were not committed. At the time of this 2026-11-17 mock-browser subset, full scenarios 1–30 had NOT RUN status; see the current execution section at the top for the 2026-11-18 real-service results.
