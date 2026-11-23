# CFR-015: Product-wide financial tenant/RBAC isolation not independently proven

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-015 — Product-wide financial tenant/RBAC isolation not independently proven

- **Severity/status:** P2 — BLOCKED for runtime negative authorization matrix; source review is incomplete, not proof of a bypass.
- **Problem / impact:** Any missing tenant/restaurant predicate, weak role gate, export/background-worker/cache scope can expose or mutate another restaurant's financial data.
- **Current / expected:** Main contains tenant-aware SQL/RLS work and finance route/permission wrappers, but this assessment has no full route-to-query matrix. Expected: trusted identity derives tenant/restaurant; every object and job/export is scope-bound; least-privilege role enforcement.
- **Root cause:** Not established as one defect. Review candidates: accounting routes in `backend/cmd/api/main.go`, financial handlers/repositories, middleware and RLS migrations `155+`.
- **Reproduction/evidence:** No cross-tenant PG/API/browser tests run. Historical Phase 2A tenant artifacts are not financial isolation proof.
- **Baseline scope:** Main and current: completeness unverified; no confirmed current cross-tenant exploit asserted. Local branch modifications mainly concern auth/tenant groundwork but are not financial-specific evidence.
- **Recommended solution / components:** Build route × role × tenant × restaurant × job/export matrix; combine application predicates, scoped transactions and RLS where viable; test cache/key/background contexts.
- **Migration:** RLS/policies/constraints only after coverage audit and safe migration plan.
- **Dependencies/policy:** Role matrix and Super Admin/auditor access approval.
- **Acceptance:** Positive and negative tests for each finance operation across tenant/restaurant; no unauthorized read/write/export; background jobs and cache isolation; RLS bypass-role tests.
- **Checklist:** BLOCKED.

## 5. Additional issues beyond the canonical 15



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
