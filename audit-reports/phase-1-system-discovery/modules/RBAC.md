# Module: RBAC (Phase 1 audit)

Read-only analysis. See `09-AUTHENTICATION-AND-RBAC.md` for architecture. Key fact: the backend does not enforce any module permission (`backend/internal/middleware/permission.go:15-22` only sets context values that nothing reads; checker commented out at `backend/cmd/api/main.go:624`). Tables below describe intended policy.

## 1. Permission levels
| Level | Value | Notes |
|---|---|---|
| NONE | 0 | |
| READ | 1 | |
| EXPORT | 2 | added by migration 122, shifts WRITE/DELETE/ADMIN up one |
| WRITE | 3 | |
| DELETE | 4 | |
| ADMIN | 5 | |

Source: `backend/internal/domain/role.go:118-126`, `migrations/101`, `migrations/122`. `role_permissions` stores the level as a string with a CHECK constraint; `permission_levels.value` rebased in 122. Frontend hierarchy matches (`frontend/apps/dashboard/src/lib/rbac.ts:42-49`).

## 2. Modules and route annotations
"Route levels" lists the levels used in `wrapWithPermission(...)` for that module in `main.go`. Modules with no route annotation have no backend route referencing them.

| ID | Module | Defined in | Route levels in main.go | Route evidence |
|---|---|---|---|---|
| 1 | PRODUCTS | 100 | READ, WRITE, DELETE | main.go:663-668 |
| 2 | CATEGORIES | 100 | (none; category routes use plain `wrapProtected`, main.go:671-675) but HR routes are annotated with 2 | see mismatch |
| 3 | HR | 100 | none (frontend uses 3) | rbac.ts:23 |
| 4 | NOTIFICATIONS | 100 | none (`wrapProtected`, main.go:1413-1458) | |
| 5 | SETTINGS | 100 | none | |
| 6 | USERS | 100 | none (users CRUD `wrapProtected`, main.go:1062-1289) | |
| 7 | ROLES | 100 | none (`wrapAuthOnly`, main.go:925-947) | |
| 8 | REPORTS | 100/120 | READ, EXPORT | main.go:1328-1352 (incl. HR reports) |
| 9 | THEMES | 100 | none (`wrapProtected`, main.go:1462-1472) | |
| 10 | ORDERS | 105 | READ, WRITE, DELETE | main.go:1294-1301,1404 |
| 11 | LEAVES | 105 | none (leave routes annotated 2) | main.go:961-967 |
| 12 | SALARIES | 105 | none (salary routes annotated 2) | main.go:971-976 |
| 13 | KDS | 126 | READ, WRITE, ADMIN | main.go:631-641,1319-1325 |
| 14 | TABLES | 132 | READ, WRITE, DELETE | main.go:1305-1314 |
| 15 | ACCOUNTING | 155 | READ, WRITE, DELETE, EXPORT | main.go:982-1051 |
| 16 | SUPPLIERS | 155 | READ, WRITE, DELETE | main.go:1053-1057 |
| 17 | FINANCIAL_REPORTS | 155 | READ, WRITE, DELETE, EXPORT | main.go:1333-1359 |
| 18 | PERFORMANCE | 155 | READ, WRITE, DELETE | main.go:1360-1367 |
| 19 | TRAINING | 155 | READ, WRITE, DELETE | main.go:1368-1373 |
| 20 | CONTRACTS | 155 | READ, WRITE, DELETE | main.go:1374-1394 (insurance, contracts, resignations) |

Module-ID mismatch: employees, attendance, leaves, salaries, HR roles routes carry module `2` with comment "HR" (main.go:757-976) while migration 100 insert order makes 2 = CATEGORIES and 3 = HR; leaves (11) and salaries (12) have dedicated modules that are never referenced. If enforcement is added without fixing this, HR would be gated by the CATEGORIES grant.

## 3. Default grants by role (module x level matrix, from seeds)

Seeds: `106_backfill_rbac_defaults.sql` (modules 1-12), `127_backfill_kds_role_permissions.sql` (13), `135_backfill_tables_role_permissions.sql` (14), `155_add_accounting_rbac_modules.sql` (15-20), `159_seed_external_accountant_role.sql`. Blank = NONE/no row.

| Module | ADMIN | MANAGER | EMPLOYEE | DRIVER | EXTERNAL_ACCOUNTANT | Kitchen roles |
|---|---|---|---|---|---|---|
| 1 PRODUCTS | ADMIN | WRITE | READ | NONE | - | |
| 2 CATEGORIES | ADMIN | WRITE | READ | NONE | - | |
| 3 HR | ADMIN | WRITE | READ | NONE | - | |
| 4 NOTIFICATIONS | ADMIN | WRITE | READ | READ | - | |
| 5 SETTINGS | ADMIN | READ | NONE | NONE | - | |
| 6 USERS | ADMIN | READ | NONE | NONE | - | |
| 7 ROLES | ADMIN | READ | NONE | NONE | - | |
| 8 REPORTS | ADMIN | READ | NONE | NONE | EXPORT | |
| 9 THEMES | ADMIN | READ | NONE | NONE | - | |
| 10 ORDERS | ADMIN | WRITE | READ | READ | - | |
| 11 LEAVES | ADMIN | WRITE | READ | NONE | - | |
| 12 SALARIES | ADMIN | READ | NONE | NONE | - | |
| 13 KDS | ADMIN | ADMIN | (EMPLOYEE WRITE via fallback `rbac_permission_handler.go:283-295`) | - | - | STAFF/LEAD/CHEF WRITE, VIEWER READ |
| 14 TABLES | ADMIN | WRITE | READ | - | - | READ |
| 15 ACCOUNTING | ADMIN | WRITE | - | - | READ | |
| 16 SUPPLIERS | ADMIN | WRITE | - | - | READ | |
| 17 FINANCIAL_REPORTS | ADMIN | READ | - | - | EXPORT | |
| 18 PERFORMANCE | ADMIN | WRITE | READ | - | READ | READ |
| 19 TRAINING | ADMIN | WRITE | READ | - | READ | READ |
| 20 CONTRACTS | ADMIN | READ | - | - | READ | |

Notes:
- 106 grants rows for every active module at the time it ran; modules created later rely on their own backfills (127, 135, 155). Backfill roles OWNER/ADMINISTRATOR are treated as ADMIN in 127/135/155 but not seeded as roles by migrations 021/106.
- Migration 021 seeds ADMIN/MANAGER/EMPLOYEE with `roles.permissions` array (a separate legacy mechanism); 106 adds DRIVER and `role_permissions`.
- Owner at registration gets ADMIN via `user_roles` backfill (106) for `users.role='owner'`; new users created via `POST /users` get EMPLOYEE (`main.go:1186-1191`).

## 4. Route coverage gaps (`backend/cmd/api/main.go`)

Routes with no permission wrapper at all (auth only via `wrapProtected`/`wrapAuthOnly`):
- Categories CRUD (671-675).
- User management: `GET/POST/PUT/DELETE /users`, `/users/{id}/settings*`, `/users/{id}/profile` (1062-1289). `POST /users` lets any tenant user create accounts; `DELETE /users/{id}` deletes any user in tenant.
- RBAC administration (927-947): all 17 routes `wrapAuthOnly`, no admin check in handlers (SEC-005).
- Drivers CRUD (1397-1403).
- Notifications (1413-1458) incl. `POST /notifications` (create).
- Themes and presets (1462-1472).
- Tenant contact info `PUT /tenant/contact-info` (725-753) and subscription actions `POST /subscription/initiate|cancel`, `GET /subscription/payments|current` (645-648): any tenant user can start payments or cancel the tenant subscription.
- `GET /me` (655).
- WebSocket `GET /kds/ws` (1318) auth handled in handler (not reviewed).

Routes with a wrapper but unenforced: every `wrapWithPermission` / `wrapKDSWithPermission` route (approximately 250). KDS routes do add per-user rate limits (`main.go:631-641`), which is the only real control on them.

Duplicate registrations: a second set of in-memory stub RBAC routes (1476-1752) is registered only when `db == nil` (`main.go:1473-1475`), using `wrapProtected` and an in-memory store seeded with user id 7 as Admin (`main.go:43-71`). Stub mode is reachable in production if the DB is down at start (SEC-012). Total `wrapWithPermission(` occurrences in main.go: 193.

## 5. Frontend vs backend comparison

| Aspect | Frontend | Backend |
|---|---|---|
| Module IDs | `lib/rbac.ts:18-39` (HR=3, CATEGORIES=2) | HR routes use 2 (mismatch) |
| Level hierarchy | NONE..ADMIN incl. EXPORT | same, but never evaluated at request time |
| Gate mechanism | `checkPermission`, `PermissionBoundary`, `DashboardAuthProvider` | `WithRequiredPermission` context annotation only |
| Gate coverage | Pages gated: HR (employees, attendance, leaves, payroll), accounting core pages, orders, KDS, tables, POS, reports (overview/sales/profitability/inventory/hr), settings/roles/users, notifications | all routes authenticated, none authorized |
| Ungated pages (no direct gate call) | accounting/advanced, categories*, drivers*, hr/contracts|insurance|performance|resignations|training, inventory, products*, rbac/roles, reports/financial/*, reports root, settings (access-control, account, appearance, billing, restaurant), theme-builder* | n/a |
| Fallback when permissions fetch fails | `RBAC_PERMISSIONS_TIMEOUT_MS = 2500` (`lib/rbac.ts:7`); behaviour on timeout not traced | n/a |

## 6. Required regression tests (see issues JSON for IDs)
1. Per route table-driven test: token for each seeded role vs every route in main.go; expect 403 when level < required (needs enforcement middleware).
2. Test that `X-Impersonation-Token` with garbage value returns 401 on all protected routes.
3. Test that a token without `restaurant_id` cannot override tenant via headers.
4. RBAC admin endpoints: EMPLOYEE token gets 403 on role create/assign/revoke; cross-tenant revoke returns 404/403.
5. Contract test: backend module IDs used in routes equal `modules.id` by name (parse migrations or query DB).
6. Existing `main_phase6_rbac_routes_test.go` only regex-checks that wrappers exist in source; extend with behavioural tests.
