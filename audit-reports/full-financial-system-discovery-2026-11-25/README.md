# Full Financial System Discovery — Audit Package

Baseline: `Yousef-Abdallah-1711/pos` `main` @ `133260db66b1143c254d79476daa7650b41a33f0`.
**ONE-FILE MAP (phases + tasks + Definition of Done + every issue): [73-SINGLE-FILE-PHASE-TASK-ISSUE-MAP.md](73-SINGLE-FILE-PHASE-TASK-ISSUE-MAP.md).**
**CURRENT MAIN ENTRYPOINT: [65-AUDIT-MASTER-OVERVIEW.md](65-AUDIT-MASTER-OVERVIEW.md)** (then [68](68-LOW-COST-AGENT-EXECUTION-GUIDE.md), queue [66](66-remaining-audit-execution-queue.csv), file map [67](67-audit-file-and-artifact-map.md), future remediation [69](69-FUTURE-REMEDIATION-OVERVIEW.md), evidence guide [70](70-EVIDENCE-AND-TRACEABILITY-GUIDE.md), acceptance checklist [71](71-FINAL-AUDIT-ACCEPTANCE-CHECKLIST.md), next-agent prompt [72](72-NEXT-AGENT-START-PROMPT.md)).
Older start point: [00-executive-summary.md](00-executive-summary.md) (discovery completion report) → [progress.md](progress.md) (single central tracker).

| # | File | Content |
|---|---|---|
| 00 | [00-executive-summary.md](00-executive-summary.md) | Discovery completion report, counts, top findings, limits |
| — | [progress.md](progress.md) | Central tracker |
| 01 | [01-environment-baseline.md](01-environment-baseline.md) | Verified repo/runtime baseline |
| 02 | [02-repository-knowledge-review.md](02-repository-knowledge-review.md) | Docs, remediation timeline, existing tests |
| 03 | [03-complete-system-architecture.md](03-complete-system-architecture.md) | Layer-by-layer architecture |
| 04 | [04-financial-model-and-policies.md](04-financial-model-and-policies.md) | Money representation & accounting policies (docs vs code vs DB vs unknown) |
| 05 | [05-module-dependency-map.md](05-module-dependency-map.md) | Package/module dependencies |
| 06 | [06-complete-page-inventory.md](06-complete-page-inventory.md) | 119 pages / route handlers |
| 07 | [07-field-level-inventory.csv](07-field-level-inventory.csv) | 200 field rows (partial coverage) |
| 08 | [08-api-inventory.csv](08-api-inventory.csv) | 374 routes |
| 09 | [09-database-schema-and-erd.md](09-database-schema-and-erd.md) | Running schema, RLS, grants, ERD, integrity probes |
| 10 | [10-financial-transaction-lifecycles.md](10-financial-transaction-lifecycles.md) | State machines & consequences |
| 11 | [11-financial-posting-matrix.csv](11-financial-posting-matrix.csv) | 36 events → journal mapping |
| 12 | [12-report-formula-catalog.md](12-report-formula-catalog.md) | Report sources & formulas |
| 13 | [13-accounting-event-architecture.md](13-accounting-event-architecture.md) | Outbox / recovery |
| 16 | [16-financial-oracle-results.csv](16-financial-oracle-results.csv) | Oracle expected values (app actuals pending) |
| 21 | [21-existing-issue-regression-matrix.csv](21-existing-issue-regression-matrix.csv) | 36 historical issues |
| 25 | [25-cross-module-dependency-matrix.md](25-cross-module-dependency-matrix.md) | Cross-module matrix |
| 26 | [26-business-workflow-catalog.md](26-business-workflow-catalog.md) | 25 workflows + NOT IMPLEMENTED register |
| 27 | [27-financial-data-lineage.md](27-financial-data-lineage.md) | Report lineage, 12 competing sources of truth |
| 28 | [28-security-and-reliability-architecture.md](28-security-and-reliability-architecture.md) | Claimed vs implemented controls |
| 29 | [29-knowledge-gap-register.md](29-knowledge-gap-register.md) | 48 knowledge gaps with resolution tests |
| — | [phase2-progress-notes.md](phase2-progress-notes.md) | Inventory coverage notes |
| — | diagrams/ (27 files), oracle/, sql/, tests/, evidence/database/ | Supporting artifacts |

| 07b/07c | [07b-form-inventory.csv](07b-form-inventory.csv), [07c-field-coverage-notes.md](07c-field-coverage-notes.md) | Form containers and coverage limits |
| 27b | [27b-financial-lineage-chains.md](27b-financial-lineage-chains.md) | 20 lineage chains |
| 30/30b | [30-report-metric-inventory.csv](30-report-metric-inventory.csv), [30b](30b-report-metric-notes.md) | 112 metric rows |
| 31 | [31-independent-finding-verification.csv](31-independent-finding-verification.csv) | 28 re-verified findings |
| 32 | [32-accounting-policy-decision-register.md](32-accounting-policy-decision-register.md) | Policy decisions (all OPEN) |
| 33/33b | [33-capability-completeness-matrix.csv](33-capability-completeness-matrix.csv), [33b](33b-competing-sources-of-truth-expanded.md) | Capabilities; sources of truth |
| 34 | [34-complete-manual-ui-test-plan.csv](34-complete-manual-ui-test-plan.csv) | 1,491 test cases (NOT EXECUTED) |
| 35/35a | [35](35-financial-scenario-test-specification.md), [35a](35a-oracle-independent-review.md) | 24 scenarios; oracle review |
| 36 | [36-security-reliability-test-plan.md](36-security-reliability-test-plan.md) | 44 tests |
| 37 | [37-isolated-environment-execution-runbook.md](37-isolated-environment-execution-runbook.md) | Runbook (not executed) |
| 38 | [38-master-traceability-matrix.csv](38-master-traceability-matrix.csv) | Traceability |
| 39 | [39-discovery-closure-and-execution-readiness.md](39-discovery-closure-and-execution-readiness.md) | Closure report and verdicts |

| 14–15, 17–20 | [14](14-manual-ui-execution-results.md) [15](15-financial-transaction-execution-ledger.md) [17](17-cross-module-integration-results.md) [18](18-negative-boundary-test-results.md) [19](19-security-concurrency-recovery-results.md) [20](20-financial-reconciliation-report.md) | Execution results |
| 16 | [16-financial-oracle-results.csv](16-financial-oracle-results.csv) | Oracle expected vs app actual (23 rows) |
| 22–24 | [22](22-confirmed-findings-register.csv) [23](23-remediation-roadmap.md) [24](24-financial-release-verdict.md) | Findings, roadmap, verdict (NO-GO) |
| 40–47 | [40](40-runtime-environment-verification.md) [41](41-fresh-tenant-onboarding-results.md) [42](42-complete-operating-month-reconciliation.md) [43](43-report-metric-runtime-results.csv) [44](44-transaction-level-reconciliation.csv) [45](45-final-coverage-and-unexecuted-tests.md) [46](46-root-cause-dependency-map.md) [47](47-environment-restoration-report.md) | Runtime verification, onboarding, coverage, root causes, restoration |
| tests/runtime | Playwright and helper scripts used in execution | |

Still pending: 21 (historical regression matrix not re-executed; unchanged from discovery).
Evidence labels: CONFIRMED-BY-SOURCE · CONFIRMED-BY-READONLY-RUNTIME · DOCUMENTED-INTENT · INFERRED · UNKNOWN · REQUIRES-EXECUTION.

## Batch 2 additions (Master Prompt 04)
[21 (runtime columns)](21-existing-issue-regression-matrix.csv) · [22 (80 findings)](22-confirmed-findings-register.csv) · [34b executed results](34b-executed-test-results.csv) · [35b scenario status](35b-scenario-execution-status.csv) · [43 (36 metric rows)](43-report-metric-runtime-results.csv) · [45 coverage](45-final-coverage-and-unexecuted-tests.md) · [48 blocked vs unexecuted + prerequisite order](48-blocked-vs-unexecuted-and-prerequisites.md) · tests/runtime (crawler, API sweep, helpers, persistent compose copy).

## Master Prompt 05 additions (RUN-003, started 2026-11-26)
[49 takeover & verification](49-existing-audit-takeover-and-verification.md) · [50 remaining execution plan](50-remaining-audit-execution-plan.md) · [51 phased audit task register](51-phased-audit-task-register.csv) (56 tasks) · [55 finding evidence & root-cause verification](55-finding-evidence-and-root-cause-verification.csv). New evidence lives under `evidence/<type>/RUN-003/PH-xx/`. New audit-only scripts: `tests/runtime/chrome_lib.mjs`, `chrome_headed_smoke.mjs`, `ph10_emp_create_headed.mjs`, `tests/build_task_register.py`, `tests/build_finding_evidence_check.py`. Resume order: README → progress.md → 50 → 51 (62 next-model handoff not yet created).

### RUN-003 deliverables added later (2026-11-26) — see [61 closure report](61-final-exhaustive-audit-closure-report.md)
- Remediation package: [`findings/`](findings/) (one folder per finding, 10 files each), [56 policy decisions](56-accounting-policy-decision-package.md), [57 future task register](57-future-remediation-task-register.csv) (all NOT STARTED), [58 spec index](58-remediation-specification-index.md), [59 regression test register](59-post-remediation-regression-test-register.csv); generators: `tests/build_finding_packages.py` (+ `tests/finding_overrides.json`), `tests/build_phase_specs.py`, `tests/task_status_overlay.py`, `tests/build_closure_report.py`, `tests/build_evidence_manifest.py`.
- Traceability: [52 money-operation trace](52-complete-money-operation-trace.md), [`phases/phase-00 … phase-22`](phases/) (spec.md + tasks.md that reference, not copy, the authoritative files), [60 evidence manifest](60-evidence-manifest.csv) (sha256 of every evidence file; `evidence/state-*.json` listed but never read), [62 handoff](62-next-model-handoff.md).
- Runtime scripts added: `tests/runtime/ph04_ui_matrix.mjs` (headed Chrome route matrix EN/AR/mobile), `ph06_*`, `ph07_*`, `ph08_*`, `ph09_customer_flow`, `ph10_hr_*`, `ph11_accounting_flow*`, `ph11_import`, `ph12_*`, `ph13_saas`, `ph14_month`, `ph16_*`, `ph17_*`, `ph18_*`, `ph19_*` — those that call the API directly are labelled API-DIAGNOSTIC in 34b (they are not UI tests).
- Environment notes: recreate frontends when the API container is recreated; storefront on tenant subdomains needs `P2A_CORS_ORIGINS` with `http://*.localhost:13113`; restarting Redis invalidates every session (a long-running browser script must re-login).

## RUN-004 (2026-11-27)
See [63-resume-baseline-run004.md](63-resume-baseline-run004.md), [53b-ui-control-interaction-results.md](53b-ui-control-interaction-results.md), [54b-field-denominator-and-form-coverage.md](54b-field-denominator-and-form-coverage.md), the addenda in 45/48, and [62-next-model-handoff.md](62-next-model-handoff.md). Evidence for the run is under `evidence/*/RUN-004/`.

## RUN-005 index (2026-11-27)
Current authoritative counts are in [61](61-final-exhaustive-audit-closure-report.md) (generated) and [progress.md](progress.md); the numbers quoted in the historical table above (for example 36 findings, 128 forms, 627 fields) belong to earlier runs. Findings now: 80.
New or extended in RUN-005: [64 baseline and continuity](64-run005-baseline.md), [06b route and API reconciliation](06b-route-and-api-reconciliation.md), [30d metric verification classes](30d-metric-verification-classes.md), [53b control interaction results](53b-ui-control-interaction-results.md), [54b field denominator and form coverage](54b-field-denominator-and-form-coverage.md), [54c isolated field-level results](54c-isolated-field-results.md), [56b policy decision detail](56b-policy-decision-detail.md), [58b package validation](58b-package-validation.md), [63 RUN-004 baseline](63-resume-baseline-run004.md). Evidence: `evidence/*/RUN-005/`. Handoff: [62](62-next-model-handoff.md).
