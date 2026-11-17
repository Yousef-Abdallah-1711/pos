# Authentication/session implementation status

Date: 2026-11-17

Base: `c052c817e114df4d060d8428e6e236f6011876f8`

Worktree branch: `codex/auth-session-island-20261117`

This is an implementation update to the eight audit artifacts in this directory. It preserves their historical findings and records only the scope exercised by this branch. It is not a claim that all roadmap phases or acceptance scenarios are complete.

## Implemented in this review slice

- `GET /api/v1/auth/session` now returns only the browser display profile needed by Landing and Dashboard, with `private, no-store` behavior and safe generic 401/503 errors. The profile query is bound to the authenticated user, tenant, role, and restaurant context, and only returns an active account. Missing or changed identity context revokes and expires the presented session cookie.
- Landing uses a hydrated client-side session island with `credentials: include`, `cache: no-store`, a five-second bound, request coalescing, 30-second in-memory freshness, cancellation, and separate guest versus temporary-unavailable states. Only safe display fields are held in memory.
- Landing desktop account disclosure and the existing mobile drawer provide localized identity, authorized restaurant display, Dashboard/account links, server-confirmed sign-out, retry, and keyboard/focus handling. Public pages remain client-personalized only after hydration; private profile data is not passed through server-rendered page props.
- Dashboard bootstrap distinguishes 401, 403, and temporary HTTP/network/malformed-response failures. It bounds and coalesces discovery, retries without clearing auth on temporary failures, clears tenant status when identity is rejected/replaced, and revalidates on focus/visibility. API cookies are attached only to the configured API origin or the same-origin `/api/v1/` proxy; legacy Authorization headers are stripped from these browser requests.
- The Dashboard keeps its existing `cookie-session` store marker and does not adopt browser bearer-token persistence.

## Traceability for P0/P1 findings

| Finding or requirement | Status in this branch | Evidence |
| --- | --- | --- |
| Safe active-user session DTO; no credential leakage; private no-store | Implemented | Repository/use-case/handler changes and backend HTTP integration assertions |
| User/tenant isolation and stale role/restaurant session invalidation | Implemented for session discovery | Tenant-bound query, cross-tenant claim rejection, role-change revocation assertions |
| Landing authenticated profile without private shared HTML | Implemented | Hydrated in-memory client island; Landing unit/component tests |
| Landing 401 versus API outage; server-confirmed logout | Implemented | Session client and account component tests |
| Dashboard 401 versus 403/429/5xx/network behavior | Implemented in bootstrap logic | Status classifier unit tests; full browser bootstrap flow remains unverified |
| Cross-origin portal synchronization | Partial | Landing sign-out hint is same-origin only; Dashboard revalidates with API on focus/visibility. Cross-origin behavior is server-authoritative, but full browser lifecycle scenarios remain unverified |
| Distributed atomic tenant-proof replay protection | Not implemented | Requires a separate Redis-focused change and multi-instance test |
| Distributed login/IP rate limiting | Not implemented | Existing process-local limiter remains |
| Immediate revocation on every protected API request after disable/role/membership change | Partial | Session discovery checks the current active row and revokes stale role/restaurant context; protected API authorization paths were not comprehensively converted to live revocation checks |
| Super Admin isolation | Preserved; full matrix not run | No privileged cookie/validator changes in this branch |
| Storefront/customer capability boundary | Not changed | Requires separate caller inventory and acceptance coverage |
| CSRF/origin behavior | Preserved; existing API/session boundary not broadened | Existing strict trusted-Origin middleware remains; full CSRF matrix not run |
| 10K concurrency and CDN/cache isolation measurements | Not run | No load test was sent to staging or production |
| 30-scenario authenticated browser acceptance matrix | Not run | Individual unit/component suites do not count as matrix execution |

## Validation results

Run results are recorded as execution evidence, not inferred from test files:

- Landing session/account/navigation unit and component suites: **PASS**, 13 tests.
- Dashboard session status-classifier suite: **PASS**, 4 tests.
- Landing and Dashboard application typechecks: **PASS** (`tsc --noEmit -p tsconfig.json`).
- Landing and Dashboard test typechecks: **PASS** (`tsc --noEmit -p tsconfig.test.json`).
- Backend package tests: **PASS** for `internal/handler/http`, `internal/repository`, and `internal/usecase` under Go 1.27.4. The HTTP session lifecycle test skipped because `POS_TEST_DATABASE_URL` was not configured; its PostgreSQL assertions are not counted as executed.
- Landing and Dashboard full lint commands: **PASS** with non-failing repository warnings. Focused lint over changed source/tests: **PASS** with no source findings (ESLint reports the app's existing `.eslintignore` configuration warning).
- Full E2E matrix, visual viewport matrix, accessibility audit, CDN cache test, load test, and deployment validation: **NOT RUN**.

The initial Landing run found and led to a JSX test-harness parse fix; the corrected run passed the three selected files. The implementation branch has not yet been reviewed by GitHub CI.

## Product decisions left unchanged

- Existing session lifetime and cookie policy were not expanded; no indefinite session or parent-domain cookie was introduced.
- Global identity/membership migration, browser-restart policy, all-device logout, customer accounts, and Super Admin MFA provider/recovery policy remain product or security decisions outside this slice.
- No database migration was added.

See [baseline](BASELINE.md) for repository, PR, CI, and migration context.
