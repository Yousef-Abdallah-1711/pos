# FIN-019: Export parity and native XLSX / CSV safety are unverified

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P3. This is a source-backed recovery draft, not acceptance.

## Original issue text

### FIN-019 — Export parity and native XLSX / CSV safety are unverified

- **Severity/status:** P3 — BLOCKED for complete export comparison.
- **Problem / impact:** Export totals/filters can disagree with on-screen reports; spreadsheet formulas may execute from user-controlled text; CSV implementation can be mistaken for XLSX support.
- **Current / expected:** CSV/PDF behavior has historical snapshots and targeted API/UI unit coverage; native XLSX parity and formula injection behavior were not demonstrated. Expected schema, exact decimal strings, filters, bounded streaming, audit provenance and safe cell escaping.
- **Root cause:** Export handlers/clients not fully verified; candidates `backend/internal/handler/http/financial_report_handler.go`, `report_handler.go`, frontend `reportsApi.ts`.
- **Reproduction/evidence:** Frontend report tests passed (8) for API parsing/filter use and CSV-export interaction; no export file adversarial fixture/browser download test.
- **Baseline scope:** Feature parity/controls unknown on main/current.
- **Recommended solution / components:** Publish versioned export contract and independent reconciliation fixture; escape formula-leading cells, cap/stream, authorize/audit.
- **Migration:** Export audit schema may exist; verify before adding.
- **Dependencies/policy:** PII/retention and maximum export size.
- **Acceptance:** CSV/PDF/XLSX, if supported, reconcile to same ledger fixture; formula payload safe; role denied; limits/memory proven.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
