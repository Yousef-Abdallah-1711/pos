# 16 — Financial observability and operations (PROPOSED)

Monitor source-to-accounting lag, pending/failed age, orphan drafts, duplicate conflicts, lease expiry, settlement variance, journal/GL mismatch, AR/AP control variance, inventory variance, suspense/rounding activity, close exceptions and report watermark age. Alert immediately on any unbalanced posted journal or journal/GL mismatch; other thresholds require approved SLOs.

Alert carries scoped event/journal/correlation IDs, amount/currency, age and retry count, never secrets. Runbook: inspect source and provider evidence; determine whether effect committed; retry idempotently or create approved correction; reconcile before closing incident; retain evidence. Daily reconciliation is a proposal, not a demonstrated frequency.
