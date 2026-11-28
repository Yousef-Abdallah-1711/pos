// RUN-003 / PH-16: tables API validation (duplicate number, negative/zero/huge capacity, status transitions, delete while occupied). Tenant 4. API-DIAGNOSTIC.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-16'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const A = async (m, p, b) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined }); const t = await x.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 200); } return { status: x.status, body: j }; };
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 300)); };
let r = await A('POST', '/tables', { table_number: 1, name: 'dup', capacity: 4 }); rec('TBL-01', 'duplicate table number 1', { http: r.status, resp: JSON.stringify(r.body).slice(0, 160) });
r = await A('POST', '/tables', { table_number: 50, name: 'neg', capacity: -2 }); const t50 = r.body?.data?.id ?? r.body?.id; rec('TBL-02', 'negative capacity', { http: r.status, resp: JSON.stringify(r.body).slice(0, 160) });
r = await A('POST', '/tables', { table_number: 51, name: 'zero', capacity: 0 }); rec('TBL-03', 'zero capacity', { http: r.status });
r = await A('POST', '/tables', { table_number: 52, name: 'huge', capacity: 2147483647 }); rec('TBL-04', 'capacity int max', { http: r.status, resp: JSON.stringify(r.body).slice(0, 160) });
r = await A('POST', '/tables', { table_number: -5, name: 'negnum', capacity: 2 }); rec('TBL-05', 'negative table number', { http: r.status, resp: JSON.stringify(r.body).slice(0, 160) });
r = await A('POST', '/tables', { table_number: 53, name: '<b>x</b>'.repeat(200), capacity: 2 }); rec('TBL-06', 'name 1,800 chars with HTML', { http: r.status, resp: JSON.stringify(r.body).slice(0, 160) });
const occ = sql("select id||'|'||status from restaurant_tables where tenant_id=4 and id=1"); rec('TBL-07', 'table 1 state after dine-in tests', { occ });
r = await A('DELETE', '/tables/1'); rec('TBL-08', 'delete table that has orders / is occupied', { http: r.status, resp: JSON.stringify(r.body).slice(0, 160), exists: sql('select count(*) from restaurant_tables where id=1') });
fs.writeFileSync(`${d.db}/AUD-TABLES.json`, JSON.stringify({ run: 'RUN-003', kind: 'API-DIAGNOSTIC', log }, null, 1));
