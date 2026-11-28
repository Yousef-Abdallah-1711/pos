# F-071 — ENV=staging is treated as development: RBAC defaults to audit-only (cashier reaches finance, HR and user endpoints), mock payments, open metrics and in-memory sessions are allowed

- **Severity:** High  ·  **Category:** Security / configuration  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Platform configuration (RBAC, payments, sessions, metrics, cookies)
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-SEC013, EX-RUN005-RBAC-MATRIX
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A staging deployment enforces the same guards as production (RBAC enforce, no mock payments, protected metrics, Secure cookies), or the weaker behaviour is explicitly documented

Observed: Enforce API: cashier gets 403 on /accounting/journal-entries, /accounting/accounts, /hr/salaries, balance sheet, profit-loss, /users. Staging API (RBAC unset): the same requests return 200 and a cashier write reaches the validator (400, not 403). ENV=production refuses to start without PAYMOB_MODE=live; staging accepts mock payments. /health/metrics answers 200 without a token and the session cookie has no Secure flag in both non-production runs. Staging does require independent signing secrets (signing_secrets.go:9) - the only guard staging shares with production RUN-005 expansion (4 roles x 28 GET endpoints + 8 write probes, same database, same users): enforce-mode API: ADMIN 27/28 GET 200, MANAGER 15/28, EMPLOYEE 10/28, DRIVER 3/28; write probes authorised 8, 2, 0, 0 of 8. ENV=staging API with RBAC unset: EVERY role (including DRIVER) gets 27/28 GET 200 with salary/bank/national-id fields in 4 responses and passes authorization on 8/8 write probes (400 body validation instead of 403). Real data was returned (sensitive-field responses 4 per role), not only status codes.

## Severity rationale
Rated **High** (category: Security / configuration; evidence: REPRODUCED). Basis recorded in the register: 0 direct; exposes all financial data and writes to any signed-in role on a staging deployment that holds real or production-like data. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Platform configuration (RBAC, payments, sessions, metrics, cookies). 0 direct; exposes all financial data and writes to any signed-in role on a staging deployment that holds real or production-like data

## Technical impact
CONFIRMED (source + runtime): config.IsProductionMode() is true only for production/prod (runtime_env.go:20-23); RBAC default mode is audit unless production (middleware/permission.go:36-42); mock Paymob allowed unless production (config/paymob_mode.go:29); metrics token and Redis/session requirements and secure cookies are gated the same way (health_routes.go:52, main.go:148-156, auth_handler.go:223)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: yes, see category. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0 direct; exposes all financial data and writes to any signed-in role on a staging deployment that holds real or production-like data

## Where it was reproduced
Start the isolated API image with ENV=staging and RBAC_ENFORCEMENT unset (same DB/Redis) next to the enforce-mode API; log in as the cashier on both and request the same endpoints; production-mode probe for contrast (API-DIAGNOSTIC)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
