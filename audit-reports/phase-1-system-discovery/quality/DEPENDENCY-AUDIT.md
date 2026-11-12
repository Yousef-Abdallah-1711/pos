# Dependency Audit

Method: manifests, frontend/pnpm-lock.yaml and backend/go.mod read only. "Latest known" was NOT checked against any registry (no lookup performed), so it is marked unverified. No CVE claims are made.

| Package | Installed (range -> lock) | Latest known | Risk | Upgrade complexity | Recommendation |
|---|---|---|---|---|---|
| next | ^15.1.0 -> 15.5.9 (all 4 apps) | unverified | Med: caret range | Med | pin; check advisories before prod |
| react / react-dom | ^18.3.1 -> 18.3.1 (shared-components ^18.2.0) | unverified (19 exists) | Low | High for 19 | stay; align shared-components |
| typescript | ^5.7.2 -> 5.9.3 (shared-components ^5.3.3) | unverified | Low | Low | align |
| pnpm | 9.0.0 pinned (packageManager, Dockerfile.prod corepack; CI `version: 9`) | unverified | Low-Med: first 9.x release | Low | bump to latest 9.x |
| turbo | ^2.3.0 | unverified | Low | Low | ok |
| tailwindcss | ^3.4.17 -> 3.4.19 | unverified (v4 exists) | Low | High for v4 | stay |
| zustand | ^4.5.0 -> 4.5.7 | unverified (v5 exists) | Low | Med | stay |
| @tanstack/react-query | ^5.59.20 (dashboard, restaurant-website, super-admin; not landing) | unverified | Low | Low | ok |
| zod | ^3.24.1 -> 3.25.76 | unverified (v4 exists) | Low | Med | stay |
| axios | ^1.7.9 | unverified | Med: frequently-advisoried lib, check | Low | verify |
| jspdf | ^4.0.0 (dashboard) | unverified | Med: check advisories | Low | verify |
| vitest | ^4.0.16 (dashboard, landing, super-admin) | unverified | Low | - | ok |
| jest | ^29.7.0 (root deps, shared-components, restaurant-website script) | unverified | Low | - | see inconsistencies |
| @playwright/test | ^1.57.0 (root only) | unverified | Low | - | no runnable tests (see TEST-COVERAGE) |
| Go | go 1.24.0 | unverified | Low | - | fix README (says 1.18+) |
| gin v1.11.0, chi v5.2.3 | go.mod direct | - | Low: gin only in non-wired code | - | prune after confirming |
| golang-jwt v4.5.2 AND v5.3.0 | both direct | - | Med: two JWT libs (v4 in internal/auth/jwt_manager.go) | Low | consolidate on v5 |
| lib/pq v1.10.9 | direct | - | Low | Med to pgx | defer |
| redis/go-redis v9.7.0 | marked `// indirect` but imported by internal/service/cache.go | - | Low | Low | `go mod tidy` (not run) |
| sentry-go v0.34.0 | indirect; no usage found in code | - | - | - | wire up or drop |
| go-sqlmock, embedded-postgres | marked indirect though used by tests/scripts | - | Low | Low | tidy |
| golang.org/x/crypto v0.46.0 | direct | unverified | Low | - | keep current |
| PostgreSQL | postgres:15-alpine (dev, prod, test compose) | unverified | Low | Med | plan 16/17 upgrade |
| Redis | Optional: REDIS_URL cache in cache.go; no Redis service in any compose file | - | - | - | document or drop |
| Docker bases | golang:1.24-alpine, alpine:3.20, node:20-alpine (tags, no digests) | unverified | Med | Low | consider Node 22 LTS, pin digests |

## Inconsistencies across apps
- Test runner: dashboard, landing, super-admin use vitest; restaurant-website uses jest but declares no jest in its own devDependencies (comes from frontend/package.json dependencies); layers/shared-components uses jest 29 with @testing-library/react ^14 vs ^16 elsewhere.
- jsdom: dashboard ^26.1.0 vs super-admin ^27.4.0 (lockfile holds both vitest instances).
- @types/node ^22.10.2 vs ^20.10.5 (shared-components).
- @radix-ui/react-slot ^1.1.1 vs ^1.2.0 (ui-components); react-dialog ^1.1.15 vs ^1.1.2 (super-admin).
- Test tooling (jest, testing-library, ts-jest) and Radix packages sit in `dependencies` of the workspace root frontend/package.json.
- Four lockfiles tracked (see 01-REPOSITORY-INVENTORY.md).
