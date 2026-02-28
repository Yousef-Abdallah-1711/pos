# Full Offline Mode Implementation Plan
## Dashboard (:3002) - Global Offline-First Architecture

**Status**: Planning Phase
**Last Updated**: 2026-02-24
**Scope**: Complete offline support for all Dashboard features across multi-tenant POS system
**Storage Limit**: 50-100MB (IndexedDB)
**Sync Strategy**: Silent background auto-sync on reconnection (Server-wins conflict resolution)

---

## 1. Architecture Overview

### 1.1 Current Codebase Snapshot (What Already Exists)

The dashboard app already has several building blocks we will reuse:

| Existing File | What It Does | How We Reuse It |
|---|---|---|
| `src/hooks/useKDSOffline.ts` | IndexedDB action queue (enqueue/replay/clear pattern) | **Generalize** into a global offline queue for all modules |
| `src/hooks/useServiceWorker.ts` | SW registration, `useOffline()`, `useBackgroundSync()`, `usePersistentStorage()` | **Use directly** — already production-ready |
| `src/components/providers/ServiceWorkerProvider.tsx` | Offline/Online indicator UI, update prompt | **Mount** in `[locale]/layout.tsx` (currently NOT mounted) |
| `src/lib/rbac.ts` → `PermissionCache` class | Persists permissions to `localStorage['rbac-permission-cache-v1']` with 5-min TTL | **Extend TTL** for offline (24h) and use as offline RBAC source |
| `src/stores/kdsStore.ts` | Zustand persist + `optimisticBump()` / `rollbackOrder()` | **Reuse pattern** for optimistic offline updates |
| `src/lib/cache-manager.ts` | LRU+TTL CacheManager (50MB limit) | **Reuse** for in-memory API response caching |
| `src/lib/api.ts` | Axios instance, auth interceptors, tenant headers | **Wrap** with offline-aware layer |
| `src/stores/authStore.ts` | Zustand persist to `localStorage['auth-storage']` | Already works offline (token/user in localStorage) |
| Toast library: `sonner` | `<Toaster>` already in `[locale]/layout.tsx` | **Use directly** for sync notifications |

### 1.2 Core Components (New)

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                     Next.js Dashboard App (:3002)                           │
│                                                                              │
│  EXISTING (reuse):                         NEW (create):                     │
│  ┌──────────────────────────┐             ┌───────────────────────────────┐  │
│  │ ServiceWorkerProvider    │             │ OfflineProvider (React ctx)   │  │
│  │ useServiceWorker()       │             │ useOfflineContext()           │  │
│  │ useOffline()             │             │ GlobalSyncManager            │  │
│  │ useBackgroundSync()      │             │ offline-db.ts (IndexedDB)    │  │
│  │ usePersistentStorage()   │             │ offline-api-layer.ts         │  │
│  │ PermissionCache (rbac)   │             │ offlineStore.ts (Zustand)    │  │
│  │ CacheManager (LRU)       │             │ GlobalStatusBar component    │  │
│  │ kdsStore optimistic()    │             │ public/sw.js (SW file)       │  │
│  │ sonner Toaster           │             │ public/manifest.json (PWA)   │  │
│  └──────────────────────────┘             └───────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │      IndexedDB ("pos-dashboard-offline-db")  — 50-100MB limit       │   │
│  │  ┌──────────┬──────────┬────────┬──────────┬──────────┬───────────┐ │   │
│  │  │ orders   │ products │ tables │ kds_queue│employees │ categories│ │   │
│  │  │ (30 day) │ (live)   │ (live) │ (1 day)  │ (live)   │ (live)    │ │   │
│  │  └──────────┴──────────┴────────┴──────────┴──────────┴───────────┘ │   │
│  │  ┌──────────────────────────────────────────────────────────────────┐│   │
│  │  │ sync_metadata  │  offline_queue (pending writes)  │ permissions ││   │
│  │  └──────────────────────────────────────────────────────────────────┘│   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
         │ HTTP / WebSocket
    ┌────▼──────────────────────────────────┐
    │  Go Backend (:8080)                   │
    │  baseURL: /api/v1                     │
    │  Health: GET /api/v1/health           │
    │  Batch sync: POST /api/v1/sync/batch  │  (NEW)
    └───────────────────────────────────────┘
```

### 1.3 Provider Chain (After Integration)

Current chain in `src/app/[locale]/layout.tsx`:
```
AuthProvider → ThemeProvider → WebVitalsProvider → PerformanceMonitor → Toaster
```

New chain:
```
AuthProvider → ThemeProvider → ServiceWorkerProvider → OfflineProvider → WebVitalsProvider → PerformanceMonitor → Toaster
```

---

## 2. Existing Code Patterns (Must Follow)

### 2.1 Zustand Store Pattern
```typescript
// From src/stores/authStore.ts — ALL stores follow this pattern:
export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({ /* state + actions */ }),
    { name: 'auth-storage' }  // localStorage key
  )
)
```
Existing localStorage keys already in use:
- `auth-storage` — auth token/user
- `kds-storage` — KDS UI prefs (uses `partialize` to only persist UI state)
- `preferences-storage` — user theme/language prefs
- `rbac-permission-cache-v1` — RBAC permissions (5-min TTL)

New key for offline store: `offline-storage`

### 2.2 IndexedDB Pattern (from `src/hooks/useKDSOffline.ts`)
```typescript
// Already production-tested pattern:
const DB_NAME = 'kds-offline-db'
const DB_VERSION = 1
const STORE_NAME = 'pending-actions'

// Functions: openQueueDB(), readAllActions(), upsertAction(), removeAction(), clearActions()
// Hook: useKDSOffline() → { isOffline, queuedActionsCount, enqueueAction, replayQueuedActions, clearQueuedActions }
```
We generalize this into `src/lib/offline-db.ts` with multiple object stores.

### 2.3 API Client Pattern (from `src/lib/api.ts`)
```typescript
// Axios instance with interceptors:
// - Request: adds Authorization Bearer, X-Tenant-ID, X-Restaurant-ID from localStorage['auth-storage']
// - Response: on 401 → clears storage → redirects to login
// - baseURL: getConfig().apiUrl (dev: http://localhost:8080/api/v1)
// - timeout: 6000ms (configurable via NEXT_PUBLIC_API_TIMEOUT_MS)
```

### 2.4 Permission Check Pattern (from every dashboard page)
```typescript
import { useHasPermission } from '@/hooks/useRbac'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
const { hasPermission, loading } = useHasPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.READ)
if (loading) return <Loader />
if (!hasPermission) return <AccessDenied />
```

### 2.5 Backend Route Registration Pattern (from `backend/cmd/api/main.go`)
```go
// Three wrappers available:
wrapAuthOnly     // Auth middleware only (reads JWT claims)
wrapProtected    // Auth + TenantContext middleware
wrapWithPermission(handler, moduleID, permLevel) // Auth + Tenant + RBAC check
```

---

## 3. MODULE_IDS (Correct Values from `src/lib/rbac.ts`)

```typescript
export const MODULE_IDS = {
  PRODUCTS: 1,
  CATEGORIES: 2,
  HR: 3,
  NOTIFICATIONS: 4,
  SETTINGS: 5,
  USERS: 6,
  ROLES: 7,
  REPORTS: 8,
  THEMES: 9,
  ORDERS: 10,
  LEAVES: 11,
  SALARIES: 12,
  KDS: 13,
  TABLES: 14,
} as const
```

---

## 4. IndexedDB Schema & Storage Management

### 4.1 Database Structure

**New file: `src/lib/offline-db.ts`**

```typescript
const DB_NAME = 'pos-dashboard-offline-db'
const DB_VERSION = 1

// Object stores with indexes:
interface OfflineDBSchema {
  orders:            { key: 'id', indexes: ['restaurantId', 'tenantId', 'status', 'createdAt'] }
  order_items:       { key: 'id', indexes: ['orderId', 'restaurantId'] }
  products:          { key: 'id', indexes: ['restaurantId', 'status'] }
  categories:        { key: 'id', indexes: ['restaurantId'] }
  tables:            { key: 'id', indexes: ['restaurantId', 'status'] }
  kds_queue:         { key: 'id', indexes: ['restaurantId', 'status', 'createdAt'] }
  employees:         { key: 'id', indexes: ['restaurantId'] }
  drivers:           { key: 'id', indexes: ['restaurantId', 'status'] }
  sync_metadata:     { key: 'moduleName', indexes: ['restaurantId', 'lastSyncAt'] }
  offline_queue:     { key: 'id', indexes: ['module', 'createdAt', 'synced'] }
  cached_permissions:{ key: 'userId', indexes: [] }
}
```

### 4.2 Retention Policy

| Store | Retention | Estimated Size |
|---|---|---|
| `orders` + `order_items` | Last 30 days | 30-40MB |
| `products` + `categories` | Live (until invalidated) | 5-8MB |
| `tables` | Live (always current) | 1-2MB |
| `kds_queue` | Current day only | 2-3MB |
| `employees` + `drivers` | Live (until invalidated) | 2-3MB |
| `sync_metadata` | Permanent (tiny) | <100KB |
| `offline_queue` | Until synced | 1-2MB |
| `cached_permissions` | Until refresh | <100KB |
| **TOTAL** | | **~50-60MB** |

### 4.3 Auto-Cleanup Strategy

```typescript
async function cleanupOldData(): Promise<void> {
  const estimate = await navigator.storage?.estimate()
  const percentUsed = ((estimate?.usage || 0) / (estimate?.quota || 100_000_000)) * 100

  if (percentUsed > 90) {
    await deleteRecordsOlderThan('orders', 30)       // Step 1
  }
  if (percentUsed > 85) {
    await deleteRecordsOlderThan('orders', 14)       // Step 2
  }
  if (percentUsed > 85) {
    await deleteRecordsOlderThan('kds_queue', 1)     // Step 3
  }
  if (percentUsed > 85) {
    toast.warning('Offline storage almost full')      // Step 4: alert user
  }
}
```

---

## 5. All Backend API Endpoints (For Offline Caching)

### 5.1 Endpoints to Cache (GET requests → IndexedDB)

| Module | Backend Endpoint | IndexedDB Store | RBAC Module ID |
|---|---|---|---|
| Orders | `GET /api/v1/admin/orders` | `orders` | 10 |
| Order Detail | `GET /api/v1/admin/orders/{id}` | `orders` (by id) | 10 |
| Order Stats | `GET /api/v1/admin/orders/stats` | `sync_metadata` | 10 |
| Products | `GET /api/v1/products` | `products` | 1 |
| Categories | `GET /api/v1/categories` | `categories` | 2 |
| KDS Queue | `GET /api/v1/kds/queue` | `kds_queue` | 13 |
| KDS Stats | `GET /api/v1/kds/stats` | `sync_metadata` | 13 |
| KDS Settings | `GET /api/v1/kds/settings` | `sync_metadata` | 13 |
| Tables | (table endpoints) | `tables` | 14 |
| Employees | `GET /api/v1/hr/employees` | `employees` | 3 |
| Attendance | `GET /api/v1/hr/attendance` | — (read-only cache) | 3 |
| Leaves | `GET /api/v1/hr/leaves` | — (read-only cache) | 11 |
| Salaries | `GET /api/v1/hr/salaries` | — (read-only cache) | 12 |
| Drivers | `GET /api/v1/admin/drivers` | `drivers` | 10 |
| Reports | `GET /api/v1/reports/*` | — (not cached offline) | 8 |
| Permissions | `GET /api/v1/rbac/me/permissions` | `cached_permissions` | — |

### 5.2 Endpoints for Write-Back (Offline Queue → POST/PUT on reconnect)

| Action | Backend Endpoint | Method |
|---|---|---|
| Create order | `POST /api/v1/admin/orders` | POST |
| Update order status | `PUT /api/v1/admin/orders/{id}/status` | PUT |
| KDS bump order | `PUT /api/v1/kds/orders/{id}/bump` | PUT |
| KDS recall order | `PUT /api/v1/kds/orders/{id}/recall` | PUT |
| Update table status | (table status endpoint) | PUT |
| **Batch sync** | `POST /api/v1/sync/batch` | POST **(NEW)** |

### 5.3 Endpoints NOT Available Offline

| Feature | Why Not Offline |
|---|---|
| `POST /api/v1/subscription/initiate` | Paymob payment requires live backend |
| `GET /api/v1/reports/*` | Requires DB aggregations |
| `POST /api/v1/auth/login` | Auth requires server verification |
| `POST /api/v1/admin/themes` | Theme creation requires server |
| `GET /api/v1/notifications` | Real-time push only |
| WebSocket `GET /api/v1/kds/ws` | Real-time only |

---

## 6. Module-by-Module Offline Strategy

### 6.1 Orders Module

**Dashboard pages affected:**
- `src/app/[locale]/dashboard/orders/page.tsx` — Order list
- `src/app/[locale]/dashboard/orders/[id]/page.tsx` — Order detail

**Offline capabilities:** View cached orders (30 days), filter/search, create new (queued), update status (queued)

**Sync strategy:** Batch POST new orders, batch PUT status changes. Server-wins on conflict.

**Files to create:**
- `src/hooks/useOfflineOrders.ts` — Read cached orders from IndexedDB
- `src/hooks/useCreateOfflineOrder.ts` — Queue new order in offline_queue

**Files to modify:**
- `src/app/[locale]/dashboard/orders/page.tsx` — Add offline data source fallback
- `src/app/[locale]/dashboard/orders/[id]/page.tsx` — Read from cache when offline

### 6.2 Drivers Module

**Dashboard page:** `src/app/[locale]/dashboard/drivers/page.tsx`

**Offline:** View cached drivers (read-only). No writes offline.

### 6.3 Products & Categories

**Dashboard pages:**
- `src/app/[locale]/dashboard/products/page.tsx`
- `src/app/[locale]/dashboard/products/new/page.tsx`
- `src/app/[locale]/dashboard/products/[id]/edit/page.tsx`
- `src/app/[locale]/dashboard/categories/page.tsx`
- `src/app/[locale]/dashboard/categories/new/page.tsx`
- `src/app/[locale]/dashboard/categories/[id]/edit/page.tsx`

**Offline:** View cached products/categories (read-only). Create/edit disabled offline.

**Files to create:**
- `src/hooks/useOfflineProducts.ts`
- `src/hooks/useOfflineCategories.ts`

### 6.4 KDS (Kitchen Display System)

**Dashboard pages:**
- `src/app/[locale]/dashboard/kds/page.tsx` — Fullscreen KDS
- `src/app/[locale]/dashboard/kds/settings/page.tsx` — KDS settings
- `src/app/[locale]/dashboard/kds/layout.tsx` — KDS layout

**Existing offline code (reuse):**
- `src/hooks/useKDSOffline.ts` — Already has IndexedDB queue for bump/recall actions
- `src/stores/kdsStore.ts` — Already has `optimisticBump()`, `optimisticRecall()`, `rollbackOrder()`
- `src/hooks/useKDSWebSocket.ts` — Handles WS disconnect/reconnect

**Offline:** View cached queue, bump/recall orders (queued via existing useKDSOffline).

**Files to modify:**
- `src/hooks/useKDS.ts` — Add fallback to IndexedDB cache when offline
- `src/hooks/useKDSOffline.ts` — Integrate with GlobalSyncManager

### 6.5 Table Management

**Offline:** View cached tables, update status (queued).

**Files to create:**
- `src/hooks/useOfflineTables.ts`

### 6.6 HR Module (Employees, Attendance, Leaves, Payroll)

**Dashboard pages:**
- `src/app/[locale]/dashboard/hr/employees/page.tsx`
- `src/app/[locale]/dashboard/hr/attendance/page.tsx`
- `src/app/[locale]/dashboard/hr/leaves/page.tsx`
- `src/app/[locale]/dashboard/hr/payroll/page.tsx`

**Offline:** All read-only. View cached employees, attendance, leaves, payroll.

**Files to create:**
- `src/hooks/useOfflineEmployees.ts`

### 6.7 Reports Module

**Dashboard pages:**
- `src/app/[locale]/dashboard/reports/overview/page.tsx`
- `src/app/[locale]/dashboard/reports/sales/page.tsx`
- `src/app/[locale]/dashboard/reports/profitability/page.tsx`
- `src/app/[locale]/dashboard/reports/inventory/page.tsx`

**Offline:** NOT supported. Reports require server-side aggregation. Show "Reports unavailable offline" message.

### 6.8 Settings Pages

**Dashboard pages:**
- `src/app/[locale]/dashboard/settings/page.tsx`
- `src/app/[locale]/dashboard/settings/account/page.tsx`
- `src/app/[locale]/dashboard/settings/appearance/page.tsx`
- `src/app/[locale]/dashboard/settings/billing/page.tsx`
- `src/app/[locale]/dashboard/settings/restaurant/page.tsx`
- `src/app/[locale]/dashboard/settings/access-control/page.tsx`
- `src/app/[locale]/dashboard/settings/roles/page.tsx`
- `src/app/[locale]/dashboard/settings/roles/new/page.tsx`
- `src/app/[locale]/dashboard/settings/roles/[id]/edit/page.tsx`
- `src/app/[locale]/dashboard/settings/users/page.tsx`

**Offline:** Read-only for appearance/account. Billing, roles, users require server. Add offline settings section for "Clear Offline Data", "Manual Sync" controls.

### 6.9 Theme Builder

**Dashboard pages:**
- `src/app/[locale]/dashboard/theme-builder/page.tsx`
- `src/app/[locale]/dashboard/theme-builder/editor/[id]/page.tsx`
- `src/app/[locale]/dashboard/theme-builder/import/page.tsx`
- `src/app/[locale]/dashboard/theme-builder/presets/page.tsx`

**Offline:** NOT supported in Phase 1. Theme editing requires server.

### 6.10 Dashboard Overview

**Dashboard page:** `src/app/[locale]/dashboard/page.tsx`

**Offline:** Show cached stats/summary. "Last updated: {timestamp}" indicator.

---

## 7. RBAC Offline Access Control

### 7.1 Current Permission System

The `PermissionCache` class in `src/lib/rbac.ts` already:
- Stores permissions per userId in `localStorage['rbac-permission-cache-v1']`
- Has 5-minute TTL (in-memory Map + localStorage fallback)
- Structure: `{ [userId]: { permissions: { [moduleId]: 'READ'|'WRITE'|... }, timestamp } }`

### 7.2 Offline Extension

**Modify `src/lib/rbac.ts`:**
- Increase TTL to 24 hours when offline (detect via `navigator.onLine`)
- Add `getOfflinePermissions(userId)` that always returns cached (no network check)
- Add `forceRefreshPermissions()` to call on reconnect

**Modify `src/hooks/useRbac.ts`:**
- In `useHasPermission()`: when offline, skip API call, use cached value directly
- Current code already falls back to cache on error — this just needs a timeout/offline skip

---

## 8. Offline-Aware API Layer

### 8.1 Wrap Existing Axios Instance

**New file: `src/lib/offline-api-layer.ts`**

```typescript
import { api } from './api'
import { offlineDB } from './offline-db'
import { offlineQueue } from './offline-db'

// Wraps existing api calls with offline detection:
export async function offlineAwareGet(url: string, storeName: string, options?: any) {
  if (navigator.onLine) {
    try {
      const response = await api.get(url, options)
      // Cache response in IndexedDB
      await offlineDB.putAll(storeName, response.data)
      return response.data
    } catch (error) {
      // Network error → fall back to cache
      return offlineDB.getAll(storeName)
    }
  }
  // Offline → read from IndexedDB
  return offlineDB.getAll(storeName)
}

export async function offlineAwarePost(url: string, data: any, module: string) {
  if (navigator.onLine) {
    return api.post(url, data)
  }
  // Offline → queue for later sync
  await offlineQueue.enqueue({ action: 'create', module, url, data, createdAt: Date.now() })
  return { data, _offline: true, _tempId: crypto.randomUUID() }
}
```

### 8.2 Integration with Existing API

The existing `src/lib/api.ts` interceptors handle:
- Auth token injection from `localStorage['auth-storage']`
- Tenant ID / Restaurant ID headers
- 401 → logout redirect

The offline layer wraps **on top** of this — it checks `navigator.onLine` before calling `api.get()` / `api.post()`.

---

## 9. Service Worker Implementation

### 9.1 Create the SW File

The existing `src/hooks/useServiceWorker.ts` already registers `/sw.js` but the file doesn't exist yet.

**New file: `public/sw.js`**

```javascript
const CACHE_NAME = 'pos-dashboard-v1'
const STATIC_ASSETS = [
  '/',
  '/en/dashboard',
  '/ar/dashboard',
]

// Install: pre-cache critical routes
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then(cache => cache.addAll(STATIC_ASSETS))
  )
  self.skipWaiting()
})

// Activate: clean old caches
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then(names =>
      Promise.all(names.filter(n => n !== CACHE_NAME).map(n => caches.delete(n)))
    )
  )
  self.clients.claim()
})

// Fetch: network-first for API, cache-first for static
self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url)

  if (url.pathname.startsWith('/api/')) {
    event.respondWith(networkFirst(event.request))
  } else if (url.pathname.startsWith('/_next/static/')) {
    event.respondWith(cacheFirst(event.request))
  } else {
    event.respondWith(networkFirst(event.request))
  }
})

// Messages from app (SKIP_WAITING, CLEAR_CACHE)
self.addEventListener('message', (event) => {
  if (event.data?.type === 'SKIP_WAITING') self.skipWaiting()
  if (event.data?.type === 'CLEAR_CACHE') {
    caches.keys().then(names => names.forEach(n => caches.delete(n)))
  }
})
```

### 9.2 Mount ServiceWorkerProvider

**Modify: `src/app/[locale]/layout.tsx`** (line ~39-48)

```diff
+ import { ServiceWorkerProvider } from '@/components/providers/ServiceWorkerProvider'

  return (
    <AuthProvider>
      <ThemeProvider attribute="class" defaultTheme="light" enableSystem>
+       <ServiceWorkerProvider>
+         <OfflineProvider>
            <WebVitalsProvider />
            <PerformanceMonitor />
            <div lang={locale} dir={dir} className={inter.className}>
              {children}
            </div>
            <Toaster position={isRTL ? 'top-left' : 'top-right'} richColors />
+         </OfflineProvider>
+       </ServiceWorkerProvider>
      </ThemeProvider>
    </AuthProvider>
  )
```

### 9.3 PWA Manifest

**New file: `public/manifest.json`**

```json
{
  "name": "POS Dashboard",
  "short_name": "Dashboard",
  "description": "Restaurant POS management dashboard",
  "start_url": "/en/dashboard",
  "display": "standalone",
  "orientation": "any",
  "theme_color": "#3B82F6",
  "background_color": "#F9FAFB",
  "icons": [
    { "src": "/icons/icon-192.png", "sizes": "192x192", "type": "image/png" },
    { "src": "/icons/icon-512.png", "sizes": "512x512", "type": "image/png" }
  ]
}
```

**Modify: `src/app/layout.tsx`** — add manifest link in `<head>`:
```html
<link rel="manifest" href="/manifest.json" />
<meta name="theme-color" content="#3B82F6" />
```

---

## 10. Sync Manager Architecture

### 10.1 Global Sync Flow

```
Connection Restored (window 'online' event)
    │
    ▼ (2 second delay for network stabilization)
    │
    ▼
Get all pending items from IndexedDB offline_queue
    │
    ▼
Group by module, sort by priority:
  1. Orders    (priority 1 - highest)
  2. KDS       (priority 2)
  3. Tables    (priority 3)
  4. Products  (priority 4 - read-only, just refresh)
  5. Employees (priority 5 - read-only, just refresh)
    │
    ▼
For each pending write:
  POST/PUT to backend → Get response
    │
    ├── Success → Remove from offline_queue, update IndexedDB
    │
    ├── Conflict → Server-wins: replace local with server data
    │
    └── Error → Keep in queue, retry later (exponential backoff)
    │
    ▼
After all writes synced:
  GET fresh data for each module → Update IndexedDB cache
    │
    ▼
Toast notification: "Data synchronized" (via sonner)
```

### 10.2 New File: `src/lib/sync-manager.ts`

```typescript
import { api } from './api'
import { offlineDB, offlineQueue } from './offline-db'
import { toast } from 'sonner'

interface SyncModule {
  name: string
  priority: number
  batchSize: number
  writeEndpoint: string     // POST /api/v1/sync/batch
  readEndpoint: string      // GET /api/v1/admin/orders etc.
  indexedDBStore: string    // 'orders', 'products', etc.
}

const SYNC_MODULES: SyncModule[] = [
  { name: 'orders',    priority: 1, batchSize: 50,  writeEndpoint: '/sync/batch', readEndpoint: '/admin/orders',  indexedDBStore: 'orders' },
  { name: 'kds',       priority: 2, batchSize: 100, writeEndpoint: '/sync/batch', readEndpoint: '/kds/queue',     indexedDBStore: 'kds_queue' },
  { name: 'tables',    priority: 3, batchSize: 50,  writeEndpoint: '/sync/batch', readEndpoint: '/tables',        indexedDBStore: 'tables' },
  { name: 'products',  priority: 4, batchSize: 100, writeEndpoint: '',            readEndpoint: '/products',      indexedDBStore: 'products' },
  { name: 'employees', priority: 5, batchSize: 100, writeEndpoint: '',            readEndpoint: '/hr/employees',  indexedDBStore: 'employees' },
]

class GlobalSyncManager {
  private syncInProgress = false

  async startSync(): Promise<void>
  async syncWrites(): Promise<{ processed: number; failed: number }>
  async refreshReads(): Promise<void>
  private handleConflict(local: any, server: any): any  // Always returns server
}

export const syncManager = new GlobalSyncManager()
```

---

## 11. Backend Batch Sync Endpoint

### 11.1 New Handler

**New file: `backend/internal/handler/http/sync_handler.go`**

```go
type SyncRequest struct {
    Action    string      `json:"action"`     // "create", "update", "delete"
    Module    string      `json:"module"`     // "orders", "kds", "tables"
    EntityID  int64       `json:"entity_id"`  // server ID (0 for new)
    OfflineID string      `json:"offline_id"` // client-generated UUID
    Data      interface{} `json:"data"`
    Timestamp int64       `json:"timestamp"`  // client modification time
}

type SyncResponse struct {
    Processed []ProcessedItem `json:"processed"`
    Failed    []FailedItem    `json:"failed"`
    Conflicts []ConflictItem  `json:"conflicts"`
}

// POST /api/v1/sync/batch
func (h *SyncHandler) BatchSync(w http.ResponseWriter, r *http.Request) {
    // 1. Parse []SyncRequest from body
    // 2. For each request:
    //    a. Check RBAC permission (module → moduleID)
    //    b. Compare timestamps for conflict detection
    //    c. Apply change or return conflict
    // 3. Return SyncResponse
}
```

### 11.2 Wire Up in Router

**Modify: `backend/cmd/api/main.go`**

```go
// Add sync handler initialization
syncHandler := http.NewSyncHandler(/* repos */)

// Add route (auth-only, no tenant middleware needed — reads from JWT)
mux.Handle("POST /api/v1/sync/batch", wrapAuthOnly(http.HandlerFunc(syncHandler.BatchSync)))
```

### 11.3 New Domain Types

**Modify: `backend/internal/domain/` — add sync types**

New file or extend existing: `backend/internal/domain/sync.go`

---

## 12. Toast Notifications (Using Existing Sonner)

The app already uses `sonner` with `<Toaster>` mounted in `src/app/[locale]/layout.tsx`:

```typescript
import { toast } from 'sonner'

// In sync manager:
toast.info('Going offline...')
toast.warning('Offline mode — some features limited')
toast.success(`Synced ${count} items`)
toast.error(`Sync failed: ${error.message}`)
```

No new component needed — just call `toast()` from sync manager.

---

## 13. i18n Translation Keys

**Modify: `src/i18n/messages/en.json`** — add offline section:

```json
{
  "offline": {
    "statusOffline": "Offline Mode - No Internet",
    "statusSyncing": "Syncing data...",
    "statusSynced": "Data synchronized",
    "statusSyncFailed": "Sync failed",
    "storageFull": "Offline storage almost full",
    "lastSynced": "Last synced",
    "pendingChanges": "Pending changes",
    "manualSync": "Sync Now",
    "clearData": "Clear Offline Data",
    "readOnly": "View only — connect to internet to edit",
    "unavailable": "This feature requires internet connection",
    "queuedForSync": "Saved offline — will sync when online"
  }
}
```

**Modify: `src/i18n/messages/ar.json`** — same keys with Arabic translations.

---

## 14. NPM Packages to Install

**Modify: `frontend/apps/dashboard/package.json`**

```bash
pnpm add idb           # Better IndexedDB API (typed, promise-based)
# Optional:
# pnpm add workbox-precaching workbox-routing workbox-strategies  # If using Workbox for SW
```

Current deps already include: `axios`, `zustand`, `sonner`, `next-themes`, `@tanstack/react-query`

**No `next-pwa` needed** — we handle SW registration manually via existing `useServiceWorker.ts`.

---

## 15. Complete File Inventory

### 15.1 New Files to Create

| # | File Path | Purpose |
|---|---|---|
| 1 | `src/lib/offline-db.ts` | IndexedDB database setup, CRUD, cleanup |
| 2 | `src/lib/offline-api-layer.ts` | Offline-aware wrapper around axios `api` |
| 3 | `src/lib/sync-manager.ts` | Global sync orchestrator |
| 4 | `src/stores/offlineStore.ts` | Zustand store for offline/sync state |
| 5 | `src/providers/OfflineProvider.tsx` | React context wrapping offline state |
| 6 | `src/hooks/useOfflineContext.ts` | Hook to access OfflineProvider context |
| 7 | `src/hooks/useOfflineOrders.ts` | Read cached orders from IndexedDB |
| 8 | `src/hooks/useCreateOfflineOrder.ts` | Queue new order for sync |
| 9 | `src/hooks/useOfflineProducts.ts` | Read cached products |
| 10 | `src/hooks/useOfflineCategories.ts` | Read cached categories |
| 11 | `src/hooks/useOfflineTables.ts` | Read/update cached tables |
| 12 | `src/hooks/useOfflineEmployees.ts` | Read cached employees |
| 13 | `src/components/offline/GlobalStatusBar.tsx` | Connection status bar + sync progress |
| 14 | `src/components/offline/OfflineIndicator.tsx` | Small badge for individual pages |
| 15 | `src/lib/feature-flags.ts` | Offline module feature flags |
| 16 | `public/sw.js` | Service Worker file |
| 17 | `public/manifest.json` | PWA manifest |
| 18 | `public/icons/icon-192.png` | PWA icon 192x192 |
| 19 | `public/icons/icon-512.png` | PWA icon 512x512 |
| 20 | `backend/internal/handler/http/sync_handler.go` | Batch sync endpoint |
| 21 | `backend/internal/domain/sync.go` | Sync domain types |

### 15.2 Existing Files to Modify

| # | File Path | What to Change |
|---|---|---|
| 1 | `src/app/[locale]/layout.tsx` | Mount `ServiceWorkerProvider` + `OfflineProvider` in provider chain |
| 2 | `src/app/layout.tsx` | Add `<link rel="manifest">` and `<meta name="theme-color">` |
| 3 | `src/components/providers/ServiceWorkerProvider.tsx` | Enhance with sync-aware status (currently basic online/offline) |
| 4 | `src/hooks/useKDSOffline.ts` | Integrate with GlobalSyncManager for queue replay |
| 5 | `src/hooks/useKDS.ts` | Add IndexedDB cache fallback when offline |
| 6 | `src/hooks/useRbac.ts` | Skip API call when offline, use cached permissions directly |
| 7 | `src/lib/rbac.ts` | Extend `PermissionCache` TTL to 24h when offline |
| 8 | `src/lib/api.ts` | (Optional) Add offline detection in response interceptor |
| 9 | `src/components/layout/DashboardLayout.tsx` | Add `<GlobalStatusBar />` before closing `</div>` |
| 10 | `src/components/layout/Sidebar.tsx` | Add offline data info section (last sync time, pending count) |
| 11 | `src/app/[locale]/dashboard/orders/page.tsx` | Add offline data source fallback |
| 12 | `src/app/[locale]/dashboard/orders/[id]/page.tsx` | Read from IndexedDB cache when offline |
| 13 | `src/app/[locale]/dashboard/products/page.tsx` | Read-only offline mode |
| 14 | `src/app/[locale]/dashboard/categories/page.tsx` | Read-only offline mode |
| 15 | `src/app/[locale]/dashboard/kds/page.tsx` | Integrate offline KDS with sync manager |
| 16 | `src/app/[locale]/dashboard/hr/employees/page.tsx` | Read-only offline mode |
| 17 | `src/app/[locale]/dashboard/hr/attendance/page.tsx` | Read-only offline mode |
| 18 | `src/app/[locale]/dashboard/hr/leaves/page.tsx` | Read-only offline mode |
| 19 | `src/app/[locale]/dashboard/hr/payroll/page.tsx` | Read-only offline mode |
| 20 | `src/app/[locale]/dashboard/drivers/page.tsx` | Read-only offline mode |
| 21 | `src/app/[locale]/dashboard/reports/overview/page.tsx` | Show "unavailable offline" |
| 22 | `src/app/[locale]/dashboard/reports/sales/page.tsx` | Show "unavailable offline" |
| 23 | `src/app/[locale]/dashboard/reports/profitability/page.tsx` | Show "unavailable offline" |
| 24 | `src/app/[locale]/dashboard/reports/inventory/page.tsx` | Show "unavailable offline" |
| 25 | `src/app/[locale]/dashboard/settings/page.tsx` | Add offline settings controls |
| 26 | `src/app/[locale]/dashboard/page.tsx` | Show cached dashboard stats |
| 27 | `src/i18n/messages/en.json` | Add `offline.*` translation keys |
| 28 | `src/i18n/messages/ar.json` | Add `offline.*` Arabic translations |
| 29 | `src/stores/kdsStore.ts` | (Minor) Integrate offline queue status |
| 30 | `frontend/apps/dashboard/next.config.js` | Add PWA headers, cache manifest route |
| 31 | `frontend/apps/dashboard/package.json` | Add `idb` dependency |
| 32 | `backend/cmd/api/main.go` | Register `POST /api/v1/sync/batch` route |

### 15.3 Test Files to Create

| # | File Path | What It Tests |
|---|---|---|
| 1 | `src/__tests__/lib/offline-db.test.ts` | IndexedDB CRUD, cleanup, quota |
| 2 | `src/__tests__/lib/sync-manager.test.ts` | Sync flow, conflict resolution, retry |
| 3 | `src/__tests__/lib/offline-api-layer.test.ts` | Online/offline switching, cache fallback |
| 4 | `src/__tests__/hooks/useOfflineOrders.test.ts` | Order cache reads |
| 5 | `backend/internal/handler/http/sync_handler_test.go` | Batch sync endpoint |

---

## 16. Multi-Tenant Isolation

All IndexedDB records include `tenantId` and `restaurantId` from the auth store:

```typescript
// From src/stores/authStore.ts, user object:
interface User {
  id: number
  tenant_id: number
  restaurant_id?: number
}

// Every IndexedDB write includes:
{ ...data, tenantId: user.tenant_id, restaurantId: user.restaurant_id }

// Every IndexedDB read filters by:
store.index('restaurantId').getAll(currentRestaurantId)
```

---

## 17. Implementation Phases

### Phase 1: Foundation (Week 1-2)
Create: `offline-db.ts`, `offlineStore.ts`, `OfflineProvider.tsx`, `sw.js`, `manifest.json`
Mount: ServiceWorkerProvider + OfflineProvider in layout
Modify: `layout.tsx`, `rbac.ts`, `useRbac.ts`
Install: `idb` package

### Phase 2: Sync Manager + Backend (Week 2-3)
Create: `sync-manager.ts`, `offline-api-layer.ts`, `sync_handler.go`
Modify: `main.go`, `useKDSOffline.ts`

### Phase 3: Orders + KDS Offline (Week 3-4)
Create: `useOfflineOrders.ts`, `useCreateOfflineOrder.ts`, `GlobalStatusBar.tsx`
Modify: Orders pages, KDS pages, DashboardLayout, Sidebar

### Phase 4: Remaining Modules (Week 4)
Create: `useOfflineProducts.ts`, `useOfflineCategories.ts`, `useOfflineTables.ts`, `useOfflineEmployees.ts`
Modify: All remaining dashboard pages

### Phase 5: Testing + Polish (Week 5)
Create: All test files
Modify: i18n files, settings page, feature flags
Add: PWA icons

---

## 18. Limitations

| Feature | Offline Status | Reason |
|---|---|---|
| Reports | NOT available | Requires server aggregation |
| Theme Builder | NOT available | Complex server operations |
| Billing/Payments | NOT available | Paymob requires live connection |
| Real-time notifications | NOT available | WebSocket only |
| Role management (RBAC admin) | NOT available | Security-critical, server-only |
| KDS WebSocket updates | Degraded | Falls back to polling on reconnect |

---

## 19. Summary

**21 new files** + **32 modified files** + **5 test files** = **58 total files**

The implementation leverages existing patterns:
- `useKDSOffline.ts` IndexedDB queue → generalized global offline queue
- `useServiceWorker.ts` / `ServiceWorkerProvider.tsx` → mounted and enhanced
- `PermissionCache` in `rbac.ts` → extended TTL for offline
- `sonner` toasts → used for sync notifications
- Zustand persist → new `offlineStore.ts` follows same pattern
- Backend `wrapAuthOnly` → used for new sync endpoint

