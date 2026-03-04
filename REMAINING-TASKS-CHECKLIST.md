# Remaining Tasks Checklist

> Generated: 2026-03-04
> Status: All completed items from previous pass are excluded.
> Only **not-done** and **partially-done** tasks remain.

---

## RBAC & Permissions (Remaining)

- [x] **R1: Migrate remaining 30 pages from `useHasPermission` to `useUserPermissions`**
  - 30 non-accounting pages still use individual `useHasPermission` hooks (each makes a separate API call)
  - Priority pages: orders, products, HR (employees, payroll, leaves, attendance), tables, POS, settings
  - Switch to single `useUserPermissions()` + `checkPermission()` pattern (like accounting pages)
  - Files: `orders/page.tsx`, `products/page.tsx`, `hr/employees/page.tsx`, `hr/payroll/page.tsx`, `hr/leaves/page.tsx`, `hr/attendance/page.tsx`, `tables/page.tsx`, `pos/page.tsx`, `settings/users/page.tsx`, `settings/roles/page.tsx`, `kds/settings/page.tsx`, `notifications/page.tsx`, `reports/overview/page.tsx`, `reports/sales/page.tsx`, `reports/profitability/page.tsx`, `reports/inventory/page.tsx`, `reports/hr/page.tsx`, and all sub-pages

- [x] **R2: Add `loading` state guard to all permission-gated pages**
  - Every page that shows "Access Denied" should first check `if (loading) return <Skeleton />`
  - Prevents flash of "Access Denied" before permissions arrive
  - Apply to all 30+ pages listed above

---

## Performance (Remaining)

- [x] **P1: Add `loading.tsx` files for all dashboard routes**
  - Only 3 routes have `loading.tsx` (sales, profitability, inventory reports)
  - Need `loading.tsx` for: orders, products, HR (4 pages), tables (4 pages), POS, settings (3 pages), accounting (9 pages), notifications, reports/overview, reports/hr, reports/financial
  - ~25 `loading.tsx` files needed (skeleton shimmer layouts matching page structure)

- [ ] **P2: True database-level pagination for list handlers** (PARTIAL: implemented for training/insurance/contract/resignation/supplier/customer_acct/bank_account + hr_employee/hr_attendance/hr_leaves/hr_salaries list endpoints using COUNT + LIMIT/OFFSET; remaining endpoints still pending)
  - Current pagination is in-memory (loads all rows, then slices)
  - Switch to SQL `LIMIT $1 OFFSET $2` in repository layer for:
    - `training_handler.go` / `training_repo.go`
    - `insurance_handler.go` / `insurance_repo.go`
    - `contract_handler.go` / `contract_repo.go`
    - `resignation_handler.go` / `resignation_repo.go`
    - `supplier_handler.go` / `supplier_repo.go`
    - `customer_acct_handler.go` / `customer_acct_repo.go`
    - `bank_account_handler.go` / `bank_account_repo.go`
  - Also add pagination to handlers that don't have it yet (check: `order_handler`, `product_handler`, `employee_handler`, etc.)

- [x] **P3: HR report batching with parallel query execution**
  - Current HR reports run queries sequentially
  - Use `errgroup` or goroutines with `sync.WaitGroup` to run independent queries in parallel
  - Files: `hr_report_handler.go`, `report_handler.go`

- [x] **P4: Context timeout propagation across all handlers**
  - Only 3 handler files use `context.WithTimeout` (report_handler, kds_handler, kds_hub_test)
  - 44 remaining handler files need timeout propagation
  - Pattern: `ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second); defer cancel()`
  - Pass `ctx` to all repo calls instead of `r.Context()`
  - Implemented globally via `TimeoutMiddleware(30s)` in `main.go`, with websocket-upgrade bypass.

- [ ] **P5: Image optimization audit**
  - Verify all `<img>` tags use `next/image` with proper `sizes`, `width`, `height`
  - Add WebP format support where missing
  - Implement blur placeholders for product images (`blurDataURL`)

- [ ] **P6: Bundle size audit & code splitting**
  - Analyze bundle with `ANALYZE=true pnpm build`
  - Identify and lazy-load heavy components (charts, editors, PDF generators)
  - Target: no single chunk > 100KB gzipped

- [ ] **P7: Preload critical data on dashboard shell**
  - Preload user permissions, tenant settings, and notification count in parallel on dashboard mount
  - Use `Promise.all` or React `use()` with Suspense for parallel fetching

- [ ] **P8: WebSocket reconnection strategy**
  - KDS WebSocket should auto-reconnect with exponential backoff
  - Show "Reconnecting..." banner when connection drops
  - Queue updates during disconnection and replay on reconnect

---

## Security (Remaining)

- [x] **S1: Full CSRF protection (token + cookie flow)**
  - No CSRF implementation exists in the backend
  - Add: Generate CSRF token on login → set as httpOnly cookie → frontend sends in header → backend validates
  - Middleware: `CSRFMiddleware` on all state-changing routes (POST, PUT, DELETE)
  - Frontend: add CSRF token header to all non-GET API requests

- [ ] **S2: Migrate JWT to httpOnly cookies**
  - Currently: JWT stored in localStorage (vulnerable to XSS)
  - Target: JWT in httpOnly, Secure, SameSite=Strict cookie
  - Backend: set cookie on login response, read from cookie in auth middleware
  - Frontend: remove localStorage token storage, let browser handle cookies
  - This is a large cross-cutting refactor — touches auth middleware, login handler, all frontend API calls
  - **Prerequisite:** S1 (CSRF) must be done first (cookies need CSRF protection)

- [ ] **S3: Audit all IP extraction usage**
  - `r.RemoteAddr` may return proxy IP, not client IP
  - Check all usages beyond super-admin path
  - Use `X-Forwarded-For` header with trusted proxy validation
  - Files to audit: all handlers that log or validate IP addresses

- [ ] **S4: CSP nonce for production**
  - Current CSP uses `'unsafe-inline' 'unsafe-eval'` for script-src
  - For production: generate per-request nonce, inject into CSP header + script tags
  - Requires Next.js custom server or middleware for nonce injection

- [x] **S5: Rate limiting on authentication endpoints**
  - Add rate limiter middleware: 5 attempts per minute per IP on `/auth/login`
  - Add rate limiter: 3 attempts per minute on `/auth/forgot-password`
  - Use in-memory token bucket (or Redis for multi-instance)

- [ ] **S6: Input validation hardening** (PARTIAL: added global JSON max body middleware (1MB) + field length validation on core HR create/update handlers; full decode/input-length audit still pending)
  - Audit all `json.Decode` calls for max body size limits
  - Add `http.MaxBytesReader` wrapper: `r.Body = http.MaxBytesReader(w, r.Body, 1<<20)` (1MB)
  - Validate string lengths on all user inputs (name, email, description fields)

---

## Frontend Quality (New)

- [x] **F1: Error boundaries for all dashboard sections**
  - Add React Error Boundaries at route segment level
  - Create `error.tsx` files for each dashboard section
  - Show user-friendly error message with "Try Again" button

- [ ] **F2: Offline/network error handling**
  - Show toast or banner when API requests fail due to network issues
  - Distinguish between 401 (redirect to login), 403 (show access denied), 500 (show error), network error (show offline banner)

- [ ] **F3: Form validation consistency**
  - Audit all forms for client-side validation
  - Ensure error messages are accessible (linked with `aria-describedby`)
  - Use consistent validation pattern (validate on blur, show on submit)

---

## Backend Quality (New)

- [ ] **B1: Structured logging**
  - Replace `log.Printf` / `fmt.Printf` with structured logger (zerolog or slog)
  - Include tenant_id, user_id, request_id in all log entries
  - Separate log levels: INFO for normal ops, WARN for degraded, ERROR for failures

- [x] **B2: Health check endpoint**
  - Add `/health` and `/ready` endpoints
  - `/health` — basic liveness (always 200)
  - `/ready` — checks DB connection, returns 503 if DB is down

- [ ] **B3: Graceful shutdown**
  - Handle SIGTERM/SIGINT signals
  - Drain active HTTP connections (30s timeout)
  - Close DB pool, WebSocket hub, background workers
  - Log shutdown progress

---

## Testing (Remaining)

- [ ] **T1: Integration tests for RBAC permission checks**
  - Test: user with READ-only cannot POST/PUT/DELETE
  - Test: sidebar shows correct items per role
  - Test: permission cache invalidation after role change

- [ ] **T2: E2E smoke tests for critical flows**
  - Login → Dashboard → Create Order → KDS → Complete
  - Login → Products → Create Product → Verify in list
  - Login → HR → Add Employee → Process Payroll

---

## Summary

| Category | Tasks | Priority |
|----------|-------|----------|
| RBAC | 2 tasks | High |
| Performance | 8 tasks | Medium-High |
| Security | 6 tasks | High |
| Frontend Quality | 3 tasks | Medium |
| Backend Quality | 3 tasks | Medium |
| Testing | 2 tasks | Medium |
| **Total** | **24 tasks** | |

### Recommended Order
1. **S1** CSRF protection (security foundation)
2. **R1+R2** Remaining permission migration (stops "Access Denied" flash)
3. **P1** Loading.tsx files (instant perceived performance gain)
4. **P4** Context timeouts (prevents hanging requests)
5. **S5** Rate limiting on auth (prevents brute force)
6. **S6** Input validation (prevents abuse)
7. **P2** True DB pagination (scalability)
8. **P3** HR report batching (performance)
9. **S2** httpOnly cookies (after S1)
10. Everything else in any order
