"""RUN-005: records the destructive-control executions (34b rows, 53b result updates, finding evidence). Idempotent."""
import csv, os, re
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
N = 'evidence/network/RUN-005/PH-04/'
FILES = 'AUD-R5-DESTRUCTIVE-A.run1.json; AUD-R5-DESTRUCTIVE-A.run2-S1S2S3S7.json; AUD-R5-DESTRUCTIVE-A.run4-S1S2S3.json; AUD-R5-DESTRUCTIVE-A.json (S4)'
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
new = [
    ('EX-RUN005-DESTR-CUSTOMER-SUPPLIER', 'UI-HEADED-CHROME', 'Delete controls on disposable customers and suppliers (clean and with an invoice)', 'owner', 'PASS',
     'clean customer 99 and clean supplier 85: browser confirm dialog, DELETE 200, row gone; customer 100 with an invoice: DELETE 400 with the raw foreign-key text "pq: update or delete on table customers_accounting violates foreign key constraint invoices_customer_customer_id_fkey" (data protected, message leak F-025); evidence ' + N + FILES),
    ('EX-RUN005-DESTR-PRODUCT', 'UI-HEADED-CHROME', 'Delete control on a disposable product', 'owner', 'PASS', 'confirm dialog, DELETE /products/51 200, status=deleted (soft delete); evidence ' + N + 'AUD-R5-DESTRUCTIVE-A.json'),
    ('EX-RUN005-DESTR-ORDER-CANCEL', 'UI-HEADED-CHROME', 'Cancel Order on an unpaid tracked order and on a paid order', 'owner', 'FAIL',
     'unpaid order 5796 (qty 2 of product 2): DELETE 200, status cancelled, stock 32 -> 30 at order -> 32 after cancel (restored), BUT the KDS queue row stays "new" (F-044 reproduced through the UI); paid order 5797: DELETE 400 "cannot cancel paid or refunded order before a supported refund and accounting reversal workflow exists", the UI shows only the generic toast "Failed to cancel order" (F-049, F-050); evidence ' + N + 'AUD-R5-DESTRUCTIVE-A.run1.json'),
    ('EX-RUN005-DESTR-JOURNAL', 'UI-HEADED-CHROME', 'Post and Void on a disposable journal (detail page)', 'owner', 'PASS',
     'journal 44792 draft -> Post: status posted, 2 general-ledger rows; Void (browser prompt for a reason accepted empty): status void, void_reason "Voided by user", 4 general-ledger rows (2 original + 2 reversal), net cash effect 0.00; evidence ' + N + 'AUD-R5-DESTRUCTIVE-A.run2-S1S2S3S7.json'),
    ('EX-RUN005-DESTR-THEME', 'UI-HEADED-CHROME', 'Duplicate then Delete a theme copy; active theme untouched', 'owner', 'PASS', 'themes 3 -> 4 -> 3, confirm dialog, active theme 4 intact; deleting the ACTIVE theme through the UI was not offered a Delete button (API-level deletion of the active theme is F-060); evidence ' + N + 'AUD-R5-DESTRUCTIVE-A.run1.json'),
]
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    for r in new:
        if r[0] not in ex:
            csv.writer(h).writerow(r)
# 53b: mark executed controls
EXEC = [('orders', r'Cancel Order', 'EXECUTED-RUN005 (cancel unpaid OK, paid rejected F-049; KDS row remains F-044)'), ('accounting/customers', r'Delete|Remove', 'EXECUTED-RUN005 (clean delete OK; with invoice blocked by FK, raw text F-025)'), ('accounting/suppliers', r'Delete|Remove', 'EXECUTED-RUN005 (clean delete OK)'),
        ('products', r'Delete|Remove', 'EXECUTED-RUN005 (soft delete OK)'), ('theme-builder', r'Delete', 'EXECUTED-RUN005 (copy deleted, active intact)'), ('accounting/journal-entries', r'Post|Void|Reverse', 'EXECUTED-RUN005 (post and void with reversal rows)')]
rows = list(csv.reader(open('53b-ui-control-interaction-results.csv', encoding='utf-8')))
hd = rows[0]; ix = {k: i for i, k in enumerate(hd)}
n = 0
for r in rows[1:]:
    if r[ix['result']] != 'SKIPPED-DESTRUCTIVE':
        continue
    for route, pat, res in EXEC:
        if route in r[ix['route']] and re.search(pat, r[ix['label']], re.I):
            r[ix['result']] = res.split(' (')[0]; r[ix['actual']] = res; r[ix['blocker']] = ''; r[ix['evidence']] = N + 'AUD-R5-DESTRUCTIVE-A*.json'; n += 1; break
csv.writer(open('53b-ui-control-interaction-results.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
# finding evidence
U = {'F-044': ' RUN-005 real-UI reproduction: cancelling unpaid order 5796 through the Orders page (confirm dialog, DELETE 200, status cancelled, stock restored 30 -> 32) left its kds_queue row in status "new".',
     'F-049': ' RUN-005 real-UI reproduction: Cancel Order on paid order 5797 -> DELETE 400 "cannot cancel paid or refunded order before a supported refund and accounting reversal workflow exists"; the order stays confirmed|paid.',
     'F-050': ' RUN-005 UI: the same rejection reaches the user only as the generic toast "Failed to cancel order".',
     'F-025': ' RUN-005 UI: deleting a customer that has an invoice returns 400 with the raw text "pq: update or delete on table customers_accounting violates foreign key constraint invoices_customer_customer_id_fkey".'}
f = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
fh = f[0]; fx = {k: i for i, k in enumerate(fh)}
for r in f[1:]:
    if r[0] in U:
        if U[r[0]].strip() not in r[fx['actual']]:
            r[fx['actual']] += U[r[0]]
        if 'AUD-R5-DESTRUCTIVE-A' not in r[fx['evidence_paths']]:
            r[fx['evidence_paths']] += '; ' + N + 'AUD-R5-DESTRUCTIVE-A.run1.json; tests/runtime/ph04_r5_destructive_a.mjs'
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(f)
print('53b rows updated', n)
