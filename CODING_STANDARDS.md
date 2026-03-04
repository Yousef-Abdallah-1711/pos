# Coding Standards & Best Practices

**Version:** 1.0.0
**Scope:** Universal — applies to ALL projects (Frontend, Backend, Database, Security, Performance)
**Last Updated:** 2026-03-04

> **Read this BEFORE writing any code.** Every rule here exists because a real bug, performance issue, or security vulnerability was found in production.

---

## Table of Contents

1. [Quick PR Self-Review Checklist](#1-quick-pr-self-review-checklist)
2. [Code Quality Rules](#2-code-quality-rules)
3. [Frontend Standards (React / Next.js)](#3-frontend-standards-react--nextjs)
4. [Backend Standards (Go)](#4-backend-standards-go)
5. [Database Standards (SQL / PostgreSQL)](#5-database-standards-sql--postgresql)
6. [Performance Rules](#6-performance-rules)
7. [Security Rules](#7-security-rules)
8. [API Design Standards](#8-api-design-standards)
9. [Testing Standards](#9-testing-standards)
10. [Git & Code Review Standards](#10-git--code-review-standards)

---

## 1. Quick PR Self-Review Checklist

Before submitting ANY pull request, verify every item:

### Performance
- [ ] No synchronous imports of heavy libraries (recharts, pdf libs, date-fns locale)
- [ ] All lists with 50+ items use pagination (frontend AND backend)
- [ ] No N+1 queries — use JOINs or batch fetches
- [ ] No `useEffect` that fires on every render (missing/wrong dependency array)
- [ ] No polling without `document.visibilityState` check
- [ ] Images use `next/image` or have explicit `width`/`height`/`loading="lazy"`

### Security
- [ ] No hardcoded secrets, tokens, or API keys
- [ ] All user input is validated (frontend: Zod, backend: explicit checks)
- [ ] No `dangerouslySetInnerHTML` without DOMPurify sanitization
- [ ] No `eval()`, `new Function()`, or template literal SQL
- [ ] Auth tokens stored in httpOnly cookies (NOT localStorage for production)
- [ ] All SQL queries use parameterized placeholders (`$1`, `$2` — NEVER string concat)

### Code Quality
- [ ] No `any` type in TypeScript (use `unknown` + type guards)
- [ ] No `eslint-disable` or `@ts-ignore` without a comment explaining WHY
- [ ] No `console.log` left in production code (use `debugLog` wrapper)
- [ ] Every function under 50 lines; every file under 400 lines
- [ ] Error states handled — no silent `catch {}` blocks
- [ ] Loading states shown — no blank screens while fetching

### Database
- [ ] New filtered columns have indexes
- [ ] Migrations are reversible (include rollback SQL)
- [ ] No `SELECT *` — list explicit columns
- [ ] `tenant_id` in EVERY WHERE clause for multi-tenant tables

### Accessibility
- [ ] Color contrast checked (4.5:1 minimum ratio)
- [ ] All interactive elements have `aria-label` or visible text
- [ ] Keyboard navigation tested (Tab, Enter, Escape — no focus traps)
- [ ] Screen reader tested on critical flows (VoiceOver / NVDA)

### Mobile
- [ ] Touch targets 48x48px minimum (no fat-finger misses)
- [ ] `safe-area-inset-bottom` tested on iOS Safari
- [ ] RTL layout verified (swipe gestures, text alignment)
- [ ] 3G throttle test passes (page load < 5s)

### Perceived Performance
- [ ] Skeleton screens match real content layout (no spinner circles)
- [ ] Optimistic updates for cart/favorites (rollback on error)
- [ ] Button feedback: idle → loading → success/error (< 500ms)
- [ ] Dynamic imports for all heavy components (> 50KB)

### Deployment
- [ ] `lint && test && build` all pass
- [ ] Database migration tested on staging
- [ ] Rollback plan documented
- [ ] Smoke tests pass on staging URL

---

## 2. Code Quality Rules

### 2.1 Naming Conventions

```
Files:
  Frontend:  PascalCase.tsx (components), camelCase.ts (utils/hooks)
  Backend:   snake_case.go
  Database:  snake_case (tables, columns)
  CSS:       kebab-case

Variables/Functions:
  TypeScript:  camelCase
  Go:          camelCase (exported: PascalCase)
  SQL:         snake_case

Constants:
  TypeScript:  UPPER_SNAKE_CASE
  Go:          PascalCase (exported), camelCase (unexported)
```

**Rule:** Names must reveal intent. If you need a comment to explain what a variable holds, rename it.

```typescript
// BAD
const d = new Date();
const arr = users.filter(u => u.a);
const temp = calculateTotal(items);

// GOOD
const orderCreatedAt = new Date();
const activeUsers = users.filter(user => user.isActive);
const cartTotal = calculateTotal(cartItems);
```

```go
// BAD
func Do(t int, r int) error { ... }

// GOOD
func AssignDriverToOrder(tenantID int, restaurantID int) error { ... }
```

### 2.2 Function Size & Complexity

**Rule:** One function = one responsibility. If you can't name it clearly, split it.

| Metric                | Limit    |
|-----------------------|----------|
| Lines per function    | ≤ 50     |
| Parameters            | ≤ 4      |
| Nesting depth         | ≤ 3      |
| Lines per file        | ≤ 400    |
| Cyclomatic complexity | ≤ 10     |

```typescript
// BAD — 80+ lines, multiple responsibilities
function handleSubmit(formData: FormData) {
  // validate... 20 lines
  // transform... 15 lines
  // API call... 10 lines
  // update state... 15 lines
  // show notification... 10 lines
  // navigate... 5 lines
}

// GOOD — each responsibility extracted
function handleSubmit(formData: FormData) {
  const validated = validateOrderForm(formData);
  if (!validated.success) return showErrors(validated.errors);

  const payload = transformToCreateOrderRequest(validated.data);
  const order = await createOrder(payload);

  showNotification('Order created successfully');
  router.push(`/orders/${order.id}`);
}
```

### 2.3 Comments — When and How

**Rule:** Code tells you WHAT. Comments tell you WHY.

```typescript
// BAD — comment says what the code already says
// Set isLoading to true
setIsLoading(true);

// BAD — no comment on non-obvious logic
if (order.total > 500 && order.paymentMethod === 'cash') {
  requireManagerApproval();
}

// GOOD — explains the business rule
// Cash orders over 500 EGP require manager approval per company policy
if (order.total > 500 && order.paymentMethod === 'cash') {
  requireManagerApproval();
}
```

```go
// GOOD — explains WHY, not WHAT
// We retry 3 times because Paymob's API occasionally returns 503
// during their daily maintenance window (02:00-02:15 UTC)
for attempt := 0; attempt < 3; attempt++ {
    resp, err = paymobClient.InitiatePayment(req)
    if err == nil { break }
    time.Sleep(time.Second * time.Duration(attempt+1))
}
```

**Always comment:**
- Business rules and domain logic
- Workarounds with links to the issue/ticket
- Regex patterns (explain what they match)
- Magic numbers (or better: extract to named constants)
- Performance-critical sections (explain why this approach)

### 2.4 Error Handling

**Rule:** Never swallow errors. Every error must be logged, reported, or surfaced to the user.

```typescript
// BAD — silent catch
try {
  await api.delete(`/orders/${id}`);
} catch (e) {}

// BAD — generic error
try {
  await api.delete(`/orders/${id}`);
} catch (e) {
  console.log(e);
}

// GOOD — specific handling + user feedback
try {
  await api.delete(`/orders/${id}`);
  toast.success('Order deleted');
} catch (error) {
  if (isAxiosError(error) && error.response?.status === 404) {
    toast.error('Order not found — it may have been already deleted');
  } else if (isAxiosError(error) && error.response?.status === 403) {
    toast.error('You do not have permission to delete orders');
  } else {
    toast.error('Failed to delete order. Please try again.');
    debugError('deleteOrder failed:', error);
  }
}
```

```go
// BAD — error ignored
json.NewDecoder(r.Body).Decode(&req)

// GOOD — error checked and returned
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
    return
}
```

### 2.5 Don't Repeat Yourself (DRY) — But Don't Over-Abstract

**Rule:** Duplicate code 2x = acceptable. Duplicate 3x = extract. But NEVER abstract prematurely.

```typescript
// BAD — premature abstraction for 1 use case
const useGenericCRUD = <T>(endpoint: string, options?: CRUDOptions<T>) => { ... } // 200 lines

// GOOD — simple, direct hooks
const useProducts = () => useSWR<Product[]>('/api/v1/products');
const useCategories = () => useSWR<Category[]>('/api/v1/categories');
// Only abstract when you have 3+ similar patterns AND they share real logic
```

### 2.6 TypeScript Strictness

**Rule:** `strict: true` in tsconfig. No exceptions.

```typescript
// BAD — using 'any' defeats the purpose of TypeScript
function processData(data: any) { return data.items.map((i: any) => i.name); }

// GOOD — explicit types
interface OrderResponse { items: OrderItem[]; total: number; }
function processData(data: OrderResponse): string[] {
  return data.items.map(item => item.name);
}

// BAD — type assertion without validation
const user = JSON.parse(token) as User;

// GOOD — runtime validation with Zod
const userSchema = z.object({ id: z.number(), name: z.string(), role: z.string() });
const user = userSchema.parse(JSON.parse(token));
```

---

## 3. Frontend Standards (React / Next.js)

### 3.1 Component Structure

**Rule:** Every component file follows this order:

```typescript
// 1. Imports (external → internal → types → styles)
import { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import type { Order } from '@/types';

// 2. Types/Interfaces (if component-specific)
interface OrderCardProps {
  order: Order;
  onStatusChange: (id: number, status: string) => void;
}

// 3. Component
export function OrderCard({ order, onStatusChange }: OrderCardProps) {
  // hooks first
  const router = useRouter();
  const [isUpdating, setIsUpdating] = useState(false);

  // derived state
  const canCancel = order.status === 'pending' && !isUpdating;

  // handlers
  const handleCancel = useCallback(async () => {
    setIsUpdating(true);
    await onStatusChange(order.id, 'cancelled');
    setIsUpdating(false);
  }, [order.id, onStatusChange]);

  // render
  return ( ... );
}
```

### 3.2 Dynamic Imports for Heavy Libraries

**Rule:** Any library >50KB must be dynamically imported. NEVER import at the top level.

```typescript
// BAD — recharts is 400KB, loaded on every page
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip } from 'recharts';

// GOOD — loaded only when the chart renders
import dynamic from 'next/dynamic';

const SalesChart = dynamic(
  () => import('@/components/reports/SalesChart'),
  {
    loading: () => <div className="h-64 animate-pulse bg-muted rounded" />,
    ssr: false,
  }
);
```

**Libraries that MUST be dynamically imported:**
| Library | Size | When to load |
|---------|------|-------------|
| recharts | ~400KB | Only on report pages |
| html2canvas + jspdf | ~350KB | Only on PDF export click |
| date-fns/locale/* | ~50-200KB | Only the active locale |
| @tanstack/react-table | ~50KB | Only on table-heavy pages |
| DOMPurify | ~30KB | Only when rendering user HTML |
| monaco-editor | ~2MB | Only in code editors |
| xlsx/exceljs | ~300KB | Only on Excel export click |

### 3.3 State Management Rules

**Rule:** Use the simplest state that works. Escalate only when needed.

```
Local state (useState)      → Component-only UI state (open/close, form fields)
Derived state (useMemo)     → Computed from props/state (filtered lists, totals)
URL state (searchParams)    → Shareable state (filters, pagination, tabs)
Zustand store               → Cross-component state (auth, cart, preferences)
Server state (React Query)  → Remote data (API responses, cache invalidation)
```

```typescript
// BAD — using Zustand for local form state
const useFormStore = create((set) => ({
  name: '', setName: (name) => set({ name }),
  email: '', setEmail: (email) => set({ email }),
}));

// GOOD — local state for forms
function CreateProductForm() {
  const { register, handleSubmit } = useForm<CreateProductInput>();
  // ...
}

// BAD — using useState for server data
const [products, setProducts] = useState<Product[]>([]);
useEffect(() => { fetchProducts().then(setProducts); }, []);

// GOOD — React Query (or SWR) for server state
const { data: products, isLoading } = useQuery({
  queryKey: ['products', tenantId],
  queryFn: () => api.get('/products').then(r => r.data),
});
```

### 3.4 useEffect Discipline

**Rule:** `useEffect` is for synchronizing with external systems. NOT for derived state, NOT for event handling.

```typescript
// BAD — derived state in useEffect
const [filteredItems, setFilteredItems] = useState([]);
useEffect(() => {
  setFilteredItems(items.filter(i => i.category === selected));
}, [items, selected]);

// GOOD — derive during render
const filteredItems = useMemo(
  () => items.filter(i => i.category === selected),
  [items, selected]
);

// BAD — missing dependency causes stale closure
useEffect(() => {
  const interval = setInterval(() => {
    fetchOrders(restaurantId); // restaurantId may be stale!
  }, 30000);
  return () => clearInterval(interval);
}, []); // Missing restaurantId in deps!

// GOOD — all dependencies declared
useEffect(() => {
  const interval = setInterval(() => fetchOrders(restaurantId), 30000);
  return () => clearInterval(interval);
}, [restaurantId]);
```

### 3.5 Loading & Error States

**Rule:** Every async operation must show: Loading → Data → Error. No blank screens.

```typescript
// BAD — blank screen while loading
function OrdersPage() {
  const { data } = useOrders();
  return <OrderList orders={data} />;
}

// GOOD — all states handled
function OrdersPage() {
  const { data, isLoading, error } = useOrders();

  if (isLoading) return <OrderListSkeleton />;
  if (error) return <ErrorState message="Failed to load orders" onRetry={refetch} />;
  if (!data?.length) return <EmptyState title="No orders yet" />;

  return <OrderList orders={data} />;
}
```

**Next.js App Router:** Always add `loading.tsx` and `error.tsx` to route segments:

```
app/[locale]/dashboard/orders/
├── page.tsx        ← Main content
├── loading.tsx     ← Skeleton shown during navigation
├── error.tsx       ← Error boundary (must be 'use client')
└── layout.tsx      ← Shared layout
```

### 3.6 Perceived Performance (Skeletons & Optimistic Updates)

**Rule:** Users judge speed by "feel", not by actual load time. Skeleton screens + optimistic updates make apps feel 3x faster.

**Skeleton screens (NOT spinners):**

```typescript
// BAD — spinner gives no spatial hint of what's loading
if (isLoading) return <Spinner />;

// GOOD — skeleton mirrors the real layout, prevents layout shift
if (isLoading) return (
  <div className="space-y-4">
    <div className="h-8 w-48 animate-pulse bg-muted rounded" />
    <div className="grid grid-cols-3 gap-4">
      {Array.from({ length: 6 }).map((_, i) => (
        <div key={i} className="h-40 animate-pulse bg-muted rounded" />
      ))}
    </div>
  </div>
);
```

**Skeleton priority order:**
1. Navigation / Sidebar — always visible first
2. Page header + breadcrumbs
3. Primary content area (table rows, card grids)
4. Secondary content (charts, stats widgets)

**Optimistic updates (feels instant):**

```typescript
// BAD — wait for API roundtrip (user sees 800ms delay)
const addToCart = async (productId: number) => {
  setIsLoading(true);
  await api.post('/cart', { productId });
  await refetchCart();
  setIsLoading(false);
};

// GOOD — update UI instantly, sync in background, rollback on error
const addToCart = async (productId: number) => {
  const previousCart = cartStore.getState().items;

  // 1. Optimistic update (0ms — user sees it immediately)
  cartStore.addItem(productId);

  // 2. Background sync
  try {
    await api.post('/cart', { productId });
  } catch {
    // 3. Rollback on failure
    cartStore.setItems(previousCart);
    toast.error('Failed to add item');
  }
};
```

**Button feedback pattern:**

```typescript
// idle → loading (spinner) → success (checkmark) → idle
const [state, setState] = useState<'idle' | 'loading' | 'success' | 'error'>('idle');

const handleClick = async () => {
  setState('loading');
  try {
    await api.post('/orders', data);
    setState('success');
    setTimeout(() => setState('idle'), 1500);
  } catch {
    setState('error');
    setTimeout(() => setState('idle'), 2000);
  }
};
```

**Parallel data fetching (Next.js App Router):**

```typescript
// BAD — sequential fetches (900ms total)
const theme = await getTheme(slug);       // 200ms
const products = await getProducts(slug);  // 300ms
const categories = await getCategories(slug); // 150ms

// GOOD — parallel fetches (300ms total — slowest wins)
const [theme, products, categories] = await Promise.all([
  getTheme(slug),
  getProducts(slug),
  getCategories(slug),
]);
```

### 3.7 Memoization — Use Sparingly, Correctly

**Rule:** Only memoize when you have MEASURED a performance problem. Premature memoization adds complexity for zero benefit.

```typescript
// WHEN to use React.memo:
// - Component re-renders often with the SAME props
// - Component render is expensive (large list, heavy computation)
// - Component is a list item rendered 50+ times

// WHEN NOT to use React.memo:
// - Component always receives new props (objects/arrays created inline)
// - Component is lightweight and renders fast
// - You haven't measured a performance issue

// BAD — memo is useless because `style` is a new object every render
const Button = React.memo(({ label, style }: Props) => (
  <button style={style}>{label}</button>
));
<Button label="Save" style={{ color: 'blue' }} /> // new object every render!

// GOOD — stable references
const buttonStyle = useMemo(() => ({ color: 'blue' }), []);
<Button label="Save" style={buttonStyle} />
```

### 3.7 Polling & Real-Time Data

**Rule:** All polling MUST pause when the tab is hidden. Wasting bandwidth on invisible tabs kills performance.

```typescript
// BAD — polls forever even when tab is hidden
useEffect(() => {
  const interval = setInterval(fetchOrders, 30000);
  return () => clearInterval(interval);
}, []);

// GOOD — pause polling when tab is not visible
useEffect(() => {
  let interval: NodeJS.Timeout;

  const startPolling = () => {
    interval = setInterval(fetchOrders, 30000);
  };

  const handleVisibility = () => {
    clearInterval(interval);
    if (document.visibilityState === 'visible') {
      fetchOrders(); // immediate refresh on tab focus
      startPolling();
    }
  };

  document.addEventListener('visibilitychange', handleVisibility);
  startPolling();

  return () => {
    clearInterval(interval);
    document.removeEventListener('visibilitychange', handleVisibility);
  };
}, []);
```

### 3.8 Barrel Exports (index.ts)

**Rule:** AVOID barrel exports for large modules. They kill tree-shaking.

```typescript
// BAD — barrel re-exports everything, bundler can't tree-shake
// components/reports/index.ts
export { SalesReport } from './SalesReport';
export { ProfitReport } from './ProfitReport';
export { InventoryReport } from './InventoryReport';
export { HRReport } from './HRReport';
// Importing SalesReport pulls ALL reports into the bundle

// GOOD — import directly from the source file
import { SalesReport } from '@/components/reports/SalesReport';
```

### 3.9 Internationalization (i18n)

**Rule:** Load only the active locale. Never bundle all locales together.

```typescript
// BAD — both locales loaded statically
import en from '@/i18n/messages/en.json';
import ar from '@/i18n/messages/ar.json';
const messages = locale === 'ar' ? ar : en;

// GOOD — dynamic import of active locale only
const messages = await import(`@/i18n/messages/${locale}.json`).then(m => m.default);

// Or in Next.js layout:
async function getMessages(locale: string) {
  return (await import(`@/i18n/messages/${locale}.json`)).default;
}
```

### 3.10 CSS & Tailwind Rules

```typescript
// BAD — inline styles that bypass the design system
<div style={{ padding: '16px', color: '#333', fontSize: '14px' }}>

// GOOD — Tailwind utilities (consistent, responsive, dark-mode aware)
<div className="p-4 text-foreground text-sm">

// BAD — duplicated Tailwind classes across components
// ButtonPrimary.tsx
<button className="px-4 py-2 bg-primary text-white rounded-md font-medium hover:bg-primary/90">
// ButtonSecondary.tsx
<button className="px-4 py-2 bg-secondary text-white rounded-md font-medium hover:bg-secondary/90">

// GOOD — extract to shared component or cn() utility
<Button variant="primary">Save</Button>
<Button variant="secondary">Cancel</Button>
```

**Responsive design rules:**
- Mobile-first: write `text-sm` then add `md:text-base`
- Test safe-area-inset on iOS Safari: `pb-[env(safe-area-inset-bottom)]`
- Never use fixed pixel widths for containers: use `max-w-screen-xl mx-auto`

---

## 4. Backend Standards (Go)

### 4.1 Handler Structure

**Rule:** Every HTTP handler follows: Parse → Validate → Execute → Respond.

```go
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
    // 1. Parse context (tenant, user from middleware)
    tenantID := r.Context().Value("tenant_id").(int)
    userID := r.Context().Value("user_id").(int)

    // 2. Parse & validate request body
    var req domain.CreateOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    if err := req.Validate(); err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    // 3. Execute business logic (delegate to usecase)
    order, err := h.orderUC.Create(r.Context(), tenantID, userID, &req)
    if err != nil {
        handleUsecaseError(w, err)
        return
    }

    // 4. Respond
    respondJSON(w, http.StatusCreated, order)
}
```

### 4.2 Always Validate Input

**Rule:** Trust NOTHING from the client. Validate every field at the handler layer.

```go
// BAD — no validation
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req domain.CreateProductRequest
    json.NewDecoder(r.Body).Decode(&req) // error ignored!
    product, _ := h.productUC.Create(req) // error ignored!
    respondJSON(w, http.StatusCreated, product)
}

// GOOD — full validation
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req domain.CreateProductRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid JSON")
        return
    }

    // Validate required fields
    if strings.TrimSpace(req.Name) == "" {
        respondError(w, http.StatusBadRequest, "name is required")
        return
    }
    if req.Price < 0 {
        respondError(w, http.StatusBadRequest, "price must be non-negative")
        return
    }
    if req.CategoryID <= 0 {
        respondError(w, http.StatusBadRequest, "valid category_id is required")
        return
    }

    product, err := h.productUC.Create(r.Context(), &req)
    if err != nil {
        handleUsecaseError(w, err)
        return
    }
    respondJSON(w, http.StatusCreated, product)
}
```

### 4.3 Pagination is Mandatory

**Rule:** Every list endpoint MUST support pagination. Default: `page=1, limit=25, max_limit=100`.

```go
// BAD — returns ALL records. Will crash with 10k+ rows.
func (r *OrderRepo) List(tenantID int) ([]domain.Order, error) {
    rows, err := r.db.Query("SELECT * FROM orders WHERE tenant_id = $1", tenantID)
    // returns ALL orders...
}

// GOOD — paginated with enforced limits
func (r *OrderRepo) List(tenantID int, page, limit int) ([]domain.Order, int, error) {
    if limit <= 0 || limit > 100 {
        limit = 25
    }
    if page <= 0 {
        page = 1
    }
    offset := (page - 1) * limit

    // Count total for pagination metadata
    var total int
    err := r.db.QueryRow(
        "SELECT COUNT(*) FROM orders WHERE tenant_id = $1", tenantID,
    ).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("count orders: %w", err)
    }

    rows, err := r.db.Query(
        `SELECT id, status, total, created_at FROM orders
         WHERE tenant_id = $1
         ORDER BY created_at DESC
         LIMIT $2 OFFSET $3`,
        tenantID, limit, offset,
    )
    // ...
    return orders, total, nil
}
```

### 4.4 Error Wrapping

**Rule:** Wrap errors with context at each layer. Use `fmt.Errorf("context: %w", err)`.

```go
// BAD — raw error, impossible to debug
return err

// GOOD — wrapped with context at each layer
// Repository
return fmt.Errorf("repo.GetOrder(id=%d): %w", orderID, err)

// UseCase
return fmt.Errorf("orderUC.ConfirmPayment: %w", err)

// Handler converts to HTTP response
func handleUsecaseError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, domain.ErrNotFound):
        respondError(w, http.StatusNotFound, "resource not found")
    case errors.Is(err, domain.ErrForbidden):
        respondError(w, http.StatusForbidden, "access denied")
    case errors.Is(err, domain.ErrConflict):
        respondError(w, http.StatusConflict, "resource already exists")
    default:
        log.Printf("internal error: %v", err)
        respondError(w, http.StatusInternalServerError, "internal server error")
    }
}
```

### 4.5 Concurrent Queries

**Rule:** When a handler needs 3+ independent queries, run them concurrently with `errgroup`.

```go
// BAD — sequential queries (3s total if each takes 1s)
employees, err := h.hrUC.ListEmployees(ctx, tenantID)
attendance, err := h.hrUC.GetAttendanceSummary(ctx, tenantID)
leaves, err := h.hrUC.GetLeaveSummary(ctx, tenantID)

// GOOD — concurrent queries (1s total)
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(r.Context())
var employees []domain.Employee
var attendance domain.AttendanceSummary
var leaves domain.LeaveSummary

g.Go(func() error {
    var err error
    employees, err = h.hrUC.ListEmployees(ctx, tenantID)
    return err
})
g.Go(func() error {
    var err error
    attendance, err = h.hrUC.GetAttendanceSummary(ctx, tenantID)
    return err
})
g.Go(func() error {
    var err error
    leaves, err = h.hrUC.GetLeaveSummary(ctx, tenantID)
    return err
})

if err := g.Wait(); err != nil {
    handleUsecaseError(w, err)
    return
}
```

### 4.6 Context & Timeouts

**Rule:** Every database query and external API call must have a timeout via `context`.

```go
// BAD — no timeout, query can hang forever
rows, err := db.Query("SELECT * FROM orders WHERE tenant_id = $1", tenantID)

// GOOD — 5-second timeout
ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
defer cancel()
rows, err := db.QueryContext(ctx, "SELECT * FROM orders WHERE tenant_id = $1", tenantID)
```

### 4.7 Connection Pool Configuration

```go
// Recommended production settings
db.SetMaxOpenConns(50)                  // Max simultaneous connections
db.SetMaxIdleConns(25)                  // Keep 25 idle for reuse
db.SetConnMaxLifetime(30 * time.Minute) // Prevent stale connections
db.SetConnMaxIdleTime(5 * time.Minute)  // Close idle connections after 5m
```

**Sizing formula:** `MaxOpenConns = (CPU cores * 2) + effective_spindle_count`
For cloud PostgreSQL with SSDs: start at 50, monitor `pg_stat_activity`.

---

## 5. Database Standards (SQL / PostgreSQL)

### 5.1 Query Rules

**Rule:** NEVER use `SELECT *`. Always list explicit columns.

```sql
-- BAD — fetches all columns including BLOBs, wastes bandwidth
SELECT * FROM products WHERE tenant_id = $1;

-- GOOD — only what you need
SELECT id, name, price, category_id, is_active, created_at
FROM products
WHERE tenant_id = $1 AND is_active = true
ORDER BY name ASC
LIMIT 25 OFFSET 0;
```

### 5.2 Index Strategy

**Rule:** Every column used in WHERE, JOIN, or ORDER BY should have an index.

```sql
-- Multi-tenant tables: ALWAYS composite index with tenant_id first
CREATE INDEX idx_orders_tenant_status ON orders(tenant_id, status);
CREATE INDEX idx_orders_tenant_date ON orders(tenant_id, restaurant_id, created_at DESC);
CREATE INDEX idx_products_tenant_category ON products(tenant_id, category_id, is_active);

-- Partial indexes for common queries (smaller, faster)
CREATE INDEX idx_orders_pending ON orders(tenant_id, created_at)
  WHERE status = 'pending';

-- Covering indexes to avoid table lookups
CREATE INDEX idx_products_listing ON products(tenant_id, category_id)
  INCLUDE (name, price, image_url);
```

**Index rules:**
- Composite indexes: most-selective column first (usually `tenant_id`)
- Add indexes BEFORE deploying features (in migration files)
- Use `EXPLAIN ANALYZE` to verify index usage
- Remove unused indexes (they slow down writes)

#### N+1 Query Anti-Pattern

```go
// BAD — N+1: 1 query for orders + N queries for items
orders, _ := repo.ListOrders(ctx, tenantID)
for _, o := range orders {
    items, _ := repo.GetOrderItems(ctx, o.ID) // N extra queries!
    o.Items = items
}

// GOOD — Single JOIN query (or 2 queries with IN clause)
query := `
  SELECT o.id, o.status, oi.product_id, oi.quantity, oi.price
  FROM orders o
  JOIN order_items oi ON oi.order_id = o.id
  WHERE o.tenant_id = $1 AND o.created_at >= $2
  ORDER BY o.created_at DESC`

// ALSO GOOD — 2-query batch approach
orderIDs := extractIDs(orders)
items, _ := repo.GetItemsByOrderIDs(ctx, orderIDs) // Single IN query
itemMap := groupByOrderID(items)
for i := range orders {
    orders[i].Items = itemMap[orders[i].ID]
}
```

#### Covering Index Design

```sql
-- A covering index includes ALL columns the query needs
-- so PostgreSQL reads ONLY the index, never touches the table

-- Query: SELECT name, price FROM products WHERE tenant_id = $1 AND category_id = $2 AND is_active = true
-- Covering index:
CREATE INDEX idx_products_listing ON products(tenant_id, category_id)
  INCLUDE (name, price) WHERE is_active = true;

-- Verify with EXPLAIN ANALYZE — look for "Index Only Scan"
EXPLAIN ANALYZE SELECT name, price FROM products WHERE tenant_id = 1 AND category_id = 5 AND is_active = true;
-- ✅ Index Only Scan using idx_products_listing
-- ❌ If you see "Index Scan" + "Heap Fetches: 5000" → index is NOT covering
```

### 5.3 Migration Standards

```sql
-- Every migration file must include:
-- 1. Forward migration (CREATE/ALTER)
-- 2. Comment with what it does and why

-- Migration: 160_add_order_notes_column.sql
-- Adds notes column to orders for special customer instructions

ALTER TABLE orders ADD COLUMN notes TEXT DEFAULT '';

-- Create index for text search on notes (if we plan to search)
-- Skipping for now — add when search feature is built

-- Rollback (document in rollback script):
-- ALTER TABLE orders DROP COLUMN notes;
```

**Migration rules:**
- One logical change per migration file
- Sequential numbering (never reuse numbers)
- Test forward AND rollback on a staging database
- Never modify a migration that has been applied to production
- Large data migrations: use batched updates (1000 rows at a time)

### 5.4 SQL Injection Prevention

**Rule:** ALWAYS use parameterized queries. NEVER string concatenation.

```go
// BAD — SQL injection vulnerability
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
rows, err := db.Query(query)

// GOOD — parameterized query
rows, err := db.Query("SELECT id, name FROM users WHERE email = $1 AND tenant_id = $2", email, tenantID)
```

```typescript
// BAD — SQL injection in any ORM/query builder
const users = await db.query(`SELECT * FROM users WHERE name = '${name}'`);

// GOOD — parameterized
const users = await db.query('SELECT * FROM users WHERE name = $1', [name]);
```

### 5.5 Transaction Rules

```go
// GOOD — transaction for multi-table writes
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return fmt.Errorf("begin tx: %w", err)
}
defer tx.Rollback() // safe: no-op after commit

// Insert order
orderID, err := insertOrder(tx, order)
if err != nil {
    return fmt.Errorf("insert order: %w", err)
}

// Insert order items
for _, item := range order.Items {
    if err := insertOrderItem(tx, orderID, item); err != nil {
        return fmt.Errorf("insert order item: %w", err)
    }
}

// Update inventory
if err := decrementStock(tx, order.Items); err != nil {
    return fmt.Errorf("decrement stock: %w", err)
}

return tx.Commit()
```

### 5.6 Multi-Tenant Query Safety

**Rule:** EVERY query on a tenant-scoped table MUST include `WHERE tenant_id = $X`. No exceptions.

```go
// BAD — cross-tenant data leak
func (r *Repo) GetOrder(orderID int) (*Order, error) {
    return r.db.QueryRow("SELECT * FROM orders WHERE id = $1", orderID)
}

// GOOD — tenant-scoped
func (r *Repo) GetOrder(tenantID, restaurantID, orderID int) (*Order, error) {
    return r.db.QueryRow(
        `SELECT id, status, total FROM orders
         WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`,
        orderID, tenantID, restaurantID,
    )
}
```

---

## 6. Performance Rules

### 6.1 Frontend Bundle Size Budget

| Metric | Budget | Action if exceeded |
|--------|--------|-------------------|
| Initial JS (gzipped) | < 150KB | Dynamic import heavy libs |
| First Contentful Paint | < 1.5s | Add `loading.tsx`, reduce blocking JS |
| Largest Contentful Paint | < 2.5s | Optimize images, preload critical assets |
| Total page JS | < 400KB | Code-split routes, check imports |
| Individual chunk | < 100KB | Split into smaller chunks |

**Check bundle size:**
```bash
# Next.js bundle analyzer
ANALYZE=true npx next build
# Or use the built-in output
npx next build  # Check .next/analyze/ output
```

### 6.2 Image Optimization

```typescript
// BAD — unoptimized image, no dimensions, blocks layout
<img src="/hero.jpg" />

// GOOD — Next.js Image with optimization
import Image from 'next/image';
<Image
  src="/hero.jpg"
  alt="Restaurant hero banner"
  width={1200}
  height={600}
  priority  // only for above-the-fold images
  sizes="(max-width: 768px) 100vw, 1200px"
/>

// For dynamic user-uploaded images:
<Image
  src={product.imageUrl}
  alt={product.name}
  width={300}
  height={300}
  loading="lazy"  // below-the-fold
  placeholder="blur"
  blurDataURL={product.thumbnailBase64}
/>
```

### 6.3 API Request Optimization

```typescript
// BAD — 16 separate API calls for permission checks on every page load
const canProducts = useHasPermission(MODULE_IDS.PRODUCTS, PERMISSION_LEVELS.READ);
const canCategories = useHasPermission(MODULE_IDS.CATEGORIES, PERMISSION_LEVELS.READ);
const canOrders = useHasPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.READ);
// ... 13 more hooks, each making a separate API call

// GOOD — single batch permission check
const permissions = useUserPermissions(); // ONE API call returns all permissions
const canProducts = permissions[MODULE_IDS.PRODUCTS] >= PERMISSION_LEVELS.READ;
const canCategories = permissions[MODULE_IDS.CATEGORIES] >= PERMISSION_LEVELS.READ;
```

```typescript
// BAD — duplicate requests on mount
useEffect(() => { fetchProducts(); }, []);
useEffect(() => { fetchProducts(); }, []); // oops, duplicate

// GOOD — React Query deduplicates automatically
const { data } = useQuery({ queryKey: ['products'], queryFn: fetchProducts });
// Same queryKey = same request, no duplicates
```

### 6.4 List Virtualization

**Rule:** Lists with 100+ items MUST use virtualization. Rendering 1000 DOM nodes kills scroll performance.

```typescript
// BAD — renders 1000 rows in the DOM
{orders.map(order => <OrderRow key={order.id} order={order} />)}

// GOOD — only renders visible rows (~20 at a time)
import { useVirtualizer } from '@tanstack/react-virtual';

function OrderList({ orders }: { orders: Order[] }) {
  const parentRef = useRef<HTMLDivElement>(null);
  const virtualizer = useVirtualizer({
    count: orders.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 60, // row height in px
    overscan: 5,
  });

  return (
    <div ref={parentRef} className="h-[600px] overflow-auto">
      <div style={{ height: virtualizer.getTotalSize() }}>
        {virtualizer.getVirtualItems().map(virtualRow => (
          <OrderRow
            key={orders[virtualRow.index].id}
            order={orders[virtualRow.index]}
            style={{
              position: 'absolute',
              top: virtualRow.start,
              height: virtualRow.size,
            }}
          />
        ))}
      </div>
    </div>
  );
}
```

### 6.5 Backend Response Time Targets

| Endpoint Type | Target | Max Acceptable |
|---------------|--------|----------------|
| Simple CRUD (GET one) | < 50ms | 200ms |
| List with pagination | < 200ms | 500ms |
| Report/aggregation | < 500ms | 2s |
| Export (PDF/CSV) | < 3s | 10s |
| Search | < 100ms | 300ms |

**If a query exceeds targets:**
1. Run `EXPLAIN ANALYZE` — check for sequential scans
2. Add missing indexes
3. Parallelize independent queries with `errgroup`
4. Add caching for data that changes rarely (< 1x/minute)

### 6.6 Caching Strategy

```go
// In-memory cache for frequently read, rarely written data
type CacheEntry struct {
    Data      interface{}
    ExpiresAt time.Time
}

// Cache invalidation: use tenant-scoped keys
cacheKey := fmt.Sprintf("products:%d:%d", tenantID, restaurantID)

// IMPORTANT: Invalidate on write
func (uc *ProductUC) Create(ctx context.Context, req *CreateProductRequest) error {
    product, err := uc.repo.Create(ctx, req)
    if err != nil {
        return err
    }
    // Invalidate cache after successful write
    uc.cache.Delete(fmt.Sprintf("products:%d:%d", req.TenantID, req.RestaurantID))
    return nil
}
```

**Cache rules:**
- Cache key MUST include `tenant_id` (prevent cross-tenant data leaks)
- Set TTL (time-to-live): 5 min for product lists, 1 min for order counts
- Invalidate on every CREATE, UPDATE, DELETE — not just on TTL expiry
- NEVER cache user-specific data in shared caches without user scoping
- Use `sync.Map` or `sync.RWMutex` for concurrent Go cache access

#### Cache Hierarchy (3 Layers)

```
Layer 1: In-Process (sync.Map) — 0ms latency, per-instance
   ↓ miss
Layer 2: Shared Cache (Redis/CDN) — 1-5ms latency, shared across instances
   ↓ miss
Layer 3: Database — 10-50ms latency, source of truth
```

```go
// Layer 1: sync.Map for hot data (menu items, settings)
var menuCache sync.Map

func (uc *MenuUC) GetMenu(ctx context.Context, tenantID, restaurantID int) ([]MenuItem, error) {
    key := fmt.Sprintf("menu:%d:%d", tenantID, restaurantID)

    // Layer 1: in-process cache
    if cached, ok := menuCache.Load(key); ok {
        entry := cached.(*CacheEntry)
        if time.Now().Before(entry.ExpiresAt) {
            return entry.Data.([]MenuItem), nil
        }
        menuCache.Delete(key)
    }

    // Layer 3: database (skip Layer 2 if no Redis)
    items, err := uc.repo.ListMenuItems(ctx, tenantID, restaurantID)
    if err != nil {
        return nil, err
    }

    menuCache.Store(key, &CacheEntry{
        Data:      items,
        ExpiresAt: time.Now().Add(5 * time.Minute),
    })
    return items, nil
}
```

```sql
-- Layer 3 optimization: Materialized views for heavy reports
CREATE MATERIALIZED VIEW mv_daily_sales AS
SELECT tenant_id, restaurant_id, DATE(created_at) AS sale_date,
       COUNT(*) AS order_count, SUM(total) AS revenue
FROM orders
WHERE status = 'completed'
GROUP BY tenant_id, restaurant_id, DATE(created_at);

-- Refresh on schedule (not on every write)
REFRESH MATERIALIZED VIEW CONCURRENTLY mv_daily_sales;

-- Index the materialized view
CREATE UNIQUE INDEX idx_mv_daily_sales
  ON mv_daily_sales(tenant_id, restaurant_id, sale_date);
```

**When to use each layer:**
| Data Type | Layer 1 (sync.Map) | Layer 2 (CDN/Redis) | Layer 3 (DB) |
|-----------|-------------------|--------------------|----|
| Menu items | 5 min TTL | CDN edge cache | Source |
| User permissions | 5 min TTL | Skip | Source |
| Order counts | 30 sec TTL | Skip | Source |
| Sales reports | Skip | Skip | Materialized view |
| Static assets | Skip | CDN with immutable | Origin |

---

## 7. Security Rules

### 7.1 Authentication

**Rule:** NEVER hardcode secrets. Use environment variables with fallback checks.

```go
// BAD — hardcoded default that will end up in production
jwtSecret := getEnvOrDefault("JWT_SECRET", "your-secret-key-change-this")

// GOOD — fail fast if secret is missing
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
    log.Fatal("JWT_SECRET environment variable is required")
}
if len(jwtSecret) < 32 {
    log.Fatal("JWT_SECRET must be at least 32 characters")
}
```

```typescript
// BAD — auth token in localStorage (XSS can steal it)
localStorage.setItem('token', response.data.token);

// GOOD — httpOnly cookie (JavaScript can't access it)
// Set from backend:
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    token,
    HttpOnly: true,   // JS can't read
    Secure:   true,   // HTTPS only
    SameSite: http.SameSiteStrictMode,
    Path:     "/",
    MaxAge:   86400,  // 24 hours
})
```

### 7.2 Input Sanitization

```typescript
// BAD — XSS vulnerability
<div dangerouslySetInnerHTML={{ __html: userContent }} />

// GOOD — sanitize first
import DOMPurify from 'dompurify';
<div dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(userContent) }} />

// EVEN BETTER — avoid dangerouslySetInnerHTML when possible
// Use a markdown renderer or structured content instead
```

```go
// BAD — user input directly in log (log injection)
log.Printf("User login: %s", username)

// GOOD — sanitize log output
log.Printf("User login: %s", strings.ReplaceAll(username, "\n", ""))
```

### 7.3 HTTP Security Headers

**Rule:** Every production deployment must set these headers:

```go
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Prevent clickjacking
        w.Header().Set("X-Frame-Options", "DENY")
        // Prevent MIME sniffing
        w.Header().Set("X-Content-Type-Options", "nosniff")
        // XSS protection (legacy browsers)
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        // HTTPS enforcement
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        // Referrer policy
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        // Permissions policy
        w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
        // CSP — customize per project
        w.Header().Set("Content-Security-Policy",
            "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:;")

        next.ServeHTTP(w, r)
    })
}
```

### 7.4 Rate Limiting

```go
// Protect sensitive endpoints
// Login:        5 attempts per minute per IP
// API:          100 requests per minute per user
// File upload:  10 per minute per user
// Password reset: 3 per hour per email

// BAD — no rate limiting on login
mux.HandleFunc("/api/v1/auth/login", authHandler.Login)

// GOOD — rate limited
mux.Handle("/api/v1/auth/login",
    rateLimiter.Limit("login", 5, time.Minute)(authHandler.Login))
```

### 7.5 CORS Configuration

```go
// BAD — allows everything
w.Header().Set("Access-Control-Allow-Origin", "*")

// GOOD — whitelist specific origins
allowedOrigins := map[string]bool{
    "https://dashboard.yourapp.com": true,
    "https://www.yourapp.com":       true,
}
// In development, also allow localhost
if isDevelopment {
    allowedOrigins["http://localhost:3001"] = true
    allowedOrigins["http://localhost:3002"] = true
}

origin := r.Header.Get("Origin")
if allowedOrigins[origin] {
    w.Header().Set("Access-Control-Allow-Origin", origin)
    w.Header().Set("Vary", "Origin")
}
```

### 7.6 File Upload Security

```go
// Rule: Validate file type, size, and name. NEVER trust the extension.

func handleUpload(w http.ResponseWriter, r *http.Request) {
    // 1. Limit request body size (10MB max)
    r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

    file, header, err := r.FormFile("image")
    if err != nil {
        respondError(w, http.StatusBadRequest, "file too large or missing")
        return
    }
    defer file.Close()

    // 2. Validate MIME type by reading magic bytes (NOT extension)
    buf := make([]byte, 512)
    file.Read(buf)
    contentType := http.DetectContentType(buf)
    file.Seek(0, 0) // reset reader

    allowedTypes := map[string]bool{
        "image/jpeg": true, "image/png": true, "image/webp": true,
    }
    if !allowedTypes[contentType] {
        respondError(w, http.StatusBadRequest, "unsupported file type")
        return
    }

    // 3. Generate safe filename (NEVER use user-provided filename)
    ext := ".jpg"
    if contentType == "image/png" { ext = ".png" }
    safeName := fmt.Sprintf("%d_%s%s", tenantID, uuid.New().String(), ext)

    // 4. Save to controlled path
    savePath := filepath.Join(uploadDir, safeName)
    // ...
}
```

### 7.7 Dependency Security

- Run `npm audit` / `go mod verify` weekly
- Pin exact dependency versions in production (`package-lock.json` committed)
- Never install packages from unknown sources
- Review changelogs before major version upgrades
- Set `typescript.ignoreBuildErrors = false` — hiding errors hides vulnerabilities

---

## 8. API Design Standards

### 8.1 RESTful URL Patterns

```
GET    /api/v1/products          → List (paginated)
GET    /api/v1/products/:id      → Get one
POST   /api/v1/products          → Create
PUT    /api/v1/products/:id      → Full update
PATCH  /api/v1/products/:id      → Partial update
DELETE /api/v1/products/:id      → Delete

Nested resources:
GET    /api/v1/orders/:id/items  → List items of an order
POST   /api/v1/orders/:id/items  → Add item to order

Actions (non-CRUD operations):
POST   /api/v1/orders/:id/confirm    → Confirm order
POST   /api/v1/orders/:id/cancel     → Cancel order
POST   /api/v1/auth/login            → Login
POST   /api/v1/auth/refresh          → Refresh token
```

**Rules:**
- Use plural nouns: `/products` not `/product`
- Version the API: `/api/v1/...`
- Use kebab-case in URLs: `/order-items` not `/orderItems`
- No verbs in URLs: `/orders` not `/getOrders`
- Nest max 2 levels: `/orders/:id/items` (never `/tenants/:id/restaurants/:id/orders/:id/items`)

### 8.2 Response Format

```json
// Success (single resource)
{
  "data": { "id": 1, "name": "Burger", "price": 50.00 },
  "message": "Product created successfully"
}

// Success (list with pagination)
{
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "limit": 25,
    "total": 150,
    "total_pages": 6
  }
}

// Error
{
  "error": "validation_error",
  "message": "Price must be a positive number",
  "details": [
    { "field": "price", "message": "must be greater than 0" }
  ]
}
```

### 8.3 HTTP Status Codes

| Code | When to use |
|------|-------------|
| 200 | Successful GET, PUT, PATCH, DELETE |
| 201 | Successful POST (resource created) |
| 204 | Successful DELETE (no body) |
| 400 | Invalid input, validation error |
| 401 | Missing or invalid authentication |
| 403 | Authenticated but not authorized |
| 404 | Resource not found |
| 409 | Conflict (duplicate, state violation) |
| 422 | Semantically invalid (valid JSON, but wrong data) |
| 429 | Rate limited |
| 500 | Server error (log it, alert on it) |

**Rule:** NEVER return 200 for errors. NEVER return 500 for client mistakes.

### 8.4 Pagination Parameters

```
Standard query parameters:
?page=1&limit=25          → Page-based
?cursor=abc123&limit=25   → Cursor-based (for real-time feeds)

Filtering:
?status=active&category_id=5
?search=burger
?created_after=2026-01-01&created_before=2026-02-01

Sorting:
?sort=created_at&order=desc
?sort=-price              → prefix minus = descending

Combined:
GET /api/v1/products?page=1&limit=25&category_id=5&sort=-price&search=burger
```

### 8.5 API Versioning

**Rule:** Never break existing API contracts. Use versioning for breaking changes.

```
v1: Original API
v2: Breaking change (new request/response format)

// Backend: version in URL prefix
mux.HandleFunc("/api/v1/products", v1ProductHandler.List)
mux.HandleFunc("/api/v2/products", v2ProductHandler.List)

// Deprecation: add header
w.Header().Set("Deprecation", "true")
w.Header().Set("Sunset", "2026-06-01")
w.Header().Set("Link", `</api/v2/products>; rel="successor-version"`)
```

---

## 9. Testing Standards

### 9.1 What to Test

| Layer | What to test | Example |
|-------|-------------|---------|
| **Unit** | Pure functions, utilities, validators | `calculateTotal()`, `formatCurrency()`, `validateEmail()` |
| **Integration** | API handlers + database | POST /api/v1/orders creates order and items |
| **Component** | React components with user interaction | "Submit button is disabled when form is invalid" |
| **E2E** | Critical user journeys | Login → Create order → Process payment |

### 9.2 Test Naming

```typescript
// BAD — describes implementation
test('calls setLoading and fetchProducts', () => { ... });

// GOOD — describes behavior
test('shows loading spinner while products are being fetched', () => { ... });
test('displays error message when product fetch fails', () => { ... });
test('renders empty state when no products exist', () => { ... });
```

```go
// BAD
func TestCreate(t *testing.T) { ... }

// GOOD — describes scenario and expectation
func TestCreateOrder_WithValidInput_ReturnsCreatedOrder(t *testing.T) { ... }
func TestCreateOrder_WithMissingName_Returns400(t *testing.T) { ... }
func TestCreateOrder_WhenUnauthorized_Returns401(t *testing.T) { ... }
```

### 9.3 Test Structure (Arrange-Act-Assert)

```typescript
test('filters orders by status', () => {
  // Arrange — set up test data
  const orders = [
    { id: 1, status: 'pending' },
    { id: 2, status: 'delivered' },
    { id: 3, status: 'pending' },
  ];

  // Act — perform the action
  const result = filterOrdersByStatus(orders, 'pending');

  // Assert — verify the result
  expect(result).toHaveLength(2);
  expect(result.every(o => o.status === 'pending')).toBe(true);
});
```

```go
func TestCalculateOrderTotal_WithDiscountAndTax(t *testing.T) {
    // Arrange
    items := []domain.OrderItem{
        {Price: 100, Quantity: 2},
        {Price: 50, Quantity: 1},
    }
    discount := 10.0 // 10%
    taxRate := 14.0   // 14% Egyptian VAT

    // Act
    total := calculateOrderTotal(items, discount, taxRate)

    // Assert
    // Subtotal: 250, After 10% discount: 225, After 14% tax: 256.50
    if total != 256.50 {
        t.Errorf("expected 256.50, got %.2f", total)
    }
}
```

### 9.4 Test Database Isolation

```go
// BAD — tests share state, order-dependent
func TestListProducts(t *testing.T) {
    // Depends on TestCreateProduct having run first!
    products, _ := repo.List(1)
    assert.Equal(t, 5, len(products))
}

// GOOD — each test sets up and cleans up its own data
func TestListProducts(t *testing.T) {
    // Setup: create test data in a transaction
    tx, _ := db.Begin()
    defer tx.Rollback() // auto-cleanup

    testRepo := NewProductRepo(tx)
    testRepo.Create(&Product{Name: "Test Burger", TenantID: 999})
    testRepo.Create(&Product{Name: "Test Fries", TenantID: 999})

    // Test
    products, err := testRepo.List(999, 1, 25)
    assert.NoError(t, err)
    assert.Equal(t, 2, len(products))
}
```

### 9.5 What NOT to Test

- Third-party library internals (React renders, Zustand stores)
- Simple getters/setters with no logic
- CSS/styling (use visual regression tools instead)
- Implementation details (state changes, function calls)

**Test behavior, not implementation.** If you refactor without changing behavior, no tests should break.

### 9.6 Coverage Targets

| Area | Minimum Coverage |
|------|-----------------|
| Business logic (usecases) | 90% |
| API handlers | 80% |
| Utility functions | 95% |
| React components | 70% (critical paths) |
| Overall | 80% |

---

## 10. Git & Code Review Standards

### 10.1 Branch Naming

```
feature/add-order-export-pdf
bugfix/fix-login-redirect-loop
hotfix/patch-sql-injection-orders
refactor/optimize-sidebar-permissions
chore/update-dependencies-march-2026
docs/add-api-documentation
```

### 10.2 Commit Messages

```
Format: <type>(<scope>): <description>

Types: feat, fix, refactor, perf, docs, test, chore, style, ci

Examples:
feat(orders): add PDF export for order history
fix(auth): prevent redirect loop on expired JWT
perf(reports): dynamic import recharts (saves 400KB)
refactor(sidebar): batch permission checks into single API call
test(payments): add integration tests for Paymob webhook
chore(deps): update next.js to 15.1.0

RULES:
- Subject line ≤ 72 characters
- Use imperative mood: "add" not "added" or "adds"
- No period at the end of subject
- Body (optional): explain WHY, not WHAT
```

### 10.3 Pull Request Standards

```markdown
## What does this PR do?
- Adds PDF export to order history page
- Dynamically imports jspdf (350KB saved from initial bundle)

## How to test
1. Navigate to /dashboard/orders
2. Click "Export PDF" button
3. Verify PDF contains all visible orders
4. Check Network tab — jspdf only loads on button click

## Checklist
- [ ] Self-review completed (used PR checklist from Section 1)
- [ ] Tests added/updated
- [ ] No TypeScript errors (`tsc --noEmit`)
- [ ] No ESLint warnings
- [ ] Tested on Chrome + Firefox
- [ ] Mobile responsive checked
```

### 10.4 Code Review Rules (for Reviewers)

**Focus on:**
1. **Correctness** — Does it do what it claims?
2. **Security** — Any injection, XSS, auth bypass, data leak?
3. **Performance** — Any N+1 queries, missing pagination, heavy imports?
4. **Readability** — Can a new developer understand this in 5 minutes?

**Don't focus on:**
- Bike-shedding (naming preferences, bracket style)
- Anything the linter should catch

**Review checklist:**
- [ ] All new queries include `tenant_id` WHERE clause
- [ ] All new list endpoints have pagination
- [ ] All new API fields are validated
- [ ] No `console.log` or debug code left behind
- [ ] Error states are handled in the UI
- [ ] Loading states are shown
- [ ] New environment variables documented

### 10.5 Protected Branches

```
main / master:
- No direct pushes
- Requires PR with at least 1 approval
- CI must pass (lint + test + build)
- No force push EVER

staging:
- Auto-deploy to staging environment
- Can merge from feature branches
- Run full test suite before merge to main
```

---

## 11. Accessibility (a11y) Rules

**Target:** WCAG 2.2 Level AA compliance for all user-facing pages.

### 11.1 Color & Contrast

```css
/* Minimum contrast ratios (WCAG 2.2 AA) */
/* Normal text (< 18px): 4.5:1 contrast ratio */
/* Large text (≥ 18px bold or ≥ 24px): 3:1 contrast ratio */
/* UI components & graphical objects: 3:1 contrast ratio */

/* BAD — light gray on white (ratio ~2:1) */
.muted-text { color: #bbbbbb; }

/* GOOD — meets 4.5:1 on white background */
.muted-text { color: #767676; }
```

**Color rules:**
- NEVER convey information by color alone — always pair with icons, text, or patterns
- Test all color themes (light + dark) with a contrast checker
- Error states: use red (#DC2626) + error icon + descriptive text (not just red border)

### 11.2 Semantic HTML & ARIA

```tsx
// BAD — div soup with no semantics
<div onClick={handleClick} className="btn">Submit</div>

// GOOD — native button with accessible label
<button onClick={handleClick} type="submit">Submit Order</button>

// BAD — icon-only button with no label
<button onClick={toggleMenu}><MenuIcon /></button>

// GOOD — icon button with aria-label
<button onClick={toggleMenu} aria-label="Toggle navigation menu">
  <MenuIcon aria-hidden="true" />
</button>

// BAD — custom dropdown with no ARIA
<div className="dropdown" onClick={toggle}>
  <div className="selected">{value}</div>
  <div className="options">{options.map(...)}</div>
</div>

// GOOD — use native select OR proper ARIA roles
<select value={value} onChange={handleChange} aria-label="Select category">
  {options.map(o => <option key={o.id} value={o.id}>{o.name}</option>)}
</select>
```

**ARIA rules:**
- Prefer native HTML elements (`<button>`, `<select>`, `<input>`) over custom ARIA
- Every `<img>` needs `alt` text (decorative images: `alt=""` + `aria-hidden="true"`)
- Form inputs need visible `<label>` elements (not just placeholder text)
- Dynamic content updates: use `aria-live="polite"` for non-urgent, `aria-live="assertive"` for errors

### 11.3 Keyboard Navigation

```tsx
// Every interactive element must be reachable by Tab and activatable by Enter/Space

// BAD — clickable div is not keyboard accessible
<div onClick={handleSelect} className="card">
  {product.name}
</div>

// GOOD — focusable and keyboard-accessible
<div
  onClick={handleSelect}
  onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') handleSelect() }}
  role="button"
  tabIndex={0}
  className="card"
>
  {product.name}
</div>

// BEST — just use a button
<button onClick={handleSelect} className="card">
  {product.name}
</button>
```

**Keyboard rules:**
- Tab order must follow visual order (no `tabIndex` > 0)
- Escape key closes modals, dropdowns, and popovers
- Focus must be trapped inside open modals (return focus to trigger on close)
- Skip-to-content link on every page: `<a href="#main-content" className="sr-only focus:not-sr-only">`
- Never remove `:focus-visible` outlines — restyle them, don't hide them

### 11.4 Screen Reader Testing

**Required testing on critical flows:**
- Login / registration
- Product browsing and ordering
- Checkout / payment
- Dashboard navigation

**Screen reader checklist:**
- [ ] Page has a single `<h1>` and headings follow logical order (h1 → h2 → h3)
- [ ] Tables have `<th>` headers with `scope="col"` or `scope="row"`
- [ ] Form errors are announced (linked with `aria-describedby` or `aria-errormessage`)
- [ ] Loading states announced: `<div aria-live="polite">{loading ? 'Loading...' : 'Content ready'}</div>`
- [ ] Route changes announced in SPAs (Next.js handles this, but verify)

### 11.5 RTL Accessibility

```tsx
// For Arabic RTL support — ensure logical properties
// BAD — physical properties break in RTL
.sidebar { margin-left: 16px; padding-right: 24px; }

// GOOD — logical properties adapt to text direction
.sidebar { margin-inline-start: 16px; padding-inline-end: 24px; }

// Ensure aria-label text matches the active locale
const { locale } = useLocale()
<button aria-label={locale === 'ar' ? 'فتح القائمة' : 'Open menu'}>
```

---

## 12. Mobile-First Standards

### 12.1 Touch Targets

```css
/* Minimum touch target: 48x48px (Google), 44x44px (Apple) — we use 48px */

/* BAD — tiny 24px icon button */
.icon-btn { width: 24px; height: 24px; }

/* GOOD — visual size can be smaller, but tap area is 48px */
.icon-btn {
  width: 24px;
  height: 24px;
  padding: 12px; /* extends tap area to 48px total */
  /* OR use min-width/min-height */
  min-width: 48px;
  min-height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
}
```

**Touch target rules:**
- All buttons, links, and interactive elements: 48×48px minimum tap area
- Space between adjacent targets: at least 8px gap (prevent mis-taps)
- Bottom-of-screen actions (mobile nav, FAB): account for `safe-area-inset-bottom`

### 12.2 Safe Area & Viewport

```css
/* iOS notch and home indicator safe areas */
.bottom-nav {
  padding-bottom: env(safe-area-inset-bottom);
}

.page-container {
  padding-left: env(safe-area-inset-left);
  padding-right: env(safe-area-inset-right);
}

/* Full viewport height that accounts for mobile browser chrome */
.full-page {
  min-height: 100dvh; /* dvh = dynamic viewport height */
}
```

```html
<!-- Required viewport meta tag -->
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover" />
```

### 12.3 Responsive Breakpoints

```css
/* Mobile-first: base styles = mobile, add complexity with min-width */

/* BAD — desktop-first (overriding down) */
.grid { display: grid; grid-template-columns: repeat(4, 1fr); }
@media (max-width: 768px) { .grid { grid-template-columns: 1fr; } }

/* GOOD — mobile-first (building up) */
.grid { display: grid; grid-template-columns: 1fr; gap: 16px; }
@media (min-width: 640px)  { .grid { grid-template-columns: repeat(2, 1fr); } }
@media (min-width: 1024px) { .grid { grid-template-columns: repeat(3, 1fr); } }
@media (min-width: 1280px) { .grid { grid-template-columns: repeat(4, 1fr); } }
```

**Tailwind breakpoints (our standard):**
| Breakpoint | Min Width | Target |
|-----------|----------|--------|
| (default) | 0px | Mobile phones |
| `sm:` | 640px | Large phones / small tablets |
| `md:` | 768px | Tablets |
| `lg:` | 1024px | Laptops |
| `xl:` | 1280px | Desktops |

### 12.4 RTL Mobile Considerations

```tsx
// Swipe gestures must respect RTL direction
const swipeDirection = locale === 'ar' ? 'right-to-left' : 'left-to-right'

// Horizontal scroll containers: set dir attribute
<div dir={locale === 'ar' ? 'rtl' : 'ltr'} className="overflow-x-auto">
  {/* Cards scroll in correct direction */}
</div>
```

### 12.5 Performance on Mobile

```tsx
// Test with Chrome DevTools Network Throttling:
// - "Slow 3G" preset (400ms RTT, 500kbps)
// - Target: interactive in < 5s on 3G

// Lazy load images that are below the fold
<Image
  src={product.imageUrl}
  alt={product.name}
  loading="lazy"
  width={300}
  height={200}
/>

// Eager load above-the-fold hero images
<Image
  src="/hero-banner.webp"
  alt=""
  loading="eager"
  priority
  width={1200}
  height={400}
/>
```

**Mobile performance rules:**
- Images: use `next/image` with WebP format and appropriate `sizes` attribute
- Fonts: subset to used characters, `font-display: swap`, preload critical fonts
- JavaScript: code-split by route, lazy-load below-fold components
- Test on real devices (not just DevTools) — at least one Android + one iOS device

---

## 13. Deployment & Production Checklist

### 13.1 Pre-Deployment Gate

Every deployment MUST pass these checks before going to production:

```bash
# 1. Code quality
pnpm lint          # Zero errors (warnings acceptable but track them)
pnpm type-check    # Zero TypeScript errors
pnpm test          # All tests pass
pnpm build         # Build succeeds with no errors

# 2. Backend
go vet ./...       # Zero issues
go test ./...      # All tests pass
go build -o api ./cmd/api  # Binary builds successfully

# 3. Database
# Migration tested on staging (forward + rollback)
# No breaking schema changes without migration path
```

### 13.2 Staging Validation

```markdown
Before promoting staging → production:

- [ ] All API endpoints respond correctly (run smoke test suite)
- [ ] Login flow works (JWT issued, refresh works, logout clears tokens)
- [ ] Multi-tenant isolation verified (Tenant A cannot see Tenant B data)
- [ ] RBAC permissions enforced (non-admin cannot access admin routes)
- [ ] Payment flow tested end-to-end (Paymob sandbox or mock mode)
- [ ] WebSocket connections work (KDS real-time updates)
- [ ] File uploads work (product images, theme assets)
- [ ] Email/notification delivery confirmed
- [ ] RTL (Arabic) layout renders correctly
- [ ] Mobile viewport tested (no horizontal overflow, touch targets work)
```

### 13.3 Rollback Plan

```markdown
Every deployment must have a documented rollback:

1. **Frontend rollback:** Revert to previous deployment
   - Vercel/Netlify: one-click rollback to previous deployment
   - Docker: `docker compose up -d --pull always` with previous image tag

2. **Backend rollback:** Revert binary + rollback migration if needed
   - Keep previous binary: `cp api api.backup` before deploying
   - Rollback: `cp api.backup api && systemctl restart api`

3. **Database rollback:** ONLY if migration was applied
   - Run rollback script: `psql -f rollback_migration_160.sql`
   - NEVER rollback a migration that has been live for > 1 hour
     (data may have been written to new columns/tables)

4. **Rollback decision:** If critical bug found within 15 minutes → rollback
   If found after 15 minutes → hotfix forward (unless data corruption)
```

### 13.4 Post-Deployment Smoke Tests

```bash
# Automated smoke tests to run after every deployment

# Health check
curl -f https://api.example.com/health || echo "FAIL: health check"

# Auth flow
TOKEN=$(curl -s -X POST https://api.example.com/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"smoke@test.com","password":"..."}' | jq -r '.token')
[ -n "$TOKEN" ] || echo "FAIL: login"

# Tenant-scoped API
curl -sf -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: 1" \
  https://api.example.com/products | jq '.data | length' || echo "FAIL: products"

# Frontend loads
curl -sf https://dashboard.example.com/en/dashboard > /dev/null || echo "FAIL: dashboard"
```

### 13.5 Environment Variables Audit

```markdown
Before every deployment, verify ALL required env vars are set:

Backend (required):
- [ ] DATABASE_URL — PostgreSQL connection string
- [ ] JWT_SECRET — Min 32 characters, unique per environment
- [ ] SUPER_ADMIN_JWT_SECRET — Different from JWT_SECRET
- [ ] PAYMOB_API_KEY — Payment gateway (or MOCK_PAYMENT=true)
- [ ] CORS_ORIGINS — Comma-separated allowed origins

Frontend (required):
- [ ] NEXT_PUBLIC_API_URL — Backend API URL (no trailing slash)
- [ ] NEXT_PUBLIC_WS_URL — WebSocket URL for KDS

Environment-specific:
- [ ] NODE_ENV=production (frontend)
- [ ] GIN_MODE=release (backend, if using Gin)
- [ ] DEBUG_LOGS=false (both — no debug output in production)
```

---

## Appendix A: Common Anti-Patterns & Fixes

| Anti-Pattern | Why It's Bad | Fix |
|-------------|-------------|-----|
| `SELECT * FROM table` | Fetches unnecessary data, breaks when columns added | List explicit columns |
| `console.log` in production | Leaks internal state, fills console | Use `debugLog()` wrapper with env check |
| `any` type in TypeScript | Defeats type checking, hides bugs | Use proper types or `unknown` + type guard |
| `eslint-disable` without comment | Hides problems | Add comment explaining WHY |
| `catch(e) {}` empty catch | Swallows errors silently | Log, report, or surface to user |
| Polling without visibility check | Wastes bandwidth on hidden tabs | Check `document.visibilityState` |
| Barrel exports (`index.ts`) | Prevents tree-shaking, bloats bundle | Import directly from source file |
| String concatenation in SQL | SQL injection vulnerability | Use parameterized queries |
| `localStorage` for auth tokens | XSS can steal tokens | Use httpOnly cookies |
| No pagination on list endpoints | Memory + CPU exhaustion at scale | Always paginate, max 100 per page |
| Sequential independent queries | Slow response times | Use `errgroup` / `Promise.all` |
| `ignoreBuildErrors: true` | Hides real bugs and type errors | Fix the errors, set to `false` |
| Hardcoded secrets in code | Secrets in git history forever | Use env variables, fail if missing |
| Missing `loading.tsx` | Blank screen during navigation | Add skeleton to every route segment |
| `dangerouslySetInnerHTML` without sanitizer | XSS vulnerability | Always use DOMPurify first |

---

## Appendix B: Recommended next.config.js Settings

```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  // NEVER set these to true in production
  typescript: { ignoreBuildErrors: false },
  eslint: { ignoreDuringBuilds: false },

  // Optimize heavy packages
  experimental: {
    optimizePackageImports: [
      'recharts',
      '@tanstack/react-table',
      '@tanstack/react-query',
      'date-fns',
      'lucide-react',
      '@radix-ui/react-icons',
    ],
  },

  // Image optimization
  images: {
    formats: ['image/avif', 'image/webp'],
    deviceSizes: [640, 750, 828, 1080, 1200, 1920],
    minimumCacheTTL: 60 * 60 * 24, // 24 hours
  },

  // Security headers
  headers: async () => [{
    source: '/(.*)',
    headers: [
      { key: 'X-Frame-Options', value: 'DENY' },
      { key: 'X-Content-Type-Options', value: 'nosniff' },
      { key: 'Referrer-Policy', value: 'strict-origin-when-cross-origin' },
    ],
  }],
};
```

---

## Appendix C: Go Project Linting Config

```bash
# Must pass before every commit:
go vet ./...                    # Static analysis
go build ./...                  # Compilation check
golangci-lint run               # Comprehensive linting

# Recommended golangci-lint config (.golangci.yml):
# linters:
#   enable:
#     - errcheck       # Check error returns
#     - gosec          # Security issues
#     - govet          # Suspicious constructs
#     - ineffassign    # Unused assignments
#     - staticcheck    # Advanced static analysis
#     - unused         # Unused code
#     - gosimple       # Simplification suggestions
#     - bodyclose      # Unclosed HTTP response bodies
#     - sqlclosecheck  # Unclosed SQL rows
```

---

## Appendix D: Performance Monitoring Checklist

**Before every release, check:**

```
Frontend:
□ Run Lighthouse: Performance score > 90
□ Check bundle size: npx next build → no chunk > 100KB gzipped
□ Verify no layout shift: CLS < 0.1
□ Test on 3G throttle: pages load < 5s
□ Check memory: no leaks after 10 min of navigation (Chrome DevTools → Memory)

Backend:
□ Run load test: k6 with 100 concurrent users
□ Check p95 response times < 500ms for list endpoints
□ Verify connection pool: no "too many connections" errors
□ Check slow query log: no queries > 1s
□ Monitor memory: no growth over time (goroutine leaks)

Database:
□ Run EXPLAIN ANALYZE on top 10 queries
□ Check for missing indexes: pg_stat_user_tables → seq_scan vs idx_scan
□ Verify no table locks during peak hours
□ Check bloat: pg_stat_user_tables → dead tuples
```

---

*This document is a living standard. Update it when you discover new anti-patterns, find better practices, or learn from production incidents. Every rule exists because someone made a mistake — don't repeat it.*
