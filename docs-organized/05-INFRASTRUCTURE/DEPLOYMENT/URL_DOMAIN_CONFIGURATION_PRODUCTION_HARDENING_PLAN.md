# URL & Domain Configuration Production Hardening Plan

## Purpose

Make URL/domain configuration production-ready across all frontend apps and backend without breaking local development (`localhost`, `demo.localhost`).

This plan focuses on removing production-risk hardcoded local URL fallbacks from runtime code while preserving:

- local development flows
- test/load scripts
- `.env.example` onboarding defaults
- local-only docs/examples (clearly labeled)

## Core Rules

- Code must read URLs/domains from environment variables.
- Do not hardcode real production domains in source files.
- `localhost` fallback is allowed only in development mode.
- Production mode must fail fast (or clearly error) when required URL env vars are missing.
- Docs/tests/scripts may keep local URLs, but must be labeled as local/dev usage.

## Scope

In scope:

- `backend`
- `frontend/apps/restaurant-website`
- `frontend/apps/dashboard`
- `frontend/apps/landing`
- `frontend/apps/super-admin`
- env templates and deployment config docs

Out of scope (for this plan's code changes):

- removing localhost examples from local onboarding docs (unless mislabeled as production)
- deleting test suites
- changing tenant local subdomain parsing logic required for local development

## Phase 0: Baseline and Policy (Planning)

### TODO

- [ ] Create/confirm URL-domain configuration policy for all apps and backend.
- [ ] Inventory all `localhost`, `.localhost`, `127.0.0.1`, and local-only URL references.
- [ ] Classify findings into:
  - runtime code
  - env templates
  - test/load scripts
  - docs
- [ ] Freeze first-pass scope to runtime code only.
- [ ] Save grep baseline for before/after comparison.

### Verification Checklist

- [ ] A reviewed list of runtime files to refactor exists.
- [ ] Non-runtime local URLs are explicitly marked to keep/defer.
- [ ] No runtime changes are made in this phase.

## Phase 1: Shared URL Resolution Pattern (Frontend + Backend)

### Goal

Define one consistent pattern for resolving URLs/domains safely in development and production.

### TODO

- [ ] Define frontend resolver behavior:
  - env var first
  - dev-only localhost fallback
  - production error if missing
- [ ] Define backend resolver/validation behavior:
  - env var first
  - dev-only localhost fallback
  - production startup validation error if missing
- [ ] Define required env vars per app (API, dashboard URL, landing URL, restaurant domain, app/webhook URLs, etc.).
- [ ] Define a shared "is production" check pattern per app/backend.
- [ ] Document allowed/forbidden fallback behavior in code comments and docs.

### Verification Checklist

- [ ] Every app has a defined source of truth for API URL and public URLs.
- [ ] Team can state which env vars are required in production for each app.
- [ ] Missing-env behavior is consistent and intentional.

## Phase 2: Restaurant Website (Highest Priority Public Runtime)

### Goal

Remove production-risk localhost fallbacks from the storefront runtime while keeping local multi-tenant dev flow (`demo.localhost`) working.

### Known Runtime Files to Refactor

- `frontend/apps/restaurant-website/app/[locale]/page.tsx`
- `frontend/apps/restaurant-website/next.config.js`
- `frontend/apps/restaurant-website/middleware.ts`
- `frontend/apps/restaurant-website/lib/config/api.ts`
- `frontend/apps/restaurant-website/lib/api/get-theme-server.ts`
- `frontend/apps/restaurant-website/lib/subdomain-context.tsx`
- `frontend/apps/restaurant-website/lib/checkTenantStatus.ts`
- `frontend/apps/restaurant-website/lib/api-client.ts`
- `frontend/apps/restaurant-website/lib/api/api-client.ts`

### TODO

- [ ] Centralize API base URL resolution in one helper used by all storefront runtime callers.
- [ ] Remove unconditional `http://localhost:8080/api/v1` fallbacks in production code paths.
- [ ] Keep localhost fallback only when `NODE_ENV !== 'production'`.
- [ ] Preserve local subdomain behavior in `frontend/apps/restaurant-website/lib/subdomain.ts`.
- [ ] Ensure Theme Builder storefront runtime still renders theme header/footer/sections correctly.
- [ ] Ensure middleware tenant checks do not rely on production localhost defaults.

### Verification Checklist

- [ ] Local storefront works via `demo.localhost:3003`.
- [ ] Theme-rendered header/footer and sections still render.
- [ ] Production build/start works when required envs are set.
- [ ] Production startup/runtime clearly fails or logs config error when required envs are missing.
- [ ] No unconditional localhost fallback remains in storefront runtime code.

## Phase 3: Dashboard App

### Goal

Remove production-risk localhost defaults from dashboard runtime and Theme Builder-related API calls.

### Known Runtime Files to Refactor

- `frontend/apps/dashboard/src/lib/config.ts`
- `frontend/apps/dashboard/src/lib/api.ts`
- `frontend/apps/dashboard/src/lib/api/themeApi.ts`
- `frontend/apps/dashboard/src/hooks/usePayment.ts`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/new/page.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/theme-builder/page.tsx`
- `frontend/apps/dashboard/Dockerfile`

### TODO

- [ ] Refactor dashboard API URL resolution to env + dev-only fallback pattern.
- [ ] Remove hardcoded localhost fetch calls in pages/components.
- [ ] Review Theme Builder preview target/domain defaults and make them production-safe.
- [ ] Prevent production Docker image from embedding localhost API defaults.

### Verification Checklist

- [ ] Dashboard login and API calls work locally.
- [ ] Theme Builder editor and preview work locally.
- [ ] Production build does not embed localhost as default API endpoint.
- [ ] No unconditional localhost fallback remains in dashboard runtime code.

## Phase 4: Landing App

### Goal

Make landing, auth, and tenant onboarding URLs environment-driven in production.

### Known Runtime Files to Refactor

- `frontend/apps/landing/src/lib/config.ts`
- `frontend/apps/landing/src/lib/api.ts`
- `frontend/apps/landing/src/components/TenantSelectorModal.tsx`
- `frontend/apps/landing/src/app/[locale]/auth/super-admin-login/page.tsx`
- `frontend/apps/landing/src/app/[locale]/auth/register/page.tsx`

### TODO

- [ ] Centralize landing/dashboard/API/restaurant-domain URL resolution.
- [ ] Keep localhost defaults only in development mode.
- [ ] Ensure tenant subdomain display/preview uses configured restaurant domain in staging/production.
- [ ] Ensure super-admin-login redirects use configured dashboard URL.

### Verification Checklist

- [ ] Landing auth flows work locally.
- [ ] Tenant URL generation is correct in staging/production.
- [ ] No unconditional localhost fallback remains in landing runtime code.

## Phase 5: Super Admin App

### Goal

Make super-admin API and cross-app redirect URLs production-safe and environment-driven.

### Known Runtime Files to Refactor

- `frontend/apps/super-admin/src/lib/api.ts`
- `frontend/apps/super-admin/src/lib/debug-info.ts`
- `frontend/apps/super-admin/src/components/impersonation/ImpersonationModal.tsx`
- `frontend/apps/super-admin/src/app/dashboard/payments/page.tsx`

### TODO

- [ ] Refactor super-admin API URL resolution to env + dev-only fallback pattern.
- [ ] Refactor dashboard redirect/impersonation URLs to env-driven values.
- [ ] Update localhost-specific error messages to be environment-aware where user-facing.
- [ ] Keep local troubleshooting hints only in debug/dev contexts.

### Verification Checklist

- [ ] Super-admin dashboard API calls work locally.
- [ ] Impersonation redirect uses configured dashboard URL.
- [ ] No unconditional localhost fallback remains in super-admin runtime code.

## Phase 6: Backend Runtime URL Hardening

### Goal

Remove production-risk localhost URL defaults from backend runtime and enforce configuration validation in production.

### Known Runtime Files to Refactor

- `backend/cmd/api/main.go`
- `backend/internal/middleware/cors.go`
- `backend/internal/repository/auth_repo.go`
- `backend/internal/service/webhook_retry_worker.go`
- `backend/internal/handler/http/super_admin_handler.go`
- `backend/internal/handler/http/mock_paymob_handler.go`

### TODO

- [ ] Add backend config validation for required public URLs/domains in production startup.
- [ ] Restrict localhost fallbacks to dev mode only.
- [ ] Make CORS localhost defaults dev-only and explicit.
- [ ] Ensure webhook base/app URLs come from env in production.
- [ ] Ensure generated dashboard redirect URLs come from env in production.
- [ ] Ensure restaurant domain generation uses `RESTAURANT_DOMAIN` in production.

### Verification Checklist

- [ ] Backend starts locally with dev defaults.
- [ ] Backend rejects unsafe production startup with missing required URL envs.
- [ ] CORS works for local dev and configured production origins.
- [ ] Redirect/webhook URLs in responses/logs are correct.

## Phase 7: Environment Templates and Deployment Configuration

### Goal

Keep developer onboarding easy while making production env requirements explicit and complete.

### Files to Review

- `backend/.env.example`
- `backend/.env.production.example`
- `frontend/apps/restaurant-website/.env.example`
- `frontend/apps/restaurant-website/.env.production.example`
- `frontend/apps/dashboard/.env.example`
- additional app env examples (landing/super-admin) if missing
- Dockerfiles / deployment configs that inject URL defaults

### TODO

- [x] Keep local values in `.env.example` files for quick local setup.
- [x] Ensure `.env.production.example` files define all required production URL/domain vars.
- [x] Add missing production env examples for apps that do not have them. (Review complete: none missing; added missing local app examples instead.)
- [x] Remove localhost defaults from production-oriented Docker/deployment configs.
- [x] Add comments to clarify local vs production examples.

### Verification Checklist

- [ ] New developers can run locally from `.env.example`.
- [ ] Production/staging deploy config can be prepared from templates without guessing.
- [ ] Production templates do not imply localhost runtime behavior.

## Phase 8: Documentation Cleanup (Local vs Production Split)

### Goal

Keep docs useful for local development while removing localhost ambiguity from production-readiness docs.

### TODO

- [ ] Label code blocks/examples as `Local Development`, `Staging`, or `Production`.
- [ ] Keep localhost examples in setup/troubleshooting docs where appropriate.
- [x] Add production examples using placeholders (e.g. `https://api.yourdomain.com`).
- [x] Fix production-readiness docs that still show localhost without context.
- [ ] Keep `demo.localhost` examples only where local tenant subdomain flow is being documented.
- [ ] Add one canonical URL/domain config matrix and link to it from major docs.

### Verification Checklist

- [ ] Production guides do not instruct localhost endpoints as production defaults.
- [ ] Local dev guides still work as written.
- [ ] Theme Builder docs correctly describe `demo.localhost` as local dev flow only.

## Phase 9: Test and Script Policy (Keep, Parameterize, Label)

### Goal

Retain local/test tooling without leaking test assumptions into production runtime code.

### Examples to Keep

- `backend/tests/load/*.js`
- `backend/tests/run_phase10_newman.ps1`
- local verification scripts and docs

### TODO

- [x] Keep localhost defaults in test/load scripts, but ensure env override support (e.g. `BASE_URL`) exists.
- [x] Add a top comment/header to scripts indicating local/test-only usage.
- [x] Verify runtime code no longer depends on test-script defaults.
- [x] Document how to point load scripts to staging/prod safely.

### Verification Checklist

- [ ] Load/Newman scripts work locally with defaults.
- [ ] Scripts support non-local target via env override.
- [ ] Runtime app/backend code does not inherit test-only localhost assumptions.

## Phase 10: Validation and CI Guardrails

### Goal

Prevent regression by blocking new hardcoded localhost fallbacks in runtime directories.

### TODO

- [x] Add grep-based CI checks for runtime code paths:
  - `backend/` runtime packages
  - `frontend/apps/*` runtime code
- [x] Create allowlist rules for docs, tests, scripts, and local tooling.
- [x] Add unit tests for frontend URL/config resolver helpers.
- [x] Add backend config validation tests.
- [x] Add smoke-check checklist for cross-app URL/domain behavior.

### Verification Checklist

- [ ] CI fails when new hardcoded localhost runtime fallback is introduced.
- [ ] CI allows localhost in docs/tests/scripts where expected.
- [x] URL/config resolver tests pass.
- [x] Backend config validation tests pass.

## Final Done Checklist (Project-Level)

- [ ] Restaurant website uses env-driven API/domain config in production.
- [ ] Dashboard uses env-driven API/domain config in production.
- [ ] Landing uses env-driven API/domain config in production.
- [ ] Super-admin uses env-driven API/domain config in production.
- [ ] Backend validates URL/domain envs in production startup.
- [ ] Theme Builder preview and storefront theme rendering still work.
- [ ] Local dev (`localhost`, `demo.localhost`) still works.
- [ ] Docs clearly separate local vs production URL usage.
- [ ] CI guardrails prevent localhost fallback regressions in runtime code.

## Recommended Execution Order

1. Phase 2 (Restaurant Website)
2. Phase 6 (Backend)
3. Phase 3 (Dashboard)
4. Phase 4 (Landing)
5. Phase 5 (Super Admin)
6. Phase 7 (Env Templates and Deployment Config)
7. Phase 8 (Docs Cleanup)
8. Phase 9 (Test/Script Policy)
9. Phase 10 (CI Guardrails)

## Notes

- `demo.localhost:3003` is part of the local multi-tenant development workflow and should be preserved for development.
- The goal is not "remove localhost everywhere"; the goal is "remove localhost as a production runtime fallback."
