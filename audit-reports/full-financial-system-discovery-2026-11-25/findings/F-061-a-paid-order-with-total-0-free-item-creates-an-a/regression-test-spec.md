# F-061 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-061-01 | regression (API/DB) | Zero-total paid order leaves the outbox completed with no journal | Expected: Zero-value sale completes without a journal (or posts a zero-line memo) and the event is completed | automated (Go integration / API) |
| REG-F-061-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-061-03 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |
| REG-F-061-04 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-061-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph16_masterdata.mjs)
- REG-F-061-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph16_masterdata.mjs)
- REG-F-061-03: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph16_masterdata.mjs)
- REG-F-061-04: authenticated session (original audit script: tests/runtime/ph16_masterdata.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
