// RUN-004 / PH-08: attributes the storefront colour mismatch: public theme API vs storefront cache (revalidate 3600 in app/[locale]/page.tsx:75,131). API-DIAGNOSTIC on tenant 4; restores Default Theme.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-08'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const A = async (m, p) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json' } }); return { s: x.status, t: await x.text() }; };
const pub = async () => { const x = await fetch(API + '/public/restaurants/audit-c-diag/theme', { headers: { Origin: 'http://localhost:13113' } }); const t = await x.text(); return { s: x.status, primary: (t.match(/"primary(?:_color)?":\s*"([^"]+)"/) || [])[1], ff: /ff0000/i.test(t), bytes: t.length }; };
const out = {};
out.before = await pub();
out.activate13 = (await A('POST', '/admin/themes/13/activate')).s; out.after13 = await pub();
out.activate4 = (await A('POST', '/admin/themes/4/activate')).s; out.after4 = await pub();
out.active = sql("select string_agg(id||':'||is_active,' ') from themes_v2 where tenant_id=4");
console.log(JSON.stringify(out, null, 1)); fs.writeFileSync(`${d.net}/AUD-R4-THEME-API-VS-STORE.json`, JSON.stringify({ run: 'RUN-004', out }, null, 1));
