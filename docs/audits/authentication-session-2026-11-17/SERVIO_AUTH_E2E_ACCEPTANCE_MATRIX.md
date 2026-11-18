# SERVIO Authentication End-to-End Acceptance Matrix

## Current execution status — 2026-11-18

This section supersedes the historical 2026-11-17 statuses below. The matrix is intentionally conservative: a scenario is PASS only when its stated acceptance checks were exercised. Partial subchecks are described in the evidence column and do not make the full scenario PASS.

**Candidate source:** PR #82 commit `badd31e45b395f37d5ab34b37765b6dd32e02865`, based on prior head `bbf18aa9d6653d50732a90062e1677940556bc52`. **Environment:** isolated local disposable Docker project `servio-auth-final-20261118`, synthetic fixtures, credentials and cookie values excluded. **Browser:** earlier Chromium lifecycle evidence was Playwright 1.64.1 / Chromium 156.0.8078.80 on the prior candidate; Firefox was not run. The current callback/provider changes passed focused unit tests, but current-source browser execution was blocked because the frontend dev servers did not serve routes under host resource pressure. **Current result:** PASS 13 · FAIL 0 · BLOCKED 3 · NOT RUN 14. The newly passing rows below are real API/service checks, not Playwright route mocks.

### Executed commands and evidence

- HTTP login/session/logout and 12-way concurrent refresh: `go test ./internal/handler/http -run '^(TestAuthHTTPLoginCookieSessionAndLogoutLifecycle|TestAuthRefreshRotatesCookieSessionExactlyOnceConcurrently)$' -count=1 -p 1 -vet=off` — passed.
- Redis sharing, 20-way atomic rotation, TTL, fail-closed outage probe, and recovery probe: `go test ./internal/service -run '^TestAuthSessionStoreRedisIntegration$' -count=1 -p 1 -vet=off` — passed using two independent Redis clients.
- Registration contract and restricted runtime role: `go test ./internal/repository -run '^(TestCreateTenantWithUserRegistrationContract|TestTenantRegistrationSeedsRBACUnderRestrictedRuntimeRole)$' -count=1 -p 1 -vet=off` — passed.
- Auth/session/tenant middleware and token unit tests: `go test ./internal/usecase ./internal/middleware ./internal/pkg/jwt -run '^Test(AuthMiddleware|LoginConfirm|SessionAuthMiddleware|TenantContext|TenantStatus|SelectionToken)_' -count=1 -p 1 -vet=off` — selected tests passed.
- Real-browser journey: `./node_modules/.bin/playwright test --config playwright.phase2a.config.ts tests/phase2a-e2e/auth-lifecycle.spec.ts --output /evidence/auth-lifecycle-final-results --trace off` — **3 passed**. `networkidle` was not used because the dev server keeps HMR WebSockets open; the test waits for the real `/auth/session` bootstrap response. Trace capture was disabled to avoid retaining credentials/cookies.
- Local frontend verification on the updated spec: dashboard `typecheck:test`, `typecheck`, and `lint` passed. The lint command reported 0 errors and 179 warnings.
- Sanitized browser log: `/tmp/servio-auth-evidence/auth-lifecycle-final.log`. Screenshots: `/tmp/servio-auth-evidence/auth-lifecycle-screenshots/`. The dashboard screenshot contains loading skeletons in chart/card regions; visual QA is therefore partial. No HAR/request-count artifact was retained.
- Current-source sanitized API evidence committed under `evidence/`: `final-api-lifecycle.json`, `final-api-security.json`, and `final-api-admin-boundary.json`. The lifecycle sequence recorded 6 `GET /auth/session`, 2 `POST /auth/login`, and 1 `POST /auth/logout`: guest 401, invalid login 401, valid login 200, authenticated session 200, same-cookie session after API restart 200, Redis-outage session 503, first retry after Redis recovery 200, logout 204, and revoked-cookie replay 401. Session responses were `private, no-store`; the host-only cookie had `HttpOnly`, `SameSite=Lax`, and `Path=/`.
- The security sequence recorded 3 `GET /auth/session`, 2 `POST /auth/login`, 3 `POST /auth/refresh`, 1 `POST /auth/login/confirm`, and 1 `GET /products/2`: missing/untrusted Origin refresh requests returned 403; trusted refresh returned 204, the old cookie returned 401, and the new one returned 200. A Tenant A request for a Tenant B product with forged tenant/restaurant headers returned 404. Multi-tenant login returned the authorized choices without setting a session cookie; confirmation bound the session to the selected tenant. The Super Admin sequence recorded one staff login, one Super Admin login, two Super Admin session requests, and one staff session request; cross-portal cookie use returned 401 in both directions and the Super Admin session returned 200.
- Historical browser response logs contain 24 response-listener events (login 200 ×3, login 401 ×1, logout 204 ×2, session 200 ×11, session 401 ×7) and 41 navigation-listener events. These are **not valid browser request counts**: `signIn` accumulated response listeners and React StrictMode replayed callback effects in development. The login and guest 401s were expected; the protected Inventory 401 after logout was the expected revoked-session response. The listener accumulation and callback duplicate-session/navigation behavior are fixed in the current source with a module-level single-flight callback helper and component redirect guard. Focused unit tests passed. A new real-browser timeline could not be captured because the current local frontend servers did not serve routes; do not use the historical listener-event totals as request counts.

### Current 30-scenario matrix

| # | Scenario | Status | Evidence and exact remaining gap |
|---:|---|---|---|
| 1 | Guest opens Landing | PASS | Chromium guest load at `/en`; SSR HTML had no synthetic account identifiers. `landing-guest.png`. |
| 2 | Guest clicks Login | PASS | Clicked Login and reached canonical localized `/en/auth/login`; browser journey log. |
| 3 | Successful login | PASS | Real API returned 200; cookie metadata showed HttpOnly and SameSite=Lax; no bearer token persisted in browser auth state. HTTP lifecycle test also passed. Password/cookie values are not retained. |
| 4 | Dashboard opens | PASS | Login redirected to Dashboard; real session endpoint returned 200 and authenticated state rendered. `dashboard-authenticated.png`. |
| 5 | Return to Landing | PASS | Opened Landing in a second tab; session island showed the authenticated profile and authorized restaurant. `landing-authenticated-account-menu.png`. |
| 6 | Profile menu | PASS | Display name and authorized restaurant shown; phone, IDs, session and token absent. Chromium screenshot reviewed. |
| 7 | Dashboard link | PASS | Account menu Dashboard and settings hrefs matched configured canonical URLs. |
| 8 | Refresh preserves session | NOT RUN | API session rotation and session preservation were verified in real HTTP tests, but browser reload in both portals, zero-flash behavior, hydration, and exact single-bootstrap/request count were not completed on the current source. |
| 9 | Dashboard in a new tab | NOT RUN | A second Landing tab was tested; a second Dashboard tab was not. |
| 10 | Browser restart policy | BLOCKED | Persistence/idle/absolute lifetime policy is undecided. **Owner:** product/security policy owner (not assigned). |
| 11 | Logout revokes session | PASS | Logout initiated from Landing and Dashboard returned 204; cookie expired; replay of old cookie returned 401. Real Go lifecycle integration also passed. |
| 12 | Landing returns to guest | PASS | After Landing logout, guest actions appeared and session endpoint returned 401. `landing-after-logout.png`. |
| 13 | Dashboard rejects revoked session | PASS | After Dashboard logout, protected Dashboard route returned to Login and old-session replay returned 401. `landing-after-dashboard-logout.png`; real Go lifecycle assertion passed. |
| 14 | Expired session and transient API outage | NOT RUN | Redis TTL and current API 503/recovery behavior were verified, but expired-cookie behavior and browser no-loop/retry recovery were not tested on the current source. **Owner:** auth E2E/QA (not assigned). |
| 15 | Concurrent refresh | PASS | 12 concurrent HTTP handler refresh attempts produced one winner; separate Redis integration used 20 concurrent rotations across two Redis clients and confirmed one winner. This does not establish independent API-process behavior or a browser HAR/retry-storm count. |
| 16 | Multiple restaurant memberships | NOT RUN | No multi-membership synthetic fixture/browser assertion. **Owner:** auth fixture/QA (not assigned). |
| 17 | Tenant switching | BLOCKED | No tenant-switch endpoint/product contract or cache-invalidation semantics are defined. Multi-tenant login selection is tested, but it is not an authorized in-session tenant switch. **Owner:** auth product/backend; define endpoint, policy, and invalidation contract before QA can test it. |
| 18 | Cross-tenant denial | NOT RUN | A real API Tenant A request for a Tenant B product with forged `X-Tenant-ID` and `X-Restaurant-ID` returned 404. The response body was not retained for a target-record disclosure assertion, so the full zero-leak criterion is not marked PASS. Evidence: `evidence/final-api-security.json`. |
| 19 | Super Admin boundary | PASS | Real API staff and Super Admin cookies were separate; staff cookie could not read the Admin session (401), Admin cookie could not read the staff session (401), and Super Admin login/session succeeded (200). Evidence: `evidence/final-api-admin-boundary.json`. |
| 20 | Anonymous storefront ordering | NOT RUN | No customer order flow/capability isolation run. **Owner:** storefront QA (not assigned). |
| 21 | EN/AR navigation | NOT RUN | Current spec covers clicking the EN→AR locale link, but current-source browser execution was blocked. Earlier Arabic rendering/direction evidence does not verify locale-switch route/state preservation. |
| 22 | Light/dark behavior | NOT RUN | Current spec covers theme persistence, but was not run in a browser on the current source. Earlier screenshots do not establish persistence or contrast conformance. |
| 23 | Mobile profile menu | NOT RUN | Current spec covers responsive/mobile drawer behavior, but no current-source browser run or full set of 320/390/768/1024/1440 screenshots was possible. Historical checks covered only a subset and do not count as full acceptance. |
| 24 | Keyboard and screen reader | NOT RUN | Escape and focus restoration passed for menu/drawer. Tab sequence, arrow keys, screen-reader semantics/announcements, and reduced-motion audit were not run. |
| 25 | CSRF rejection | NOT RUN | Current real API checks rejected missing and untrusted Origin on refresh (403) and accepted trusted Origin (204); all unsafe routes and CSRF-token cases were not exhaustively tested. |
| 26 | Malicious redirect rejection | NOT RUN | No encoded, protocol-relative, alternate-host, userinfo, or malformed-locale redirect cases run. **Owner:** frontend security/QA (not assigned). |
| 27 | CDN cache isolation | NOT RUN | Guest SSR did not contain identity and auth responses were private/no-store in local checks; no CDN warm A/B cache-poisoning run. **Owner:** release/platform (not assigned). |
| 28 | Backend restart with shared Redis | PASS | Real API container was restarted in the isolated Compose project while Redis/PostgreSQL remained available; the same session cookie resolved to 200 afterward. Evidence: `evidence/final-api-lifecycle.json`. |
| 29 | Redis outage and recovery | NOT RUN | Real API session resolution returned 503/private-no-store while Redis was paused and 200 on the first retry after recovery. Browser retry UI/no-loop behavior was not exercised on current source. Evidence: `evidence/final-api-lifecycle.json`. |
| 30 | High concurrency/load | BLOCKED | No dedicated capacity-approved environment or k6 run. Host root volume was 48 GiB with 897 MiB free (99% used) after the isolated attempt; frontend servers failed to serve routes under pressure. No 10k-session/333/2,000-RPS claim. **Owner:** performance/release engineering; provide a dedicated approved environment with sufficient disk/CPU and test-data isolation. |

**Release interpretation:** functional login/logout across Landing and Dashboard is demonstrated on a disposable Chromium stack. This matrix does not establish full tenant isolation, Super Admin/storefront boundaries, outage behavior, cache safety, accessibility conformance, Firefox support, or capacity readiness. Do not use it as staging or production acceptance.

**Historical audit date:** 2026-11-17 UTC · At that time, no authenticated live browser tests had been executed.
**Source SHA:** `18ac7a4d8fcb918ad14a15930e4718fc30fca0e9` · **Live staging SHA:** `96773ef22b818113327fb65ad044fdf7624c5656`.

## Historical status definitions and count (2026-11-17)

- **PASS:** actually executed against the named build and evidence stored.
- **FAIL:** executed and assertion failed.
- **BLOCKED:** cannot execute due a specific missing dependency/authorization/fixture.
- **NOT RUN:** defined, but not executed.
- Source tests existing in repository are labeled **SOURCE TEST ONLY**, not PASS.

**Historical execution count:** PASS 0 · FAIL 0 · BLOCKED 0 · NOT RUN 30. Source tests were read, not executed. Authenticated live journeys were unverified at that point. Public landing/API health endpoints were read-only checked and are documented in the architecture audit; they did not count as this matrix's authentication tests.

**VERIFIED CURRENT BEHAVIOR:** Source tests define coverage for cookie login/logout, single-winner refresh and tenant-header enforcement. **OBSERVED DEFECT:** Landing profile/session acceptance cannot pass yet because the authenticated header island is absent from current source. **PROPOSED DESIGN:** Execute the matrix below with isolated synthetic accounts and evidence recording. **UNVERIFIED ASSUMPTION:** A safe synthetic identity provisioning path is available on a disposable staging dataset. **REQUIRED BUSINESS DECISION:** Define persistence lifetime, customer identity boundary, evidence retention and accepted load SLOs.

## Test environment record required for every run

Record exact deployed/source SHA and image digests; environment and canonical URLs; browser/version/OS; timestamp UTC; user fixture ID (synthetic only); role and allowed tenant set; locale/direction/theme; viewport; trace/video/screenshot path; request/response HAR with credentials redacted; backend log/request IDs; DB fixture assertions; cache state and headers. Never put passwords, access tokens, cookies, or raw order capabilities in reports.

## Matrix

| # | Scenario | Steps/assertions | Evidence | Current status |
|---:|---|---|---|---|
| 1 | Guest opens Landing | Load `/en` and `/ar`; guest header; no private identity in HTML; page cache may be public | Browser DOM, response headers, HAR | NOT RUN |
| 2 | Guest clicks Login | Correct canonical localized auth route and safe return path | URL trace, screenshot | NOT RUN |
| 3 | Successful login | Synthetic staff account; cookie attributes; response body has no JWT/bearer | HAR redacted, cookie metadata via browser context | NOT RUN |
| 4 | Dashboard opens | Login completes; dashboard callback receives cookie and authenticated `/auth/session` 200 | Trace, backend request IDs, dashboard screenshot | NOT RUN |
| 5 | Return to Landing | New tab opens Landing; session island shows the signed-in display after API check | DOM/screenshot and one session request | NOT RUN; session island is implemented in this branch, no authenticated browser evidence |
| 6 | Profile menu | Opens by pointer/keyboard, shows safe name/current authorized tenant, no sensitive fields | Screenshot + accessibility tree | NOT RUN; account disclosure is implemented in this branch, no browser accessibility evidence |
| 7 | Dashboard link | Uses canonical configured URL and correct locale, no open redirect | Redirect trace | NOT RUN |
| 8 | Refresh preserves session | Reload dashboard; one bootstrap check; no persistent login flash/hydration mismatch | Video/trace, request count | NOT RUN; bounded single-flight bootstrap is implemented, browser refresh behavior unverified |
| 9 | New tab preserves session | Open dashboard in new same browser context; API session validates | Cookie metadata and HAR | NOT RUN |
| 10 | Browser restart policy | Close/reopen persistent context; behavior matches approved idle/absolute lifetime | Browser storage state without cookie values; timestamps | NOT RUN; business policy undecided |
| 11 | Logout revokes session | POST logout; Redis key removed; cookie expired; 204 | Backend integration assertion and response metadata | NOT RUN live; **SOURCE TEST ONLY** test exists |
| 12 | Landing returns to guest state | After logout on Landing, local UI clears and API returns 401 | Screenshot, API trace | NOT RUN; server-confirmed sign-out is implemented, no browser lifecycle evidence |
| 13 | Dashboard rejects revoked session | Reuse old browser context after logout; protected API returns 401 and UI handles | HAR, server assertion | NOT RUN live; source session test asserts 401 |
| 14 | Expired session | Expire test TTL; 401 transitions to login; no loop; transient 503 is not treated as expiry | Playwright clock/test fixture and trace | NOT RUN; Dashboard status classification unit tests pass, full browser flow unverified |
| 15 | Concurrent refresh | 10+ simultaneous refresh requests; exactly one old ID rotation accepted; no retry storm | Redis-backed two-process integration test and HAR | NOT RUN live; **SOURCE TEST ONLY** `auth_refresh_test.go` expects one winner |
| 16 | Multiple restaurant memberships | Same synthetic principal sees only DB-authorized memberships and clear names | DB fixture + UI screenshot | NOT RUN; source tenant choices use placeholder `Tenant <id>` |
| 17 | Tenant switching | Switch only to authorized tenant; session tenant changes; stale caches clear | DB assertion, API trace, UI screenshot | NOT RUN; endpoint/business semantics absent |
| 18 | Cross-tenant denial | Tenant A requests tenant B row and forged headers; zero data leaks | API integration tests and DB/RLS evidence | NOT RUN end-to-end; candidate source adds tenant-bound session-profile assertions; middleware tests cover some header behavior |
| 19 | Super Admin boundary | Staff cookie/bearer cannot read Admin; admin session cannot authorize staff; unauthorized role denied | API logs, cookie names, browser contexts | NOT RUN |
| 20 | Anonymous storefront ordering | Browse/order without staff login; order capability scoped to created order only | Synthetic order fixture and trace | NOT RUN |
| 21 | EN/AR navigation | Locale switch retains route/state safely; Arabic RTL page and profile menu | Screenshots, DOM `lang/dir` | NOT RUN |
| 22 | Light/dark behavior | Theme persists by existing preference; contrast/focus pass in guest/authenticated state | Screenshots and contrast report | NOT RUN |
| 23 | Mobile profile menu | 320/390 px; touch target, no overflow, focus order and Escape | Screenshots, accessibility tree | NOT RUN; account actions are integrated in the existing drawer, no viewport/accessibility run |
| 24 | Keyboard/screen-reader | Tab, Enter, Space, arrows, Escape; menu labeling and status announced | Accessibility scan/manual AT report | NOT RUN |
| 25 | CSRF rejection | Cookie request with absent/disallowed Origin/token denied; allowed request accepted | API integration result | NOT RUN live; source tests exercise untrusted logout Origin |
| 26 | Malicious redirect rejection | `//attacker`, encoded URL, userinfo, alternate host/port and malformed locale rejected | Playwright URL assertions | NOT RUN |
| 27 | CDN cache isolation | Warm public cache as A then B; no user-specific API/HTML leak; auth endpoint no-store | Cache header/body comparison | NOT RUN; live public pages observed cacheable and noindex |
| 28 | Backend restart | With Redis shared, session works after API restart; without Redis, fail closed and recover | Container restart in isolated staging, logs | NOT RUN |
| 29 | Redis outage | Protected route 503/401 policy and user retry; never silently authenticates via local state | Fault-injection trace and alert | NOT RUN |
| 30 | High-concurrency validation | 10k cookie sessions, 333 and 2,000 protected API RPS envelopes, then burst/recovery; zero tenant leakage | Load report, histograms, Redis/DB metrics | NOT RUN; scenarios defined in scalability plan |

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
