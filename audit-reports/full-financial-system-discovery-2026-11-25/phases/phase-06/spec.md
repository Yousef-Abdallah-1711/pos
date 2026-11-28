# Phase 06 — Products, categories, inventory and suppliers

**Closure status:** DERIVED FROM REGISTER 51 (RUN-004): 3 EXECUTED-PARTIAL; legacy label from RUN-003: PARTIAL

## Objective and scope
See the master prompt section for Phase 06 (`writing-block.md`) — this file only references existing artefacts.

## Authoritative reports
- [26-business-workflow-catalog.md](../../26-business-workflow-catalog.md)

## Inventories
file 33 (see the audit root README for names).

## Existing test scripts
- `tests/runtime/ph06_inventory_controlled.mjs`
- `tests/runtime/ph06_r4_form_matrix.mjs`
- `tests/runtime/ph06_r4_probe_missing.mjs`
- `tests/runtime/ph06_r5_f074.mjs`
- `tests/runtime/ph06_r5_field_isolation.mjs`
- `tests/runtime/ph06_r5_supplier_balance.mjs`
- `tests/runtime/ph06_supplier_flow.mjs`
- `tests/runtime/ph06_ui_category.mjs`
- `tests/runtime/ph06_ui_forms.mjs`
- `tests/runtime/ph06_ui_forms2.mjs`

## Evidence
- `evidence/database/RUN-003/PH-06`
- `evidence/database/RUN-004/PH-06`
- `evidence/database/RUN-005/PH-06`
- `evidence/network/RUN-003/PH-06`
- `evidence/network/RUN-004/PH-06`
- `evidence/network/RUN-005/PH-06`
- `evidence/screenshots/RUN-003/PH-06`
- `evidence/screenshots/RUN-004/PH-06`
- `evidence/screenshots/RUN-005/PH-06`

## Existing test IDs (34b)
EX-PG-dash-owner-001, EX-PG-dash-owner-002, EX-PG-dash-owner-003, EX-PG-dash-owner-004, EX-PG-dash-owner-005, EX-PG-dash-owner-006, EX-PG-dash-owner-007, EX-PG-dash-owner-008, EX-PG-dash-owner-009, EX-PG-dash-owner-010, EX-PG-dash-owner-011, EX-PG-dash-owner-012, EX-PG-dash-owner-013, EX-PG-dash-owner-014, EX-PG-dash-owner-015, EX-PG-dash-owner-016, EX-PG-dash-owner-017, EX-PG-dash-owner-018, EX-PG-dash-owner-019, EX-PG-dash-owner-020, EX-PG-dash-owner-021, EX-PG-dash-owner-022, EX-PG-dash-owner-023, EX-PG-dash-owner-024, EX-PG-dash-owner-025, EX-PG-dash-owner-026, EX-PG-dash-owner-027, EX-PG-dash-owner-028, EX-PG-dash-owner-029, EX-PG-dash-owner-030, EX-PG-dash-owner-031, EX-PG-dash-owner-032, EX-PG-dash-owner-033, EX-PG-dash-owner-034, EX-PG-dash-owner-035, EX-PG-dash-owner-036, EX-PG-dash-owner-037, EX-PG-dash-owner-038, EX-PG-dash-owner-039, EX-PG-dash-owner-040, EX-PG-dash-owner-041, EX-PG-dash-owner-042, EX-PG-dash-owner-043, EX-PG-dash-owner-044, EX-PG-dash-owner-045, EX-PG-dash-owner-046, EX-PG-dash-owner-047, EX-PG-dash-owner-048, EX-PG-dash-owner-049, EX-PG-dash-owner-050, EX-PG-dash-owner-051, EX-PG-dash-owner-052, EX-PG-dash-owner-053, EX-PG-dash-owner-054, EX-PG-dash-owner-055, EX-PG-dash-owner-056, EX-PG-dash-owner-057, EX-PG-dash-owner-058, EX-PG-dash-owner-059, EX-PG-dash-owner-060 …

## Actual results
- Controlled stock-vs-ledger experiment (OBS-INV-001 not reproduced on a fresh tracked product); supplier invoice/payment/concurrency matrix.
- Limit: product/category/supplier flows were driven by API-diagnostic scripts, not by UI forms.

## Findings
- [F-011](../../findings/F-011-inventory-valuation-and-turnover-do-not-tie-to-l/finding.md) — Inventory valuation and turnover do not tie to ledger; opening stock has no journal (High)
- [F-034](../../findings/F-034-waste-shrinkage-adjustments-change-stock-but-pos/finding.md) — Waste/shrinkage adjustments change stock but post nothing and the turnover report shows wa (Medium)
- [F-040](../../findings/F-040-no-cogs-or-inventory-movement-is-posted-for-prod/finding.md) — No COGS or inventory movement is posted for products with track_inventory=false even when  (Medium)
- [F-041](../../findings/F-041-supplier-invoice-vat-is-computed-on-the-pre-disc/finding.md) — Supplier invoice VAT is computed on the pre-discount amount, overstating input VAT and acc (High)

## Remaining blockers and unexecuted work
Recorded per task in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and in [48-blocked-vs-unexecuted-and-prerequisites.md](../../48-blocked-vs-unexecuted-and-prerequisites.md).

## Required new tests (if the phase is re-opened)
Any row of 51 with status PLANNED/BLOCKED for this phase; plus the regression tests in [59](../../59-post-remediation-regression-test-register.csv) after remediation.

## Required deliverables
Evidence under `evidence/{screenshots,network,database}/RUN-003/PH-06/`, updated registers (22, 34b, 43, 44), progress.md entry.

## Completion criteria
Every executable task in 51 for this phase is EXECUTED-* or has a documented blocker; findings recorded in 22 with evidence; no product code changed.
