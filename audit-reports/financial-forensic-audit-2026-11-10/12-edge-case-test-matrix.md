# Adversarial edge-case matrix

| Scenario | Result | Status |
|---|---|---|
| Balanced manual journal | Created with attachments [] and posted | PASS in disposable tenant |
| UI submit with attachments absent | HTTP 400 invalid JSON input | FAIL; FIN-002 |
| One-line journal | HTTP 400, at least two lines required | PASS reject |
| Two-line journal debit 100 / credit 90 | HTTP 400, journal entry is not balanced | PASS reject |
| Balance Sheet after balanced posting | equation_balanced false | FAIL; FIN-001 |
| Trial balance totals | 100/100, balanced true | PASS for fixture |
| Duplicate idempotency / repeated POST | Not executed | NOT RUN |
| Concurrent post/payment | Not executed | NOT RUN |
| Crash after journal before ledger | Not injected | NOT RUN |
| Refund / partial refund / chargeback | Not executed | NOT RUN |
| Cross-tenant reads/writes | Not executed | NOT RUN |
| Negative permission / export access | Not executed | NOT RUN |
| Large export / large ledger | Not executed | NOT RUN |
| Date boundary / timezone / closed period | Not executed | NOT RUN |
| Precision / currency / extreme amount | Not executed | NOT RUN |

No destructive tests were run against protected databases.

## Additional probes — 2026-11-11

| Scenario | Result | Status |
|---|---|---|
| Actual supplier auto-post Dr 100 / Cr 90 | Failed before persistence: PostgreSQL invalid JSON from nil attachments; reference row count 0 | FAIL; FIN-008 |
| Normalizer Dr 100 / Cr 90 | Final revenue credit changed 90 → 100 | VERIFIED behavior; FIN-003 |
| Normalized lines, valid JSON, manual repo route | Posted 100/100 journal; two ledger rows | CONDITIONAL PASS; not the current auto-method path |
| Reverse imbalance Dr 100 / Cr 110 | Normalization yields a two-sided line and journal repository rejects | PASS reject for this direction |
| Supplier auto entry type with valid JSON | `auto_supplier_invoice` exceeds VARCHAR(20) | FAIL; FIN-009 |
| Paid order with posting failure | Caller source logs and suppresses error | SOURCE VERIFIED; runtime NOT TESTED |
| Four balance-sheet fixtures | Equation off by EGP 100 initially and EGP 60 for remaining three; net profit matches gap | FAIL; FIN-001 |

The positive imbalance was not posted through the actual automatic supplier method because the independent attachments defect stops it before persistence. Do not report that as an end-to-end FIN-003 posting success. All remaining matrix scenarios are in `19-coverage-matrix.md`.
