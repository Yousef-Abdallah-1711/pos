# Drivers / delivery (dashboard) - static audit

Status: **Partial** - driver CRUD, status toggle and order assignment present; tracking/location, stats and RBAC missing.

## Pages
`drivers/page.tsx` (323), `drivers/new/page.tsx` (365, react-hook-form), `drivers/[id]/page.tsx` (625). Assignment: `components/driver-assignment-modal.tsx`. **None of the three pages has a permission gate**; sidebar entry is gated by ORDERS (`Sidebar.tsx:118`).

## Contract
Used (match main.go:1397-1404): POST/GET `/admin/drivers`, GET `/admin/drivers/available`, GET/PUT/DELETE `/admin/drivers/{id}`, PUT `/admin/drivers/{id}/status`, POST `/admin/orders/{orderId}/assign`.
Backend routes with no UI: GET `/admin/drivers/{id}/stats`, POST/GET `/admin/drivers/{id}/location` (main.go:1405-1407) - no live map or tracking.
Backend protection: all driver routes are `wrapProtected` (auth only, no RBAC module) - any authenticated staff role can create/delete drivers or assign orders (FE-DASH-014).

## Behaviour
- Completion rate computed client-side `completed_deliveries/total_deliveries*100` (`drivers/page.tsx:240`, `[id]/page.tsx:215`, modal:99) - division guarded by ternary.
- Active orders total via `reduce` (`drivers/page.tsx:316`).
- RTL: the most RTL-aware pages (17 refs in list page).
- Customer-facing delivery tracking is not in the dashboard (restaurant-website scope).

## Findings
FE-DASH-014, FE-DASH-026.
