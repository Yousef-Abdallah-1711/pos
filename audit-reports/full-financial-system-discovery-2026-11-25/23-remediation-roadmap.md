# 23 — Dependency-ordered remediation roadmap (direction only; nothing implemented)

Phase 0 — Decisions (no code): approve accounting policies at minimum TAX-01 (rate source/ordering), EARN-01/CLOSE-01, revenue recognition and payment timing, INV-01/inventory capitalisation, PAY-01, ROUND-01 ([32](32-accounting-policy-decision-register.md)). Define roles (cashier, accountant, inventory, HR).
Phase 1 — Foundations: (a) tenant bootstrap: COA, posting rules, fiscal period, tax config created atomically at registration (F-001); (b) one server-authoritative pricing/tax API consumed by POS and invoices (F-004, F-007, F-024); (c) payment/tender model for orders incl. card/wallet recording and clearing (F-005, F-026); (d) idempotency for order and payment creation (F-017).
Phase 2 — Integrity: customer settlement serialisation (F-006); fiscal-period lifecycle and enforcement incl. voids (F-013); payroll approval/accrual (F-015); inventory receiving and adjustment postings (F-011, F-021); invoice posting on issue not draft (F-008).
Phase 3 — Reporting: fix Balance Sheet earnings (F-002), Cash Flow on cash accounts (F-003), VAT from GL (F-010), inventory valuation as-of, `products.cost` in advanced reports (F-009); single definitions and labels for revenue/profit/cash (F-012).
Phase 4 — Platform hardening: RLS expansion and immutability for orders/payments (F-019); Redis timeouts and dependency-aware health (F-018); sanitised errors (F-025); role seeds (F-020); UI contract fixes (F-014, F-016, F-023); storefront setup and menu defect (F-022).
Phase 5 — Verification: rerun this audit's suite (tests/runtime + oracle) on a production-equivalent environment (`ENV=production`), complete the unexecuted tests in [45](45-final-coverage-and-unexecuted-tests.md), run the 36 historical regressions, and the operating-month scenario.

---
## RUN-003 addendum — sequenced waves covering F-001 … F-055 (direction only; nothing implemented)
Detailed per-finding phases and tasks: [57](57-future-remediation-task-register.csv) (all NOT STARTED), specs: [58](58-remediation-specification-index.md), regression tests: [59](59-post-remediation-regression-test-register.csv), policy decisions: [56](56-accounting-policy-decision-package.md).

| Wave | Goal | Findings | Gate to start | Exit criterion |
|---|---|---|---|---|
| 0 | Decisions (no code) | POL-01…POL-15 in 56 | owner/accountant availability | decisions recorded; P0 tasks in 57 become READY |
| 1 | Unblock workflows | F-001 (bootstrap), F-014 (employee create), F-022 (public menu API null-category; storefront needs wildcard CORS in non-prod), F-013 (period lifecycle) | POL-05 | fresh tenant posts without SQL fixtures; employee created via UI; storefront renders; periods closable |
| 2 | Money integrity | F-006 + F-052 (serialize receipts, idempotency for money POSTs), F-017, F-028, F-049/F-050 (refund workflow), F-005/F-026 (tender model), F-007 + F-041 (discount-first tax in both triggers), F-008 (draft posting), F-015 (payroll approval/accrual), F-029 (bank opening journal), F-040/F-034/F-011 (inventory accounting), F-047 | POL-02, 03, 06, 07, 08, 09, 10, 13 | oracle reconciliation zero for each; concurrency tests pass |
| 3 | Reporting alignment | F-002, F-003, F-010, F-009, F-012, F-042, F-043, F-045, F-046, F-004, F-023 | POL-01, 04, 11, 12 | metric register 30c: every FAIL row PASS against the GL oracle |
| 4 | Security and platform hardening | F-019, F-018, F-020/F-033/F-027, F-025/F-031/F-050 (error mapping), F-035, F-036, F-032, F-055, F-030 | POL-15 | negative authorization matrix passes; no raw DB text in any 4xx/5xx |
| 5 | Validation and UX | F-037, F-038, F-053, F-054, F-048, F-016, F-021, F-024, F-039, F-044 | POL-14 | boundary/negative matrix passes |
| 6 | Verification | re-run entire audit suite on `ENV=production`-equivalent stack with TLS (cookie Secure flag, SEC-013), 100-VU load and 50k-line report benchmark (not run in RUN-003) | all waves | release gate in 24 re-evaluated |

Unlock order for blocked audit workflows: F-001 → F-014 → F-022 → F-013 → (F-005/F-026, F-049) → full operating-month rerun (42).
