# Root cause to fix matrix

| Finding | Evidence | Engineering action | Approval dependency | Current status |
|---|---|---|---|---|
| 5,000-cent synthetic residual | Existing synthetic report integration test demonstrates unclosed earnings can leave a Balance Sheet equation difference. This does not identify the original operating cause. | Preserve existing report semantics; add/extend diagnostics only after recording the accounting contract. | EARN-01 / CLOSE-01 | Mechanism confirmed synthetically; production treatment pending |
| Off-ledger synthetic opening cash | Existing synthetic fixture demonstrates a separate opening/current balance source that is absent from GL. Original opening evidence is unavailable. | Add read-only opening source, cutover, account-mapping and duplicate checks. No auto-posting. | OPEN-01 | Safeguard implementation remains open |
| Journal-to-GL effect mismatch risk | Prior implementation did not provide direct row-level scoped reconciliation. | Added exact-NUMERIC read-only scoped reconciliation counters, deterministic status, reversal checks, link validation and journal header balance check. | None for read-only detection | Implemented; PostgreSQL acceptance and independent review passed for candidate |
| Aggregate SQL header query rejected by PostgreSQL | Run 44953642408 on `59e5d92` failed at query execution with aggregate-in-FILTER error. | Replaced invalid expression with grouped `HAVING`; run 44953908598 passed on `e68d70d`. | None | Defect corrected and verified |
| Void reversal timestamp boundary | A split timestamp could disagree across stored void time and reversal date. | Use one application timestamp for both values. | None | Implemented; reversal scenarios passed in PostgreSQL |
| Original operating discrepancy | Original signed opening TB, source balances, scoped journals/GL, reports, close state, and provenance were not present in authorized evidence. | No records fabricated or changed; source recovery remains required. | Finance data owner / qualified accountant | **NOT ESTABLISHED** |

All PostgreSQL fixtures are synthetic. Successful engineering acceptance does not convert synthetic results into original-data evidence or close CFR-010.
