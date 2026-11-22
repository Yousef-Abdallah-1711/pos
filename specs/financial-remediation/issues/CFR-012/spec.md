# CFR-012: Cash Flow report uses account classification heuristics, not demonstrated cash movement policy

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-012 — Cash Flow report uses account classification heuristics, not demonstrated cash movement policy

- **Severity/status:** P2 — SOURCE-CONFIRMED heuristic; report accuracy not independently reproduced.
- **Problem / impact:** Operating/investing/financing classification may misclassify noncash entries and accrual movements as cash flows, yielding misleading statement.
- **Current / expected:** `GenerateCashFlow` classifies using account types/code patterns, including equity or code prefix heuristics. Expected: statement derived from cash/bank movements and approved direct/indirect method with noncash adjustments.
- **Root cause:** `backend/internal/usecase/financial_report_usecase.go:340 onward` (classification SQL around 346+; code-prefix filter in same function). Exact source at fetched SHA.
- **Reproduction/evidence:** Static SQL review; no independent cash-flow fixture.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Adopt accountant-approved direct or indirect method; identify cash equivalents and map activities; reconcile opening + net movement = closing cash.
- **Migration:** Optional cash-flow classification metadata/account mapping.
- **Dependencies/policy:** Cash-equivalent scope and method need accountant approval.
- **Acceptance:** Independently calculated fixture with noncash depreciation, working-capital movements, asset purchase and financing; closing cash ties to GL/bank; comparative export parity.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
