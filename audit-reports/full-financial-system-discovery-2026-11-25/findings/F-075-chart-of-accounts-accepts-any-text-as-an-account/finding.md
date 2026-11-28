# F-075 — Chart of accounts accepts any text as an account code or name (HTML and SQL characters, non-numeric codes) with no format or length rule in the UI or API

- **Severity:** Low  ·  **Category:** Software defect (validation)  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Accounting > Chart of accounts
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN004-FORMMATRIX
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 400 with a field error: codes follow the numeric/hierarchical format of the chart; names are plain text within the column limit

Observed: POST /accounting/accounts returns 201 and stores account id 389 with code <b>x</b>'; -- (the response HTML-escapes it as unicode); a 300-character value returns 400 with the raw database length error; no maxlength on the inputs

## Severity rationale
Rated **Low** (category: Software defect (validation); evidence: REPRODUCED). Basis recorded in the register: 0 direct; report columns, exports and code-prefix classification (for example the 12% and 22% cash-flow rules) depend on well-formed codes. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Accounting > Chart of accounts. 0 direct; report columns, exports and code-prefix classification (for example the 12% and 22% cash-flow rules) depend on well-formed codes

## Technical impact
CONFIRMED (source + runtime): CreateAccount only trims the strings and passes them on (backend/internal/handler/http/account_handler.go:124-170: no regexp, length or format check); repository errors are returned to the client as 400 text (F-025)

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0 direct; report columns, exports and code-prefix classification (for example the 12% and 22% cash-flow rules) depend on well-formed codes

## Where it was reproduced
Headed Chrome tenant 4: open the account form, enter account code and names "<b>x</b>'; --", submit

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
