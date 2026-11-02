# Orders module (dashboard) - static audit

Status: **Verified implemented** (list/detail/status/cancel/driver assign); several backend capabilities unused; no manual order creation outside POS.

## Pages
- `dashboard/orders/page.tsx` (627 lines): GET `/admin/orders` with page_size 20, filters status/payment_status/table_id/search; DELETE `/admin/orders/{id}` (cancel); CSV export of the CURRENT PAGE via `lib/pdf-export.ts exportOrdersToCSV` (client-side).
- `dashboard/orders/[id]/page.tsx` (643): GET `/admin/orders/{id}`, PUT `/admin/orders/{id}/status`; client-side allowed-transition map `STATUS_TRANSITIONS` (pending->confirmed|cancelled, confirmed->preparing|cancelled, preparing->ready ...) at lines 92-106; driver assignment via `components/driver-assignment-modal.tsx` (GET /admin/drivers/available, POST /admin/orders/{orderId}/assign).
- Gating: `useHasPermission` ORDERS (6 refs per page). Backend: `wrapWithPermission(...,10,READ|WRITE|DELETE)` main.go:1294-1301.

## Contract
All frontend paths match backend. Backend routes with NO frontend caller: GET `/admin/orders/stats`, `/admin/orders/export`, `/admin/orders/{id}/history` (order audit history not shown). `PUT /admin/orders/{id}/payment-status` only used by POS.

## Notes
- Money rendered with hard-coded `$` and `toFixed(2)` (`orders/page.tsx:503`, `[id]/page.tsx:362-392`) while accounting uses EGP (FE-DASH-018).
- Transition rules duplicated on client; server is authority (not verified that server enforces same map). Falls back to ALL options when no transitions match (`[id]/page.tsx:103-106`).
- Overview dashboard counts orders from list totals (`dashboard/page.tsx`), revenue tile = sum of one page (FE-DASH-018).
- No polling/WebSocket on orders list: stale until refresh (not verified for auto-refresh; no setInterval in the file).
- Public order creation by customers happens in restaurant-website app (out of scope).

## Findings
FE-DASH-018, FE-DASH-026, FE-DASH-024 (POS creation path).
