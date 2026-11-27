# RUN-004 designs, batch 3 (F-024, F-042..F-050). Policy IDs refer to 56 (POL-01..POL-15); POL-16+ are PROPOSED additions.
R4 = {}

R4['F-024'] = dict(
    design='Proposed: compute money in integer minor units (cents) in frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx (calculateLineValues :620-635 and mapLineToPayload :637-650): convert unit price and quantity to cents once, apply discount and tax with an explicit rounding rule (half-up to 2 decimals, the same rule the backend trigger uses), and send values rounded to two decimals; better, send only quantity, unit price, discount rate and tax rate and let the server return the authoritative amounts for display.',
    alternatives='A: integer cents + one rounding rule shared with the server (recommended). B: toFixed(2) at the end (still accumulates float error across lines). C: server-only calculation with preview endpoint (best, larger change).',
    acceptance=['The payload for 5,000 x 1 with 10% discount and 14% VAT contains tax_amount 630 and line_total 5130 exactly (no 630.0000000000001).', 'For 1,000 random lines the UI total equals the server-stored total to the cent.'],
    proposed_files=['frontend/apps/dashboard/src/lib/money.ts (new)'],
    api=['unchanged (values are rounded)'],
    policy='Rounding rule (half-up per line vs per document) under POL-02.',
    migration='No.',
    historical='Existing rows with float artefacts are rounded in the numeric(…,2) columns already; list any that differ from the rounded value.',
)

R4['F-042'] = dict(
    design='Proposed: make order statistics follow the approved revenue basis. In backend/internal/repository/order_repo.go GetOrderStats (:859-870) filter total revenue by status <> cancelled and payment_status = paid (or the basis approved in POL-01) and expose both values (booked vs collected); compute delivered_today from status = delivered and a delivered_at timestamp (add the column if missing) instead of COUNT(created_at >= CURRENT_DATE); use the restaurant time zone for the day boundary (see F-070).',
    alternatives='A: filter by basis and add delivered_at (recommended). B: filter only (delivered_today stays wrong).',
    acceptance=['For a dataset with cancelled and unpaid orders the stats revenue equals the approved basis computed independently in SQL.', 'delivered_today counts only orders delivered today.', 'The dashboard tiles use the same endpoint and show the same numbers.'],
    proposed_files=[],
    api=['GET /api/v1/admin/orders/stats: adds booked_revenue and collected_revenue (additive); total_revenue follows POL-01'],
    policy='POL-01 (revenue recognition basis) in 56.',
    migration='Possibly orders.delivered_at (nullable, back-filled from status history if available).',
    historical='None stored; numbers change on deploy (communicate to owners).',
)

R4['F-043'] = dict(
    design='Proposed: derive the KPI dashboard from the GL. In backend/internal/repository/accounting_advanced_repo.go:838-880 replace revenue = SUM(orders.total_amount) with the same GL revenue query used by the P&L (financial_report_usecase.go:208), expenses from the GL expense accounts, net profit = revenue - expenses, and cash position = balances of the cash and bank GL accounts (1110/1120-family) instead of bank_accounts_ledger.current_balance; show a reconciliation line when the bank module differs from the GL (see F-029).',
    alternatives='A: GL-derived KPIs (recommended; one source with the P&L). B: keep operational KPIs and label them (does not fix net profit mixing tax-inclusive revenue with GL expenses).',
    acceptance=['KPI revenue, expenses and net profit equal the P&L for the same period.', 'KPI cash position equals the GL cash and bank balance.', 'A difference between the bank module and the GL is displayed, not hidden.'],
    proposed_files=[],
    api=['GET KPI endpoint: same fields; values change; adds reconciliation block (additive)'],
    policy='POL-12 (single source of truth for dashboards/KPIs) in 56.',
    migration='No.',
    historical='None.',
)

R4['F-044'] = dict(
    design='Proposed: propagate cancellation to the kitchen. When an order is cancelled (usecase/order_usecase.go cancellation path and repository/order_repo.go status update) update the matching kds_queue row to status cancelled inside the same transaction (or add a trigger on orders.status that does so); exclude cancelled rows from /kds/queue, /kds/stats active_orders and analytics; broadcast the change over the KDS channel so open screens remove the card.',
    alternatives='A: update in the cancellation transaction + broadcast (recommended). B: DB trigger only (covers every writer, no broadcast). C: filter at read time by joining orders (hides the data inconsistency).',
    acceptance=['Cancelling an order removes it from the KDS board within 2 s and from active_orders.', 'kds_queue has no active row whose order is cancelled (SQL assertion over the fixture).', 'Analytics exclude cancelled orders.'],
    proposed_files=['backend/migrations/<next>_kds_cancel_sync.sql (new, if trigger approach)'],
    api=['KDS queue/stats exclude cancelled (behaviour change)'],
    policy='None.',
    migration='Yes if the trigger option is used; plus a one-off update of existing active rows whose order is cancelled (3 in the evidence).',
    historical='Existing stale rows are listed and set to cancelled after review.',
)

R4['F-045'] = dict(
    design='Proposed: build the forecast history from the ledger. In backend/internal/repository/accounting_advanced_repo.go GetPredictiveCashFlow (:1105-1175) read monthly cash inflows and outflows from the GL cash and bank accounts (debits and credits on the cash/bank accounts by month, excluding transfers between them) so POS and order receipts are included; keep customer receipts and supplier payments as forecast drivers for open invoices.',
    alternatives='A: GL cash/bank history (recommended). B: add orders to the existing two sources (double-counts receipts already posted to cash).',
    acceptance=['For a month with 100,000 of POS receipts and 20,000 of supplier payments the historical series shows net +80,000 equal to the GL cash movement.', 'The forecast changes when open invoices change.'],
    proposed_files=[],
    api=['response shape unchanged'],
    policy='POL-04 (cash-flow method) in 56 (same method must feed both reports).',
    migration='No.',
    historical='None.',
)

R4['F-046'] = dict(
    design='Proposed: compute turnover from the inventory movement ledger. In backend/internal/handler/http/report_handler.go (:557-645 and :2480-2609) take the opening quantity from inventory_movements at the period start (not current stock or the legacy inventory table) and waste from movements of type waste/adjustment_out; average inventory = (opening value + closing value)/2 at product cost; document the cost basis.',
    alternatives='A: movement-ledger based (recommended). B: keep legacy table and label (stale source).',
    acceptance=['For a fixture with known opening stock, purchases, sales and waste the turnover equals an independent calculation to 2 decimals.', 'Waste in the report equals the sum of waste movements (see F-034).'],
    proposed_files=[],
    api=['response shape unchanged'],
    policy='POL-03 (inventory accounting) in 56.',
    migration='No.',
    historical='None.',
)

R4['F-047'] = dict(
    design='Proposed: make recurring journals idempotent per period. In backend/internal/repository/accounting_advanced_repo.go RunRecurringJournalTemplate (:396-420) lock the template row (SELECT ... FOR UPDATE), reject a run when runDate < next_run_date unless an explicit force flag with a reason is supplied, store a unique key (template_id, period) on the generated journal and advance next_run_date/last_run_date in the same transaction.',
    alternatives='A: lock + unique (template, period) + next_run_date check (recommended). B: unique key only (concurrency handled, but manual re-runs still blocked silently).',
    acceptance=['Running the same template twice for the same period creates one journal; the second call returns 409 with the existing journal id.', 'Five concurrent runs create one journal.', 'A forced run records the reason and creates a second journal only when explicitly requested.'],
    proposed_files=['backend/migrations/<next>_recurring_journal_unique_period.sql (new)'],
    api=['POST /api/v1/accounting/recurring-journals/{id}/run: 409 on duplicate period; optional force+reason'],
    policy='POL-13 (recurring-journal idempotency) in 56.',
    migration='Yes: unique index on (template_id, period); existing duplicates listed first.',
    historical='Duplicate recurring journals in the evidence (2 runs, 2 journals) are reversed by approved journals.',
)

R4['F-048'] = dict(
    design='Proposed: validate currency pairs. In backend/internal/handler/http/accounting_advanced_handler.go (:475-482) require ISO-4217 uppercase codes from an allow-list, base != quote, and a rate within a sane range (> 0 and < configurable maximum); enforce uniqueness of (tenant, base, quote, effective_date) in the repository (accounting_advanced_repo.go:596).',
    alternatives='A: handler validation + unique index (recommended). B: handler only.',
    acceptance=['USD/USD returns 400.', 'An unknown code returns 400.', 'Two rates for the same pair and date return 409.'],
    proposed_files=[],
    api=['POST exchange-rates: 400/409 for the listed cases'],
    policy='POL-11 (currency policy) in 56.',
    migration='Optional unique index; list existing duplicates first.',
    historical='Existing same-currency pairs are listed and deactivated.',
)

R4['F-049'] = dict(
    design='Proposed: design and build a refund/void/return workflow for paid orders (the largest missing capability). Remove the blanket rejection at usecase/order_usecase.go:341 and repository/order_repo.go:634 only behind the new workflow: (1) a refund request document (order, items/quantities, reason, tender, approver) with partial and full refunds; (2) approval gated by a new RBAC action; (3) one transaction that records the refund, returns stock for tracked products at their original cost, and posts reversing journals (Dr revenue and VAT payable, Cr cash/bank/card clearing; Dr inventory / Cr COGS for returned goods); (4) idempotency key on the refund; (5) the order moves to refunded or partially_refunded, reports exclude refunded amounts per POL-06, and KDS/stock are updated.',
    alternatives='A: refund document with approval and reversing journals (recommended). B: cancel-and-recreate (loses audit trail). C: manual correcting journals only (stock and orders stay wrong).',
    acceptance=['A full refund of a paid cash order returns the cash, reverses revenue/VAT/COGS exactly (sum of journals for the order = 0) and restores stock.', 'A partial refund of one item posts proportional reversals.', 'The same refund request twice posts once; concurrent refunds cannot exceed the paid amount.', 'The P&L, balance sheet and sales reports reflect the refund per POL-06.'],
    proposed_files=['backend/internal/handler/http/order_refund_handler.go (new)', 'backend/internal/usecase/order_refund_usecase.go (new)', 'backend/migrations/<next>_order_refunds.sql (new)'],
    api=['POST /api/v1/admin/orders/{id}/refunds (new)', 'GET /api/v1/admin/orders/{id}/refunds (new)'],
    policy='POL-06 (refund / void / return policy) in 56; blocked until approved.',
    migration='Yes: refunds and refund_items tables, status values, refunded totals on orders.',
    historical='None (new capability).',
    idempotency='Mandatory Idempotency-Key on refund creation, stored with the refund (see F-017).',
    authz='New permission for refund approval, separate from order write.',
)

R4['F-050'] = dict(
    design='Proposed: map domain guard errors to client errors. Introduce sentinel errors (for example ErrPaidOrderCannotBeCancelled) in the usecase layer and translate them in handler/http/admin_order_handler.go:410 to 409 with a stable code; keep 500 only for unexpected errors; add table-driven tests for every guard.',
    alternatives='A: sentinel errors + mapper (recommended; shares the F-025 helper). B: string matching on messages (fragile).',
    acceptance=['Cancelling a paid order via PUT /admin/orders/{id}/status returns 409 with code ORDER_PAID_CANNOT_CANCEL.', 'Unexpected failures still return 500 with a request id.'],
    proposed_files=[],
    api=['PUT /api/v1/admin/orders/{id}/status: 409 for guard violations (was 500)'],
    policy='None.',
    migration='No.',
    historical='None.',
)
