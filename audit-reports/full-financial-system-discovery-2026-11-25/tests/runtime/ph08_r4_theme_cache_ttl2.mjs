// RUN-004 / PH-08: CLEAN public-theme cache test. Editor Save (#00ff00 on theme 13) -> activate -> poll public API every 45 s for 6 min. Prior primed cache entry (call before activation) must expire (TTL 5 min) because invalidatePublicThemeCacheByRestaurant has no caller. Restores theme 4 and deletes the theme-13 copy.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH, sql, API } from './chrome_lib.mjs';
const d = evidenceDirs('PH-08'); const O = 'http://localhost:13112';
const pub = async () => { const t = await (await fetch(API + '/public/restaurants/audit-c-diag/theme', { headers: { Origin: 'http://localhost:13113' } })).text(); return (t.match(/"primary":"([^"]+)"/) || [])[1]; };
const { ctx, page } = await openChrome();
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/theme-builder/editor/13`); await page.waitForTimeout(8000);
const primed = await pub();
await page.waitForTimeout(6000); await page.locator('input[type=text]').first().fill('#00ff00'); await page.waitForTimeout(1500); await page.locator('button:has-text("Save"), button:has-text("Publish")').first().click({ timeout: 20000 }); await page.waitForTimeout(3500);
const ck = (await ctx.cookies()).filter((c) => c.name === 'poos_session').map((c) => `${c.name}=${c.value}`).join('; ');
const act = async (id) => (await fetch(`${API}/admin/themes/${id}/activate`, { method: 'POST', headers: { Origin: O, Cookie: ck } })).status;
const samples = []; const t0 = Date.now();
samples.push({ s: 0, event: 'primed before edit', primary: primed, db13: sql("select config->'colors'->>'primary' from themes_v2 where id=13") });
samples.push({ s: 0, event: 'activate 13', http: await act(13), primary: await pub() });
for (let i = 1; i <= 8; i++) { await new Promise((r) => setTimeout(r, 45000)); const s = { s: Math.round((Date.now() - t0) / 1000), primary: await pub(), db13: sql("select config->'colors'->>'primary' from themes_v2 where id=13") }; samples.push(s); console.log(JSON.stringify(s)); if (s.primary === '#00ff00') break; }
samples.push({ event: 'restore theme 4', http: await act(4), active: sql("select string_agg(id||':'||is_active,' ') from themes_v2 where tenant_id=4") });
fs.writeFileSync(`${d.net}/AUD-R4-THEME-CACHE-TTL2.json`, JSON.stringify({ run: 'RUN-004', samples }, null, 1)); console.log(JSON.stringify(samples[1]), JSON.stringify(samples.slice(-1)));
await ctx.close();
