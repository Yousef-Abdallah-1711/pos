# F-071 — ENV=staging is treated as development: RBAC defaults to audit-only (cashier reaches finance, HR and user endpoints), mock payments, open metrics and in-memory sessions are allowed

- **Severity:** High  ·  **Category:** Security / configuration  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Platform configuration (RBAC, payments, sessions, metrics, cookies)
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-SEC013
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A staging deployment enforces the same guards as production (RBAC enforce, no mock payments, protected metrics, Secure cookies), or the weaker behaviour is explicitly documented

Observed: Enforce API: cashier gets 403 on /accounting/journal-entries, /accounting/accounts, /hr/salaries, balance sheet, profit-loss, /users. Staging API (RBAC unset): the same requests return 200 and a cashier write reaches the validator (400, not 403). ENV=production refuses to start without PAYMOB_MODE=live; staging accepts mock payments. /health/metrics answers 200 without a token and the session cookie has no Secure flag in both non-production runs. Staging does require independent signing secrets (signing_secrets.go:9) - the only guard staging shares with production

## Monetary / operational impact
0 direct; exposes all financial data and writes to any signed-in role on a staging deployment that holds real or production-like data

## Where it was reproduced
Start the isolated API image with ENV=staging and RBAC_ENFORCEMENT unset (same DB/Redis) next to the enforce-mode API; log in as the cashier on both and request the same endpoints; production-mode probe for contrast (API-DIAGNOSTIC)

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
