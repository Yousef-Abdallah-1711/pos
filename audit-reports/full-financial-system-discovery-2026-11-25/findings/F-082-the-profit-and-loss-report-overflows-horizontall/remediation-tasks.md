# F-082 — Implementation phases and tasks (all NOT STARTED)

Finding: The Profit and Loss report overflows horizontally at tablet width (768 px): content is 803 px wide and the page scrolls sideways · Severity Low · Root-cause confidence PROBABLE

Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.


## P1 Failing test first

### REM-F-082-01
- **Finding:** F-082
- **Phase:** P1 Failing test first
- **Objective:** Add automated test that reproduces F-082 at code level
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Write a failing unit/integration test using the evidence dataset (tests/runtime/ph04_r5_lang_viewport.mjs).
- **Dependencies:** none
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** failing test committed on a branch
- **Required tests:** the new test must fail before the fix
- **Acceptance criteria:** Test fails with the observed Actual value.
- **Risks:** low
- **Status:** NOT STARTED


## P3 Backend implementation

### REM-F-082-02
- **Finding:** F-082
- **Phase:** P3 Backend implementation
- **Objective:** Fix root cause of F-082
- **Existing source paths:** frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Proposed (root cause PROBABLE: first isolate the 803 px element with the browser inspector at 768 px): wrap the report table or summary grid in frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx with an overflow-x-auto container, or let the grid collapse to one column at the md breakpoint; add a viewport test at 768, 1024 and 1440 px asserting scrollWidth <= innerWidth.
- **Dependencies:** REM-F-082-01
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** backend change
- **Required tests:** unit + integration
- **Acceptance criteria:** At 768 px the page has no horizontal scroll (document scrollWidth <= 768). | The table remains fully readable by scrolling inside its container. | No regression at 360, 390, 1024 and 1440 px, English and Arabic.
- **Risks:** medium
- **Status:** NOT STARTED


## P8 Verification

### REM-F-082-03
- **Finding:** F-082
- **Phase:** P8 Verification
- **Objective:** Run regression suite for F-082 (API, real Chrome, reconciliation)
- **Existing source paths:** tests/runtime (audit scripts) + new automated tests
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** At 768, 1024 and 1440 px the page has no horizontal scroll (scrollWidth <= innerWidth)
- **Dependencies:** all previous
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** test report
- **Required tests:** regression-test-spec.md
- **Acceptance criteria:** Every REG-F-082-* test in regression-test-spec.md passes on the isolated stack, including the expected result: No horizontal page scroll; wide tables scroll inside their own container; evidence stored under the release evidence folder
- **Risks:** low
- **Status:** NOT STARTED


## P9 Release and rollback

### REM-F-082-04
- **Finding:** F-082
- **Phase:** P9 Release and rollback
- **Objective:** Release verification for F-082
- **Existing source paths:** deployment pipeline
- **Proposed paths:** tests next to the affected package; new migration if required
- **Required changes:** Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).
- **Dependencies:** P8
- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)
- **Expected deliverables:** release note
- **Required tests:** staging smoke + reproduction
- **Acceptance criteria:** On staging the reproduction steps of F-082 no longer reproduce (expected: No horizontal page scroll; wide tables scroll inside their own container); the rollback in remediation-spec section 19 is rehearsed and its duration recorded
- **Risks:** low
- **Status:** NOT STARTED
