# 00 — Blueprint executive summary

## Recommendation

Adopt a **hybrid financial event architecture**: persist each source transition and durable outbox event in one PostgreSQL transaction; process accounting effects idempotently; permit synchronous posting only when all effects share the same database and a measured latency budget. Never call an external payment provider while a database transaction is open. A source event is not financially complete until its accounting outcome is visible and reconciled. This is PROPOSED, not implemented.

## Immediate financial gates

1. Do not use the tested Balance Sheet as a close statement. Historical fixtures showed the equation gap equals net income; the source query omits revenue/expense/current earnings from equity. HISTORICAL AUDIT EVIDENCE.
2. Remove material auto-normalization from release paths. It can change a final business line to force balance. HISTORICAL AUDIT EVIDENCE.
3. Resolve automatic journal insertion blockers: nil JSONB attachments and supplier entry type longer than VARCHAR(20). HISTORICAL AUDIT EVIDENCE.
4. Do not equate paid order status with completed accounting; auto-post errors are logged/suppressed by source. A missing-ledger paid-order runtime case is NOT VERIFIED.
5. Approve jurisdiction, accounting basis, tax, currency, inventory, refund and close policies with a qualified accountant.
6. Preserve primary local branch ecadb5a and its extensive existing edits; fetched origin/main remains e519f52.

## Current architecture

Go REST API and PostgreSQL implement orders, SaaS subscription payments, invoices, accounting payments, auto-posting, journals, GL and report handlers. Next.js Dashboard contains POS, orders, accounting and financial-report screens. The restaurant storefront/public API and KDS affect order lifecycle. Order money and payment state live on the order, while subscription payments are separate cents-based Paymob records. Customer/supplier accounting payments have distinct records. Report endpoints read different sources: journals, GL, invoice subledgers and product stock. There is no demonstrated end-to-end source-to-ledger completeness contract.

## Source-of-truth proposal

- Accepted order revision: accepted items, quantity, price, discount, tax and charges.
- Payment aggregate: internal payment lifecycle; provider-signed events and settlement statements establish external evidence.
- Inventory movement: physical stock movement.
- Posted journal: recognized accounting effect.
- GL: immutable projection of journal lines.
- Reports/exports: versioned projections over a reconciled posting watermark.

## First safe phase and readiness

First approve money and posting contracts, then add independent fixture-driven accounting/report invariants. Resolve FIN-001, FIN-003, FIN-008 and FIN-009 before production reliance. Keep policy/docs, test harness and implementation changes in independent PRs. Commercial readiness stays **NO-GO** until reports reconcile, posting recovery is durable, monetary representation is consistent, and tenant/RBAC, payments/refunds, exports, concurrency and recovery pass. Capacity is planning-only.
