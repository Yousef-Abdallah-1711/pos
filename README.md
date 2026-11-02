# POS SaaS Platform

**Multi-tenant Point-of-Sale SaaS platform** with restaurant management, online ordering, HR, theme customization, and subscription management.

**Last Updated:** March 3, 2026

## Recent Updates (March 3, 2026)

- Phase 6 (HR + Accounting) completed end-to-end:
  - Added migration range `136..157` and RBAC modules `15..20`.
  - Implemented accounting APIs/UI (accounts, journal entries, ledger, bank accounts, suppliers, customers, invoices, payments, tax, financial reports).
  - Implemented advanced HR APIs/UI (performance, training, insurance, contracts, resignations, HR reports).
  - Added validation/deployment artifacts:
    - `docs-organized/04-MODULES/updated-hr-accounting-system/05-PHASE-6-VALIDATION-REPORT.md`
    - `docs-organized/04-MODULES/updated-hr-accounting-system/06-PHASE-6-PRODUCTION-DEPLOYMENT-CHECKLIST.md`
    - `backend/scripts/rollback_phase6_accounting_hr.sql`

- Theme-builder and storefront integration is now aligned: storefront header/footer render from active theme config instead of shared static layout.
- Backend global compile path restored: `cd backend && go build ./...` now passes after isolating deprecated legacy-only packages behind `legacy` build tags.
- Frontend inline-style hardening completed across active code paths:
  - Removed React inline style props (`style={{...}}` and `style={...}`) from `frontend/**` active files.
  - Converted theme-builder editor previews to scoped class-based styling (`<style jsx>` with generated classes).
  - Converted `frontend/apps/dashboard/src/lib/pdf-export.ts` invoice HTML template from inline `style="..."` attributes to class-based CSS.
- Current style-string exceptions are intentional and limited to sanitizer/testing literals:
  - `frontend/apps/dashboard/src/components/editors/CustomCSSEditor.tsx`
  - `frontend/apps/dashboard/src/components/editors/__tests__/CustomCSSEditor.test.tsx`
- Added Paymob Phase 9 DB-backed verification path:
  - Script: `backend/scripts/phase9_db_verify.ps1`
  - Seed SQL: `backend/scripts/paymob_phase9_seed.sql`
  - Runbook: `docs-organized/11-PAYMOB-PAYMENT-INTEGRATION/PHASE9-DB-VERIFICATION-PATH.md`

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

# Docker
docker compose up --build                                   # DEV: Postgres + API (hot reload) + all 4 frontends (3001-3004, 8080)
cp .env.prod.example .env.prod                              # PROD: fill in real values first
docker compose --env-file .env.prod -f docker-compose.prod.yml up -d --build
```

### Default Credentials
- **Super Admin:** Created via migration `109_seed_super_admin.sql`
- **Regular Users:** Register via Landing app (Port 3001) at `/en/auth/register`

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.18+ (stdlib `net/http`, no framework) |
| **Database** | PostgreSQL 15 (157+ migrations) |
| **Auth** | JWT (HS256, 24h expiry), bcrypt passwords |
| **Frontend** | Next.js 15, React 18, TypeScript 5.7 |
| **Monorepo** | Turbo 2.3 + pnpm 9.0 workspaces |
| **State** | Zustand 4.5 (with persist middleware) |
| **Styling** | Tailwind CSS 3.4, Radix UI, Lucide icons |
| **Forms** | React Hook Form + Zod validation |
| **HTTP** | Axios 1.7 (frontend), stdlib (backend) |
| **i18n** | next-intl (Dashboard/Restaurant) + URL-based routing (Landing) â€” EN + AR with RTL support |
| **Container** | Docker + Docker Compose |

---

## Architecture Overview

```
                    â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
                    â”‚           PostgreSQL 15                  â”‚
                    â”‚    (157+ migrations, row-level tenancy)   â”‚
                    â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¬â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
                                       â”‚
                    â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”´â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
                    â”‚       Go Backend API (:8080)             â”‚
                    â”‚  Clean Architecture (Handlerâ†’UseCaseâ†’Repo)â”‚
                    â”‚  JWT Auth | RBAC | Multi-Tenant Middlewareâ”‚
                    â””â”€â”€â”¬â”€â”€â”€â”€â”€â”€â”€â”€â”¬â”€â”€â”€â”€â”€â”€â”€â”€â”¬â”€â”€â”€â”€â”€â”€â”€â”€â”¬â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
                       â”‚        â”‚        â”‚        â”‚
          â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜   â”Œâ”€â”€â”€â”€â”˜   â”Œâ”€â”€â”€â”€â”˜   â”Œâ”€â”€â”€â”€â”˜
          â–¼                â–¼        â–¼        â–¼
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
- **Restaurant Website:** Subdomain extraction in Next.js middleware (`demo.localhost:3003` -> slug `demo`)
- **Security:** Triple-check validation (JWT -> context -> query WHERE clause)

---

## Project Structure

```
POS/
â”œâ”€â”€ backend/                          Go API server
â”‚   â”œâ”€â”€ cmd/api/main.go              Entry point (1133 lines, all routes)
â”‚   â”œâ”€â”€ internal/
â”‚   â”‚   â”œâ”€â”€ domain/                  Entity models (User, Product, Order, Tenant, etc.)
â”‚   â”‚   â”œâ”€â”€ handler/http/            38+ HTTP handlers
â”‚   â”‚   â”œâ”€â”€ middleware/              11 middleware (auth, tenant, RBAC, status, CORS)
â”‚   â”‚   â”œâ”€â”€ repository/             15+ data access repos
â”‚   â”‚   â”œâ”€â”€ usecase/                9 business logic use cases
â”‚   â”‚   â”œâ”€â”€ service/                Theme + cache services
â”‚   â”‚   â””â”€â”€ pkg/jwt/                Token service
â”‚   â”œâ”€â”€ migrations/                 157+ SQL migration files
â”‚   â”œâ”€â”€ Dockerfile
â”‚   â””â”€â”€ .env.example
â”‚
â”œâ”€â”€ frontend/                         Next.js monorepo
â”‚   â”œâ”€â”€ apps/
â”‚   â”‚   â”œâ”€â”€ dashboard/               Admin dashboard (:3002)
â”‚   â”‚   â”‚   â””â”€â”€ src/
â”‚   â”‚   â”‚       â”œâ”€â”€ app/[locale]/dashboard/   Page routes
â”‚   â”‚   â”‚       â”œâ”€â”€ components/               Feature-organized components
â”‚   â”‚   â”‚       â”œâ”€â”€ stores/                   Zustand stores (auth, tenant, preferences)
â”‚   â”‚   â”‚       â”œâ”€â”€ hooks/                    Custom hooks (RBAC, notifications, etc.)
â”‚   â”‚   â”‚       â”œâ”€â”€ lib/                      API client, debug, validation, utils
â”‚   â”‚   â”‚       â”‚   â”œâ”€â”€ api.ts                Main axios client + interceptors
â”‚   â”‚   â”‚       â”‚   â”œâ”€â”€ api/themeApi.ts        Theme endpoints
â”‚   â”‚   â”‚       â”‚   â””â”€â”€ debug.ts              Debug logging wrapper
â”‚   â”‚   â”‚       â”œâ”€â”€ i18n/                     EN + AR translations
â”‚   â”‚   â”‚       â””â”€â”€ types/                    TypeScript definitions
â”‚   â”‚   â”‚
â”‚   â”‚   â”œâ”€â”€ landing/                 Auth & marketing (:3001)
â”‚   â”‚   â”‚   â””â”€â”€ src/
â”‚   â”‚   â”‚       â”œâ”€â”€ middleware.ts     Locale detection + redirect (/en, /ar)
â”‚   â”‚   â”‚       â”œâ”€â”€ app/
â”‚   â”‚   â”‚       â”‚   â”œâ”€â”€ layout.tsx   Root layout (ThemeProvider only)
â”‚   â”‚   â”‚       â”‚   â”œâ”€â”€ page.tsx     Redirect â†’ /en
â”‚   â”‚   â”‚       â”‚   â””â”€â”€ [locale]/    URL-based i18n (en | ar)
â”‚   â”‚   â”‚       â”‚       â”œâ”€â”€ layout.tsx       LanguageProvider + Toaster
â”‚   â”‚   â”‚       â”‚       â”œâ”€â”€ page.tsx         Landing page (all sections)
â”‚   â”‚   â”‚       â”‚       â””â”€â”€ auth/
â”‚   â”‚   â”‚       â”‚           â”œâ”€â”€ login/       Login page
â”‚   â”‚   â”‚       â”‚           â”œâ”€â”€ register/    Register page
â”‚   â”‚   â”‚       â”‚           â””â”€â”€ super-admin-login/
â”‚   â”‚   â”‚       â”œâ”€â”€ components/      Navbar, Hero, Features, HowItWorks,
â”‚   â”‚   â”‚       â”‚                    Pricing, Testimonials, CTA, Footer
â”‚   â”‚   â”‚       â”œâ”€â”€ contexts/        LanguageContext (URL-based locale)
â”‚   â”‚   â”‚       â”œâ”€â”€ lib/             translations.ts, config.ts
â”‚   â”‚   â”‚       â””â”€â”€ stores/          authStore (Zustand)
â”‚   â”‚   â”‚
â”‚   â”‚   â”œâ”€â”€ restaurant-website/      Public storefront (:3003)
â”‚   â”‚   â”‚   â”œâ”€â”€ middleware.ts        Subdomain + locale + tenant status
â”‚   â”‚   â”‚   â””â”€â”€ lib/
â”‚   â”‚   â”‚       â”œâ”€â”€ store/           Cart + Theme stores
â”‚   â”‚   â”‚       â”œâ”€â”€ api/             API clients
â”‚   â”‚   â”‚       â”œâ”€â”€ hooks/           Theme, web vitals
â”‚   â”‚   â”‚       â”œâ”€â”€ debug.ts         Debug logging wrapper
â”‚   â”‚   â”‚       â””â”€â”€ checkTenantStatus.ts
â”‚   â”‚   â”‚
â”‚   â”‚   â””â”€â”€ super-admin/             Platform admin (:3004)
â”‚   â”‚       â””â”€â”€ src/
â”‚   â”‚           â”œâ”€â”€ components/tenants/   Extend trial, disable, delete
â”‚   â”‚           â”œâ”€â”€ components/impersonation/  Impersonation UI
â”‚   â”‚           â””â”€â”€ lib/api.ts           SuperAdminAPI class
â”‚   â”‚
â”‚   â”œâ”€â”€ packages/
â”‚   â”‚   â”œâ”€â”€ ui/                      Radix UI + CVA components
â”‚   â”‚   â”œâ”€â”€ ui-components/           Atomic design library
â”‚   â”‚   â””â”€â”€ theme/                   next-themes provider
â”‚   â”‚
â”‚   â””â”€â”€ turbo.json                   Build pipeline config
â”‚
â”œâ”€â”€ docker-compose.yml               PostgreSQL + Backend
â”œâ”€â”€ docs-organized/                   70+ documentation files
â””â”€â”€ README.md                         This file
```

---

## Backend Architecture

### Clean Architecture Layers

```
HTTP Request â†’ CORS â†’ Auth â†’ TenantContext â†’ [RBAC] â†’ [TenantStatus] â†’ Handler â†’ UseCase â†’ Repository â†’ PostgreSQL
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
wrapProtected := AuthMiddleware â†’ TenantContextMiddleware â†’ Handler

// RBAC-protected route
wrapWithPermission := AuthMiddleware â†’ TenantContext â†’ WithRequiredPermission(moduleID, level) â†’ Handler

// With subscription check
wrapProtectedWithStatus := AuthMiddleware â†’ TenantContext â†’ TenantStatus â†’ Handler

// Super admin route
wrapSuperAdmin := RateLimit â†’ SuperAdminAuth â†’ IPWhitelist â†’ RequireSuperAdmin â†’ Handler
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
  â†’ Create Tenant (trial 15 days) â†’ Restaurant â†’ User (role=owner) â†’ JWT

Login: POST /auth/login
  â†’ Verify bcrypt password â†’ Check multi-tenant
  â†’ If 1 tenant: return JWT
  â†’ If N tenants: return tenant list â†’ POST /auth/login/confirm â†’ JWT

JWT Claims: { user_id, tenant_id, restaurant_id, email, role, exp }
Super Admin JWT: Separate secret (JWT_SECRET + "_super_admin"), 15min expiry
```

### Database Schema (Key Tables)

```
tenants          â†’ id, name, slug, subscription_status, trial_dates, is_disabled
users            â†’ id, tenant_id, restaurant_id, email, password_hash, role, is_super_admin
restaurants      â†’ id, tenant_id, name, slug, theme, feature_flags
products         â†’ id, tenant_id, restaurant_id, name_en/ar, price, inventory
categories       â†’ id, tenant_id, restaurant_id, name_en/ar, display_order
orders           â†’ id, tenant_id, restaurant_id, order_number, customer_info, items, status
drivers          â†’ id, tenant_id, name, status, availability, rating
employees        â†’ id, tenant_id, employee_number, salary, department
rbac_roles       â†’ id, tenant_id, code, name, is_system
role_permissions â†’ id, role_id, module_id, permission_level
user_roles       â†’ id, user_id, role_id
super_admins     â†’ id, user_id, ip_whitelist, secret_key
impersonation_sessions â†’ id, super_admin_id, target_user_id, token_hash, expires_at
subscription_history   â†’ id, tenant_id, action, old/new_status, days_added, reason
themes           â†’ id, tenant_id, restaurant_id, name, config, is_active
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
demo.localhost:3003/en/menu
  â†“ middleware.ts
  1. Extract slug from subdomain: "demo"
  2. Extract locale from path: "en"
  3. Validate locale (redirect if invalid)
  4. Check tenant status via API
  5. If blocked â†’ redirect to /maintenance
  6. Set headers: x-restaurant-slug, x-restaurant-locale
```

**Theme System:**
- Server-side theme loading (`get-theme-server.ts`)
- Theme payload normalization across public API response variants (`data.data`, `data.config`, `theme.config`, JSON string config)
- Multi-layer cache: memory â†’ localStorage â†’ API (with retry + backoff)
- Theme data: colors, typography, identity, header, footer, components array
- Storefront header/footer are rendered directly from active theme config (no static global header/footer)
- Applied via CSS variables and Tailwind classes

**Cart/Checkout:**
- Zustand cart store persisted to localStorage
- Cart â†’ Checkout form â†’ `POST /public/orders` â†’ Order confirmation
- Inventory-aware quantity limits

### Landing App (:3001)

**URL-based i18n routing** â€” all pages under `/en/...` or `/ar/...`. Middleware redirects bare paths to `/en`.

**Landing Page Sections:**
- `Navbar` â€” Sticky glassmorphism, orange/red gradient CTA, language toggle (URL-based), mobile menu
- `Hero` â€” Dark POS mockup, animated floating stats cards, badge "30-Day Free Trial"
- `Features` â€” 6 feature cards (POS, Kitchen Display, Online Ordering, Analytics, Multi-Restaurant, Mobile Apps)
- `HowItWorks` â€” 3-step process: Sign Up â†’ Set Up â†’ Start Selling
- `Pricing` â€” Single plan: **500 EGP/month** after 30-day free trial (no credit card required)
- `Testimonials` â€” 3 testimonials with star ratings
- `CTA` â€” Red gradient, full-width call-to-action
- `Footer` â€” All link text from translations (EN + AR)

**Auth Pages:** `/en/auth/login`, `/en/auth/register`, `/en/auth/super-admin-login`

**i18n Architecture:**
- `middleware.ts` â€” Validates locale segment, redirects to `/en` if missing
- `LanguageContext.tsx` â€” Accepts `locale` prop from server, toggle uses `router.push()` to swap URL segment
- Language toggle: `setLanguage('ar')` â†’ navigates to `/ar/current-path`
- RTL: `document.documentElement.dir` set via `useEffect` in LanguageProvider

**Multi-tenant login flow:**
- `TenantSelectorModal` â€” shown when user has multiple tenants
- On tenant select: `POST /auth/login/confirm` â†’ redirect to Dashboard auth callback

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

**Dashboard + Restaurant Website:**
- Locale from URL path: `/en/dashboard` or `/ar/dashboard`
- `createTranslator(locale)` returns `t()` function
- `isRTL = locale === 'ar'` for conditional RTL classes
- Root layout sets `dir="rtl"` attribute
- CSS handles RTL via `[dir="rtl"]` selectors

**Landing App:**
- Locale from URL segment: `/en/...` or `/ar/...`
- `middleware.ts` redirects root and locale-less paths to `/en`
- `LanguageProvider` accepts `locale` prop (from server, no localStorage)
- Language toggle calls `router.push()` swapping URL segment (`/en/` â†” `/ar/`)
- RTL applied via `useEffect` â†’ `document.documentElement.dir`
- Translations in `src/lib/translations.ts` (EN + AR, inline TypeScript object)

### Tenant Status Flow

```
Dashboard mounts â†’ TenantStatusProvider calls GET /tenant/status
  â†’ If canAccess: true â†’ show dashboard normally
  â†’ If canAccess: false â†’ show BlockedOverlay (trial_expired / disabled)

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

### Pricing

- **Free Trial:** 30 days (seeded via migration `112_seed_trial_data.sql`)
- **Subscription:** 500 EGP/month (all features included, no tiers)
- **No credit card required** to start trial

### Lifecycle

```
Tenant registers â†’ trial (30 days) â†’ [extend trial / subscribe]
                                    â†’ trial_expired â†’ blocked
                                    â†’ [admin disables] â†’ disabled â†’ blocked
                                    â†’ [admin enables] â†’ restored
                                    â†’ [admin deletes] â†’ soft delete
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

1. **Backend:** Domain model â†’ Repository â†’ UseCase â†’ Handler â†’ Register route in `main.go`
2. **Frontend:** Page in `app/[locale]/dashboard/` â†’ Components â†’ Hook/API call â†’ Add to Sidebar
3. **Translation:** Add keys to `i18n/messages/en.json` and `ar.json`
4. **RBAC:** Add module in migration â†’ Register permission checks

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
- **Subdomain testing on localhost:** If `demo.localhost` DNS is unavailable, test using Host header (`curl -H "Host: demo.localhost:3003" http://localhost:3003/en`)
- **Debug logs:** All `console.log/error/warn` wrapped in `debugLog/debugError/debugWarn` - enable with `NEXT_PUBLIC_DEBUG_LOGS=true`
- **Super admin JWT:** Different secret (`JWT_SECRET + "_super_admin"`), 15min expiry
- **Windows file-lock issue (Next.js):** If dev fails with `EPERM ... .next/trace`, stop node processes using that app, delete its `.next` folder, then restart `pnpm dev`.
- **Monorepo dev run:** `turbo run dev --parallel` can fail from non-critical workspace package errors. For focused work use package filters (for example `pnpm --filter @pos-saas/dashboard dev`).

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

