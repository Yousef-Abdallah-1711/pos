# Authenticated E2E results

## Isolated browser acceptance — 2026-11-11

- The disposable seven-service stack is running under Compose project `pos-remediation-20261111`: PostgreSQL (tmpfs-backed), Redis, API, landing, dashboard, restaurant website, and super-admin.
- Full Phase 2A Playwright suite: **15 passed, 0 failed** (`frontend/apps/dashboard/playwright.phase2a.config.ts`). The final Playwright `.last-run.json` reports `status: passed` and an empty `failedTests` list.
- Coverage includes invalid and valid login, logout, protected-route redirect, authenticated dashboard CLS (**0.047** in final run), tenant A/B product isolation, user CRUD and foreign-ID side effects, sibling restaurant tenant-wide user access, POS order creation, inventory persistence and concurrent/rejected adjustment behavior, storefront public order creation, storefront browser cart/checkout, super-admin login/access controls, and fail-closed impersonation attribution.
- The first run exposed a host/container networking mismatch. After using the actual isolated ports, two app-level issues were reproduced: SSR used the browser-only API host, and the test tenant subdomain was absent from disposable CORS origins. The server-only API base URL and fixture CORS allowlist were corrected; storefront order and browser checkout both pass in the final full run.
- Restaurant website API configuration unit tests: **5 passed**. Restaurant website TypeScript check: **PASS**. Restaurant website production build: **PASS**.
- Registration Playwright suite: **12 passed, 0 failed** on the same isolated API/database stack. This includes a real registration creating disposable account/tenant data, persisted HttpOnly session, logout and re-login, duplicate email/subdomain rejection, phone validation, network failure recovery, keyboard/accessibility checks, and English LTR/Arabic RTL registration presentation.
- The browser suite intentionally observes expected 401 responses for invalid credentials and unauthenticated session checks. These are expected assertions, not unhandled application failures.

- KDS WebSocket hook unit tests: 2 passed, including cookie-backed session connection eligibility and unauthenticated non-connection.
- Browser loaded the landing page at `/en` successfully with no browser-console errors. This was a UI render check, not an authentication E2E.
- PostgreSQL backend integration suites passed on disposable PostgreSQL 17.7: restricted-role RLS pilot, auth audit/settings/preferences RLS, payment recovery/idempotency/concurrency (14 tests), webhook retry persistence/backoff/lease recovery (4 tests), and repository tenant/parent-reference isolation (27 tests). These tests exercise backend logic and data constraints, not browser authentication.

## Prior bounded evidence retained

The prior isolated run reported passing disposable checks for signup (201), logout (204), valid login (200), invalid login (401), unauthenticated reads (401), own settings (200), cross-user settings denial (403), category and product creation (201), cash order creation (201), payment status update (200), and one KDS status transition. Those results apply only to their exact paths and prior disposable stack.

## Still not covered by this browser suite

Session refresh and expiration, tenant selection workflows, full KDS WebSocket event delivery, and suspended-tenant behavior remain unverified. The suites are scoped to the described paths and do not establish complete application security or production readiness. PostgreSQL-backed payment and isolation tests remain separate evidence from browser acceptance.
