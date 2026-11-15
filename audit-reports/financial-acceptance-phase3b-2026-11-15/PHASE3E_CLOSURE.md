# SERVIO Phase 3E Closure Evidence

**Evidence date:** 2026-11-15 (Africa/Cairo)  
**Current main:** `4a5482d3aeec70e31bd10aa6c0e0b6e23fd79145`  
**Fixture correction:** PR [#66](https://github.com/Yousef-Abdallah-1711/pos/pull/66), merged as the current main commit  
**Financial API acceptance:** [PR-head run 43763455210](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43763455210) and [main run 43763556472](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43763556472), both successful. PR-head artifact: [financial acceptance logs](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43763455210/artifacts/14180799894).

## Change and execution boundary

PR #57 had merged to main with a failing financial job. Its authenticated AR lifecycle failed before making an API request because the synthetic customer fixture appended a full nanosecond timestamp to a `VARCHAR(20)` `customer_code`. PR #66 constrained only the synthetic fixture value to the existing schema. No production column, accounting calculation, or customer data was changed.

Local `go test ./cmd/api` passed. The local Docker engine was unavailable, so PostgreSQL-backed acceptance ran in the repository's GitHub Actions workflow. That workflow started disposable PostgreSQL and Redis with ephemeral credentials, created restricted runtime roles, applied migrations, started the real Go API, used authenticated synthetic identities and mocked payments, and destroyed the services afterward. No public API, staging, production, or real customer data was used.

The exact-head and post-merge financial jobs both passed, including the live API step and the broader PostgreSQL-backed handler, use-case, repository, and API test suites. Full `SERVIO verified CI` also passed on the exact PR head in [run 43763455259](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43763455259) and on merged main `4a5482d` in [run 43763556436](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43763556436), including the five application-image build, Chromium smoke, and required aggregate. This general CI does not replace the authenticated financial browser acceptance that remains open.

## Reconciled API scenarios

All values below are integer cents, independently read from PostgreSQL or parsed from the downloaded export as identified by the test evidence.

| Scenario | Expected | Actual | Difference / result | Status |
|---|---:|---:|---:|---|
| Paid order to GL | Source subtotal 5,000 + tax 500 = total 5,500; cash 5,500; revenue 5,000; tax payable 500 | Order 5,500; cash 5,500; revenue 5,000; tax payable 500; one event and one journal | Order, cash, revenue, tax differences 0; repeated confirmation HTTP 200; foreign tenant read HTTP 404 | PASS for this workflow |
| Accountant journal | Debit 10,000 = credit 10,000 | Debit 10,000; credit 10,000; 2 GL rows | Difference 0; cashier post 403; accountant post 200; duplicate post 400; 3 invalid cases rejected; no partial writes | PASS for this workflow |
| Positive AR lifecycle | Invoice 10,000 less payment 4,000 = open 6,000 | Open subledger 6,000; receivable GL and aging each 6,000 | Subledger, GL, aging differences 0; foreign invoice HTTP 404 | PASS for this workflow |
| Payroll | Gross/net 150,000; expense and cash GL 150,000 each | Gross/net 150,000; expense/cash GL 150,000; one journal, 2 GL rows | Source-to-GL difference 0; injected failure rolled back; retry succeeded; paid record edits/deletes and foreign tenant denied | PASS for this workflow; statutory treatment remains pending |
| Supplier lifecycle | Invoice 22,000 = expense 20,000 + input VAT 2,000; payment clears AP | Invoice and payment journals posted; invoice paid with balance 0 | AP invoice credit and payment debit each 22,000; cash credit 22,000 | PASS for supported `supplier_invoice` / `payment_made` routes |
| Supplier concurrency and rollback | Four 5,000 payments against 22,000 invoice leave 2,000 due | Four successes; 20,000 paid; 2,000 due and AP net credit 2,000 | Repeated payment and overpayment HTTP 400; injected failures rolled back | PASS for this workflow |
| P&L and Trial Balance | Revenue equals scoped GL; debits equal credits | P&L/SQL revenue 5,000; Trial Balance debits and credits 5,500 each | Differences 0; Trial Balance balanced | PASS for measured totals |
| Balance Sheet | Assets = liabilities + equity under approved opening and earnings presentation | Assets 5,500; liabilities 500; equity 0 | Residual 5,000; API integrity state `UNRECONCILED` | PARTIAL; no balancing entry created |
| Cash Flow | API net equals independent SQL net for the fixture | API net 5,000; SQL classified net 5,000 | Arithmetic difference 0; operating 5,000, investing 0, financing 0 | PARTIAL; classification policy not approved |
| Sales CSV and PDF | CSV/PDF source gross 5,500 for one order | CSV gross 5,500; extracted PDF gross 5,500 and one order | Both monetary differences 0 | PASS for measured download content |
| Foreign-tenant access | No cross-tenant reads or writes | Order/invoice reads returned 404; payroll foreign access denied | No cross-tenant access in measured scenarios | PASS for measured endpoints |

The supplier tests use the supported current API path. They do not certify the separate legacy `PostSupplierInvoiceEntry` and `PostSupplierPaymentEntry` methods.

## Financial finding tracker refresh

This table updates the relevant engineering evidence while retaining each finding's overall acceptance status. A passing scenario does not approve accounting policy or establish a complete subsystem.

| Finding | Current evidence and disposition |
|---|---|
| CFR-001 | Unbalanced auto-posting validation remains merged; authenticated order and accountant examples reconcile. **PARTIAL**: broader source/mapping coverage is still required. |
| CFR-002 | Journal/GL database constraints and atomic-posting tests remain merged; current live scenarios persisted reconciled journals. **PARTIAL**: historical and cross-source integrity coverage remains broader than this run. |
| CFR-003 | Numbering/concurrency improvements remain merged. Supplier payment serialization passed; concurrent journal/document-number allocation was not measured here. **PARTIAL**. |
| CFR-004 | Durable event, atomic posting, and worker recovery controls remain merged; one order event completed once. Process restart/replay evidence is in PR #64; this run did not restart the process. **PARTIAL**. |
| CFR-005 | Paid-order transition protections remain merged. A complete refund and reversal lifecycle was not exercised. **PARTIAL**; behavior remains subject to REFUND-01. |
| CFR-006 | No representative financial dataset or p50/p95 query/export benchmark was run. **NOT MEASURED**. |
| CFR-007 | Non-finite validation exists on selected posting/pricing paths; financial domain models still expose floating-point monetary fields. **PARTIAL**; ROUND-01 blocks semantic precision changes. |
| CFR-008 | The tested order reconciled 5,000 revenue + 500 tax to 5,500 gross with no discount. **PARTIAL**; TAX-01 blocks claims about tax basis, jurisdiction, and discount ordering. |
| CFR-009 | Authenticated salary payment, rollback, retry, immutability, and tenant denial passed for gross = net. **PARTIAL**; deductions and statutory liabilities require PAY-01. |
| CFR-010 | Balance Sheet figures matched raw GL, and the 5,000 residual was visibly reported as `UNRECONCILED`. **PARTIAL / POLICY BLOCKED** pending OPEN-01 and EARN-01. |
| CFR-011 | Current supplier invoice/payment API lifecycle and concurrency passed. The legacy auto-posting methods remain without a verified production caller. **PARTIAL / CONTRACT PENDING**. |
| CFR-012 | Cash Flow net matched SQL at 5,000. **PARTIAL / POLICY BLOCKED** pending CF-01 classification approval. |
| CFR-013 | Invoice void/tenant exclusions are covered in prior repository evidence; closed-period correction and linked reversal acceptance were not run. **PARTIAL / POLICY BLOCKED** pending CLOSE-01 and REFUND-01. |
| CFR-014 | Worker recovery and one sales/COGS scenario have prior PostgreSQL evidence; this fixture's order product has no COGS. **PARTIAL / POLICY BLOCKED** pending INV-01 valuation approval. |
| CFR-015 | Order, AR, supplier, payroll and journal foreign-tenant/role denials passed on measured API scenarios. **PARTIAL**; not every financial endpoint or UI route was tested. |
| FIN-016 | Gift cards, deposits, and deferred revenue were not exercised. **NOT RUN**. |
| FIN-017 | Fixed assets and depreciation were not exercised. **NOT RUN**. |
| FIN-018 | Fiscal close and retained-earnings roll-forward remain **POLICY BLOCKED** pending CLOSE-01 and EARN-01. |
| FIN-019 | Authenticated CSV and PDF monetary values reconciled for the sales fixture; native XLSX support/scope was not established. **PARTIAL**. |
| FIN-020 | FX and consolidation were not exercised. **NOT RUN**. |
| FIN-021 | Financial UI locale/direction was not exercised. **NOT RUN**; EN/AR authenticated browser acceptance remains open. |

The older master audit's 30-subsystem counts are historical. The current tracked report set provides a 21-finding matrix and the handoff names 20 core modules, but does not provide the authoritative row mapping for all 30 subsystems. This update therefore refreshes the 21 canonical CFR/FIN findings and the measured subsystem scenarios above without inventing ten unmapped rows or recalculating a 30-row completion percentage.

## Accounting decision register

All ten decisions remain **PENDING APPROVAL**. No approval was supplied or inferred:

| Decision | Status | Implementation boundary |
|---|---|---|
| CFR-011 supplier auto-posting contract | PENDING APPROVAL | Keep the legacy methods out of the live route until product/accounting approves lifecycle and mappings. |
| OPEN-01 opening balances | PENDING APPROVAL | Do not synthesize opening entries or balancing plugs. |
| EARN-01 current/retained earnings | PENDING APPROVAL | Keep the 5,000-cent residual visible; do not create close entries. |
| CLOSE-01 fiscal close/reopen | PENDING APPROVAL | No unapproved close/reopen transitions. |
| TAX-01 taxes and discounts | PENDING APPROVAL | Do not change tax basis, rates, discount ordering, or rounding point. |
| ROUND-01 currency precision | PENDING APPROVAL | Do not change currency precision or monetary semantics. |
| REFUND-01 refunds and reversals | PENDING APPROVAL | Do not invent refund allocation, reversal, tax, inventory, or closed-period rules. |
| INV-01 inventory valuation | PENDING APPROVAL | Do not choose FIFO, weighted average, or another valuation basis. |
| PAY-01 payroll accounting | PENDING APPROVAL | Do not infer statutory deductions or liability mappings. |
| CF-01 Cash Flow | PENDING APPROVAL | Do not approve operating/investing/financing classification from arithmetic parity alone. |

The full options and approver register remains in [ACCOUNTING_POLICY_DECISION_REGISTER.md](../financial-acceptance-phase3-2026-11-15/ACCOUNTING_POLICY_DECISION_REGISTER.md).

## Remaining acceptance work and verdict

- **Engineering acceptance:** the Phase 3E authenticated API scenarios above pass on disposable PostgreSQL/Redis, including the corrected positive AR lifecycle, supplier concurrency/rollback, payroll, report totals, and PDF text extraction.
- **Not executed:** authenticated EN/AR financial browser workflows at desktop and 375px; a manual visual session; a live worker restart within this exact acceptance job; XLSX acceptance; and performance benchmarking on realistic financial data.
- **Policy gates:** all ten decisions above remain pending. The Balance Sheet residual is not treated as an arithmetic fix request; current earnings/opening presentation must follow approved policy.
- **Release verdict:** **NO-GO** for financial production readiness. Green isolated API acceptance establishes the tested engineering workflows only and does not approve accounting presentation or the remaining release gates.
