# F-080 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-080-01 | regression (API/DB) | Creating a customer or supplier with a header or non-receivable account returns 400; the dropdown never lists header accounts | Expected: Only postable receivable (or payable) accounts are offered, and the API rejects a header or non-receivable account | automated (Go integration / API) |
| REG-F-080-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-080-03 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-080-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph09_r5_invoice_form.mjs)
- REG-F-080-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph09_r5_invoice_form.mjs)
- REG-F-080-03: authenticated session (original audit script: tests/runtime/ph09_r5_invoice_form.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
