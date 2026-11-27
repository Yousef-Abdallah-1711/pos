// RUN-004 / PH-08: persistence of an INVALID colour through the editor Save, then valid colour + activation + storefront computed style in Chrome; restores state.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-08'); const id = process.env.THEME_ID || '13';
const { ctx, page } = await openChrome();
const calls = []; page.on('response', async (r) => { if (r.url().includes(':18082') && /theme/i.test(r.url()) && r.request().method() !== 'GET') { let b = ''; try { b = (await r.text()).slice(0, 200); } catch {} calls.push({ m: r.request().method(), u: r.url().replace('http://localhost:18082/api/v1', ''), s: r.status(), req: (r.request().postData() || '').slice(0, 300), b }); } });
const log = []; const rec = (k, desc, extra = {}) => { log.push({ k, desc, ...extra }); console.log(k, desc, JSON.stringify(extra).slice(0, 420)); };
const cfg = () => sql(`select coalesce(config::text,'')||' || '||coalesce(theme_json::text,'') from themes_v2 where id=${id}`);
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/theme-builder/editor/${id}`); await page.waitForTimeout(8000);
const t1 = page.locator('input[type=text]').first(); const save = page.locator('button:has-text("Save"), button:has-text("Publish")').first();
for (const [k, v] of [['P1', 'not-a-color'], ['P2', 'url(javascript:alert(1))'], ['P3', '#ff0000']]) {
  await t1.fill(v); const b = calls.length; await save.click({ timeout: 5000 }).catch(() => {}); await page.waitForTimeout(4000);
  const c = cfg(); rec(k, `save primary colour "${v}"`, { calls: calls.slice(b).map((x) => `${x.m} ${x.u} ${x.s} ${x.b.slice(0, 80)}`), persisted: c.includes(v) || c.toLowerCase().includes(v.toLowerCase()), snippet: c.slice(c.toLowerCase().indexOf('primary'), c.toLowerCase().indexOf('primary') + 80) });
}
await page.goto(`${DASH}/en/dashboard/theme-builder`, { waitUntil: 'domcontentloaded' }); await page.waitForTimeout(6000);
await page.locator('button:has-text("Activate Theme")').last().click({ timeout: 6000 }).catch(() => {}); await page.waitForTimeout(4000);
rec('A1', 'activated', { active: sql("select string_agg(id||':'||is_active,' ') from themes_v2 where tenant_id=4") });
const sp = await ctx.newPage(); let sfInfo = {};
try { await sp.goto('http://audit-c-diag.localhost:13113/en', { waitUntil: 'domcontentloaded', timeout: 120000 }); await sp.waitForTimeout(8000); sfInfo = await sp.evaluate(() => { const cs = getComputedStyle(document.documentElement); const vars = ['--primary', '--color-primary', '--theme-primary', '--primary-color'].map((v) => v + '=' + cs.getPropertyValue(v).trim()); const btn = document.querySelector('button, a[class*=btn], a[class*=button]'); return { vars, btnBg: btn ? getComputedStyle(btn).backgroundColor : null, title: document.title, len: document.body.innerText.length }; }); await sp.screenshot({ path: `${d.shot}/AUD-R4-THEME-storefront-active-copy.png` }).catch(() => {}); } catch (e) { sfInfo = { error: String(e.message).slice(0, 100) }; }
rec('A2', 'storefront in Chrome while the edited copy is active', sfInfo);
await sp.close();
await page.goto(`${DASH}/en/dashboard/theme-builder`, { waitUntil: 'domcontentloaded' }); await page.waitForTimeout(6000);
await page.locator('div:has(h3:has-text("Default Theme")):not(:has(h3:has-text("(Copy)"))) button:has-text("Activate Theme")').first().click({ timeout: 6000 }).catch(() => {}); await page.waitForTimeout(4000);
rec('R1', 'restore Default Theme active', { active: sql("select string_agg(id||':'||is_active,' ') from themes_v2 where tenant_id=4") });
fs.writeFileSync(`${d.net}/AUD-R4-THEME-PERSIST.json`, JSON.stringify({ run: 'RUN-004', tenant: 4, themeId: id, log, calls }, null, 1));
await ctx.close();
