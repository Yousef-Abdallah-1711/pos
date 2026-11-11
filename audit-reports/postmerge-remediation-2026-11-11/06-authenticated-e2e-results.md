# Authenticated E2E results

## Isolated browser acceptance — 2026-11-11

- The disposable seven-service stack is running under Compose project `pos-remediation-20261111`: PostgreSQL (tmpfs-backed), Redis, API, landing, dashboard, restaurant website, and super-admin.
- Full Phase 2A Playwright suite: **17 passed, 0 failed** (`frontend/apps/dashboard/playwright.phase2a.config.ts`). The final completed run reports all 17 tests passing.
- The final run used the isolated stack's published host URLs (`19111`–`19114`, API `18082`). Earlier attempts with stale `131xx` defaults failed at connection setup; a repeat during an active one-minute rate-limit window also failed storefront tracking. The final run after using the current ports and allowing the rate window to reset passed all tests.
- Coverage includes invalid and valid login, logout, active session refresh and cookie rotation, old-cookie replay denial, refresh denial after logout, protected-route redirect, authenticated dashboard CLS (0.047 in the latest full run; 0.051 in a focused repeat), dashboard responsive widths 320/390/768/1440 with no overflow, tenant A/B product isolation, user CRUD and foreign-ID side effects, sibling restaurant tenant-wide user access, POS order creation, authenticated KDS WebSocket connection and live order delivery, inventory persistence and concurrent/rejected adjustment behavior, storefront public order creation, storefront browser cart/checkout, super-admin login/access controls, and fail-closed impersonation attribution.
- The first run exposed a host/container networking mismatch. After using the actual isolated ports, two app-level issues were reproduced: SSR used the browser-only API host, and the test tenant subdomain was absent from disposable CORS origins. The server-only API base URL and fixture CORS allowlist were corrected; storefront order and browser checkout both pass in the final full run.
- Restaurant website API configuration unit tests: **5 passed**. Restaurant website TypeScript check: **PASS**. Restaurant website production build: **PASS**.
- Registration Playwright suite: **12 passed, 0 failed** on the same isolated API/database stack. This includes a real registration creating disposable account/tenant data, persisted HttpOnly session, logout and re-login, duplicate email/subdomain rejection, phone validation, network failure recovery, keyboard/accessibility checks, and English LTR/Arabic RTL registration presentation.
- The browser suite intentionally observes expected 401 responses for invalid credentials and unauthenticated session checks. These are expected assertions, not unhandled application failures.

- KDS WebSocket hook unit tests: 2 passed, including cookie-backed session connection eligibility and unauthenticated non-connection.
- Live KDS browser verification: authenticated dashboard session connected to `/kds/ws`; a POS-created order was received by KDS with its matching product name.
- Browser loaded the landing page at `/en` successfully with no browser-console errors. This was a UI render check, not an authentication E2E.
- PostgreSQL backend integration suites passed on disposable PostgreSQL 17.7: restricted-role RLS pilot, auth audit/settings/preferences RLS, payment recovery/idempotency/concurrency (14 tests), webhook retry persistence/backoff/lease recovery (4 tests), and repository tenant/parent-reference isolation (27 tests). These tests exercise backend logic and data constraints, not browser authentication.

## Prior bounded evidence retained

The prior isolated run reported passing disposable checks for signup (201), logout (204), valid login (200), invalid login (401), unauthenticated reads (401), own settings (200), cross-user settings denial (403), category and product creation (201), cash order creation (201), payment status update (200), and one KDS status transition. Those results apply only to their exact paths and prior disposable stack.

## Still not covered by this browser suite

Server-side session TTL expiration through the live API, tenant selection workflows, and suspended-tenant behavior remain unverified. The session store has separate unit coverage for TTL expiry; the browser suite verifies refresh, rotation, replay rejection, and logout invalidation. The suites are scoped to the described paths and do not establish complete application security or production readiness. PostgreSQL-backed payment and isolation tests remain separate evidence from browser acceptance.
