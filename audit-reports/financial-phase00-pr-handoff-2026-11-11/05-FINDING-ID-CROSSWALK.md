# Finding ID crosswalk

The forensic source is `H:\pos-financial-forensic-audit\audit-reports\financial-forensic-audit-2026-11-10\18-findings.json`. The discovery source is `H:\pos-financial-phase00\audit-reports\phase-1-system-discovery\ISSUE-REGISTER.json`, with narrative at `quality\FINANCIAL-AND-HR-CORRECTNESS.md`. Original registers were read directly. Aliases are source-qualified because IDs overlap.

| Alias | Exact source title | Canonical | Relation |
|---|---|---|---|
| FORENSIC:FIN-001 | Balance Sheet equation fails after balanced revenue posting | CFR-001 | unique |
| FORENSIC:FIN-002 | Manual journal UI omits JSONB attachments and API rejects request | CFR-002 | unique |
| FORENSIC:FIN-003 | Auto-posting normalizer absorbs material imbalance into final line | CFR-003 | exact duplicate with DISCOVERY:FIN-001 |
| FORENSIC:FIN-004 | Failed auto-post may leave a draft journal requiring review | CFR-004 | related, distinct from discovery recovery finding |
| FORENSIC:FIN-005 | Order tax rate is hard-coded | CFR-005 | partial overlap with DISCOVERY:FIN-006 |
| FORENSIC:FIN-006 | Development dashboard shows duplicate requests and slow one-shot loads | CFR-006 | unique; root cause unresolved |
| FORENSIC:FIN-007 | **Absent from source register** | — | unmapped; no finding invented |
| FORENSIC:FIN-008 | Automatic journals fail insertion because JSONB attachments are nil | CFR-014 | unique |
| FORENSIC:FIN-009 | Supplier automatic entry type exceeds journal schema width | CFR-015 | unique |
| FORENSIC:FIN-010 | Paid-order auto-post errors are logged and suppressed | CFR-010 | partial overlap with DISCOVERY:FIN-005 |
| DISCOVERY:FIN-001 | Auto-posting silently plugs any journal imbalance into the last line | CFR-003 | exact duplicate with FORENSIC:FIN-003 |
| DISCOVERY:FIN-002 | No database-level debit=credit enforcement; balance trigger only recomputes totals; float tolerance | CFR-007 | unique |
| DISCOVERY:FIN-003 | No row locking: concurrent post/void, running-balance updates and number generation race | CFR-008 | unique |
| DISCOVERY:FIN-004 | float64 used for money in Go while DB uses NUMERIC/DECIMAL | CFR-009 | unique |
| DISCOVERY:FIN-005 | Order accounting not reversed on refund/cancel; posting failures only logged; create and post not atomic | CFR-010 | broader parent; forensic overlap is partial |
| DISCOVERY:FIN-006 | Tax/discount calculation defects: hardcoded 10% order tax; invoice tax ignores discounts | CFR-005 + CFR-011 | split into two distinct canonical defects; partial overlap |
| DISCOVERY:FIN-007 | Payroll records editable/deletable after payment and posting; tax/insurance are manual inputs | CFR-012 | unique |
| DISCOVERY:FIN-008 | Void and account deletion edge cases | CFR-013 | unique |

Machine-readable details and per-source evidence references are in `canonical-findings.json`; CSV aliases are in `finding-id-crosswalk.csv`.
