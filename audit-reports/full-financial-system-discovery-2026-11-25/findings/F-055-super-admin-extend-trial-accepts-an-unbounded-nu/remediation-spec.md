# F-055 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Super Admin extend-trial accepts an unbounded number of days (100,000 days moves the trial end to the year 2300).

## 2. Current behavior (verified)
200 days_added 100000, new_trial_end_date 2300-10-02 (negative days correctly rejected) | RUN-003: extend-trial 100000 days accepted -> 2300-10-02 (REPRODUCED; tenant 8 restored by fixture). 

## 3. Required behavior
Upper bound / confirmation and audit reason

## 4. Reproduction conditions
POST /super-admin/tenants/{id}/extend-trial {days:100000}

## 5. Verified evidence
- evidence/database/RUN-003/PH-13/AUD-SAAS-MATRIX.json
- tests/runtime/ph13_saas.mjs

## 6. Root cause and confidence
CONFIRMED — ExtendTenantTrial validates only days > 0 (subscription_handler.go:128-138); there is no upper bound

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/subscription_handler.go`:137, `backend/internal/handler/http/subscription_handler.go`:128-138
- Frontend: `frontend/apps/super-admin`
- API contracts: POST /api/v1/super-admin/tenants/{id}/extend-trial: reason required; days capped
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: Unbounded free access; no billing record
- Operational: Super Admin tenant management; SaaS access control workflows behave as in section 2.
- Security: Yes — see category: Validation gap
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: bound and record trial extensions. In backend/internal/handler/http/subscription_handler.go ExtendTenantTrial (:128-138) enforce 1 <= days <= 90 (configurable), require a reason, and write a subscription history row (actor, days, previous end, new end); show the history in Super Admin.

## 10. Alternative designs
A: cap + reason + history (recommended). B: cap only.

## 11. Required policy decisions
Maximum extension is a business decision (proposed 90).

## 12. Required data migrations
Possibly a subscription_history table if none exists.

## 13. Historical data implications
Tenants with trial ends beyond the cap are listed.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
Keep RBAC module/level checks unchanged; add negative tests for roles without permission.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/subscription_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
days above the cap rejected; audit row created (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- days 100000 returns 400.
- A valid extension records actor, reason and dates and appears in the tenant history.
