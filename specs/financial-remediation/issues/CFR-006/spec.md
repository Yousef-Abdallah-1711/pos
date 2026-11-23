# CFR-006: Historical report performance symptom not remeasured

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-006 — Historical report performance symptom not remeasured

- **Severity/status:** P2 — BLOCKED (no current performance fixture/measurement).
- **Problem / impact:** Slow or high-memory financial reports can delay close and time out at larger tenant volumes; no current numeric performance claim is made.
- **Current / expected:** Historical audit carried a performance symptom; this session did not benchmark current main. Financial queries scan posted-ledger/report paths; proper bounded pagination, query plans and dataset size have not been validated here.
- **Root cause:** Not established. Candidate query locations: `backend/internal/usecase/financial_report_usecase.go`, `backend/internal/repository/general_ledger_repo.go`, `accounting_advanced_repo.go`.
- **Reproduction/evidence:** No load test or `EXPLAIN (ANALYZE, BUFFERS)` this session. Therefore symptom and root cause remain unverified.
- **Baseline scope:** Historical finding assessed against main; persistence of performance symptom is unknown. Current branch changes are unrelated.
- **Recommended solution / components:** Establish realistic synthetic dataset, baseline p50/p95/query plans and export memory; optimize only measured hotspots.
- **Migration:** Index/partition only after evidence and migration plan.
- **Dependencies/policy:** Workload SLO, data retention and report freshness.
- **Acceptance:** Recorded dataset/command, explain plans, bounded memory, SLO gates; no unsupported scale claim.
- **Checklist:** BLOCKED.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
