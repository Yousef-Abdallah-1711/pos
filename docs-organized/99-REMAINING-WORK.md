# Remaining Work (2026-01-26)

This file tracks open issues and incomplete tasks found during the latest verification. It is the active checklist until all items below are closed. Each task includes target file paths and a minimal test plan.

## Critical blockers

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
- `README.md`

## Verification checklist

- `GET /api/v1/health` returns 200
- `GET /api/v1/admin/orders` returns list and count without scan errors
- `PUT /api/v1/admin/orders/{id}/status` returns 200 with valid transition
- `GET /api/v1/admin/drivers/available` returns list (even if empty)
- `POST /api/v1/admin/orders/{id}/assign` assigns a driver
- Order tracking works for `/en/` and `/ar/` routes (including uppercase locale redirect)
- Create user shows in users list and can be assigned roles
- Inventory: out-of-stock products cannot be added to cart when tracking is enabled
