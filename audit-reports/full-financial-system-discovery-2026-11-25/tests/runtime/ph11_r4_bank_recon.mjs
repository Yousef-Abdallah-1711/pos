// RUN-004 / PH-11: bank reconciliation behaviour (API-DIAGNOSTIC, tenant 4 fixture). Questions: is book balance read from the ledger or supplied by the client? are line journal ids validated (exist, same tenant, bank account)?
// can a session be completed without matching real ledger movements? what does completion write? Also confirms depreciation capability status (fixed_assets table without endpoints).
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-11'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const A = async (m, p, b) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined }); const t = await x.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 200); } return { status: x.status, body: j }; };
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 420)); };
const bank = sql("select id||'|'||account_name||'|'||current_balance||'|'||coalesce(gl_account_id::text,'-') from bank_accounts_ledger where tenant_id=4 order by id limit 3");
rec('BR-00', 'bank accounts (tenant 4)', { bank });
const bid = Number(sql('select id from bank_accounts_ledger where tenant_id=4 order by id limit 1'));
const lrBefore = sql(`select coalesce(last_reconciled_date::text,'-')||'|'||coalesce(last_reconciled_balance::text,'-') from bank_accounts_ledger where id=${bid}`);
rec('BR-01', 'bank account reconciliation state before', { bid, lrBefore, glNote: 'ledger balance of the linked GL account is the true book balance' });
let r = await A('POST', '/accounting/bank-reconciliation/sessions', { bank_account_id: bid, statement_date: '2026-11-26', statement_balance: 123456.78, book_balance: 0, notes: 'RUN-004 probe' });
const sid = r.body?.data?.id ?? r.body?.id; rec('BR-02', 'create session with arbitrary statement balance 123,456.78', { http: r.status, sid, resp: JSON.stringify(r.body).slice(0, 220) });
r = await A('POST', `/accounting/bank-reconciliation/sessions/${sid}/lines`, { journal_entry_id: 99999999, line_description: 'nonexistent journal', statement_amount: 123456.78, book_amount: 123456.78, match_status: 'matched' });
rec('BR-03', 'add line with NONEXISTENT journal_entry_id and invented book_amount', { http: r.status, resp: JSON.stringify(r.body).slice(0, 200) });
const foreign = sql("select id from journal_entries where tenant_id=2 order by id limit 1");
r = await A('POST', `/accounting/bank-reconciliation/sessions/${sid}/lines`, { journal_entry_id: Number(foreign), line_description: 'other tenant journal', statement_amount: 0, book_amount: 0, match_status: 'matched' });
rec('BR-04', 'add line referencing ANOTHER TENANT journal', { http: r.status, foreignJournal: foreign, resp: JSON.stringify(r.body).slice(0, 200) });
rec('BR-05', 'session after lines (book balance from client-supplied amounts?)', { row: sql(`select statement_balance||'|'||book_balance||'|'||difference||'|'||status from bank_reconciliation_sessions where id=${sid}`) });
r = await A('POST', `/accounting/bank-reconciliation/sessions/${sid}/complete`, {});
rec('BR-06', 'complete the session', { http: r.status, resp: JSON.stringify(r.body).slice(0, 200), session: sql(`select status||'|'||coalesce(completed_at::text,'-') from bank_reconciliation_sessions where id=${sid}`), bankAfter: sql(`select coalesce(last_reconciled_date::text,'-')||'|'||coalesce(last_reconciled_balance::text,'-') from bank_accounts_ledger where id=${bid}`) });
rec('BR-07', 'true ledger balance of the bank GL account vs reconciled balance', { gl: sql(`select coalesce(sum(g.debit_amount-g.credit_amount),0) from general_ledger g where g.tenant_id=4 and g.account_id=(select gl_account_id from bank_accounts_ledger where id=${bid})`) });
// second session: difference not zero, then complete
r = await A('POST', '/accounting/bank-reconciliation/sessions', { bank_account_id: bid, statement_date: '2026-11-27', statement_balance: 500, book_balance: 0 }); const sid2 = r.body?.data?.id ?? r.body?.id;
r = await A('POST', `/accounting/bank-reconciliation/sessions/${sid2}/complete`, {});
rec('BR-08', 'complete a session with difference <> 0 (no lines)', { http: r.status, session: sql(`select status||'|'||difference||'|'||coalesce(completed_at::text,'-') from bank_reconciliation_sessions where id=${sid2}`) });
r = await A('POST', '/accounting/bank-reconciliation/sessions', { bank_account_id: 99999, statement_date: '2026-11-27', statement_balance: 1 });
rec('BR-09', 'session for a bank account that does not exist / other tenant', { http: r.status, resp: JSON.stringify(r.body).slice(0, 200) });
rec('DEP-01', 'depreciation capability', { fixed_assets_table_rows: sql('select count(*) from fixed_assets'), routes: 'no /fixed-assets or /depreciation route registered in cmd/api/main.go (grep asset = 0 hits); PostDepreciationEntries has no caller outside the engine definition' });
fs.writeFileSync(`${d.db}/AUD-R4-BANK-RECON.json`, JSON.stringify({ run: 'RUN-004', kind: 'API-DIAGNOSTIC', tenant: 4, log }, null, 1));
