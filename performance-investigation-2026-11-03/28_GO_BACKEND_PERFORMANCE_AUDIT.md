# Go Backend Performance Audit

## Runtime architecture

Go stdlib `net/http` server. Main route setup in `backend/cmd/api/main.go`, health routes in `health_routes.go`; global middleware includes request ID, recovery, request metrics, security headers, gzip and CORS. Request metrics store bounded in-memory route/status durations (up to 240 samples per metric); selected KDS repository operations record database timing. No complete spans dividing middleware/auth/handler/repository/Redis/serialization were found. No pprof, Prometheus or OpenTelemetry path was found in source search.

## Static latency risks

- Admin order listing calls `ListOrders` then `GetOrder` per result (`internal/handler/http/admin_order_handler.go:108–225`), yielding a likely N+1 pattern. Export also hydrates records per order. No measured query count or latency.
- Repository direct `Query`/`QueryRow`/`Exec` call sites outnumber context-aware variants in static scan (~375 vs ~211); this is not unique method or runtime use. Cancellation/timeout behavior must be traced by hot endpoint.
- Product repository builds ORDER BY using filter strings (`internal/repository/product_repo.go:~285–293`); audit validation and allowlist in addition to performance. Do not benchmark arbitrary sort expressions.
- DB pool is configurable (75 open / 25 idle defaults; 5m max lifetime, 2m idle timeout). No pool saturation runtime evidence was collected.
- Redis is not on all request paths and falls back to process memory at init time. Its latency was not attributed per request.

## Runtime evidence

Only GET `/api/v1/health` was measured. 25 sequential requests to isolated WS05 E2E stack returned HTTP 200, 69 bytes; total latency median 2.18 ms, p75 2.46 ms, p95 4.94 ms, max 126.93 ms. These are host-to-loopback E2E development measurements for a health endpoint, not a business API and not production. Endpoint response body was not stored. No endpoint-level DB count/timing was observable.

`29_BACKEND_LATENCY_BREAKDOWN.csv` records measured and unavailable segments.
