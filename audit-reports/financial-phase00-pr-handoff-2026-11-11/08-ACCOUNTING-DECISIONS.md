# Accounting decisions required

No tax rate, statutory rule, or accounting framework is inferred here. Defaults below are provisional engineering guardrails only.

| Decision | Options / trade-off | Provisional default | Approval / dependency |
|---|---|---|---|
| Jurisdiction and legal entity | One jurisdiction simplifies rules; multi-country requires localization | Unknown; block tax certification | Owner + qualified accountant/legal adviser |
| Accounting framework | Local GAAP / IFRS / other; changes recognition and disclosure | None selected | Accountant |
| Cash vs accrual | Cash simpler; accrual matches earned/incurred obligations | Design ledger for accrual capability; policy unapproved | Accountant + owner |
| Revenue recognition | On order, fulfillment, payment, or delivery; impacts liabilities/refunds | Do not recognize on order creation alone without policy | Accountant |
| VAT/taxes | Inclusive/exclusive, rates, exemptions, rounding, point of supply | No rate assumed; versioned configurable rules required | Tax adviser / accountant |
| Inventory costing and COGS | FIFO, weighted average, specific ID; timing at sale/close | No costing method selected | Accountant + inventory owner |
| Fiscal periods and close | Calendar/fiscal year; hard close/reopen approvals | Explicit period records and controlled corrections | Owner + accountant |
| Currency and FX | Functional currency, transaction currency, rate source/date | Preserve transaction currency and approved rate provenance | Owner + accountant |
| Rounding | Half-even/half-up/etc., precision by currency, materiality | Reject imbalance; any permitted rounding explicit and bounded | Accountant; policy before precision migration |
| Tips/service charges | Payable to staff, revenue, or split; withholding varies | Segregated liability pending policy | Accountant + payroll owner |
| Gift cards/deposits | Liability at issuance; recognition at redemption/breakage rules | Liability model; breakage disabled pending policy | Accountant + legal adviser |
| Refunds/chargebacks | Original-period vs current-period reversal; fees separate | Linked reversing event; never mutate posted entry | Accountant + payments owner |
| Marketplace/gateway fees | Gross vs net settlement; fee/VAT classification | Separate clearing and fee lines | Accountant + payment provider terms |
| Financial approval roles | Maker/checker, posting, void, reopen | Separate draft/post/void/close permissions | Owner/security |
| Audit retention | Retention periods and legal holds vary | Append-only auditable trail; retention duration unresolved | Owner + legal adviser |

These decisions block policy-sensitive Phase 01 implementation. Technical work on invariants, event identities, and failure-injection fixtures can proceed without choosing tax rates or recognition policy.
