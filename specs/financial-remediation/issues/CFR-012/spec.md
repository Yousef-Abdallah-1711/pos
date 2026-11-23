# CFR-012: Cash Flow report uses account classification heuristics, not demonstrated cash movement policy

Status: SOURCE-CONFIRMED / POLICY_PENDING. Original severity: P2. The report methodology is not accepted.

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

At main `3fa3702148771f39af5c9f7f25f377ecdd3f0a98`, `GenerateCashFlow` computes operating activity from revenue/expense and other categories from account-code prefixes and equity. That is not a demonstrated cash-movement method. Open [#114](https://github.com/Yousef-Abdallah-1711/pos/pull/114) changes movement dates for voids but retains the classification formula; it does not approve that formula.

The candidate's synthetic taxed sale test independently observes cash GL +55 and reversal -55 while the existing operating formula gives +50 and -50. Its period and same-day void tests passed in [disposable PostgreSQL](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45077468898/job/137377801133) at candidate `f9ef0156c7df5c850140cefb08d7722bcdd53754`. This demonstrates a difference for that fixture, not a complete cash-flow statement oracle. CF-01 method and cash-equivalent definitions remain pending.

## Definition of done

- Accountant approves direct or indirect method, cash equivalents, treatment of tax, noncash and activity mapping.
- Synthetic sale/refund, accrual, depreciation, working capital, asset purchase and financing fixtures reconcile opening cash plus movements to closing cash and independent GL/bank oracle; voids and period boundaries preserve approved dates.
- Validate tenant and restaurant isolation, performance, API, comparative views and export parity with exact-SHA PostgreSQL evidence and required review. Keep CF-01 pending until signed approval.
