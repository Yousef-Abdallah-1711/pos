import json, csv, re, collections
r = json.load(open('tests/_routes_raw.json'))
MOD = {1:'PRODUCTS',2:'CATEGORIES',3:'HR',4:'NOTIFICATIONS',5:'SETTINGS',6:'USERS',7:'ROLES/RBAC',8:'REPORTS',9:'THEMES',10:'ORDERS',11:'LEAVES',12:'SALARIES',13:'KDS',14:'TABLES',15:'ACCOUNTING',16:'SUPPLIERS',17:'FINANCIAL_REPORTS',18:'PERFORMANCE',19:'TRAINING',20:'CONTRACTS'}
POST = {
 ('POST','/accounting/journal-entries/{id}/post'):('posts journal','Posts draft to GL via JournalEntryRepository.PostEntry (accounting pool)'),
 ('POST','/accounting/journal-entries/{id}/void'):('posts journal','Void/reversal via VoidEntryTx'),
 ('POST','/accounting/customer-invoices'):('posts journal','Invoice insert + journal create + post in ONE tx (invoice_customer_repo.go:452-560); no draft state'),
 ('POST','/accounting/customer-invoices/{id}/void'):('posts journal','VoidEntryTx + customer balance adjust in one tx (invoice_customer_repo.go:323)'),
 ('PUT','/accounting/customer-invoices/{id}'):('writes money','Update invoice (invoice_customer_repo.go:161); journal re-sync behaviour UNVERIFIED in this pass'),
 ('POST','/accounting/supplier-invoices/{id}/approve'):('posts journal','Approve = create+post journal in one tx (invoice_supplier_repo.go:295)'),
 ('POST','/accounting/supplier-invoices/{id}/void'):('posts journal','Void invoice; journal reversal path in invoice_supplier_repo.go:350'),
 ('POST','/accounting/payments/received'):('posts journal','Payment + journal create/post + invoice apply in one tx (payment_acct_repo.go:29-175)'),
 ('POST','/accounting/payments/made'):('posts journal','Payment + journal create/post + invoice apply in one tx (payment_acct_repo.go:175-320)'),
 ('POST','/hr/salaries/{id}/mark-paid'):('posts journal','Handler-orchestrated tx on accounting pool: MarkAsPaidTx + PostPayrollEntryTx (hr_salary_handler.go:476-565); gated by SALARIES WRITE only'),
 ('PUT','/admin/orders/{id}/payment-status'):('posts journal','paid => enqueue accounting_event_outbox in same tx (order_repo.go:575-598); journals posted ASYNC by worker (eventual); gated by ORDERS WRITE only'),
 ('POST','/accounting/recurring-journals/{id}/run'):('posts journal','Creates entry; posts if template.AutoPost (accounting_advanced_repo.go:396-480); NOT one tx (CreateJournalEntry then PostEntry)'),
 ('POST','/accounting/csv/journal-entries/import'):('posts journal','CreateImportedJournalEntry; posts if autoPost (accounting_advanced_repo.go:1921-1970); non-atomic create+post'),
 ('POST','/accounting/excel/journal-entries/import'):('posts journal','Alias of csv import (same handler)'),
 ('POST','/accounting/journal-entries'):('writes money','Creates DRAFT journal (journal_entry_handler.go:127-162); no GL until /post'),
 ('PUT','/accounting/journal-entries/{id}'):('writes money','Edits draft journal'),
 ('POST','/accounting/paid-order-events/{id}/redrive'):('writes money','Requeues failed outbox event + audit row (accounting_event_admin_repo.go:98); later worker posts journals; gated only by ACCOUNTING WRITE'),
 ('GET','/accounting/paid-order-events/failed'):('reads money','Lists failed outbox events'),
 ('POST','/inventory/adjust'):('writes money','Stock qty + inventory_movements row on OPERATIONAL pool; NO unit cost, NO journal entry (inventory_adjustment_handler.go:30-110)'),
 ('POST','/pos/orders'):('writes money','Creates order (payment_status forced pending, order_usecase.go:101); registered pos_routes.go:16; statusGuard + ORDERS WRITE'),
 ('POST','/public/orders'):('writes money','Unauthenticated storefront order creation; payment_status forced pending'),
 ('POST','/public/orders/quote'):('reads money','Server price quote'),
 ('POST','/public/orders/validate'):('reads money','Order validation'),
 ('DELETE','/public/orders/{id}'):('writes money','Public cancel; paid/refunded orders refused (order_repo.go transitionOrderStatus)'),
 ('POST','/admin/orders/{id}/items/{item_id}/cancel'):('writes money','Partial cancellation (migration 181_partial_order_cancellations); accounting reversal for already-paid orders UNVERIFIED'),
 ('PUT','/admin/orders/{id}/status'):('writes money','Status transition; no direct journal'),
 ('DELETE','/admin/orders/{id}'):('writes money','Cancel; refused if paid/refunded (no reversal workflow, order_repo.go)'),
 ('POST','/webhook/paymob'):('writes money','SaaS SUBSCRIPTION billing webhook (platform revenue), NOT tenant restaurant ledger'),
 ('POST','/subscription/initiate'):('writes money','SaaS subscription payment initiation (platform billing)'),
 ('POST','/subscription/cancel'):('writes money','SaaS subscription cancel'),
 ('GET','/subscription/payments'):('reads money','SaaS billing history'),
 ('POST','/super-admin/payments/{id}/refund'):('writes money','SaaS payment refund (platform); super-admin'),
 ('GET','/super-admin/payment-metrics'):('reads money','Platform payment metrics'),
 ('POST','/mock/paymob/process'):('writes money','MOCK payment processor; only registered when paymobMode==mock (main.go:666)'),
 ('POST','/reports/export'):('reads money','Report export'),
 ('GET','/audit/export'):('reads money','Audit trail CSV'),
}
READ_RE = re.compile(r'/(accounting|reports/financial|hr/salaries|hr/contracts|hr/insurance|hr/resignations/\{id\}/settlement|reports/(sales|overview|profitability|inventory|hr)|admin/orders|suppliers|public/orders|public/restaurants/\{slug\}/(menu|products|search|categories/\{categoryId\}/products)|subscription|audit)')
def rel(x):
    m, p = x['method'], x['path'].replace('/api/v1', '')
    if (m, p) in POST: return POST[(m, p)]
    if m == 'GET':
        if READ_RE.search(p) or p in ('/products', '/hr/salaries') or p.startswith('/products'):
            return ('reads money', '')
        return ('none', '')
    if re.search(r'/accounting/(accounts|bank-accounts|customers|exchange-rates|budgets|tax|supplier-payment|einvoice|bank-reconciliation|csv|excel)', p): return ('writes money', 'Accounting master/derived data')
    if p.startswith('/suppliers'): return ('writes money', 'Supplier master (accounting pool)')
    if p.startswith('/products'): return ('writes money', 'Product price/cost master (operational pool)')
    if re.match(r'^/hr/(salaries|contracts|insurance|resignations)', p): return ('writes money', 'HR financial terms')
    return ('none', '')
cols = ['method','path','handler_func','handler_file_line','middleware_wrappers','module','financial_relevance','notes']
seen = collections.Counter((x['method'], x['path']) for x in r)
rows = []
for x in r:
    fr, note = rel(x)
    mw = x['mw'] or 'none'
    parts = ['auth=' + x['auth'], 'wrappers=' + mw.replace('|', '+')]
    if x['status']: parts.append('tenantStatus=' + x['status'])
    if x['perm']: parts.append('permission ' + x['perm'].replace('=', ':'))
    elif x['auth'] == 'tenant-session-auth': parts.append('permission=none (identity/tenant only)')
    mi = int(x['mod']) if x['mod'] else None
    m = MOD.get(mi, '') if mi else ''
    notes = [note] if note else []
    if seen[(x['method'], x['path'])] > 1: notes.append('DUPLICATE pattern registered in alternate branch (db!=nil vs stub/legacy); see registration lines')
    if x['handler_func'] == 'inline closure': notes.append('inline closure handler in main.go (no usecase/repo layering)')
    notes.append('registered at ' + x['reg'])
    hf, hl = x['handler_func'], x['handler_file_line']
    if x['path'] == '/api/v1/pos/orders':
        hf = 'publicOrderHandler.CreateOrder (via registerPOSOrderRoute)'
        hl = 'backend/cmd/api/pos_routes.go:16'
    if x['path'] == '/uploads/{path...}':
        hf = 'productImages (handler.NewProductImageHandler)'
        hl = 'backend/cmd/api/main.go:537'
    rows.append([x['method'], x['path'], hf, hl, '; '.join(parts), (f"{x['mod']}={m}" if m else ''), fr, ' | '.join(notes)])
with open('08-api-inventory.csv', 'w', newline='', encoding='utf-8') as f:
    w = csv.writer(f); w.writerow(cols); w.writerows(rows)
print(len(rows), collections.Counter(r[6] for r in rows))
