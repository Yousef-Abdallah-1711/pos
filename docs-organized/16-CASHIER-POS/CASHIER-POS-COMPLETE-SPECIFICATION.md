# Cashier POS — Complete Feature Specification

> **Version:** 1.0
> **Status:** Planning
> **Target:** Production-ready, high-performance cashier terminal
> **Metric Goals:** 3 clicks = full order in ~3 seconds, <200ms per action, 200+ orders/hour, fully offline-capable

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Architecture Overview](#2-architecture-overview)
3. [Backend API Contract](#3-backend-api-contract)
4. [Phase 1 — Core POS UI (MVP)](#4-phase-1--core-pos-ui-mvp)
5. [Phase 2 — Offline & Caching](#5-phase-2--offline--caching)
6. [Phase 3 — Hardware Integration](#6-phase-3--hardware-integration)
7. [Phase 4 — Advanced Performance & PWA](#7-phase-4--advanced-performance--pwa)
8. [File Map (All Phases)](#8-file-map-all-phases)
9. [Data Flow & State Architecture](#9-data-flow--state-architecture)
10. [Security Considerations](#10-security-considerations)
11. [Performance Budget](#11-performance-budget)
12. [Accessibility & i18n](#12-accessibility--i18n)
13. [Dependencies](#13-dependencies)

---

## 1. Executive Summary

Build a dedicated cashier POS terminal page at `/dashboard/pos` within the existing Next.js dashboard app (`frontend/apps/dashboard`). The POS is designed for restaurant cashiers taking orders at the counter or tableside on tablets/kiosk devices.

**Key decisions:**
- **Architecture:** React + Zustand + aggressive optimizations (memoization, IndexedDB, Service Worker, keyboard shortcuts) — within the existing dashboard app, not a separate app.
- **Backend:** Zero backend changes. Reuse existing `POST /api/v1/public/orders` with `order_source: 'pos'`. Reuse `PUT /api/v1/admin/orders/{id}/payment-status` for auto-marking cash payments as paid.
- **RBAC:** Uses existing `MODULE_IDS.ORDERS` (10) with `PERMISSION_LEVELS.WRITE` — no new migration needed.
- **Scope:** 4 phased implementations, each independently deployable.

---

## 2. Architecture Overview

```
┌──────────────────────────────────────────────────────────────────┐
│                        POS Terminal Page                         │
│  /[locale]/dashboard/pos                                         │
├──────────────┬───────────────────────────────────────────────────┤
│  Product     │  Cart Panel (35%)                                 │
│  Panel (65%) │                                                   │
│              │  ┌─ POSCartItem[] ──────────────────────────────┐ │
│  ┌─ Header ─┐│  │  name | qty +/- | price | remove             │ │
│  │ Search   ││  └────────────────────────────────────────────-─┘ │
│  │ Category ││  ┌─ Order Details ──────────────────────────────┐ │
│  │ Tabs     ││  │  Table selector | Customer | Payment method  │ │
│  └──────────┘│  └──────────────────────────────────────────────┘ │
│  ┌──────────┐│  ┌─ Summary ────────────────────────────────────┐ │
│  │ Product  ││  │  Subtotal | Tax (14%) | TOTAL                │ │
│  │ Grid     ││  └──────────────────────────────────────────────┘ │
│  │ (CSS     ││  ┌──────────────────────────────────────────────┐ │
│  │  Grid)   ││  │  [ PLACE ORDER ]  (Space)                    │ │
│  └──────────┘│  └──────────────────────────────────────────────┘ │
├──────────────┴───────────────────────────────────────────────────┤
│  State: Zustand posStore (persisted to localStorage)             │
│  API: posApi → existing api.ts (Bearer + tenant auto-injected)   │
│  Offline: IndexedDB queue (Phase 2) → replay on reconnect        │
│  Audio: Web Audio API (Phase 3) → add/remove/success/error beeps │
└──────────────────────────────────────────────────────────────────┘
```

**Rendering strategy:**
- Fullscreen mode (sidebar + header hidden, same as KDS pattern)
- Own header bar with clock, cashier name, online/offline indicator, fullscreen toggle
- RTL-aware layout (Arabic locale flips product panel and cart panel)
- CSS `contain: content` on product grid for layout isolation
- `React.memo` on all product cards and cart items
- Zustand selectors for minimal re-renders

---

## 3. Backend API Contract

**No backend changes required.** All endpoints already exist.

### 3.1 Endpoints Used

| Method | Endpoint | Purpose | Auth |
|--------|----------|---------|------|
| `GET` | `/api/v1/products?page_size=500` | Fetch all products | Bearer + Tenant |
| `GET` | `/api/v1/categories` | Fetch all categories | Bearer + Tenant |
| `GET` | `/api/v1/tables?status=available` | Fetch available tables | Bearer + Tenant |
| `POST` | `/api/v1/public/orders` | Create order | Bearer + Tenant |
| `PUT` | `/api/v1/admin/orders/{id}/payment-status` | Mark payment paid | Bearer + Tenant |
| `GET` | `/api/v1/admin/orders` | Orders list (existing page) | Bearer + Tenant |

### 3.2 Create Order Payload

```json
{
  "order_source": "pos",
  "customer_name": "Walk-in",
  "customer_phone": "0000000000",
  "payment_method": "cash",
  "table_id": 5,
  "notes": "Extra ketchup",
  "items": [
    {
      "product_id": 42,
      "quantity": 2,
      "special_instructions": "No onions"
    }
  ]
}
```

**Backend domain struct (`CreateOrderRequest`):**
- `customer_name` — required (default "Walk-in" for POS)
- `customer_phone` — required (default empty or "0000000000")
- `payment_method` — required: `'cash'`, `'card'`, `'online'`, `'wallet'`
- `order_source` — optional, defaults to `'website'`, POS sends `'pos'`
- `table_id` — optional, nullable (null = takeaway)
- `notes` — optional free-text
- `items[]` — required, min 1 item, each with `product_id`, `quantity`, optional `special_instructions`

### 3.3 Create Order Response

```json
{
  "success": true,
  "data": {
    "id": 123,
    "order_number": "ORD-20260225-001",
    "customer_name": "Walk-in",
    "table_id": 5,
    "total_amount": 45.50,
    "status": "pending",
    "created_at": "2026-02-25T14:30:00Z"
  }
}
```

### 3.4 Payment Status Update

After order creation, if `payment_method === 'cash'`, auto-call:
```
PUT /api/v1/admin/orders/{id}/payment-status
{ "payment_status": "paid" }
```

### 3.5 Order Source Badge in Orders Page

The existing orders page (`/dashboard/orders`) needs a new badge entry for `'pos'` in the `getOrderSourceBadge()` function (currently supports: `website`, `dine_in`, `mobile_app`, `phone`).

---

## 4. Phase 1 — Core POS UI (MVP)

### 4.1 New Files (14 files)

| # | File Path (relative to `frontend/apps/dashboard/`) | Purpose |
|---|------------------------------------------------------|---------|
| 1 | `src/app/[locale]/dashboard/pos/layout.tsx` | Viewport meta for kiosk/tablet (no zoom) |
| 2 | `src/app/[locale]/dashboard/pos/page.tsx` | Top-level page: RBAC check, data loading, keyboard shortcuts, render terminal |
| 3 | `src/stores/posStore.ts` | Zustand store with persist middleware for cart state |
| 4 | `src/lib/api/posApi.ts` | API functions: getProducts, getCategories, createOrder, updatePaymentStatus |
| 5 | `src/hooks/usePOSData.ts` | Data loading hook: products, categories, tables with localStorage cache |
| 6 | `src/hooks/usePOSKeyboard.ts` | Keyboard shortcut handler (F1-F12, Space, Escape, Ctrl+H, Ctrl+R) |
| 7 | `src/app/[locale]/dashboard/pos/components/POSTerminal.tsx` | Main flex layout (65% products + 35% cart), RTL-aware |
| 8 | `src/app/[locale]/dashboard/pos/components/POSHeader.tsx` | Clock, cashier name, online/offline, fullscreen, sound toggle, exit |
| 9 | `src/app/[locale]/dashboard/pos/components/POSCategoryTabs.tsx` | Horizontal scrollable category filter tabs |
| 10 | `src/app/[locale]/dashboard/pos/components/POSSearchBar.tsx` | Product search with 50ms debounce |
| 11 | `src/app/[locale]/dashboard/pos/components/POSProductGrid.tsx` | CSS grid of product cards, filtered by category + search |
| 12 | `src/app/[locale]/dashboard/pos/components/POSProductCard.tsx` | Product tile: name, price, image, availability badge, cart qty badge |
| 13 | `src/app/[locale]/dashboard/pos/components/POSCartPanel.tsx` | Right panel: cart items, table selector, customer, payment, summary, submit |
| 14 | `src/app/[locale]/dashboard/pos/components/POSCartItem.tsx` | Cart line: name, qty +/-, price, remove button |
| 15 | `src/app/[locale]/dashboard/pos/components/POSTableSelector.tsx` | Dropdown of available tables |
| 16 | `src/app/[locale]/dashboard/pos/components/POSCustomerInput.tsx` | Inline name + phone fields |
| 17 | `src/app/[locale]/dashboard/pos/components/POSPaymentSelector.tsx` | Cash/Card/Wallet toggle buttons |
| 18 | `src/app/[locale]/dashboard/pos/components/POSOrderSummary.tsx` | Subtotal, tax (14%), total |
| 19 | `src/app/[locale]/dashboard/pos/components/POSSubmitButton.tsx` | Large "Place Order" button with loading state |

### 4.2 Modified Files (4 files)

| # | File Path | Change |
|---|-----------|--------|
| 1 | `src/components/layout/Sidebar.tsx` | Add POS entry: `{ label: 'Cashier POS', href: '/dashboard/pos', icon: CreditCard, moduleId: 10, requiredPermission: WRITE }` |
| 2 | `src/components/layout/DashboardLayout.tsx` | Add `isPOSFullscreenRoute` detection, replace `isKDSFullscreenRoute` checks with `(isKDSFullscreenRoute \|\| isPOSFullscreenRoute)` |
| 3 | `src/i18n/messages/en.json` + `ar.json` | Add `pos.*` translation keys (~30 keys) and `navigation.pos` key |
| 4 | `src/app/[locale]/dashboard/orders/page.tsx` | Add `pos` entry to `getOrderSourceBadge()` config |

### 4.3 Zustand Store Shape (`posStore.ts`)

```typescript
interface POSState {
  // Cart
  items: POSCartItem[]               // { productId, name, price, quantity, notes, image, trackInventory, quantityInStock }
  tableId: number | null             // Selected table (null = takeaway)
  tableName: string | null
  customerName: string               // Default: '' (sent as 'Walk-in' on submit)
  customerPhone: string
  paymentMethod: 'cash' | 'card' | 'wallet'  // Default: 'cash'
  orderNotes: string

  // Held orders
  heldOrders: HeldOrder[]            // Parked orders for recall

  // UI state (not persisted)
  selectedCategoryId: number | null
  searchQuery: string

  // Settings
  soundEnabled: boolean

  // Actions
  addItem(product): void             // Inventory-capped, quantity auto-increment
  removeItem(productId): void
  updateQuantity(productId, qty): void
  updateItemNotes(productId, notes): void
  clearCart(): void
  setTableContext(id, name): void
  setCustomer(name, phone): void
  setPaymentMethod(method): void
  holdCurrentOrder(): string | null  // Returns held order ID
  recallOrder(id): void
  deleteHeldOrder(id): void

  // Computed
  getSubtotal(): number
  getTax(): number                   // 14% tax rate
  getTotal(): number
  getTotalItems(): number
}
```

**Persisted fields:** `items`, `tableId`, `tableName`, `customerName`, `customerPhone`, `paymentMethod`, `heldOrders`, `soundEnabled`

**Pattern references:**
- `src/stores/kdsStore.ts` — persist middleware with custom storage, validation, error recovery
- `restaurant-website/lib/store/cart-store.ts` — addItem with inventory capping

### 4.4 API Layer (`posApi.ts`)

```typescript
posApi.getProducts()                    // GET /products?page_size=500
posApi.getCategories()                  // GET /categories
posApi.createOrder(payload)             // POST /public/orders { order_source: 'pos', ... }
posApi.updateOrderPaymentStatus(id, s)  // PUT /admin/orders/{id}/payment-status
```

Uses existing `api` instance from `src/lib/api.ts` which auto-injects Bearer token + `X-Tenant-ID` + `X-Restaurant-ID` headers.

### 4.5 Data Loading Hook (`usePOSData.ts`)

```
Load priority: localStorage cache → API fetch (background)
1. Serve cached products/categories instantly from localStorage
2. Background-fetch fresh data from API (Promise.allSettled)
3. Update state + localStorage cache on success
4. Tables fetched from tablesApi.list('available')
Returns: { products, categories, tables, loading, error, refresh }
```

### 4.6 Keyboard Shortcuts (`usePOSKeyboard.ts`)

| Key | Action | Condition |
|-----|--------|-----------|
| F1–F12 | Quick-add product (index 0–11 of visible products) | Product exists + available |
| Space | Submit order | Cart has items + not submitting |
| Escape | Clear cart | Always |
| Ctrl+H | Hold current order | Cart has items |
| Ctrl+R | Recall last held order | Has held orders |

Skips when focus is on `INPUT`, `TEXTAREA`, or `SELECT` elements. Uses `useRef` for deps to avoid stale closures.

### 4.7 Order Submission Flow

```
1. Build payload from posStore state
   - customer_name: store.customerName || 'Walk-in'
   - customer_phone: store.customerPhone || undefined
   - payment_method: store.paymentMethod
   - table_id: store.tableId || undefined
   - order_source: 'pos'
   - items: store.items.map(i => ({ product_id, quantity, special_instructions }))

2. POST /api/v1/public/orders → { id, order_number, status, total_amount }

3. If payment_method === 'cash' && order.id exists:
   PUT /api/v1/admin/orders/{id}/payment-status { payment_status: 'paid' }
   (Non-critical — order still created if this fails)

4. Toast success: "Order placed! #ORD-20260225-001"

5. posStore.clearCart()
```

### 4.8 i18n Keys (~30 keys)

**English (`pos.*`):**
```
title, searchProducts, allCategories, cart, emptyCart, tapToAdd,
subtotal, tax, total, placeOrder, processing, customerName,
customerPhone, selectTable, noTable, seats, paymentMethod, cash,
card, wallet, orderSuccess, orderNumber, holdOrder, recallOrder,
orderHeld, orderRecalled, clearCart, quantity, notes, orderNotes,
outOfStock, online, offline, cashier, exitPOS, accessDenied,
accessDeniedText, noProducts, shortcutsHelp
```

**Arabic (`pos.*`):** Full Arabic translations for all keys.

**Navigation:** `navigation.pos: "Cashier POS"` / `"كاشير نقاط البيع"`

---

## 5. Phase 2 — Offline & Caching

### 5.1 New Files (3 files)

| # | File Path | Purpose |
|---|-----------|---------|
| 1 | `src/hooks/usePOSOffline.ts` | IndexedDB offline order queue with replay on reconnect |
| 2 | `src/hooks/usePOSProductCache.ts` | IndexedDB product catalog cache (5-min TTL, stale-while-revalidate) |
| 3 | `src/hooks/usePOSProductIndex.ts` | In-memory search index for <5ms product lookups |

### 5.2 Modified Files (2 files)

| # | File Path | Change |
|---|-----------|--------|
| 1 | `src/hooks/usePOSData.ts` | Enhanced priority: in-memory → IndexedDB → API |
| 2 | `src/app/[locale]/dashboard/pos/page.tsx` | Offline banner, queued order count, offline/online routing |

### 5.3 Offline Order Queue (`usePOSOffline.ts`)

**Pattern reference:** `src/hooks/useKDSOffline.ts`

```
IndexedDB database: 'pos-offline-db'
Object store: 'pending-orders'

Each queued order:
{
  id: string (uuid),
  payload: CreatePOSOrderPayload,
  queuedAt: ISO string,
  retryCount: number
}

Flow:
- If navigator.onLine === false → enqueueOrder(payload) to IndexedDB
- Toast: "Order queued — will sync when online"
- On 'online' event → replayQueuedOrders()
  - For each queued order: POST to API, if success delete from IndexedDB
  - If fail: increment retryCount, stop replay (FIFO order preservation)
- Badge shows queued count in POS header
```

### 5.4 Product Cache (`usePOSProductCache.ts`)

```
IndexedDB database: 'pos-product-db'
Object stores: 'products', 'categories'

Strategy: Stale-While-Revalidate
- On load: serve from IndexedDB (~1-5ms)
- Background-fetch from API
- If API succeeds: update IndexedDB + state
- TTL: 5 minutes (serve stale immediately, refresh in background)
- If API fails + IndexedDB has data: use stale data
- If API fails + no IndexedDB: show error
```

### 5.5 In-Memory Search Index (`usePOSProductIndex.ts`)

```
- On products load: build Map<barcode, product> + lowercase name array
- search(query): products.filter(p => p._lowerName.includes(query.toLowerCase()))
- getByBarcode(barcode): Map.get(barcode) — O(1) lookup for Phase 3
- Performance: <5ms for <500 products with includes() matching
```

---

## 6. Phase 3 — Hardware Integration

### 6.1 New Files (4 files + sound assets)

| # | File Path | Purpose |
|---|-----------|---------|
| 1 | `src/hooks/usePOSBarcodeScanner.ts` | Barcode scanner detection via rapid keystroke pattern |
| 2 | `src/hooks/usePOSAudio.ts` | Sound effects: add, remove, success, error (Web Audio API) |
| 3 | `src/app/[locale]/dashboard/pos/components/POSReceipt.tsx` | Hidden receipt layout for 80mm thermal printing |
| 4 | `src/hooks/usePOSPrint.ts` | Print handler via hidden iframe + window.print() |

### 6.2 Sound Assets

| File | Purpose |
|------|---------|
| `public/sounds/pos-add.mp3` | Item added to cart beep |
| `public/sounds/pos-remove.mp3` | Item removed beep |
| `public/sounds/pos-success.mp3` | Order submitted success chime |
| `public/sounds/pos-error.mp3` | Error/warning tone |

### 6.3 Modified Files (1 file)

| # | File Path | Change |
|---|-----------|--------|
| 1 | `src/app/[locale]/dashboard/pos/page.tsx` | Wire barcode scanner, audio hooks, auto-print on order success |

### 6.4 Barcode Scanner Detection (`usePOSBarcodeScanner.ts`)

```
Barcode scanners emulate rapid keyboard input followed by Enter.

Detection algorithm:
1. Track keystrokes arriving <50ms apart
2. Buffer characters until:
   a. Enter key received → process buffer as barcode
   b. >500ms gap → discard buffer (normal typing)
   c. Buffer < 6 chars on Enter → discard (too short for barcode)
3. Look up barcode in product index (Phase 2 usePOSProductIndex)
4. If found: posStore.addItem(product) + play add sound
5. If not found: toast "Product not found: {barcode}"

Edge cases:
- Ignore when focus is on INPUT/TEXTAREA/SELECT
- Clear buffer on Escape
- Max buffer length: 50 chars
```

### 6.5 Audio System (`usePOSAudio.ts`)

**Pattern reference:** `src/hooks/useKDSAudio.ts`

```
- Web Audio API primary, HTMLAudioElement fallback, oscillator tone last resort
- User gesture unlock pattern (pointerdown/keydown/touchstart)
- Sound mapping:
  - 'add': pos-add.mp3 (short beep)
  - 'remove': pos-remove.mp3 (lower beep)
  - 'success': pos-success.mp3 (chime)
  - 'error': pos-error.mp3 (warning tone)
- Respects posStore.soundEnabled toggle
- Preloads all audio files on mount
```

### 6.6 Receipt Printing

**POSReceipt.tsx — 80mm thermal paper layout:**
```
┌─────────────────────────────┐
│     RESTAURANT NAME          │
│     Date: 2026-02-25 14:30  │
│     Order: #ORD-001         │
│     Table: #5 - Patio       │
├─────────────────────────────┤
│  2x  Chicken Shawarma  $14 │
│  1x  Hummus            $5  │
│  3x  Fresh Juice       $12 │
├─────────────────────────────┤
│  Subtotal:         $31.00   │
│  Tax (14%):        $4.34    │
│  TOTAL:            $35.34   │
├─────────────────────────────┤
│  Payment: Cash - PAID       │
│  Cashier: Ahmed             │
├─────────────────────────────┤
│     Thank you!              │
│     شكراً لزيارتكم          │
└─────────────────────────────┘
```

**usePOSPrint.ts:**
- Creates hidden iframe with receipt HTML
- Applies `@media print` CSS
- Calls `iframe.contentWindow.print()`
- Auto-print option: triggers after successful order submission
- Removes iframe after print dialog closes

---

## 7. Phase 4 — Advanced Performance & PWA

### 7.1 New Files (4 files)

| # | File Path | Purpose |
|---|-----------|---------|
| 1 | `src/app/[locale]/dashboard/pos/components/POSVirtualGrid.tsx` | Virtual scrolling grid for 500+ products |
| 2 | `public/pos-sw.js` | Service Worker: cache-first static, network-first API |
| 3 | `src/hooks/usePOSBackgroundSync.ts` | Background Sync API for offline order replay |
| 4 | `src/hooks/usePOSPerformance.ts` | Performance monitoring using existing lib/performance.ts |

### 7.2 Modified Files (2 files)

| # | File Path | Change |
|---|-----------|--------|
| 1 | `src/app/[locale]/dashboard/pos/page.tsx` | Register Service Worker, conditional virtual grid, perf badge |
| 2 | `src/app/[locale]/dashboard/pos/components/POSProductGrid.tsx` | Swap to POSVirtualGrid when products > 50 |

### 7.3 Virtual Scrolling Grid (`POSVirtualGrid.tsx`)

```
- Uses @tanstack/react-virtual
- Fixed card size: 130x140px (matches CSS grid minmax)
- ResizeObserver for container width → compute columns
- Only renders visible rows ± 2 overscan
- Falls back to regular CSS grid if <50 products
- Smooth scrolling at 60fps even with 500+ products
```

### 7.4 Service Worker (`public/pos-sw.js`)

```
Strategy:
- Cache-first: JS, CSS, fonts, images, sound files
- Network-first: API calls (/api/v1/*)
- Precache: POS page bundle, product images
- Cache name: 'pos-cache-v1'

Background Sync:
- Register sync event 'pos-order-sync'
- On sync: read IndexedDB pending-orders, replay via fetch()
- On success: delete from IndexedDB
- On failure: re-register sync

Lifecycle:
- Install: precache critical assets
- Activate: clean old caches
- Fetch: route based on URL pattern
```

### 7.5 Performance Monitoring (`usePOSPerformance.ts`)

```
Uses existing src/lib/performance.ts recordPerformanceMetric()

Metrics tracked:
- pos.page-load: Time from navigation to interactive
- pos.add-to-cart: Time from click to store update
- pos.search: Time from keystroke to filtered results
- pos.order-submit: Time from click to API response
- pos.product-fetch: Time to load products from cache/API

Dev-mode: floating performance badge overlay showing last 5 metrics
Production: metrics silently recorded, available in console
```

---

## 8. File Map (All Phases)

### New Files Total: ~25 files + 4 sound assets

```
frontend/apps/dashboard/
├── src/
│   ├── app/[locale]/dashboard/pos/
│   │   ├── layout.tsx                          # Phase 1
│   │   ├── page.tsx                            # Phase 1 (modified in 2, 3, 4)
│   │   └── components/
│   │       ├── POSTerminal.tsx                 # Phase 1
│   │       ├── POSHeader.tsx                   # Phase 1
│   │       ├── POSCategoryTabs.tsx             # Phase 1
│   │       ├── POSSearchBar.tsx                # Phase 1
│   │       ├── POSProductGrid.tsx              # Phase 1 (modified in 4)
│   │       ├── POSProductCard.tsx              # Phase 1
│   │       ├── POSCartPanel.tsx                # Phase 1
│   │       ├── POSCartItem.tsx                 # Phase 1
│   │       ├── POSTableSelector.tsx            # Phase 1
│   │       ├── POSCustomerInput.tsx            # Phase 1
│   │       ├── POSPaymentSelector.tsx          # Phase 1
│   │       ├── POSOrderSummary.tsx             # Phase 1
│   │       ├── POSSubmitButton.tsx             # Phase 1
│   │       ├── POSReceipt.tsx                  # Phase 3
│   │       └── POSVirtualGrid.tsx              # Phase 4
│   ├── stores/
│   │   └── posStore.ts                         # Phase 1
│   ├── lib/api/
│   │   └── posApi.ts                           # Phase 1
│   └── hooks/
│       ├── usePOSData.ts                       # Phase 1 (modified in 2)
│       ├── usePOSKeyboard.ts                   # Phase 1
│       ├── usePOSOffline.ts                    # Phase 2
│       ├── usePOSProductCache.ts               # Phase 2
│       ├── usePOSProductIndex.ts               # Phase 2
│       ├── usePOSBarcodeScanner.ts             # Phase 3
│       ├── usePOSAudio.ts                      # Phase 3
│       ├── usePOSPrint.ts                      # Phase 3
│       ├── usePOSBackgroundSync.ts             # Phase 4
│       └── usePOSPerformance.ts                # Phase 4
├── public/
│   ├── sounds/
│   │   ├── pos-add.mp3                         # Phase 3
│   │   ├── pos-remove.mp3                      # Phase 3
│   │   ├── pos-success.mp3                     # Phase 3
│   │   └── pos-error.mp3                       # Phase 3
│   └── pos-sw.js                               # Phase 4
```

### Modified Files Total: ~6 files

```
frontend/apps/dashboard/
├── src/
│   ├── components/layout/
│   │   ├── Sidebar.tsx                         # Phase 1: Add POS nav entry
│   │   └── DashboardLayout.tsx                 # Phase 1: Add POS fullscreen detection
│   ├── i18n/messages/
│   │   ├── en.json                             # Phase 1: Add pos.* keys
│   │   └── ar.json                             # Phase 1: Add pos.* keys (Arabic)
│   └── app/[locale]/dashboard/orders/
│       └── page.tsx                            # Phase 1: Add 'pos' to order source badge
```

---

## 9. Data Flow & State Architecture

### 9.1 State Layers

```
Layer 1: Zustand posStore (in-memory, persisted to localStorage)
  → Cart items, table, customer, payment method, held orders, UI prefs
  → Instant read/write, survives page refresh

Layer 2: localStorage cache (pos-products-cache, pos-categories-cache)
  → Product catalog + categories
  → Served on page load before API fetch
  → Phase 1 only

Layer 3: IndexedDB (Phase 2)
  → pos-product-db: Full product catalog with 5-min TTL
  → pos-offline-db: Pending orders queue
  → Survives storage pressure, larger capacity than localStorage

Layer 4: Service Worker cache (Phase 4)
  → Static assets (JS, CSS, images, sounds)
  → Enables true offline PWA experience
```

### 9.2 Data Loading Priority (Phase 2+)

```
Products/Categories:
  1. In-memory (Zustand/React state) — 0ms
  2. IndexedDB (pos-product-db) — 1-5ms
  3. localStorage (pos-products-cache) — 1-3ms
  4. API fetch (GET /products) — 200-500ms
  → Always serve fastest available, background-refresh from API

Orders:
  Online: POST /api/v1/public/orders
  Offline: IndexedDB pos-offline-db → replay on reconnect
```

---

## 10. Security Considerations

### 10.1 Authentication & Authorization
- RBAC check on page load: `MODULE_IDS.ORDERS` (10) with `PERMISSION_LEVELS.WRITE`
- All API calls go through `api.ts` interceptor which auto-injects Bearer token + tenant headers
- If token expires (401), user is redirected to login (existing behavior)
- Access denied screen shown if user lacks ORDERS WRITE permission

### 10.2 Input Validation
- Product prices come from backend — never trust client-side prices for order totals
- Backend recalculates totals from product_id + quantity (client totals are display-only)
- Customer name/phone sanitized (trim whitespace)
- Order notes: free-text, max length enforced by backend
- Table ID validated against available tables from backend

### 10.3 XSS Prevention
- All user inputs rendered via React JSX (auto-escaped)
- No `dangerouslySetInnerHTML` used
- Product names, customer names displayed via text content only
- Receipt HTML generated via React components (not string concatenation)

### 10.4 Offline Queue Security
- Queued orders stored in IndexedDB (browser-sandboxed, same-origin)
- No sensitive data in IndexedDB beyond order payloads
- Queue replayed with same Bearer token (may fail if token expired — user re-authenticates)

### 10.5 Service Worker Security
- Service Worker scoped to `/dashboard/pos` only
- Cache-first only for static assets (immutable)
- API calls always network-first (never serve stale API responses)
- SW registered only in production builds

### 10.6 CSRF Protection
- API uses Bearer token auth (not cookies) — inherently CSRF-resistant
- No state-changing requests use cookies

---

## 11. Performance Budget

| Metric | Target | Technique |
|--------|--------|-----------|
| Page load (cached) | <200ms | localStorage/IndexedDB product cache |
| Page load (cold) | <1s | Parallel API fetches, code splitting |
| Add to cart | <20ms | React.memo on cards, Zustand selectors |
| Search products | <50ms | useMemo filter, 50ms debounce |
| Order submission | <500ms | Direct API call, optimistic UI |
| Product grid scroll | 60fps | CSS `contain: content`, lazy images |
| Product grid (500+ items) | 60fps | Virtual scrolling (Phase 4) |
| Memory usage | <50MB | Only visible DOM nodes rendered |

### Optimization Techniques

| Technique | Target | Implementation |
|-----------|--------|---------------|
| `React.memo` | Product cards, cart items | Only re-render when own data changes |
| Zustand selectors | All components | Subscribe to specific state slices, not entire store |
| `useMemo` | Filtered product list | Recompute only when category/search changes |
| `useCallback` | All event handlers | Stable references prevent child re-renders |
| CSS `contain: content` | Product grid container | Layout isolation — grid doesn't trigger full-page reflow |
| `loading="lazy"` | Product images | Load images on scroll, not on mount |
| Debounced search (50ms) | Search input | Avoid re-filtering every keystroke |
| localStorage snapshots | Products, categories | Serve cached data in <3ms on page load |
| IndexedDB catalog | Products (Phase 2) | 1-5ms lookup, larger capacity |
| `@tanstack/react-virtual` | Product grid (Phase 4) | Render only visible rows |
| Service Worker | Static assets (Phase 4) | Cache-first for JS/CSS/images |
| `useRef` for callbacks | Keyboard shortcuts | Avoid stale closures without re-registering listeners |

---

## 12. Accessibility & i18n

### 12.1 RTL Support
- `flex-row-reverse` on all horizontal layouts when `locale === 'ar'`
- Product panel and cart panel swap sides in RTL
- Text alignment: `text-right` for Arabic
- Chevron icons swap direction
- All hardcoded margins/paddings use RTL-aware classes

### 12.2 Keyboard Navigation
- F1–F12 for quick product add (visible shortcut badges on cards)
- Space for order submission
- Escape for cart clear
- Ctrl+H / Ctrl+R for hold/recall
- Tab navigation through form fields

### 12.3 Screen Reader
- Semantic HTML: `<header>`, `<main>`, `<nav>`, `<button>`
- ARIA labels on icon-only buttons
- Cart item count announced on change
- Toast notifications use `role="status"`

### 12.4 Touch/Tablet
- Minimum touch target: 44x44px on all buttons
- `active:scale-95` feedback on product cards
- No hover-only interactions
- Viewport meta prevents zoom (kiosk mode)

---

## 13. Dependencies

### Existing (no new installs for Phase 1-2)
- `zustand` — state management
- `axios` — HTTP client (via existing `api.ts`)
- `sonner` — toast notifications
- `lucide-react` — icons
- `next` — framework
- `react` — UI library

### New Dependencies

| Package | Phase | Purpose |
|---------|-------|---------|
| `@tanstack/react-virtual` | Phase 4 | Virtual scrolling for 500+ product grids |

No other new dependencies needed. All offline/caching uses native browser APIs (IndexedDB, Service Worker, Web Audio API).

---

## Key Existing Files to Reuse

| File | Reuse For |
|------|-----------|
| `src/lib/api.ts` | Axios instance with auth/tenant interceptors |
| `src/lib/api/tablesApi.ts` | Table listing API |
| `src/hooks/useTables.ts` | Table data patterns |
| `src/stores/kdsStore.ts` | Zustand persist + custom storage + optimistic update pattern |
| `src/hooks/useKDSOffline.ts` | IndexedDB offline queue pattern |
| `src/hooks/useKDSAudio.ts` | Web Audio API + user gesture unlock pattern |
| `src/hooks/useRbac.ts` | `useHasPermission()` hook |
| `src/lib/rbac.ts` | `MODULE_IDS`, `PERMISSION_LEVELS` constants |
| `src/components/layout/DashboardLayout.tsx` | Fullscreen route detection pattern |
| `src/components/layout/Sidebar.tsx` | Navigation entry pattern |
| `restaurant-website/lib/store/cart-store.ts` | Cart addItem with inventory capping pattern |
| `src/lib/performance.ts` | `recordPerformanceMetric()` utility |
| `src/lib/debug.ts` | `debugLog()`, `debugError()` utilities |
| `src/lib/translations.ts` | `createTranslator()`, `getLocaleFromPath()` |
