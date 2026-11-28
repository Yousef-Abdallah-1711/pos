# F-035 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Standards deviations: login limit 20/min (standard 5/min); CSP allows unsafe-eval and unsafe-inline.

## 2. Current behavior (verified)
429 after 18-20 attempts; CSP default includes script-src unsafe-inline unsafe-eval RUN-005: 25 bad logins -> 19 x 401 then 429 from the 20th attempt; the dashboard Content-Security-Policy (frontend/apps/dashboard/next.config.js:155) contains unsafe-inline and unsafe-eval in the isolated build; the API /health response carries no CSP header and the landing response carried none in this run.

## 3. Required behavior
Per CODING_STANDARDS 7.3/7.4

## 4. Reproduction conditions
25 rapid bad logins; response headers

## 5. Verified evidence
- evidence/runtime/31
- evidence/network/RUN-005/PH-17/AUD-R5-OBSERVED-REPRO.json
- tests/runtime/ph17_r5_observed_repro.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/middleware/cache.go:312-314 builds a CSP with script-src unsafe-inline and unsafe-eval; login limit is configured in middleware/rate_limit.go (observed 20/min, 14 failed attempts then 429 in SEC-012c)

## 7. Affected paths
- Backend/database: `backend/internal/middleware/cache.go`:312-314, `backend/internal/middleware/rate_limit.go`:32, `backend/internal/middleware/rate_limit.go`
- Frontend: none
- API contracts: unchanged
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: n/a
- Operational: Auth; Platform workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: tighten the Content-Security-Policy built at backend/internal/middleware/cache.go:312-314 (remove unsafe-eval; replace unsafe-inline with nonces or hashes generated per response for scripts and styles that need them) and set the login rate limit in middleware/rate_limit.go to the documented standard (5 per minute per IP+account) with a lockout response; keep the Redis-backed store (see F-018).

## 10. Alternative designs
A: nonce-based CSP + standard limit (recommended). B: report-only CSP first, then enforce (safe rollout, recommended as step 1).

## 11. Required policy decisions
Security standard (CODING_STANDARDS.md) defines the targets.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/middleware/cache.go, backend/internal/middleware/rate_limit.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Header and limit tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Response CSP contains no unsafe-eval and no unsafe-inline for scripts; the dashboard, landing and storefront still render in Chrome without CSP violations.
- The sixth failed login in a minute from one IP returns 429.
