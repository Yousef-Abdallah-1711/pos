# 69 - Future remediation overview (Workstream B: planning only, NOT STARTED)

Verified counts 2026-11-27: 81 findings, 472 implementation tasks (57), 243 regression tests (59), 22 policy decisions all OPEN (56b). Do not execute anything here; Workstream A (audit) is separate (65, 66).

Authoritative sources: task details in `57-future-remediation-task-register.csv` (do not copy them); specifications in `findings/F-xxx/remediation-spec.md`; index in `58-remediation-specification-index.md`; order and dependencies in `23-remediation-roadmap.md` and `46-root-cause-dependency-map.md`; shared-root-cause owner tasks in the packages (SHARED_OWNER map in `tests/build_finding_packages.py`).

## Implementation phases as defined in register 57
| Implementation phase | Tasks | Findings | Severity mix |
|---|---|---|---|
| P0 Requirements and policy | 37 | 37 | {'Critical': 1, 'High': 15, 'Medium': 16, 'Low': 5} |
| P1 Failing test first | 81 | 81 | {'Critical': 1, 'High': 20, 'Medium': 35, 'Low': 25} |
| P2 Database and data integrity | 21 | 21 | {'Critical': 1, 'High': 7, 'Medium': 11, 'Low': 2} |
| P3 Backend implementation | 81 | 81 | {'Critical': 1, 'High': 20, 'Medium': 35, 'Low': 25} |
| P4 Frontend and API contract | 18 | 18 | {'High': 3, 'Low': 5, 'Medium': 10} |
| P5 Posting and reporting integration | 28 | 28 | {'Critical': 1, 'High': 14, 'Medium': 11, 'Low': 2} |
| P6 Security and authorization | 14 | 14 | {'Critical': 1, 'Medium': 6, 'Low': 6, 'High': 1} |
| P7 Historical data remediation | 30 | 30 | {'Critical': 1, 'High': 14, 'Medium': 13, 'Low': 2} |
| P8 Verification | 81 | 81 | {'Critical': 1, 'High': 20, 'Medium': 35, 'Low': 25} |
| P9 Release and rollback | 81 | 81 | {'Critical': 1, 'High': 20, 'Medium': 35, 'Low': 25} |

## How to read each implementation phase
For a phase P: filter 57 by `phase`; included finding IDs are its `finding_id` values; prerequisites are the `dependencies` column (shared-root-cause owner tasks first); acceptance and Definition of Done are the `acceptance` and `required_tests` columns; regression tests are the 59 rows with the same finding_id; affected paths are `existing_paths`; policy decisions are listed per finding in `findings/F-xxx/remediation-spec.md` and `56b`. Database, financial and security implications are in `financial-impact.md` and `root-cause.md` of each package.

## Policy gates (all OPEN; the owner decides, the audit never approves)
POL-01..POL-22 in 56b. Any task whose package says "policy" must wait for that decision. Unlocked blocked workflows (F-001 new-tenant accounting bootstrap, F-013 period close, F-014 employee create, others in 48) are listed in the `Blocked workflows unlocked` line of the relevant package.

## Sequencing rule
Critical and High findings first (F-001 bootstrap is the single Critical, release-blocking); fix shared root causes once; after each implementation phase rerun its regression tests (59) and the original audit script named in each test row.
