# F-023 — POS displays $ for EGP amounts

- **Severity:** Low  ·  **Category:** UX defect  ·  **Evidence status:** REPRODUCED (RUN-005 headed Chrome)
- **Affected modules:** POS
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-POS-001, EX-RUN003-AUD-POS-002, EX-RUN005-OBSERVED-REPRO
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: EGP symbol

Observed: Prices shown as $10000.00 | RUN-003 (AUD-EVD-F023): POS shows '$10,000.00' for an EGP product; no EGP label on page (STILL REPRODUCIBLE).  RUN-005: the POS page shows $12.00, $10.00, $100.00, $10,000.00 while the dashboard home shows "EGP 110.00" for the same tenant (currency EGP).

## Severity rationale
Rated **Low** (category: UX defect; evidence: REPRODUCED (RUN-005 headed Chrome)). Basis recorded in the register: n/a. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: POS. n/a

## Technical impact
CONFIRMED (source): frontend/apps/dashboard/src/lib/rtl.ts:81-91 formatCurrencyWithLocale defaults currency to USD (Intl.NumberFormat style currency), so any caller that omits the tenant currency prints $; the storefront passes EGP and shows "EGP" (AUD-STORE-004) while the POS shows $

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
n/a

## Where it was reproduced
Open POS

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
