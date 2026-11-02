# Phase 2A completion roadmap

| Order / task | Priority | Missing work and approach | Dependencies | Required verification / acceptance | Complexity |
|---|---|---|---|---|---|
| 1. Restrict dev DB role | P0 | Rehearse current owner/grant transfer on copied disposable DB; configure distinct migration/runtime roles | Migration grant history, backup/restore rehearsal | API connects as non-superuser; DDL and migration table denied; migration + app smoke pass | L |
| 2. Close resource isolation gaps | P0/P1 | Recount driver methods; test foreign tenant and restaurant resources across products/categories/orders/inventory/users/roles/themes/settings | Stable tenant fixtures and route map | Negative GET/PUT/DELETE and side-effect checks pass for all table families | XL |
| 3. Complete payment callback integrity | P1 | Implement/verify transactional, idempotent success activation and concurrent/out-of-order callback behavior | Disposable payment fixtures/mock | Duplicate/reordered/concurrent callbacks result in one correct state transition and balanced records | L |
| 4. Recover frontend test baseline | P1 | Correct alias/config and stale contracts; triage failures without disabling assertions | Workspace install lockfile, maintainers for expected UX | All intended tests run with zero failures; app and test typechecks pass | L |
| 5. Expand isolated fixtures/stack | P1 | Add Redis and four apps, two tenants × two restaurants, users/roles/products/orders/subscription/super-admin fixtures | Docker compose and deterministic seed design | Repeatable full stack startup + teardown isolated from `pos-dev`; smoke/E2E stable | XL |
| 6. CI gates | P1 | Add Go build/vet/test, frontend checks, disposable security smoke, migration validation and scans | Test repairs and documented toolchain | Required checks pass on PR and failures block merge | M/L |
| 7. Architecture extraction | P2 | Extract route modules then migrate handlers to use case/repo with contract-level route diff | Stable Go test/security baseline | No method/path loss; middleware-order checks and full regression pass each module | XL |
| 8. Frontend contracts | P2 | Select authoritative theme model; resolve package ownership and API client duplication | Cross-app inventory | No casts bridging incompatible model at boundaries; all app checks pass | L |
| 9. Error/observability foundations | P2 | Typed errors, common HTTP mapping, request IDs/redaction, worker shutdown | Backend boundary extraction | Unit/integration tests for mappings, logging redaction and shutdown | L |
| 10. Staged RLS | P2/2C | Migrate module access to `WithTenantTx`, then enable policies table-family by family | Resource contract and restricted role | Two-tenant tests prove fail-closed RLS on each family, pooled connections do not leak context | XL |

Priorities: P0 containment before shared/production-like data; P1 blocks foundation closure; P2 architecture/hardening. Full RLS rollout belongs to 2C, but safe transaction boundary and least-privilege foundation are required now.
