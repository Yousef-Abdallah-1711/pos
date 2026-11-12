# Observability Gap Analysis

| Capability | State | Evidence |
|---|---|---|
| Structured logging | Absent | 0 slog/zap/logrus/zerolog uses; 1325 log.Printf/Println |
| Request ID / correlation | Absent | grep for X-Request-ID found none in middleware (only unrelated "requestID" path param and GDPR ids) |
| Access log | Absent | RequestMetricsMiddleware (middleware/request_metrics.go) records latency only |
| Metrics | Partial | internal/observability/perf_metrics.go: in-memory ring of 240 samples per metric (avg/p95/max); not Prometheus format; resets on restart |
| Metrics endpoint exposure | Risk | GET /api/v1/health/metrics registered with plain mux.HandleFunc (main.go:547) with no auth wrapper; also a test route `/api/v1/public/test-routing/{slug}` (main.go:556) left in |
| Health | Present | /api/v1/health, /health/ready (db_healthy), DB ping monitor every interval logs state changes |
| Error tracking | Absent | sentry-go only an indirect dep, no usage |
| Log redaction | Not found | no redaction helper; raw err text logged and returned |
| Tracing | Absent | no OpenTelemetry |
| Frontend | Partial | web-vitals, lib/monitoring.ts, error-logger.ts with optional remote endpoint; @vercel/analytics in dashboard and restaurant-website |
| Alerting/dashboards | Not found | none in repo |
