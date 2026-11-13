# Inventory and subledgers

**Status: DRAFT / PROPOSED.**

Define inventory cost policy before COGS. Reconcile receipt, consumption, waste, return and adjustment movement to inventory control. Supplier invoice/payments reconcile AP; customer invoices/receipts reconcile AR. FIN-009 blocks supplier auto-entry; no claim that other invoice paths work end-to-end. Use distinct approvals and correction entries.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
