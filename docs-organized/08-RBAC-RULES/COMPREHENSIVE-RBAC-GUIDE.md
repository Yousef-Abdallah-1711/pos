# Comprehensive RBAC (Role-Based Access Control) Guide

**Version**: 1.0
**Last Updated**: February 24, 2026
**Status**: Complete

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Module & Permission Matrix](#module--permission-matrix)
3. [Permission Levels](#permission-levels)
4. [Default Role Templates](#default-role-templates)
5. [JWT Claims & Permission Checking](#jwt-claims--permission-checking)
6. [Backend Permission Enforcement](#backend-permission-enforcement)
7. [Frontend Permission Checking](#frontend-permission-checking)
8. [Feature-Specific RBAC Rules](#feature-specific-rbac-rules)
9. [Best Practices](#best-practices)
10. [Creating Custom Roles](#creating-custom-roles)
11. [Common Role Combinations](#common-role-combinations)
12. [Troubleshooting](#troubleshooting)

---

## Executive Summary

The POS SaaS platform uses a **Module-Based RBAC System** to control user access to features. Users are assigned **Roles**, and each role has specific **Permissions** on different **Modules**. Permissions are hierarchical (READ < EXPORT < WRITE < DELETE < ADMIN).

### Key Principles

✅ **Principle of Least Privilege**: Users only get access to what they need
✅ **Role-Based**: Assign roles to users, not individual permissions
✅ **Granular Control**: Admin can create custom roles with specific permissions
✅ **Frontend + Backend**: Permissions enforced on both UI (hide features) and API (return 403)
✅ **JWT-Based**: Permissions encoded in JWT token, no additional API calls needed

### Quick Answer: "Where do users get access?"

**Admin Tools**: Settings > Administration > Roles (create/edit roles and assign permissions)
**Sidebar Filtering**: Sidebar shows only items user has permission for
**Dashboard Widgets**: Widgets hidden if user lacks READ permission
**Direct URL Access**: Accessing `/dashboard/products` without PRODUCTS READ permission shows "Access Denied"

---

## Module & Permission Matrix

### Complete Module List

| ID | Module Name | Purpose | Sidebar Item | Example Widget | Min Permission to See | Default Roles |
|----|------------|---------|-------------|-----------------|----------------------|--------------|
| 1 | PRODUCTS | Menu items, inventory | Products | Products Count | READ | Manager, Admin |
| 2 | CATEGORIES | Menu sections | Categories | Categories Count | READ | Manager, Admin |
| 3 | HR | Employees, attendance | HR | Staff Count | READ | Manager, Admin |
| 4 | NOTIFICATIONS | System alerts | (Hidden) | Unread Notifications | READ | All Users |
| 5 | SETTINGS | System configuration | Settings | - | - | Owner, Admin |
| 6 | USERS | User management | Settings > Admin | - | WRITE | Owner, Admin |
| 7 | ROLES | Role management | Settings > Admin | - | WRITE | Owner, Admin |
| 8 | REPORTS | Analytics & revenue | Reports | Recent Paid Revenue | READ | Manager, Viewer |
| 9 | THEMES | Storefront customization | Theme Builder | - | READ | Manager, Admin |
| 10 | ORDERS | Customer orders, delivery | Orders | Orders Count, Drivers | READ | All Roles except Viewer |
| 11 | LEAVES | Employee leave requests | HR > Leaves | - | READ | Manager, Admin |
| 12 | SALARIES | Payroll management | HR > Payroll | - | READ | Owner, Admin |
| 13 | KDS | Kitchen display system | KDS | - | READ | Kitchen Staff |
| 14 | TABLES | Table management | Tables | - | READ | Manager |

### Permission Matrix Legend

```
✓ = Feature visible/available
✗ = Feature hidden/unavailable
- = Not applicable
```

### Role × Module Permission Matrix

| Role / Module | 1:PROD | 2:CAT | 3:HR | 4:NOTIF | 5:SET | 6:USER | 7:ROLE | 8:REP | 9:THEME | 10:ORD | 11:LEAVE | 12:SAL | 13:KDS | 14:TABLE |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| **Admin** | ADMIN | ADMIN | ADMIN | READ | ADMIN | ADMIN | ADMIN | ADMIN | ADMIN | ADMIN | ADMIN | ADMIN | ADMIN | ADMIN |
| **Owner** | WRITE | WRITE | WRITE | READ | ADMIN | WRITE | WRITE | WRITE | WRITE | WRITE | READ | WRITE | - | WRITE |
| **Manager** | WRITE | WRITE | READ | READ | READ | - | - | READ | WRITE | WRITE | READ | READ | - | WRITE |
| **Kitchen Staff** | - | - | - | READ | - | - | - | - | - | READ | - | - | WRITE | - |
| **Viewer** | - | - | - | READ | - | - | - | READ | - | - | - | - | - | - |

---

## Permission Levels

### Hierarchical Permission Model

Permissions are **hierarchical** - higher permission levels include lower ones:

```
NONE (0) < READ (1) < EXPORT (2) < WRITE (3) < DELETE (4) < ADMIN (5)
```

### Permission Meanings by Module

#### PRODUCTS Module
- **READ**: View product list, see inventory
- **EXPORT**: Export product data to CSV/Excel
- **WRITE**: Create, edit, update products
- **DELETE**: Delete products
- **ADMIN**: Manage product categories, pricing rules, variants

#### ORDERS Module
- **READ**: View orders, see customer details, track delivery
- **EXPORT**: Export orders to CSV, generate reports
- **WRITE**: Create orders, update order status, assign drivers
- **DELETE**: Cancel orders
- **ADMIN**: Access advanced order settings, refunds, adjustments

#### REPORTS Module
- **READ**: View reports, see analytics charts
- **EXPORT**: Download report data as PDF/Excel
- **WRITE**: Create custom reports, save report templates
- **DELETE**: Delete saved reports
- **ADMIN**: Access reporting settings, data retention policies

#### KDS Module (Kitchen Display System)
- **READ**: View order queue on kitchen display
- **WRITE**: Bump orders through preparation stages
- **DELETE**: Recall/cancel orders from kitchen
- **ADMIN**: Configure KDS settings, manage multiple stations

#### TABLES Module
- **READ**: View tables, see floor plan
- **WRITE**: Create/edit tables, modify floor plan, assign orders to tables
- **DELETE**: Delete tables, clear table assignments
- **ADMIN**: QR code generation, billing settings, advanced floor plan

#### Settings/RBAC Modules
- **READ**: View current settings
- **WRITE**: Modify settings, create roles, assign permissions
- **ADMIN**: Access all settings, system-wide configurations

---

## Default Role Templates

### Admin Role

**Access Level**: Full system access

```
All modules: ADMIN permission
```

**Use Case**: System administrators, super users
**Dashboard View**: All widgets visible (Revenue, Orders, Products, Staff, etc.)
**Can Do**: Everything including user management, role creation, system settings

**Permissions**:
- Create/edit/delete any user
- Create/edit/delete roles
- Access all modules
- Configure system settings
- View analytics

---

### Owner Role

**Access Level**: Business owner - most operational features

```
PRODUCTS: WRITE
CATEGORIES: WRITE
ORDERS: WRITE
HR: WRITE
REPORTS: WRITE
THEMES: WRITE
TABLES: WRITE
KDS: -
SETTINGS: ADMIN
USERS: WRITE
ROLES: WRITE
```

**Use Case**: Restaurant owner, franchise owner
**Dashboard View**: All operational widgets (Orders, Products, Revenue, Staff)
**Can Do**: Manage entire restaurant operations, but NOT kitchen display

**Permissions**:
- View and create orders
- Edit menu items and categories
- Manage staff and payroll
- Create roles for other staff members
- View and customize reports
- Access theme builder

---

### Manager Role

**Access Level**: Operations manager - core features with some restrictions

```
PRODUCTS: WRITE
CATEGORIES: WRITE
ORDERS: WRITE
HR: READ
REPORTS: READ
THEMES: WRITE
TABLES: WRITE
KDS: -
SETTINGS: READ
```

**Use Case**: Restaurant manager, shift supervisor
**Dashboard View**: Orders, Products, Categories, Staff count (read-only), Revenue (read-only)
**Can Do**: Day-to-day operations, but cannot manage staff or create new roles

**Permissions**:
- Create/edit orders and assign drivers
- Edit menu items
- View staff list (cannot add/edit/delete)
- View analytics reports (cannot create custom ones)
- Customize storefront theme
- Manage tables and floor plan

---

### Kitchen Staff Role

**Access Level**: Very limited - KDS only

```
KDS: WRITE
ORDERS: READ
```

**Use Case**: Kitchen prep staff, expediter
**Dashboard View**: Order count (read-only), Notifications only
**Can Do**: View incoming orders on KDS, bump orders through preparation stages

**Permissions**:
- See incoming orders in KDS queue
- Bump orders from one stage to next
- View order details (ingredients, special instructions)
- Cannot modify order status or cancel orders

---

### Viewer Role

**Access Level**: Read-only - analytics only

```
REPORTS: READ
```

**Use Case**: Investor, consultant, analytics team
**Dashboard View**: Recent Paid Revenue widget only
**Can Do**: View reports and analytics, nothing else

**Permissions**:
- View all reports and analytics
- See revenue metrics
- Cannot make any changes
- Cannot access operational features

---

## JWT Claims & Permission Checking

### JWT Token Structure

When a user logs in, the backend generates a JWT token with embedded permissions:

```json
{
  "user_id": 123,
  "email": "manager@restaurant.com",
  "tenant_id": 1,
  "restaurant_id": 5,
  "role": "manager",
  "permissions": {
    "1": "WRITE",     // PRODUCTS: WRITE
    "2": "WRITE",     // CATEGORIES: WRITE
    "3": "READ",      // HR: READ
    "8": "READ",      // REPORTS: READ
    "9": "WRITE",     // THEMES: WRITE
    "10": "WRITE"     // ORDERS: WRITE
  },
  "iat": 1708760000,
  "exp": 1708846400
}
```

### Token Encoding Rules

- Stored in `authStore.token` on frontend
- Sent in `Authorization: Bearer <token>` header on API calls
- Decoded by middleware on backend
- Frontend extracts claims without needing API calls

### Hierarchical Permission Checking

The system uses integer comparison:

```
READ = 1
EXPORT = 2
WRITE = 3
DELETE = 4
ADMIN = 5

User has WRITE (3) on PRODUCTS
User requests to DELETE (4)
3 >= 4? NO → Access Denied
```

---

## Backend Permission Enforcement

### Middleware Pattern

All protected routes use the `wrapWithPermission` middleware:

```go
// backend/cmd/api/main.go

// Route: GET /api/v1/products (requires READ permission on PRODUCTS)
router.GET("/products", wrapWithPermission(productHandler.GetProducts, 1, "READ"))

// Route: POST /api/v1/products (requires WRITE permission on PRODUCTS)
router.POST("/products", wrapWithPermission(productHandler.CreateProduct, 1, "WRITE"))

// Route: DELETE /api/v1/products/:id (requires DELETE permission on PRODUCTS)
router.DELETE("/products/:id", wrapWithPermission(productHandler.DeleteProduct, 1, "DELETE"))
```

### How It Works

1. **Extract JWT**: Middleware reads `Authorization: Bearer <token>` header
2. **Parse Claims**: Decodes JWT and extracts user permissions object
3. **Check Permission**: Compares requested permission level with user's permission
4. **Enforce**: Returns 403 Forbidden if user lacks permission
5. **Allow**: Passes to handler if permission check passes

### Example: Order Update Endpoint

```
User: manager@restaurant.com
Role: Manager
Permissions: { "10": "WRITE" }  // ORDERS: WRITE

Request: PUT /api/v1/orders/123
Module: 10 (ORDERS)
Required: WRITE (3)

Check: permissions["10"] >= 3?
Result: WRITE (3) >= 3? YES → ✅ Allow

Response: 200 OK - Order updated
```

### Example: Unauthorized Access

```
User: kitchen@restaurant.com
Role: Kitchen Staff
Permissions: { "13": "WRITE", "10": "READ" }  // KDS: WRITE, ORDERS: READ

Request: DELETE /api/v1/products/123
Module: 1 (PRODUCTS)
Required: DELETE (4)

Check: permissions["1"] >= 4?
Result: No permission for module 1 → ❌ Deny

Response: 403 Forbidden - Access Denied
```

---

## Frontend Permission Checking

### useHasPermission Hook

**Location**: `frontend/apps/dashboard/src/hooks/useRbac.ts`

The `useHasPermission` hook allows components to check permissions:

```typescript
import { useHasPermission } from '@/hooks/useRbac'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'

export function ProductsPage() {
  const { hasPermission: canViewProducts } = useHasPermission(
    MODULE_IDS.PRODUCTS,
    PERMISSION_LEVELS.READ
  )

  if (!canViewProducts) {
    return <AccessDenied reason="You don't have permission to view products" />
  }

  return <ProductsList />
}
```

### MODULE_IDS Constants

**Location**: `frontend/apps/dashboard/src/lib/rbac.ts`

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
}
```

### PERMISSION_LEVELS Constants

```typescript
export const PERMISSION_LEVELS = {
  NONE: 0,
  READ: 1,
  EXPORT: 2,
  WRITE: 3,
  DELETE: 4,
  ADMIN: 5,
}
```

### Examples

#### Check for READ permission

```typescript
const { hasPermission: canViewOrders } = useHasPermission(
  MODULE_IDS.ORDERS,
  PERMISSION_LEVELS.READ
)

// Usage:
{canViewOrders && <OrdersList />}
```

#### Check for WRITE permission (edit button)

```typescript
const { hasPermission: canEditOrders } = useHasPermission(
  MODULE_IDS.ORDERS,
  PERMISSION_LEVELS.WRITE
)

// Usage:
<button disabled={!canEditOrders}>Edit Order</button>
```

#### Check for ADMIN permission (settings)

```typescript
const { hasPermission: isKdsAdmin } = useHasPermission(
  MODULE_IDS.KDS,
  PERMISSION_LEVELS.ADMIN
)

// Usage:
{isKdsAdmin && <KdsSettingsButton />}
```

### Sidebar Filtering (Already Implemented)

**File**: `frontend/apps/dashboard/src/components/layout/Sidebar.tsx`

Sidebar items are automatically hidden if user lacks permission:

```typescript
const sidebarItems = useMemo(() => {
  const permissionMap = {
    [MODULE_IDS.REPORTS]: canViewReports,
    [MODULE_IDS.ORDERS]: canViewOrders,
    [MODULE_IDS.PRODUCTS]: canViewProducts,
    [MODULE_IDS.HR]: canViewHR,
    [MODULE_IDS.THEMES]: canViewThemes,
  }

  return allSidebarItems.filter((item) => {
    if (!item.moduleId) return true  // Always show items without module requirement
    return permissionMap[item.moduleId] ?? false
  })
}, [canViewReports, canViewOrders, ...])
```

### Dashboard Widget Filtering (Already Implemented)

**File**: `frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx`

Dashboard only shows widgets the user has permission for:

```typescript
const metricCards = useMemo(() => {
  const allCards = [/* all metrics */]

  return allCards.filter((card) => {
    if (card.key === 'recentRevenue' && !canViewReports) return false
    if (card.key === 'orders' && !canViewOrders) return false
    if (card.key === 'products' && !canViewProducts) return false
    // ... more filters
    return true
  })
}, [canViewReports, canViewOrders, ...])
```

---

## Feature-Specific RBAC Rules

### Products & Categories

**Sidebar Item**: "Products" menu (under Products section)
**Module ID**: 1 (PRODUCTS), 2 (CATEGORIES)
**Required Permission**: READ

**Visibility**:
- Sidebar hidden if no PRODUCTS READ permission
- Products page shows "Access Denied" if no READ permission
- Create/Edit buttons hidden if no WRITE permission
- Delete buttons hidden if no DELETE permission

**Widget**: "Products Count" metric on dashboard
- Hidden if user lacks PRODUCTS READ permission

### Orders

**Sidebar Item**: "Orders" menu (under Orders section)
**Module ID**: 10 (ORDERS)
**Required Permission**: READ

**Visibility**:
- Sidebar hidden if no ORDERS READ permission
- Orders page shows list if READ, shows "Access Denied" if not
- Create order button hidden if no WRITE permission
- Edit/cancel buttons hidden if no DELETE or WRITE permission
- Assign driver button hidden if no WRITE permission

**Widgets**: "Orders Count", "Pending Orders", "Drivers Count", "Recent Orders"
- All hidden if user lacks ORDERS READ permission

### HR & Employees

**Sidebar Item**: "HR" menu → "Employees" (under HR section)
**Module ID**: 3 (HR)
**Required Permission**: READ

**Visibility**:
- HR sidebar collapsed if no HR READ permission
- Employees page shows "Access Denied" if no READ permission
- Add/edit employee buttons hidden if no WRITE permission
- Delete buttons hidden if no DELETE permission

**Widget**: "Staff Count" metric on dashboard
- Hidden if user lacks HR READ permission

### Reports & Analytics

**Sidebar Item**: "Reports" menu (in main sidebar)
**Module ID**: 8 (REPORTS)
**Required Permission**: READ

**Visibility**:
- Sidebar hidden if no REPORTS READ permission
- Reports page shows "Access Denied" if no READ permission
- Export button hidden if no EXPORT permission

**Widget**: "Recent Paid Revenue" metric on dashboard
- Hidden if user lacks REPORTS READ permission

### Theme Builder

**Sidebar Item**: "Theme Builder" menu (in main sidebar)
**Module ID**: 9 (THEMES)
**Required Permission**: READ

**Visibility**:
- Sidebar hidden if no THEMES READ permission
- Theme builder shows "Access Denied" if no READ permission
- Save/publish buttons hidden if no WRITE permission

### KDS (Kitchen Display System)

**Sidebar Item**: "KDS" menu (if permission granted)
**Module ID**: 13 (KDS)
**Required Permission**: READ

**Visibility**:
- Sidebar hidden if no KDS READ permission
- KDS page shows "Access Denied" if no READ permission
- Bump button visible if WRITE permission
- Recall button visible if DELETE permission
- Settings button visible if ADMIN permission

**Use Case**: Kitchen staff can view orders and bump them through prep stages

### Tables

**Sidebar Item**: "Tables" menu (if permission granted)
**Module ID**: 14 (TABLES)
**Required Permission**: READ

**Visibility**:
- Sidebar hidden if no TABLES READ permission
- Tables page shows "Access Denied" if no READ permission
- Add/edit buttons hidden if no WRITE permission
- Delete buttons hidden if no DELETE permission
- Floor plan editor hidden if no WRITE permission

**Use Case**: Restaurant managers can manage physical tables and floor plan

### Settings

**Sidebar Item**: Always visible, shows "Settings"
**Module ID**: 5 (SETTINGS) and submodules

**Visibility**:
- All users can access Settings
- Account/billing sections visible to all
- Restaurant settings hidden if no SETTINGS WRITE permission
- Role management visible only to ROLES WRITE users
- User management visible only to USERS WRITE users

---

## Best Practices

### 1. Principle of Least Privilege

✅ **DO**: Create roles with minimal permissions
```
Kitchen Staff: Only KDS WRITE + ORDERS READ
```

❌ **DON'T**: Give admin access to junior staff
```
Kitchen Staff: ADMIN on all modules (too much power)
```

### 2. Role Clarity

✅ **DO**: Name roles by job function
```
"Kitchen Lead", "Cashier", "Delivery Manager"
```

❌ **DON'T**: Create overlapping roles
```
"Manager1", "Manager2", "Supervisor1"
```

### 3. Permission Hierarchy

✅ **DO**: Use hierarchical permissions
```
Manager has WRITE (which includes READ)
Kitchen Staff has READ (view only)
```

❌ **DON'T**: Grant multiple conflicting permissions
```
User has both WRITE and DELETE separately (confusing)
```

### 4. Frontend + Backend Enforcement

✅ **DO**: Check permissions in both places
```
Frontend: Hide button if no permission
Backend: Return 403 if no permission
```

❌ **DON'T**: Rely only on frontend hiding
```
Frontend: Hide button, but API call still works (security issue)
```

### 5. User Audit Trail

✅ **DO**: Log role changes
```
"Admin changed Kitchen Staff role to add WRITE permission - Feb 24 2:30 PM"
```

❌ **DON'T**: Change permissions without tracking
```
(Lose track of who has what access)
```

---

## Creating Custom Roles

### Step 1: Determine Job Function

Ask: "What job does this person do?"

Example: "Order fulfillment supervisor"

### Step 2: Identify Required Modules

Map job function to modules:

```
Order fulfillment supervisor needs:
✓ ORDERS (view and update orders)
✓ DRIVERS (assign delivery)
✗ PRODUCTS (don't need to edit menu)
✗ REPORTS (don't need analytics)
```

### Step 3: Determine Permission Levels

For each module, choose the minimum permission needed:

```
ORDERS: WRITE (create orders, update status, assign drivers)
DRIVERS: WRITE (manage delivery assignments)
KDS: READ (view incoming orders from kitchen)
```

### Step 4: Create Role via UI

**Path**: Dashboard > Settings > Administration > Roles

1. Click "Create New Role"
2. Enter role name: "Order Fulfillment Supervisor"
3. Enter description: "Manages order fulfillment and delivery"
4. Set module permissions:
   - ORDERS: WRITE
   - KDS: READ
5. Save
6. Assign to users via Settings > User Assignments

### Step 5: Test

1. Create test user with new role
2. Test all expected features work
3. Test restricted features are hidden
4. Check for API "403 Forbidden" errors on restricted operations

---

## Common Role Combinations

### Fast Food Franchise

**Owner**: WRITE on Products, Orders, HR, Reports, Themes, Tables, Settings ADMIN
**Shift Manager**: WRITE on Orders, READ on Products/Reports/HR, WRITE on Tables
**Kitchen Lead**: WRITE on KDS, READ on Orders
**Cashier**: WRITE on Orders only
**Delivery Coordinator**: WRITE on Orders + Drivers

### Fine Dining Restaurant

**Owner**: Full ADMIN access
**General Manager**: WRITE on all operational modules, READ on Reports, ADMIN on KDS
**Chef**: READ on Orders, WRITE on KDS, READ on Inventory (Products)
**Front of House Manager**: WRITE on Orders, WRITE on Tables, READ on Reports
**Sommelier**: READ on Products, can create wine product variants
**Accountant**: READ on Reports + Invoices

### Cloud-Based SaaS (Multi-tenant)

**Super Admin**: ADMIN on all modules for all restaurants
**Restaurant Owner**: WRITE on their restaurant, ADMIN on ROLES/USERS for their staff
**Multi-Restaurant Manager**: WRITE on Products/Orders/HR for assigned restaurants
**Finance Team**: READ on Reports across all restaurants
**Support Agent**: READ-only access to specific restaurants for troubleshooting

---

## Troubleshooting

### Issue: User can't see menu item

**Symptoms**: Sidebar doesn't show "Products" menu for user

**Root Causes**:
1. User doesn't have PRODUCTS READ permission
2. Role not assigned to user
3. Role doesn't have permission configured

**Solution**:
1. Check user's role: Settings > User Assignments
2. Check role's permissions: Settings > Roles > (role name)
3. Ensure MODULE_IDS.PRODUCTS = 1 in codebase
4. Check useHasPermission hook is called correctly in Sidebar.tsx

### Issue: User can click button but gets "Access Denied"

**Symptoms**: Button is visible, but clicking it returns error

**Root Causes**:
1. Frontend missing permission check (button shown but shouldn't be)
2. Backend permission check stricter than frontend check
3. Token expired mid-session

**Solution**:
1. Add permission check to button:
   ```typescript
   {hasPermission && <button>Edit</button>}
   ```
2. Verify `wrapWithPermission` middleware uses same MODULE_IDS
3. Check token in browser DevTools > Application > Cookies
4. Ask user to log out and back in

### Issue: "Access Denied" page appears but user should have access

**Symptoms**: User claims they can't access feature they should be able to

**Root Causes**:
1. Wrong permission level specified
2. Module ID mismatch
3. Permission hierarchy misunderstood

**Solution**:
1. Check permission: `const { hasPermission } = useHasPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.WRITE)`
   - Does the module ID match? (ORDERS = 10?)
   - Does the permission level match? (WRITE = 3?)
2. Verify user's role has the permission:
   ```
   Role "Manager" → ORDERS: WRITE ✓
   ```
3. Check permission hierarchy: If WRITE (3) is needed, ADMIN (5) also works

### Issue: API returns 403 Forbidden

**Symptoms**: Browser network tab shows 403 on API request

**Root Causes**:
1. Backend `wrapWithPermission` module ID incorrect
2. JWT token claims missing permission
3. Permission level insufficient

**Solution**:
1. Check backend route:
   ```go
   router.POST("/products", wrapWithPermission(handler, 1, "WRITE"))
   //                                                   ^ correct module?
   ```
2. Decode JWT: Copy token to jwt.io, check `permissions` object
3. Check permission level: Is the user trying to DELETE when they only have READ?

### Issue: Permission changes don't take effect immediately

**Symptoms**: User role changed but user still sees old features

**Root Causes**:
1. User session still has old token
2. Browser cache showing old page
3. Token not refreshed

**Solution**:
1. Ask user to log out completely
2. Clear browser cache (Ctrl+Shift+Delete)
3. Log back in - should get new token with updated permissions
4. Hard refresh page (Ctrl+Shift+R)

### Issue: New module added but permissions not working

**Symptoms**: Created new module 15 (LOYALTY), but permission checks fail

**Root Causes**:
1. MODULE_IDS not updated in frontend
2. Migration not seeded in backend
3. Module not added to role seeding

**Solution**:
1. Add to `frontend/apps/dashboard/src/lib/rbac.ts`:
   ```typescript
   LOYALTY: 15,
   ```
2. Run migration to seed module in backend:
   ```sql
   INSERT INTO rbac_modules (id, name, description)
   VALUES (15, 'LOYALTY', 'Loyalty program management')
   ```
3. Ensure role seeding includes the new module
4. Clear browser cache and log in again

---

## Summary

The RBAC system provides:

✅ **Granular Control**: Permissions down to module level with 5 permission levels
✅ **Easy Management**: Admin interface to create roles and assign permissions
✅ **Security**: Enforced on both frontend (UX) and backend (API)
✅ **Scalability**: Support for custom roles and multi-restaurant access
✅ **Auditability**: Track which user has which permissions

**Remember**: When unsure about access, check:
1. User's assigned role
2. Role's module permissions
3. Module ID and permission level match between frontend and backend
4. Token hasn't expired

---

## Related Documents

- **[RBAC-IMPLEMENTATION-STATUS.md](../RBAC-RULES/RBAC-IMPLEMENTATION-STATUS.md)** - Implementation progress and current state
- **[KDS-IMPLEMENTATION-PLAN.md](../../12-KDS-KITCHEN-DISPLAY-SYSTEM/KDS-IMPLEMENTATION-PLAN.md)** - KDS-specific RBAC rules
- **[TABLE-MANAGEMENT-IMPLEMENTATION-PLAN.md](../../13-TABLE-MANAGEMENT/TABLE-MANAGEMENT-IMPLEMENTATION-PLAN.md)** - Table Management-specific RBAC rules
- **[Backend Routes](../../../backend/cmd/api/main.go)** - See `wrapWithPermission` usage
- **[Frontend Hooks](../../../frontend/apps/dashboard/src/hooks/useRbac.ts)** - See `useHasPermission` hook

---

**Questions?** Check the [Troubleshooting](#troubleshooting) section or review feature-specific guides above.
