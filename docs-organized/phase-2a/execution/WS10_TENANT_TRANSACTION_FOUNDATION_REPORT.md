# WS10 — Tenant Transaction Foundation Report

**Status:** PARTIAL — integration proof freshly rerun; application-module adoption remains open
**Date:** 2026-11-02

The existing audit records `backend/internal/pkg/database/tenant_tx.go` and a disposable pilot. In this execution, `go test -tags integration -count=1 ./internal/pkg/database` passed against a newly created PostgreSQL 15 container with a tmpfs data directory, a unique container name, and restricted `pos_app` credentials. The pilot exercises transaction-local tenant settings, connection reuse, fail-closed access without context, cross-tenant write denial, and restricted-role attributes. The temporary database container was stopped after the run.

No RLS policy was enabled on development data. A repository usage survey and selected business-module adoption remain open, so this workstream is not complete.

## Fresh verification — 2026-11-03

The tagged integration test passed again against a uniquely named PostgreSQL 15 service with tmpfs storage and a non-superuser, `NOBYPASSRLS` `pos_app` role. The test checked transaction-local context, forced connection reuse, fail-closed behavior without context, rejected cross-tenant inserts, and unchanged rows for unauthorized update/delete attempts. The isolated Compose project was torn down by its exact project name. No development DB, volume, or RLS policy was touched.

The helper still has no selected production business-module adoption or broader concurrent tenant request/pool-reuse matrix. WS10 remains **PARTIAL**.

## Isolated inventory transaction — 2026-11-03

The new WS05 inventory adjustment handler performs a row-locked, tenant/restaurant-scoped product update and its movement audit insert in the same SQL transaction. Browser execution against the disposable WS05 database returned HTTP 200, and PostgreSQL confirmed both the quantity delta and exactly one movement row. This is application transaction evidence for one non-payment business operation; it does not enable RLS on business tables or complete the requested concurrent request, rollback-fault injection, and pool-reuse matrix. WS10 remains **PARTIAL**.

## Inventory adjustment rollback fault injection — 2026-11-03

Added `TestInventoryAdjustmentPostgresRollsBackWhenMovementInsertFails`. On the isolated PostgreSQL closure database, a scoped test trigger injects an error only for the synthetic fixture movement after the handler has updated stock. The authenticated handler returned HTTP 500; direct PostgreSQL assertions confirmed product stock remained at its original value and no movement row was committed. The trigger/function are dropped by test cleanup. The test passed with `go test -count=1 ./internal/handler/http -run '^TestInventoryAdjustmentPostgresRollsBackWhenMovementInsertFails$'`. This adds rollback evidence to the existing concurrent-adjustment evidence, but broader business-module adoption and connection-pool tenant/restaurant isolation remain open.
