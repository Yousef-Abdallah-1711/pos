# WebSocket / realtime (dashboard side) - static audit

Scope: dashboard client only. Server pieces cited as cross-reference.

## Inventory
Only one WebSocket in the dashboard: KDS (`hooks/useKDSWebSocket.ts`). Notifications are REST-polled (`hooks/useNotifications.ts:254`: 5-minute interval + refetch on tab visible); `TenantStatusProvider` (unused) would poll tenant status; `useReports` polls on an interval (`:81`); KDS has REST polling fallback. No SSE. Orders list, POS stock, tables status have no realtime channel.

## KDS socket
- URL: derived from `apiUrl` (`http->ws`, `https->wss`) + `/kds/ws?token=<JWT>` (`:20-35,131`). Backend route `GET /api/v1/kds/ws` (main.go:1318) -> `internal/handler/ws/kds_ws_handler.go`, hub keyed per tenant+restaurant (`:34`), validates JWT from query (`:98-106`), server ping every 54s (`kds_hub.go:18`).
- Lifecycle: opens only when `token` exists; closes on logout/unmount; Strict-Mode-safe close; reconnect with exponential backoff 1s..10s, no jitter, unlimited attempts (`:163-171`).
- Messages: `order_new|order_updated|order_bumped|order_recalled` with `data` single order or array (replace-all); newline-delimited frames split (`:142-147`). Unknown types ignored. No ack/sequence numbers.
- `order_new` triggers audio alert callback.
- Fallback: REST poll of `/kds/queue` + `/kds/stats` only while disconnected and online after a delay (`useKDS.ts:446-480`).

## Issues
- FE-DASH-006: JWT in query string; full URL written by `debugLog` (`:132`).
- FE-DASH-007: under impersonation the token is the literal 'impersonated' so WS auth cannot succeed.
- FE-DASH-019: no client stale detection; no refetch when socket reconnects (events during the gap are only recovered if the poll ran); server-side no KDS permission check on the socket and `CheckOrigin` returns true (`kds_ws_handler.go:83-86`) - any authenticated tenant user can subscribe.
- Token expiry: socket is not re-established with a refreshed token (no refresh flow found; auth uses a single JWT in localStorage).
- No multi-tab dedupe: each KDS tab opens its own socket and runs its own auto-bump timer.

## Not verified at runtime
Latency, hub fan-out, message loss, behaviour behind proxies (nginx upgrade headers) - runtime auditor.
