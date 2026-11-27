// RUN-004 / PH-08: does ACTIVATING a theme overwrite colour edits saved in the editor? UI Save (#ff0000) -> DB -> API activate -> DB -> public API -> restore.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH, sql, API } from './chrome_lib.mjs';
const d = evidenceDirs('PH-08');
const { ctx, page } = await openChrome();
const log = []; const rec = (k, desc, extra = {}) => { log.push({ k, desc, ...extra }); console.log(k, desc, JSON.stringify(extra).slice(0, 300)); };
const prim = () => sql("select (config->'colors'->>'primary')||' | col '||primary_color||' | updated '||updated_at from themes_v2 where id=13");
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/theme-builder/editor/13`); await page.waitForTimeout(8000);
rec('S0', 'before', { row: prim() });
await page.locator('input[type=text]').first().fill('#ff0000'); await page.locator('button:has-text("Save"), button:has-text("Publish")').first().click({ timeout: 5000 }); await page.waitForTimeout(4000);
rec('S1', 'after editor Save', { row: prim() });
const cookies = (await ctx.cookies()).filter((c) => c.name === 'poos_session'); const ck = cookies.map((c) => `${c.name}=${c.value}`).join('; ');
const x = await fetch(`${API}/admin/themes/13/activate`, { method: 'POST', headers: { Origin: 'http://localhost:13112', Cookie: ck } }); rec('S2', 'API activate theme 13', { http: x.status, row: prim() });
const t = await (await fetch(API + '/public/restaurants/audit-c-diag/theme', { headers: { Origin: 'http://localhost:13113' } })).text(); rec('S3', 'public API after activation', { primary: (t.match(/"primary":"([^"]+)"/) || [])[1] });
await fetch(`${API}/admin/themes/4/activate`, { method: 'POST', headers: { Origin: 'http://localhost:13112', Cookie: ck } }); rec('S4', 'restored', { active: sql("select string_agg(id||':'||is_active,' ') from themes_v2 where tenant_id=4"), row13: prim() });
fs.writeFileSync(`${d.net}/AUD-R4-THEME-SAVE-VS-ACTIVATE.json`, JSON.stringify({ run: 'RUN-004', log }, null, 1));
await ctx.close();
