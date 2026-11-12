# WS07 — Backend Architecture Report

**Status:** PARTIAL — method/path parity baseline added; no route extraction completed
**Date:** 2026-11-03

## Previous checkpoint — 2026-11-02

The current architecture audit and route authorization matrix remain the source inventory. No route-registration extraction or handler/use-case/repository restructuring was attempted in that security-foundation session. This avoided broad route-contract changes while WS02's cross-resource isolation gate remained incomplete.

Backend `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` passed in network-disabled Go 1.24 containers with source and module-cache mounts read-only. No route-manifest comparison was performed, so architecture acceptance remained open.

## Current evidence

- `backend/cmd/api/main.go` remains the composition root and contains extensive dependency construction and route registration.
- The router is Go's `http.ServeMux`; no machine-readable inventory of method, path, handler, authentication, authorization, tenant context, and restaurant context was captured before a refactor.
- No route registration was moved and no payment-specific route was modified.
- The full Go 1.24.13 build, vet, and default test suite passed after separate WS09 middleware/server lifecycle changes. This is regression evidence only; it does not satisfy WS07 route parity.

## Remaining work

1. Capture a complete route/middleware manifest from the current source or a safe instrumented registration path.
2. Extract coherent non-payment route modules and reduce the composition root incrementally.
3. Compare route/middleware/tenant-context manifests after each change and run focused plus full backend regression.
4. Inspect direct SQL in HTTP handlers and move only inappropriate data access into repositories/application services while preserving transaction behavior.

No WS07 completion claim is made.

## Continuation check — 2026-11-03

The earlier audit reported 363 `mux.Handle` registrations in `backend/cmd/api/main.go`. That was a source-only count and is not directly comparable to the active literal method/path patterns in the later manifest.

## Route parity baseline — 2026-11-03

Added `cmd/api/routes.manifest` for the 359 active route registrations and `TestRegisteredRouteManifestParity`, which parses active `mux.Handle`/`mux.HandleFunc` calls across the package and compares the sorted method/path patterns against the checked-in manifest. The test passed before adding the authenticated inventory-adjustment route. The manifest was updated to 360 active registrations afterward; final parity verification is included in the current Go regression run. This establishes method/path parity only; middleware and authorization parity are not represented. No route extraction is claimed.

## Continuation verification — 2026-11-03

The full Go 1.24.13 `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` regression passed in a network-disabled container with source and module-cache volumes read-only. The full suite includes `TestRegisteredRouteManifestParity` and passed against all 360 checked-in active method/path patterns after adding the inventory adjustment route. The 360-pattern baseline was generated from current source; it does not prove historical equivalence with the older 363 source count, nor does it include middleware, authentication, tenant context, or authorization. No route was extracted in this continuation, so no extraction parity claim is made. WS07 remains **PARTIAL**.

## Current route count reconciliation — 2026-11-03

The checked-in manifest contains 360 active literal method/path registrations, and `TestRegisteredRouteManifestParity` passed in the latest full Go test run. The previously cited 363 count was a historical source count without a checked-in method/path manifest, so it is not directly comparable and cannot be proven as the current active registration count. The route-manifest check validates only method/path patterns. It does not establish historical parity or middleware/auth/tenant/permission parity. Extracted the three health, readiness, and metrics registrations into `cmd/api/health_routes.go`; the original metrics internal-token middleware and public health handlers are preserved. `TestRegisteredRouteManifestParity` passed with the extracted registrations included. The route module test suite verifies healthy and unavailable readiness responses. This small extraction does not establish middleware and authorization parity for the other route families; WS07 remains partial.

## Extracted health route module verification — 2026-11-03

Moved the unchanged public health and readiness endpoints plus the internally protected metrics endpoint from `cmd/api/main.go` to `cmd/api/health_routes.go`, with route wiring still at the same composition point. Added readiness tests for both healthy and unavailable database responses. The full backend build, vet, and test suite passed after extraction; the checked-in 360-pattern method/path manifest test passed. The older 363 count remains historical and not comparable without its original manifest, and middleware/authorization parity outside this extracted module remains unverified.

## Health-route middleware parity regression — 2026-11-03

Added a route-level test proving the extracted metrics endpoint still returns 401 without its configured internal token and 200 with the token. The full backend gate passed with this test and the existing 360-pattern method/path manifest check. This verifies middleware behavior for the health/metrics extraction only; authentication, tenant, restaurant, and permission parity across the remaining route families remains open.

## Translation-route extraction and authorization parity — 2026-11-03

Extracted the five translation route registrations from `cmd/api/main.go` into `cmd/api/translation_routes.go`. Public translation health remains unauthenticated, and the upload route remains protected by Super Admin token middleware. The focused registration test passed for the public route, missing-token denial, and valid Super Admin token; `TestRegisteredRouteManifestParity` also passed with all 360 method/path patterns unchanged. Middleware and authorization parity for the rest of the router remains open; WS07 is still **PARTIAL**.
