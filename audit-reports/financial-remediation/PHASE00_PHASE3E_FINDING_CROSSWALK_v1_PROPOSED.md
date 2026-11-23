# Phase 00 / Phase 3E finding crosswalk — v1 proposal

**Status:** Proposed; owner approval required. This file does not alter, renumber, close, or supersede either historical register.

## Source provenance

- Phase 00 canonical findings: `audit-reports/financial-phase00-pr-handoff-2026-11-11/04-CANONICAL-FINANCIAL-FINDINGS.md`.
- Phase 00 alias/source map: `audit-reports/financial-phase00-pr-handoff-2026-11-11/05-FINDING-ID-CROSSWALK.md`.
- Phase 3E closure snapshot: `audit-reports/financial-acceptance-phase3b-2026-11-15/PHASE3E_CLOSURE.md`.
- Current-state engineering dispositions: `SERVIO_CURRENT_FINANCIAL_AUDIT_2026-11-16.md`.
- Current decision register: `audit-reports/financial-acceptance-phase3-2026-11-15/ACCOUNTING_POLICY_DECISION_REGISTER.md`.

The registers use overlapping `CFR-*` identifiers for different findings. Any join must include the register version and source path. “Related” below means a reviewer may compare evidence; it does not mean the findings are duplicates or may share a closure disposition.

## Proposed mapping

| Phase 00 alias | Phase 00 meaning | Phase 3E alias / meaning | Proposed relation |
|---|---|---|---|
| `PHASE00:CFR-001` | Balance Sheet omits current-period earnings/equity | `PHASE3E:CFR-010` — Balance Sheet residual and earnings boundary | Related; policy-dependent. Keep separate finding records. |
| `PHASE00:CFR-002` | Manual journal UI omits JSONB attachments and API rejects the request | `PHASE3E:CFR-002` — journal/GL constraints and atomic persistence | Partial overlap. Preserve the UI/API-contract defect separately. |
| `PHASE00:CFR-003` | Auto-post normalizer absorbs material imbalance | `PHASE3E:CFR-001` — unbalanced auto-post validation | Exact semantic match; preserve both source aliases. |
| `PHASE00:CFR-004` | Failed auto-post may leave a draft requiring review | `PHASE3E:CFR-004` — durable event/recovery controls | Related but distinct lifecycle scopes. |
| `PHASE00:CFR-005` | Order tax rate is hard-coded | `PHASE3E:CFR-008` — tax/discount basis and order-to-GL | Partial overlap; do not collapse rate and discount-base questions. |
| `PHASE00:CFR-006` | Duplicate dashboard requests / slow development route | `PHASE3E:CFR-006` — representative financial performance | Related only; frontend request behavior and PostgreSQL financial query performance are separate. |
| `PHASE00:CFR-007` | Missing database debit/credit and posted-entry immutability enforcement | `PHASE3E:CFR-002` — database constraints/atomic posting | Partial overlap; retain constraint-specific acceptance gates. |
| `PHASE00:CFR-008` | Concurrent posting, void, balance, and numbering races | `PHASE3E:CFR-003` — numbering/concurrency | Partial overlap; retain each operation and isolation scope. |
| `PHASE00:CFR-009` | Go monetary `float64` versus database `NUMERIC` | `PHASE3E:CFR-007` — monetary precision and rounding | Strong semantic match; rounding policy remains a separate approval dependency. |
| `PHASE00:CFR-010` | Paid-order accounting reversal/recovery behavior | `PHASE3E:CFR-005` — paid-order transition/correction safety | Related; full refund/reversal policy remains separate. |
| `PHASE00:CFR-011` | Invoice tax does not apply discounts to the tax base | `PHASE3E:CFR-008` — tax/discount basis | Partial overlap; retain the exact invoice-tax scenario under its original alias. |
| `PHASE00:CFR-012` | Paid payroll edit/delete and manual statutory deductions | `PHASE3E:CFR-009` — payroll source-to-GL | Partial overlap; immutability and statutory mapping are separate gates. |
| `PHASE00:CFR-013` | Void and account-deletion edge cases | `PHASE3E:CFR-013` — voids, reversals, and closed-period corrections | Related; account deletion and reversal policies remain distinct. |
| `PHASE00:CFR-014` | Automatic journal insertion fails on nil JSONB attachments | `PHASE3E:CFR-002` — journal/GL persistence constraints | Partial overlap; keep the nil-attachment regression traceable. |
| `PHASE00:CFR-015` | Supplier automatic entry type exceeds the historical schema width | `PHASE3E:CFR-011` — supported supplier automatic-posting contract | Related; schema-width remediation is not approval of the legacy supplier posting contract. |
| — | No equivalent Phase 00 finding located | `PHASE3E:CFR-012` — Cash Flow method/classification | New Phase 3E scope; retain as such. |
| — | No equivalent Phase 00 finding located | `PHASE3E:CFR-014` — inventory/COGS/valuation | New Phase 3E scope; retain as such. |
| — | No equivalent Phase 00 finding located | `PHASE3E:CFR-015` — tenant boundaries and role authorization | New Phase 3E scope; retain as such. |
| — | No equivalent Phase 00 finding located | `PHASE3E:FIN-016` through `FIN-021` — deposits/deferred revenue, fixed assets, close, exports, FX/consolidation, authenticated UI | New Phase 3E scope; retain each identifier and its historical evidence. |

## Register rules proposed for approval

1. Store keys as `(register_version, source_id)`, never as `source_id` alone.
2. Keep original titles, source files, evidence links, severity, and historical status immutable; append dispositions with date and author.
3. Use a separate canonical issue key only after an owner approves the semantic relation and duplicate/partial-overlap classification.
4. Do not infer missing source findings. Phase 00 records a forensic `FIN-007` source row as absent/unmapped; no finding title is supplied for it.
5. Keep the 30-subsystem acceptance scope separate. No authoritative 30-row register was located in the reviewed files; this crosswalk does not define or count that scope.

**Approval needed:** financial/product owner to accept or revise the proposed semantic relations and canonical-key rules before this crosswalk is used to merge, close, reprioritize, or count findings.
