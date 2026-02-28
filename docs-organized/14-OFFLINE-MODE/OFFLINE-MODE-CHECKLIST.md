# Offline Mode Implementation Checklist
## Dashboard (:3002) - Complete Task List with Exact File Paths

**Status**: Ready for Implementation
**Total Phases**: 5 | **Total Tasks**: 58
**Last Updated**: 2026-02-24

---

## Phase 1: Foundation & Core Infrastructure (Week 1-2)

### 1.1 Install Dependencies

#### Task 1.1.1: Add `idb` Package
- [ ] Run `pnpm add idb` in `frontend/apps/dashboard/`
- [ ] Verify `idb` appears in `frontend/apps/dashboard/package.json` dependencies
- [ ] Run `pnpm install` from `frontend/` root to update lockfile

**File**: `frontend/apps/dashboard/package.json`
**Acceptance**: `idb` in dependencies, no install errors

---

### 1.2 IndexedDB Database Setup

#### Task 1.2.1: Create `offline-db.ts`
- [ ] Create file: `frontend/apps/dashboard/src/lib/offline-db.ts`
- [ ] Use `idb` library (promise-based IndexedDB wrapper)
- [ ] Define database name: `pos-dashboard-offline-db`, version 1
- [ ] Create object stores with indexes:
  - [ ] `orders` — keyPath: `id`, indexes: `restaurantId`, `tenantId`, `status`, `createdAt`
  - [ ] `order_items` — keyPath: `id`, indexes: `orderId`, `restaurantId`
  - [ ] `products` — keyPath: `id`, indexes: `restaurantId`, `status`
  - [ ] `categories` — keyPath: `id`, indexes: `restaurantId`
  - [ ] `tables` — keyPath: `id`, indexes: `restaurantId`, `status`
  - [ ] `kds_queue` — keyPath: `id`, indexes: `restaurantId`, `status`, `createdAt`
  - [ ] `employees` — keyPath: `id`, indexes: `restaurantId`
  - [ ] `drivers` — keyPath: `id`, indexes: `restaurantId`, `status`
  - [ ] `sync_metadata` — keyPath: `moduleName`, indexes: `restaurantId`, `lastSyncAt`
  - [ ] `offline_queue` — keyPath: `id`, indexes: `module`, `createdAt`, `synced`
  - [ ] `cached_permissions` — keyPath: `userId`
- [ ] Implement CRUD: `getAll(store, index?, value?)`, `getById(store, id)`, `put(store, data)`, `putAll(store, items[])`, `delete(store, id)`, `clear(store)`
- [ ] Implement `checkStorageQuota()` — returns `{ used, quota, percentUsed }`
- [ ] Implement `cleanupOldData()` — cascading cleanup (30 days → 14 days → 1 day)
- [ ] All records must include `tenantId` and `restaurantId` from auth store
- [ ] Follow pattern from existing `src/hooks/useKDSOffline.ts` (graceful null returns if IndexedDB unavailable)

**File**: `frontend/apps/dashboard/src/lib/offline-db.ts`
**Pattern reference**: `frontend/apps/dashboard/src/hooks/useKDSOffline.ts` (lines 23-57: `openQueueDB()`)
**Acceptance**: All stores created, CRUD works, quota detection working, cleanup tested

---

### 1.3 Offline Zustand Store

#### Task 1.3.1: Create `offlineStore.ts`
- [ ] Create file: `frontend/apps/dashboard/src/stores/offlineStore.ts`
- [ ] Follow Zustand persist pattern from `src/stores/authStore.ts`
- [ ] Use persist key: `offline-storage`
- [ ] State:
  - [ ] `isOnline: boolean`
  - [ ] `isSyncing: boolean`
  - [ ] `lastSyncTime: number | null`
  - [ ] `pendingChangesCount: number`
  - [ ] `syncProgress: { current: number; total: number; module: string } | null`
- [ ] Actions:
  - [ ] `setOnline(value: boolean)`
  - [ ] `setSyncing(value: boolean)`
  - [ ] `setLastSyncTime(timestamp: number)`
  - [ ] `setPendingCount(count: number)`
  - [ ] `setSyncProgress(progress)`
- [ ] Use `partialize` to only persist `lastSyncTime` (like `src/stores/kdsStore.ts` pattern)

**File**: `frontend/apps/dashboard/src/stores/offlineStore.ts`
**Pattern reference**: `frontend/apps/dashboard/src/stores/kdsStore.ts` (partialize pattern)
**Acceptance**: Store works, persists `lastSyncTime`, no re-render loops

---

### 1.4 Offline Provider Context

#### Task 1.4.1: Create `OfflineProvider.tsx`
- [ ] Create file: `frontend/apps/dashboard/src/providers/OfflineProvider.tsx`
- [ ] Create React Context with `useOfflineContext()` hook
- [ ] Listens to `window.addEventListener('online'/'offline')` events
- [ ] Periodically pings `GET /api/v1/health` every 30 seconds (heartbeat)
- [ ] On connection restored: triggers GlobalSyncManager (2-second delay)
- [ ] Provides: `{ isOnline, isSyncing, lastSyncTime, pendingChangesCount, syncProgress, triggerManualSync, clearOfflineData }`
- [ ] Uses `offlineStore` (Zustand) internally

**File**: `frontend/apps/dashboard/src/providers/OfflineProvider.tsx`

#### Task 1.4.2: Create `useOfflineContext.ts` Hook
- [ ] Create file: `frontend/apps/dashboard/src/hooks/useOfflineContext.ts`
- [ ] Simple context consumer hook that throws if used outside OfflineProvider

**File**: `frontend/apps/dashboard/src/hooks/useOfflineContext.ts`
**Acceptance**: Hook accessible from any dashboard component

---

### 1.5 Mount Providers in Layout

#### Task 1.5.1: Update `[locale]/layout.tsx`
- [ ] Open: `frontend/apps/dashboard/src/app/[locale]/layout.tsx`
- [ ] Import `ServiceWorkerProvider` from `@/components/providers/ServiceWorkerProvider`
- [ ] Import `OfflineProvider` from `@/providers/OfflineProvider`
- [ ] Wrap inside ThemeProvider:
  ```
  AuthProvider → ThemeProvider → ServiceWorkerProvider → OfflineProvider → [rest]
  ```
- [ ] **Important**: `ServiceWorkerProvider` is currently NOT mounted (lines 39-48) — must add it
- [ ] Ensure `<Toaster>` from sonner stays at bottom (used for sync notifications)

**File**: `frontend/apps/dashboard/src/app/[locale]/layout.tsx` (currently lines 39-48)
**Acceptance**: App loads without errors, `useOfflineContext()` works from any page

#### Task 1.5.2: Add PWA Meta Tags to Root Layout
- [ ] Open: `frontend/apps/dashboard/src/app/layout.tsx`
- [ ] Add inside `<head>`:
  ```html
  <link rel="manifest" href="/manifest.json" />
  <meta name="theme-color" content="#3B82F6" />
  ```

**File**: `frontend/apps/dashboard/src/app/layout.tsx`

---

### 1.6 Service Worker File

#### Task 1.6.1: Create `sw.js`
- [ ] Create file: `frontend/apps/dashboard/public/sw.js`
- [ ] Implement cache strategies:
  - [ ] Network-first for `/api/*` (try network, fallback cache)
  - [ ] Cache-first for `/_next/static/*` (JS/CSS bundles)
  - [ ] Network-first for HTML pages
- [ ] Pre-cache critical routes on install: `/`, `/en/dashboard`, `/ar/dashboard`
- [ ] Handle messages: `SKIP_WAITING`, `CLEAR_CACHE` (already sent by existing `useServiceWorker.ts`)
- [ ] Version-based cache name (`pos-dashboard-v1`)
- [ ] Clean up old caches on activate

**File**: `frontend/apps/dashboard/public/sw.js`
**Integration**: Already registered by `src/hooks/useServiceWorker.ts` at path `/sw.js`
**Acceptance**: SW registers, assets cached, API falls back to cache offline

#### Task 1.6.2: Create PWA Manifest
- [ ] Create file: `frontend/apps/dashboard/public/manifest.json`
- [ ] Set `start_url: "/en/dashboard"`, `display: "standalone"`
- [ ] Add icons (192x192, 512x512)

**File**: `frontend/apps/dashboard/public/manifest.json`

#### Task 1.6.3: Create PWA Icons
- [ ] Create: `frontend/apps/dashboard/public/icons/icon-192.png`
- [ ] Create: `frontend/apps/dashboard/public/icons/icon-512.png`

---

### 1.7 RBAC Offline Extension

#### Task 1.7.1: Extend Permission Cache TTL
- [ ] Open: `frontend/apps/dashboard/src/lib/rbac.ts`
- [ ] In `PermissionCache` class (line 294):
  - [ ] Add method `getOfflineTTL()` that returns 24 hours when `!navigator.onLine`
  - [ ] Modify `isExpired()` (line 372) to use longer TTL when offline
  - [ ] Add `forceRefreshPermissions()` method for use on reconnect

**File**: `frontend/apps/dashboard/src/lib/rbac.ts` (class at lines 294-417)
**Acceptance**: Permissions work offline for 24h, refresh on reconnect

#### Task 1.7.2: Update `useRbac.ts` for Offline
- [ ] Open: `frontend/apps/dashboard/src/hooks/useRbac.ts`
- [ ] In `useHasPermission()`: when `!navigator.onLine`, skip API call entirely, return cached value
- [ ] Current code already falls back to cache on network error — make offline skip explicit

**File**: `frontend/apps/dashboard/src/hooks/useRbac.ts`
**Acceptance**: No API calls when offline, permissions still work from cache

---

### 1.8 Global Status Bar UI

#### Task 1.8.1: Create GlobalStatusBar Component
- [ ] Create file: `frontend/apps/dashboard/src/components/offline/GlobalStatusBar.tsx`
- [ ] Use `useOfflineContext()` hook for state
- [ ] Display states:
  - [ ] Offline: Red bar at top "Offline Mode - {pendingCount} pending changes"
  - [ ] Syncing: Blue bar with progress "Syncing... {current}/{total}"
  - [ ] Back online: Green bar "Back online" (auto-dismiss 3s)
- [ ] Support RTL (check `locale === 'ar'`)
- [ ] Use translation keys from `offline.*` namespace

**File**: `frontend/apps/dashboard/src/components/offline/GlobalStatusBar.tsx`
**Note**: Existing `ServiceWorkerProvider.tsx` already has `OfflineIndicator` (line 113) — we replace it with this richer version

#### Task 1.8.2: Mount in DashboardLayout
- [ ] Open: `frontend/apps/dashboard/src/components/layout/DashboardLayout.tsx`
- [ ] Import and render `<GlobalStatusBar />` before the closing `</div>` (after line 303)

**File**: `frontend/apps/dashboard/src/components/layout/DashboardLayout.tsx` (line ~303)
**Acceptance**: Status bar visible on all dashboard pages

---

### 1.9 i18n Translation Keys

#### Task 1.9.1: Add English Translations
- [ ] Open: `frontend/apps/dashboard/src/i18n/messages/en.json`
- [ ] Add `"offline": { ... }` section with all keys:
  - `statusOffline`, `statusSyncing`, `statusSynced`, `statusSyncFailed`
  - `storageFull`, `lastSynced`, `pendingChanges`, `manualSync`, `clearData`
  - `readOnly`, `unavailable`, `queuedForSync`

**File**: `frontend/apps/dashboard/src/i18n/messages/en.json`

#### Task 1.9.2: Add Arabic Translations
- [ ] Open: `frontend/apps/dashboard/src/i18n/messages/ar.json`
- [ ] Add same `"offline": { ... }` section with Arabic translations

**File**: `frontend/apps/dashboard/src/i18n/messages/ar.json`

---

## Phase 2: Sync Manager & Backend (Week 2-3)

### 2.1 Sync Manager

#### Task 2.1.1: Create `sync-manager.ts`
- [ ] Create file: `frontend/apps/dashboard/src/lib/sync-manager.ts`
- [ ] Define `SyncModule` interface (name, priority, batchSize, endpoints, indexedDBStore)
- [ ] Define module list with priorities (orders=1, kds=2, tables=3, products=4, employees=5)
- [ ] Implement `GlobalSyncManager` class:
  - [ ] `startSync()` — orchestrate full sync (writes first, then reads)
  - [ ] `syncWrites()` — read offline_queue from IndexedDB, batch POST to `/api/v1/sync/batch`
  - [ ] `refreshReads()` — GET fresh data for each module, update IndexedDB
  - [ ] `handleConflict(local, server)` — always return server data (server-wins)
  - [ ] Retry logic: exponential backoff (2s, 4s, 8s, max 3 retries per item)
- [ ] Emit events via offlineStore (setSyncing, setSyncProgress, setLastSyncTime)
- [ ] Use `toast()` from sonner for user notifications
- [ ] Export singleton: `export const syncManager = new GlobalSyncManager()`

**File**: `frontend/apps/dashboard/src/lib/sync-manager.ts`
**Acceptance**: Syncs all modules, handles errors gracefully, retries, shows toast

---

#### Task 2.1.2: Create `offline-api-layer.ts`
- [ ] Create file: `frontend/apps/dashboard/src/lib/offline-api-layer.ts`
- [ ] Import `api` from `./api` (existing axios instance)
- [ ] Import `offlineDB` from `./offline-db`
- [ ] Implement:
  - [ ] `offlineAwareGet(url, storeName, options?)` — online: fetch + cache, offline: read from IndexedDB
  - [ ] `offlineAwarePost(url, data, module)` — online: post normally, offline: queue in offline_queue
  - [ ] `offlineAwarePut(url, data, module, entityId)` — same as post but for updates
- [ ] All functions read tenant context from `localStorage['auth-storage']` (same as `src/lib/api.ts` interceptor)

**File**: `frontend/apps/dashboard/src/lib/offline-api-layer.ts`
**Integration**: Wraps existing `src/lib/api.ts` axios instance
**Acceptance**: Seamlessly switches between online API and offline cache

---

### 2.2 Backend Batch Sync Endpoint

#### Task 2.2.1: Create Domain Types
- [ ] Create file: `backend/internal/domain/sync.go`
- [ ] Define structs:
  ```go
  type SyncRequest struct { Action, Module, OfflineID string; EntityID int64; Data json.RawMessage; Timestamp int64 }
  type SyncResponse struct { Processed, Failed, Conflicts []SyncItem }
  type SyncItem struct { OfflineID, ServerID string; Status string; Error string }
  ```

**File**: `backend/internal/domain/sync.go`

#### Task 2.2.2: Create Sync Handler
- [ ] Create file: `backend/internal/handler/http/sync_handler.go`
- [ ] Implement `BatchSync(w, r)`:
  - [ ] Parse `[]SyncRequest` from body
  - [ ] For each: check RBAC permission (map module name → MODULE_ID)
  - [ ] For each: compare timestamps, apply or return conflict
  - [ ] Return `SyncResponse` with results per item
  - [ ] Max request size: 10MB
  - [ ] All changes in single DB transaction
- [ ] Inject required repos (orderRepo, kdsRepo, tableRepo, etc.)

**File**: `backend/internal/handler/http/sync_handler.go`
**Pattern reference**: `backend/internal/handler/http/kds_handler.go` (handler structure)
**Acceptance**: Endpoint processes batch, returns clear success/failure per item

#### Task 2.2.3: Register Route in Main
- [ ] Open: `backend/cmd/api/main.go`
- [ ] Initialize sync handler with required repos
- [ ] Add route: `mux.Handle("POST /api/v1/sync/batch", wrapAuthOnly(http.HandlerFunc(syncHandler.BatchSync)))`
- [ ] Use `wrapAuthOnly` (auth-only, no tenant middleware — reads tenant from JWT claims)

**File**: `backend/cmd/api/main.go`
**Pattern reference**: Lines where `wrapAuthOnly` is used for RBAC routes
**Acceptance**: `POST /api/v1/sync/batch` returns 200 with valid auth token

---

### 2.3 Integrate KDS Offline with Sync Manager

#### Task 2.3.1: Update `useKDSOffline.ts`
- [ ] Open: `frontend/apps/dashboard/src/hooks/useKDSOffline.ts`
- [ ] In `replayQueuedActions()` (line 191): integrate with syncManager instead of standalone replay
- [ ] On reconnect: delegate to `syncManager.startSync()` instead of replaying individually

**File**: `frontend/apps/dashboard/src/hooks/useKDSOffline.ts` (line 191)
**Acceptance**: KDS offline actions sync via global sync manager

---

## Phase 3: Orders + KDS Offline (Week 3-4)

### 3.1 Orders Offline Hooks

#### Task 3.1.1: Create `useOfflineOrders.ts`
- [ ] Create file: `frontend/apps/dashboard/src/hooks/useOfflineOrders.ts`
- [ ] Import `offlineAwareGet` from `@/lib/offline-api-layer`
- [ ] Fetch orders: online → `GET /api/v1/admin/orders` + cache, offline → IndexedDB
- [ ] Support filtering: status, dateRange, search (client-side on cached data)
- [ ] Support pagination (50 items per page, client-side)
- [ ] Return: `{ orders, isLoading, isOffline, lastSyncTime, refetch }`
- [ ] Filter by `restaurantId` from auth store

**File**: `frontend/apps/dashboard/src/hooks/useOfflineOrders.ts`
**Pattern reference**: `frontend/apps/dashboard/src/hooks/useKDS.ts` (data fetching hook pattern)

#### Task 3.1.2: Create `useCreateOfflineOrder.ts`
- [ ] Create file: `frontend/apps/dashboard/src/hooks/useCreateOfflineOrder.ts`
- [ ] Online: `POST /api/v1/admin/orders` as normal
- [ ] Offline: save to IndexedDB `orders` + `offline_queue`
- [ ] Generate temp ID: `crypto.randomUUID()`
- [ ] Return: `{ createOrder, isOffline }`

**File**: `frontend/apps/dashboard/src/hooks/useCreateOfflineOrder.ts`

---

### 3.2 Orders Pages Offline Support

#### Task 3.2.1: Update Orders List Page
- [ ] Open: `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/page.tsx`
- [ ] Import `useOfflineContext` and `useOfflineOrders`
- [ ] When offline: use `useOfflineOrders()` instead of direct API calls
- [ ] Show offline indicator badge in page header
- [ ] Show "Last synced: {time}" below header
- [ ] Disable "New Order" button when offline (or show "queued for sync" mode)

**File**: `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/page.tsx`

#### Task 3.2.2: Update Order Detail Page
- [ ] Open: `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/[id]/page.tsx`
- [ ] When offline: read order from IndexedDB by ID
- [ ] Show offline indicator
- [ ] Status update: queue in offline_queue, show "pending sync" badge

**File**: `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/[id]/page.tsx`

---

### 3.3 KDS Offline Integration

#### Task 3.3.1: Update KDS Data Hook
- [ ] Open: `frontend/apps/dashboard/src/hooks/useKDS.ts`
- [ ] Add fallback: when offline, read from IndexedDB `kds_queue` store
- [ ] Show cached KDS queue with "Offline" badge

**File**: `frontend/apps/dashboard/src/hooks/useKDS.ts`

#### Task 3.3.2: Update KDS Page
- [ ] Open: `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/page.tsx`
- [ ] Show offline indicator when disconnected
- [ ] Bump/recall actions still work (already in `useKDSOffline.ts`)

**File**: `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/page.tsx`

---

### 3.4 Layout Integration

#### Task 3.4.1: Update Sidebar with Offline Info
- [ ] Open: `frontend/apps/dashboard/src/components/layout/Sidebar.tsx`
- [ ] When offline: show section at bottom:
  - [ ] "Offline Mode" label
  - [ ] "Last sync: {time ago}"
  - [ ] "{N} pending changes"
  - [ ] "Sync Now" button (calls `triggerManualSync()`)
- [ ] Use `useOfflineContext()` hook

**File**: `frontend/apps/dashboard/src/components/layout/Sidebar.tsx` (above `<SidebarSubscriptionBadge />` at line 284)

---

## Phase 4: Remaining Modules Offline (Week 4)

### 4.1 Products & Categories

#### Task 4.1.1: Create `useOfflineProducts.ts`
- [ ] Create: `frontend/apps/dashboard/src/hooks/useOfflineProducts.ts`
- [ ] Online: `GET /api/v1/products` + cache to IndexedDB
- [ ] Offline: read from IndexedDB `products` store
- [ ] Return: `{ products, isOffline, isLoading }`

**File**: `frontend/apps/dashboard/src/hooks/useOfflineProducts.ts`

#### Task 4.1.2: Create `useOfflineCategories.ts`
- [ ] Create: `frontend/apps/dashboard/src/hooks/useOfflineCategories.ts`
- [ ] Same pattern as products

**File**: `frontend/apps/dashboard/src/hooks/useOfflineCategories.ts`

#### Task 4.1.3: Update Products Page
- [ ] Open: `frontend/apps/dashboard/src/app/[locale]/dashboard/products/page.tsx`
- [ ] Use offline products hook, show read-only when offline
- [ ] Disable "Add Product", "Edit", "Delete" buttons when offline

**File**: `frontend/apps/dashboard/src/app/[locale]/dashboard/products/page.tsx`

#### Task 4.1.4: Update Categories Page
- [ ] Open: `frontend/apps/dashboard/src/app/[locale]/dashboard/categories/page.tsx`
- [ ] Same as products — read-only offline

**File**: `frontend/apps/dashboard/src/app/[locale]/dashboard/categories/page.tsx`

---

### 4.2 Table Management

#### Task 4.2.1: Create `useOfflineTables.ts`
- [ ] Create: `frontend/apps/dashboard/src/hooks/useOfflineTables.ts`
- [ ] Cache tables in IndexedDB, queue status changes

**File**: `frontend/apps/dashboard/src/hooks/useOfflineTables.ts`

---

### 4.3 HR Module (Read-Only Offline)

#### Task 4.3.1: Create `useOfflineEmployees.ts`
- [ ] Create: `frontend/apps/dashboard/src/hooks/useOfflineEmployees.ts`
- [ ] Cache employees, read-only offline

**File**: `frontend/apps/dashboard/src/hooks/useOfflineEmployees.ts`

#### Task 4.3.2: Update HR Pages (Read-Only)
- [ ] Update: `src/app/[locale]/dashboard/hr/employees/page.tsx` — read-only offline
- [ ] Update: `src/app/[locale]/dashboard/hr/attendance/page.tsx` — read-only offline
- [ ] Update: `src/app/[locale]/dashboard/hr/leaves/page.tsx` — read-only offline
- [ ] Update: `src/app/[locale]/dashboard/hr/payroll/page.tsx` — read-only offline

**Files**: 4 HR page files
**Acceptance**: All HR pages show cached data offline, editing disabled

---

### 4.4 Drivers Page

#### Task 4.4.1: Update Drivers Page
- [ ] Update: `src/app/[locale]/dashboard/drivers/page.tsx` — read-only offline

**File**: `frontend/apps/dashboard/src/app/[locale]/dashboard/drivers/page.tsx`

---

### 4.5 Reports (Unavailable Offline)

#### Task 4.5.1: Add Offline Message to Report Pages
- [ ] Update: `src/app/[locale]/dashboard/reports/overview/page.tsx`
- [ ] Update: `src/app/[locale]/dashboard/reports/sales/page.tsx`
- [ ] Update: `src/app/[locale]/dashboard/reports/profitability/page.tsx`
- [ ] Update: `src/app/[locale]/dashboard/reports/inventory/page.tsx`
- [ ] When offline: show "Reports require internet connection" with icon

**Files**: 4 report page files
**Acceptance**: Clear offline message on all report pages

---

### 4.6 Dashboard Overview

#### Task 4.6.1: Update Dashboard Home
- [ ] Open: `src/app/[locale]/dashboard/page.tsx`
- [ ] When offline: show cached stats with "Last updated: {time}" label
- [ ] Cache dashboard stats in sync_metadata store

**File**: `frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx`

---

### 4.7 Settings Page

#### Task 4.7.1: Add Offline Controls to Settings
- [ ] Open: `src/app/[locale]/dashboard/settings/page.tsx`
- [ ] Add "Offline Data" section:
  - [ ] Storage used: "{X} MB / {Y} MB"
  - [ ] Last sync: "{datetime}"
  - [ ] Pending changes: "{N}"
  - [ ] Button: "Sync Now" (calls `triggerManualSync()`)
  - [ ] Button: "Clear Offline Data" (calls `clearOfflineData()` with confirmation)
  - [ ] Button: "Refresh Product Catalog" (invalidates products cache)

**File**: `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/page.tsx`

---

### 4.8 Feature Flags

#### Task 4.8.1: Create Feature Flag File
- [ ] Create: `frontend/apps/dashboard/src/lib/feature-flags.ts`
- [ ] Define:
  ```typescript
  export const OFFLINE_FLAGS = {
    enabled: true,
    modules: {
      orders: true,
      kds: true,
      tables: true,
      products: true,
      categories: true,
      employees: true,
      drivers: true,
      reports: false,    // never offline
      themes: false,     // never offline
      billing: false,    // never offline
    }
  }
  ```
- [ ] Check in each offline hook: `if (!OFFLINE_FLAGS.modules[moduleName]) return onlineOnlyFallback()`

**File**: `frontend/apps/dashboard/src/lib/feature-flags.ts`

---

### 4.9 Next.js Config Updates

#### Task 4.9.1: Update `next.config.js`
- [ ] Open: `frontend/apps/dashboard/next.config.js`
- [ ] Add header for `/manifest.json`: `Cache-Control: public, max-age=0, must-revalidate`
- [ ] Ensure `/sw.js` is not cached by Next.js (served as static file)

**File**: `frontend/apps/dashboard/next.config.js`

---

## Phase 5: Testing, Optimization & Polish (Week 5)

### 5.1 Unit Tests

#### Task 5.1.1: Test IndexedDB Operations
- [ ] Create: `frontend/apps/dashboard/src/__tests__/lib/offline-db.test.ts`
- [ ] Test: database creation with all stores
- [ ] Test: CRUD operations (put, get, getAll, delete, clear)
- [ ] Test: batch operations
- [ ] Test: index queries (by restaurantId, status)
- [ ] Test: storage quota detection
- [ ] Test: cleanup (delete records older than N days)
- [ ] Use `fake-indexeddb` package for testing

**File**: `frontend/apps/dashboard/src/__tests__/lib/offline-db.test.ts`

#### Task 5.1.2: Test Sync Manager
- [ ] Create: `frontend/apps/dashboard/src/__tests__/lib/sync-manager.test.ts`
- [ ] Test: full sync flow (writes then reads)
- [ ] Test: conflict resolution (server-wins)
- [ ] Test: retry logic (exponential backoff)
- [ ] Test: error handling (partial failure)
- [ ] Test: module prioritization
- [ ] Mock: axios, IndexedDB, toast

**File**: `frontend/apps/dashboard/src/__tests__/lib/sync-manager.test.ts`

#### Task 5.1.3: Test Offline API Layer
- [ ] Create: `frontend/apps/dashboard/src/__tests__/lib/offline-api-layer.test.ts`
- [ ] Test: online mode → calls API, caches response
- [ ] Test: offline mode → reads from IndexedDB
- [ ] Test: offline write → queues in offline_queue
- [ ] Mock: `navigator.onLine`, axios, IndexedDB

**File**: `frontend/apps/dashboard/src/__tests__/lib/offline-api-layer.test.ts`

#### Task 5.1.4: Test Offline Orders Hook
- [ ] Create: `frontend/apps/dashboard/src/__tests__/hooks/useOfflineOrders.test.ts`
- [ ] Test: returns cached orders when offline
- [ ] Test: filters work on cached data

**File**: `frontend/apps/dashboard/src/__tests__/hooks/useOfflineOrders.test.ts`

#### Task 5.1.5: Test Backend Sync Handler
- [ ] Create: `backend/internal/handler/http/sync_handler_test.go`
- [ ] Test: batch sync with valid requests
- [ ] Test: RBAC enforcement per request
- [ ] Test: conflict detection
- [ ] Test: partial failure (some succeed, some fail)

**File**: `backend/internal/handler/http/sync_handler_test.go`

---

### 5.2 Integration Testing

#### Task 5.2.1: Manual E2E Test Script
- [ ] Document manual test steps:
  1. Login to dashboard at `http://localhost:3002`
  2. Navigate to Orders page — data loads
  3. Open DevTools → Network → check "Offline"
  4. Verify orders still display (from cache)
  5. Verify "Offline Mode" indicator appears
  6. Navigate to Products → cached products display
  7. Navigate to Reports → "unavailable offline" message
  8. Create new order (offline) → "queued for sync" notification
  9. Uncheck "Offline" in DevTools
  10. Verify sync starts automatically (toast notification)
  11. Verify new order appears with server ID
  12. Verify "Back online" indicator

---

### 5.3 Browser Compatibility

#### Task 5.3.1: Test All Target Browsers
- [ ] Chrome (latest) — Desktop
- [ ] Firefox (latest) — Desktop
- [ ] Safari (latest) — Desktop
- [ ] Edge (latest) — Desktop
- [ ] Chrome (Android) — Mobile
- [ ] Safari (iOS) — Mobile

---

### 5.4 Performance

#### Task 5.4.1: Performance Benchmarks
- [ ] Measure: IndexedDB write speed (1000 orders)
- [ ] Measure: IndexedDB read speed (filtered query)
- [ ] Measure: Sync time (100 pending changes)
- [ ] Measure: Sync time (1000 pending changes)
- [ ] Measure: Memory usage during sync
- [ ] Target: Sync <30s for 1000 items, no memory leaks

---

### 5.5 Documentation

#### Task 5.5.1: User Guide
- [ ] Create: `docs-organized/14-OFFLINE-MODE/OFFLINE-MODE-USER-GUIDE.md`
- [ ] Cover: what offline mode is, how it works, limitations, troubleshooting

**File**: `docs-organized/14-OFFLINE-MODE/OFFLINE-MODE-USER-GUIDE.md`

#### Task 5.5.2: Developer Guide
- [ ] Create: `docs-organized/14-OFFLINE-MODE/OFFLINE-MODE-DEVELOPER-GUIDE.md`
- [ ] Cover: architecture, adding offline to new modules, hook APIs, testing

**File**: `docs-organized/14-OFFLINE-MODE/OFFLINE-MODE-DEVELOPER-GUIDE.md`

---

## Quick Reference: Complete File Inventory

### NEW Files (21)

| # | Path | Phase |
|---|---|---|
| 1 | `frontend/apps/dashboard/src/lib/offline-db.ts` | 1 |
| 2 | `frontend/apps/dashboard/src/lib/offline-api-layer.ts` | 2 |
| 3 | `frontend/apps/dashboard/src/lib/sync-manager.ts` | 2 |
| 4 | `frontend/apps/dashboard/src/lib/feature-flags.ts` | 4 |
| 5 | `frontend/apps/dashboard/src/stores/offlineStore.ts` | 1 |
| 6 | `frontend/apps/dashboard/src/providers/OfflineProvider.tsx` | 1 |
| 7 | `frontend/apps/dashboard/src/hooks/useOfflineContext.ts` | 1 |
| 8 | `frontend/apps/dashboard/src/hooks/useOfflineOrders.ts` | 3 |
| 9 | `frontend/apps/dashboard/src/hooks/useCreateOfflineOrder.ts` | 3 |
| 10 | `frontend/apps/dashboard/src/hooks/useOfflineProducts.ts` | 4 |
| 11 | `frontend/apps/dashboard/src/hooks/useOfflineCategories.ts` | 4 |
| 12 | `frontend/apps/dashboard/src/hooks/useOfflineTables.ts` | 4 |
| 13 | `frontend/apps/dashboard/src/hooks/useOfflineEmployees.ts` | 4 |
| 14 | `frontend/apps/dashboard/src/components/offline/GlobalStatusBar.tsx` | 1 |
| 15 | `frontend/apps/dashboard/public/sw.js` | 1 |
| 16 | `frontend/apps/dashboard/public/manifest.json` | 1 |
| 17 | `frontend/apps/dashboard/public/icons/icon-192.png` | 1 |
| 18 | `frontend/apps/dashboard/public/icons/icon-512.png` | 1 |
| 19 | `backend/internal/handler/http/sync_handler.go` | 2 |
| 20 | `backend/internal/domain/sync.go` | 2 |
| 21 | `docs-organized/14-OFFLINE-MODE/OFFLINE-MODE-USER-GUIDE.md` | 5 |

### MODIFIED Files (32)

| # | Path | What Changes | Phase |
|---|---|---|---|
| 1 | `frontend/apps/dashboard/package.json` | Add `idb` dep | 1 |
| 2 | `frontend/apps/dashboard/src/app/[locale]/layout.tsx` | Mount providers | 1 |
| 3 | `frontend/apps/dashboard/src/app/layout.tsx` | PWA meta tags | 1 |
| 4 | `frontend/apps/dashboard/src/lib/rbac.ts` | Extend cache TTL | 1 |
| 5 | `frontend/apps/dashboard/src/hooks/useRbac.ts` | Skip API when offline | 1 |
| 6 | `frontend/apps/dashboard/src/components/layout/DashboardLayout.tsx` | Add GlobalStatusBar | 1 |
| 7 | `frontend/apps/dashboard/src/i18n/messages/en.json` | Add offline.* keys | 1 |
| 8 | `frontend/apps/dashboard/src/i18n/messages/ar.json` | Add offline.* keys | 1 |
| 9 | `frontend/apps/dashboard/src/hooks/useKDSOffline.ts` | Integrate with sync manager | 2 |
| 10 | `backend/cmd/api/main.go` | Register sync route | 2 |
| 11 | `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/page.tsx` | Offline fallback | 3 |
| 12 | `frontend/apps/dashboard/src/app/[locale]/dashboard/orders/[id]/page.tsx` | Offline detail | 3 |
| 13 | `frontend/apps/dashboard/src/hooks/useKDS.ts` | Cache fallback | 3 |
| 14 | `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/page.tsx` | Offline indicator | 3 |
| 15 | `frontend/apps/dashboard/src/components/layout/Sidebar.tsx` | Offline info section | 3 |
| 16 | `frontend/apps/dashboard/src/app/[locale]/dashboard/products/page.tsx` | Read-only offline | 4 |
| 17 | `frontend/apps/dashboard/src/app/[locale]/dashboard/categories/page.tsx` | Read-only offline | 4 |
| 18 | `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/employees/page.tsx` | Read-only offline | 4 |
| 19 | `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/attendance/page.tsx` | Read-only offline | 4 |
| 20 | `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/leaves/page.tsx` | Read-only offline | 4 |
| 21 | `frontend/apps/dashboard/src/app/[locale]/dashboard/hr/payroll/page.tsx` | Read-only offline | 4 |
| 22 | `frontend/apps/dashboard/src/app/[locale]/dashboard/drivers/page.tsx` | Read-only offline | 4 |
| 23 | `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/overview/page.tsx` | Unavailable msg | 4 |
| 24 | `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/sales/page.tsx` | Unavailable msg | 4 |
| 25 | `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/profitability/page.tsx` | Unavailable msg | 4 |
| 26 | `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/inventory/page.tsx` | Unavailable msg | 4 |
| 27 | `frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx` | Cached stats | 4 |
| 28 | `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/page.tsx` | Offline controls | 4 |
| 29 | `frontend/apps/dashboard/src/stores/kdsStore.ts` | Integrate offline queue | 3 |
| 30 | `frontend/apps/dashboard/src/components/providers/ServiceWorkerProvider.tsx` | Enhance indicators | 1 |
| 31 | `frontend/apps/dashboard/next.config.js` | PWA headers | 4 |
| 32 | `frontend/pnpm-lock.yaml` | Updated by `pnpm add` | 1 |

### TEST Files (5)

| # | Path | Phase |
|---|---|---|
| 1 | `frontend/apps/dashboard/src/__tests__/lib/offline-db.test.ts` | 5 |
| 2 | `frontend/apps/dashboard/src/__tests__/lib/sync-manager.test.ts` | 5 |
| 3 | `frontend/apps/dashboard/src/__tests__/lib/offline-api-layer.test.ts` | 5 |
| 4 | `frontend/apps/dashboard/src/__tests__/hooks/useOfflineOrders.test.ts` | 5 |
| 5 | `backend/internal/handler/http/sync_handler_test.go` | 5 |

---

## Summary Statistics

| Metric | Count |
|---|---|
| **New files** | 21 |
| **Modified files** | 32 |
| **Test files** | 5 |
| **Total files** | **58** |
| **NPM packages to add** | 1 (`idb`) |
| **Backend endpoints** | 1 (`POST /sync/batch`) |
| **Dashboard pages affected** | 22 |
| **New hooks** | 7 |
| **New lib modules** | 4 |
| **Phases** | 5 |
| **Estimated duration** | 5 weeks |

---

## Phase Completion Criteria

### Phase 1 Done When:
- [ ] IndexedDB creates all 11 stores without errors
- [ ] OfflineProvider mounted, `useOfflineContext()` works from any page
- [ ] ServiceWorkerProvider mounted (previously missing)
- [ ] `sw.js` registers, caches assets
- [ ] RBAC permissions work offline (24h TTL)
- [ ] GlobalStatusBar shows offline/online/syncing states
- [ ] Translation keys added for EN and AR

### Phase 2 Done When:
- [ ] Sync manager syncs all module writes to backend
- [ ] `POST /api/v1/sync/batch` endpoint works
- [ ] Server-wins conflict resolution tested
- [ ] KDS offline queue integrated with sync manager

### Phase 3 Done When:
- [ ] Orders page loads from cache when offline
- [ ] Order detail page loads from cache
- [ ] New orders created offline, synced on reconnect
- [ ] KDS page works offline with cached queue
- [ ] Sidebar shows offline info section

### Phase 4 Done When:
- [ ] All 22 dashboard pages handle offline mode
- [ ] Read-only modules show "view only" indicator
- [ ] Reports show "unavailable offline" message
- [ ] Settings page has offline data controls
- [ ] Feature flags control which modules are offline-enabled

### Phase 5 Done When:
- [ ] All 5 test files passing
- [ ] Tested on Chrome, Firefox, Safari, Edge, mobile
- [ ] Sync <30s for 1000 items
- [ ] No memory leaks during sync
- [ ] User guide and developer guide written

---

## Known Risks & Mitigations

| Risk | Severity | Mitigation |
|---|---|---|
| IndexedDB quota exceeded | HIGH | Auto-cleanup cascade, user alert, manual clear option |
| Sync conflicts corrupt data | HIGH | Server-wins always, log all conflicts, no data merge |
| Service Worker serves stale cache | MEDIUM | Version-based cache names, periodic update checks (existing `useServiceWorker.ts`) |
| IndexedDB unavailable (private browsing) | MEDIUM | Graceful fallback: online-only mode (already in `useKDSOffline.ts` pattern) |
| Large sync takes too long | MEDIUM | Batch processing (50-100 per request), progress indicator, cancellable |
| Offline RBAC bypass attempt | HIGH | Validate all writes server-side in sync endpoint, permissions cached not trusted for writes |
| Auth token expires while offline | MEDIUM | On reconnect: refresh token before sync, redirect to login if expired |

