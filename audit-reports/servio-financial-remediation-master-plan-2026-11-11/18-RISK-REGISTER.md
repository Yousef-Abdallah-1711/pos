# Risk register

**Status: DRAFT / PROPOSED.**

Risks include false statements, missing/duplicate effects, wrong account, provider mismatch, migration/data damage, tenant leakage, unauthorized export, precision drift, unreconciled history, backlog, failed recovery, unapproved policy and unmeasured capacity. Qualitative probability only where evidence permits; otherwise unknown. See matrices/risk-register.csv for owner, detection, mitigation and gate.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
