# Finding disposition

**Status: DRAFT / PROPOSED.**

FIN-001 verified P1 report equity omission; FIN-002 verified nil attachments on manual UI; FIN-003 verified material auto-normalizer mutates final line; FIN-004 inferred draft orphan on failure; FIN-005 verified hard-coded 10% tax but jurisdiction unknown; FIN-006 dev duplicate/slow observations, root cause unverified; FIN-008 verified automatic JSONB nil; FIN-009 verified entry type exceeds VARCHAR(20); FIN-010 source verified log-only errors, runtime scenario untested. Each maps to tasks in remediation-tasks.json.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
