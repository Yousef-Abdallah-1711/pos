# 62 - Next-model handoff (updated at the end of RUN-004, 2026-11-27)

Resume order: README.md -> progress.md -> 63-resume-baseline-run004.md -> this file -> 61 -> 51 -> 45/48 addenda -> current phase tasks.

## State
- Audit directory: `...\audit-reports\full-financial-system-discovery-2026-11-25\` (the folder is untracked by git: it exists only on this machine; back it up).
- Git HEAD `0de839b665982a2072f520c20a1aaf9f502f4bea` (audit baseline for RUN-001..003 was `133260db...`). No product code modified by the audit. Two stashes preserved.
- Findings 78 (Critical 1, High 20, Medium 33, Low 24; F-051 unused); root cause 74 CONFIRMED / 4 PROBABLE (F-018, F-021, F-038, F-074: each states the exact gap). Retracted claim: duplicate-table-number part of F-064.
- Future tasks 459 (NOT STARTED; only applicable phases per finding; shared root causes depend on one owner task via `SHARED_OWNER` in `tests/build_finding_packages.py`); future regression tests 233; policy decisions 19 (POL-16..19 are proposals; none approved).
- Register 51: {'EXECUTED-PASS': 16, 'PLANNED': 2, 'EXECUTED-PARTIAL': 35, 'EXECUTED-FAIL': 4, 'BLOCKED-CAPABILITY': 4}. Overall classification: PARTIAL. Financial release: NO-GO.

## Environment
- Compose project `pos-p2a-ws05-e2e`, persistent volume `pos-p2a-ws05-e2e_p2a_audit_pgdata`, DB `pos_phase2a_e2e`. Start from the repo root with `--project-directory .` and `P2A_CORS_ORIGINS` including `http://*.localhost:13113` (compose file `tests/runtime/docker-compose.audit-persistent.yml`); recreating the API requires recreating the four frontends. URLs: API 18082, landing 13111, dashboard 13112, storefront 13113, super-admin 13114.
- The isolated stack is currently STOPPED. Dev containers are stopped because ports 3001/3004 are held by `pos-landing-ci-check` and `pos-landing-ci-recheck` (see 61 section 8 for the exact restoration commands).
- Tenants: 4 diagnostic fixture (heavily mutated: 5,794 orders, themes 4/11/13, products 12-18, customers/suppliers/accounts/roles/users created by the form matrix); 8 and 9 earlier fixtures; **10 = synthetic benchmark tenant (54,000 posted lines, label SYNTHETIC BENCHMARK FIXTURE, owner audit.g.owner@example.test)**; 5, 6, 7 protected unseeded (F-001 evidence); 3 unseeded.
- Snapshots (sha256 files in `evidence/database/RUN-004/`): start, mid, final 2026-11-27; earlier dumps untouched.
- Never open `evidence/state-*.json` (session cookies).

## Regenerate derived files (order matters; r4_* scripts are RUN-004 overlays)
`python -I tests/build_finding_evidence_check.py` -> `python -I tests/build_finding_packages.py` (reads `tests/finding_overrides.json` and `tests/finding_overrides_r4*.py`) -> `python -I tests/task_status_overlay.py` then `python -I tests/r4_task_overlay.py` (51) -> `python -I -X utf8 tests/build_phase_specs.py` -> `python -I -X utf8 tests/build_field_denominator.py` -> `python -I -X utf8 tests/build_control_interaction_results.py` -> `python -I tests/build_evidence_manifest.py` -> `python -I tests/build_closure_report.py` -> `python -I -X utf8 tests/r4_finalize.py`. Register edits for RUN-004 are in `tests/r4_register_update*.py` and `tests/r4_rootcause_update.py` (idempotent).
- Tooling notes: heredocs containing apostrophes break the Bash tool, write files with the Write tool; `python -I` ignores PYTHONIOENCODING, use `-X utf8`; Chrome profiles for the audit live in the session scratchpad (separate profile per concurrent runner); Chrome helper `tests/runtime/chrome_lib.mjs` honours `RUN_ID` (default RUN-004).

## Remaining work (feasible, in priority order)
1. Finish environment restoration (61 section 8) and re-run the landing smoke if the landing changes again.
2. Single-field isolation tests for the 531 logical fields (54b), starting with money forms; extend the form matrix to invoices, payments, bank accounts, POS, settings sub-forms, storefront checkout, super-admin; add the forms missing from the field inventory (driver, user, role, employee fields).
3. Arabic and mobile interaction passes for the control matrix; rerun the 58 tooling-blocked controls (increase timeouts, one runner at a time); purpose-built fixtures for the 88 destructive controls.
4. Confirm root causes F-018, F-021, F-038, F-074; audit other date-bounded reports for the UTC cut-off (F-070 candidate); measure browser rendering on the 54k tenant.
5. Reconcile route counts 119 vs 117 (T-P03-002) and the 10 dynamic routes.
6. After remediation exists (not started): execute the 233 regression tests and re-run blocked workflows (payroll, period close, refunds, subscription billing, depreciation).

## Do not
Seed tenants 5, 6 or 7; commit; modify product code; trust SQL agreement alone as financial correctness; mark future regression tests as passed.
