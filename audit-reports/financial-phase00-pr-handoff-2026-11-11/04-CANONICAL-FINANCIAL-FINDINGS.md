# Canonical financial findings

Canonical IDs below are stable identifiers for the handoff. Source IDs remain unchanged in the historical registers and are retained as qualified aliases in `canonical-findings.json` and `finding-id-crosswalk.csv`. Severity/status reflect source evidence and are not a fresh production assessment.

| Canonical | Finding | Aliases | Severity / evidence status | Proposed phase |
|---|---|---|---|---|
| CFR-001 | Balance Sheet omits current-period earnings/equity | FORENSIC:FIN-001 | P1, runtime verified | Phase 01 policy + reporting design |
| CFR-002 | Manual journal UI omits JSONB attachments; API rejects default request | FORENSIC:FIN-002 | P2, runtime verified | Phase 01 ledger hardening |
| CFR-003 | Auto-post normalizer absorbs material imbalance into final business line | FORENSIC:FIN-003; DISCOVERY:FIN-001 | P1, runtime verified; exact duplicate | Phase 01 immediate safety |
| CFR-004 | Auto-post failure can leave recoverable draft | FORENSIC:FIN-004 | P2, inferred; failure not injected | Phase 01 reliability |
| CFR-005 | Order tax uses hard-coded rate | FORENSIC:FIN-005 | P2, source verified | Blocked on tax policy |
| CFR-006 | Duplicate dashboard requests / slow dev route observed | FORENSIC:FIN-006 | P2, root cause unverified | Investigation / performance |
| CFR-007 | Database lacks debit=credit and posted-entry immutability enforcement | DISCOVERY:FIN-002 | P2, source inspection | Phase 01 ledger hardening |
| CFR-008 | Concurrent posting, void, running-balance and numbering races | DISCOVERY:FIN-003 | P1, source inspection | Phase 01 concurrency |
| CFR-009 | `float64` used for monetary values while DB columns are NUMERIC/DECIMAL | DISCOVERY:FIN-004 | P2, source inventory | Policy and migration design |
| CFR-010 | Paid-order accounting may lack reversals and recoverable posting outcome | DISCOVERY:FIN-005; FORENSIC:FIN-010 | P1/P2, source verified; no runtime failure injection | Phase 01 reliability + refund policy |
| CFR-011 | Invoice tax calculation does not apply discounts to tax base | DISCOVERY:FIN-006 | P2, source inspection | Blocked on tax policy |
| CFR-012 | Paid payroll records editable/deletable; statutory deductions supplied manually | DISCOVERY:FIN-007 | P2, source inspection | Policy-dependent HR phase |
| CFR-013 | Void-period and account-deletion edge cases | DISCOVERY:FIN-008 | P3, source inspection; FK effect unverified | Phase 01 controls |
| CFR-014 | Automatic journals fail on nil JSONB attachments | FORENSIC:FIN-008 | P2, disposable DB verified | Phase 01 ledger hardening |
| CFR-015 | Supplier automatic entry type exceeds DB column width | FORENSIC:FIN-009 | P2, disposable DB verified | Phase 01 ledger schema design |

## Identity/relation notes

- `DISCOVERY:FIN-001` and `FORENSIC:FIN-003` describe the same material imbalance normalization behavior; they map to one canonical finding.
- `FORENSIC:FIN-010` is the suppressed paid-order posting-error part of broader `DISCOVERY:FIN-005`; the relationship is partial, not a blanket duplicate.
- `DISCOVERY:FIN-006` includes both hard-coded order tax and invoice-discount tax-base behavior. It partially overlaps `FORENSIC:FIN-005`; separate CFR-005 and CFR-011 preserve the distinct defects.
- `FORENSIC:FIN-004` is a draft recovery risk related to, but distinct from, order posting failure/recovery in `DISCOVERY:FIN-005`.
- Source A has no `FIN-007` entry. It is recorded as absent/unmapped; no title or defect is inferred.

Highest financial risks are incorrect statement equity (CFR-001), silent material journal alteration (CFR-003), concurrency integrity (CFR-008), and paid-order omissions/reversals (CFR-010). Phase 00 does not resolve them.
