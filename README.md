# POS SaaS Platform

**Multi-tenant Point-of-Sale SaaS platform** with restaurant management, online ordering, HR, theme customization, and subscription management.

**Last Updated:** February 11, 2026

---

## Quick Start

```bash
# Backend (Go API - Port 8080)
cd backend
cp .env.example .env        # Configure DB, JWT_SECRET
go run cmd/api/main.go

# Frontend (All 4 apps via Turbo)
cd frontend
pnpm install
pnpm dev                    # Starts all apps in parallel

# Or run individual apps:
pnpm dev:dashboard          # Port 3002
pnpm dev:landing            # Port 3001
pnpm dev:restaurant         # Port 3003
pnpm dev:super-admin        # Port 3004

# Docker (Backend + PostgreSQL)
docker-compose up -d
```

### Default Credentials
- **Super Admin:** Created via migration `109_seed_super_admin.sql`
- **Regular Users:** Register via Landing app (Port 3001) `/auth/register`

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.18+ (stdlib `net/http`, no framework) |
| **Database** | PostgreSQL 15 (58+ migrations) |
| **Auth** | JWT (HS256, 24h expiry), bcrypt passwords |
| **Frontend** | Next.js 15, React 18, TypeScript 5.7 |
| **Monorepo** | Turbo 2.3 + pnpm 9.0 workspaces |
| **State** | Zustand 4.5 (with persist middleware) |
| **Styling** | Tailwind CSS 3.4, Radix UI, Lucide icons |
| **Forms** | React Hook Form + Zod validation |
| **HTTP** | Axios 1.7 (frontend), stdlib (backend) |
| **i18n** | next-intl (EN + AR with RTL support) |
| **Container** | Docker + Docker Compose |

---

## Architecture Overview

```
                    ┌─────────────────────────────────────────┐
                    │           PostgreSQL 15                  │
                    │    (58+ migrations, row-level tenancy)   │
                    └──────────────────┬──────────────────────┘
                                       │
                    ┌──────────────────┴──────────────────────┐
                    │       Go Backend API (:8080)             │
                    │  Clean Architecture (Handler→UseCase→Repo)│
                    │  JWT Auth | RBAC | Multi-Tenant Middleware│
                    └──┬────────┬────────┬────────┬───────────┘
                       │        │        │        │
          ┌────────────┘   ┌────┘   ┌────┘   ┌────┘
          ▼                ▼        ▼        ▼
   Landing (:3001)   Dashboard   Restaurant  Super-Admin
   Auth/Register     (:3002)     Website     (:3004)
   Marketing page    Admin UI    (:3003)     Tenant mgmt
                     CRUD ops    Subdomain   Impersonation
                     HR/Orders   multi-tenant Audit logs
                     Theme edit  Cart/Checkout Subscriptions
```

### Multi-Tenancy Model

- **Database:** Row-level isolation via `tenant_id` on every table
- **Backend:** `TenantContextMiddleware` extracts tenant from JWT claims
- **Dashboard:** `X-Tenant-ID` / `X-Restaurant-ID` headers on every API call
- **Restaurant Website:** Subdomain extraction in Next.js middleware (`demo.localhost:3000` -> slug `demo`)
- **Security:** Triple-check validation (JWT -> context -> query WHERE clause)

---

## Project Structure

```
POS/
├── backend/                          Go API server
│   ├── cmd/api/main.go              Entry point (1133 lines, all routes)
│   ├── internal/
│   │   ├── domain/                  Entity models (User, Product, Order, Tenant, etc.)
│   │   ├── handler/http/            38+ HTTP handlers
│   │   ├── middleware/              11 middleware (auth, tenant, RBAC, status, CORS)
│   │   ├── repository/             15+ data access repos
│   │   ├── usecase/                9 business logic use cases
│   │   ├── service/                Theme + cache services
│   │   └── pkg/jwt/                Token service
│   ├── migrations/                  58+ SQL migration files
│   ├── Dockerfile
│   └── .env.example
│
├── frontend/                         Next.js monorepo
│   ├── apps/
│   │   ├── dashboard/               Admin dashboard (:3002)
│   │   │   └── src/
│   │   │       ├── app/[locale]/dashboard/   Page routes
│   │   │       ├── components/               Feature-organized components
│   │   │       ├── stores/                   Zustand stores (auth, tenant, preferences)
│   │   │       ├── hooks/                    Custom hooks (RBAC, notifications, etc.)
│   │   │       ├── lib/                      API client, debug, validation, utils
│   │   │       │   ├── api.ts                Main axios client + interceptors
│   │   │       │   ├── api/themeApi.ts        Theme endpoints
│   │   │       │   └── debug.ts              Debug logging wrapper
│   │   │       ├── i18n/                     EN + AR translations
│   │   │       └── types/                    TypeScript definitions
│   │   │
│   │   ├── landing/                 Auth & marketing (:3001)
│   │   │   └── src/app/auth/        Login, Register pages
│   │   │
│   │   ├── restaurant-website/      Public storefront (:3003)
│   │   │   ├── middleware.ts        Subdomain + locale + tenant status
│   │   │   └── lib/
│   │   │       ├── store/           Cart + Theme stores
│   │   │       ├── api/             API clients
│   │   │       ├── hooks/           Theme, web vitals
│   │   │       ├── debug.ts         Debug logging wrapper
│   │   │       └── checkTenantStatus.ts
│   │   │
│   │   └── super-admin/             Platform admin (:3004)
│   │       └── src/
│   │           ├── components/tenants/   Extend trial, disable, delete
│   │           ├── components/impersonation/  Impersonation UI
│   │           └── lib/api.ts           SuperAdminAPI class
│   │
│   ├── packages/
│   │   ├── ui/                      Radix UI + CVA components
│   │   ├── ui-components/           Atomic design library
│   │   └── theme/                   next-themes provider
│   │
│   └── turbo.json                   Build pipeline config
│
├── docker-compose.yml               PostgreSQL + Backend
├── docs-organized/                   70+ documentation files
└── README.md                         This file
```

---

## Backend Architecture

### Clean Architecture Layers

```
HTTP Request → CORS → Auth → TenantContext → [RBAC] → [TenantStatus] → Handler → UseCase → Repository → PostgreSQL
```

### Middleware Chain

| Middleware | Purpose |
|-----------|---------|
| `CORSMiddleware` | Cross-origin headers |
| `AuthMiddleware` | JWT validation (HS256), stores claims in context |
| `TenantContextMiddleware` | Extracts tenant_id + restaurant_id from JWT claims |
| `WithRequiredPermission` | RBAC check (module_id + permission_level) |
| `TenantStatusMiddleware` | Blocks expired trials / disabled tenants |
| `SuperAdminAuthMiddleware` | Separate JWT secret (15min expiry) |
| `RateLimitSuperAdmin` | 30 req/min for super admin |
| `ValidateIPWhitelist` | IP allowlist for super admin |
| `ImpersonationMiddleware` | Validates `X-Impersonation-Token` header |

### Route Wrapping Pattern (main.go)

```go
// Standard protected route
wrapProtected := AuthMiddleware → TenantContextMiddleware → Handler

// RBAC-protected route
wrapWithPermission := AuthMiddleware → TenantContext → WithRequiredPermission(moduleID, level) → Handler

// With subscription check
wrapProtectedWithStatus := AuthMiddleware → TenantContext → TenantStatus → Handler

// Super admin route
wrapSuperAdmin := RateLimit → SuperAdminAuth → IPWhitelist → RequireSuperAdmin → Handler
```

### API Endpoints (120+)

**Auth (Public):**
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/login/confirm        (multi-tenant selection)
POST   /api/v1/auth/token-exchange        (super admin impersonation)
GET    /api/v1/auth/check-subdomain
```

**Public Restaurant API (No Auth):**
```
GET    /api/v1/public/restaurants/{slug}/menu
GET    /api/v1/public/restaurants/{slug}/products
GET    /api/v1/public/restaurants/{slug}/categories
GET    /api/v1/public/restaurants/{slug}/info
GET    /api/v1/public/restaurants/{slug}/homepage
GET    /api/v1/public/restaurants/{slug}/theme
POST   /api/v1/public/orders
GET    /api/v1/public/orders/{id}
GET    /api/v1/public/orders/number/{number}/lookup
GET    /api/v1/public/tenant-status/{slug}
```

**Protected Admin API:**
```
GET/POST/PUT/DELETE  /api/v1/products        (RBAC: module 1)
GET/POST/PUT/DELETE  /api/v1/categories
GET/PUT/DELETE       /api/v1/admin/orders     (RBAC: module 10)
GET/POST/PUT/DELETE  /api/v1/admin/drivers
GET/POST             /api/v1/notifications
GET/POST/PUT/DELETE  /api/v1/rbac/roles
POST/DELETE/GET      /api/v1/rbac/permissions
GET/PUT              /api/v1/users/{id}/settings
GET/POST/PUT/DELETE  /api/v1/admin/themes
GET                  /api/v1/tenant/status
```

**Super Admin API:**
```
POST   /api/v1/super-admin/auth/login
GET    /api/v1/super-admin/dashboard/stats
GET    /api/v1/super-admin/tenants
GET    /api/v1/super-admin/tenants/{id}
POST   /api/v1/super-admin/tenants/{id}/extend-trial
POST   /api/v1/super-admin/tenants/{id}/disable
POST   /api/v1/super-admin/tenants/{id}/enable
DELETE /api/v1/super-admin/tenants/{id}
GET    /api/v1/super-admin/tenants/{id}/subscription-history
POST   /api/v1/super-admin/impersonate/{userID}
GET    /api/v1/super-admin/audit-logs
```

### Authentication Flow

```
Register: POST /auth/register
  → Create Tenant (trial 15 days) → Restaurant → User (role=owner) → JWT

Login: POST /auth/login
  → Verify bcrypt password → Check multi-tenant
  → If 1 tenant: return JWT
  → If N tenants: return tenant list → POST /auth/login/confirm → JWT

JWT Claims: { user_id, tenant_id, restaurant_id, email, role, exp }
Super Admin JWT: Separate secret (JWT_SECRET + "_super_admin"), 15min expiry
```

### Database Schema (Key Tables)

```
tenants          → id, name, slug, subscription_status, trial_dates, is_disabled
users            → id, tenant_id, restaurant_id, email, password_hash, role, is_super_admin
restaurants      → id, tenant_id, name, slug, theme, feature_flags
products         → id, tenant_id, restaurant_id, name_en/ar, price, inventory
categories       → id, tenant_id, restaurant_id, name_en/ar, display_order
orders           → id, tenant_id, restaurant_id, order_number, customer_info, items, status
drivers          → id, tenant_id, name, status, availability, rating
employees        → id, tenant_id, employee_number, salary, department
rbac_roles       → id, tenant_id, code, name, is_system
role_permissions → id, role_id, module_id, permission_level
user_roles       → id, user_id, role_id
super_admins     → id, user_id, ip_whitelist, secret_key
impersonation_sessions → id, super_admin_id, target_user_id, token_hash, expires_at
subscription_history   → id, tenant_id, action, old/new_status, days_added, reason
themes           → id, tenant_id, restaurant_id, name, config, is_active
```

---

## Frontend Architecture

### Dashboard App (:3002)

**Page Routes:** (`/[locale]/dashboard/...`)
```
/dashboard                    Overview (stats, charts)
/dashboard/products           Product CRUD (list, new, [id]/edit)
/dashboard/categories         Category CRUD
/dashboard/orders             Order management ([id] detail)
/dashboard/drivers            Driver management (new, [id])
/dashboard/inventory          Stock management
/dashboard/hr/employees       Employee CRUD
/dashboard/hr/attendance      Attendance tracking
/dashboard/hr/leaves          Leave management
/dashboard/hr/payroll         Payroll/salary
/dashboard/hr/roles           HR roles
/dashboard/settings           User settings (profile, appearance, access-control)
/dashboard/settings/users     User management
/dashboard/settings/roles     Role CRUD
/dashboard/theme-builder      Theme editor (presets, editor/[id])
/dashboard/rbac/roles         RBAC role management
```

**Zustand Stores:**
| Store | Key | Persisted | Purpose |
|-------|-----|-----------|---------|
| `authStore` | `auth-storage` | Yes (localStorage) | User, token, login/logout |
| `tenantStatusStore` | N/A | No | Trial/subscription status |
| `preferencesStore` | `preferences-storage` | Yes | UI settings (theme, locale) |
| `themeBuilderStore` | N/A | No | Theme editor state |

**API Client (`lib/api.ts`):**
- Axios instance with base URL `http://localhost:8080/api/v1`
- Request interceptor: Adds `Authorization: Bearer {token}`, `X-Tenant-ID`, `X-Restaurant-ID`
- Impersonation: Uses `X-Impersonation-Token` from sessionStorage instead of Bearer
- Response interceptor: Auto-logout on 401 (redirects to landing app)

**Key Hooks:**
| Hook | Purpose |
|------|---------|
| `useTenantStatus()` | Returns `{ status, daysRemaining, isBlocked, isLoading }` |
| `useRbac()` / `useHasPermission()` | RBAC permission checking |
| `useNotifications()` | Notification CRUD + stats |
| `useWebVitals()` | Performance monitoring (LCP, CLS, FID) |
| `useServiceWorker()` | Offline support |

**Layout Components:**
- `DashboardLayout.tsx` - Main wrapper (sidebar + header + content)
- `Sidebar.tsx` - Navigation with expandable groups + `SidebarSubscriptionBadge`
- `TrialWarningBanner.tsx` - Amber banner when trial <= 7 days remaining

### Restaurant Website (:3003)

**Subdomain Multi-Tenancy:**
```
demo.localhost:3000/en/menu
  ↓ middleware.ts
  1. Extract slug from subdomain: "demo"
  2. Extract locale from path: "en"
  3. Validate locale (redirect if invalid)
  4. Check tenant status via API
  5. If blocked → redirect to /maintenance
  6. Set headers: x-restaurant-slug, x-restaurant-locale
```

**Theme System:**
- Server-side theme loading (`get-theme-server.ts`)
- Multi-layer cache: memory → localStorage → API (with retry + backoff)
- Theme data: colors, typography, identity, header, footer, components array
- Applied via CSS variables and Tailwind classes

**Cart/Checkout:**
- Zustand cart store persisted to localStorage
- Cart → Checkout form → `POST /public/orders` → Order confirmation
- Inventory-aware quantity limits

### Landing App (:3001)

- Marketing homepage (Hero, Features, Pricing, Testimonials)
- Auth pages: `/auth/login`, `/auth/register`
- Multi-tenant login flow (tenant selector modal)
- Redirects to Dashboard after successful login

### Super-Admin App (:3004)

- Platform-wide tenant management
- Tenant list with search/pagination
- Tenant detail: extend trial, disable/enable, delete
- Subscription history with audit trail
- User impersonation (generates token, opens Dashboard in new tab)
- Audit log viewer (filter by action, date range)
- Dashboard stats (total tenants, active trials, etc.)

---

## Key Patterns & Conventions

### Debug Logging

All console calls go through debug wrappers that respect environment variables:

```typescript
// Dashboard: import { debugLog, debugError, debugWarn } from '@/lib/debug'
// Restaurant: import { debugLog, debugError, debugWarn } from '@/lib/debug'

// Controlled by: NEXT_PUBLIC_DEBUG_LOGS=true (client) or DEBUG_LOGS=true (server)
// Production: disabled by default, zero console output
```

### i18n & RTL

- Locale from URL path: `/en/dashboard` or `/ar/dashboard`
- `createTranslator(locale)` returns `t()` function
- `isRTL = locale === 'ar'` for conditional RTL classes
- Root layout sets `dir="rtl"` attribute
- CSS handles RTL via `[dir="rtl"]` selectors

### Tenant Status Flow

```
Dashboard mounts → TenantStatusProvider calls GET /tenant/status
  → If canAccess: true → show dashboard normally
  → If canAccess: false → show BlockedOverlay (trial_expired / disabled)

Sidebar shows: SidebarSubscriptionBadge (plan name + days remaining)
Header shows: TrialWarningBanner (when <= 7 days remaining, dismissible)
```

### RBAC

```
Modules: Products(1), Categories(2), Orders(10), HR(5), etc.
Levels: READ, WRITE, DELETE, ADMIN

Backend: wrapWithPermission(handler, moduleID, "WRITE")
Frontend: useHasPermission(MODULE_IDS.PRODUCTS, PERMISSION_LEVELS.EDIT)
```

### Error Handling

```json
// Standard API response
{ "success": true, "data": { ... } }
{ "success": false, "error": "message", "code": "ERROR_CODE" }

// HTTP Status Codes
200 - Success (GET/PUT/DELETE)
201 - Created (POST)
400 - Bad Request (validation)
401 - Unauthorized (missing/invalid JWT)
403 - Forbidden (RBAC / tenant blocked)
404 - Not Found
500 - Server Error
```

---

## SaaS Trial & Subscription System

### Lifecycle

```
Tenant registers → trial (15 days) → [extend trial / upgrade]
                                    → trial_expired → blocked
                                    → [admin disables] → disabled → blocked
                                    → [admin enables] → restored
                                    → [admin deletes] → soft delete
```

### Backend (Fully Implemented)

- 8 API endpoints for subscription management
- `TenantStatusMiddleware` blocks expired/disabled tenants
- `subscription_history` table for audit trail
- Domain methods: `IsTrialActive()`, `CanAccess()`, `DaysRemainingInTrial()`

### Frontend Status

| Feature | Dashboard (:3002) | Restaurant (:3003) | Super-Admin (:3004) |
|---------|-------------------|---------------------|---------------------|
| Blocking overlay | Done | Done (maintenance page) | N/A |
| Sidebar badge | Done | N/A | N/A |
| Trial warning | Done (banner) | N/A | N/A |
| Extend trial | N/A | N/A | Done |
| Disable/Enable | N/A | N/A | Done |
| Delete tenant | N/A | N/A | Done |
| History view | N/A | N/A | Done |
| Impersonation | Accepts token | N/A | Done (generates token) |

---

## Environment Variables

### Backend (.env)
```bash
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=pos_saas
DB_SSLMODE=disable
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=http://localhost:3001,http://localhost:3002,http://localhost:3003,http://localhost:3004
```

### Frontend (.env.local)
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_DASHBOARD_URL=http://localhost:3002
NEXT_PUBLIC_LANDING_URL=http://localhost:3001
NEXT_PUBLIC_DEBUG_LOGS=false          # Set true for console output
```

---

## Database Migrations

```
001-005   Core: tenants, restaurants, users, categories, products
012-017   Products: images, inventory, audit logs
020-027   HR: employees, roles, attendance, salaries, leaves, reviews
040       User settings
050       Notifications
060-066   Orders: orders, items, status_history, indexes
067-090   Themes: themes, presets, components, history, sections
070-072   Drivers: drivers, assignments, location_history
099-106   RBAC: modules, permission_levels, roles, permissions, user_roles
107-109   Super Admin: super_admins, impersonation, audit_logs
110-112   Subscription: tenant fields, history, trial seeding
```

Run migrations: `go run cmd/migrate/main.go` or via Docker entrypoint.

---

## Development Notes

### Adding a New Feature

1. **Backend:** Domain model → Repository → UseCase → Handler → Register route in `main.go`
2. **Frontend:** Page in `app/[locale]/dashboard/` → Components → Hook/API call → Add to Sidebar
3. **Translation:** Add keys to `i18n/messages/en.json` and `ar.json`
4. **RBAC:** Add module in migration → Register permission checks

### File Naming Conventions

```
Backend:   snake_case.go (product_handler.go, auth_repo.go)
Frontend:  PascalCase.tsx (ProductForm.tsx), camelCase.ts (authStore.ts)
Pages:     page.tsx (Next.js convention)
Hooks:     use*.ts (useRbac.ts, useTenantStatus.ts)
Stores:    *Store.ts (authStore.ts, tenantStatusStore.ts)
```

### Common Gotchas

- **Tenant context:** Every API call needs `X-Tenant-ID` header (handled by axios interceptor)
- **Impersonation:** Uses `sessionStorage` (not localStorage) - clears on tab close
- **Locale:** URL must include locale (`/en/dashboard`, not `/dashboard`)
- **Debug logs:** All `console.log/error/warn` wrapped in `debugLog/debugError/debugWarn` - enable with `NEXT_PUBLIC_DEBUG_LOGS=true`
- **Super admin JWT:** Different secret (`JWT_SECRET + "_super_admin"`), 15min expiry

---

## Detailed Documentation

See [docs-organized/](docs-organized/) for in-depth guides:

| Topic | Path |
|-------|------|
| Getting Started | `docs-organized/00-GETTING-STARTED/` |
| Theme System | `docs-organized/01-THEME-SYSTEM/` |
| API Reference | `docs-organized/02-API-BACKEND/API_REFERENCE_V2.0.0.md` |
| Frontend Architecture | `docs-organized/03-FRONTEND-DASHBOARD/` |
| HR Module | `docs-organized/04-MODULES/HR/` |
| Multi-Tenancy | `docs-organized/05-INFRASTRUCTURE/MULTI_TENANCY/` |
| Deployment | `docs-organized/05-INFRASTRUCTURE/DEPLOYMENT/` |
| SaaS Subscription | `docs-organized/10-Sass-Trial-Subscription/` |
| Remaining Work | `docs-organized/99-REMAINING-WORK.md` |
