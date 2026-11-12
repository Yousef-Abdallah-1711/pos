# Endpoint Response Contracts (Static, Partial)

The current source defines success/error response helpers and endpoint-specific handler contracts, but they are not uniform: existing documentation describes `{success,data}` while API route handlers include legacy shapes and inline fallbacks. The full contract must be derived per handler; no generic response schema should be assumed.

## Static baseline

- Current checkout route registration manifest has 360 registration expressions (GET 183, POST 99, PUT 44, DELETE 34), with 335 unique method/path pairs. The 25 repeated pairs include alternative conditional handler registrations and are not proof of simultaneous runtime routes.
- `08_COMPLETE_BACKEND_ENDPOINT_INVENTORY.csv` lists all manifest occurrences; it does not claim all conditional handlers are present at runtime.
- Handler response bodies, status paths, validation schemas and error branches have not been extracted for all routes. Endpoint runtime data is only available for GET `/api/v1/health` in the isolated stack.

## Known examples requiring per-handler contract extraction

- Health route returns a 69-byte body in sampled E2E runtime (body intentionally not stored). Status was 200 for 25/25 requests.
- Product list handler parses paging/filter/sort query values; source implementation should be checked for constraints before clients rely on contract. Static hotspot: sort identifiers are interpolated into SQL in `product_repo.go`; validate allowlist semantics.
- Admin orders listing returns a list then hydrates each record separately; there is no runtime body sample.

No sensitive response body was saved. Runtime measurements are in `10_ENDPOINT_RUNTIME_PERFORMANCE.csv` and summary JSON.
