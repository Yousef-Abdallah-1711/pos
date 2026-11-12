# Phase 2A exit criteria

| Area | Measurable exit criterion | Current status |
|---|---|---|
| Architecture | Module ownership and dependency direction documented; composition root has bounded extraction plan; no unknown route loss; contracts documented | PARTIAL: documented target; implementation still centralized and mixed-layer. |
| Security | All route classes inventoried; auth/RBAC/public boundaries tested; no uncontained critical runtime privilege; priority tenant/restaurant paths have negative tests | FAIL: route inventory/selected tests exist; dev superuser remains; isolation incomplete. |
| Database | Runtime/migration roles separated and rehearsed; 91-table ownership manually reviewed; RLS staged behind transaction adoption | FAIL/PARTIAL: test role works, dev still superuser, RLS absent. |
| Backend tests | Build/vet/default tests reproducible and required in CI | PARTIAL: checkpoint all pass, but no CI gates. |
| Frontend tests | All four intended suites runnable; app/test typechecks visible; CI test gate stable | FAIL: dashboard red, landing alias failures, storefront 8 failures; only super-admin green. |
| Integration | Isolated fixtures cover two tenants and multiple restaurants; all four apps connect to isolated backend; browser journeys and security smoke repeat | FAIL: test stack is API+Postgres only; 73 selected smoke checks. |
| Engineering | Supported commands/config and migration strategy documented; local service health/URLs known; no dev stack mutation during audit | PARTIAL: topology and URLs verified; configuration/security drift remains. |
| Error/observability | Error mapping, request correlation, redaction and worker shutdown criteria defined and tested | NOT VERIFIED/PARTIAL. |

Later phases may own broad RLS rollout, load testing, production alerting and manual cross-browser tests. Phase 2A still requires a restricted runtime role strategy, complete security boundaries, stable test gates and reproducible isolated regression before closure.
