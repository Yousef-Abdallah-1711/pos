# F-054 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Leave approval does not enforce entitlement and total_days is not validated against the requested dates.

## 2. Current behavior (verified)
30-day leave approved; report shows used 30, remaining 0 (clamped) vs allowance 21; 3-day range stored total_days 40 | RUN-003: 30-day annual leave approved over allowance 21; total_days 40 stored for a 3-day range (REPRODUCED). 

## 3. Required behavior
Reject or flag leave beyond remaining entitlement; total_days derived from dates

## 4. Reproduction conditions
Create annual leave 2026-12-01..12-30 total_days 30 for an employee with allowance 21 and approve; create a 3-day range with total_days 40

## 5. Verified evidence
- evidence/database/RUN-003/PH-10/AUD-HR-MATRIX.json
- evidence/database/RUN-003/PH-10/AUD-HR-MATRIX2.json
- tests/runtime/ph10_hr_matrix.mjs
- tests/runtime/ph10_hr_matrix2.mjs

## 6. Root cause and confidence
CONFIRMED — total_days is taken from the client when non-zero (hr_leave_handler.go:207-210; update :328-329) and no allowance or date-consistency check exists

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/hr_leave_handler.go`:132-200, `backend/internal/handler/http/hr_leave_handler.go`:207, `backend/internal/handler/http/hr_leave_handler.go`:328, `backend/internal/handler/http/hr_leave_handler.go`:207-210
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: POST/PUT /api/v1/hr/leaves: total_days ignored on input (returned computed)
- Accounting logic / posting rules: none named
- Reports affected: HR leaves; leave-balance report

## 8. Impact
- Financial: Unpaid-leave exposure / payroll deduction not triggered
- Operational: HR leaves; leave-balance report workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: derive and enforce leave entitlement. In backend/internal/handler/http/hr_leave_handler.go (:207-210 create, :328-329 update, approve near :412) ignore client total_days and compute working days from start/end (excluding configured weekends/holidays); at approval compare with the remaining annual allowance and either reject, or approve the excess as unpaid leave according to POL-14, recording the unpaid portion; the leave balances report reads the same ledger (it currently clamps remaining to 0 when used exceeds the allowance, hiding the overdraw).

## 10. Alternative designs
A: server-derived days + balance check at approval (recommended). B: keep client days and only warn.

## 11. Required policy decisions
POL-14 (leave entitlement) in 56.

## 12. Required data migrations
Optional unpaid_days column.

## 13. Historical data implications
Existing over-allowance approvals (30 vs 21 in the evidence) are listed for HR.

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
none and shared-file findings F-079

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/hr_leave_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Over-allowance and mismatched total_days rejected or converted to unpaid (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Requesting 3 calendar days stores total_days 3 regardless of the client value.
- Approving 30 days against a remaining allowance of 21 follows POL-14 (reject, or 21 paid + 9 unpaid recorded).
- The balances report shows the true remaining balance including negative values or the unpaid portion.
