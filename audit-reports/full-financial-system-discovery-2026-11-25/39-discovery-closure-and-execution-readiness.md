# 39 — Discovery closure and execution readiness (Master Prompt 02)

Date 2026-11-25. Read-only phase: no product code, git, docker, DB or tenant change; no form submitted; no transaction. Baseline re-verified: local HEAD = `133260db66b1143c254d79476daa7650b41a33f0` (= upstream `main`); `git status` shows only the same 6 tooling-modified files plus untracked audit files; all 6 dev containers still up.

## Verdicts
- **DISCOVERY AND TEST-PLANNING READINESS: PARTIAL.** Test-planning artefacts exist for every inventoried item, but exhaustive discovery is **not** achieved (see limitations).
- **FINANCIAL RELEASE READINESS: NO-GO** until execution evidence says otherwise.

## What this phase produced (all in this folder)
| File | Content | Quality |
|---|---|---|
| 30 / 30b | 112 metric rows (105 verified-formula + 7 placeholders) | medium: formulas from catalogue; UI labels not verified |
| 31 | 28 independently re-read findings | high for the rows marked CONFIRMED |
| 32 | Policy decision register (25 decisions + 2 cross-cutting, all OPEN) | written by a sub-agent that was stopped; file was complete; content not independently re-audited (recommendations unapproved) |
| 33 | 25 capabilities reclassified | medium |
| 33b | 22 financial concepts × sources | medium (source-reading) |
| 27b | 20 lineage chains with exact breaks | medium |
| 34 | **1,491 test cases**, all NOT EXECUTED | **templated**: see limits |
| 35 / 35a | 24 scenarios (oracle self-consistency 24/24) + oracle review (no errors found) | scenarios file generated from `oracle/scenarios.py` written by a stopped sub-agent; ran clean, assumptions tagged |
| 36 | 44 security/reliability tests | medium |
| 37 | Isolated-environment runbook | high (compose read in full) |
| 38 | 1,491-row traceability matrix | finding links keyword-heuristic |
| 07 / 07b / 07c | 627 field rows; 128 input-bearing files | see 07c |
| 29 | 48 gap dispositions | — |

## Counts (actual, reconciled)
Applications 4 frontends + Go API + PostgreSQL (+Redis); pages/route handlers 119; API routes 374; DB tables 102, FKs 287, triggers 47, RLS tables **20** (now CONFIRMED-BY-READONLY-RUNTIME; KG-SEC-003's "13" was wrong), policies 23, tables with tenant_id 70; posting-matrix events 36; workflows 25; capability checks 25; competing-source concepts 22 (was 12); report metric rows 112; field rows 627 (was reported 200; original file really had 201); form containers 128 files; manual UI test cases 1,491; scenarios 24; security/reliability tests 44; historical issues 36; knowledge gaps 48.

## Verification outcome (31)
28 findings: CONFIRMED-BY-SOURCE 20 (VER-011 also runtime-confirmed), PARTIALLY-CONFIRMED 5, INFERRED 2, UNKNOWN 1, CONTRADICTED 0. **Corrections discovered** (earlier statements that were wrong or incomplete):
1. Paymob handles **SaaS subscription** payments only (VER-010); customer orders have no gateway — manual `payment-status` marking is the only way to pay a card/storefront order.
2. Missing `0 posting rules` alone does not prove first-order failure: defaults exist in code (VER-012); a missing **chart of accounts** is the real risk.
3. Bank accounts/transactions/reconciliation exist (N9 partially implemented); recurring journals exist (N23); FX rates exist (N24); attendance exists but is not linked to payroll (N11).
4. Outbox redrive endpoint exists (`POST /accounting/paid-order-events/{id}/redrive`, module 15 WRITE).
5. "Suppliers.current_balance decremented on payment" could not be re-found (VER-023).
6. RLS table count = 20 (not 13).

## Missing / partial capabilities (33)
Truly missing (clear): N1, N6, N8, N14, N16, N18, N19, N21, N22, N25, N10 (API). Partial/other: N2, N3, N4, N5, N9, N11, N13, N17, N23, N24, N7 (not integrated), N12 (not enforced), N20 (dead/unreachable), N15 UNKNOWN. Five verdicts changed from "missing" (N2, N3, N9, N15, N23).

## Accounting decisions requiring approval (32)
All 25 are OPEN; **none approved**. Owner/accountant must decide at least: VAT source and recoverability, revenue recognition trigger, inventory capitalisation and costing, current-period earnings and closing, payroll accrual, refunds, rounding. Scenarios tag these as PROVISIONAL.

## Knowledge gaps (29)
48 total: REQUIRES-RUNTIME 18, PARTIALLY-RESOLVED 13, REQUIRES-OWNER-DECISION 8, REQUIRES-ENVIRONMENT-CHANGE 6, REQUIRES-ACCOUNTANT-DECISION 2, OUT-OF-SCOPE 1, CLOSED-BY-DISCOVERY 0.

## Limitations — why discovery is not "exhaustive"
1. **Fields**: regex extraction, 427 of 627 rows low-fidelity; landing/super-admin almost certainly under-counted; true form count unknown (07c).
2. **Metrics**: KDS and Super Admin formulas unread; UI labels not checked (30b).
3. **Test cases (34)**: 1,491 rows are generated from templates per category; only the scenario, financial-operation, metric and historical rows carry specific expectations. Buttons, filters, search, pagination, status transitions and exports are **not enumerated as separate tests**; field tests use generic boundary sets, not exact per-field values. API tests (374) are permission-boundary tests, not functional tests. Dependency ordering is stated below but individual cases are not sequenced.
4. **Traceability (38)**: finding links derive from keyword rules (459 rows linked); database column and workflow links are blank for field/page rows; workflows (25) are not individually linked to test cases; module→page→field chain is only partial.
5. **Verification (31)**: 28 findings re-read, not every claim in 00 §4/12/26/27; items marked PARTIAL/INFERRED/UNKNOWN stay open.
6. **Capability matrix (33)**: backend routes checked through the manifest/grep; frontend only by grep.
7. Nothing was run; no runtime claim beyond read-only catalog queries.
8. Two deliverables (32 and `oracle/scenarios.py`) were produced by sub-agents that were stopped mid-run; I reviewed their structure and ran the scenarios but did not re-derive every scenario by hand.

## Execution prerequisites and order
Prerequisites: owner authorization; ~RAM decision (runbook §6); approval list in runbook §12; policy decisions or acceptance that scenarios run under PROVISIONAL assumptions; accounting initialization approach for new tenants (no seeding at registration).
Order (dependency-based): 1 environment verification (SEC-001) → 2 tenant registration → 3 accounting initialization (COA; gate for everything financial; if it fails, run diagnostic fixtures with an imported COA to keep testing downstream while recording the defect) → 4 roles/permissions → 5 opening balance & capital (manual journals) → 6 suppliers & inventory → 7 products → 8 POS cash order then card/manual-paid → 9 storefront → 10 invoices & receipts → 11 HR/payroll → 12 expenses/bank → 13 reports → 14 reconciliation → 15 negative/boundary → 16 concurrency → 17 recovery → 18 historical regression → 19 final reconciliation → 20 verdict.

## Changes made this phase
Product code: NONE. Environment: NONE (the dev DB was only queried read-only). Financial transactions: ZERO. Only files under `audit-reports/full-financial-system-discovery-2026-11-25/` (+ scratch scripts in the session scratchpad) were created or edited; two files (07 and 29) were extended in place with the original backed up (07) / appended (29).
