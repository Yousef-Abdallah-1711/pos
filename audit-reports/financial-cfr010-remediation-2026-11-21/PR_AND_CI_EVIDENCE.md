# PR and CI evidence

## Candidate state

- Repository: `Yousef-Abdallah-1711/pos`
- Baseline main: `c527e08bb47ec54e4d2f6fd14ecbd9b30e7cadc7`
- Candidate branch: `codex/cfr010-complete-balance-sheet-remediation`
- Candidate accepted for engineering checks: `afc0da5e8659ddea9c331c309173c58d6d9be04a` (evidence-only commit atop code SHA `e68d70d45ccf507ee0474449fbf5e3f19a96f76e`)
- Worktree: `H:\Projects\pos-cfr010-complete-20261121`
- PR: [#104](https://github.com/Yousef-Abdallah-1711/pos/pull/104), open against main. No merge or deployment occurred.

## Hosted run

- [Run 44953908598](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44953908598), exact SHA `e68d70d45ccf507ee0474449fbf5e3f19a96f76e`, conclusion **success**.
- [Run 44954541491](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44954541491), exact evidence SHA `afc0da5e8659ddea9c331c309173c58d6d9be04a`, conclusion **success**.
- Passed disposable PostgreSQL migration and reconciliation; real API build/start and authenticated API workflow; independent Balance Sheet oracle; EN/AR desktop/mobile financial browser checks; and broad payment/auth/order/accounting/reporting integration tests.
- Workflow explicitly removed API, dashboard, PostgreSQL, and Redis test services after acceptance.
- Evidence upload was skipped because the workflow restricts artifact upload to a designated private acceptance repository. The run logs and check status are linked above.

## Independent review

The final code diff at `e68d70d` received an independent read-only review. Reviewer reported no concrete actionable defects, specifically noting grouped journal-header exact totals, tenant/restaurant predicates, and the repeatable-read read-only transaction. Reviewer ran no tests; run 44953908598 supplies the test evidence.

## Current update

PR #104 checks are in progress on the evidence commit. Independent local review found no code defects; no GitHub review decision is recorded. PR merge status remains pending. Accounting approval is not inferred from engineering checks.
