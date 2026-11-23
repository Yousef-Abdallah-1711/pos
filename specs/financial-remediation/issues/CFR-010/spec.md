# CFR-010: Balance Sheet omits current-period income from equity; mainline has no reliable warning/export gate

Status: POLICY_PENDING / REVIEW_PENDING; synthetic equation unreconciled. Original severity: P1. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-010 — Balance Sheet omits current-period income from equity; mainline has no reliable warning/export gate

- **Severity/status:** P1 — SOURCE-CONFIRMED.
- **Problem / impact:** Balance Sheet can fail Assets = Liabilities + Equity whenever period earnings are not closed into retained earnings; report/export may still appear usable.
- **Current / expected:** Main query selects balance-sheet account types only; it does not derive current earnings from revenue less expenses. `EquationBalanced` is calculated, but main has no Phase 00 warning/export guard. Expected: include correctly derived current earnings or use a completed close/retained earnings model and block misleading exports.
- **Root cause:** `backend/internal/usecase/financial_report_usecase.go:261-336`, especially account filter at 288-293 and equation at 336. Mainline is checked at fetched SHA.
- **Reproduction/evidence:** Query and equation source confirms omission path. The prior local Phase 00 candidate's warning/export guard is not on origin/main; no current DB fixture demonstrated a numeric imbalance.
- **Baseline scope:** Defect and absence of guard on origin/main and current branch. A local-only Phase 00 candidate/dirty artifacts have limited guardrails but are not this branch's fetched main behavior and do not calculate missing equity.
- **Recommended solution / components:** Compute as-of net income from posted revenue/expense balances respecting close dates, or enforce close journal and report closed/available state. Preserve equity account presentation and expose equation residual.
- **Migration:** Possible close/retained-earnings fields only if lifecycle design requires them.
- **Dependencies/policy:** Accounting basis, close schedule, prior-period adjustment and report cut-off require accountant approval.
- **Acceptance:** Independent fixtures prove Assets = Liabilities + Equity for open and closed periods; reports/export refuse invalid state; P&L, TB and GL agree.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
