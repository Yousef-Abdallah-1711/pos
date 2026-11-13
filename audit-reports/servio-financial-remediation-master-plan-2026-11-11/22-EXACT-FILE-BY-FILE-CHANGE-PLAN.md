# File-level changes

**Status: DRAFT / PROPOSED.**

Known evidence locations: financial_report_usecase.go:261-336 (FIN-001); auto_posting_engine.go:633-660 (FIN-003); journal_entry_repo.go and journal handler/UI/API (FIN-002/008); auto_posting.go and migration 137 (FIN-009); order_usecase.go:492-520 and admin/public order handlers (FIN-010); order_usecase.go:173 (FIN-005). Additional migrations/routes/UI locations are DISCOVERY REQUIRED against approved implementation base. No path is invented.

## Acceptance and dependencies

All changes require a stable task ID, exact source and function confirmation on the implementation baseline, tenant-scoped authorization, independent expected financial results, unit and PostgreSQL integration coverage, failure/retry coverage where applicable, and a documented rollback or forward-correction path. Unknown exact locations are DISCOVERY REQUIRED. Accounting-policy choices remain provisional until an authorized accountant/owner records approval. No implementation has been performed.
