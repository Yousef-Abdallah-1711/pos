# F-070 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
KDS busiest-hours and "today" boundaries are computed in UTC, not the restaurant local time.

## 2. Current behavior (verified)
API hours equal the UTC buckets exactly (1:7, 17:4, 19:4, 20:18, 22:11); the Cairo buckets differ (0:11, 3:7, 19:4, 21:4, 22:18); kds_queue.received_at is timestamp without time zone; no tenant time-zone setting was found

## 3. Required behavior
Hours and the "today" cut-off follow the restaurant time zone (Egypt UTC+2/+3)

## 4. Reproduction conditions
Tenant 4: compare /kds/analytics busiest_hours with SQL bucketed by UTC and by Africa/Cairo for the same kds_queue rows

## 5. Verified evidence
- evidence/network/RUN-004/PH-12/AUD-R4-KDS-after-read.json
- evidence/network/RUN-004/PH-12/AUD-R4-KDS.json
- evidence/screenshots/RUN-004/PH-12/AUD-R4-KDS-after-bump.png

## 6. Root cause and confidence
CONFIRMED — kds_repo.go:418-430 EXTRACT(HOUR FROM received_at) and CURRENT_DATE at :314-317; DB timezone is UTC

## 7. Affected paths
- Backend/database: `backend/internal/repository/kds_repo.go`:418, `backend/internal/repository/kds_repo.go`:314, `backend/internal/repository/kds_repo.go`:418-430
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/kds`
- API contracts: settings gain timezone (additive)
- Accounting logic / posting rules: none named
- Reports affected: KDS analytics and stats

## 8. Impact
- Financial: 0 (scheduling signal shifted 2-3 hours); other date-bounded reports may share the cut-off issue (CANDIDATE, not tested)
- Operational: KDS analytics and stats workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: introduce a restaurant time zone. Add a tenant/restaurant timezone setting (default Africa/Cairo) and use it for KDS "today" and busiest-hours buckets (backend/internal/repository/kds_repo.go:314-317 and :418-430: convert with AT TIME ZONE before CURRENT_DATE / EXTRACT(HOUR ...)); store timestamps as timestamptz going forward; audit other date-bounded reports (sales, dashboard, orders stats) for the same UTC cut-off (currently a CANDIDATE, not tested).

## 10. Alternative designs
A: time-zone setting + conversion in queries (recommended). B: convert in the frontend only (wrong boundaries in aggregates).

## 11. Required policy decisions
Default time zone per tenant.

## 12. Required data migrations
Yes: setting column; timestamps conversion plan.

## 13. Historical data implications
Historical aggregates are recomputed on read; no stored data changes.

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
F-044 and shared-file findings F-044

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/repository/kds_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Orders at 23:30 and 00:30 local land in the right day and hour (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- An order received at 23:30 local appears in the previous local day and in hour 23 (not hour 21).
- KDS completed_today uses local midnight.
- The audit of other reports lists every query that uses CURRENT_DATE or date_trunc on a UTC timestamp.
