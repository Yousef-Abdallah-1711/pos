// RUN-003 / PH-11 part 2: recurring journals, budgets, FX (API-DIAGNOSTIC, tenant 4).
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-11'); let cookie = '';
const call = async (method, path, body) => { const r = await fetch(API + path, { method, headers: { Origin: 'http://localhost:13112', Cookie: cookie, 'Content-Type': 'application/json' }, body: body ? JSON.stringify(body) : undefined }); const sc = r.headers.getSetCookie?.() ?? []; if (sc.length) cookie = sc.map((c) => c.split(';')[0]).join('; '); const t = await r.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 200); } return { status: r.status, body: j }; };
const acct = (code) => Number(sql(`select id from chart_of_accounts where tenant_id=4 and account_code='${code}'`));
const bal = (code) => Number(sql(`select coalesce(sum(debit_amount-credit_amount),0) from general_ledger g join chart_of_accounts a on a.id=g.account_id where g.tenant_id=4 and a.account_code='${code}'`));
const log = []; const rec = (id, desc, r, extra = {}) => { log.push({ id, desc, status: r?.status, resp: JSON.stringify(r?.body).slice(0, 300), ...extra }); console.log(id, r?.status, desc, JSON.stringify(extra), JSON.stringify(r?.body).slice(0, 200)); };
await call('POST', '/auth/login', { email: 'audit.c.owner@example.test', password: PASSWORD });
const rent0 = bal('5310');
let r = await call('POST', '/accounting/recurring-journals', { template_code: 'RENT-M', template_name: 'Monthly rent probe', frequency: 'monthly', start_date: '2026-11-26', is_active: true, auto_post: true, lines: [{ account_id: acct('5310'), debit_amount: 700, credit_amount: 0, line_order: 1 }, { account_id: acct('1110'), debit_amount: 0, credit_amount: 700, line_order: 2 }] }); rec('REC-01', 'create recurring template (auto_post true)', r);
const rid = r.body?.data?.id ?? r.body?.id;
if (rid) { r = await call('POST', `/accounting/recurring-journals/${rid}/run`, {}); rec('REC-02', 'run recurring template', r, { rentDelta: bal('5310') - rent0 }); r = await call('POST', `/accounting/recurring-journals/${rid}/run`, {}); rec('REC-03', 'run again same period (duplicate run guard?)', r, { rentDelta: bal('5310') - rent0 }); }
r = await call('POST', '/accounting/budgets', { account_id: acct('5310'), fiscal_year: 2026, period_month: 11, budget_amount: 3000, currency: 'EGP' }); rec('BUD-01', 'create budget rent Nov 3000', r);
r = await call('GET', '/accounting/budgets?fiscal_year=2026'); rec('BUD-02', 'list budgets (variance vs GL actual)', r, { glRent: bal('5310') });
r = await call('POST', '/accounting/exchange-rates', { base_currency: 'USD', quote_currency: 'EGP', rate: 48.5, effective_date: '2026-11-26', source: 'audit-manual', is_active: true }); rec('FX-01', 'create USD/EGP 48.5', r);
r = await call('GET', '/accounting/exchange-rates/convert?base_currency=USD&quote_currency=EGP&amount=100&date=2026-11-26'); rec('FX-02', 'convert 100 USD (oracle 4,850)', r);
r = await call('POST', '/accounting/exchange-rates', { base_currency: 'USD', quote_currency: 'EGP', rate: -1, effective_date: '2026-11-26' }); rec('FX-03', 'negative rate rejected?', r);
r = await call('POST', '/accounting/exchange-rates', { base_currency: 'USD', quote_currency: 'USD', rate: 1, effective_date: '2026-11-26' }); rec('FX-04', 'same-currency pair rejected?', r);
fs.writeFileSync(`${d.db}/AUD-ACC-FLOW2.json`, JSON.stringify({ run: 'RUN-003', kind: 'API-DIAGNOSTIC', log }, null, 1));
