# 08 — Runtime Environment & Browser Pass

**Environment audited:** local Docker dev stack (`docker compose up`, project `pos-dev`) on Windows 10, Docker 29.3.1. Database is a throwaway Postgres 15.19 volume created for this audit (1 tenant, 1 user, 1 super admin, all migration-seeded; no real data). Go is not installed on the host, so nothing Go was compiled or tested. Redis is not part of any compose file (cache.go reads optional `REDIS_URL`).

## Services (verified running)
| Service | URL | Result |
|---|---|---|
| API | :8080 `/api/v1/health` | 200, `db_healthy:true` |
| Landing | :3001 | 200 (via /en) |
| Dashboard | :3002 | 200 after redirect; unauth users bounced to landing login |
| Storefront | :3003 (`demo.localhost:3003` resolves in Chrome) | tenant `demo` does not exist → redirected to `/maintenance` |
| Super admin | :3004 | 200 `/login`; `/dashboard` unauth → `/login` |
| Postgres | internal (5432 published in dev) | healthy; 91 tables, 504 indexes, 262 FKs, 53 triggers, 0 RLS policies, only superuser role |

Startup: backend ran migrations on boot (107 recorded as applied vs 123 files in `backend/migrations`; cause = duplicate numeric prefixes, DB-014, not confirmed). Frontend first install of the workspace took 2m18s; Next dev compiles on first hit (up to ~13 s per route).

## Logs (evidence/backend-logs/*.log)
- Postgres: 0 errors, 2 warnings. Landing/dashboard/storefront/super-admin: 0 error lines.
- Backend: `pos-backend.log` 565 lines. The dominant ERROR lines are the *expected* "tenant not found – slug: demo" (a business-rule 404 logged at ERROR level — see OBSERVABILITY-GAPS: business rejections and operational failures are indistinguishable). No request IDs, unstructured `log.Printf`.
- Storefront (dev): Next logs a browserslist-age notice only.

## Browser pass (Playwright MCP, real Chromium) — PARTIAL
39 page loads = 13 URLs × 3 viewports (320×720, 768×1024, 1440×900). Screenshots: `evidence/screenshots/` (33 files; the 6 dashboard loads that never settled have none).

| Page | Result |
|---|---|
| Landing `/en`, `/ar`, `/en/auth/login`, `/en/auth/register`, `/ar/auth/login` | Render at all 3 viewports; no console errors; no horizontal overflow; Arabic sets `lang=ar dir=rtl` |
| Landing `/login` (bare) | 307 → `/en/login` → **404** at all viewports (dead redirect; confirms FE-OTH-009) |
| Dashboard `/en/dashboard`, `/en/dashboard/pos` (anonymous) | Client-side redirect to landing login; page title "Loading http://localhost:3001/…"; DOM destroyed mid-navigation so no layout metrics (auth is client-side only — no server middleware) |
| Storefront `/en`, `/en/menu`, `/ar` | All → `/maintenance` because tenant `demo` is absent (API 404 on `/public/restaurants/demo/info`). `/ar` lands on a page with `lang=en dir=ltr` — confirms FE-OTH-014 |
| Super admin `/login`, `/dashboard` | Render; `/dashboard` unauth → `/login`; `<html>` has no `dir` attribute |

**Not covered (marked NOT_TESTED / BLOCKED):** every authenticated dashboard and super-admin route (no approved test accounts/fixtures; the only login is the seeded super admin with a documented password, which I did not use); role-based contexts (owner/cashier/kitchen/driver/accountant/HR); storefront with real menu/cart/checkout (needs an isolated tenant fixture); viewports 375/390/1024/1280/1920; accessibility tooling; Core Web Vitals (dev-mode server, not representative). Counts in FINAL-MASTER-AUDIT-REPORT §C.
