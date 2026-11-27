# F-071 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
ENV=staging is treated as development: RBAC defaults to audit-only (cashier reaches finance, HR and user endpoints), mock payments, open metrics and in-memory sessions are allowed.

## 2. Current behavior (verified)
Enforce API: cashier gets 403 on /accounting/journal-entries, /accounting/accounts, /hr/salaries, balance sheet, profit-loss, /users. Staging API (RBAC unset): the same requests return 200 and a cashier write reaches the validator (400, not 403). ENV=production refuses to start without PAYMOB_MODE=live; staging accepts mock payments. /health/metrics answers 200 without a token and the session cookie has no Secure flag in both non-production runs. Staging does require independent signing secrets (signing_secrets.go:9) - the only guard staging shares with production

## 3. Required behavior
A staging deployment enforces the same guards as production (RBAC enforce, no mock payments, protected metrics, Secure cookies), or the weaker behaviour is explicitly documented

## 4. Reproduction conditions
Start the isolated API image with ENV=staging and RBAC_ENFORCEMENT unset (same DB/Redis) next to the enforce-mode API; log in as the cashier on both and request the same endpoints; production-mode probe for contrast (API-DIAGNOSTIC)

## 5. Verified evidence
- evidence/network/RUN-004/PH-17/AUD-R4-SEC013.json
- evidence/network/RUN-004/PH-17/AUD-R4-SEC013-production-probe.log
- tests/runtime/ph17_sec013.mjs

## 6. Root cause and confidence
CONFIRMED — config.IsProductionMode() is true only for production/prod (runtime_env.go:20-23); RBAC default mode is audit unless production (middleware/permission.go:36-42); mock Paymob allowed unless production (config/paymob_mode.go:29); metrics token and Redis/session requirements and secure cookies are gated the same way (health_routes.go:52, main.go:148-156, auth_handler.go:223)

## 7. Affected paths
- Backend/database: `backend/internal/config/runtime_env.go`:20, `backend/internal/middleware/permission.go`:36, `backend/internal/config/paymob_mode.go`:29, `backend/cmd/api/health_routes.go`:52, `backend/cmd/api/main.go`:705, `backend/cmd/api/signing_secrets.go`:9, `backend/internal/config/runtime_env.go`:20-23, `backend/internal/middleware/permission.go`:36-42, `backend/cmd/api/main.go`:148-156, `backend/internal/handler/http/auth_handler.go`:223
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: /health/metrics requires a token outside development
- Accounting logic / posting rules: none named
- Reports affected: Platform configuration (RBAC, payments, sessions, metrics, cookies)

## 8. Impact
- Financial: 0 direct; exposes all financial data and writes to any signed-in role on a staging deployment that holds real or production-like data
- Operational: Platform configuration (RBAC, payments, sessions, metrics, cookies) workflows behave as in section 2.
- Security: Yes — see category: Security / configuration
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: invert the environment check. Replace config.IsProductionMode() (backend/internal/config/runtime_env.go:20-23) usage with an explicit IsDevelopmentMode() allow-list (development, dev, local, test) so that staging and any unknown value behave like production for: RBAC default and refusing RBAC_ENFORCEMENT=audit (middleware/permission.go:36-52), mock payment refusal (config/paymob_mode.go:29), Redis/session and rate-limit store requirements (cmd/api/main.go:148-160), secure cookies (handler/http/auth_handler.go:223) and the metrics token (cmd/api/health_routes.go:52); fail startup on an unrecognised ENV value; keep the independent-secret check (cmd/api/signing_secrets.go:9).

## 10. Alternative designs
A: development allow-list (recommended; unknown values become safe). B: add staging to IsProductionMode (fixes only the named value).

## 11. Required policy decisions
List of environments considered development (owner decision).

## 12. Required data migrations
No.

## 13. Historical data implications
Review staging data for exposure while RBAC was in audit mode.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
Core of the fix: enforcement mode and role matrix tests.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-020, F-033 and shared-file findings F-018, F-020, F-025, F-030, F-031, F-033, F-036, F-037

## 21. Risks
- Regression in modules sharing the same source files (backend/cmd/api/health_routes.go, backend/cmd/api/main.go, backend/cmd/api/signing_secrets.go, backend/internal/config/paymob_mode.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Start with ENV=staging and no RBAC variable: cashier requests to finance endpoints are 403; mock payment mode is rejected; metrics require a token (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Starting with ENV=staging and no RBAC variable enforces RBAC: a cashier gets 403 on finance, HR and user endpoints (re-run tests/runtime/ph17_sec013.mjs).
- ENV=staging with PAYMOB_MODE unset or mock fails to start.
- /health/metrics returns 401 without a token outside development.
- An unrecognised ENV value stops the API with a clear message.
