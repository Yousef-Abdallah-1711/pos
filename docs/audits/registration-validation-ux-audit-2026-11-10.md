# Restaurant POS Registration Validation, Error Handling, Database Contract, and UX Audit

**Audit date:** 2026-11-10  
**Scope:** Read-only source/schema/test review plus read-only local browser inspection.  
**Checkout:** `codex/phase2a-continuation`; working tree was already substantially dirty. The registration page, auth handler/use case, and many other files have pre-existing local edits. Findings describe the current working tree unless explicitly called committed-baseline comparison. No dirty work was cleaned or overwritten.  
**Source changes / application fixes:** None. No database writes or registration submissions were made.  
**Evidence labels:** VERIFIED = directly observed in code/schema/runtime; INFERRED = supported but not fully reproduced; UNVERIFIED = unavailable evidence; PROPOSED = recommendation.

## 1. Executive Summary

**VERIFIED:** The registration form uses React Hook Form with Zod. Its visible local rules are substantially broader than the backend's rules, and neither layer checks phone length. The backend inserts the phone into `tenants.phone` first; migration 001 defines that column as `VARCHAR(20)`. The reported PostgreSQL exception is therefore most consistent with an over-20-character phone submitted in registration. The exact request value and the live database schema were unavailable, so runtime root-cause confirmation remains **INFERRED**, not VERIFIED.

**VERIFIED:** Registration creation is grouped in one SQL transaction covering tenant, owner user, restaurant, default theme, role/permission seeds, and role assignment. Any insert/seed/commit failure before commit rolls back these writes. However, token generation and session creation occur after the transaction commits; failures there can leave a successfully created account while the browser receives an error.

**VERIFIED:** The form has a 500 ms subdomain debounce, but no request cancellation or stale-response guard. The API's availability result is advisory; the tenant slug unique constraint and insert are the final arbiter. The UI currently requires an `available` state before submit. The current working-tree form sends `restaurant_slug`; the API and repository support server-side slug validation and conflict mapping.

**VERIFIED:** API errors are a flat `{ "error": string }`. The handler returns arbitrary use-case/database error text, usually as HTTP 400, so SQL details can reach clients. The frontend displays that raw field in a toast. No field-level mapping or request ID is present in this contract.

**UNVERIFIED:** Current live browser inspection cannot reproduce the historical database error or reported JavaScript `toLowerCase` exception because no registration was submitted and those errors were not present in this browser session. No persistent test data was created.

## 2. Verified Root Cause

### Assessment

- **Exact failing column:** **INFERRED: `tenants.phone`**. Migration `backend/migrations/001_create_tenants.sql:10` declares `phone VARCHAR(20)`. `backend/internal/repository/auth_repo.go:61-65` inserts `req.Phone` into `tenants.phone` in the first tenant insert. The reported message context says “failed to create tenant,” consistent with failure at this statement. The same phone is later inserted into `users.phone` (`003_create_users.sql:11`, repository lines 93-99), also `VARCHAR(20)`, but that insert is reached only if tenant creation succeeds.
- **Exact failing value category:** **INFERRED: user-supplied phone string**. It is the variable written to the relevant first insert. No request body was captured in the available evidence.
- **Maximum permitted length:** **VERIFIED in migration: 20 characters** for tenant phone; also 20 for user phone. Actual runtime schema not queried.
- **Actual submitted/generated length:** **UNVERIFIED**. The original safe/redacted payload and DB error detail are not available.
- **Root cause:** **INFERRED** unvalidated phone input exceeds the migration-defined 20-character limit. `VARCHAR(20)` is not implicated by the displayed subdomain limit of 63; tenant slug/domain/name are `VARCHAR(255)` in the initial migration.
- **Runtime reproduction:** Not attempted. The local app was reachable, but reproducing registration writes shared development data and the task forbids creating persistent test tenants. No verified disposable DB was established.

### Evidence and correction

The backend validates required values using exact empty-string comparisons and enforces password minimum length 6. It does not trim or limit phone. Phone is passed unchanged to the tenant insert. This diagnosis should be confirmed against a safe, redacted failing request and the actual target schema before choosing a migration. The minimum correction is to define and enforce a product-approved phone contract consistently: normalize/trim, allow empty as the optional representation, validate maximum length before SQL, and map invalid phone to a field error. If valid international phone representations exceed 20 characters, use a deliberate schema migration for both tenant and user columns rather than broad, unrelated length increases. Do not silently truncate.

## 3. Current Registration Architecture

| Stage | Current implementation / evidence | Finding |
|---|---|---|
| Form state / validation | `frontend/apps/landing/src/app/[locale]/auth/register/page.tsx:3-61` | React Hook Form + Zod resolver; no explicit `mode`/`reValidateMode` (library defaults apply). |
| Availability | Same file:65-106; API handler `backend/internal/handler/http/auth_handler.go:167-180` | 500 ms timeout; GET `/api/v1/auth/check-subdomain`; query slug lowercased in browser. |
| Submit | Register page:108-145 | Blocks unless current subdomain state equals `available`; posts name, email, password, restaurant name, slug, optional phone. Double submit is guarded through `isLoading` button disable, without handler-level early guard. |
| Route/middleware | `backend/cmd/api/main.go:483-492` | POST registration is wrapped in trusted-origin validation and IP rate limit of 10; availability endpoint has no rate-limit wrapper in this route block. |
| Binding/error | `backend/internal/handler/http/auth_handler.go:27-47,208-217` | JSON decode; no explicit body-size limit or DTO validation; errors returned as flat string. Conflict detection uses string matching. |
| Business checks | `backend/internal/usecase/auth_usecase.go:74-127` | Required fields, password minimum, subdomain format/availability, and email existence checks. |
| Persistence | `backend/internal/repository/auth_repo.go:31-215` | `Begin()` transaction; tenant → bcrypt hash → owner user → restaurant → theme → seed roles/permissions → assign ADMIN role → commit. |
| Auth/session | `backend/internal/usecase/auth_usecase.go:129-141`; handler:43-47 and current dirty session code | JWT created after repository commit, then session established by handler; response is `{user}` in current working tree. |
| Client success | Register page:126-131 (current working tree) | Stores user with empty token and redirects to dashboard locale auth callback. Exact end-to-end session/callback completion was not exercised. |

**Atomicity:** The multi-table DB creation is atomic for failures before successful `Commit`; deferred rollback covers early return. No explicit rollback is needed after commit. Post-commit JWT or session failure is outside that transaction and can leave a valid tenant/user without successful browser onboarding. The transaction uses `Begin()` rather than `BeginTx` with request context. No idempotency key exists. Concurrent slug races are ultimately constrained by tenant slug uniqueness and mapped to conflict, but the preliminary check is race-prone by design. Email duplicate precheck is not backed by a global unique constraint; concurrent duplicate emails across tenants may both pass because users are unique only on `(tenant_id,email)`.

## 4. Current Live Validation

The frontend uses `react-hook-form`, `@hookform/resolvers/zod`, and `zod`; `useForm` sets only `resolver`. Thus errors are surfaced on submit/default revalidation behavior, not a configured blur-first contract. The subdomain availability check runs during changes after a 500 ms timeout. Local browser inspection showed the form with no initial errors, consistent with untouched/default state.

The availability implementation clears the prior timer, but once a request starts it cannot be aborted and has no request sequence check. An older response can overwrite the status/message for newer text. Local network failure is represented as `unavailable`, indistinguishable in its status enum from taken/invalid values. Submit sends the currently entered slug in the current tree and performs a stale-able status check; backend uniqueness remains authoritative.

The intended neutral → blur validation → revalidate corrected fields → submit validation flow is a **PROPOSED** UX. For this small form, retain the existing RHF/Zod stack; explicitly configure mode/revalidation, connect errors to inputs with IDs/ARIA, focus the first invalid field, and use an abortable/sequence-safe async availability check. Keep server validation independent and authoritative.

## 5. Field Validation Matrix

| Field | Frontend (current) | Backend (current) | Database (migration evidence) | Gaps / assessment |
|---|---|---|---|---|
| Full name | Zod string min 2; no max/trim; input type text | Required iff exact empty string; no length/character validation | `users.name VARCHAR(255) NOT NULL` | Whitespace-only passes backend; client/backend minimum differ; Unicode is accepted by storage but not intentionally normalized; no max contract. |
| Email | Zod `.email()`; no trim/case normalization/max | Required iff exact empty; duplicate precheck; no format/max/case normalization | `users.email VARCHAR(255) NOT NULL`; unique `(tenant_id,email)` only | Cross-tenant duplicate is allowed by DB though app precheck treats it globally duplicate; race can bypass check. Case behavior depends on DB collation/query and is not normalized. |
| Restaurant name | Zod min 2; no max/trim | Required iff exact empty; no length validation | `tenants.name VARCHAR(255) NOT NULL UNIQUE`; `restaurants.name VARCHAR(255) NOT NULL` | Whitespace-only passes; max absent; tenant global unique name is imposed by DB even though no product rule was evidenced; errors can be misclassified as “restaurant name already exists.” |
| Subdomain | Zod min 2/max 63; availability on change | `ValidateAndFormat`: trims/lowercases; 2–63 bytes; regex ASCII alphanumeric/hyphen, must start/end alphanumeric; rejects consecutive hyphens/reserved values; checks tenant slug before create | tenant slug/domain `VARCHAR(255) UNIQUE`; restaurant `UNIQUE(tenant_id,slug)` | DB permits broader values than app. Frontend max counts UTF-16 code units, Go `len` counts bytes. Availability query can race/stale. Final unique tenant constraint is authoritative. |
| Phone | Optional string; no max/format/normalization; `type=tel` only | No phone validation or normalization; unchanged value inserted | `tenants.phone VARCHAR(20)` and `users.phone VARCHAR(20)` | Likely root-cause field. Empty string (not NULL) sent/stored when omitted by form; international format accepted only if ≤20 characters by DB. |
| Password | Zod min 6; no max/strength; show/hide | Required exact nonempty; min 6; bcrypt default cost | `users.password_hash VARCHAR(255) NOT NULL` | Policy is minimal and no documented maximum; byte-length/bcrypt limits not aligned or translated. Hashing occurs inside transaction after tenant insert. No source evidence of password in registration logging, but broad operational logs were not exhaustively audited. |
| Confirm password | Must equal password via Zod refine; not explicitly required but empty fails mismatch unless password empty | Not part of request DTO or backend contract | No column | Frontend-only confirmation is appropriate; field is not sent. Revalidation mode unspecified. |

Additional schema observations: tenant plan is inserted as literal `free`, while `subscription_status` is migration-defaulted to `trial` (migration 110:5); no explicit trial dates are inserted in the shown registration repository. No evidence of a registration-created subscription-history row. Runtime schema/drift was not queried.

## 6. Error Handling Findings

- Registration handler returns HTTP 400 for most use-case errors, including wrapped unexpected database errors; only strings containing “already exists”/“already taken” become 409 (`auth_handler.go:34-41`). This means internal PostgreSQL text may be disclosed as the 400 body. **VERIFIED.**
- Response shape is `{ "error": "..." }`, not field-keyed and has no stable code/request ID (`auth_handler.go:215-217`). **VERIFIED.**
- Frontend renders `error.response.data.error` in a toast (`register/page.tsx:132-141`); no field-specific server mapping or form error summary. Entered values are not cleared by explicit reset, so typical RHF behavior preserves them. **VERIFIED from code.**
- No explicit register body size limit, strict JSON unknown-field policy, request correlation ID handling, or database error classification was found in the inspected handler path. **VERIFIED within inspected files; wider middleware coverage should be rechecked in implementation.**
- Registration origin/rate protections exist: trusted origin + IP limit 10 (`main.go:484`). Availability GET has no visible rate limiter in nearby registration route declarations. Rate limiter storage/distribution scope was not established.
- Proposed compatible rollout: preserve current `error` string for existing clients while adding stable `code`, `fields`, and `requestId` fields; stop serializing raw wrapped errors; classify validation (400/422 per chosen convention), conflicts (409), rate limit (429), and unexpected dependency/internal failures (500/503) centrally. Use an opaque correlation ID in safe logs. Keep messages generic for duplicate email where account enumeration is a concern, balancing signup UX with existing product behavior.

| Case | Current evidence / proposed handling |
|---|---|
| Missing/invalid fields | Current 400 string; proposed stable validation code + field map, no SQL text. |
| Duplicate email/subdomain | Email usecase says exists; slug conflict mapped; handler 409 only via string matching. Replace with typed errors. Avoid unnecessary account disclosure. |
| Phone/length/weak password/password mismatch | Client-only or absent validation; backend must enforce authoritative contract (confirm password can remain client-only). |
| Rate limit | Registration middleware exists; 429 behavior must be confirmed in middleware tests. |
| Network/timeout/unavailable | Frontend handles `ERR_NETWORK` only; timeout/server status falls to generic/raw response. Preserve values and offer retry; do not blindly retry non-idempotent create without idempotency strategy. |
| DB unavailable/internal | Currently likely exposed as 400 with internal wrapped text. Return safe 500/503 and log correlation ID. |
| Partial write/concurrent conflict | Transaction handles pre-commit partial writes; unique constraints + typed 409. Post-commit session failure can create account despite failed response. |

## 7. Console Error Classification

| Console message | Source | Severity | Registration impact | Action |
|---|---|---:|---|---|
| `Cannot read properties of undefined (reading 'toLowerCase')` from reported `content.*.js` | **UNKNOWN**; filename alone does not identify application vs extension | Unknown | **UNKNOWN**. App has a guarded `restaurantSlug.toLowerCase()` in preview, but a truthy non-string could still fail; ordinary input is string. No matching error occurred in inspected session. | Capture stack and script origin with extensions disabled and reproduce. Do not attribute yet. |
| `POST /api/v1/auth/login 401` | Could be expected invalid credentials or a separate auth/session problem | Unknown | Not enough evidence; no current reproduction | Inspect request purpose/body safely and response contract; a 401 alone is not a registration defect. |
| `POST /api/v1/auth/register 400` | Registration API; reported response indicates tenant insert failure | Major | Direct | Root cause above; return classified safe error and field feedback. No POST was replayed. |
| `scroll-behavior: smooth` warning | Framework/browser compatibility warning per report; exact origin/version unverified | Low | No demonstrated impact | Inspect computed CSS/Next version; only change if framework guidance confirms a project rule. |
| `Can't add file system: <illegal path>` | **UNKNOWN**; likely tooling/browser/extension, insufficient evidence | Unknown | No demonstrated impact | Capture emitting component/stack. |
| React DevTools suggestion | React development tooling | Informational | None | No action. |
| Next.js Fast Refresh logs | Next.js development runtime | Informational | None unless accompanied by failed refresh/error | No action based on report alone. |
| Current inspected session: 3 unused `i.pravatar.cc` preload warnings | Page document resource preload/runtime | Low, not registration-blocking | None observed | Remove unrelated preload or scope it appropriately if confirmed in document source. |

During live inspection the registration URL returned successfully, the accessibility snapshot showed all seven fields and controls, and Playwright reported **0 console errors** and 3 preload warnings. No user data was submitted.

## 8. UI/UX Review

**Runtime inspected:** `http://localhost:38111/en/auth/register`. Desktop viewport snapshot: 2532×1041. Mobile emulation: 390×844. The form is centered, single-column, uses visible labels and input placeholders, differentiates primary action, offers password visibility toggles, and supports page `dir` from locale. Current design uses shared `Button`, `Input`, and `Label` components.

**Observed issues:**

- Mobile rendered content measured 416px wide at a 390px viewport; screenshot showed horizontal overflow / clipped right background. **VERIFIED in browser snapshot.**
- No required indicators or field error summary. Inputs are labelled, but errors are rendered as plain `<p>` without `aria-describedby`, `aria-invalid`, or live announcement. **VERIFIED in source.**
- Password eye buttons lack accessible names/`aria-label`. **VERIFIED.**
- Availability success/error communicate using icon/color/text, but the checking/error distinction is incomplete: request failure is shown in the same red “unavailable” state as a taken slug; no stale-response protection. **VERIFIED.**
- Generic submission/server errors use toast only; no inline summary and no focus management. **VERIFIED.**
- Loading label and disabled button exist; success toast/redirect exist. Browser redirect/session flow was not tested.
- Phone label says optional; omission becomes empty string. Password placeholder suggests six chars, matching minimum only.
- Arabic locale `dir` is applied to the wrapper and icon positions account for RTL; Arabic rendering and keyboard/screen-reader behavior were not directly exercised.

**PROPOSED states:** Default neutral; focus visible; valid/error accompanied by text and semantic state; subdomain checking/available/taken/format/network states distinct; disabled/loading announced; preserve values on server failure. Use inline field errors and an accessible summary, `aria-describedby`/`aria-invalid`, focus first invalid input after submit, and `aria-live=polite` for async status. Keep existing design system, fix mobile width/spacing, test touch targets and Arabic RTL.

## 9. Security Review

- **Input boundary:** Frontend schema is not authoritative; handler has no DTO validation, body limit, or explicit field lengths. Backend has minimal required/password/subdomain checks. **VERIFIED.**
- **SQL injection:** The inspected registration SQL uses parameter placeholders for values. However, the default theme config is assembled through string concatenation using `RestaurantName` and passed as `$5::jsonb`; quotes/backslashes in names can make invalid JSON and special input is not JSON-escaped. This is a reliability/encoding flaw, not SQL interpolation. **VERIFIED in `auth_repo.go:127-179`.** Construct JSON with a serializer.
- **Password storage:** bcrypt is used with `bcrypt.DefaultCost`; registration returns user after clearing hash. No direct password log found in inspected auth path. Maximum length and request/body limits are not established. **VERIFIED / UNVERIFIED as scoped.**
- **Tenant ownership/isolation:** Owner record is inserted with tenant ID; restaurant belongs to tenant; seeded ADMIN role is assigned to that user. This is consistent in shown creation transaction. No registration integration test was found to prove isolation after creation. **VERIFIED code / test unverified.**
- **Abuse:** Registration IP rate limit and trusted-origin wrapper exist. Check-subdomain is an unthrottled public GET in route block. Shared/distributed limiter behavior not established. No evidence CAPTCHA should be added.
- **Enumeration:** Registration duplicate email returns an explicit “email already exists” business error. Assess whether this is acceptable against signup threat model; no generic anti-enumeration policy found.
- **Concurrency:** Slug precheck has TOCTOU race, final tenant slug unique constraint prevents duplicates; transaction retries generated slugs only when slug was not explicitly supplied. Duplicate email precheck can race because DB uniqueness is tenant-scoped, not global. No idempotency key.
- **Cookie/session:** Current local dirty code establishes an HttpOnly cookie session after DB commit. SameSite Lax and Secure in TLS/production are visible; logout has trusted-origin check. Full cookie/domain/CORS/CSRF deployment configuration was not audited end-to-end. The registration redirect/callback is not verified.

## 10. Database Contract Findings

| Frontend field | DTO | Backend validation/normalization | Domain persistence | DB column/constraint |
|---|---|---|---|---|
| Full name | `RegisterRequest.Name` (`domain/user.go:32-38`) | Required exact nonempty; none else | `users.name` insert | `VARCHAR(255) NOT NULL` |
| Email | `RegisterRequest.Email` | Required; repository existence check | tenant email + user email | `tenants.email VARCHAR(255)`; `users.email VARCHAR(255) NOT NULL`; unique `(tenant_id,email)` |
| Restaurant name | `RegisterRequest.RestaurantName` | Required exact nonempty | tenant name + restaurant name | tenant name `VARCHAR(255) NOT NULL UNIQUE`; restaurant name `VARCHAR(255) NOT NULL` |
| Subdomain | `RegisterRequest.RestaurantSlug` | Trim/lowercase; ASCII slug validation; availability + insert conflict | tenants slug/domain and restaurants slug | tenants slug/domain `VARCHAR(255) UNIQUE`; restaurants unique `(tenant_id,slug)` |
| Phone | `RegisterRequest.Phone` | None | tenant phone then user phone | both `VARCHAR(20)` nullable in DDL, but request uses plain string and insert sends `''` when omitted |
| Password | `RegisterRequest.Password` | Required; min 6; bcrypt hash | users password_hash | `VARCHAR(255) NOT NULL` |
| Confirm password | Frontend only | Not sent | None | None |

The first migration declares global uniqueness on `tenants.name` and tenant slug/domain, whereas a restaurant name need not be globally unique by any evidenced product requirement. `users` only has per-tenant email uniqueness despite a cross-tenant duplicate check. Migration files provide intended schema, not proof that the running database applied the same sequence. The live schema was not queried because no verified disposable database was available.

The reported `VARCHAR(20)` failure should be fixed through a shared documented contract and matching tests; do not assume the database should be widened until accepted phone format is decided. For future schema sharing, contract tests or generated API schemas are lower-risk than trying to share Go and TypeScript validators directly.

## 11. Test Coverage

**Existing coverage found:** Broad backend auth tests exist elsewhere, including `backend/internal/pkg/jwt/*_test.go`, middleware/rate-limit tests, and general backend suite. Landing app has login tests and a multi-tenant login integration test. No registration-page/form test or registration-specific backend handler/use-case/repository integration test was found in the targeted paths. No registration Playwright acceptance spec was found. CI workflow configuration exists but registration-specific jobs were not identified.

**Not run:** No tests or build checks were run; this is a source/report audit and the current worktree contains broad unrelated edits. No test was marked passed.

**Recommended acceptance coverage:** Valid English and Arabic name; optional phone omitted and international phone supplied; valid subdomain/session/redirect; missing and whitespace-only fields; invalid/duplicate email; overlong name/restaurant/subdomain/phone; invalid slug characters/reserved slug; weak password/mismatch; duplicate and simultaneous slug; duplicate email race; DB constraint rejection and rollback; network/timeout/unavailable; double-submit/retry behavior; no partial tenant/user/restaurant/theme/role creation; safe errors; tenant isolation; server rejection independent of frontend; field placement/clear-on-correction/first focus; keyboard and screen reader; 390px mobile without overflow; RTL; visible loading; preserved fields; stale availability response and availability conflict. These are **PROPOSED**, not existing/passed.

## 12. Proposed Future Architecture

Use one explicit registration contract: trim strings where semantically valid; define Unicode-aware name rules without excluding international names; normalize email consistently with product policy; formalize phone as optional/null and international-friendly; define slug ASCII/case/reserved/max contract; bound password length with a security rationale; keep confirm-password client-only. Enforce authoritative server rules before repository work and retain DB constraints as final concurrency enforcement.

Keep React Hook Form + Zod and define `mode: onBlur`/revalidation behavior intentionally. Share rule examples through contract tests or generated schema rather than drift-prone duplicated prose. Make availability abortable/sequence-safe, distinguish network from taken, and always revalidate uniqueness in the transaction/insert path.

Use typed domain errors and a structured response with stable code, message, field errors, and request ID. Preserve compatibility during migration by retaining the existing `error` string temporarily. Never expose driver/SQL details. Use a request-scoped transaction and consider an idempotency token for safe retry; explicitly handle post-commit session/token failures so the UX does not imply creation failed when it succeeded.

## 13. Prioritized Remediation Backlog

- **P0 — Registration failure:** Confirm redacted phone value and actual disposable DB schema; align phone contract across frontend/server/database; map overlength input to actionable field error and safe status. Ensure unexpected SQL errors do not return as 400/raw text.
- **P1 — Contract and consistency:** Validate all DTO fields server-side; align trimming/length/normalization/uniqueness rules; replace string-matched conflicts with typed errors; review global tenant-name uniqueness and cross-tenant email policy; JSON-encode theme config; request-context transaction; handle post-commit failures/idempotent retry.
- **P1 — Async reliability/security:** Abort/stale-guard availability checks; final uniqueness enforcement; rate limit availability endpoint; test concurrent registration and duplicate email race.
- **P2 — UX/accessibility:** Fix 390px overflow; inline field mapping/summary; accessible invalid/error/async semantics and password button labels; focus first invalid field; distinct availability network/taken states; verify Arabic RTL.
- **P3 — Refinements:** Product-approved password strength hints, success/onboarding copy and optional availability suggestions; optimize unrelated preload warnings only if trace confirms they are emitted by this page.

## 14. Implementation Plan

Affected paths below are **expected candidates, not final scope**; implementation should re-check the changed checkout.

| Phase | Scope / likely files | Migration, tests, rollback / compatibility | Definition of Done | Effort |
|---|---|---|---|---|
| A — Root cause | `backend/internal/domain/user.go`; `backend/internal/usecase/auth_usecase.go`; `backend/internal/handler/http/auth_handler.go`; `backend/internal/repository/auth_repo.go`; landing register page; possibly migration | First establish verified valid phone contract/schema. If column changes, additive/reversible migration and disposable DB test; unit + API integration. Keep old inputs valid; no truncation. Roll back app/migration only after confirming no longer values were stored. | Overlong invalid input gets field-level response before persistence; valid agreed phone formats create account; DB error never leaks. | Small–medium |
| B — Contract alignment | Same DTO/usecase/repository + frontend register schema; API contract tests | Migration only for deliberate constraints/index changes. Boundary tests for all fields; staged compatibility review. | Matrix enforced at UI, API, DB; no undocumented mismatch. | Medium |
| C — Live validation | `frontend/apps/landing/src/app/[locale]/auth/register/page.tsx` and focused component/hooks/tests | No migration. Component tests + mocked delayed/reordered availability responses. Rollback UI-only; preserve endpoint contract. | Blur/submit behavior, request cancellation/stale guard, pending/taken/network distinction all covered. | Small–medium |
| D — Error contract | `backend/internal/handler/http/auth_handler.go`, shared error response/middleware; landing `src/lib/api.ts` and register form | Additive response shape first; contract tests for consumers; preserve `error` during migration. Rollback response extension without reverting safe server mapping. | Stable codes/fields/request ID; safe logs; retryability handled; no raw internals. | Medium |
| E — UI/UX | Register page, shared UI primitives only as needed, localized messages | No migration. Accessibility/component tests, desktop/mobile/RTL visual review. | No overflow; clear field states/errors/loading/success; keyboard/screen reader semantics. | Medium |
| F — Regression/security | Backend auth tests; landing component tests; Playwright registration spec/config; isolated DB fixture | Disposable DB only; add rollback/cleanup safeguards; CI jobs. | Acceptance matrix passes, including concurrency, rollback, tenant isolation, no leakage. | Large |
| G — Delivery verification | CI workflows and deployment docs/config as needed | Forward/rollback migration in staging disposable clone; smoke tests; no production-readiness claim until evidence complete. | CI green and staged end-to-end evidence on approved environment; rollback tested. | Medium |

## 15. Acceptance Criteria

Registration is eligible for production-readiness review only when: all fields have documented aligned contracts; invalid payloads are rejected server-side before writes; valid English/Arabic names and supported phone formats work; subdomain availability is race-safe and unique constraint conflicts are actionable; signup writes are atomic and retry-safe; post-commit session failures are recoverable; errors expose no SQL/secrets and correlate to safe logs; owner/role/session/tenant isolation is verified; rate limits are exercised; frontend is accessible and responsive at mobile and RTL widths; focused unit/integration/E2E tests cover the matrix above; disposable DB migrations and rollback are verified; CI passes. None of these full gates were run as part of this audit.

## 16. Risks and Dependencies

- Original registration request payload and PostgreSQL error diagnostic are unavailable; exact phone length and live failing column remain unverified.
- Current checkout is dirty, including the primary registration page and auth code. Findings may include unfinished local edits and must be reconciled before implementation.
- No disposable database was identified; actual migration drift is unknown.
- The business contract for phone formats/max length, global restaurant-name uniqueness, cross-tenant email reuse, email verification, and post-registration onboarding/session behavior requires product confirmation.
- `CheckEmailExists` precheck and tenant-name unique constraint semantics may change conflict behavior when aligned.
- Any schema widening/constraints may affect existing data and need a forward/backward compatibility plan.
- Current default theme JSON interpolation can reject valid names containing JSON-significant characters; fix should use JSON serialization.
- Production routing/domain/cookie configuration was not verified against the observed local `localhost:3003` preview.

## 17. Evidence Appendix

### Core source references

- `frontend/apps/landing/src/app/[locale]/auth/register/page.tsx:3-61` form library/schema/config.
- Same file `:63-106` watched slug/debounce/availability.
- Same file `:108-145` submit, request payload, success/error behavior.
- Same file `:167-295` fields, semantics, controls, loading and sign-in link.
- `backend/internal/domain/user.go:32-38` registration DTO.
- `backend/internal/usecase/auth_usecase.go:74-141` validation, uniqueness checks, create, post-commit token.
- `backend/internal/usecase/subdomain_validator.go:15-126` reserved words, format, length and normalization.
- `backend/internal/handler/http/auth_handler.go:27-47,167-180,208-217` registration binding/status, availability, response contract.
- `backend/internal/repository/auth_repo.go:31-83` transaction start and tenant insert; `:85-124` user/restaurant creation; `:126-215` theme/role creation and commit.
- `backend/cmd/api/main.go:483-492` auth routes and registration limiter.
- `backend/migrations/001_create_tenants.sql:4-20` tenant columns/limits/uniqueness.
- `backend/migrations/002_create_restaurants.sql:4-31` restaurant columns and tenant-scoped slug uniqueness.
- `backend/migrations/003_create_users.sql:4-20` user columns and tenant-scoped email uniqueness.
- `backend/migrations/110_add_tenant_subscription_fields.sql:5-29` subscription status default and allowed values.

### Commands and runtime evidence

- Read full `CODING_STANDARDS.md` from the repository copy after the AGENTS-specified external path was found unavailable. No code review began before this read.
- `git status --short --branch`: branch `codex/phase2a-continuation`; extensive pre-existing modified/untracked files, including the relevant registration/auth sources. No cleanup performed.
- `rg -n 'VARCHAR(20)' backend/migrations backend/internal/db/migrations`: directly related `tenants.phone`, `restaurants.phone`, and `users.phone` definitions are in migrations 001-003; many unrelated VARCHAR(20) fields also exist.
- Playwright read-only navigation to `http://localhost:38111/en/auth/register`: page title and form available. Desktop snapshot 2532×1041; mobile viewport set to 390×844 produced 416px content width. Current session console: zero errors, three unused `i.pravatar.cc` preload warnings.
- No POST registration, DB query, test, migration, source edit, commit, or deployment was performed.

### Report status

`ROOT CAUSE VERIFIED: NO` — overlength phone in `tenants.phone` is the best-supported diagnosis, but no original payload or verified disposable DB runtime failure was available.  
`VALIDATION CONTRACT CONSISTENT: NO`  
`LIVE VALIDATION COMPLETE: NO`  
`ERROR HANDLING PRODUCTION-SAFE: NO`  
`REGISTRATION UX ACCEPTABLE: NO`  
`REGISTRATION PRODUCTION-READY: NOT ESTABLISHED`  
`READY FOR REMEDIATION PHASE: YES` — pending confirmation of phone format contract and reconciliation of the dirty working tree.
