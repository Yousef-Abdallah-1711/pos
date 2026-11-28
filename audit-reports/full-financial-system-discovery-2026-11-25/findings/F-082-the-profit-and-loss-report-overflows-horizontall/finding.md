# F-082 — The Profit and Loss report overflows horizontally at tablet width (768 px): content is 803 px wide and the page scrolls sideways

- **Severity:** Low  ·  **Category:** UX defect (responsive layout)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Reports > Financial > Profit and Loss
- **Root-cause confidence:** PROBABLE
- **Executed audit tests referencing this finding:** EX-RUN005-MOBILE-TABLET
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: No horizontal page scroll; wide tables scroll inside their own container

Observed: document.documentElement.scrollWidth is 803 against a 768 px viewport; the same report fits at 1024 px, 390 px and 360 px and in Arabic at 1440 px

## Severity rationale
Rated **Low** (category: UX defect (responsive layout); evidence: REPRODUCED). Basis recorded in the register: 0. Root-cause confidence: PROBABLE. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Reports > Financial > Profit and Loss. 0

## Technical impact
PROBABLE (measured layout fact; the element that sets the 803 px width was not isolated): a fixed-width element or table in the report page without an overflow wrapper at the md breakpoint

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0

## Where it was reproduced
Headed Chrome, viewport 768x1024, English, tenant 4 owner: open /en/dashboard/reports/financial/profit-loss and compare document scrollWidth with the viewport width

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
