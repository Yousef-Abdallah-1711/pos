// RUN-004 / PH-08: public theme cache staleness after activation (TTL 5 min, invalidatePublicThemeCacheByRestaurant has no caller). Activates theme 13 (#ff0000), samples the public API every 60 s for 6 min, restores theme 4.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-08'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const act = async (id) => (await fetch(`${API}/admin/themes/${id}/activate`, { method: 'POST', headers: { Origin: O, Cookie: ck } })).status;
const pub = async () => { const t = await (await fetch(API + '/public/restaurants/audit-c-diag/theme', { headers: { Origin: 'http://localhost:13113' } })).text(); return (t.match(/"primary(?:_color)?":\s*"([^"]+)"/) || [])[1]; };
const samples = []; const t0 = Date.now();
samples.push({ s: 0, event: 'before (theme 4 active)', primary: await pub() });
samples.push({ s: 0, event: 'activate theme 13', status: await act(13), primary: await pub() });
for (let i = 1; i <= 6; i++) { await new Promise((r) => setTimeout(r, 60000)); samples.push({ s: Math.round((Date.now() - t0) / 1000), primary: await pub(), active: sql("select string_agg(id||':'||is_active,' ') from themes_v2 where tenant_id=4") }); console.log(JSON.stringify(samples[samples.length - 1])); }
samples.push({ event: 'restore theme 4', status: await act(4), primary: await pub() });
fs.writeFileSync(`${d.net}/AUD-R4-THEME-CACHE-TTL.json`, JSON.stringify({ run: 'RUN-004', samples }, null, 1)); console.log(JSON.stringify(samples.slice(-1)));
