# Table Management - Master Implementation Document

Last Updated: 2026-02-25  
Status: Active single source of truth

This file replaces the old split docs for Table Management.  
Folder policy for `docs-organized/13-TABLE-MANAGEMENT`: keep only this file.

## 1) Feature Summary

Table Management is fully integrated across:
- Backend domain/repository/usecase/handlers/routes
- RBAC (module permissions)
- Dashboard UI (list, details, edit, floor plan, QR tools)
- Restaurant website dine-in flow (table QR to menu and checkout)
- Order pipeline (table validation + auto-occupy behavior)

## 2) Scope Delivered

Implemented:
- Restaurant tables CRUD
- Table status management (`available`, `occupied`, `reserved`, `maintenance`)
- Table layout persistence and bulk position update
- Floor-plan editor scaffold and table rendering
- Table details page and actions (including free-table flow)
- Table QR generation tools (studio + modal + preview)
- Public table-info endpoint for QR flow
- Dine-in order integration (`table_id`, `order_source`) with backend validation
- Auto-occupy table when valid dine-in order is created
- Dashboard and restaurant-site UI integration
- RBAC gating + owner role permission backfill migration
- Unit/integration tests for core flows

## 3) Data Model and Migrations

Primary migrations:
- `backend/migrations/130_create_restaurant_tables.sql`
- `backend/migrations/131_add_table_id_to_orders.sql`
- `backend/migrations/132_add_tables_rbac_module.sql`
- `backend/migrations/133_create_table_layouts.sql`
- `backend/migrations/135_backfill_tables_role_permissions.sql`

Core entities:
- Restaurant table entity and validation rules in domain layer
- Order entity extended with table context
- Layout entity for floor-plan metadata and table positions

## 4) RBAC Model

Module:
- `TABLES` module added in RBAC and wired to permission checks

Enforcement points:
- Backend routes under protected table endpoints
- Dashboard feature visibility and page-level guards

Owner access:
- Backfill migration ensures owner role receives required table permissions in existing tenants
- Dashboard sidebar/menu visibility depends on role permissions

Key RBAC implementation files:
- `backend/migrations/132_add_tables_rbac_module.sql`
- `backend/migrations/135_backfill_tables_role_permissions.sql`
- `frontend/apps/dashboard/src/lib/rbac.ts`
- `frontend/apps/dashboard/src/components/layout/Sidebar.tsx`

## 5) Backend Architecture

Domain:
- Table model and constraints:
  - `backend/internal/domain/table.go`
- Order model includes table fields:
  - `backend/internal/domain/order.go`

Repository:
- Table CRUD, status updates, stats, layout persistence, bulk positions:
  - `backend/internal/repository/table_repo.go`
- Order queries include table joins/context:
  - `backend/internal/repository/order_repo.go`
- KDS query enrichment includes table context where needed:
  - `backend/internal/repository/kds_repo.go`

Use cases:
- Table business logic:
  - `backend/internal/usecase/table_usecase.go`
- Order use case table integration:
  - validates table existence/tenant scope
  - blocks invalid table state (for example maintenance)
  - auto-occupies table for dine-in path
  - `backend/internal/usecase/order_usecase.go`

HTTP handlers/routes:
- Table endpoints:
  - `backend/internal/handler/http/table_handler.go`
- Route registration:
  - `backend/cmd/api/main.go`

Related order handlers using table data:
- `backend/internal/handler/http/admin_order_handler.go`
- `backend/internal/handler/http/public_order_handler.go`

## 6) API Surface (Implemented)

Protected table management endpoints (dashboard side):
- List tables
- Get table by id
- Create table
- Update table
- Delete table
- Update table status
- Table stats summary
- Bulk update table positions
- Layout CRUD endpoints

Public endpoint for website QR entry:
- `GET /api/v1/public/table/{id}/info`

Notes:
- Public table endpoint is rate-limited in route setup
- Tenant isolation applies for protected endpoints via auth context

## 7) Dashboard UX Flow (Admin/Owner)

Entry points:
- Sidebar `Tables` under Orders section
- Pages:
  - Tables list
  - Table details
  - Edit table
  - QR page
  - Floor-plan layout page

Main flow:
1. Open Tables list and review capacity/status/location data.
2. Create or edit table metadata.
3. Open table details for operational actions.
4. Manage floor layout positions in layout page.
5. Generate QR from QR studio/modal and distribute.
6. Use status actions to free or update table state when operations change.

UX safeguards:
- Permission boundary (access denied for missing permissions)
- Invalid table handling states
- Status badge consistency and quick actions

## 8) End-User UX Flow (Restaurant Website)

QR entry flow:
1. Customer scans table QR.
2. Lands on `/<locale>/table/<tableId>/menu`.
3. Website resolves public table info.
4. Table context is stored for session/cart.
5. Customer browses menu and proceeds to checkout.
6. Checkout submits order with `table_id` + dine-in source.
7. Backend validates table and creates order.
8. Table status auto-updates to occupied (when applicable).

Occupied/invalid behavior:
- UI warning path exists for occupied-table context
- Invalid/missing table context is guarded to prevent broken checkout submission

## 9) Test and Hardening Coverage

Backend tests:
- `backend/internal/repository/table_repo_test.go`
- `backend/internal/usecase/table_usecase_test.go`
- `backend/internal/usecase/order_usecase_table_integration_test.go`
- `backend/internal/handler/http/table_handler_test.go`

Dashboard tests:
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/AddTableForm.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/TableList.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/TableStatusBadge.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/FloorPlanErrorBoundary.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/QRBuilderModal.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/QRBuilderModal.fallback.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSOrderCard.test.tsx` (table context display touchpoint)
- `frontend/apps/dashboard/src/lib/table-qr.test.ts`
- `frontend/apps/dashboard/src/__tests__/rbac.test.ts`

Restaurant website tests:
- `frontend/apps/restaurant-website/lib/table-ordering.test.ts`
- `frontend/apps/restaurant-website/lib/menu-page-content.table-warning.test.tsx`
- `frontend/apps/restaurant-website/lib/checkout-page-content.dine-in.test.tsx`
- `frontend/apps/restaurant-website/lib/store/cart-store.table-context.test.ts`

## 10) Full File Map (Current)

Docs:
- `docs-organized/13-TABLE-MANAGEMENT/TABLE-MANAGEMENT-IMPLEMENTATION-PLAN.md`

Backend migrations:
- `backend/migrations/130_create_restaurant_tables.sql`
- `backend/migrations/131_add_table_id_to_orders.sql`
- `backend/migrations/132_add_tables_rbac_module.sql`
- `backend/migrations/133_create_table_layouts.sql`
- `backend/migrations/135_backfill_tables_role_permissions.sql`

Backend domain:
- `backend/internal/domain/table.go`
- `backend/internal/domain/order.go`

Backend repository:
- `backend/internal/repository/table_repo.go`
- `backend/internal/repository/table_repo_test.go`
- `backend/internal/repository/order_repo.go`
- `backend/internal/repository/kds_repo.go`

Backend usecase:
- `backend/internal/usecase/table_usecase.go`
- `backend/internal/usecase/table_usecase_test.go`
- `backend/internal/usecase/order_usecase.go`
- `backend/internal/usecase/order_usecase_table_integration_test.go`

Backend handlers/routes:
- `backend/internal/handler/http/table_handler.go`
- `backend/internal/handler/http/table_handler_test.go`
- `backend/internal/handler/http/admin_order_handler.go`
- `backend/internal/handler/http/public_order_handler.go`
- `backend/cmd/api/main.go`

Dashboard frontend pages:
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/page.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/layout/page.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/[id]/page.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/[id]/edit/page.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/[id]/qr/page.tsx`

Dashboard frontend table components:
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/AddTableForm.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/TableList.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/TableStatusBadge.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/FloorPlanEditor.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/FloorPlanErrorBoundary.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/FloorPlanTable.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/QRBuilderModal.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/QRPreview.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/QRStudio.tsx`

Dashboard frontend hooks/libs/rbac/nav:
- `frontend/apps/dashboard/src/hooks/useTables.ts`
- `frontend/apps/dashboard/src/lib/api/tablesApi.ts`
- `frontend/apps/dashboard/src/lib/table-qr.ts`
- `frontend/apps/dashboard/src/lib/rbac.ts`
- `frontend/apps/dashboard/src/components/layout/Sidebar.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/page.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/[id]/page.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSOrderCard.tsx`
- `frontend/apps/dashboard/src/i18n/messages/en.json`
- `frontend/apps/dashboard/src/i18n/messages/ar.json`

Dashboard frontend tests:
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/AddTableForm.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/TableList.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/TableStatusBadge.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/FloorPlanErrorBoundary.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/QRBuilderModal.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/tables/components/QRBuilderModal.fallback.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSOrderCard.test.tsx`
- `frontend/apps/dashboard/src/lib/table-qr.test.ts`
- `frontend/apps/dashboard/src/__tests__/rbac.test.ts`

Restaurant website integration:
- `frontend/apps/restaurant-website/app/[locale]/table/[tableId]/menu/page.tsx`
- `frontend/apps/restaurant-website/app/[locale]/menu/menu-page-content.tsx`
- `frontend/apps/restaurant-website/app/[locale]/checkout/checkout-page-content.tsx`
- `frontend/apps/restaurant-website/lib/table-ordering.ts`
- `frontend/apps/restaurant-website/lib/hooks/use-api-queries.ts`
- `frontend/apps/restaurant-website/lib/store/cart-store.ts`
- `frontend/apps/restaurant-website/middleware.ts`

Restaurant website tests:
- `frontend/apps/restaurant-website/lib/table-ordering.test.ts`
- `frontend/apps/restaurant-website/lib/menu-page-content.table-warning.test.tsx`
- `frontend/apps/restaurant-website/lib/checkout-page-content.dine-in.test.tsx`
- `frontend/apps/restaurant-website/lib/store/cart-store.table-context.test.ts`

## 11) Remaining Optional Enhancements (Not Blocking Core Delivery)

- Advanced QR template gallery/presets for different restaurant brands
- Expanded QR design controls (logo masking presets, frame packs, print sheets)
- Analytics around QR scans and table conversion funnel

---

Change note (2026-02-25):
- Consolidated old plan/checklist docs into this one master file.
- Updated to match current implemented backend/frontend/test state and real paths.
