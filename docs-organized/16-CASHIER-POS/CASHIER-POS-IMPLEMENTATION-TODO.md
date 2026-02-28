# Cashier POS — Implementation TODO (Production-Ready)

> **Companion doc:** [CASHIER-POS-COMPLETE-SPECIFICATION.md](./CASHIER-POS-COMPLETE-SPECIFICATION.md)
> **Each phase is independently deployable.** Complete all items in a phase before moving to the next.
> **Legend:** ⬜ = TODO, 🔲 = In Progress, ✅ = Done

---

## Phase 1: Core POS UI (MVP)

### 1.1 Store & State Management

- ⬜ Create `src/stores/posStore.ts`
  - ⬜ Define `POSCartItem` interface (productId, name, price, quantity, notes, image, trackInventory, quantityInStock)
  - ⬜ Define `HeldOrder` interface (id, items, tableId, tableName, customerName, customerPhone, paymentMethod, orderNotes, heldAt)
  - ⬜ Define `PaymentMethod` type: `'cash' | 'card' | 'wallet'`
  - ⬜ Implement Zustand store with `persist` middleware (storage key: `'pos-storage'`)
  - ⬜ Implement custom `PersistStorage` with validation + error recovery (pattern: kdsStore.ts)
  - ⬜ Implement `normalizePersistedState()` for corrupted storage recovery
  - ⬜ Implement `addItem()` with inventory capping (pattern: cart-store.ts)
  - ⬜ Implement `removeItem()`, `updateQuantity()`, `updateItemNotes()`
  - ⬜ Implement `clearCart()` — resets items + tableId + customer + payment + notes
  - ⬜ Implement `setTableContext()`, `setCustomer()`, `setPaymentMethod()`, `setOrderNotes()`
  - ⬜ Implement `holdCurrentOrder()` — moves current cart to heldOrders array, returns ID
  - ⬜ Implement `recallOrder(id)` — moves held order back to active cart
  - ⬜ Implement `deleteHeldOrder(id)`
  - ⬜ Implement computed: `getSubtotal()`, `getTax()` (14%), `getTotal()`, `getTotalItems()`
  - ⬜ Implement partialize: persist only items, tableId, tableName, customerName, customerPhone, paymentMethod, heldOrders, soundEnabled
  - ⬜ Export optimized selectors for each state slice

### 1.2 API Layer

- ⬜ Create `src/lib/api/posApi.ts`
  - ⬜ Define `POSProduct` interface (id, name, name_ar, description, price, image, category_id, category_name, is_available, track_inventory, quantity_in_stock, barcode)
  - ⬜ Define `POSCategory` interface (id, name, name_ar, description, image, product_count)
  - ⬜ Define `CreatePOSOrderPayload` interface (order_source, customer_name, customer_phone, payment_method, table_id, notes, items[])
  - ⬜ Define `POSOrderResponse` interface (id, order_number, status, total)
  - ⬜ Implement `getProducts()` — `GET /products?page_size=500`, handle response.data.products || response.data.data || response.data
  - ⬜ Implement `getCategories()` — `GET /categories`, same response handling
  - ⬜ Implement `createOrder(payload)` — `POST /public/orders`
  - ⬜ Implement `updateOrderPaymentStatus(orderId, status)` — `PUT /admin/orders/{id}/payment-status`
  - ⬜ Use existing `api` instance from `src/lib/api.ts` (auto Bearer + tenant headers)
  - ⬜ Use `debugError()` for all catch blocks

### 1.3 Data Loading Hook

- ⬜ Create `src/hooks/usePOSData.ts`
  - ⬜ Implement localStorage cache read/write helpers (`pos-products-cache`, `pos-categories-cache`)
  - ⬜ Initialize state from localStorage cache (instant data on page load)
  - ⬜ Fetch products, categories, tables in parallel via `Promise.allSettled`
  - ⬜ Use `tablesApi.list('available')` from existing `src/lib/api/tablesApi.ts`
  - ⬜ Update state + localStorage cache on successful fetch
  - ⬜ Handle partial failures (show data from successful fetches)
  - ⬜ Return `{ products, categories, tables, loading, error, refresh }`
  - ⬜ Guard against double-fetching with `useRef` flag

### 1.4 Keyboard Shortcuts Hook

- ⬜ Create `src/hooks/usePOSKeyboard.ts`
  - ⬜ Single `useEffect` with `keydown` listener
  - ⬜ Skip when target is INPUT/TEXTAREA/SELECT
  - ⬜ F1–F12: Quick-add visible product at index 0–11
  - ⬜ Space: Submit order (when cart has items + not submitting)
  - ⬜ Escape: Clear cart
  - ⬜ Ctrl+H: Hold current order
  - ⬜ Ctrl+R: Recall last held order
  - ⬜ Use `useRef` for options to avoid stale closures
  - ⬜ Clean up listener on unmount

### 1.5 Page & Layout

- ⬜ Create `src/app/[locale]/dashboard/pos/layout.tsx`
  - ⬜ Export `viewport` meta: `userScalable: false`, `maximumScale: 1`, `viewportFit: 'cover'`
  - ⬜ Pass-through children (same as KDS layout)

- ⬜ Create `src/app/[locale]/dashboard/pos/page.tsx`
  - ⬜ `'use client'` directive
  - ⬜ RBAC check: `useHasPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.WRITE)`
  - ⬜ Show loading spinner during permission check
  - ⬜ Show access denied screen if no permission (Lock icon + message)
  - ⬜ Load data via `usePOSData()`
  - ⬜ Register keyboard shortcuts via `usePOSKeyboard()`
  - ⬜ Compute `visibleProducts` via `useMemo` (filter by category + search + availability)
  - ⬜ Implement `handleAddProduct()` with locale-aware name
  - ⬜ Implement `handleSubmitOrder()` — full order submission flow with payment status update
  - ⬜ Implement `handleClearCart()`, `handleHoldOrder()`, `handleRecallOrder()`
  - ⬜ Render `<POSTerminal />` with all props

### 1.6 Components

- ⬜ Create `POSTerminal.tsx`
  - ⬜ Full-height flex layout: product panel (65%) + cart panel (35%)
  - ⬜ RTL-aware: `flex-row-reverse` when `locale === 'ar'`
  - ⬜ POSHeader at top, panels below
  - ⬜ Loading state: spinner when no cached data
  - ⬜ Error state: error message + retry button

- ⬜ Create `POSHeader.tsx`
  - ⬜ Real-time clock (updates every second)
  - ⬜ Cashier name from authStore
  - ⬜ Online/offline indicator (navigator.onLine + event listeners)
  - ⬜ Held orders count badge
  - ⬜ Sound toggle button
  - ⬜ Fullscreen toggle (document.requestFullscreen/exitFullscreen)
  - ⬜ Exit POS button (navigates to dashboard)
  - ⬜ Dark header bar (bg-gray-900)
  - ⬜ RTL-aware layout

- ⬜ Create `POSCategoryTabs.tsx`
  - ⬜ `React.memo` wrapper
  - ⬜ Horizontal scrollable container with `overflow-x-auto`
  - ⬜ "All" button (selectedCategoryId === null)
  - ⬜ Category buttons from props
  - ⬜ Active state: primary bg + white text
  - ⬜ Click toggles category (click same = deselect)
  - ⬜ Locale-aware: show name_ar when Arabic

- ⬜ Create `POSSearchBar.tsx`
  - ⬜ 50ms debounced input updating posStore.searchQuery
  - ⬜ Search icon (left for LTR, right for RTL)
  - ⬜ Clear button when query is non-empty
  - ⬜ `defaultValue` from store for controlled-uncontrolled hybrid (performance)

- ⬜ Create `POSProductGrid.tsx`
  - ⬜ CSS grid: `repeat(auto-fill, minmax(130px, 1fr))`
  - ⬜ `contain: content` for layout isolation
  - ⬜ Empty state: icon + "No products found"
  - ⬜ Renders POSProductCard for each product

- ⬜ Create `POSProductCard.tsx`
  - ⬜ `React.memo` wrapper
  - ⬜ Product image or initial-letter placeholder
  - ⬜ Product name (line-clamp-2)
  - ⬜ Price display
  - ⬜ Cart quantity badge (top-right corner)
  - ⬜ Keyboard shortcut badge (F1-F12, top-left corner)
  - ⬜ Out-of-stock disabled state with label
  - ⬜ `active:scale-95` touch feedback
  - ⬜ `loading="lazy"` on images
  - ⬜ Zustand selector for own cart item only (minimal re-renders)

- ⬜ Create `POSCartPanel.tsx`
  - ⬜ Cart header with item count badge
  - ⬜ Hold button (Ctrl+H hint)
  - ⬜ Recall button (Ctrl+R hint) — shown only when held orders exist
  - ⬜ Clear cart button (Esc hint)
  - ⬜ Scrollable cart items area
  - ⬜ Empty cart state (icon + message + hint)
  - ⬜ Order details section: table + customer + payment + notes
  - ⬜ Order summary section
  - ⬜ Submit button section
  - ⬜ Only shows bottom section when cart has items

- ⬜ Create `POSCartItem.tsx`
  - ⬜ `React.memo` wrapper
  - ⬜ Item name + unit price
  - ⬜ Quantity controls: minus (-) / display / plus (+)
  - ⬜ Line total (price × quantity)
  - ⬜ Remove button (X)
  - ⬜ RTL-aware layout

- ⬜ Create `POSTableSelector.tsx`
  - ⬜ Dropdown select element
  - ⬜ Default option: "No table (Takeaway)"
  - ⬜ Options: `#{table_number} - {name} ({capacity} seats)` for each available table
  - ⬜ Updates posStore.setTableContext on change

- ⬜ Create `POSCustomerInput.tsx`
  - ⬜ Two inline inputs: name + phone
  - ⬜ Updates posStore.setCustomer on change
  - ⬜ Phone input type="tel"

- ⬜ Create `POSPaymentSelector.tsx`
  - ⬜ Three toggle buttons: Cash (Banknote icon), Card (CreditCard icon), Wallet (Wallet icon)
  - ⬜ Active state: primary bg
  - ⬜ Updates posStore.setPaymentMethod on click

- ⬜ Create `POSOrderSummary.tsx`
  - ⬜ Subtotal line
  - ⬜ Tax (14%) line
  - ⬜ Total line (large font, primary color)
  - ⬜ Zustand selectors for each value
  - ⬜ `tabular-nums` for number alignment

- ⬜ Create `POSSubmitButton.tsx`
  - ⬜ Large button: "Place Order (N)" where N = total items
  - ⬜ "Space" keyboard hint
  - ⬜ Loading state: spinner + "Processing..."
  - ⬜ Disabled when cart empty or submitting
  - ⬜ Green bg (emerald-500) with shadow

### 1.7 Modify Existing Files

- ⬜ Modify `src/components/layout/Sidebar.tsx`
  - ⬜ Add `CreditCard` to lucide-react import
  - ⬜ Add POS entry after KDS: `{ label: t('navigation.pos', 'Cashier POS'), href: '/dashboard/pos', icon: <CreditCard />, moduleId: MODULE_IDS.ORDERS, requiredPermission: PERMISSION_LEVELS.WRITE }`

- ⬜ Modify `src/components/layout/DashboardLayout.tsx`
  - ⬜ Add `const isPOSFullscreenRoute = pathname.includes('/dashboard/pos')`
  - ⬜ Add `const isFullscreenRoute = isKDSFullscreenRoute || isPOSFullscreenRoute`
  - ⬜ Replace all `isKDSFullscreenRoute` references (sidebar hide, header hide, trial banner hide, content class, main class, bg class) with `isFullscreenRoute`
  - ⬜ Keep `isKDSFullscreenRoute` for the trial banner check: `!isBlocked && !isFullscreenRoute`
  - ⬜ Keep `isKDSFullscreenRoute` for main class: use `isFullscreenRoute` instead

- ⬜ Modify `src/i18n/messages/en.json`
  - ⬜ Add `"navigation.pos": "Cashier POS"`
  - ⬜ Add `"pos"` object with all ~30 translation keys

- ⬜ Modify `src/i18n/messages/ar.json`
  - ⬜ Add `"navigation.pos": "كاشير نقاط البيع"`
  - ⬜ Add `"pos"` object with all Arabic translations

- ⬜ Modify `src/app/[locale]/dashboard/orders/page.tsx`
  - ⬜ Add `pos` entry to `getOrderSourceBadge()` config:
    ```
    pos: {
      label: t('orders.sources.pos', 'POS'),
      className: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300',
    }
    ```

### 1.8 Phase 1 Testing

- ⬜ **Manual Testing Checklist:**
  - ⬜ Navigate to `/en/dashboard/pos` — page loads with product grid and empty cart
  - ⬜ Click a product — appears in cart with correct price
  - ⬜ Click same product again — quantity increments (not duplicate line)
  - ⬜ Click + / - on cart item — quantity changes correctly
  - ⬜ Click X on cart item — removed from cart
  - ⬜ Select a table from dropdown — reflected in store
  - ⬜ Enter customer name + phone — reflected in store
  - ⬜ Toggle Cash/Card/Wallet — highlighted correctly
  - ⬜ Verify subtotal, tax (14%), total calculations are correct
  - ⬜ Press Space — order submitted, toast shows order number, cart clears
  - ⬜ Check `/dashboard/orders` — new order appears with "POS" source badge
  - ⬜ Verify cash orders are auto-marked as "paid"
  - ⬜ Press F1–F12 — products added via keyboard (visible shortcut badges)
  - ⬜ Press Escape — cart clears
  - ⬜ Press Ctrl+H — order held, cart clears, held count badge shows in header
  - ⬜ Press Ctrl+R — held order recalled back to cart
  - ⬜ Verify keyboard shortcuts DON'T fire when focused on search input
  - ⬜ Verify sidebar shows "Cashier POS" entry with CreditCard icon
  - ⬜ Click POS sidebar entry — navigates to POS page
  - ⬜ Verify sidebar + header are hidden on POS page (fullscreen mode)
  - ⬜ Test fullscreen toggle button in POS header
  - ⬜ Test exit button — navigates back to dashboard
  - ⬜ Test AR locale — layout flips RTL, all labels translated
  - ⬜ Test with user lacking ORDERS WRITE permission — access denied screen shown
  - ⬜ Test with out-of-stock product — card disabled, cannot add
  - ⬜ Test with inventory-tracked product — quantity capped at stock level
  - ⬜ Page refresh — cart items persist (localStorage), products load from cache then refresh
  - ⬜ Test category filter — clicking a category shows only its products
  - ⬜ Test search — typing filters products in <50ms
  - ⬜ Test held orders persist across page refresh

- ⬜ **Error Handling Tests:**
  - ⬜ API error on products fetch — error state shown, retry button works
  - ⬜ API error on order submission — toast shows error message, cart NOT cleared
  - ⬜ Network timeout on order submission — toast shows timeout error
  - ⬜ Payment status update fails — order still created, no crash
  - ⬜ Corrupted localStorage — store recovers to defaults (normalizePersistedState)

- ⬜ **Performance Checks (Phase 1):**
  - ⬜ Product grid renders 100+ products without visible lag
  - ⬜ Adding item to cart: no visible delay
  - ⬜ Search input: results appear immediately (50ms debounce)
  - ⬜ Page load with cached products: no loading spinner visible

---

## Phase 2: Offline & Caching

### 2.1 Offline Order Queue

- ⬜ Create `src/hooks/usePOSOffline.ts`
  - ⬜ IndexedDB database: `'pos-offline-db'`, store: `'pending-orders'`
  - ⬜ Implement `openQueueDB()` with upgrade handler (pattern: useKDSOffline.ts)
  - ⬜ Implement `enqueueOrder(payload)` — saves order to IndexedDB with uuid, timestamp, retryCount
  - ⬜ Implement `readAllPendingOrders()` — returns sorted by queuedAt
  - ⬜ Implement `removePendingOrder(id)` — deletes single order after successful replay
  - ⬜ Implement `replayQueuedOrders()` — replays queue in FIFO order via posApi.createOrder
    - ⬜ For each order: create → if cash: update payment status → remove from queue
    - ⬜ On failure: increment retryCount, stop replay (preserve order)
    - ⬜ Guard against concurrent replay with ref flag
  - ⬜ Online/offline detection via `navigator.onLine` + event listeners
  - ⬜ Auto-trigger replay on 'online' event
  - ⬜ Return `{ isOffline, queuedOrderCount, isReplaying, enqueueOrder, replayQueuedOrders, clearQueue }`

### 2.2 Product Catalog Cache

- ⬜ Create `src/hooks/usePOSProductCache.ts`
  - ⬜ IndexedDB database: `'pos-product-db'`, stores: `'products'`, `'categories'`
  - ⬜ Implement `readProducts()` — reads all products from IndexedDB
  - ⬜ Implement `writeProducts(products)` — bulk writes to IndexedDB with timestamp
  - ⬜ Implement `readCategories()` — reads all categories from IndexedDB
  - ⬜ Implement `writeCategories(categories)` — bulk writes to IndexedDB with timestamp
  - ⬜ Implement TTL check: 5 minutes (serve stale immediately, refresh in background)
  - ⬜ Return `{ products, categories, isStale, refresh }`

### 2.3 In-Memory Search Index

- ⬜ Create `src/hooks/usePOSProductIndex.ts`
  - ⬜ On products change: build lowercase name index (`_lowerName` field)
  - ⬜ On products change: build barcode Map (`Map<string, POSProduct>`)
  - ⬜ Implement `search(query)` — `includes()` matching on _lowerName, <5ms for <500 products
  - ⬜ Implement `getByBarcode(barcode)` — O(1) Map lookup
  - ⬜ Memoize index rebuild (only when products array reference changes)

### 2.4 Modify Existing Files

- ⬜ Modify `src/hooks/usePOSData.ts`
  - ⬜ Enhanced loading priority: in-memory → IndexedDB → localStorage → API
  - ⬜ On load: serve from IndexedDB first (1-5ms), then background-fetch
  - ⬜ On API success: update IndexedDB + localStorage + state
  - ⬜ On API failure with IndexedDB data: use stale data, show warning

- ⬜ Modify `src/app/[locale]/dashboard/pos/page.tsx`
  - ⬜ Wire `usePOSOffline()` hook
  - ⬜ Show offline banner when `isOffline === true`
  - ⬜ Show queued order count badge in header
  - ⬜ Route order submission: online → API, offline → IndexedDB queue
  - ⬜ Toast "Order queued — will sync when online" for offline orders
  - ⬜ Auto-replay on reconnect with success/failure toasts

### 2.5 Phase 2 Testing

- ⬜ **Offline Testing Checklist:**
  - ⬜ Load POS page — products appear instantly from IndexedDB (no loading spinner)
  - ⬜ Go offline (DevTools > Network > Offline) — "Offline" indicator shows in header
  - ⬜ Submit order while offline — toast says "Order queued"
  - ⬜ Verify order is in IndexedDB (DevTools > Application > IndexedDB > pos-offline-db)
  - ⬜ Submit another order while offline — both queued
  - ⬜ Go back online — queued orders auto-sync in FIFO order
  - ⬜ Toast confirms each synced order with order number
  - ⬜ Verify orders appear in `/dashboard/orders` list
  - ⬜ Verify IndexedDB pending-orders store is empty after sync
  - ⬜ Search products while offline — results appear in <50ms

- ⬜ **Cache Testing:**
  - ⬜ First load (cold): products fetched from API, stored in IndexedDB
  - ⬜ Second load: products served from IndexedDB (~1-5ms), then background refresh
  - ⬜ TTL expired (>5min): stale data shown immediately, fresh data fetched
  - ⬜ IndexedDB cleared manually: falls back to localStorage, then API
  - ⬜ Both IndexedDB + localStorage cleared: full API fetch with loading spinner

- ⬜ **Edge Cases:**
  - ⬜ Go offline → queue 3 orders → first replay fails → stops, retains remaining 2
  - ⬜ Go offline → go online → go offline immediately → queue preserved
  - ⬜ Token expires while offline → replay fails with 401 → user redirected to login
  - ⬜ IndexedDB quota exceeded — graceful fallback to localStorage cache

- ⬜ **Performance Checks:**
  - ⬜ Products load from IndexedDB in <5ms (measure with DevTools Performance)
  - ⬜ Product search: <5ms for 500 products (measure with console.time)
  - ⬜ Offline order enqueue: <10ms to IndexedDB write
  - ⬜ Online replay: orders sent sequentially, no duplicate submissions

---

## Phase 3: Hardware Integration

### 3.1 Barcode Scanner

- ⬜ Create `src/hooks/usePOSBarcodeScanner.ts`
  - ⬜ Track keystrokes arriving <50ms apart in buffer
  - ⬜ On Enter + buffer >= 6 chars + elapsed < 500ms → treat as barcode
  - ⬜ Look up barcode via `usePOSProductIndex.getByBarcode()`
  - ⬜ If found: `posStore.addItem(product)` + play 'add' sound
  - ⬜ If not found: toast "Product not found: {barcode}"
  - ⬜ Clear buffer on Escape or >500ms gap
  - ⬜ Max buffer length: 50 chars
  - ⬜ Ignore when focus is on INPUT/TEXTAREA/SELECT
  - ⬜ Clean up listener on unmount

### 3.2 Audio System

- ⬜ Create `src/hooks/usePOSAudio.ts`
  - ⬜ Pattern: follow `src/hooks/useKDSAudio.ts` exactly
  - ⬜ Sound file paths: `/sounds/pos-add.mp3`, `/sounds/pos-remove.mp3`, `/sounds/pos-success.mp3`, `/sounds/pos-error.mp3`
  - ⬜ Implement `ensureAudioContext()` with WebKit fallback
  - ⬜ Implement `unlockAudio()` on first user gesture (pointerdown/keydown/touchstart)
  - ⬜ Implement `playSound(type)` where type = 'add' | 'remove' | 'success' | 'error'
  - ⬜ HTMLAudioElement primary, oscillator tone fallback
  - ⬜ Respect `posStore.soundEnabled` toggle
  - ⬜ Preload all audio files on mount
  - ⬜ Clean up audio elements + AudioContext on unmount

- ⬜ Create/source sound assets:
  - ⬜ `public/sounds/pos-add.mp3` — short beep (~100ms)
  - ⬜ `public/sounds/pos-remove.mp3` — lower beep (~100ms)
  - ⬜ `public/sounds/pos-success.mp3` — success chime (~300ms)
  - ⬜ `public/sounds/pos-error.mp3` — warning tone (~200ms)

### 3.3 Receipt Printing

- ⬜ Create `src/app/[locale]/dashboard/pos/components/POSReceipt.tsx`
  - ⬜ Hidden component (display: none, rendered for print only)
  - ⬜ 80mm thermal paper layout (width: 72mm printable)
  - ⬜ Sections: restaurant name, date/time, order number, table, items table, totals, payment method, footer
  - ⬜ `@media print` CSS to isolate receipt and hide other content
  - ⬜ Monospace font for alignment
  - ⬜ Bilingual footer: "Thank you! / شكراً لزيارتكم"
  - ⬜ Props: order data (number, items, totals, payment, table, cashier name)

- ⬜ Create `src/hooks/usePOSPrint.ts`
  - ⬜ Implement `printReceipt(orderData)` — creates hidden iframe with receipt HTML
  - ⬜ Apply print-specific CSS via style tag in iframe
  - ⬜ Call `iframe.contentWindow.print()`
  - ⬜ Remove iframe after print dialog closes (onafterprint or timeout)
  - ⬜ Implement `autoPrintEnabled` setting (stored in posStore or localStorage)

### 3.4 Wire Into POS Page

- ⬜ Modify `src/app/[locale]/dashboard/pos/page.tsx`
  - ⬜ Wire `usePOSBarcodeScanner()` hook with product index
  - ⬜ Wire `usePOSAudio()` hook
  - ⬜ Play 'add' sound on addItem
  - ⬜ Play 'remove' sound on removeItem
  - ⬜ Play 'success' sound on order submission success
  - ⬜ Play 'error' sound on order submission failure
  - ⬜ Auto-print receipt after successful order (if autoPrint enabled)

### 3.5 Phase 3 Testing

- ⬜ **Barcode Scanner Testing:**
  - ⬜ Simulate barcode: type "123456" rapidly (<50ms between keys) + Enter → product added
  - ⬜ Simulate slow typing: type "123456" slowly (>50ms between keys) + Enter → NOT treated as barcode
  - ⬜ Unknown barcode: type valid-length unknown code + Enter → toast "Product not found"
  - ⬜ Short input: type "123" + Enter → NOT treated as barcode (< 6 chars)
  - ⬜ Scanner while in search input → NOT captured (ignored in INPUT elements)
  - ⬜ Real barcode scanner device test (if available)

- ⬜ **Audio Testing:**
  - ⬜ First click on page → audio unlocked (no errors in console)
  - ⬜ Add item → beep sound plays
  - ⬜ Remove item → lower beep plays
  - ⬜ Submit order → success chime plays
  - ⬜ Submission error → error tone plays
  - ⬜ Toggle sound off → no sounds play
  - ⬜ Toggle sound on → sounds resume
  - ⬜ Test on mobile browser → sounds work after first tap
  - ⬜ Sound files missing → fallback oscillator tone plays

- ⬜ **Receipt Printing:**
  - ⬜ Submit order → browser print dialog opens
  - ⬜ Verify receipt layout: restaurant name, date, order #, items, totals, payment
  - ⬜ Verify 80mm width formatting
  - ⬜ Test RTL receipt (Arabic locale)
  - ⬜ Test auto-print toggle: off → no print dialog, on → dialog opens
  - ⬜ Test canceling print dialog → no errors

---

## Phase 4: Advanced Performance & PWA

### 4.1 Virtual Scrolling Grid

- ⬜ Install `@tanstack/react-virtual` dependency
- ⬜ Create `src/app/[locale]/dashboard/pos/components/POSVirtualGrid.tsx`
  - ⬜ Use `useVirtualizer` from @tanstack/react-virtual
  - ⬜ Fixed card size: 130x140px (matches CSS grid minmax)
  - ⬜ ResizeObserver on container to compute column count
  - ⬜ Calculate row count: Math.ceil(products.length / columns)
  - ⬜ Render only visible rows + 2 overscan rows
  - ⬜ Each row renders N product cards (N = column count)
  - ⬜ Smooth scrolling at 60fps
  - ⬜ Keyboard shortcut badges still work (based on overall index)

- ⬜ Modify `src/app/[locale]/dashboard/pos/components/POSProductGrid.tsx`
  - ⬜ If products.length > 50: render `<POSVirtualGrid />` instead of CSS grid
  - ⬜ If products.length <= 50: keep regular CSS grid (simpler, no overhead)

### 4.2 Service Worker

- ⬜ Create `public/pos-sw.js`
  - ⬜ Cache name: `'pos-cache-v1'`
  - ⬜ Install event: precache POS page bundle, sound files, critical CSS
  - ⬜ Activate event: clean old caches (different version)
  - ⬜ Fetch handler:
    - ⬜ Static assets (JS/CSS/fonts/images/sounds): cache-first strategy
    - ⬜ API calls (`/api/v1/*`): network-first, fallback to cache for GET only
    - ⬜ Navigation requests: network-first
  - ⬜ Background Sync: register `'pos-order-sync'` event
  - ⬜ Sync handler: read IndexedDB pending-orders, replay via fetch()

- ⬜ Modify `src/app/[locale]/dashboard/pos/page.tsx`
  - ⬜ Register Service Worker on mount: `navigator.serviceWorker.register('/pos-sw.js')`
  - ⬜ Only register in production (check `process.env.NODE_ENV`)
  - ⬜ Handle registration errors gracefully (log, don't crash)

### 4.3 Background Sync

- ⬜ Create `src/hooks/usePOSBackgroundSync.ts`
  - ⬜ Register Background Sync when orders queued offline: `registration.sync.register('pos-order-sync')`
  - ⬜ Feature detection: check `'sync' in registration`
  - ⬜ Fallback: if Background Sync not supported, use manual replay (Phase 2)
  - ⬜ SW sync event handler replays IndexedDB queue via fetch()
  - ⬜ On sync success: clear queue, post message to client
  - ⬜ On sync failure: sync will be retried by browser

### 4.4 Performance Monitoring

- ⬜ Create `src/hooks/usePOSPerformance.ts`
  - ⬜ Use existing `recordPerformanceMetric()` from `src/lib/performance.ts`
  - ⬜ Track `pos.page-load`: mark from navigation start to first render
  - ⬜ Track `pos.add-to-cart`: mark around addItem calls
  - ⬜ Track `pos.search`: mark around search filter computation
  - ⬜ Track `pos.order-submit`: mark around API call
  - ⬜ Track `pos.product-fetch`: mark around data loading
  - ⬜ Dev-mode: floating performance badge showing last 5 metrics
  - ⬜ Production: metrics silently recorded, available via `window.__posMetrics`

### 4.5 Phase 4 Testing

- ⬜ **Virtual Scrolling:**
  - ⬜ Load POS with 500+ products — smooth scrolling at 60fps (DevTools Performance tab)
  - ⬜ Scroll quickly — no blank areas visible
  - ⬜ Resize browser window — grid reflows correctly
  - ⬜ Category filter with 500+ products — instant filter + scroll to top
  - ⬜ Search with 500+ products — results appear in <50ms
  - ⬜ F1–F12 shortcuts work correctly with virtual grid
  - ⬜ Products < 50 — regular CSS grid used (no virtual overhead)

- ⬜ **Service Worker:**
  - ⬜ SW installed (DevTools > Application > Service Workers)
  - ⬜ Static assets cached (DevTools > Application > Cache Storage)
  - ⬜ Page loads offline (after first visit) — products from cache
  - ⬜ API calls go network-first (verify in Network tab)
  - ⬜ SW update: old cache cleaned on activate

- ⬜ **Background Sync:**
  - ⬜ Queue order offline → close browser → reopen → order synced
  - ⬜ If Background Sync not supported → manual replay works as fallback
  - ⬜ Multiple queued orders → all synced in order

- ⬜ **Performance Metrics:**
  - ⬜ Dev mode: performance badge visible with metrics
  - ⬜ Page load metric: < 200ms (cached), < 1s (cold)
  - ⬜ Add-to-cart metric: < 20ms
  - ⬜ Search metric: < 50ms
  - ⬜ Order submit metric: < 500ms
  - ⬜ No console errors in production mode

---

## Cross-Phase: Security Checklist

- ⬜ **Authentication:**
  - ⬜ RBAC check on page load: ORDERS module + WRITE permission
  - ⬜ All API calls use Bearer token via api.ts interceptor
  - ⬜ 401 response → user redirected to login
  - ⬜ Access denied screen for unauthorized users
  - ⬜ Token expiry during offline → replay fails gracefully, user re-authenticates

- ⬜ **Input Validation:**
  - ⬜ Backend recalculates totals from product_id + quantity (client totals are display-only)
  - ⬜ Customer name/phone: trimmed, no HTML injection possible (React auto-escapes)
  - ⬜ Order notes: free-text, backend enforces max length
  - ⬜ Table ID validated against available tables from backend
  - ⬜ Product quantities: positive integers only, capped at inventory stock
  - ⬜ Payment method: validated against allowed values ('cash', 'card', 'wallet')

- ⬜ **XSS Prevention:**
  - ⬜ No `dangerouslySetInnerHTML` anywhere in POS code
  - ⬜ All user inputs rendered via React JSX (auto-escaped)
  - ⬜ Product names from API rendered as text content
  - ⬜ Receipt HTML generated via React components (not string templates)
  - ⬜ Search query used only for includes() filtering (not injected into DOM)

- ⬜ **Data Security:**
  - ⬜ No passwords or tokens stored in IndexedDB
  - ⬜ IndexedDB data is browser-sandboxed (same-origin policy)
  - ⬜ localStorage data: only cart items, product cache (no sensitive data)
  - ⬜ Service Worker: scoped to POS routes only
  - ⬜ No sensitive data in SW cache

- ⬜ **CSRF:**
  - ⬜ Bearer token auth (not cookies) — inherently CSRF-resistant
  - ⬜ No cookie-based state-changing requests

- ⬜ **Rate Limiting:**
  - ⬜ Backend rate limits apply to POS API calls (existing middleware)
  - ⬜ Client-side: prevent double-submit with `isSubmitting` flag
  - ⬜ Offline queue replay: sequential (not parallel) to avoid overwhelming backend

---

## Cross-Phase: Performance Optimization Checklist

- ⬜ **Rendering:**
  - ⬜ `React.memo` on: POSProductCard, POSCartItem, POSCategoryTabs
  - ⬜ Zustand selectors on all components (subscribe to specific slices)
  - ⬜ `useMemo` for filtered product list (recompute only on category/search change)
  - ⬜ `useCallback` for all event handlers passed as props
  - ⬜ `useRef` for callback deps in keyboard/scanner hooks (avoid stale closures)
  - ⬜ No unnecessary re-renders: verify with React DevTools Profiler

- ⬜ **Layout:**
  - ⬜ CSS `contain: content` on product grid container
  - ⬜ Fixed card sizes for predictable layout
  - ⬜ `loading="lazy"` on product images
  - ⬜ `tabular-nums` for price/number displays (prevent layout shift)

- ⬜ **Data:**
  - ⬜ Products served from cache in <5ms (IndexedDB or localStorage)
  - ⬜ Background refresh: UI never blocks on API fetch
  - ⬜ 50ms search debounce (avoid excessive re-filtering)
  - ⬜ Single fetch for all products (`page_size=500`)
  - ⬜ Parallel API fetches via Promise.allSettled

- ⬜ **Network:**
  - ⬜ Service Worker caches static assets (Phase 4)
  - ⬜ Network-first for API (always fresh data when online)
  - ⬜ Offline queue: no unnecessary network requests
  - ⬜ Product images: lazy-loaded, cached by SW

- ⬜ **Memory:**
  - ⬜ Virtual scrolling for 500+ products (Phase 4)
  - ⬜ Audio elements properly cleaned up on unmount
  - ⬜ IndexedDB connections closed after each operation
  - ⬜ No memory leaks: verify with DevTools Memory tab

---

## Cross-Phase: Production Readiness Checklist

- ⬜ **Error Handling:**
  - ⬜ All API calls wrapped in try/catch with user-friendly error toasts
  - ⬜ IndexedDB operations: graceful fallback if unavailable
  - ⬜ localStorage operations: graceful fallback if full/unavailable
  - ⬜ Service Worker: graceful degradation if not supported
  - ⬜ Audio: fallback oscillator if sound files fail to load
  - ⬜ Barcode scanner: no crash on unexpected input patterns
  - ⬜ Corrupted persisted state: normalizePersistedState() recovers to defaults

- ⬜ **Logging:**
  - ⬜ All errors logged via `debugError()` (development only)
  - ⬜ Performance metrics recorded via `recordPerformanceMetric()`
  - ⬜ No console.log in production (only debugLog/debugError which are dev-only)

- ⬜ **i18n:**
  - ⬜ All user-visible strings use `t()` translation function
  - ⬜ Full Arabic translations for all ~30 keys
  - ⬜ RTL layout tested for all components
  - ⬜ Date/time formatting uses locale-aware API
  - ⬜ Receipt supports bilingual footer

- ⬜ **Responsive Design:**
  - ⬜ POS optimized for tablet (768px+) and desktop
  - ⬜ Minimum touch targets: 44x44px
  - ⬜ Product grid: auto-fill responsive columns
  - ⬜ Cart panel: scrollable when many items
  - ⬜ Viewport meta prevents zoom (kiosk mode)

- ⬜ **Browser Compatibility:**
  - ⬜ Chrome 90+ (primary target for POS terminals)
  - ⬜ Safari 15+ (iPad POS use case)
  - ⬜ Firefox 90+ (fallback)
  - ⬜ WebKit AudioContext fallback for Safari
  - ⬜ IndexedDB polyfill not needed (all target browsers support it)

- ⬜ **Deployment:**
  - ⬜ No new environment variables required
  - ⬜ No new backend migrations required
  - ⬜ No new npm dependencies for Phases 1-3
  - ⬜ Phase 4: add `@tanstack/react-virtual` to package.json
  - ⬜ Service Worker file in `public/` directory (served as static asset)
  - ⬜ Sound files in `public/sounds/` directory
  - ⬜ All new files are within existing dashboard app — no new build targets

---

## Summary

| Phase | New Files | Modified Files | New Dependencies | Estimated Complexity |
|-------|-----------|----------------|------------------|--------------------|
| **1. Core UI** | 19 | 5 | 0 | High (core feature) |
| **2. Offline** | 3 | 2 | 0 | Medium (IndexedDB patterns exist in KDS) |
| **3. Hardware** | 4 + sounds | 1 | 0 | Medium (audio patterns exist in KDS) |
| **4. Performance** | 4 | 2 | 1 (@tanstack/react-virtual) | Medium (virtual grid + SW) |
| **Total** | ~30 files | ~6 files | 1 package | — |
