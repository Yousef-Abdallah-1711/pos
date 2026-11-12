# Evidence Method and Redaction

The package stores aggregate measurements, production Compose configuration, build fingerprints, and sanitized browser console/page snapshots. It does not include auth storage state, cookies, tokens, sensitive headers, credential-bearing HAR, or raw trace payloads. No production data was used.

Reproduction: Compose project pos-p2b-perf-20261103; synthetic fixture counts in 01_PRODUCTION_BROWSER_BASELINE.md; production-build Next.js 15.5.9. Admin order samples: 25 sequential authenticated GETs per page size 1, 5, 25, 50. Dashboard readiness: five hard reload samples. Anonymous wait timeout: three 8-second networkidle attempts. SQL plans used EXPLAIN ANALYZE only on verified read queries against the disposable synthetic DB.

Audit-log evidence keeps only status, row count, field keys, and null/object counts; row content was not retained. No optimization experiment or application source modification was performed.
