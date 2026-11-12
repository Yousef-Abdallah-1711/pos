# Restaurant POS financial architecture blueprint

**Status: PARTIAL decision package; no remediation implemented.** Core documents 00–23 now exist; 14 proposed ADR files and diagrams/matrices/evidence index are present. Cross-reference and full content quality validation is still required. This is a repository-grounded proposal for owner, qualified-accountant, security and engineering review. It is not accepted policy, compliance advice or a production-readiness claim.

## Navigation

- 00-BLUEPRINT-EXECUTIVE-SUMMARY.md
- 01-CURRENT-FINANCIAL-ARCHITECTURE.md
- 02-FINANCIAL-SOURCE-OF-TRUTH-CONTRACT.md
- 03-DOUBLE-ENTRY-LEDGER-INVARIANTS.md
- 04-MONEY-PRECISION-AND-ROUNDING.md
- 05-CHART-OF-ACCOUNTS-BLUEPRINT.md
- 06-FINANCIAL-EVENT-CATALOG.md
- 07-ORDER-TO-LEDGER-LIFECYCLE.md
- 08-TRANSACTION-ATOMICITY-AND-RECOVERY.md
- 09-FINANCIAL-STATE-MACHINES.md
- 10-ACCOUNTING-POLICY-DECISION-REGISTER.md
- 11-FINANCIAL-REPORTING-CONTRACTS.md
- 12-FINANCIAL-SECURITY-AND-TENANT-ISOLATION.md
- 13-FINANCIAL-SCALE-AND-PERFORMANCE.md
- 14-FINANCIAL-EXPORT-CONTRACT.md
- 15-FINANCIAL-TESTING-AND-INVARIANTS.md
- 16-FINANCIAL-OBSERVABILITY-OPERATIONS.md
- 17-TARGET-ARCHITECTURE-AND-MIGRATION.md
- 18-PHASED-IMPLEMENTATION-ROADMAP.md
- 19-ARCHITECTURE-DECISION-RECORDS.md
- 20-OWNER-DECISION-WORKSHOP.md
- 21-EVIDENCE-AND-SOURCE-INDEX.md
- 22-CURRENT-VS-TARGET-GAP-MATRIX.md
- 23-IMPLEMENTATION-READINESS-GATES.md
- architecture/ contains diagrams; adrs/ contains individual proposed ADRs; matrices/ contains machine-readable matrices; evidence/ records evidence classification and freshness.

## Source snapshot

- Fetched origin/main: e519f52eb0badfb80c283f79a52c51acc1c663c6 as of 2026-11-11 Cairo client date; equals the prior historical audit baseline.
- Primary workspace HEAD: ecadb5ab7665df3eb9f69886a9b186c710eff62d, branch codex/phase2a-continuation, with committed differences from the audit baseline and extensive uncommitted changes. It is not origin/main.
- Blueprint writing worktree: H:\pos-financial-forensic-audit, detached at e519f52eb0badfb80c283f79a52c51acc1c663c6. Prior audit reports remain untouched.
- Key finance source files inspected. Primary main.go has pre-existing uncommitted auth/session/POS/inventory/worker wiring. Other user changes and peer worktrees were preserved.
- Prior runtime results are HISTORICAL AUDIT EVIDENCE. No new protected-DB access or runtime test was performed for this blueprint. Recommendations are PROPOSED. Missing business/accounting inputs are OWNER DECISION REQUIRED or ACCOUNTANT APPROVAL REQUIRED.

## Prompt completeness and safety

The first prompt ended mid-sentence at “Respect a”; the continuation through Section 27 was provided and integrated. Source modified: NO. Protected DB modified: NO. Commits/pushes/merges/deployments: NO. Accounting/tax compliance certified: NO. 1,000-restaurant capacity demonstrated: NO. Blueprint acceptance remains PROPOSED.
