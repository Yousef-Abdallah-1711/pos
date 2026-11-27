"""RUN-004 oracle for MET-PROF-007 and MET-CF-002 (tenant 4, isolated DB). Independent Decimal computation from SQL rows, compared with the API body captured by headed Chrome."""
import json, subprocess, decimal
from decimal import Decimal as D
decimal.getcontext().prec = 28
def q(sql):
    r = subprocess.run(['docker','exec','pos-p2a-ws05-e2e-postgres-1','psql','-U','postgres','-d','pos_phase2a_e2e','-At','-F','|','-c',sql],capture_output=True,text=True,encoding='utf-8')
    if r.returncode: raise SystemExit(r.stderr)
    return [l.split('|') for l in r.stdout.strip().splitlines() if l]
cap = json.load(open('evidence/network/RUN-004/PH-12/AUD-R4-MET2-CAPTURE.json',encoding='utf-8'))['bodies']
out = {}
# ---- PROF-007
k = [x for x in cap if x.startswith('/reports/profitability')][0]
prof = json.loads(cap[k])['data']
rows = q("""select oi.product_id, oi.quantity, oi.total_price, coalesce(p.cost,0), (p.cost is null), coalesce(nullif(trim(p.name_en),''),nullif(trim(p.name),''),'Product #'||oi.product_id)
 from order_items oi join orders o on o.id=oi.order_id left join products p on p.id=oi.product_id
 where oi.tenant_id=4 and oi.restaurant_id=4 and o.status<>'cancelled' and oi.created_at>='2026-10-29' and oi.created_at<'2026-11-28'""")
agg = {}
for pid, qty, tp, cost, miss, name in rows:
    a = agg.setdefault((pid, cost), {'name': name, 'qty': 0, 'rev': D(0), 'cogs': D(0)})
    a['qty'] += int(qty); a['rev'] += D(tp); a['cogs'] += int(qty) * D(cost)
items = sorted(({'pid': k[0], 'cost': k[1], **v, 'gp': v['rev'] - v['cogs']} for k, v in agg.items()), key=lambda x: (x['gp'], x['rev']), reverse=True)
tot_rev = sum(i['rev'] for i in items); tot_cogs = sum(i['cogs'] for i in items)
print('PROF summary API', prof['summary']); print('PROF summary SQL rev', tot_rev, 'cogs', tot_cogs, 'gp', tot_rev - tot_cogs)
print('API top', [(i.get('name') or i.get('product_name'), i.get('revenue'), i.get('cogs'), i.get('gross_profit')) for i in prof.get('top_items', [])][:5])
print('SQL top', [(i['name'], str(i['rev']), str(i['cogs']), str(i['gp'])) for i in items[:5]])
print('API lowest', [(i.get('name') or i.get('product_name'), i.get('revenue'), i.get('gross_profit')) for i in prof.get('lowest_items', [])][:5])
print('SQL lowest', [(i['name'], str(i['rev']), str(i['gp'])) for i in items[::-1][:5]])
print('API item-key sample', list(prof['top_items'][0].keys()) if prof.get('top_items') else None)
out['prof'] = {'api_summary': prof['summary'], 'sql_rev': str(tot_rev), 'sql_cogs': str(tot_cogs), 'n_items_sql': len(items), 'n_top_api': len(prof.get('top_items', [])), 'n_low_api': len(prof.get('lowest_items', []))}
# also negative-price/zero-price lines & cancelled orders effect
print('lines on cancelled orders in range (excluded):', q("select count(*), coalesce(sum(oi.total_price),0) from order_items oi join orders o on o.id=oi.order_id where oi.tenant_id=4 and o.status='cancelled' and oi.created_at>='2026-10-29'"))
print('lines with unpaid orders included:', q("select count(*), coalesce(sum(oi.total_price),0) from order_items oi join orders o on o.id=oi.order_id where oi.tenant_id=4 and o.status<>'cancelled' and o.payment_status<>'paid' and oi.created_at>='2026-10-29'"))
# ---- CF-002
cf = json.loads(cap[[x for x in cap if x.startswith('/reports/financial/cash-flow')][0]])['data']
lines = q("""select je.id, je.entry_number, je.entry_date, je.status, coa.account_code, coa.account_name, jel.debit_amount, jel.credit_amount
 from journal_entry_lines jel join journal_entries je on je.id=jel.journal_entry_id join chart_of_accounts coa on coa.id=jel.account_id
 where je.tenant_id=4 and je.restaurant_id=4 and coa.account_code like '12%' and je.entry_date between '2026-10-29' and '2026-11-27' order by je.id""")
inv = sum(D(l[6]) - D(l[7]) for l in lines if l[3] in ('posted', 'void'))
print('CF API', cf); print('12% lines', len(lines)); [print(l) for l in lines[:12]]; print('SQL investing (debit-credit, posted+void-reversal ignored)', inv)
out['cf'] = {'api': cf, 'sql_investing_debit_minus_credit': str(inv), 'lines': lines}
json.dump(out, open('evidence/database/RUN-004/PH-12/AUD-R4-ORACLE-PROF-CF.json', 'w'), indent=1, default=str)
