// RUN-004 / PH-19: revalidation of F-064 duplicate table number against the live isolated DB. Tenant 4. API-DIAGNOSTIC. Cleans up any table it creates (diagnostic fixture).
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-19'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const A = async (m, p, b) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined }); const t = await x.text(); return { status: x.status, text: t.slice(0, 300) }; };
const log = [];
for (const n of [4, 1, 12]) { const r = await A('POST', '/tables', { table_number: n, name: 'dup-r4-' + n, capacity: 4 }); log.push({ n, http: r.status, resp: r.text }); console.log('dup table_number', n, r.status, r.text.slice(0, 200)); }
console.log('tables now:', sql('select id,table_number from restaurant_tables where tenant_id=4 order by id').replace(/\n/g, ' ; '));
console.log('constraint:', sql("select pg_get_constraintdef(oid) from pg_constraint where conname='uq_restaurant_tables_number'"));
console.log('first migration that created it:', sql("select version from schema_migrations order by version desc limit 1"));
fs.writeFileSync(`${d.db}/AUD-R4-F064-dup-table.json`, JSON.stringify({ run: 'RUN-004', log, tables: sql('select id,table_number,created_at from restaurant_tables where tenant_id=4 order by id') }, null, 1));
