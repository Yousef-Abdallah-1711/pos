# Offline and synchronization (dashboard) - static audit

## Mechanisms present
| Mechanism | Storage | Entry | Replay |
|---|---|---|---|
| POS pending orders | IndexedDB `pos-offline-db` v1 store `pending-orders` (`usePOSOffline.ts:28-31`) | `enqueueOrder` when `navigator.onLine` false at submit (`pos/page.tsx:288-299`) | in-page: on `online` event/mount (`usePOSOffline.ts:159-198,233-251`); SW Background Sync tag `pos-order-sync` (`pos-sw.js:231+`) |
| POS product cache | IndexedDB via `usePOSProductCache.ts` | on data load | read when offline |
| POS cart/held orders | localStorage `pos-storage` (custom storage in `posStore.ts`) | always | n/a (local only) |
| KDS pending actions | IndexedDB `kds-offline-db` store `pending-actions` | bump/recall while offline (`useKDS.ts:615,652`) | on online (`useKDS.ts:495-545`) |
| KDS snapshot | localStorage `kds-bootstrap-snapshot-v1` | after each fetch | boot only, max-age limited |
| App SW `public/sw.js` | Cache API | never registered (ServiceWorkerProvider unused) | n/a |
| POS SW `public/pos-sw.js` | Cache API `pos-cache-v1` | registered by POS page in production only, scope `/{locale}/dashboard/pos` (`pos/page.tsx:396-405`) | - |

## Conflict and integrity analysis
1. **No idempotency**: queue items carry no client request id (`POSQueuedOrder {id,payload,queuedAt,retryCount}`); server create has no dedupe key (`CreatePOSOrderPayload` lacks one). Two replay owners exist (hook + SW sync) - if the SW ever succeeded it deletes the item and notifies the page, but race windows exist (page replay reading the queue before SW delete) (FE-DASH-002/003).
2. **Replay breaks on first failure** (`usePOSOffline.ts:187`): later orders blocked behind a poison item; `retryCount` increments forever, no max/dead-letter and no user-visible queue management beyond count (`clearQueue` exists).
3. **Price/stock staleness**: offline orders are priced by the server at replay time, but the offline cart/receipt used cached prices and 14% tax (FE-DASH-017); out-of-stock at replay surfaces only as a failed replay.
4. **Payment**: offline cash orders are marked paid at replay by the hook only (SW path never marks paid).
5. **SW origin bug**: SW posts to relative `/api/v1/public/orders` without auth/tenant headers (FE-DASH-003).
6. **Cache isolation**: SW caches authenticated GETs by URL only, no logout purge (FE-DASH-004). KDS snapshot, POS storage and IndexedDB queues are not cleared on logout (FE-DASH-005) - on shared tills the next user can see and replay the previous user's queued orders under a different token.
7. **KDS actions** replay uses `withRetry` and re-fetches after success (good); bump/recall are idempotent-ish by order id, server rate limit 5/s may reject bursts.
8. `sw.js` (unused) would fail `cache.addAll(['/layout.js','/styles.css',...])` if ever registered (FE-DASH-022).
9. No offline support for reports, accounting, HR, theme-builder (not expected). `public/offline.html` exists but only referenced by unused `sw.js`.

## Verdict
POS offline: Partial (queue + in-page replay implemented; SW background sync Stub/broken; integrity gaps). KDS offline: Verified implemented (queue + snapshot), with privacy gap on logout.

## Unverified
Actual online/offline transition behaviour in browsers; IndexedDB quota handling (code resolves null on failure and treats as not queued, `usePOSOffline.ts:61-63`).
