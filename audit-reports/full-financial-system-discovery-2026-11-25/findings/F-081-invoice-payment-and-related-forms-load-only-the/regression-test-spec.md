# F-081 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-081-01 | regression (API/DB) | With more than 20 customers every customer can be found and selected in the invoice and receipt forms | Expected: Every active customer can be selected (search or paging) | automated (Go integration / API) |
| REG-F-081-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-081-03 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-081-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph09_r5_invoice_form.mjs)
- REG-F-081-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph09_r5_invoice_form.mjs)
- REG-F-081-03: authenticated session (original audit script: tests/runtime/ph09_r5_invoice_form.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
