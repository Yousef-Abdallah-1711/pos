# POS SaaS - Project Knowledge Base

**Version:** 3.0.0
**Status:** Production Ready (Multi-tenant, RBAC, Real-time, Payment Integration)
**Last Updated:** 2026-02-25
**Target Scale:** 1M+ Users / Multi-tenant SaaS

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture & Tech Stack](#architecture--tech-stack)
3. [Multi-Tenant Architecture](#multi-tenant-architecture)
4. [Frontend Monorepo](#frontend-monorepo)
5. [Backend Structure](#backend-structure)
6. [RBAC (Role-Based Access Control)](#rbac-role-based-access-control)
7. [Database Schema & Migrations](#database-schema--migrations)
8. [API Reference](#api-reference)
9. [Feature Modules](#feature-modules)
10. [Development Guidelines](#development-guidelines)
11. [Adding New Features](#adding-new-features)
12. [Scalability & Performance](#scalability--performance)
13. [Security Best Practices](#security-best-practices)
14. [Testing Strategy](#testing-strategy)
15. [Deployment & Infrastructure](#deployment--infrastructure)
16. [Common Issues & Solutions](#common-issues--solutions)

---

## Project Overview

**POS SaaS** is a multi-tenant, production-grade Point-of-Sale system for restaurant management:

- **Multi-tenant Support** - Isolated data per restaurant chain/franchise
- **RBAC** - 14 modules with 6 hierarchical permission levels
- **4 Frontend Apps** - Dashboard, Restaurant Website, Landing Page, Super Admin
- **Real-time KDS** - WebSocket-powered Kitchen Display System
- **Table Management** - Dine-in table ordering with QR codes
- **Theme Builder** - WYSIWYG section-based website builder
- **Payment Processing** - Paymob integration (live + mock mode) with webhook retry queue
- **Subscription/Billing** - Trial/paid plans with automatic expiry
- **Reporting** - Sales, profitability, inventory reports with PDF export
- **HR Module** - Employees, attendance, leaves, payroll
- **Internationalization** - English + Arabic (RTL) with next-intl
- **Dark Mode** - System/light/dark theme switching
- **Offline Support** - IndexedDB queues for KDS offline operations

**Planned:** Cashier POS terminal (see `docs-organized/16-CASHIER-POS/`)

---

## Architecture & Tech Stack

### Backend
- **Language:** Go 1.24
- **Framework:** stdlib `net/http` with `http.NewServeMux()` (no heavy frameworks)
- **Database:** PostgreSQL 15 (row-level tenant isolation)
- **Cache:** Redis (sessions, tokens, rate limiting)
- **WebSocket:** Custom KDS hub for real-time kitchen events
- **Payment:** Paymob integration (switchable live/mock via `PAYMOB_MODE` env)
- **Authentication:** JWT (24h expiry) - separate secrets for users vs super admin
- **Background Workers:** Webhook retry (exponential backoff), subscription expiry checker
- **PDF Generation:** Embedded fonts for report export

### Frontend
- **Monorepo:** pnpm workspace with 4 apps + 2 shared layers
- **Framework:** Next.js 15 (React 18) - TypeScript
- **State Management:** Zustand (with persist middleware + IndexedDB)
- **Styling:** Tailwind CSS + dark mode
- **Components:** Shadcn/ui (headless, customizable)
- **i18n:** next-intl for routing-level internationalization
- **Forms:** React Hook Form + Zod validation
- **Data Fetching:** React Query (TanStack Query) for restaurant-website

### DevOps
- **Containerization:** Docker + Docker Compose
- **CI/CD:** GitHub Actions (localhost URL guardrail + domain config smoke tests)
- **Monitoring:** Request metrics middleware, health/ready/metrics endpoints

---

## Multi-Tenant Architecture

### Tenant Hierarchy

```
Super Admin Layer              (Master auth, tenant management, impersonation)
  |
  +-- Tenant 1 (Company A)    (franchise/chain)
  |   +-- Restaurant 1
  |   +-- Restaurant 2
  |   +-- Restaurant 3
  |
  +-- Tenant 2 (Company B)
      +-- Restaurant A
      +-- Restaurant B
```

### Data Isolation Strategy

**Database Level (PostgreSQL):**
- Single schema, row-level isolation
- `tenant_id` + `restaurant_id` in every query WHERE clause
- Composite indexes on `(tenant_id, restaurant_id)` for performance
- `TenantStatusMiddleware` blocks suspended/expired tenants

**Application Level (Go):**
- `TenantContextMiddleware` extracts tenant_id + restaurant_id from JWT claims
- All repository queries automatically scoped by tenant
- `WithRequiredPermission(moduleID, level)` middleware gates endpoints

**Frontend Level:**
- Restaurant website uses subdomain-based routing (`{slug}.domain.com`)
- Dashboard injects `X-Tenant-ID` + `X-Restaurant-ID` via Axios interceptors
- Auth store persists JWT with tenant context

---

## Frontend Monorepo

### Apps (`frontend/apps/`)

| App | Port | Purpose |
|-----|------|---------|
| `dashboard` | 3002 | Admin dashboard (products, orders, KDS, tables, reports, HR, theme builder, settings) |
| `restaurant-website` | 3003 | Customer-facing storefront per restaurant (menu, cart, checkout, order tracking) |
| `landing` | 3001 | Marketing landing page + tenant registration |
| `super-admin` | 3004 | Super admin portal (tenant management, impersonation, payment monitoring) |

### Shared Layers (`frontend/layers/`)

| Layer | Purpose |
|-------|---------|
| `ui-themes` | Theme rendering engine + section components (shared by dashboard theme builder & restaurant website) |
| `shared-components` | Shared UI primitives |

### Dashboard App Structure

```
frontend/apps/dashboard/src/
+-- app/[locale]/dashboard/
|   +-- page.tsx              (home/overview)
|   +-- orders/               (order list + detail)
|   +-- products/             (product list + form)
|   +-- categories/
|   +-- kds/                  (Kitchen Display System - fullscreen)
|   +-- tables/               (table management + layout)
|   +-- reports/              (overview, sales, profitability, inventory)
|   +-- hr/                   (employees, attendance, leaves, payroll)
|   +-- drivers/
|   +-- inventory/
|   +-- settings/             (restaurant, billing, roles, users, access-control, account, appearance)
|   +-- theme-builder/        (editor, presets, import, QR studio)
|   +-- pos/                  (Cashier POS - planned)
|   +-- notifications/
|   +-- rbac/
+-- components/
|   +-- layout/               (Sidebar, DashboardLayout, Header, LanguageSwitcher, TrialWarningBanner)
|   +-- products/
|   +-- billing/              (PlanCard, PaymentHistoryTable)
|   +-- reports/
|   +-- settings/             (RBACSettings)
|   +-- editors/              (ColorPicker, FontPicker, HeaderEditor, FooterEditor, ContrastChecker)
|   +-- hr/
|   +-- auth/, blocked/, common/, form/, navigation/, notifications/, providers/
|   +-- rbac/, skeletons/, theme/, ui/, virtualization/
+-- hooks/
|   +-- useRbac.ts            (permission checks)
|   +-- useKDS.ts             (KDS data)
|   +-- useKDSAnalytics.ts    (KDS performance metrics)
|   +-- useKDSAudio.ts        (Web Audio API with HTMLAudio fallback)
|   +-- useKDSOffline.ts      (IndexedDB offline queue)
|   +-- useKDSWebSocket.ts    (real-time KDS events)
|   +-- useTables.ts          (table management)
|   +-- useReports.ts         (report data)
|   +-- usePayment.ts         (Paymob payment flow)
|   +-- useTenantStatus.ts    (subscription/block status)
|   +-- useNotifications.ts
|   +-- useThemeStore.ts
|   +-- useCache.ts, useDebounce.ts, useFormValidation.ts
|   +-- usePreloadRoute.ts, useServiceWorker.ts, useWebVitals.ts
+-- stores/
|   +-- authStore.ts          (user + token + tenant context)
|   +-- kdsStore.ts           (KDS state with Zustand persist + IndexedDB)
|   +-- preferencesStore.ts   (user preferences: theme, language, colors)
|   +-- tenantStatusStore.ts  (trial/blocked state)
|   +-- themeBuilderStore.ts  (theme editor state)
+-- lib/
|   +-- api.ts                (Axios instance with auth/tenant interceptors)
|   +-- rbac.ts               (MODULE_IDS, PERMISSION_LEVELS constants)
|   +-- config.ts             (environment-aware config)
|   +-- translations.ts       (i18n helpers)
|   +-- api/themeApi.ts       (theme CRUD)
|   +-- api/reportsApi.ts     (report endpoints)
|   +-- api/tablesApi.ts      (table CRUD)
|   +-- api/posApi.ts         (POS API - planned)
|   +-- pdf-export.ts         (client-side PDF generation)
|   +-- code-splitting.ts     (lazy loading utils)
|   +-- performance.ts        (Web Vitals tracking)
+-- i18n/messages/
|   +-- en.json, ar.json
+-- types/
    +-- index.ts, reports.ts
```

### Restaurant Website Structure

```
frontend/apps/restaurant-website/
+-- app/[locale]/
|   +-- page.tsx              (dynamic homepage)
|   +-- menu/                 (product catalog)
|   +-- cart/                 (shopping cart)
|   +-- checkout/             (order submission + payment)
|   +-- orders/               (order history + tracking)
|   +-- orders/track/         (public order tracking)
|   +-- table/[tableId]/      (QR dine-in table ordering)
|   +-- settings/
+-- components/
|   +-- commerce/             (cart, order-list, order-status-tracker, product-cards)
|   +-- homepage/             (dynamic-home-page, featured-products-section, error-boundary)
|   +-- layout/               (header-ssr, footer-ssr, ThemeHeader, ThemeFooter)
|   +-- providers/            (theme-provider, query-provider, web-vitals-provider)
|   +-- system/               (MaintenancePage)
+-- lib/
|   +-- api/, config/, store/, hooks/
|   +-- table-ordering.ts     (QR table context)
|   +-- subdomain-context.tsx  (multi-tenant subdomain detection)
|   +-- image-optimization.ts, seo-structured-data.ts
+-- middleware.ts              (subdomain routing, locale detection)
```

### UI Themes Layer

```
frontend/layers/ui-themes/src/
+-- components/
|   +-- layout/               (Header.tsx, Footer.tsx)
|   +-- runtime/              (ThemeProvider.tsx, SectionRenderer.tsx)
+-- sections/
|   +-- marketing/            (HeroSection, CTASection, TestimonialsSection)
|   +-- commerce/             (ProductsSection)
|   +-- content/              (WhyChooseUsSection)
|   +-- forms/                (ContactSection)
|   +-- custom/               (CustomSection)
+-- registry/componentRegistry.ts
+-- context/MockDataContext.tsx
```

---

## Backend Structure

### Go Project Layout

```
backend/
+-- cmd/api/main.go           (route registration, middleware wiring, WebSocket hub)
+-- internal/
|   +-- config/               (runtime_env.go - environment detection)
|   +-- domain/               (all entity structs)
|   +-- handler/
|   |   +-- http/             (REST handlers for every module)
|   |   +-- ws/               (WebSocket: kds_hub.go, kds_ws_handler.go)
|   +-- middleware/            (auth, tenant, permission, CORS, rate-limit, cache, gzip, metrics)
|   +-- repository/           (38 repo files - PostgreSQL data access)
|   +-- usecase/              (business logic layer)
|   +-- observability/        (perf_metrics.go)
+-- migrations/               (98 migration files, numbered 1-135 + legacy)
+-- tests/load/               (k6 load tests: baseline, normal, peak)
```

### Domain Models (`internal/domain/`)

| File | Key Entities |
|------|-------------|
| `order.go` | `Order`, `OrderItem`, `OrderAddOn`, `CreateOrderRequest`, `OrderListFilters`, `OrderStats` |
| `role.go` | `SimpleRole`, `ModuleDefinition`, `PermissionLevel`, `RolePermission`, `UserRole` |
| `product.go` | `Product`, variants, addons |
| `kds.go` | KDS queue items, settings, analytics |
| `table.go` | Restaurant tables, layouts |
| `payment.go` | Payment transactions |
| `subscription.go` | SaaS trial/subscription lifecycle |
| `theme.go` | Theme builder domain (versioned) |
| `homepage_builder.go` | Homepage section builder |
| `notification.go` | Notification types, priorities |
| `hr.go` | HR employees |
| `driver.go` | Delivery drivers |
| `user.go`, `user_settings.go` | User + per-user settings |
| `restaurant.go`, `tenant.go` | Multi-tenancy |
| `super_admin.go` | Super admin portal |
| `permission_audit_log.go` | RBAC audit trail |

### Order Domain

**Order Sources:** `website`, `mobile_app`, `phone`, `in_store` (planned: `pos`)
**Order Statuses:** `pending` -> `confirmed` -> `preparing` -> `ready` -> `out_for_delivery` -> `delivered` / `cancelled`
**Payment Methods:** `cash`, `card`, `online`, `wallet`
**Payment Statuses:** `pending`, `paid`, `failed`, `refunded`

Orders support: optional `TableID` (dine-in), delivery address with geo-coordinates, item variants, add-ons, and notes. Status transitions enforced via `CanTransitionStatus()` state machine.

### HTTP Handlers (`internal/handler/http/`)

| Handler | Purpose |
|---------|---------|
| `auth_handler.go` | Register, Login, LoginConfirm, TokenExchange, ForgotPassword, CheckSubdomain |
| `product_handler.go` | Product CRUD, low-stock alerts |
| `category_handler.go` | Category CRUD |
| `public_menu_handler.go` | Public menu/products/categories/search (no auth) |
| `public_order_handler.go` | Public order create/validate/track/cancel |
| `admin_order_handler.go` | Admin order management, status/payment updates, export, driver assignment |
| `kds_handler.go` | KDS queue, bump/recall, stats, analytics, settings |
| `table_handler.go` | Table CRUD, layout management, status updates |
| `report_handler.go` | Overview/sales/profitability/inventory reports, PDF export |
| `hr_employee_handler.go` | Employee CRUD |
| `hr_attendance_handler.go` | Attendance, clock-in |
| `hr_leave_handler.go` | Leave requests, approve/reject |
| `hr_salary_handler.go` | Payroll, mark-paid |
| `rbac_role_handler.go` | RBAC role CRUD |
| `rbac_permission_handler.go` | Permission assign/revoke/check |
| `rbac_user_role_handler.go` | Role assignment to users |
| `admin_theme_handler.go` | Theme CRUD, activate, duplicate, history |
| `admin_theme_handler_v2.go` | Theme presets |
| `public_homepage_handler.go` | Public homepage data, theme, sections |
| `notification_handler.go` | Notification CRUD, mark read/unread |
| `payment_handler.go` | Paymob initiate payment, webhook, history |
| `payment_monitoring_handler.go` | Payment observability |
| `super_admin_handler.go` | Super admin login, tenant management |
| `impersonation_handler.go` | Tenant impersonation |
| `subscription_handler.go` | Subscription status |
| `user_settings_handler.go` | User settings, profile |
| `admin_driver_handler.go` | Driver CRUD, assignment |
| `mock_paymob_handler.go` | Mock payment (dev only) |

### WebSocket (`internal/handler/ws/`)

- `kds_hub.go` - Broadcasts kitchen events to connected KDS clients
- `kds_ws_handler.go` - WebSocket upgrade handler
- Wired via `orderUC.SetKDSBroadcaster(kdsHub)` so order creation/status changes auto-broadcast

### Middleware (`internal/middleware/`)

| Middleware | Purpose |
|-----------|---------|
| `auth.go` | JWT validation, sets user claims in context |
| `tenant.go` | Extracts tenant/restaurant context from JWT |
| `tenant_status.go` | Blocks suspended/expired tenants |
| `permission.go` | `WithRequiredPermission(moduleID, level)` - RBAC gate |
| `cors.go` | CORS handling |
| `rate_limit.go` | Rate limiting (by IP, by user, by user+window) |
| `cache.go` | Response caching |
| `gzip.go` | Compression |
| `impersonation.go` | Super admin impersonation |
| `super_admin_auth.go` | Super admin JWT verification |
| `request_metrics.go` | Request duration/count metrics |

---

## RBAC (Role-Based Access Control)

### Module IDs (1-14)

| ID | Module | Purpose | Typical Users |
|----|--------|---------|---------------|
| 1 | PRODUCTS | Menu management | Manager, Admin |
| 2 | CATEGORIES | Category organization | Manager, Admin |
| 3 | HR | Employee management | Admin, Manager |
| 4 | NOTIFICATIONS | System alerts | All |
| 5 | SETTINGS | Business configuration | Admin only |
| 6 | USERS | User account management | Admin only |
| 7 | ROLES | Role/permission management | Super Admin |
| 8 | REPORTS | Analytics & PDF exports | Admin, Manager |
| 9 | THEMES | Website theme builder | Admin, Designer |
| 10 | ORDERS | Order management | All roles |
| 11 | LEAVES | Employee leave tracking | Admin, Manager |
| 12 | SALARIES | Payroll management | Admin, Accountant |
| 13 | KDS | Kitchen Display System | Kitchen Staff |
| 14 | TABLES | Table management | Admin, Manager, Staff |

### Permission Levels (Hierarchical)

```
ADMIN (5)     - Full control, no restrictions
  |
DELETE (4)    - Create, Read, Update, Delete
  |
WRITE (3)     - Create, Read, Update
  |
EXPORT (2)    - Read + Export to PDF/CSV
  |
READ (1)      - Read-only access
  |
NONE (0)      - No access (hidden from UI)
```

### System Roles (Predefined)

| Role | ID | Key Permissions |
|------|-----|----------------|
| ADMIN | 4 | All 14 modules with ADMIN level |
| MANAGER | 9 | Products, Categories, Orders, HR, Leaves, Reports - WRITE |
| EMPLOYEE | 10 | Products, Categories, Orders, HR, Leaves - READ |
| DRIVER | 11 | Orders - READ |
| KITCHEN_STAFF | 13 | KDS - READ/WRITE, Orders - READ |

### RBAC Database Tables

```sql
modules         { id, name, description, created_at }
roles           { id, tenant_id, role_code, role_name, description, is_active, is_system_role }
role_permissions { id, role_id, module_id, permission_level }
user_roles       { id, tenant_id, user_id, role_id, assigned_by, created_at }
```

### Backend RBAC Flow

```
1. User Login -> JWT with user_id, tenant_id, restaurant_id, role
2. Request with Authorization header
3. AuthMiddleware validates JWT
4. TenantContextMiddleware sets tenant scope
5. WithRequiredPermission(moduleID, level) checks role_permissions
6. Handler executes if authorized
```

### Frontend RBAC Flow

```
1. useHasPermission(MODULE_IDS.PRODUCTS, PERMISSION_LEVELS.READ)
2. Hook calls: GET /api/v1/rbac/check-permission
3. Sidebar filters menu items based on permissions
4. Page components render permission gates
```

### Key RBAC Files

**Frontend:**
- `lib/rbac.ts` - `MODULE_IDS` (1-14), `PERMISSION_LEVELS`, `rbacApi`
- `hooks/useRbac.ts` - `useHasPermission()`, `useUserPermissions()`
- `components/layout/Sidebar.tsx` - Permission-filtered menu items

**Backend:**
- `middleware/permission.go` - `WithRequiredPermission(moduleID, level)`
- `handler/http/rbac_permission_handler.go` - Permission check endpoints
- `repository/role_permission_repo.go` - Permission queries

---

## Database Schema & Migrations

### Core Tables

```
tenants         - id, name, slug, status, subscription info
restaurants     - id, tenant_id, name, address, phone, settings_json
users           - id, tenant_id, restaurant_id, email, password_hash, name, role, status
```

### Key Indexes

```sql
CREATE INDEX idx_users_tenant_restaurant ON users(tenant_id, restaurant_id);
CREATE INDEX idx_orders_tenant_restaurant ON orders(tenant_id, restaurant_id);
CREATE INDEX idx_products_tenant_restaurant ON products(tenant_id, restaurant_id);
CREATE INDEX idx_users_email_tenant ON users(email, tenant_id);
CREATE INDEX idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX idx_user_roles_user ON user_roles(user_id, tenant_id);
```

### Migration History

98 migration files (numbered 1-135). Key milestones:

| Range | Feature |
|-------|---------|
| 001-010 | Core tables (tenants, users, restaurants, products, categories, orders) |
| 020-040 | HR module (employees, attendance, leaves, salaries) |
| 050-070 | Theme system, homepage builder |
| 080-100 | RBAC (roles, modules, permissions, user_roles) |
| 101-119 | Notifications, payments, subscriptions, audit logs |
| 120-123 | Reports module + export permissions |
| 124-129 | KDS (queue, settings, analytics, performance indexes) |
| 130-135 | Tables (restaurant_tables, table_id on orders, layouts, RBAC module) |

---

## API Reference

### Public Endpoints (No Auth)

```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/login/confirm
POST   /api/v1/auth/token-exchange
POST   /api/v1/auth/forgot-password
GET    /api/v1/auth/check-subdomain

GET    /api/v1/public/restaurants/{slug}/menu
GET    /api/v1/public/restaurants/{slug}/products
GET    /api/v1/public/restaurants/{slug}/categories
GET    /api/v1/public/restaurants/{slug}/search
GET    /api/v1/public/restaurants/{slug}/homepage
GET    /api/v1/public/restaurants/{slug}/theme
GET    /api/v1/public/restaurants/{slug}/settings
GET    /api/v1/public/table/{id}/info

POST   /api/v1/public/orders                    (tenant context required)
GET    /api/v1/public/orders/{orderNumber}/track

GET    /api/v1/subscription/plans
POST   /api/v1/webhook/paymob
GET    /api/v1/health
```

### Protected Endpoints (Auth + Tenant + RBAC)

```
# Products (module 1)
GET/POST     /api/v1/products
GET/PUT/DEL  /api/v1/products/{id}

# Categories (module 2)
GET/POST     /api/v1/categories
GET/PUT/DEL  /api/v1/categories/{id}

# Orders (module 10)
GET          /api/v1/admin/orders
GET          /api/v1/admin/orders/{id}
PUT          /api/v1/admin/orders/{id}/status
PUT          /api/v1/admin/orders/{id}/payment-status
PUT          /api/v1/admin/orders/{id}/assign-driver
GET          /api/v1/admin/orders/export

# KDS (module 13)
GET          /api/v1/kds/queue
PUT          /api/v1/kds/orders/{id}/bump
PUT          /api/v1/kds/orders/{id}/recall
GET          /api/v1/kds/stats
GET          /api/v1/kds/analytics
GET/PUT      /api/v1/kds/settings
WS           /api/v1/kds/ws

# Tables (module 14)
GET/POST     /api/v1/tables
GET/PUT/DEL  /api/v1/tables/{id}
PUT          /api/v1/tables/{id}/status
GET/POST     /api/v1/tables/layouts

# Reports (module 8)
GET          /api/v1/reports/overview
GET          /api/v1/reports/sales
GET          /api/v1/reports/profitability
GET          /api/v1/reports/inventory
GET          /api/v1/reports/export/pdf

# HR (module 3)
GET/POST     /api/v1/employees
GET/PUT/DEL  /api/v1/employees/{id}
POST         /api/v1/attendance/clock-in
GET/POST     /api/v1/leaves
PUT          /api/v1/leaves/{id}/approve
GET/POST     /api/v1/salaries
PUT          /api/v1/salaries/{id}/mark-paid

# RBAC
GET          /api/v1/rbac/check-permission
GET          /api/v1/rbac/me/permissions
GET/POST     /api/v1/rbac/roles
GET/PUT/DEL  /api/v1/rbac/roles/{id}
POST/DEL     /api/v1/rbac/users/{id}/roles

# Themes
GET/POST     /api/v1/themes
GET/PUT/DEL  /api/v1/themes/{id}
POST         /api/v1/themes/{id}/activate
POST         /api/v1/themes/{id}/duplicate
GET          /api/v1/themes/{id}/history

# Notifications
GET          /api/v1/notifications
PUT          /api/v1/notifications/{id}/read
GET          /api/v1/notifications/stats

# Drivers
GET/POST     /api/v1/drivers
GET/PUT/DEL  /api/v1/drivers/{id}

# Payment & Subscription
POST         /api/v1/payments/initiate
GET          /api/v1/payments/history
GET          /api/v1/subscription/status

# User Settings
GET/PUT      /api/v1/users/{id}/settings
```

### Super Admin Endpoints (Separate Auth)

```
POST   /api/v1/super-admin/login
GET    /api/v1/super-admin/tenants
POST   /api/v1/super-admin/tenants/{id}/impersonate
GET    /api/v1/super-admin/payments
```

### JWT Token Structure

```json
{
  "user_id": 12,
  "email": "user@example.com",
  "tenant_id": 8,
  "restaurant_id": 8,
  "role": "user",
  "exp": 1771990950,
  "iat": 1771904550
}
```

---

## Feature Modules

### Kitchen Display System (KDS)

Real-time kitchen order management with WebSocket broadcasting.

**Key files:**
- Backend: `kds_handler.go`, `kds_usecase.go`, `kds_repo.go`, `ws/kds_hub.go`
- Frontend: `hooks/useKDS.ts`, `useKDSWebSocket.ts`, `useKDSAudio.ts`, `useKDSOffline.ts`, `stores/kdsStore.ts`
- Dashboard route: `/dashboard/kds` (fullscreen, hides sidebar/header)

**Features:** Queue management, order bump/recall, analytics, audio alerts (Web Audio API with HTMLAudio fallback), offline mode (IndexedDB), configurable settings.

### Table Management

Dine-in table ordering with QR code generation and visual layout editor.

**Key files:**
- Backend: `table_handler.go`, `table_usecase.go`, `table_repo.go`
- Frontend: `hooks/useTables.ts`, `lib/api/tablesApi.ts`, `lib/table-qr.ts`
- Dashboard route: `/dashboard/tables`
- Restaurant website: `/table/[tableId]/` (QR scan landing)

**Features:** Table CRUD, status tracking (available/occupied/reserved), visual layout editor, QR code generation, dine-in order association.

### Reports

Business analytics with PDF export capability.

**Key files:**
- Backend: `report_handler.go` (with embedded report fonts)
- Frontend: `hooks/useReports.ts`, `lib/api/reportsApi.ts`, `types/reports.ts`
- Dashboard route: `/dashboard/reports` (overview, sales, profitability, inventory)

### Theme Builder

WYSIWYG website builder with section-based homepage construction.

**Key files:**
- Backend: `admin_theme_handler.go`, `public_homepage_handler.go`
- Frontend: `stores/themeBuilderStore.ts`, `lib/api/themeApi.ts`
- Shared layer: `layers/ui-themes/` (sections, registry, renderer)
- Dashboard route: `/dashboard/theme-builder/editor/[id]` (fullscreen editor)

**Features:** Section types (Hero, CTA, Testimonials, Products, WhyChooseUs, Contact, Custom), color palette editor, font picker, header/footer customization, contrast checker, version history, theme presets, QR studio.

### Payment Integration (Paymob)

**Key files:**
- Backend: `payment_handler.go`, `payment_usecase.go`, `mock_paymob_handler.go`
- Frontend: `hooks/usePayment.ts`, `components/billing/`

**Features:** Live + mock mode, webhook with retry queue (exponential backoff: 1m -> 5m -> 30m -> 2h -> 12h), subscription plans, payment history.

### Subscription & Billing

**Features:** Trial period, paid plans, automatic expiry checker (background worker), blocked overlay for expired tenants (bypasses for billing page), subscription status API.

### Cashier POS (Planned)

High-performance cashier terminal targeting 3 clicks = full order in ~3 seconds. Full specification at `docs-organized/16-CASHIER-POS/CASHIER-POS-COMPLETE-SPECIFICATION.md`. Implementation TODO at `docs-organized/16-CASHIER-POS/CASHIER-POS-IMPLEMENTATION-TODO.md`.

---

## Development Guidelines

### Coding Standards

**TypeScript:**
- Strict mode enabled
- Interfaces for all data structures
- No `any` types (use `unknown` then narrow)

**Go:**
- Standard Go conventions (CamelCase exported, snake_case DB)
- All errors explicitly handled (no panic in production)
- Parameterized SQL queries (always `$1`, `$2` placeholders)
- Repository pattern for all DB access

**React:**
- Functional components only
- Hooks for all stateful logic
- Memoization for expensive computations
- Stable unique keys on lists

### Naming Conventions

- Components: PascalCase (`ProductCard.tsx`)
- Hooks: camelCase with `use` prefix (`useProducts.ts`)
- Stores: camelCase with `Store` suffix (`kdsStore.ts`)
- API files: camelCase with `Api` suffix (`reportsApi.ts`)
- Constants: UPPER_SNAKE_CASE

### API Integration Pattern

Frontend uses a centralized Axios instance (`lib/api.ts`) that auto-injects:
- `Authorization: Bearer {token}` from authStore
- `X-Tenant-ID` from authStore
- `X-Restaurant-ID` from authStore

```typescript
import api from '@/lib/api'

// All calls automatically include auth + tenant headers
const products = await api.get('/api/v1/products')
const order = await api.post('/api/v1/public/orders', payload)
```

### Fullscreen Route Pattern

KDS and POS use a fullscreen pattern that hides the sidebar and header:

```typescript
// In DashboardLayout.tsx
const isKDSFullscreenRoute = pathname.includes('/dashboard/kds') && !pathname.includes('/dashboard/kds/settings')
// When true: hides sidebar, header, trial banner; uses full-width layout
```

### Zustand Store Pattern

Stores use persist middleware with custom storage for resilience:

```typescript
// Pattern from kdsStore.ts
const useStore = create<StoreType>()(
  persist(
    (set, get) => ({ /* state + actions */ }),
    {
      name: 'store-key',
      storage: createJSONStorage(() => customPersistStorage),
      partialize: (state) => ({ /* only persist these fields */ }),
    }
  )
)
```

---

## Adding New Features

### Feature Checklist

1. **Backend Domain** - `internal/domain/[feature].go` (include `tenant_id`, `restaurant_id`)
2. **Migration** - `migrations/NNN_create_[feature].sql` (indexes on tenant+restaurant)
3. **Repository** - `internal/repository/[feature]_repo.go` (all queries scoped by tenant)
4. **Use Case** - `internal/usecase/[feature]_usecase.go`
5. **Handler** - `internal/handler/http/[feature]_handler.go`
6. **Routes** - Register in `cmd/api/main.go` with `wrapProtected` or `wrapWithPermission`
7. **RBAC Migration** - Seed module + default role permissions (if new module)
8. **Frontend Types** - `types/[feature].ts`
9. **Frontend API** - `lib/api/[feature]Api.ts`
10. **Frontend Hook** - `hooks/use[Feature].ts`
11. **Frontend Components** - `components/[feature]/`
12. **Frontend Page** - `app/[locale]/dashboard/[feature]/page.tsx`
13. **Sidebar Entry** - Add to `components/layout/Sidebar.tsx` with RBAC filtering
14. **i18n** - Add keys to `en.json` and `ar.json`

### RBAC Integration for New Module

1. Add module ID to `lib/rbac.ts`:
```typescript
export const MODULE_IDS = { ..., MY_FEATURE: 15 }
```

2. Create migration to seed module and permissions:
```sql
INSERT INTO modules (id, name, description) VALUES (15, 'MY_FEATURE', 'Description');
INSERT INTO role_permissions (role_id, module_id, permission_level)
VALUES (4, 15, 'ADMIN'), (9, 15, 'WRITE'), (10, 15, 'READ');
```

3. Use `WithRequiredPermission(15, "READ")` in route registration
4. Use `useHasPermission(MODULE_IDS.MY_FEATURE, PERMISSION_LEVELS.READ)` in frontend

---

## Scalability & Performance

### Database
- Connection pooling: Max 200 connections (PostgreSQL config)
- Composite indexes on `(tenant_id, restaurant_id)` for all tenant-scoped tables
- Pagination: Default 20 items, max 100

### API
- Rate limiting: IP-based + user-based (Redis-backed)
- Response caching middleware
- gzip compression middleware
- Request metrics tracking

### Frontend
- Code splitting by route (Next.js automatic)
- Lazy loading with `code-splitting.ts` utilities
- Zustand selectors for minimal re-renders
- Web Vitals monitoring
- Service Worker support

---

## Security Best Practices

- JWT tokens with 24h expiry
- Separate JWT secrets for regular users vs super admin
- Tenant isolation at database level (every query filtered)
- bcrypt password hashing
- HTTPS only in production
- CORS: specific origins only (configured per environment)
- Parameterized SQL queries (no injection)
- Rate limiting per IP and per user
- Permission audit logging
- Webhook signature verification (Paymob)
- Tenant status enforcement (blocks suspended tenants)
- CI guardrails: no hardcoded localhost URLs in production code

---

## Testing Strategy

### Backend Tests
- Unit tests: Repository, use case, domain logic
- Integration tests: Handler tests with database
- Load tests: k6 scripts (`tests/load/baseline.js`, `normal-load.js`, `peak-load.js`)

### Frontend Tests
- Unit: Vitest configs per feature (`vitest.config.unit.mts`, `vitest.config.reports.mts`)
- Component: React Testing Library
- Config tests: `config.test.ts` in each app
- API tests: `reportsApi.test.ts`, `table-qr.test.ts`, `website-qr.test.ts`
- Store tests: `kdsStore.test.ts`

### CI/CD
- `runtime-localhost-guardrail.yml` - Ensures no hardcoded localhost URLs
- `url-domain-config-smoke.yml` - Domain config smoke tests

---

## Deployment & Infrastructure

### Docker Compose Services

| Service | Image | Port | Purpose |
|---------|-------|------|---------|
| postgres | postgres:15-alpine | 5432 | Database (max 200 connections) |
| backend | pos-backend:latest | 8080 | Go API server |

Frontend apps run separately (not in docker-compose).

### Environment Variables

```bash
# Backend
DATABASE_URL=postgresql://user:pass@localhost:5432/pos_saas
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret
SUPER_ADMIN_JWT_SECRET=your-super-admin-secret
CORS_ORIGINS=http://localhost:3001,http://localhost:3002,http://localhost:3003,http://localhost:3004
PAYMOB_MODE=mock  # or 'live'
LOG_LEVEL=info

# Frontend (each app)
NEXT_PUBLIC_API_URL=https://api.example.com
```

### Database Migrations

```bash
# Run migrations (auto-applied via init scripts in Docker)
./migrate -path ./migrations -database "$DATABASE_URL" up
```

### Stub Mode

Backend can run with `db = nil` for development without PostgreSQL - handlers return mocked JSON responses.

---

## Common Issues & Solutions

### User has no sidebar items with ADMIN role
- JWT not refreshed after role assignment -> logout and login again
- Role permissions not seeded -> check `role_permissions` table

### Permission check returning 400
- User missing `restaurant_id` -> verify `users` table
- `TenantContextMiddleware` failing -> check JWT claims

### Cross-tenant data leakage
- Repository query missing `tenant_id` filter -> always filter by `tenant_id AND restaurant_id`

### KDS not receiving real-time updates
- WebSocket connection not established -> check `/api/v1/kds/ws` endpoint
- `kdsHub` not wired to order use case -> verify `orderUC.SetKDSBroadcaster(kdsHub)` in main.go

### Restaurant website showing wrong tenant data
- Subdomain routing misconfigured -> check `middleware.ts` subdomain detection
- Missing `X-Tenant-ID` header -> verify API client configuration

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 3.0.0 | 2026-02-25 | KDS, Tables, Reports, Paymob, Theme Builder, full 14-module RBAC, 4 frontend apps |
| 2.0.0 | 2026-02-24 | RBAC system complete, frontend filtering, 1M user readiness |
| 1.5.0 | 2026-02-20 | Multi-tenant isolation, backend RBAC endpoints |
| 1.0.0 | 2026-01-01 | Initial project setup |

---

**Last Updated:** 2026-02-25
**Maintained By:** Development Team
