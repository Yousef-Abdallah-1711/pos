# CFR-001: Auto-posting silently forces unbalanced lines to balance

Status: IMPLEMENTATION MERGED OR OPEN — ACCEPTANCE PENDING; verify scope and PRs. Original severity: P1. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-001 — Auto-posting silently forces unbalanced lines to balance

- **Severity/status:** P1 — SOURCE-CONFIRMED (source behavior); defect was not re-executed through an order-posting regression this session.
- **Problem / impact:** A material tax, discount, quantity or mapping error can be concealed by changing the final accounting line; this can misstate revenue, VAT, expense or cash while producing a balanced journal.
- **Current / expected:** `AutoPostingEngine.createAndPostEntry` passes input through `normalizeBalancedLines`; the latter adds any imbalance over 0.009 to the last line. Expected: reject material differences; only a policy-approved immaterial residual may go to an explicit rounding account within a defined currency tolerance.
- **Root cause:** `backend/internal/usecase/auto_posting_engine.go:522-553,633-663`.
- **Reproduction/evidence:** Source-confirmed by `math.Abs(diff) <= 0.009` early return followed by unconditional last-line mutation. Test evidence this session: related auto-post unit tests are included in the 33 passing Go cases, but none asserts rejection of a material imbalance; therefore no runtime reproduction claim.
- **Baseline scope:** Present on `origin/main`; no financial source delta on current branch. Local Phase 00 diagnostics/guardrails do not establish root-cause correction.
- **Recommended solution / components:** Replace implicit plug with explicit validation and an approved rounding-account policy; preserve source business amounts and expose posting failure/alert. Use shared money type and currency precision (FIN-007). Components: order, payroll, supplier, depreciation posting use cases; journal repository; monitoring.
- **Migration:** Only if adding controlled rounding-account mapping/config; no migration until accountant approves policy. Do not add a balancing plug.
- **Dependencies/policy:** Currency, permitted residual maximum, rounding mode and account require owner/accountant approval.
- **Acceptance:** 5.00 mismatch rejected and no journal/GL row persists; debit-heavy and credit-heavy cases; approved one-minor-unit case posts to explicit rounding account; currency-specific boundary/property tests; failure visible and retryable.
- **Checklist:** OPEN.



## Current evidence and open questions

At main `afda1f91c67922a92f66507adc208e6e1943d631`,
`createAndPostEntryTx` calls `validateAutoPostingLines` before either
`CreateAndPostAutoEntry` or `CreateAndPostAutoEntryTx`. PR
[#29](https://github.com/Yousef-Abdallah-1711/pos/pull/29) merged as
`57c57373ddd1f0b335369ab65717a036695bddcb`. Main contains
`backend/internal/usecase/auto_posting_validation_test.go`, including tests
named `TestValidateAutoPostingLinesRejectsBothDirectionsOfImbalance` and
`TestAutoPostingRejectsUnbalancedEntriesBeforeJournalOrLedgerWrites`.
These observations support **FIXED_IN_MAIN_UNVERIFIED** for the original
automatic balancing behavior only: the tests have not been rerun against this
exact main SHA in this session, and the full original acceptance, including
approved rounding policy and visible recovery, is not established. Verify the
validator's boundary behavior and run the tests on an exact-main checkout;
also prove persistence and retry behavior with disposable PostgreSQL.

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
