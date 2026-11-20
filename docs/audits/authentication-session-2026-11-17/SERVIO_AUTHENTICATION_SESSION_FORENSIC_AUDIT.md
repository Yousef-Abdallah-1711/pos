# SERVIO Authentication and Session Forensic Audit

**Audit date:** 2026-11-17 UTC · **Source SHA:** `18ac7a4d8fcb918ad14a15930e4718fc30fca0e9` · **Live staging SHA:** `96773ef22b818113327fb65ad044fdf7624c5656`.

Labels used in this report: **VERIFIED CURRENT BEHAVIOR**, **OBSERVED DEFECT**, **PROPOSED DESIGN**, **UNVERIFIED ASSUMPTION**, **REQUIRED BUSINESS DECISION**. Current behavior means checked-in source unless explicitly tagged “live observation.” This is an audit, not an implementation or live authentication test.

## 1. Executive result

**VERIFIED CURRENT BEHAVIOR:** Staff browser login now has a cookie-backed session path. The API returns user data without the JWT, stores the access JWT in Redis behind an opaque random session ID, and sets host-only `poos_session` with `HttpOnly`, `SameSite=Lax`, `Secure` in production, and a TTL equal to configured JWT expiry. API middleware resolves the cookie, validates the stored JWT, and places claims in request context. Bearer JWT remains accepted as a compatibility route for endpoints wrapped by `SessionAuthMiddleware`.

**VERIFIED CURRENT BEHAVIOR:** Staging config sets `JWT_EXPIRY=2h`, requires Redis in production mode, and falls back to process-local sessions only in non-production. Redis stores a SHA-256-derived key for opaque IDs, but the value is the signed JWT credential. Sessions are revocable one at a time by deleting the key. No PostgreSQL browser session/refresh-token table was found.

**OBSERVED DEFECT:** Cross-portal login works by having each app call the central API with credentials, but cross-portal *UI state synchronization* is absent. Landing Navbar never calls `/auth/session`, so it does not show profile state while Dashboard and Super Admin each do a bootstrap request. Logout revokes backend session but only clears the current app’s in-memory/persisted state; other tabs/apps discover revocation only on a later API request. Dashboard's provider clears legacy token storage and refetches session on each page load, but a failed bootstrap redirects rather than distinguishing transient API/Redis outage from an expired session.

**OBSERVED DEFECT:** The tenant-selection proof replay store is a process-local map in `AuthUseCase`; replicas would not share consumed-state. The existing exact-origin + SameSite defenses reduce browser CSRF risk, but a separately defined CSRF double-submit middleware exists and is not wired into the registered auth routes in `main.go`. Do not assume that middleware is active.

**UNVERIFIED ASSUMPTION:** The currently live backend image matches all checked-out source details; it does not. Release receipt reports an older SHA. Source code establishes intended current route behavior; a protected synthetic-account browser run against the running image is needed to confirm live behavior.

## 2. Reachable routes and handlers

`backend/cmd/api/main.go:531-540` registers the following main user auth routes:

| Method and route | Handler/middleware | Observed contract |
|---|---|---|
| `POST /api/v1/auth/register` | `RequireTrustedOrigin` + IP rate limit 10/min → `AuthHandler.Register` | Validates/creates tenant+restaurant+owner; establishes cookie session; returns user and `session_established`, no bearer token. If session setup fails after account creation it returns 201 with `session_established:false`. |
| `POST /api/v1/auth/login` | `RequireTrustedOrigin` + IP rate limit 20/min → `Login` | Password check; one tenant returns user and cookie. Multiple tenant rows return a selection token/list and do not establish final session until confirm. Generic invalid-credentials message. |
| `POST /api/v1/auth/login/confirm` | `RequireTrustedOrigin` + IP rate limit 20/min → `LoginConfirm` | Validates email/tenant-bound proof, consumes it once (currently process-local), sets session cookie. |
| `GET /api/v1/auth/session` | `SessionAuthMiddleware` → `Session` | Resolves cookie or legacy bearer; returns ID/email/tenant/restaurant/role, no-store outer API middleware. |
| `POST /api/v1/auth/refresh` | `RequireTrustedOrigin` + IP rate limit 20/min + session middleware | Reissues access JWT from claims, atomically rotates opaque session ID, sets replacement cookie. Old ID should be consumed once when Redis atomic rotation is available. |
| `POST /api/v1/auth/logout` | `Logout` | If a session cookie is present, requires trusted Origin, deletes session key, expires cookies. No-cookie logout returns 204 and cookie clear. Does not revoke all devices. |
| `POST /api/v1/auth/forgot-password` | IP rate limit 5/min | Handler is wired; end-to-end email/reset completion and generic anti-enumeration response were not verified. |

Separate Super Admin routes: login `POST /api/v1/super-admin/auth/login`; session `GET .../session`; logout `POST .../logout`; each protected route goes through super-admin JWT/session resolution, active-account check, IP whitelist middleware, and rate limit. Super Admin has a separate token validator and cookie name. See `backend/cmd/api/main.go:1867-1883`.

An additional older-looking implementation exists at `backend/internal/handlers/auth_api.go`, but `main.go` wires `backend/internal/handler/http.AuthHandler`; do not mistake the former for the registered browser handler without another call site.

## 3. Login and registration flow

**VERIFIED CURRENT BEHAVIOR:** Landing login uses Axios `withCredentials:true` against configured API. `AuthHandler.Login` validates, calls `AuthUseCase.Login`, stores a cookie-backed session and returns user only. Landing then changes `window.location` to dashboard `/auth/callback`; the callback requests `/auth/session` with `credentials: include`, sets `cookie-session` marker in dashboard Zustand store, and navigates to the dashboard. Login form is client-rendered. Registration uses the same API/session mechanism.

**VERIFIED CURRENT BEHAVIOR:** Password verification uses bcrypt. The repository checks normalized email and password; registration creates tenant, restaurant and owner transactionally. Users table has tenant-scoped `(tenant_id,email)` uniqueness; repository login scans all same-email tenant rows. Account state handling and lockout controls were not exhaustively verified; route-level IP limits exist but appear in-process rather than shared fleet-wide.

### Login sequence

```mermaid
sequenceDiagram
  actor U as User
  participant L as Landing browser app
  participant API as Go API
  participant PG as PostgreSQL
  participant Redis
  participant DB as Dashboard browser app
  U->>L: Enter email/password
  L->>API: POST /auth/login + Origin, credentials
  API->>PG: Find same-email tenant-scoped users
  API->>PG: Compare bcrypt hash
  alt one tenant
    API->>Redis: SET hashed opaque ID -> access JWT with TTL
    API-->>L: 200 user + Set-Cookie poos_session
  else multiple tenant rows
    API-->>L: 200 tenant choices + signed selection proof
    U->>L: Choose tenant
    L->>API: POST /auth/login/confirm with proof and tenant ID
    API->>PG: Re-read eligible user
    API->>Redis: Consume proof in process memory; SET session key
    API-->>L: 200 user + Set-Cookie poos_session
  end
  L->>DB: Full navigation to dashboard callback
  DB->>API: GET /auth/session with API-host cookie
  API->>Redis: Resolve session ID
  API->>API: Validate JWT claims
  API-->>DB: 200 session identity
```

## 4. Cookie attributes, validation and refresh

| Cookie | Created by | Domain | Path | HttpOnly | Secure | SameSite | Lifetime |
|---|---|---|---|---|---|---|---|
| `poos_session` | Staff auth handler | Omitted, therefore host-only API origin | `/` | yes | TLS or production mode | Lax | JWT expiry (`2h` in staging compose) |
| `poos_super_admin_session` | Super Admin auth handler | Omitted, host-only API origin | `/` | yes | TLS or production mode | Lax | session store TTL, initialized from staff JWT expiry at API startup |
| `poos_impersonation` | impersonation flow | Omitted | `/` | yes | secure in production | Lax | explicit impersonation TTL |
| `csrf_token` | Helper exists; no current registered call found | Omitted | `/` | no | production | Strict | 24h |

There is deliberately no `Domain=.rakmyat.online` on the session cookies. The UI documents make cross-origin credentialed requests to the API; browser policy allows this only where CORS returns the exact allowed origin and `Access-Control-Allow-Credentials: true`. Since the origins are HTTPS subdomains of the same registrable domain, `SameSite=Lax` remains same-site in normal deployments; exact CORS and Origin checks still matter.

**VERIFIED CURRENT BEHAVIOR:** Login rotates a prior valid cookie if present; refresh atomically consumes the old ID and sets a new one. The Redis adapter supports shared atomic rotation and tests exercise concurrent rotation. On a stale cookie, login creates a new random identifier instead of adopting the stale one. Logout revokes only the presented browser session ID. The user may have other live sessions.

**OBSERVED DEFECT:** API credential values in Redis are signed JWTs; storage of opaque IDs is hashed, not session credential values. Redis access controls, TLS, key eviction policy, backup/restore behavior, and in-production HA were not verified here. Losing non-persistent Redis sessions logs everyone out; Redis outage makes session auth unavailable in production (fail closed).

## 5. Tenant resolution and membership

**VERIFIED CURRENT BEHAVIOR:** `AuthUseCase.Login` treats same-email rows across tenants as memberships and lists a generic `Tenant <id>` label, then `LoginConfirm` validates a signed selection proof whose tenant list is embedded in claims. `TenantContextMiddleware` sources authenticated tenant from JWT claims, disregards mismatching `X-Tenant-ID`, and only accepts a client restaurant selector if the restaurant verifier confirms ownership. Anonymous storefront requests must provide valid positive tenant and restaurant headers and pass DB ownership verification.

**OBSERVED DEFECT:** `LoginConfirm` proof consumption is a process-local `sync.Map`-style mechanism, which does not coordinate replica-level replay. Current tenant switch is not provided; switching likely requires another login selection flow or a new session design. Stale JWT tenant/role claims can persist until refresh or expiration unless the relevant endpoint rechecks an account/role status.

**RECOMMENDED ACTION:** Move single-use login challenge consumption to Redis with atomic compare-and-set/consume keyed by `jti`; include normalized tenant display names; add a server-validated tenant-switch endpoint only after business policy confirms identity/membership semantics. Re-evaluate authorization from fresh membership/role version on sensitive calls or short session lifetime.

## 6. Logout, expiration and browser lifecycle

| Journey | Current behavior established by source | Not yet verified live |
|---|---|---|
| Logout | API deletes current opaque session key and expires the staff and impersonation cookies; Super Admin logout invokes its use-case logout, revokes its API session key, clears its own cookie. | Browser tabs/app state synchronization; actual screenshot/header changes; cross-portal revocation timing. |
| Expiration | API JWT/session TTL are configured from JWT expiry; staging value 2 hours. Missing/expired cookie resolves 401. | Client recovery after transient Redis outage; refresh scheduling behavior from deployed application; simultaneous real-browser refresh. |
| Refresh | Explicit POST route rotates ID exactly once via Redis atomic operation; returns 204. | Whether production UI actually invokes refresh routinely and whether 401 retry storms are coalesced. Search of current UI found no normal `/auth/refresh` invocation outside tests. |
| Browser restart | Cookie includes `Max-Age` and `Expires`, so it is persistent for server TTL; browser restart should retain it unless user/session policy or browser cleanup removes it. | Actual browser profile persistence settings and policy meaning of “remember me”; no explicit remember-me control found. |
| Password/account change | Credential claims are short-lived and session key persists until TTL/revocation. | Account-wide key enumeration/revocation and password-change invalidation were not found in this audit. |
| Multiple tabs | Cookie jar is shared within API origin; API server state is authoritative. | Other application stores are not automatically notified; no BroadcastChannel/local-storage auth event found. |

## 7. Authentication request security

**VERIFIED CURRENT BEHAVIOR:** `RequireTrustedOrigin` checks exact configured origins for register/login/login-confirm/refresh/Super Admin login. Cookie-authenticated unsafe methods check Origin in `SessionAuthMiddleware`. CORS returns exact allowed origin, credentials true and `Vary: Origin`; it rejects preflight from unconfigured origins. API response middleware applies no-store to API routes. Login errors use generic invalid email/password text. IP rate limiting exists for login and register.

**OBSERVED DEFECT:** Rate limit maps are process-local. Multiple replicas can multiply limits; map growth/cleanup and trustworthy client IP forwarding need review. Super Admin handler writes some underlying auth error strings to response and logs, which can reveal more diagnostic detail than staff login. The helper CSRF token middleware is not wired into `main.go`; current exact Origin enforcement is the active defense. A CSRF cookie without `Domain` is also API-host-only, and the browser UI cannot read it from its own origin; current dashboard `authFetch` attempts to read `csrf_token` on the dashboard origin, so it will not see an API-host cookie if one were set there.

**VERIFIED CURRENT BEHAVIOR:** Root Caddy UI proxies strip `Authorization` before forwarding to Next.js, reducing accidental forwarding of legacy credentials to app servers. API proxy forwards `Set-Cookie` with Secure suffix in the active hardened deployment.

## 8. Landing and cross-portal root cause

Landing `Navbar.tsx` has no auth store subscription, no API session request and no profile menu: it renders Login and trial CTA unconditionally. Its zustand auth store (`landing/src/stores/authStore.ts`) is not persisted. Landing login correctly receives an HttpOnly cookie from the API because its Axios client uses `withCredentials:true`, but the cookie belongs to `api.pos.rakmyat.online`, is HttpOnly and cannot be read by JavaScript on `pos.rakmyat.online`. The Dashboard performs a `/auth/session` call on its own bootstrap. Therefore the correct future fix is a deliberate server/API or small client-session bridge for landing, not widening cookie Domain.

**PROPOSED DESIGN:** Keep API-host-only session cookie; have landing call a no-store `/auth/session` with `credentials: include` in a bounded client island (or have a landing BFF securely forward cookies server-side). Render a neutral/loading header until the result arrives, never assume persisted UI data proves authentication, and keep public shell cacheable. Any navigation decision remains backend-authorized.

## 9. Root findings with actions

| Severity | Finding | Evidence | Recommended action | Verification |
|---|---|---|---|---|
| High | Landing lacks authenticated header/session bootstrap | `Navbar.tsx:50-87`; landing auth store is memory-only | Add bounded no-store session island or BFF in planned work; keep host-only API cookie | Browser login→dashboard→landing, DOM/header, network trace, cache isolation |
| High | Tenant-selection single-use proof is process-local | `auth_usecase.go:235-278` | Shared atomic proof store with `jti`, expiry and one-time consume | Multi-replica replay integration test |
| High | Super Admin remains a distinct privileged session but currently same API parent host and shared session store | Super Admin cookie and route middleware in handler/main | Preserve distinct host-only privileged cookie, harden MFA/step-up and policy | Privileged boundary matrix and credential isolation test |
| Medium | Legacy bearer and localStorage paths remain in client/source | dashboard `lib/api.ts`, storefront `lib/api-client.ts`, `authFetch.ts` | Inventory route-by-route; migrate/remove only with compatibility proof | Static scan + browser storage/cookie inspection + API integration |
| Medium | Logout state is local to each app; no cross-origin synchronization | no cross-app event channel; independent app stores | Backend-authoritative check on navigation/focus and optional same-origin BroadcastChannel | Multi-app/multi-tab revocation test |
| Medium | Client redirects on any dashboard session-check failure | `DashboardAuthProvider.tsx:55-69` | Separate 401 expiry from 5xx/network failure; preserve protected screen shell with retry | Inject 401, 503, timeout in browser |
| Medium | In-process IP limits do not coordinate replicas | `RateLimitByIP` implementation and main route wrapping | Shared edge or Redis token-bucket limiter, trusted proxy chain | Multi-replica burst test |
| Medium | Tenant-switch API absent; account rows are duplicated per tenant | login/usecase + `users` schema | Resolve identity/membership business rule before schema design | API and data model review |
| Medium | Canonical dashboard domain alias mismatch | compose/Caddy use `pos.app`; live `app.pos` TLS fails | Owner confirms canonical URL and redirects/retire aliases | DNS, TLS, redirect, CORS, email link checks |
| Low | Staging app check returns older signed receipt than source HEAD | live receipt above | Treat staging as old deployment; no deployment performed in audit | Release verification workflow |

## 10. Tests and evidence status

**VERIFIED FROM TEST SOURCE (not run in this audit):** `auth_http_session_integration_test.go` asserts no JWT in login JSON, host-only HttpOnly cookie attributes, trusted-origin enforcement, server revocation and 401 after logout. `auth_refresh_test.go` exercises concurrent refresh and expects exactly one success while old ID is consumed. Session store tests include shared Redis and expiry. Tenant middleware tests cover mismatch rejection/ignore/verification. These test definitions are valuable but do not equal live acceptance or prove all deployment gates.

**NOT RUN:** No synthetic credentials/session were supplied or present in an approved mechanism for a real authenticated live browser test. This report does not mark authenticated portal journeys as passed. The safe targeted Go test commands were attempted but the audit environment has no `go` executable (`go: command not found`), so local backend tests could not run. See the acceptance matrix.

**PROPOSED DESIGN:** Use one-time login proof state in shared Redis; keep opaque host-only cookie; make logout server-authoritative; add session version/revocation semantics for role/password changes; use bounded refresh coordination; expose a minimal session response with no cache; retain explicit role/tenant boundaries.

**REQUIRED BUSINESS DECISION:** Remember-me/persistent session duration, whether logout should revoke current device or all devices by default, MFA/step-up for Super Admin, and whether customer accounts may ever share staff identity.

See [security threat model](SERVIO_AUTH_SECURITY_THREAT_MODEL.md) and [acceptance matrix](SERVIO_AUTH_E2E_ACCEPTANCE_MATRIX.md).

## Implementation update — 2026-11-17

This report was reconciled against `origin/main` at `c052c817e114df4d060d8428e6e236f6011876f8`. The current implementation slice, P0/P1 traceability, test evidence, and outstanding work are recorded in [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md); the refreshed repository/CI baseline is in [BASELINE.md](BASELINE.md). Historical audit conclusions remain intact. This update does not mark unimplemented Redis replay/rate-limit work, the full E2E matrix, Super Admin/storefront acceptance, cache-boundary testing, or 10K measurements as complete.
