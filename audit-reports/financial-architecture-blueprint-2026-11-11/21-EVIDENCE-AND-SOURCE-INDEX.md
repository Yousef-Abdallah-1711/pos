# 21 — Evidence and source index

Historical evidence: `../financial-forensic-audit-2026-11-10/18-findings.json`, `19-coverage-matrix.md`, reports 05–08, 12, 14 and evidence/. FIN-001/002/003/008/009 historical runtime verified in disposable audit context; FIN-004 inferred; FIN-005/006 source/symptom; FIN-010 source-verified, runtime untested.

Source files cited include `backend/internal/usecase/financial_report_usecase.go`, `auto_posting_engine.go`, `order_usecase.go`, repositories `journal_entry_repo.go` and `general_ledger_repo.go`, journal UI/API and migrations 136–156. Primary checkout HEAD ecadb5a has extensive dirty changes, distinct from audit SHA e519f52; fetched origin/main equals historical SHA. No fresh probes were executed for this blueprint.
