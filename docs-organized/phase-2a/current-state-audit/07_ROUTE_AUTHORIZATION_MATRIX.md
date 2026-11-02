# Route authorization matrix

Current route census is derived from `backend/cmd/api/main.go` and existing generated inventory: 363 registrations, 338 distinct method/path pairs. Existing method/path-level matrix is retained at `docs-organized/phase-2a/route-authorization-matrix.json` (and Markdown at `docs-organized/phase-2a/ROUTE_AUTHORIZATION_MATRIX.md`); those files list grouped route classes and per-module counts. Counts are registration categories, not a test count.

| Authorization class | Registrations | Required validation |
|---|---:|---|
| Tenant permission | 274 | JWT + validated tenant/restaurant + permission module/level + resource ownership |
| Public | 38 | Explicit public contract, rate limits where appropriate, tenant/restaurant resolution, minimal serialization |
| Authenticated only | 20 | JWT + self-resource or billing/status ownership checks |
| Super-admin | 26 | Super-admin authentication and platform scope; audit/impersonation controls |
| Session token | 3 | Server-side impersonation-session validation, expiry and revocation |
| Service internal | 1 | Internal service token, not client identity |
| Webhook | 1 | Provider signature, idempotency and transaction/state validation |

## Findings

- Route wrappers are centralized in `cmd/api/main.go`; middleware sequence in current architecture notes is authentication → tenant context → tenant status → permission → handler, with public/internal/webhook exceptions.
- `RBAC_ENFORCEMENT=enforce` is explicit in `docker-compose.local.yml` and isolated test compose. Production is forced to enforce. Generic non-production fallback is `audit` according to the mode parser/current task record; this is not a fail-closed default.
- Auth-only paths include self-service profile/settings, current user/permission introspection, subscription operations and tenant status. Each requires object-level ownership, even where no module permission is appropriate.
- The public bucket includes login, storefront reads/order creation/tracking, translation resources, health, mock payment endpoints in mock mode and KDS WebSocket; public classification is not equivalent to unrestricted data exposure.
- KDS WebSocket module permission is reported as wired in current code; actual handshake, tenant binding and browser flow remain unverified.
- Route matrix smoke proved select denied/allowed route cases using disposable fixtures; it did not exercise all 338 distinct routes or all roles/permission levels.

## Verification evidence

Classification: CODE VERIFIED through route registration wrappers and generated matrix; INTEGRATION TEST VERIFIED for selected routes by 73/73 disposable smoke; NOT VERIFIED exhaustively for all resource IDs/roles. Preserve the generated JSON as the complete machine-readable listing and rerun generator whenever route registration changes.
