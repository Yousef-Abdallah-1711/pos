# Settings, RBAC, users, billing (dashboard) - static audit

Status: **Verified implemented** (roles/users/permissions UI, account, appearance, restaurant contact, billing); implementation style is the least consistent in the app (raw fetch).

## Pages (10 settings + rbac)
`settings` (redirect index), `account` (238; raw fetch profile + change-password), `appearance` (328; raw fetch theme/theme-colors), `restaurant` (170; api GET/PUT /tenant/contact-info), `users` (753 lines; 9 raw fetch; gate via useHasPermission), `roles` (230), `roles/new` (357; 24 gate refs), `roles/[id]/edit` (435; 23), `access-control` (41, wrapper), `billing` (454; Paymob iframe via `iframeURL`, no gate - bypasses tenant-blocked overlay by design `DashboardLayout` blockedBypassPaths). Plus `rbac/roles/page.tsx` (351, ungated duplicate of settings/roles using `lib/rbac.ts`).

## Contract
Match main.go: `/users` (list/create/update/delete), `/users/{id}/profile`, `/users/{id}/settings/{theme,theme-colors,change-password}`, `/users/{id}/settings`, `/rbac/roles` CRUD, `/rbac/roles/{id}/permissions`, `/rbac/permissions/assign|{role}/{module}`, `/rbac/users/{id}/roles`, `/rbac/users-by-email/*`, `/rbac/me/permissions`, `/rbac/check-permission`, `/subscription/{plans,current,payments,initiate,cancel}`, `/tenant/status`, `/tenant/contact-info`. Unused: `PUT /users/{id}/settings/language`, `GET /me`.
Server: `/users*` and `/tenant/contact-info` use `wrapProtected` only (main.go:1062-1264, 725) - management of users relies on handler-internal checks (not read; backend auditors).

## RBAC model (client)
`lib/rbac.ts`: MODULE_IDS 1-20, levels NONE<READ<EXPORT<WRITE<DELETE<ADMIN. `hooks/useRbac.ts`: `useHasPermission` (one /rbac/check-permission request per hook instance), `useUserPermissions` (all perms, cached), `useRoles`. Cache key constant 'current_user' (FE-DASH-015). Sidebar uses `useUserPermissions().checkPermission` for 17 modules.

## Auth flow and stores
Login lives in landing app; dashboard accepts `?auth_token&auth_user` (DashboardAuthProvider), `/auth/callback?token&user`, `/auth/super-admin-login?token` (POST /auth/token-exchange) or `?impersonation_token` (GET /impersonation/validate). Token in `auth-storage` localStorage. Stores: `authStore` (persist), `preferencesStore` (persist 'preferences-storage', sets documentElement dir/lang), `tenantStatusStore` (not persisted; fail-open `canAccess=true` on non-403 errors, `tenantStatusStore.ts:25,89-93`), `posStore`, `kdsStore`, `themeBuilderStore` (persist). Blocked tenants: overlay/sidebar in DashboardLayout/`components/blocked/*`.

## Findings
FE-DASH-005,006,007,008,014,015,021; users page error handling = raw strings (FE-DASH-026).
