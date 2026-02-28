# KDS (Kitchen Display System) - Single Source of Truth

**Document Type:** Implementation + Operations + File Manifest  
**Status:** Implemented and hardened  
**Last Updated:** 2026-02-25

---

## 1) Scope and Goal

KDS is the kitchen execution layer for restaurant orders.

Primary outcomes delivered:
- Real-time kitchen queue for confirmed orders.
- Status transitions: `new -> in_progress -> ready -> completed` (+ `recalled -> ready`).
- RBAC protection using **Module 13 (KDS)**.
- WebSocket live updates per tenant/restaurant room.
- KDS settings + analytics endpoints.
- Offline-safe and resilience behaviors on frontend (snapshot + persisted UI preferences).
- Table context shown on KDS cards for dine-in orders.

---

## 2) Phase Status (As Built)

### Phase 1: DB + RBAC Foundation
Completed.
- KDS queue table created.
- KDS settings table created.
- Module 13 seeded into RBAC modules.
- Role backfill migration applied for KDS permissions.
- Performance indexes and analytics-panel default migration applied.

### Phase 2: Backend (Repo/Usecase/Handler/Routes + Order Integration)
Completed.
- Repository, use case, HTTP handler implemented.
- Order integration: auto-enqueue when order transitions to `confirmed`.
- Status synchronization back into order lifecycle on bump transitions.
- Main API wiring completed with KDS permission wrappers.

### Phase 3: Frontend KDS Core UI
Completed.
- KDS page + fullscreen shell behavior.
- Header, grid, timer, order card components.
- Status filters and actions.

### Phase 4: Hardening
Completed.
- API timeout handling and graceful fallbacks.
- KDS-focused per-user rate limiting.
- Error boundary and degraded-state behavior.
- Local storage hardening for malformed persisted payloads.

### Phase 5/6/7: UX, Offline, Settings, Analytics, and Integration
Completed.
- WebSocket + offline queue replay support.
- Audio notifications + sound toggle.
- Settings page with ADMIN guard.
- Analytics panel support and backend analytics endpoint.
- Table context appears in KDS queue/order cards.

### Phase 8: Tests + Hardening Checklist
Implemented with current KDS test suite and integration coverage present in backend and frontend.

---

## 3) Runtime Architecture

### Backend flow
1. Order reaches `confirmed`.
2. `OrderUseCase` calls KDS repository enqueue.
3. KDS queue row is inserted/upserted.
4. WebSocket broadcast (`order_new`) sent to room `tenant:{id}:restaurant:{id}`.
5. KDS clients receive updates; fallback polling still works.

### Frontend flow
1. User opens `/dashboard/kds`.
2. RBAC checks READ permission for module 13.
3. `useKDS` fetches queue/stats/settings, starts polling fallback.
4. `useKDSWebSocket` streams real-time updates.
5. Actions (bump/recall) are optimistic with rollback support.
6. Offline actions are queued and replayed when online.

---

## 4) RBAC and Security Rules

- **Module ID:** `13` (`KDS`)
- **READ:** view queue/stats/settings
- **WRITE:** bump and recall orders
- **ADMIN:** update settings, read analytics

Route protection uses KDS-specific wrapper in API main routing with:
- Auth middleware
- Tenant context middleware
- `WithRequiredPermission(13, level)`
- Rate limits:
  - General KDS endpoints: `RateLimitByUser(60)`
  - Bump endpoint extra burst control: `5 requests / second`

Frontend access control:
- Sidebar KDS item appears only when module 13 READ is granted.
- KDS page shows access denied when user lacks permission.
- KDS settings page requires ADMIN.

---

## 5) API Surface (Implemented)

- `GET /api/v1/kds/ws` (websocket upgrade with JWT query token)
- `GET /api/v1/kds/queue` (READ)
- `PUT /api/v1/kds/orders/{id}/bump` (WRITE)
- `PUT /api/v1/kds/orders/{id}/recall` (WRITE)
- `GET /api/v1/kds/stats` (READ)
- `GET /api/v1/kds/analytics` (ADMIN)
- `GET /api/v1/kds/settings` (READ)
- `PUT /api/v1/kds/settings` (ADMIN)

---

## 6) Data Model Notes

### `kds_queue`
Key fields:
- Tenant + restaurant scoping
- Order linkage (`order_id`, `order_number`)
- Status, station, priority
- Denormalized `items` JSONB
- Time fields (`received_at`, `started_at`, `completed_at`)
- `bumped_by`

### `kds_settings`
Key fields:
- `auto_bump_minutes`
- `alert_sound_enabled`
- `show_analytics_panel`
- `display_density`
- Warning/danger minute thresholds

### Table context in KDS
KDS queue read joins `orders` and `restaurant_tables` so KDS order payload can include:
- `table_id`
- `table_number`
- `table_name`

---

## 7) Hardening Implemented

### Backend hardening
- Timeout wrappers per endpoint type.
- Graceful fallback payloads when KDS storage schema is temporarily unavailable.
- Rate limiting for KDS endpoints and bump spam control.
- WebSocket room isolation and connection limit per tenant/restaurant.

### Frontend hardening
- Snapshot bootstrap with max-age guard.
- Persisted store safety:
  - malformed `kds-storage` payload is now auto-cleared.
  - malformed bootstrap snapshot is auto-cleared.
- Optimistic updates with rollback.
- Offline action queue + replay.

---

## 8) KDS UX Flows

### Kitchen operator flow
1. Open `Dashboard -> Kitchen Display`.
2. Filter queue (`all`, `new`, `in_progress`, `ready`).
3. Open card details and bump through statuses.
4. Recall completed item if needed.
5. Use fullscreen mode in kitchen screens.

### Manager/admin flow
1. Open KDS settings.
2. Configure auto-bump, density, audio, color thresholds.
3. Toggle analytics panel visibility.
4. Monitor analytics trends and busiest hours.

### Failure/degraded flow
- If websocket disconnects: polling continues.
- If browser offline: actions queue locally and replay on reconnection.
- If local persisted JSON is corrupted: state auto-recovers to defaults (no runtime crash).

---

## 9) Test Coverage Snapshot

Backend tests include:
- Repository tests (`kds_repo_test.go`)
- Use case tests (`kds_usecase_test.go`)
- HTTP handler tests (`kds_handler_test.go`)
- Analytics integration tests (`kds_analytics_integration_test.go`)
- Flow E2E tests (`kds_flow_e2e_test.go`)
- WebSocket hub/handler tests

Frontend tests include:
- KDS component tests (`KDSGrid`, `KDSHeader`, `KDSOrderCard`)
- Hook tests (`useKDSAnalytics`)
- Store persistence hardening tests (`kdsStore.test.ts`)

---

## 10) Complete File Manifest (Created/Edited for KDS Feature)

### A) Backend migrations
- `backend/migrations/124_create_kds_queue.sql`
- `backend/migrations/125_create_kds_settings.sql`
- `backend/migrations/126_add_kds_rbac_module.sql`
- `backend/migrations/127_backfill_kds_role_permissions.sql`
- `backend/migrations/128_add_kds_show_analytics_setting.sql`
- `backend/migrations/129_add_kds_performance_indexes.sql`
- `backend/migrations/134_set_kds_analytics_default_hidden.sql`

### B) Backend domain/repository/usecase
- `backend/internal/domain/kds.go`
- `backend/internal/repository/kds_repo.go`
- `backend/internal/repository/kds_repo_test.go`
- `backend/internal/usecase/kds_usecase.go`
- `backend/internal/usecase/kds_usecase_test.go`
- `backend/internal/usecase/order_usecase.go` (order -> KDS enqueue integration)

### C) Backend HTTP/WS and wiring
- `backend/internal/handler/http/kds_handler.go`
- `backend/internal/handler/http/kds_handler_test.go`
- `backend/internal/handler/http/kds_analytics_integration_test.go`
- `backend/internal/handler/http/kds_flow_e2e_test.go`
- `backend/internal/handler/ws/kds_hub.go`
- `backend/internal/handler/ws/kds_hub_test.go`
- `backend/internal/handler/ws/kds_ws_handler.go`
- `backend/internal/handler/ws/kds_ws_handler_test.go`
- `backend/cmd/api/main.go` (dependency wiring + route registration + KDS rate-limited wrapper)

### D) Frontend KDS pages/components
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/layout.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/page.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/settings/page.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSAnalyticsPanel.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSErrorBoundary.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSGrid.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSHeader.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSOrderCard.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSTimer.tsx`

### E) Frontend hooks/store
- `frontend/apps/dashboard/src/hooks/useKDS.ts`
- `frontend/apps/dashboard/src/hooks/useKDSAnalytics.ts`
- `frontend/apps/dashboard/src/hooks/useKDSAudio.ts`
- `frontend/apps/dashboard/src/hooks/useKDSOffline.ts`
- `frontend/apps/dashboard/src/hooks/useKDSWebSocket.ts`
- `frontend/apps/dashboard/src/stores/kdsStore.ts`

### F) Frontend RBAC/layout/i18n integration
- `frontend/apps/dashboard/src/components/layout/DashboardLayout.tsx`
- `frontend/apps/dashboard/src/components/layout/Sidebar.tsx`
- `frontend/apps/dashboard/src/hooks/useRbac.ts`
- `frontend/apps/dashboard/src/lib/rbac.ts`
- `frontend/apps/dashboard/src/i18n/messages/en.json`
- `frontend/apps/dashboard/src/i18n/messages/ar.json`

### G) Frontend tests
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSGrid.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSHeader.test.tsx`
- `frontend/apps/dashboard/src/app/[locale]/dashboard/kds/components/KDSOrderCard.test.tsx`
- `frontend/apps/dashboard/src/hooks/useKDSAnalytics.test.ts`
- `frontend/apps/dashboard/src/stores/kdsStore.test.ts`

### H) KDS docs (this folder)
- `docs-organized/12-KDS-KITCHEN-DISPLAY-SYSTEM/KDS-IMPLEMENTATION-PLAN.md` (single source of truth)

---

## 11) Recent Critical Fix

Issue observed on `http://localhost:3002/en/dashboard/kds`:
- Runtime error: `Unexpected end of JSON input`.

Root cause:
- Corrupted localStorage payload for KDS persisted state.

Fix applied:
- Safe persisted storage parser in `kdsStore.ts`.
- Auto-remove malformed persisted `kds-storage` value.
- Auto-remove malformed `kds-bootstrap-snapshot-v1` in `useKDS.ts`.
- Regression test added in `kdsStore.test.ts`.

---

## 12) Directory Policy

This KDS docs directory is intentionally reduced to **one file only**:
- `KDS-IMPLEMENTATION-PLAN.md`

This file replaces the previous split between implementation plan and checklist.
