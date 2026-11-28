# Phase 16 — Negative, boundary and concurrency testing

**Closure status:** DERIVED FROM REGISTER 51 (RUN-004): 1 EXECUTED-PARTIAL; legacy label from RUN-003: COMPLETE-WITH-LIMITS

## Objective and scope
See the master prompt section for Phase 16 (`writing-block.md`) — this file only references existing artefacts.

## Authoritative reports
- [18-negative-boundary-test-results.md](../../18-negative-boundary-test-results.md)
- [35-financial-scenario-test-specification.md](../../35-financial-scenario-test-specification.md)

## Inventories
file 36 (see the audit root README for names).

## Existing test scripts
- `tests/runtime/ph16_masterdata.mjs`
- `tests/runtime/ph16_negprice.mjs`
- `tests/runtime/ph16_tables.mjs`

## Evidence
- `evidence/database/RUN-003/PH-16`
- `evidence/network/RUN-003/PH-16`
- `evidence/screenshots/RUN-003/PH-16`

## Existing test IDs (34b)
EX-PG-dash-owner-001, EX-PG-dash-owner-002, EX-PG-dash-owner-003, EX-PG-dash-owner-004, EX-PG-dash-owner-005, EX-PG-dash-owner-006, EX-PG-dash-owner-007, EX-PG-dash-owner-008, EX-PG-dash-owner-009, EX-PG-dash-owner-010, EX-PG-dash-owner-011, EX-PG-dash-owner-012, EX-PG-dash-owner-013, EX-PG-dash-owner-014, EX-PG-dash-owner-015, EX-PG-dash-owner-016, EX-PG-dash-owner-017, EX-PG-dash-owner-018, EX-PG-dash-owner-019, EX-PG-dash-owner-020, EX-PG-dash-owner-021, EX-PG-dash-owner-022, EX-PG-dash-owner-023, EX-PG-dash-owner-024, EX-PG-dash-owner-025, EX-PG-dash-owner-026, EX-PG-dash-owner-027, EX-PG-dash-owner-028, EX-PG-dash-owner-029, EX-PG-dash-owner-030, EX-PG-dash-owner-031, EX-PG-dash-owner-032, EX-PG-dash-owner-033, EX-PG-dash-owner-034, EX-PG-dash-owner-035, EX-PG-dash-owner-036, EX-PG-dash-owner-037, EX-PG-dash-owner-038, EX-PG-dash-owner-039, EX-PG-dash-owner-040, EX-PG-dash-owner-041, EX-PG-dash-owner-042, EX-PG-dash-owner-043, EX-PG-dash-owner-044, EX-PG-dash-owner-045, EX-PG-dash-owner-046, EX-PG-dash-owner-047, EX-PG-dash-owner-048, EX-PG-dash-owner-049, EX-PG-dash-owner-050, EX-PG-dash-owner-051, EX-PG-dash-owner-052, EX-PG-dash-owner-053, EX-PG-dash-owner-054, EX-PG-dash-owner-055, EX-PG-dash-owner-056, EX-PG-dash-owner-057, EX-PG-dash-owner-058, EX-PG-dash-owner-059, EX-PG-dash-owner-060 …

## Actual results
- Receipt/supplier/stock/redrive/recurring/mark-paid concurrency; validation matrices for journals, accounts, HR, FX.
- Limit: no k6 100-VU profile (see Phase 18).

## Findings
- [F-006](../../findings/F-006-concurrent-customer-receipts-over-settle-an-invo/finding.md) — Concurrent customer receipts over-settle an invoice (High)
- [F-017](../../findings/F-017-duplicate-pos-submissions-create-duplicate-order/finding.md) — Duplicate POS submissions create duplicate orders; Idempotency-Key is ignored on order cre (Medium)
- [F-047](../../findings/F-047-recurring-journal-template-can-be-run-repeatedly/finding.md) — Recurring journal template can be run repeatedly for the same period and posts a duplicate (Medium)
- [F-052](../../findings/F-052-customer-receipts-and-supplier-payments-ignore-i/finding.md) — Customer receipts and supplier payments ignore Idempotency-Key: a client retry records the (High)

## Remaining blockers and unexecuted work
Recorded per task in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and in [48-blocked-vs-unexecuted-and-prerequisites.md](../../48-blocked-vs-unexecuted-and-prerequisites.md).

## Required new tests (if the phase is re-opened)
Any row of 51 with status PLANNED/BLOCKED for this phase; plus the regression tests in [59](../../59-post-remediation-regression-test-register.csv) after remediation.

## Required deliverables
Evidence under `evidence/{screenshots,network,database}/RUN-003/PH-16/`, updated registers (22, 34b, 43, 44), progress.md entry.

## Completion criteria
Every executable task in 51 for this phase is EXECUTED-* or has a documented blocker; findings recorded in 22 with evidence; no product code changed.
