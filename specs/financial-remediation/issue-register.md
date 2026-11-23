# Financial issue register — recovery draft

Source: original tracker at `fa2e301`; main observed 2026-11-23 `3fa3702148771f39af5c9f7f25f377ecdd3f0a98`. Statuses are provisional and must be reconciled to dated evidence.

| ID | Original title | Severity | Current disposition | PRs |
|---|---|---|---|---|
| [CFR-001](issues/CFR-001/spec.md) | Auto-posting silently forces unbalanced lines to balance | P1 | FIXED_IN_MAIN_UNVERIFIED for automatic plug; full acceptance pending | 29 |
| [CFR-002](issues/CFR-002/spec.md) | Journal integrity relies too heavily on application checks | P1 | PARTIALLY_FIXED — migration merged; current acceptance and reversal policy pending | 34 |
| [CFR-003](issues/CFR-003/spec.md) | Concurrent journal posting, balance update and numbering races | P1 | IMPLEMENTATION MERGED OR OPEN — ACCEPTANCE PENDING; verify scope and PRs | 41 |
| [CFR-004](issues/CFR-004/spec.md) | Journal draft creation and GL posting have a recoverable partial-failure window | P1 | PARTIALLY IMPLEMENTED — atomic posting and paid-order outbox merged; current restart and operational acceptance pending | 33, 35, 37 |
| [CFR-005](issues/CFR-005/spec.md) | Paid orders can suppress posting failures and have no proven cancellation/refund reversal | P1 | PARTIALLY IMPLEMENTED — paid intent and unsupported-refund guard merged; refund lifecycle and acceptance pending | 35, 39 |
| [CFR-006](issues/CFR-006/spec.md) | Historical report performance symptom not remeasured | P2 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [CFR-007](issues/CFR-007/spec.md) | Binary floating point is used across monetary domain and APIs | P2 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [CFR-008](issues/CFR-008/spec.md) | Order tax hardcoded; invoice tax computed before discount | P2 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [CFR-009](issues/CFR-009/spec.md) | Paid payroll remains editable/deletable and duplicate mark-paid guard is unproven | P1 | PARTIALLY IMPLEMENTED — paid guards and atomic posting merged; correction and current acceptance pending | 40 |
| [CFR-010](issues/CFR-010/spec.md) | Balance Sheet omits current-period income from equity; mainline has no reliable warning/export gate | P1 | POLICY_PENDING / REVIEW_PENDING; synthetic equation unreconciled | 42, 46, 104, 108, 109, 110, 112, 114 |
| [CFR-011](issues/CFR-011/spec.md) | Automatic supplier entry type exceeds journal schema width | P2 | SCHEMA FIX MERGED — constructed-entry PostgreSQL tests pass; supported supplier lifecycle and acceptance pending | 44 |
| [CFR-012](issues/CFR-012/spec.md) | Cash Flow report uses account classification heuristics, not demonstrated cash movement policy | P2 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [CFR-013](issues/CFR-013/spec.md) | Voiding a journal writes a reversal dated today without adequate closed-period handling | P2 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [CFR-014](issues/CFR-014/spec.md) | Inventory valuation, COGS timing and refunds/adjustments are not proven to reconcile | P2 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [CFR-015](issues/CFR-015/spec.md) | Product-wide financial tenant/RBAC isolation not independently proven | P2 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [FIN-016](issues/FIN-016/spec.md) | No complete gift-card/deposit/deferred-revenue lifecycle established | P3 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [FIN-017](issues/FIN-017/spec.md) | Fixed-asset acquisition, depreciation and disposal lifecycle not verified end to end | P3 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [FIN-018](issues/FIN-018/spec.md) | No complete fiscal close and retained-earnings roll-forward established | P2 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [FIN-019](issues/FIN-019/spec.md) | Export parity and native XLSX / CSV safety are unverified | P3 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [FIN-020](issues/FIN-020/spec.md) | Multi-currency / FX and consolidated financial reporting policy is not established | P3 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |
| [FIN-021](issues/FIN-021/spec.md) | Root HTML language/direction does not reflect the active locale | P3 | HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION | unmapped |

## CFR-010 engineering review — 2026-11-23

PR [#109](https://github.com/Yousef-Abdallah-1711/pos/pull/109) head `244c44c8f7c53fc81adbd83453edf2875acce9fb` corrects fiscal-period diagnostic classification when a closed period overlaps another period or a covering period has an unrecognized status. Two disposable PostgreSQL regression scenarios were added. The local attempt failed before fixture creation because no PostgreSQL server was listening; the [financial PostgreSQL CI job](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663) passed, while its full workflow was still in progress at this snapshot. Main advanced afterward, so update the candidate and rerun applicable checks. This diagnostic work does not approve earnings presentation or close CFR-010. Independent review and accounting decisions remain outstanding.

The linked PostgreSQL job log confirms PASS for the closed-only, overlapping-period and unknown-status diagnostic tests. It also records PASS for CFR-004 duplicate posting, tenant/restaurant scope and all three injected failure boundaries; four CFR-005 payment transition/scope/race tests; and seven CFR-009 payroll scenarios. These results are for the #109 candidate, not final main. Outbox retry, lease recovery, missing-intent reconciliation and paid-status/intent rollback tests also passed in that job. Remaining issue acceptance is described in each task list.

On recovery checkout `8667c036a5c2c3e12b8a7d953d25fd5553f15bf9`, `go test ./internal/domain -run '^TestCanTransitionOrderPaymentStatus$' -count=1 -v` passed all eight scenarios. It validates the CFR-005 transition function, not provider refund or database acceptance.

PR [#110](https://github.com/Yousef-Abdallah-1711/pos/pull/110) head `605c9aaa46cc71e7395f9c2b332e5d9a6177ef86` changes the production reconciliation query as well as adding benchmark infrastructure. The recorded separate-run large-scope warm baseline p50/p95/max was 174.044/175.441/175.524 ms and candidate was 180.726/182.565/182.660 ms. The line-query plan improved and temporary I/O was eliminated, but overall latency regressed and no performance budget has been approved. Benchmark CI success is not performance acceptance. See the PR's `PERFORMANCE_RESULTS.md` for dataset, environment and limitations.
