# Inventory (dashboard) - static audit

Status: **Partial / Stub for adjustments** - viewing works through products; stock adjustment calls a non-existent endpoint.

## Evidence
- `dashboard/inventory/page.tsx` (314 lines): loads `GET /products` (`api.get`) and computes low-stock client-side `quantity_in_stock <= low_stock_threshold` (line ~118). No MODULE_IDS gate.
- Adjust dialog (react-hook-form, fields product_id, quantity_change, reason default 'adjustment', notes) submits `POST /inventory/adjust` (`inventory/page.tsx:107`). **No such backend route** in `backend/cmd/api/main.go` (FE-DASH-001, P1). Error path toasts `error.response?.data?.error` (line 116).
- Backend inventory-related routes available: `GET /products/alerts/low-stock` (unused by UI), `GET /reports/inventory` (used by reports/inventory page via `reportsApi`), `GET /reports/financial/inventory-valuation`.
- Backend has `inventoryMovementRepo` (main.go:136) - suggests movements exist server-side but no HTTP surface was found for manual adjustment.
- Navigation: Sidebar item `navigation.inventory` gated by PRODUCTS (`Sidebar.tsx:129`).
- Reports: `reports/inventory/page.tsx` (gate REPORTS) and financial `inventory-valuation` page present.

## Not present
Stock counts, receiving/purchase-order to stock, recipes/BOM, warehouses, movement history UI, supplier-linked receiving.

## Findings
FE-DASH-001 (orphan adjust), FE-DASH-014 (ungated page), FE-DASH-026.
