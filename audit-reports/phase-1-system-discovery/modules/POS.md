# POS module (dashboard) - static audit

Status: **Partial** (core cart-to-order flow implemented; payments, offline sync and server authorization have gaps).

## Files
- Pages: `dashboard/pos/page.tsx` (524 lines), `pos/settings/page.tsx`, `pos/layout.tsx` (no-zoom viewport). Components: `pos/components/*` (POSTerminal, POSCartPanel, POSProductGrid/VirtualGrid, POSPaymentSelector, POSReceipt ...).
- State: `stores/posStore.ts` (647 lines, zustand + custom persist storage `pos-storage`).
- Data: `hooks/usePOSData.ts`, `usePOSProductCache.ts` (IndexedDB), `usePOSProductIndex.ts`, `usePOSOffline.ts`, `usePOSBackgroundSync.ts`, `usePOSPrint.ts`, `usePOSBarcodeScanner.ts`, `usePOSKeyboard.ts`, `usePOSAudio.ts`; `lib/api/posApi.ts`; `public/pos-sw.js`.
- Tests: 5 vitest files (audio, background sync, barcode, performance, posApi) - not run by this audit.

## Gating
`useHasPermission(MODULE_IDS.ORDERS, WRITE)` at `pos/page.tsx:76` (UI only). Sidebar entry gated by ORDERS.

## API contract
| Call | Backend | Match |
|---|---|---|
| GET /products?page_size=500 (`posApi.ts:68`) | main.go:664 wrapWithPermission | OK, silently capped at 500 |
| GET /categories | main.go:671 wrapProtected | OK |
| POST /public/orders (`posApi.ts:228`, order_source 'pos') | main.go:472 TenantContextMiddleware only | OK but uses the PUBLIC route (FE-DASH-024) |
| PUT /admin/orders/{id}/payment-status | main.go:1298 ORDERS WRITE | OK |
Tables for the selector come from `/tables` (usePOSData).

## Behaviour verified in source
- Order payload carries items/qty/notes only, no prices (`posApi.ts:188-213`): server is price authority. Good.
- Display totals computed client-side: subtotal = sum(price*qty), tax = 14% (`posStore.ts:106,595-600`), `roundMoney`. Receipt (`pos/page.tsx buildReceiptData`) uses these client values.
- Cash orders are auto-marked `paid` after create (`pos/page.tsx:304-311`); failure is swallowed. No tendered/change, no split payment, no discount, no refund/void UI. Payment methods offered: cash/card/wallet (`posStore.ts:109`).
- Hold/recall orders implemented in store (`holdCurrentOrder`, `recallOrder`); persisted in localStorage.
- Barcode scanner, keyboard shortcuts, audio feedback, receipt print (browser print) implemented.
- Offline: IndexedDB queue `pos-offline-db/pending-orders` (`usePOSOffline.ts:28-31`), replay on `online` event; SW `pos-sw.js` registered with scope `/{locale}/dashboard/pos` (`pos/page.tsx:402`) in production only.

## Findings
FE-DASH-002 (duplicate/unpaid on replay), FE-DASH-003 (SW sync wrong origin), FE-DASH-004 (SW caches auth API GETs), FE-DASH-017 (hard-coded tax, no tender), FE-DASH-024 (public order route trusts tenant headers), FE-DASH-005 (pos-storage survives logout).

## Unverified
Runtime behaviour (printing, SW registration in prod build), server-side stock decrement on POS orders, whether POS customer data is PII-purged.
