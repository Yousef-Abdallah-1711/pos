// RUN-004 / PH-11: reconciliation completion with invented book amounts tied to a REAL own-tenant journal (API-DIAGNOSTIC, tenant 4 fixture, bank account 2).
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-11'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const A = async (m, p, b) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined }); const t = await x.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 200); } return { status: x.status, body: j }; };
const bid = 2; const own = sql("select id from journal_entries where tenant_id=4 and status='posted' order by id limit 1");
const gl = sql(`select coalesce(sum(g.debit_amount-g.credit_amount),0) from general_ledger g where g.tenant_id=4 and g.account_id=(select gl_account_id from bank_accounts_ledger where id=${bid})`);
let r = await A('POST', '/accounting/bank-reconciliation/sessions', { bank_account_id: bid, statement_date: '2026-11-27', statement_balance: 777777.77, notes: 'RUN-004 invented amounts' }); const sid = r.body?.data?.id ?? r.body?.id;
r = await A('POST', `/accounting/bank-reconciliation/sessions/${sid}/lines`, { journal_entry_id: Number(own), line_description: 'unrelated real journal', statement_amount: 777777.77, book_amount: 777777.77, match_status: 'matched' });
const lineHttp = r.status; r = await A('POST', `/accounting/bank-reconciliation/sessions/${sid}/complete`, {});
const out = { bank: bid, ownJournal: own, glBalanceOfBankAccount: gl, lineHttp, completeHttp: r.status, session: sql(`select status||'|stmt '||statement_balance||'|book '||book_balance||'|diff '||difference||'|completed_at '||coalesce(completed_at::text,'-') from bank_reconciliation_sessions where id=${sid}`), bankAfter: sql(`select coalesce(last_reconciled_date::text,'-')||'|'||coalesce(last_reconciled_balance::text,'-')||'|ledger current_balance '||current_balance from bank_accounts_ledger where id=${bid}`), journalBankLines: sql(`select count(*) from journal_entry_lines l join bank_accounts_ledger b on b.gl_account_id=l.account_id where l.journal_entry_id=${own} and b.id=${bid}`) };
console.log(JSON.stringify(out, null, 1)); fs.writeFileSync(`${d.db}/AUD-R4-BANK-RECON2.json`, JSON.stringify({ run: 'RUN-004', ...out }, null, 1));
