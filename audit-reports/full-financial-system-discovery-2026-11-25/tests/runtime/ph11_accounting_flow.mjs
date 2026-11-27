// RUN-003 / PH-11: accounting workflow matrix on isolated tenant 4 (API-DIAGNOSTIC: fetch + SQL oracle). No product code touched.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-11');
let cookie = '';
const call = async (method, path, body) => {
  const r = await fetch(API + path, { method, headers: { Origin: 'http://localhost:13112', Cookie: cookie, 'Content-Type': 'application/json' }, body: body ? JSON.stringify(body) : undefined });
  const sc = r.headers.getSetCookie?.() ?? []; if (sc.length) cookie = sc.map((c) => c.split(';')[0]).join('; ');
  const t = await r.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 200); } return { status: r.status, body: j };
};
const acct = (code) => Number(sql(`select id from chart_of_accounts where tenant_id=4 and account_code='${code}'`));
const bal = (code) => Number(sql(`select coalesce(sum(debit_amount-credit_amount),0) from general_ledger g join chart_of_accounts a on a.id=g.account_id where g.tenant_id=4 and a.account_code='${code}'`));
const log = [];
const rec = (id, desc, r, extra = {}) => { const e = { id, desc, status: r?.status, resp: JSON.stringify(r?.body).slice(0, 260), ...extra }; log.push(e); console.log(id, r?.status, desc, JSON.stringify(extra), JSON.stringify(r?.body).slice(0, 140)); return e; };
await call('POST', '/auth/login', { email: 'audit.c.owner@example.test', password: PASSWORD });
const je = (date, desc, lines, type = 'manual') => ({ entry_date: date, description: desc, entry_type: type, lines: lines.map(([code, dr, cr], i) => ({ account_id: acct(code), debit_amount: dr, credit_amount: cr, line_order: i + 1 })) });
const post = async (id) => call('POST', `/accounting/journal-entries/${id}/post`, {});
const jid = (r) => r.body?.data?.id ?? r.body?.id;

// A. account management
let r = await call('POST', '/accounting/accounts', { account_code: '3150', account_name: 'Owner Drawings', account_type: 'equity', normal_balance: 'debit', currency: 'EGP' }); rec('ACC-01', 'create equity account 3150 Owner Drawings', r);
r = await call('POST', '/accounting/accounts', { account_code: '3150', account_name: 'Duplicate', account_type: 'equity', normal_balance: 'debit' }); rec('ACC-02', 'duplicate account code rejected?', r);
r = await call('POST', '/accounting/accounts', { account_code: '9999', account_name: 'Bad type', account_type: 'nonsense', normal_balance: 'debit' }); rec('ACC-03', 'invalid account_type rejected?', r);
r = await call('POST', '/accounting/accounts', { account_code: '1250', account_name: 'Delivery Van', account_type: 'asset', normal_balance: 'debit' }); rec('ACC-04', 'create asset account 1250', r);

// B. journal validation
r = await call('POST', '/accounting/journal-entries', je('2026-11-26', 'unbalanced probe', [['1110', 100, 0], ['3100', 0, 90]])); rec('JE-01', 'unbalanced journal rejected?', r);
r = await call('POST', '/accounting/journal-entries', je('2026-11-26', 'both sides probe', [['1110', 100, 100], ['3100', 0, 0]])); rec('JE-02', 'line with debit and credit both set / zero line', r);
r = await call('POST', '/accounting/journal-entries', je('2026-11-26', 'negative probe', [['1110', -100, 0], ['3100', 0, -100]])); rec('JE-03', 'negative amounts rejected?', r);

// C. capital contribution (balanced draft -> post)
const cap0 = bal('1110');
r = await call('POST', '/accounting/journal-entries', je('2026-11-26', 'Capital contribution probe', [['1110', 50000, 0], ['3100', 0, 50000]])); const cj = jid(r); rec('JE-04', 'capital contribution draft', r, { status_after: sql(`select status from journal_entries where id=${cj}`) });
r = await post(cj); rec('JE-05', 'post capital contribution', r, { cashDelta: bal('1110') - cap0 });
// D. owner withdrawal
const w0 = bal('1110');
r = await call('POST', '/accounting/journal-entries', je('2026-11-26', 'Owner withdrawal probe', [['3150', 8000, 0], ['1110', 0, 8000]])); const wj = jid(r); r = await post(wj); rec('JE-06', 'post owner withdrawal via custom 3150', r, { cashDelta: bal('1110') - w0 });
// E. fixed asset purchase (investing) + depreciation
r = await call('POST', '/accounting/journal-entries', je('2026-11-26', 'Van purchase probe', [['1250', 20000, 0], ['1110', 0, 20000]])); const aj = jid(r); r = await post(aj); rec('JE-07', 'post fixed-asset purchase', r, { van: bal('1250') });
const cf = await call('GET', '/reports/financial/cash-flow?date_from=2026-10-28&date_to=2026-11-26'); rec('CF-INV', 'cash-flow investing after asset purchase (cash fell by 20,000)', cf, { investing: cf.body?.data?.investing_cash_flow, net: cf.body?.data?.net_cash_flow });
// F. expense (rent) and edit / re-post rules
r = await call('POST', '/accounting/journal-entries', je('2026-11-26', 'Rent probe', [['5310', 1500, 0], ['1110', 0, 1500]])); const rj = jid(r); r = await post(rj); rec('JE-08', 'post rent expense', r);
r = await call('PUT', `/accounting/journal-entries/${rj}`, je('2026-11-26', 'Rent probe edited', [['5310', 1, 0], ['1110', 0, 1]])); rec('JE-09', 'edit a POSTED journal rejected?', r, { rentNow: bal('5310') });
r = await post(rj); rec('JE-10', 'post an already posted journal (idempotency)', r, { rentAfter: bal('5310') });
// G. void and reversal
const v0 = bal('5310');
r = await call('POST', `/accounting/journal-entries/${rj}/void`, { reason: 'probe void' }); rec('JE-11', 'void posted journal creates reversal', r, { rentAfterVoid: bal('5310'), delta: bal('5310') - v0, status: sql(`select status from journal_entries where id=${rj}`) });
r = await call('POST', `/accounting/journal-entries/${rj}/void`, { reason: 'again' }); rec('JE-12', 'void an already void journal', r);
// H. backdated and future dated postings
r = await call('POST', '/accounting/journal-entries', je('2019-01-15', 'backdated 2019 probe', [['1110', 10, 0], ['3100', 0, 10]])); const bj = jid(r); r = await post(bj); rec('JE-13', 'post a journal dated 2019 (no fiscal period)', r, { status: sql(`select status from journal_entries where id=${bj}`) });
r = await call('POST', '/accounting/journal-entries', je('2030-06-01', 'future 2030 probe', [['1110', 10, 0], ['3100', 0, 10]])); const fj = jid(r); r = await post(fj); rec('JE-14', 'post a journal dated 2030', r, { status: sql(`select status from journal_entries where id=${fj}`) });
// I. posting to header / inactive accounts
const hdr = sql("select id from chart_of_accounts where tenant_id=4 and is_header limit 1");
if (hdr) { r = await call('POST', '/accounting/journal-entries', { entry_date: '2026-11-26', description: 'header account probe', entry_type: 'manual', lines: [{ account_id: Number(hdr), debit_amount: 5, credit_amount: 0, line_order: 1 }, { account_id: acct('3100'), debit_amount: 0, credit_amount: 5, line_order: 2 }] }); const hj = jid(r); const pr = hj ? await post(hj) : r; rec('JE-15', 'post to a header (non-postable) account', pr, { created: r.status }); }
// J. period handling
r = await call('GET', '/accounting/fiscal-periods'); rec('PER-01', 'list fiscal periods', r);
// K. bank account opening balance (F-029 revalidation) + transactions
r = await call('POST', '/accounting/bank-accounts', { account_name: 'Audit Savings B', bank_name: 'Audit Bank', account_number: '900100', currency: 'EGP', account_type: 'savings', opening_balance: 12000, opening_date: '2026-11-26', gl_account_id: acct('1112') }); const bid = r.body?.data?.id ?? r.body?.id; rec('BNK-01', 'create bank account opening 12,000 linked to GL 1112', r, { gl1112: bal('1112'), journalsForBank: sql(`select count(*) from journal_entries where tenant_id=4 and description ilike '%opening%'`) });
r = await call('GET', `/accounting/bank-accounts/${bid}/transactions`); rec('BNK-02', 'bank transactions list', r);
// L. recurring, budget, FX
r = await call('POST', '/accounting/recurring-journals', { name: 'Rent monthly probe', frequency: 'monthly', start_date: '2026-11-26', template: je('2026-11-26', 'recurring rent', [['5310', 100, 0], ['1110', 0, 100]]) }); rec('REC-01', 'create recurring journal', r);
r = await call('POST', '/accounting/budgets', { name: 'Audit budget', fiscal_year: 2026, lines: [{ account_id: acct('5310'), amount: 36000 }] }); rec('BUD-01', 'create budget', r);
r = await call('POST', '/accounting/exchange-rates', { from_currency: 'USD', to_currency: 'EGP', rate: 48.5, effective_date: '2026-11-26' }); rec('FX-01', 'create exchange rate', r);
r = await call('GET', '/accounting/exchange-rates/convert?from=USD&to=EGP&amount=100&date=2026-11-26'); rec('FX-02', 'convert 100 USD', r);
// M. integrity
const tb = await call('GET', '/accounting/trial-balance'); rec('TB-01', 'trial balance balanced after all postings', tb, { balanced: tb.body?.data?.is_balanced, debit: tb.body?.data?.total_debit, credit: tb.body?.data?.total_credit });
const gl = sql("select coalesce(sum(debit_amount),0)||'|'||coalesce(sum(credit_amount),0) from general_ledger where tenant_id=4");
log.push({ glTotals: gl });
fs.writeFileSync(`${d.db}/AUD-ACC-FLOW.json`, JSON.stringify({ run: 'RUN-003', kind: 'API-DIAGNOSTIC', tenant: 4, log }, null, 1));
