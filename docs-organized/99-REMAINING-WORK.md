# Remaining Work (Updated 2026-02-21)

## Completed Since Last Update (Feb 2026)

### Landing App — Full Refactor ✅
- **URL-based i18n routing**: Migrated from localStorage/context-based language switching to Next.js `[locale]` URL segments (`/en/...`, `/ar/...`)
- **Middleware**: New `src/middleware.ts` — detects locale, redirects bare paths to `/en`
- **New folder structure**: All pages under `src/app/[locale]/` — landing, auth, register, super-admin-login
- **LanguageContext refactor**: Now accepts `locale` prop from server (no localStorage), language toggle uses `router.push()`
- **Modern redesign**: Navbar (glassmorphism), Hero (dark POS mockup, stats), Features (gradient cards), HowItWorks (3-step, NEW section), Pricing (single plan), Testimonials, CTA, Footer
- **Pricing updated**: Single plan — 500 EGP/month after 30-day free trial (replaced 3-tier Free/Pro/Enterprise)
- **All translations**: Full EN + AR parity in `src/lib/translations.ts`

### Dashboard — Arabic Locale Enabled ✅
- Added `ar` to `generateStaticParams()` in dashboard locale layout
- Arabic UI now accessible at `/ar/dashboard/...`

### HR Module — Backend Handlers ✅
- New files: `hr_employee_handler.go`, `hr_attendance_handler.go`, `hr_leave_handler.go`, `hr_salary_handler.go`
- Routes registered in `main.go`

### HR Module — Frontend Translations ✅
- `AttendanceForm.tsx`, `LeaveForm.tsx`, `SalaryForm.tsx` — full EN + AR translations

### Restaurant Website Fixes ✅
- Middleware: Removed redirect loop, properly handles locale + maintenance routing
- All pages updated to use `next/link` instead of `<a>` tags for navigation
- Order tracking, cart, checkout, menu, settings pages — translation fixes

### Super-Admin ExtendTrialModal ✅
- Fixed 404 error on trial extension API call
- Corrected endpoint URL construction in `ExtendTrialModal.tsx`

### Frontend Styling Hardening (Inline Style Cleanup) [Done]
- Removed React inline style props from active frontend code paths (`style={{...}}`, `style={...}`) across dashboard, landing, restaurant website, shared layers, and ui-themes.
- Refactored dynamic editor previews to class-based styling with scoped `<style jsx>` blocks (no inline style objects).
- Refactored invoice HTML generator in `frontend/apps/dashboard/src/lib/pdf-export.ts` to class-based CSS (removed inline `style="..."` attributes).
- Validation snapshot:
  - `rg -n "style=\\{\\{" frontend -S -g "!**/*.disabled"` -> 0
  - `rg -n "style=\\{" frontend -S -g "!**/*.disabled"` -> 0
  - Remaining `style=` occurrences are intentional sanitizer/test literals in:
    - `frontend/apps/dashboard/src/components/editors/CustomCSSEditor.tsx`
    - `frontend/apps/dashboard/src/components/editors/__tests__/CustomCSSEditor.test.tsx`

### Paymob Phase 9 Progress (Feb 21, 2026) ✅
- `go build ./...` backend blocker resolved by isolating deprecated legacy-only packages behind `legacy` build tags.
- Phase 9 DB-backed verification path added:
  - `backend/scripts/phase9_db_verify.ps1`
  - `backend/scripts/paymob_phase9_seed.sql`
  - `docs-organized/11-PAYMOB-PAYMENT-INTEGRATION/PHASE9-DB-VERIFICATION-PATH.md`
- Remaining Paymob closeout tasks now focus on running the DB-backed script in an environment with working Docker/PostgreSQL and recording final `9.2/9.3` runtime results.

---

# Open Issues

This file tracks open issues and incomplete tasks found during the latest verification. It is the active checklist until all items below are closed. Each task includes target file paths and a minimal test plan.

## Critical Blockers

### 1) Docker healthcheck + DB connectivity
- Fix runtime healthcheck (missing `curl`).
  - Paths: `backend/Dockerfile`, `docker-compose.yml`
  - Test:
    - `docker compose build backend`
    - `docker compose up -d backend`
    - `docker compose ps` shows backend `healthy`
- Verify DB connectivity (no stub mode).
  - Paths: `backend/cmd/api/main.go`, `backend/internal/repository/*`, `backend/migrations/*.sql`
  - Test:
    - `curl -s http://localhost:8080/api/v1/health` returns 200
    - Backend logs show DB connected

### 2) Orders list count scan mismatch
- Align count query scan arguments or use a proper count-only scan.
  - Paths: `backend/internal/repository/order_repo.go`
  - Test:
    - `GET /api/v1/admin/orders?page=1&page_size=20` returns 200 and includes `total`
    - Dashboard orders page loads without SQL scan errors

### 3) Order status update 400
- Validate status transitions and ensure payload shape matches backend expectations.
  - Paths: `backend/internal/handler/http/admin_order_handler.go`, `backend/internal/usecase/order_usecase.go`, `backend/internal/domain/order.go`
  - Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/[id]/page.tsx`
  - Test:
    - Update status via dashboard for a real order (no 400)
    - `PUT /api/v1/admin/orders/{id}/status` returns 200

### 4) Driver assignment 404
- Ensure available drivers endpoint is registered and returns an array (even empty).
  - Paths: `backend/internal/handler/http/admin_driver_handler.go`, `backend/internal/repository/driver_repo.go`, `backend/cmd/api/main.go`
  - Frontend: `frontend/apps/dashboard/src/components/driver-assignment-modal.tsx`
  - Test:
    - `GET /api/v1/admin/drivers/available` returns 200 with JSON array
    - Assign driver in dashboard order details works

### 5) Order tracking locale 404
- Normalize locale to lowercase or redirect.
  - Paths: `frontend/apps/restaurant-website/middleware.ts`, `frontend/apps/restaurant-website/app/[locale]/orders/track/[orderNumber]/page.tsx`
  - Test:
    - `/AR/orders/track/{number}` redirects or loads correctly
    - `/ar/orders/track/{number}` loads correctly

### 6) Shared-components TypeScript drift in watch mode
- `pnpm dev` currently reports high-volume TS errors from `@pos/shared-components` tests and component typings (manifest shape mismatch, outdated API assumptions, missing module paths, duplicate exports).
  - Paths: `frontend/layers/shared-components/__tests__/*`, `frontend/layers/shared-components/components/*`, `frontend/layers/shared-components/utils/*`
  - Test:
    - `pnpm --filter @pos/shared-components dev` runs with 0 TS errors
    - `pnpm --filter @pos/shared-components test` passes

### 7) Next.js EPERM trace lock on Windows
- Landing app (and sometimes other Next apps) can fail with: `EPERM: operation not permitted, open .../.next/trace`.
  - Paths: `frontend/apps/landing/.next` (runtime artifact)
  - Test:
    - after cleaning `.next` and restarting dev, `pnpm --filter @pos-saas/landing dev` stays running
    - full `pnpm dev` no longer exits due to trace file lock

## Inventory end-to-end

### Backend
- Expose `track_inventory`, `quantity_in_stock`, and `low_stock_threshold` in public product list/search queries.
  - Paths: `backend/internal/repository/product_repo.go`
  - Test:
    - `GET /api/v1/public/restaurants/{slug}/products` includes inventory fields
- Persist inventory, category, featured, and description in update flow.
  - Paths: `backend/internal/handler/http/product_handler.go`, `backend/internal/usecase/product_usecase.go`, `backend/internal/repository/product_repo.go`, `backend/internal/domain/product.go`
  - Test:
    - Update product from dashboard and verify fields saved

### Frontend (Dashboard + Restaurant)
- Parse numeric inputs before submit.
  - Path: `frontend/apps/dashboard/src/components/products/ProductForm.tsx`
  - Test: Save product with numeric strings; backend accepts values
- Show out-of-stock badge and block add-to-cart when tracked stock <= 0.
  - Paths: `frontend/apps/restaurant-website/components/product-card-simple.tsx`, `frontend/apps/restaurant-website/lib/store/cart-store.ts`
  - Test: Product with `track_inventory=true` and `quantity_in_stock=0` cannot be added
- Surface insufficient inventory errors on checkout.
  - Path: `frontend/apps/restaurant-website/app/[locale]/checkout/checkout-page-content.tsx`
  - Test: Attempt checkout with out-of-stock item shows toast with reason

## Dashboard orders and drivers

### Drivers UI
- Ensure Orders -> Drivers UI exists and only shows "No available drivers" after a successful empty API response.
  - Paths: `frontend/apps/dashboard/src/app/[locale]/dashboard/drivers/page.tsx`, `frontend/apps/dashboard/src/components/driver-assignment-modal.tsx`
  - Test: Drivers page loads; assignment modal handles empty list gracefully

### CSV export
- Prevent `toFixed` on undefined totals.
  - Path: `frontend/apps/dashboard/src/lib/pdf-export.ts`
  - Test: Export CSV from orders page with no crash

### Logging noise
- Reduce debug logs in stable modules.
  - Paths: `frontend/apps/dashboard/src/lib/api.ts`, `frontend/apps/dashboard/src/lib/apiCache.ts`, `frontend/apps/dashboard/src/components/providers/DashboardAuthProvider.tsx`, `frontend/apps/dashboard/src/lib/monitoring.ts`
  - Test: Normal browsing logs only warnings/errors, not verbose debug traces

## RBAC / users

### Users list visibility
- Verify created users appear in users list; confirm tenant headers and API responses.
  - Paths: `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/users/page.tsx`, `backend/internal/handler/http/user_settings_handler.go`, `backend/internal/repository/user_repo.go`, `backend/internal/middleware/tenant.go`
  - Test:
    - Create user in dashboard
    - Users list shows new user without refresh issues

### Access control validation
- Validate access-control and roles pages against real permissions and tenant scope.
  - Paths: `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/access-control/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx`, `backend/internal/handler/http/rbac_helpers.go`, `backend/internal/handler/http/rbac_role_handler.go`
  - Test:
    - User with limited permissions sees only allowed sections

### Seed default roles (Driver role)
- Seed default roles/permissions per tenant, including a Driver role limited to assigned orders only.
  - Paths: `backend/migrations/104_create_role_permissions_table.sql`, `backend/migrations/105_add_rbac_modules.sql`, `backend/migrations/106_backfill_rbac_defaults.sql`, `backend/internal/repository/rbac_role_repository.go`
  - Test:
    - New tenant has default roles
    - Driver role can read assigned orders only

## Docs to update after fixes

- `docs-organized/07-website-ecommerce-plan/06-CURRENT-STATUS-AND-FINDINGS.md`
- `docs-organized/07-website-ecommerce-plan/08-REMAINING-PHASES-AND-ROADMAP.md`
- `docs-organized/07-website-ecommerce-plan/09-ECOMMERCE-PRODUCTION-READY-VERIFICATION.md`
- `docs-organized/08-rules-role-based/01-RBAC_CURRENT_STATUS.md`
- `docs-organized/08-rules-role-based/02-RBAC_NEXT_PHASES.md`
- `docs-organized/ORGANIZATION_SUMMARY.md`
- `docs-organized/03-FRONTEND-DASHBOARD/COMPLETION_REPORT.md` (keep historical status clearly labeled)
- `docs-organized/01-THEME-SYSTEM/06-THEME-BUILDER-PRODUCTION-TODO.md`

## Verification checklist

- `GET /api/v1/health` returns 200
- `GET /api/v1/admin/orders` returns list and count without scan errors
- `PUT /api/v1/admin/orders/{id}/status` returns 200 with valid transition
- `GET /api/v1/admin/drivers/available` returns list (even if empty)
- `POST /api/v1/admin/orders/{id}/assign` assigns a driver
- Order tracking works for `/en/` and `/ar/` routes (including uppercase locale redirect)
- Create user shows in users list and can be assigned roles
- Inventory: out-of-stock products cannot be added to cart when tracking is enabled
- `rg -n "style=\\{\\{" frontend -S -g "!**/*.disabled"` returns 0
- `rg -n "style=\\{" frontend -S -g "!**/*.disabled"` returns 0
