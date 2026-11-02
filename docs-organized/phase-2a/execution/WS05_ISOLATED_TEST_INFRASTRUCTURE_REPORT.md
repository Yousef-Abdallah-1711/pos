# WS05 — Isolated Test Infrastructure Report

**Status:** PARTIAL / NOT VERIFIED
**Date:** 2026-11-02

## Current state

The existing disposable test project contains `pos-test-backend` and `pos-test-postgres`; Docker reports the test PostgreSQL healthy. It does not currently provide the required Redis and four frontend services or the deterministic two-tenant/two-restaurants-per-tenant fixture and browser journeys.

The `pos-dev` Compose project was inspected read-only using the base and local Compose files plus `.env.local`. All seven services were up; PostgreSQL and Redis were healthy. API health, Landing `/en`, Dashboard `/en`, Restaurant Website `/en` (including a restaurant host header), and Super Admin `/login` returned HTTP 200. Direct DNS resolution of `restaurant-1.localhost` is unavailable in this shell, so the browser's custom-host resolution remains unverified. No service was started, stopped, recreated, or removed by this workstream.

## Safety boundary

No fixture seeding or destructive test was run against `pos-dev`. Existing test/dev volumes were not removed. WS05 changes were not made because the complete disposable topology and ownership fixture require further verified wiring and cleanup boundaries.

## Blockers

- Add Redis and all four frontends to the disposable project with distinct project/network/volume names.
- Add deterministic synthetic tenants, restaurants, roles and representative resources, including super-admin and storefront domains.
- Prove no connection path reaches `pos-dev` or production.
- Add repeatable startup, smoke/browser execution, and teardown evidence.

## Closure execution update — 2026-11-02

Added `docker-compose.security-closure.yml` for the minimum WS02/WS03 database-backed tests. Its unique project runs PostgreSQL 15 on loopback port 15434 with `/var/lib/postgresql/data` on tmpfs and disposable labels. The full migration directory was applied to this database. Tests used the cached Go 1.24 container, a read-only module cache and `GOPROXY=off`.

This project completed the driver ownership and payment concurrency/rollback tests. It does not replace the broader WS05 target: no isolated Redis, four frontend apps or browser journey were added. The project was torn down by its exact Compose project name after verification; the Compose file and repeat-run instructions remain.
