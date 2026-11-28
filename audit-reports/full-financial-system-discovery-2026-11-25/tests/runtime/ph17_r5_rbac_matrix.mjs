// RUN-005 / PH-17 (Phase K): authorization matrix on the ENFORCE-mode isolated API for roles ADMIN (owner), MANAGER, EMPLOYEE (cashier), DRIVER. Records status, number of rows returned and presence of sensitive fields (actual data, not only status).
// Write probes send EMPTY bodies: 403 = denied, 400 = authorised but invalid, nothing is created. Uses synthetic users created through the API.
import fs from 'fs';
import { sql, API as API0, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const API = process.env.API_URL || API0; const TAG = process.env.OUT_TAG || 'ENFORCE';
const d = evidenceDirs('PH-17'); const O = 'http://localhost:13112';
const raw = async (m, p, ck, b) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, ...(ck ? { Cookie: ck } : {}), 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined, signal: AbortSignal.timeout(60000) }); const t = await x.text(); return { status: x.status, text: t, ck: x.headers.getSetCookie?.().map((c) => c.split(';')[0]).join('; ') }; };
const login = async (email) => (await raw('POST', '/auth/login', null, { email, password: PASSWORD })).ck;
const owner = await login('audit.c.owner@example.test');
const roleId = (code) => Number(sql(`select id from roles where tenant_id=4 and role_code='${code}'`));
const mk = async (email, name, code) => {
  let uid = sql(`select id from users where email='${email}'`);
  if (!uid) { const r = await raw('POST', '/users', owner, { name, email, phone: '+201000000' + String(Math.floor(Math.random() * 900) + 100), password: PASSWORD }); uid = sql(`select id from users where email='${email}'`); }
  const a = await raw('POST', `/rbac/users/${uid}/roles`, owner, { role_id: roleId(code) });
  const rm = code !== 'EMPLOYEE' ? await raw('DELETE', `/rbac/users/${uid}/roles/${roleId('EMPLOYEE')}`, owner) : { status: 'n/a' };
  return { uid, assign: a.status, removeEmployee: rm.status, roles: sql(`select string_agg(r.role_code,',') from user_roles ur join roles r on r.id=ur.role_id where ur.user_id=${uid}`) };
};
const setup = { MANAGER: await mk('manager.r5@example.test', 'Manager R5', 'MANAGER'), DRIVER: await mk('driver.r5@example.test', 'Driver R5', 'DRIVER') };
console.log('setup', JSON.stringify(setup));
const users = { ADMIN: 'audit.c.owner@example.test', MANAGER: 'manager.r5@example.test', EMPLOYEE: 'cashier.c@example.test', DRIVER: 'driver.r5@example.test' };
const cookies = {}; for (const [r, e] of Object.entries(users)) cookies[r] = await login(e);
const GETS = ['/accounting/journal-entries?page=1&limit=5', '/accounting/accounts', '/accounting/trial-balance?as_of=2026-12-31', '/accounting/customers', '/suppliers', '/accounting/customer-invoices', '/accounting/supplier-invoices', '/accounting/bank-accounts', '/hr/employees', '/hr/salaries', '/hr/attendance', '/hr/leaves', '/hr/contracts', '/users', '/rbac/roles', '/reports/financial/profit-loss?date_from=2026-01-01&date_to=2026-12-31', '/reports/financial/balance-sheet?as_of=2026-12-31', '/reports/sales?date_from=2026-11-01&date_to=2026-11-30&granularity=day', '/admin/orders?page=1&page_size=5', '/products?page=1&page_size=5', '/categories', '/kds/queue', '/tables', '/admin/drivers', '/admin/themes?pageSize=5&page=1', '/payments/subscription', '/notifications', '/accounting/csv/journal-entries/export'];
const WRITES = [['POST', '/accounting/accounts'], ['POST', '/accounting/journal-entries'], ['POST', '/hr/employees'], ['POST', '/suppliers'], ['POST', '/accounting/customers'], ['POST', '/admin/drivers'], ['POST', '/rbac/roles'], ['POST', '/accounting/bank-reconciliation/sessions']];
const SENS = /base_salary|salary_amount|password_hash|iban|national_id|bank_account/i;
const cells = [];
for (const role of Object.keys(users)) {
  for (const p of GETS) { const r = await raw('GET', p, cookies[role]); let n = null; try { const j = JSON.parse(r.text); const arr = j?.data?.items ?? j?.data ?? j?.items ?? j?.products ?? j?.orders ?? j?.entries ?? j; n = Array.isArray(arr) ? arr.length : (arr && typeof arr === 'object' ? Object.keys(arr).length : null); } catch { /* not json */ } cells.push({ role, method: 'GET', path: p.split('?')[0], status: r.status, rows_or_keys: n, bytes: r.text.length, sensitive_fields: SENS.test(r.text) }); }
  for (const [m, p] of WRITES) { const r = await raw(m, p, cookies[role], {}); cells.push({ role, method: m, path: p, status: r.status, outcome: r.status === 403 ? 'DENIED' : r.status === 400 || r.status === 422 ? 'AUTHORISED (rejected as invalid body)' : 'OTHER', body: r.text.slice(0, 80) }); }
}
const table = {}; for (const c of cells) { const k = `${c.method} ${c.path}`; (table[k] ||= {})[c.role] = c.status + (c.sensitive_fields ? '*' : ''); }
console.log('role columns: ADMIN MANAGER EMPLOYEE DRIVER  (* = response contains salary/bank/national-id fields)'); for (const [k, v] of Object.entries(table)) console.log(k.padEnd(52), ['ADMIN', 'MANAGER', 'EMPLOYEE', 'DRIVER'].map((r) => String(v[r] ?? '-').padEnd(6)).join(' '));
fs.writeFileSync(`${d.net}/AUD-R5-RBAC-MATRIX-${TAG}.json`, JSON.stringify({ run: 'RUN-005', api: API + ' (' + TAG + ')', setup, cells }, null, 1));
