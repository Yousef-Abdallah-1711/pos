# Lead-auditor verification notes (live checks on local throwaway dev stack)

Checks were read-only GETs against http://localhost:8080 (dev DB, seed data only). Raw output: network-logs/sec-repro-local.txt

| Claim | Result |
|---|---|
| SEC-001 / BE-001: any non-empty X-Impersonation-Token skips JWT auth | CONFIRMED (partially). No header: 401. With header: request passes authentication. |
| ...and then reaches data with header-chosen tenant | NOT CONFIRMED for /products: got 403 "Insufficient permissions" from product_handler.go:59 (an in-handler role check). Impact depends on each handler's own role check; needs per-route Phase 2 testing. |
| SEC-003 / BE-002: "no permission enforced on any route" | OVERSTATED. middleware.WithRequiredPermission (permission.go:16-24) only stores context values (confirmed), but individual handlers perform ad-hoc role checks (e.g. product_handler.go). Correct wording: the *middleware layer* enforces nothing; enforcement is inconsistent and handler-local. Severity stays P0 until a route-by-route test shows which handlers lack a check. |
| QA-001: /health/metrics public | CONFIRMED (HTTP 200 without auth). |
| DB-002: public order lookup | Order id 1 returned 404 (no orders in dev DB); sequential-id read risk NOT reproducible without fixtures -> BLOCKED REQUIRES ISOLATED TEST DATA. |
| DB-005: no RLS | CONFIRMED live: 0 policies, 0 RLS-enabled tables, only superuser role (database-evidence/live-catalog.txt). |
