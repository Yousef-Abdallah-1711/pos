# SERVIO Authentication End-to-End Acceptance Matrix

**Audit date:** 2026-11-17 UTC · **Current result:** No authenticated live browser tests were executed.
**Source SHA:** `18ac7a4d8fcb918ad14a15930e4718fc30fca0e9` · **Live staging SHA:** `96773ef22b818113327fb65ad044fdf7624c5656`.

## Status definitions and count

- **PASS:** actually executed against the named build and evidence stored.
- **FAIL:** executed and assertion failed.
- **BLOCKED:** cannot execute due a specific missing dependency/authorization/fixture.
- **NOT RUN:** defined, but not executed.
- Source tests existing in repository are labeled **SOURCE TEST ONLY**, not PASS.

**Execution count for this audit:** PASS 0 · FAIL 0 · BLOCKED 0 · NOT RUN 30. Source tests were read, not executed. Authenticated live journeys remain unverified because no approved synthetic account/session was provided through a secure mechanism. Public landing/API health endpoints were read-only checked and are documented in the architecture audit; they do not count as this matrix's authentication tests.

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

The E2E matrix remains **PASS 0 · FAIL 0 · BLOCKED 0 · NOT RUN 30** for browser-level scenarios. Selected unit/component tests were run for this branch and do not count as matrix executions.
