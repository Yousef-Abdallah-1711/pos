# Executive plan

**Status: DRAFT / PROPOSED.**

Recommended order is safety containment and FIN-001 correction criteria, policy/ledger contracts, auto-post correctness, durable recovery, order/payment and inventory integrations, reporting, security, exports, measured scale, then release reconciliation. FIN-001/003 are P1 release blockers; FIN-008/009 block automatic supplier posting; FIN-010 may conceal missing sales/COGS effects. Security is a continuous gate. No dates or compliance claims.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
