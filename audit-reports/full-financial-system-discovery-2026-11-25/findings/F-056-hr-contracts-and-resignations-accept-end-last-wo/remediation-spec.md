# F-056 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
HR contracts and resignations accept end/last-working dates earlier than their start/resignation dates.

## 2. Current behavior (verified)
both return 201 and the invalid periods are stored (contract id 1, resignation id 1); duplicate contract number returns 400 with raw pq constraint text (F-031 class); attendance/leave/salary periods in the same module family are protected only by DB check constraints that surface as 500 | RUN-003: contract end 2026-10-01 < start 2026-11-01 stored (201); resignation last day before resignation date stored (201) (REPRODUCED). 

## 3. Required behavior
400 with a clear message

## 4. Reproduction conditions
POST /hr/contracts with start_date 2026-11-01 and end_date 2026-10-01; POST /hr/resignations with resignation_date 2026-11-26 and last_working_date 2026-11-01

## 5. Verified evidence
- evidence/database/RUN-003/PH-10/AUD-HR-MODULES.json
- tests/runtime/ph10_hr_modules.mjs

## 6. Root cause and confidence
CONFIRMED — contract create/update validate only the date format (contract_handler.go:165-170, 221-226); no end >= start rule; resignation create has the same gap (resignation_repo.go:117-143)

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/contract_handler.go`, `backend/internal/handler/http/resignation_handler.go`, `backend/internal/handler/http/contract_handler.go`:165, `backend/internal/repository/resignation_repo.go`:117, `backend/internal/handler/http/contract_handler.go`:165-170, `backend/internal/repository/resignation_repo.go`:117-143
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: contracts/resignations create and update: 400 for inverted ranges
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0 directly; wrong notice/tenure calculations downstream
- Operational: HR contracts; HR resignations workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: validate ordered date ranges. In backend/internal/handler/http/contract_handler.go (create :165-170, update :221-226) and the resignation create/update (backend/internal/repository/resignation_repo.go:117-143 and its handler) require end_date >= start_date and last_working_date >= resignation_date, returning 400; add CHECK constraints so direct writers are covered.

## 10. Alternative designs
A: handler validation + CHECK constraints (recommended). B: handler only.

## 11. Required policy decisions
None.

## 12. Required data migrations
CHECK constraints after listing and correcting existing inverted rows (contract 1 and resignation 1 in the evidence).

## 13. Historical data implications
Existing inverted rows are listed for HR correction.

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
F-031, F-053, F-054 and shared-file findings F-069

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/contract_handler.go, backend/internal/handler/http/resignation_handler.go, backend/internal/repository/resignation_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Invalid date ranges rejected with 400 for contracts and resignations (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A contract ending before its start returns 400; a valid one returns 201.
- A resignation with last working day before the resignation date returns 400.
- The CHECK constraints reject the same data written directly.
