// RUN-004 / PH-05: does the user-level Appearance > "Theme Colors" action overwrite the ACTIVE public storefront theme? HEADED Chrome, tenant 4 fixture. Sets theme 13 colour to a marker first via API-free path (DB before/after only).
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-05');
const { ctx, page } = await openChrome();
const calls = []; page.on('response', async (r) => { if (r.url().includes(':18082') && r.request().method() !== 'GET') { let b = ''; try { b = (await r.text()).slice(0, 160); } catch {} calls.push({ m: r.request().method(), u: r.url().replace('http://localhost:18082/api/v1', ''), s: r.status(), req: (r.request().postData() || '').slice(0, 240), b }); } });
const row = () => sql("select id||':'||name||' primary='||(config->'colors'->>'primary')||' active='||is_active||' upd='||updated_at from themes_v2 where tenant_id=4 and id in (4,13) order by id").replace(/\n/g, ' ; ');
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/settings/appearance`); await page.waitForTimeout(7000);
const before = row(); console.log('before', before);
await page.screenshot({ path: `${d.shot}/AUD-R4-APPEARANCE-before.png`, fullPage: true }).catch(() => {});
const controls = await page.evaluate(() => [...document.querySelectorAll('main button')].map((b) => (b.innerText || '').trim().slice(0, 40)).filter(Boolean));
console.log('controls', JSON.stringify(controls));
const b0 = calls.length; await page.locator('main button:has-text("Theme Colors")').first().click({ timeout: 6000 }).catch((e) => console.log('click failed', e.message.slice(0, 80))); await page.waitForTimeout(5000);
const after = row(); console.log('after', after); console.log('calls', JSON.stringify(calls.slice(b0).map((c) => `${c.m} ${c.u} ${c.s} ${c.req.slice(0, 120)}`)));
await page.screenshot({ path: `${d.shot}/AUD-R4-APPEARANCE-after.png`, fullPage: true }).catch(() => {});
fs.writeFileSync(`${d.net}/AUD-R4-APPEARANCE-THEME.json`, JSON.stringify({ run: 'RUN-004', before, after, controls, calls: calls.slice(b0) }, null, 1));
await ctx.close();
