# F-078 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-078-01 | regression (API/DB) | After activate or save the public API returns the new colour within 2 s | Expected: The public theme reflects an activation or save immediately (or within seconds) | automated (Go integration / API) |
| REG-F-078-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-078-03 | boundary/negative | Invalid, empty, boundary and oversized inputs for the affected field/endpoint | 4xx with a clear message (never 500/raw DB text); no partial record | API + Chrome |

## Preconditions per test
- REG-F-078-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph08_r4_theme_cache_ttl2.mjs)
- REG-F-078-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph08_r4_theme_cache_ttl2.mjs)
- REG-F-078-03: authenticated session (original audit script: tests/runtime/ph08_r4_theme_cache_ttl2.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
