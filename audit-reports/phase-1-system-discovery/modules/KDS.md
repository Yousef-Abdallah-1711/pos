# KDS module (dashboard) - static audit

Status: **Verified implemented (code present), with reliability/security gaps**; runtime not exercised.

## Files
`dashboard/kds/page.tsx` (315 lines), `kds/settings/page.tsx` (495, uses axios directly), `kds/layout.tsx`; components `kds/components/*` (KDSGrid, KDSOrderCard, KDSHeader, KDSTimer, KDSAnalyticsPanel, KDSErrorBoundary; 3 test files). Hooks: `useKDS.ts`, `useKDSWebSocket.ts`, `useKDSOffline.ts`, `useKDSAudio.ts`, `useKDSAnalytics.ts`; store `stores/kdsStore.ts` (in-memory, not persisted).

## Gating
`useHasPermission(MODULE_IDS.KDS ...)` (bump needs WRITE). Client role fallback grants KDS=ADMIN for privileged roles read from localStorage (`useRbac.ts:100-157`) - UI only. Backend: `wrapKDSWithPermission` module 13 on REST (main.go:631-640), rate limits (bump 5/s, 60/min).

## API / realtime contract
REST (all match main.go): GET /kds/queue, /kds/stats, /kds/analytics, /kds/settings; PUT /kds/settings; PUT /kds/orders/{id}/bump, /recall. WebSocket `GET /api/v1/kds/ws?token=<JWT>` (main.go:1318 -> `internal/handler/ws/kds_ws_handler.go`). Client message types handled: `order_new|order_updated|order_bumped|order_recalled` (data = order or array) (`useKDSWebSocket.ts:108-122`); newline-delimited batches split on `\n`.

## Offline
- IndexedDB `kds-offline-db/pending-actions` stores bump/recall actions (`useKDSOffline.ts:19-21`), replayed when online (`useKDS.ts:495-545`) with `withRetry`.
- localStorage snapshot `kds-bootstrap-snapshot-v1` (orders+stats+settings, max-age constant) for fast boot (`useKDS.ts:58,166-190`).
- No KDS service worker; `/pos-sw.js` is POS-scoped only.
- REST polling fallback runs only while WS disconnected and after a delay (`useKDS.ts:446-480`).

## Findings
FE-DASH-019 (no client resync/stale detection, server WS lacks RBAC and CheckOrigin=true, browser-side auto-bump, snapshot persists after logout), FE-DASH-006 (JWT in WS query and logged at `useKDSWebSocket.ts:132`), FE-DASH-007 (impersonation sends dummy token to WS), FE-DASH-005.

## Positive
Exponential reconnect, Strict-Mode safe socket close, audio alert on `order_new`, error boundary, pending-order set prevents double bump in UI.

## Unverified
Event ordering/loss under reconnect, multi-station behaviour, backend hub tenant isolation (`buildKDSConnectionKey` is per tenant+restaurant, kds_ws_handler.go:34).
